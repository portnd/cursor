package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// RegisterRoutes registers all Sentinel module routes
func RegisterRoutes(router *gin.RouterGroup, usecase domain.SentinelUsecase, projectFinanceUsecase domain.ProjectFinanceUsecase, googleAPIKey ...string) {
	key := ""
	if len(googleAPIKey) > 0 {
		key = googleAPIKey[0]
	}
	handler := NewSentinelHandler(usecase, key)
	financeHandler := NewProjectFinanceHandler(projectFinanceUsecase)
	b2bHandler := NewB2BHandler(usecase)
	backupHandler := NewProjectBackupHandler(usecase)

	sentinelGroup := router.Group("/sentinel")
	// sentinelGroup.Use(middleware.Auth()) // Uncomment if you have auth middleware
	{
		// Project Management
		sentinelGroup.POST("/projects", handler.CreateProject)
		sentinelGroup.POST("/projects/import-backup", backupHandler.ImportProjectFromBackup) // must be before /:id
		sentinelGroup.GET("/projects", handler.GetProjects)
		sentinelGroup.GET("/projects/:id", handler.GetProjectByID)
		sentinelGroup.PATCH("/projects/:id", handler.UpdateProject)
		sentinelGroup.DELETE("/projects/:id", handler.DeleteProject)
		sentinelGroup.POST("/projects/:id/ai-schedule", handler.ScheduleProjectWithAI) // AI Agent: estimate + schedule existing tasks (CEO/PM)
		sentinelGroup.POST("/projects/:id/ai-plan", handler.GenerateProjectPlan)         // (Optional) generate new epics/sprints/milestones/tasks from scratch
		sentinelGroup.POST("/projects/:id/clear-plan", handler.ClearProjectPlan)        // Clear plan: remove all tasks, sprints, milestones, epics (CEO/PM)
		sentinelGroup.PATCH("/projects/:id/assign-team", handler.AssignProjectTeam)     // Squad Model: assign project to team (CEO only)

		// Task Management
		sentinelGroup.POST("/tasks", handler.CreateTask)                   // Create new task (CEO/PM)
		sentinelGroup.GET("/tasks/my", handler.GetMyTasks)                 // Get my assigned tasks (DEV)
		sentinelGroup.GET("/tasks/my-global-active", handler.GetGlobalActiveTasks) // Get TASK/BUG in active sprints across ALL projects (DEV Board)
		sentinelGroup.GET("/tasks/team-active", handler.GetTeamActiveTasks)         // All ACTIVE-sprint TASK/BUG items in caller's team (cross-dev Quick Log Time)
		sentinelGroup.GET("/tasks/features", handler.GetActiveFeatures)             // FEATURE items for PM/CEO Roadmap Board (must be before /:id)
		sentinelGroup.GET("/tasks/unassigned", handler.GetUnassignedTasks) // Get unassigned tasks (CEO/PM)
		sentinelGroup.GET("/tasks/approvals", handler.GetApprovals)        // Get approvals inbox (CEO/PM only)
		sentinelGroup.GET("/tasks/gantt", handler.GetGantt)                // Get all tasks + dependencies for Gantt chart
		sentinelGroup.GET("/tasks/ready-for-test", handler.GetTasksReadyForTest) // Continuous UAT queue (must be before /:id)
		sentinelGroup.GET("/tasks", handler.GetAllTasks)                    // Get all tasks (ADMIN/PM overview)
		sentinelGroup.GET("/tasks/:id", handler.GetTaskByID)                // Get single task with submission history
		sentinelGroup.PATCH("/tasks/:id", handler.UpdateTask)                      // Update task (Creator or CEO only, triggers AI re-estimation)
		sentinelGroup.PATCH("/tasks/:id/slide-resources", handler.UpdateTaskSlideResources) // Update task resource_urls (slide images/annotations)
		sentinelGroup.POST("/tasks/:id/estimate", handler.EstimateTask)                // AI estimate time (Creator/CEO/PM)
		sentinelGroup.DELETE("/tasks/:id", handler.DeleteTask)                     // Delete task (Creator or CEO only)
		sentinelGroup.POST("/tasks/:id/split", handler.SplitTask)                   // Split task into N sub-tasks (PM/CEO/Creator)
		sentinelGroup.POST("/tasks/:id/assign", handler.AssignTask)        // Assign task to developer (PM)
		sentinelGroup.POST("/tasks/:id/submit", handler.SubmitWork)        // Handover: Dev submits PR/Commit URL for review
		sentinelGroup.POST("/tasks/:id/submit-uat", handler.SubmitUAT)     // UAT: Dev submits staging URL + release notes for FEATURE review
		sentinelGroup.POST("/tasks/:id/negotiate", handler.NegotiateTime)  // Developer negotiates AI time estimate
		sentinelGroup.POST("/tasks/:id/approve", handler.ApproveTask)      // Approve task after review (PM/CEO only)
		sentinelGroup.POST("/tasks/:id/reject", handler.RejectTask)        // Reject task and return to IN_PROGRESS (PM/CEO/MANAGER)

		// Continuous UAT: sub-task testing lane (READY_FOR_TEST)
		sentinelGroup.POST("/tasks/:id/ready-for-test", handler.MarkReadyForTest)
		sentinelGroup.POST("/tasks/:id/approve-sub", handler.ApproveSubTask)
		sentinelGroup.POST("/tasks/:id/reject-sub", handler.RejectSubTask)

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

		// Project Finance (Internal VC — per-project capital)
		sentinelGroup.GET("/projects/:id/finance/capital", financeHandler.GetProjectCapital)
		sentinelGroup.POST("/projects/:id/finance/inject", financeHandler.InjectProjectCapital)
		sentinelGroup.PUT("/projects/:id/finance/capital", financeHandler.EditProjectCapital)
		sentinelGroup.POST("/projects/:id/finance/close-cycle", financeHandler.CloseProjectCycleAndPayout)
		sentinelGroup.DELETE("/projects/:id/finance/transactions/:txID", financeHandler.DeleteProjectTransaction)

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
		sentinelGroup.GET("/projects/:id/timeline/export-pdf", handler.ExportTimelinePDF)

		// Internal B2B Outsource Requests
		sentinelGroup.POST("/b2b/requests", b2bHandler.CreateB2BRequest)                 // PM/CEO creates outsource request
		sentinelGroup.GET("/b2b/requests", b2bHandler.GetB2BRequests)                    // ?direction=inbound|outbound
		sentinelGroup.PATCH("/b2b/requests/:id", b2bHandler.UpdateB2BRequest)            // Counter-offer or Reject
		sentinelGroup.POST("/b2b/requests/:id/accept", b2bHandler.AcceptB2BRequest)      // Accept → creates Task

		// Project Backups (Disaster Recovery)
		sentinelGroup.GET("/projects/:id/backups", backupHandler.ListProjectBackups)
		sentinelGroup.POST("/projects/:id/backups", backupHandler.CreateProjectBackup)
		sentinelGroup.GET("/projects/:id/backups/:backupId/payload", backupHandler.GetProjectBackupPayload)
		sentinelGroup.POST("/projects/:id/backups/:backupId/restore", backupHandler.RestoreProjectBackup)
		sentinelGroup.DELETE("/projects/:id/backups/:backupId", backupHandler.DeleteProjectBackup)
	}

	// Admin/CEO Configuration Management
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/config", handler.GetSystemConfig)        // Get current AI config
		adminGroup.PUT("/config", handler.UpdateSystemConfig)     // Update AI config (CEO only)
		adminGroup.GET("/models", handler.GetAvailableModels)     // Get available Gemini models
		adminGroup.GET("/ai-usage", handler.GetAIUsage)           // Approximate Gemini quota usage
	}
}
