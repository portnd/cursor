// Admin AI configuration: system config, model list, usage counters.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	c.JSON(http.StatusOK, gin.H{"message": "Available GLM models", "data": models})
}

// GetAIUsage handles GET /admin/ai-usage — approximate AI API usage and remaining quota (from our request counter).
func (h *SentinelHandler) GetAIUsage(c *gin.Context) {
	usage := h.usecase.GetAIUsage()
	c.JSON(http.StatusOK, gin.H{"message": "AI usage (approximate)", "data": usage})
}
