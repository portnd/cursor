package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/komgrip/starter-kit/internal/core/config"
	"github.com/komgrip/starter-kit/internal/core/database"
	healthHttp "github.com/komgrip/starter-kit/internal/modules/health/delivery/http"
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

	healthHttp.RegisterRoutes(router, db, mongoClient, redisClient)

	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	log.Printf("🌐 Health endpoint: http://localhost:%s/health", cfg.AppPort)

	if err := router.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
