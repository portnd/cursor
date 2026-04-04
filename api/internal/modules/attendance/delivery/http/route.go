package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
)

// RegisterRoutes mounts attendance endpoints on the given router group (already behind auth).
func RegisterRoutes(r *gin.RouterGroup, uc domain.AttendanceUsecase) {
	h := newAttendanceHandler(uc)

	g := r.Group("/attendance")
	{
		g.POST("/check-in", h.checkIn)
		g.POST("/check-out", h.checkOut)
		g.GET("/today", h.today)
		g.GET("/history", h.history)

		admin := g.Group("/admin")
		{
			admin.GET("/config", h.adminGetConfig)
			admin.PUT("/config", h.adminPutConfig)
			admin.GET("/records", h.adminRecords)
		}
	}
}
