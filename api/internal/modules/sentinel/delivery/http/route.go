package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// RegisterRoutes registers all Sentinel module routes
func RegisterRoutes(router *gin.RouterGroup, usecase domain.SentinelUsecase, googleAPIKey ...string) {
	key := ""
	if len(googleAPIKey) > 0 {
		key = googleAPIKey[0]
	}
	handler := NewSentinelHandler(usecase, key)

	sentinelGroup := router.Group("/sentinel")
	// sentinelGroup.Use(middleware.Auth()) // Uncomment if you have auth middleware
	{
		// Project Management
		sentinelGroup.POST("/projects", handler.CreateProject)
		sentinelGroup.GET("/projects", handler.GetProjects)
		sentinelGroup.GET("/projects/:id", handler.GetProjectByID)
		sentinelGroup.DELETE("/projects/:id", handler.DeleteProject)

		// Task Management
		sentinelGroup.POST("/tasks", handler.CreateTask)                   // Create new task (CEO/PM)
		sentinelGroup.GET("/tasks/my", handler.GetMyTasks)                 // Get my assigned tasks (DEV)
		sentinelGroup.GET("/tasks/unassigned", handler.GetUnassignedTasks) // Get unassigned tasks (CEO/PM)
		sentinelGroup.GET("/tasks/approvals", handler.GetApprovals)        // Get approvals inbox (CEO/PM only)
		sentinelGroup.GET("/tasks/gantt", handler.GetGantt)                // Get all tasks + dependencies for Gantt chart
		sentinelGroup.GET("/tasks", handler.GetAllTasks)                    // Get all tasks (ADMIN/PM overview)
		sentinelGroup.GET("/tasks/:id", handler.GetTaskByID)                // Get single task with submission history
		sentinelGroup.PATCH("/tasks/:id", handler.UpdateTask)                      // Update task (Creator or CEO only, triggers AI re-estimation)
		sentinelGroup.PATCH("/tasks/:id/slide-resources", handler.UpdateTaskSlideResources) // Update task resource_urls (slide images/annotations)
		sentinelGroup.DELETE("/tasks/:id", handler.DeleteTask)                     // Delete task (Creator or CEO only)
		sentinelGroup.POST("/tasks/:id/assign", handler.AssignTask)        // Assign task to developer (PM)
		sentinelGroup.POST("/tasks/:id/submit", handler.SubmitWork)        // Submit code for task (DEV)
		sentinelGroup.POST("/tasks/:id/negotiate", handler.NegotiateTime)  // Developer negotiates AI time estimate
		sentinelGroup.POST("/tasks/:id/approve", handler.ApproveTask)      // Approve task after review (PM/CEO only)

		// Task Dependencies (Gantt links)
		sentinelGroup.POST("/tasks/dependencies", handler.CreateDependency)
		sentinelGroup.DELETE("/tasks/dependencies/:id", handler.DeleteDependency)

		// Task Comments
		sentinelGroup.POST("/tasks/:id/comments", handler.AddComment)
		sentinelGroup.GET("/tasks/:id/comments", handler.GetComments)

		// Time Logging
		sentinelGroup.POST("/tasks/:id/time-logs", handler.LogTime)
		sentinelGroup.GET("/tasks/:id/time-logs", handler.GetTimeLogs)

		// Appeal System
		sentinelGroup.POST("/submissions/:id/appeal", handler.SubmitAppeal)
		sentinelGroup.POST("/appeals/:id/resolve", handler.ResolveAppeal)

		// Sprint Management
		sentinelGroup.POST("/sprints", handler.CreateSprint)
		sentinelGroup.GET("/sprints", handler.GetSprintsByProject) // ?project_id=UUID
		sentinelGroup.PATCH("/sprints/:id", handler.UpdateSprint)
		sentinelGroup.DELETE("/sprints/:id", handler.DeleteSprint)
		sentinelGroup.POST("/sprints/:id/start", handler.StartSprint)
		sentinelGroup.POST("/sprints/:id/complete", handler.CompleteSprint)
		sentinelGroup.POST("/sprints/:id/reopen", handler.ReopenSprint)
		sentinelGroup.POST("/sprints/:id/tasks", handler.AddTasksToSprint)

		// Milestone Management
		sentinelGroup.POST("/milestones", handler.CreateMilestone)
		sentinelGroup.GET("/milestones", handler.GetMilestonesByProject) // ?project_id=UUID
		sentinelGroup.PATCH("/milestones/:id", handler.UpdateMilestone)
		sentinelGroup.DELETE("/milestones/:id", handler.DeleteMilestone)

		// Project Analytics
		sentinelGroup.GET("/projects/:id/analytics", handler.GetProjectAnalytics)

		// Bulk Operations
		sentinelGroup.PATCH("/tasks/bulk-status", handler.BulkUpdateTaskStatus)

		// Google Slides Import
		sentinelGroup.POST("/import/google-slides/preview", handler.PreviewGoogleSlides)
		sentinelGroup.POST("/import/google-slides", handler.ImportGoogleSlides)

		// Epic Management (Hierarchy Dimension 1)
		sentinelGroup.POST("/epics", handler.CreateEpic)
		sentinelGroup.GET("/epics", handler.GetEpicsByProject)   // ?project_id=UUID
		sentinelGroup.PATCH("/epics/:id", handler.UpdateEpic)
		sentinelGroup.DELETE("/epics/:id", handler.DeleteEpic)

		// Timeline Views (Matrix Dimension)
		sentinelGroup.GET("/projects/:id/timeline/epic-view", handler.GetEpicTimeline)
		sentinelGroup.GET("/projects/:id/timeline/sprint-view", handler.GetSprintTimeline)
	}

	// Admin/CEO Configuration Management
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/config", handler.GetSystemConfig)        // Get current AI config
		adminGroup.PUT("/config", handler.UpdateSystemConfig)     // Update AI config (CEO only)
		adminGroup.GET("/models", handler.GetAvailableModels)     // Get available Gemini models
	}
}
