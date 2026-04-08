// Komgrip: project-less personal/misc tasks accessible to all employees.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type createKomgripTaskReq struct {
	Title            string `json:"title" binding:"required"`
	Description      string `json:"description"`
	Priority         string `json:"priority"`
	EstimatedMinutes int    `json:"estimated_minutes"`
}

type updateKomgripStatusReq struct {
	Status string `json:"status" binding:"required"`
}

// GetKomgripTasks handles GET /sentinel/komgrip/tasks
func (h *SentinelHandler) GetKomgripTasks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	tasks, err := h.usecase.GetKomgripTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get komgrip tasks", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// CreateKomgripTask handles POST /sentinel/komgrip/tasks
func (h *SentinelHandler) CreateKomgripTask(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	var req createKomgripTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	task, err := h.usecase.CreateKomgripTask(req.Title, req.Description, userID, req.Priority, req.EstimatedMinutes)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create komgrip task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Komgrip task created successfully", "data": task})
}

// UpdateKomgripTaskStatus handles PATCH /sentinel/komgrip/tasks/:id/status
func (h *SentinelHandler) UpdateKomgripTaskStatus(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "task id must be a valid UUID"})
		return
	}

	var req updateKomgripStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	task, err := h.usecase.UpdateKomgripTaskStatus(taskID, req.Status, userID)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update komgrip task status", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "data": task})
}

// DeleteKomgripTask handles DELETE /sentinel/komgrip/tasks/:id
func (h *SentinelHandler) DeleteKomgripTask(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "task id must be a valid UUID"})
		return
	}

	task, err := h.usecase.GetTaskByID(taskID)
	if err != nil || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "task not found"})
		return
	}
	if !task.IsKomgrip {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "task is not a Komgrip task"})
		return
	}

	if err := h.usecase.DeleteTask(taskID, userID, getUserRoleFromContext(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete komgrip task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komgrip task deleted"})
}
