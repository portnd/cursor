// Appeals, negotiation, task updates, estimates, lifecycle (approve / reject / UAT lanes).
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)

// SubmitAppeal handles POST /api/v1/submissions/:id/appeal
// Allows a developer to appeal an AI FAIL verdict
func (h *SentinelHandler) SubmitAppeal(c *gin.Context) {
	submissionIDStr := c.Param("id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid submission UUID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	var req submitAppealReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	appeal, err := h.usecase.SubmitAppeal(submissionID, userID, req.Reason)
	if err != nil {
		// Check for specific error types
		if err.Error() == "unauthorized: only the developer who submitted can appeal" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "appeal already exists for this submission" ||
			err.Error() == "can only appeal FAIL verdicts" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit appeal",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appeal submitted successfully",
		"data":    appeal,
	})
}

// ResolveAppeal handles POST /api/v1/appeals/:id/resolve
// Allows Product Owner/CEO to approve or reject an appeal
func (h *SentinelHandler) ResolveAppeal(c *gin.Context) {
	appealIDStr := c.Param("id")
	appealID, err := uuid.Parse(appealIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid appeal UUID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	var req resolveAppealReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.usecase.ResolveAppeal(appealID, userID, req.Status, req.Note); err != nil {
		// Check for specific error types
		errMsg := err.Error()

		// Forbidden: Role-based access control
		if errMsg == "status must be APPROVED or REJECTED" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		// Check if error contains "forbidden" or "role"
		if contains(errMsg, "forbidden") || contains(errMsg, "only CEO or PM") || contains(errMsg, "only CEO or Product Owner") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Check if error contains "not found"
		if contains(errMsg, "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}

		// Default: Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to resolve appeal",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appeal resolved successfully",
	})
}

// NegotiateTime handles POST /api/v1/sentinel/tasks/:id/negotiate
// :id can be UUID or task code
func (h *SentinelHandler) NegotiateTime(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	var req negotiateTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.usecase.NegotiateTime(task.ID, userID, req.Minutes, req.Reason); err != nil {
		// Check for specific error types
		if err.Error() == "unauthorized: only the assigned developer can negotiate time" ||
			err.Error() == "unauthorized: only the task creator can negotiate time for unassigned tasks" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		if err.Error() == "time negotiation already pending review" ||
			strings.Contains(err.Error(), "proposed time must be greater") ||
			err.Error() == "negotiation reason must be at least 20 characters" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit time negotiation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Time negotiation submitted successfully and pending Product Owner/CEO review",
	})
}

// UpdateTask handles PATCH /api/v1/sentinel/tasks/:id
// :id can be UUID or task code
func (h *SentinelHandler) UpdateTask(c *gin.Context) {
	taskResolved, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	// Get requesting user info from context
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found",
		})
		return
	}

	// 3. Parse request body
	var req updateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// 4. Validate: At least one field must be provided
	hasTitle := req.Title != ""
	hasDesc := req.Description != ""
	hasTaskType := req.TaskType != ""
	hasParent := req.ParentID != nil && *req.ParentID != ""
	hasEpicKey := req.EpicID != nil
	hasDueAt := req.DueAt != nil && *req.DueAt != ""
	hasStart := req.StartDate != nil && *req.StartDate != ""
	hasEnd := req.EndDate != nil && *req.EndDate != ""
	hasProgress := req.Progress != nil
	hasPriority := req.Priority != ""
	hasSP := req.StoryPoints != nil
	hasSprint := req.SprintID != nil
	hasMilestone := req.MilestoneID != nil
	hasSortOrder := req.SortOrder != nil
	hasEstMins := req.EstimatedMinutes != nil
	if !hasTitle && !hasDesc && !hasTaskType && !hasParent && !hasEpicKey && !hasDueAt && !hasStart && !hasEnd && !hasProgress && !hasPriority && !hasSP && !hasSprint && !hasMilestone && !hasSortOrder && !hasEstMins {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "At least one field must be provided",
		})
		return
	}

	var parentID *uuid.UUID
	if hasParent {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "parent_id must be a valid UUID"})
			return
		}
		parentID = &pid
	}
	var epicIDUpd *uuid.UUID
	if hasEpicKey {
		if *req.EpicID == "" {
			epicIDUpd = nil
		} else {
			eid, err := uuid.Parse(*req.EpicID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "epic_id must be a valid UUID"})
				return
			}
			epicIDUpd = &eid
		}
	}
	var dueAt, startDate, endDate *time.Time
	if hasDueAt {
		t, err := time.Parse(time.RFC3339, *req.DueAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "due_at must be ISO8601/RFC3339"})
			return
		}
		dueAt = &t
	}
	if hasStart {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "start_date must be ISO8601/RFC3339"})
			return
		}
		startDate = &t
	}
	if hasEnd {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "end_date must be ISO8601/RFC3339"})
			return
		}
		endDate = &t
	}
	var progress *int
	if hasProgress {
		p := *req.Progress
		if p < 0 || p > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "progress must be between 0 and 100"})
			return
		}
		progress = &p
	}
	var sprintIDUpd *uuid.UUID
	if hasSprint {
		if *req.SprintID != "" {
			sid, err := uuid.Parse(*req.SprintID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "sprint_id must be a valid UUID"})
				return
			}
			sprintIDUpd = &sid
		}
		// empty string = unassign from sprint (sprintIDUpd stays nil)
	}
	var milestoneIDUpd *uuid.UUID
	if hasMilestone {
		if *req.MilestoneID != "" {
			mid, err := uuid.Parse(*req.MilestoneID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "milestone_id must be a valid UUID"})
				return
			}
			milestoneIDUpd = &mid
		}
		// empty string = unassign from milestone (milestoneIDUpd stays nil)
	}
	var sortOrderUpd *int
	if hasSortOrder {
		sortOrderUpd = req.SortOrder
	}

	task, err := h.usecase.UpdateTask(taskResolved.ID, userID, userRole, req.Title, req.Description, req.TaskType, parentID, dueAt, startDate, endDate, progress, req.Priority, req.StoryPoints, sprintIDUpd, hasSprint, milestoneIDUpd, epicIDUpd, hasEpicKey, sortOrderUpd, req.EstimatedMinutes)
	if err != nil {
		// Check for authorization error
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Check for not found error
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"data":    task,
	})
}

