package usecase

import (
	"fmt"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	perfDomain "github.com/portnd/the-sentinel-core/internal/modules/performance/domain"
	"math"
)

type performanceUsecase struct {
	repo     perfDomain.Repository
	authRepo authDomain.Repository
}

// NewPerformanceUsecase creates the performance usecase with repo and auth repo (for team list)
func NewPerformanceUsecase(repo perfDomain.Repository, authRepo authDomain.Repository) perfDomain.Usecase {
	return &performanceUsecase{repo: repo, authRepo: authRepo}
}

func (u *performanceUsecase) GetPersonalKPIs(userID uint, role string) (*perfDomain.PersonalKPIs, error) {
	email, _, healthScore, err := u.repo.GetUserEmailAndRole(userID)
	if err != nil {
		return nil, err
	}
	out := &perfDomain.PersonalKPIs{
		UserID:      userID,
		Email:       email,
		Role:        role,
		HealthScore: healthScore,
	}
	if !authDomain.IsEngineerRole(role) {
		return out, nil
	}
	tasksWithDue, completedOnTime, err := u.repo.GetUserTaskDeliveryStats(userID)
	if err != nil {
		return nil, err
	}
	if tasksWithDue > 0 {
		out.DeliveryRatePct = float64(completedOnTime) / float64(tasksWithDue) * 100
	}
	avgScore, _, _, err := u.repo.GetUserSubmissionStats(userID)
	if err != nil {
		return nil, err
	}
	out.CodeQualityIndex = avgScore
	jobDoneCount, reworkCount, err := u.repo.GetUserReworkStats(userID)
	if err != nil {
		return nil, err
	}
	if totalOutcomes := jobDoneCount + reworkCount; totalOutcomes > 0 {
		out.ReworkRatePct = float64(reworkCount) / float64(totalOutcomes) * 100
	}
	accPct, _, err := u.repo.GetUserTimeAccuracy(userID)
	if err != nil {
		return nil, err
	}
	out.TimeAccuracyPct = accPct
	avgSP, trend, err := u.repo.GetUserSprintVelocity(userID, 3)
	if err != nil {
		return nil, err
	}
	out.SprintVelocitySP = avgSP
	out.VelocityTrend = trend
	return out, nil
}

// compositeScore computes 0-100 score from delivery, quality, rework, velocity, time accuracy
func compositeScore(deliveryPct, codeQuality, reworkRatePct, timeAccuracyPct float64, sprintVelocitySP float64) float64 {
	qualityNorm := codeQuality
	if qualityNorm > 100 {
		qualityNorm = 100
	}
	reworkNorm := 100 - reworkRatePct
	if reworkNorm < 0 {
		reworkNorm = 0
	}
	velocityNorm := sprintVelocitySP * 5
	if velocityNorm > 100 {
		velocityNorm = 100
	}
	return 0.30*deliveryPct + 0.25*qualityNorm + 0.20*reworkNorm + 0.15*velocityNorm + 0.10*timeAccuracyPct
}

