package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
)

type attendanceHandler struct {
	uc domain.AttendanceUsecase
}

func newAttendanceHandler(uc domain.AttendanceUsecase) *attendanceHandler {
	return &attendanceHandler{uc: uc}
}

func userIDFromContext(c *gin.Context) (uint, bool) {
	raw, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	switch v := raw.(type) {
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	default:
		return 0, false
	}
}

func roleFromContext(c *gin.Context) string {
	raw, ok := c.Get("role")
	if !ok {
		return ""
	}
	s, _ := raw.(string)
	return s
}

func respondAttendanceErr(c *gin.Context, err error) {
	code, status := mapAttendanceErr(err)
	c.JSON(status, gin.H{"error": code})
}

func mapAttendanceErr(err error) (string, int) {
	switch {
	case errors.Is(err, domain.ErrOutsideOffice):
		return err.Error(), http.StatusForbidden
	case errors.Is(err, domain.ErrNoOfficeConfig):
		return err.Error(), http.StatusServiceUnavailable
	case errors.Is(err, domain.ErrAlreadyCheckedIn):
		return err.Error(), http.StatusConflict
	case errors.Is(err, domain.ErrNotCheckedIn):
		return err.Error(), http.StatusConflict
	case errors.Is(err, domain.ErrAlreadyCheckedOut):
		return err.Error(), http.StatusConflict
	case errors.Is(err, domain.ErrNotWorkDay):
		return err.Error(), http.StatusForbidden
	case errors.Is(err, domain.ErrForbiddenAdmin):
		return err.Error(), http.StatusForbidden
	case errors.Is(err, domain.ErrInvalidCursor):
		return err.Error(), http.StatusBadRequest
	case errors.Is(err, domain.ErrInvalidSchedule):
		return err.Error(), http.StatusBadRequest
	default:
		return err.Error(), http.StatusInternalServerError
	}
}

// POST /api/v1/attendance/check-in
func (h *attendanceHandler) checkIn(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		Lat float64 `json:"lat" binding:"required"`
		Lng float64 `json:"lng" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := h.uc.CheckIn(userID, req.Lat, req.Lng, c.ClientIP())
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rec)
}

// POST /api/v1/attendance/check-out
func (h *attendanceHandler) checkOut(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	rec, err := h.uc.CheckOut(userID)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, rec)
}

// GET /api/v1/attendance/today
func (h *attendanceHandler) today(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	rec, cfg, err := h.uc.GetTodayStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"record":        rec,
		"office_config": cfg,
	})
}

// GET /api/v1/attendance/history?cursor=&limit=
func (h *attendanceHandler) history(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit := 20
	if ls := c.Query("limit"); ls != "" {
		if n, err := strconv.Atoi(ls); err == nil {
			limit = n
		}
	}
	resp, err := h.uc.GetHistory(userID, c.Query("cursor"), limit)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GET /api/v1/attendance/admin/config
func (h *attendanceHandler) adminGetConfig(c *gin.Context) {
	cfg, err := h.uc.GetOfficeConfigForAdmin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

// PUT /api/v1/attendance/admin/config
func (h *attendanceHandler) adminPutConfig(c *gin.Context) {
	role := roleFromContext(c)
	var req domain.UpsertOfficeConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, err := h.uc.UpsertOfficeConfig(role, &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

// GET /api/v1/attendance/admin/records?date=YYYY-MM-DD
func (h *attendanceHandler) adminRecords(c *gin.Context) {
	role := roleFromContext(c)
	ds := c.Query("date")
	if ds == "" {
		ds = time.Now().Format("2006-01-02")
	}
	d, err := time.Parse("2006-01-02", ds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be YYYY-MM-DD"})
		return
	}
	rows, err := h.uc.ListAdminRecordsByDate(role, d)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"date": ds, "records": rows})
}