// UpdateTaskSlideResources handles PATCH /api/v1/sentinel/tasks/:id/slide-resources (resource_urls for slide images/annotations)
func (h *SentinelHandler) UpdateTaskSlideResources(c *gin.Context) {
	taskResolved, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	var req struct {
		ResourceURLs json.RawMessage `json:"resource_urls" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	task, err := h.usecase.UpdateTaskResourceURLs(taskResolved.ID, userID, userRole, datatypes.JSON(req.ResourceURLs))
	if err != nil {
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Slide resources updated", "data": task})
}

// EstimateTask handles POST /api/v1/sentinel/tasks/:id/estimate — AI estimates time and updates task.estimated_minutes (used by ScheduleProjectWithAI)
func (h *SentinelHandler) EstimateTask(c *gin.Context) {
	taskResolved, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	updated, err := h.usecase.EstimateTask(taskResolved.ID, userID, userRole)
	if err != nil {
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Estimate failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "AI estimate updated", "data": updated})
}

// DeleteTask handles DELETE /api/v1/sentinel/tasks/:id
// :id can be UUID or task code
func (h *SentinelHandler) DeleteTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found",
		})
		return
	}

	if err := h.usecase.DeleteTask(task.ID, userID, userRole); err != nil {
		// Check for authorization error
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Check for not found error
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}

		// Task has sub-tasks: must delete sub-tasks first
		if contains(err.Error(), "task_has_sub_tasks") || contains(err.Error(), "violates foreign key constraint \"fk_tasks_sub_tasks\"") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Cannot delete task",
				"message": "มี sub task ไม่สามารถลบได้ หากต้องการลบ ต้องลบ sub task ก่อน",
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// ApproveTask marks a task as COMPLETED after human verification (Product Owner/CEO only)
// POST /api/v1/sentinel/tasks/:id/approve ; :id can be UUID or task code
func (h *SentinelHandler) ApproveTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	approverID := getUserIDFromContext(c)
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	approverRole := getUserRoleFromContext(c)
	if approverRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found",
		})
		return
	}

	if err := h.usecase.ApproveTask(task.ID, approverID, approverRole); err != nil {
		// Check for authorization error
		if contains(err.Error(), "access denied") || contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Check for not found error
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}

		// Check for invalid status error
		if contains(err.Error(), "not pending review") || contains(err.Error(), "current status") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid Status",
				"message": err.Error(),
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to approve task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task approved and marked as COMPLETED",
	})
}

// RejectTask returns a task to IN_PROGRESS with a rejection reason comment (Product Owner/CEO/MANAGER only)
// POST /api/v1/sentinel/tasks/:id/reject ; :id can be UUID or task code
func (h *SentinelHandler) RejectTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	rejectorID := getUserIDFromContext(c)
	if rejectorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	rejectorRole := getUserRoleFromContext(c)
	if rejectorRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found",
		})
		return
	}

	var req rejectTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.usecase.RejectTask(task.ID, rejectorID, rejectorRole, req.Reason); err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not pending review") || contains(err.Error(), "current status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Status", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reject task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task rejected and returned to IN_PROGRESS",
	})
}

// MarkReadyForTest moves a TASK/BUG from IN_PROGRESS to READY_FOR_TEST (Dev action).
// POST /api/v1/sentinel/tasks/:id/ready-for-test
func (h *SentinelHandler) MarkReadyForTest(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	devID := getUserIDFromContext(c)
	if devID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	if err := h.usecase.MarkReadyForTest(task.ID, devID); err != nil {
		var badReq *domain.ErrBadRequest
		if errors.As(err, &badReq) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as READY_FOR_TEST"})
}

// PMApproveSubTask is the Product Owner's first-stage approval: READY_FOR_TEST → READY_FOR_UAT.
// The Product Owner must supply a test URL and detailed test steps for the CEO to follow.
// POST /api/v1/sentinel/tasks/:id/pm-approve-sub
func (h *SentinelHandler) PMApproveSubTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	pmID := getUserIDFromContext(c)
	if pmID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	pmRole := getUserRoleFromContext(c)

	var req struct {
		TestURL   string `json:"test_url"`
		TestSteps string `json:"test_steps"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid request body"})
		return
	}

	if err := h.usecase.PMApproveSubTask(task.ID, pmID, pmRole, req.TestURL, req.TestSteps); err != nil {
		var badReq *domain.ErrBadRequest
		if errors.As(err, &badReq) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": badReq.Msg})
			return
		}
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		if contains(err.Error(), "READY_FOR_TEST") || contains(err.Error(), "current status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Status", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit for CEO approval", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Test evidence submitted — task forwarded to CEO for final approval"})
}

