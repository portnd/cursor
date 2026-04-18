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
	attendanceHttp "github.com/portnd/the-sentinel-core/internal/modules/attendance/delivery/http"
	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	attendanceRepo "github.com/portnd/the-sentinel-core/internal/modules/attendance/repository"
	attendanceUsecase "github.com/portnd/the-sentinel-core/internal/modules/attendance/usecase"
	deploymentHttp "github.com/portnd/the-sentinel-core/internal/modules/deployment/delivery/http"
	deploymentDomain "github.com/portnd/the-sentinel-core/internal/modules/deployment/domain"
	deploymentRepo "github.com/portnd/the-sentinel-core/internal/modules/deployment/repository"
	deploymentUsecase "github.com/portnd/the-sentinel-core/internal/modules/deployment/usecase"
	authHttp "github.com/portnd/the-sentinel-core/internal/modules/auth/delivery/http"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	authRepo "github.com/portnd/the-sentinel-core/internal/modules/auth/repository"
	authUsecase "github.com/portnd/the-sentinel-core/internal/modules/auth/usecase"
	financeHttp "github.com/portnd/the-sentinel-core/internal/modules/finance/delivery/http"
	financeDomain "github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
	financeRepo "github.com/portnd/the-sentinel-core/internal/modules/finance/repository"
	financeUsecase "github.com/portnd/the-sentinel-core/internal/modules/finance/usecase"
	healthHttp "github.com/portnd/the-sentinel-core/internal/modules/health/delivery/http"
	perfHttp "github.com/portnd/the-sentinel-core/internal/modules/performance/delivery/http"
	perfRepo "github.com/portnd/the-sentinel-core/internal/modules/performance/repository"
	perfUsecase "github.com/portnd/the-sentinel-core/internal/modules/performance/usecase"
	pricingHttp "github.com/portnd/the-sentinel-core/internal/modules/pricing/delivery/http"
	pricingDomain "github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
	pricingRepo "github.com/portnd/the-sentinel-core/internal/modules/pricing/repository"
	pricingUsecase "github.com/portnd/the-sentinel-core/internal/modules/pricing/usecase"
	pulseHttp "github.com/portnd/the-sentinel-core/internal/modules/pulse/delivery/http"
	pulseDomain "github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
	pulseRepo "github.com/portnd/the-sentinel-core/internal/modules/pulse/repository"
	pulseUsecase "github.com/portnd/the-sentinel-core/internal/modules/pulse/usecase"
	sentinelHttp "github.com/portnd/the-sentinel-core/internal/modules/sentinel/delivery/http"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	sentinelRepo "github.com/portnd/the-sentinel-core/internal/modules/sentinel/repository"
	sentinelUsecase "github.com/portnd/the-sentinel-core/internal/modules/sentinel/usecase"
)

