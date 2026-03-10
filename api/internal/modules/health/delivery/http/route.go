package http

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, mongoClient *mongo.Client, redisClient *redis.Client) {
	handler := NewHealthHandler(db, mongoClient, redisClient)

	router.GET("/health", handler.Check)
	router.HEAD("/health", handler.Check)
	router.GET("/live", handler.Live)
	router.HEAD("/live", handler.Live)
}
