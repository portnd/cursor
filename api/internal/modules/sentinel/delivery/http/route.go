package http

import (
	"github.com/gin-gonic/gin"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// RegisterRoutes registers all Sentinel module routes (handlers are split across sentinel_handler*.go in this package).
func RegisterRoutes(
	router *gin.RouterGroup,
	usecase domain.SentinelUsecase,
	projectFinanceUsecase domain.ProjectFinanceUsecase,
	googleAPIKey, canvaAccessToken string,
	authRepo authDomain.Repository,
	attendanceRepo attendanceDomain.AttendanceRepository,
	sentinelRepo domain.SentinelRepository,
	discordSvc domain.DiscordNotifier,
) {
	handler := NewSentinelHandler(usecase, googleAPIKey, canvaAccessToken)
	financeHandler := NewProjectFinanceHandler(projectFinanceUsecase)
	b2bHandler := NewB2BHandler(usecase)
	backupHandler := NewProjectBackupHandler(usecase)
	discordTestHandler := NewDiscordTestHandler(usecase, authRepo, attendanceRepo, sentinelRepo, discordSvc)

	sentinelGroup := router.Group("/sentinel")
	// sentinelGroup.Use(middleware.Auth()) // Uncomment if you have auth middleware
	{
		// Project Management
		sentinelGroup.POST("/projects", handler.CreateProject)
		sentinelGroup.POST("/projects/import-backup", backupHandler.ImportProjectFromBackup) // must be before /:id
		sentinelGroup.GET("/projects", handler.GetProjects)
		sentinelGroup.GET("/projects/:id/details", handler.GetProjectDetails) // Combined payload (1 round-trip) — must be before /:id
		sentinelGroup.GET("/projects/:id/tasks", handler.GetProjectTasksPage) // Page 2+ loading for tasks via cursor/offset
		sentinelGroup.GET("/projects/:id", handler.GetProjectByID)
		sentinelGroup.PATCH("/projects/:id", handler.UpdateProject)
		sentinelGroup.DELETE("/projects/:id", handler.DeleteProject)
		sentinelGroup.POST("/projects/:id/ai-schedule", handler.ScheduleProjectWithAI) // AI Agent: estimate + schedule existing tasks (CEO / Product Owner)
		sentinelGroup.POST("/projects/:id/ai-plan", handler.GenerateProjectPlan)       // (Optional) generate new epics/sprints/milestones/tasks from scratch
		sentinelGroup.POST("/projects/:id/clear-plan", handler.ClearProjectPlan)       // Clear plan: remove all tasks, sprints, milestones, epics (CEO / Product Owner)
		sentinelGroup.PATCH("/projects/:id/assign-team", handler.AssignProjectTeam)    // Squad Model: assign project to team (CEO only)
		sentinelGroup.PATCH("/projects/:id/pm-owners", handler.AssignProjectPmOwners)  // No-squads mode: CEO/MANAGER assigns Product Owners per project (path kept for compatibility)

		// Task Management
		sentinelGroup.POST("/tasks", handler.CreateTask)                                    // Create new task (CEO / Product Owner)
		sentinelGroup.GET("/tasks/my", handler.GetMyTasks)                                  // Get my assigned tasks (engineer)
		sentinelGroup.GET("/tasks/my-global-active", handler.GetGlobalActiveTasks)          // Active sprints; CEO/MANAGER all projects; Product Owner / engineer per rules
		sentinelGroup.GET("/tasks/team-active", handler.GetTeamActiveTasks)                 // All ACTIVE-sprint TASK/BUG items in caller's team (cross-engineer Quick Log Time)
		sentinelGroup.GET("/tasks/features", handler.GetActiveFeatures)                     // FEATURE items for Product Owner/CEO Roadmap Board (must be before /:id)
		sentinelGroup.GET("/tasks/unassigned", handler.GetUnassignedTasks)                  // Get unassigned tasks (CEO / Product Owner)
		sentinelGroup.GET("/tasks/approvals", handler.GetApprovals)                         // Get approvals inbox (CEO / Product Owner only)
		sentinelGroup.GET("/tasks/gantt", handler.GetGantt)                                 // Get all tasks + dependencies for Gantt chart
		sentinelGroup.GET("/tasks/ready-for-test", handler.GetTasksReadyForTest)            // Continuous UAT queue (must be before /:id)
		sentinelGroup.GET("/tasks", handler.GetAllTasks)                                    // Get all tasks (ADMIN / Product Owner overview)
		sentinelGroup.GET("/tasks/:id/activity", handler.GetTaskActivity)                   // Immutable task lifecycle timeline
		sentinelGroup.GET("/tasks/:id/summary", handler.GetTaskSummary)                     // Lightweight task summary for initial render
		sentinelGroup.GET("/tasks/:id/detail", handler.GetTaskDetail)                       // Full task payload with rich content for on-demand loading
		sentinelGroup.GET("/tasks/:id", handler.GetTaskDetail)                              // Full task payload (legacy alias → detail)
		sentinelGroup.PATCH("/tasks/:id", handler.UpdateTask)                               // Update task (Creator or CEO only, triggers AI re-estimation)
		sentinelGroup.PATCH("/tasks/:id/slide-resources", handler.UpdateTaskSlideResources) // Update task resource_urls (slide images/annotations)
		sentinelGroup.POST("/tasks/:id/estimate", handler.EstimateTask)                     // AI estimate time (Creator / CEO / Product Owner)
		sentinelGroup.DELETE("/tasks/:id", handler.DeleteTask)                              // Delete task (Creator or CEO only)
		sentinelGroup.POST("/tasks/:id/split", handler.SplitTask)                           // Split task into N sub-tasks (Product Owner/CEO/Creator)
		sentinelGroup.POST("/tasks/:id/assign", handler.AssignTask)                         // Assign task (Product Owner/CEO/Manager; parent assignee for subtasks)
		sentinelGroup.POST("/tasks/:id/submit", handler.SubmitWork)                         // Handover: engineer submits PR/Commit URL for review
		sentinelGroup.POST("/tasks/:id/submit-uat", handler.SubmitUAT)                      // UAT: engineer submits staging URL + release notes for FEATURE review
		sentinelGroup.POST("/tasks/:id/negotiate", handler.NegotiateTime)                   // Engineer negotiates AI time estimate
		sentinelGroup.POST("/tasks/:id/approve", handler.ApproveTask)                       // Approve task after review (Product Owner / CEO / Manager)
		sentinelGroup.POST("/tasks/:id/reject", handler.RejectTask)                         // Reject task and return to IN_PROGRESS (Product Owner/CEO/MANAGER)

		// Continuous UAT: sub-task testing lane
		sentinelGroup.POST("/tasks/:id/ready-for-test", handler.MarkReadyForTest)
		sentinelGroup.POST("/tasks/:id/pm-approve-sub", handler.PMApproveSubTask) // Product Owner: READY_FOR_TEST → READY_FOR_UAT (with test evidence); path kept for compatibility
		sentinelGroup.POST("/tasks/:id/approve-sub", handler.ApproveSubTask)      // CEO / Manager: READY_FOR_UAT → COMPLETED (final approval)
		sentinelGroup.POST("/tasks/:id/reject-sub", handler.RejectSubTask)
		sentinelGroup.GET("/tasks/ceo-approval-queue", handler.GetTasksReadyForCEOApproval) // CEO / Manager: tasks awaiting final approval

		// Task Dependencies (Gantt links)
		sentinelGroup.POST("/tasks/dependencies", handler.CreateDependency)
		sentinelGroup.DELETE("/tasks/dependencies/:id", handler.DeleteDependency)

		// Task Comments
		sentinelGroup.POST("/tasks/:id/comments", handler.AddComment)
		sentinelGroup.GET("/tasks/:id/comments", handler.GetComments)
		sentinelGroup.PATCH("/comments/:commentId", handler.EditComment)
		sentinelGroup.DELETE("/comments/:commentId", handler.DeleteComment)

		// Time Logging
		sentinelGroup.POST("/tasks/:id/time-logs", handler.LogTime)
		sentinelGroup.GET("/tasks/:id/time-logs", handler.GetTimeLogs)
		sentinelGroup.PATCH("/time-logs/:logId", handler.EditTimeLog)
		sentinelGroup.DELETE("/time-logs/:logId", handler.DeleteTimeLog)
		sentinelGroup.GET("/users/me/time-logs", handler.GetMyDailyTimeLogs)

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
		sentinelGroup.GET("/projects/finance/capital", financeHandler.GetProjectCapitals) // bulk: ?project_ids=uuid1,uuid2
		sentinelGroup.GET("/projects/:id/finance/capital", financeHandler.GetProjectCapital)
		sentinelGroup.POST("/projects/:id/finance/inject", financeHandler.InjectProjectCapital)
		sentinelGroup.PUT("/projects/:id/finance/capital", financeHandler.EditProjectCapital)
		sentinelGroup.POST("/projects/:id/finance/close-cycle", financeHandler.CloseProjectCycleAndPayout)
		sentinelGroup.DELETE("/projects/:id/finance/transactions/:txID", financeHandler.DeleteProjectTransaction)

		// Bulk Operations
		sentinelGroup.PATCH("/tasks/bulk-status", handler.BulkUpdateTaskStatus)
		sentinelGroup.POST("/time-logs/bulk", handler.BulkLogTime)

		// Google Slides Import
		sentinelGroup.POST("/import/google-slides/preview", handler.PreviewGoogleSlides)
		sentinelGroup.POST("/import/google-slides", handler.ImportGoogleSlides)

		// Google Sheets Import
		sentinelGroup.POST("/import/google-sheets/preview", handler.PreviewGoogleSheets)
		sentinelGroup.POST("/import/google-sheets", handler.ImportGoogleSheets)

		// Canva Import
		sentinelGroup.POST("/import/canva/preview", handler.PreviewCanva)
		sentinelGroup.POST("/import/canva", handler.ImportCanva)

		// PPTX File Upload Import (multipart)
		sentinelGroup.POST("/import/pptx/preview", handler.PreviewPPTX)
		sentinelGroup.POST("/import/pptx", handler.ImportPPTX)

		// Epic Management (Hierarchy Dimension 1)
		sentinelGroup.POST("/epics", handler.CreateEpic)
		sentinelGroup.GET("/epics", handler.GetEpicsByProject) // ?project_id=UUID
		sentinelGroup.PATCH("/epics/:id", handler.UpdateEpic)
		sentinelGroup.DELETE("/epics/:id", handler.DeleteEpic)

		// Timeline Views (Matrix Dimension)
		sentinelGroup.GET("/projects/:id/timeline/epic-view", handler.GetEpicTimeline)
		sentinelGroup.GET("/projects/:id/timeline/sprint-view", handler.GetSprintTimeline)
		sentinelGroup.GET("/projects/:id/timeline/export-pdf", handler.ExportTimelinePDF)

		// Internal B2B Outsource Requests
		sentinelGroup.POST("/b2b/requests", b2bHandler.CreateB2BRequest)            // Product Owner/CEO creates outsource request
		sentinelGroup.GET("/b2b/requests", b2bHandler.GetB2BRequests)               // ?direction=inbound|outbound
		sentinelGroup.PATCH("/b2b/requests/:id", b2bHandler.UpdateB2BRequest)       // Counter-offer or Reject
		sentinelGroup.POST("/b2b/requests/:id/accept", b2bHandler.AcceptB2BRequest) // Accept → creates Task

		// Project Backups (Disaster Recovery)
		sentinelGroup.GET("/projects/:id/backups", backupHandler.ListProjectBackups)
		sentinelGroup.POST("/projects/:id/backups", backupHandler.CreateProjectBackup)
		sentinelGroup.GET("/projects/:id/backups/:backupId/payload", backupHandler.GetProjectBackupPayload)
		sentinelGroup.POST("/projects/:id/backups/:backupId/restore", backupHandler.RestoreProjectBackup)
		sentinelGroup.DELETE("/projects/:id/backups/:backupId", backupHandler.DeleteProjectBackup)

		// Komgrip: project-less personal/misc tasks (all employees)
		sentinelGroup.GET("/komgrip/tasks", handler.GetKomgripTasks)
		sentinelGroup.POST("/komgrip/tasks", handler.CreateKomgripTask)
		sentinelGroup.PATCH("/komgrip/tasks/:id/status", handler.UpdateKomgripTaskStatus)
		sentinelGroup.DELETE("/komgrip/tasks/:id", handler.DeleteKomgripTask)
	}

	// Admin/CEO Configuration Management
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/config", handler.GetSystemConfig)    // Get current AI config
		adminGroup.PUT("/config", handler.UpdateSystemConfig) // Update AI config (CEO only)
		adminGroup.GET("/models", handler.GetAvailableModels) // Get available GLM models
		adminGroup.GET("/ai-usage", handler.GetAIUsage)       // Approximate AI quota usage

		// Discord Notification Testing (CEO only)
		adminGroup.POST("/discord/test-missing-log", discordTestHandler.TestMissingLogNotification) // Test missing log notification
		adminGroup.POST("/discord/test-leave", discordTestHandler.TestLeaveNotification)             // Test leave notification
	}
}