// Version and BuildTime are injected at link time via -ldflags (-X main.Version=...).
var (
	Version   = "dev"
	BuildTime = "unknown"
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

	log.Println("🔌 Connecting to Redis...")
	redisClient, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Redis connected")

	// Auto-migrate database schemas
	log.Println("🔄 Running database migrations...")
	if err := db.AutoMigrate(
		&authDomain.Team{},
		&authDomain.User{},
		&authDomain.TeamTransaction{},
		&sentinelDomain.SystemConfig{},
		&sentinelDomain.Project{},
		&sentinelDomain.ProjectPmAssignment{},
		&sentinelDomain.Sprint{},
		&sentinelDomain.Milestone{},
		&sentinelDomain.Epic{},
		&sentinelDomain.Task{},
		&sentinelDomain.Submission{},
		&sentinelDomain.Appeal{},
		&sentinelDomain.TaskDependency{},
		&sentinelDomain.TaskComment{},
		&sentinelDomain.TaskActivityEvent{},
		&sentinelDomain.TimeLog{},
		&sentinelDomain.ProjectTransaction{},
		&financeDomain.MonthlyEntry{},
		&pricingDomain.EmployeeSalary{},
		&pricingDomain.CompanyCostConfig{},
		&pricingDomain.ProjectCostSnapshot{},
		&pricingDomain.ProjectExpense{},
		&pulseDomain.DailyStandup{},
		&attendanceDomain.OfficeConfig{},
		&attendanceDomain.AttendanceRecord{},
		&attendanceDomain.LeaveRequest{},
		&attendanceDomain.LeavePolicy{},
		&attendanceDomain.HolidayCalendar{},
		&attendanceDomain.LeaveAuditLog{},
		&attendanceDomain.LeaveNotification{},
		&deploymentDomain.DeploymentRequest{},
	); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Initialize Auth Module (Hexagonal Architecture Wiring)
	log.Println("🔐 Initializing Auth Module...")
	authRepository := authRepo.NewPostgresRepository(db)
	authUsecaseInstance := authUsecase.NewAuthUsecase(authRepository)
	log.Println("✅ Auth Module initialized")

	// Initialize Sentinel Module (Hexagonal Architecture Wiring)
	log.Println("🛡️  Initializing Sentinel Module...")

	sentinelRepository := sentinelRepo.NewPostgresRepository(db)
	aiUsageTracker := sentinelRepo.NewMemoryUsageTracker()
	var aiService sentinelDomain.AIService
	switch {
	case cfg.UseNoopAI:
		aiService = sentinelRepo.NewNoopAIService()
		log.Println("✅ AI service: noop (USE_NOOP_AI=true — no external API calls)")
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

	// Initialize Pricing Module (fully loaded costing & quotation)
	log.Println("💰 Initializing Pricing Module...")
	pricingRepository := pricingRepo.NewPostgresRepository(db)
	pricingUsecaseInstance := pricingUsecase.NewCostingUsecase(pricingRepository)
	log.Println("✅ Pricing Module initialized")

	// Initialize Team Finance Usecase (Internal VC model — depends on auth + pricing repos)
	log.Println("🏦 Initializing Team Finance Usecase...")
	teamFinanceUsecaseInstance := authUsecase.NewTeamFinanceUsecase(authRepository, pricingRepository)
	log.Println("✅ Team Finance Usecase initialized")

	// Initialize Project Finance Usecase (Internal VC model — per-project capital)
	log.Println("💼 Initializing Project Finance Usecase...")
	projectFinanceUsecaseInstance := sentinelUsecase.NewProjectFinanceUsecase(sentinelRepository, authRepository, pricingRepository)
	log.Println("✅ Project Finance Usecase initialized")

	// Initialize Pulse Module (async daily standup & activity tracker)
	log.Println("📡 Initializing Pulse Module...")
	pulseRepository := pulseRepo.NewPostgresRepository(db)
	pulseUsecaseInstance := pulseUsecase.NewPulseUsecase(pulseRepository, authRepository)
	log.Println("✅ Pulse Module initialized")

	log.Println("🕐 Initializing Attendance Module...")
	attendanceRepository := attendanceRepo.NewPostgresRepository(db)
	attendanceUsecaseInstance := attendanceUsecase.NewAttendanceUsecase(attendanceRepository, cfg.JWTSecret)
	log.Println("✅ Attendance Module initialized")

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// Large PPTX uploads (Canva exports): default 32 MiB multipart buffer is low; spill threshold for multipart parsing.
	router.MaxMultipartMemory = 128 << 20 // 128 MiB

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Content-Encoding", "Cache-Control"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.GzipResponse())

	// Register module routes
	healthHttp.RegisterRoutes(router, db, redisClient)

	// API v1 Group - Uniform Interface for all endpoints
	apiGroup := router.Group("/api/v1")
	// Auth middleware (used by protected routes)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Auth routes (includes user management endpoints)
	authHttp.RegisterRoutes(apiGroup, authUsecaseInstance, teamFinanceUsecaseInstance, authMiddleware)

	// Sentinel routes (protected by auth middleware)
	sentinelGroup := apiGroup.Group("")
	sentinelGroup.Use(authMiddleware)
	sentinelHttp.RegisterRoutes(sentinelGroup, sentinelUsecaseInstance, projectFinanceUsecaseInstance, cfg.GoogleAPIKey, cfg.CanvaAccessToken)

	// Performance routes (protected by auth middleware)
	perfGroup := apiGroup.Group("")
	perfGroup.Use(authMiddleware)
	perfHttp.RegisterRoutes(perfGroup, perfUsecaseInstance)

	// Finance routes (CEO: accounting entries + summary)
	finGroup := apiGroup.Group("")
	finGroup.Use(authMiddleware)
	financeHttp.RegisterRoutes(finGroup, financeUsecaseInstance)

	// Pricing routes (costing & quotation — protected by auth middleware)
	pricingGroup := apiGroup.Group("")
	pricingGroup.Use(authMiddleware)
	pricingHttp.RegisterRoutes(pricingGroup, pricingUsecaseInstance)

	// Pulse routes (daily standup & team activity — protected by auth middleware)
	pulseGroup := apiGroup.Group("")
	pulseGroup.Use(authMiddleware)
	pulseHttp.RegisterRoutes(pulseGroup, pulseUsecaseInstance)

	attendanceGroup := apiGroup.Group("")
	attendanceGroup.Use(authMiddleware)
	attendanceHttp.RegisterRoutes(attendanceGroup, attendanceUsecaseInstance)

	// Deployment routes (code review & deploy pipeline — protected by auth middleware)
	log.Println("🚀 Initializing Deployment Module...")
	deploymentRepository := deploymentRepo.NewPostgresRepository(db, authRepository)
	// sentinelUsecaseInstance implements deployment.domain.TaskStatusAdvancer (AdvanceTaskAfterDeploy)
	deploymentUsecaseInstance := deploymentUsecase.NewDeploymentUsecase(deploymentRepository, sentinelUsecaseInstance)
	deploymentGroup := apiGroup.Group("")
	deploymentGroup.Use(authMiddleware)
	deploymentHttp.RegisterRoutes(deploymentGroup, deploymentUsecaseInstance)
	log.Println("✅ Deployment Module initialized")

	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	log.Printf("🔗 Listening on http://0.0.0.0:%s (all interfaces)", cfg.AppPort)
	log.Printf("🌐 Health endpoint: http://localhost:%s/health", cfg.AppPort)
	log.Printf("🔐 Auth endpoint: http://localhost:%s/api/v1/auth/login (self-registration disabled)", cfg.AppPort)
	log.Printf("👥 User Management (CEO): GET /api/v1/auth/users | POST /api/v1/auth/users | POST /api/v1/auth/users/import | PATCH /api/v1/auth/users/:id/role")
	log.Printf("🛡️  Sentinel endpoints: http://localhost:%s/api/v1/sentinel/tasks | /api/v1/sentinel/tasks/my", cfg.AppPort)
	log.Printf("⚙️  AI Config endpoints (CEO): GET/PUT /api/v1/admin/config | GET /api/v1/admin/models")

	if err := router.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
