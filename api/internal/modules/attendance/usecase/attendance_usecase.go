package usecase

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"gorm.io/datatypes"
)

const lateGraceMinutes = 15

type attendanceUsecase struct {
	repo   domain.AttendanceRepository
	secret string
	loc    *time.Location
}

// NewAttendanceUsecase wires attendance business logic.
func NewAttendanceUsecase(repo domain.AttendanceRepository, jwtSecret string) domain.AttendanceUsecase {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	return &attendanceUsecase{
		repo:   repo,
		secret: jwtSecret,
		loc:    loc,
	}
}

func (u *attendanceUsecase) CheckIn(userID uint, lat, lng float64, clientIP string) (*domain.AttendanceRecord, error) {
	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, domain.ErrNoOfficeConfig
	}

	now := time.Now().In(u.loc)
	if !u.isAttendanceRequiredDay(cfg, now) {
		return nil, domain.ErrNotWorkDay
	}

	attDate := calendarDateUTC(now)

	existing, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.CheckInAt != nil {
		return nil, domain.ErrAlreadyCheckedIn
	}

	wfhToday := u.isWFHDay(cfg, now)
	var method string
	if wfhToday {
		method = "wfh"
	} else {
		gpsOK := u.gpsWithinOffice(lat, lng, cfg)
		if !gpsOK {
			return nil, domain.ErrOutsideOfficeGPS
		}
		method = "gps"
	}

	halfDaySession, err := u.getApprovedHalfDaySessionForDate(userID, attDate)
	if err != nil {
		return nil, err
	}

	workStart, err := u.workStartOnDate(cfg, now)
	if err != nil {
		return nil, fmt.Errorf("work_start_time: %w", err)
	}
	isLate := now.After(workStart.Add(lateGraceMinutes * time.Minute))
	if halfDaySession == "AM" {
		deadline := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, u.loc)
		isLate = !now.Before(deadline)
	}
	status := "present"
	if isLate {
		status = "late"
	}

	latCopy, lngCopy := lat, lng
	nowUTC := time.Now().UTC()
	rec := &domain.AttendanceRecord{
		UserID:         userID,
		OfficeConfigID: cfg.ID,
		AttendanceDate: attDate,
		CheckInAt:      &nowUTC,
		CheckInLat:     &latCopy,
		CheckInLng:     &lngCopy,
		CheckInMethod:  method,
		CheckInIP:      strings.TrimSpace(clientIP),
		IsLate:         isLate,
		EarlyCheckout:  false,
		Status:         status,
	}

	if err := u.repo.SaveRecord(rec); err != nil {
		return nil, err
	}
	saved, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (u *attendanceUsecase) CheckOut(userID uint) (*domain.AttendanceRecord, error) {
	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, domain.ErrNoOfficeConfig
	}

	now := time.Now().In(u.loc)
	attDate := calendarDateUTC(now)

	rec, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.CheckInAt == nil {
		return nil, domain.ErrNotCheckedIn
	}
	if rec.CheckOutAt != nil {
		return nil, domain.ErrAlreadyCheckedOut
	}

	halfDaySession, err := u.getApprovedHalfDaySessionForDate(userID, attDate)
	if err != nil {
		return nil, err
	}

	workEnd, err := u.workEndOnDate(cfg, now)
	if err != nil {
		return nil, fmt.Errorf("work_end_time: %w", err)
	}
	nowUTC := time.Now().UTC()
	early := now.In(u.loc).Before(workEnd)
	if halfDaySession == "PM" {
		minCheckout := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, u.loc)
		early = now.Before(minCheckout)
	}

	rec.CheckOutAt = &nowUTC
	rec.EarlyCheckout = early
	if err := u.repo.SaveRecord(rec); err != nil {
		return nil, err
	}
	return u.repo.GetRecordByUserAndDate(userID, attDate)
}

func (u *attendanceUsecase) GetTodayStatus(userID uint) (*domain.AttendanceRecord, *domain.OfficeConfig, error) {
	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().In(u.loc)
	attDate := calendarDateUTC(now)
	rec, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, cfg, err
	}
	return rec, cfg, nil
}

func (u *attendanceUsecase) GetTodayOffsiteCheckInRequest(userID uint) (*domain.OffsiteCheckInRequest, error) {
	now := time.Now().In(u.loc)
	attDate := calendarDateUTC(now)
	return u.repo.GetLatestOffsiteCheckInRequestByUserAndDate(userID, attDate)
}

func (u *attendanceUsecase) GetTodayOffsiteCheckOutRequest(userID uint) (*domain.OffsiteCheckOutRequest, error) {
	now := time.Now().In(u.loc)
	attDate := calendarDateUTC(now)
	return u.repo.GetLatestOffsiteCheckOutRequestByUserAndDate(userID, attDate)
}

func (u *attendanceUsecase) GetHistory(userID uint, cursor string, limit int) (*domain.AttendanceHistoryResponse, error) {
	limit = domain.CapAttendanceLimit(limit)
	afterID, err := domain.DecodeAttendanceCursor(cursor, u.secret)
	if err != nil {
		return nil, err
	}

	fetch := limit + 1
	rows, err := u.repo.ListUserRecordsAfterID(userID, afterID, fetch)
	if err != nil {
		return nil, err
	}

	resp := &domain.AttendanceHistoryResponse{Items: rows}
	if len(rows) > limit {
		resp.Items = rows[:limit]
		last := resp.Items[len(resp.Items)-1]
		next, encErr := domain.EncodeAttendanceCursor(last.ID, u.secret)
		if encErr != nil {
			return nil, encErr
		}
		resp.NextCursor = next
	}
	return resp, nil
}

