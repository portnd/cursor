// Task listing, Gantt, dependencies, and approvals inbox.
package http

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// CreateTask handles POST /api/v1/tasks
// Top-level tasks: CEO / Product Owner / Manager. Sub-tasks: same, or parent assignee/creator.
func (h *SentinelHandler) CreateTask(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	var req createTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid date format",
				"message": "due_date must be in ISO8601/RFC3339 format (e.g., '2026-01-30T15:00:00Z')",
			})
			return
		}
		dueDate = &parsedTime
	}

	// Parse optional project_id (task belongs to project)
	var projectID *uuid.UUID
	if req.ProjectID != nil && *req.ProjectID != "" {
		pid, err := uuid.Parse(*req.ProjectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"message": "project_id must be a valid UUID",
			})
			return
		}
		projectID = &pid
	}
	// Parse Gantt/WBS optional fields
	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"message": "parent_id must be a valid UUID",
			})
			return
		}
		parentID = &pid
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid date format",
				"message": "start_date must be ISO8601/RFC3339",
			})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid date format",
				"message": "end_date must be ISO8601/RFC3339",
			})
			return
		}
		endDate = &t
	}

	var sprintID *uuid.UUID
	if req.SprintID != nil && *req.SprintID != "" {
		sid, err := uuid.Parse(*req.SprintID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "sprint_id must be a valid UUID"})
			return
		}
		sprintID = &sid
	}
	var milestoneID *uuid.UUID
	if req.MilestoneID != nil && *req.MilestoneID != "" {
		mid, err := uuid.Parse(*req.MilestoneID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "milestone_id must be a valid UUID"})
			return
		}
		milestoneID = &mid
	}
	var epicID *uuid.UUID
	if req.EpicID != nil && *req.EpicID != "" {
		eid, err := uuid.Parse(*req.EpicID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "epic_id must be a valid UUID"})
			return
		}
		epicID = &eid
	}

	userRole := getUserRoleFromContext(c)
	if parentID == nil {
		if !isPrivilegedTaskRole(userRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO, Product Owner, or Manager can create top-level tasks"})
			return
		}
	} else {
		parent, err := h.usecase.GetTaskByID(*parentID)
		if err != nil || parent == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "parent task not found"})
			return
		}
		if !canUserCreateSubtask(parent, userID, userRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO, Product Owner, Manager, the parent task assignee, or creator can create sub-tasks"})
			return
		}
	}

	task, err := h.usecase.CreateTask(req.Title, req.Description, req.TaskType, userID, dueDate, projectID, parentID, startDate, endDate, req.Priority, req.StoryPoints, sprintID, milestoneID, epicID, req.EstimatedMinutes)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		log.Printf("[CreateTask] usecase error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create task",
			"message": err.Error(),
		})
		return
	}

	if req.AssignedTo != nil && *req.AssignedTo != 0 {
		if assignErr := h.usecase.AssignTask(task.ID, *req.AssignedTo, userID, getUserRoleFromContext(c)); assignErr != nil {
			log.Printf("[CreateTask] auto-assign warning: %v", assignErr)
		} else {
			task.AssignedTo = req.AssignedTo
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    task,
	})
}

// AssignTask handles POST /api/v1/tasks/:id/assign
// :id can be UUID or task code. The requesting user (Product Owner/CEO) is recorded as assigned_by for leaderboard scope.
func (h *SentinelHandler) AssignTask(c *gin.Context) {
	assignerID := getUserIDFromContext(c)
	if assignerID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	var req assignTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.usecase.AssignTask(task.ID, req.DevID, assignerID, getUserRoleFromContext(c)); err != nil {
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to assign task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task assigned successfully",
	})
}

// SubmitWork handles POST /api/v1/tasks/:id/submit
// :id can be UUID or task code
func (h *SentinelHandler) SubmitWork(c *gin.Context) {
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

	var req submitWorkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	sub, err := h.usecase.SubmitWork(task.ID, userID, req.ReferenceURL, req.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit work",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Work handed over for review",
		"data":    sub,
	})
}

// SubmitUAT stores UAT payload on a FEATURE task and promotes it to REVIEW_PENDING for Product Owner/CEO review.
func (h *SentinelHandler) SubmitUAT(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	var req domain.UATPayloadData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	if err := h.usecase.SubmitUAT(task.ID, userID, req); err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit UAT", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "UAT submitted — feature is now pending Product Owner/CEO review"})
}

