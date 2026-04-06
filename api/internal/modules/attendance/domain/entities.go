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

// LeaveRequest is employee leave workflow record.
type LeaveRequest struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint       `json:"user_id" gorm:"not null;index"`
	StartDate        time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate          time.Time  `json:"end_date" gorm:"type:date;not null"`
	DaysRequested    int        `json:"days_requested" gorm:"not null;default:1"`
	LeaveType        string     `json:"leave_type" gorm:"type:varchar(20);not null;default:ANNUAL"`
	Reason           string     `json:"reason" gorm:"type:text;not null;default:''"`
	Status           string     `json:"status" gorm:"type:varchar(20);not null;default:PENDING"`
	ApproverID       *uint      `json:"approver_id,omitempty" gorm:"index"`
	ManagerComment   string     `json:"manager_comment" gorm:"type:text;default:''"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	UserEmail        string     `json:"user_email,omitempty" gorm:"-"`
	UserDisplayName  string     `json:"user_display_name,omitempty" gorm:"-"`
	ApproverEmail    string     `json:"approver_email,omitempty" gorm:"-"`
	ApproverName     string     `json:"approver_name,omitempty" gorm:"-"`
}

func (LeaveRequest) TableName() string { return "leave_requests" }

const (
	LeaveStatusPending  = "PENDING"
	LeaveStatusApproved = "APPROVED"
	LeaveStatusRejected = "REJECTED"
)

// LeavePolicy defines annual quota and carry-forward policy by leave type.
type LeavePolicy struct {
	ID                  int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	LeaveType           string    `json:"leave_type" gorm:"type:varchar(20);not null;uniqueIndex"`
	AnnualQuotaDays     int       `json:"annual_quota_days" gorm:"not null;default:10"`
	MaxCarryForwardDays int       `json:"max_carry_forward_days" gorm:"not null;default:0"`
	IsActive            bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (LeavePolicy) TableName() string { return "leave_policies" }

// HolidayCalendar is company holiday configuration.
type HolidayCalendar struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Date      time.Time `json:"date" gorm:"type:date;not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (HolidayCalendar) TableName() string { return "holiday_calendars" }

// LeaveAuditLog tracks all state transitions and comments for compliance.
type LeaveAuditLog struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	LeaveID     int64     `json:"leave_id" gorm:"not null;index"`
	Action      string    `json:"action" gorm:"type:varchar(50);not null"`
	ActorID     *uint     `json:"actor_id,omitempty" gorm:"index"`
	ActorRole   string    `json:"actor_role" gorm:"type:varchar(20);not null;default:''"`
	OldStatus   string    `json:"old_status" gorm:"type:varchar(20);not null;default:''"`
	NewStatus   string    `json:"new_status" gorm:"type:varchar(20);not null;default:''"`
	Comment     string    `json:"comment" gorm:"type:text;not null;default:''"`
	Metadata    string    `json:"metadata" gorm:"type:text;not null;default:''"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	ActorEmail  string    `json:"actor_email,omitempty" gorm:"-"`
	ActorName   string    `json:"actor_name,omitempty" gorm:"-"`
}

func (LeaveAuditLog) TableName() string { return "leave_audit_logs" }

// LeaveNotification records in-app notifications and delivery attempts.
type LeaveNotification struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	LeaveID     int64      `json:"leave_id" gorm:"not null;index"`
	Channel     string     `json:"channel" gorm:"type:varchar(20);not null;default:IN_APP"`
	Event       string     `json:"event" gorm:"type:varchar(50);not null"`
	Title       string     `json:"title" gorm:"type:varchar(200);not null"`
	Message     string     `json:"message" gorm:"type:text;not null"`
	IsRead      bool       `json:"is_read" gorm:"not null;default:false"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (LeaveNotification) TableName() string { return "leave_notifications" }

type LeaveBalanceSummary struct {
	LeaveType         string `json:"leave_type"`
	AnnualQuotaDays   int    `json:"annual_quota_days"`
	CarryForwardDays  int    `json:"carry_forward_days"`
	ApprovedDaysTaken int    `json:"approved_days_taken"`
	RemainingDays     int    `json:"remaining_days"`
}