func (u *attendanceUsecase) RequestOffsiteCheckIn(userID uint, lat, lng float64, reason string) (*domain.OffsiteCheckInRequest, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, domain.ErrOffsiteReasonRequired
	}

	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, domain.ErrNoOfficeConfig
	}

	now := time.Now().In(u.loc)
	if !u.isAttendanceRequiredDay(cfg, now) {
		return nil, domain.ErrNotWorkDay
	}
	if u.isWFHDay(cfg, now) {
		return nil, domain.ErrOffsiteWFHNotAllowed
	}
	if u.gpsWithinOffice(lat, lng, cfg) {
		return nil, domain.ErrOffsiteApprovalNotRequired
	}

	attDate := calendarDateUTC(now)
	rec, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if rec != nil && rec.CheckInAt != nil {
		return nil, domain.ErrAlreadyCheckedIn
	}

	existing, err := u.repo.GetLatestOffsiteCheckInRequestByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == domain.OffsiteStatusPending {
		return nil, domain.ErrOffsiteAlreadyRequested
	}

	item := &domain.OffsiteCheckInRequest{
		UserID:         userID,
		OfficeConfigID: cfg.ID,
		AttendanceDate: attDate,
		RequestLat:     lat,
		RequestLng:     lng,
		Reason:         reason,
		Status:         domain.OffsiteStatusPending,
		RequestedAt:    time.Now().UTC(),
	}
	if err := u.repo.CreateOffsiteCheckInRequest(item); err != nil {
		return nil, err
	}
	return u.repo.GetOffsiteCheckInRequestByID(item.ID)
}

func (u *attendanceUsecase) ListPendingOffsiteCheckInRequests(role string) ([]domain.OffsiteCheckInRequest, error) {
	if role != authDomain.RoleCEO {
		return nil, domain.ErrForbiddenCEOOnly
	}
	return u.repo.ListPendingOffsiteCheckInRequests()
}

func (u *attendanceUsecase) ReviewOffsiteCheckInRequest(role string, approverID uint, requestID int64, status, note string) (*domain.OffsiteCheckInRequest, error) {
	if role != authDomain.RoleCEO {
		return nil, domain.ErrForbiddenCEOOnly
	}
	if requestID <= 0 {
		return nil, domain.ErrOffsiteRequestNotFound
	}

	item, err := u.repo.GetOffsiteCheckInRequestByID(requestID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, domain.ErrOffsiteRequestNotFound
	}
	if item.Status != domain.OffsiteStatusPending {
		return nil, domain.ErrOffsiteRequestNotPending
	}

	status = strings.ToUpper(strings.TrimSpace(status))
	if status != domain.OffsiteStatusApproved && status != domain.OffsiteStatusRejected {
		return nil, domain.ErrOffsiteRequestNotPending
	}

	now := time.Now().UTC()
	item.Status = status
	item.ApproverID = &approverID
	item.ApproverNote = strings.TrimSpace(note)
	item.ApprovedAt = &now
	if err := u.repo.UpdateOffsiteCheckInRequest(item); err != nil {
		return nil, err
	}

	if status == domain.OffsiteStatusApproved {
		existingRec, err := u.repo.GetRecordByUserAndDate(item.UserID, item.AttendanceDate)
		if err != nil {
			return nil, err
		}
		if existingRec == nil || existingRec.CheckInAt == nil {
			cfg, err := u.repo.GetOfficeConfigByID(item.OfficeConfigID)
			if err != nil {
				return nil, err
			}
			if cfg == nil {
				return nil, domain.ErrNoOfficeConfig
			}
			reqLocal := item.RequestedAt.In(u.loc)
			workStart, err := u.workStartOnDate(cfg, reqLocal)
			if err != nil {
				return nil, fmt.Errorf("work_start_time: %w", err)
			}
			isLate := reqLocal.After(workStart.Add(lateGraceMinutes * time.Minute))
			recStatus := "present"
			if isLate {
				recStatus = "late"
			}
			latCopy, lngCopy := item.RequestLat, item.RequestLng
			checkInAt := item.RequestedAt.UTC()
			newRec := &domain.AttendanceRecord{
				UserID:         item.UserID,
				OfficeConfigID: item.OfficeConfigID,
				AttendanceDate: item.AttendanceDate,
				CheckInAt:      &checkInAt,
				CheckInLat:     &latCopy,
				CheckInLng:     &lngCopy,
				CheckInMethod:  "offsite_approved",
				CheckInIP:      "",
				IsLate:         isLate,
				EarlyCheckout:  false,
				Status:         recStatus,
			}
			if err := u.repo.SaveRecord(newRec); err != nil {
				return nil, err
			}
		}
	}

	return u.repo.GetOffsiteCheckInRequestByID(item.ID)
}

