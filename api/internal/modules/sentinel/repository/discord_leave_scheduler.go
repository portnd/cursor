package repository

import (
	"log"
	"sync"
	"time"

	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// DiscordLeaveScheduler handles scheduled Discord notifications for leave
type DiscordLeaveScheduler struct {
	discordSvc     *DiscordService
	attendanceRepo attendanceDomain.AttendanceRepository
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

// NewDiscordLeaveScheduler creates a new scheduler for leave notifications
func NewDiscordLeaveScheduler(discordSvc *DiscordService, attendanceRepo attendanceDomain.AttendanceRepository) *DiscordLeaveScheduler {
	return &DiscordLeaveScheduler{
		discordSvc:     discordSvc,
		attendanceRepo: attendanceRepo,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the scheduler loop
func (s *DiscordLeaveScheduler) Start() {
	if s.discordSvc == nil || !s.discordSvc.IsEnabled() {
		log.Println("Discord leave scheduler: disabled (webhook not configured)")
		return
	}

	s.wg.Add(1)
	go s.run()
	log.Println("✅ Discord leave notification scheduler started (daily check at 08:30 Bangkok time)")
}

// Stop stops the scheduler
func (s *DiscordLeaveScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Println("Discord leave scheduler stopped")
}

// run is the main scheduler loop
func (s *DiscordLeaveScheduler) run() {
	defer s.wg.Done()

	// Calculate next 8:30 AM Bangkok time (UTC+7)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Discord leave scheduler: failed to load Bangkok timezone: %v", err)
		loc = time.FixedZone("Bangkok", 7*60*60)
	}

	for {
		now := time.Now().In(loc)
		// Calculate next 8:30 AM
		next830AM := time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, loc)
		if now.After(next830AM) {
			// Already past 8:30 AM today, schedule for tomorrow
			next830AM = next830AM.AddDate(0, 0, 1)
		}

		durationUntil830AM := next830AM.Sub(now)
		log.Printf("Discord leave scheduler: next check in %v (at %s)", durationUntil830AM, next830AM.Format("2006-01-02 15:04:05"))

		select {
		case <-time.After(durationUntil830AM):
			s.checkAndNotifyLeaves()
		case <-s.stopChan:
			return
		}
	}
}

// checkAndNotifyLeaves checks for employees on leave today and sends notification
func (s *DiscordLeaveScheduler) checkAndNotifyLeaves() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	if loc == nil {
		loc = time.FixedZone("Bangkok", 7*60*60)
	}

	// Get today's date in Bangkok timezone
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayStr := now.Format("2006-01-02")

	log.Printf("Discord leave scheduler: checking for leaves on %s", todayStr)

	// Get approved leave requests for today
	leaves, err := s.attendanceRepo.ListApprovedLeaveRequestsByDate(today)
	if err != nil {
		log.Printf("Discord leave scheduler: failed to get leaves: %v", err)
		return
	}

	if len(leaves) == 0 {
		log.Println("Discord leave scheduler: no leaves today")
		return
	}

	// Build notification entries
	var entries []domain.LeaveEntryForDiscord
	for _, l := range leaves {
		displayName := l.UserDisplayName
		if displayName == "" {
			displayName = l.UserEmail
		}

		entries = append(entries, domain.LeaveEntryForDiscord{
			DisplayName: displayName,
			LeaveType:   l.LeaveType,
			StartDate:   l.StartDate.Format("2006-01-02"),
			EndDate:     l.EndDate.Format("2006-01-02"),
			IsHalfDay:   l.IsHalfDay,
		})
	}

	// Send notification
	if err := s.discordSvc.SendLeaveNotification(entries, todayStr); err != nil {
		log.Printf("Discord leave scheduler: failed to send notification: %v", err)
	} else {
		log.Printf("Discord leave scheduler: sent leave notification for %d employees", len(leaves))
	}
}
