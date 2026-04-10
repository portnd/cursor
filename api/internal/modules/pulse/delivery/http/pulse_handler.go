package http

import (
	"errors"
	"net/http"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
)

type pulseHandler struct {
	usecase domain.PulseUsecase
}

func newPulseHandler(uc domain.PulseUsecase) *pulseHandler {
	return &pulseHandler{usecase: uc}
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

// POST /api/v1/pulse/standup
func (h *pulseHandler) submitStandup(c *gin.Context) {
	rawUID, _ := c.Get("user_id")
	userID := uint(rawUID.(float64))

	var req struct {
		Date             string   `json:"date" binding:"required"` // YYYY-MM-DD
		YesterdaySummary string   `json:"yesterday_summary" binding:"required"`
		TodayTaskIDs     []string `json:"today_task_ids"`
		Blocker          string   `json:"blocker"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be YYYY-MM-DD"})
		return
	}

	if req.TodayTaskIDs == nil {
		req.TodayTaskIDs = []string{}
	}

	standup, err := h.usecase.SubmitStandup(userID, date, req.YesterdaySummary, req.Blocker, req.TodayTaskIDs)
	if err != nil {
		if errors.Is(err, domain.ErrStandupNotRequiredForRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, standup)
}

// GET /api/v1/pulse/daily?date=YYYY-MM-DD
func (h *pulseHandler) getDailyPulse(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date must be YYYY-MM-DD"})
			return
		}
	}

	role := getRole(c)
	pulse, err := h.usecase.GetDailyCompanyPulse(date, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pulse)
}

// GET /api/v1/pulse/hidden-users
func (h *pulseHandler) getHiddenUsers(c *gin.Context) {
	role := getRole(c)
	ids, err := h.usecase.GetHiddenPulseUserIDs(role)
	if err != nil {
		if errors.Is(err, domain.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_ids": ids})
}

// PUT /api/v1/pulse/hidden-users
func (h *pulseHandler) setHiddenUsers(c *gin.Context) {
	role := getRole(c)
	if role != authDomain.RoleCEO {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
		return
	}
	var req struct {
		UserIDs []uint `json:"user_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid request body"})
		return
	}
	if req.UserIDs == nil {
		req.UserIDs = []uint{}
	}
	if err := h.usecase.SetHiddenPulseUserIDs(role, req.UserIDs); err != nil {
		if errors.Is(err, domain.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "CEO only"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_ids": req.UserIDs})
}