func (u *attendanceUsecase) RequestOffsiteCheckOut(userID uint, lat, lng float64, reason string) (*domain.OffsiteCheckOutRequest, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, domain.ErrOffsiteCheckOutReasonRequired
	}

	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, domain.ErrNoOfficeConfig
	}

	now := time.Now().In(u.loc)
	attDate := calendarDateUTC(now)

	rec, err := u.repo.GetRecordByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.CheckInAt == nil {
		return nil, domain.ErrNotCheckedIn
	}
	if rec.CheckOutAt != nil {
		return nil, domain.ErrAlreadyCheckedOut
	}

	existing, err := u.repo.GetLatestOffsiteCheckOutRequestByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == domain.OffsiteStatusPending {
		return nil, domain.ErrOffsiteCheckOutAlreadyRequested
	}

	item := &domain.OffsiteCheckOutRequest{
		UserID:         userID,
		OfficeConfigID: cfg.ID,
		AttendanceDate: attDate,
		RequestLat:     lat,
		RequestLng:     lng,
		Reason:         reason,
		Status:         domain.OffsiteStatusPending,
		RequestedAt:    time.Now().UTC(),
	}
	if err := u.repo.CreateOffsiteCheckOutRequest(item); err != nil {
		return nil, err
	}
	return u.repo.GetOffsiteCheckOutRequestByID(item.ID)
}

func (u *attendanceUsecase) ListPendingOffsiteCheckOutRequests(role string) ([]domain.OffsiteCheckOutRequest, error) {
	if role != authDomain.RoleCEO {
		return nil, domain.ErrForbiddenCEOOnly
	}
	return u.repo.ListPendingOffsiteCheckOutRequests()
}

func (u *attendanceUsecase) ReviewOffsiteCheckOutRequest(role string, approverID uint, requestID int64, status, note string) (*domain.OffsiteCheckOutRequest, error) {
	if role != authDomain.RoleCEO {
		return nil, domain.ErrForbiddenCEOOnly
	}
	if requestID <= 0 {
		return nil, domain.ErrOffsiteCheckOutRequestNotFound
	}

	item, err := u.repo.GetOffsiteCheckOutRequestByID(requestID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, domain.ErrOffsiteCheckOutRequestNotFound
	}
	if item.Status != domain.OffsiteStatusPending {
		return nil, domain.ErrOffsiteCheckOutRequestNotPending
	}

	status = strings.ToUpper(strings.TrimSpace(status))
	if status != domain.OffsiteStatusApproved && status != domain.OffsiteStatusRejected {
		return nil, domain.ErrOffsiteCheckOutRequestNotPending
	}

	now := time.Now().UTC()
	item.Status = status
	item.ApproverID = &approverID
	item.ApproverNote = strings.TrimSpace(note)
	item.ApprovedAt = &now
	if err := u.repo.UpdateOffsiteCheckOutRequest(item); err != nil {
		return nil, err
	}

	if status == domain.OffsiteStatusApproved {
		rec, err := u.repo.GetRecordByUserAndDate(item.UserID, item.AttendanceDate)
		if err != nil {
			return nil, err
		}
		if rec == nil || rec.CheckInAt == nil {
			return nil, domain.ErrNotCheckedIn
		}
		if rec.CheckOutAt == nil {
			cfg, err := u.repo.GetOfficeConfigByID(item.OfficeConfigID)
			if err != nil {
				return nil, err
			}
			if cfg == nil {
				return nil, domain.ErrNoOfficeConfig
			}

			halfDaySession, err := u.getApprovedHalfDaySessionForDate(item.UserID, item.AttendanceDate)
			if err != nil {
				return nil, err
			}

			reqLocal := item.RequestedAt.In(u.loc)
			workEnd, err := u.workEndOnDate(cfg, reqLocal)
			if err != nil {
				return nil, fmt.Errorf("work_end_time: %w", err)
			}
			early := reqLocal.Before(workEnd)
			if halfDaySession == "PM" {
				minCheckout := time.Date(reqLocal.Year(), reqLocal.Month(), reqLocal.Day(), 12, 0, 0, 0, u.loc)
				early = reqLocal.Before(minCheckout)
			}

			checkOutAt := item.RequestedAt.UTC()
			rec.CheckOutAt = &checkOutAt
			rec.EarlyCheckout = early
			if err := u.repo.SaveRecord(rec); err != nil {
				return nil, err
			}
		}
	}

	return u.repo.GetOffsiteCheckOutRequestByID(item.ID)
}

func (u *attendanceUsecase) GetOfficeConfigForAdmin() (*domain.OfficeConfig, error) {
	cfg, err := u.repo.GetActiveOfficeConfig()
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		return cfg, nil
	}
	return u.repo.GetFirstOfficeConfig()
}

