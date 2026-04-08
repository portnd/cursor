// Sprints, milestones, comments, time logs, project analytics, and bulk task status.
package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

const (
	maxCommentAttachmentCount = 5
	maxCommentAttachmentBytes = 8 * 1024 * 1024 // 8 MiB per file
)

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
	content := ""
	attachments := make([]domain.TaskCommentAttachment, 0)
	contentType := c.ContentType()
	if strings.HasPrefix(contentType, "multipart/form-data") {
		content = strings.TrimSpace(c.PostForm("content"))
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "invalid multipart payload"})
			return
		}
		files := form.File["attachments"]
		if len(files) > maxCommentAttachmentCount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": fmt.Sprintf("最多上传 %d 个附件", maxCommentAttachmentCount)})
			return
		}
		for _, fileHeader := range files {
			attachment, err := toTaskCommentAttachment(fileHeader)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment", "message": err.Error()})
				return
			}
			attachments = append(attachments, attachment)
		}
	} else {
		var req addCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
			return
		}
		content = req.Content
	}

	comment, err := h.usecase.AddComment(task.ID, userID, content, attachments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment added", "data": comment})
}

func toTaskCommentAttachment(fileHeader *multipart.FileHeader) (domain.TaskCommentAttachment, error) {
	if fileHeader == nil {
		return domain.TaskCommentAttachment{}, errors.New("empty file")
	}
	if fileHeader.Size <= 0 {
		return domain.TaskCommentAttachment{}, errors.New("附件为空文件")
	}
	if fileHeader.Size > maxCommentAttachmentBytes {
		return domain.TaskCommentAttachment{}, fmt.Errorf("附件 %s 超过 8MB 限制", fileHeader.Filename)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return domain.TaskCommentAttachment{}, fmt.Errorf("无法读取附件 %s", fileHeader.Filename)
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return domain.TaskCommentAttachment{}, fmt.Errorf("读取附件失败: %s", fileHeader.Filename)
	}
	mimeType := fileHeader.Header.Get("Content-Type")
	if strings.TrimSpace(mimeType) == "" {
		mimeType = mime.TypeByExtension(strings.ToLower(filepath.Ext(fileHeader.Filename)))
	}
	if strings.TrimSpace(mimeType) == "" {
		mimeType = "application/octet-stream"
	}
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(raw)
	return domain.TaskCommentAttachment{
		FileName: fileHeader.Filename,
		MimeType: mimeType,
		Size:     fileHeader.Size,
		DataURL:  dataURL,
		IsImage:  strings.HasPrefix(mimeType, "image/"),
	}, nil
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
	var loggedDate *time.Time
	if req.LoggedDate != "" {
		d, err := time.Parse("2006-01-02", req.LoggedDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "logged_date must be YYYY-MM-DD format"})
			return
		}
		loggedDate = &d
	}
	entry, err := h.usecase.LogTime(task.ID, userID, req.Minutes, req.Description, req.WorkType, loggedDate, req.IsTimerSession)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log time", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Time logged", "data": entry})
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

// EditTimeLog updates minutes, description, or work_type of a log (owner only, within 24h).
// PATCH /api/v1/sentinel/time-logs/:logId
func (h *SentinelHandler) EditTimeLog(c *gin.Context) {
	logIDStr := c.Param("logId")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid log ID"})
		return
	}
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req editTimeLogReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	var taskID *uuid.UUID
	if req.TaskID != nil && *req.TaskID != "" {
		parsed, err := uuid.Parse(*req.TaskID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid task_id"})
			return
		}
		taskID = &parsed
	}
	updated, err := h.usecase.EditTimeLog(logID, userID, req.Minutes, req.Description, req.WorkType, taskID)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update time log", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time log updated", "data": updated})
}

// DeleteTimeLog removes a log (owner only, within 24h).
// DELETE /api/v1/sentinel/time-logs/:logId
func (h *SentinelHandler) DeleteTimeLog(c *gin.Context) {
	logIDStr := c.Param("logId")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid log ID"})
		return
	}
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := h.usecase.DeleteTimeLog(logID, userID); err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete time log", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time log deleted"})
}

// GetMyDailyTimeLogs returns the caller's time logs for a given date (defaults today).
// GET /api/v1/sentinel/users/me/time-logs?date=YYYY-MM-DD
func (h *SentinelHandler) GetMyDailyTimeLogs(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	dateStr := c.Query("date")
	date := time.Now().UTC()
	if dateStr != "" {
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "date must be YYYY-MM-DD"})
			return
		}
		date = d
	}
	summary, err := h.usecase.GetMyDailyTimeLogs(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch daily time logs", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Daily time logs retrieved", "data": summary})
}

// BulkLogTime processes multiple time log entries in one request (EOD batch).
// POST /api/v1/sentinel/time-logs/bulk
func (h *SentinelHandler) BulkLogTime(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var entries []domain.BulkLogEntry
	if err := c.ShouldBindJSON(&entries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	results, err := h.usecase.BulkLogTime(entries, userID)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process bulk log", "message": err.Error()})
		return
	}
	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Processed %d entries, %d succeeded", len(results), successCount),
		"success_count": successCount,
		"total":         len(results),
		"results":       results,
	})
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
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user role not found"})
		return
	}
	normalizedRole := strings.ToUpper(strings.TrimSpace(userRole))
	if req.Status == "COMPLETED" && normalizedRole != authDomain.RoleCEO && normalizedRole != authDomain.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "cannot drag task to COMPLETED — only CEO or Manager can perform final approval"})
		return
	}
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}
	if err := h.usecase.BulkUpdateTaskStatus(taskIDs, req.Status, userID); err != nil {
		if contains(err.Error(), "invalid status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		if contains(err.Error(), "cannot drag task to COMPLETED") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		if contains(err.Error(), "READY_FOR_UAT") || contains(err.Error(), "no tasks provided") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bulk update status", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated %d tasks to %s", len(taskIDs), req.Status)})
}
