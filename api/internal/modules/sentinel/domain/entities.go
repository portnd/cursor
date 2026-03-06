package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

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
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Tasks       []Task    `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
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
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID    uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	UserEmail string    `json:"user_email,omitempty" gorm:"-"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (TaskComment) TableName() string { return "task_comments" }

// TimeLog represents actual time logged by a developer on a task
type TimeLog struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID      uuid.UUID `json:"task_id" gorm:"type:uuid;not null;index"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	UserEmail   string    `json:"user_email,omitempty" gorm:"-"`
	Minutes     int       `json:"minutes" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	LoggedAt    time.Time `json:"logged_at" gorm:"autoCreateTime"`
}

func (TimeLog) TableName() string { return "time_logs" }

// ProjectAnalytics holds computed analytics for a project
type ProjectAnalytics struct {
	ProjectID       uuid.UUID          `json:"project_id"`
	TotalTasks      int                `json:"total_tasks"`
	CompletedTasks  int                `json:"completed_tasks"`
	TotalStoryPoints int               `json:"total_story_points"`
	CompletedSP     int                `json:"completed_story_points"`
	TotalLoggedMinutes int             `json:"total_logged_minutes"`
	AvgCycleTimeDays   float64         `json:"avg_cycle_time_days"`
	Burndown        []BurndownPoint    `json:"burndown"`
	Velocity        []VelocityPoint    `json:"velocity"`
	TeamCapacity    []TeamCapacityRow  `json:"team_capacity"`
}

type BurndownPoint struct {
	Day       string `json:"day"`
	Ideal     float64 `json:"ideal"`
	Remaining float64 `json:"remaining"`
}

type VelocityPoint struct {
	SprintName  string `json:"sprint_name"`
	CompletedSP int    `json:"completed_sp"`
	PlannedSP   int    `json:"planned_sp"`
}

type TeamCapacityRow struct {
	UserID          uint    `json:"user_id"`
	UserEmail       string  `json:"user_email"`
	AssignedTasks   int     `json:"assigned_tasks"`
	EstimatedHours  float64 `json:"estimated_hours"`
	LoggedHours     float64 `json:"logged_hours"`
	Utilization     float64 `json:"utilization_pct"`
}

func (Project) TableName() string { return "projects" }

// Task represents a work assignment
type Task struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code               string         `json:"code" gorm:"type:varchar(64);uniqueIndex"` // e.g. mims-hdmap-main-001 (empty for legacy tasks)
	Title              string         `json:"title" gorm:"not null"`
	Description        string         `json:"description"`
	ResourceURLs       datatypes.JSON `json:"resource_urls" gorm:"type:jsonb;default:'{}'"`
	AIEstimatedMinutes int            `json:"ai_estimated_minutes"`

	// Project: every task belongs to a project
	ProjectID *uuid.UUID `json:"project_id" gorm:"type:uuid;index"`

	// Epic linking (Hierarchy Dimension 1)
	EpicID *uuid.UUID `json:"epic_id,omitempty" gorm:"type:uuid;index"`
	Epic   *Epic      `json:"epic,omitempty" gorm:"foreignKey:EpicID"`

	// Sprint & Milestone linking
	SprintID    *uuid.UUID `json:"sprint_id,omitempty" gorm:"type:uuid;index"`
	MilestoneID *uuid.UUID `json:"milestone_id,omitempty" gorm:"type:uuid;index"`

	// Priority & Estimation
	Priority    string `json:"priority" gorm:"default:'MEDIUM'"` // CRITICAL, HIGH, MEDIUM, LOW
	StoryPoints int    `json:"story_points" gorm:"default:0"`

	// WBS / Gantt: sub-tasks and planned dates
	ParentID  *uuid.UUID `json:"parent_id" gorm:"type:uuid;index"`             // For sub-tasks (Work Breakdown Structure)
	SortOrder int        `json:"sort_order" gorm:"default:0"`                 // Order within epic (backlog drag-and-drop)
	StartDate *time.Time `json:"start_date"`                                    // Planned start date
	EndDate   *time.Time `json:"end_date"`                                      // Planned end date
	Progress  int        `json:"progress" gorm:"default:0"`                     // 0 to 100%
	SubTasks  []Task     `json:"sub_tasks,omitempty" gorm:"foreignKey:ParentID"` // Has many sub-tasks

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

	Status string `json:"status" gorm:"default:'PENDING';index"`

	// Relationships (uint to match User ID)
	AssignedTo   *uint `json:"assigned_to"`   // Developer assigned to this task
	AssignedByID *uint `json:"assigned_by_id" gorm:"index"` // PM/CEO who assigned the task (for PM-scoped leaderboard)
	CreatedBy    *uint `json:"created_by"`

	// Enriched from auth (not stored in DB)
	CreatedByRole  string `json:"created_by_role,omitempty" gorm:"-"`
	CreatedByEmail string `json:"created_by_email,omitempty" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationship: Submissions history (loaded via Preload)
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:TaskID"`
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

