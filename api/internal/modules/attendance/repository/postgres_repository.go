package repository

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository returns an AttendanceRepository backed by PostgreSQL.
func NewPostgresRepository(db *gorm.DB) domain.AttendanceRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetActiveOfficeConfig() (*domain.OfficeConfig, error) {
	var cfg domain.OfficeConfig
	err := r.db.Where("is_active = ?", true).Order("id ASC").First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *postgresRepository) GetFirstOfficeConfig() (*domain.OfficeConfig, error) {
	var cfg domain.OfficeConfig
	err := r.db.Order("id ASC").First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *postgresRepository) GetOfficeConfigByID(id uint) (*domain.OfficeConfig, error) {
	var cfg domain.OfficeConfig
	err := r.db.First(&cfg, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *postgresRepository) CreateOfficeConfig(cfg *domain.OfficeConfig) error {
	return r.db.Create(cfg).Error
}

func (r *postgresRepository) UpdateOfficeConfig(cfg *domain.OfficeConfig) error {
	return r.db.Save(cfg).Error
}

func (r *postgresRepository) DeactivateAllOfficeConfigs() error {
	return r.db.Model(&domain.OfficeConfig{}).Where("1 = 1").Update("is_active", false).Error
}

func (r *postgresRepository) DeactivateOfficeConfigsExcept(id uint) error {
	return r.db.Model(&domain.OfficeConfig{}).Where("id <> ?", id).Update("is_active", false).Error
}

func (r *postgresRepository) GetRecordByUserAndDate(userID uint, attendanceDate time.Time) (*domain.AttendanceRecord, error) {
	d := attendanceDate.UTC()
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)

	var rec domain.AttendanceRecord
	err := r.db.Where("user_id = ? AND attendance_date = ?", userID, day.Format("2006-01-02")).First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *postgresRepository) SaveRecord(rec *domain.AttendanceRecord) error {
	dateStr := rec.AttendanceDate.UTC()
	day := time.Date(dateStr.Year(), dateStr.Month(), dateStr.Day(), 0, 0, 0, 0, time.UTC)
	key := day.Format("2006-01-02")

	var existing domain.AttendanceRecord
	err := r.db.Where("user_id = ? AND attendance_date = ?", rec.UserID, key).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rec.AttendanceDate = day
		return r.db.Create(rec).Error
	}
	if err != nil {
		return err
	}

	rec.ID = existing.ID
	rec.CreatedAt = existing.CreatedAt
	rec.AttendanceDate = day
	return r.db.Model(&domain.AttendanceRecord{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
		"office_config_id": rec.OfficeConfigID,
		"attendance_date":  key,
		"check_in_at":      rec.CheckInAt,
		"check_out_at":     rec.CheckOutAt,
		"check_in_lat":     rec.CheckInLat,
		"check_in_lng":     rec.CheckInLng,
		"check_in_method":  rec.CheckInMethod,
		"check_in_ip":      rec.CheckInIP,
		"is_late":          rec.IsLate,
		"early_checkout":   rec.EarlyCheckout,
		"status":           rec.Status,
	}).Error
}

