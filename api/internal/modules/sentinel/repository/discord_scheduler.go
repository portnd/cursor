package repository

import (
	"log"
	"strings"
	"sync"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// DiscordScheduler handles scheduled Discord notifications
type DiscordScheduler struct {
	discordSvc   *DiscordService
	authRepo     authDomain.Repository
	sentinelRepo domain.SentinelRepository
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewDiscordScheduler creates a new scheduler for Discord notifications
func NewDiscordScheduler(discordSvc *DiscordService, authRepo authDomain.Repository, sentinelRepo domain.SentinelRepository) *DiscordScheduler {
	return &DiscordScheduler{
		discordSvc:   discordSvc,
		authRepo:     authRepo,
		sentinelRepo: sentinelRepo,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the scheduler loop
func (s *DiscordScheduler) Start() {
	if s.discordSvc == nil || !s.discordSvc.IsEnabled() {
		log.Println("Discord scheduler: disabled (webhook not configured)")
		return
	}

	s.wg.Add(1)
	go s.run()
	log.Println("✅ Discord notification scheduler started (daily check at 08:00 Bangkok time)")
}

// Stop stops the scheduler
func (s *DiscordScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Println("Discord scheduler stopped")
}

// run is the main scheduler loop
func (s *DiscordScheduler) run() {
	defer s.wg.Done()

	// Calculate next 8:00 AM Bangkok time (UTC+7)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Discord scheduler: failed to load Bangkok timezone: %v", err)
		loc = time.FixedZone("Bangkok", 7*60*60)
	}

	for {
		now := time.Now().In(loc)
		// Calculate next 8:00 AM
		next8AM := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, loc)
		if now.After(next8AM) {
			// Already past 8:00 AM today, schedule for tomorrow
			next8AM = next8AM.AddDate(0, 0, 1)
		}

		durationUntil8AM := next8AM.Sub(now)
		log.Printf("Discord scheduler: next check in %v (at %s)", durationUntil8AM, next8AM.Format("2006-01-02 15:04:05"))

		select {
		case <-time.After(durationUntil8AM):
			s.checkAndNotifyMissingLogs()
		case <-s.stopChan:
			return
		}
	}
}

// checkAndNotifyMissingLogs checks for users who didn't log time yesterday
func (s *DiscordScheduler) checkAndNotifyMissingLogs() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	if loc == nil {
		loc = time.FixedZone("Bangkok", 7*60*60)
	}

	// Get yesterday's date in Bangkok timezone
	now := time.Now().In(loc)
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
	yesterdayUTC := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayStr := yesterday.Format("2006-01-02")

	log.Printf("Discord scheduler: checking for missing logs on %s", yesterdayStr)

	// Get all users
	users, err := s.authRepo.GetAllUsers()
	if err != nil {
		log.Printf("Discord scheduler: failed to get users: %v", err)
		return
	}

	// Build list of users without logs
	var usersWithoutLogs []domain.UserWithoutLogForDiscord

	log.Printf("Discord scheduler: checking %d users for logs", len(users))

	for _, user := range users {
		// Skip non-engineer roles (CEO, MANAGER, PRODUCT_OWNER, SUPPORT typically don't log time)
		if !s.shouldCheckUserForLogs(user.Role) {
			continue
		}

		// Check if user has time logs for yesterday
		logs, err := s.sentinelRepo.GetTimeLogsByUserAndDate(user.ID, yesterdayUTC)
		if err != nil {
			log.Printf("Discord scheduler: failed to get logs for user %d: %v", user.ID, err)
			continue
		}

		// Calculate total minutes logged
		totalMinutes := 0
		for _, log := range logs {
			totalMinutes += log.Minutes
		}
		totalHours := float64(totalMinutes) / 60.0

		// If user has less than 1 hour logged, consider as missing log
		if totalHours < 1.0 {
			// Build display name - prioritize real name
			displayName := strings.TrimSpace(user.FirstName + " " + user.LastName)
			if displayName == "" {
				displayName = user.DisplayName
			}
			if displayName == "" {
				displayName = strings.Split(user.Email, "@")[0]
			}

			usersWithoutLogs = append(usersWithoutLogs, domain.UserWithoutLogForDiscord{
				DisplayName: displayName,
				TotalHours:  totalHours,
			})
		}
	}

	// Send notification if there are users without logs
	if len(usersWithoutLogs) > 0 {
		if err := s.discordSvc.SendMissingLogNotification(usersWithoutLogs, yesterdayStr); err != nil {
			log.Printf("Discord scheduler: failed to send missing log notification: %v", err)
		} else {
			log.Printf("Discord scheduler: sent missing log notification for %d users", len(usersWithoutLogs))
		}
	} else {
		log.Println("Discord scheduler: all users logged sufficient time yesterday")
	}
}

// shouldCheckUserForLogs returns true if the user role should be checked for time logs
func (s *DiscordScheduler) shouldCheckUserForLogs(role string) bool {
	switch role {
	case authDomain.RoleEngineer, authDomain.RoleChiefEngineer:
		return true
	case authDomain.RoleProductOwner, authDomain.RoleManager:
		// Product Owners and Managers may also log time
		return true
	case authDomain.RoleCEO, authDomain.RoleSupport:
		// CEO and Support typically don't log development time
		return false
	default:
		return false
	}
}