func (u *performanceUsecase) GetTeamKPIs(requestingUserID uint, requestingRole string) (*perfDomain.TeamKPIsResponse, error) {
	if requestingRole != "CEO" && requestingRole != authDomain.RoleProductOwner {
		return &perfDomain.TeamKPIsResponse{Members: []perfDomain.TeamMemberKPI{}}, nil
	}

	if requestingRole == authDomain.RoleProductOwner {
		// Product Owner sees only engineers assigned at least one task by this Product Owner; scores from those tasks only.
		devIDs, err := u.repo.GetDevUserIDsAssignedByPM(requestingUserID)
		if err != nil {
			return nil, err
		}
		members := make([]perfDomain.TeamMemberKPI, 0, len(devIDs))
		for _, devID := range devIDs {
			email, devRole, healthScore, _ := u.repo.GetUserEmailAndRole(devID)
			if !authDomain.IsEngineerRole(devRole) {
				continue
			}
			m := perfDomain.TeamMemberKPI{
				UserID:      devID,
				Email:       email,
				Role:        devRole,
				HealthScore: healthScore,
			}
			tasksWithDue, completedOnTime, _ := u.repo.GetUserTaskDeliveryStatsForAssignedBy(devID, requestingUserID)
			if tasksWithDue > 0 {
				m.DeliveryRatePct = float64(completedOnTime) / float64(tasksWithDue) * 100
			}
			avgScore, _, _, _ := u.repo.GetUserSubmissionStatsForAssignedBy(devID, requestingUserID)
			m.CodeQualityIndex = avgScore
			jobDoneCount, reworkCount, _ := u.repo.GetUserReworkStatsForAssignedBy(devID, requestingUserID)
			if totalOutcomes := jobDoneCount + reworkCount; totalOutcomes > 0 {
				m.ReworkRatePct = float64(reworkCount) / float64(totalOutcomes) * 100
			}
			accPct, _, _ := u.repo.GetUserTimeAccuracyForAssignedBy(devID, requestingUserID)
			m.TimeAccuracyPct = accPct
			avgSP, _, _ := u.repo.GetUserSprintVelocityForAssignedBy(devID, requestingUserID, 3)
			m.SprintVelocitySP = avgSP
			m.CompositeScore = compositeScore(m.DeliveryRatePct, m.CodeQualityIndex, m.ReworkRatePct, m.TimeAccuracyPct, m.SprintVelocitySP)
			members = append(members, m)
		}
		// Sort by composite score descending (leaderboard)
		for i := 0; i < len(members); i++ {
			for j := i + 1; j < len(members); j++ {
				if members[j].CompositeScore > members[i].CompositeScore {
					members[i], members[j] = members[j], members[i]
				}
			}
		}
		return &perfDomain.TeamKPIsResponse{Members: members}, nil
	}

	// CEO: all users. Engineer = total score (all tasks); Product Owner = team score (avg of assigned-engineer composites).
	users, err := u.authRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	members := make([]perfDomain.TeamMemberKPI, 0, len(users))
	for _, usr := range users {
		m := perfDomain.TeamMemberKPI{
			UserID:      usr.ID,
			Email:       usr.Email,
			Role:        usr.Role,
			HealthScore: usr.HealthScore,
		}
		if authDomain.IsEngineerRole(usr.Role) {
			tasksWithDue, completedOnTime, _ := u.repo.GetUserTaskDeliveryStats(usr.ID)
			if tasksWithDue > 0 {
				m.DeliveryRatePct = float64(completedOnTime) / float64(tasksWithDue) * 100
			}
			avgScore, _, _, _ := u.repo.GetUserSubmissionStats(usr.ID)
			m.CodeQualityIndex = avgScore
			jobDoneCount, reworkCount, _ := u.repo.GetUserReworkStats(usr.ID)
			if totalOutcomes := jobDoneCount + reworkCount; totalOutcomes > 0 {
				m.ReworkRatePct = float64(reworkCount) / float64(totalOutcomes) * 100
			}
			accPct, _, _ := u.repo.GetUserTimeAccuracy(usr.ID)
			m.TimeAccuracyPct = accPct
			avgSP, _, _ := u.repo.GetUserSprintVelocity(usr.ID, 3)
			m.SprintVelocitySP = avgSP
			m.CompositeScore = compositeScore(m.DeliveryRatePct, m.CodeQualityIndex, m.ReworkRatePct, m.TimeAccuracyPct, m.SprintVelocitySP)
		} else if usr.Role == authDomain.RoleProductOwner {
			// Product Owner's score = team score from "My Accountability KPIs" (average of composites of engineers assigned by this Product Owner)
			devIDs, _ := u.repo.GetDevUserIDsAssignedByPM(usr.ID)
			var sumComposite float64
			var count int
			for _, devID := range devIDs {
				_, devRole, _, _ := u.repo.GetUserEmailAndRole(devID)
				if !authDomain.IsEngineerRole(devRole) {
					continue
				}
				tasksWithDue, completedOnTime, _ := u.repo.GetUserTaskDeliveryStatsForAssignedBy(devID, usr.ID)
				var deliveryPct float64
				if tasksWithDue > 0 {
					deliveryPct = float64(completedOnTime) / float64(tasksWithDue) * 100
				}
				avgScore, _, _, _ := u.repo.GetUserSubmissionStatsForAssignedBy(devID, usr.ID)
				var reworkPct float64
				jobDoneCount, reworkCount, _ := u.repo.GetUserReworkStatsForAssignedBy(devID, usr.ID)
				if totalOutcomes := jobDoneCount + reworkCount; totalOutcomes > 0 {
					reworkPct = float64(reworkCount) / float64(totalOutcomes) * 100
				}
				accPct, _, _ := u.repo.GetUserTimeAccuracyForAssignedBy(devID, usr.ID)
				avgSP, _, _ := u.repo.GetUserSprintVelocityForAssignedBy(devID, usr.ID, 3)
				sumComposite += compositeScore(deliveryPct, avgScore, reworkPct, accPct, avgSP)
				count++
			}
			if count > 0 {
				m.CompositeScore = sumComposite / float64(count)
			} else {
				m.CompositeScore = usr.HealthScore
			}
		} else {
			m.CompositeScore = usr.HealthScore
		}
		members = append(members, m)
	}
	// Sort by composite score descending (leaderboard)
	for i := 0; i < len(members); i++ {
		for j := i + 1; j < len(members); j++ {
			if members[j].CompositeScore > members[i].CompositeScore {
				members[i], members[j] = members[j], members[i]
			}
		}
	}
	return &perfDomain.TeamKPIsResponse{Members: members}, nil
}

