package usecase

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	chromepdf "github.com/portnd/the-sentinel-core/internal/core/pdf"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)

type sentinelUsecase struct {
	repo        domain.SentinelRepository
	aiService   domain.AIService
	authRepo    authDomain.Repository
	usageTracker domain.UsageTracker
	aiLimitRPM  int
	aiLimitRPD  int
	timeout     time.Duration
}

// NewSentinelUsecase creates the usecase. usageTracker may be nil (GetAIUsage will return zeros). aiLimitRPM/aiLimitRPD 0 = use tracker defaults.
func NewSentinelUsecase(repo domain.SentinelRepository, aiService domain.AIService, authRepo authDomain.Repository, usageTracker domain.UsageTracker, aiLimitRPM, aiLimitRPD int) domain.SentinelUsecase {
	return &sentinelUsecase{
		repo:         repo,
		aiService:    aiService,
		authRepo:     authRepo,
		usageTracker: usageTracker,
		aiLimitRPM:   aiLimitRPM,
		aiLimitRPD:   aiLimitRPD,
		timeout:      time.Second * 10,
	}
}

// --- Project Operations ---

// projectNameEnglishOnly matches letters, digits, spaces, hyphens, underscores (English only)
var projectNameEnglishOnly = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)

func (u *sentinelUsecase) CreateProject(name, description, status string) (*domain.Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("project name is required")
	}
	if !projectNameEnglishOnly.MatchString(name) {
		return nil, errors.New("project name must be in English only (letters, numbers, spaces, hyphens)")
	}
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "COMPLETED" && status != "ON_HOLD" {
		return nil, fmt.Errorf("invalid project status: %s (allowed: ACTIVE, COMPLETED, ON_HOLD)", status)
	}
	code := slugify(name)
	p := &domain.Project{
		Code:        code,
		Name:        name,
		Description: description,
		Status:      status,
	}
	if err := u.repo.CreateProject(p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return p, nil
}

func (u *sentinelUsecase) GetProjects() ([]domain.Project, error) {
	return u.repo.GetAllProjects()
}

func (u *sentinelUsecase) GetProjectDetails(id uuid.UUID) (*domain.Project, error) {
	p, err := u.repo.GetProjectByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	return p, nil
}

// GetProjectByIDOrCode retrieves a project by UUID or by code (e.g. mims-hdmap-main).
func (u *sentinelUsecase) GetProjectByIDOrCode(idOrCode string) (*domain.Project, error) {
	idOrCode = strings.TrimSpace(idOrCode)
	if idOrCode == "" {
		return nil, errors.New("project id or code is required")
	}
	if id, err := uuid.Parse(idOrCode); err == nil {
		return u.GetProjectDetails(id)
	}
	p, err := u.repo.GetProjectByCode(idOrCode)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	return p, nil
}

// UpdateProject updates project name, description, and status. If updateCode is true, also sets project.Code to slugify(name) and updates all task codes in the project to use the new prefix (so they match the new name).
func (u *sentinelUsecase) UpdateProject(projectID uuid.UUID, name, description, status string, updateCode bool) (*domain.Project, error) {
	p, err := u.repo.GetProjectByID(projectID)
	if err != nil || p == nil {
		return nil, errors.New("project not found")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("project name is required")
	}
	if !projectNameEnglishOnly.MatchString(name) {
		return nil, errors.New("project name must be in English only (letters, numbers, spaces, hyphens)")
	}
	p.Name = name
	p.Description = strings.TrimSpace(description)
	if status != "" {
		if status != "ACTIVE" && status != "COMPLETED" && status != "ON_HOLD" {
			return nil, fmt.Errorf("invalid project status: %s", status)
		}
		p.Status = status
	}
	if updateCode {
		newCode := slugify(name)
		if newCode != "" {
			p.Code = newCode
			// Update all task codes in this project to use the new prefix (globally unique suffixes).
			tasks, err := u.repo.GetTasksByProjectID(projectID)
			if err != nil {
				return nil, fmt.Errorf("failed to get project tasks: %w", err)
			}
			maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(newCode)
			for i := range tasks {
				maxSuffix++
				tasks[i].Code = fmt.Sprintf("%s-%03d", newCode, maxSuffix)
				if err := u.repo.UpdateTask(&tasks[i]); err != nil {
					return nil, fmt.Errorf("failed to update task code: %w", err)
				}
			}
		}
	}
	if err := u.repo.UpdateProject(p); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return p, nil
}

func (u *sentinelUsecase) DeleteProject(id uuid.UUID) error {
	return u.repo.DeleteProject(id)
}

// slugify converts project name to code prefix, e.g. "MIMS HDMap Main" -> "mims-hdmap-main"
func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "task"
	}
	return s
}

var validPriorities = map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}

func (u *sentinelUsecase) CreateTask(title, desc string, creatorID uint, dueDate *time.Time, projectID, parentID *uuid.UUID, startDate, endDate *time.Time, priority string, storyPoints int, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID) (*domain.Task, error) {
	const defaultEstimatedMinutes = 0

	if priority == "" {
		priority = "MEDIUM"
	}
	if !validPriorities[priority] {
		return nil, fmt.Errorf("invalid priority: %s (allowed: CRITICAL, HIGH, MEDIUM, LOW)", priority)
	}
	if storyPoints < 0 {
		return nil, errors.New("story_points cannot be negative")
	}

	// Sub-tasks (have a parent_id) inherit dates from their parent — clear any provided dates
	if parentID != nil {
		parent, err := u.repo.GetTaskByID(*parentID)
		if err == nil && parent != nil {
			startDate = parent.StartDate
			endDate = parent.EndDate
		} else {
			startDate = nil
			endDate = nil
		}
	}

	slug := "task"
	if projectID != nil {
		proj, err := u.repo.GetProjectByID(*projectID)
		if err == nil && proj != nil {
			slug = slugify(proj.Name)
		}
	}
	maxSuffix, err := u.repo.GetMaxTaskCodeSuffix(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get next task code: %w", err)
	}
	code := fmt.Sprintf("%s-%03d", slug, maxSuffix+1)

	task := &domain.Task{
		ID:                 uuid.New(),
		Code:               code,
		Title:              title,
		Description:        desc,
		CreatedBy:          &creatorID,
		Status:             "PENDING",
		AIEstimatedMinutes: defaultEstimatedMinutes,
		DueAt:              dueDate,
		ProjectID:          projectID,
		ParentID:           parentID,
		EpicID:             epicID,
		StartDate:          startDate,
		EndDate:            endDate,
		Priority:           priority,
		StoryPoints:        storyPoints,
		SprintID:           sprintID,
		MilestoneID:        milestoneID,
	}

	if err := u.repo.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

// AssignTask assigns a developer to a task. assignerID is the PM/CEO who performs the assign (for PM-scoped leaderboard).
func (u *sentinelUsecase) AssignTask(taskID uuid.UUID, devID uint, assignerID uint) error {
	// 1. Validate if task exists
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	// 2. Update assignment
	task.AssignedTo = &devID
	task.AssignedByID = &assignerID // PM/CEO who assigned (drives PM Team Leaderboard scope)
	task.Status = "IN_PROGRESS"

	// 3. ⏰ Start Time Tracking: Set StartedAt = NOW()
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
	// No AI code review: submission is always PASS, human approval still required
	sub := &domain.Submission{
		ID:         uuid.New(),
		TaskID:     taskID,
		DevID:      devID,
		CommitHash: commitHash,
		Diff:       diff,
		AIVerdict:  "PASS",
		AIScore:    100,
		AIFeedback: []byte(`{"feedback": ""}`),
	}

	if err := u.repo.CreateSubmission(sub); err != nil {
		return nil, err
	}

	// Move to REVIEW_PENDING for PM/CEO approval (human quality gate)
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		fmt.Printf("⚠️  Failed to get task for review queue: %v\n", err)
	} else {
		task.Status = "REVIEW_PENDING"
		if err := u.repo.UpdateTask(task); err != nil {
			fmt.Printf("⚠️  Failed to update task status to REVIEW_PENDING: %v\n", err)
		} else {
			fmt.Printf("🚦 Task %s moved to REVIEW_PENDING - awaiting PM/CEO approval\n", task.ID)
		}
	}

	return sub, nil
}

// GetTaskByID retrieves a single task with full submission history.
// Enriches task with created_by_role and created_by_email from auth.
func (u *sentinelUsecase) GetTaskByID(taskID uuid.UUID) (*domain.Task, error) {
	return u.getTaskByIDAndEnrich(taskID, nil)
}

// GetTaskByIDOrCode retrieves a task by UUID or by code (e.g. mims-hdmap-main-001).
func (u *sentinelUsecase) GetTaskByIDOrCode(idOrCode string) (*domain.Task, error) {
	idOrCode = strings.TrimSpace(idOrCode)
	if idOrCode == "" {
		return nil, errors.New("task id or code is required")
	}
	if id, err := uuid.Parse(idOrCode); err == nil {
		return u.getTaskByIDAndEnrich(id, nil)
	}
	task, err := u.repo.GetTaskByCode(idOrCode)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	return u.getTaskByIDAndEnrich(task.ID, task)
}