func (r *postgresRepository) ListUserRecordsAfterID(userID uint, afterID int64, limit int) ([]domain.AttendanceRecord, error) {
	q := r.db.Where("user_id = ?", userID).Order("id ASC").Limit(limit)
	if afterID > 0 {
		q = q.Where("id > ?", afterID)
	}
	var rows []domain.AttendanceRecord
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *postgresRepository) ListRecordsByDate(attendanceDate time.Time) ([]domain.AttendanceRecord, error) {
	d := attendanceDate.UTC()
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	dateStr := day.Format("2006-01-02")

	type row struct {
		ID              int64      `gorm:"column:id"`
		UserID          uint       `gorm:"column:user_id"`
		OfficeConfigID  uint       `gorm:"column:office_config_id"`
		AttendanceDate  time.Time  `gorm:"column:attendance_date"`
		CheckInAt       *time.Time `gorm:"column:check_in_at"`
		CheckOutAt      *time.Time `gorm:"column:check_out_at"`
		CheckInLat      *float64   `gorm:"column:check_in_lat"`
		CheckInLng      *float64   `gorm:"column:check_in_lng"`
		CheckInMethod   string     `gorm:"column:check_in_method"`
		CheckInIP       string     `gorm:"column:check_in_ip"`
		IsLate          bool       `gorm:"column:is_late"`
		EarlyCheckout   bool       `gorm:"column:early_checkout"`
		Status          string     `gorm:"column:status"`
		CreatedAt       time.Time  `gorm:"column:created_at"`
		UpdatedAt       time.Time  `gorm:"column:updated_at"`
		UserEmail       string     `gorm:"column:user_email"`
		UserDisplayName string     `gorm:"column:user_display_name"`
	}

	var raw []row
	err := r.db.Raw(`
		SELECT ar.id, ar.user_id, ar.office_config_id, ar.attendance_date, ar.check_in_at, ar.check_out_at,
		       ar.check_in_lat, ar.check_in_lng, ar.check_in_method, ar.check_in_ip, ar.is_late, ar.early_checkout,
		       ar.status, ar.created_at, ar.updated_at,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name
		FROM attendance_records ar
		JOIN users u ON u.id = ar.user_id
		WHERE ar.attendance_date = ?
		ORDER BY ar.check_in_at NULLS LAST, ar.id ASC
	`, dateStr).Scan(&raw).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.AttendanceRecord, 0, len(raw))
	for _, rw := range raw {
		out = append(out, domain.AttendanceRecord{
			ID:              rw.ID,
			UserID:          rw.UserID,
			OfficeConfigID:  rw.OfficeConfigID,
			AttendanceDate:  rw.AttendanceDate,
			CheckInAt:       rw.CheckInAt,
			CheckOutAt:      rw.CheckOutAt,
			CheckInLat:      rw.CheckInLat,
			CheckInLng:      rw.CheckInLng,
			CheckInMethod:   rw.CheckInMethod,
			CheckInIP:       rw.CheckInIP,
			IsLate:          rw.IsLate,
			EarlyCheckout:   rw.EarlyCheckout,
			Status:          rw.Status,
			CreatedAt:       rw.CreatedAt,
			UpdatedAt:       rw.UpdatedAt,
			UserEmail:       rw.UserEmail,
			UserDisplayName: rw.UserDisplayName,
		})
	}
	return out, nil
}

