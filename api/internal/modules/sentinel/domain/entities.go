package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ErrBadRequest wraps domain validation errors that should map to HTTP 400.
type ErrBadRequest struct{ Msg string }

func (e *ErrBadRequest) Error() string { return e.Msg }

// IsBadRequest returns true when err (or its cause) is an ErrBadRequest.
func IsBadRequest(err error) bool {
	var e *ErrBadRequest
	return errors.As(err, &e)
}

// Epic represents a large feature or goal within a project (Hierarchy Dimension 1)
type Epic struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	Status      string     `json:"status" gorm:"default:'PLANNING'"` // PLANNING, IN_PROGRESS, DONE
	Color       string     `json:"color" gorm:"default:'#6366f1'"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Tasks       []Task     `json:"tasks,omitempty" gorm:"foreignKey:EpicID"`
}

func (Epic) TableName() string { return "epics" }

// EpicTimelineData is the payload for the Epic Roadmap view
type EpicTimelineData struct {
	Epics []Epic `json:"epics"` // Each Epic has Tasks preloaded
}

// SprintTimelineData is the payload for the Sprint Execution view
type SprintTimelineData struct {
	Sprints []Sprint `json:"sprints"` // Each Sprint has Tasks preloaded
}

// Project groups tasks (multi-project architecture)
type Project struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string    `json:"code" gorm:"type:varchar(64);uniqueIndex"` // slug from name e.g. mims-hdmap-main (empty for legacy)
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"default:'ACTIVE'"` // ACTIVE, COMPLETED, ON_HOLD
	Color       string    `json:"color" gorm:"type:varchar(16);default:'#6366f1'"`
	// Squad Model: which team owns this project
	TeamID   *uint  `json:"team_id" gorm:"index"`
	TeamName string `json:"team_name" gorm:"-"` // populated via JOIN, not stored
	// Internal VC — Project Capital
	CapitalBalance  float64   `json:"capital_balance" gorm:"column:capital_balance;type:decimal(15,2);default:0"`
	BonusPercentage float64   `json:"bonus_percentage" gorm:"type:decimal(5,2);default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Tasks           []Task    `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	// Task counts for list view (populated by repo, not stored)
	TaskTotal     int `json:"task_total" gorm:"-"`
	TaskCompleted int `json:"task_completed" gorm:"-"`
	TaskOverdue   int `json:"task_overdue" gorm:"-"`
	// When teams feature is off, CEO assigns Product Owners here (multiple allowed); populated by repo only (JSON field pm_owners kept for compatibility)
	PmOwners []ProjectPmOwner `json:"pm_owners,omitempty" gorm:"-"`
}

// ProjectPmAssignment links a project to Product Owner users (used when squads/teams feature is disabled; table name kept for compatibility).
type ProjectPmAssignment struct {
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;primaryKey"`
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
}

func (ProjectPmAssignment) TableName() string { return "project_pm_assignments" }

// ProjectPmOwner is returned on project payloads for display and CEO editing.
type ProjectPmOwner struct {
	UserID      uint   `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
}

// ProjectTransactionType defines the type of a project capital transaction
type ProjectTransactionType string

const (
	ProjTxInjection   ProjectTransactionType = "INJECTION"
	ProjTxBurn        ProjectTransactionType = "BURN"
	ProjTxBonusPayout ProjectTransactionType = "BONUS_PAYOUT"
	ProjTxAdjustment  ProjectTransactionType = "ADJUSTMENT"
)

// ProjectTransaction records every capital movement for a project
type ProjectTransaction struct {
	ID        int64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID uuid.UUID              `json:"project_id" gorm:"type:uuid;not null;index"`
	Type      ProjectTransactionType `json:"type" gorm:"type:varchar(20);not null"`
	Amount    float64                `json:"amount" gorm:"type:decimal(15,2);not null"`
	Reference string                 `json:"reference" gorm:"type:text;default:''"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
}

func (ProjectTransaction) TableName() string { return "project_transactions" }

// ProjectCapitalResponse is the response DTO for GetProjectCapital
type ProjectCapitalResponse struct {
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectName string    `json:"project_name"`
	TeamID      *uint     `json:"team_id"`
	// Burn rate = team's full monthly cost (entire team, not split by project)
	TeamMonthlyCost float64              `json:"team_monthly_cost"`
	CapitalBalance  float64              `json:"capital_balance"`
	BonusPercentage float64              `json:"bonus_percentage"`
	RunwayMonths    float64              `json:"runway_months"`
	Transactions    []ProjectTransaction `json:"transactions,omitempty"`
}

// InjectProjectCapitalRequest is the DTO for injecting capital into a project
type InjectProjectCapitalRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	BonusPercentage float64 `json:"bonus_percentage" binding:"gte=0,lte=100"`
	Note            string  `json:"note"`
}

// EditProjectCapitalRequest is the DTO for directly setting the project capital balance
type EditProjectCapitalRequest struct {
	NewBalance      float64  `json:"new_balance" binding:"gte=0"`
	BonusPercentage *float64 `json:"bonus_percentage" binding:"omitempty,gte=0,lte=100"`
	Note            string   `json:"note"`
}

// CloseProjectCycleResponse is the response DTO for CloseProjectCycleAndPayout
type CloseProjectCycleResponse struct {
	ProjectID       uuid.UUID `json:"project_id"`
	BalanceBefore   float64   `json:"balance_before"`
	BonusPercentage float64   `json:"bonus_percentage"`
	BonusAmount     float64   `json:"bonus_amount"`
	BalanceAfter    float64   `json:"balance_after"`
}

// ProjectFinanceUsecase defines business logic for per-project Internal VC model
type ProjectFinanceUsecase interface {
	GetProjectCapital(projectID uuid.UUID) (*ProjectCapitalResponse, error)
	GetProjectCapitals(projectIDs []uuid.UUID) ([]ProjectCapitalResponse, error)
	InjectProjectCapital(projectID uuid.UUID, req *InjectProjectCapitalRequest) (*Project, error)
	EditProjectCapital(projectID uuid.UUID, req *EditProjectCapitalRequest) (*Project, error)
	CloseProjectCycleAndPayout(projectID uuid.UUID) (*CloseProjectCycleResponse, error)
	DeleteProjectTransaction(txID int64, projectID uuid.UUID) error
}

// GlobalActiveTask is the payload for GET /tasks/my-global-active.
// It embeds a Task enriched with the project identity fields so a developer
// can identify which project each task belongs to without an extra round-trip.
type GlobalActiveTask struct {
	Task
	ProjectName  string `json:"project_name" gorm:"column:project_name"`
	ProjectColor string `json:"project_color" gorm:"column:project_color"`
}

// Sprint represents a time-boxed iteration within a project
type Sprint struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Name      string     `json:"name" gorm:"not null"`
	Goal      string     `json:"goal" gorm:"type:text"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Status    string     `json:"status" gorm:"default:'PLANNING'"` // PLANNING, ACTIVE, COMPLETED
	SortOrder int        `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Tasks     []Task     `json:"tasks,omitempty" gorm:"foreignKey:SprintID"`
	// GetActiveSprintsForUser: joined project identity for dev dashboard links
	ProjectName string `json:"project_name,omitempty" gorm:"-"`
	ProjectCode string `json:"project_code,omitempty" gorm:"-"`
}

func (Sprint) TableName() string { return "sprints" }

// Milestone represents a key deliverable deadline within a project
type Milestone struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status" gorm:"default:'PENDING'"` // PENDING, REACHED, MISSED
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Milestone) TableName() string { return "milestones" }

// TaskComment represents a discussion message on a task
type TaskComment struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID          uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	UserEmail       string    `json:"user_email,omitempty" gorm:"-"`
	UserDisplayName string    `json:"user_display_name,omitempty" gorm:"-"`
	UserAvatarURL   string    `json:"user_avatar_url,omitempty" gorm:"-"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	// Attachments stores serialized []TaskCommentAttachment in JSONB.
	Attachments datatypes.JSON `json:"attachments,omitempty" gorm:"type:jsonb;default:'[]'"`
	// EditHistory stores serialized []TaskCommentEditHistoryItem in JSONB.
	EditHistory datatypes.JSON `json:"edit_history,omitempty" gorm:"type:jsonb;default:'[]'"`
	EditedAt    *time.Time     `json:"edited_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (TaskComment) TableName() string { return "task_comments" }

