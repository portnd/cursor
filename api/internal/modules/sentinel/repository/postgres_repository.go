package repository

import (
	"errors"
	"fmt"
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

func (r *postgresRepository) CreateProject(p *domain.Project) error {
	return r.db.Create(p).Error
}

func (r *postgresRepository) GetAllProjects() ([]domain.Project, error) {
	var projects []domain.Project
	err := r.db.Order("created_at desc").Find(&projects).Error
	return projects, err
}

func (r *postgresRepository) GetProjectByID(id uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	err := r.db.First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *postgresRepository) GetProjectByCode(code string) (*domain.Project, error) {
	var project domain.Project
	err := r.db.First(&project, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *postgresRepository) GetTasksByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("project_id = ?", projectID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
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
		return db.Order("created_at desc").Limit(100)
	}).First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *postgresRepository) GetTaskByCode(code string) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(100)
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

func (r *postgresRepository) GetUnassignedTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("assigned_to IS NULL").
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

// maxAllTasksLimit caps GET /tasks when no project_id (CEO/PM dashboard) to avoid slow full scan
const maxAllTasksLimit = 2000

func (r *postgresRepository) GetAllTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Order("created_at desc").
		Limit(maxAllTasksLimit).
		Find(&tasks).Error
	return tasks, err
}

// GetTasksRequiringApproval returns tasks that need PM/CEO attention
// Logic:
// - Tasks with negotiation_status = 'PENDING' (dev wants more time)
// - Tasks with submissions that have appeals where status = 'PENDING'
func (r *postgresRepository) GetTasksRequiringApproval() ([]domain.Task, error) {
	var tasks []domain.Task

	// Subquery: Find task IDs with PENDING appeals
	var taskIDsWithPendingAppeals []string
	r.db.Table("tasks").
		Select("DISTINCT tasks.id").
		Joins("JOIN submissions ON submissions.task_id = tasks.id").
		Joins("JOIN appeals ON appeals.submission_id = submissions.id").
		Where("appeals.status = ?", "PENDING").
		Pluck("tasks.id", &taskIDsWithPendingAppeals)

	// Main query: Find tasks with PENDING time negotiations OR PENDING appeals
	// ALWAYS preload submissions and appeals for ALL tasks
	query := r.db.
		Preload("Submissions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("Submissions.Appeal") // Eager load appeals within submissions

	// Build WHERE condition
	if len(taskIDsWithPendingAppeals) > 0 {
		query = query.Where("negotiation_status = ? OR id IN ?", "PENDING", taskIDsWithPendingAppeals)
	} else {
		query = query.Where("negotiation_status = ?", "PENDING")
	}

	err := query.Order("created_at desc").Find(&tasks).Error

	return tasks, err
}

func (r *postgresRepository) UpdateTask(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *postgresRepository) DeleteTask(id uuid.UUID) error {
	return r.db.Delete(&domain.Task{}, "id = ?", id).Error
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
	err := r.db.Where("project_id = ?", projectID).Order("created_at desc").Find(&sprints).Error
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

// --- Time Log Operations ---

func (r *postgresRepository) CreateTimeLog(t *domain.TimeLog) error {
	return r.db.Create(t).Error
}

func (r *postgresRepository) GetTimeLogsByTaskID(taskID uuid.UUID) ([]domain.TimeLog, error) {
	var logs []domain.TimeLog
	err := r.db.Where("task_id = ?", taskID).Order("logged_at desc").Find(&logs).Error
	return logs, err
}

func (r *postgresRepository) GetTotalLoggedMinutes(taskID uuid.UUID) (int, error) {
	var total int64
	err := r.db.Model(&domain.TimeLog{}).Where("task_id = ?", taskID).Select("COALESCE(SUM(minutes), 0)").Scan(&total).Error
	return int(total), err
}

// --- Bulk Operations ---

func (r *postgresRepository) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error {
	if len(taskIDs) == 0 {
		return nil
	}
	return r.db.Model(&domain.Task{}).Where("id IN ?", taskIDs).Update("status", status).Error
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

	// Team capacity
	type capacityRow struct {
		UserID         uint
		AssignedTasks  int
		EstimatedMins  int
		LoggedMins     int
	}
	var capRows []capacityRow
	r.db.Raw(`
		SELECT 
			t.assigned_to as user_id,
			COUNT(t.id) as assigned_tasks,
			COALESCE(SUM(t.ai_estimated_minutes), 0) as estimated_mins,
			COALESCE(SUM(tl_sum.total_logged), 0) as logged_mins
		FROM tasks t
		LEFT JOIN (
			SELECT task_id, SUM(minutes) as total_logged FROM time_logs GROUP BY task_id
		) tl_sum ON tl_sum.task_id = t.id
		WHERE t.project_id = ? AND t.assigned_to IS NOT NULL
		GROUP BY t.assigned_to
	`, projectID).Scan(&capRows)

	for _, row := range capRows {
		util := 0.0
		if row.EstimatedMins > 0 {
			util = float64(row.LoggedMins) / float64(row.EstimatedMins) * 100
		}
		analytics.TeamCapacity = append(analytics.TeamCapacity, domain.TeamCapacityRow{
			UserID:         row.UserID,
			UserEmail:      fmt.Sprintf("user-%d", row.UserID),
			AssignedTasks:  row.AssignedTasks,
			EstimatedHours: float64(row.EstimatedMins) / 60,
			LoggedHours:    float64(row.LoggedMins) / 60,
			Utilization:    util,
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

// GetEpicTimelineData returns all epics for a project with their tasks preloaded, ordered by start_date.
func (r *postgresRepository) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	var epics []domain.Epic
	err := r.db.Where("project_id = ?", projectID).
		Order("sort_order asc, start_date asc NULLS LAST, created_at asc").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").Order("start_date asc NULLS LAST, created_at asc")
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
