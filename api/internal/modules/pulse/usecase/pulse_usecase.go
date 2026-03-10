package usecase

import (
	"fmt"
	"math"
	"sort"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
)

type pulseUsecase struct {
	repo     domain.PulseRepository
	authRepo authDomain.Repository
}

// NewPulseUsecase wires a PulseUsecase with its dependencies.
func NewPulseUsecase(repo domain.PulseRepository, authRepo authDomain.Repository) domain.PulseUsecase {
	return &pulseUsecase{repo: repo, authRepo: authRepo}
}

// SubmitStandup creates or updates the calling user's daily standup.
func (u *pulseUsecase) SubmitStandup(
	userID uint,
	date time.Time,
	yesterday, blocker string,
	todayTaskIDs []string,
) (*domain.DailyStandup, error) {
	// Normalise to UTC date (strip time component)
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	standup := &domain.DailyStandup{
		UserID:           userID,
		Date:             d,
		YesterdaySummary: yesterday,
		TodayTaskIDs:     todayTaskIDs,
		Blocker:          blocker,
	}

	if err := u.repo.SaveStandup(standup); err != nil {
		return nil, fmt.Errorf("pulse: submit standup: %w", err)
	}

	// Enrich with user info for the response
	user, err := u.authRepo.FindByID(userID)
	if err == nil && user != nil {
		standup.UserEmail = user.Email
		standup.UserDisplayName = user.DisplayName
	}

	return standup, nil
}

// GetDailyCompanyPulse aggregates standups, time-logs, and submissions for a date
// into a structured per-user board.
func (u *pulseUsecase) GetDailyCompanyPulse(date time.Time) (*domain.CompanyPulseResponse, error) {
	// ── Fetch all source data in parallel-friendly sequence ────────────────────
	standups, err := u.repo.GetStandupsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("pulse: fetch standups: %w", err)
	}

	timeLogs, err := u.repo.GetTimeLogsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("pulse: fetch time logs: %w", err)
	}

	submissions, err := u.repo.GetSubmissionsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("pulse: fetch submissions: %w", err)
	}

	// ── Load all team members to ensure every user appears ────────────────────
	users, err := u.authRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("pulse: fetch users: %w", err)
	}

	// ── Build lookup maps ─────────────────────────────────────────────────────
	standupMap := make(map[uint]*domain.DailyStandup, len(standups))
	for i := range standups {
		s := &standups[i]
		user, err2 := u.authRepo.FindByID(s.UserID)
		if err2 == nil && user != nil {
			s.UserEmail = user.Email
			s.UserDisplayName = user.DisplayName
		}
		standupMap[s.UserID] = s
	}

	// Group time-logs by userID
	type timeAccum struct {
		totalMinutes int
		items        []domain.ActivityItem
	}
	timeMap := make(map[uint]*timeAccum)
	for _, tl := range timeLogs {
		acc, ok := timeMap[tl.UserID]
		if !ok {
			acc = &timeAccum{}
			timeMap[tl.UserID] = acc
		}
		acc.totalMinutes += tl.Minutes
		acc.items = append(acc.items, domain.ActivityItem{
			Type:        "time_log",
			Description: tl.Description,
			Minutes:     tl.Minutes,
			OccurredAt:  tl.LoggedAt,
		})
	}

	// Group submissions by devID
	subMap := make(map[uint][]domain.ActivityItem)
	for _, sub := range submissions {
		desc := "Handover submitted"
		if sub.ReferenceURL != "" {
			desc = sub.ReferenceURL
		}
		subMap[sub.DevID] = append(subMap[sub.DevID], domain.ActivityItem{
			Type:         "submission",
			Description:  desc,
			ReferenceURL: sub.ReferenceURL,
			OccurredAt:   sub.CreatedAt,
		})
	}

	// ── Build UserPulse entries ────────────────────────────────────────────────
	const maxActivities = 5

	memberMap := make(map[uint]*domain.UserPulse, len(users))
	for _, usr := range users {
		p := &domain.UserPulse{
			UserID:          usr.ID,
			UserEmail:       usr.Email,
			UserDisplayName: usr.DisplayName,
			Standup:         standupMap[usr.ID],
			HasBlocker:      standupMap[usr.ID] != nil && standupMap[usr.ID].Blocker != "",
		}

		// Aggregate time logs
		if acc, ok := timeMap[usr.ID]; ok {
			p.TotalLoggedMin = acc.totalMinutes
			p.TotalLoggedHrs = math.Round(float64(acc.totalMinutes)/60*100) / 100
		}

		// Merge and sort activities (time_logs + submissions) newest-first
		var activities []domain.ActivityItem
		if acc, ok := timeMap[usr.ID]; ok {
			activities = append(activities, acc.items...)
		}
		activities = append(activities, subMap[usr.ID]...)
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].OccurredAt.After(activities[j].OccurredAt)
		})
		if len(activities) > maxActivities {
			activities = activities[:maxActivities]
		}
		p.LatestActivities = activities

		memberMap[usr.ID] = p
	}

	// ── Assemble response ─────────────────────────────────────────────────────
	members := make([]domain.UserPulse, 0, len(memberMap))
	for _, p := range memberMap {
		members = append(members, *p)
	}
	// Sort: blockers first, then by user_id for stable order
	sort.Slice(members, func(i, j int) bool {
		if members[i].HasBlocker != members[j].HasBlocker {
			return members[i].HasBlocker
		}
		return members[i].UserID < members[j].UserID
	})

	totalLoggedMin := 0
	checkedIn := 0
	for _, m := range members {
		totalLoggedMin += m.TotalLoggedMin
		if m.Standup != nil {
			checkedIn++
		}
	}

	return &domain.CompanyPulseResponse{
		Date:           date.Format("2006-01-02"),
		TotalMembers:   len(members),
		CheckedIn:      checkedIn,
		TotalMinLogged: totalLoggedMin,
		Members:        members,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
