package usecase

import (
	"encoding/json"
	"fmt"
	"net"
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
		netOK := clientIPAllowed(strings.TrimSpace(clientIP), []string(cfg.AllowedIPs))
		if !gpsOK && !netOK {
			return nil, domain.ErrOutsideOffice
		}
		method = checkInMethod(gpsOK, netOK)
	}

	workStart, err := u.workStartOnDate(cfg, now)
	if err != nil {
		return nil, fmt.Errorf("work_start_time: %w", err)
	}
	isLate := now.After(workStart.Add(lateGraceMinutes * time.Minute))
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

	workEnd, err := u.workEndOnDate(cfg, now)
	if err != nil {
		return nil, fmt.Errorf("work_end_time: %w", err)
	}
	nowUTC := time.Now().UTC()
	early := now.In(u.loc).Before(workEnd)

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
	if role != authDomain.RoleCEO && role != authDomain.RoleManager {
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
	if role != authDomain.RoleCEO && role != authDomain.RoleManager {
		return nil, domain.ErrForbiddenAdmin
	}
	return u.repo.ListRecordsByDate(date)
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

func checkInMethod(gpsOK, netOK bool) string {
	switch {
	case gpsOK && netOK:
		return "both"
	case gpsOK:
		return "gps"
	case netOK:
		return "network"
	default:
		return ""
	}
}

func clientIPAllowed(clientIP string, cidrs []string) bool {
	if len(cidrs) == 0 {
		return false
	}
	host, _, err := net.SplitHostPort(clientIP)
	if err == nil {
		clientIP = host
	}
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}
	for _, c := range cidrs {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		_, network, perr := net.ParseCIDR(c)
		if perr == nil {
			if network.Contains(ip) {
				return true
			}
			continue
		}
		single := net.ParseIP(c)
		if single != nil && single.Equal(ip) {
			return true
		}
	}
	return false
}