type TaskCommentAttachment struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	DataURL  string `json:"data_url"`
	IsImage  bool   `json:"is_image"`
}

type TaskCommentEditHistoryItem struct {
	EditedAt   time.Time `json:"edited_at"`
	EditedBy   uint      `json:"edited_by"`
	OldContent string    `json:"old_content"`
	NewContent string    `json:"new_content"`
}

// TimeLog represents actual time logged by a developer on a task
type TimeLog struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID          uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	UserEmail       string    `json:"user_email,omitempty"  gorm:"column:user_email;<-:false"`
	TaskCode        string    `json:"task_code,omitempty"   gorm:"column:task_code;<-:false"`
	TaskTitle       string    `json:"task_title,omitempty"  gorm:"column:task_title;<-:false"`
	Minutes         int       `json:"minutes" gorm:"not null"`
	Description     string    `json:"description" gorm:"type:text"`
	WorkType        string    `json:"work_type" gorm:"default:'DEV'"`
	LoggedDate      time.Time `json:"logged_date" gorm:"type:date;default:CURRENT_DATE"`
	IsTimerSession  bool      `json:"is_timer_session" gorm:"default:false"`
	LoggedAt        time.Time `json:"logged_at" gorm:"autoCreateTime"`
}

func (TimeLog) TableName() string { return "time_logs" }

// Valid work types for TimeLog
var ValidWorkTypes = map[string]bool{
	"DEV": true, "REVIEW": true, "TESTING": true,
	"MEETING": true, "RESEARCH": true, "OTHER": true,
}

// DailyTimeLogSummary is the per-user daily log summary for /users/me/time-logs
type DailyTimeLogSummary struct {
	Date         string    `json:"date"`
	TotalMinutes int       `json:"total_minutes"`
	Entries      []TimeLog `json:"entries"`
}

// BulkLogEntry is a single item in a bulk log request
type BulkLogEntry struct {
	TaskID      string  `json:"task_id" binding:"required"`
	Minutes     int     `json:"minutes" binding:"required,min=1"`
	Description string  `json:"description"`
	WorkType    string  `json:"work_type"`
	LoggedDate  *string `json:"logged_date"` // YYYY-MM-DD, optional
}

// BulkLogResult is the outcome for each entry in a bulk log request
type BulkLogResult struct {
	TaskID  string    `json:"task_id"`
	Success bool      `json:"success"`
	Log     *TimeLog  `json:"log,omitempty"`
	Error   string    `json:"error,omitempty"`
}

// ProjectAnalytics holds computed analytics for a project
type ProjectAnalytics struct {
	ProjectID          uuid.UUID         `json:"project_id"`
	TotalTasks         int               `json:"total_tasks"`
	CompletedTasks     int               `json:"completed_tasks"`
	TotalStoryPoints   int               `json:"total_story_points"`
	CompletedSP        int               `json:"completed_story_points"`
	TotalLoggedMinutes int               `json:"total_logged_minutes"`
	AvgCycleTimeDays   float64           `json:"avg_cycle_time_days"`
	Burndown           []BurndownPoint   `json:"burndown"`
	Velocity           []VelocityPoint   `json:"velocity"`
	TeamCapacity       []TeamCapacityRow `json:"team_capacity"`
}

type BurndownPoint struct {
	Day       string  `json:"day"`
	Ideal     float64 `json:"ideal"`
	Remaining float64 `json:"remaining"`
}

type VelocityPoint struct {
	SprintName  string `json:"sprint_name"`
	CompletedSP int    `json:"completed_sp"`
	PlannedSP   int    `json:"planned_sp"`
}

type TeamCapacityRow struct {
	UserID          uint   `json:"user_id"`
	UserEmail       string `json:"user_email"`
	UserDisplayName string `json:"user_display_name,omitempty"`
	AssignedTasks    int     `json:"assigned_tasks"`
	EstimatedHours   float64 `json:"estimated_hours"`
	LoggedHours      float64 `json:"logged_hours"`
	Utilization      float64 `json:"utilization_pct"`
}

func (Project) TableName() string { return "projects" }

// ProjectDetailsTasksMeta describes task paging metadata for project details response.
type ProjectDetailsTasksMeta struct {
	Limit    int  `json:"limit"`
	Returned int  `json:"returned"`
	HasMore  bool `json:"has_more"`
}

// ProjectTaskPageCursor is cursor metadata for project task paging.
type ProjectTaskPageCursor struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
}

// ProjectTasksPageResponse is paginated task response for project details page 2+ loading.
type ProjectTasksPageResponse struct {
	Tasks       []TaskSummary          `json:"tasks"`
	Limit       int                   `json:"limit"`
	Returned    int                   `json:"returned"`
	HasMore     bool                  `json:"has_more"`
	NextCursor  *ProjectTaskPageCursor `json:"next_cursor,omitempty"`
	NextOffset  *int                  `json:"next_offset,omitempty"`
}

// ProjectDetailsResponse is the combined payload for GET /projects/:id/details (project page - 1 round-trip).
// Note: Tasks can be large, so details responses expose a minimal task shape and separate lightweight sprint/epic shapes.
type ProjectDetailsResponse struct {
	Project    *Project                `json:"project"`
	Tasks      []ProjectDetailsTask    `json:"tasks"`
	TasksMeta  ProjectDetailsTasksMeta `json:"tasks_meta"`
	Sprints    []ProjectDetailsSprint  `json:"sprints"`
	Milestones []Milestone             `json:"milestones"`
	Epics      []ProjectDetailsEpic    `json:"epics"`
}

// ProjectDetailsEpic is a lightweight epic projection for the project details payload.
type ProjectDetailsEpic struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Color       string     `json:"color"`
	SortOrder   int        `json:"sort_order"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ProjectDetailsSprint is a lightweight sprint projection for the project details payload.
type ProjectDetailsSprint struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Name      string     `json:"name"`
	Goal      string     `json:"goal"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Status    string     `json:"status"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ProjectDetailsTask is a lightweight task projection used by the project details payload.
type ProjectDetailsTask struct {
	ID                    uuid.UUID  `json:"id"`
	Code                  string     `json:"code"`
	Title                 string     `json:"title"`
	ProjectID             *uuid.UUID `json:"project_id"`
	EpicID                *uuid.UUID `json:"epic_id,omitempty"`
	SprintID              *uuid.UUID `json:"sprint_id,omitempty"`
	MilestoneID           *uuid.UUID `json:"milestone_id,omitempty"`
	TaskType              string     `json:"task_type"`
	Priority              string     `json:"priority"`
	StoryPoints           float64    `json:"story_points"`
	ParentID              *uuid.UUID `json:"parent_id"`
	SortOrder             int        `json:"sort_order"`
	StartDate             *time.Time `json:"start_date"`
	EndDate               *time.Time `json:"end_date"`
	Progress              int        `json:"progress"`
	DueAt                 *time.Time `json:"due_at"`
	Status                string     `json:"status"`
	AssignedTo            *uint      `json:"assigned_to"`
	AssignedToDisplayName string     `json:"assigned_to_display_name,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

