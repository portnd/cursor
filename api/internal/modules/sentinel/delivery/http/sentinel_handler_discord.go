// Discord notification test handlers for CEO
package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// DiscordTestHandler handles Discord notification testing (CEO only)
type DiscordTestHandler struct {
	usecase         domain.SentinelUsecase
	authRepo        authDomain.Repository
	attendanceRepo  attendanceDomain.AttendanceRepository
	sentinelRepo    domain.SentinelRepository
	discordSvc      domain.DiscordNotifier
}

// NewDiscordTestHandler creates a new Discord test handler
func NewDiscordTestHandler(
	usecase domain.SentinelUsecase,
	authRepo authDomain.Repository,
	attendanceRepo attendanceDomain.AttendanceRepository,
	sentinelRepo domain.SentinelRepository,
	discordSvc domain.DiscordNotifier,
) *DiscordTestHandler {
	return &DiscordTestHandler{
		usecase:        usecase,
		authRepo:       authRepo,
		attendanceRepo: attendanceRepo,
		sentinelRepo:   sentinelRepo,
		discordSvc:     discordSvc,
	}
}

// TestMissingLogNotification handles POST /admin/discord/test-missing-log
// Sends a test notification for users who didn't log time yesterday (CEO only)
func (h *DiscordTestHandler) TestMissingLogNotification(c *gin.Context) {
	// Check CEO role
	userRole := getUserRoleFromContext(c)
	if userRole != authDomain.RoleCEO {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "Only CEO can test Discord notifications",
		})
		return
	}

	// Check if Discord is enabled
	if h.discordSvc == nil || !h.discordSvc.IsEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Discord not configured",
			"message": "DISCORD_WEBHOOK_URL is not set",
		})
		return
	}

	// Get yesterday's date in Bangkok timezone
	loc, _ := time.LoadLocation("Asia/Bangkok")
	if loc == nil {
		loc = time.FixedZone("Bangkok", 7*60*60)
	}

	now := time.Now().In(loc)
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
	yesterdayUTC := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayStr := yesterday.Format("2006-01-02")

	// Get all users
	users, err := h.authRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users",
			"message": err.Error(),
		})
		return
	}

	// Build list of users without logs
	var usersWithoutLogs []domain.UserWithoutLogForDiscord

	for _, user := range users {
		// Skip non-engineer roles
		if !h.shouldCheckUserForLogs(user.Role) {
			continue
		}

		// Check if user has time logs for yesterday
		logs, err := h.sentinelRepo.GetTimeLogsByUserAndDate(user.ID, yesterdayUTC)
		if err != nil {
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

	// Send notification
	if len(usersWithoutLogs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":       "All users logged sufficient time yesterday",
			"date":          yesterdayStr,
			"users_checked": len(users),
		})
		return
	}

	if err := h.discordSvc.SendMissingLogNotification(usersWithoutLogs, yesterdayStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send Discord notification",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Discord notification sent successfully",
		"date":              yesterdayStr,
		"users_without_logs": len(usersWithoutLogs),
		"users":             usersWithoutLogs,
	})
}

// TestLeaveNotification handles POST /admin/discord/test-leave
// Sends a test notification for users on leave today (CEO only)
func (h *DiscordTestHandler) TestLeaveNotification(c *gin.Context) {
	// Check CEO role
	userRole := getUserRoleFromContext(c)
	if userRole != authDomain.RoleCEO {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "Only CEO can test Discord notifications",
		})
		return
	}

	// Check if Discord is enabled
	if h.discordSvc == nil || !h.discordSvc.IsEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Discord not configured",
			"message": "DISCORD_WEBHOOK_URL is not set",
		})
		return
	}

	// Get today's date
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayStr := now.Format("2006-01-02")

	// Get approved leave requests for today
	leaves, err := h.attendanceRepo.ListApprovedLeaveRequestsByDate(today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get leaves",
			"message": err.Error(),
		})
		return
	}

	if len(leaves) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No approved leaves found for today",
			"date":    todayStr,
		})
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
	if err := h.discordSvc.SendLeaveNotification(entries, todayStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send Discord notification",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Discord notification sent successfully",
		"date":          todayStr,
		"leaves_count":  len(leaves),
		"leaves":        entries,
	})
}

// shouldCheckUserForLogs returns true if the user role should be checked for time logs
func (h *DiscordTestHandler) shouldCheckUserForLogs(role string) bool {
	switch role {
	case authDomain.RoleEngineer, authDomain.RoleChiefEngineer:
		return true
	case authDomain.RoleProductOwner, authDomain.RoleManager:
		return true
	case authDomain.RoleCEO, authDomain.RoleSupport:
		return false
	default:
		return false
	}
}