func (u *sentinelUsecase) getTaskByIDAndEnrich(taskID uuid.UUID, task *domain.Task) (*domain.Task, error) {
	if task == nil {
		var err error
		task, err = u.repo.GetTaskByID(taskID)
		if err != nil {
			return nil, fmt.Errorf("failed to get task: %w", err)
		}
		if task == nil {
			return nil, errors.New("task not found")
		}
	}
	if task.CreatedBy != nil {
		creator, err := u.authRepo.FindByID(*task.CreatedBy)
		if err == nil && creator != nil {
			task.CreatedByRole = creator.Role
			task.CreatedByEmail = creator.Email
		}
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

// GetTasksByProjectID returns tasks for a single project (for project Board/Backlog/Overview).
func (u *sentinelUsecase) GetTasksByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	return u.repo.GetTasksByProjectID(projectID)
}

// GetGanttData returns tasks and dependencies for Gantt chart rendering.
// If projectID is set, only tasks (and dependencies between them) for that project are returned.
func (u *sentinelUsecase) GetGanttData(projectID *uuid.UUID) (*domain.GanttData, error) {
	var tasks []domain.Task
	var err error
	if projectID != nil {
		tasks, err = u.repo.GetTasksByProjectID(*projectID)
	} else {
		tasks, err = u.repo.GetAllTasks()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	deps, err := u.repo.GetAllTaskDependencies()
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}
	// If filtering by project, keep only dependencies where both tasks are in the task set
	if projectID != nil {
		taskIDSet := make(map[uuid.UUID]bool)
		for _, t := range tasks {
			taskIDSet[t.ID] = true
		}
		var filtered []domain.TaskDependency
		for _, d := range deps {
			if taskIDSet[d.PredecessorID] && taskIDSet[d.SuccessorID] {
				filtered = append(filtered, d)
			}
		}
		deps = filtered
	}
	return &domain.GanttData{Tasks: tasks, Dependencies: deps}, nil
}

// Valid dependency types (MS Project–style)
var validDependencyTypes = map[string]bool{"FS": true, "SS": true, "FF": true, "SF": true}

// AddDependency creates a link from predecessor to successor (e.g. Task B waits for Task A)
func (u *sentinelUsecase) AddDependency(predecessorID, successorID uuid.UUID, depType string) (*domain.TaskDependency, error) {
	if predecessorID == successorID {
		return nil, errors.New("predecessor and successor must be different tasks (no self-linking)")
	}
	if depType == "" {
		depType = "FS"
	}
	if !validDependencyTypes[depType] {
		return nil, fmt.Errorf("invalid dependency type %q (allowed: FS, SS, FF, SF)", depType)
	}
	dep := &domain.TaskDependency{
		ID:            uuid.New(),
		PredecessorID: predecessorID,
		SuccessorID:   successorID,
		Type:          depType,
	}
	if err := u.repo.CreateTaskDependency(dep); err != nil {
		return nil, fmt.Errorf("failed to create dependency: %w", err)
	}
	return dep, nil
}

// RemoveDependency deletes a task dependency by ID
func (u *sentinelUsecase) RemoveDependency(id uuid.UUID) error {
	return u.repo.DeleteTaskDependency(id)
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

	// No AI: appeal is human-only; CEO/PM decide without AI advisory
	appeal := &domain.Appeal{
		ID:                uuid.New(),
		SubmissionID:      submissionID,
		DeveloperID:       devID,
		Reason:            reason,
		Status:            "PENDING",
		AIRecommendation:  "",
		AIConfidence:      0,
		AIReasoning:       "",
	}

	if err := u.repo.CreateAppeal(appeal); err != nil {
		return nil, fmt.Errorf("failed to create appeal: %w", err)
	}

	fmt.Printf("⚖️  Appeal created (human review only)\n")

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
	if task.AIEstimatedMinutes > 0 && minutes <= task.AIEstimatedMinutes {
		return errors.New("proposed time must be greater than current estimated minutes")
	}

	// 4. Validate reason
	if len(reason) < 20 {
		return errors.New("negotiation reason must be at least 20 characters")
	}

	// 5. Check if already negotiating
	if task.NegotiationStatus == "PENDING" {
		return errors.New("time negotiation already pending review")
	}

	// No AI: store negotiation for PM/CEO to approve manually
	task.NegotiationStatus = "PENDING"
	task.ProposedMinutes = minutes
	task.NegotiationReason = reason
	task.NegotiationAIRecommendation = ""
	task.NegotiationAIConfidence = 0
	task.NegotiationAIReasoning = ""

	if err := u.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to submit time negotiation: %w", err)
	}

	fmt.Printf("⏰ Time Negotiation Submitted: Task %s | Proposed: %d min\n", task.ID, minutes)

	return nil
}

// UpdateTask updates a task with access control (no AI).
// Creator, CEO, or PM can update. Gantt fields applied when provided.
func (u *sentinelUsecase) UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string, parentID *uuid.UUID, startDate, endDate *time.Time, progress *int, priority string, storyPoints *int, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID, applyEpic bool, sortOrder *int) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isPM := requestingUserRole == "PM"

	if !isCreator && !isCEO && !isPM {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, or PM can update this task")
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if parentID != nil {
		task.ParentID = parentID
	}
	if startDate != nil {
		task.StartDate = startDate
	}
	if endDate != nil {
		task.EndDate = endDate
		task.DueAt = endDate
	}
	if progress != nil {
		if *progress < 0 || *progress > 100 {
			return nil, fmt.Errorf("progress must be between 0 and 100")
		}
		task.Progress = *progress
	}
	if priority != "" {
		if !validPriorities[priority] {
			return nil, fmt.Errorf("invalid priority: %s", priority)
		}
		task.Priority = priority
	}
	if storyPoints != nil {
		if *storyPoints < 0 {
			return nil, errors.New("story_points cannot be negative")
		}
		task.StoryPoints = *storyPoints
	}
	if sprintID != nil {
		task.SprintID = sprintID
	}
	if milestoneID != nil {
		task.MilestoneID = milestoneID
	}
	if applyEpic {
		task.EpicID = epicID
	}
	if sortOrder != nil {
		task.SortOrder = *sortOrder
	}

	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	fmt.Printf("✅ Task Updated: %s by %s (User ID: %d)\n", taskID, requestingUserRole, requestingUserID)

	return task, nil
}

// UpdateTaskResourceURLs updates only task.resource_urls (e.g. slide images/annotations). Same permission as UpdateTask.
func (u *sentinelUsecase) UpdateTaskResourceURLs(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, resourceURLs datatypes.JSON) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isPM := requestingUserRole == "PM"
	if !isCreator && !isCEO && !isPM {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, or PM can update this task")
	}
	task.ResourceURLs = resourceURLs
	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task resource_urls: %w", err)
	}
	return task, nil
}

// EstimateTask uses AI to estimate task effort (title + description) and updates task.ai_estimated_minutes.
// Only task creator, CEO, or PM can run estimate.
func (u *sentinelUsecase) EstimateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isPM := requestingUserRole == "PM"
	if !isCreator && !isCEO && !isPM {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, or PM can run AI estimate")
	}
	minutes, _, err := u.aiService.EstimateEffort(task.Title, task.Description)
	if err != nil {
		return nil, fmt.Errorf("AI estimate failed: %w", err)
	}
	task.AIEstimatedMinutes = minutes
	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task estimate: %w", err)
	}
	fmt.Printf("✅ AI Estimate: Task %s → %d minutes by %s (User ID: %d)\n", taskID, minutes, requestingUserRole, requestingUserID)
	return task, nil
}

// GenerateProjectPlan uses AI to generate a full work plan (epics, milestones, sprints, tasks) and creates them in the project.
// Only CEO or PM can run this.
func (u *sentinelUsecase) GenerateProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (*domain.AIGeneratedPlan, error) {
	if requestingUserRole != "CEO" && requestingUserRole != "PM" {
		return nil, fmt.Errorf("unauthorized: only CEO or PM can generate AI work plan")
	}
	proj, err := u.repo.GetProjectByID(projectID)
	if err != nil || proj == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	plan, err := u.aiService.GenerateWorkPlan(proj.Name, proj.Description)
	if err != nil {
		return nil, fmt.Errorf("AI plan failed: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("AI plan failed: empty plan returned")
	}
	parseDate := func(s string) *time.Time {
		if s == "" {
			return nil
		}
		t, e := time.Parse("2006-01-02", strings.TrimSpace(s))
		if e != nil {
			return nil
		}
		return &t
	}
	epicIDs := make([]uuid.UUID, 0, len(plan.Epics))
	for _, e := range plan.Epics {
		color := e.Color
		if color == "" {
			color = "#6366f1"
		}
		epic, err := u.CreateEpic(projectID, e.Title, e.Description, color, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("create epic %q: %w", e.Title, err)
		}
		epicIDs = append(epicIDs, epic.ID)
	}
	milestoneIDs := make([]uuid.UUID, 0, len(plan.Milestones))
	for _, m := range plan.Milestones {
		due := parseDate(m.DueDate)
		milestone, err := u.CreateMilestone(projectID, m.Title, m.Description, due)
		if err != nil {
			return nil, fmt.Errorf("create milestone %q: %w", m.Title, err)
		}
		milestoneIDs = append(milestoneIDs, milestone.ID)
	}
	sprintIDs := make([]uuid.UUID, 0, len(plan.Sprints))
	for _, s := range plan.Sprints {
		start := parseDate(s.StartDate)
		end := parseDate(s.EndDate)
		sprint, err := u.CreateSprint(projectID, s.Name, s.Goal, start, end)
		if err != nil {
			return nil, fmt.Errorf("create sprint %q: %w", s.Name, err)
		}
		sprintIDs = append(sprintIDs, sprint.ID)
	}
	for _, t := range plan.Tasks {
		priority := t.Priority
		if priority == "" {
			priority = "MEDIUM"
		}
		if priority != "CRITICAL" && priority != "HIGH" && priority != "MEDIUM" && priority != "LOW" {
			priority = "MEDIUM"
		}
		if t.StoryPoints < 0 {
			t.StoryPoints = 0
		}
		var epicID, sprintID, milestoneID *uuid.UUID
		if t.EpicIndex != nil && *t.EpicIndex >= 0 && *t.EpicIndex < len(epicIDs) {
			id := epicIDs[*t.EpicIndex]
			epicID = &id
		}
		if t.SprintIndex != nil && *t.SprintIndex >= 0 && *t.SprintIndex < len(sprintIDs) {
			id := sprintIDs[*t.SprintIndex]
			sprintID = &id
		}
		if t.MilestoneIndex != nil && *t.MilestoneIndex >= 0 && *t.MilestoneIndex < len(milestoneIDs) {
			id := milestoneIDs[*t.MilestoneIndex]
			milestoneID = &id
		}
		startDate := parseDate(t.StartDate)
		endDate := parseDate(t.EndDate)
		var dueDate *time.Time
		if endDate != nil {
			dueDate = endDate
		}
		_, err := u.CreateTask(t.Title, t.Description, requestingUserID, dueDate, &projectID, nil, startDate, endDate, priority, t.StoryPoints, sprintID, milestoneID, epicID)
		if err != nil {
			return nil, fmt.Errorf("create task %q: %w", t.Title, err)
		}
	}
	fmt.Printf("✅ AI Plan created: %d epics, %d milestones, %d sprints, %d tasks (project %s)\n",
		len(epicIDs), len(milestoneIDs), len(sprintIDs), len(plan.Tasks), projectID)
	return plan, nil
}

// ClearProjectPlan removes all tasks, sprints, milestones, and epics for the project. Only CEO or PM.
func (u *sentinelUsecase) ClearProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
	if requestingUserRole != "CEO" && requestingUserRole != "PM" {
		return fmt.Errorf("unauthorized: only CEO or PM can clear project plan")
	}
	if _, err := u.repo.GetProjectByID(projectID); err != nil {
		return fmt.Errorf("project not found: %w", err)
	}
	if err := u.repo.DeleteProjectPlan(projectID); err != nil {
		return fmt.Errorf("failed to clear plan: %w", err)
	}
	fmt.Printf("🗑️  Project plan cleared: %s (by %s)\n", projectID, requestingUserRole)
	return nil
}

