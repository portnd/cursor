package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/deployment/domain"
)

type deploymentUsecase struct {
	repo     domain.Repository
	advancer domain.TaskStatusAdvancer // optional; nil means no auto-advance
}

// NewDeploymentUsecase wires up the deployment usecase.
// advancer may be nil (feature is simply disabled if nil).
func NewDeploymentUsecase(repo domain.Repository, advancer domain.TaskStatusAdvancer) domain.Usecase {
	return &deploymentUsecase{repo: repo, advancer: advancer}
}

// isChiefEngineer returns true only for CHIEF_ENGINEER.
func isChiefEngineer(role string) bool {
	return strings.ToUpper(strings.TrimSpace(role)) == authDomain.RoleChiefEngineer
}

// canSubmit returns true for roles allowed to create deployment requests.
func canSubmit(role string) bool {
	r := strings.ToUpper(strings.TrimSpace(role))
	switch r {
	case authDomain.RoleCEO,
		authDomain.RoleManager,
		authDomain.RoleProductOwner,
		authDomain.RoleEngineer,
		authDomain.RoleChiefEngineer:
		return true
	}
	return false
}

// canViewAll returns true for roles that see every request.
func canViewAll(role string) bool {
	r := strings.ToUpper(strings.TrimSpace(role))
	switch r {
	case authDomain.RoleCEO, authDomain.RoleChiefEngineer:
		return true
	}
	return false
}

// ─── Usecase methods ─────────────────────────────────────────────────────────

func (u *deploymentUsecase) CreateRequest(callerID uint, callerRole string, in domain.CreateDeploymentRequest) (*domain.DeploymentRequest, error) {
	if !canSubmit(callerRole) {
		return nil, domain.ErrForbidden
	}
	if strings.TrimSpace(in.Title) == "" {
		return nil, domain.ErrMissingTitle
	}
	if strings.TrimSpace(in.Branch) == "" {
		return nil, domain.ErrMissingBranch
	}

	env := strings.ToUpper(strings.TrimSpace(in.Environment))
	if env == "" {
		env = domain.EnvStaging
	}
	if env != domain.EnvStaging && env != domain.EnvPreProd && env != domain.EnvProduction {
		return nil, domain.ErrInvalidEnv
	}

	req := &domain.DeploymentRequest{
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		Branch:      strings.TrimSpace(in.Branch),
		PRUrl:       strings.TrimSpace(in.PRUrl),
		Environment: env,
		Status:      domain.StatusPending,
		RequesterID: callerID,
		ReviewerID:  in.ReviewerID,
		TaskRef:     strings.TrimSpace(in.TaskRef),
	}

	// Parse optional task UUID link
	if tidStr := strings.TrimSpace(in.TaskID); tidStr != "" {
		if parsed, err := uuid.Parse(tidStr); err == nil {
			req.TaskID = &parsed
		}
	}

	if err := u.repo.Create(req); err != nil {
		return nil, fmt.Errorf("deployment: create: %w", err)
	}
	return req, nil
}

func (u *deploymentUsecase) ListRequests(callerID uint, callerRole string, status string) ([]domain.DeploymentRequest, error) {
	if canViewAll(callerRole) {
		return u.repo.List(nil, status)
	}
	return u.repo.List(&callerID, status)
}

func (u *deploymentUsecase) GetRequest(callerID uint, callerRole string, id uint) (*domain.DeploymentRequest, error) {
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !canViewAll(callerRole) && req.RequesterID != callerID {
		return nil, domain.ErrForbidden
	}
	return req, nil
}

func (u *deploymentUsecase) GetRequestByTaskID(callerID uint, callerRole string, taskIDStr string) (*domain.DeploymentRequest, error) {
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	req, err := u.repo.GetByTaskID(taskID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !canViewAll(callerRole) && req.RequesterID != callerID {
		return nil, domain.ErrForbidden
	}
	return req, nil
}

func (u *deploymentUsecase) PickForReview(callerID uint, callerRole string, id uint) (*domain.DeploymentRequest, error) {
	if !isChiefEngineer(callerRole) {
		return nil, domain.ErrForbidden
	}
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if req.Status != domain.StatusPending {
		return nil, domain.ErrAlreadyReviewing
	}
	req.Status = domain.StatusReviewing
	req.ReviewerID = &callerID
	if err := u.repo.Update(req); err != nil {
		return nil, fmt.Errorf("deployment: pick-review: %w", err)
	}
	return req, nil
}

func (u *deploymentUsecase) Approve(callerID uint, callerRole string, id uint, in domain.ReviewActionRequest) (*domain.DeploymentRequest, error) {
	if !isChiefEngineer(callerRole) {
		return nil, domain.ErrForbidden
	}
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if req.Status != domain.StatusPending && req.Status != domain.StatusReviewing {
		return nil, domain.ErrInvalidStatus
	}
	req.Status = domain.StatusApproved
	req.ReviewerID = &callerID
	req.ReviewNotes = strings.TrimSpace(in.Notes)
	if err := u.repo.Update(req); err != nil {
		return nil, fmt.Errorf("deployment: approve: %w", err)
	}
	return req, nil
}

func (u *deploymentUsecase) Reject(callerID uint, callerRole string, id uint, in domain.ReviewActionRequest) (*domain.DeploymentRequest, error) {
	if !isChiefEngineer(callerRole) {
		return nil, domain.ErrForbidden
	}
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if req.Status != domain.StatusPending && req.Status != domain.StatusReviewing {
		return nil, domain.ErrInvalidStatus
	}
	req.Status = domain.StatusRejected
	req.ReviewerID = &callerID
	req.ReviewNotes = strings.TrimSpace(in.Notes)
	req.RejectionReason = strings.TrimSpace(in.Reason)
	if err := u.repo.Update(req); err != nil {
		return nil, fmt.Errorf("deployment: reject: %w", err)
	}
	return req, nil
}

func (u *deploymentUsecase) MarkDeployed(callerID uint, callerRole string, id uint, in domain.ReviewActionRequest) (*domain.DeploymentRequest, error) {
	if !isChiefEngineer(callerRole) {
		return nil, domain.ErrForbidden
	}
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if req.Status != domain.StatusApproved {
		return nil, domain.ErrInvalidStatus
	}
	now := time.Now().UTC()
	req.Status = domain.StatusDeployed
	req.ReviewerID = &callerID
	req.ReviewNotes = strings.TrimSpace(in.Notes)
	req.DeployedAt = &now
	if err := u.repo.Update(req); err != nil {
		return nil, fmt.Errorf("deployment: mark-deployed: %w", err)
	}

	// Auto-advance linked task: WAIT_FOR_DEPLOY → READY_FOR_UAT
	if req.TaskID != nil && u.advancer != nil {
		if advErr := u.advancer.AdvanceTaskAfterDeploy(*req.TaskID, callerID); advErr != nil {
			// Log but don't fail — deployment already saved; task advance can be retried manually
			fmt.Printf("⚠️  deployment: advance task %s after deploy: %v\n", req.TaskID, advErr)
		}
	}

	return req, nil
}

func (u *deploymentUsecase) GetStats(callerID uint, callerRole string) (*domain.DeploymentStats, error) {
	return u.repo.GetStats()
}