// TaskSummary is the lightweight payload used for project boards, task lists and the task detail
// header/sidebar (pre-hydration shell). Heavy fields such as description, resource_urls, sub_tasks
// and submissions are intentionally omitted — callers that need them must hit the /detail endpoint.
type TaskSummary struct {
	ID                    uuid.UUID  `json:"id"`
	Code                  string     `json:"code"`
	Title                 string     `json:"title"`
	ProjectID             *uuid.UUID `json:"project_id"`
	EpicID                *uuid.UUID `json:"epic_id,omitempty"`
	SprintID              *uuid.UUID `json:"sprint_id,omitempty"`
	MilestoneID           *uuid.UUID `json:"milestone_id,omitempty"`
	TaskType              string     `json:"task_type"`
	Priority              string     `json:"priority"`
	StoryPoints           float64    `json:"story_points"`
	EstimatedMinutes      int        `json:"estimated_minutes"`
	ParentID              *uuid.UUID `json:"parent_id"`
	SortOrder             int        `json:"sort_order"`
	StartDate             *time.Time `json:"start_date"`
	EndDate               *time.Time `json:"end_date"`
	Progress              int        `json:"progress"`
	DueAt                 *time.Time `json:"due_at"`
	StartedAt             *time.Time `json:"started_at,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	Status                string     `json:"status"`
	NegotiationStatus     string     `json:"negotiation_status"`
	AssignedTo            *uint      `json:"assigned_to"`
	AssignedToDisplayName string     `json:"assigned_to_display_name,omitempty"`
	AssignedToEmail       string     `json:"assigned_to_email,omitempty"`
	AssignedToAvatarURL   string     `json:"assigned_to_avatar_url,omitempty"`
	IsKomgrip             bool       `json:"is_komgrip"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// Task represents a work assignment
type Task struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code             string         `json:"code" gorm:"type:varchar(64);uniqueIndex"` // e.g. mims-hdmap-main-001 (empty for legacy tasks)
	Title            string         `json:"title" gorm:"not null"`
	Description      string         `json:"description"`
	ResourceURLs     datatypes.JSON `json:"resource_urls" gorm:"type:jsonb;default:'{}'"`
	EstimatedMinutes int            `json:"estimated_minutes"`

	// Project: every task belongs to a project
	ProjectID *uuid.UUID `json:"project_id" gorm:"type:uuid;index"`

	// Epic linking (Hierarchy Dimension 1)
	EpicID *uuid.UUID `json:"epic_id,omitempty" gorm:"type:uuid;index"`
	Epic   *Epic      `json:"epic,omitempty" gorm:"foreignKey:EpicID"`

	// Sprint & Milestone linking
	SprintID    *uuid.UUID `json:"sprint_id,omitempty" gorm:"type:uuid;index"`
	MilestoneID *uuid.UUID `json:"milestone_id,omitempty" gorm:"type:uuid;index"`

	// Task Typology: distinguishes client-facing Features from dev Tasks and Bugs
	TaskType string `json:"task_type" gorm:"type:varchar(20);default:'TASK'"` // FEATURE, TASK, BUG

	// Priority & Estimation
	Priority    string  `json:"priority" gorm:"default:'MEDIUM'"` // CRITICAL, HIGH, MEDIUM, LOW
	StoryPoints float64 `json:"story_points" gorm:"type:decimal(5,1);default:0"`

	// WBS / Gantt: sub-tasks and planned dates
	ParentID   *uuid.UUID `json:"parent_id" gorm:"type:uuid;index"`                               // For sub-tasks (Work Breakdown Structure)
	ParentTask *Task      `json:"parent_task,omitempty" gorm:"foreignKey:ParentID;references:ID"` // Loaded on demand
	SortOrder  int        `json:"sort_order" gorm:"default:0"`                                    // Order within epic (backlog drag-and-drop)
	StartDate  *time.Time `json:"start_date"`                                                     // Planned start date
	EndDate    *time.Time `json:"end_date"`                                                       // Planned end date
	Progress   int        `json:"progress" gorm:"default:0"`                                      // 0 to 100%
	SubTasks   []Task     `json:"sub_tasks,omitempty" gorm:"foreignKey:ParentID"`                 // Has many sub-tasks

	// Time Negotiation System
	NegotiationStatus string `json:"negotiation_status" gorm:"default:'NONE'"` // NONE, PENDING, APPROVED, REJECTED
	ProposedMinutes   int    `json:"proposed_minutes"`                         // Dev's proposed time
	NegotiationReason string `json:"negotiation_reason" gorm:"type:text"`      // Why dev needs more time

	// AI Advisory for Time Negotiation
	NegotiationAIRecommendation string `json:"negotiation_ai_recommendation"` // APPROVE, REJECT
	NegotiationAIConfidence     int    `json:"negotiation_ai_confidence"`     // 0-100
	NegotiationAIReasoning      string `json:"negotiation_ai_reasoning" gorm:"type:text"`

	// Timestamps (Pointers allow NULL in DB)
	DueAt       *time.Time `json:"due_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`

	// Status values: PENDING | IN_PROGRESS | READY_FOR_TEST | REVIEW_PENDING | READY_FOR_UAT | COMPLETED | CANCELLED
	Status string `json:"status" gorm:"default:'PENDING';index"`

	// UAT payload (only set on FEATURE tasks after dev submits for UAT)
	UATPayload datatypes.JSON `json:"uat_payload,omitempty" gorm:"type:jsonb"`

	// Relationships (uint to match User ID)
	AssignedTo   *uint `json:"assigned_to"`                 // Developer assigned to this task
	AssignedByID *uint `json:"assigned_by_id" gorm:"index"` // Product Owner/CEO who assigned the task (for assigner-scoped leaderboard)
	CreatedBy    *uint `json:"created_by"`

	// Enriched from auth (not stored in DB)
	CreatedByRole        string `json:"created_by_role,omitempty" gorm:"-"`
	CreatedByEmail       string `json:"created_by_email,omitempty" gorm:"-"`
	CreatedByDisplayName string `json:"created_by_display_name,omitempty" gorm:"-"`

	AssignedToDisplayName string `json:"assigned_to_display_name,omitempty" gorm:"-"`
	AssignedToEmail       string `json:"assigned_to_email,omitempty" gorm:"-"`
	AssignedToAvatarURL   string `json:"assigned_to_avatar_url,omitempty" gorm:"-"`

	// My Board (GET /tasks/my): populated in repository, not stored
	ProjectName       string     `json:"project_name,omitempty" gorm:"-"`
	ProjectColor      string     `json:"project_color,omitempty" gorm:"-"`
	SprintName        string     `json:"sprint_name,omitempty" gorm:"-"`
	EffectiveSprintID *uuid.UUID `json:"effective_sprint_id,omitempty" gorm:"-"`

	// Komgrip: task not tied to any project (personal/misc tasks)
	IsKomgrip bool `json:"is_komgrip" gorm:"default:false;index"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationship: Submissions history (loaded via Preload)
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:TaskID"`
}

// UATPayloadData holds the UAT submission details stored in the uat_payload JSONB column.
type UATPayloadData struct {
	StagingURL      string `json:"staging_url"`
	TestCredentials string `json:"test_credentials"`
	ReleaseNotes    string `json:"release_notes"`
}

// TaskDependency links tasks (e.g. Task B cannot start until Task A finishes)
type TaskDependency struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PredecessorID uuid.UUID `json:"predecessor_id" gorm:"type:uuid;not null;index"` // Task that must happen first
	SuccessorID   uuid.UUID `json:"successor_id" gorm:"type:uuid;not null;index"`   // Task that waits
	Type          string    `json:"type" gorm:"type:varchar(2);default:'FS'"`       // FS (Finish-to-Start), SS (Start-to-Start), etc.
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName overrides default for task_dependencies
func (TaskDependency) TableName() string { return "task_dependencies" }

// TableName overrides default
func (Task) TableName() string { return "tasks" }

