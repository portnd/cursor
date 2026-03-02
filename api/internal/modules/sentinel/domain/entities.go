package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Task represents a work assignment
type Task struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title              string         `json:"title" gorm:"not null"`
	Description        string         `json:"description"`
	ResourceURLs       datatypes.JSON `json:"resource_urls" gorm:"type:jsonb;default:'{}'"`
	AIEstimatedMinutes int            `json:"ai_estimated_minutes"`

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
	AssignedTo *uint `json:"assigned_to"`
	CreatedBy  *uint `json:"created_by"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationship: Submissions history (loaded via Preload)
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:TaskID"`
}

// TableName overrides default
func (Task) TableName() string { return "tasks" }

// Submission represents a code push for review
type Submission struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
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
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
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
	CreateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	GetTasksByAssignee(userID uint) ([]Task, error)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetTasksRequiringApproval() ([]Task, error) // Tasks with PENDING appeals or time negotiations
	UpdateTask(task *Task) error
	DeleteTask(id uuid.UUID) error // Delete a task by ID

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

	// System Configuration (Singleton)
	GetSystemConfig() (*SystemConfig, error)
	UpdateSystemConfig(config *SystemConfig) error
}

// SentinelUsecase defines the business logic
type SentinelUsecase interface {
	CreateTask(title, desc string, creatorID uint, dueDate *time.Time) (*Task, error)
	AssignTask(taskID uuid.UUID, devID uint) error
	SubmitWork(taskID uuid.UUID, devID uint, commitHash, diff string) (*Submission, error)
	GetTaskByID(taskID uuid.UUID) (*Task, error)
	GetMyTasks(userID uint) ([]Task, error)
	GetUnassignedTasks() ([]Task, error)
	GetAllTasks() ([]Task, error)
	GetPendingApprovals(userRole string) ([]Task, error) // Approvals inbox for PM/CEO
	
	// Task Management with Access Control
	UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string) (*Task, error)
	DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error

	// Appeal System
	SubmitAppeal(submissionID uuid.UUID, devID uint, reason string) (*Appeal, error)
	ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error

	// Time Negotiation
	NegotiateTime(taskID uuid.UUID, devID uint, minutes int, reason string) error

	// Human Quality Gate
	ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error

	// System Configuration Management
	GetSystemConfig() (*SystemConfig, error)
	UpdateSystemConfig(activeModel string, temperature float32, cursorAssistance int, userRole string) (*SystemConfig, error)
	GetAvailableModels() []string
}

// AIService defines the interface for AI operations (Port)
type AIService interface {
	// EstimateEffort รับ Title/Desc แล้วคืนค่าเป็น นาที (minutes) และเหตุผล
	EstimateEffort(title, description string) (minutes int, reasoning string, err error)

	// ReviewCode วิเคราะห์ code diff และคืนค่า verdict (PASS/FAIL), score (0-100), feedback
	ReviewCode(diff string) (verdict string, score int, feedback string, err error)

	// AnalyzeAppeal วิเคราะห์ความน่าเชื่อถือของ Appeal เพื่อแนะนำ CEO/PM
	// Returns: recommendation (UPHOLD/OVERTURN), confidence (0-100), reasoning, error
	AnalyzeAppeal(diff string, originalFeedback string, appealReason string) (recommendation string, confidence int, reasoning string, err error)

	// AnalyzeTimeNegotiation วิเคราะห์คำขอเจรจาเวลาจากนักพัฒนา
	// Returns: recommendation (APPROVE/REJECT), confidence (0-100), reasoning, error
	AnalyzeTimeNegotiation(taskTitle, taskDescription string, aiEstimate, devProposal int, devReason string) (recommendation string, confidence int, reasoning string, err error)
}