func (u *attendanceUsecase) UpsertOfficeConfig(role string, req *domain.UpsertOfficeConfigRequest) (*domain.OfficeConfig, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}

	ws, err := normalizeTimeString(req.WorkStartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid work_start_time: %w", err)
	}
	we, err := normalizeTimeString(req.WorkEndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid work_end_time: %w", err)
	}

	workDays := req.WorkDays
	if workDays == nil {
		workDays = []int{}
	}
	wfhDays := req.WfhDays
	if wfhDays == nil {
		wfhDays = []int{}
	}
	if len(workDays) == 0 && len(wfhDays) == 0 {
		return nil, domain.ErrInvalidSchedule
	}

	wdJSON, err := json.Marshal(workDays)
	if err != nil {
		return nil, err
	}
	wfhJSON, err := json.Marshal(wfhDays)
	if err != nil {
		return nil, err
	}

	primary, err := u.repo.GetFirstOfficeConfig()
	if err != nil {
		return nil, err
	}

	ips := pq.StringArray{}
	if len(req.AllowedIPs) > 0 {
		ips = append(ips, req.AllowedIPs...)
	}

	cfg := &domain.OfficeConfig{
		Name:          strings.TrimSpace(req.Name),
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		RadiusMeters:  req.RadiusMeters,
		AllowedIPs:    ips,
		WorkStartTime: ws,
		WorkEndTime:   we,
		WorkDays:      datatypes.JSON(wdJSON),
		WfhDays:       datatypes.JSON(wfhJSON),
		IsActive:      req.IsActive,
	}

	if primary == nil {
		if err := u.repo.CreateOfficeConfig(cfg); err != nil {
			return nil, err
		}
		if cfg.IsActive {
			if err := u.repo.DeactivateOfficeConfigsExcept(cfg.ID); err != nil {
				return nil, err
			}
		}
		out, gerr := u.repo.GetOfficeConfigByID(cfg.ID)
		if gerr != nil {
			return nil, gerr
		}
		if out == nil {
			return cfg, nil
		}
		return out, nil
	}

	cfg.ID = primary.ID
	cfg.CreatedAt = primary.CreatedAt
	if err := u.repo.UpdateOfficeConfig(cfg); err != nil {
		return nil, err
	}
	if cfg.IsActive {
		if err := u.repo.DeactivateOfficeConfigsExcept(cfg.ID); err != nil {
			return nil, err
		}
	}
	out, gerr := u.repo.GetOfficeConfigByID(cfg.ID)
	if gerr != nil {
		return nil, gerr
	}
	if out == nil {
		return cfg, nil
	}
	return out, nil
}

func (u *attendanceUsecase) ListAdminRecordsByDate(role string, date time.Time) ([]domain.AttendanceRecord, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListRecordsByDate(date)
}

func (u *attendanceUsecase) DeleteAdminRecordByID(role string, id int64) error {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return domain.ErrForbiddenAdmin
	}
	if id <= 0 {
		return domain.ErrAttendanceRecordNotFound
	}
	return u.repo.DeleteRecordByID(id)
}

