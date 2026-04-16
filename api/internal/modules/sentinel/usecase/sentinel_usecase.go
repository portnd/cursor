package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type sentinelUsecase struct {
	repo         domain.SentinelRepository
	aiService    domain.AIService
	authRepo     authDomain.Repository
	usageTracker domain.UsageTracker
	aiLimitRPM   int
	aiLimitRPD   int
	timeout      time.Duration
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

// projectNameAllowedChars matches unicode letters/marks/digits, spaces, hyphens, underscores.
var projectNameAllowedChars = regexp.MustCompile(`^[\p{L}\p{M}\p{N}\s\-_]+$`)

func (u *sentinelUsecase) CreateProject(name, description, status string, ctx domain.CallerContext) (*domain.Project, error) {
	ctx = u.withCallerScope(ctx)
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("project name is required")
	}
	if !projectNameAllowedChars.MatchString(name) {
		return nil, errors.New("project name contains invalid characters (allowed: letters, numbers, spaces, hyphens, underscores)")
	}
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "COMPLETED" && status != "ON_HOLD" {
		return nil, fmt.Errorf("invalid project status: %s (allowed: ACTIVE, COMPLETED, ON_HOLD)", status)
	}
	code, err := u.generateUniqueProjectCode(slugify(name), uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project code: %w", err)
	}
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
	if ctx.TeamsFeatureDisabled && ctx.Role == domain.RoleProductOwner && ctx.UserID != 0 {
		if err := u.repo.ReplaceProjectPmAssignments(p.ID, []uint{ctx.UserID}); err != nil {
			return nil, fmt.Errorf("failed to register creating Product Owner as project owner: %w", err)
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
func (u *sentinelUsecase) GetProjectDetailsPage(idOrCode string, taskLimit int, ctx domain.CallerContext) (*domain.ProjectDetailsResponse, error) {
	startedAt := time.Now()
	log.Printf("[ProjectDetails] start id=%s taskLimit=%d role=%s teamDisabled=%t", idOrCode, taskLimit, ctx.Role, ctx.TeamsFeatureDisabled)
	p, err := u.GetProjectByIDOrCode(idOrCode, ctx)
	if err != nil || p == nil {
		log.Printf("[ProjectDetails] project lookup failed id=%s elapsed=%s err=%v", idOrCode, time.Since(startedAt), err)
		return nil, err
	}
	if taskLimit <= 0 {
		taskLimit = 600
	}
	// Fetch all child data in parallel (4 queries → 1 DB round-trip per type; network already 1 round-trip).
	type result struct {
		tasks      []domain.Task
		sprints    []domain.ProjectDetailsSprint
		milestones []domain.Milestone
		epics      []domain.ProjectDetailsEpic
	}
	var res result
	var errTasks, errSprints, errMilestones, errEpics error
	var wg sync.WaitGroup
	queriesStartedAt := time.Now()
	log.Printf("[ProjectDetails] loading children projectID=%s tasksLimit=%d", p.ID, taskLimit)
	wg.Add(4)
	go func() {
		defer wg.Done()
		stepStartedAt := time.Now()
		res.tasks, errTasks = u.repo.GetTasksByProjectIDForProjectPage(p.ID, taskLimit)
		log.Printf("[ProjectDetails] tasks loaded projectID=%s count=%d elapsed=%s err=%v", p.ID, len(res.tasks), time.Since(stepStartedAt), errTasks)
	}()
	go func() {
		defer wg.Done()
		stepStartedAt := time.Now()
		sprints, err := u.repo.GetSprintsByProjectID(p.ID)
		if err == nil {
			res.sprints = make([]domain.ProjectDetailsSprint, len(sprints))
			for i := range sprints {
				s := sprints[i]
				res.sprints[i] = domain.ProjectDetailsSprint{
					ID: s.ID, ProjectID: s.ProjectID, Name: s.Name, Goal: s.Goal, StartDate: s.StartDate, EndDate: s.EndDate, Status: s.Status, SortOrder: s.SortOrder, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
				}
			}
		}
		errSprints = err
		log.Printf("[ProjectDetails] sprints loaded projectID=%s count=%d elapsed=%s err=%v", p.ID, len(res.sprints), time.Since(stepStartedAt), errSprints)
	}()
	go func() {
		defer wg.Done()
		stepStartedAt := time.Now()
		res.milestones, errMilestones = u.repo.GetMilestonesByProjectID(p.ID)
		log.Printf("[ProjectDetails] milestones loaded projectID=%s count=%d elapsed=%s err=%v", p.ID, len(res.milestones), time.Since(stepStartedAt), errMilestones)
	}()
	go func() {
		defer wg.Done()
		stepStartedAt := time.Now()
		epics, err := u.repo.GetEpicsByProjectID(p.ID)
		if err == nil {
			res.epics = make([]domain.ProjectDetailsEpic, len(epics))
			for i := range epics {
				e := epics[i]
				res.epics[i] = domain.ProjectDetailsEpic{
					ID: e.ID, ProjectID: e.ProjectID, Title: e.Title, Description: e.Description, Status: e.Status, Color: e.Color, SortOrder: e.SortOrder, StartDate: e.StartDate, EndDate: e.EndDate, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
				}
			}
		}
		errEpics = err
		log.Printf("[ProjectDetails] epics loaded projectID=%s count=%d elapsed=%s err=%v", p.ID, len(res.epics), time.Since(stepStartedAt), errEpics)
	}()
	wg.Wait()
	log.Printf("[ProjectDetails] children finished projectID=%s totalElapsed=%s queryElapsed=%s", p.ID, time.Since(startedAt), time.Since(queriesStartedAt))
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
	returned := len(res.tasks)
	hasMore := returned > taskLimit
	if hasMore {
		res.tasks = res.tasks[:taskLimit]
		returned = taskLimit
	}
	lightTasks := make([]domain.ProjectDetailsTask, len(res.tasks))
	for i := range res.tasks {
		t := res.tasks[i]
		lightTasks[i] = domain.ProjectDetailsTask{
			ID:                    t.ID,
			Code:                  t.Code,
			Title:                 t.Title,
			EstimatedMinutes:      t.EstimatedMinutes,
			ProjectID:             t.ProjectID,
			EpicID:                t.EpicID,
			SprintID:              t.SprintID,
			MilestoneID:           t.MilestoneID,
			TaskType:              t.TaskType,
			Priority:              t.Priority,
			StoryPoints:           t.StoryPoints,
			ParentID:              t.ParentID,
			SortOrder:             t.SortOrder,
			StartDate:             t.StartDate,
			EndDate:               t.EndDate,
			Progress:              t.Progress,
			DueAt:                 t.DueAt,
			StartedAt:             t.StartedAt,
			CompletedAt:           t.CompletedAt,
			Status:                t.Status,
			AssignedTo:            t.AssignedTo,
			AssignedToDisplayName: t.AssignedToDisplayName,
			AssignedToEmail:       t.AssignedToEmail,
			AssignedToAvatarURL:   t.AssignedToAvatarURL,
			CreatedAt:             t.CreatedAt,
			UpdatedAt:             t.UpdatedAt,
		}
	}
	log.Printf("[ProjectDetails] done projectID=%s taskCount=%d sprintCount=%d milestoneCount=%d epicCount=%d totalElapsed=%s", p.ID, returned, len(res.sprints), len(res.milestones), len(res.epics), time.Since(startedAt))
	return &domain.ProjectDetailsResponse{
		Project:    p,
		Tasks:      lightTasks,
		TasksMeta:  domain.ProjectDetailsTasksMeta{Limit: taskLimit, Returned: returned, HasMore: hasMore},
		Sprints:    res.sprints,
		Milestones: res.milestones,
		Epics:      res.epics,
	}, nil
}

func (u *sentinelUsecase) GetProjectTasksPage(idOrCode string, limit int, cursorCreatedAt, cursorID string, offset int, ctx domain.CallerContext) (*domain.ProjectTasksPageResponse, error) {
	p, err := u.GetProjectByIDOrCode(idOrCode, ctx)
	if err != nil || p == nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 600
	}

	var cursorTime *time.Time
	var cursorTaskID *uuid.UUID
	if strings.TrimSpace(cursorCreatedAt) != "" || strings.TrimSpace(cursorID) != "" {
		if strings.TrimSpace(cursorCreatedAt) == "" || strings.TrimSpace(cursorID) == "" {
			return nil, &domain.ErrBadRequest{Msg: "cursor_created_at and cursor_id must be provided together"}
		}
		parsedTime, parseErr := time.Parse(time.RFC3339Nano, strings.TrimSpace(cursorCreatedAt))
		if parseErr != nil {
			return nil, &domain.ErrBadRequest{Msg: "cursor_created_at must be RFC3339 timestamp"}
		}
		parsedID, parseErr := uuid.Parse(strings.TrimSpace(cursorID))
		if parseErr != nil {
			return nil, &domain.ErrBadRequest{Msg: "cursor_id must be valid UUID"}
		}
		cursorTime = &parsedTime
		cursorTaskID = &parsedID
	}

	tasks, err := u.repo.GetTasksByProjectIDForProjectPageCursor(p.ID, limit, cursorTime, cursorTaskID, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks page: %w", err)
	}
	returned := len(tasks)
	hasMore := returned > limit
	if hasMore {
		tasks = tasks[:limit]
		returned = limit
	}

	resp := &domain.ProjectTasksPageResponse{
		Tasks:    tasks,
		Limit:    limit,
		Returned: returned,
		HasMore:  hasMore,
	}
	if hasMore && returned > 0 {
		last := tasks[returned-1]
		if !last.CreatedAt.IsZero() {
			resp.NextCursor = &domain.ProjectTaskPageCursor{
				CreatedAt: last.CreatedAt.UTC().Format(time.RFC3339Nano),
				ID:        last.ID.String(),
			}
		}
		nextOffset := offset + returned
		resp.NextOffset = &nextOffset
	}
	return resp, nil
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

// AssignProjectPmOwners sets which Product Owner users own a project when squads are disabled (replaces the whole list).
func (u *sentinelUsecase) AssignProjectPmOwners(projectID uuid.UUID, pmUserIDs []uint, requesterRole string) (*domain.Project, error) {
	if requesterRole != domain.RoleCEO && requesterRole != domain.RoleManager {
		return nil, fmt.Errorf("unauthorized: only CEO or MANAGER can assign project Product Owners")
	}
	if !u.isTeamsFeatureDisabled() {
		return nil, &domain.ErrBadRequest{Msg: "project Product Owner assignments can only be edited when the teams feature is disabled"}
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
		if !authDomain.IsProductOwnerAssignableRole(uu.Role) {
			return nil, &domain.ErrBadRequest{Msg: fmt.Sprintf("user %d must have role PRODUCT_OWNER or PM (current role: %s)", id, uu.Role)}
		}
		clean = append(clean, id)
	}
	if err := u.repo.ReplaceProjectPmAssignments(projectID, clean); err != nil {
		return nil, fmt.Errorf("failed to save Product Owner assignments: %w", err)
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
	if !projectNameAllowedChars.MatchString(name) {
		return nil, errors.New("project name contains invalid characters (allowed: letters, numbers, spaces, hyphens, underscores)")
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
		newCode, err := u.generateUniqueProjectCode(slugify(name), p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate project code: %w", err)
		}
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

func (u *sentinelUsecase) generateUniqueProjectCode(base string, currentProjectID uuid.UUID) (string, error) {
	candidate := strings.TrimSpace(base)
	if candidate == "" {
		candidate = "task"
	}
	lookupCtx := domain.CallerContext{Role: domain.RoleCEO}
	for i := 0; i < 1000; i++ {
		code := candidate
		if i > 0 {
			code = fmt.Sprintf("%s-%d", candidate, i+1)
		}
		existing, err := u.repo.GetProjectByCode(code, lookupCtx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
		if existing == nil || existing.ID == currentProjectID {
			return code, nil
		}
	}
	return "", errors.New("unable to allocate unique project code")
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

	if projectID == nil {
		return nil, &domain.ErrBadRequest{Msg: "project_id is required"}
	}
	proj, err := u.repo.GetProjectByID(*projectID, domain.CallerContext{Role: domain.RoleCEO})
	if err != nil || proj == nil {
		return nil, &domain.ErrBadRequest{Msg: "project not found"}
	}
	slug := slugify(proj.Name)
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
	u.recordTaskActivity(task.ID, domain.TaskActivityCreated, &creatorID, map[string]interface{}{
		"title":     title,
		"task_type": taskType,
	})
	return task, nil
}

// authorizeTaskAssign allows CEO / Product Owner / MANAGER for any task.
// Engineer / Chief Engineer can claim an unassigned task only for themselves.
// For subtasks, parent assignee/creator or current subtask assignee can assign/reassign.
func (u *sentinelUsecase) authorizeTaskAssign(task *domain.Task, devID uint, assignerID uint, assignerRole string) error {
	role := strings.ToUpper(strings.TrimSpace(assignerRole))
	if role == domain.RoleCEO || role == domain.RoleProductOwner || role == domain.RoleManager {
		return nil
	}

	// Claim flow: Engineer / Chief Engineer can self-assign only when task is currently unassigned.
	if (role == domain.RoleEngineer || role == domain.RoleChiefEngineer || role == "CHIEF") && devID == assignerID && task.AssignedTo == nil {
		return nil
	}

	if task.ParentID == nil {
		return fmt.Errorf("unauthorized: only CEO, Product Owner, or Manager can assign top-level tasks")
	}
	if task.AssignedTo != nil && *task.AssignedTo == assignerID {
		return nil
	}
	parent, err := u.repo.GetTaskByID(*task.ParentID)
	if err != nil || parent == nil {
		return fmt.Errorf("parent task not found")
	}
	if parent.AssignedTo != nil && *parent.AssignedTo == assignerID {
		return nil
	}
	if parent.CreatedBy != nil && *parent.CreatedBy == assignerID {
		return nil
	}
	return fmt.Errorf("unauthorized: you cannot assign this task")
}

// AssignTask assigns a developer to a task. assignerID is recorded as assigned_by for Product Owner–scoped leaderboard when set by Product Owner/CEO.
// Parent-task assignees may assign subtasks; subtask assignees may reassign their own subtask.
// devID = 0 means unassign (set AssignedTo = nil, revert to PENDING).
func (u *sentinelUsecase) AssignTask(taskID uuid.UUID, devID uint, assignerID uint, assignerRole string) error {
	// 1. Validate if task exists
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	if err := u.authorizeTaskAssign(task, devID, assignerID, assignerRole); err != nil {
		return err
	}

	if devID == 0 {
		var prev uint
		if task.AssignedTo != nil {
			prev = *task.AssignedTo
		}
		// Unassign: clear assignee, revert status to PENDING
		task.AssignedTo = nil
		task.AssignedByID = nil
		task.Status = "PENDING"
		task.StartedAt = nil
		if err := u.repo.UpdateTask(task); err != nil {
			return err
		}
		aid := assignerID
		u.recordTaskActivity(taskID, domain.TaskActivityUnassigned, &aid, map[string]interface{}{
			"previous_assignee_user_id": prev,
			"to_status":                 "PENDING",
		})
		return nil
	}

	// 2. Update assignment only — status remains unchanged so dev can move the card themselves
	task.AssignedTo = &devID
	task.AssignedByID = &assignerID // Product Owner/CEO who assigned (drives Product Owner team leaderboard scope)

	// 3. Persist changes
	if err := u.repo.UpdateTask(task); err != nil {
		return err
	}
	aid := assignerID
	u.recordTaskActivity(taskID, domain.TaskActivityAssigned, &aid, map[string]interface{}{
		"assignee_user_id": devID,
		"assigner_user_id": assignerID,
	})

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
	u.recordTaskActivity(taskID, domain.TaskActivitySubmittedReview, &devID, map[string]interface{}{
		"reference_url": referenceURL,
		"to_status":     "REVIEW_PENDING",
	})

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
			if creator.DisplayName != "" {
				task.CreatedByDisplayName = creator.DisplayName
			}
		}
	}
	if task.AssignedTo != nil {
		assignee, err := u.authRepo.FindByID(*task.AssignedTo)
		if err == nil && assignee != nil {
			if assignee.DisplayName != "" {
				task.AssignedToDisplayName = assignee.DisplayName
			}
			task.AssignedToEmail = assignee.Email
			task.AssignedToAvatarURL = assignee.AvatarURL
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

// GetGlobalActiveTasks returns TASK/BUG in ACTIVE sprints, enriched with project name and color.
// CEO/MANAGER see all projects; when teams off: Product Owner → project_pm_assignments; engineer → projects with any task assigned to caller.
func (u *sentinelUsecase) GetGlobalActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	ctx = u.withCallerScope(ctx)
	return u.repo.GetGlobalActiveTasks(ctx)
}

// GetTeamActiveTasks returns active-sprint TASK/BUG items visible to the caller.
// Teams on: caller is team-scoped (except CEO/MANAGER). Teams off: visibility follows project assignment rules.
func (u *sentinelUsecase) GetTeamActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	ctx = u.withCallerScope(ctx)
	tasks, err := u.repo.GetTeamActiveTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetActiveFeatures returns FEATURE-type tasks for the Product Owner/CEO Feature Roadmap Board.
// Each feature carries a roll-up progress (0–100%) computed from child TASK/BUG completion.
// CEO/MANAGER see all teams; Product Owner is scoped to their own team.
// Optional projectID narrows result set to a single project.
func (u *sentinelUsecase) GetActiveFeatures(callerTeamID *uint, callerRole string, projectID *uuid.UUID) ([]domain.FeatureRoadmapItem, error) {
	teamID := uint(0)
	if callerRole != domain.RoleCEO && callerRole != domain.RoleManager && callerTeamID != nil {
		teamID = *callerTeamID
	}
	return u.repo.GetActiveFeatures(teamID, projectID)
}

// GetUnassignedTasks retrieves all tasks that are not assigned to anyone
func (u *sentinelUsecase) GetUnassignedTasks() ([]domain.Task, error) {
	tasks, err := u.repo.GetUnassignedTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetAllTasks retrieves all tasks in the system (for ADMIN / Product Owner view)
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

func (u *sentinelUsecase) GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]domain.Task, error) {
	return u.repo.GetTasksByProjectIDs(projectIDs)
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

// GetPendingApprovals returns tasks requiring Product Owner/CEO/MANAGER attention
// Includes: REVIEW_PENDING handovers, time negotiations (PENDING), appeals (PENDING)
func (u *sentinelUsecase) GetPendingApprovals(userRole string) ([]domain.Task, error) {
	if userRole != "CEO" && userRole != authDomain.RoleProductOwner && userRole != "MANAGER" {
		return nil, fmt.Errorf("access denied: only CEO, Product Owner, or MANAGER can view approvals inbox")
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

// ResolveAppeal allows Product Owner/CEO to approve or reject an appeal
func (u *sentinelUsecase) ResolveAppeal(appealID uuid.UUID, resolverID uint, status string, note string) error {
	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("status must be APPROVED or REJECTED")
	}

	resolver, err := u.authRepo.FindByID(resolverID)
	if err != nil {
		return fmt.Errorf("unauthorized: resolver user not found: %w", err)
	}

	if resolver.Role != authDomain.RoleCEO && resolver.Role != authDomain.RoleManager && resolver.Role != authDomain.RoleProductOwner {
		return fmt.Errorf("forbidden: only CEO, MANAGER, or Product Owner can resolve appeals (current role: %s)", resolver.Role)
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
			log.Printf("⚠️  Warning: Failed to get task for auto-completion: %v\n", err)
		} else if task.Status != "COMPLETED" {
			task.Status = "COMPLETED"
			now := time.Now()
			task.CompletedAt = &now
			if task.StartedAt == nil {
				task.StartedAt = &now
			}
			if err := u.repo.UpdateTask(task); err != nil {
				log.Printf("⚠️  Warning: Failed to auto-complete task: %v\n", err)
			} else {
				rid := resolverID
				u.recordTaskActivity(task.ID, domain.TaskActivityAppealComplete, &rid, map[string]interface{}{
					"appeal_id":   appealID.String(),
					"submission_id": submission.ID.String(),
					"to_status":   "COMPLETED",
				})
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

	// No AI: store negotiation for Product Owner/CEO to approve manually
	prevEst := task.EstimatedMinutes
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
	did := devID
	u.recordTaskActivity(taskID, domain.TaskActivityNegotiation, &did, map[string]interface{}{
		"proposed_minutes": minutes,
		"previous_minutes": prevEst,
	})

	return nil
}

// UpdateTask updates a task with access control (no AI).
// Creator, CEO, MANAGER, or Product Owner can update all fields.
// Assigned user can update description only.
func (u *sentinelUsecase) UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description, taskType string, parentID *uuid.UUID, dueAt, startDate, endDate *time.Time, progress *int, priority string, storyPoints *int, sprintID *uuid.UUID, applySprint bool, milestoneID *uuid.UUID, epicID *uuid.UUID, applyEpic bool, sortOrder *int, estimatedMinutes *int) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isManager := requestingUserRole == authDomain.RoleManager
	isSelectedProjectPO := u.isSelectedProjectProductOwner(task.ProjectID, requestingUserID, requestingUserRole)
	isAssignee := task.AssignedTo != nil && *task.AssignedTo == requestingUserID

	onlyDescriptionUpdate := description != "" &&
		title == "" &&
		taskType == "" &&
		parentID == nil &&
		dueAt == nil &&
		startDate == nil &&
		endDate == nil &&
		progress == nil &&
		priority == "" &&
		storyPoints == nil &&
		!applySprint &&
		milestoneID == nil &&
		!applyEpic &&
		sortOrder == nil &&
		estimatedMinutes == nil

	canFullyUpdate := isCreator || isCEO || isManager || isSelectedProjectPO
	canAssigneeUpdateDescription := isAssignee && onlyDescriptionUpdate

	if !canFullyUpdate && !canAssigneeUpdateDescription {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, MANAGER, selected Product Owner, or assignee (description only) can update this task")
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
	if dueAt != nil {
		task.DueAt = dueAt
	}
	if startDate != nil {
		task.StartDate = startDate
	}
	if endDate != nil {
		task.EndDate = endDate
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
	isManager := requestingUserRole == authDomain.RoleManager
	isSelectedProjectPO := u.isSelectedProjectProductOwner(task.ProjectID, requestingUserID, requestingUserRole)
	if !isCreator && !isCEO && !isManager && !isSelectedProjectPO {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, MANAGER, or selected Product Owner can update this task")
	}
	task.ResourceURLs = resourceURLs
	if err := u.repo.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task resource_urls: %w", err)
	}
	return task, nil
}

// EstimateTask uses AI to estimate task effort (title + description) and updates task.estimated_minutes.
// Used internally by ScheduleProjectWithAI. Task creator, CEO, MANAGER, or Product Owner can run estimate.
func (u *sentinelUsecase) EstimateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) (*domain.Task, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isManager := requestingUserRole == authDomain.RoleManager
	isSelectedProjectPO := u.isSelectedProjectProductOwner(task.ProjectID, requestingUserID, requestingUserRole)
	if !isCreator && !isCEO && !isManager && !isSelectedProjectPO {
		return nil, fmt.Errorf("unauthorized: only the task creator, CEO, MANAGER, or selected Product Owner can run AI estimate")
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
// Only CEO or Product Owner can run this.
func (u *sentinelUsecase) GenerateProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (*domain.AIGeneratedPlan, error) {
	if requestingUserRole != "CEO" && requestingUserRole != authDomain.RoleProductOwner {
		return nil, fmt.Errorf("unauthorized: only CEO or Product Owner can generate AI work plan")
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

// ClearProjectPlan removes all tasks, sprints, milestones, and epics for the project. Only CEO or Product Owner.
func (u *sentinelUsecase) ClearProjectPlan(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
	if requestingUserRole != "CEO" && requestingUserRole != authDomain.RoleProductOwner {
		return fmt.Errorf("unauthorized: only CEO or Product Owner can clear project plan")
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

// ScheduleProjectWithAI ประเมินเวลาและจัดเรียง timeline ของ task ที่มีอยู่แล้ว (ไม่สร้าง task ใหม่). เฉพาะ CEO / Product Owner.
func (u *sentinelUsecase) ScheduleProjectWithAI(projectID uuid.UUID, requestingUserID uint, requestingUserRole string) (int, error) {
	if requestingUserRole != "CEO" && requestingUserRole != authDomain.RoleProductOwner {
		return 0, fmt.Errorf("unauthorized: only CEO or Product Owner can run AI schedule")
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
	ordered := make([]struct {
		taskIdx int
		res     domain.TaskEstimateAndOrder
	}, 0, len(byIndex))
	for idx, res := range byIndex {
		ordered = append(ordered, struct {
			taskIdx int
			res     domain.TaskEstimateAndOrder
		}{idx, res})
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

// DeleteTask deletes a task with access control.
// Creator, CEO, MANAGER, or selected Product Owner can delete.
func (u *sentinelUsecase) DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
	// 1️⃣ Fetch the task to check ownership
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2️⃣ ACCESS CONTROL: Creator, CEO, MANAGER, or Product Owner can delete.
	isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
	isCEO := requestingUserRole == "CEO"
	isManager := requestingUserRole == authDomain.RoleManager
	isProductOwner := requestingUserRole == authDomain.RoleProductOwner

	if !isCreator && !isCEO && !isManager && !isProductOwner {
		return fmt.Errorf("unauthorized: only the task creator, CEO, MANAGER, or Product Owner can delete this task")
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

// ApproveTask marks a task as COMPLETED after human verification (Product Owner, CEO, or Manager)
func (u *sentinelUsecase) ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error {
	if approverRole != authDomain.RoleCEO && approverRole != authDomain.RoleProductOwner && approverRole != authDomain.RoleManager {
		return fmt.Errorf("access denied: only Product Owner, CEO, or Manager can approve tasks (your role: %s)", approverRole)
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
	u.recordTaskActivity(taskID, domain.TaskActivityApprovedReview, &approverID, map[string]interface{}{
		"from_status": "REVIEW_PENDING",
		"to_status":   "COMPLETED",
	})

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
					u.recordTaskActivity(parent.ID, domain.TaskActivityParentRollupStatus, &approverID, map[string]interface{}{
						"child_task_id": taskID.String(),
						"to_status":     "READY_FOR_UAT",
					})
				}
			}
		}
	}

	return nil
}

// SubmitUAT stores the UAT payload on a FEATURE task and moves it to REVIEW_PENDING for Product Owner/CEO to review.
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
	u.recordTaskActivity(taskID, domain.TaskActivitySubmitUAT, &devID, map[string]interface{}{
		"from_status":  "READY_FOR_UAT",
		"to_status":    "REVIEW_PENDING",
		"staging_url":  payload.StagingURL,
		"release_notes": strings.TrimSpace(payload.ReleaseNotes),
	})
	return nil
}

// RejectTask returns a task to IN_PROGRESS and logs rejection reason as a comment (Product Owner/CEO/MANAGER only)
func (u *sentinelUsecase) RejectTask(taskID uuid.UUID, rejectorID uint, rejectorRole string, reason string) error {
	if rejectorRole != "CEO" && rejectorRole != authDomain.RoleProductOwner && rejectorRole != "MANAGER" {
		return fmt.Errorf("access denied: only Product Owner, CEO or MANAGER can reject tasks (your role: %s)", rejectorRole)
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

	if err := u.repo.RejectTask(taskID, rejectorID, reason); err != nil {
		return err
	}
	rid := rejectorID
	u.recordTaskActivity(taskID, domain.TaskActivityRejectedReview, &rid, map[string]interface{}{
		"from_status": "REVIEW_PENDING",
		"to_status":   "IN_PROGRESS",
	})
	return nil
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
	u.recordTaskActivity(taskID, domain.TaskActivityReadyForTest, &devID, map[string]interface{}{
		"from_status": "IN_PROGRESS",
		"to_status":   "READY_FOR_TEST",
	})
	return nil
}

// PMApproveSubTask is the Product Owner's first-stage approval: READY_FOR_TEST → WAIT_FOR_DEPLOY.
// The PO provides test evidence (URL + steps); the task waits for a Chief Engineer deployment request.
func (u *sentinelUsecase) PMApproveSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, testURL string, testSteps string) error {
	if pmRole != authDomain.RoleProductOwner && pmRole != "MANAGER" {
		return fmt.Errorf("access denied: only Product Owner or MANAGER can approve test (your role: %s)", pmRole)
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
	if err := u.repo.SetTaskWaitForDeploy(taskID, []byte(payload)); err != nil {
		return fmt.Errorf("failed to submit for deployment: %w", err)
	}
	pmid := pmUserID
	u.recordTaskActivity(taskID, domain.TaskActivityPMApprovedDeploy, &pmid, map[string]interface{}{
		"from_status": "READY_FOR_TEST",
		"to_status":   "WAIT_FOR_DEPLOY",
		"test_url":    strings.TrimSpace(testURL),
	})
	return nil
}

// AdvanceTaskAfterDeploy moves a task from WAIT_FOR_DEPLOY → READY_FOR_UAT.
// Called automatically by the deployment module when the Chief Engineer marks a request as deployed.
func (u *sentinelUsecase) AdvanceTaskAfterDeploy(taskID uuid.UUID, deployedByUserID uint) error {
	if err := u.repo.AdvanceTaskToReadyForUAT(taskID); err != nil {
		return fmt.Errorf("advance-after-deploy: %w", err)
	}
	deployer := deployedByUserID
	u.recordTaskActivity(taskID, domain.TaskActivityDeployed, &deployer, map[string]interface{}{
		"from_status": "WAIT_FOR_DEPLOY",
		"to_status":   "READY_FOR_UAT",
	})
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
		return fmt.Errorf("task is not in READY_FOR_UAT status (current: %s) — Product Owner must submit test evidence first", task.Status)
	}

	if err := u.repo.ApproveTask(taskID); err != nil {
		return fmt.Errorf("failed to approve sub-task: %w", err)
	}
	ceoid := ceoUserID
	u.recordTaskActivity(taskID, domain.TaskActivityCEOFinalApproved, &ceoid, map[string]interface{}{
		"from_status": "READY_FOR_UAT",
		"to_status":   "COMPLETED",
	})

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
					u.recordTaskActivity(parent.ID, domain.TaskActivityParentRollupStatus, &ceoUserID, map[string]interface{}{
						"child_task_id": taskID.String(),
						"to_status":     "READY_FOR_UAT",
					})
				}
			}
		}
	}

	return nil
}

// RejectSubTask returns a task to IN_PROGRESS with a reason log.
// Accepts READY_FOR_TEST (Product Owner rejecting engineer work) or READY_FOR_UAT (CEO rejecting Product Owner evidence).
func (u *sentinelUsecase) RejectSubTask(taskID uuid.UUID, pmUserID uint, pmRole string, reason string) error {
	if pmRole != "CEO" && pmRole != authDomain.RoleProductOwner && pmRole != "MANAGER" {
		return fmt.Errorf("access denied: only Product Owner, CEO or MANAGER can reject sub-tasks (your role: %s)", pmRole)
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
	oldSt := task.Status
	if err := u.repo.RejectTask(taskID, pmUserID, reason); err != nil {
		return err
	}
	rj := pmUserID
	reasonExcerpt := reason
	if len(reasonExcerpt) > 160 {
		reasonExcerpt = reasonExcerpt[:160] + "…"
	}
	u.recordTaskActivity(taskID, domain.TaskActivityWorkflowReject, &rj, map[string]interface{}{
		"from_status":   oldSt,
		"to_status":     "IN_PROGRESS",
		"reason_excerpt": reasonExcerpt,
	})
	return nil
}

// GetTasksReadyForTest returns all TASK/BUG items in READY_FOR_TEST status, scoped to the caller's team (Product Owner/CEO).
func (u *sentinelUsecase) GetTasksReadyForTest(callerTeamID *uint, callerRole string) ([]domain.GlobalActiveTask, error) {
	if callerRole != "CEO" && callerRole != authDomain.RoleProductOwner && callerRole != "MANAGER" {
		return nil, fmt.Errorf("access denied: only Product Owner, CEO or MANAGER can view the test queue (your role: %s)", callerRole)
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

func (u *sentinelUsecase) AddComment(taskID uuid.UUID, userID uint, content string, attachments []domain.TaskCommentAttachment) (*domain.TaskComment, error) {
	content = strings.TrimSpace(content)
	if content == "" && len(attachments) == 0 {
		return nil, errors.New("comment content or attachments are required")
	}

	attachmentsJSON := datatypes.JSON([]byte("[]"))
	if len(attachments) > 0 {
		raw, err := json.Marshal(attachments)
		if err != nil {
			return nil, fmt.Errorf("failed to encode attachments: %w", err)
		}
		attachmentsJSON = datatypes.JSON(raw)
	}
	c := &domain.TaskComment{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      userID,
		Content:     content,
		Attachments: attachmentsJSON,
	}
	if err := u.repo.CreateTaskComment(c); err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}
	user, err := u.authRepo.FindByID(userID)
	if err == nil && user != nil {
		c.UserEmail = user.Email
		c.UserDisplayName = user.DisplayName
		c.UserAvatarURL = user.AvatarURL
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
			comments[i].UserDisplayName = user.DisplayName
			comments[i].UserAvatarURL = user.AvatarURL
		}
	}
	return comments, nil
}

func (u *sentinelUsecase) EditComment(commentID uuid.UUID, editorUserID uint, content string) (*domain.TaskComment, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, &domain.ErrBadRequest{Msg: "comment content is required"}
	}
	comment, err := u.repo.GetTaskCommentByID(commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, &domain.ErrBadRequest{Msg: "comment not found"}
	}
	if comment.UserID != editorUserID {
		return nil, &domain.ErrBadRequest{Msg: "you can only edit your own comments"}
	}
	if comment.Content == trimmed {
		user, _ := u.authRepo.FindByID(comment.UserID)
		if user != nil {
			comment.UserEmail = user.Email
			comment.UserDisplayName = user.DisplayName
			comment.UserAvatarURL = user.AvatarURL
		}
		return comment, nil
	}

	now := time.Now().UTC()
	history := make([]domain.TaskCommentEditHistoryItem, 0)
	if len(comment.EditHistory) > 0 {
		_ = json.Unmarshal(comment.EditHistory, &history)
	}
	history = append(history, domain.TaskCommentEditHistoryItem{
		EditedAt:   now,
		EditedBy:   editorUserID,
		OldContent: comment.Content,
		NewContent: trimmed,
	})
	rawHistory, err := json.Marshal(history)
	if err != nil {
		return nil, fmt.Errorf("failed to encode edit history: %w", err)
	}

	comment.Content = trimmed
	comment.EditedAt = &now
	comment.EditHistory = datatypes.JSON(rawHistory)
	if err := u.repo.UpdateTaskComment(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}
	user, err := u.authRepo.FindByID(comment.UserID)
	if err == nil && user != nil {
		comment.UserEmail = user.Email
		comment.UserDisplayName = user.DisplayName
		comment.UserAvatarURL = user.AvatarURL
	}
	return comment, nil
}

// --- Time Log Operations ---

func (u *sentinelUsecase) LogTime(taskID uuid.UUID, userID uint, minutes int, description, workType string, loggedDate *time.Time, isTimer bool) (*domain.TimeLog, error) {
	if minutes <= 0 {
		return nil, &domain.ErrBadRequest{Msg: "minutes must be greater than 0"}
	}
	if minutes > 960 {
		return nil, &domain.ErrBadRequest{Msg: "cannot log more than 16 hours (960 minutes) in a single entry"}
	}

	// Validate and default work_type
	wt := strings.ToUpper(strings.TrimSpace(workType))
	if wt == "" {
		wt = "DEV"
	}
	if !domain.ValidWorkTypes[wt] {
		return nil, &domain.ErrBadRequest{Msg: "invalid work_type; allowed: DEV, REVIEW, TESTING, MEETING, RESEARCH, OTHER"}
	}

	// Validate logged_date: default today, allow up to 7 days back, no future.
	// We allow the client date to be up to 1 calendar day ahead of UTC today to
	// accommodate users in UTC+ timezones (e.g. UTC+7 Bangkok, UTC+14 max).
	now := time.Now().UTC()
	todayUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	date := todayUTC
	if loggedDate != nil {
		d := time.Date(loggedDate.Year(), loggedDate.Month(), loggedDate.Day(), 0, 0, 0, 0, time.UTC)
		if d.After(todayUTC.AddDate(0, 0, 1)) {
			return nil, &domain.ErrBadRequest{Msg: "cannot log time for a future date"}
		}
		if todayUTC.Sub(d).Hours() > 7*24 {
			return nil, &domain.ErrBadRequest{Msg: "cannot backfill time logs older than 7 days"}
		}
		date = d
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
		ID:             uuid.New(),
		TaskID:         taskID,
		UserID:         userID,
		Minutes:        minutes,
		Description:    strings.TrimSpace(description),
		WorkType:       wt,
		LoggedDate:     date,
		IsTimerSession: isTimer,
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
	// Repository now does JOIN — no N+1 here
	return u.repo.GetTimeLogsByTaskID(taskID)
}

func (u *sentinelUsecase) EditTimeLog(logID uuid.UUID, callerID uint, minutes int, description, workType string, taskID *uuid.UUID) (*domain.TimeLog, error) {
	existing, err := u.repo.GetTimeLogByID(logID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time log: %w", err)
	}
	if existing == nil {
		return nil, &domain.ErrBadRequest{Msg: "time log not found"}
	}
	if existing.UserID != callerID {
		return nil, &domain.ErrBadRequest{Msg: "you can only edit your own time logs"}
	}
	if time.Since(existing.LoggedAt) > 24*time.Hour {
		return nil, &domain.ErrBadRequest{Msg: "time logs can only be edited within 24 hours of creation"}
	}
	if minutes <= 0 {
		return nil, &domain.ErrBadRequest{Msg: "minutes must be greater than 0"}
	}
	if minutes > 960 {
		return nil, &domain.ErrBadRequest{Msg: "cannot log more than 16 hours (960 minutes) in a single entry"}
	}
	wt := strings.ToUpper(strings.TrimSpace(workType))
	if wt == "" {
		wt = existing.WorkType
	}
	if !domain.ValidWorkTypes[wt] {
		return nil, &domain.ErrBadRequest{Msg: "invalid work_type"}
	}
	existing.Minutes = minutes
	existing.Description = strings.TrimSpace(description)
	existing.WorkType = wt
	if taskID != nil {
		existing.TaskID = *taskID
	}
	if err := u.repo.UpdateTimeLog(existing); err != nil {
		return nil, fmt.Errorf("failed to update time log: %w", err)
	}
	return existing, nil
}

func (u *sentinelUsecase) DeleteTimeLog(logID uuid.UUID, callerID uint) error {
	existing, err := u.repo.GetTimeLogByID(logID)
	if err != nil {
		return fmt.Errorf("failed to fetch time log: %w", err)
	}
	if existing == nil {
		return &domain.ErrBadRequest{Msg: "time log not found"}
	}
	if existing.UserID != callerID {
		return &domain.ErrBadRequest{Msg: "you can only delete your own time logs"}
	}
	if time.Since(existing.LoggedAt) > 24*time.Hour {
		return &domain.ErrBadRequest{Msg: "time logs can only be deleted within 24 hours of creation"}
	}
	return u.repo.DeleteTimeLog(logID)
}

func (u *sentinelUsecase) GetMyDailyTimeLogs(userID uint, date time.Time) (*domain.DailyTimeLogSummary, error) {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	logs, err := u.repo.GetTimeLogsByUserAndDate(userID, d)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily logs: %w", err)
	}
	total := 0
	for _, l := range logs {
		total += l.Minutes
	}
	if logs == nil {
		logs = []domain.TimeLog{}
	}
	return &domain.DailyTimeLogSummary{
		Date:         d.Format("2006-01-02"),
		TotalMinutes: total,
		Entries:      logs,
	}, nil
}

// BulkLogTime processes multiple time log entries in one call (EOD batch).
// Each entry is validated independently — failures don't abort the batch.
func (u *sentinelUsecase) BulkLogTime(entries []domain.BulkLogEntry, userID uint) ([]domain.BulkLogResult, error) {
	if len(entries) == 0 {
		return nil, &domain.ErrBadRequest{Msg: "entries must not be empty"}
	}
	if len(entries) > 20 {
		return nil, &domain.ErrBadRequest{Msg: "bulk log is limited to 20 entries per request"}
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	results := make([]domain.BulkLogResult, 0, len(entries))
	var logsToCreate []domain.TimeLog

	for _, entry := range entries {
		res := domain.BulkLogResult{TaskID: entry.TaskID}

		taskID, err := uuid.Parse(entry.TaskID)
		if err != nil {
			res.Error = "invalid task_id format"
			results = append(results, res)
			continue
		}
		if entry.Minutes <= 0 || entry.Minutes > 960 {
			res.Error = "minutes must be between 1 and 960"
			results = append(results, res)
			continue
		}

		wt := strings.ToUpper(strings.TrimSpace(entry.WorkType))
		if wt == "" {
			wt = "DEV"
		}
		if !domain.ValidWorkTypes[wt] {
			res.Error = "invalid work_type"
			results = append(results, res)
			continue
		}

		logDate := today
		if entry.LoggedDate != nil {
			d, err := time.Parse("2006-01-02", *entry.LoggedDate)
			if err != nil {
				res.Error = "invalid logged_date format; expected YYYY-MM-DD"
				results = append(results, res)
				continue
			}
			d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
			if d.After(today) {
				res.Error = "cannot log time for a future date"
				results = append(results, res)
				continue
			}
			if today.Sub(d).Hours() > 7*24 {
				res.Error = "cannot backfill time logs older than 7 days"
				results = append(results, res)
				continue
			}
			logDate = d
		}

		// Guard: no logging on parent tasks
		childCount, err := u.repo.CountChildTasks(taskID)
		if err != nil || childCount > 0 {
			res.Error = "cannot log time on a parent task; use sub-tasks"
			results = append(results, res)
			continue
		}

		tl := domain.TimeLog{
			ID:          uuid.New(),
			TaskID:      taskID,
			UserID:      userID,
			Minutes:     entry.Minutes,
			Description: strings.TrimSpace(entry.Description),
			WorkType:    wt,
			LoggedDate:  logDate,
		}
		res.Success = true
		res.Log = &tl
		logsToCreate = append(logsToCreate, tl)
		results = append(results, res)
	}

	if len(logsToCreate) > 0 {
		if err := u.repo.BulkCreateTimeLogs(logsToCreate); err != nil {
			return nil, fmt.Errorf("failed to save bulk time logs: %w", err)
		}
	}
	return results, nil
}

// --- Analytics ---

func (u *sentinelUsecase) GetProjectAnalytics(projectID uuid.UUID) (*domain.ProjectAnalytics, error) {
	return u.repo.GetProjectAnalytics(projectID)
}

// --- Bulk Operations ---

func (u *sentinelUsecase) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string, actorID uint) error {
	if len(taskIDs) == 0 {
		return errors.New("no tasks provided")
	}

	validStatuses := map[string]bool{
		"PENDING": true, "IN_PROGRESS": true, "READY_FOR_TEST": true,
		"REVIEW_PENDING": true, "WAIT_FOR_DEPLOY": true, "BLOCKED": true, "COMPLETED": true,
		// COMPLETED permission is enforced at handler level (CEO or Manager)
		// READY_FOR_UAT is intentionally excluded — set automatically on deployment
	}
	if !validStatuses[status] {
		if status == "READY_FOR_UAT" {
			return fmt.Errorf("READY_FOR_UAT is set automatically when the Chief Engineer marks the deployment as deployed")
		}
		return fmt.Errorf("invalid status: %s", status)
	}

	type bulkPrev struct {
		Status   string
		TaskType string
	}
	prevByID := make(map[uuid.UUID]bulkPrev, len(taskIDs))
	for _, id := range taskIDs {
		t, err := u.repo.GetTaskByID(id)
		if err != nil || t == nil {
			continue
		}
		prevByID[id] = bulkPrev{Status: t.Status, TaskType: t.TaskType}
	}

	if err := u.repo.BulkUpdateTaskStatus(taskIDs, status); err != nil {
		return err
	}

	aid := actorID
	for _, id := range taskIDs {
		prev, ok := prevByID[id]
		if !ok || prev.Status == status {
			continue
		}
		u.recordTaskActivity(id, domain.TaskActivityStatusChanged, &aid, map[string]interface{}{
			"from_status": prev.Status,
			"to_status":   status,
			"source":      "kanban_bulk",
		})
		// Discipline "Rework" counts task_comments LIKE '[REJECTED]%'. Dragging from Ready for Test → In Progress is rework.
		if status == "IN_PROGRESS" && u.prevStatusWasReadyForTestColumn(prev.Status, prev.TaskType) {
			u.recordKanbanReworkComment(id, actorID)
		}
	}
	return nil
}

// prevStatusWasReadyForTestColumn is true when the task was shown in the Kanban "Ready for Test" column before the move.
func (u *sentinelUsecase) prevStatusWasReadyForTestColumn(oldStatus, taskType string) bool {
	if oldStatus == "READY_FOR_TEST" {
		return true
	}
	if oldStatus == "REVIEW_PENDING" && (taskType == "TASK" || taskType == "BUG") {
		return true
	}
	return false
}

// recordKanbanReworkComment persists a [REJECTED] comment so performance/discipline rework metrics include Kanban send-backs.
func (u *sentinelUsecase) recordKanbanReworkComment(taskID uuid.UUID, actorID uint) {
	c := &domain.TaskComment{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      actorID,
		Content:     "[REJECTED] Moved back to In Progress from Ready for Test (Kanban).",
		Attachments: datatypes.JSON([]byte("[]")),
	}
	if err := u.repo.CreateTaskComment(c); err != nil {
		log.Printf("kanban rework comment: failed task=%s: %v", taskID, err)
	}
}

// SplitTask decomposes one task into N new sub-tasks (inheriting same parent_id, project, epic, sprint)
// then deletes the original task. The caller must be CEO, Product Owner, or the task creator.
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

	// Access control: CEO / Product Owner / MANAGER / creator / parent assignee / current subtask assignee
	role := strings.ToUpper(strings.TrimSpace(requestingUserRole))
	allowed := role == "CEO" || role == authDomain.RoleProductOwner || role == "MANAGER"
	if !allowed && orig.CreatedBy != nil && *orig.CreatedBy == requestingUserID {
		allowed = true
	}
	if !allowed && orig.AssignedTo != nil && *orig.AssignedTo == requestingUserID {
		allowed = true
	}
	if !allowed && orig.ParentID != nil {
		if parent, perr := u.repo.GetTaskByID(*orig.ParentID); perr == nil && parent != nil && parent.AssignedTo != nil && *parent.AssignedTo == requestingUserID {
			allowed = true
		}
	}
	if !allowed {
		return nil, &domain.ErrBadRequest{Msg: "only CEO, Product Owner, Manager, the task creator, the parent task assignee, or the subtask assignee can split this task"}
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
			ParentID:         orig.ParentID, // same parent as the original
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

func (u *sentinelUsecase) isSelectedProjectProductOwner(projectID *uuid.UUID, userID uint, role string) bool {
	if strings.ToUpper(strings.TrimSpace(role)) != authDomain.RoleProductOwner {
		return false
	}
	if projectID == nil || userID == 0 {
		return false
	}
	proj, err := u.repo.GetProjectByID(*projectID, domain.CallerContext{Role: domain.RoleCEO})
	if err != nil || proj == nil {
		return false
	}
	for _, owner := range proj.PmOwners {
		if owner.UserID == userID {
			return true
		}
	}
	return false
}

func (u *sentinelUsecase) GetB2BRequests(callerTeamID uint, direction string) ([]domain.B2BRequest, error) {
	if direction != "inbound" && direction != "outbound" {
		direction = "inbound"
	}
	reqs, err := u.repo.GetB2BRequests(callerTeamID, direction)
	if err != nil {
		return nil, err
	}

	// Enrich team names via authRepo (load once, then O(1) lookup)
	teamCache := map[uint]string{}
	if teams, err2 := u.authRepo.GetAllTeams(); err2 == nil {
		for _, t := range teams {
			teamCache[t.ID] = t.Name
		}
	}
	for i := range reqs {
		reqs[i].RequesterTeamName = teamCache[reqs[i].RequesterTeamID]
		reqs[i].TargetTeamName = teamCache[reqs[i].TargetTeamID]
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

// AcceptB2BRequest is called by the target team's Product Owner.
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
		nil, // no due date
		&project.ID,
		nil,      // no parent
		nil, nil, // no start/end date
		"MEDIUM",
		0,        // no story points
		nil, nil, // no sprint/milestone
		nil, // no epic
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

// --- Komgrip Operations ---

func (u *sentinelUsecase) CreateKomgripTask(title, description string, creatorID uint, priority string, estimatedMinutes int) (*domain.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, &domain.ErrBadRequest{Msg: "title is required"}
	}
	if priority == "" {
		priority = "MEDIUM"
	}
	if !validPriorities[priority] {
		return nil, &domain.ErrBadRequest{Msg: fmt.Sprintf("invalid priority: %s (allowed: CRITICAL, HIGH, MEDIUM, LOW)", priority)}
	}
	if estimatedMinutes < 0 {
		estimatedMinutes = 0
	}

	maxSuffix, err := u.repo.GetMaxTaskCodeSuffix("komgrip")
	if err != nil {
		return nil, fmt.Errorf("failed to get next komgrip code: %w", err)
	}
	code := fmt.Sprintf("komgrip-%03d", maxSuffix+1)

	task := &domain.Task{
		ID:               uuid.New(),
		Code:             code,
		Title:            strings.TrimSpace(title),
		Description:      description,
		TaskType:         string(domain.TaskTypeTask),
		CreatedBy:        &creatorID,
		Status:           "PENDING",
		EstimatedMinutes: estimatedMinutes,
		Priority:         priority,
		IsKomgrip:        true,
	}

	if err := u.repo.CreateTask(task); err != nil {
		return nil, err
	}
	u.recordTaskActivity(task.ID, domain.TaskActivityCreated, &creatorID, map[string]interface{}{
		"title":      title,
		"is_komgrip": true,
	})
	return task, nil
}

func (u *sentinelUsecase) GetKomgripTasks(userID uint) ([]domain.Task, error) {
	return u.repo.GetKomgripTasks(userID)
}

func (u *sentinelUsecase) UpdateKomgripTaskStatus(taskID uuid.UUID, status string, userID uint) (*domain.Task, error) {
	switch status {
	case "PENDING", "COMPLETED":
		// valid for Komgrip
	default:
		return nil, &domain.ErrBadRequest{Msg: "komgrip tasks only support status: PENDING, COMPLETED"}
	}

	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, &domain.ErrBadRequest{Msg: "task not found"}
	}
	if !task.IsKomgrip {
		return nil, &domain.ErrBadRequest{Msg: "task is not a Komgrip task"}
	}

	task.Status = status
	if status == "COMPLETED" {
		now := time.Now()
		task.CompletedAt = &now
		// Credit the person who marked it done so discipline tracker counts it as Job Done
		if task.AssignedTo == nil {
			task.AssignedTo = &userID
		}
	} else {
		task.CompletedAt = nil
	}

	if err := u.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	u.recordTaskActivity(task.ID, domain.TaskActivityStatusChanged, &userID, map[string]interface{}{
		"status":     status,
		"is_komgrip": true,
	})
	return task, nil
}
