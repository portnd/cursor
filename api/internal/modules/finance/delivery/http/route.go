package http

import (
	"github.com/gin-gonic/gin"
	financeDomain "github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
)

// RegisterRoutes registers finance module routes (expects auth middleware on router).
func RegisterRoutes(router *gin.RouterGroup, usecase financeDomain.Usecase) {
	handler := NewFinanceHandler(usecase)
	fin := router.Group("/finance")
	{
		fin.POST("/entries", handler.CreateOrUpdateEntry)
		fin.GET("/entries", handler.ListEntries)
		fin.GET("/summary", handler.GetSummary)
	}
}
