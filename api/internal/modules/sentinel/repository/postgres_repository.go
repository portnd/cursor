package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository is the constructor
func NewPostgresRepository(db *gorm.DB) domain.SentinelRepository {
	return &postgresRepository{db: db}
}

// --- Project Operations ---

// teamScope applies a team-based filter for non-CEO/MANAGER users.
// CEO and MANAGER bypass all filters; Product Owner / engineer are restricted to their team's projects.
func teamScope(db *gorm.DB, ctx domain.CallerContext) *gorm.DB {
	if ctx.Role == domain.RoleCEO || ctx.Role == domain.RoleManager || ctx.TeamID == nil {
		return db
	}
	return db.Where("team_id = ?", *ctx.TeamID)
}

// projectAccessScope restricts project rows for list/detail queries.
// When teams feature is disabled: Product Owners see project_pm_assignments; engineers see projects with at least one task assigned_to them.
func projectAccessScope(db *gorm.DB, ctx domain.CallerContext) *gorm.DB {
	if ctx.Role == domain.RoleCEO || ctx.Role == domain.RoleManager {
		return db
	}
	if ctx.TeamsFeatureDisabled && ctx.UserID != 0 {
		if ctx.Role == domain.RoleProductOwner {
			return db.Where("id IN (SELECT project_id FROM project_pm_assignments WHERE user_id = ?)", ctx.UserID)
		}
		if domain.IsEngineerRole(ctx.Role) {
			return db.Where("id IN (SELECT DISTINCT project_id FROM tasks WHERE assigned_to = ?)", ctx.UserID)
		}
	}
	return teamScope(db, ctx)
}

func (r *postgresRepository) CreateProject(p *domain.Project) error {
	return r.db.Create(p).Error
}

func (r *postgresRepository) fillProjectPmOwners(projects []domain.Project) error {
	if len(projects) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, len(projects))
	for i := range projects {
		ids[i] = projects[i].ID
	}
	type row struct {
		ProjectID   uuid.UUID `gorm:"column:project_id"`
		UserID      uint      `gorm:"column:user_id"`
		Email       string
		DisplayName string `gorm:"column:display_name"`
	}
	var rows []row
	err := r.db.Raw(`
		SELECT ppa.project_id, u.id AS user_id, u.email, u.display_name
		FROM project_pm_assignments ppa
		JOIN users u ON u.id = ppa.user_id
		WHERE ppa.project_id IN ?
		ORDER BY ppa.project_id, u.email
	`, ids).Scan(&rows).Error
	if err != nil {
		return err
	}
	byProj := make(map[uuid.UUID][]domain.ProjectPmOwner)
	for _, row := range rows {
		byProj[row.ProjectID] = append(byProj[row.ProjectID], domain.ProjectPmOwner{
			UserID: row.UserID, Email: row.Email, DisplayName: row.DisplayName,
		})
	}
	for i := range projects {
		if owners, ok := byProj[projects[i].ID]; ok {
			projects[i].PmOwners = owners
		} else {
			projects[i].PmOwners = nil
		}
	}
	return nil
}

func (r *postgresRepository) GetAllProjects(ctx domain.CallerContext) ([]domain.Project, error) {
	var projects []domain.Project
	err := projectAccessScope(r.db.Model(&domain.Project{}), ctx).
		Order("created_at desc").
		Find(&projects).Error
	if err != nil {
		return projects, err
	}
	if err := r.fillProjectPmOwners(projects); err != nil {
		return projects, err
	}
	if len(projects) == 0 {
		return projects, nil
	}
	projectIDs := make([]uuid.UUID, len(projects))
	for i := range projects {
		projectIDs[i] = projects[i].ID
	}
	var countRows []struct {
		ProjectID uuid.UUID
		Total     int
		Completed int
		Overdue   int
	}
	err = r.db.Raw(`
		SELECT project_id,
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE status = 'COMPLETED')::int AS completed,
			COUNT(*) FILTER (WHERE status != 'COMPLETED' AND due_at IS NOT NULL AND due_at < NOW())::int AS overdue
		FROM tasks
		WHERE project_id IN ?
		GROUP BY project_id
	`, projectIDs).Scan(&countRows).Error
	if err != nil {
		return projects, err
	}
	countByID := make(map[uuid.UUID]struct{ Total, Completed, Overdue int })
	for _, row := range countRows {
		countByID[row.ProjectID] = struct{ Total, Completed, Overdue int }{row.Total, row.Completed, row.Overdue}
	}
	for i := range projects {
		if c, ok := countByID[projects[i].ID]; ok {
			projects[i].TaskTotal = c.Total
			projects[i].TaskCompleted = c.Completed
			projects[i].TaskOverdue = c.Overdue
		}
	}
	return projects, nil
}

func (r *postgresRepository) GetProjectByID(id uuid.UUID, ctx domain.CallerContext) (*domain.Project, error) {
	var project domain.Project
	err := projectAccessScope(r.db, ctx).First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	slice := []domain.Project{project}
	if err := r.fillProjectPmOwners(slice); err != nil {
		return nil, err
	}
	project = slice[0]
	return &project, nil
}

