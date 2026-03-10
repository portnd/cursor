package usecase

import (
	"fmt"

	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	pricingDomain "github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type projectFinanceUsecase struct {
	sentinelRepo sentinelDomain.SentinelRepository
	authRepo     authDomain.Repository
	pricingRepo  pricingDomain.Repository
}

// NewProjectFinanceUsecase constructs the ProjectFinanceUsecase.
func NewProjectFinanceUsecase(
	sentinelRepo sentinelDomain.SentinelRepository,
	authRepo authDomain.Repository,
	pricingRepo pricingDomain.Repository,
) sentinelDomain.ProjectFinanceUsecase {
	return &projectFinanceUsecase{
		sentinelRepo: sentinelRepo,
		authRepo:     authRepo,
		pricingRepo:  pricingRepo,
	}
}

// teamMonthlyCost calculates the fully-loaded monthly burn rate for the team that owns the project.
// It reuses the same formula as TeamFinanceUsecase.CalculateTeamMonthlyCost.
func (u *projectFinanceUsecase) teamMonthlyCost(teamID uint) (float64, error) {
	team, err := u.authRepo.GetTeamByID(teamID)
	if err != nil {
		return 0, fmt.Errorf("team not found: %w", err)
	}

	memberIDs := make([]uint, 0, len(team.Users))
	for _, m := range team.Users {
		memberIDs = append(memberIDs, m.ID)
	}

	var memberCost float64
	if len(memberIDs) > 0 {
		salaries, err := u.pricingRepo.GetEmployeeSalaries(memberIDs)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch member salaries: %w", err)
		}
		for _, s := range salaries {
			memberCost += s.MonthlySalary + pricingDomain.SSCost(s.MonthlySalary)
		}
	}

	cfg, err := u.pricingRepo.GetCostConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch cost config: %w", err)
	}

	overheadRoles, err := u.pricingRepo.GetActiveSalariesByRoles([]string{"MANAGER", "SUPPORT"})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch overhead role salaries: %w", err)
	}
	var overheadRoleCost float64
	for _, s := range overheadRoles {
		overheadRoleCost += s.MonthlySalary + pricingDomain.SSCost(s.MonthlySalary)
	}

	allTeams, err := u.authRepo.GetAllTeams()
	if err != nil {
		return 0, fmt.Errorf("failed to count teams: %w", err)
	}
	totalTeams := float64(len(allTeams))
	if totalTeams == 0 {
		totalTeams = 1
	}

	globalOverhead := cfg.ExecutiveExpense + cfg.CompanyExpense + overheadRoleCost
	return memberCost + globalOverhead/totalTeams, nil
}