func (u *attendanceUsecase) CreateLeaveRequest(userID uint, req *domain.CreateLeaveRequest) (*domain.LeaveRequest, error) {
	start, err := time.Parse("2006-01-02", strings.TrimSpace(req.StartDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	end, err := time.Parse("2006-01-02", strings.TrimSpace(req.EndDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	if end.Before(start) {
		return nil, domain.ErrInvalidDateRange
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	isHalfDay := req.IsHalfDay
	halfDaySession := strings.ToUpper(strings.TrimSpace(req.HalfDaySession))
	if isHalfDay {
		if !start.Equal(end) {
			return nil, domain.ErrInvalidDateRange
		}
		if halfDaySession != "AM" && halfDaySession != "PM" {
			return nil, domain.ErrInvalidDateRange
		}
	}

	days, err := u.calculateRequestedLeaveDays(start, end, isHalfDay)
	if err != nil {
		return nil, err
	}
	if days <= 0 || days > 365 {
		return nil, domain.ErrInvalidDateRange
	}

	leaveType := strings.ToUpper(strings.TrimSpace(req.LeaveType))
	if leaveType == "" {
		leaveType = "ANNUAL"
	}
	if leaveType == "ANNUAL" {
		summary, berr := u.GetLeaveBalanceSummary(userID, start.Year())
		if berr == nil {
			for _, s := range summary {
				if s.LeaveType == "ANNUAL" && days > s.RemainingDays {
					return nil, domain.ErrInvalidDateRange
				}
			}
		}
	}

	leave := &domain.LeaveRequest{
		UserID:         userID,
		StartDate:      start,
		EndDate:        end,
		DaysRequested:  days,
		LeaveType:      leaveType,
		IsHalfDay:      isHalfDay,
		HalfDaySession: halfDaySession,
		Reason:         strings.TrimSpace(req.Reason),
		Status:         domain.LeaveStatusPending,
	}
	if err := u.repo.CreateLeaveRequest(leave); err != nil {
		return nil, err
	}
	created, gerr := u.repo.GetLeaveRequestByID(leave.ID)
	if gerr != nil {
		return nil, gerr
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leave.ID, Action: "LEAVE_CREATED", ActorID: &userID, ActorRole: "EMPLOYEE", OldStatus: "", NewStatus: domain.LeaveStatusPending, Comment: strings.TrimSpace(req.Reason)})
	_ = u.notifyLeaveRequested(leave.ID, userID, leaveType, days)
	return created, nil
}

func (u *attendanceUsecase) ListMyLeaveRequests(userID uint) ([]domain.LeaveRequest, error) {
	return u.repo.ListLeaveRequestsByUser(userID)
}

func (u *attendanceUsecase) ListPendingLeaveRequests(role string) ([]domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListPendingLeaveRequests()
}

func (u *attendanceUsecase) ListAdminLeaveRequests(role string) ([]domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListAllLeaveRequests()
}

func (u *attendanceUsecase) ReviewLeaveRequest(role string, approverID uint, leaveID int64, req *domain.ReviewLeaveRequest) (*domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	leave, err := u.repo.GetLeaveRequestByID(leaveID)
	if err != nil {
		return nil, err
	}
	if leave == nil {
		return nil, domain.ErrLeaveNotFound
	}
	if leave.Status != domain.LeaveStatusPending {
		return nil, domain.ErrLeaveNotPending
	}
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status != domain.LeaveStatusApproved && status != domain.LeaveStatusRejected {
		return nil, domain.ErrLeaveNotPending
	}
	now := time.Now().UTC()
	old := leave.Status
	leave.Status = status
	leave.ApproverID = &approverID
	leave.ManagerComment = strings.TrimSpace(req.Comment)
	leave.ApprovedAt = &now
	if err := u.repo.UpdateLeaveRequest(leave); err != nil {
		return nil, err
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leaveID, Action: "LEAVE_REVIEWED", ActorID: &approverID, ActorRole: strings.ToUpper(strings.TrimSpace(role)), OldStatus: old, NewStatus: status, Comment: leave.ManagerComment})
	_ = u.notifyLeaveReviewed(leaveID, leave.UserID, status)
	return u.repo.GetLeaveRequestByID(leaveID)
}

func (u *attendanceUsecase) UpdateAdminLeaveRequest(role string, actorID uint, leaveID int64, req *domain.UpdateLeaveRequest) (*domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	leave, err := u.repo.GetLeaveRequestByID(leaveID)
	if err != nil {
		return nil, err
	}
	if leave == nil {
		return nil, domain.ErrLeaveNotFound
	}

	start, err := time.Parse("2006-01-02", strings.TrimSpace(req.StartDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	end, err := time.Parse("2006-01-02", strings.TrimSpace(req.EndDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	if end.Before(start) {
		return nil, domain.ErrInvalidDateRange
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	isHalfDay := req.IsHalfDay
	halfDaySession := strings.ToUpper(strings.TrimSpace(req.HalfDaySession))
	if isHalfDay {
		if !start.Equal(end) {
			return nil, domain.ErrInvalidDateRange
		}
		if halfDaySession != "AM" && halfDaySession != "PM" {
			return nil, domain.ErrInvalidDateRange
		}
	}

	days, err := u.calculateRequestedLeaveDays(start, end, isHalfDay)
	if err != nil {
		return nil, err
	}
	if days <= 0 || days > 365 {
		return nil, domain.ErrInvalidDateRange
	}

	leaveType := strings.ToUpper(strings.TrimSpace(req.LeaveType))
	if leaveType == "" {
		leaveType = "ANNUAL"
	}
	old := leave.Status
	leave.StartDate = start
	leave.EndDate = end
	leave.DaysRequested = days
	leave.LeaveType = leaveType
	leave.IsHalfDay = isHalfDay
	leave.HalfDaySession = halfDaySession
	leave.Reason = strings.TrimSpace(req.Reason)
	if err := u.repo.UpdateLeaveRequest(leave); err != nil {
		return nil, err
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leaveID, Action: "LEAVE_UPDATED_BY_ADMIN", ActorID: &actorID, ActorRole: strings.ToUpper(strings.TrimSpace(role)), OldStatus: old, NewStatus: leave.Status, Comment: leave.Reason})
	return u.repo.GetLeaveRequestByID(leaveID)
}

func (u *attendanceUsecase) CancelAdminLeaveRequest(role string, actorID uint, leaveID int64, req *domain.CancelLeaveRequest) (*domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	leave, err := u.repo.GetLeaveRequestByID(leaveID)
	if err != nil {
		return nil, err
	}
	if leave == nil {
		return nil, domain.ErrLeaveNotFound
	}
	old := leave.Status
	leave.Status = domain.LeaveStatusRejected
	leave.ManagerComment = strings.TrimSpace(req.Comment)
	now := time.Now().UTC()
	leave.ApprovedAt = &now
	leave.ApproverID = &actorID
	if err := u.repo.UpdateLeaveRequest(leave); err != nil {
		return nil, err
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leaveID, Action: "LEAVE_CANCELLED_BY_ADMIN", ActorID: &actorID, ActorRole: strings.ToUpper(strings.TrimSpace(role)), OldStatus: old, NewStatus: leave.Status, Comment: leave.ManagerComment})
	_ = u.notifyLeaveReviewed(leaveID, leave.UserID, leave.Status)
	return u.repo.GetLeaveRequestByID(leaveID)
}

func (u *attendanceUsecase) DeleteAdminLeaveRequest(role string, actorID uint, leaveID int64) error {
	if role != authDomain.RoleCEO && role != authDomain.RoleSupport {
		return domain.ErrForbiddenAdmin
	}
	leave, err := u.repo.GetLeaveRequestByID(leaveID)
	if err != nil {
		return err
	}
	if leave == nil {
		return domain.ErrLeaveNotFound
	}
	if err := u.repo.DeleteLeaveRequestByID(leaveID); err != nil {
		return err
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leaveID, Action: "LEAVE_DELETED_BY_ADMIN", ActorID: &actorID, ActorRole: strings.ToUpper(strings.TrimSpace(role)), OldStatus: leave.Status, NewStatus: "", Comment: ""})
	return nil
}

func (u *attendanceUsecase) GetLeaveBalanceSummary(userID uint, year int) ([]domain.LeaveBalanceSummary, error) {
	if year <= 0 {
		year = time.Now().Year()
	}
	policies, err := u.repo.ListLeavePolicies()
	if err != nil {
		return nil, err
	}
	items, err := u.repo.ListLeaveRequestsByUser(userID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveBalanceSummary, 0, len(policies))
	for _, p := range policies {
		if !p.IsActive {
			continue
		}
		taken := 0.0
		for _, it := range items {
			if it.Status != domain.LeaveStatusApproved || it.LeaveType != p.LeaveType {
				continue
			}
			if it.StartDate.Year() == year {
				taken += it.DaysRequested
			}
		}
		carry := float64(p.MaxCarryForwardDays)
		remaining := float64(p.AnnualQuotaDays) + carry - taken
		if remaining < 0 {
			remaining = 0
		}
		out = append(out, domain.LeaveBalanceSummary{
			LeaveType:         p.LeaveType,
			AnnualQuotaDays:   p.AnnualQuotaDays,
			CarryForwardDays:  p.MaxCarryForwardDays,
			ApprovedDaysTaken: taken,
			RemainingDays:     remaining,
		})
	}
	return out, nil
}

func (u *attendanceUsecase) ListLeavePolicies(role string) ([]domain.LeavePolicy, error) {
	if role == "" {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListLeavePolicies()
}

func (u *attendanceUsecase) UpsertLeavePolicy(role string, req *domain.LeavePolicyUpsertRequest) (*domain.LeavePolicy, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	item := &domain.LeavePolicy{LeaveType: strings.ToUpper(strings.TrimSpace(req.LeaveType)), AnnualQuotaDays: req.AnnualQuotaDays, MaxCarryForwardDays: req.MaxCarryForwardDays, IsActive: req.IsActive}
	if err := u.repo.UpsertLeavePolicy(item); err != nil {
		return nil, err
	}
	all, err := u.repo.ListLeavePolicies()
	if err != nil {
		return item, nil
	}
	for _, p := range all {
		if p.LeaveType == item.LeaveType {
			cp := p
			return &cp, nil
		}
	}
	return item, nil
}

func (u *attendanceUsecase) ListHolidayCalendars(role string, fromDate, toDate time.Time) ([]domain.HolidayCalendar, error) {
	if role == "" {
		return nil, domain.ErrForbiddenAdmin
	}
	if fromDate.IsZero() || toDate.IsZero() || toDate.Before(fromDate) {
		fromDate = time.Now().UTC().AddDate(0, -1, 0)
		toDate = time.Now().UTC().AddDate(0, 6, 0)
	}
	return u.repo.ListHolidayCalendars(fromDate, toDate)
}

func (u *attendanceUsecase) UpsertHolidayCalendar(role string, req *domain.HolidayUpsertRequest) (*domain.HolidayCalendar, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	d, err := time.Parse("2006-01-02", strings.TrimSpace(req.Date))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	item := &domain.HolidayCalendar{Date: d.UTC(), Name: strings.TrimSpace(req.Name)}
	if err := u.repo.UpsertHolidayCalendar(item); err != nil {
		return nil, err
	}
	list, _ := u.repo.ListHolidayCalendars(item.Date, item.Date)
	if len(list) > 0 {
		return &list[0], nil
	}
	return item, nil
}

func (u *attendanceUsecase) ListLeaveAuditLogs(role string, leaveID int64) ([]domain.LeaveAuditLog, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListLeaveAuditLogs(leaveID)
}

func (u *attendanceUsecase) ListMyNotifications(userID uint, unreadOnly bool) ([]domain.LeaveNotification, error) {
	return u.repo.ListLeaveNotifications(userID, unreadOnly)
}

func (u *attendanceUsecase) MarkMyNotificationRead(userID uint, notificationID int64) error {
	return u.repo.MarkLeaveNotificationRead(userID, notificationID)
}

func (u *attendanceUsecase) GetLeaveTrend(role string, fromDate, toDate time.Time) ([]domain.LeaveTrendPoint, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	if fromDate.IsZero() || toDate.IsZero() || toDate.Before(fromDate) {
		fromDate = time.Now().UTC().AddDate(0, -11, 0)
		toDate = time.Now().UTC()
	}
	return u.repo.GetLeaveTrendByMonth(role, fromDate, toDate)
}

func (u *attendanceUsecase) BackfillLeave(role string, actorID uint, req *domain.LeaveBackfillRequest) (*domain.LeaveRequest, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	it := req.Item
	userID, err := u.repo.FindUserIDByEmail(strings.TrimSpace(it.EmployeeEmail))
	if err != nil {
		return nil, err
	}
	start, err := time.Parse("2006-01-02", strings.TrimSpace(it.StartDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	end, err := time.Parse("2006-01-02", strings.TrimSpace(it.EndDate))
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}
	if end.Before(start) {
		return nil, domain.ErrInvalidDateRange
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	isHalfDay := it.IsHalfDay
	halfDaySession := strings.ToUpper(strings.TrimSpace(it.HalfDaySession))
	if isHalfDay {
		if !start.Equal(end) {
			return nil, domain.ErrInvalidDateRange
		}
		if halfDaySession != "AM" && halfDaySession != "PM" {
			return nil, domain.ErrInvalidDateRange
		}
	}
	days, err := u.calculateRequestedLeaveDays(start, end, isHalfDay)
	if err != nil {
		return nil, err
	}
	if days <= 0 || days > 365 {
		return nil, domain.ErrInvalidDateRange
	}
	status := strings.ToUpper(strings.TrimSpace(it.Status))
	if status == "" {
		status = domain.LeaveStatusApproved
	}
	if status != domain.LeaveStatusPending && status != domain.LeaveStatusApproved && status != domain.LeaveStatusRejected {
		return nil, domain.ErrInvalidDateRange
	}
	leaveType := strings.ToUpper(strings.TrimSpace(it.LeaveType))
	if leaveType == "" {
		leaveType = "ANNUAL"
	}
	now := time.Now().UTC()
	leave := &domain.LeaveRequest{
		UserID:         userID,
		StartDate:      start,
		EndDate:        end,
		DaysRequested:  days,
		LeaveType:      leaveType,
		IsHalfDay:      isHalfDay,
		HalfDaySession: halfDaySession,
		Reason:         strings.TrimSpace(it.Reason),
		Status:         status,
		ManagerComment: strings.TrimSpace(it.Comment),
	}
	if status == domain.LeaveStatusApproved || status == domain.LeaveStatusRejected {
		leave.ApproverID = &actorID
		leave.ApprovedAt = &now
	}
	if err := u.repo.CreateLeaveRequest(leave); err != nil {
		return nil, err
	}
	_ = u.repo.CreateLeaveAuditLog(&domain.LeaveAuditLog{LeaveID: leave.ID, Action: "LEAVE_BACKFILL_CREATED", ActorID: &actorID, ActorRole: strings.ToUpper(strings.TrimSpace(role)), OldStatus: "", NewStatus: status, Comment: strings.TrimSpace(it.Comment), Metadata: "source=admin-backfill"})
	return u.repo.GetLeaveRequestByID(leave.ID)
}

func (u *attendanceUsecase) BackfillLeaveBulk(role string, actorID uint, req *domain.LeaveBackfillBulkRequest) (*domain.LeaveBackfillBulkResponse, error) {
	if role != authDomain.RoleCEO && role != authDomain.RoleManager && role != authDomain.RoleSupport {
		return nil, domain.ErrForbiddenAdmin
	}
	resp := &domain.LeaveBackfillBulkResponse{Total: len(req.Items), Results: make([]domain.LeaveBackfillBulkResultItem, 0, len(req.Items))}
	for i, it := range req.Items {
		out, err := u.BackfillLeave(role, actorID, &domain.LeaveBackfillRequest{Item: it})
		if err != nil {
			resp.Failed++
			resp.Results = append(resp.Results, domain.LeaveBackfillBulkResultItem{Index: i, Email: it.EmployeeEmail, Status: "failed", Error: err.Error()})
			continue
		}
		resp.Succeeded++
		resp.Results = append(resp.Results, domain.LeaveBackfillBulkResultItem{Index: i, Email: it.EmployeeEmail, Status: "created", LeaveID: out.ID})
	}
	return resp, nil
}

func (u *attendanceUsecase) notifyLeaveRequested(leaveID int64, requesterID uint, leaveType string, days float64) error {
	approvers, err := u.repo.ListAdminApproverUserIDs()
	if err != nil {
		return err
	}
	for _, uid := range approvers {
		if uid == requesterID {
			continue
		}
		title := "Leave request pending approval"
		msg := fmt.Sprintf("Employee #%d requested %s leave (%.1f day(s)).", requesterID, leaveType, days)
		_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: uid, LeaveID: leaveID, Channel: "IN_APP", Event: "LEAVE_REQUESTED", Title: title, Message: msg})
		_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: uid, LeaveID: leaveID, Channel: "EMAIL", Event: "LEAVE_REQUESTED", Title: title, Message: msg})
		_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: uid, LeaveID: leaveID, Channel: "LINE", Event: "LEAVE_REQUESTED", Title: title, Message: msg})
	}
	return nil
}

func (u *attendanceUsecase) notifyLeaveReviewed(leaveID int64, requesterID uint, status string) error {
	title := "Your leave request has been reviewed"
	msg := fmt.Sprintf("Leave request #%d is now %s.", leaveID, status)
	_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: requesterID, LeaveID: leaveID, Channel: "IN_APP", Event: "LEAVE_REVIEWED", Title: title, Message: msg})
	_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: requesterID, LeaveID: leaveID, Channel: "EMAIL", Event: "LEAVE_REVIEWED", Title: title, Message: msg})
	_ = u.repo.CreateLeaveNotification(&domain.LeaveNotification{UserID: requesterID, LeaveID: leaveID, Channel: "LINE", Event: "LEAVE_REVIEWED", Title: title, Message: msg})
	return nil
}