func (r *postgresRepository) GetProjectByCode(code string, ctx domain.CallerContext) (*domain.Project, error) {
	var project domain.Project
	err := projectAccessScope(r.db, ctx).First(&project, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	slice := []domain.Project{project}
	if err := r.fillProjectPmOwners(slice); err != nil {
		return nil, err
	}
	project = slice[0]
	return &project, nil
}

func (r *postgresRepository) AssignProjectTeam(projectID uuid.UUID, teamID *uint) error {
	result := r.db.Model(&domain.Project{}).Where("id = ?", projectID).Update("team_id", teamID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (r *postgresRepository) ReplaceProjectPmAssignments(projectID uuid.UUID, userIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectID).Delete(&domain.ProjectPmAssignment{}).Error; err != nil {
			return err
		}
		for _, uid := range userIDs {
			if uid == 0 {
				continue
			}
			row := domain.ProjectPmAssignment{ProjectID: projectID, UserID: uid}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *postgresRepository) UpdateProject(p *domain.Project) error {
	return r.db.Save(p).Error
}

func (r *postgresRepository) GetTasksByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	type taskRow struct {
		domain.Task
		DisplayName string `gorm:"column:assigned_to_display_name"`
		Email       string `gorm:"column:assigned_to_email"`
		AvatarURL   string `gorm:"column:assigned_to_avatar_url"`
	}
	var rows []taskRow
	err := r.db.Table("tasks").
		Select(`tasks.*,
			COALESCE(NULLIF(u.display_name, ''), SPLIT_PART(u.email, '@', 1), '') AS assigned_to_display_name,
			COALESCE(u.email, '') AS assigned_to_email,
			COALESCE(u.avatar_url, '') AS assigned_to_avatar_url`).
		Joins("LEFT JOIN users u ON u.id = tasks.assigned_to").
		Where("tasks.project_id = ?", projectID).
		Order("tasks.created_at desc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(rows))
	for i := range rows {
		tasks[i] = rows[i].Task
		tasks[i].AssignedToDisplayName = rows[i].DisplayName
		tasks[i].AssignedToEmail = rows[i].Email
		tasks[i].AssignedToAvatarURL = rows[i].AvatarURL
	}
	return tasks, nil
}

// projectPageTaskColumns is tasks.* minus heavy TEXT/JSONB columns not needed for board/backlog/overview lists.
const projectPageTaskColumns = `tasks.id, tasks.code, tasks.title, tasks.estimated_minutes, tasks.project_id, tasks.epic_id, tasks.sprint_id, tasks.milestone_id, tasks.task_type, tasks.priority, tasks.story_points, tasks.parent_id, tasks.sort_order, tasks.start_date, tasks.end_date, tasks.progress, tasks.due_at, tasks.started_at, tasks.completed_at, tasks.status, tasks.assigned_to, tasks.created_at, tasks.updated_at`

func (r *postgresRepository) GetTasksByProjectIDForProjectPage(projectID uuid.UUID, limit int) ([]domain.Task, error) {
	return r.GetTasksByProjectIDForProjectPageCursor(projectID, limit, nil, nil, 0)
}

func (r *postgresRepository) GetTasksByProjectIDForProjectPageCursor(projectID uuid.UUID, limit int, cursorCreatedAt *time.Time, cursorID *uuid.UUID, offset int) ([]domain.Task, error) {
	type taskRow struct {
		domain.Task
		DisplayName string `gorm:"column:assigned_to_display_name"`
	}
	const defaultLimit = 600
	if limit <= 0 {
		limit = defaultLimit
	}
	// Hard cap to protect DB/response size for very large projects.
	if limit > 5000 {
		limit = 5000
	}
	if offset < 0 {
		offset = 0
	}

	startedAt := time.Now()
	q := r.db.Table("tasks").
		Select(projectPageTaskColumns+`,
			COALESCE(NULLIF(u.display_name, ''), SPLIT_PART(u.email, '@', 1), '') AS assigned_to_display_name`).
		Joins("LEFT JOIN users u ON u.id = tasks.assigned_to").
		Where("tasks.project_id = ?", projectID)

	if cursorCreatedAt != nil && cursorID != nil {
		q = q.Where("(tasks.created_at < ?) OR (tasks.created_at = ? AND tasks.id < ?)", *cursorCreatedAt, *cursorCreatedAt, *cursorID)
	} else if offset > 0 {
		q = q.Offset(offset)
	}

	var rows []taskRow
	err := q.Order("tasks.created_at desc, tasks.id desc").
		Limit(limit + 1).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(rows))
	for i := range rows {
		tasks[i] = rows[i].Task
		tasks[i].AssignedToDisplayName = rows[i].DisplayName
	}
	fmt.Printf("[ProjectDetails] repo tasks query projectID=%s limit=%d rows=%d cursor=%t offset=%d elapsed=%s\n", projectID, limit, len(tasks), cursorCreatedAt != nil && cursorID != nil, offset, time.Since(startedAt))
	return tasks, nil
}

// DeleteProjectPlan removes all tasks, sprints, milestones, and epics for a project (for "clear plan" / reset before AI plan).
func (r *postgresRepository) DeleteProjectPlan(projectID uuid.UUID) error {
	// Remove task dependencies that reference any task in this project
	r.db.Exec("DELETE FROM task_dependencies WHERE predecessor_id IN (SELECT id FROM tasks WHERE project_id = ?) OR successor_id IN (SELECT id FROM tasks WHERE project_id = ?)", projectID, projectID)
	if err := r.db.Where("project_id = ?", projectID).Delete(&domain.Task{}).Error; err != nil {
		return err
	}
	if err := r.db.Where("project_id = ?", projectID).Delete(&domain.Sprint{}).Error; err != nil {
		return err
	}
	if err := r.db.Where("project_id = ?", projectID).Delete(&domain.Milestone{}).Error; err != nil {
		return err
	}
	if err := r.db.Where("project_id = ?", projectID).Delete(&domain.Epic{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) DeleteProject(id uuid.UUID) error {
	if err := r.db.Where("project_id = ?", id).Delete(&domain.Task{}).Error; err != nil {
		return err
	}
	result := r.db.Delete(&domain.Project{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

// --- Task Operations ---

func (r *postgresRepository) CreateTask(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *postgresRepository) GetTaskByID(id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, task_id, dev_id, reference_url, note, created_at").
			Order("created_at desc").Limit(20)
	}).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, status, priority, task_type, story_points, parent_id, project_id, assigned_to, sort_order")
	}).Preload("ParentTask", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, status, task_type, project_id, parent_id")
	}).First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *postgresRepository) GetTaskByCode(code string) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, task_id, dev_id, reference_url, note, created_at").
			Order("created_at desc").Limit(20)
	}).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, status, priority, task_type, story_points, parent_id, project_id, assigned_to, sort_order")
	}).Preload("ParentTask", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, status, task_type, project_id, parent_id")
	}).First(&task, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *postgresRepository) CountTasksForCode(projectID *uuid.UUID) (int, error) {
	var count int64
	q := r.db.Model(&domain.Task{})
	if projectID != nil {
		q = q.Where("project_id = ?", *projectID)
	} else {
		q = q.Where("project_id IS NULL")
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetMaxTaskCodeSuffix returns the maximum numeric suffix among tasks with code like "prefix-001", "prefix-002", ...
// so import can generate globally unique codes (task code has a global unique index).
func (r *postgresRepository) GetMaxTaskCodeSuffix(prefix string) (int, error) {
	if prefix == "" {
		return 0, nil
	}
	var maxSuffix int
	err := r.db.Raw(
		`SELECT COALESCE(MAX((regexp_match(code, '[0-9]+$'))[1]::int), 0) FROM tasks WHERE code ~ ?`,
		`^`+prefix+`-[0-9]+$`,
	).Scan(&maxSuffix).Error
	if err != nil {
		return 0, err
	}
	return maxSuffix, nil
}

func (r *postgresRepository) GetTasksByAssignee(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("assigned_to = ?", userID).
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func uuidSetToSlice(m map[uuid.UUID]struct{}) []uuid.UUID {
	if len(m) == 0 {
		return nil
	}
	out := make([]uuid.UUID, 0, len(m))
	for id := range m {
		out = append(out, id)
	}
	return out
}

// enrichMyBoardTasks fills project/sprint display fields for GET /tasks/my.
func (r *postgresRepository) enrichMyBoardTasks(tasks []domain.Task) {
	if len(tasks) == 0 {
		return
	}
	projIDSet := make(map[uuid.UUID]struct{})
	parentIDSet := make(map[uuid.UUID]struct{})
	for i := range tasks {
		if tasks[i].ProjectID != nil {
			projIDSet[*tasks[i].ProjectID] = struct{}{}
		}
		if tasks[i].SprintID == nil && tasks[i].ParentID != nil {
			parentIDSet[*tasks[i].ParentID] = struct{}{}
		}
	}
	var projects []domain.Project
	if ids := uuidSetToSlice(projIDSet); len(ids) > 0 {
		_ = r.db.Select("id", "name", "code", "color").Where("id IN ?", ids).Find(&projects).Error
	}
	projByID := make(map[uuid.UUID]domain.Project, len(projects))
	for _, p := range projects {
		projByID[p.ID] = p
	}
	parentByID := make(map[uuid.UUID]domain.Task)
	if pids := uuidSetToSlice(parentIDSet); len(pids) > 0 {
		var parents []domain.Task
		_ = r.db.Select("id", "sprint_id").Where("id IN ?", pids).Find(&parents).Error
		for _, p := range parents {
			parentByID[p.ID] = p
		}
	}
	sprintIDSet := make(map[uuid.UUID]struct{})
	for i := range tasks {
		var sid *uuid.UUID
		if tasks[i].SprintID != nil {
			sid = tasks[i].SprintID
		} else if tasks[i].ParentID != nil {
			if par, ok := parentByID[*tasks[i].ParentID]; ok && par.SprintID != nil {
				sid = par.SprintID
			}
		}
		if sid != nil {
			sprintIDSet[*sid] = struct{}{}
		}
	}
	var sprints []domain.Sprint
	if sids := uuidSetToSlice(sprintIDSet); len(sids) > 0 {
		_ = r.db.Select("id", "name").Where("id IN ?", sids).Find(&sprints).Error
	}
	sprintByID := make(map[uuid.UUID]domain.Sprint, len(sprints))
	for _, s := range sprints {
		sprintByID[s.ID] = s
	}
	for i := range tasks {
		if tasks[i].ProjectID != nil {
			if p, ok := projByID[*tasks[i].ProjectID]; ok {
				tasks[i].ProjectName = p.Name
				tasks[i].ProjectColor = p.Color
			}
		}
		var eff *uuid.UUID
		if tasks[i].SprintID != nil {
			eff = tasks[i].SprintID
		} else if tasks[i].ParentID != nil {
			if par, ok := parentByID[*tasks[i].ParentID]; ok && par.SprintID != nil {
				eff = par.SprintID
			}
		}
		if eff != nil {
			tasks[i].EffectiveSprintID = eff
			if s, ok := sprintByID[*eff]; ok {
				tasks[i].SprintName = s.Name
			}
		}
	}
}

// GetActiveSprintTasksByAssignee returns only tasks that belong to an ACTIVE sprint
// and are assigned to the given user. DEV role sees strictly their active battlefield.
// Also includes child tasks whose parent (FEATURE) is in an active sprint.
func (r *postgresRepository) GetActiveSprintTasksByAssignee(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.
		Where("tasks.assigned_to = ?", userID).
		Where(`
			EXISTS (
				SELECT 1 FROM sprints s
				WHERE s.id = tasks.sprint_id AND s.status = 'ACTIVE'
			)
			OR EXISTS (
				SELECT 1 FROM tasks parent
				JOIN sprints s ON s.id = parent.sprint_id
				WHERE parent.id = tasks.parent_id AND s.status = 'ACTIVE'
			)
		`).
		Order("tasks.created_at desc").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	r.enrichMyBoardTasks(tasks)
	return tasks, nil
}

// GetActiveSprintsForUser returns distinct ACTIVE sprints that have tasks assigned to the given user.
func (r *postgresRepository) GetActiveSprintsForUser(userID uint) ([]domain.Sprint, error) {
	var sprints []domain.Sprint
	err := r.db.
		Distinct("sprints.*").
		Table("sprints").
		Joins("JOIN tasks ON tasks.sprint_id = sprints.id").
		Where("sprints.status = ? AND tasks.assigned_to = ?", "ACTIVE", userID).
		Order("sprints.created_at desc").
		Find(&sprints).Error
	if err != nil {
		return nil, err
	}
	projIDSet := make(map[uuid.UUID]struct{})
	for i := range sprints {
		projIDSet[sprints[i].ProjectID] = struct{}{}
	}
	if ids := uuidSetToSlice(projIDSet); len(ids) > 0 {
		var projs []domain.Project
		if qerr := r.db.Select("id", "name", "code").Where("id IN ?", ids).Find(&projs).Error; qerr == nil {
			pmap := make(map[uuid.UUID]domain.Project, len(projs))
			for _, p := range projs {
				pmap[p.ID] = p
			}
			for i := range sprints {
				if p, ok := pmap[sprints[i].ProjectID]; ok {
					sprints[i].ProjectName = p.Name
					sprints[i].ProjectCode = p.Code
				}
			}
		}
	}
	return sprints, nil
}

// GetGlobalActiveTasks returns TASK and BUG items scoped to ACTIVE sprints only.
// A task is included when:
//   a) it is directly assigned to an ACTIVE sprint, OR
//   b) its parent task (FEATURE or TASK) is assigned to an ACTIVE sprint (sub-task inclusion rule).
// FEATURE types are excluded from results — they belong to the Product Owner/CEO Feature Roadmap Board.
//
// CEO/MANAGER: all projects (company-wide). Teams off: Product Owner → project_pm_assignments; engineer → assigned-task projects.
// Otherwise: projects whose team_id matches the caller's user.team_id.
func (r *postgresRepository) GetGlobalActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	var results []domain.GlobalActiveTask
	q := r.db.Table("tasks").
		Select("tasks.id, tasks.code, tasks.title, tasks.description, tasks.status, tasks.priority, tasks.task_type, tasks.story_points, tasks.estimated_minutes, tasks.assigned_to, tasks.project_id, tasks.sprint_id, tasks.epic_id, tasks.parent_id, tasks.due_at, tasks.started_at, tasks.created_at, tasks.updated_at, projects.name AS project_name, projects.color AS project_color").
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Where("tasks.task_type IN ? AND tasks.status NOT IN ?",
			[]string{"TASK", "BUG"},
			[]string{"COMPLETED", "CANCELLED"},
		).
		Where(`
			EXISTS (
				SELECT 1 FROM sprints s
				WHERE s.id = tasks.sprint_id AND s.status = 'ACTIVE'
			)
			OR EXISTS (
				SELECT 1 FROM tasks parent
				JOIN sprints s ON s.id = parent.sprint_id
				WHERE parent.id = tasks.parent_id AND s.status = 'ACTIVE'
			)
		`)

	switch {
	case ctx.Role == domain.RoleCEO || ctx.Role == domain.RoleManager:
		// no project scope — entire company
	case ctx.TeamsFeatureDisabled && ctx.UserID != 0:
		switch {
		case ctx.Role == domain.RoleProductOwner:
			q = q.Where("projects.id IN (SELECT project_id FROM project_pm_assignments WHERE user_id = ?)", ctx.UserID)
		case domain.IsEngineerRole(ctx.Role):
			q = q.Where("projects.id IN (SELECT DISTINCT project_id FROM tasks WHERE assigned_to = ?)", ctx.UserID)
		default:
			q = q.Joins("JOIN users ON users.id = ? AND users.team_id = projects.team_id", ctx.UserID)
		}
	default:
		q = q.Joins("JOIN users ON users.id = ? AND users.team_id = projects.team_id", ctx.UserID)
	}

	err := q.Order("tasks.created_at DESC").Scan(&results).Error
	return results, err
}

// GetTeamActiveTasks returns TASK and BUG items in ACTIVE sprints visible to caller.
// Includes direct sprint tasks AND sub-tasks whose parent FEATURE is in an ACTIVE sprint.
// Excludes COMPLETED and CANCELLED tasks.
// Joins users to populate AssignedToDisplayName / AssignedToEmail for Quick Log Time UI.
// Teams on: non-CEO/MANAGER are scoped by team_id.
// Teams off: Product Owner sees PM-assigned projects; engineers see projects where they have assignments.
func (r *postgresRepository) GetTeamActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	type row struct {
		domain.Task
		ProjectName           string `gorm:"column:project_name"`
		ProjectColor          string `gorm:"column:project_color"`
		AssignedToDisplayName string `gorm:"column:assigned_to_display_name"`
		AssignedToEmail       string `gorm:"column:assigned_to_email"`
	}
	var rows []row
	q := r.db.Table("tasks").
		Select(`tasks.*, projects.name AS project_name, projects.color AS project_color,
			COALESCE(NULLIF(u.display_name, ''), SPLIT_PART(u.email, '@', 1), '') AS assigned_to_display_name,
			COALESCE(u.email, '')               AS assigned_to_email`).
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("LEFT JOIN users u ON u.id = tasks.assigned_to").
		Where("tasks.task_type IN ? AND tasks.status NOT IN ?",
			[]string{"TASK", "BUG"},
			[]string{"COMPLETED", "CANCELLED"},
		).
		Where(`
			EXISTS (SELECT 1 FROM sprints s WHERE s.id = tasks.sprint_id AND s.status = 'ACTIVE')
			OR EXISTS (
				SELECT 1 FROM tasks parent
				JOIN sprints s ON s.id = parent.sprint_id
				WHERE parent.id = tasks.parent_id AND s.status = 'ACTIVE'
			)
		`).
		Order("tasks.created_at DESC")

	switch {
	case ctx.Role == domain.RoleCEO || ctx.Role == domain.RoleManager:
		// company-wide
	case ctx.TeamsFeatureDisabled && ctx.UserID != 0:
		switch {
		case ctx.Role == domain.RoleProductOwner:
			q = q.Where("projects.id IN (SELECT project_id FROM project_pm_assignments WHERE user_id = ?)", ctx.UserID)
		case domain.IsEngineerRole(ctx.Role):
			q = q.Where("projects.id IN (SELECT DISTINCT project_id FROM tasks WHERE assigned_to = ?)", ctx.UserID)
		default:
			q = q.Joins("JOIN users ON users.id = ? AND users.team_id = projects.team_id", ctx.UserID)
		}
	default:
		q = q.Joins("JOIN users ON users.id = ? AND users.team_id = projects.team_id", ctx.UserID)
	}
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]domain.GlobalActiveTask, len(rows))
	for i, r := range rows {
		results[i] = domain.GlobalActiveTask{
			Task:         r.Task,
			ProjectName:  r.ProjectName,
			ProjectColor: r.ProjectColor,
		}
		results[i].AssignedToDisplayName = r.AssignedToDisplayName
		results[i].AssignedToEmail = r.AssignedToEmail
	}
	return results, nil
}

