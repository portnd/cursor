package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Task activity audit — immutable events for task lifecycle (compliance & UX timeline).
const (
	TaskActivityCreated            = "CREATED"
	TaskActivityAssigned           = "ASSIGNED"
	TaskActivityUnassigned         = "UNASSIGNED"
	TaskActivityStatusChanged      = "STATUS_CHANGED"
	TaskActivitySubmittedReview    = "SUBMITTED_REVIEW"
	TaskActivityApprovedReview     = "APPROVED_REVIEW"
	TaskActivityRejectedReview     = "REJECTED_REVIEW"
	TaskActivityReadyForTest       = "READY_FOR_TEST"
	TaskActivityPMApprovedDeploy   = "PM_APPROVED_TEST"
	TaskActivityDeployed           = "DEPLOYED"
	TaskActivityCEOFinalApproved   = "CEO_FINAL_APPROVED"
	TaskActivityWorkflowReject     = "WORKFLOW_REJECT"
	TaskActivitySubmitUAT          = "SUBMIT_UAT"
	TaskActivityNegotiation        = "NEGOTIATION_SUBMITTED"
	TaskActivityAppealComplete     = "APPEAL_APPROVED_COMPLETE"
	TaskActivityParentRollupStatus = "PARENT_ROLLUP_STATUS"
)

// TaskActivityEvent is persisted for each task lifecycle action.
type TaskActivityEvent struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID    uuid.UUID      `json:"task_id" gorm:"type:uuid;not null;index"`
	Action    string         `json:"action" gorm:"type:varchar(64);not null;index"`
	ActorID   *uint          `json:"actor_id" gorm:"index"`
	Payload   datatypes.JSON `json:"payload" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

func (TaskActivityEvent) TableName() string { return "task_activity_events" }

// TaskActivityItem is an API row with optional actor enrichment and inference flag.
type TaskActivityItem struct {
	ID               string          `json:"id"`
	Action           string          `json:"action"`
	At               time.Time       `json:"at"`
	ActorUserID      *uint           `json:"actor_user_id,omitempty"`
	ActorEmail       string          `json:"actor_email,omitempty"`
	ActorDisplayName string          `json:"actor_display_name,omitempty"`
	Payload          json.RawMessage `json:"payload,omitempty"`
	Inferred         bool            `json:"inferred"`
}
