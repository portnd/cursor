package http

import (
	"github.com/gin-gonic/gin"
	perfDomain "github.com/portnd/the-sentinel-core/internal/modules/performance/domain"
)

// RegisterRoutes registers performance module routes (expects auth middleware on router)
func RegisterRoutes(router *gin.RouterGroup, usecase perfDomain.Usecase) {
	handler := NewPerformanceHandler(usecase)
	perf := router.Group("/performance")
	{
		perf.GET("/me", handler.GetMe)
		perf.GET("/team", handler.GetTeam)
		perf.GET("/overview", handler.GetOverview)
		perf.POST("/users/:id/reset-rework", handler.ResetReworkRate) // CEO only
	}
}