// GetActiveFeatures returns FEATURE-type tasks for roadmap view.
// teamID=0 means no team filter (CEO/MANAGER). Optional projectID scopes to one project.
func (r *postgresRepository) GetActiveFeatures(teamID uint, projectID *uuid.UUID) ([]domain.FeatureRoadmapItem, error) {
	// Step 1: fetch all FEATURE tasks with project details
	type featureRow struct {
		domain.Task
		ProjectName  string `gorm:"column:project_name"`
		ProjectColor string `gorm:"column:project_color"`
		ProjectCode  string `gorm:"column:project_code"`
	}
	var features []featureRow
	q := r.db.Table("tasks").
		Select("tasks.*, projects.name AS project_name, projects.color AS project_color, projects.code AS project_code").
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Where("tasks.task_type = ?", "FEATURE").
		Order("tasks.created_at DESC")
	if teamID != 0 {
		q = q.Where("projects.team_id = ?", teamID)
	}
	if projectID != nil {
		q = q.Where("tasks.project_id = ?", *projectID)
	}
	if err := q.Scan(&features).Error; err != nil {
		return nil, err
	}
	if len(features) == 0 {
		return []domain.FeatureRoadmapItem{}, nil
	}

	// Step 2: collect feature IDs, then bulk-fetch child tasks
	featureIDs := make([]uuid.UUID, 0, len(features))
	for _, f := range features {
		featureIDs = append(featureIDs, f.ID)
	}
	var children []domain.Task
	if err := r.db.Table("tasks").
		Where("parent_id IN ? AND task_type IN ?", featureIDs, []string{"TASK", "BUG"}).
		Order("created_at ASC").
		Scan(&children).Error; err != nil {
		return nil, err
	}

	// Step 3: index children by parent_id
	childMap := make(map[string][]domain.Task)
	for _, c := range children {
		if c.ParentID != nil {
			pid := c.ParentID.String()
			childMap[pid] = append(childMap[pid], c)
		}
	}

	// Step 4: build result items with roll-up progress
	results := make([]domain.FeatureRoadmapItem, 0, len(features))
	for _, f := range features {
		kids := childMap[f.ID.String()]
		var rollup int
		if len(kids) > 0 {
			completed := 0
			for _, k := range kids {
				if k.Status == "COMPLETED" {
					completed++
				}
			}
			rollup = int(float64(completed) / float64(len(kids)) * 100)
		}
		results = append(results, domain.FeatureRoadmapItem{
			Task:           f.Task,
			ProjectName:    f.ProjectName,
			ProjectColor:   f.ProjectColor,
			ProjectCode:    f.ProjectCode,
			RollupProgress: rollup,
			ChildTasks:     kids,
		})
	}
	return results, nil
}