// GetProjectCapital returns current capital state + burn rate for a project.
func (u *projectFinanceUsecase) GetProjectCapital(projectID uuid.UUID) (*sentinelDomain.ProjectCapitalResponse, error) {
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	project, err := u.sentinelRepo.GetProjectByID(projectID, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	txns, err := u.sentinelRepo.GetProjectTransactions(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	var teamMonthlyCost float64
	if project.TeamID != nil {
		teamMonthlyCost, err = u.teamMonthlyCost(*project.TeamID)
		if err != nil {
			teamMonthlyCost = 0
		}
	}

	var runwayMonths float64
	if teamMonthlyCost > 0 {
		runwayMonths = project.CapitalBalance / teamMonthlyCost
	}

	return &sentinelDomain.ProjectCapitalResponse{
		ProjectID:       projectID,
		ProjectName:     project.Name,
		TeamID:          project.TeamID,
		TeamMonthlyCost: teamMonthlyCost,
		CapitalBalance:  project.CapitalBalance,
		BonusPercentage: project.BonusPercentage,
		RunwayMonths:    runwayMonths,
		Transactions:    txns,
	}, nil
}

// InjectProjectCapital adds capital to the project and records an INJECTION transaction.
func (u *projectFinanceUsecase) InjectProjectCapital(projectID uuid.UUID, req *sentinelDomain.InjectProjectCapitalRequest) (*sentinelDomain.Project, error) {
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	project, err := u.sentinelRepo.GetProjectByID(projectID, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	newBalance := project.CapitalBalance + req.Amount

	var bonusPct *float64
	if req.BonusPercentage >= 0 {
		bonusPct = &req.BonusPercentage
	}

	if err := u.sentinelRepo.UpdateProjectCapital(projectID, newBalance, bonusPct); err != nil {
		return nil, fmt.Errorf("failed to update project capital: %w", err)
	}

	note := req.Note
	if note == "" {
		note = "Capital injection"
	}
	tx := &sentinelDomain.ProjectTransaction{
		ProjectID: projectID,
		Type:      sentinelDomain.ProjTxInjection,
		Amount:    req.Amount,
		Reference: note,
	}
	if err := u.sentinelRepo.CreateProjectTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record transaction: %w", err)
	}

	project.CapitalBalance = newBalance
	if bonusPct != nil {
		project.BonusPercentage = *bonusPct
	}
	return project, nil
}

// EditProjectCapital sets the project capital balance to an exact value and records an ADJUSTMENT.
func (u *projectFinanceUsecase) EditProjectCapital(projectID uuid.UUID, req *sentinelDomain.EditProjectCapitalRequest) (*sentinelDomain.Project, error) {
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	project, err := u.sentinelRepo.GetProjectByID(projectID, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	prevBalance := project.CapitalBalance
	if err := u.sentinelRepo.UpdateProjectCapital(projectID, req.NewBalance, req.BonusPercentage); err != nil {
		return nil, fmt.Errorf("failed to update project capital: %w", err)
	}

	note := req.Note
	if note == "" {
		note = fmt.Sprintf("Manual adjustment: ฿%.2f → ฿%.2f", prevBalance, req.NewBalance)
	}
	delta := req.NewBalance - prevBalance
	tx := &sentinelDomain.ProjectTransaction{
		ProjectID: projectID,
		Type:      sentinelDomain.ProjTxAdjustment,
		Amount:    delta,
		Reference: note,
	}
	if err := u.sentinelRepo.CreateProjectTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record adjustment: %w", err)
	}

	project.CapitalBalance = req.NewBalance
	if req.BonusPercentage != nil {
		project.BonusPercentage = *req.BonusPercentage
	}
	return project, nil
}

// CloseProjectCycleAndPayout calculates the bonus from remaining balance,// records a BONUS_PAYOUT transaction, and resets the project capital to 0.
func (u *projectFinanceUsecase) CloseProjectCycleAndPayout(projectID uuid.UUID) (*sentinelDomain.CloseProjectCycleResponse, error) {
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	project, err := u.sentinelRepo.GetProjectByID(projectID, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	balanceBefore := project.CapitalBalance
	bonusAmount := balanceBefore * (project.BonusPercentage / 100.0)

	if err := u.sentinelRepo.UpdateProjectCapital(projectID, 0, nil); err != nil {
		return nil, fmt.Errorf("failed to reset project capital: %w", err)
	}

	ref := fmt.Sprintf("Milestone close — %.2f%% bonus payout (฿%.2f)", project.BonusPercentage, bonusAmount)
	tx := &sentinelDomain.ProjectTransaction{
		ProjectID: projectID,
		Type:      sentinelDomain.ProjTxBonusPayout,
		Amount:    bonusAmount,
		Reference: ref,
	}
	if err := u.sentinelRepo.CreateProjectTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record payout transaction: %w", err)
	}

	return &sentinelDomain.CloseProjectCycleResponse{
		ProjectID:       projectID,
		BalanceBefore:   balanceBefore,
		BonusPercentage: project.BonusPercentage,
		BonusAmount:     bonusAmount,
		BalanceAfter:    0,
	}, nil
}

// DeleteProjectTransaction removes a single transaction record for a project.
func (u *projectFinanceUsecase) DeleteProjectTransaction(txID int64, projectID uuid.UUID) error {
	return u.sentinelRepo.DeleteProjectTransaction(txID, projectID)
}