func (r *postgresRepository) DeleteRecordByID(id int64) error {
	res := r.db.Delete(&domain.AttendanceRecord{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrAttendanceRecordNotFound
	}
	return nil
}

func (r *postgresRepository) CreateOffsiteCheckInRequest(req *domain.OffsiteCheckInRequest) error {
	return r.db.Create(req).Error
}

func (r *postgresRepository) GetLatestOffsiteCheckInRequestByUserAndDate(userID uint, attendanceDate time.Time) (*domain.OffsiteCheckInRequest, error) {
	d := attendanceDate.UTC()
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	var item domain.OffsiteCheckInRequest
	err := r.db.Where("user_id = ? AND attendance_date = ?", userID, day.Format("2006-01-02")).
		Order("id DESC").
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepository) ListPendingOffsiteCheckInRequests() ([]domain.OffsiteCheckInRequest, error) {
	type row struct {
		domain.OffsiteCheckInRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT o.*,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM offsite_checkin_requests o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN users au ON au.id = o.approver_id
		WHERE o.status = ?
		ORDER BY o.requested_at ASC, o.id ASC
	`, domain.OffsiteStatusPending).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.OffsiteCheckInRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.OffsiteCheckInRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) GetOffsiteCheckInRequestByID(id int64) (*domain.OffsiteCheckInRequest, error) {
	type row struct {
		domain.OffsiteCheckInRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var rw row
	err := r.db.Raw(`
		SELECT o.*,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM offsite_checkin_requests o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN users au ON au.id = o.approver_id
		WHERE o.id = ?
	`, id).Scan(&rw).Error
	if err != nil {
		return nil, err
	}
	if rw.ID == 0 {
		return nil, nil
	}
	item := rw.OffsiteCheckInRequest
	item.UserEmail = rw.UserEmail
	item.UserDisplayName = rw.UserDisplayName
	item.ApproverEmail = rw.ApproverEmail
	item.ApproverName = rw.ApproverName
	return &item, nil
}

func (r *postgresRepository) UpdateOffsiteCheckInRequest(req *domain.OffsiteCheckInRequest) error {
	return r.db.Save(req).Error
}

func (r *postgresRepository) CreateOffsiteCheckOutRequest(req *domain.OffsiteCheckOutRequest) error {
	return r.db.Create(req).Error
}

func (r *postgresRepository) GetLatestOffsiteCheckOutRequestByUserAndDate(userID uint, attendanceDate time.Time) (*domain.OffsiteCheckOutRequest, error) {
	d := attendanceDate.UTC()
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	var item domain.OffsiteCheckOutRequest
	err := r.db.Where("user_id = ? AND attendance_date = ?", userID, day.Format("2006-01-02")).
		Order("id DESC").
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepository) ListPendingOffsiteCheckOutRequests() ([]domain.OffsiteCheckOutRequest, error) {
	type row struct {
		domain.OffsiteCheckOutRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT o.*,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM offsite_checkout_requests o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN users au ON au.id = o.approver_id
		WHERE o.status = ?
		ORDER BY o.requested_at ASC, o.id ASC
	`, domain.OffsiteStatusPending).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.OffsiteCheckOutRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.OffsiteCheckOutRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) GetOffsiteCheckOutRequestByID(id int64) (*domain.OffsiteCheckOutRequest, error) {
	type row struct {
		domain.OffsiteCheckOutRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var rw row
	err := r.db.Raw(`
		SELECT o.*,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM offsite_checkout_requests o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN users au ON au.id = o.approver_id
		WHERE o.id = ?
	`, id).Scan(&rw).Error
	if err != nil {
		return nil, err
	}
	if rw.ID == 0 {
		return nil, nil
	}
	item := rw.OffsiteCheckOutRequest
	item.UserEmail = rw.UserEmail
	item.UserDisplayName = rw.UserDisplayName
	item.ApproverEmail = rw.ApproverEmail
	item.ApproverName = rw.ApproverName
	return &item, nil
}

func (r *postgresRepository) UpdateOffsiteCheckOutRequest(req *domain.OffsiteCheckOutRequest) error {
	return r.db.Save(req).Error
}

func (r *postgresRepository) CreateLeaveRequest(req *domain.LeaveRequest) error {
	return r.db.Create(req).Error
}

func (r *postgresRepository) GetLeaveRequestByID(id int64) (*domain.LeaveRequest, error) {
	var req domain.LeaveRequest
	err := r.db.First(&req, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *postgresRepository) ListLeaveRequestsByUser(userID uint) ([]domain.LeaveRequest, error) {
	type row struct {
		domain.LeaveRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT lr.*, 
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM leave_requests lr
		JOIN users u ON u.id = lr.user_id
		LEFT JOIN users au ON au.id = lr.approver_id
		WHERE lr.user_id = ?
		ORDER BY lr.created_at DESC, lr.id DESC
	`, userID).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.LeaveRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) ListPendingLeaveRequests() ([]domain.LeaveRequest, error) {
	type row struct {
		domain.LeaveRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT lr.*, 
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM leave_requests lr
		JOIN users u ON u.id = lr.user_id
		LEFT JOIN users au ON au.id = lr.approver_id
		WHERE lr.status = ?
		ORDER BY lr.created_at ASC, lr.id ASC
	`, domain.LeaveStatusPending).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.LeaveRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) ListAllLeaveRequests() ([]domain.LeaveRequest, error) {
	type row struct {
		domain.LeaveRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT lr.*, 
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM leave_requests lr
		JOIN users u ON u.id = lr.user_id
		LEFT JOIN users au ON au.id = lr.approver_id
		ORDER BY lr.created_at DESC, lr.id DESC
	`).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.LeaveRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) ListApprovedLeaveRequestsByDate(date time.Time) ([]domain.LeaveRequest, error) {
	type row struct {
		domain.LeaveRequest
		UserEmail       string `gorm:"column:user_email"`
		UserDisplayName string `gorm:"column:user_display_name"`
		ApproverEmail   string `gorm:"column:approver_email"`
		ApproverName    string `gorm:"column:approver_name"`
	}
	var raw []row
	dateStr := date.Format("2006-01-02")
	err := r.db.Raw(`
		SELECT lr.*,
		       u.email AS user_email,
		       COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), NULLIF(TRIM(u.display_name), ''), u.email) AS user_display_name,
		       au.email AS approver_email,
		       COALESCE(NULLIF(TRIM(au.display_name), ''), au.email) AS approver_name
		FROM leave_requests lr
		JOIN users u ON u.id = lr.user_id
		LEFT JOIN users au ON au.id = lr.approver_id
		WHERE lr.status = ?
		  AND lr.start_date <= ?
		  AND lr.end_date >= ?
		ORDER BY lr.start_date ASC, lr.id ASC
	`, domain.LeaveStatusApproved, dateStr, dateStr).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveRequest, 0, len(raw))
	for _, rw := range raw {
		item := rw.LeaveRequest
		item.UserEmail = rw.UserEmail
		item.UserDisplayName = rw.UserDisplayName
		item.ApproverEmail = rw.ApproverEmail
		item.ApproverName = rw.ApproverName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) UpdateLeaveRequest(req *domain.LeaveRequest) error {
	return r.db.Save(req).Error
}

func (r *postgresRepository) DeleteLeaveRequestByID(id int64) error {
	res := r.db.Delete(&domain.LeaveRequest{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrLeaveNotFound
	}
	return nil
}

func (r *postgresRepository) ListLeavePolicies() ([]domain.LeavePolicy, error) {
	var out []domain.LeavePolicy
	err := r.db.Order("leave_type ASC").Find(&out).Error
	return out, err
}

func (r *postgresRepository) UpsertLeavePolicy(req *domain.LeavePolicy) error {
	var existing domain.LeavePolicy
	err := r.db.Where("leave_type = ?", req.LeaveType).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(req).Error
	}
	if err != nil {
		return err
	}
	req.ID = existing.ID
	req.CreatedAt = existing.CreatedAt
	return r.db.Save(req).Error
}

func (r *postgresRepository) ListHolidayCalendars(fromDate, toDate time.Time) ([]domain.HolidayCalendar, error) {
	var out []domain.HolidayCalendar
	err := r.db.Where("date >= ? AND date <= ?", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Order("date ASC").Find(&out).Error
	return out, err
}

func (r *postgresRepository) UpsertHolidayCalendar(item *domain.HolidayCalendar) error {
	day := time.Date(item.Date.Year(), item.Date.Month(), item.Date.Day(), 0, 0, 0, 0, time.UTC)
	var existing domain.HolidayCalendar
	err := r.db.Where("date = ?", day.Format("2006-01-02")).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item.Date = day
		return r.db.Create(item).Error
	}
	if err != nil {
		return err
	}
	item.ID = existing.ID
	item.Date = day
	item.CreatedAt = existing.CreatedAt
	return r.db.Save(item).Error
}

func (r *postgresRepository) CreateLeaveAuditLog(logItem *domain.LeaveAuditLog) error {
	return r.db.Create(logItem).Error
}

func (r *postgresRepository) ListLeaveAuditLogs(leaveID int64) ([]domain.LeaveAuditLog, error) {
	type row struct {
		domain.LeaveAuditLog
		ActorEmail string `gorm:"column:actor_email"`
		ActorName  string `gorm:"column:actor_name"`
	}
	var raw []row
	err := r.db.Raw(`
		SELECT l.*, u.email AS actor_email,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS actor_name
		FROM leave_audit_logs l
		LEFT JOIN users u ON u.id = l.actor_id
		WHERE l.leave_id = ?
		ORDER BY l.created_at DESC, l.id DESC
	`, leaveID).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.LeaveAuditLog, 0, len(raw))
	for _, rw := range raw {
		item := rw.LeaveAuditLog
		item.ActorEmail = rw.ActorEmail
		item.ActorName = rw.ActorName
		out = append(out, item)
	}
	return out, nil
}

func (r *postgresRepository) CreateLeaveNotification(item *domain.LeaveNotification) error {
	now := time.Now().UTC()
	item.DeliveredAt = &now
	return r.db.Create(item).Error
}

func (r *postgresRepository) ListLeaveNotifications(userID uint, unreadOnly bool) ([]domain.LeaveNotification, error) {
	q := r.db.Where("user_id = ?", userID).Order("created_at DESC, id DESC")
	if unreadOnly {
		q = q.Where("is_read = ?", false)
	}
	var out []domain.LeaveNotification
	err := q.Find(&out).Error
	return out, err
}

func (r *postgresRepository) MarkLeaveNotificationRead(userID uint, notificationID int64) error {
	return r.db.Model(&domain.LeaveNotification{}).Where("id = ? AND user_id = ?", notificationID, userID).Update("is_read", true).Error
}

func (r *postgresRepository) GetLeaveTrendByMonth(role string, fromDate, toDate time.Time) ([]domain.LeaveTrendPoint, error) {
	var out []domain.LeaveTrendPoint
	isAdmin := strings.EqualFold(role, authDomain.RoleCEO) || strings.EqualFold(role, authDomain.RoleManager) || strings.EqualFold(role, authDomain.RoleSupport)
	if !isAdmin {
		return out, nil
	}

	teamFeatureEnabled := true
	var settingValue string
	if err := r.db.Raw(`SELECT value FROM app_settings WHERE key = 'teams_feature_enabled'`).Scan(&settingValue).Error; err == nil {
		if parsed, perr := strconv.ParseBool(strings.TrimSpace(settingValue)); perr == nil {
			teamFeatureEnabled = parsed
		}
	}

	if teamFeatureEnabled {
		err := r.db.Raw(`
			SELECT to_char(date_trunc('month', lr.start_date), 'YYYY-MM') AS month,
			       u.team_id AS team_id,
			       COALESCE(t.name, 'Unassigned') AS team_name,
			       COUNT(*) AS requested,
			       SUM(CASE WHEN lr.status = 'APPROVED' THEN 1 ELSE 0 END) AS approved,
			       SUM(CASE WHEN lr.status = 'REJECTED' THEN 1 ELSE 0 END) AS rejected,
			       COALESCE(SUM(lr.days_requested), 0) AS total_days
			FROM leave_requests lr
			JOIN users u ON u.id = lr.user_id
			LEFT JOIN teams t ON t.id = u.team_id
			WHERE lr.start_date >= ? AND lr.start_date <= ?
			GROUP BY 1, 2, 3
			ORDER BY 1 ASC, 3 ASC
		`, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Scan(&out).Error
		return out, err
	}

	err := r.db.Raw(`
		SELECT to_char(date_trunc('month', lr.start_date), 'YYYY-MM') AS month,
		       NULL::bigint AS team_id,
		       '' AS team_name,
		       u.id AS user_id,
		       COALESCE(NULLIF(TRIM(u.display_name), ''), u.email) AS user_name,
		       u.email AS user_email,
		       COUNT(*) AS requested,
		       SUM(CASE WHEN lr.status = 'APPROVED' THEN 1 ELSE 0 END) AS approved,
		       SUM(CASE WHEN lr.status = 'REJECTED' THEN 1 ELSE 0 END) AS rejected,
		       COALESCE(SUM(lr.days_requested), 0) AS total_days
		FROM leave_requests lr
		JOIN users u ON u.id = lr.user_id
		WHERE lr.start_date >= ? AND lr.start_date <= ?
		GROUP BY 1, 4, 5, 6
		ORDER BY 1 ASC, 5 ASC
	`, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Scan(&out).Error
	return out, err
}

func (r *postgresRepository) ListAdminApproverUserIDs() ([]uint, error) {
	var ids []uint
	err := r.db.Model(&authDomain.User{}).Where("upper(role) IN ?", []string{authDomain.RoleCEO, authDomain.RoleManager}).Pluck("id", &ids).Error
	return ids, err
}

func (r *postgresRepository) FindUserIDByEmail(email string) (uint, error) {
	var u authDomain.User
	err := r.db.Where("lower(email) = lower(?)", strings.TrimSpace(email)).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.ErrUserNotFound
		}
		return 0, err
	}
	return u.ID, nil
}

func (r *postgresRepository) IsUserRemote(userID uint) (bool, error) {
	var u authDomain.User
	err := r.db.Select("is_remote").First(&u, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return u.IsRemote, nil
}