// ScheduleProjectWithAI ประเมินเวลาและจัดเรียง timeline ของ task ที่มีอยู่แล้ว (ไม่สร้าง task ใหม่). เฉพาะ CEO/PM.
func (u *sentinelUsecase) ScheduleProjectWithAI(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (int, error) {
	if requestingUserRole != "CEO" && requestingUserRole != "PM" {
		return 0, fmt.Errorf("unauthorized: only CEO or PM can run AI schedule")
	}
	if _, err := u.repo.GetProjectByID(projectID); err != nil {
		return 0, fmt.Errorf("project not found: %w", err)
	}
	tasks, err := u.repo.GetTasksByProjectID(projectID)
	if err != nil {
		return 0, fmt.Errorf("failed to get tasks: %w", err)
	}
	if len(tasks) == 0 {
		return 0, fmt.Errorf("no tasks to schedule: create tasks first")
	}
	inputs := make([]domain.TaskEstimateInput, 0, len(tasks))
	for i, t := range tasks {
		inputs = append(inputs, domain.TaskEstimateInput{
			Index:       i,
			Title:       t.Title,
			Description: t.Description,
			Priority:    t.Priority,
			StoryPoints: t.StoryPoints,
		})
	}
	results, err := u.aiService.EstimateAndScheduleTasks(inputs)
	if err != nil {
		return 0, fmt.Errorf("AI estimate failed: %w", err)
	}
	// Build map task_index -> result (use first occurrence per task_index if duplicate)
	byIndex := make(map[int]domain.TaskEstimateAndOrder)
	for _, r := range results {
		if r.TaskIndex >= 0 && r.TaskIndex < len(tasks) {
			byIndex[r.TaskIndex] = r
		}
	}
	// Sort by Order (1, 2, 3...) to assign timeline
	ordered := make([]struct{ taskIdx int; res domain.TaskEstimateAndOrder }, 0, len(byIndex))
	for idx, res := range byIndex {
		ordered = append(ordered, struct{ taskIdx int; res domain.TaskEstimateAndOrder }{idx, res})
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].res.Order < ordered[j].res.Order })
	// Start from start of today (UTC) or next midnight local; use simple "today" for cursor
	now := time.Now()
	cursor := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if cursor.Before(now) {
		cursor = now
	}
	updated := 0
	for _, o := range ordered {
		task := &tasks[o.taskIdx]
		minutes := o.res.Minutes
		if minutes <= 0 {
			minutes = 60
		}
		end := cursor.Add(time.Duration(minutes) * time.Minute)
		task.AIEstimatedMinutes = minutes
		task.StartDate = &cursor
		task.EndDate = &end
		if err := u.repo.UpdateTask(task); err != nil {
			return updated, fmt.Errorf("update task %s: %w", task.ID, err)
		}
		updated++
		cursor = end
	}
	fmt.Printf("✅ AI Schedule: %d tasks updated (project %s)\n", updated, projectID)
	return updated, nil
}

// DeleteTask deletes a task with access control
// Only the Creator OR CEO can delete a task
func (u *sentinelUsecase) DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
	// 1️⃣ Fetch the task to check ownership
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2️⃣ ACCESS CONTROL: Creator, CEO, or PM can delete
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isPM := requestingUserRole == "PM"

	if !isCreator && !isCEO && !isPM {
		return fmt.Errorf("unauthorized: only the task creator, CEO, or PM can delete this task")
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

var fallbackModels = []string{
	"gemini-1.5-flash",
	"gemini-1.5-pro",
	"gemini-2.0-flash-exp",
	"gemini-2.5-flash-lite",
	"gemini-exp-1206",
	"gemini-flash-lite-latest",
	"gemini-pro-latest",
	"gemini-flash-latest",
}

// GetAvailableModels returns list of Gemini models from List Models API, or fallback if API key missing/fails.
func (u *sentinelUsecase) GetAvailableModels() []string {
	list, err := u.aiService.ListModels()
	if err != nil || len(list) == 0 {
		return fallbackModels
	}
	sort.Strings(list)
	return list
}

// GetAIUsage returns approximate Gemini API usage (requests last minute, today) and remaining quota. Uses in-memory tracker; limits from config or default.
func (u *sentinelUsecase) GetAIUsage() domain.AIUsage {
	if u.usageTracker == nil {
		return domain.AIUsage{}
	}
	return u.usageTracker.GetUsage(u.aiLimitRPM, u.aiLimitRPD)
}

// --- Sprint Operations ---

func (u *sentinelUsecase) CreateSprint(projectID uuid.UUID, name, goal string, startDate, endDate *time.Time) (*domain.Sprint, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("sprint name is required")
	}
	existing, _ := u.repo.GetSprintsByProjectID(projectID)
	sortOrder := 0
	for _, s := range existing {
		if s.SortOrder >= sortOrder {
			sortOrder = s.SortOrder + 1
		}
	}
	sprint := &domain.Sprint{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		Goal:      goal,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "PLANNING",
		SortOrder: sortOrder,
	}
	if err := u.repo.CreateSprint(sprint); err != nil {
		return nil, fmt.Errorf("failed to create sprint: %w", err)
	}
	return sprint, nil
}

func (u *sentinelUsecase) GetSprintsByProject(projectID uuid.UUID) ([]domain.Sprint, error) {
	return u.repo.GetSprintsByProjectID(projectID)
}

func (u *sentinelUsecase) StartSprint(sprintID uuid.UUID) (*domain.Sprint, error) {
	sprint, err := u.repo.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	if sprint.Status != "PLANNING" {
		return nil, fmt.Errorf("sprint is already %s", sprint.Status)
	}
	active, err := u.repo.GetActiveSprintByProjectID(sprint.ProjectID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, fmt.Errorf("project already has an active sprint: %s", active.Name)
	}
	sprint.Status = "ACTIVE"
	now := time.Now()
	if sprint.StartDate == nil {
		sprint.StartDate = &now
	}
	if err := u.repo.UpdateSprint(sprint); err != nil {
		return nil, fmt.Errorf("failed to start sprint: %w", err)
	}
	return sprint, nil
}

func (u *sentinelUsecase) CompleteSprint(sprintID uuid.UUID) (*domain.Sprint, error) {
	sprint, err := u.repo.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	if sprint.Status != "ACTIVE" {
		return nil, fmt.Errorf("only ACTIVE sprints can be completed (current: %s)", sprint.Status)
	}
	sprint.Status = "COMPLETED"
	now := time.Now()
	if sprint.EndDate == nil {
		sprint.EndDate = &now
	}
	if err := u.repo.UpdateSprint(sprint); err != nil {
		return nil, fmt.Errorf("failed to complete sprint: %w", err)
	}
	return sprint, nil
}

func (u *sentinelUsecase) ReopenSprint(sprintID uuid.UUID) (*domain.Sprint, error) {
	sprint, err := u.repo.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	if sprint.Status != "COMPLETED" {
		return nil, fmt.Errorf("only COMPLETED sprints can be reopened (current: %s)", sprint.Status)
	}
	active, err := u.repo.GetActiveSprintByProjectID(sprint.ProjectID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, fmt.Errorf("project already has an active sprint: %s (complete or reopen it first)", active.Name)
	}
	sprint.Status = "ACTIVE"
	if err := u.repo.UpdateSprint(sprint); err != nil {
		return nil, fmt.Errorf("failed to reopen sprint: %w", err)
	}
	return sprint, nil
}

func (u *sentinelUsecase) AddTasksToSprint(sprintID uuid.UUID, taskIDs []uuid.UUID) error {
	if len(taskIDs) == 0 {
		return errors.New("no tasks provided")
	}
	sprint, err := u.repo.GetSprintByID(sprintID)
	if err != nil {
		return fmt.Errorf("sprint not found: %w", err)
	}
	for _, tid := range taskIDs {
		task, err := u.repo.GetTaskByID(tid)
		if err != nil {
			continue
		}
		task.SprintID = &sprint.ID
		u.repo.UpdateTask(task)
	}
	return nil
}

func (u *sentinelUsecase) UpdateSprint(sprintID uuid.UUID, name, goal string, startDate, endDate *time.Time, sortOrder *int) (*domain.Sprint, error) {
	sprint, err := u.repo.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	if name != "" {
		sprint.Name = strings.TrimSpace(name)
	}
	if goal != "" {
		sprint.Goal = goal
	}
	if startDate != nil {
		sprint.StartDate = startDate
	}
	if endDate != nil {
		sprint.EndDate = endDate
	}
	if sortOrder != nil {
		sprint.SortOrder = *sortOrder
	}
	if err := u.repo.UpdateSprint(sprint); err != nil {
		return nil, fmt.Errorf("failed to update sprint: %w", err)
	}
	return sprint, nil
}

func (u *sentinelUsecase) DeleteSprint(sprintID uuid.UUID) error {
	return u.repo.DeleteSprint(sprintID)
}

// --- Milestone Operations ---

func (u *sentinelUsecase) CreateMilestone(projectID uuid.UUID, title, description string, dueDate *time.Time) (*domain.Milestone, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("milestone title is required")
	}
	m := &domain.Milestone{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      "PENDING",
	}
	if err := u.repo.CreateMilestone(m); err != nil {
		return nil, fmt.Errorf("failed to create milestone: %w", err)
	}
	return m, nil
}

func (u *sentinelUsecase) GetMilestonesByProject(projectID uuid.UUID) ([]domain.Milestone, error) {
	return u.repo.GetMilestonesByProjectID(projectID)
}

func (u *sentinelUsecase) UpdateMilestone(id uuid.UUID, title, description, status string, dueDate *time.Time) (*domain.Milestone, error) {
	m, err := u.repo.GetMilestoneByID(id)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %w", err)
	}
	if title != "" {
		m.Title = strings.TrimSpace(title)
	}
	if description != "" {
		m.Description = description
	}
	if status != "" {
		if status != "PENDING" && status != "REACHED" && status != "MISSED" {
			return nil, fmt.Errorf("invalid milestone status: %s", status)
		}
		m.Status = status
	}
	if dueDate != nil {
		m.DueDate = dueDate
	}
	if err := u.repo.UpdateMilestone(m); err != nil {
		return nil, fmt.Errorf("failed to update milestone: %w", err)
	}
	return m, nil
}

func (u *sentinelUsecase) DeleteMilestone(id uuid.UUID) error {
	return u.repo.DeleteMilestone(id)
}

// --- Comment Operations ---

func (u *sentinelUsecase) AddComment(taskID uuid.UUID, userID uint, content string) (*domain.TaskComment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment content cannot be empty")
	}
	c := &domain.TaskComment{
		ID:      uuid.New(),
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}
	if err := u.repo.CreateTaskComment(c); err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}
	user, err := u.authRepo.FindByID(userID)
	if err == nil && user != nil {
		c.UserEmail = user.Email
	}
	return c, nil
}

func (u *sentinelUsecase) GetComments(taskID uuid.UUID) ([]domain.TaskComment, error) {
	comments, err := u.repo.GetCommentsByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	for i := range comments {
		user, err := u.authRepo.FindByID(comments[i].UserID)
		if err == nil && user != nil {
			comments[i].UserEmail = user.Email
		}
	}
	return comments, nil
}

