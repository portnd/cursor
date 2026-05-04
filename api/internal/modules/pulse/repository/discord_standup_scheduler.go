package repository

import (
	"log"
	"strings"
	"sync"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
)

// DiscordStandupScheduler sends Discord notifications for missing daily standups.
// Runs every workday (Mon-Fri, excluding company holidays) at 09:30 Bangkok time.
type DiscordStandupScheduler struct {
	discordSvc     sentinelDomain.DiscordNotifier
	pulseRepo      domain.PulseRepository
	authRepo       authDomain.Repository
	attendanceRepo attendanceDomain.AttendanceRepository
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

// NewDiscordStandupScheduler creates a new scheduler for missing standup notifications.
func NewDiscordStandupScheduler(
	discordSvc sentinelDomain.DiscordNotifier,
	pulseRepo domain.PulseRepository,
	authRepo authDomain.Repository,
	attendanceRepo attendanceDomain.AttendanceRepository,
) *DiscordStandupScheduler {
	return &DiscordStandupScheduler{
		discordSvc:     discordSvc,
		pulseRepo:      pulseRepo,
		authRepo:       authRepo,
		attendanceRepo: attendanceRepo,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the scheduler loop.
func (s *DiscordStandupScheduler) Start() {
	if s.discordSvc == nil || !s.discordSvc.IsEnabled() {
		log.Println("Discord standup scheduler: disabled (webhook not configured)")
		return
	}
	s.wg.Add(1)
	go s.run()
	log.Println("✅ Discord standup notification scheduler started (daily check at 09:30 Bangkok time, workdays only)")
}

// Stop stops the scheduler.
func (s *DiscordStandupScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Println("Discord standup scheduler stopped")
}

func (s *DiscordStandupScheduler) run() {
	defer s.wg.Done()
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("Bangkok", 7*60*60)
	}
	for {
		now := time.Now().In(loc)
		next := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, loc)
		if now.After(next) {
			next = next.AddDate(0, 0, 1)
		}
		// Skip weekends
		for next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
			next = next.AddDate(0, 0, 1)
		}
		dur := next.Sub(now)
		log.Printf("Discord standup scheduler: next check in %v (at %s)", dur, next.Format("2006-01-02 15:04:05"))
		select {
		case <-time.After(dur):
			s.checkAndNotify()
		case <-s.stopChan:
			return
		}
	}
}

// isWorkday returns true if the date is a weekday and not a company holiday.
func (s *DiscordStandupScheduler) isWorkday(d time.Time) bool {
	if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		return false
	}
	holidays, err := s.attendanceRepo.ListHolidayCalendars(d, d)
	if err != nil {
		log.Printf("Discord standup scheduler: failed to check holidays: %v", err)
		return true // assume workday if we can't check
	}
	return len(holidays) == 0
}

func (s *DiscordStandupScheduler) checkAndNotify() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	if loc == nil {
		loc = time.FixedZone("Bangkok", 7*60*60)
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayStr := now.Format("2006-01-02")

	if !s.isWorkday(now) {
		log.Printf("Discord standup scheduler: skipping — %s is not a workday", todayStr)
		return
	}

	// Get today's standups
	standups, err := s.pulseRepo.GetStandupsByDate(today)
	if err != nil {
		log.Printf("Discord standup scheduler: failed to get standups: %v", err)
		return
	}
	standupUIDs := make(map[uint]bool, len(standups))
	for _, su := range standups {
		standupUIDs[su.UserID] = true
	}

	// Get approved leaves (exclude from missing list)
	leaves, err := s.pulseRepo.GetApprovedLeavesByDate(today)
	if err != nil {
		log.Printf("Discord standup scheduler: failed to get leaves: %v", err)
	}
	leaveUIDs := make(map[uint]bool, len(leaves))
	for _, l := range leaves {
		leaveUIDs[l.UserID] = true
	}

	// Get all users
	users, err := s.authRepo.GetAllUsers()
	if err != nil {
		log.Printf("Discord standup scheduler: failed to get users: %v", err)
		return
	}

	var missing []sentinelDomain.UserWithoutStandupForDiscord
	for _, user := range users {
		role := strings.ToUpper(strings.TrimSpace(user.Role))
		// CEO and SUPPORT are exempt from standup
		if role == authDomain.RoleCEO || role == authDomain.RoleSupport {
			continue
		}
		// Skip if user already submitted standup
		if standupUIDs[user.ID] {
			continue
		}
		// Skip if user is on approved leave
		if leaveUIDs[user.ID] {
			continue
		}

		displayName := strings.TrimSpace(user.FirstName + " " + user.LastName)
		if displayName == "" {
			displayName = user.DisplayName
		}
		if displayName == "" {
			displayName = strings.Split(user.Email, "@")[0]
		}
		missing = append(missing, sentinelDomain.UserWithoutStandupForDiscord{
			DisplayName: displayName,
		})
	}

	if len(missing) == 0 {
		log.Println("Discord standup scheduler: all users submitted standup today")
		return
	}

	if err := s.discordSvc.SendMissingStandupNotification(missing, todayStr); err != nil {
		log.Printf("Discord standup scheduler: failed to send notification: %v", err)
	} else {
		log.Printf("Discord standup scheduler: sent missing standup notification for %d users", len(missing))
	}
}