// Submission represents a code push for review
type Submission struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID     uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	DevID      uint      `json:"dev_id" gorm:"not null"`
	CommitHash string    `json:"commit_hash" gorm:"not null"`
	Diff       string    `json:"diff" gorm:"type:text"` // Code diff for appeal analysis

	// AI Logic
	AIVerdict  string         `json:"ai_verdict"` // PASS, FAIL, PENDING
	AIScore    int            `json:"ai_score"`
	AIFeedback datatypes.JSON `json:"ai_feedback" gorm:"type:jsonb;default:'{}'"`

	// Appeal System
	IsOverridden bool    `json:"is_overridden" gorm:"default:false"` // True if appeal approved
	Appeal       *Appeal `json:"appeal,omitempty" gorm:"foreignKey:SubmissionID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
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
	AIReasoning      string `json:"ai_reasoning" gorm:"type:text"`      // Explanation for CEO/PM

	ResolverID   *uint  `json:"resolver_id"`   // PM/CEO who resolved
	ResolverNote string `json:"resolver_note"` // Admin's decision note

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Appeal) TableName() string { return "appeals" }

// SystemConfig represents dynamic AI configuration (singleton pattern)
type SystemConfig struct {
	ID               uint    `json:"-" gorm:"primaryKey"`
	ActiveModel      string  `json:"active_model" gorm:"default:'gemini-2.5-flash-lite'"`
	Temperature      float32 `json:"temperature" gorm:"default:0.4"` // 0.0 (Stable) to 1.0 (Creative)
	CursorAssistance int     `json:"cursor_assistance" gorm:"default:80"` // 0 to 100% (AI assistance level)
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (SystemConfig) TableName() string { return "system_configs" }

// SentinelRepository defines the interface (Port)
type SentinelRepository interface {
	// Projects
	CreateProject(p *Project) error
	GetAllProjects() ([]Project, error)
	GetProjectByID(id uuid.UUID) (*Project, error)
	GetProjectByCode(code string) (*Project, error)
	GetTasksByProjectID(projectID uuid.UUID) ([]Task, error)
	UpdateProject(p *Project) error
	DeleteProject(id uuid.UUID) error
	DeleteProjectPlan(projectID uuid.UUID) error // Remove all tasks, sprints, milestones, epics for the project

	CreateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	GetTaskByCode(code string) (*Task, error)
	CountTasksForCode(projectID *uuid.UUID) (int, error)
	GetMaxTaskCodeSuffix(prefix string) (int, error) // max numeric suffix for code like "prefix-001"; used for globally unique codes
	GetTasksByAssignee(userID uint) ([]Task, error)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksRequiringApproval() ([]Task, error) // Tasks with PENDING appeals or time negotiations
	UpdateTask(task *Task) error
	DeleteTask(id uuid.UUID) error // Delete a task by ID

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
	ApproveTask(id uuid.UUID) error // Mark task as COMPLETED and set CompletedAt

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

	// Time Logs
	CreateTimeLog(t *TimeLog) error
	GetTimeLogsByTaskID(taskID uuid.UUID) ([]TimeLog, error)
	GetTotalLoggedMinutes(taskID uuid.UUID) (int, error)

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
}

// GanttData is the payload for GET /sentinel/tasks/gantt (all tasks + dependencies)
type GanttData struct {
	Tasks         []Task           `json:"tasks"`
	Dependencies  []TaskDependency `json:"dependencies"`
}

// SentinelUsecase defines the business logic
type SentinelUsecase interface {
	// Projects
	CreateProject(name, description, status string) (*Project, error)
	GetProjects() ([]Project, error)
	GetProjectDetails(id uuid.UUID) (*Project, error)
	GetProjectByIDOrCode(idOrCode string) (*Project, error) // UUID or project code (e.g. mims-hdmap-main)
	UpdateProject(projectID uuid.UUID, name, description, status string, updateCode bool) (*Project, error)
	DeleteProject(id uuid.UUID) error

	CreateTask(title, desc string, creatorID uint, dueDate *time.Time, projectID, parentID *uuid.UUID, startDate, endDate *time.Time, priority string, storyPoints int, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID) (*Task, error)
	AssignTask(taskID uuid.UUID, devID uint, assignerID uint) error
	SubmitWork(taskID uuid.UUID, devID uint, commitHash, diff string) (*Submission, error)
	GetTaskByID(taskID uuid.UUID) (*Task, error)
	GetTaskByIDOrCode(idOrCode string) (*Task, error) // idOrCode is UUID or task code (e.g. mims-hdmap-main-001)
	GetMyTasks(userID uint) ([]Task, error)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksByProjectID(projectID uuid.UUID) ([]Task, error)
	GetGanttData(projectID *uuid.UUID) (*GanttData, error) // All tasks + dependencies; if projectID set, filter by project
	GetPendingApprovals(userRole string) ([]Task, error)   // Approvals inbox for PM/CEO

	// Task Dependencies (Gantt links)
	AddDependency(predecessorID, successorID uuid.UUID, depType string) (*TaskDependency, error)
	RemoveDependency(id uuid.UUID) error

	// Task Management with Access Control
	UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string, parentID *uuid.UUID, startDate, endDate *time.Time, progress *int, priority string, storyPoints *int, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID, applyEpic bool, sortOrder *int) (*Task, error)
	UpdateTaskResourceURLs(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, resourceURLs datatypes.JSON) (*Task, error)
	EstimateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) (*Task, error) // AI estimate time and update task.ai_estimated_minutes
	GenerateProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (*AIGeneratedPlan, error) // AI generates epics, milestones, sprints, tasks
	ClearProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) error                     // Remove all tasks, sprints, milestones, epics (CEO/PM)
	ScheduleProjectWithAI(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (updatedCount int, err error) // Estimate + schedule existing tasks (CEO/PM)
	DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error

	// Appeal System
	SubmitAppeal(submissionID uuid.UUID, devID uint, reason string) (*Appeal, error)
	ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error

	// Time Negotiation
	NegotiateTime(taskID uuid.UUID, devID uint, minutes int, reason string) error

	// Human Quality Gate
	ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error

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
	AddComment(taskID uuid.UUID, userID uint, content string) (*TaskComment, error)
	GetComments(taskID uuid.UUID) ([]TaskComment, error)

	// Time Logging
	LogTime(taskID uuid.UUID, userID uint, minutes int, description string) (*TimeLog, error)
	GetTimeLogs(taskID uuid.UUID) ([]TimeLog, error)

	// Analytics
	GetProjectAnalytics(projectID uuid.UUID) (*ProjectAnalytics, error)

	// Bulk Operations
	BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error

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
	Images         []string       `json:"images"`         // base64 data URLs of embedded images
	SlideURL       string         `json:"slide_url"`
	Source         string         `json:"source"`
	SlideIndex     int            `json:"slide_index"`
	PresentationID string         `json:"presentation_id"`
	Comments       []SlideComment `json:"comments"`
}

// PreviewGoogleSlidesRequest is the payload for POST /sentinel/import/google-slides/preview
type PreviewGoogleSlidesRequest struct {
	PresentationURL string `json:"presentation_url" binding:"required"`
	APIKey         string `json:"api_key"`
}

// PreviewSlideItem is one slide in the preview list (index, title, and hidden flag)
type PreviewSlideItem struct {
	Index  int    `json:"index"`
	Title  string `json:"title"`
	Hidden bool   `json:"hidden"` // true when slide is hidden/skipped in presentation (PPTX show="0"); import should skip by default
}

// PreviewGoogleSlidesResult is the response for the preview endpoint
type PreviewGoogleSlidesResult struct {
	PresentationTitle string            `json:"presentation_title"`
	PresentationID     string            `json:"presentation_id"`      // extracted from URL for "already imported" check
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

// ImportGoogleSlidesRequest is the payload for POST /sentinel/import/google-slides
// Either sprint_id (import into sprint) or epic_id (import into backlog/epic) or neither (import into backlog unassigned).
type ImportGoogleSlidesRequest struct {
	PresentationURL string `json:"presentation_url" binding:"required"`
	SprintID        string `json:"sprint_id"`  // optional: when set, tasks go to this sprint
	EpicID         string `json:"epic_id"`    // optional: when set, tasks go to this epic (backlog)
	ProjectID      string `json:"project_id" binding:"required"`
	APIKey         string `json:"api_key"`
	Priority       string `json:"priority"`
	StoryPoints    int    `json:"story_points"`
	// SlideIndices: 1-based indices of slides to import. If empty or nil, import all.
	SlideIndices []int `json:"slide_indices"`
}

// ImportGoogleSlidesResult is the response for POST /sentinel/import/google-slides
type ImportGoogleSlidesResult struct {
	CreatedCount    int     `json:"created_count"`
	SlideCount      int     `json:"slide_count"`
	PresentationTitle string `json:"presentation_title"`
	Tasks           []*Task `json:"tasks"`
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
	Title           string `json:"title"`
	Description     string `json:"description"`
	Priority        string `json:"priority"`
	StoryPoints     int    `json:"story_points"`
	EpicIndex       *int   `json:"epic_index"`        // 0-based index into epics list
	SprintIndex     *int   `json:"sprint_index"`     // 0-based index into sprints list
	MilestoneIndex  *int   `json:"milestone_index"`  // 0-based index into milestones list
	StartDate       string `json:"start_date"`       // YYYY-MM-DD
	EndDate         string `json:"end_date"`         // YYYY-MM-DD
}

// TaskEstimateInput is a minimal task info sent to AI for batch estimate + order
type TaskEstimateInput struct {
	Index       int    `json:"task_index"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	StoryPoints int    `json:"story_points"`
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

	// AnalyzeAppeal วิเคราะห์ความน่าเชื่อถือของ Appeal เพื่อแนะนำ CEO/PM
	// Returns: recommendation (UPHOLD/OVERTURN), confidence (0-100), reasoning, error
	AnalyzeAppeal(diff string, originalFeedback string, appealReason string) (recommendation string, confidence int, reasoning string, err error)

	// AnalyzeTimeNegotiation วิเคราะห์คำขอเจรจาเวลาจากนักพัฒนา
	// Returns: recommendation (APPROVE/REJECT), confidence (0-100), reasoning, error
	AnalyzeTimeNegotiation(taskTitle, taskDescription string, aiEstimate, devProposal int, devReason string) (recommendation string, confidence int, reasoning string, err error)
}
