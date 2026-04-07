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
		g.POST("/leaves", h.createLeaveRequest)
		g.GET("/leaves/my", h.listMyLeaves)
		g.GET("/leaves/balance", h.myLeaveBalance)
		g.GET("/leaves/notifications", h.myLeaveNotifications)
		g.PATCH("/leaves/notifications/:id/read", h.markMyLeaveNotificationRead)
		g.GET("/holidays", h.listHolidays)

		admin := g.Group("/admin")
		{
			admin.GET("/config", h.adminGetConfig)
			admin.PUT("/config", h.adminPutConfig)
			admin.GET("/records", h.adminRecords)
			admin.DELETE("/records/:id", h.adminDeleteRecord)
			admin.GET("/leaves/pending", h.listPendingLeaves)
			admin.PATCH("/leaves/:id/review", h.reviewLeave)
			admin.GET("/leaves/policies", h.adminListLeavePolicies)
			admin.PUT("/leaves/policies", h.adminUpsertLeavePolicy)
			admin.GET("/leaves/:id/audit", h.adminLeaveAudit)
			admin.GET("/leaves/trend", h.adminLeaveTrend)
			admin.POST("/leaves/backfill", h.adminBackfillLeave)
			admin.POST("/leaves/backfill/bulk", h.adminBackfillLeaveBulk)
			admin.PUT("/holidays", h.adminUpsertHoliday)
		}
	}
}
