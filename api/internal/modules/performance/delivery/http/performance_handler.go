package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	perfDomain "github.com/portnd/the-sentinel-core/internal/modules/performance/domain"
)

// Handler handles performance KPI HTTP requests
type Handler struct {
	usecase perfDomain.Usecase
}

// NewPerformanceHandler creates the performance HTTP handler
func NewPerformanceHandler(usecase perfDomain.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func getUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if uid, ok := userID.(float64); ok {
		return uint(uid)
	}
	if uid, ok := userID.(uint); ok {
		return uid
	}
	if uid, ok := userID.(int); ok {
		return uint(uid)
	}
	return 0
}

func getRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	if s, ok := role.(string); ok {
		return s
	}
	return ""
}

// GetMe returns personal KPIs for the authenticated user
// GET /api/v1/performance/me
func (h *Handler) GetMe(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	role := getRole(c)
	kpis, err := h.usecase.GetPersonalKPIs(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kpis)
}

// GetTeam returns team KPIs (CEO + PM only)
// GET /api/v1/performance/team
func (h *Handler) GetTeam(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	role := getRole(c)
	resp, err := h.usecase.GetTeamKPIs(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetOverview returns company-level overview KPIs (CEO only)
// GET /api/v1/performance/overview
func (h *Handler) GetOverview(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	role := getRole(c)
	overview, err := h.usecase.GetOverviewKPIs(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	if overview == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}
	c.JSON(http.StatusOK, overview)
}

// ResetReworkRate resets a developer's Rework Rate by setting rework_reset_at = NOW().
// POST /api/v1/performance/users/:id/reset-rework (CEO only)
func (h *Handler) ResetReworkRate(c *gin.Context) {
	requesterRole := getRole(c)
	if requesterRole != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO can reset rework rate"})
		return
	}

	devID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || devID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid user id"})
		return
	}

	if err := h.usecase.ResetReworkRate(uint(devID), requesterRole); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rework rate reset successfully — counter starts from now"})
}
