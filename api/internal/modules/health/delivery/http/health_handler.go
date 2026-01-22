package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db          *gorm.DB
	mongoClient *mongo.Client
	redisClient *redis.Client
}

func NewHealthHandler(db *gorm.DB, mongoClient *mongo.Client, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		mongoClient: mongoClient,
		redisClient: redisClient,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	services := make(map[string]string)

	postgresStatus := "UP"
	if sqlDB, err := h.db.DB(); err != nil {
		postgresStatus = "DOWN"
	} else if err := sqlDB.PingContext(ctx); err != nil {
		postgresStatus = "DOWN"
	}
	services["postgres"] = postgresStatus

	mongoStatus := "UP"
	if err := h.mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		mongoStatus = "DOWN"
	}
	services["mongodb"] = mongoStatus

	redisStatus := "UP"
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		redisStatus = "DOWN"
	}
	services["redis"] = redisStatus

	overallStatus := "UP"
	for _, status := range services {
		if status == "DOWN" {
			overallStatus = "DEGRADED"
			break
		}
	}

	statusCode := http.StatusOK
	if overallStatus == "DEGRADED" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	})
}
