package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	financeDomain "github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
)

// Handler handles finance HTTP requests (accounting entries + CEO summary).
type Handler struct {
	usecase financeDomain.Usecase
}

// NewFinanceHandler creates the finance HTTP handler.
func NewFinanceHandler(usecase financeDomain.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func getRole(c *gin.Context) string {
	role, _ := c.Get("role")
	if s, ok := role.(string); ok {
		return s
	}
	return ""
}

// CreateOrUpdateEntry upserts a monthly accounting entry (CEO or accountant).
// POST /api/v1/finance/entries
func (h *Handler) CreateOrUpdateEntry(c *gin.Context) {
	if getRole(c) != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}
	var req financeDomain.CreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	entry, err := h.usecase.CreateOrUpdateEntry(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// ListEntries returns monthly entries for the accounting page.
// GET /api/v1/finance/entries?limit=24
func (h *Handler) ListEntries(c *gin.Context) {
	if getRole(c) != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}
	limit := 24
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 60 {
			limit = n
		}
	}
	entries, err := h.usecase.ListEntries(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entries})
}

// GetSummary returns computed finance summary for CEO dashboard.
// GET /api/v1/finance/summary
func (h *Handler) GetSummary(c *gin.Context) {
	if getRole(c) != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}
	summary, err := h.usecase.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
