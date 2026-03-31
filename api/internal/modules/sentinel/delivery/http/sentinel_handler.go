package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
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
	StoryPoints      int     `json:"story_points"`
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
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	TaskType         string  `json:"task_type"` // FEATURE, TASK, BUG
	ParentID         *string `json:"parent_id"`
	EpicID           *string `json:"epic_id"`
	SortOrder        *int    `json:"sort_order"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
	Progress         *int    `json:"progress"`
	Priority         string  `json:"priority"`
	StoryPoints      *int    `json:"story_points"`
	SprintID         *string `json:"sprint_id"`
	MilestoneID      *string `json:"milestone_id"`
	EstimatedMinutes *int    `json:"estimated_minutes"` // Manual estimate; feeds Costing Engine
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

type logTimeReq struct {
	Minutes     int    `json:"minutes" binding:"required,gt=0"`
	Description string `json:"description"`
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

// --- Handlers ---

// callerCtx extracts CallerContext (role + team_id + user_id) from the Gin context (set by AuthMiddleware).
func callerCtx(c *gin.Context) domain.CallerContext {
	role, _ := c.Get("role")
	teamID, _ := c.Get("team_id")
	ctx := domain.CallerContext{UserID: getUserIDFromContext(c)}
	if r, ok := role.(string); ok {
		ctx.Role = r
	}
	if t, ok := teamID.(*uint); ok {
		ctx.TeamID = t
	}
	return ctx
}

// CreateProject handles POST /sentinel/projects
func (h *SentinelHandler) CreateProject(c *gin.Context) {
	var req createProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	project, err := h.usecase.CreateProject(req.Name, req.Description, status, callerCtx(c))
	if err != nil {
		if err.Error() == "project name is required" || err.Error() == "project name must be in English only (letters, numbers, spaces, hyphens)" || contains(err.Error(), "invalid project status") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create project",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"data":    project,
	})
}

// GetProjects handles GET /sentinel/projects
func (h *SentinelHandler) GetProjects(c *gin.Context) {
	projects, err := h.usecase.GetProjects(callerCtx(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve projects",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Projects retrieved successfully",
		"data":    projects,
	})
}

// GetProjectDetails handles GET /sentinel/projects/:id/details — returns project + tasks + sprints + milestones + epics (1 round-trip).
func (h *SentinelHandler) GetProjectDetails(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Project id or code is required"})
		return
	}
	data, err := h.usecase.GetProjectDetailsPage(idStr, callerCtx(c))
	if err != nil {
		if err.Error() == "project not found" || contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project details retrieved successfully", "data": data})
}

// GetProjectByID handles GET /sentinel/projects/:id (id may be UUID or project code e.g. mims-hdmap-main)
func (h *SentinelHandler) GetProjectByID(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Project id or code is required",
		})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil {
		if err.Error() == "project not found" || contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve project",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Project retrieved successfully",
		"data":    project,
	})
}

type updateProjectReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdateCode  bool   `json:"update_code"` // if true, set project code to slugify(name) and update all task codes to new prefix
}

// UpdateProject handles PATCH /sentinel/projects/:id (id may be UUID or project code)
func (h *SentinelHandler) UpdateProject(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Project id or code is required"})
		return
	}
	existing, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
		return
	}
	var req updateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	status := req.Status
	if status == "" {
		status = existing.Status
	}
	project, err := h.usecase.UpdateProject(existing.ID, req.Name, req.Description, status, req.UpdateCode)
	if err != nil {
		if err.Error() == "project name is required" || err.Error() == "project name must be in English only (letters, numbers, spaces, hyphens)" || contains(err.Error(), "invalid project status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "data": project})
}

// DeleteProject handles DELETE /sentinel/projects/:id (id may be UUID or project code)
func (h *SentinelHandler) DeleteProject(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Project id or code is required",
		})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil {
		if err.Error() == "project not found" || contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Project not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	if err := h.usecase.DeleteProject(project.ID); err != nil {
		if err.Error() == "project not found" || contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete project",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
	})
}

// assignProjectTeamReq is the DTO for PATCH /sentinel/projects/:id/assign-team
type assignProjectTeamReq struct {
	TeamID *uint `json:"team_id"` // null = unassign from team
}

// AssignProjectTeam handles PATCH /sentinel/projects/:id/assign-team (CEO only)
func (h *SentinelHandler) AssignProjectTeam(c *gin.Context) {
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if roleStr != "CEO" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO can assign teams to projects"})
		return
	}
	idStr := strings.TrimSpace(c.Param("id"))
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "project not found"})
		return
	}
	var req assignProjectTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	updated, err := h.usecase.AssignProjectTeam(project.ID, req.TeamID, roleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign team", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team assigned successfully", "data": updated})
}

type assignProjectPmOwnersReq struct {
	PmUserIDs []uint `json:"pm_user_ids"` // empty = clear all PM owners for this project
}

// AssignProjectPmOwners handles PATCH /sentinel/projects/:id/pm-owners (CEO or MANAGER; only when teams feature is disabled).
func (h *SentinelHandler) AssignProjectPmOwners(c *gin.Context) {
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if roleStr != "CEO" && roleStr != "MANAGER" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO or MANAGER can assign project PM owners"})
		return
	}
	idStr := strings.TrimSpace(c.Param("id"))
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "project not found"})
		return
	}
	var req assignProjectPmOwnersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	updated, err := h.usecase.AssignProjectPmOwners(project.ID, req.PmUserIDs, roleStr)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign PM owners", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project PM owners updated", "data": updated})
}

func (h *SentinelHandler) GenerateProjectPlan(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project id required"})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "project not found"})
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
	plan, err := h.usecase.GenerateProjectPlan(project.ID, userID, userRole)
	if err != nil {
		log.Printf("[GenerateProjectPlan] error: %v", err)
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI plan failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "AI work plan created successfully",
		"data":    plan,
	})
}

// ScheduleProjectWithAI handles POST /sentinel/projects/:id/ai-schedule — estimate time + arrange timeline for existing tasks (CEO/PM only).
func (h *SentinelHandler) ScheduleProjectWithAI(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project id required"})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "project not found"})
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
	updatedCount, err := h.usecase.ScheduleProjectWithAI(project.ID, userID, userRole)
	if err != nil {
		log.Printf("[ScheduleProjectWithAI] error: %v", err)
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		if contains(err.Error(), "no tasks to schedule") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI schedule failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "AI estimated and scheduled tasks successfully",
		"updated": updatedCount,
	})
}

// ClearProjectPlan handles POST /sentinel/projects/:id/clear-plan — removes all tasks, sprints, milestones, epics (CEO/PM only).
func (h *SentinelHandler) ClearProjectPlan(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project id required"})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "project not found"})
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
	if err := h.usecase.ClearProjectPlan(project.ID, userID, userRole); err != nil {
		if contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear plan", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project plan cleared successfully"})
}

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

	task, err := h.usecase.CreateTask(req.Title, req.Description, req.TaskType, userID, dueDate, projectID, parentID, startDate, endDate, req.Priority, req.StoryPoints, sprintID, milestoneID, epicID, req.EstimatedMinutes)
	if err != nil {
		log.Printf("[CreateTask] usecase error: %v", err)
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
// :id can be UUID or task code. The requesting user (PM/CEO) is recorded as assigned_by for leaderboard scope.
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

	if err := h.usecase.AssignTask(task.ID, req.DevID, assignerID); err != nil {
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

// SubmitUAT stores UAT payload on a FEATURE task and promotes it to REVIEW_PENDING for PM/CEO review.
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

	c.JSON(http.StatusOK, gin.H{"message": "UAT submitted — feature is now pending PM/CEO review"})
}

// GetTaskByID handles GET /api/v1/sentinel/tasks/:id
// :id can be UUID or task code (e.g. mims-hdmap-main-001)
func (h *SentinelHandler) GetTaskByID(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "task id or code is required",
		})
		return
	}
	task, err := h.usecase.GetTaskByIDOrCode(idStr)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "task id or code is required" ||
			strings.Contains(errMsg, "task not found") ||
			strings.Contains(errMsg, "record not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve task",
			"message": errMsg,
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
// TASK/BUG in ACTIVE sprints: CEO/MANAGER = company-wide; others team-scoped; teams off → PM/DEV assignment rules.
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

// GetTeamActiveTasks handles GET /api/v1/sentinel/tasks/team-active
// Returns all tasks in ACTIVE sprints within the caller's team.
// Used by the "Quick Log Time" global modal to let devs log time against any teammate's task.
// CEO/MANAGER bypass team restriction and see all active-sprint tasks.
func (h *SentinelHandler) GetTeamActiveTasks(c *gin.Context) {
	ctx := callerCtx(c)
	tasks, err := h.usecase.GetTeamActiveTasks(ctx.TeamID, ctx.Role)
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
// Returns all FEATURE-type tasks for the PM/CEO Feature Roadmap Board.
// Each feature includes a roll-up progress (0-100%) and its child TASK/BUG items for the accordion.
// PM is scoped to their team; CEO/MANAGER see all teams.
func (h *SentinelHandler) GetActiveFeatures(c *gin.Context) {
	ctx := callerCtx(c)
	items, err := h.usecase.GetActiveFeatures(ctx.TeamID, ctx.Role)
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
// Without project_id: returns ALL tasks (ADMIN/PM overview).
func (h *SentinelHandler) GetAllTasks(c *gin.Context) {
	var tasks []domain.Task
	var err error
	var message string
	if idStr := c.Query("project_id"); idStr != "" {
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
		"message": "Time negotiation submitted successfully and pending PM/CEO review",
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
	hasStart := req.StartDate != nil && *req.StartDate != ""
	hasEnd := req.EndDate != nil && *req.EndDate != ""
	hasProgress := req.Progress != nil
	hasPriority := req.Priority != ""
	hasSP := req.StoryPoints != nil
	hasSprint := req.SprintID != nil
	hasMilestone := req.MilestoneID != nil
	hasSortOrder := req.SortOrder != nil
	hasEstMins := req.EstimatedMinutes != nil
	if !hasTitle && !hasDesc && !hasTaskType && !hasParent && !hasEpicKey && !hasStart && !hasEnd && !hasProgress && !hasPriority && !hasSP && !hasSprint && !hasMilestone && !hasSortOrder && !hasEstMins {
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
	var startDate, endDate *time.Time
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

	task, err := h.usecase.UpdateTask(taskResolved.ID, userID, userRole, req.Title, req.Description, req.TaskType, parentID, startDate, endDate, progress, req.Priority, req.StoryPoints, sprintIDUpd, hasSprint, milestoneIDUpd, epicIDUpd, hasEpicKey, sortOrderUpd, req.EstimatedMinutes)
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

// ApproveTask marks a task as COMPLETED after human verification (PM/CEO only)
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

// RejectTask returns a task to IN_PROGRESS with a rejection reason comment (PM/CEO/MANAGER only)
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

// PMApproveSubTask is the PM's first-stage approval: READY_FOR_TEST → READY_FOR_UAT.
// The PM must supply a test URL and detailed test steps for the CEO to follow.
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

// RejectSubTask moves a READY_FOR_TEST task back to IN_PROGRESS and logs the reason (PM/CEO action).
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

func (h *SentinelHandler) resolveTaskIDOrCode(c *gin.Context) (*domain.Task, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, fmt.Errorf("task id or code is required")
	}
	return h.usecase.GetTaskByIDOrCode(idStr)
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
func (h *SentinelHandler) GetAvailableModels(c *gin.Context) {
	models := h.usecase.GetAvailableModels()
	c.JSON(http.StatusOK, gin.H{"message": "Available Gemini models", "data": models})
}

// GetAIUsage handles GET /admin/ai-usage — approximate Gemini API usage and remaining quota (from our request counter).
func (h *SentinelHandler) GetAIUsage(c *gin.Context) {
	usage := h.usecase.GetAIUsage()
	c.JSON(http.StatusOK, gin.H{"message": "AI usage (approximate)", "data": usage})
}

// --- Sprint Handlers ---

func (h *SentinelHandler) CreateSprint(c *gin.Context) {
	var req createSprintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "start_date must be ISO8601"})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "end_date must be ISO8601"})
			return
		}
		endDate = &t
	}
	sprint, err := h.usecase.CreateSprint(projectID, req.Name, req.Goal, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Sprint created", "data": sprint})
}

func (h *SentinelHandler) GetSprintsByProject(c *gin.Context) {
	idStr := c.Query("project_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project_id query param required"})
		return
	}
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	sprints, err := h.usecase.GetSprintsByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sprints", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprints retrieved", "data": sprints})
}

func (h *SentinelHandler) StartSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	sprint, err := h.usecase.StartSprint(sprintID)
	if err != nil {
		if contains(err.Error(), "already") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Conflict", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint started", "data": sprint})
}

func (h *SentinelHandler) CompleteSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	sprint, err := h.usecase.CompleteSprint(sprintID)
	if err != nil {
		if contains(err.Error(), "only ACTIVE") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint completed", "data": sprint})
}

func (h *SentinelHandler) ReopenSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	sprint, err := h.usecase.ReopenSprint(sprintID)
	if err != nil {
		if contains(err.Error(), "only COMPLETED") || contains(err.Error(), "already has an active sprint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reopen sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint reopened", "data": sprint})
}

func (h *SentinelHandler) AddTasksToSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	var req addTasksToSprintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var taskIDs []uuid.UUID
	for _, idStr := range req.TaskIDs {
		tid, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "task_ids must all be valid UUIDs"})
			return
		}
		taskIDs = append(taskIDs, tid)
	}
	if err := h.usecase.AddTasksToSprint(sprintID, taskIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tasks to sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tasks added to sprint"})
}

func (h *SentinelHandler) UpdateSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	var req updateSprintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "start_date must be ISO8601"})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "end_date must be ISO8601"})
			return
		}
		endDate = &t
	}
	sprint, err := h.usecase.UpdateSprint(sprintID, req.Name, req.Goal, startDate, endDate, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint updated", "data": sprint})
}

func (h *SentinelHandler) DeleteSprint(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid sprint UUID"})
		return
	}
	if err := h.usecase.DeleteSprint(sprintID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sprint", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint deleted"})
}

// --- Milestone Handlers ---

func (h *SentinelHandler) CreateMilestone(c *gin.Context) {
	var req createMilestoneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "due_date must be ISO8601"})
			return
		}
		dueDate = &t
	}
	milestone, err := h.usecase.CreateMilestone(projectID, req.Title, req.Description, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Milestone created", "data": milestone})
}

func (h *SentinelHandler) GetMilestonesByProject(c *gin.Context) {
	idStr := c.Query("project_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project_id query param required"})
		return
	}
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	milestones, err := h.usecase.GetMilestonesByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get milestones", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestones retrieved", "data": milestones})
}

func (h *SentinelHandler) UpdateMilestone(c *gin.Context) {
	milestoneID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid milestone UUID"})
		return
	}
	var req updateMilestoneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date", "message": "due_date must be ISO8601"})
			return
		}
		dueDate = &t
	}
	milestone, err := h.usecase.UpdateMilestone(milestoneID, req.Title, req.Description, req.Status, dueDate)
	if err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated", "data": milestone})
}

func (h *SentinelHandler) DeleteMilestone(c *gin.Context) {
	milestoneID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Invalid milestone UUID"})
		return
	}
	if err := h.usecase.DeleteMilestone(milestoneID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete milestone", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted"})
}

// --- Comment Handlers ---

func (h *SentinelHandler) AddComment(c *gin.Context) {
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
	var req addCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	comment, err := h.usecase.AddComment(task.ID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment added", "data": comment})
}

func (h *SentinelHandler) GetComments(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	comments, err := h.usecase.GetComments(task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comments retrieved", "data": comments})
}

// --- Time Log Handlers ---

func (h *SentinelHandler) LogTime(c *gin.Context) {
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
	var req logTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	log, err := h.usecase.LogTime(task.ID, userID, req.Minutes, req.Description)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log time", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Time logged", "data": log})
}

func (h *SentinelHandler) GetTimeLogs(c *gin.Context) {
	task, err := h.resolveTaskIDOrCode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
		return
	}
	logs, err := h.usecase.GetTimeLogs(task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get time logs", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time logs retrieved", "data": logs})
}

// --- Analytics Handler ---

func (h *SentinelHandler) GetProjectAnalytics(c *gin.Context) {
	idStr := c.Param("id")
	var projectID uuid.UUID
	var err error
	if projectID, err = uuid.Parse(idStr); err != nil {
		project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		projectID = project.ID
	}
	analytics, err := h.usecase.GetProjectAnalytics(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Analytics retrieved", "data": analytics})
}

// --- Bulk Status Handler ---

func (h *SentinelHandler) BulkUpdateTaskStatus(c *gin.Context) {
	var req bulkStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var taskIDs []uuid.UUID
	for _, idStr := range req.TaskIDs {
		tid, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "task_ids must all be valid UUIDs"})
			return
		}
		taskIDs = append(taskIDs, tid)
	}
	if err := h.usecase.BulkUpdateTaskStatus(taskIDs, req.Status); err != nil {
		if contains(err.Error(), "invalid status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bulk update status", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated %d tasks to %s", len(taskIDs), req.Status)})
}

// --- Google Slides Import ---

func (h *SentinelHandler) PreviewGoogleSlides(c *gin.Context) {
	var req domain.PreviewGoogleSlidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewGoogleSlides(&req, h.googleAPIKey)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "API error 403") || contains(err.Error(), "API error 401") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportGoogleSlides(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportGoogleSlidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromGoogleSlides(&req, h.googleAPIKey, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "API error 403") || contains(err.Error(), "API error 401") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.PresentationTitle),
		"data":    result,
	})
}

// --- Google Sheets Import ---

func (h *SentinelHandler) PreviewGoogleSheets(c *gin.Context) {
	var req domain.PreviewGoogleSheetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewGoogleSheets(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "no importable") || contains(err.Error(), "no data rows") {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportGoogleSheets(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportGoogleSheetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromGoogleSheets(&req, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one row") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") {
			statusCode = http.StatusBadRequest
		}
		if domain.IsBadRequest(err) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.SheetTitle),
		"data":    result,
	})
}

// --- Canva Import ---

func (h *SentinelHandler) PreviewCanva(c *gin.Context) {
	var req domain.PreviewCanvaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewCanva(&req, h.canvaAccessToken)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "empty") || contains(err.Error(), "could not find") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "HTTP 401") || contains(err.Error(), "HTTP 403") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportCanva(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportCanvaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromCanva(&req, h.canvaAccessToken, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") || contains(err.Error(), "no matching") {
			statusCode = http.StatusBadRequest
		}
		if domain.IsBadRequest(err) {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "HTTP 401") || contains(err.Error(), "HTTP 403") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.DesignTitle),
		"data":    result,
	})
}

// --- PPTX File Upload Import ---

// maxPPTXUploadBytes: Canva / rich decks with images often exceed 50 MB.
const maxPPTXUploadBytes = 250 * 1024 * 1024 // 250 MiB

func (h *SentinelHandler) PreviewPPTX(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file", "message": "multipart field 'file' is required"})
		return
	}
	if fileHeader.Size > maxPPTXUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": fmt.Sprintf("PPTX file must be under %d MB (got %.1f MB)", maxPPTXUploadBytes/(1024*1024), float64(fileHeader.Size)/(1024*1024)),
		})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": err.Error()})
		return
	}
	defer f.Close()
	data, err := readAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewPPTX(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PPTX", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportPPTX(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file", "message": "multipart field 'file' is required"})
		return
	}
	if fileHeader.Size > maxPPTXUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": fmt.Sprintf("PPTX file must be under %d MB (got %.1f MB)", maxPPTXUploadBytes/(1024*1024), float64(fileHeader.Size)/(1024*1024)),
		})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": err.Error()})
		return
	}
	defer f.Close()
	data, err := readAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "message": err.Error()})
		return
	}
	payloadStr := c.PostForm("payload")
	if payloadStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing payload", "message": "form field 'payload' (JSON) is required"})
		return
	}
	var req domain.ImportPPTXRequest
	if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "message": err.Error()})
		return
	}
	result, err := h.usecase.ImportFromPPTX(data, &req, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if domain.IsBadRequest(err) || contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") || contains(err.Error(), "no matching") {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.Title),
		"data":    result,
	})
}

// readAll reads all bytes from a ReadCloser (alias for io.ReadAll used in handlers).
func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	return io.ReadAll(r)
}

// --- Epic Handlers (Hierarchy Dimension 1) ---

// CreateEpic handles POST /sentinel/epics
func (h *SentinelHandler) CreateEpic(c *gin.Context) {
	var req createEpicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "start_date must be ISO8601/RFC3339"})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "end_date must be ISO8601/RFC3339"})
			return
		}
		endDate = &t
	}
	epic, err := h.usecase.CreateEpic(projectID, req.Title, req.Description, req.Color, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create epic", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Epic created successfully", "data": epic})
}

// GetEpicsByProject handles GET /sentinel/epics?project_id=UUID
func (h *SentinelHandler) GetEpicsByProject(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "project_id query param is required"})
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project_id must be a valid UUID"})
		return
	}
	epics, err := h.usecase.GetEpicsByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve epics", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Epics retrieved successfully", "data": epics})
}

// UpdateEpic handles PATCH /sentinel/epics/:id
func (h *SentinelHandler) UpdateEpic(c *gin.Context) {
	epicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "epic id must be a valid UUID"})
		return
	}
	var req updateEpicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "start_date must be ISO8601/RFC3339"})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "message": "end_date must be ISO8601/RFC3339"})
			return
		}
		endDate = &t
	}
	epic, err := h.usecase.UpdateEpic(epicID, req.Title, req.Description, req.Status, req.Color, req.SortOrder, startDate, endDate)
	if err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update epic", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Epic updated successfully", "data": epic})
}

// DeleteEpic handles DELETE /sentinel/epics/:id
func (h *SentinelHandler) DeleteEpic(c *gin.Context) {
	epicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "epic id must be a valid UUID"})
		return
	}
	if err := h.usecase.DeleteEpic(epicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete epic", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Epic deleted successfully"})
}

// --- Timeline View Handlers (Matrix Dimension) ---

// GetEpicTimeline handles GET /sentinel/projects/:id/timeline/epic-view — :id may be UUID or project code (e.g. mims-hd-map).
func (h *SentinelHandler) GetEpicTimeline(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project id or code is required"})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		if err != nil && (err.Error() == "project not found" || contains(err.Error(), "not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project", "message": err.Error()})
		return
	}
	data, err := h.usecase.GetEpicTimelineData(project.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve epic timeline", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Epic timeline retrieved successfully", "data": data})
}

// GetSprintTimeline handles GET /sentinel/projects/:id/timeline/sprint-view — :id may be UUID or project code (e.g. mims-hd-map).
func (h *SentinelHandler) GetSprintTimeline(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project id or code is required"})
		return
	}
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	if err != nil || project == nil {
		if err != nil && (err.Error() == "project not found" || contains(err.Error(), "not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project", "message": err.Error()})
		return
	}
	data, err := h.usecase.GetSprintTimelineData(project.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sprint timeline", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sprint timeline retrieved successfully", "data": data})
}

// ExportTimelinePDF handles GET /sentinel/projects/:id/timeline/export-pdf?mode=epic|sprint
// Returns raw PDF bytes (application/pdf) — same pattern as mims-api-service ExportPDF.
// Frontend fetches with Authorization header, receives blob, opens in new tab.
func (h *SentinelHandler) ExportTimelinePDF(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	mode := c.DefaultQuery("mode", "epic")
	if mode != "epic" && mode != "sprint" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode must be 'epic' or 'sprint'"})
		return
	}

	// Template dir: from env TEMPLATE_DIR or default to /app/templates/
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = "./templates/"
	}
	if templateDir[len(templateDir)-1] != '/' {
		templateDir += "/"
	}

	log.Printf("ExportTimelinePDF: projectID=%s mode=%s templateDir=%s", projectID, mode, templateDir)

	pdfBytes, filename, err := h.usecase.ExportTimelinePDF(projectID, mode, templateDir)
	if err != nil {
		log.Printf("ExportTimelinePDF error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PDF generation failed", "message": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	c.Header("Content-Description", "Timeline PDF Report")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// SplitTask decomposes one task into N sub-tasks and deletes the original.
func (h *SentinelHandler) SplitTask(c *gin.Context) {
	requesterID := getUserIDFromContext(c)
	if requesterID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	requesterRole := getUserRoleFromContext(c)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Splits []domain.SplitTaskItem `json:"splits" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	created, err := h.usecase.SplitTask(taskID, req.Splits, requesterID, requesterRole)
	if err != nil {
		status := http.StatusInternalServerError
		if domain.IsBadRequest(err) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": "Split failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Task split into %d sub-tasks", len(created)),
		"data":    created,
	})
}
