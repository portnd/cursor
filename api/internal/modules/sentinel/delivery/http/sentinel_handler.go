// SentinelHandler request DTOs and constructor; HTTP handlers live in sentinel_handler_*.go.
package http

import (
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type SentinelHandler struct {
	usecase          domain.SentinelUsecase
	googleAPIKey     string
	canvaAccessToken string
}

func NewSentinelHandler(usecase domain.SentinelUsecase, googleAPIKey, canvaAccessToken string) *SentinelHandler {
	return &SentinelHandler{usecase: usecase, googleAPIKey: googleAPIKey, canvaAccessToken: canvaAccessToken}
}

// Request DTOs
type createProjectReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"` // Optional: ACTIVE, COMPLETED, ON_HOLD (default ACTIVE)
}

type createTaskReq struct {
	Title            string  `json:"title" binding:"required"`
	Description      string  `json:"description"`
	TaskType         string  `json:"task_type"` // FEATURE, TASK, BUG (default: TASK)
	DueDate          *string `json:"due_date"`
	ProjectID        *string `json:"project_id"`
	ParentID         *string `json:"parent_id"`
	EpicID           *string `json:"epic_id"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
	Priority         string  `json:"priority"`
	StoryPoints      float64 `json:"story_points"`
	SprintID         *string `json:"sprint_id"`
	MilestoneID      *string `json:"milestone_id"`
	EstimatedMinutes *int    `json:"estimated_minutes"` // Manual estimate; stored for Costing Engine (mandatory from frontend)
}

type assignTaskReq struct {
	DevID uint `json:"dev_id"` // 0 = unassign
}

type submitWorkReq struct {
	ReferenceURL string `json:"reference_url" binding:"required"`
	Note         string `json:"note"`
}

type rejectTaskReq struct {
	Reason string `json:"reason" binding:"required,min=10"`
}

type submitAppealReq struct {
	Reason string `json:"reason" binding:"required"`
}

type resolveAppealReq struct {
	Status string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	Note   string `json:"note"`
}

type negotiateTimeReq struct {
	Minutes int    `json:"minutes" binding:"required,gt=0"`
	Reason  string `json:"reason" binding:"required,min=20"`
}

type updateTaskReq struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	TaskType         string   `json:"task_type"` // FEATURE, TASK, BUG
	ParentID         *string  `json:"parent_id"`
	EpicID           *string  `json:"epic_id"`
	SortOrder        *int     `json:"sort_order"`
	DueAt            *string  `json:"due_at"`
	StartDate        *string  `json:"start_date"`
	EndDate          *string  `json:"end_date"`
	Progress         *int     `json:"progress"`
	Priority         string   `json:"priority"`
	StoryPoints      *float64 `json:"story_points"`
	SprintID         *string  `json:"sprint_id"`
	MilestoneID      *string  `json:"milestone_id"`
	EstimatedMinutes *int     `json:"estimated_minutes"` // Manual estimate; feeds Costing Engine
}

type taskSummaryResponse struct {
	Summary  *domain.TaskSummary `json:"summary"`
	HasRichContent bool          `json:"has_rich_content"`
}

type taskDetailResponse struct {
	Task          *domain.Task `json:"task"`
	AttachmentCount int        `json:"attachment_count"`
	HasRichContent bool        `json:"has_rich_content"`
}

type createEpicReq struct {
	ProjectID   string  `json:"project_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type updateEpicReq struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Color       string  `json:"color"`
	SortOrder   *int    `json:"sort_order"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type createSprintReq struct {
	ProjectID string  `json:"project_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Goal      string  `json:"goal"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type updateSprintReq struct {
	Name      string  `json:"name"`
	Goal      string  `json:"goal"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	SortOrder *int    `json:"sort_order"`
}

type addTasksToSprintReq struct {
	TaskIDs []string `json:"task_ids" binding:"required"`
}

type createMilestoneReq struct {
	ProjectID   string  `json:"project_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
}

type updateMilestoneReq struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type addCommentReq struct {
	Content string `json:"content" binding:"required"`
}

type editCommentReq struct {
	Content string `json:"content" binding:"required"`
}

type logTimeReq struct {
	Minutes        int    `json:"minutes" binding:"required,gt=0"`
	Description    string `json:"description"`
	WorkType       string `json:"work_type"`
	LoggedDate     string `json:"logged_date"` // YYYY-MM-DD, optional (defaults to today)
	IsTimerSession bool   `json:"is_timer_session"`
}

type editTimeLogReq struct {
	Minutes     int     `json:"minutes" binding:"required,gt=0"`
	Description string  `json:"description"`
	WorkType    string  `json:"work_type"`
	TaskID      *string `json:"task_id"` // optional — change task assignment
}

type bulkStatusReq struct {
	TaskIDs []string `json:"task_ids" binding:"required"`
	Status  string   `json:"status" binding:"required"`
}

type createDependencyReq struct {
	PredecessorID string `json:"predecessor_id" binding:"required"` // Task that must happen first
	SuccessorID   string `json:"successor_id" binding:"required"`   // Task that waits
	Type          string `json:"type"`                              // FS, SS, FF, SF (default FS)
}

type updateConfigReq struct {
	ActiveModel      string  `json:"active_model" binding:"required"`
	Temperature      float32 `json:"temperature" binding:"required,gte=0,lte=1"`
	CursorAssistance int     `json:"cursor_assistance" binding:"required,gte=0,lte=100"`
}
