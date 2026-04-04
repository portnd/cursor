package repository

import (
	"errors"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
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