// --- Time Log Operations ---

func (u *sentinelUsecase) LogTime(taskID uuid.UUID, userID uint, minutes int, description string) (*domain.TimeLog, error) {
	if minutes <= 0 {
		return nil, errors.New("minutes must be greater than 0")
	}
	if _, err := u.repo.GetTaskByID(taskID); err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	t := &domain.TimeLog{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      userID,
		Minutes:     minutes,
		Description: strings.TrimSpace(description),
	}
	if err := u.repo.CreateTimeLog(t); err != nil {
		return nil, fmt.Errorf("failed to log time: %w", err)
	}
	user, err := u.authRepo.FindByID(userID)
	if err == nil && user != nil {
		t.UserEmail = user.Email
	}
	return t, nil
}

func (u *sentinelUsecase) GetTimeLogs(taskID uuid.UUID) ([]domain.TimeLog, error) {
	logs, err := u.repo.GetTimeLogsByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	for i := range logs {
		user, err := u.authRepo.FindByID(logs[i].UserID)
		if err == nil && user != nil {
			logs[i].UserEmail = user.Email
		}
	}
	return logs, nil
}

// --- Analytics ---

func (u *sentinelUsecase) GetProjectAnalytics(projectID uuid.UUID) (*domain.ProjectAnalytics, error) {
	return u.repo.GetProjectAnalytics(projectID)
}

// --- Bulk Operations ---

func (u *sentinelUsecase) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error {
	validStatuses := map[string]bool{
		"PENDING": true, "IN_PROGRESS": true, "REVIEW_PENDING": true, "COMPLETED": true, "BLOCKED": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	if len(taskIDs) == 0 {
		return errors.New("no tasks provided")
	}
	return u.repo.BulkUpdateTaskStatus(taskIDs, status)
}

// --- Google Slides Import ---

// slideInfo is the unified internal representation of a parsed slide
type slideInfo struct {
	Index        int
	Title        string
	Body         string
	Notes        string
	ThumbnailURL string   // URL (from Slides API) or empty
	Images       []string // base64 data URLs (from PPTX media, no API key needed)
	SlideObjID   string   // populated when using Slides API
}

// ---- PPTX parsing (no API key required for public presentations) ----

type pptxPresentation struct {
	Title    string
	Slides   []pptxSlideData
	SlideIDs []string // object IDs from presentation.xml, same order as Slides (for export?format=jpeg&slide=id.XXX)
}

type pptxSlideData struct {
	Index  int
	Title  string
	Body   string
	Notes  string
	Images []string // base64 data URLs, largest-first
	Hidden bool     // true when slide has show="0" in presentation.xml (hidden/skipped in show)
}

type pptxImageEntry struct {
	DataURL string
	Size    int
}

var (
	presentationIDRegex = regexp.MustCompile(`/presentation/d/([a-zA-Z0-9_-]+)`)
	pptxSlideNumRegex   = regexp.MustCompile(`slide(\d+)\.xml$`)
	pptxTitlePhRe       = regexp.MustCompile(`<p:ph[^>]*type="(?:title|ctrTitle)"`)
	pptxSystemPhRe      = regexp.MustCompile(`<p:ph[^>]*type="(?:dt|ftr|sldNum|hdr)"`)
	pptxNotesBodyPhRe   = regexp.MustCompile(`<p:ph(?:[^>]*idx="1"|[^>]*type="body")`)
	pptxATextRe         = regexp.MustCompile(`<a:t(?:\s[^>]*)?>([^<]*)</a:t>`)
	pptxNotesRelRe      = regexp.MustCompile(`Type="[^"]*notesSlide"[^>]*Target="([^"]+)"`)
	pptxImageRelRe       = regexp.MustCompile(`Type="[^"]*image"[^>]*Target="([^"]+)"`)
	pptxSlideLayoutRelRe = regexp.MustCompile(`Type="[^"]*slideLayout"[^>]*Target="([^"]+)"`)
	pptxDocTitleRe    = regexp.MustCompile(`<dc:title>([^<]*)</dc:title>`)
	pptxSldIdRe       = regexp.MustCompile(`<p:sldId[^>]*id="(\d+)"`)
	pptxSldShowAttrRe = regexp.MustCompile(`<p:sld\s*([^>]+)>`) // opening tag of slide: check for show="0" or show="false" (hidden)
)

const (
	pptxMaxImageBytes     = 2 * 1024 * 1024 // 2MB max per image — เก็บรูปใหญ่ให้ครบเหมือนต้นฉบับ
	pptxMinImageBytes     = 256              // skip tiny icons/bullets
	pptxMaxImagesPerSlide = 30               // เก็บได้หลายรูปต่อหน้า (รูปฝัง + รูป export)
)

func pptxMimeType(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "" // EMF, WMF, SVG etc. — skip
	}
}

// extractImagesFromRels reads a _rels file and extracts all image targets from relDir.
// relDir is the directory containing the .rels file (e.g. ppt/slides for slide1.xml.rels).
func extractImagesFromRels(zr *zip.Reader, relDir string, relsData []byte) []pptxImageEntry {
	matches := pptxImageRelRe.FindAllStringSubmatch(string(relsData), -1)
	seen := make(map[string]bool)
	var entries []pptxImageEntry

	for _, m := range matches {
		target := m[1]
		mediaPath := path.Clean(relDir + "/" + target)
		if seen[mediaPath] {
			continue
		}
		seen[mediaPath] = true

		ext := path.Ext(mediaPath)
		mime := pptxMimeType(ext)
		if mime == "" {
			continue
		}

		imgData := readZipEntry(zr, mediaPath)
		if len(imgData) < pptxMinImageBytes || len(imgData) > pptxMaxImageBytes {
			continue
		}

		b64 := base64.StdEncoding.EncodeToString(imgData)
		dataURL := fmt.Sprintf("data:%s;base64,%s", mime, b64)
		entries = append(entries, pptxImageEntry{DataURL: dataURL, Size: len(imgData)})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})
	return entries
}

// extractSlideImages extracts embedded images from a PPTX slide via its _rels file.
// If the slide has no images, tries the slide layout's _rels (background/placeholders).
// Returns base64 data URLs sorted by image size (largest first).
func extractSlideImages(zr *zip.Reader, slideFile string) []string {
	dir := path.Dir(slideFile)
	base := path.Base(slideFile)
	relsPath := dir + "/_rels/" + base + ".rels"

	relsData := readZipEntry(zr, relsPath)
	if relsData == nil {
		return nil
	}

	entries := extractImagesFromRels(zr, dir, relsData)

	// Fallback: if slide has no images, try the slide layout (title slides often have only layout art)
	if len(entries) == 0 {
		if layoutMatches := pptxSlideLayoutRelRe.FindStringSubmatch(string(relsData)); len(layoutMatches) > 1 {
			layoutTarget := layoutMatches[1]
			layoutPath := path.Clean(dir + "/" + layoutTarget)
			layoutDir := path.Dir(layoutPath)
			layoutBase := path.Base(layoutPath)
			layoutRelsPath := layoutDir + "/_rels/" + layoutBase + ".rels"
			if layoutRelsData := readZipEntry(zr, layoutRelsPath); layoutRelsData != nil {
				entries = extractImagesFromRels(zr, layoutDir, layoutRelsData)
			}
		}
	}

	result := make([]string, 0, len(entries))
	for i, e := range entries {
		if i >= pptxMaxImagesPerSlide {
			break
		}
		result = append(result, e.DataURL)
	}
	return result
}

// fetchSlideExportImage gets the full-slide image from Google's export (no API key).
// slideID is the object ID from presentation.xml (e.g. "256"). Returns base64 data URL or empty.
func fetchSlideExportImage(presentationID, slideID string) string {
	if presentationID == "" || slideID == "" {
		return ""
	}
	// Google: /export?format=jpeg&slide=id.[id] — try both "id.256" and "256"
	for _, slideParam := range []string{"id." + slideID, slideID} {
		exportURL := fmt.Sprintf(
			"https://docs.google.com/presentation/d/%s/export?format=jpeg&slide=%s",
			presentationID, slideParam,
		)
		client := &http.Client{
			Timeout: 8 * time.Second,
			CheckRedirect: func(req *http.Request, _ []*http.Request) error {
				if strings.Contains(req.URL.Host, "accounts.google.com") {
					return fmt.Errorf("auth required")
				}
				return nil
			},
		}
		resp, err := client.Get(exportURL) //nolint:noctx
		if err != nil {
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		// Must look like image (JPEG magic bytes or reasonable size)
		if len(data) < 500 {
			continue
		}
		b64 := base64.StdEncoding.EncodeToString(data)
		return fmt.Sprintf("data:image/jpeg;base64,%s", b64)
	}
	return ""
}

func extractPresentationID(rawURL string) (string, error) {
	m := presentationIDRegex.FindStringSubmatch(rawURL)
	if len(m) < 2 {
		return "", errors.New("invalid Google Slides URL: could not extract presentation ID")
	}
	return m[1], nil
}

var xmlEntities = strings.NewReplacer(
	"&amp;", "&", "&lt;", "<", "&gt;", ">", "&apos;", "'", "&quot;", `"`,
)

func pptxUnescapeXML(s string) string { return xmlEntities.Replace(s) }

func readZipEntry(r *zip.Reader, name string) []byte {
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil
			}
			defer rc.Close()
			data, _ := io.ReadAll(rc)
			return data
		}
	}
	return nil
}

func pptxSlideFiles(r *zip.Reader) []string {
	var names []string
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") &&
			!strings.Contains(f.Name[len("ppt/slides/"):], "/") {
			names = append(names, f.Name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		ni, nj := 0, 0
		if m := pptxSlideNumRegex.FindStringSubmatch(names[i]); len(m) > 1 {
			ni, _ = strconv.Atoi(m[1])
		}
		if m := pptxSlideNumRegex.FindStringSubmatch(names[j]); len(m) > 1 {
			nj, _ = strconv.Atoi(m[1])
		}
		return ni < nj
	})
	return names
}

// pptxSlideIsHidden returns true if the slide XML has show="0" or show="false" on the root <p:sld> (OOXML: hidden/skip in show).
func pptxSlideIsHidden(slideXML []byte) bool {
	m := pptxSldShowAttrRe.FindSubmatch(slideXML)
	if len(m) < 2 {
		return false
	}
	attr := string(m[1])
	return strings.Contains(attr, `show="0"`) || strings.Contains(attr, `show='0'`) ||
		strings.Contains(attr, `show="false"`) || strings.Contains(attr, `show='false'`)
}

