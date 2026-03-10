package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
)

// RegisterRoutes mounts all pulse endpoints onto the provided router group.
func RegisterRoutes(r *gin.RouterGroup, uc domain.PulseUsecase) {
	h := newPulseHandler(uc)

	pulse := r.Group("/pulse")
	{
		pulse.POST("/standup", h.submitStandup)
		pulse.GET("/daily", h.getDailyPulse)
	}
}