func (u *performanceUsecase) GetOverviewKPIs(requestingUserID uint, requestingRole string) (*perfDomain.OverviewKPIs, error) {
	if requestingRole != "CEO" {
		return nil, nil
	}
	out := &perfDomain.OverviewKPIs{}
	sprintRate, err := u.repo.GetSprintSuccessRate()
	if err != nil {
		return nil, err
	}
	out.ProjectOnTrackRatePct, err = u.repo.GetProjectOnTrackRate()
	if err != nil {
		return nil, err
	}
	reached, missed, err := u.repo.GetMilestoneHitRate()
	if err != nil {
		return nil, err
	}
	if reached+missed > 0 {
		out.MilestoneHitRatePct = float64(reached) / float64(reached+missed) * 100
	}
	out.CursorAdoptionScore, err = u.repo.GetCursorAdoptionScore()
	if err != nil {
		return nil, err
	}
	out.TeamVelocityTrendPct, err = u.repo.GetTeamVelocityTrend()
	if err != nil {
		return nil, err
	}
	avgDelivery, avgQuality, err := u.repo.GetCompanyWideDeliveryAndQuality()
	if err != nil {
		return nil, err
	}
	avgRework, avgTimeAcc, err := u.repo.GetCompanyWideReworkAndTimeAccuracy()
	if err != nil {
		return nil, err
	}
	// Engineering Health Index: blend of delivery, quality, rework, velocity trend, time accuracy
	velocityNorm := 50 + out.TeamVelocityTrendPct*0.5
	if velocityNorm < 0 {
		velocityNorm = 0
	}
	if velocityNorm > 100 {
		velocityNorm = 100
	}
	out.SprintSuccessRatePct = sprintRate
	out.EngineeringHealthIndex = 0.30*avgDelivery + 0.25*avgQuality + 0.20*(100-avgRework) + 0.15*velocityNorm + 0.10*avgTimeAcc
	out.EngineeringHealthIndex = math.Round(out.EngineeringHealthIndex*100) / 100
	return out, nil
}

// GetDisciplineDayDetail returns drill-down activity for one user on one day.
func (u *performanceUsecase) GetDisciplineDayDetail(userID uint, date string) (*perfDomain.DisciplineDayDetail, error) {
	return u.repo.GetDisciplineDayDetail(userID, date)
}

// GetDiscipline returns daily discipline stats for all users in the date range.
// from/to format: YYYY-MM-DD. Accessible to CEO and Product Owner.
func (u *performanceUsecase) GetDiscipline(from, to string) (*perfDomain.DisciplineResponse, error) {
	return u.repo.GetDisciplineStats(from, to)
}

func (u *performanceUsecase) GetDisciplineStartDate() (*perfDomain.DisciplineStartDateResponse, error) {
	return u.repo.GetDisciplineStartDate()
}

func (u *performanceUsecase) SetDisciplineStartDate(startDate string) (*perfDomain.DisciplineStartDateResponse, error) {
	return u.repo.SetDisciplineStartDate(startDate)
}

// ResetReworkRate clears the rework history for a developer by setting rework_reset_at = NOW().
// Only CEO is allowed to perform this action.
func (u *performanceUsecase) ResetReworkRate(devUserID uint, requesterRole string) error {
	if requesterRole != "CEO" {
		return fmt.Errorf("access denied: only CEO can reset rework rate (your role: %s)", requesterRole)
	}
	// Verify target user exists
	targetUser, err := u.authRepo.FindByID(devUserID)
	if err != nil || targetUser == nil {
		return fmt.Errorf("user not found")
	}
	return u.authRepo.ResetReworkRate(devUserID)
}