func (r *postgresRepository) GetUnassignedTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("assigned_to IS NULL").
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

// maxAllTasksLimit caps GET /tasks when no project_id (CEO / Product Owner dashboard) to avoid slow full scan
const maxAllTasksLimit = 2000

func (r *postgresRepository) GetAllTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Order("created_at desc").
		Limit(maxAllTasksLimit).
		Find(&tasks).Error
	return tasks, err
}

func (r *postgresRepository) GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]domain.Task, error) {
	if len(projectIDs) == 0 {
		return []domain.Task{}, nil
	}
	var tasks []domain.Task
	err := r.db.Where("project_id IN ?", projectIDs).
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

// GetTasksRequiringApproval returns tasks that need Product Owner/CEO attention:
// - Tasks with status REVIEW_PENDING (awaiting handover approval)
// - Tasks with negotiation_status = 'PENDING' (engineer wants more time)
func (r *postgresRepository) GetTasksRequiringApproval() ([]domain.Task, error) {
	var tasks []domain.Task

	err := r.db.
		Preload("Submissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, task_id, dev_id, reference_url, note, created_at").
				Order("created_at desc")
		}).
		Where("status = ? OR negotiation_status = ?", "REVIEW_PENDING", "PENDING").
		Order("created_at desc").
		Find(&tasks).Error

	return tasks, err
}

