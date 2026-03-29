package usecase

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
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

const appSettingTeamsFeature = "teams_feature_enabled"

func (u *sentinelUsecase) isTeamsFeatureDisabled() bool {
	val, err := u.authRepo.GetAppSetting(appSettingTeamsFeature)
	if err != nil {
		return false
	}
	return val == "false"
}

func (u *sentinelUsecase) withCallerScope(ctx domain.CallerContext) domain.CallerContext {
	ctx.TeamsFeatureDisabled = u.isTeamsFeatureDisabled()
	return ctx
}

// projectNameEnglishOnly matches letters, digits, spaces, hyphens, underscores (English only)
var projectNameEnglishOnly = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)

func (u *sentinelUsecase) CreateProject(name, description, status string, ctx domain.CallerContext) (*domain.Project, error) {
	ctx = u.withCallerScope(ctx)
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
	// Auto-assign team_id so the project is visible to the creator's team
	if ctx.TeamID != nil {
		p.TeamID = ctx.TeamID
	}
	if err := u.repo.CreateProject(p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	if ctx.TeamsFeatureDisabled && ctx.Role == domain.RolePM && ctx.UserID != 0 {
		if err := u.repo.ReplaceProjectPmAssignments(p.ID, []uint{ctx.UserID}); err != nil {
			return nil, fmt.Errorf("failed to register creating PM as project owner: %w", err)
		}
	}
	return p, nil
}

func (u *sentinelUsecase) GetProjects(ctx domain.CallerContext) ([]domain.Project, error) {
	ctx = u.withCallerScope(ctx)
	return u.repo.GetAllProjects(ctx)
}

func (u *sentinelUsecase) GetProjectDetails(id uuid.UUID, ctx domain.CallerContext) (*domain.Project, error) {
	p, err := u.repo.GetProjectByID(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	return p, nil
}

// GetProjectByIDOrCode retrieves a project by UUID or by code (e.g. mims-hdmap-main).
func (u *sentinelUsecase) GetProjectByIDOrCode(idOrCode string, ctx domain.CallerContext) (*domain.Project, error) {
	idOrCode = strings.TrimSpace(idOrCode)
	if idOrCode == "" {
		return nil, errors.New("project id or code is required")
	}
	ctx = u.withCallerScope(ctx)
	if id, err := uuid.Parse(idOrCode); err == nil {
		return u.GetProjectDetails(id, ctx)
	}
	p, err := u.repo.GetProjectByCode(idOrCode, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	if p == nil {
		return nil, errors.New("project not found")
	}
	return p, nil
}

// GetProjectDetailsPage returns project + tasks + sprints + milestones + epics in one call (1 round-trip vs 5).
func (u *sentinelUsecase) GetProjectDetailsPage(idOrCode string, ctx domain.CallerContext) (*domain.ProjectDetailsResponse, error) {
	p, err := u.GetProjectByIDOrCode(idOrCode, ctx)
	if err != nil || p == nil {
		return nil, err
	}
	// Fetch all child data in parallel (4 queries → 1 DB round-trip per type; network already 1 round-trip).
	type result struct {
		tasks     []domain.Task
		sprints   []domain.Sprint
		milestones []domain.Milestone
		epics     []domain.Epic
	}
	var res result
	var errTasks, errSprints, errMilestones, errEpics error
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		res.tasks, errTasks = u.repo.GetTasksByProjectID(p.ID)
	}()
	go func() {
		defer wg.Done()
		res.sprints, errSprints = u.repo.GetSprintsByProjectID(p.ID)
	}()
	go func() {
		defer wg.Done()
		res.milestones, errMilestones = u.repo.GetMilestonesByProjectID(p.ID)
	}()
	go func() {
		defer wg.Done()
		res.epics, errEpics = u.repo.GetEpicsByProjectID(p.ID)
	}()
	wg.Wait()
	if errTasks != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", errTasks)
	}
	if errSprints != nil {
		return nil, fmt.Errorf("failed to load sprints: %w", errSprints)
	}
	if errMilestones != nil {
		return nil, fmt.Errorf("failed to load milestones: %w", errMilestones)
	}
	if errEpics != nil {
		return nil, fmt.Errorf("failed to load epics: %w", errEpics)
	}
	return &domain.ProjectDetailsResponse{
		Project:    p,
		Tasks:      res.tasks,
		Sprints:    res.sprints,
		Milestones: res.milestones,
		Epics:      res.epics,
	}, nil
}

// AssignProjectTeam sets or clears the team for a project (CEO or MANAGER only).
func (u *sentinelUsecase) AssignProjectTeam(projectID uuid.UUID, teamID *uint, requesterRole string) (*domain.Project, error) {
	if requesterRole != domain.RoleCEO && requesterRole != domain.RoleManager {
		return nil, fmt.Errorf("unauthorized: only CEO or MANAGER can assign projects to teams")
	}
	if err := u.repo.AssignProjectTeam(projectID, teamID); err != nil {
		return nil, fmt.Errorf("failed to assign team: %w", err)
	}
	ceoCtx := domain.CallerContext{Role: domain.RoleCEO}
	return u.repo.GetProjectByID(projectID, ceoCtx)
}

// AssignProjectPmOwners sets which PM users own a project when squads are disabled (replaces the whole list).
func (u *sentinelUsecase) AssignProjectPmOwners(projectID uuid.UUID, pmUserIDs []uint, requesterRole string) (*domain.Project, error) {
	if requesterRole != domain.RoleCEO && requesterRole != domain.RoleManager {
		return nil, fmt.Errorf("unauthorized: only CEO or MANAGER can assign project PM owners")
	}
	if !u.isTeamsFeatureDisabled() {
		return nil, &domain.ErrBadRequest{Msg: "project PM owners can only be edited when the teams feature is disabled"}
	}
	ceoCtx := domain.CallerContext{Role: domain.RoleCEO}
	if _, err := u.repo.GetProjectByID(projectID, ceoCtx); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	seen := make(map[uint]struct{})
	var clean []uint
	for _, id := range pmUserIDs {
		if id == 0 {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		uu, err := u.authRepo.FindByID(id)
		if err != nil || uu == nil {
			return nil, &domain.ErrBadRequest{Msg: fmt.Sprintf("user %d not found", id)}
		}
		if uu.Role != authDomain.RolePM {
			return nil, &domain.ErrBadRequest{Msg: fmt.Sprintf("user %d must have role PM (current role: %s)", id, uu.Role)}
		}
		clean = append(clean, id)
	}
	if err := u.repo.ReplaceProjectPmAssignments(projectID, clean); err != nil {
		return nil, fmt.Errorf("failed to save PM assignments: %w", err)
	}
	return u.repo.GetProjectByID(projectID, ceoCtx)
}

// UpdateProject updates project name, description, and status. If updateCode is true, also sets project.Code to slugify(name) and updates all task codes in the project to use the new prefix (so they match the new name).
func (u *sentinelUsecase) UpdateProject(projectID uuid.UUID, name, description, status string, updateCode bool) (*domain.Project, error) {
	ceoCtx := domain.CallerContext{Role: domain.RoleCEO}
	p, err := u.repo.GetProjectByID(projectID, ceoCtx)
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

func (u *sentinelUsecase) CreateTask(title, desc, taskType string, creatorID uint, dueDate *time.Time, projectID, parentID *uuid.UUID, startDate, endDate *time.Time, priority string, storyPoints int, sprintID, milestoneID *uuid.UUID, epicID *uuid.UUID, estimatedMinutes *int) (*domain.Task, error) {
	defaultEstimatedMinutes := 0
	if estimatedMinutes != nil && *estimatedMinutes >= 0 {
		defaultEstimatedMinutes = *estimatedMinutes
	}

	// Validate and normalise task type
	switch domain.TaskType(taskType) {
	case domain.TaskTypeFeature, domain.TaskTypeTask, domain.TaskTypeBug:
		// valid
	case "":
		taskType = string(domain.TaskTypeTask)
	default:
		return nil, fmt.Errorf("invalid task_type: %s (allowed: FEATURE, TASK, BUG)", taskType)
	}

	if priority == "" {
		priority = "MEDIUM"
	}
	if !validPriorities[priority] {
		return nil, fmt.Errorf("invalid priority: %s (allowed: CRITICAL, HIGH, MEDIUM, LOW)", priority)
	}
	if storyPoints < 0 {
		return nil, errors.New("story_points cannot be negative")
	}

	// Sub-tasks (have a parent_id) inherit from parent: dates, and project/epic/sprint when not provided.
	// Max nesting depth is 2 levels (A → B → C). Level C tasks cannot have children.
	if parentID != nil {
		parent, err := u.repo.GetTaskByID(*parentID)
		if err == nil && parent != nil {
			// Reject if the parent itself already has a parent (would create level D)
			if parent.ParentID != nil {
				return nil, &domain.ErrBadRequest{Msg: "cannot create sub-tasks beyond level C — maximum nesting depth is 2 levels"}
			}
			startDate = parent.StartDate
			endDate = parent.EndDate
			if projectID == nil && parent.ProjectID != nil {
				projectID = parent.ProjectID
			}
			if epicID == nil && parent.EpicID != nil {
				epicID = parent.EpicID
			}
			if sprintID == nil && parent.SprintID != nil {
				sprintID = parent.SprintID
			}
		} else {
			startDate = nil
			endDate = nil
		}
	}

	slug := "task"
	if projectID != nil {
		proj, err := u.repo.GetProjectByID(*projectID, domain.CallerContext{Role: domain.RoleCEO})
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
		ID:               uuid.New(),
		Code:             code,
		Title:            title,
		Description:      desc,
		TaskType:         taskType,
		CreatedBy:        &creatorID,
		Status:           "PENDING",
		EstimatedMinutes: defaultEstimatedMinutes,
		DueAt:            dueDate,
		ProjectID:        projectID,
		ParentID:         parentID,
		EpicID:           epicID,
		StartDate:        startDate,
		EndDate:          endDate,
		Priority:         priority,
		StoryPoints:      storyPoints,
		SprintID:         sprintID,
		MilestoneID:      milestoneID,
	}

	if err := u.repo.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

// AssignTask assigns a developer to a task. assignerID is the PM/CEO who performs the assign (for PM-scoped leaderboard).
// devID = 0 means unassign (set AssignedTo = nil, revert to PENDING).
func (u *sentinelUsecase) AssignTask(taskID uuid.UUID, devID uint, assignerID uint) error {
	// 1. Validate if task exists
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	if devID == 0 {
		// Unassign: clear assignee, revert status to PENDING
		task.AssignedTo = nil
		task.AssignedByID = nil
		task.Status = "PENDING"
		task.StartedAt = nil
		return u.repo.UpdateTask(task)
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

// SubmitWork records a handover: Dev submits a PR/Commit URL and moves the task to REVIEW_PENDING
func (u *sentinelUsecase) SubmitWork(taskID uuid.UUID, devID uint, referenceURL, note string) (*domain.Submission, error) {
	sub := &domain.Submission{
		ID:           uuid.New(),
		TaskID:       taskID,
		DevID:        devID,
		ReferenceURL: referenceURL,
		Note:         note,
	}

	if err := u.repo.CreateSubmission(sub); err != nil {
		return nil, err
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	task.Status = "REVIEW_PENDING"
	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
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
	if task.AssignedTo != nil {
		assignee, err := u.authRepo.FindByID(*task.AssignedTo)
		if err == nil && assignee != nil {
			if assignee.DisplayName != "" {
				task.AssignedToDisplayName = assignee.DisplayName
			}
			task.AssignedToEmail = assignee.Email
		}
	}
	return task, nil
}

// GetMyTasks retrieves tasks assigned to a user that belong to an ACTIVE sprint.
// DEV role must only see their active battlefield — tasks outside active sprints are excluded.
func (u *sentinelUsecase) GetMyTasks(userID uint) ([]domain.Task, error) {
	tasks, err := u.repo.GetActiveSprintTasksByAssignee(userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetMyActiveSprints returns the ACTIVE sprints that contain tasks assigned to the user.
func (u *sentinelUsecase) GetMyActiveSprints(userID uint) ([]domain.Sprint, error) {
	return u.repo.GetActiveSprintsForUser(userID)
}

// GetGlobalActiveTasks returns all tasks assigned to the user across ALL projects
// where the sprint is ACTIVE, enriched with project name and color.
func (u *sentinelUsecase) GetGlobalActiveTasks(userID uint) ([]domain.GlobalActiveTask, error) {
	return u.repo.GetGlobalActiveTasksByUser(userID)
}

// GetTeamActiveTasks returns all ACTIVE-sprint tasks within the caller's team.
// CEO/MANAGER (no team restriction) get all active-sprint tasks across all teams.
func (u *sentinelUsecase) GetTeamActiveTasks(callerTeamID *uint, callerRole string) ([]domain.GlobalActiveTask, error) {
	teamID := uint(0)
	if callerRole != domain.RoleCEO && callerRole != domain.RoleManager && callerTeamID != nil {
		teamID = *callerTeamID
	}
	tasks, err := u.repo.GetTeamActiveTasks(teamID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetActiveFeatures returns all FEATURE-type tasks for the PM/CEO Feature Roadmap Board.
// Each feature carries a roll-up progress (0–100%) computed from child TASK/BUG completion.
// CEO/MANAGER see all teams; PM is scoped to their own team.
func (u *sentinelUsecase) GetActiveFeatures(callerTeamID *uint, callerRole string) ([]domain.FeatureRoadmapItem, error) {
	teamID := uint(0)
	if callerRole != domain.RoleCEO && callerRole != domain.RoleManager && callerTeamID != nil {
		teamID = *callerTeamID
	}
	return u.repo.GetActiveFeatures(teamID)
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

// GetPendingApprovals returns tasks requiring PM/CEO/MANAGER attention
// Includes: REVIEW_PENDING handovers, time negotiations (PENDING), appeals (PENDING)
func (u *sentinelUsecase) GetPendingApprovals(userRole string) ([]domain.Task, error) {
	if userRole != "CEO" && userRole != "PM" && userRole != "MANAGER" {
		return nil, fmt.Errorf("access denied: only CEO, PM, or MANAGER can view approvals inbox")
	}

	tasks, err := u.repo.GetTasksRequiringApproval()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending approvals: %w", err)
	}

	return tasks, nil
}

// SubmitAppeal allows a developer to appeal a rejected task
func (u *sentinelUsecase) SubmitAppeal(submissionID uuid.UUID, devID uint, reason string) (*domain.Appeal, error) {
	submission, err := u.repo.GetSubmissionByID(submissionID)
	if err != nil {
		return nil, fmt.Errorf("submission not found: %w", err)
	}

	if submission.DevID != devID {
		return nil, errors.New("unauthorized: only the developer who submitted can appeal")
	}

	existingAppeal, err := u.repo.GetAppealBySubmissionID(submissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing appeal: %w", err)
	}
	if existingAppeal != nil {
		return nil, errors.New("appeal already exists for this submission")
	}

	appeal := &domain.Appeal{
		ID:           uuid.New(),
		SubmissionID: submissionID,
		DeveloperID:  devID,
		Reason:       reason,
		Status:       "PENDING",
	}

	if err := u.repo.CreateAppeal(appeal); err != nil {
		return nil, fmt.Errorf("failed to create appeal: %w", err)
	}

	return appeal, nil
}

// ResolveAppeal allows PM/CEO to approve or reject an appeal
func (u *sentinelUsecase) ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error {
	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("status must be APPROVED or REJECTED")
	}

	resolver, err := u.authRepo.FindByID(resolverID)
	if err != nil {
		return fmt.Errorf("unauthorized: resolver user not found: %w", err)
	}

	if resolver.Role != authDomain.RoleCEO && resolver.Role != authDomain.RoleManager && resolver.Role != authDomain.RolePM {
		return fmt.Errorf("forbidden: only CEO, MANAGER, or PM can resolve appeals (current role: %s)", resolver.Role)
	}

	appeal, err := u.repo.GetAppealByID(appealID)
	if err != nil {
		return fmt.Errorf("appeal not found: %w", err)
	}

	if appeal.Status != "PENDING" {
		return fmt.Errorf("appeal already resolved with status: %s", appeal.Status)
	}

	appeal.Status = status
	appeal.ResolverID = &resolverID
	appeal.ResolverNote = note

	if err := u.repo.UpdateAppeal(appeal); err != nil {
		return fmt.Errorf("failed to update appeal: %w", err)
	}

	// If APPROVED, auto-complete the task
	if status == "APPROVED" {
		submission, err := u.repo.GetSubmissionByID(appeal.SubmissionID)
		if err != nil {
			return fmt.Errorf("failed to get submission: %w", err)
		}

		task, err := u.repo.GetTaskByID(submission.TaskID)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to get task for auto-completion: %v\n", err)
		} else if task.Status != "COMPLETED" {
			task.Status = "COMPLETED"
			now := time.Now()
			task.CompletedAt = &now
			if task.StartedAt == nil {
				task.StartedAt = &now
			}
			if err := u.repo.UpdateTask(task); err != nil {
				fmt.Printf("⚠️  Warning: Failed to auto-complete task: %v\n", err)
			}
		}
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
	if task.EstimatedMinutes > 0 && minutes <= task.EstimatedMinutes {
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
func (u *sentinelUsecase) UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description, taskType string, parentID *uuid.UUID, startDate, endDate *time.Time, progress *int, priority string, storyPoints *int, sprintID *uuid.UUID, applySprint bool, milestoneID *uuid.UUID, epicID *uuid.UUID, applyEpic bool, sortOrder *int, estimatedMinutes *int) (*domain.Task, error) {
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
	if taskType != "" {
		switch domain.TaskType(taskType) {
		case domain.TaskTypeFeature, domain.TaskTypeTask, domain.TaskTypeBug:
			task.TaskType = taskType
		default:
			return nil, fmt.Errorf("invalid task_type: %s (allowed: FEATURE, TASK, BUG)", taskType)
		}
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
	if applySprint {
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
	if estimatedMinutes != nil && *estimatedMinutes >= 0 {
		task.EstimatedMinutes = *estimatedMinutes
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

// EstimateTask uses AI to estimate task effort (title + description) and updates task.estimated_minutes.
// Used internally by ScheduleProjectWithAI. Only task creator, CEO, or PM can run estimate.
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
	task.EstimatedMinutes = minutes
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
	proj, err := u.repo.GetProjectByID(projectID, domain.CallerContext{Role: domain.RoleCEO})
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
		_, err := u.CreateTask(t.Title, t.Description, string(domain.TaskTypeTask), requestingUserID, dueDate, &projectID, nil, startDate, endDate, priority, t.StoryPoints, sprintID, milestoneID, epicID, nil)
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
	if _, err := u.repo.GetProjectByID(projectID, domain.CallerContext{Role: domain.RoleCEO}); err != nil {
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
	if _, err := u.repo.GetProjectByID(projectID, domain.CallerContext{Role: domain.RoleCEO}); err != nil {
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
		task.EstimatedMinutes = minutes
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

	// 3️⃣ Cannot delete if task has sub-tasks (FK constraint)
	childCount, err := u.repo.CountChildTasks(taskID)
	if err != nil {
		return fmt.Errorf("failed to check sub-tasks: %w", err)
	}
	if childCount > 0 {
		return fmt.Errorf("task_has_sub_tasks: cannot delete task with sub-tasks, delete sub-tasks first")
	}

	// 4️⃣ Delete from database
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

	// 4️⃣ Roll-up: if this is a child TASK/BUG and has a parent FEATURE,
	//    check if all siblings are now COMPLETED → auto-promote parent to READY_FOR_UAT
	if (task.TaskType == "TASK" || task.TaskType == "BUG") && task.ParentID != nil {
		siblings, err := u.repo.GetChildTasksByParentID(*task.ParentID)
		if err == nil && len(siblings) > 0 {
			allDone := true
			for _, s := range siblings {
				effectiveStatus := s.Status
				if s.ID == taskID {
					effectiveStatus = "COMPLETED" // just approved above
				}
				if effectiveStatus != "COMPLETED" {
					allDone = false
					break
				}
			}
			if allDone {
				parent, err := u.repo.GetTaskByID(*task.ParentID)
				if err == nil && parent != nil && parent.TaskType == "FEATURE" &&
					parent.Status != "COMPLETED" && parent.Status != "READY_FOR_UAT" && parent.Status != "REVIEW_PENDING" {
					parent.Status = "READY_FOR_UAT"
					_ = u.repo.UpdateTask(parent)
				}
			}
		}
	}

	return nil
}

// SubmitUAT stores the UAT payload on a FEATURE task and moves it to REVIEW_PENDING for PM/CEO to review.
func (u *sentinelUsecase) SubmitUAT(taskID uuid.UUID, devID uint, payload domain.UATPayloadData) error {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}

	if task.TaskType != "FEATURE" {
		return &domain.ErrBadRequest{Msg: "submit-uat is only allowed for FEATURE tasks"}
	}
	if task.Status != "READY_FOR_UAT" {
		return &domain.ErrBadRequest{Msg: fmt.Sprintf("task must be in READY_FOR_UAT status to submit UAT (current: %s)", task.Status)}
	}

	if payload.StagingURL == "" || (!strings.HasPrefix(payload.StagingURL, "http://") && !strings.HasPrefix(payload.StagingURL, "https://")) {
		return &domain.ErrBadRequest{Msg: "staging_url must be a valid http or https URL"}
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode UAT payload: %w", err)
	}
	task.UATPayload = datatypes.JSON(raw)
	task.Status = "REVIEW_PENDING"

	if err := u.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to save UAT payload: %w", err)
	}
	return nil
}

// RejectTask returns a task to IN_PROGRESS and logs rejection reason as a comment (PM/CEO/MANAGER only)
func (u *sentinelUsecase) RejectTask(taskID uuid.UUID, rejectorID uint, rejectorRole string, reason string) error {
	if rejectorRole != "CEO" && rejectorRole != "PM" && rejectorRole != "MANAGER" {
		return fmt.Errorf("access denied: only PM, CEO or MANAGER can reject tasks (your role: %s)", rejectorRole)
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != "REVIEW_PENDING" {
		return fmt.Errorf("task is not pending review (current status: %s)", task.Status)
	}

	return u.repo.RejectTask(taskID, rejectorID, reason)
}

// MarkReadyForTest moves a TASK or BUG from IN_PROGRESS to READY_FOR_TEST (Dev action).
func (u *sentinelUsecase) MarkReadyForTest(taskID uuid.UUID, devID uint) error {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.TaskType != "TASK" && task.TaskType != "BUG" {
		return &domain.ErrBadRequest{Msg: "only TASK or BUG type tasks can be moved to READY_FOR_TEST"}
	}
	if task.Status != "IN_PROGRESS" {
		return &domain.ErrBadRequest{Msg: fmt.Sprintf("task must be IN_PROGRESS to mark ready for test (current: %s)", task.Status)}
	}
	task.Status = "READY_FOR_TEST"
	if err := u.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

// PMApproveSubTask is the PM's first-stage approval: READY_FOR_TEST → READY_FOR_UAT.
// The PM must provide a test URL and detailed test steps that the CEO will follow to do the final UAT.
func (u *sentinelUsecase) PMApproveSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, testURL string, testSteps string) error {
	if pmRole != "PM" && pmRole != "MANAGER" {
		return fmt.Errorf("access denied: only PM or MANAGER can submit for CEO approval (your role: %s)", pmRole)
	}
	if strings.TrimSpace(testURL) == "" {
		return &domain.ErrBadRequest{Msg: "test_url is required"}
	}
	if len(strings.TrimSpace(testSteps)) < 20 {
		return &domain.ErrBadRequest{Msg: "test_steps must be at least 20 characters"}
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != "READY_FOR_TEST" {
		return fmt.Errorf("task is not in READY_FOR_TEST status (current: %s)", task.Status)
	}

	payload := fmt.Sprintf(`{"test_url":%q,"test_steps":%q}`, strings.TrimSpace(testURL), strings.TrimSpace(testSteps))
	if err := u.repo.SetTaskReadyForUAT(taskID, []byte(payload)); err != nil {
		return fmt.Errorf("failed to submit for CEO approval: %w", err)
	}
	return nil
}

// ApproveSubTask is the CEO's final approval: READY_FOR_UAT → COMPLETED.
func (u *sentinelUsecase) ApproveSubTask(taskID uuid.UUID, ceoUserID uint, ceoRole string) error {
	if ceoRole != "CEO" && ceoRole != "MANAGER" {
		return fmt.Errorf("access denied: only CEO or MANAGER can give final approval (your role: %s)", ceoRole)
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != "READY_FOR_UAT" {
		return fmt.Errorf("task is not in READY_FOR_UAT status (current: %s) — PM must submit test evidence first", task.Status)
	}

	if err := u.repo.ApproveTask(taskID); err != nil {
		return fmt.Errorf("failed to approve sub-task: %w", err)
	}

	// Roll-up: if all siblings are COMPLETED, promote parent FEATURE to READY_FOR_UAT
	if (task.TaskType == "TASK" || task.TaskType == "BUG") && task.ParentID != nil {
		siblings, err := u.repo.GetChildTasksByParentID(*task.ParentID)
		if err == nil && len(siblings) > 0 {
			allDone := true
			for _, s := range siblings {
				effectiveStatus := s.Status
				if s.ID == taskID {
					effectiveStatus = "COMPLETED"
				}
				if effectiveStatus != "COMPLETED" {
					allDone = false
					break
				}
			}
			if allDone {
				parent, err := u.repo.GetTaskByID(*task.ParentID)
				if err == nil && parent != nil && parent.TaskType == "FEATURE" &&
					parent.Status != "COMPLETED" && parent.Status != "READY_FOR_UAT" && parent.Status != "REVIEW_PENDING" {
					parent.Status = "READY_FOR_UAT"
					_ = u.repo.UpdateTask(parent)
				}
			}
		}
	}

	return nil
}

// RejectSubTask returns a task to IN_PROGRESS with a reason log.
// Accepts READY_FOR_TEST (PM rejecting dev work) or READY_FOR_UAT (CEO rejecting PM evidence).
func (u *sentinelUsecase) RejectSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, reason string) error {
	if pmRole != "CEO" && pmRole != "PM" && pmRole != "MANAGER" {
		return fmt.Errorf("access denied: only PM, CEO or MANAGER can reject sub-tasks (your role: %s)", pmRole)
	}
	if len(reason) < 10 {
		return &domain.ErrBadRequest{Msg: "rejection reason must be at least 10 characters"}
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != "READY_FOR_TEST" && task.Status != "READY_FOR_UAT" {
		return fmt.Errorf("task must be in READY_FOR_TEST or READY_FOR_UAT status to be rejected (current: %s)", task.Status)
	}

	return u.repo.RejectTask(taskID, pmUserID, reason)
}

// GetTasksReadyForTest returns all TASK/BUG items in READY_FOR_TEST status, scoped to the caller's team (PM/CEO).
func (u *sentinelUsecase) GetTasksReadyForTest(callerTeamID *uint, callerRole string) ([]domain.GlobalActiveTask, error) {
	if callerRole != "CEO" && callerRole != "PM" && callerRole != "MANAGER" {
		return nil, fmt.Errorf("access denied: only PM, CEO or MANAGER can view the test queue (your role: %s)", callerRole)
	}
	var teamID uint
	if callerTeamID != nil {
		teamID = *callerTeamID
	}
	return u.repo.GetTasksReadyForTest(teamID)
}

// GetTasksReadyForCEOApproval returns TASK/BUG items in READY_FOR_UAT status awaiting CEO final approval.
func (u *sentinelUsecase) GetTasksReadyForCEOApproval(callerTeamID *uint, callerRole string) ([]domain.GlobalActiveTask, error) {
	if callerRole != "CEO" && callerRole != "MANAGER" {
		return nil, fmt.Errorf("access denied: only CEO or MANAGER can view the CEO approval queue (your role: %s)", callerRole)
	}
	var teamID uint
	if callerTeamID != nil {
		teamID = *callerTeamID
	}
	return u.repo.GetTasksReadyForCEOApproval(teamID)
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
	childCount, err := u.repo.CountChildTasks(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify task hierarchy: %w", err)
	}
	if childCount > 0 {
		return nil, &domain.ErrBadRequest{Msg: "time cannot be logged against a Parent Task; log time on its Sub-tasks instead"}
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
		"PENDING": true, "IN_PROGRESS": true, "READY_FOR_TEST": true, "REVIEW_PENDING": true, "COMPLETED": true, "BLOCKED": true,
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
	spreadsheetIDRegex  = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	spreadsheetGIDRegex = regexp.MustCompile(`[#?&]gid=(\d+)`)
	presentationIDRegex = regexp.MustCompile(`/presentation/d/([a-zA-Z0-9_-]+)`)
	pptxSlideNumRegex   = regexp.MustCompile(`slide(\d+)\.xml$`)
	pptxTitlePhRe       = regexp.MustCompile(`<p:ph[^>]*type="(?:title|ctrTitle)"`)
	pptxSystemPhRe      = regexp.MustCompile(`<p:ph[^>]*type="(?:dt|ftr|sldNum|hdr)"`)
	pptxNotesBodyPhRe   = regexp.MustCompile(`<p:ph(?:[^>]*idx="1"|[^>]*type="body")`)
	pptxATextRe         = regexp.MustCompile(`<a:t(?:\s[^>]*)?>([^<]*)</a:t>`)
	pptxRelationshipTagRe = regexp.MustCompile(`(?i)<Relationship\s+([^>]+)/\s*>`)
	pptxRelTypeDq         = regexp.MustCompile(`(?i)\bType\s*=\s*"([^"]*)"`)
	pptxRelTypeSq         = regexp.MustCompile(`(?i)\bType\s*=\s*'([^']*)'`)
	pptxRelTargetDq       = regexp.MustCompile(`(?i)\bTarget\s*=\s*"([^"]*)"`)
	pptxRelTargetSq       = regexp.MustCompile(`(?i)\bTarget\s*=\s*'([^']*)'`)
	pptxDocTitleRe    = regexp.MustCompile(`<dc:title>([^<]*)</dc:title>`)
	pptxSldIdRe       = regexp.MustCompile(`<p:sldId[^>]*id="(\d+)"`)
	pptxSldShowAttrRe = regexp.MustCompile(`<p:sld\s*([^>]+)>`) // opening tag of slide: check for show="0" or show="false" (hidden)
)

const (
	pptxMaxImageBytes     = 12 * 1024 * 1024 // cap per embedded file (design exports can be large)
	pptxMinImageBytes     = 256               // skip tiny icons/bullets
	pptxMaxImagesPerSlide = 30                // เก็บได้หลายรูปต่อหน้า (รูปฝัง + รูป export)
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
	case ".svg":
		return "image/svg+xml"
	default:
		return "" // EMF, WMF, etc. — skip (no lightweight decoder)
	}
}

func normalizeZipEntryName(name string) string {
	return path.Clean(strings.ReplaceAll(strings.TrimSpace(name), "\\", "/"))
}

// pptxRelAttr reads Type or Target from a Relationship attribute fragment; order-independent.
func pptxRelAttr(attrs, key string) string {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "type":
		if m := pptxRelTypeDq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
		if m := pptxRelTypeSq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
	case "target":
		if m := pptxRelTargetDq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
		if m := pptxRelTargetSq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}

func isPPTXImageRelationshipType(typ string) bool {
	t := strings.ToLower(strings.TrimSpace(typ))
	return strings.Contains(t, "relationships/image") || strings.HasSuffix(t, "/image")
}

// pptxFirstRelationshipTargetByTypeContains returns Target for the first Relationship whose Type contains hint (e.g. notesSlide, slideLayout).
func pptxFirstRelationshipTargetByTypeContains(relsData []byte, typeHint string) string {
	h := strings.ToLower(typeHint)
	for _, m := range pptxRelationshipTagRe.FindAllStringSubmatch(string(relsData), -1) {
		if len(m) < 2 {
			continue
		}
		typ := pptxRelAttr(m[1], "Type")
		if typ == "" || !strings.Contains(strings.ToLower(typ), h) {
			continue
		}
		if tgt := pptxRelAttr(m[1], "Target"); tgt != "" {
			return tgt
		}
	}
	return ""
}

// extractImagesFromRels reads a _rels file and extracts all image targets from relDir.
// relDir is the directory containing the .rels file (e.g. ppt/slides for slide1.xml.rels).
func extractImagesFromRels(zr *zip.Reader, relDir string, relsData []byte) []pptxImageEntry {
	seen := make(map[string]bool)
	var entries []pptxImageEntry

	for _, m := range pptxRelationshipTagRe.FindAllStringSubmatch(string(relsData), -1) {
		if len(m) < 2 {
			continue
		}
		typ := pptxRelAttr(m[1], "Type")
		target := pptxRelAttr(m[1], "Target")
		if target == "" || !isPPTXImageRelationshipType(typ) {
			continue
		}
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
		if layoutTarget := pptxFirstRelationshipTargetByTypeContains(relsData, "slideLayout"); layoutTarget != "" {
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

// readZipEntry reads a zip member by logical path. Matching is case-insensitive so PPTX from
// macOS/Canva (mixed-case paths like PPT/Media/...) still resolves on Linux/Docker.
func readZipEntry(r *zip.Reader, name string) []byte {
	want := normalizeZipEntryName(name)
	var match *zip.File
	for _, f := range r.File {
		if strings.EqualFold(normalizeZipEntryName(f.Name), want) {
			match = f
			break
		}
	}
	if match == nil {
		return nil
	}
	rc, err := match.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}
	return data
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
	notesTarget := pptxFirstRelationshipTargetByTypeContains(relsData, "notesSlide")
	if notesTarget == "" {
		return ""
	}
	// Target is relative to the slide directory, e.g. "../notesSlides/notesSlide1.xml"
	return path.Clean(dir + "/" + notesTarget)
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
			slides = append(slides, domain.PreviewSlideItem{
				Index:              s.Index,
				Title:              t,
				SuggestedTaskTitle: suggestedTaskTitleFromSlideText(s.Body, s.Index),
				Hidden:             s.Hidden,
			})
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
					body := extractAPISlideBody(slide)
					slides = append(slides, domain.PreviewSlideItem{
						Index:              i + 1,
						Title:              t,
						SuggestedTaskTitle: suggestedTaskTitleFromSlideText(body, i+1),
						Hidden:             false,
					})
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
	var parentUUID *uuid.UUID
	if req.ParentID != "" {
		parsed, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentUUID = &parsed
	}

	// Build triage map: slide_index -> TriagedSlide for per-slide overrides.
	// Also derive SlideIndices from triage data when provided.
	triagedMap := make(map[int]domain.TriagedSlide)
	if len(req.Slides) > 0 {
		for _, ts := range req.Slides {
			triagedMap[ts.SlideIndex] = ts
		}
		// Derive SlideIndices from triage data so filtering works normally downstream.
		if len(req.SlideIndices) == 0 {
			for idx := range triagedMap {
				req.SlideIndices = append(req.SlideIndices, idx)
			}
		}
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
	proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO})
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
			// Block-level <img> — not inside <p> (TipTap/ProseMirror block image cannot nest in paragraph).
			htmlParts = append(htmlParts, "<img src=\""+html.EscapeString(imgSrc)+"\" class=\"editor-image\" alt=\"Slide\" />")
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

		// Apply per-slide triage overrides when provided; default task title from slide body text, not placeholder title.
		taskTitle := suggestedTaskTitleFromSlideText(slide.Body, slide.Index)
		taskPriority := priority
		taskEstimatedMinutes := 0
		var taskAssigneeID *uint
		if triage, ok := triagedMap[slide.Index]; ok {
			if strings.TrimSpace(triage.Title) != "" {
				taskTitle = strings.TrimSpace(triage.Title)
			}
			p := strings.ToUpper(strings.TrimSpace(triage.Priority))
			if map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}[p] {
				taskPriority = p
			}
			if triage.EstimatedMinutes > 0 {
				taskEstimatedMinutes = triage.EstimatedMinutes
			}
			taskAssigneeID = triage.AssigneeID
		}

		task := &domain.Task{
			ID:               uuid.New(),
			Code:             code,
			Title:            taskTitle,
			Description:      description,
			TaskType:         string(domain.TaskTypeTask),
			CreatedBy:        &creatorID,
			Status:           "PENDING",
			Priority:         taskPriority,
			StoryPoints:      storyPoints,
			EstimatedMinutes: taskEstimatedMinutes,
			SprintID:         sprintUUID,
			EpicID:           epicUUID,
			ProjectID:        &projectIDCopy,
			ParentID:         parentUUID,
			ResourceURLs:     datatypes.JSON(resourceURLsJSON),
		}
		if taskAssigneeID != nil {
			task.AssignedTo = taskAssigneeID
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

var thaiMonthAbbrevToMonth = map[string]time.Month{
	"ม.ค.": time.January, "ก.พ.": time.February, "มี.ค.": time.March, "เม.ย.": time.April,
	"พ.ค.": time.May, "มิ.ย.": time.June, "ก.ค.": time.July, "ส.ค.": time.August,
	"ก.ย.": time.September, "ต.ค.": time.October, "พ.ย.": time.November, "ธ.ค.": time.December,
}

var validSheetImportStatuses = map[string]bool{
	"PENDING": true, "IN_PROGRESS": true, "READY_FOR_TEST": true, "READY_FOR_UAT": true, "COMPLETED": true, "CANCELLED": true,
}

func sheetCSVCell(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(row[col], "\u200b", ""))
}

func sheetTitleFromContentDisposition(cd string) string {
	cd = strings.TrimSpace(cd)
	if cd == "" {
		return ""
	}
	if m := regexp.MustCompile(`filename\*=UTF-8''([^;\s]+)`).FindStringSubmatch(cd); len(m) > 1 {
		if dec, err := url.PathUnescape(strings.Trim(m[1], `"`)); err == nil {
			return strings.TrimSuffix(strings.TrimSpace(dec), ".csv")
		}
		return strings.TrimSuffix(strings.TrimSpace(m[1]), ".csv")
	}
	if m := regexp.MustCompile(`filename="([^"]+)"`).FindStringSubmatch(cd); len(m) > 1 {
		return strings.TrimSuffix(m[1], ".csv")
	}
	return ""
}

func parseGoogleSheetURL(raw string) (sheetID, gid string, err error) {
	u := strings.TrimSpace(raw)
	if u == "" {
		return "", "", errors.New("empty sheet URL")
	}
	m := spreadsheetIDRegex.FindStringSubmatch(u)
	if len(m) < 2 {
		return "", "", errors.New("invalid Google Sheets URL: missing spreadsheet id")
	}
	sheetID = m[1]
	gid = "0"
	if gm := spreadsheetGIDRegex.FindStringSubmatch(u); len(gm) > 1 {
		gid = gm[1]
	}
	return sheetID, gid, nil
}

func fetchGoogleSheetCSVRecords(sheetID, gid string) (records [][]string, sheetTitle string, err error) {
	exportURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", sheetID, gid)
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodGet, exportURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Sentinel/1.0)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download sheet CSV: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Google Sheets returned HTTP %d: ensure the spreadsheet is shared as \"Anyone with the link can view\"", resp.StatusCode)
	}
	sheetTitle = sheetTitleFromContentDisposition(resp.Header.Get("Content-Disposition"))
	r := csv.NewReader(bytes.NewReader(body))
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	records, err = r.ReadAll()
	if err != nil {
		return nil, sheetTitle, fmt.Errorf("invalid CSV: %w", err)
	}
	return records, sheetTitle, nil
}

func parseThaiBuddhistShortDate(s string) (time.Time, bool) {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) < 3 {
		return time.Time{}, false
	}
	day, err1 := strconv.Atoi(parts[0])
	month, ok := thaiMonthAbbrevToMonth[parts[1]]
	yearBE, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || !ok {
		return time.Time{}, false
	}
	if day < 1 || day > 31 {
		return time.Time{}, false
	}
	yearCE := yearBE - 543
	if yearCE < 1900 || yearCE > 2100 {
		return time.Time{}, false
	}
	return time.Date(yearCE, month, day, 0, 0, 0, 0, time.UTC), true
}

func parseSlashSheetDate(s string) (time.Time, bool) {
	parts := strings.Split(strings.TrimSpace(s), "/")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	d, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil || d < 1 || m < 1 || m > 12 {
		return time.Time{}, false
	}
	if y < 100 {
		y += 2000
	}
	if y < 1900 || y > 2100 {
		return time.Time{}, false
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), true
}

func parseSheetDueRaw(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\u200b", ""))
	if s == "" {
		return ""
	}
	if t, ok := parseThaiBuddhistShortDate(s); ok {
		return t.Format("2006-01-02")
	}
	if t, ok := parseSlashSheetDate(s); ok {
		return t.Format("2006-01-02")
	}
	return ""
}

func mapKGSheetStatus(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "\u200b", ""))
	if raw == "" {
		return "PENDING"
	}
	if idx := strings.IndexAny(raw, "\n\r"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	switch raw {
	case "แก้ไขแล้ว", "นำขึ้น Prod แล้ว":
		return "COMPLETED"
	case "ยังไม่แก้":
		return "PENDING"
	case "กำลังแก้ไข", "แก้แล้วแต่ไม่ถูกต้อง", "แก้ไขอีกครั้ง":
		return "IN_PROGRESS"
	case "ทดสอบอีกครั้ง":
		return "READY_FOR_TEST"
	case "รอนำขึ้น Prod":
		return "READY_FOR_UAT"
	default:
		return "PENDING"
	}
}

func (u *sentinelUsecase) PreviewGoogleSheets(req *domain.PreviewGoogleSheetsRequest) (*domain.PreviewGoogleSheetsResult, error) {
	if req == nil || strings.TrimSpace(req.SheetURL) == "" {
		return nil, errors.New("sheet_url is required")
	}
	sheetID, gid, err := parseGoogleSheetURL(req.SheetURL)
	if err != nil {
		return nil, err
	}
	records, sheetTitle, err := fetchGoogleSheetCSVRecords(sheetID, gid)
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, errors.New("sheet has no data rows (only header or empty)")
	}
	if sheetTitle == "" {
		sheetTitle = "Google Sheet"
	}
	var rows []domain.SheetRowPreviewItem
	for i := 1; i < len(records); i++ {
		row := records[i]
		title := sheetCSVCell(row, 1)
		if title == "" {
			continue
		}
		dueRaw := sheetCSVCell(row, 0)
		statusRaw := sheetCSVCell(row, 5)
		notes := sheetCSVCell(row, 10)
		dueStr := parseSheetDueRaw(dueRaw)
		rawFirst := statusRaw
		if idx := strings.IndexAny(rawFirst, "\n\r"); idx >= 0 {
			rawFirst = strings.TrimSpace(rawFirst[:idx])
		}
		rows = append(rows, domain.SheetRowPreviewItem{
			RowIndex:  i + 1,
			Title:     title,
			DueDate:   dueStr,
			Status:    mapKGSheetStatus(statusRaw),
			RawStatus: rawFirst,
			Notes:     notes,
		})
	}
	if len(rows) == 0 {
		return nil, errors.New("no importable rows: column B (รายละเอียด) is empty for all data rows")
	}
	return &domain.PreviewGoogleSheetsResult{
		SheetTitle: sheetTitle,
		SheetID:    sheetID,
		Rows:       rows,
	}, nil
}

func (u *sentinelUsecase) ImportFromGoogleSheets(req *domain.ImportGoogleSheetsRequest, creatorID uint) (*domain.ImportGoogleSheetsResult, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if len(req.Rows) == 0 {
		return nil, errors.New("at least one row is required to import")
	}
	_, _, err := parseGoogleSheetURL(req.SheetURL)
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
	var parentUUID *uuid.UUID
	if req.ParentID != "" {
		parsed, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentUUID = &parsed
	}

	if parentUUID != nil {
		parent, err := u.repo.GetTaskByID(*parentUUID)
		if err != nil || parent == nil {
			return nil, errors.New("parent task not found")
		}
		if parent.ParentID != nil {
			return nil, &domain.ErrBadRequest{Msg: "cannot attach sheet import under a nested sub-task"}
		}
		if parent.ProjectID == nil || *parent.ProjectID != projectUUID {
			return nil, &domain.ErrBadRequest{Msg: "parent task must belong to the same project"}
		}
	}

	slug := "task"
	if proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO}); err == nil && proj != nil {
		slug = slugify(proj.Name)
	}
	maxSuffix, err := u.repo.GetMaxTaskCodeSuffix(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get next task code: %w", err)
	}

	var created []*domain.Task
	for i, tr := range req.Rows {
		title := strings.TrimSpace(tr.Title)
		if title == "" {
			return nil, fmt.Errorf("row %d: title is required", tr.RowIndex)
		}
		estMins := tr.EstimatedMinutes
		if estMins < 0 {
			return nil, fmt.Errorf("row %d: estimated_minutes cannot be negative", tr.RowIndex)
		}
		priority := strings.ToUpper(strings.TrimSpace(tr.Priority))
		if priority == "" {
			priority = "MEDIUM"
		}
		if !validPriorities[priority] {
			return nil, fmt.Errorf("row %d: invalid priority %q", tr.RowIndex, tr.Priority)
		}
		st := strings.ToUpper(strings.TrimSpace(tr.Status))
		if st == "" {
			st = "PENDING"
		}
		if !validSheetImportStatuses[st] {
			return nil, fmt.Errorf("row %d: invalid status %q", tr.RowIndex, tr.Status)
		}

		desc := strings.TrimSpace(tr.Notes)
		if desc != "" {
			desc = "<p>" + html.EscapeString(desc) + "</p>"
		}

		var duePtr *time.Time
		if ds := strings.TrimSpace(tr.DueDate); ds != "" {
			t, err := time.Parse("2006-01-02", ds)
			if err != nil {
				return nil, fmt.Errorf("row %d: invalid due_date %q (use YYYY-MM-DD)", tr.RowIndex, tr.DueDate)
			}
			duePtr = &t
		}

		projectIDCopy := projectUUID
		task := &domain.Task{
			ID:               uuid.New(),
			Code:             fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i),
			Title:            title,
			Description:      desc,
			TaskType:         string(domain.TaskTypeBug),
			CreatedBy:        &creatorID,
			Status:           st,
			Priority:         priority,
			StoryPoints:      0,
			EstimatedMinutes: estMins,
			DueAt:            duePtr,
			SprintID:         sprintUUID,
			EpicID:           epicUUID,
			ProjectID:        &projectIDCopy,
			ParentID:         parentUUID,
		}

		if err := u.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("failed to create task for sheet row %d: %w", tr.RowIndex, err)
		}
		created = append(created, task)
	}

	titleOut := strings.TrimSpace(req.SheetTitle)
	if titleOut == "" {
		titleOut = "Google Sheet"
	}
	return &domain.ImportGoogleSheetsResult{
		CreatedCount: len(created),
		SheetTitle:   titleOut,
		Tasks:        created,
	}, nil
}

// SplitTask decomposes one task into N new sub-tasks (inheriting same parent_id, project, epic, sprint)
// then deletes the original task. The caller must be CEO, PM, or the task creator.
func (u *sentinelUsecase) SplitTask(taskID uuid.UUID, splits []domain.SplitTaskItem, requestingUserID uint, requestingUserRole string) ([]*domain.Task, error) {
	if len(splits) < 2 {
		return nil, &domain.ErrBadRequest{Msg: "split requires at least 2 items"}
	}
	for i, s := range splits {
		if strings.TrimSpace(s.Title) == "" {
			return nil, &domain.ErrBadRequest{Msg: fmt.Sprintf("split item %d has empty title", i+1)}
		}
	}

	// Load original task
	orig, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Access control: CEO / PM / creator
	role := strings.ToUpper(strings.TrimSpace(requestingUserRole))
	if role != "CEO" && role != "PM" && role != "MANAGER" {
		if orig.CreatedBy == nil || *orig.CreatedBy != requestingUserID {
			return nil, &domain.ErrBadRequest{Msg: "only CEO, PM, or the task creator can split a task"}
		}
	}

	// Determine slug for code generation
	slug := "task"
	if orig.ProjectID != nil {
		if proj, err := u.repo.GetProjectByID(*orig.ProjectID, domain.CallerContext{Role: domain.RoleCEO}); err == nil && proj != nil {
			slug = slugify(proj.Name)
		}
	}
	maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(slug)

	var created []*domain.Task
	for i, item := range splits {
		priority := strings.ToUpper(strings.TrimSpace(item.Priority))
		if !map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}[priority] {
			priority = orig.Priority
			if priority == "" {
				priority = "MEDIUM"
			}
		}
		code := fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i)
		t := &domain.Task{
			ID:               uuid.New(),
			Code:             code,
			Title:            strings.TrimSpace(item.Title),
			Description:      orig.Description,
			TaskType:         orig.TaskType,
			Status:           "PENDING",
			Priority:         priority,
			EstimatedMinutes: item.EstimatedMinutes,
			ProjectID:        orig.ProjectID,
			ParentID:         orig.ParentID,   // same parent as the original
			EpicID:           orig.EpicID,
			SprintID:         orig.SprintID,
			ResourceURLs:     orig.ResourceURLs,
			CreatedBy:        &requestingUserID,
		}
		if item.AssigneeID != nil {
			t.AssignedTo = item.AssigneeID
		}
		if err := u.repo.CreateTask(t); err != nil {
			return nil, fmt.Errorf("failed to create split task %d: %w", i+1, err)
		}
		created = append(created, t)
	}

	// Delete the original task
	if err := u.repo.DeleteTask(taskID); err != nil {
		return nil, fmt.Errorf("split tasks created but failed to delete original: %w", err)
	}

	return created, nil
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
	project, err := u.repo.GetProjectByID(projectID, domain.CallerContext{Role: domain.RoleCEO})
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

// --- Internal B2B Outsource Usecase Methods ---

func (u *sentinelUsecase) CreateB2BRequest(title, description string, estimatedMinutes int, requesterTeamID, targetTeamID, requesterUserID uint) (*domain.B2BRequest, error) {
	if title == "" {
		return nil, &domain.ErrBadRequest{Msg: "title is required"}
	}
	if estimatedMinutes <= 0 {
		return nil, &domain.ErrBadRequest{Msg: "estimated_minutes must be greater than 0"}
	}
	if requesterTeamID == targetTeamID {
		return nil, &domain.ErrBadRequest{Msg: "cannot send a B2B request to your own team"}
	}
	req := &domain.B2BRequest{
		Title:            title,
		Description:      description,
		EstimatedMinutes: estimatedMinutes,
		Status:           "PENDING",
		RequesterTeamID:  requesterTeamID,
		TargetTeamID:     targetTeamID,
		RequesterUserID:  requesterUserID,
	}
	if err := u.repo.CreateB2BRequest(req); err != nil {
		return nil, fmt.Errorf("create b2b request: %w", err)
	}
	return req, nil
}

func (u *sentinelUsecase) GetB2BRequests(callerTeamID uint, direction string) ([]domain.B2BRequest, error) {
	if direction != "inbound" && direction != "outbound" {
		direction = "inbound"
	}
	reqs, err := u.repo.GetB2BRequests(callerTeamID, direction)
	if err != nil {
		return nil, err
	}

	// Enrich team names via authRepo
	teamCache := map[uint]string{}
	getTeamName := func(id uint) string {
		if name, ok := teamCache[id]; ok {
			return name
		}
		teams, err2 := u.authRepo.GetAllTeams()
		if err2 != nil {
			return ""
		}
		for _, t := range teams {
			teamCache[t.ID] = t.Name
		}
		return teamCache[id]
	}

	for i := range reqs {
		reqs[i].RequesterTeamName = getTeamName(reqs[i].RequesterTeamID)
		reqs[i].TargetTeamName = getTeamName(reqs[i].TargetTeamID)
	}
	return reqs, nil
}

func (u *sentinelUsecase) CounterOfferB2BRequest(id uuid.UUID, callerTeamID uint, proposedMinutes int, reason string) (*domain.B2BRequest, error) {
	req, err := u.repo.GetB2BRequestByID(id)
	if err != nil {
		return nil, fmt.Errorf("b2b request not found: %w", err)
	}
	if req.TargetTeamID != callerTeamID {
		return nil, &domain.ErrBadRequest{Msg: "only the target team can counter-offer"}
	}
	if req.Status != "PENDING" && req.Status != "COUNTER_OFFERED" {
		return nil, &domain.ErrBadRequest{Msg: "cannot counter-offer a request that is already accepted or rejected"}
	}
	if proposedMinutes <= 0 {
		return nil, &domain.ErrBadRequest{Msg: "proposed_minutes must be greater than 0"}
	}
	req.Status = "COUNTER_OFFERED"
	req.ProposedMinutes = proposedMinutes
	req.NegotiationReason = reason
	if err := u.repo.UpdateB2BRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (u *sentinelUsecase) RejectB2BRequest(id uuid.UUID, callerTeamID uint) (*domain.B2BRequest, error) {
	req, err := u.repo.GetB2BRequestByID(id)
	if err != nil {
		return nil, fmt.Errorf("b2b request not found: %w", err)
	}
	// Either team can reject
	if req.RequesterTeamID != callerTeamID && req.TargetTeamID != callerTeamID {
		return nil, &domain.ErrBadRequest{Msg: "not authorized to reject this request"}
	}
	if req.Status == "ACCEPTED" {
		return nil, &domain.ErrBadRequest{Msg: "cannot reject an already accepted request"}
	}
	req.Status = "REJECTED"
	if err := u.repo.UpdateB2BRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

// AcceptB2BRequest is called by the target team's PM.
// It creates a Task in the first project of the target team and marks the request as ACCEPTED.
func (u *sentinelUsecase) AcceptB2BRequest(id uuid.UUID, callerTeamID uint, accepterUserID uint) (*domain.Task, error) {
	req, err := u.repo.GetB2BRequestByID(id)
	if err != nil {
		return nil, fmt.Errorf("b2b request not found: %w", err)
	}
	if req.TargetTeamID != callerTeamID {
		return nil, &domain.ErrBadRequest{Msg: "only the target team can accept this request"}
	}
	if req.Status != "PENDING" && req.Status != "COUNTER_OFFERED" {
		return nil, &domain.ErrBadRequest{Msg: "request is not in a state that can be accepted"}
	}

	// Resolve the minutes to use (if counter-offered, use proposed; else original estimate)
	minutes := req.EstimatedMinutes
	if req.Status == "COUNTER_OFFERED" && req.ProposedMinutes > 0 {
		minutes = req.ProposedMinutes
	}

	// Find the first project of the target team
	project, err := u.repo.GetFirstProjectByTeamID(callerTeamID)
	if err != nil {
		return nil, fmt.Errorf("target team has no projects: %w", err)
	}

	// Create the task in that project
	task, err := u.CreateTask(
		req.Title,
		req.Description,
		"TASK",
		accepterUserID,
		nil,           // no due date
		&project.ID,
		nil,           // no parent
		nil, nil,      // no start/end date
		"MEDIUM",
		0,             // no story points
		nil, nil,      // no sprint/milestone
		nil,           // no epic
		&minutes,
	)
	if err != nil {
		return nil, fmt.Errorf("create task for b2b request: %w", err)
	}

	req.Status = "ACCEPTED"
	req.CreatedTaskID = &task.ID
	if err := u.repo.UpdateB2BRequest(req); err != nil {
		return nil, fmt.Errorf("update b2b request: %w", err)
	}
	return task, nil
}

