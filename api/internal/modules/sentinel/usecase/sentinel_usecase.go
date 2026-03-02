package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type sentinelUsecase struct {
	repo       domain.SentinelRepository
	aiService  domain.AIService
	authRepo   authDomain.Repository // 👈 Add auth repo for role validation
	timeout    time.Duration
}

// Update Constructor
func NewSentinelUsecase(repo domain.SentinelRepository, aiService domain.AIService, authRepo authDomain.Repository) domain.SentinelUsecase {
	return &sentinelUsecase{
		repo:      repo,
		aiService: aiService,
		authRepo:  authRepo, // 👈 Inject auth repo
		timeout:   time.Second * 10,
	}
}

func (u *sentinelUsecase) CreateTask(title, desc string, creatorID uint, dueDate *time.Time) (*domain.Task, error) {
	// 1. Ask AI for Estimation 🧠
	minutes, reasoning, err := u.aiService.EstimateEffort(title, desc)
	if err != nil {
		return nil, fmt.Errorf("AI estimation failed: %w", err)
	}

	// Log AI reasoning for debugging
	fmt.Printf("📊 AI Reasoning: %s\n", reasoning)

	// 2. Prepare the entity
	task := &domain.Task{
		ID:                 uuid.New(),
		Title:              title,
		Description:        desc,
		CreatedBy:          &creatorID,
		Status:             "PENDING",
		AIEstimatedMinutes: minutes, // 👈 ใช้ค่าจาก AI
		DueAt:              dueDate, // 👈 Deadline set by PM/CEO
		// เราอาจจะเก็บ reasoning ไว้ใน Description หรือ Field ใหม่ก็ได้ (ในอนาคต)
		// ตอนนี้แปะไว้ท้าย Description ก่อนเพื่อให้เห็นว่า AI คิดยังไง
		// Description: desc + "\n\n[AI Reasoning]: " + reasoning,
	}

	// 3. Persist to DB
	if err := u.repo.CreateTask(task); err != nil {
		return nil, err
	}

	return task, nil
}

// AssignTask assigns a developer to a task
func (u *sentinelUsecase) AssignTask(taskID uuid.UUID, devID uint) error {
	// 1. Validate if task exists
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	// 2. Update assignment (In a real app, we might check if Dev is overloaded here)
	task.AssignedTo = &devID
	task.Status = "IN_PROGRESS"

	// 3. ⏰ Start Time Tracking: Set StartedAt = NOW() (Assuming assignment starts work)
	now := time.Now()
	task.StartedAt = &now
	fmt.Printf("⏰ Time Tracking: Task %s started at %s\n", task.ID, now.Format(time.RFC3339))

	// 4. Persist changes
	if err := u.repo.UpdateTask(task); err != nil {
		return err
	}

	return nil
}

