package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type projectBackupHandler struct {
	usecase domain.SentinelUsecase
}

func NewProjectBackupHandler(usecase domain.SentinelUsecase) *projectBackupHandler {
	return &projectBackupHandler{usecase: usecase}
}

// ListProjectBackups GET /sentinel/projects/:id/backups
func (h *projectBackupHandler) ListProjectBackups(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	backups, err := h.usecase.GetProjectBackups(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": backups})
}

// CreateProjectBackup POST /sentinel/projects/:id/backups
func (h *projectBackupHandler) CreateProjectBackup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req domain.CreateProjectBackupRequest
	_ = c.ShouldBindJSON(&req) // optional body — label may be empty

	// Extract caller user ID from auth middleware claims (optional).
	var createdBy *uint
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			createdBy = &id
		}
	}

	backup, err := h.usecase.CreateProjectBackup(projectID, req.Label, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": backup})
}

// RestoreProjectBackup POST /sentinel/projects/:id/backups/:backupId/restore
func (h *projectBackupHandler) RestoreProjectBackup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	backupID, err := uuid.Parse(c.Param("backupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup id"})
		return
	}

	if err := h.usecase.RestoreProjectBackup(backupID, projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "restored"})
}

// DeleteProjectBackup DELETE /sentinel/projects/:id/backups/:backupId
func (h *projectBackupHandler) DeleteProjectBackup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	backupID, err := uuid.Parse(c.Param("backupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup id"})
		return
	}

	if err := h.usecase.DeleteProjectBackup(backupID, projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetProjectBackupPayload GET /sentinel/projects/:id/backups/:backupId/payload
func (h *projectBackupHandler) GetProjectBackupPayload(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	backupID, err := uuid.Parse(c.Param("backupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup id"})
		return
	}
	payload, err := h.usecase.GetProjectBackupPayload(projectID, backupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": payload})
}

// ImportProjectFromBackup POST /sentinel/projects/import-backup
func (h *projectBackupHandler) ImportProjectFromBackup(c *gin.Context) {
	var req domain.ImportProjectFromBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createdBy *uint
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			createdBy = &id
		}
	}

	project, err := h.usecase.ImportProjectFromBackup(req.Name, &req.Payload, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": project})
}