// splitPPTXShapes extracts individual <p:sp>...</p:sp> blocks from slide XML.
func splitPPTXShapes(content string) []string {
	var shapes []string
	rest := content
	for {
		start := strings.Index(rest, "<p:sp")
		if start == -1 {
			break
		}
		// Verify it's really a <p:sp> tag (not <p:spTree> etc.)
		if len(rest) > start+5 {
			c := rest[start+5]
			if c != ' ' && c != '>' && c != '\t' && c != '\n' && c != '\r' {
				rest = rest[start+5:]
				continue
			}
		}
		end := strings.Index(rest[start:], "</p:sp>")
		if end == -1 {
			break
		}
		shapes = append(shapes, rest[start:start+end+7])
		rest = rest[start+end+7:]
	}
	return shapes
}

// parsePPTXSlideText extracts title and body text from a PPTX slide XML.
func parsePPTXSlideText(data []byte) (title, body string) {
	content := string(data)
	shapes := splitPPTXShapes(content)

	var bodyParts []string
	for _, shape := range shapes {
		// Skip system placeholders (date, footer, slide number)
		if pptxSystemPhRe.MatchString(shape) {
			continue
		}
		matches := pptxATextRe.FindAllStringSubmatch(shape, -1)
		var texts []string
		for _, m := range matches {
			t := strings.TrimSpace(pptxUnescapeXML(m[1]))
			if t != "" {
				texts = append(texts, t)
			}
		}
		if len(texts) == 0 {
			continue
		}
		shapeText := strings.Join(texts, "")
		if pptxTitlePhRe.MatchString(shape) && title == "" {
			title = strings.TrimSpace(shapeText)
		} else {
			bodyParts = append(bodyParts, strings.TrimSpace(shapeText))
		}
	}
	body = strings.Join(bodyParts, "\n")
	return
}

// parsePPTXNotesText extracts speaker notes from a notesSlide XML.
func parsePPTXNotesText(data []byte) string {
	content := string(data)
	shapes := splitPPTXShapes(content)
	for _, shape := range shapes {
		if pptxSystemPhRe.MatchString(shape) {
			continue
		}
		// Notes body is typically idx="1" or type="body"
		if !pptxNotesBodyPhRe.MatchString(shape) {
			continue
		}
		matches := pptxATextRe.FindAllStringSubmatch(shape, -1)
		var texts []string
		for _, m := range matches {
			t := strings.TrimSpace(pptxUnescapeXML(m[1]))
			if t != "" {
				texts = append(texts, t)
			}
		}
		if len(texts) > 0 {
			return strings.TrimSpace(strings.Join(texts, "\n"))
		}
	}
	return ""
}

// findNotesFileInZip resolves the notes slide file for a given slide via its _rels file.
func findNotesFileInZip(r *zip.Reader, slideFile string) string {
	// e.g. slideFile = "ppt/slides/slide1.xml"
	// rels  = "ppt/slides/_rels/slide1.xml.rels"
	dir := path.Dir(slideFile)
	base := path.Base(slideFile)
	relsPath := dir + "/_rels/" + base + ".rels"

	relsData := readZipEntry(r, relsPath)
	if relsData == nil {
		return ""
	}
	m := pptxNotesRelRe.FindStringSubmatch(string(relsData))
	if len(m) < 2 {
		return ""
	}
	// Target is relative to the slide directory, e.g. "../notesSlides/notesSlide1.xml"
	notesPath := path.Clean(dir + "/" + m[1])
	return notesPath
}

func downloadAndParsePPTX(presentationID string) (*pptxPresentation, error) {
	exportURL := fmt.Sprintf(
		"https://docs.google.com/presentation/d/%s/export/pptx",
		presentationID,
	)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			if strings.Contains(req.URL.Host, "accounts.google.com") {
				return fmt.Errorf("presentation requires authentication — make sure it is shared as 'Anyone with the link can view'")
			}
			return nil
		},
	}
	resp, err := client.Get(exportURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("failed to download presentation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("access denied (HTTP %d) — share the presentation as 'Anyone with the link can view'", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed (HTTP %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PPTX data: %w", err)
	}

	return parsePPTX(data)
}

func parsePPTX(data []byte) (*pptxPresentation, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid PPTX file: %w", err)
	}

	// Get presentation title from docProps/core.xml
	title := "Imported Presentation"
	if coreXML := readZipEntry(zr, "docProps/core.xml"); coreXML != nil {
		if m := pptxDocTitleRe.FindSubmatch(coreXML); len(m) > 1 {
			if t := strings.TrimSpace(pptxUnescapeXML(string(m[1]))); t != "" {
				title = t
			}
		}
	}

	// Slide object IDs from ppt/presentation.xml (same order as slide files)
	var slideIDs []string
	if presXML := readZipEntry(zr, "ppt/presentation.xml"); presXML != nil {
		for _, m := range pptxSldIdRe.FindAllStringSubmatch(string(presXML), -1) {
			if len(m) > 1 {
				slideIDs = append(slideIDs, m[1])
			}
		}
	}

	slideFiles := pptxSlideFiles(zr)
	var slides []pptxSlideData
	for i, sf := range slideFiles {
		slideXML := readZipEntry(zr, sf)
		if slideXML == nil {
			continue
		}
		// Hidden: OOXML stores show="0" or show="false" on <p:sld> in each slide XML (not in presentation.xml)
		hidden := pptxSlideIsHidden(slideXML)

		slideTitle, slideBody := parsePPTXSlideText(slideXML)

		var notes string
		notesFile := findNotesFileInZip(zr, sf)
		if notesFile != "" {
			if notesXML := readZipEntry(zr, notesFile); notesXML != nil {
				notes = parsePPTXNotesText(notesXML)
			}
		}

		images := extractSlideImages(zr, sf)

		slides = append(slides, pptxSlideData{
			Index:  i + 1,
			Title:  slideTitle,
			Body:   slideBody,
			Notes:  notes,
			Images: images,
			Hidden: hidden,
		})
	}

	return &pptxPresentation{Title: title, Slides: slides, SlideIDs: slideIDs}, nil
}

// ---- Google Slides REST API (optional — used when API key provided for thumbnails & comments) ----

type gSlidePresentation struct {
	PresentationID string   `json:"presentationId"`
	Title          string   `json:"title"`
	Slides         []gSlide `json:"slides"`
}

type gSlide struct {
	ObjectID        string           `json:"objectId"`
	PageElements    []gPageElement   `json:"pageElements"`
	SlideProperties gSlideProperties `json:"slideProperties"`
}

type gSlideProperties struct {
	NotesPage gNotesPage `json:"notesPage"`
}

type gNotesPage struct {
	PageElements []gPageElement `json:"pageElements"`
}

type gPageElement struct {
	ObjectID string  `json:"objectId"`
	Shape    *gShape `json:"shape"`
}

type gShape struct {
	ShapeType   string        `json:"shapeType"`
	Placeholder *gPlaceholder `json:"placeholder"`
	Text        *gTextContent `json:"text"`
}

type gPlaceholder struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type gTextContent struct {
	TextElements []gTextElement `json:"textElements"`
}

type gTextElement struct {
	TextRun *gTextRun `json:"textRun"`
}

type gTextRun struct {
	Content string `json:"content"`
}

type gThumbnail struct {
	ContentURL string `json:"contentUrl"`
}

type gDriveCommentsResponse struct {
	Comments []gDriveComment `json:"comments"`
}

type gDriveComment struct {
	ID       string        `json:"id"`
	Content  string        `json:"content"`
	Author   gDriveAuthor  `json:"author"`
	Resolved bool          `json:"resolved"`
	Replies  []gDriveReply `json:"replies"`
}

type gDriveAuthor struct {
	DisplayName string `json:"displayName"`
}

type gDriveReply struct {
	Content string       `json:"content"`
	Author  gDriveAuthor `json:"author"`
}

// googleAPIError is the JSON shape of Google API error responses.
type googleAPIError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func fetchGoogleSlidesAPI(presentationID, apiKey string) (*gSlidePresentation, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s?key=%s", presentationID, apiKey)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var apiErr googleAPIError
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.Error.Message != "" {
			msg := apiErr.Error.Message
			// Google Slides API does not accept API keys; only OAuth2. Return a clear message.
			if strings.Contains(msg, "API keys are not supported") || strings.Contains(msg, "Expected OAuth2") {
				msg = "Google Slides API ไม่รองรับ API Key — ต้องใช้ OAuth2 (ล็อกอินด้วย Google). ตอนนี้ใช้โหมด PPTX ได้โดยไม่ต้องใส่ Key (ได้ text, notes, รูปฝังใน slide)."
			}
			return nil, fmt.Errorf("%s", msg)
		}
		return nil, fmt.Errorf("Google Slides API error %d: %s", resp.StatusCode, string(body))
	}
	var p gSlidePresentation
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &p, nil
}

func fetchSlideThumbnail(presentationID, pageObjectID, apiKey string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://slides.googleapis.com/v1/presentations/%s/pages/%s/thumbnail?key=%s&thumbnailProperties.mimeType=PNG&thumbnailProperties.thumbnailSize=LARGE",
		presentationID, pageObjectID, apiKey,
	)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("thumbnail API error %d", resp.StatusCode)
	}
	var t gThumbnail
	if err := json.Unmarshal(body, &t); err != nil {
		return "", err
	}
	return t.ContentURL, nil
}

// downloadThumbnailAsDataURL fetches the thumbnail image from Google's contentUrl and returns a data URL
// so it can be stored in the task and displayed without expiry (contentUrl lifetime is ~30 min).
// The thumbnail includes the full rendered slide (shapes, lines, drawings).
// apiKey is optional; when set, it is appended to the URL so Google may allow the server-side download.
func downloadThumbnailAsDataURL(contentURL, apiKey string) (string, error) {
	if contentURL == "" {
		return "", fmt.Errorf("empty content URL")
	}
	downloadURL := contentURL
	if apiKey != "" {
		if strings.Contains(contentURL, "?") {
			downloadURL = contentURL + "&key=" + apiKey
		} else {
			downloadURL = contentURL + "?key=" + apiKey
		}
	}
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Sentinel-Slides-Import/1.0 (https://github.com/portnd/the-sentinel-core)")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Slides Import] thumbnail download request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		snippet := string(body)
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		log.Printf("[Slides Import] thumbnail download failed: status=%d body=%q", resp.StatusCode, snippet)
		return "", fmt.Errorf("thumbnail download returned %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Slides Import] thumbnail response read failed: %v", err)
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(body)
	return "data:image/png;base64," + b64, nil
}

func fetchDriveComments(fileID, apiKey string) ([]gDriveComment, error) {
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/drive/v3/files/%s/comments?key=%s&fields=comments(id,content,anchor,author,resolved,replies)&pageSize=100",
		fileID, apiKey,
	)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Drive API error %d: %s", resp.StatusCode, string(body))
	}
	var result gDriveCommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Comments, nil
}