// SubmitWork handles code submission with AI Code Review
func (u *sentinelUsecase) SubmitWork(taskID uuid.UUID, devID uint, commitHash, diff string) (*domain.Submission, error) {
	// Initialize submission
	sub := &domain.Submission{
		ID:         uuid.New(),
		TaskID:     taskID,
		DevID:      devID,
		CommitHash: commitHash,
		Diff:       diff,       // Store diff for appeal analysis
		AIVerdict:  "PENDING", // Default
	}

	// If diff is provided, run AI Code Review
	if diff != "" {
		verdict, score, feedback, err := u.aiService.ReviewCode(diff)
		if err != nil {
			// If AI fails, mark as PENDING but don't block submission
			fmt.Printf("⚠️  AI Code Review failed: %v. Marking as PENDING.\n", err)
			sub.AIVerdict = "PENDING"
			sub.AIScore = 0
			
			// Properly encode error message as JSON to avoid PostgreSQL syntax errors
			errorMap := map[string]string{"error": err.Error()}
			errorJSON, marshalErr := json.Marshal(errorMap)
			if marshalErr != nil {
				// Fallback to safe JSON
				sub.AIFeedback = []byte(`{"error": "AI review service unavailable"}`)
			} else {
				sub.AIFeedback = errorJSON
			}
		} else {
			// AI Review successful
			sub.AIVerdict = verdict
			sub.AIScore = score
			
			// Store feedback as JSON (properly encode to avoid escaping issues)
			feedbackMap := map[string]string{"feedback": feedback}
			feedbackJSON, err := json.Marshal(feedbackMap)
			if err != nil {
				// Fallback to simple JSON
				sub.AIFeedback = []byte(`{"feedback": "Failed to encode feedback"}`)
			} else {
				sub.AIFeedback = feedbackJSON
			}
			
			fmt.Printf("🤖 AI Review Complete: %s (%d/100)\n", verdict, score)
		}
	} else {
		// No diff provided, skip review
		fmt.Printf("⚠️  No diff provided. Skipping AI Code Review.\n")
	}

	// Persist submission to DB
	if err := u.repo.CreateSubmission(sub); err != nil {
		return nil, err
	}

	// 🚦 Human Quality Gate: If AI PASS, move to REVIEW_PENDING (not COMPLETED)
	if sub.AIVerdict == "PASS" {
		task, err := u.repo.GetTaskByID(taskID)
		if err != nil {
			fmt.Printf("⚠️  Failed to get task for review queue: %v\n", err)
		} else {
			task.Status = "REVIEW_PENDING" // 🔒 Requires PM/CEO approval
			// Do NOT set CompletedAt yet - human approval required
			
			if err := u.repo.UpdateTask(task); err != nil {
				fmt.Printf("⚠️  Failed to update task status to REVIEW_PENDING: %v\n", err)
			} else {
				fmt.Printf("🚦 Task %s moved to REVIEW_PENDING - awaiting PM/CEO approval\n", task.ID)
				fmt.Printf("📋 AI Review: PASS (%d/100) - Ready for human verification\n", sub.AIScore)
			}
		}
	}

	return sub, nil
}

