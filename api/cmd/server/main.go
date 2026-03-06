package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/core/config"
	"github.com/portnd/the-sentinel-core/internal/core/database"
	"github.com/portnd/the-sentinel-core/internal/core/middleware"
	authHttp "github.com/portnd/the-sentinel-core/internal/modules/auth/delivery/http"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	authRepo "github.com/portnd/the-sentinel-core/internal/modules/auth/repository"
	authUsecase "github.com/portnd/the-sentinel-core/internal/modules/auth/usecase"
	financeDomain "github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
	financeHttp "github.com/portnd/the-sentinel-core/internal/modules/finance/delivery/http"
	financeRepo "github.com/portnd/the-sentinel-core/internal/modules/finance/repository"
	financeUsecase "github.com/portnd/the-sentinel-core/internal/modules/finance/usecase"
	healthHttp "github.com/portnd/the-sentinel-core/internal/modules/health/delivery/http"
	perfHttp "github.com/portnd/the-sentinel-core/internal/modules/performance/delivery/http"
	perfRepo "github.com/portnd/the-sentinel-core/internal/modules/performance/repository"
	perfUsecase "github.com/portnd/the-sentinel-core/internal/modules/performance/usecase"
	sentinelHttp "github.com/portnd/the-sentinel-core/internal/modules/sentinel/delivery/http"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	sentinelRepo "github.com/portnd/the-sentinel-core/internal/modules/sentinel/repository"
	sentinelUsecase "github.com/portnd/the-sentinel-core/internal/modules/sentinel/usecase"
	walletHttp "github.com/portnd/the-sentinel-core/internal/modules/wallet/delivery/http"
	walletDomain "github.com/portnd/the-sentinel-core/internal/modules/wallet/domain"
	walletRepo "github.com/portnd/the-sentinel-core/internal/modules/wallet/repository"
	walletUsecase "github.com/portnd/the-sentinel-core/internal/modules/wallet/usecase"
)

