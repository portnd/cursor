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
	case uint:
		return v, true
	case int:
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
	case errors.Is(err, domain.ErrInvalidSchedule), errors.Is(err, domain.ErrInvalidDateRange):
		return err.Error(), http.StatusBadRequest
	case errors.Is(err, domain.ErrLeaveNotFound), errors.Is(err, domain.ErrUserNotFound):
		return err.Error(), http.StatusNotFound
	case errors.Is(err, domain.ErrLeaveNotPending):
		return err.Error(), http.StatusConflict
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

// POST /api/v1/attendance/leaves
func (h *attendanceHandler) createLeaveRequest(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req domain.CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.CreateLeaveRequest(userID, &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// GET /api/v1/attendance/leaves/my
func (h *attendanceHandler) listMyLeaves(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	items, err := h.uc.ListMyLeaveRequests(userID)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.LeaveListResponse{Items: items})
}

// GET /api/v1/attendance/admin/leaves/pending
func (h *attendanceHandler) listPendingLeaves(c *gin.Context) {
	role := roleFromContext(c)
	items, err := h.uc.ListPendingLeaveRequests(role)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.LeaveListResponse{Items: items})
}

// PATCH /api/v1/attendance/admin/leaves/:id/review
func (h *attendanceHandler) reviewLeave(c *gin.Context) {
	role := roleFromContext(c)
	approverID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	leaveID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || leaveID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leave id"})
		return
	}
	var req domain.ReviewLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, uerr := h.uc.ReviewLeaveRequest(role, approverID, leaveID, &req)
	if uerr != nil {
		respondAttendanceErr(c, uerr)
		return
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/v1/attendance/leaves/balance?year=YYYY
func (h *attendanceHandler) myLeaveBalance(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	year := time.Now().Year()
	if ys := c.Query("year"); ys != "" {
		if y, err := strconv.Atoi(ys); err == nil && y > 1900 {
			year = y
		}
	}
	items, err := h.uc.GetLeaveBalanceSummary(userID, year)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"year": year, "items": items})
}

// GET /api/v1/attendance/admin/leaves/policies
func (h *attendanceHandler) adminListLeavePolicies(c *gin.Context) {
	items, err := h.uc.ListLeavePolicies(roleFromContext(c))
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// PUT /api/v1/attendance/admin/leaves/policies
func (h *attendanceHandler) adminUpsertLeavePolicy(c *gin.Context) {
	var req domain.LeavePolicyUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.UpsertLeavePolicy(roleFromContext(c), &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/v1/attendance/holidays?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *attendanceHandler) listHolidays(c *gin.Context) {
	from := time.Now().UTC().AddDate(0, -1, 0)
	to := time.Now().UTC().AddDate(0, 6, 0)
	if fs := c.Query("from"); fs != "" {
		if d, err := time.Parse("2006-01-02", fs); err == nil {
			from = d.UTC()
		}
	}
	if ts := c.Query("to"); ts != "" {
		if d, err := time.Parse("2006-01-02", ts); err == nil {
			to = d.UTC()
		}
	}
	items, err := h.uc.ListHolidayCalendars(roleFromContext(c), from, to)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// PUT /api/v1/attendance/admin/holidays
func (h *attendanceHandler) adminUpsertHoliday(c *gin.Context) {
	var req domain.HolidayUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.UpsertHolidayCalendar(roleFromContext(c), &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/v1/attendance/admin/leaves/:id/audit
func (h *attendanceHandler) adminLeaveAudit(c *gin.Context) {
	leaveID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || leaveID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leave id"})
		return
	}
	items, uerr := h.uc.ListLeaveAuditLogs(roleFromContext(c), leaveID)
	if uerr != nil {
		respondAttendanceErr(c, uerr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GET /api/v1/attendance/leaves/notifications?unread_only=true
func (h *attendanceHandler) myLeaveNotifications(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	unreadOnly := c.Query("unread_only") == "true"
	items, err := h.uc.ListMyNotifications(userID, unreadOnly)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// PATCH /api/v1/attendance/leaves/notifications/:id/read
func (h *attendanceHandler) markMyLeaveNotificationRead(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}
	if err := h.uc.MarkMyNotificationRead(userID, id); err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /api/v1/attendance/admin/leaves/trend?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *attendanceHandler) adminLeaveTrend(c *gin.Context) {
	from := time.Now().UTC().AddDate(0, -11, 0)
	to := time.Now().UTC()
	if fs := c.Query("from"); fs != "" {
		if d, err := time.Parse("2006-01-02", fs); err == nil {
			from = d.UTC()
		}
	}
	if ts := c.Query("to"); ts != "" {
		if d, err := time.Parse("2006-01-02", ts); err == nil {
			to = d.UTC()
		}
	}
	items, err := h.uc.GetLeaveTrend(roleFromContext(c), from, to)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.LeaveTrendResponse{Items: items})
}

// POST /api/v1/attendance/admin/leaves/backfill
func (h *attendanceHandler) adminBackfillLeave(c *gin.Context) {
	role := roleFromContext(c)
	actorID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req domain.LeaveBackfillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.BackfillLeave(role, actorID, &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// POST /api/v1/attendance/admin/leaves/backfill/bulk
func (h *attendanceHandler) adminBackfillLeaveBulk(c *gin.Context) {
	role := roleFromContext(c)
	actorID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req domain.LeaveBackfillBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.BackfillLeaveBulk(role, actorID, &req)
	if err != nil {
		respondAttendanceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