func extractAPISlideTitle(slide gSlide) string {
	for _, el := range slide.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		if el.Shape.Placeholder != nil {
			pt := el.Shape.Placeholder.Type
			if pt == "CENTERED_TITLE" || pt == "TITLE" {
				var sb strings.Builder
				for _, te := range el.Shape.Text.TextElements {
					if te.TextRun != nil {
						sb.WriteString(te.TextRun.Content)
					}
				}
				if t := strings.TrimSpace(strings.ReplaceAll(sb.String(), "\n", " ")); t != "" {
					return t
				}
			}
		}
	}
	return ""
}

func extractAPISlideBody(slide gSlide) string {
	var parts []string
	for _, el := range slide.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		var sb strings.Builder
		for _, te := range el.Shape.Text.TextElements {
			if te.TextRun != nil {
				sb.WriteString(te.TextRun.Content)
			}
		}
		if t := strings.TrimSpace(sb.String()); t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, "\n")
}

func extractAPISpeakerNotes(slide gSlide) string {
	for _, el := range slide.SlideProperties.NotesPage.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		if el.Shape.Placeholder != nil && el.Shape.Placeholder.Type == "BODY" {
			var sb strings.Builder
			for _, te := range el.Shape.Text.TextElements {
				if te.TextRun != nil {
					sb.WriteString(te.TextRun.Content)
				}
			}
			if t := strings.TrimSpace(sb.String()); t != "" {
				return t
			}
		}
	}
	return ""
}

// ---- Preview: get slide list only (no thumbnails, no task creation) ----

func (u *sentinelUsecase) PreviewGoogleSlides(req *domain.PreviewGoogleSlidesRequest, serverAPIKey string) (*domain.PreviewGoogleSlidesResult, error) {
	presentationID, err := extractPresentationID(req.PresentationURL)
	if err != nil {
		return nil, err
	}
	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = serverAPIKey
	}

	title, items, importMode, apiKeyStatus, apiKeyErrMsg, err := getSlidesListOnly(presentationID, apiKey)
	if err != nil {
		return nil, err
	}
	alreadyImported, _ := u.repo.GetImportedSlideIndicesByPresentationID(presentationID)
	return &domain.PreviewGoogleSlidesResult{
		PresentationTitle:           title,
		PresentationID:               presentationID,
		Slides:                       items,
		AlreadyImportedSlideIndices:  alreadyImported,
		ImportMode:                   importMode,
		APIKeyStatus:                 apiKeyStatus,
		APIKeyError:                  apiKeyErrMsg,
	}, nil
}

// getSlidesListOnly returns presentation title, slide list, import mode, API key status, and optional error message. No thumbnails.
func getSlidesListOnly(presentationID, apiKey string) (title string, slides []domain.PreviewSlideItem, importMode, apiKeyStatus, apiKeyError string, err error) {
	var presentationTitle string
	apiKeyProvided := apiKey != ""

	pptxData, pptxErr := downloadAndParsePPTX(presentationID)
	if pptxErr != nil && !apiKeyProvided {
		return "", nil, "", "", "", fmt.Errorf("failed to download presentation: %w\nTip: ensure the presentation is shared as 'Anyone with the link can view'", pptxErr)
	}
	if pptxErr == nil {
		presentationTitle = pptxData.Title
		for _, s := range pptxData.Slides {
			t := s.Title
			if t == "" {
				t = fmt.Sprintf("Slide %d", s.Index)
			}
			slides = append(slides, domain.PreviewSlideItem{Index: s.Index, Title: t, Hidden: s.Hidden})
		}
	}
	if apiKeyProvided {
		apiPresentation, apiErr := fetchGoogleSlidesAPI(presentationID, apiKey)
		if apiErr == nil {
			apiKeyStatus = "valid"
			if pptxErr != nil {
				importMode = "api_only"
				presentationTitle = apiPresentation.Title
				slides = nil
				for i, slide := range apiPresentation.Slides {
					t := extractAPISlideTitle(slide)
					if t == "" {
						t = fmt.Sprintf("Slide %d", i+1)
					}
					slides = append(slides, domain.PreviewSlideItem{Index: i + 1, Title: t, Hidden: false})
				}
			} else {
				importMode = "pptx_with_api"
			}
		} else {
			apiKeyStatus = "invalid"
			apiKeyError = apiErr.Error()
			if pptxErr != nil {
				return "", nil, "", "", "", fmt.Errorf("failed to download PPTX (%v) and Slides API failed: %w", pptxErr, apiErr)
			}
			importMode = "pptx_only"
		}
	} else {
		apiKeyStatus = "not_provided"
		importMode = "pptx_only"
	}
	if len(slides) == 0 {
		return "", nil, "", "", "", errors.New("no slides found in the presentation")
	}
	return presentationTitle, slides, importMode, apiKeyStatus, apiKeyError, nil
}

// ---- Main ImportFromGoogleSlides usecase ----

func (u *sentinelUsecase) ImportFromGoogleSlides(req *domain.ImportGoogleSlidesRequest, serverAPIKey string, creatorID uint) (*domain.ImportGoogleSlidesResult, error) {
	presentationID, err := extractPresentationID(req.PresentationURL)
	if err != nil {
		return nil, err
	}

	var sprintUUID *uuid.UUID
	if req.SprintID != "" {
		parsed, err := uuid.Parse(req.SprintID)
		if err != nil {
			return nil, fmt.Errorf("invalid sprint_id: %w", err)
		}
		sprintUUID = &parsed
	}
	var epicUUID *uuid.UUID
	if req.EpicID != "" {
		parsed, err := uuid.Parse(req.EpicID)
		if err != nil {
			return nil, fmt.Errorf("invalid epic_id: %w", err)
		}
		epicUUID = &parsed
	}
	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = serverAPIKey
	}

	// Step 1: Get slide content.
	// Primary: PPTX export — works for any public "anyone with link" presentation, no API key needed.
	// If API key provided: also fetch thumbnails and Drive comments.
	var slides []slideInfo
	var presentationTitle string

	pptxData, pptxErr := downloadAndParsePPTX(presentationID)
	if pptxErr != nil && apiKey == "" {
		return nil, fmt.Errorf("failed to download presentation: %w\nTip: ensure the presentation is shared as 'Anyone with the link can view'", pptxErr)
	}

	if pptxErr == nil {
		// Successfully parsed PPTX
		presentationTitle = pptxData.Title
		for _, s := range pptxData.Slides {
			title := s.Title
			if title == "" {
				title = fmt.Sprintf("Slide %d", s.Index)
			}
			// ใช้เฉพาะรูปที่ฝังใน slide นั้น (จาก PPTX media) — ไม่เรียก export ต่อ slide เพราะมักได้รูป slide แรกซ้ำทุกหน้า
			slides = append(slides, slideInfo{
				Index:  s.Index,
				Title:  title,
				Body:   s.Body,
				Notes:  s.Notes,
				Images: s.Images,
			})
		}
	}

		// Step 2: If API key available, enhance with thumbnails (via Slides API) and comments (via Drive API).
	var allComments []domain.SlideComment
	if apiKey != "" {
		log.Printf("[Slides Import] API key present, calling Slides API for thumbnails/comments")
		apiPresentation, apiErr := fetchGoogleSlidesAPI(presentationID, apiKey)
		if apiErr == nil {
			log.Printf("[Slides Import] Slides API OK, slides=%d", len(apiPresentation.Slides))
			// If PPTX also failed, use the API data as the content source
			if pptxErr != nil {
				presentationTitle = apiPresentation.Title
				slides = nil
				for i, slide := range apiPresentation.Slides {
					title := extractAPISlideTitle(slide)
					if title == "" {
						title = fmt.Sprintf("Slide %d", i+1)
					}
					slides = append(slides, slideInfo{
						Index:      i + 1,
						Title:      title,
						Body:       extractAPISlideBody(slide),
						Notes:      extractAPISpeakerNotes(slide),
						SlideObjID: slide.ObjectID,
					})
				}
			} else {
				// Merge: fill in slideObjID from API data (same order)
				for i := range slides {
					if i < len(apiPresentation.Slides) {
						slides[i].SlideObjID = apiPresentation.Slides[i].ObjectID
					}
				}
			}

			// Fetch per-slide thumbnails and download as base64 so drawings/lines are persisted (contentUrl expires in ~30 min)
			withObjID := 0
			for i := range slides {
				if slides[i].SlideObjID != "" {
					withObjID++
				}
			}
			log.Printf("[Slides Import] fetching thumbnails: %d slides with SlideObjID", withObjID)
			for i := range slides {
				if slides[i].SlideObjID == "" {
					log.Printf("[Slides Import] slide %d: skip (no SlideObjID)", i+1)
					continue
				}
				url, err := fetchSlideThumbnail(presentationID, slides[i].SlideObjID, apiKey)
				if err != nil {
					log.Printf("[Slides Import] slide %d: thumbnail URL failed: %v", i+1, err)
					continue
				}
				slides[i].ThumbnailURL = url
				dataURL, dlErr := downloadThumbnailAsDataURL(url, apiKey)
				if dlErr != nil && apiKey != "" {
					dataURL, dlErr = downloadThumbnailAsDataURL(url, "") // fallback: CDN may not accept key param
				}
				if dlErr == nil {
					// Prepend so the first image shown is the full slide with shapes/lines (กรอบแดง, เส้นวาด)
					slides[i].Images = append([]string{dataURL}, slides[i].Images...)
					log.Printf("[Slides Import] slide %d: thumbnail OK (base64 prepended)", i+1)
				} else {
					log.Printf("[Slides Import] slide %d: thumbnail download failed (กรอบ/เส้น will be missing): %v", i+1, dlErr)
				}
			}

			// Fetch Drive comments (non-fatal)
			driveComments, _ := fetchDriveComments(presentationID, apiKey)
			for _, c := range driveComments {
				comment := domain.SlideComment{
					Content:  c.Content,
					Author:   c.Author.DisplayName,
					Resolved: c.Resolved,
				}
				for _, r := range c.Replies {
					comment.Content += fmt.Sprintf("\n  ↳ [%s]: %s", r.Author.DisplayName, r.Content)
				}
				allComments = append(allComments, comment)
			}
		} else if pptxErr != nil {
			// Both PPTX and API failed
			return nil, fmt.Errorf("failed to download PPTX (%v) and Slides API also failed: %w", pptxErr, apiErr)
		} else {
			log.Printf("[Slides Import] Slides API failed (thumbnails/comments skipped): %v", apiErr)
		}
		// If API key provided but API call failed, we still continue with PPTX data (no thumbnails/comments)
	} else {
		log.Printf("[Slides Import] no API key: thumbnails (กรอบ/เส้น) and Drive comments will not be fetched")
	}

	if len(slides) == 0 {
		return nil, errors.New("no slides found in the presentation")
	}

	// Step 3: Validate priority and story points
	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if !map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}[priority] {
		priority = "MEDIUM"
	}
	storyPoints := req.StoryPoints
	if storyPoints < 0 {
		storyPoints = 0
	}

	slug := "task"
	proj, err := u.repo.GetProjectByID(projectUUID)
	if err == nil && proj != nil {
		slug = slugify(proj.Name)
	}

	// Filter by selected slide indices if provided (1-based)
	if len(req.SlideIndices) > 0 {
		allowed := make(map[int]bool)
		for _, idx := range req.SlideIndices {
			allowed[idx] = true
		}
		filtered := slides[:0]
		for _, s := range slides {
			if allowed[s.Index] {
				filtered = append(filtered, s)
			}
		}
		slides = filtered
	}
	if len(slides) == 0 {
		return nil, errors.New("no slides selected to import")
	}

	// Next code numbers: use global max suffix so codes are unique across all projects (idx_tasks_code is global).
	maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(slug)

	// Step 4: Create one task per (filtered) slide
	var createdTasks []*domain.Task
	for i, slide := range slides {
		// Build description as HTML so images appear in Description (single place); no separate Slide Images section.
		var htmlParts []string
		if slide.Body != "" {
			htmlParts = append(htmlParts, "<p>"+html.EscapeString(slide.Body)+"</p>")
		}
		for _, imgSrc := range slide.Images {
			// imgSrc is either data URL (base64) or https URL; safe to use in src
			htmlParts = append(htmlParts, "<p><img src=\""+html.EscapeString(imgSrc)+"\" alt=\"Slide\" /></p>")
		}
		if slide.Notes != "" {
			htmlParts = append(htmlParts, "<p><em>Speaker Notes:</em> "+html.EscapeString(slide.Notes)+"</p>")
		}
		description := strings.Join(htmlParts, "\n")
		if description == "" {
			description = "<p></p>"
		}

		var slideURL string
		if slide.SlideObjID != "" {
			// Prefer object ID so link stays correct if slides are reordered
			slideURL = fmt.Sprintf("https://docs.google.com/presentation/d/%s/edit#slide=id.%s", presentationID, slide.SlideObjID)
		} else {
			// Fallback: use 1-based slide index; Google rewrites to internal ID and navigates to that position
			slideURL = fmt.Sprintf("https://docs.google.com/presentation/d/%s/edit#slide=%d", presentationID, slide.Index)
		}

		// resource_urls: keep only metadata for "Open in Slides"; images are now in description
		resourceURLs := domain.SlideResourceURLs{
			ThumbnailURL:   "",   // no longer used; images in description
			Images:         nil,  // no duplicate; images in description
			SlideURL:       slideURL,
			Source:         "google_slides",
			SlideIndex:     slide.Index,
			PresentationID: presentationID,
			Comments:       allComments,
		}
		resourceURLsJSON, err := json.Marshal(resourceURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource URLs for slide %d: %w", slide.Index, err)
		}

		code := fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i)

		projectIDCopy := projectUUID
		taskTitle := fmt.Sprintf("Slide %d: %s", slide.Index, slide.Title)
		task := &domain.Task{
			ID:           uuid.New(),
			Code:         code,
			Title:        taskTitle,
			Description:  description,
			CreatedBy:    &creatorID,
			Status:       "PENDING",
			Priority:     priority,
			StoryPoints:  storyPoints,
			SprintID:     sprintUUID, // nil when importing to backlog
			EpicID:       epicUUID,  // set when importing to a specific epic
			ProjectID:    &projectIDCopy,
			ResourceURLs: datatypes.JSON(resourceURLsJSON),
		}

		if err := u.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("failed to create task for slide %d: %w", slide.Index, err)
		}
		createdTasks = append(createdTasks, task)
	}

	return &domain.ImportGoogleSlidesResult{
		CreatedCount:      len(createdTasks),
		SlideCount:        len(slides),
		PresentationTitle: presentationTitle,
		Tasks:             createdTasks,
	}, nil
}