// GetTaskSummary handles GET /api/v1/sentinel/tasks/:id/summary
// Returns a lightweight shape for the task detail header and sidebar.
func (h *SentinelHandler) GetTaskSummary(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "task id or code is required"})
		return
	}
	task, err := h.usecase.GetTaskByIDOrCode(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Task not found"})
		return
	}
	summary := domain.TaskSummary{
		ID: task.ID, Code: task.Code, Title: task.Title, ProjectID: task.ProjectID, EpicID: task.EpicID, SprintID: task.SprintID, PreviousSprintID: task.PreviousSprintID,
		MilestoneID: task.MilestoneID, TaskType: task.TaskType, Priority: task.Priority, StoryPoints: task.StoryPoints,
		EstimatedMinutes: task.EstimatedMinutes, ParentID: task.ParentID, SortOrder: task.SortOrder,
		StartDate: task.StartDate, EndDate: task.EndDate, Progress: task.Progress,
		DueAt: task.DueAt, StartedAt: task.StartedAt, CompletedAt: task.CompletedAt,
		Status: task.Status, NegotiationStatus: task.NegotiationStatus,
		AssignedTo: task.AssignedTo, AssignedToDisplayName: task.AssignedToDisplayName,
		AssignedToEmail: task.AssignedToEmail, AssignedToAvatarURL: task.AssignedToAvatarURL,
		IsKomgrip: task.IsKomgrip, CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
	}
	hasRichContent := task.Description != "" || len(task.ResourceURLs) > 2 || len(task.Submissions) > 0
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{"message": "Task summary retrieved successfully", "data": gin.H{"summary": summary, "has_rich_content": hasRichContent}})
}

// GetTaskDetail handles GET /api/v1/sentinel/tasks/:id/detail
// Returns the full task payload including rich description, attachments and images.
func (h *SentinelHandler) GetTaskDetail(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "task id or code is required"})
		return
	}
	task, err := h.usecase.GetTaskByIDOrCode(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Task not found"})
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{"message": "Task detail retrieved successfully", "data": gin.H{"task": task, "attachment_count": len(task.Submissions), "has_rich_content": task.Description != "" || len(task.ResourceURLs) > 2 || len(task.Submissions) > 0}})
}


// GetTaskActivity handles GET /api/v1/sentinel/tasks/:id/activity
func (h *SentinelHandler) GetTaskActivity(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "task id or code is required" ||
			strings.Contains(errMsg, "task not found") ||
			strings.Contains(errMsg, "record not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Task not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": errMsg})
		return
	}
	items, err := h.usecase.GetTaskActivityTimeline(task.ID)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load activity", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Activity retrieved",
		"data":    items,
	})
}

// GetMyTasks handles GET /api/v1/tasks/my
// Retrieves all tasks assigned to the authenticated user
func (h *SentinelHandler) GetMyTasks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	tasks, err := h.usecase.GetMyTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve tasks",
			"message": err.Error(),
		})
		return
	}

	activeSprints, err := h.usecase.GetMyActiveSprints(userID)
	if err != nil {
		activeSprints = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Tasks retrieved successfully",
		"data":           tasks,
		"active_sprints": activeSprints,
	})
}

// GetGlobalActiveTasks handles GET /api/v1/sentinel/tasks/my-global-active
// TASK/BUG in ACTIVE sprints: CEO/MANAGER = company-wide; others team-scoped; teams off → Product Owner / engineer assignment rules.
func (h *SentinelHandler) GetGlobalActiveTasks(c *gin.Context) {
	if getUserIDFromContext(c) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	tasks, err := h.usecase.GetGlobalActiveTasks(callerCtx(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve global active tasks",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Global active tasks retrieved successfully",
		"data":    tasks,
	})
}

// GetUnassignedTasks handles GET /api/v1/sentinel/tasks/unassigned
// Retrieves all tasks that are not assigned to anyone (for CEO / Product Owner to assign)
func (h *SentinelHandler) GetUnassignedTasks(c *gin.Context) {
	tasks, err := h.usecase.GetUnassignedTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve unassigned tasks",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unassigned tasks retrieved successfully",
		"data":    tasks,
	})
}

// GetTeamActiveTasks handles GET /api/v1/sentinel/tasks/team-active
// Returns all tasks in ACTIVE sprints within the caller's team.
// Used by the "Quick Log Time" global modal to let devs log time against any teammate's task.
// CEO/MANAGER bypass team restriction and see all active-sprint tasks.
func (h *SentinelHandler) GetTeamActiveTasks(c *gin.Context) {
	ctx := callerCtx(c)
	tasks, err := h.usecase.GetTeamActiveTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve team active tasks",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Team active tasks retrieved successfully",
		"data":    tasks,
	})
}

// GetActiveFeatures handles GET /api/v1/sentinel/tasks/features
// Returns all FEATURE-type tasks for the Product Owner/CEO Feature Roadmap Board.
// Each feature includes a roll-up progress (0-100%) and its child TASK/BUG items for the accordion.
// Product Owner is scoped to their team; CEO/MANAGER see all teams.
func (h *SentinelHandler) GetActiveFeatures(c *gin.Context) {
	ctx := callerCtx(c)

	var projectID *uuid.UUID
	if idStr := strings.TrimSpace(c.Query("project_id")); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"message": "project_id must be a valid UUID",
			})
			return
		}
		projectID = &id
	}

	items, err := h.usecase.GetActiveFeatures(ctx.TeamID, ctx.Role, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve feature roadmap",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Feature roadmap retrieved successfully",
		"data":    items,
	})
}

