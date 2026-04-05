package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/deployment/domain"
)

type postgresRepo struct {
	db       *gorm.DB
	authRepo authDomain.Repository
}

// NewPostgresRepository creates a new PostgreSQL-backed deployment repository.
func NewPostgresRepository(db *gorm.DB, authRepo authDomain.Repository) domain.Repository {
	return &postgresRepo{db: db, authRepo: authRepo}
}

func (r *postgresRepo) Create(req *domain.DeploymentRequest) error {
	if err := r.db.Create(req).Error; err != nil {
		return fmt.Errorf("deployment: create: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetByID(id uint) (*domain.DeploymentRequest, error) {
	var req domain.DeploymentRequest
	if err := r.db.First(&req, id).Error; err != nil {
		return nil, fmt.Errorf("deployment: get by id: %w", err)
	}
	r.enrichSingle(&req)
	return &req, nil
}

// GetByTaskID returns the most recent deployment request linked to a sentinel task.
func (r *postgresRepo) GetByTaskID(taskID uuid.UUID) (*domain.DeploymentRequest, error) {
	var req domain.DeploymentRequest
	if err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").First(&req).Error; err != nil {
		return nil, fmt.Errorf("deployment: get by task_id: %w", err)
	}
	r.enrichSingle(&req)
	return &req, nil
}

// List returns requests filtered by optional requesterID and/or status.
// CEO / CHIEF_ENGINEER sees all; engineers see only their own (handled in usecase).
func (r *postgresRepo) List(requesterID *uint, status string) ([]domain.DeploymentRequest, error) {
	q := r.db.Order("created_at DESC")
	if requesterID != nil {
		q = q.Where("requester_id = ?", *requesterID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var reqs []domain.DeploymentRequest
	if err := q.Find(&reqs).Error; err != nil {
		return nil, fmt.Errorf("deployment: list: %w", err)
	}
	for i := range reqs {
		r.enrichSingle(&reqs[i])
	}
	return reqs, nil
}

func (r *postgresRepo) Update(req *domain.DeploymentRequest) error {
	if err := r.db.Save(req).Error; err != nil {
		return fmt.Errorf("deployment: update: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetStats() (*domain.DeploymentStats, error) {
	type row struct {
		Status string
		Count  int
	}
	var rows []row
	if err := r.db.Raw(`
		SELECT status, COUNT(*) AS count
		FROM deployment_requests
		GROUP BY status
	`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("deployment: stats: %w", err)
	}

	stats := &domain.DeploymentStats{}
	for _, rw := range rows {
		switch rw.Status {
		case domain.StatusPending:
			stats.TotalPending = rw.Count
		case domain.StatusReviewing:
			stats.TotalReviewing = rw.Count
		case domain.StatusApproved:
			stats.TotalApproved = rw.Count
		case domain.StatusRejected:
			stats.TotalRejected = rw.Count
		case domain.StatusDeployed:
			stats.TotalDeployed = rw.Count
		}
	}

	// Deployed today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var deployedToday int64
	r.db.Model(&domain.DeploymentRequest{}).
		Where("status = ? AND deployed_at >= ?", domain.StatusDeployed, today).
		Count(&deployedToday)
	stats.DeployedToday = int(deployedToday)

	return stats, nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func (r *postgresRepo) enrichSingle(req *domain.DeploymentRequest) {
	if u, err := r.authRepo.FindByID(req.RequesterID); err == nil && u != nil {
		req.RequesterEmail = u.Email
		req.RequesterDisplayName = u.DisplayName
	}
	if req.ReviewerID != nil {
		if u, err := r.authRepo.FindByID(*req.ReviewerID); err == nil && u != nil {
			req.ReviewerEmail = u.Email
			req.ReviewerDisplayName = u.DisplayName
		}
	}
}
