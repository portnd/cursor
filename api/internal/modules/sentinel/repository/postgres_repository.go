package repository

import (
	"errors"

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

// --- Task Operations ---

func (r *postgresRepository) CreateTask(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *postgresRepository) GetTaskByID(id uuid.UUID) (*domain.Task, error) {
	var task domain.Task

	// ✅ UPGRADE: Preload "Submissions" ordered by newest first
	err := r.db.Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&task, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &task, nil
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

func (r *postgresRepository) GetAllTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Order("created_at desc").
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
	// Force ID=1 to maintain singleton pattern
	config.ID = 1
	
	// Use Save to update all fields
	return r.db.Save(config).Error
}