// GetTaskByID retrieves a single task with full submission history
func (u *sentinelUsecase) GetTaskByID(taskID uuid.UUID) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// GetMyTasks retrieves all tasks assigned to a user
func (u *sentinelUsecase) GetMyTasks(userID uint) ([]domain.Task, error) {
	tasks, err := u.repo.GetTasksByAssignee(userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetUnassignedTasks retrieves all tasks that are not assigned to anyone
func (u *sentinelUsecase) GetUnassignedTasks() ([]domain.Task, error) {
	tasks, err := u.repo.GetUnassignedTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetAllTasks retrieves all tasks in the system (for ADMIN/PM view)
func (u *sentinelUsecase) GetAllTasks() ([]domain.Task, error) {
	tasks, err := u.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetPendingApprovals returns tasks requiring PM/CEO attention
// Includes: Time negotiations (PENDING) and Appeals (PENDING)
// Access: CEO and PM roles only
func (u *sentinelUsecase) GetPendingApprovals(userRole string) ([]domain.Task, error) {
	// 🔒 ROLE VALIDATION: Only CEO and PM can view approvals inbox
	if userRole != "CEO" && userRole != "PM" {
		return nil, fmt.Errorf("access denied: only CEO and PM can view approvals inbox")
	}

	// Fetch tasks requiring approval
	tasks, err := u.repo.GetTasksRequiringApproval()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending approvals: %w", err)
	}

	return tasks, nil
}

// SubmitAppeal allows a developer to appeal an AI FAIL verdict
func (u *sentinelUsecase) SubmitAppeal(submissionID uuid.UUID, devID uint, reason string) (*domain.Appeal, error) {
	// 1. Validate submission exists
	submission, err := u.repo.GetSubmissionByID(submissionID)
	if err != nil {
		return nil, fmt.Errorf("submission not found: %w", err)
	}

	// 2. Ensure only the developer who submitted can appeal
	if submission.DevID != devID {
		return nil, errors.New("unauthorized: only the developer who submitted can appeal")
	}

	// 3. Check if appeal already exists
	existingAppeal, err := u.repo.GetAppealBySubmissionID(submissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing appeal: %w", err)
	}
	if existingAppeal != nil {
		return nil, errors.New("appeal already exists for this submission")
	}

	// 4. Validate that submission is a FAIL (only FAILs can be appealed)
	if submission.AIVerdict != "FAIL" {
		return nil, errors.New("can only appeal FAIL verdicts")
	}

	// 5. 🤖 AI ADVISORY: Analyze the appeal validity
	// Extract original feedback from AIFeedback JSON
	originalFeedback := ""
	if len(submission.AIFeedback) > 0 {
		var feedbackMap map[string]interface{}
		if err := json.Unmarshal(submission.AIFeedback, &feedbackMap); err == nil {
			if feedback, ok := feedbackMap["feedback"].(string); ok {
				originalFeedback = feedback
			}
		}
		// Fallback: use raw JSON string
		if originalFeedback == "" {
			originalFeedback = string(submission.AIFeedback)
		}
	}

	// Call AI to analyze the appeal
	recommendation, confidence, reasoning, err := u.aiService.AnalyzeAppeal(
		submission.Diff,      // Code diff
		originalFeedback,     // Original AI failure reason
		reason,               // Developer's appeal reason
	)
	
	// If AI analysis fails, continue with default values (don't block appeal)
	if err != nil {
		fmt.Printf("⚠️  AI Appeal Analysis failed: %v. Proceeding with default values.\n", err)
		recommendation = "UPHOLD" // Conservative default
		confidence = 0
		reasoning = "⚠️ AI ไม่สามารถวิเคราะห์ได้ในขณะนี้ - กรุณาตรวจสอบด้วยตนเองอย่างละเอียด"
	}

	// 6. Create appeal with AI advisory
	appeal := &domain.Appeal{
		ID:           uuid.New(),
		SubmissionID: submissionID,
		DeveloperID:  devID,
		Reason:       reason,
		Status:       "PENDING",
		
		// AI Advisory System
		AIRecommendation: recommendation, // OVERTURN or UPHOLD
		AIConfidence:     confidence,     // 0-100
		AIReasoning:      reasoning,      // Advice for CEO/PM
	}

	if err := u.repo.CreateAppeal(appeal); err != nil {
		return nil, fmt.Errorf("failed to create appeal: %w", err)
	}

	fmt.Printf("⚖️  Appeal created with AI Advisory:\n")
	fmt.Printf("   • Recommendation: %s\n", appeal.AIRecommendation)
	fmt.Printf("   • Confidence: %d%%\n", appeal.AIConfidence)
	fmt.Printf("   • Reasoning: %s\n", appeal.AIReasoning)

	return appeal, nil
}

// ResolveAppeal allows PM/CEO to approve or reject an appeal
func (u *sentinelUsecase) ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error {
	// 1. Validate status
	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("status must be APPROVED or REJECTED")
	}

	// 2. 🔒 ROLE VALIDATION: Only CEO/PM can resolve appeals
	resolver, err := u.authRepo.FindByID(resolverID)
	if err != nil {
		return fmt.Errorf("unauthorized: resolver user not found: %w", err)
	}

	if resolver.Role != authDomain.RoleCEO && resolver.Role != authDomain.RolePM {
		return fmt.Errorf("forbidden: only CEO or PM can resolve appeals (current role: %s)", resolver.Role)
	}

	// 3. Get appeal
	appeal, err := u.repo.GetAppealByID(appealID)
	if err != nil {
		return fmt.Errorf("appeal not found: %w", err)
	}

	// 4. Check if already resolved
	if appeal.Status != "PENDING" {
		return fmt.Errorf("appeal already resolved with status: %s", appeal.Status)
	}

	// 5. Update appeal
	appeal.Status = status
	appeal.ResolverID = &resolverID
	appeal.ResolverNote = note

	if err := u.repo.UpdateAppeal(appeal); err != nil {
		return fmt.Errorf("failed to update appeal: %w", err)
	}

	// 6. If APPROVED, override the submission verdict AND auto-complete task
	if status == "APPROVED" {
		submission, err := u.repo.GetSubmissionByID(appeal.SubmissionID)
		if err != nil {
			return fmt.Errorf("failed to get submission: %w", err)
		}

		// Override submission verdict
		submission.AIVerdict = "PASS"
		submission.IsOverridden = true

		if err := u.repo.UpdateSubmission(submission); err != nil {
			return fmt.Errorf("failed to override submission: %w", err)
		}

		fmt.Printf("✅ Appeal APPROVED: Submission %s overridden to PASS\n", submission.ID)

		// 🎯 AUTO-COMPLETE TASK
		task, err := u.repo.GetTaskByID(submission.TaskID)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to get task for auto-completion: %v\n", err)
		} else if task.Status != "COMPLETED" {
			task.Status = "COMPLETED"
			now := time.Now()
			task.CompletedAt = &now

			// 🔧 FIX: Ensure started_at is set (required by DB constraint)
			if task.StartedAt == nil {
				task.StartedAt = &now
				fmt.Printf("⚠️  Task had no started_at, setting to now\n")
			}

			if err := u.repo.UpdateTask(task); err != nil {
				fmt.Printf("⚠️  Warning: Failed to auto-complete task: %v\n", err)
			} else {
				fmt.Printf("🎉 Task %s auto-completed via appeal approval\n", task.ID)

				// Calculate actual time taken
				if task.StartedAt != nil {
					duration := now.Sub(*task.StartedAt)
					fmt.Printf("📊 Actual Time: %.2f hours (AI Estimated: %.2f hours)\n",
						duration.Hours(), float64(task.AIEstimatedMinutes)/60.0)
				}
			}
		}
	} else {
		fmt.Printf("❌ Appeal REJECTED by %s (%s): Submission remains FAIL\n", resolver.Email, resolver.Role)
	}

	return nil
}

// NegotiateTime allows a developer to negotiate/dispute the AI-estimated time
func (u *sentinelUsecase) NegotiateTime(taskID uuid.UUID, devID uint, minutes int, reason string) error {
	// 1. Get the task
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2. Validate: Only the assigned developer (or creator for unassigned tasks) can negotiate
	if task.AssignedTo != nil {
		// Task is assigned - only assignee can negotiate
		if *task.AssignedTo != devID {
			return errors.New("unauthorized: only the assigned developer can negotiate time")
		}
	} else {
		// Task is unassigned - only creator can negotiate (before assignment)
		if task.CreatedBy == nil || *task.CreatedBy != devID {
			return errors.New("unauthorized: only the task creator can negotiate time for unassigned tasks")
		}
	}

	// 3. Validate proposed minutes
	if minutes <= 0 {
		return errors.New("proposed minutes must be greater than 0")
	}

	if minutes <= task.AIEstimatedMinutes {
		return errors.New("proposed time must be greater than AI estimate (why negotiate if you need less time?)")
	}

	// 4. Validate reason
	if len(reason) < 20 {
		return errors.New("negotiation reason must be at least 20 characters")
	}

	// 5. Check if already negotiating
	if task.NegotiationStatus == "PENDING" {
		return errors.New("time negotiation already pending review")
	}

	// 🤖 AI Analysis: Analyze time negotiation request
	recommendation, confidence, reasoning, aiErr := u.aiService.AnalyzeTimeNegotiation(
		task.Title,
		task.Description,
		task.AIEstimatedMinutes,
		minutes,
		reason,
	)

	// If AI analysis fails, proceed with defaults but log the error
	if aiErr != nil {
		fmt.Printf("⚠️  AI Time Negotiation Analysis Failed: %v\n", aiErr)
		recommendation = "REJECT"
		confidence = 0
		reasoning = "ระบบ AI ขัดข้อง - กรุณาใช้ดุลยพินิจในการตัดสิน"
	}

	fmt.Printf("⏱️  AI Recommendation: %s (%d%% confidence)\n", recommendation, confidence)
	fmt.Printf("   Reasoning: %s\n", reasoning)

	// 6. Update task with negotiation data + AI advisory
	task.NegotiationStatus = "PENDING"
	task.ProposedMinutes = minutes
	task.NegotiationReason = reason
	task.NegotiationAIRecommendation = recommendation
	task.NegotiationAIConfidence = confidence
	task.NegotiationAIReasoning = reasoning

	if err := u.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to submit time negotiation: %w", err)
	}

	fmt.Printf("⏰ Time Negotiation Submitted: Task %s | AI: %d min → Proposed: %d min\n",
		task.ID, task.AIEstimatedMinutes, minutes)

	return nil
}