func (r *postgresRepository) UpdateTask(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *postgresRepository) DeleteTask(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete in dependency order to avoid FK violations (works even if DB has no ON DELETE CASCADE)
		// 1. Appeals reference submissions; delete appeals for submissions of this task
		if err := tx.Exec("DELETE FROM appeals WHERE submission_id IN (SELECT id FROM submissions WHERE task_id = ?)", id).Error; err != nil {
			return err
		}
		// 2. Submissions reference tasks
		if err := tx.Exec("DELETE FROM submissions WHERE task_id = ?", id).Error; err != nil {
			return err
		}
		// 3. Task dependencies reference tasks (predecessor_id, successor_id)
		if err := tx.Exec("DELETE FROM task_dependencies WHERE predecessor_id = ? OR successor_id = ?", id, id).Error; err != nil {
			return err
		}
		// 4. Task activity audit trail
		if err := tx.Exec("DELETE FROM task_activity_events WHERE task_id = ?", id).Error; err != nil {
			return err
		}
		// 5. Delete the task
		return tx.Delete(&domain.Task{}, "id = ?", id).Error
	})
}

func (r *postgresRepository) CreateTaskActivity(e *domain.TaskActivityEvent) error {
	return r.db.Create(e).Error
}

func (r *postgresRepository) ListTaskActivitiesByTaskID(taskID uuid.UUID) ([]domain.TaskActivityEvent, error) {
	var rows []domain.TaskActivityEvent
	err := r.db.Where("task_id = ?", taskID).Order("created_at asc, id asc").Find(&rows).Error
	return rows, err
}

// GetAllTaskDependencies returns all task dependencies for Gantt chart rendering
func (r *postgresRepository) GetAllTaskDependencies() ([]domain.TaskDependency, error) {
	var deps []domain.TaskDependency
	err := r.db.Order("created_at asc").Find(&deps).Error
	return deps, err
}

// CreateTaskDependency persists a new task dependency (Gantt link)
func (r *postgresRepository) CreateTaskDependency(dep *domain.TaskDependency) error {
	return r.db.Create(dep).Error
}

// DeleteTaskDependency removes a task dependency by ID
func (r *postgresRepository) DeleteTaskDependency(id uuid.UUID) error {
	return r.db.Delete(&domain.TaskDependency{}, "id = ?", id).Error
}

// ApproveTask marks a task as COMPLETED and sets CompletedAt timestamp
func (r *postgresRepository) ApproveTask(id uuid.UUID) error {
	// Use SQL UPDATE with NOW() to ensure atomic operation
	result := r.db.Exec(`
		UPDATE tasks 
		SET status = 'COMPLETED', 
		    completed_at = NOW(),
		    updated_at = NOW()
		WHERE id = ?
	`, id)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	
	return nil
}

// RejectTask returns a task to IN_PROGRESS and logs the rejection reason as a comment
func (r *postgresRepository) RejectTask(taskID uuid.UUID, rejectorID uint, reason string) error {
	result := r.db.Exec(`
		UPDATE tasks
		SET status = 'IN_PROGRESS',
		    updated_at = NOW()
		WHERE id = ?
	`, taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	comment := &domain.TaskComment{
		TaskID:  taskID,
		UserID:  rejectorID,
		Content: fmt.Sprintf("[REJECTED] %s", reason),
	}
	return r.db.Create(comment).Error
}

// --- Submission Operations ---

func (r *postgresRepository) CreateSubmission(sub *domain.Submission) error {
	return r.db.Create(sub).Error
}

func (r *postgresRepository) GetSubmissionByID(id uuid.UUID) (*domain.Submission, error) {
	var sub domain.Submission
	err := r.db.Preload("Appeal").First(&sub, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *postgresRepository) UpdateSubmission(sub *domain.Submission) error {
	return r.db.Save(sub).Error
}

func (r *postgresRepository) GetLatestSubmission(taskID uuid.UUID) (*domain.Submission, error) {
	var sub domain.Submission
	// Get the most recent submission for a task
	err := r.db.Where("task_id = ?", taskID).
		Order("created_at desc").
		First(&sub).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil if no submissions yet
	}
	return &sub, err
}

// --- Appeal Operations ---

func (r *postgresRepository) CreateAppeal(appeal *domain.Appeal) error {
	return r.db.Create(appeal).Error
}

func (r *postgresRepository) GetAppealBySubmissionID(subID uuid.UUID) (*domain.Appeal, error) {
	var appeal domain.Appeal
	err := r.db.Where("submission_id = ?", subID).First(&appeal).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // No appeal exists
	}
	return &appeal, err
}

