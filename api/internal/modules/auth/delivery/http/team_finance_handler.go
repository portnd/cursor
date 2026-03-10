package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
)

// TeamFinanceHandler handles HTTP requests for the Internal VC finance model
type TeamFinanceHandler struct {
	usecase domain.TeamFinanceUsecase
}

// NewTeamFinanceHandler creates a new TeamFinanceHandler
func NewTeamFinanceHandler(usecase domain.TeamFinanceUsecase) *TeamFinanceHandler {
	return &TeamFinanceHandler{usecase: usecase}
}

// GetTeamMonthlyCost handles GET /auth/teams/:id/finance/cost
// Returns the fully loaded monthly burn rate for the team, plus runway information.
func (h *TeamFinanceHandler) GetTeamMonthlyCost(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	result, err := h.usecase.CalculateTeamMonthlyCost(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// InjectCapital handles POST /auth/teams/:id/finance/inject
// Body: { "amount": 500000, "bonus_percentage": 20, "note": "งวดที่ 1 MIMS HD-MAP" }
func (h *TeamFinanceHandler) InjectCapital(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var req domain.InjectCapitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.usecase.InjectCapital(uint(teamID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": team})
}

// CloseCycleAndPayout handles POST /auth/teams/:id/finance/close-cycle
// Calculates bonus from remaining balance, records BONUS_PAYOUT, resets balance to 0.
func (h *TeamFinanceHandler) CloseCycleAndPayout(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	result, err := h.usecase.CloseCycleAndPayout(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// EditCapital handles PUT /auth/teams/:id/finance/capital
// Body: { "new_balance": 250000, "bonus_percentage": 20, "note": "แก้ไขยอด" }
// Sets the capital balance to an exact value and records an ADJUSTMENT transaction.
func (h *TeamFinanceHandler) EditCapital(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var req domain.EditCapitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.usecase.EditCapital(uint(teamID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": team})
}