// UpdateTask updates a task with access control and AI re-estimation
// Only the Creator OR CEO can update a task
// If title/description changes, AI will re-estimate automatically
func (u *sentinelUsecase) UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string) (*domain.Task, error) {
	// 1️⃣ Fetch the existing task
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// 2️⃣ ACCESS CONTROL: Only Creator or CEO can update
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"

	if !isCreator && !isCEO {
		return nil, fmt.Errorf("unauthorized: only the task creator or CEO can update this task")
	}

	// 3️⃣ Check if title or description changed
	titleChanged := title != "" && title != task.Title
	descriptionChanged := description != "" && description != task.Description

	needsAIReEstimation := titleChanged || descriptionChanged

	// 4️⃣ If content changed, trigger AI re-estimation
	if needsAIReEstimation {
		newTitle := title
		if newTitle == "" {
			newTitle = task.Title
		}

		newDescription := description
		if newDescription == "" {
			newDescription = task.Description
		}

		fmt.Printf("🔄 Task content changed. Triggering AI re-estimation...\n")
		fmt.Printf("   Old: [%s] %s\n", task.Title, task.Description)
		fmt.Printf("   New: [%s] %s\n", newTitle, newDescription)

		// Call Gemini AI for new estimation
		estimatedMinutes, reasoning, err := u.aiService.EstimateEffort(newTitle, newDescription)
		if err != nil {
			// Log warning but don't block the update
			fmt.Printf("⚠️  AI Re-estimation failed: %v (continuing with manual update)\n", err)
		} else {
			// Update AI estimation
			task.AIEstimatedMinutes = estimatedMinutes
			fmt.Printf("✅ AI Re-estimation Complete: %d minutes (%.1f hours)\n", estimatedMinutes, float64(estimatedMinutes)/60)
			fmt.Printf("   AI Reasoning: %s\n", reasoning)

			// 🔄 RESET NEGOTIATION STATUS (since AI has new estimate)
			if task.NegotiationStatus == "PENDING" || task.NegotiationStatus == "APPROVED" {
				fmt.Printf("🔄 Resetting negotiation status (AI has new estimate)\n")
				task.NegotiationStatus = "NONE"
				task.ProposedMinutes = 0
				task.NegotiationReason = ""
			}
		}
	}

	// 5️⃣ Apply updates
	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}

	// 6️⃣ Save to database
	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	fmt.Printf("✅ Task Updated: %s by %s (User ID: %d)\n", taskID, requestingUserRole, requestingUserID)

	return task, nil
}