// Submission represents a handover record: engineer submits a PR/Commit URL for Product Owner review
type Submission struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID       uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	DevID        uint      `json:"dev_id" gorm:"not null"`
	ReferenceURL string    `json:"reference_url" gorm:"type:varchar(512)"`
	Note         string    `json:"note" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Submission) TableName() string { return "submissions" }

// Appeal represents a developer's appeal against AI verdict
type Appeal struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubmissionID uuid.UUID `json:"submission_id" gorm:"type:uuid;not null"`
	DeveloperID  uint      `json:"developer_id" gorm:"not null"`
	Reason       string    `json:"reason" gorm:"type:text;not null"`
	Status       string    `json:"status" gorm:"default:'PENDING'"` // PENDING, APPROVED, REJECTED

	// AI Advisory System (Generated on submission)
	AIRecommendation string `json:"ai_recommendation" gorm:"type:text"` // "UPHOLD" (reject appeal) or "OVERTURN" (approve appeal)
	AIConfidence     int    `json:"ai_confidence"`                      // 0-100 confidence score
	AIReasoning      string `json:"ai_reasoning" gorm:"type:text"`      // Explanation for CEO / Product Owner

	ResolverID   *uint  `json:"resolver_id"`   // Product Owner/CEO who resolved
	ResolverNote string `json:"resolver_note"` // Admin's decision note

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Appeal) TableName() string { return "appeals" }

// SystemConfig represents dynamic AI configuration (singleton pattern)
type SystemConfig struct {
	ID               uint      `json:"-" gorm:"primaryKey"`
	ActiveModel      string    `json:"active_model" gorm:"default:'gemini-2.5-flash-lite'"`
	Temperature      float32   `json:"temperature" gorm:"default:0.4"`      // 0.0 (Stable) to 1.0 (Creative)
	CursorAssistance int       `json:"cursor_assistance" gorm:"default:80"` // 0 to 100% (AI assistance level)
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (SystemConfig) TableName() string { return "system_configs" }

// CallerContext carries the authenticated user's role and team for data-isolation enforcement.
type CallerContext struct {
	Role                 string
	TeamID               *uint
	UserID               uint
	TeamsFeatureDisabled bool // true when app setting teams_feature_enabled is "false"
}

// Role constants for CallerContext (mirrors auth/domain constants)
const (
	RoleCEO            = "CEO"
	RoleManager        = "MANAGER"
	RoleProductOwner   = "PRODUCT_OWNER"
	RoleEngineer       = "ENGINEER"
	RoleChiefEngineer  = "CHIEF_ENGINEER"
)

// IsEngineerRole is true for roles that share engineer project/task visibility rules.
func IsEngineerRole(role string) bool {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case RoleEngineer, RoleChiefEngineer:
		return true
	default:
		return false
	}
}

// TaskType defines the typology of a task item
type TaskType string

const (
	TaskTypeFeature TaskType = "FEATURE" // Client-facing feature; acts as a parent container (no assignee/estimate)
	TaskTypeTask    TaskType = "TASK"    // Standard developer task; assignee and estimate are mandatory
	TaskTypeBug     TaskType = "BUG"     // Defect/bug report; assignee and estimate are mandatory
)

// SentinelRepository defines the interface (Port)
type SentinelRepository interface {
	// Projects
	CreateProject(p *Project) error
	GetAllProjects(ctx CallerContext) ([]Project, error)
	GetProjectByID(id uuid.UUID, ctx CallerContext) (*Project, error)
	GetProjectByCode(code string, ctx CallerContext) (*Project, error)
	GetTasksByProjectID(projectID uuid.UUID) ([]Task, error)
	// GetTasksByProjectIDForProjectPage returns tasks without large columns (description, resource_urls, negotiation text, uat_payload)
	// for GET /projects/:id/details and supports limit+1 pagination probe (returned tasks are at most limit+1 before caller trims).
	GetTasksByProjectIDForProjectPage(projectID uuid.UUID, limit int) ([]Task, error)
	GetTasksByProjectIDForProjectPageCursor(projectID uuid.UUID, limit int, cursorCreatedAt *time.Time, cursorID *uuid.UUID, offset int) ([]Task, error)
	UpdateProject(p *Project) error
	DeleteProject(id uuid.UUID) error
	DeleteProjectPlan(projectID uuid.UUID) error               // Remove all tasks, sprints, milestones, epics for the project
	AssignProjectTeam(projectID uuid.UUID, teamID *uint) error // Squad Model: assign/unassign team
	ReplaceProjectPmAssignments(projectID uuid.UUID, userIDs []uint) error

	CreateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	GetTaskByCode(code string) (*Task, error)
	CountTasksForCode(projectID *uuid.UUID) (int, error)
	GetMaxTaskCodeSuffix(prefix string) (int, error) // max numeric suffix for code like "prefix-001"; used for globally unique codes
	GetTasksByAssignee(userID uint) ([]Task, error)
	GetActiveSprintTasksByAssignee(userID uint) ([]Task, error)         // Only tasks in ACTIVE sprint (for engineer role)
	GetActiveSprintsForUser(userID uint) ([]Sprint, error)              // ACTIVE sprints that have tasks assigned to user
	GetGlobalActiveTasks(ctx CallerContext) ([]GlobalActiveTask, error) // TASK/BUG in ACTIVE sprints; CEO/MANAGER: all projects; else team / Product Owner–engineer rules when teams off
	GetTeamActiveTasks(ctx CallerContext) ([]GlobalActiveTask, error)   // ACTIVE-sprint tasks visible to caller (team-scoped or project-visibility when teams disabled)
	GetActiveFeatures(teamID uint, projectID *uuid.UUID) ([]FeatureRoadmapItem, error)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]Task, error)
	GetTasksRequiringApproval() ([]Task, error) // Tasks with PENDING appeals or time negotiations
	UpdateTask(task *Task) error
	DeleteTask(id uuid.UUID) error // Delete a task by ID

	// Task activity audit (immutable timeline)
	CreateTaskActivity(e *TaskActivityEvent) error
	ListTaskActivitiesByTaskID(taskID uuid.UUID) ([]TaskActivityEvent, error)

	GetAllTaskDependencies() ([]TaskDependency, error) // For Gantt chart
	CreateTaskDependency(dep *TaskDependency) error
	DeleteTaskDependency(id uuid.UUID) error

	CreateSubmission(sub *Submission) error
	GetSubmissionByID(id uuid.UUID) (*Submission, error)
	UpdateSubmission(sub *Submission) error
	GetLatestSubmission(taskID uuid.UUID) (*Submission, error)

	CreateAppeal(appeal *Appeal) error
	GetAppealBySubmissionID(subID uuid.UUID) (*Appeal, error)
	GetAppealByID(id uuid.UUID) (*Appeal, error)
	UpdateAppeal(appeal *Appeal) error

	// Human Quality Gate
	ApproveTask(id uuid.UUID) error                                    // Mark task as COMPLETED and set CompletedAt
	RejectTask(taskID uuid.UUID, rejectorID uint, reason string) error // Return task to IN_PROGRESS with comment

	// Continuous UAT: sub-task level testing queue
	GetTasksReadyForTest(teamID uint) ([]GlobalActiveTask, error)                // All TASK/BUG in READY_FOR_TEST status, scoped to team
	SetTaskReadyForUAT(taskID uuid.UUID, uatPayload []byte) error                // (legacy) Product Owner approves: READY_FOR_TEST → READY_FOR_UAT with test evidence
	SetTaskWaitForDeploy(taskID uuid.UUID, uatPayload []byte) error              // Product Owner approves: READY_FOR_TEST → WAIT_FOR_DEPLOY (pending Chief Engineer deployment)
	AdvanceTaskToReadyForUAT(taskID uuid.UUID) error                             // Chief Engineer deployed: WAIT_FOR_DEPLOY → READY_FOR_UAT
	GetTasksReadyForCEOApproval(teamID uint) ([]GlobalActiveTask, error)         // All TASK/BUG in READY_FOR_UAT status for CEO final approval

	// Sprints
	CreateSprint(sprint *Sprint) error
	GetSprintByID(id uuid.UUID) (*Sprint, error)
	GetSprintsByProjectID(projectID uuid.UUID) ([]Sprint, error)
	GetActiveSprintByProjectID(projectID uuid.UUID) (*Sprint, error)
	UpdateSprint(sprint *Sprint) error
	DeleteSprint(id uuid.UUID) error

	// Milestones
	CreateMilestone(m *Milestone) error
	GetMilestoneByID(id uuid.UUID) (*Milestone, error)
	GetMilestonesByProjectID(projectID uuid.UUID) ([]Milestone, error)
	UpdateMilestone(m *Milestone) error
	DeleteMilestone(id uuid.UUID) error

	// Task Comments
	CreateTaskComment(c *TaskComment) error
	GetCommentsByTaskID(taskID uuid.UUID) ([]TaskComment, error)
	GetTaskCommentByID(commentID uuid.UUID) (*TaskComment, error)
	UpdateTaskComment(c *TaskComment) error

	// Time Logs
	CreateTimeLog(t *TimeLog) error
	GetTimeLogsByTaskID(taskID uuid.UUID) ([]TimeLog, error)
	GetTimeLogByID(logID uuid.UUID) (*TimeLog, error)
	UpdateTimeLog(t *TimeLog) error
	DeleteTimeLog(logID uuid.UUID) error
	GetTimeLogsByUserAndDate(userID uint, date time.Time) ([]TimeLog, error)
	BulkCreateTimeLogs(logs []TimeLog) error
	GetTotalLoggedMinutes(taskID uuid.UUID) (int, error)
	CountChildTasks(parentID uuid.UUID) (int, error)             // Leaf-node guard: returns number of direct children
	GetChildTasksByParentID(parentID uuid.UUID) ([]Task, error)  // For UAT roll-up: fetch all direct children of a feature

	// Analytics
	GetProjectAnalytics(projectID uuid.UUID) (*ProjectAnalytics, error)

	// Bulk operations
	BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error

	// System Configuration (Singleton)
	GetSystemConfig() (*SystemConfig, error)
	UpdateSystemConfig(config *SystemConfig) error

	// Epics
	CreateEpic(epic *Epic) error
	GetEpicByID(id uuid.UUID) (*Epic, error)
	GetEpicsByProjectID(projectID uuid.UUID) ([]Epic, error)
	UpdateEpic(epic *Epic) error
	DeleteEpic(id uuid.UUID) error

	// Timeline views (Matrix Dimension)
	GetEpicTimelineData(projectID uuid.UUID) (*EpicTimelineData, error)
	GetSprintTimelineData(projectID uuid.UUID) (*SprintTimelineData, error)

	// Google Slides Import: which slide indices were already imported from a presentation
	GetImportedSlideIndicesByPresentationID(presentationID string) ([]int, error)

	// Project Finance (Internal VC — per-project capital)
	UpdateProjectCapital(projectID uuid.UUID, newBalance float64, bonusPct *float64) error
	CreateProjectTransaction(tx *ProjectTransaction) error
	GetProjectTransactions(projectID uuid.UUID) ([]ProjectTransaction, error)
	DeleteProjectTransaction(txID int64, projectID uuid.UUID) error

	// Internal B2B Outsource Requests
	CreateB2BRequest(req *B2BRequest) error
	GetB2BRequests(teamID uint, direction string) ([]B2BRequest, error) // direction: "inbound" | "outbound"
	GetB2BRequestByID(id uuid.UUID) (*B2BRequest, error)
	UpdateB2BRequest(req *B2BRequest) error
	GetFirstProjectByTeamID(teamID uint) (*Project, error)

	// Project Backups
	CreateProjectBackup(backup *ProjectBackup) error
	GetProjectBackups(projectID uuid.UUID) ([]ProjectBackup, error)
	GetProjectBackupByID(id uuid.UUID) (*ProjectBackup, error)
	DeleteProjectBackup(id uuid.UUID, projectID uuid.UUID) error

	// Komgrip: project-less personal/misc tasks
	GetKomgripTasks(userID uint) ([]Task, error)
}

// GanttData is the payload for GET /sentinel/tasks/gantt (all tasks + dependencies)
type GanttData struct {
	Tasks        []Task           `json:"tasks"`
	Dependencies []TaskDependency `json:"dependencies"`
}

// SentinelUsecase defines the business logic
type SentinelUsecase interface {
	// Projects
	CreateProject(name, description, status string, ctx CallerContext) (*Project, error)
	GetProjects(ctx CallerContext) ([]Project, error)
	GetProjectDetails(id uuid.UUID, ctx CallerContext) (*Project, error)
	GetProjectByIDOrCode(idOrCode string, ctx CallerContext) (*Project, error)                 // UUID or project code (e.g. mims-hdmap-main)
	GetProjectDetailsPage(idOrCode string, taskLimit int, ctx CallerContext) (*ProjectDetailsResponse, error) // Combined project + tasks + sprints + milestones + epics (1 round-trip)
	GetProjectTasksPage(idOrCode string, limit int, cursorCreatedAt, cursorID string, offset int, ctx CallerContext) (*ProjectTasksPageResponse, error)
	UpdateProject(projectID uuid.UUID, name, description, status string, updateCode bool) (*Project, error)
	DeleteProject(id uuid.UUID) error
	AssignProjectTeam(projectID uuid.UUID, teamID *uint, requesterRole string) (*Project, error) // CEO only
	AssignProjectPmOwners(projectID uuid.UUID, pmUserIDs []uint, requesterRole string) (*Project, error) // CEO/MANAGER when teams feature off (param name kept for compatibility)

	CreateTask(title, desc, taskType string, creatorID uint, dueDate *time.Time, projectID, parentID *uuid.UUID, startDate, endDate *time.Time, priority string, storyPoints float64, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID, estimatedMinutes *int) (*Task, error)
	AssignTask(taskID uuid.UUID, devID uint, assignerID uint, assignerRole string) error
	SubmitWork(taskID uuid.UUID, devID uint, referenceURL, note string) (*Submission, error)
	SubmitUAT(taskID uuid.UUID, devID uint, payload UATPayloadData) error // Engineer submits UAT payload for a FEATURE (READY_FOR_UAT → REVIEW_PENDING)
	GetTaskByID(taskID uuid.UUID) (*Task, error)
	GetTaskByIDOrCode(idOrCode string) (*Task, error) // idOrCode is UUID or task code (e.g. mims-hdmap-main-001)
	GetTaskActivityTimeline(taskID uuid.UUID) ([]TaskActivityItem, error)
	GetMyTasks(userID uint) ([]Task, error)
	GetMyActiveSprints(userID uint) ([]Sprint, error)                                      // Active sprints containing user's tasks
	GetGlobalActiveTasks(ctx CallerContext) ([]GlobalActiveTask, error)                    // TASK/BUG in ACTIVE sprints; CEO/MANAGER company-wide; Product Owner/engineer per project list rules
	GetTeamActiveTasks(ctx CallerContext) ([]GlobalActiveTask, error)    // ACTIVE-sprint TASK/BUG items visible to caller
	GetActiveFeatures(callerTeamID *uint, callerRole string, projectID *uuid.UUID) ([]FeatureRoadmapItem, error) // FEATURE items for Product Owner/CEO Roadmap Board (optional project scope)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksByProjectID(projectID uuid.UUID) ([]Task, error)
	GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]Task, error)
	GetGanttData(projectID *uuid.UUID) (*GanttData, error) // All tasks + dependencies; if projectID set, filter by project
	GetPendingApprovals(userRole string) ([]Task, error)   // Approvals inbox for Product Owner/CEO

	// Task Dependencies (Gantt links)
	AddDependency(predecessorID, successorID uuid.UUID, depType string) (*TaskDependency, error)
	RemoveDependency(id uuid.UUID) error

	// Task Management with Access Control
	UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description, taskType string, parentID *uuid.UUID, dueAt, startDate, endDate *time.Time, progress *int, priority string, storyPoints *float64, sprintID *uuid.UUID, applySprint bool, milestoneID *uuid.UUID, epicID *uuid.UUID, applyEpic bool, sortOrder *int, estimatedMinutes *int) (*Task, error)
	UpdateTaskResourceURLs(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, resourceURLs datatypes.JSON) (*Task, error)
	EstimateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) (*Task, error)                            // Kept for AI scheduling (ScheduleProjectWithAI)
	GenerateProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (*AIGeneratedPlan, error)       // AI generates epics, milestones, sprints, tasks
	ClearProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) error                              // Remove all tasks, sprints, milestones, epics (CEO / Product Owner)
	ScheduleProjectWithAI(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (updatedCount int, err error) // Estimate + schedule existing tasks (CEO / Product Owner)
	DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error

	// Split Task: decompose one task into N smaller sub-tasks, then delete the original
	SplitTask(taskID uuid.UUID, splits []SplitTaskItem, requestingUserID uint, requestingUserRole string) ([]*Task, error)

	// Appeal System
	SubmitAppeal(submissionID uuid.UUID, devID uint, reason string) (*Appeal, error)
	ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error

	// Time Negotiation
	NegotiateTime(taskID uuid.UUID, devID uint, minutes int, reason string) error

	// Human Quality Gate
	ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error
	RejectTask(taskID uuid.UUID, rejectorID uint, rejectorRole string, reason string) error

	// Continuous UAT: sub-task level testing queue (READY_FOR_TEST lane)
	MarkReadyForTest(taskID uuid.UUID, devID uint) error                                                     // Engineer marks TASK/BUG as ready for Product Owner to test
	PMApproveSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, testURL string, testSteps string) error // Product Owner approves READY_FOR_TEST → WAIT_FOR_DEPLOY
	ApproveSubTask(taskID uuid.UUID, ceoUserID uint, ceoRole string) error                                   // CEO final approves READY_FOR_UAT → COMPLETED
	RejectSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, reason string) error                       // Product Owner/CEO rejects → IN_PROGRESS
	AdvanceTaskAfterDeploy(taskID uuid.UUID, deployedByUserID uint) error                                     // Chief Engineer marks deployed: WAIT_FOR_DEPLOY → READY_FOR_UAT
	GetTasksReadyForTest(callerTeamID *uint, callerRole string) ([]GlobalActiveTask, error)                  // Product Owner/CEO: fetch READY_FOR_TEST tasks for team
	GetTasksReadyForCEOApproval(callerTeamID *uint, callerRole string) ([]GlobalActiveTask, error)           // CEO: fetch READY_FOR_UAT tasks awaiting final approval

	// Sprints
	CreateSprint(projectID uuid.UUID, name, goal string, startDate, endDate *time.Time) (*Sprint, error)
	GetSprintsByProject(projectID uuid.UUID) ([]Sprint, error)
	StartSprint(sprintID uuid.UUID) (*Sprint, error)
	CompleteSprint(sprintID uuid.UUID) (*Sprint, error)
	ReopenSprint(sprintID uuid.UUID) (*Sprint, error)
	AddTasksToSprint(sprintID uuid.UUID, taskIDs []uuid.UUID) error
	UpdateSprint(sprintID uuid.UUID, name, goal string, startDate, endDate *time.Time, sortOrder *int) (*Sprint, error)
	DeleteSprint(sprintID uuid.UUID) error

	// Milestones
	CreateMilestone(projectID uuid.UUID, title, description string, dueDate *time.Time) (*Milestone, error)
	GetMilestonesByProject(projectID uuid.UUID) ([]Milestone, error)
	UpdateMilestone(id uuid.UUID, title, description, status string, dueDate *time.Time) (*Milestone, error)
	DeleteMilestone(id uuid.UUID) error

	// Comments
	AddComment(taskID uuid.UUID, userID uint, content string, attachments []TaskCommentAttachment) (*TaskComment, error)
	GetComments(taskID uuid.UUID) ([]TaskComment, error)
	EditComment(commentID uuid.UUID, editorUserID uint, content string) (*TaskComment, error)

	// Time Logging
	LogTime(taskID uuid.UUID, userID uint, minutes int, description, workType string, loggedDate *time.Time, isTimer bool) (*TimeLog, error)
	GetTimeLogs(taskID uuid.UUID) ([]TimeLog, error)
	EditTimeLog(logID uuid.UUID, callerID uint, minutes int, description, workType string, taskID *uuid.UUID) (*TimeLog, error)
	DeleteTimeLog(logID uuid.UUID, callerID uint) error
	GetMyDailyTimeLogs(userID uint, date time.Time) (*DailyTimeLogSummary, error)
	BulkLogTime(entries []BulkLogEntry, userID uint) ([]BulkLogResult, error)

	// Analytics
	GetProjectAnalytics(projectID uuid.UUID) (*ProjectAnalytics, error)

	// Bulk Operations
	BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string, actorID uint) error

	// System Configuration Management
	GetSystemConfig() (*SystemConfig, error)
	UpdateSystemConfig(activeModel string, temperature float32, cursorAssistance int, userRole string) (*SystemConfig, error)
	GetAvailableModels() []string
	GetAIUsage() AIUsage

	// Epics (Hierarchy Dimension 1)
	CreateEpic(projectID uuid.UUID, title, description, color string, startDate, endDate *time.Time) (*Epic, error)
	GetEpicsByProject(projectID uuid.UUID) ([]Epic, error)
	UpdateEpic(epicID uuid.UUID, title, description, status, color string, sortOrder *int, startDate, endDate *time.Time) (*Epic, error)
	DeleteEpic(epicID uuid.UUID) error

	// Timeline Views (Matrix Dimension)
	GetEpicTimelineData(projectID uuid.UUID) (*EpicTimelineData, error)
	GetSprintTimelineData(projectID uuid.UUID) (*SprintTimelineData, error)

	// Timeline PDF Export (chromedp → PDF bytes, same pattern as mims-api-service)
	ExportTimelinePDF(projectID uuid.UUID, mode string, templateDir string) ([]byte, string, error)

	// Google Slides Import
	PreviewGoogleSlides(req *PreviewGoogleSlidesRequest, serverAPIKey string) (*PreviewGoogleSlidesResult, error)
	ImportFromGoogleSlides(req *ImportGoogleSlidesRequest, serverAPIKey string, creatorID uint) (*ImportGoogleSlidesResult, error)

	// Google Sheets Import
	PreviewGoogleSheets(req *PreviewGoogleSheetsRequest) (*PreviewGoogleSheetsResult, error)
	ImportFromGoogleSheets(req *ImportGoogleSheetsRequest, creatorID uint) (*ImportGoogleSheetsResult, error)

	// Canva Import (Connect API → PPTX → same pipeline as Slides)
	PreviewCanva(req *PreviewCanvaRequest, accessToken string) (*PreviewCanvaResult, error)
	ImportFromCanva(req *ImportCanvaRequest, accessToken string, creatorID uint) (*ImportCanvaResult, error)

	// PPTX File Upload Import (no API key required — user exports PPTX manually)
	PreviewPPTX(data []byte) (*PreviewPPTXResult, error)
	ImportFromPPTX(data []byte, req *ImportPPTXRequest, creatorID uint) (*ImportPPTXResult, error)

	// Internal B2B Outsource
	CreateB2BRequest(title, description string, estimatedMinutes int, requesterTeamID, targetTeamID, requesterUserID uint) (*B2BRequest, error)
	GetB2BRequests(callerTeamID uint, direction string) ([]B2BRequest, error)
	CounterOfferB2BRequest(id uuid.UUID, callerTeamID uint, proposedMinutes int, reason string) (*B2BRequest, error)
	RejectB2BRequest(id uuid.UUID, callerTeamID uint) (*B2BRequest, error)
	AcceptB2BRequest(id uuid.UUID, callerTeamID uint, accepterUserID uint) (*Task, error)

	// Project Backups
	CreateProjectBackup(projectID uuid.UUID, label string, createdBy *uint) (*ProjectBackup, error)
	GetProjectBackups(projectID uuid.UUID) ([]ProjectBackup, error)
	GetProjectBackupPayload(projectID uuid.UUID, backupID uuid.UUID) (*ProjectBackupPayload, error)
	RestoreProjectBackup(backupID uuid.UUID, projectID uuid.UUID) error
	DeleteProjectBackup(backupID uuid.UUID, projectID uuid.UUID) error
	ImportProjectFromBackup(newName string, payload *ProjectBackupPayload, createdBy *uint) (*Project, error)

	// Komgrip: project-less personal/misc tasks (all employees)
	CreateKomgripTask(title, description string, creatorID uint, priority string, estimatedMinutes int) (*Task, error)
	GetKomgripTasks(userID uint) ([]Task, error)
	UpdateKomgripTaskStatus(taskID uuid.UUID, status string, userID uint) (*Task, error)
}

// SlideComment represents a comment/annotation on a Google Slides slide
type SlideComment struct {
	Content  string `json:"content"`
	Author   string `json:"author"`
	Resolved bool   `json:"resolved"`
}

// SlideResourceURLs is stored in task.resource_urls for imported slides
type SlideResourceURLs struct {
	ThumbnailURL   string         `json:"thumbnail_url"`
	Images         []string       `json:"images"` // base64 data URLs of embedded images
	SlideURL       string         `json:"slide_url"`
	Source         string         `json:"source"`
	SlideIndex     int            `json:"slide_index"`
	PresentationID string         `json:"presentation_id"`
	Comments       []SlideComment `json:"comments"`
}

// PreviewGoogleSlidesRequest is the payload for POST /sentinel/import/google-slides/preview
type PreviewGoogleSlidesRequest struct {
	PresentationURL string `json:"presentation_url" binding:"required"`
	APIKey          string `json:"api_key"`
}

// PreviewSlideItem is one slide in the preview list (index, title, and hidden flag)
type PreviewSlideItem struct {
	Index  int    `json:"index"`
	Title  string `json:"title"` // structural slide title from deck (for reference / secondary label)
	// SuggestedTaskTitle is derived from visible text on the slide (body), not the placeholder title — default for Task Title in triage.
	SuggestedTaskTitle string `json:"suggested_task_title,omitempty"`
	Hidden             bool   `json:"hidden"` // true when slide is hidden/skipped in presentation (PPTX show="0"); import should skip by default
}

// PreviewGoogleSlidesResult is the response for the preview endpoint
type PreviewGoogleSlidesResult struct {
	PresentationTitle string             `json:"presentation_title"`
	PresentationID    string             `json:"presentation_id"` // extracted from URL for "already imported" check
	Slides            []PreviewSlideItem `json:"slides"`
	// AlreadyImportedSlideIndices: 1-based slide indices that were previously imported from this presentation (so UI can pre-uncheck them and select only new slides)
	AlreadyImportedSlideIndices []int `json:"already_imported_slide_indices"`
	// ImportMode: "pptx_only" (no API), "pptx_with_api" (PPTX + API key valid), "api_only" (PPTX failed, API used)
	ImportMode string `json:"import_mode"`
	// APIKeyStatus: "not_provided", "valid", "invalid"
	APIKeyStatus string `json:"api_key_status"`
	// APIKeyError: when api_key_status is "invalid", short reason from Google (e.g. API not enabled, key invalid)
	APIKeyError string `json:"api_key_error,omitempty"`
}

// TriagedSlide holds per-slide triage data filled by Product Owner during manual triage before import.
type TriagedSlide struct {
	SlideIndex       int    `json:"slide_index"`       // 1-based index (matches PreviewSlideItem.Index)
	Title            string `json:"title"`             // editable task title
	AssigneeID       *uint  `json:"assignee_id"`       // optional: user id to assign the task to
	EstimatedMinutes int    `json:"estimated_minutes"` // mandatory manual estimate
	Priority         string `json:"priority"`          // CRITICAL | HIGH | MEDIUM | LOW
}

// ImportGoogleSlidesRequest is the payload for POST /sentinel/import/google-slides
// Either sprint_id (import into sprint) or epic_id (import into backlog/epic) or neither (import into backlog unassigned).
type ImportGoogleSlidesRequest struct {
	PresentationURL string `json:"presentation_url" binding:"required"`
	SprintID        string `json:"sprint_id"` // optional: when set, tasks go to this sprint
	EpicID          string `json:"epic_id"`   // optional: when set, tasks go to this epic (backlog)
	ProjectID       string `json:"project_id" binding:"required"`
	ParentID        string `json:"parent_id"` // optional: when set, tasks are created as sub-tasks under this parent
	APIKey          string `json:"api_key"`
	Priority        string `json:"priority"`
	StoryPoints     float64 `json:"story_points"`
	// SlideIndices: 1-based indices of slides to import. If empty or nil, import all.
	// Deprecated in favour of Slides when triage data is provided.
	SlideIndices []int `json:"slide_indices"`
	// Slides: per-slide triage data (title, assignee, estimated_minutes, priority).
	// When provided, SlideIndices is derived from this array automatically.
	Slides []TriagedSlide `json:"slides"`
}

// ImportGoogleSlidesResult is the response for POST /sentinel/import/google-slides
type ImportGoogleSlidesResult struct {
	CreatedCount      int     `json:"created_count"`
	SlideCount        int     `json:"slide_count"`
	PresentationTitle string  `json:"presentation_title"`
	Tasks             []*Task `json:"tasks"`
}

// SheetRowPreviewItem is one data row from a Google Sheet CSV preview.
type SheetRowPreviewItem struct {
	RowIndex  int    `json:"row_index"`
	Title     string `json:"title"`
	DueDate   string `json:"due_date"`   // YYYY-MM-DD or empty
	Status    string `json:"status"`     // mapped Sentinel status
	RawStatus string `json:"raw_status"` // original cell text
	Notes     string `json:"notes"`
	// IOD-specific extra fields (empty for non-IOD sheets)
	Header        string   `json:"header"`         // Section / page name (col Header)
	HeaderLink    string   `json:"header_link"`    // URL of the page being tested
	RequestMethod string   `json:"request_method"` // GET / POST / etc.
	Payload       string   `json:"payload"`        // Request / response payload
	ImageRef      string   `json:"image_ref"`      // Image column text (URL if available)
	DetailLinks   []string `json:"detail_links"`   // Additional URLs extracted from Detail cell
}

// PreviewGoogleSheetsRequest is the payload for POST /sentinel/import/google-sheets/preview
type PreviewGoogleSheetsRequest struct {
	SheetURL string `json:"sheet_url" binding:"required"`
}

// PreviewGoogleSheetsResult is the response for the sheets preview endpoint
type PreviewGoogleSheetsResult struct {
	SheetTitle string                `json:"sheet_title"`
	SheetID    string                `json:"sheet_id"`
	Rows       []SheetRowPreviewItem `json:"rows"`
}

// TriagedSheetRow holds per-row overrides before import from Google Sheets.
type TriagedSheetRow struct {
	RowIndex         int    `json:"row_index"`
	Title            string `json:"title"`
	Priority         string `json:"priority"` // CRITICAL | HIGH | MEDIUM | LOW
	EstimatedMinutes int    `json:"estimated_minutes"`
	DueDate          string `json:"due_date"` // YYYY-MM-DD
	Status           string `json:"status"`
	Notes            string `json:"notes"`
	// IOD-specific extra fields
	Header        string   `json:"header"`
	HeaderLink    string   `json:"header_link"`
	RequestMethod string   `json:"request_method"`
	Payload       string   `json:"payload"`
	ImageRef      string   `json:"image_ref"`
	DetailLinks   []string `json:"detail_links"` // Additional URLs extracted from Detail cell
}

// ImportGoogleSheetsRequest is the payload for POST /sentinel/import/google-sheets
type ImportGoogleSheetsRequest struct {
	SheetURL   string            `json:"sheet_url" binding:"required"`
	SheetTitle string            `json:"sheet_title"` // optional; from preview for result display
	ProjectID  string            `json:"project_id" binding:"required"`
	SprintID   string            `json:"sprint_id"`
	EpicID     string            `json:"epic_id"`
	ParentID   string            `json:"parent_id"`
	Rows       []TriagedSheetRow `json:"rows" binding:"required"`
}

// ImportGoogleSheetsResult is the response for POST /sentinel/import/google-sheets
type ImportGoogleSheetsResult struct {
	CreatedCount int     `json:"created_count"`
	SheetTitle   string  `json:"sheet_title"`
	Tasks        []*Task `json:"tasks"`
}

// PreviewCanvaRequest is the payload for POST /sentinel/import/canva/preview
type PreviewCanvaRequest struct {
	DesignURL string `json:"design_url" binding:"required"`
}

// PreviewCanvaResult lists pages (1..page_count) for triage before import
type PreviewCanvaResult struct {
	DesignTitle string             `json:"design_title"`
	DesignID    string             `json:"design_id"`
	Pages       []PreviewSlideItem `json:"pages"`
}

// ImportCanvaRequest is the payload for POST /sentinel/import/canva
type ImportCanvaRequest struct {
	DesignURL   string `json:"design_url" binding:"required"`
	ProjectID   string `json:"project_id" binding:"required"`
	SprintID    string `json:"sprint_id"`
	EpicID      string `json:"epic_id"`
	ParentID    string `json:"parent_id"`
	Priority    string `json:"priority"`
	StoryPoints float64    `json:"story_points"`
	// Pages: per-page triage (slide_index = 1-based page index, reuses TriagedSlide)
	Pages []TriagedSlide `json:"pages" binding:"required"`
}

// ImportCanvaResult is the response for POST /sentinel/import/canva
type ImportCanvaResult struct {
	CreatedCount int     `json:"created_count"`
	PageCount    int     `json:"page_count"`
	DesignTitle  string  `json:"design_title"`
	Tasks        []*Task `json:"tasks"`
}

// PreviewPPTXResult is the response for POST /sentinel/import/pptx/preview (multipart file upload)
type PreviewPPTXResult struct {
	Title  string             `json:"title"`
	Slides []PreviewSlideItem `json:"slides"`
}

// ImportPPTXRequest is the JSON payload (sent as form field "payload") for POST /sentinel/import/pptx
type ImportPPTXRequest struct {
	ProjectID   string         `json:"project_id"`
	SprintID    string         `json:"sprint_id"`
	EpicID      string         `json:"epic_id"`
	ParentID    string         `json:"parent_id"`
	Priority    string         `json:"priority"`
	StoryPoints float64         `json:"story_points"`
	Pages       []TriagedSlide `json:"pages"`
}

// ImportPPTXResult is the response for POST /sentinel/import/pptx
type ImportPPTXResult struct {
	CreatedCount int     `json:"created_count"`
	PageCount    int     `json:"page_count"`
	Title        string  `json:"title"`
	Tasks        []*Task `json:"tasks"`
}

// FeatureRoadmapItem is a FEATURE-type task enriched with project identity and roll-up progress.
// Used by the Feature Roadmap Board (Product Owner/CEO management view).
type FeatureRoadmapItem struct {
	Task
	ProjectName    string `json:"project_name"`
	ProjectColor   string `json:"project_color"`
	ProjectCode    string `json:"project_code"`    // URL slug (e.g. mims-hd-map); empty for legacy
	RollupProgress int    `json:"rollup_progress"` // 0-100 percentage of completed child tasks
	ChildTasks     []Task `json:"child_tasks"`
}

// SplitTaskItem describes one slice that the original task is split into
type SplitTaskItem struct {
	Title            string `json:"title" binding:"required"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	AssigneeID       *uint  `json:"assignee_id"`
	Priority         string `json:"priority"` // inherit from parent if empty
}