// --- Epic Operations (Hierarchy Dimension 1) ---

func (u *sentinelUsecase) CreateEpic(projectID uuid.UUID, title, description, color string, startDate, endDate *time.Time) (*domain.Epic, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("epic title is required")
	}
	if color == "" {
		color = "#6366f1"
	}
	epic := &domain.Epic{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      "PLANNING",
		Color:       color,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	if err := u.repo.CreateEpic(epic); err != nil {
		return nil, fmt.Errorf("failed to create epic: %w", err)
	}
	return epic, nil
}

func (u *sentinelUsecase) GetEpicsByProject(projectID uuid.UUID) ([]domain.Epic, error) {
	return u.repo.GetEpicsByProjectID(projectID)
}

func (u *sentinelUsecase) UpdateEpic(epicID uuid.UUID, title, description, status, color string, sortOrder *int, startDate, endDate *time.Time) (*domain.Epic, error) {
	epic, err := u.repo.GetEpicByID(epicID)
	if err != nil {
		return nil, fmt.Errorf("epic not found: %w", err)
	}
	if title != "" {
		epic.Title = title
	}
	if description != "" {
		epic.Description = description
	}
	if status != "" {
		epic.Status = status
	}
	if color != "" {
		epic.Color = color
	}
	if sortOrder != nil {
		epic.SortOrder = *sortOrder
	}
	if startDate != nil {
		epic.StartDate = startDate
	}
	if endDate != nil {
		epic.EndDate = endDate
	}
	if err := u.repo.UpdateEpic(epic); err != nil {
		return nil, fmt.Errorf("failed to update epic: %w", err)
	}
	return epic, nil
}

func (u *sentinelUsecase) DeleteEpic(epicID uuid.UUID) error {
	return u.repo.DeleteEpic(epicID)
}

func (u *sentinelUsecase) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	return u.repo.GetEpicTimelineData(projectID)
}

func (u *sentinelUsecase) GetSprintTimelineData(projectID uuid.UUID) (*domain.SprintTimelineData, error) {
	return u.repo.GetSprintTimelineData(projectID)
}

// ─── 2-Page Client Report:
//   Page 1 (Portrait): งวด/Milestone → Epic → Task list with dates
//   Page 2 (Landscape): Gantt chart (Sprint mode, exact system colours)
// ──────────────────────────────────────────────────────────────────────────────

// p1Task is one deliverable row under an epic.
type p1Task struct {
	Title     string
	StartDate string
	EndDate   string
}

// p1Epic is an epic block (header + tasks) inside a milestone group.
type p1Epic struct {
	Title string
	Tasks []p1Task
}

// p1Milestone is one งวด/delivery group.
type p1Milestone struct {
	Number  int
	Title   string
	DueDate string
	Epics   []p1Epic
	Count   int // total tasks in this milestone
}

// p2GanttCol is a month header column with percentage-based positioning.
type p2GanttCol struct {
	Label    string
	LeftPct  float64
	WidthPct float64
}

// p2GanttRow is one row (sprint header or task) in the Gantt chart.
type p2GanttRow struct {
	IsSprint  bool
	Label     string
	BarLeft   float64 // % from left edge of chart area
	BarWidth  float64 // % width
	BarLabel  string
	HasBar    bool
}

// ganttMonthRow is one task row for the month-timeframe Gantt chart.
// Bar position uses actual dates so bar width reflects real duration (e.g. 2 weeks ≠ full 2 months).
type ganttMonthRow struct {
	Label      string
	EpicTitle  string
	StartMonth int     // 0-based index (kept for sort/display)
	EndMonth   int
	StartDate  string
	EndDate    string
	BarLeftPct float64 // 0–100: position of bar start as % of chart timeline
	BarWidthPct float64 // 0–100: bar width as % of chart timeline (real duration)
}

// clientReportData is the full payload injected into the HTML template.
type clientReportData struct {
	ProjectName string
	GeneratedAt string
	// Page 1
	MilestoneGroups []p1Milestone
	HasUnassigned   bool
	UnassignedEpics []p1Epic
	// Page 2: day-scale Gantt (sprint)
	GanttCols    []p2GanttCol
	GanttRows    []p2GanttRow
	HasGanttData bool
	// Gantt by month (timeframe เดือน)
	GanttMonthLabels []string
	GanttMonthRows   []ganttMonthRow
	HasGanttMonth    bool
}

func fmtDate(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("02 Jan 2006")
}

