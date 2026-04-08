// Project CRUD, team / PM owner assignment, and AI plan / schedule / clear-plan endpoints.
package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

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
		if err.Error() == "project name is required" || err.Error() == "project name contains invalid characters (allowed: letters, numbers, spaces, hyphens, underscores)" || contains(err.Error(), "invalid project status") {
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
	taskLimit := 600
	if rawLimit := strings.TrimSpace(c.Query("tasks_limit")); rawLimit != "" {
		if parsed, parseErr := strconv.Atoi(rawLimit); parseErr == nil && parsed > 0 {
			taskLimit = parsed
		}
	}
	data, err := h.usecase.GetProjectDetailsPage(idStr, taskLimit, callerCtx(c))
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

// GetProjectTasksPage handles GET /sentinel/projects/:id/tasks for page 2+ task loading.
func (h *SentinelHandler) GetProjectTasksPage(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Project id or code is required"})
		return
	}
	limit := 600
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		if parsed, parseErr := strconv.Atoi(rawLimit); parseErr == nil && parsed > 0 {
			limit = parsed
		}
	}
	offset := 0
	if rawOffset := strings.TrimSpace(c.Query("offset")); rawOffset != "" {
		parsed, parseErr := strconv.Atoi(rawOffset)
		if parseErr != nil || parsed < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "offset must be a non-negative integer"})
			return
		}
		offset = parsed
	}
	cursorCreatedAt := strings.TrimSpace(c.Query("cursor_created_at"))
	cursorID := strings.TrimSpace(c.Query("cursor_id"))
	data, err := h.usecase.GetProjectTasksPage(idStr, limit, cursorCreatedAt, cursorID, offset, callerCtx(c))
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		if err.Error() == "project not found" || contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project tasks", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project tasks page retrieved successfully", "data": data})
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
		if err.Error() == "project name is required" || err.Error() == "project name contains invalid characters (allowed: letters, numbers, spaces, hyphens, underscores)" || contains(err.Error(), "invalid project status") {
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
	PmUserIDs []uint `json:"pm_user_ids"` // empty = clear all Product Owners for this project (JSON key kept for API compatibility)
}

// AssignProjectPmOwners handles PATCH /sentinel/projects/:id/pm-owners (CEO or MANAGER; only when teams feature is disabled).
func (h *SentinelHandler) AssignProjectPmOwners(c *gin.Context) {
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if roleStr != "CEO" && roleStr != "MANAGER" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "only CEO or MANAGER can assign project Product Owners"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign Product Owners", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project Product Owners updated", "data": updated})
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

// ScheduleProjectWithAI handles POST /sentinel/projects/:id/ai-schedule — estimate time + arrange timeline for existing tasks (CEO / Product Owner only).
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

// ClearProjectPlan handles POST /sentinel/projects/:id/clear-plan — removes all tasks, sprints, milestones, epics (CEO / Product Owner only).
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
