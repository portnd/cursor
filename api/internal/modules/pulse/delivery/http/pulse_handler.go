package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
)

type pulseHandler struct {
	usecase domain.PulseUsecase
}

func newPulseHandler(uc domain.PulseUsecase) *pulseHandler {
	return &pulseHandler{usecase: uc}
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

	pulse, err := h.usecase.GetDailyCompanyPulse(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pulse)
}
