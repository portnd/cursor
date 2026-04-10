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
		perf.GET("/discipline", handler.GetDiscipline)                        // CEO + Product Owner + Manager + Engineer + Chief Engineer + Support
		perf.GET("/discipline/detail", handler.GetDisciplineDayDetail)          // CEO + Product Owner + Manager + Engineer + Chief Engineer + Support
		perf.GET("/discipline/start-date", handler.GetDisciplineStartDate)      // all discipline viewers
		perf.PUT("/discipline/start-date", handler.SetDisciplineStartDate)      // CEO only
		perf.POST("/users/:id/reset-rework", handler.ResetReworkRate)          // CEO only
	}
}