func (u *attendanceUsecase) getApprovedHalfDaySessionForDate(userID uint, attDate time.Time) (string, error) {
	items, err := u.repo.ListLeaveRequestsByUser(userID)
	if err != nil {
		return "", err
	}
	for _, it := range items {
		if it.Status != domain.LeaveStatusApproved || !it.IsHalfDay {
			continue
		}
		sy, sm, sd := it.StartDate.In(time.UTC).Date()
		ey, em, ed := it.EndDate.In(time.UTC).Date()
		ay, am, ad := attDate.In(time.UTC).Date()
		if sy == ay && sm == am && sd == ad && ey == ay && em == am && ed == ad {
			s := strings.ToUpper(strings.TrimSpace(it.HalfDaySession))
			if s == "AM" || s == "PM" {
				return s, nil
			}
		}
	}
	return "", nil
}

func (u *attendanceUsecase) calculateRequestedLeaveDays(startUTC, endUTC time.Time, isHalfDay bool) (float64, error) {
	if isHalfDay {
		return 0.5, nil
	}
	if endUTC.Before(startUTC) {
		return 0, domain.ErrInvalidDateRange
	}

	holidays, err := u.repo.ListHolidayCalendars(startUTC, endUTC)
	if err != nil {
		return 0, err
	}
	holidaySet := make(map[string]struct{}, len(holidays))
	for _, h := range holidays {
		hy, hm, hd := h.Date.In(time.UTC).Date()
		holidaySet[fmt.Sprintf("%04d-%02d-%02d", hy, hm, hd)] = struct{}{}
	}

	days := 0.0
	for d := startUTC; !d.After(endUTC); d = d.AddDate(0, 0, 1) {
		local := d.In(u.loc)
		wd := local.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			continue
		}
		dy, dm, dd := d.In(time.UTC).Date()
		key := fmt.Sprintf("%04d-%02d-%02d", dy, dm, dd)
		if _, isHoliday := holidaySet[key]; isHoliday {
			continue
		}
		days += 1
	}

	return days, nil
}

