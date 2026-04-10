package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
)

// RegisterRoutes mounts all pricing endpoints under the provided router group.
//
// Quotation (per-project):
//
//	POST /sentinel/projects/:id/quotation/calculate → JSON breakdown
//	POST /sentinel/projects/:id/quotation/export    → PDF download
//
// Admin (company-wide cost settings):
//
//	GET  /pricing/config           → current CompanyCostConfig
//	PUT  /pricing/config           → update CompanyCostConfig
//	GET  /pricing/salaries         → list all EmployeeSalary records (with user info)
//	POST /pricing/salaries         → upsert a salary record
//	DELETE /pricing/salaries/:salaryId → delete a salary record
//	POST /pricing/report/export    → Cost Analysis Report PDF
//	GET  /pricing/manday-rate      → company-wide fully loaded manday rate
func RegisterRoutes(r *gin.RouterGroup, uc domain.Usecase) {
	h := newPricingHandler(uc)

	// Quotation endpoints
	q := r.Group("/sentinel/projects/:id/quotation")
	{
		q.POST("/calculate", h.Calculate)
		q.POST("/export", h.Export)
	}

	// MA Quotation PDF export
	r.POST("/sentinel/projects/:id/ma-quotation/export", h.ExportMA)

	// Admin cost config & salary endpoints
	p := r.Group("/pricing")
	{
		p.GET("/config", h.GetConfig)
		p.PUT("/config", h.UpdateConfig)
		p.GET("/salaries", h.ListSalaries)
		p.POST("/salaries", h.UpsertSalary)
		p.DELETE("/salaries/:salaryId", h.DeleteSalary)
		p.POST("/report/export", h.ExportCostReport)
		p.GET("/manday-rate", h.GetCompanyMandayRate)
	}
}