func fmtDateStr(s string) string {
	if s == "" {
		return "—"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.Format("02 Jan 2006")
}

// daysBetween returns the number of days between two times (can be negative).
func daysBetween(a, b time.Time) float64 {
	return b.Sub(a).Hours() / 24
}

// truncToMonth returns the first day of the month containing t.
func truncToMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ExportTimelinePDF generates the 2-page client report via chromedp (same pattern as mims).
func (u *sentinelUsecase) ExportTimelinePDF(projectID uuid.UUID, _ string, templateDir string) ([]byte, string, error) {
	project, err := u.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("project not found: %w", err)
	}

	// ── Fetch data ──────────────────────────────────────────────────────────────
	epicData, err := u.repo.GetEpicTimelineData(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("get epic timeline: %w", err)
	}
	sprintData, err := u.repo.GetSprintTimelineData(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("get sprint timeline: %w", err)
	}
	milestones, _ := u.repo.GetMilestonesByProjectID(projectID) // optional

	// ── Page 1: group tasks by Milestone → Epic ─────────────────────────────────
	// Sort milestones by DueDate ascending
	sort.Slice(milestones, func(i, j int) bool {
		if milestones[i].DueDate == nil {
			return false
		}
		if milestones[j].DueDate == nil {
			return true
		}
		return milestones[i].DueDate.Before(*milestones[j].DueDate)
	})

	// milestoneMap: milestoneID → index in milestones slice (-1 = unassigned)
	milestoneIndexByID := map[string]int{}
	for i, m := range milestones {
		milestoneIndexByID[m.ID.String()] = i
	}

	// Collect all tasks from epics, annotated with their epic title.
	type taskWithEpic struct {
		epicTitle string
		task      domain.Task
	}
	var allTasksWithEpic []taskWithEpic
	for _, e := range epicData.Epics {
		for _, t := range e.Tasks {
			allTasksWithEpic = append(allTasksWithEpic, taskWithEpic{epicTitle: e.Title, task: t})
		}
	}

	// For each milestone, maintain ordered epic groups.
	type epicGroupBuilder struct {
		order  []string         // epic title order
		tasks  map[string][]p1Task
	}
	milestoneBuilders := make([]epicGroupBuilder, len(milestones))
	for i := range milestoneBuilders {
		milestoneBuilders[i].tasks = map[string][]p1Task{}
	}
	var unassignedBuilder epicGroupBuilder
	unassignedBuilder.tasks = map[string][]p1Task{}

	assignTask := func(idx int, epicTitle string, t domain.Task) {
		taskEnd := coalesce(t.EndDate, t.DueAt)
		row := p1Task{Title: t.Title, StartDate: fmtDate(t.StartDate), EndDate: fmtDate(taskEnd)}
		if idx < 0 {
			if _, seen := unassignedBuilder.tasks[epicTitle]; !seen {
				unassignedBuilder.order = append(unassignedBuilder.order, epicTitle)
			}
			unassignedBuilder.tasks[epicTitle] = append(unassignedBuilder.tasks[epicTitle], row)
			return
		}
		b := &milestoneBuilders[idx]
		if _, seen := b.tasks[epicTitle]; !seen {
			b.order = append(b.order, epicTitle)
		}
		b.tasks[epicTitle] = append(b.tasks[epicTitle], row)
	}

	for _, tw := range allTasksWithEpic {
		t := tw.task
		// 1) Direct milestone_id link
		if t.MilestoneID != nil {
			if idx, ok := milestoneIndexByID[t.MilestoneID.String()]; ok {
				assignTask(idx, tw.epicTitle, t)
				continue
			}
		}
		// 2) Date-based: earliest milestone whose DueDate >= task end date
		taskEnd := coalesce(t.EndDate, t.DueAt)
		assigned := false
		if taskEnd != nil {
			for idx, m := range milestones {
				if m.DueDate != nil && !m.DueDate.Before(*taskEnd) {
					assignTask(idx, tw.epicTitle, t)
					assigned = true
					break
				}
			}
		}
		if !assigned {
			assignTask(-1, tw.epicTitle, t)
		}
	}

	// Convert builders → []p1Milestone
	var mGroups []p1Milestone
	for i, m := range milestones {
		b := milestoneBuilders[i]
		var epics []p1Epic
		count := 0
		for _, et := range b.order {
			tasks := b.tasks[et]
			epics = append(epics, p1Epic{Title: et, Tasks: tasks})
			count += len(tasks)
		}
		mGroups = append(mGroups, p1Milestone{
			Number:  i + 1,
			Title:   m.Title,
			DueDate: fmtDate(m.DueDate),
			Epics:   epics,
			Count:   count,
		})
	}
	var unassignedEpics []p1Epic
	for _, et := range unassignedBuilder.order {
		unassignedEpics = append(unassignedEpics, p1Epic{Title: et, Tasks: unassignedBuilder.tasks[et]})
	}

	// ── Gantt by month (timeframe เดือน) ───────────────────────────────────────────
	var ganttMonthLabels []string
	var ganttMonthRows []ganttMonthRow
	hasGanttMonth := false
	var minMonth, maxMonth time.Time
	for _, e := range epicData.Epics {
		for _, t := range e.Tasks {
			start := t.StartDate
			end := coalesce(t.EndDate, t.DueAt)
			if start == nil || end == nil {
				continue
			}
			sM, eM := truncToMonth(*start), truncToMonth(*end)
			if !hasGanttMonth {
				minMonth, maxMonth = sM, eM
				hasGanttMonth = true
			} else {
				if sM.Before(minMonth) {
					minMonth = sM
				}
				if eM.After(maxMonth) {
					maxMonth = eM
				}
			}
		}
	}
	if hasGanttMonth {
		var monthList []time.Time
		for m := minMonth; !m.After(maxMonth); m = m.AddDate(0, 1, 0) {
			monthList = append(monthList, m)
		}
		monthIndex := make(map[string]int)
		for i, m := range monthList {
			monthIndex[m.Format("2006-01")] = i
		}
		for _, lab := range monthList {
			ganttMonthLabels = append(ganttMonthLabels, lab.Format("Jan 06"))
		}
		// Chart range for proportional bar: first day of first month → first day of month after last
		chartStart := minMonth
		chartEnd := maxMonth.AddDate(0, 1, 0)
		chartDurationDays := daysBetween(chartStart, chartEnd)
		if chartDurationDays <= 0 {
			chartDurationDays = 1
		}
		for _, e := range epicData.Epics {
			for _, t := range e.Tasks {
				start := t.StartDate
				end := coalesce(t.EndDate, t.DueAt)
				if start == nil || end == nil {
					continue
				}
				sM, eM := truncToMonth(*start), truncToMonth(*end)
				si, okS := monthIndex[sM.Format("2006-01")]
				ei, okE := monthIndex[eM.Format("2006-01")]
				if !okS || !okE {
					continue
				}
				if ei < si {
					ei = si
				}
				// Bar position by actual dates so width = real duration (e.g. 2 weeks, not full 2 months)
				leftPct := daysBetween(chartStart, *start) / chartDurationDays * 100
				widthPct := daysBetween(*start, *end) / chartDurationDays * 100
				if widthPct <= 0 {
					widthPct = 2 // min visible bar
				}
				if leftPct < 0 {
					leftPct = 0
				}
				if leftPct+widthPct > 100 {
					widthPct = 100 - leftPct
				}
				ganttMonthRows = append(ganttMonthRows, ganttMonthRow{
					Label:       t.Title,
					EpicTitle:   e.Title,
					StartMonth:  si,
					EndMonth:    ei,
					StartDate:   fmtDate(start),
					EndDate:     fmtDate(end),
					BarLeftPct:  leftPct,
					BarWidthPct: widthPct,
				})
			}
		}
		// Sort by start month then by epic/title
		sort.Slice(ganttMonthRows, func(i, j int) bool {
			if ganttMonthRows[i].StartMonth != ganttMonthRows[j].StartMonth {
				return ganttMonthRows[i].StartMonth < ganttMonthRows[j].StartMonth
			}
			if ganttMonthRows[i].EpicTitle != ganttMonthRows[j].EpicTitle {
				return ganttMonthRows[i].EpicTitle < ganttMonthRows[j].EpicTitle
			}
			return ganttMonthRows[i].Label < ganttMonthRows[j].Label
		})
	}

	// ── Page 2: Sprint Gantt ─────────────────────────────────────────────────────
	var chartStart, chartEnd time.Time
	hasGanttData := false

	// Find overall date range
	for _, sp := range sprintData.Sprints {
		for _, d := range []*time.Time{sp.StartDate, sp.EndDate} {
			if d == nil {
				continue
			}
			if !hasGanttData || d.Before(chartStart) {
				chartStart = *d
			}
			if !hasGanttData || d.After(chartEnd) {
				chartEnd = *d
			}
			hasGanttData = true
		}
		for _, t := range sp.Tasks {
			for _, d := range []*time.Time{t.StartDate, t.EndDate, t.DueAt} {
				if d == nil {
					continue
				}
				if d.Before(chartStart) {
					chartStart = *d
				}
				if d.After(chartEnd) {
					chartEnd = *d
				}
			}
		}
	}

	var ganttCols []p2GanttCol
	var ganttRows []p2GanttRow

	if hasGanttData {
		// Padding: 2 days on each side for visual breathing room
		chartStart = chartStart.AddDate(0, 0, -2)
		chartEnd = chartEnd.AddDate(0, 0, 2)
		totalDays := daysBetween(chartStart, chartEnd)

		pct := func(d time.Time) float64 {
			return daysBetween(chartStart, d) / totalDays * 100
		}

		// Month header columns
		for m := truncToMonth(chartStart); !m.After(chartEnd); m = m.AddDate(0, 1, 0) {
			mEnd := m.AddDate(0, 1, 0)
			colStart := m
			if colStart.Before(chartStart) {
				colStart = chartStart
			}
			colEnd := mEnd
			if colEnd.After(chartEnd) {
				colEnd = chartEnd
			}
			leftPct := pct(colStart)
			widthPct := pct(colEnd) - leftPct
			if widthPct < 0.1 {
				continue
			}
			ganttCols = append(ganttCols, p2GanttCol{
				Label:    m.Format("Jan 2006"),
				LeftPct:  leftPct,
				WidthPct: widthPct,
			})
		}

		// Gantt rows: sprint header + tasks
		for _, sp := range sprintData.Sprints {
			// Sprint row
			spRow := p2GanttRow{IsSprint: true, Label: sp.Name}
			if sp.StartDate != nil && sp.EndDate != nil {
				l := pct(*sp.StartDate)
				w := pct(*sp.EndDate) - l
				if w < 0.3 {
					w = 0.3
				}
				spRow.BarLeft, spRow.BarWidth, spRow.BarLabel, spRow.HasBar = l, w, sp.Name, true
			}
			ganttRows = append(ganttRows, spRow)

			// Task rows
			for _, t := range sp.Tasks {
				taskEnd := coalesce(t.EndDate, t.DueAt)
				tRow := p2GanttRow{IsSprint: false, Label: t.Title}
				if t.StartDate != nil && taskEnd != nil {
					l := pct(*t.StartDate)
					w := pct(*taskEnd) - l
					if w < 0.3 {
						w = 0.3
					}
					tRow.BarLeft, tRow.BarWidth, tRow.BarLabel, tRow.HasBar = l, w, t.Title, true
				}
				ganttRows = append(ganttRows, tRow)
			}
		}
	}

	data := clientReportData{
		ProjectName:       project.Name,
		GeneratedAt:       time.Now().Format("2 January 2006"),
		MilestoneGroups:   mGroups,
		HasUnassigned:     len(unassignedEpics) > 0,
		UnassignedEpics:   unassignedEpics,
		GanttCols:         ganttCols,
		GanttRows:         ganttRows,
		HasGanttData:      hasGanttData,
		GanttMonthLabels:  ganttMonthLabels,
		GanttMonthRows:    ganttMonthRows,
		HasGanttMonth:     hasGanttMonth,
	}

	// ── Render HTML template (same as mims InitDataToHtml) ─────────────────────
	tmplPath := templateDir + "timeline_report.html"
	funcMap := template.FuncMap{
		"add":     func(a, b int) int { return a + b },
		"mod":     func(a, b int) int { return a % b },
		"ge":      func(a, b int) bool { return a >= b },
		"le":      func(a, b int) bool { return a <= b },
		"printf2": func(f float64) string { return fmt.Sprintf("%.4f", f) },
	}
	tmpl, err := template.New("timeline_report.html").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return nil, "", fmt.Errorf("parse template: %w", err)
	}
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, "", fmt.Errorf("execute template: %w", err)
	}

	// ── Generate PDF via chromedp (same as mims PrintToPDF) ─────────────────────
	ctx, cancel := chromepdf.NewChromedpContext(context.Background())
	defer cancel()

	var pdfBuf []byte
	if err := chromedp.Run(ctx, chromepdf.PrintToPDF(htmlBuf.String(), &pdfBuf, true)); err != nil {
		return nil, "", fmt.Errorf("chromedp print to pdf: %w", err)
	}

	filename := fmt.Sprintf("project-plan-%s-%s.pdf",
		slugify(project.Name),
		time.Now().Format("20060102"),
	)
	return pdfBuf, filename, nil
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}

func coalesce(a, b *time.Time) *time.Time {
	if a != nil {
		return a
	}
	return b
}
