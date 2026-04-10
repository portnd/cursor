package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// DailyStandup represents an async daily check-in submitted by a team member.
type DailyStandup struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uint           `json:"user_id" gorm:"not null;uniqueIndex:uq_daily_standups_user_date"`
	Date             time.Time      `json:"date" gorm:"type:date;not null;uniqueIndex:uq_daily_standups_user_date"` // UTC date of the standup
	YesterdaySummary string         `json:"yesterday_summary" gorm:"type:text;not null"`
	TodayTaskIDs     pq.StringArray `json:"today_task_ids" gorm:"type:text[];default:'{}'"`
	Blocker          string         `json:"blocker" gorm:"type:text"` // empty string means no blocker
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`

	// Enriched at read-time (not stored)
	UserEmail       string `json:"user_email,omitempty" gorm:"-"`
	UserDisplayName string `json:"user_display_name,omitempty" gorm:"-"`
}

func (DailyStandup) TableName() string { return "daily_standups" }

// UserPulse aggregates all activity for a single user on a given day.
type UserPulse struct {
	UserID           uint           `json:"user_id"`
	UserEmail        string         `json:"user_email"`
	UserDisplayName  string         `json:"user_display_name"`
	UserAvatarURL    string         `json:"user_avatar_url,omitempty"`
	Standup          *DailyStandup  `json:"standup"`
	IsOnLeave        bool           `json:"is_on_leave"`
	LeaveType        string         `json:"leave_type,omitempty"`
	LeaveSession     string         `json:"leave_session,omitempty"` // AM|PM|FULL
	TotalLoggedMin   int            `json:"total_logged_minutes"`
	TotalLoggedHrs   float64        `json:"total_logged_hours"`
	LatestActivities []ActivityItem `json:"latest_activities"`
	HasBlocker       bool           `json:"has_blocker"`
}

// ActivityItem is a lightweight summary of a time-log or submission event.
type ActivityItem struct {
	Type         string    `json:"type"` // "time_log" | "submission"
	Description  string    `json:"description"`
	Minutes      int       `json:"minutes,omitempty"`      // only for time_log
	ReferenceURL string    `json:"reference_url,omitempty"` // only for submission
	OccurredAt   time.Time `json:"occurred_at"`
}

// CompanyPulseResponse is the full daily pulse board payload.
type CompanyPulseResponse struct {
	Date           string      `json:"date"` // YYYY-MM-DD
	TotalMembers   int         `json:"total_members"`
	CheckedIn      int         `json:"checked_in"`
	OnLeaveCount   int         `json:"on_leave_count"`
	TotalMinLogged int         `json:"total_minutes_logged"`
	Members        []UserPulse `json:"members"`
}

// ─── Ports (interfaces) ────────────────────────────────────────────────────────

// PulseRepository defines data-access operations.
type PulseRepository interface {
	SaveStandup(s *DailyStandup) error
	GetStandupByUserAndDate(userID uint, date time.Time) (*DailyStandup, error)
	GetStandupsByDate(date time.Time) ([]DailyStandup, error)

	// Cross-module reads (read-only queries into sentinel tables)
	GetTimeLogsByDate(date time.Time) ([]TimeLogSummary, error)
	GetSubmissionsByDate(date time.Time) ([]SubmissionSummary, error)
	GetApprovedLeavesByDate(date time.Time) ([]LeaveSummary, error)

	// CEO-configured hidden users for pulse board visibility
	GetHiddenPulseUserIDs() ([]uint, error)
	SetHiddenPulseUserIDs(userIDs []uint) error
}

// PulseUsecase defines business operations.
type PulseUsecase interface {
	SubmitStandup(userID uint, date time.Time, yesterday, blocker string, todayTaskIDs []string) (*DailyStandup, error)
	GetDailyCompanyPulse(date time.Time, viewerRole string) (*CompanyPulseResponse, error)
	GetHiddenPulseUserIDs(requesterRole string) ([]uint, error)
	SetHiddenPulseUserIDs(requesterRole string, userIDs []uint) error
}

// ─── Lightweight cross-module read models ─────────────────────────────────────

// TimeLogSummary is a minimal projection of the time_logs table.
type TimeLogSummary struct {
	UserID      uint      `json:"user_id"`
	Minutes     int       `json:"minutes"`
	Description string    `json:"description"`
	LoggedAt    time.Time `json:"logged_at"`
}

// SubmissionSummary is a minimal projection of the submissions table.
type SubmissionSummary struct {
	DevID        uint      `json:"dev_id"`
	ReferenceURL string    `json:"reference_url"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"created_at"`
}

// LeaveSummary is a minimal projection of approved leave covering a date.
type LeaveSummary struct {
	UserID       uint   `json:"user_id"`
	LeaveType    string `json:"leave_type"`
	IsHalfDay    bool   `json:"is_half_day"`
	HalfSession  string `json:"half_day_session"` // AM|PM|""
}