// AIGeneratedPlan is the structured output from AI for generating a project work plan
type AIGeneratedPlan struct {
	Epics      []AIPlanEpic      `json:"epics"`
	Milestones []AIPlanMilestone `json:"milestones"`
	Sprints    []AIPlanSprint    `json:"sprints"`
	Tasks      []AIPlanTask      `json:"tasks"`
}
type AIPlanEpic struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
type AIPlanMilestone struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // YYYY-MM-DD
}
type AIPlanSprint struct {
	Name      string `json:"name"`
	Goal      string `json:"goal"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
type AIPlanTask struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Priority       string `json:"priority"`
	StoryPoints    float64 `json:"story_points"`
	EpicIndex      *int   `json:"epic_index"`      // 0-based index into epics list
	SprintIndex    *int   `json:"sprint_index"`    // 0-based index into sprints list
	MilestoneIndex *int   `json:"milestone_index"` // 0-based index into milestones list
	StartDate      string `json:"start_date"`      // YYYY-MM-DD
	EndDate        string `json:"end_date"`        // YYYY-MM-DD
}

// TaskEstimateInput is a minimal task info sent to AI for batch estimate + order
type TaskEstimateInput struct {
	Index       int    `json:"task_index"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	StoryPoints float64 `json:"story_points"`
}

