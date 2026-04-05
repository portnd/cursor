package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ─── Status & Environment constants ───────────────────────────────────────────

const (
	StatusPending   = "PENDING"
	StatusReviewing = "REVIEWING"
	StatusApproved  = "APPROVED"
	StatusRejected  = "REJECTED"
	StatusDeployed  = "DEPLOYED"

	EnvStaging    = "STAGING"
	EnvPreProd    = "PRE-PROD"
	EnvProduction = "PRODUCTION"
)

// ─── Domain errors ────────────────────────────────────────────────────────────

var (
	ErrNotFound         = errors.New("deployment request not found")
	ErrForbidden        = errors.New("insufficient permissions")
	ErrInvalidStatus    = errors.New("invalid status transition")
	ErrInvalidEnv       = errors.New("environment must be STAGING, PRE-PROD, or PRODUCTION")
	ErrMissingTitle     = errors.New("title is required")
	ErrMissingBranch    = errors.New("branch is required")
	ErrAlreadyReviewing = errors.New("request is already being reviewed")
)

// ─── Core entity ──────────────────────────────────────────────────────────────

// DeploymentRequest represents a request to deploy code to a target environment.
// Engineers submit requests; Chief Engineers review, approve, and mark deployed.
type DeploymentRequest struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string     `json:"title" gorm:"type:varchar(200);not null"`
	Description     string     `json:"description" gorm:"type:text"`
	Branch          string     `json:"branch" gorm:"type:varchar(300);not null"`
	PRUrl           string     `json:"pr_url" gorm:"column:pr_url;type:varchar(512)"`
	Environment     string     `json:"environment" gorm:"type:varchar(20);not null;default:'STAGING'"` // STAGING | PRE-PROD | PRODUCTION
	Status          string     `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`      // PENDING | REVIEWING | APPROVED | REJECTED | DEPLOYED
	RequesterID     uint       `json:"requester_id" gorm:"not null"`
	ReviewerID      *uint      `json:"reviewer_id,omitempty"`
	TaskID          *uuid.UUID `json:"task_id,omitempty" gorm:"type:uuid;index"` // optional link to a sentinel task (auto-advances task to READY_FOR_UAT on deploy)
	TaskRef         string     `json:"task_ref" gorm:"type:varchar(200)"`
	RejectionReason string     `json:"rejection_reason" gorm:"type:text"`
	ReviewNotes     string     `json:"review_notes" gorm:"type:text"`
	DeployedAt      *time.Time `json:"deployed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Enriched at query-time, not persisted
	RequesterEmail       string `json:"requester_email,omitempty" gorm:"-"`
	RequesterDisplayName string `json:"requester_display_name,omitempty" gorm:"-"`
	ReviewerEmail        string `json:"reviewer_email,omitempty" gorm:"-"`
	ReviewerDisplayName  string `json:"reviewer_display_name,omitempty" gorm:"-"`
}

func (DeploymentRequest) TableName() string { return "deployment_requests" }

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateDeploymentRequest struct {
	Title      string `json:"title" binding:"required"`
	Description string `json:"description"`
	Branch      string `json:"branch" binding:"required"`
	PRUrl       string `json:"pr_url"`
	Environment string `json:"environment"` // defaults to STAGING
	TaskID      string `json:"task_id"`     // optional UUID string linking to a sentinel task
	TaskRef     string `json:"task_ref"`
	ReviewerID  *uint  `json:"reviewer_id"` // optional pre-assign to a specific Chief Engineer
}

type ReviewActionRequest struct {
	Notes  string `json:"notes"`
	Reason string `json:"reason"` // only for reject
}

// ─── Stats ────────────────────────────────────────────────────────────────────

type DeploymentStats struct {
	TotalPending   int `json:"total_pending"`
	TotalReviewing int `json:"total_reviewing"`
	TotalApproved  int `json:"total_approved"`
	TotalDeployed  int `json:"total_deployed"`
	TotalRejected  int `json:"total_rejected"`
	DeployedToday  int `json:"deployed_today"`
}

// ─── TaskStatusAdvancer ───────────────────────────────────────────────────────

// TaskStatusAdvancer is a narrow interface allowing the deployment module to advance a linked
// sentinel task to READY_FOR_UAT when the Chief Engineer marks a deployment as deployed.
// Implemented by the sentinel usecase; injected at wire-up time in main.go.
type TaskStatusAdvancer interface {
	AdvanceTaskAfterDeploy(taskID uuid.UUID) error
}

// ─── Ports ────────────────────────────────────────────────────────────────────

type Repository interface {
	Create(req *DeploymentRequest) error
	GetByID(id uint) (*DeploymentRequest, error)
	GetByTaskID(taskID uuid.UUID) (*DeploymentRequest, error)
	List(requesterID *uint, status string) ([]DeploymentRequest, error)
	Update(req *DeploymentRequest) error
	GetStats() (*DeploymentStats, error)
}

type Usecase interface {
	CreateRequest(callerID uint, callerRole string, in CreateDeploymentRequest) (*DeploymentRequest, error)
	ListRequests(callerID uint, callerRole string, status string) ([]DeploymentRequest, error)
	GetRequest(callerID uint, callerRole string, id uint) (*DeploymentRequest, error)
	GetRequestByTaskID(callerID uint, callerRole string, taskIDStr string) (*DeploymentRequest, error)
	PickForReview(callerID uint, callerRole string, id uint) (*DeploymentRequest, error)
	Approve(callerID uint, callerRole string, id uint, in ReviewActionRequest) (*DeploymentRequest, error)
	Reject(callerID uint, callerRole string, id uint, in ReviewActionRequest) (*DeploymentRequest, error)
	MarkDeployed(callerID uint, callerRole string, id uint, in ReviewActionRequest) (*DeploymentRequest, error)
	GetStats(callerID uint, callerRole string) (*DeploymentStats, error)
}