// GetAllTasks handles GET /api/v1/sentinel/tasks
// Optional query: ?project_id=UUID to get only that project's tasks (for project Board/Backlog).
// Optional query: ?task_type=FEATURE|TASK|BUG to filter by task typology.
// Without project_id: returns ALL tasks (ADMIN / Product Owner overview).
func (h *SentinelHandler) GetAllTasks(c *gin.Context) {
	var tasks []domain.Task
	var err error
	var message string

	if idsCSV := strings.TrimSpace(c.Query("project_ids")); idsCSV != "" {
		parts := strings.Split(idsCSV, ",")
		projectIDs := make([]uuid.UUID, 0, len(parts))
		seen := make(map[uuid.UUID]struct{}, len(parts))
		for _, part := range parts {
			idStr := strings.TrimSpace(part)
			if idStr == "" {
				continue
			}
			id, parseErr := uuid.Parse(idStr)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_ids must be comma-separated UUIDs"})
				return
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			projectIDs = append(projectIDs, id)
		}
		tasks, err = h.usecase.GetTasksByProjectIDs(projectIDs)
		message = "Project-scoped tasks retrieved successfully"
	} else if idStr := c.Query("project_id"); idStr != "" {
		projectID, parseErr := uuid.Parse(idStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
			return
		}
		tasks, err = h.usecase.GetTasksByProjectID(projectID)
		message = "Project tasks retrieved successfully"
	} else {
		tasks, err = h.usecase.GetAllTasks()
		message = "All tasks retrieved successfully"
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve tasks",
			"message": err.Error(),
		})
		return
	}
	// Optional client-side filter by task_type (avoids breaking the repository interface)
	if taskTypeFilter := c.Query("task_type"); taskTypeFilter != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if t.TaskType == taskTypeFilter {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	if sprintIDStr := strings.TrimSpace(c.Query("sprint_id")); sprintIDStr != "" {
		sprintID, parseErr := uuid.Parse(sprintIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "sprint_id must be a valid UUID"})
			return
		}
		filtered := tasks[:0]
		for _, t := range tasks {
			if t.SprintID != nil && *t.SprintID == sprintID {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": tasks})
}

// GetGantt handles GET /api/v1/sentinel/tasks/gantt
// Returns tasks and dependencies for Gantt chart. Optional query: ?project_id=xxx to filter by project.
func (h *SentinelHandler) GetGantt(c *gin.Context) {
	var projectID *uuid.UUID
	if idStr := c.Query("project_id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"message": "project_id must be a valid UUID",
			})
			return
		}
		projectID = &id
	}
	data, err := h.usecase.GetGanttData(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve Gantt data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gantt data retrieved successfully",
		"data":    data,
	})
}

// CreateDependency handles POST /api/v1/sentinel/tasks/dependencies
// Creates a link between tasks (e.g. successor cannot start until predecessor finishes)
func (h *SentinelHandler) CreateDependency(c *gin.Context) {
	var req createDependencyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	predID, err := uuid.Parse(req.PredecessorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "predecessor_id must be a valid UUID",
		})
		return
	}
	succID, err := uuid.Parse(req.SuccessorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "successor_id must be a valid UUID",
		})
		return
	}
	dep, err := h.usecase.AddDependency(predID, succID, req.Type)
	if err != nil {
		if contains(err.Error(), "self-linking") || contains(err.Error(), "predecessor and successor") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if contains(err.Error(), "invalid dependency type") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create dependency",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Dependency created successfully",
		"data":    dep,
	})
}

// DeleteDependency handles DELETE /api/v1/sentinel/tasks/dependencies/:id
func (h *SentinelHandler) DeleteDependency(c *gin.Context) {
	idStr := c.Param("id")
	depID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid dependency UUID",
		})
		return
	}
	if err := h.usecase.RemoveDependency(depID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete dependency",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Dependency deleted successfully",
	})
}

// GetApprovals handles GET /api/v1/sentinel/tasks/approvals
// Returns tasks requiring Product Owner/CEO/Manager attention (PENDING appeals or time negotiations)
// Access: CEO, Manager, and Product Owner
func (h *SentinelHandler) GetApprovals(c *gin.Context) {
	// 1️⃣ Extract user role from JWT context
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found in token",
		})
		return
	}

	// 2️⃣ Call usecase with role validation
	tasks, err := h.usecase.GetPendingApprovals(userRole)
	if err != nil {
		// Check if it's an authorization error
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve approvals inbox",
			"message": err.Error(),
		})
		return
	}

	// 3️⃣ Return the approvals inbox
	c.JSON(http.StatusOK, gin.H{
		"message": "Approvals inbox retrieved successfully",
		"data":    tasks,
		"count":   len(tasks),
	})
}
