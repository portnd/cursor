package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
)

// RegisterRoutes registers all authentication routes
// This follows the Dependency Injection pattern - all dependencies are passed in
// authMiddleware is optional - pass nil if you want to set it up differently
func RegisterRoutes(router *gin.RouterGroup, usecase domain.Usecase, teamFinanceUsecase domain.TeamFinanceUsecase, authMiddleware gin.HandlerFunc) {
	// Create handler with usecase dependency
	handler := NewAuthHandler(usecase)
	financeHandler := NewTeamFinanceHandler(teamFinanceUsecase)

	// Auth routes group
	authGroup := router.Group("/auth")
	{
		// Public routes (no authentication required)
		// NOTE: /register is intentionally removed — user accounts are created by CEO only via POST /auth/users
		authGroup.POST("/login", handler.Login)

		// Protected routes (authentication required)
		if authMiddleware != nil {
			// Profile (any authenticated user)
			authGroup.GET("/me", authMiddleware, handler.GetMe)
			authGroup.PATCH("/me", authMiddleware, handler.UpdateProfile)
			authGroup.PATCH("/me/password", authMiddleware, handler.ChangePassword)
			// User Management endpoints (CEO only - authorization checked in usecase)
			authGroup.GET("/users", authMiddleware, handler.GetTeam)
			authGroup.POST("/users/import", authMiddleware, handler.ImportUsers)
			authGroup.POST("/users", authMiddleware, handler.CreateUser)
			authGroup.DELETE("/users/:id", authMiddleware, handler.DeleteUser)
			authGroup.PATCH("/users/:id/password", authMiddleware, handler.ResetPassword)
			authGroup.PATCH("/users/:id/role", authMiddleware, handler.ChangeRole)
			// Team / Squad management (CEO only)
			authGroup.GET("/teams", authMiddleware, handler.GetTeams)
			authGroup.POST("/teams", authMiddleware, handler.CreateTeam)
			authGroup.PATCH("/teams/:id", authMiddleware, handler.UpdateTeam)
			authGroup.DELETE("/teams/:id", authMiddleware, handler.DeleteTeam)
			authGroup.PATCH("/users/:id/assign-team", authMiddleware, handler.AssignUserToTeam)
			// Teams feature flag (GET: any auth, PUT: CEO/MANAGER)
			authGroup.GET("/settings/teams-feature", authMiddleware, handler.GetTeamsFeature)
			authGroup.PUT("/settings/teams-feature", authMiddleware, handler.SetTeamsFeature)
			// Team Finance / Internal VC model (CEO only)
			authGroup.GET("/teams/:id/finance/cost", authMiddleware, financeHandler.GetTeamMonthlyCost)
			authGroup.POST("/teams/:id/finance/inject", authMiddleware, financeHandler.InjectCapital)
			authGroup.PUT("/teams/:id/finance/capital", authMiddleware, financeHandler.EditCapital)
			authGroup.POST("/teams/:id/finance/close-cycle", authMiddleware, financeHandler.CloseCycleAndPayout)
		}
	}
}
