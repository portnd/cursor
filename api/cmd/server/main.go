package main

import (
	"fmt"
	"log"
	"time"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/core/config"
	"github.com/portnd/the-sentinel-core/internal/core/database"
	"github.com/portnd/the-sentinel-core/internal/core/middleware"
	authHttp "github.com/portnd/the-sentinel-core/internal/modules/auth/delivery/http"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	authRepo "github.com/portnd/the-sentinel-core/internal/modules/auth/repository"
	authUsecase "github.com/portnd/the-sentinel-core/internal/modules/auth/usecase"
	healthHttp "github.com/portnd/the-sentinel-core/internal/modules/health/delivery/http"
	sentinelHttp "github.com/portnd/the-sentinel-core/internal/modules/sentinel/delivery/http"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	sentinelRepo "github.com/portnd/the-sentinel-core/internal/modules/sentinel/repository"
	sentinelUsecase "github.com/portnd/the-sentinel-core/internal/modules/sentinel/usecase"
	walletHttp "github.com/portnd/the-sentinel-core/internal/modules/wallet/delivery/http"
	walletDomain "github.com/portnd/the-sentinel-core/internal/modules/wallet/domain"
	walletRepo "github.com/portnd/the-sentinel-core/internal/modules/wallet/repository"
	walletUsecase "github.com/portnd/the-sentinel-core/internal/modules/wallet/usecase"
)

// 👈 อย่าลืมเพิ่ม os

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
		&sentinelDomain.SystemConfig{}, // Dynamic AI Configuration
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

	// Init Repository first (needed by AI service for dynamic config)
	sentinelRepository := sentinelRepo.NewPostgresRepository(db)

	// Init Gemini AI Service with dynamic config support
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("❌ FATAL: GEMINI_API_KEY is required. Please set it in .env file.")
	}

	aiService, err := sentinelRepo.NewGeminiService(geminiKey, sentinelRepository)
	if err != nil {
		log.Fatalf("❌ FATAL: Failed to initialize Gemini AI: %v", err)
	}

	// Init Usecase (Inject AI + Auth Repo for role validation)
	sentinelUsecaseInstance := sentinelUsecase.NewSentinelUsecase(sentinelRepository, aiService, authRepository)

	log.Println("✅ Sentinel Module initialized with Dynamic AI Configuration + Role Validation")

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
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
	sentinelHttp.RegisterRoutes(sentinelGroup, sentinelUsecaseInstance)

	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	log.Printf("🌐 Health endpoint: http://localhost:%s/health", cfg.AppPort)
	log.Printf("🔐 Auth endpoints: http://localhost:%s/api/v1/auth/register | /api/v1/auth/login", cfg.AppPort)
	log.Printf("👥 User Management (CEO): GET /api/v1/auth/users | PATCH /api/v1/auth/users/:id/role")
	log.Printf("💰 Wallet endpoints: http://localhost:%s/api/v1/wallets/me | /api/v1/wallets/transfer", cfg.AppPort)
	log.Printf("🛡️  Sentinel endpoints: http://localhost:%s/api/v1/sentinel/tasks | /api/v1/sentinel/tasks/my", cfg.AppPort)
	log.Printf("⚙️  AI Config endpoints (CEO): GET/PUT /api/v1/admin/config | GET /api/v1/admin/models")

	if err := router.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
