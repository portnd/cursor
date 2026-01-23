package http

import (
	"github.com/gin-gonic/gin"
	"github.com/komgrip/starter-kit/internal/modules/auth/domain"
)

// RegisterRoutes registers all authentication routes
// This follows the Dependency Injection pattern - all dependencies are passed in
func RegisterRoutes(router *gin.RouterGroup, usecase domain.Usecase) {
	// Create handler with usecase dependency
	handler := NewAuthHandler(usecase)

	// Auth routes group
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
	}
}
