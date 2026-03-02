package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
)

// RegisterRoutes registers all authentication routes
// This follows the Dependency Injection pattern - all dependencies are passed in
// authMiddleware is optional - pass nil if you want to set it up differently
func RegisterRoutes(router *gin.RouterGroup, usecase domain.Usecase, authMiddleware gin.HandlerFunc) {
	// Create handler with usecase dependency
	handler := NewAuthHandler(usecase)

	// Auth routes group
	authGroup := router.Group("/auth")
	{
		// Public routes (no authentication required)
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)

		// Protected routes (authentication required)
		if authMiddleware != nil {
			// User Management endpoints (CEO only - authorization checked in usecase)
			authGroup.GET("/users", authMiddleware, handler.GetTeam)
			authGroup.PATCH("/users/:id/role", authMiddleware, handler.ChangeRole)
		}
	}
}
