package http

import (
	"net/http"
	"strconv"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
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

// GetTeam returns team KPIs (CEO + Product Owner only)
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

// GetDiscipline returns daily discipline stats for all employees in a date range.
// GET /api/v1/performance/discipline?from=YYYY-MM-DD&to=YYYY-MM-DD (CEO + Product Owner + Manager + Engineer + Chief Engineer + Support)
func (h *Handler) GetDiscipline(c *gin.Context) {
	role := getRole(c)
	if role != "CEO" && role != authDomain.RoleProductOwner && role != authDomain.RoleManager && role != authDomain.RoleEngineer && role != authDomain.RoleChiefEngineer && role != authDomain.RoleSupport {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO, Product Owner, Manager, Engineer, Chief Engineer and Support only"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "from and to query params required (YYYY-MM-DD)"})
		return
	}

	resp, err := h.usecase.GetDiscipline(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetDisciplineDayDetail returns drill-down activity for one user on one day.
// GET /api/v1/performance/discipline/detail?user_id=X&date=YYYY-MM-DD (CEO + Product Owner + Manager + Engineer + Chief Engineer + Support)
func (h *Handler) GetDisciplineDayDetail(c *gin.Context) {
	role := getRole(c)
	if role != "CEO" && role != authDomain.RoleProductOwner && role != authDomain.RoleManager && role != authDomain.RoleEngineer && role != authDomain.RoleChiefEngineer && role != authDomain.RoleSupport {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO, Product Owner, Manager, Engineer, Chief Engineer and Support only"})
		return
	}

	userIDStr := c.Query("user_id")
	date := c.Query("date")
	if userIDStr == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "user_id and date required"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid user_id"})
		return
	}

	detail, err := h.usecase.GetDisciplineDayDetail(uint(userID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// GetDisciplineStartDate returns configured global start date for discipline data.
// GET /api/v1/performance/discipline/start-date
func (h *Handler) GetDisciplineStartDate(c *gin.Context) {
	role := getRole(c)
	if role != "CEO" && role != authDomain.RoleProductOwner && role != authDomain.RoleManager && role != authDomain.RoleEngineer && role != authDomain.RoleChiefEngineer && role != authDomain.RoleSupport {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO, Product Owner, Manager, Engineer, Chief Engineer and Support only"})
		return
	}

	resp, err := h.usecase.GetDisciplineStartDate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// SetDisciplineStartDate sets configured global start date for discipline data.
// PUT /api/v1/performance/discipline/start-date (CEO only)
func (h *Handler) SetDisciplineStartDate(c *gin.Context) {
	role := getRole(c)
	if role != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}

	var req struct {
		StartDate string `json:"start_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid request body"})
		return
	}
	if req.StartDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "start_date is required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "start_date must be YYYY-MM-DD"})
		return
	}

	resp, err := h.usecase.SetDisciplineStartDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
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