func (r *postgresRepository) GetAppealByID(id uuid.UUID) (*domain.Appeal, error) {
	var appeal domain.Appeal
	err := r.db.First(&appeal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &appeal, nil
}

func (r *postgresRepository) UpdateAppeal(appeal *domain.Appeal) error {
	return r.db.Save(appeal).Error
}

// --- System Configuration Operations (Singleton Pattern) ---

// GetSystemConfig fetches the single system configuration record (ID=1)
// If not exists, creates a default configuration
func (r *postgresRepository) GetSystemConfig() (*domain.SystemConfig, error) {
	var config domain.SystemConfig
	
	err := r.db.First(&config, "id = ?", 1).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create default config if not exists
		config = domain.SystemConfig{
			ID:               1,
			ActiveModel:      "gemini-2.5-flash-lite",
			Temperature:      0.4,
			CursorAssistance: 80,
		}
		
		if createErr := r.db.Create(&config).Error; createErr != nil {
			return nil, createErr
		}
		return &config, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

// UpdateSystemConfig updates the system configuration (always updates ID=1)
func (r *postgresRepository) UpdateSystemConfig(config *domain.SystemConfig) error {
	config.ID = 1
	return r.db.Save(config).Error
}

// --- Sprint Operations ---

func (r *postgresRepository) CreateSprint(sprint *domain.Sprint) error {
	return r.db.Create(sprint).Error
}

func (r *postgresRepository) GetSprintByID(id uuid.UUID) (*domain.Sprint, error) {
	var sprint domain.Sprint
	err := r.db.First(&sprint, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (r *postgresRepository) GetSprintsByProjectID(projectID uuid.UUID) ([]domain.Sprint, error) {
	var sprints []domain.Sprint
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc, created_at asc").Find(&sprints).Error
	return sprints, err
}

func (r *postgresRepository) GetActiveSprintByProjectID(projectID uuid.UUID) (*domain.Sprint, error) {
	var sprint domain.Sprint
	err := r.db.Where("project_id = ? AND status = ?", projectID, "ACTIVE").First(&sprint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sprint, err
}

func (r *postgresRepository) UpdateSprint(sprint *domain.Sprint) error {
	return r.db.Save(sprint).Error
}

func (r *postgresRepository) DeleteSprint(id uuid.UUID) error {
	// Unlink tasks before deleting sprint
	r.db.Model(&domain.Task{}).Where("sprint_id = ?", id).Update("sprint_id", nil)
	return r.db.Delete(&domain.Sprint{}, "id = ?", id).Error
}

// --- Milestone Operations ---

func (r *postgresRepository) CreateMilestone(m *domain.Milestone) error {
	return r.db.Create(m).Error
}

func (r *postgresRepository) GetMilestoneByID(id uuid.UUID) (*domain.Milestone, error) {
	var m domain.Milestone
	err := r.db.First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *postgresRepository) GetMilestonesByProjectID(projectID uuid.UUID) ([]domain.Milestone, error) {
	var milestones []domain.Milestone
	err := r.db.Where("project_id = ?", projectID).Order("due_date asc").Find(&milestones).Error
	return milestones, err
}

func (r *postgresRepository) UpdateMilestone(m *domain.Milestone) error {
	return r.db.Save(m).Error
}

func (r *postgresRepository) DeleteMilestone(id uuid.UUID) error {
	return r.db.Delete(&domain.Milestone{}, "id = ?", id).Error
}

// --- Task Comment Operations ---

func (r *postgresRepository) CreateTaskComment(c *domain.TaskComment) error {
	return r.db.Create(c).Error
}

func (r *postgresRepository) GetCommentsByTaskID(taskID uuid.UUID) ([]domain.TaskComment, error) {
	var comments []domain.TaskComment
	err := r.db.Where("task_id = ?", taskID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *postgresRepository) GetTaskCommentByID(commentID uuid.UUID) (*domain.TaskComment, error) {
	var c domain.TaskComment
	err := r.db.First(&c, "id = ?", commentID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *postgresRepository) UpdateTaskComment(c *domain.TaskComment) error {
	return r.db.Save(c).Error
}

// --- Time Log Operations ---

func (r *postgresRepository) CreateTimeLog(t *domain.TimeLog) error {
	return r.db.Create(t).Error
}

func (r *postgresRepository) GetTimeLogsByTaskID(taskID uuid.UUID) ([]domain.TimeLog, error) {
	var logs []domain.TimeLog
	err := r.db.Raw(`
		SELECT tl.id, tl.task_id, tl.user_id, u.email AS user_email,
		       tl.minutes, tl.description, tl.work_type, tl.logged_date, tl.is_timer_session, tl.logged_at
		FROM time_logs tl
		JOIN users u ON u.id = tl.user_id
		WHERE tl.task_id = ?
		ORDER BY tl.logged_date DESC, tl.logged_at DESC
	`, taskID).Scan(&logs).Error
	return logs, err
}

func (r *postgresRepository) GetTimeLogByID(logID uuid.UUID) (*domain.TimeLog, error) {
	var t domain.TimeLog
	err := r.db.Raw(`
		SELECT tl.id, tl.task_id, tl.user_id, u.email AS user_email,
		       tl.minutes, tl.description, tl.work_type, tl.logged_date, tl.is_timer_session, tl.logged_at
		FROM time_logs tl
		JOIN users u ON u.id = tl.user_id
		WHERE tl.id = ?
	`, logID).Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.ID == uuid.Nil {
		return nil, nil
	}
	return &t, nil
}

func (r *postgresRepository) UpdateTimeLog(t *domain.TimeLog) error {
	return r.db.Model(&domain.TimeLog{}).Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"task_id":     t.TaskID,
			"minutes":     t.Minutes,
			"description": t.Description,
			"work_type":   t.WorkType,
		}).Error
}

func (r *postgresRepository) DeleteTimeLog(logID uuid.UUID) error {
	return r.db.Where("id = ?", logID).Delete(&domain.TimeLog{}).Error
}

func (r *postgresRepository) GetTimeLogsByUserAndDate(userID uint, date time.Time) ([]domain.TimeLog, error) {
	var logs []domain.TimeLog
	err := r.db.Raw(`
		SELECT tl.id, tl.task_id, tl.user_id, u.email AS user_email,
		       COALESCE(t.code, '') AS task_code,
		       COALESCE(t.title, '') AS task_title,
		       tl.minutes, tl.description, tl.work_type, tl.logged_date, tl.is_timer_session, tl.logged_at
		FROM time_logs tl
		JOIN users u ON u.id = tl.user_id
		LEFT JOIN tasks t ON t.id = tl.task_id
		WHERE tl.user_id = ? AND tl.logged_date = ?
		ORDER BY tl.logged_at DESC
	`, userID, date.Format("2006-01-02")).Scan(&logs).Error
	return logs, err
}

func (r *postgresRepository) BulkCreateTimeLogs(logs []domain.TimeLog) error {
	return r.db.Create(&logs).Error
}

func (r *postgresRepository) GetTotalLoggedMinutes(taskID uuid.UUID) (int, error) {
	var total int64
	err := r.db.Model(&domain.TimeLog{}).Where("task_id = ?", taskID).Select("COALESCE(SUM(minutes), 0)").Scan(&total).Error
	return int(total), err
}

func (r *postgresRepository) CountChildTasks(parentID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&domain.Task{}).Where("parent_id = ?", parentID).Count(&count).Error
	return int(count), err
}

func (r *postgresRepository) GetChildTasksByParentID(parentID uuid.UUID) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("parent_id = ?", parentID).Find(&tasks).Error
	return tasks, err
}

func (r *postgresRepository) GetTasksReadyForTest(teamID uint) ([]domain.GlobalActiveTask, error) {
	var results []domain.GlobalActiveTask
	q := r.db.Table("tasks").
		Select(`tasks.*,
			projects.name AS project_name,
			projects.color AS project_color,
			COALESCE(NULLIF(u.display_name, ''), SPLIT_PART(u.email, '@', 1), '') AS assigned_to_display_name,
			COALESCE(u.email, '') AS assigned_to_email`).
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("LEFT JOIN users u ON u.id = tasks.assigned_to").
		Where("tasks.status = ? AND tasks.task_type IN ?", "READY_FOR_TEST", []string{"TASK", "BUG"}).
		Order("tasks.created_at DESC")
	if teamID != 0 {
		q = q.Where("projects.team_id = ?", teamID)
	}
	err := q.Scan(&results).Error
	return results, err
}

// SetTaskReadyForUAT transitions a task to READY_FOR_UAT and saves the Product Owner test evidence payload.
func (r *postgresRepository) SetTaskReadyForUAT(taskID uuid.UUID, uatPayload []byte) error {
	result := r.db.Exec(`
		UPDATE tasks
		SET status = 'READY_FOR_UAT',
		    uat_payload = ?,
		    updated_at = NOW()
		WHERE id = ?
	`, uatPayload, taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

// SetTaskWaitForDeploy transitions a task to WAIT_FOR_DEPLOY (pending Chief Engineer deployment).
// The Product Owner test evidence payload is saved alongside the status change.
func (r *postgresRepository) SetTaskWaitForDeploy(taskID uuid.UUID, uatPayload []byte) error {
	result := r.db.Exec(`
		UPDATE tasks
		SET status = 'WAIT_FOR_DEPLOY',
		    uat_payload = ?,
		    updated_at = NOW()
		WHERE id = ?
	`, uatPayload, taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

// AdvanceTaskToReadyForUAT transitions a task from WAIT_FOR_DEPLOY → READY_FOR_UAT.
// Called automatically when the Chief Engineer marks a linked deployment request as deployed.
func (r *postgresRepository) AdvanceTaskToReadyForUAT(taskID uuid.UUID) error {
	result := r.db.Exec(`
		UPDATE tasks
		SET status = 'READY_FOR_UAT',
		    updated_at = NOW()
		WHERE id = ? AND status = 'WAIT_FOR_DEPLOY'
	`, taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or not in WAIT_FOR_DEPLOY status")
	}
	return nil
}

// GetTasksReadyForCEOApproval returns TASK/BUG items in READY_FOR_UAT status awaiting CEO final approval.
func (r *postgresRepository) GetTasksReadyForCEOApproval(teamID uint) ([]domain.GlobalActiveTask, error) {
	var results []domain.GlobalActiveTask
	q := r.db.Table("tasks").
		Select(`tasks.*,
			projects.name AS project_name,
			projects.color AS project_color,
			COALESCE(NULLIF(u.display_name, ''), SPLIT_PART(u.email, '@', 1), '') AS assigned_to_display_name,
			COALESCE(u.email, '') AS assigned_to_email`).
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Joins("LEFT JOIN users u ON u.id = tasks.assigned_to").
		Where("tasks.status = ? AND tasks.task_type IN ?", "READY_FOR_UAT", []string{"TASK", "BUG"}).
		Order("tasks.created_at DESC")
	if teamID != 0 {
		q = q.Where("projects.team_id = ?", teamID)
	}
	err := q.Scan(&results).Error
	return results, err
}

// --- Bulk Operations ---

func (r *postgresRepository) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error {
	if len(taskIDs) == 0 {
		return nil
	}
	updates := map[string]interface{}{"status": status}
	if status == "IN_PROGRESS" {
		now := time.Now()
		updates["started_at"] = now
	}
	return r.db.Model(&domain.Task{}).Where("id IN ?", taskIDs).Updates(updates).Error
}

// --- Analytics ---

func (r *postgresRepository) GetProjectAnalytics(projectID uuid.UUID) (*domain.ProjectAnalytics, error) {
	analytics := &domain.ProjectAnalytics{ProjectID: projectID}

	type taskStats struct {
		Total     int
		Completed int
	}
	var ts taskStats
	r.db.Raw(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'COMPLETED' THEN 1 ELSE 0 END) as completed
		FROM tasks WHERE project_id = ? AND parent_id IS NULL
	`, projectID).Scan(&ts)
	analytics.TotalTasks = ts.Total
	analytics.CompletedTasks = ts.Completed

	// Story points
	type spStats struct {
		Total     int
		Completed int
	}
	var sp spStats
	r.db.Raw(`
		SELECT 
			COALESCE(SUM(story_points), 0) as total,
			COALESCE(SUM(CASE WHEN status = 'COMPLETED' THEN story_points ELSE 0 END), 0) as completed
		FROM tasks WHERE project_id = ? AND parent_id IS NULL
	`, projectID).Scan(&sp)
	analytics.TotalStoryPoints = sp.Total
	analytics.CompletedSP = sp.Completed

	// Total logged time
	var totalLogged struct{ Total int }
	r.db.Raw(`
		SELECT COALESCE(SUM(tl.minutes), 0) as total
		FROM time_logs tl
		JOIN tasks t ON t.id = tl.task_id
		WHERE t.project_id = ?
	`, projectID).Scan(&totalLogged)
	analytics.TotalLoggedMinutes = totalLogged.Total

	// Average cycle time (days from started_at to completed_at)
	var avgCycle struct{ Avg float64 }
	r.db.Raw(`
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - started_at)) / 86400), 0) as avg
		FROM tasks
		WHERE project_id = ? AND status = 'COMPLETED' AND started_at IS NOT NULL AND completed_at IS NOT NULL
	`, projectID).Scan(&avgCycle)
	analytics.AvgCycleTimeDays = avgCycle.Avg

	// Burndown: for active sprint
	activeSprint, _ := r.GetActiveSprintByProjectID(projectID)
	if activeSprint != nil && activeSprint.StartDate != nil && activeSprint.EndDate != nil {
		type burnRow struct {
			Day       time.Time
			Remaining float64
		}
		var burnRows []burnRow
		r.db.Raw(`
			SELECT 
				gs.day::date as day,
				COALESCE(SUM(CASE WHEN t.status != 'COMPLETED' OR t.completed_at::date > gs.day THEN t.story_points ELSE 0 END), 0) as remaining
			FROM generate_series(?::date, ?::date, '1 day'::interval) gs(day)
			LEFT JOIN tasks t ON t.sprint_id = ? AND t.project_id = ?
			GROUP BY gs.day
			ORDER BY gs.day
		`, activeSprint.StartDate, activeSprint.EndDate, activeSprint.ID, projectID).Scan(&burnRows)

		totalSP := float64(sp.Total)
		sprintDays := activeSprint.EndDate.Sub(*activeSprint.StartDate).Hours() / 24
		for i, row := range burnRows {
			ideal := totalSP - (totalSP * float64(i) / sprintDays)
			if ideal < 0 {
				ideal = 0
			}
			analytics.Burndown = append(analytics.Burndown, domain.BurndownPoint{
				Day:       row.Day.Format("2006-01-02"),
				Ideal:     ideal,
				Remaining: row.Remaining,
			})
		}
	}

	// Velocity: completed story points per sprint (last 6)
	type velocityRow struct {
		SprintName  string
		CompletedSP int
		PlannedSP   int
	}
	var vRows []velocityRow
	r.db.Raw(`
		SELECT 
			s.name as sprint_name,
			COALESCE(SUM(CASE WHEN t.status = 'COMPLETED' THEN t.story_points ELSE 0 END), 0) as completed_sp,
			COALESCE(SUM(t.story_points), 0) as planned_sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id
		WHERE s.project_id = ? AND s.status IN ('ACTIVE', 'COMPLETED')
		GROUP BY s.id, s.name, s.created_at
		ORDER BY s.created_at DESC
		LIMIT 6
	`, projectID).Scan(&vRows)

	for i := len(vRows) - 1; i >= 0; i-- {
		analytics.Velocity = append(analytics.Velocity, domain.VelocityPoint{
			SprintName:  vRows[i].SprintName,
			CompletedSP: vRows[i].CompletedSP,
			PlannedSP:   vRows[i].PlannedSP,
		})
	}

	// Team capacity (join users for real email / display name)
	type capacityRow struct {
		UserID          uint
		UserEmail       string
		UserDisplayName string
		AssignedTasks   int
		EstimatedMins   int
		LoggedMins      int
	}
	var capRows []capacityRow
	r.db.Raw(`
		SELECT 
			t.assigned_to as user_id,
			COALESCE(u.email, '') as user_email,
			COALESCE(u.display_name, '') as user_display_name,
			COUNT(t.id) as assigned_tasks,
			COALESCE(SUM(t.estimated_minutes), 0) as estimated_mins,
			COALESCE(SUM(tl_sum.total_logged), 0) as logged_mins
		FROM tasks t
		LEFT JOIN users u ON u.id = t.assigned_to
		LEFT JOIN (
			SELECT task_id, SUM(minutes) as total_logged FROM time_logs GROUP BY task_id
		) tl_sum ON tl_sum.task_id = t.id
		WHERE t.project_id = ? AND t.assigned_to IS NOT NULL
		GROUP BY t.assigned_to, u.email, u.display_name
	`, projectID).Scan(&capRows)

	for _, row := range capRows {
		util := 0.0
		if row.EstimatedMins > 0 {
			util = float64(row.LoggedMins) / float64(row.EstimatedMins) * 100
		}
		analytics.TeamCapacity = append(analytics.TeamCapacity, domain.TeamCapacityRow{
			UserID:          row.UserID,
			UserEmail:       row.UserEmail,
			UserDisplayName: strings.TrimSpace(row.UserDisplayName),
			AssignedTasks:   row.AssignedTasks,
			EstimatedHours:  float64(row.EstimatedMins) / 60,
			LoggedHours:     float64(row.LoggedMins) / 60,
			Utilization:     util,
		})
	}

	return analytics, nil
}

// --- Epic Operations ---

func (r *postgresRepository) CreateEpic(epic *domain.Epic) error {
	return r.db.Create(epic).Error
}

func (r *postgresRepository) GetEpicByID(id uuid.UUID) (*domain.Epic, error) {
	var epic domain.Epic
	err := r.db.First(&epic, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &epic, nil
}

func (r *postgresRepository) GetEpicsByProjectID(projectID uuid.UUID) ([]domain.Epic, error) {
	var epics []domain.Epic
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc, created_at asc").Find(&epics).Error
	return epics, err
}

func (r *postgresRepository) UpdateEpic(epic *domain.Epic) error {
	return r.db.Save(epic).Error
}

func (r *postgresRepository) DeleteEpic(id uuid.UUID) error {
	// Unlink tasks before deleting epic
	r.db.Model(&domain.Task{}).Where("epic_id = ?", id).Update("epic_id", nil)
	return r.db.Delete(&domain.Epic{}, "id = ?", id).Error
}

// GetEpicTimelineData returns all epics for a project with their tasks preloaded (including sub_tasks for epic bar span), ordered by start_date.
func (r *postgresRepository) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	var epics []domain.Epic
	err := r.db.
		Select("id", "project_id", "title", "description", "status", "color", "sort_order", "start_date", "end_date", "created_at", "updated_at").
		Where("project_id = ?", projectID).
		Order("sort_order asc, start_date asc NULLS LAST, created_at asc").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("id", "code", "title", "project_id", "epic_id", "milestone_id", "task_type", "priority", "story_points", "parent_id", "sort_order", "start_date", "end_date", "progress", "due_at", "status", "assigned_to", "created_at", "updated_at").
				Where("parent_id IS NULL").
				Order("start_date asc NULLS LAST, created_at asc").
				Preload("SubTasks", func(subDB *gorm.DB) *gorm.DB {
					return subDB.Select("id", "title", "project_id", "epic_id", "parent_id", "task_type", "priority", "story_points", "sort_order", "start_date", "end_date", "progress", "due_at", "status", "assigned_to", "created_at", "updated_at").Order("start_date asc NULLS LAST, created_at asc")
				})
		}).
		Find(&epics).Error
	if err != nil {
		return nil, err
	}
	return &domain.EpicTimelineData{Epics: epics}, nil
}

// GetSprintTimelineData returns all sprints for a project with their tasks preloaded, ordered by start_date.
func (r *postgresRepository) GetSprintTimelineData(projectID uuid.UUID) (*domain.SprintTimelineData, error) {
	var sprints []domain.Sprint
	err := r.db.Where("project_id = ?", projectID).
		Order("start_date asc NULLS LAST, created_at asc").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").Order("start_date asc NULLS LAST, created_at asc")
		}).
		Find(&sprints).Error
	if err != nil {
		return nil, err
	}
	return &domain.SprintTimelineData{Sprints: sprints}, nil
}

// GetImportedSlideIndicesByPresentationID returns 1-based slide indices of tasks already imported from this presentation.
func (r *postgresRepository) GetImportedSlideIndicesByPresentationID(presentationID string) ([]int, error) {
	if presentationID == "" {
		return nil, nil
	}
	type row struct {
		SlideIndex *int `gorm:"column:slide_index"`
	}
	var rows []row
	// resource_urls is JSONB; slide_index is stored as number in SlideResourceURLs
	err := r.db.Raw(`
		SELECT (resource_urls->>'slide_index')::int AS slide_index
		FROM tasks
		WHERE resource_urls->>'source' = 'google_slides' AND resource_urls->>'presentation_id' = ?
	`, presentationID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	var indices []int
	for _, row := range rows {
		if row.SlideIndex != nil && *row.SlideIndex > 0 {
			indices = append(indices, *row.SlideIndex)
		}
	}
	return indices, nil
}

// --- Project Finance Operations ---

func (r *postgresRepository) UpdateProjectCapital(projectID uuid.UUID, newBalance float64, bonusPct *float64) error {
	updates := map[string]interface{}{"capital_balance": newBalance}
	if bonusPct != nil {
		updates["bonus_percentage"] = *bonusPct
	}
	result := r.db.Model(&domain.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (r *postgresRepository) CreateProjectTransaction(tx *domain.ProjectTransaction) error {
	return r.db.Create(tx).Error
}

func (r *postgresRepository) GetProjectTransactions(projectID uuid.UUID) ([]domain.ProjectTransaction, error) {
	var txns []domain.ProjectTransaction
	err := r.db.Where("project_id = ?", projectID).Order("created_at desc").Find(&txns).Error
	return txns, err
}

func (r *postgresRepository) DeleteProjectTransaction(txID int64, projectID uuid.UUID) error {
	result := r.db.Where("id = ? AND project_id = ?", txID, projectID).Delete(&domain.ProjectTransaction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

// --- Internal B2B Outsource Request Operations ---

func (r *postgresRepository) CreateB2BRequest(req *domain.B2BRequest) error {
	return r.db.Create(req).Error
}

func (r *postgresRepository) GetB2BRequestByID(id uuid.UUID) (*domain.B2BRequest, error) {
	var req domain.B2BRequest
	if err := r.db.First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *postgresRepository) UpdateB2BRequest(req *domain.B2BRequest) error {
	return r.db.Save(req).Error
}

// GetB2BRequests returns requests filtered by team and direction.
// direction "inbound"  → requests where target_team_id = teamID
// direction "outbound" → requests where requester_team_id = teamID
func (r *postgresRepository) GetB2BRequests(teamID uint, direction string) ([]domain.B2BRequest, error) {
	var reqs []domain.B2BRequest
	var q *gorm.DB
	if direction == "outbound" {
		q = r.db.Where("requester_team_id = ?", teamID)
	} else {
		q = r.db.Where("target_team_id = ?", teamID)
	}
	if err := q.Order("created_at desc").Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

// GetFirstProjectByTeamID returns the first active project that belongs to the given team.
func (r *postgresRepository) GetFirstProjectByTeamID(teamID uint) (*domain.Project, error) {
	var project domain.Project
	err := r.db.Where("team_id = ?", teamID).
		Order("created_at asc").
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// --- Project Backup Operations ---

func (r *postgresRepository) CreateProjectBackup(backup *domain.ProjectBackup) error {
	return r.db.Create(backup).Error
}

func (r *postgresRepository) GetProjectBackups(projectID uuid.UUID) ([]domain.ProjectBackup, error) {
	var backups []domain.ProjectBackup
	err := r.db.Where("project_id = ?", projectID).
		Order("created_at desc").
		Find(&backups).Error
	return backups, err
}

func (r *postgresRepository) GetProjectBackupByID(id uuid.UUID) (*domain.ProjectBackup, error) {
	var backup domain.ProjectBackup
	if err := r.db.First(&backup, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &backup, nil
}

func (r *postgresRepository) DeleteProjectBackup(id uuid.UUID, projectID uuid.UUID) error {
	result := r.db.Where("id = ? AND project_id = ?", id, projectID).Delete(&domain.ProjectBackup{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("backup not found")
	}
	return nil
}

// --- Komgrip Operations ---

func (r *postgresRepository) GetKomgripTasks(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("is_komgrip = TRUE").
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}
