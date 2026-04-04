package domain

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// OfficeConfig is the configurable office location and work schedule.
type OfficeConfig struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"type:varchar(200);not null;default:Main Office"`
	Latitude      float64        `json:"latitude" gorm:"not null;default:0"`
	Longitude     float64        `json:"longitude" gorm:"not null;default:0"`
	RadiusMeters  float64        `json:"radius_meters" gorm:"column:radius_meters;not null;default:100"`
	AllowedIPs    pq.StringArray `json:"allowed_ips" gorm:"type:text[]"`
	WorkStartTime string         `json:"work_start_time" gorm:"type:varchar(8);not null;default:09:00:00"`
	WorkEndTime   string         `json:"work_end_time" gorm:"type:varchar(8);not null;default:18:00:00"`
	WorkDays      datatypes.JSON `json:"work_days" gorm:"type:jsonb"`
	WfhDays       datatypes.JSON `json:"wfh_days" gorm:"type:jsonb"`
	IsActive      bool           `json:"is_active" gorm:"not null;default:true"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (OfficeConfig) TableName() string { return "office_configs" }

// AttendanceRecord is one row per user per calendar day.
type AttendanceRecord struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	OfficeConfigID uint       `json:"office_config_id" gorm:"not null"`
	AttendanceDate time.Time  `json:"attendance_date" gorm:"type:date;column:attendance_date;not null"`
	CheckInAt      *time.Time `json:"check_in_at,omitempty"`
	CheckOutAt     *time.Time `json:"check_out_at,omitempty"`
	CheckInLat     *float64   `json:"check_in_lat,omitempty"`
	CheckInLng     *float64   `json:"check_in_lng,omitempty"`
	CheckInMethod  string     `json:"check_in_method" gorm:"type:varchar(20);not null;default:''"`
	CheckInIP      string     `json:"check_in_ip" gorm:"type:varchar(64);not null;default:''"`
	IsLate         bool       `json:"is_late" gorm:"not null;default:false"`
	EarlyCheckout  bool       `json:"early_checkout" gorm:"not null;default:false"`
	Status         string     `json:"status" gorm:"type:varchar(20);not null;default:absent"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Enriched (not stored)
	UserEmail       string `json:"user_email,omitempty" gorm:"-"`
	UserDisplayName string `json:"user_display_name,omitempty" gorm:"-"`
}

func (AttendanceRecord) TableName() string { return "attendance_records" }

// UpsertOfficeConfigRequest is the admin payload for office settings.
type UpsertOfficeConfigRequest struct {
	Name          string   `json:"name" binding:"required,min=1,max=200"`
	Latitude      float64  `json:"latitude" binding:"required"`
	Longitude     float64  `json:"longitude" binding:"required"`
	RadiusMeters  float64  `json:"radius_meters" binding:"required,gt=0,lte=5000"`
	AllowedIPs    []string `json:"allowed_ips"`
	WorkStartTime string   `json:"work_start_time" binding:"required"` // HH:MM or HH:MM:SS
	WorkEndTime   string   `json:"work_end_time" binding:"required"`
	WorkDays      []int    `json:"work_days" binding:"max=7,dive,gte=1,lte=7"`
	WfhDays       []int    `json:"wfh_days" binding:"max=7,dive,gte=1,lte=7"`
	IsActive      bool     `json:"is_active"`
}

// AttendanceHistoryResponse is cursor-paginated history for the current user.
type AttendanceHistoryResponse struct {
	Items      []AttendanceRecord `json:"items"`
	NextCursor string             `json:"next_cursor,omitempty"`
}

// AttendanceRepository defines persistence for attendance.
type AttendanceRepository interface {
	GetActiveOfficeConfig() (*OfficeConfig, error)
	GetFirstOfficeConfig() (*OfficeConfig, error)
	GetOfficeConfigByID(id uint) (*OfficeConfig, error)
	CreateOfficeConfig(cfg *OfficeConfig) error
	UpdateOfficeConfig(cfg *OfficeConfig) error
	DeactivateAllOfficeConfigs() error
	DeactivateOfficeConfigsExcept(id uint) error

	GetRecordByUserAndDate(userID uint, attendanceDate time.Time) (*AttendanceRecord, error)
	SaveRecord(rec *AttendanceRecord) error
	ListUserRecordsAfterID(userID uint, afterID int64, limit int) ([]AttendanceRecord, error)
	ListRecordsByDate(attendanceDate time.Time) ([]AttendanceRecord, error)
}

// AttendanceUsecase defines attendance operations.
type AttendanceUsecase interface {
	CheckIn(userID uint, lat, lng float64, clientIP string) (*AttendanceRecord, error)
	CheckOut(userID uint) (*AttendanceRecord, error)
	GetTodayStatus(userID uint) (*AttendanceRecord, *OfficeConfig, error)
	GetHistory(userID uint, cursor string, limit int) (*AttendanceHistoryResponse, error)

	GetOfficeConfigForAdmin() (*OfficeConfig, error)
	UpsertOfficeConfig(role string, req *UpsertOfficeConfigRequest) (*OfficeConfig, error)
	ListAdminRecordsByDate(role string, date time.Time) ([]AttendanceRecord, error)
}
