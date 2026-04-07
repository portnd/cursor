// Epics, sprint/epic timeline JSON, timeline PDF export, task split.
package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

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
	startedAt := time.Now()
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "project id or code is required"})
		return
	}

	projectLookupStartedAt := time.Now()
	project, err := h.usecase.GetProjectByIDOrCode(idStr, callerCtx(c))
	projectLookupElapsedMs := time.Since(projectLookupStartedAt).Milliseconds()
	if err != nil || project == nil {
		if err != nil && (err.Error() == "project not found" || contains(err.Error(), "not found")) {
			log.Printf("[timeline][epic] project_not_found id_or_code=%s project_lookup_ms=%d total_ms=%d", idStr, projectLookupElapsedMs, time.Since(startedAt).Milliseconds())
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Project not found"})
			return
		}
		if err != nil {
			log.Printf("[timeline][epic] project_lookup_error id_or_code=%s error=%v project_lookup_ms=%d total_ms=%d", idStr, err, projectLookupElapsedMs, time.Since(startedAt).Milliseconds())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project", "message": err.Error()})
		return
	}

	timelineFetchStartedAt := time.Now()
	data, err := h.usecase.GetEpicTimelineData(project.ID)
	timelineFetchElapsedMs := time.Since(timelineFetchStartedAt).Milliseconds()
	if err != nil {
		log.Printf("[timeline][epic] timeline_fetch_error project_id=%s id_or_code=%s error=%v project_lookup_ms=%d timeline_fetch_ms=%d total_ms=%d", project.ID, idStr, err, projectLookupElapsedMs, timelineFetchElapsedMs, time.Since(startedAt).Milliseconds())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve epic timeline", "message": err.Error()})
		return
	}

	totalElapsedMs := time.Since(startedAt).Milliseconds()
	slowThresholdMs := int64(300)
	if totalElapsedMs >= slowThresholdMs {
		epicCount := len(data.Epics)
		taskCount := 0
		subTaskCount := 0
		for _, e := range data.Epics {
			taskCount += len(e.Tasks)
			for _, t := range e.Tasks {
				subTaskCount += len(t.SubTasks)
			}
		}
		log.Printf("[timeline][epic] slow project_id=%s id_or_code=%s epics=%d tasks=%d subtasks=%d project_lookup_ms=%d timeline_fetch_ms=%d total_ms=%d threshold_ms=%d", project.ID, idStr, epicCount, taskCount, subTaskCount, projectLookupElapsedMs, timelineFetchElapsedMs, totalElapsedMs, slowThresholdMs)
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