func main() {
	log.Println("🛡️  KOMGRIP API - Starting...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	log.Printf("📋 Environment: %s", cfg.AppEnv)

	log.Println("🔌 Connecting to PostgreSQL...")
	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := database.InitMongo(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	log.Println("✅ MongoDB connected")

	log.Println("🔌 Connecting to Redis...")
	redisClient, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Redis connected")

	// Auto-migrate database schemas
	log.Println("🔄 Running database migrations...")
	if err := db.AutoMigrate(
		&authDomain.User{},
		&walletDomain.Wallet{},
		&walletDomain.Transaction{},
		&sentinelDomain.SystemConfig{},
		&sentinelDomain.Project{},
		&sentinelDomain.Sprint{},
		&sentinelDomain.Milestone{},
		&sentinelDomain.Epic{},
		&sentinelDomain.Task{},
		&sentinelDomain.Submission{},
		&sentinelDomain.Appeal{},
		&sentinelDomain.TaskDependency{},
		&sentinelDomain.TaskComment{},
		&sentinelDomain.TimeLog{},
		&financeDomain.MonthlyEntry{},
	); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Initialize Auth Module (Hexagonal Architecture Wiring)
	log.Println("🔐 Initializing Auth Module...")
	authRepository := authRepo.NewPostgresRepository(db)
	authUsecaseInstance := authUsecase.NewAuthUsecase(authRepository)
	log.Println("✅ Auth Module initialized")

	// Initialize Wallet Module (Hexagonal Architecture Wiring)
	log.Println("💰 Initializing Wallet Module...")
	walletRepository := walletRepo.NewPostgresRepository(db)
	walletUsecaseInstance := walletUsecase.NewWalletUsecase(walletRepository, db)
	walletHandlerInstance := walletHttp.NewWalletHandler(walletUsecaseInstance)
	log.Println("✅ Wallet Module initialized")

	// Initialize Sentinel Module (Hexagonal Architecture Wiring)
	log.Println("🛡️  Initializing Sentinel Module...")

	sentinelRepository := sentinelRepo.NewPostgresRepository(db)
	aiUsageTracker := sentinelRepo.NewMemoryUsageTracker()
	var aiService sentinelDomain.AIService
	switch {
	case cfg.GroqAPIKey != "":
		var errAI error
		aiService, errAI = sentinelRepo.NewGroqService(cfg.GroqAPIKey, sentinelRepository, aiUsageTracker)
		if errAI != nil {
			log.Printf("⚠️  Groq AI init failed, using no-op: %v", errAI)
			aiService = sentinelRepo.NewNoopAIService()
		} else {
			log.Println("✅ Groq AI service enabled (estimate & code review)")
		}
	case cfg.GeminiAPIKey != "":
		var errAI error
		aiService, errAI = sentinelRepo.NewGeminiService(cfg.GeminiAPIKey, sentinelRepository, aiUsageTracker)
		if errAI != nil {
			log.Printf("⚠️  Gemini AI init failed, using no-op: %v", errAI)
			aiService = sentinelRepo.NewNoopAIService()
		} else {
			log.Println("✅ Gemini AI service enabled (estimate & code review)")
		}
	default:
		aiService = sentinelRepo.NewNoopAIService()
		log.Println("⚠️  GROQ_API_KEY / GEMINI_API_KEY not set; AI estimate/code review disabled")
	}
	sentinelUsecaseInstance := sentinelUsecase.NewSentinelUsecase(sentinelRepository, aiService, authRepository, aiUsageTracker, cfg.AILimitRPM, cfg.AILimitRPD)

	log.Println("✅ Sentinel Module initialized")

	// Initialize Performance Module
	log.Println("📊 Initializing Performance Module...")
	perfRepository := perfRepo.NewPostgresRepository(db)
	perfUsecaseInstance := perfUsecase.NewPerformanceUsecase(perfRepository, authRepository)
	log.Println("✅ Performance Module initialized")

	// Initialize Finance Module (accounting entries + CEO summary)
	log.Println("📒 Initializing Finance Module...")
	financeRepository := financeRepo.NewPostgresRepository(db)
	financeUsecaseInstance := financeUsecase.NewFinanceUsecase(financeRepository)
	log.Println("✅ Finance Module initialized")

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register module routes
	healthHttp.RegisterRoutes(router, db, mongoClient, redisClient)

	// API v1 Group - Uniform Interface for all endpoints
	apiGroup := router.Group("/api/v1")
	// Auth middleware (used by protected routes)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	
	// Auth routes (includes user management endpoints)
	authHttp.RegisterRoutes(apiGroup, authUsecaseInstance, authMiddleware)

	// Wallet routes (protected by auth middleware)
	walletHttp.RegisterRoutes(apiGroup, walletHandlerInstance, authMiddleware)

	// Sentinel routes (protected by auth middleware)
	sentinelGroup := apiGroup.Group("")
	sentinelGroup.Use(authMiddleware)
	sentinelHttp.RegisterRoutes(sentinelGroup, sentinelUsecaseInstance, cfg.GoogleAPIKey)

	// Performance routes (protected by auth middleware)
	perfGroup := apiGroup.Group("")
	perfGroup.Use(authMiddleware)
	perfHttp.RegisterRoutes(perfGroup, perfUsecaseInstance)

	// Finance routes (CEO: accounting entries + summary)
	finGroup := apiGroup.Group("")
	finGroup.Use(authMiddleware)
	financeHttp.RegisterRoutes(finGroup, financeUsecaseInstance)

	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	log.Printf("🔗 Listening on http://0.0.0.0:%s (all interfaces)", cfg.AppPort)
	log.Printf("🌐 Health endpoint: http://localhost:%s/health", cfg.AppPort)
	log.Printf("🔐 Auth endpoints: http://localhost:%s/api/v1/auth/register | /api/v1/auth/login", cfg.AppPort)
	log.Printf("👥 User Management (CEO): GET /api/v1/auth/users | POST /api/v1/auth/users | POST /api/v1/auth/users/import | PATCH /api/v1/auth/users/:id/role")
	log.Printf("💰 Wallet endpoints: http://localhost:%s/api/v1/wallets/me | /api/v1/wallets/transfer", cfg.AppPort)
	log.Printf("🛡️  Sentinel endpoints: http://localhost:%s/api/v1/sentinel/tasks | /api/v1/sentinel/tasks/my", cfg.AppPort)
	log.Printf("⚙️  AI Config endpoints (CEO): GET/PUT /api/v1/admin/config | GET /api/v1/admin/models")

	if err := router.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
