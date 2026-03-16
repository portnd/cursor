package http

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	handler := NewHealthHandler(db, redisClient)

	router.GET("/health", handler.Check)
	router.HEAD("/health", handler.Check)
	router.GET("/live", handler.Live)
	router.HEAD("/live", handler.Live)
}