// TaskEstimateAndOrder is AI output per task: estimated minutes and suggested execution order (1-based)
type TaskEstimateAndOrder struct {
	TaskIndex int `json:"task_index"` // 0-based, matches input
	Minutes   int `json:"minutes"`
	Order     int `json:"order"` // 1-based execution order (1 = first task to do)
}

// AIUsage holds approximate Gemini API usage for display (we track calls; limits from env or default).
type AIUsage struct {
	RequestsLastMinute int `json:"requests_last_minute"`
	RequestsToday      int `json:"requests_today"`
	LimitRPM           int `json:"limit_rpm"`
	LimitRPD           int `json:"limit_rpd"`
	RemainingRPM       int `json:"remaining_rpm"`
	RemainingRPD       int `json:"remaining_rpd"`
}

// UsageTracker records Gemini API calls and returns approximate usage (thread-safe).
type UsageTracker interface {
	RecordRequest()
	GetUsage(limitRPM, limitRPD int) AIUsage
}

// AIService defines the interface for AI operations (Port)
type AIService interface {
	// ListModels returns model IDs from Gemini API (e.g. gemini-2.5-flash-lite). Empty or error = use fallback list.
	ListModels() ([]string, error)

	// EstimateEffort รับ Title/Desc แล้วคืนค่าเป็น นาที (minutes) และเหตุผล
	EstimateEffort(title, description string) (minutes int, reasoning string, err error)

	// EstimateAndScheduleTasks ประเมินเวลาและลำดับการทำของแต่ละ task จากข้อมูลที่มี คืนค่า minutes และ order (1-based)
	EstimateAndScheduleTasks(inputs []TaskEstimateInput) ([]TaskEstimateAndOrder, error)

	// GenerateWorkPlan สร้างแผนงาน (epics, milestones, sprints, tasks) จากชื่อและคำอธิบายโปรเจกต์
	GenerateWorkPlan(projectName, projectDescription string) (*AIGeneratedPlan, error)

	// ReviewCode วิเคราะห์ code diff และคืนค่า verdict (PASS/FAIL), score (0-100), feedback
	ReviewCode(diff string) (verdict string, score int, feedback string, err error)

	// AnalyzeAppeal วิเคราะห์ความน่าเชื่อถือของ Appeal เพื่อแนะนำ CEO / Product Owner
	// Returns: recommendation (UPHOLD/OVERTURN), confidence (0-100), reasoning, error
	AnalyzeAppeal(diff string, originalFeedback string, appealReason string) (recommendation string, confidence int, reasoning string, err error)

	// AnalyzeTimeNegotiation วิเคราะห์คำขอเจรจาเวลาจากนักพัฒนา
	// Returns: recommendation (APPROVE/REJECT), confidence (0-100), reasoning, error
	AnalyzeTimeNegotiation(taskTitle, taskDescription string, aiEstimate, devProposal int, devReason string) (recommendation string, confidence int, reasoning string, err error)
}