// DeleteTask deletes a task with access control
// Only the Creator OR CEO can delete a task
func (u *sentinelUsecase) DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
	// 1️⃣ Fetch the task to check ownership
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2️⃣ ACCESS CONTROL: Only Creator or CEO can delete
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"

	if !isCreator && !isCEO {
		return fmt.Errorf("unauthorized: only the task creator or CEO can delete this task")
	}

	// 3️⃣ Delete from database
	if err := u.repo.DeleteTask(taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	fmt.Printf("🗑️  Task Deleted: %s by %s (User ID: %d)\n", taskID, requestingUserRole, requestingUserID)

	return nil
}

// --- Human Quality Gate ---

// ApproveTask marks a task as COMPLETED after human verification (PM/CEO only)
func (u *sentinelUsecase) ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error {
	// 🔒 ROLE VALIDATION: Only PM or CEO can approve tasks
	if approverRole != "CEO" && approverRole != "PM" {
		return fmt.Errorf("access denied: only PM or CEO can approve tasks (your role: %s)", approverRole)
	}

	// 1️⃣ Get the task
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}

	// 2️⃣ Verify task is in REVIEW_PENDING status
	if task.Status != "REVIEW_PENDING" {
		return fmt.Errorf("task is not pending review (current status: %s)", task.Status)
	}

	// 3️⃣ Approve the task (repository handles status change and CompletedAt)
	if err := u.repo.ApproveTask(taskID); err != nil {
		return fmt.Errorf("failed to approve task: %w", err)
	}

	// 4️⃣ Reload task to get updated CompletedAt for logging
	task, _ = u.repo.GetTaskByID(taskID)
	
	fmt.Printf("✅ Task %s APPROVED by %s (ID: %d)\n", taskID, approverRole, approverID)
	if task != nil && task.CompletedAt != nil {
		fmt.Printf("🎉 Task marked COMPLETED at %s\n", task.CompletedAt.Format(time.RFC3339))
		
		// Calculate actual time taken
		if task.StartedAt != nil {
			duration := task.CompletedAt.Sub(*task.StartedAt)
			fmt.Printf("📊 Actual Time: %.2f hours (AI Estimated: %.2f hours)\n", 
				duration.Hours(), float64(task.AIEstimatedMinutes)/60.0)
		}
	}

	return nil
}

