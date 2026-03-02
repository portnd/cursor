package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type SentinelHandler struct {
	usecase domain.SentinelUsecase
}

func NewSentinelHandler(usecase domain.SentinelUsecase) *SentinelHandler {
	return &SentinelHandler{usecase: usecase}
}

// Request DTOs
type createTaskReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"` // Optional: ISO8601/RFC3339 format (e.g., "2026-01-30T15:00:00Z")
}

type assignTaskReq struct {
	DevID uint `json:"dev_id" binding:"required"`
}

type submitWorkReq struct {
	CommitHash string `json:"commit_hash" binding:"required"`
	Diff       string `json:"diff"` // Optional: Code diff for AI review
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
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateConfigReq struct {
	ActiveModel      string  `json:"active_model" binding:"required"`
	Temperature      float32 `json:"temperature" binding:"required,gte=0,lte=1"`
	CursorAssistance int     `json:"cursor_assistance" binding:"required,gte=0,lte=100"`
}

// --- Handlers ---

// CreateTask handles POST /api/v1/tasks
// Creates a new task (CEO/PM only)
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

	task, err := h.usecase.CreateTask(req.Title, req.Description, userID, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    task,
	})
}

// AssignTask handles POST /api/v1/tasks/:id/assign
// Assigns a task to a developer (PM only)
func (h *SentinelHandler) AssignTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
		})
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

	if err := h.usecase.AssignTask(taskID, req.DevID); err != nil {
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
// Submits code for a task (Developer only)
func (h *SentinelHandler) SubmitWork(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
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

	var req submitWorkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	sub, err := h.usecase.SubmitWork(taskID, userID, req.CommitHash, req.Diff)
	if err != nil {
		// Check for duplicate commit hash error (PostgreSQL unique constraint violation)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") || 
		   strings.Contains(err.Error(), "idx_submissions_task_commit") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Duplicate submission",
				"message": "This commit hash has already been submitted for this task. Please use a different commit hash or check your previous submissions.",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit work",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Work submitted successfully",
		"data":    sub,
	})
}

// GetTaskByID handles GET /api/v1/sentinel/tasks/:id
// Retrieves a single task with full submission history
func (h *SentinelHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
		})
		return
	}

	task, err := h.usecase.GetTaskByID(taskID)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "task not found" || err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task retrieved successfully",
		"data":    task,
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks retrieved successfully",
		"data":    tasks,
	})
}

// GetUnassignedTasks handles GET /api/v1/sentinel/tasks/unassigned
// Retrieves all tasks that are not assigned to anyone (for CEO/PM to assign)
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

// GetAllTasks handles GET /api/v1/sentinel/tasks
// Retrieves ALL tasks in the system (for ADMIN/PM overview)
func (h *SentinelHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.usecase.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve all tasks",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All tasks retrieved successfully",
		"data":    tasks,
	})
}

// GetApprovals handles GET /api/v1/sentinel/tasks/approvals
// Returns tasks requiring PM/CEO attention (PENDING appeals or time negotiations)
// Access: CEO and PM only
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
// Allows PM/CEO to approve or reject an appeal
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
		if contains(errMsg, "forbidden") || contains(errMsg, "only CEO or PM") {
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
// Allows a developer to dispute/negotiate the AI-estimated time
func (h *SentinelHandler) NegotiateTime(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
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

	var req negotiateTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := h.usecase.NegotiateTime(taskID, userID, req.Minutes, req.Reason); err != nil {
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
			err.Error() == "proposed time must be greater than AI estimate (why negotiate if you need less time?)" ||
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
		"message": "Time negotiation submitted successfully and pending PM/CEO review",
	})
}

// UpdateTask handles PATCH /api/v1/sentinel/tasks/:id
// Updates task title/description with access control and AI re-estimation
// Only Creator or CEO can update
func (h *SentinelHandler) UpdateTask(c *gin.Context) {
	// 1. Parse Task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
		})
		return
	}

	// 2. Get requesting user info from context
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
	if req.Title == "" && req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "At least one field (title or description) must be provided",
		})
		return
	}

	// 5. Call usecase
	task, err := h.usecase.UpdateTask(taskID, userID, userRole, req.Title, req.Description)
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

// DeleteTask handles DELETE /api/v1/sentinel/tasks/:id
// Deletes a task with access control
// Only Creator or CEO can delete
func (h *SentinelHandler) DeleteTask(c *gin.Context) {
	// 1. Parse Task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
		})
		return
	}

	// 2. Get requesting user info from context
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

	// 3. Call usecase
	if err := h.usecase.DeleteTask(taskID, userID, userRole); err != nil {
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
			"error":   "Failed to delete task",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// ApproveTask marks a task as COMPLETED after human verification (PM/CEO only)
// POST /api/v1/sentinel/tasks/:id/approve
func (h *SentinelHandler) ApproveTask(c *gin.Context) {
	// 1. Parse Task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Invalid task UUID",
		})
		return
	}

	// 2. Get approver info from context
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

	// 3. Call usecase to approve task
	if err := h.usecase.ApproveTask(taskID, approverID, approverRole); err != nil {
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

// Helper to extract UserID safely from context
// The auth middleware sets "user_id" as float64 (from JSON)
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// JWT claims store numbers as float64
	if uid, ok := userID.(float64); ok {
		return uint(uid)
	}

	// Fallback for other types
	if uid, ok := userID.(uint); ok {
		return uid
	}

	if uid, ok := userID.(int); ok {
		return uint(uid)
	}

	return 0
}

// getUserRoleFromContext extracts the user's role from JWT context
func getUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}

	if roleStr, ok := role.(string); ok {
		return roleStr
	}

	return ""
}

// Helper to check if error message contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || (len(str) > len(substr) && 
		(str[:len(substr)] == substr || contains(str[1:], substr))))
}

// --- System Configuration Handlers (Admin/CEO Only) ---

// GetSystemConfig handles GET /admin/config
// Returns current AI configuration
func (h *SentinelHandler) GetSystemConfig(c *gin.Context) {
	config, err := h.usecase.GetSystemConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve system configuration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "System configuration retrieved successfully",
		"data":    config,
	})
}

// UpdateSystemConfig handles PUT /admin/config
// Updates AI configuration (CEO only)
func (h *SentinelHandler) UpdateSystemConfig(c *gin.Context) {
	// 1. Get user role from context
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found",
		})
		return
	}

	// 2. Parse request body
	var req updateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// 3. Call usecase (with CEO role validation)
	config, err := h.usecase.UpdateSystemConfig(
		req.ActiveModel,
		req.Temperature,
		req.CursorAssistance,
		userRole,
	)
	if err != nil {
		// Check for authorization error
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Check for validation error
		if contains(err.Error(), "must be") || contains(err.Error(), "invalid model") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update system configuration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "System configuration updated successfully. Changes take effect immediately.",
		"data":    config,
	})
}

// GetAvailableModels handles GET /admin/models
// Returns list of available Gemini models
func (h *SentinelHandler) GetAvailableModels(c *gin.Context) {
	models := h.usecase.GetAvailableModels()

	c.JSON(http.StatusOK, gin.H{
		"message": "Available Gemini models",
		"data":    models,
	})
}
