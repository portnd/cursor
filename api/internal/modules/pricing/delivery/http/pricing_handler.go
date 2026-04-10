package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
)

type pricingHandler struct {
	usecase domain.Usecase
}

func newPricingHandler(uc domain.Usecase) *pricingHandler {
	return &pricingHandler{usecase: uc}
}

// ─── Quotation endpoints ──────────────────────────────────────────────────────

// Calculate godoc — POST /api/v1/sentinel/projects/:id/quotation/calculate
func (h *pricingHandler) Calculate(c *gin.Context) {
	startedAt := time.Now()
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project id is required"})
		return
	}
	var req domain.QuotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}

	result, err := h.usecase.CalculateQuotation(projectID, &req)
	elapsedMs := time.Since(startedAt).Milliseconds()
	if err != nil {
		log.Printf("[pricing][quotation] calculate_error project_id=%s task_ids=%d epic_ids=%d dev_ids=%d total_ms=%d error=%v", projectID, len(req.TaskIDs), len(req.EpicIDs), len(req.DevUserIDs), elapsedMs, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[pricing][quotation] calculate_success project_id=%s task_ids=%d epic_ids=%d dev_ids=%d result_tasks=%d total_mandays=%.2f total_ms=%d", projectID, len(req.TaskIDs), len(req.EpicIDs), len(req.DevUserIDs), len(result.Tasks), result.TotalMandays, elapsedMs)
	c.JSON(http.StatusOK, result)
}

// Export godoc — POST /api/v1/sentinel/projects/:id/quotation/export
func (h *pricingHandler) Export(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project id is required"})
		return
	}
	var req domain.QuotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	pdfBytes, err := h.usecase.ExportQuotationPDF(projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("quotation-%s.pdf", projectID)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportMA godoc — POST /api/v1/sentinel/projects/:id/ma-quotation/export
func (h *pricingHandler) ExportMA(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project id is required"})
		return
	}
	var req domain.MAQuotationExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	pdfBytes, err := h.usecase.ExportMAQuotationPDF(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("ma-quotation-%s.pdf", projectID)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportCostReport godoc — POST /pricing/report/export
func (h *pricingHandler) ExportCostReport(c *gin.Context) {
	var req domain.CostReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	if req.Period == "" {
		req.Period = "all"
	}
	pdfBytes, err := h.usecase.GenerateCostReport(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := "cost-analysis-report.pdf"
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}



// GetConfig godoc — GET /pricing/config
func (h *pricingHandler) GetConfig(c *gin.Context) {
	cfg, err := h.usecase.GetCostConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

// UpdateConfig godoc — PUT /pricing/config
func (h *pricingHandler) UpdateConfig(c *gin.Context) {
	var req domain.UpdateCostConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	cfg, err := h.usecase.UpdateCostConfig(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

// ─── Admin: Employee Salaries ─────────────────────────────────────────────────

// ListSalaries godoc — GET /pricing/salaries
func (h *pricingHandler) ListSalaries(c *gin.Context) {
	salaries, err := h.usecase.ListAllSalaries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": salaries})
}

// UpsertSalary godoc — POST /pricing/salaries
func (h *pricingHandler) UpsertSalary(c *gin.Context) {
	var req domain.UpsertSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	sal, err := h.usecase.UpsertSalary(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sal)
}

// DeleteSalary godoc — DELETE /pricing/salaries/:salaryId
func (h *pricingHandler) DeleteSalary(c *gin.Context) {
	idStr := c.Param("salaryId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid salary id"})
		return
	}
	if err := h.usecase.DeleteSalary(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetCompanyMandayRate godoc — GET /pricing/manday-rate
// Returns the company-wide fully loaded cost per manday, derived from
// all active employee salaries + company_expense + executive_expense in company_cost_configs.
func (h *pricingHandler) GetCompanyMandayRate(c *gin.Context) {
	rate, err := h.usecase.GetCompanyMandayRate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rate)
}