type LeaveTrendPoint struct {
	Month      string `json:"month"`
	TeamID     *uint  `json:"team_id,omitempty"`
	TeamName   string `json:"team_name,omitempty"`
	Requested  int    `json:"requested"`
	Approved   int    `json:"approved"`
	Rejected   int    `json:"rejected"`
	TotalDays  int    `json:"total_days"`
}

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

	CreateLeaveRequest(req *LeaveRequest) error
	GetLeaveRequestByID(id int64) (*LeaveRequest, error)
	ListLeaveRequestsByUser(userID uint) ([]LeaveRequest, error)
	ListPendingLeaveRequests() ([]LeaveRequest, error)
	UpdateLeaveRequest(req *LeaveRequest) error

	ListLeavePolicies() ([]LeavePolicy, error)
	UpsertLeavePolicy(req *LeavePolicy) error
	ListHolidayCalendars(fromDate, toDate time.Time) ([]HolidayCalendar, error)
	UpsertHolidayCalendar(item *HolidayCalendar) error
	CreateLeaveAuditLog(log *LeaveAuditLog) error
	ListLeaveAuditLogs(leaveID int64) ([]LeaveAuditLog, error)
	CreateLeaveNotification(item *LeaveNotification) error
	ListLeaveNotifications(userID uint, unreadOnly bool) ([]LeaveNotification, error)
	MarkLeaveNotificationRead(userID uint, notificationID int64) error
	GetLeaveTrendByMonth(role string, fromDate, toDate time.Time) ([]LeaveTrendPoint, error)
	ListAdminApproverUserIDs() ([]uint, error)
	FindUserIDByEmail(email string) (uint, error)
}

// CreateLeaveRequest is employee payload to request leave.
type CreateLeaveRequest struct {
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate   string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	LeaveType string `json:"leave_type" binding:"required,oneof=ANNUAL SICK PERSONAL UNPAID"`
	Reason    string `json:"reason" binding:"required,min=3,max=1000"`
}

// ReviewLeaveRequest is manager payload to approve/reject leave request.
type ReviewLeaveRequest struct {
	Status  string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	Comment string `json:"comment" binding:"max=1000"`
}

// LeaveListResponse is list wrapper for leave requests.
type LeaveListResponse struct {
	Items []LeaveRequest `json:"items"`
}

type LeavePolicyUpsertRequest struct {
	LeaveType           string `json:"leave_type" binding:"required,oneof=ANNUAL SICK PERSONAL UNPAID"`
	AnnualQuotaDays     int    `json:"annual_quota_days" binding:"required,gte=0,lte=365"`
	MaxCarryForwardDays int    `json:"max_carry_forward_days" binding:"required,gte=0,lte=365"`
	IsActive            bool   `json:"is_active"`
}

type HolidayUpsertRequest struct {
	Date string `json:"date" binding:"required"` // YYYY-MM-DD
	Name string `json:"name" binding:"required,min=2,max=200"`
}

type LeaveTrendResponse struct {
	Items []LeaveTrendPoint `json:"items"`
}

type LeaveBackfillItem struct {
	EmployeeEmail string `json:"employee_email"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	LeaveType     string `json:"leave_type"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Comment       string `json:"comment"`
}

type LeaveBackfillRequest struct {
	Item LeaveBackfillItem `json:"item" binding:"required"`
}

type LeaveBackfillBulkRequest struct {
	Items []LeaveBackfillItem `json:"items" binding:"required,dive"`
}

type LeaveBackfillBulkResultItem struct {
	Index   int    `json:"index"`
	Email   string `json:"email"`
	Status  string `json:"status"`
	LeaveID int64  `json:"leave_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LeaveBackfillBulkResponse struct {
	Total     int                           `json:"total"`
	Succeeded int                           `json:"succeeded"`
	Failed    int                           `json:"failed"`
	Results   []LeaveBackfillBulkResultItem `json:"results"`
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

	CreateLeaveRequest(userID uint, req *CreateLeaveRequest) (*LeaveRequest, error)
	ListMyLeaveRequests(userID uint) ([]LeaveRequest, error)
	ListPendingLeaveRequests(role string) ([]LeaveRequest, error)
	ReviewLeaveRequest(role string, approverID uint, leaveID int64, req *ReviewLeaveRequest) (*LeaveRequest, error)

	GetLeaveBalanceSummary(userID uint, year int) ([]LeaveBalanceSummary, error)
	ListLeavePolicies(role string) ([]LeavePolicy, error)
	UpsertLeavePolicy(role string, req *LeavePolicyUpsertRequest) (*LeavePolicy, error)
	ListHolidayCalendars(role string, fromDate, toDate time.Time) ([]HolidayCalendar, error)
	UpsertHolidayCalendar(role string, req *HolidayUpsertRequest) (*HolidayCalendar, error)
	ListLeaveAuditLogs(role string, leaveID int64) ([]LeaveAuditLog, error)
	ListMyNotifications(userID uint, unreadOnly bool) ([]LeaveNotification, error)
	MarkMyNotificationRead(userID uint, notificationID int64) error
	GetLeaveTrend(role string, fromDate, toDate time.Time) ([]LeaveTrendPoint, error)
	BackfillLeave(role string, actorID uint, req *LeaveBackfillRequest) (*LeaveRequest, error)
	BackfillLeaveBulk(role string, actorID uint, req *LeaveBackfillBulkRequest) (*LeaveBackfillBulkResponse, error)
}
