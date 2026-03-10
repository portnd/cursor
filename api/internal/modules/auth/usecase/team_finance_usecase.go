package usecase

import (
	"fmt"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	pricingDomain "github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
)

type teamFinanceUsecase struct {
	authRepo    authDomain.Repository
	pricingRepo pricingDomain.Repository
}

// NewTeamFinanceUsecase constructs a TeamFinanceUsecase with both auth and pricing repo dependencies.
func NewTeamFinanceUsecase(authRepo authDomain.Repository, pricingRepo pricingDomain.Repository) authDomain.TeamFinanceUsecase {
	return &teamFinanceUsecase{
		authRepo:    authRepo,
		pricingRepo: pricingRepo,
	}
}

// CalculateTeamMonthlyCost computes the fully loaded monthly burn rate for a team.
//
// Formula:
//
//	globalOverhead = cfg.ExecutiveExpense + cfg.CompanyExpense
//	               + Σ (mgr.salary + SS(mgr.salary))  // ALL MANAGER + SUPPORT company-wide
//
//	monthly_cost(team) = Σ (member.salary + SS(member.salary))
//	                   + globalOverhead / totalTeams
func (u *teamFinanceUsecase) CalculateTeamMonthlyCost(teamID uint) (*authDomain.TeamMonthlyCostResponse, error) {
	team, err := u.authRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Collect user IDs of all team members
	memberIDs := make([]uint, 0, len(team.Users))
	for _, u := range team.Users {
		memberIDs = append(memberIDs, u.ID)
	}

	// Sum loaded salary (salary + SS) for members in this team
	var memberCost float64
	if len(memberIDs) > 0 {
		salaries, err := u.pricingRepo.GetEmployeeSalaries(memberIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch member salaries: %w", err)
		}
		for _, s := range salaries {
			memberCost += s.MonthlySalary + pricingDomain.SSCost(s.MonthlySalary)
		}
	}

	// Global overhead: company config costs
	cfg, err := u.pricingRepo.GetCostConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cost config: %w", err)
	}

	// Global overhead: MANAGER + SUPPORT salaries company-wide
	overheadRoles, err := u.pricingRepo.GetActiveSalariesByRoles([]string{"MANAGER", "SUPPORT"})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch overhead role salaries: %w", err)
	}
	var overheadRoleCost float64
	for _, s := range overheadRoles {
		overheadRoleCost += s.MonthlySalary + pricingDomain.SSCost(s.MonthlySalary)
	}

	// Total teams for overhead allocation
	allTeams, err := u.authRepo.GetAllTeams()
	if err != nil {
		return nil, fmt.Errorf("failed to count teams: %w", err)
	}
	totalTeams := float64(len(allTeams))
	if totalTeams == 0 {
		totalTeams = 1
	}

	globalOverhead := cfg.ExecutiveExpense + cfg.CompanyExpense + overheadRoleCost
	sharedOverhead := globalOverhead / totalTeams
	totalMonthlyCost := memberCost + sharedOverhead

	var runwayMonths float64
	if totalMonthlyCost > 0 {
		runwayMonths = team.CapitalBalance / totalMonthlyCost
	}

	return &authDomain.TeamMonthlyCostResponse{
		TeamID:           teamID,
		MemberCost:       memberCost,
		SharedOverhead:   sharedOverhead,
		TotalMonthlyCost: totalMonthlyCost,
		CapitalBalance:   team.CapitalBalance,
		BonusPercentage:  team.BonusPercentage,
		RunwayMonths:     runwayMonths,
	}, nil
}

// InjectCapital adds capital to a team and records an INJECTION transaction.
// Optionally updates the team's bonus percentage target.
func (u *teamFinanceUsecase) InjectCapital(teamID uint, req *authDomain.InjectCapitalRequest) (*authDomain.Team, error) {
	team, err := u.authRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	newBalance := team.CapitalBalance + req.Amount

	var bonusPctPtr *float64
	if req.BonusPercentage >= 0 {
		bonusPctPtr = &req.BonusPercentage
	}

	if err := u.authRepo.UpdateTeamCapital(teamID, newBalance, bonusPctPtr); err != nil {
		return nil, fmt.Errorf("failed to update capital: %w", err)
	}

	note := req.Note
	if note == "" {
		note = "Capital injection"
	}
	tx := &authDomain.TeamTransaction{
		TeamID:    teamID,
		Type:      authDomain.TxInjection,
		Amount:    req.Amount,
		Reference: note,
	}
	if err := u.authRepo.CreateTeamTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record transaction: %w", err)
	}

	// Return refreshed team
	team.CapitalBalance = newBalance
	if bonusPctPtr != nil {
		team.BonusPercentage = *bonusPctPtr
	}
	return team, nil
}

// EditCapital directly sets a team's capital balance to an exact value and records an ADJUSTMENT.
// Use this to correct balance after manual reconciliation or CEO override.
func (u *teamFinanceUsecase) EditCapital(teamID uint, req *authDomain.EditCapitalRequest) (*authDomain.Team, error) {
	team, err := u.authRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	prevBalance := team.CapitalBalance
	if err := u.authRepo.UpdateTeamCapital(teamID, req.NewBalance, req.BonusPercentage); err != nil {
		return nil, fmt.Errorf("failed to update capital: %w", err)
	}

	note := req.Note
	if note == "" {
		note = fmt.Sprintf("Manual adjustment: ฿%.2f → ฿%.2f", prevBalance, req.NewBalance)
	}
	// Record the delta as the transaction amount (can be negative for reduction)
	delta := req.NewBalance - prevBalance
	tx := &authDomain.TeamTransaction{
		TeamID:    teamID,
		Type:      authDomain.TxAdjustment,
		Amount:    delta,
		Reference: note,
	}
	if err := u.authRepo.CreateTeamTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record adjustment transaction: %w", err)
	}

	team.CapitalBalance = req.NewBalance
	if req.BonusPercentage != nil {
		team.BonusPercentage = *req.BonusPercentage
	}
	return team, nil
}

// CloseCycleAndPayout calculates the bonus from remaining balance, records a BONUS_PAYOUT
// transaction, and resets the capital balance to 0.
func (u *teamFinanceUsecase) CloseCycleAndPayout(teamID uint) (*authDomain.CloseCycleResponse, error) {
	team, err := u.authRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	balanceBefore := team.CapitalBalance
	bonusAmount := balanceBefore * (team.BonusPercentage / 100.0)

	if err := u.authRepo.UpdateTeamCapital(teamID, 0, nil); err != nil {
		return nil, fmt.Errorf("failed to reset capital balance: %w", err)
	}

	ref := fmt.Sprintf("Milestone close — %.2f%% bonus payout (฿%.2f)", team.BonusPercentage, bonusAmount)
	tx := &authDomain.TeamTransaction{
		TeamID:    teamID,
		Type:      authDomain.TxBonusPayout,
		Amount:    bonusAmount,
		Reference: ref,
	}
	if err := u.authRepo.CreateTeamTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to record payout transaction: %w", err)
	}

	return &authDomain.CloseCycleResponse{
		TeamID:          teamID,
		BalanceBefore:   balanceBefore,
		BonusPercentage: team.BonusPercentage,
		BonusAmount:     bonusAmount,
		BalanceAfter:    0,
	}, nil
}
