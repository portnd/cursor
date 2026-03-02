package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// RegisterRoutes registers all Sentinel module routes
func RegisterRoutes(router *gin.RouterGroup, usecase domain.SentinelUsecase) {
	handler := NewSentinelHandler(usecase)

	sentinelGroup := router.Group("/sentinel")
	// sentinelGroup.Use(middleware.Auth()) // Uncomment if you have auth middleware
	{
		// Task Management
		sentinelGroup.POST("/tasks", handler.CreateTask)                   // Create new task (CEO/PM)
		sentinelGroup.GET("/tasks/my", handler.GetMyTasks)                 // Get my assigned tasks (DEV)
		sentinelGroup.GET("/tasks/unassigned", handler.GetUnassignedTasks) // Get unassigned tasks (CEO/PM)
		sentinelGroup.GET("/tasks/approvals", handler.GetApprovals)        // Get approvals inbox (CEO/PM only)
		sentinelGroup.GET("/tasks", handler.GetAllTasks)                   // Get all tasks (ADMIN/PM overview)
		sentinelGroup.GET("/tasks/:id", handler.GetTaskByID)               // Get single task with submission history
		sentinelGroup.PATCH("/tasks/:id", handler.UpdateTask)              // Update task (Creator or CEO only, triggers AI re-estimation)
		sentinelGroup.DELETE("/tasks/:id", handler.DeleteTask)             // Delete task (Creator or CEO only)
		sentinelGroup.POST("/tasks/:id/assign", handler.AssignTask)        // Assign task to developer (PM)
		sentinelGroup.POST("/tasks/:id/submit", handler.SubmitWork)        // Submit code for task (DEV)
		sentinelGroup.POST("/tasks/:id/negotiate", handler.NegotiateTime)  // Developer negotiates AI time estimate
		sentinelGroup.POST("/tasks/:id/approve", handler.ApproveTask)      // Approve task after review (PM/CEO only)

		// Appeal System
		sentinelGroup.POST("/submissions/:id/appeal", handler.SubmitAppeal)  // Developer appeals AI verdict
		sentinelGroup.POST("/appeals/:id/resolve", handler.ResolveAppeal)    // PM/CEO resolves appeal
	}

	// Admin/CEO Configuration Management
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/config", handler.GetSystemConfig)        // Get current AI config
		adminGroup.PUT("/config", handler.UpdateSystemConfig)     // Update AI config (CEO only)
		adminGroup.GET("/models", handler.GetAvailableModels)     // Get available Gemini models
	}
}