// ApproveSubTask is the CEO's final approval: READY_FOR_UAT → COMPLETED.
// POST /api/v1/sentinel/tasks/:id/approve-sub
func (h *SentinelHandler) ApproveSubTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	ceoID := getUserIDFromContext(c)
	if ceoID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	ceoRole := getUserRoleFromContext(c)
	if ceoRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	if err := h.usecase.ApproveSubTask(task.ID, ceoID, ceoRole); err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		if contains(err.Error(), "READY_FOR_UAT") || contains(err.Error(), "current status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Status", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve sub-task", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sub-task approved and marked as COMPLETED"})
}

// RejectSubTask moves a READY_FOR_TEST task back to IN_PROGRESS and logs the reason (Product Owner/CEO action).
// POST /api/v1/sentinel/tasks/:id/reject-sub
func (h *SentinelHandler) RejectSubTask(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	pmID := getUserIDFromContext(c)
	if pmID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	pmRole := getUserRoleFromContext(c)
	if pmRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	var req rejectTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	if err := h.usecase.RejectSubTask(task.ID, pmID, pmRole, req.Reason); err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		var badReq *domain.ErrBadRequest
		if errors.As(err, &badReq) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		if contains(err.Error(), "READY_FOR_TEST") || contains(err.Error(), "current status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Status", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject sub-task", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sub-task rejected and returned to IN_PROGRESS"})
}

// GetTasksReadyForTest returns all TASK/BUG items in READY_FOR_TEST status for the caller's team.
// GET /api/v1/sentinel/tasks/ready-for-test
func (h *SentinelHandler) GetTasksReadyForTest(c *gin.Context) {
	ctx := callerCtx(c)
	if ctx.Role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	tasks, err := h.usecase.GetTasksReadyForTest(ctx.TeamID, ctx.Role)
	if err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test queue", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTasksReadyForCEOApproval returns TASK/BUG in READY_FOR_UAT status awaiting CEO final approval.
// GET /api/v1/sentinel/tasks/ceo-approval-queue
func (h *SentinelHandler) GetTasksReadyForCEOApproval(c *gin.Context) {
	ctx := callerCtx(c)
	if ctx.Role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	tasks, err := h.usecase.GetTasksReadyForCEOApproval(ctx.TeamID, ctx.Role)
	if err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CEO approval queue", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}