// --- System Configuration Management ---

// GetSystemConfig retrieves the current AI system configuration
func (u *sentinelUsecase) GetSystemConfig() (*domain.SystemConfig, error) {
	config, err := u.repo.GetSystemConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get system config: %w", err)
	}
	return config, nil
}

// UpdateSystemConfig updates AI configuration (CEO only)
func (u *sentinelUsecase) UpdateSystemConfig(activeModel string, temperature float32, cursorAssistance int, userRole string) (*domain.SystemConfig, error) {
	// 🔒 ROLE VALIDATION: Only CEO can change system config
	if userRole != "CEO" {
		return nil, fmt.Errorf("access denied: only CEO can modify system configuration")
	}

	// Validate temperature (0.0 to 1.0)
	if temperature < 0.0 || temperature > 1.0 {
		return nil, fmt.Errorf("temperature must be between 0.0 and 1.0")
	}

	// Validate cursor assistance (0 to 100)
	if cursorAssistance < 0 || cursorAssistance > 100 {
		return nil, fmt.Errorf("cursor_assistance must be between 0 and 100")
	}

	// Validate model (must be in available list)
	availableModels := u.GetAvailableModels()
	validModel := false
	for _, m := range availableModels {
		if m == activeModel {
			validModel = true
			break
		}
	}
	if !validModel {
		return nil, fmt.Errorf("invalid model: %s (must be one of: %v)", activeModel, availableModels)
	}

	// Update config
	config := &domain.SystemConfig{
		ID:               1, // Singleton
		ActiveModel:      activeModel,
		Temperature:      temperature,
		CursorAssistance: cursorAssistance,
	}

	if err := u.repo.UpdateSystemConfig(config); err != nil {
		return nil, fmt.Errorf("failed to update system config: %w", err)
	}

	fmt.Printf("⚙️  System Config Updated: Model=%s, Temp=%.2f, Cursor=%d%% (by CEO)\n",
		activeModel, temperature, cursorAssistance)

	return config, nil
}

// GetAvailableModels returns list of supported Gemini models
func (u *sentinelUsecase) GetAvailableModels() []string {
	return []string{
		"gemini-1.5-flash",
		"gemini-1.5-pro",
		"gemini-2.0-flash-exp",
		"gemini-2.5-flash-lite",
		"gemini-exp-1206",
		"gemini-flash-lite-latest", // 🆕 Latest lite version
		"gemini-pro-latest",        // 🆕 Latest pro version
		"gemini-flash-latest",      // 🆕 Latest flash version
	}
}