func calendarDateUTC(nowLocal time.Time) time.Time {
	y, m, d := nowLocal.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func isoWeekday(t time.Time) int {
	w := int(t.Weekday())
	if w == 0 {
		return 7
	}
	return w
}

func dayInList(days []int, wd int) bool {
	for _, d := range days {
		if d == wd {
			return true
		}
	}
	return false
}

// isAttendanceRequiredDay is true if today is an onsite office day or a WFH day (check-in required).
func (u *attendanceUsecase) isAttendanceRequiredDay(cfg *domain.OfficeConfig, nowLocal time.Time) bool {
	wd := isoWeekday(nowLocal)
	office, err := parseWorkDaysJSON(cfg.WorkDays)
	if err != nil {
		return false
	}
	wfh, err := parseWfhDaysJSON(cfg.WfhDays)
	if err != nil {
		return false
	}
	return dayInList(office, wd) || dayInList(wfh, wd)
}

// isWFHDay is true if today is configured as WFH (geofence skipped). If a day is both office and WFH, WFH wins.
func (u *attendanceUsecase) isWFHDay(cfg *domain.OfficeConfig, nowLocal time.Time) bool {
	wd := isoWeekday(nowLocal)
	wfh, err := parseWfhDaysJSON(cfg.WfhDays)
	if err != nil || len(wfh) == 0 {
		return false
	}
	return dayInList(wfh, wd)
}

func parseWorkDaysJSON(raw datatypes.JSON) ([]int, error) {
	b := []byte(raw)
	if len(b) == 0 {
		return []int{1, 2, 3, 4, 5}, nil
	}
	var days []int
	if err := json.Unmarshal(b, &days); err != nil {
		return nil, err
	}
	return days, nil
}

func parseWfhDaysJSON(raw datatypes.JSON) ([]int, error) {
	b := []byte(raw)
	if len(b) == 0 {
		return []int{}, nil
	}
	var days []int
	if err := json.Unmarshal(b, &days); err != nil {
		return nil, err
	}
	return days, nil
}

func (u *attendanceUsecase) workStartOnDate(cfg *domain.OfficeConfig, dayLocal time.Time) (time.Time, error) {
	h, m, s, err := parseClockParts(cfg.WorkStartTime)
	if err != nil {
		return time.Time{}, err
	}
	y, mo, d := dayLocal.Date()
	return time.Date(y, mo, d, h, m, s, 0, u.loc), nil
}

func (u *attendanceUsecase) workEndOnDate(cfg *domain.OfficeConfig, dayLocal time.Time) (time.Time, error) {
	h, m, s, err := parseClockParts(cfg.WorkEndTime)
	if err != nil {
		return time.Time{}, err
	}
	y, mo, d := dayLocal.Date()
	return time.Date(y, mo, d, h, m, s, 0, u.loc), nil
}

func parseClockParts(s string) (h, m, sec int, err error) {
	s = strings.TrimSpace(s)
	var t time.Time
	for _, layout := range []string{"15:04:05", "15:04"} {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t.Hour(), t.Minute(), t.Second(), nil
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid time %q", s)
}

func normalizeTimeString(s string) (string, error) {
	h, m, sec, err := parseClockParts(s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, sec), nil
}

func (u *attendanceUsecase) gpsWithinOffice(lat, lng float64, cfg *domain.OfficeConfig) bool {
	if cfg.RadiusMeters <= 0 {
		return false
	}
	if cfg.Latitude == 0 && cfg.Longitude == 0 {
		return false
	}
	d := domain.HaversineDistanceMeters(lat, lng, cfg.Latitude, cfg.Longitude)
	return d <= cfg.RadiusMeters
}