// B2BRequest represents an internal cross-team outsource request.
// Team A (requester) asks Team B (target) to handle work described by this request.
// Status flow: PENDING → COUNTER_OFFERED → ACCEPTED | REJECTED
type B2BRequest struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title             string    `json:"title" gorm:"not null"`
	Description       string    `json:"description" gorm:"type:text"`
	EstimatedMinutes  int       `json:"estimated_minutes" gorm:"not null;default:0"`
	ProposedMinutes   int       `json:"proposed_minutes" gorm:"default:0"`
	NegotiationReason string    `json:"negotiation_reason" gorm:"type:text"`
	Status            string    `json:"status" gorm:"default:'PENDING';not null"`

	RequesterTeamID uint `json:"requester_team_id" gorm:"not null"`
	TargetTeamID    uint `json:"target_team_id" gorm:"not null"`
	RequesterUserID uint `json:"requester_user_id" gorm:"not null"`

	CreatedTaskID *uuid.UUID `json:"created_task_id,omitempty" gorm:"type:uuid"`

	// Enriched fields (not in DB)
	RequesterTeamName string `json:"requester_team_name,omitempty" gorm:"-"`
	TargetTeamName    string `json:"target_team_name,omitempty" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (B2BRequest) TableName() string { return "b2b_requests" }

// ProjectBackupPayload is the full snapshot of a project stored as JSONB.
type ProjectBackupPayload struct {
	Project    Project     `json:"project"`
	Epics      []Epic      `json:"epics"`
	Sprints    []Sprint    `json:"sprints"`
	Milestones []Milestone `json:"milestones"`
	Tasks      []Task      `json:"tasks"`
}

// ProjectBackup is a point-in-time snapshot of a project and its child records.
type ProjectBackup struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index"`
	Label     string         `json:"label" gorm:"type:varchar(255);default:''"`
	Payload   datatypes.JSON `json:"payload" gorm:"type:jsonb;not null;default:'{}'"`
	CreatedBy *uint          `json:"created_by"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

func (ProjectBackup) TableName() string { return "project_backups" }

// CreateProjectBackupRequest is the DTO for creating a backup.
type CreateProjectBackupRequest struct {
	Label string `json:"label"`
}

// ImportProjectFromBackupRequest is the DTO for POST /sentinel/projects/import-backup
type ImportProjectFromBackupRequest struct {
	Name    string               `json:"name" binding:"required"`
	Payload ProjectBackupPayload `json:"payload" binding:"required"`
}
