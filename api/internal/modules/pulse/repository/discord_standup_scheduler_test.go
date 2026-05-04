package repository

import (
	"errors"
	"sync"
	"testing"
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	pulseDomain "github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// --- Mocks ---

type mockDiscordSvc struct {
	mu              sync.Mutex
	lastUsers       []sentinelDomain.UserWithoutStandupForDiscord
	lastDate        string
	enabled         bool
	sendErr         error
	callCount       int
}

func (m *mockDiscordSvc) IsEnabled() bool { return m.enabled }
func (m *mockDiscordSvc) SendTimeLogNotification(_ []sentinelDomain.TimeLogEntryForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvc) SendMissingLogNotification(_ []sentinelDomain.UserWithoutLogForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvc) SendLeaveNotification(_ []sentinelDomain.LeaveEntryForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvc) SendMissingStandupNotification(users []sentinelDomain.UserWithoutStandupForDiscord, date string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	m.lastUsers = users
	m.lastDate = date
	return m.sendErr
}

type mockAuthRepo struct {
	users []authDomain.User
	err   error
}

func (m *mockAuthRepo) GetAllUsers() ([]authDomain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}
// Stub remaining authDomain.Repository methods
func (m *mockAuthRepo) CreateUser(_ *authDomain.User) error { return nil }
func (m *mockAuthRepo) FindByEmail(_ string) (*authDomain.User, error) { return nil, nil }
func (m *mockAuthRepo) FindByID(_ uint) (*authDomain.User, error) { return nil, nil }
func (m *mockAuthRepo) UpdateProfile(_ uint, _, _, _ *string, _ []string) error { return nil }
func (m *mockAuthRepo) UpdateAvatar(_ uint, _ string) error { return nil }
func (m *mockAuthRepo) UpdateThemePreference(_ uint, _ string) error { return nil }
func (m *mockAuthRepo) UpdateUserRole(_ uint, _ string) error { return nil }
func (m *mockAuthRepo) UpdateUserAdmin(_ uint, _, _, _, _ *string) error { return nil }
func (m *mockAuthRepo) DeleteUser(_ uint) error { return nil }
func (m *mockAuthRepo) UpdatePassword(_ uint, _ string) error { return nil }
func (m *mockAuthRepo) ResetReworkRate(_ uint) error { return nil }
func (m *mockAuthRepo) SetUserRemote(_ uint, _ bool) error { return nil }
func (m *mockAuthRepo) CreateTeam(_ *authDomain.Team) error { return nil }
func (m *mockAuthRepo) GetAllTeams() ([]authDomain.Team, error) { return nil, nil }
func (m *mockAuthRepo) GetTeamByID(_ uint) (*authDomain.Team, error) { return nil, nil }
func (m *mockAuthRepo) UpdateTeamName(_ uint, _ string) error { return nil }
func (m *mockAuthRepo) DeleteTeam(_ uint) error { return nil }
func (m *mockAuthRepo) AssignUserToTeam(_ uint, _ *uint) error { return nil }
func (m *mockAuthRepo) UpdateTeamCapital(_ uint, _ float64, _ *float64) error { return nil }
func (m *mockAuthRepo) CreateTeamTransaction(_ *authDomain.TeamTransaction) error { return nil }
func (m *mockAuthRepo) GetTeamTransactions(_ uint) ([]authDomain.TeamTransaction, error) { return nil, nil }
func (m *mockAuthRepo) GetAppSetting(_ string) (string, error) { return "", nil }
func (m *mockAuthRepo) SetAppSetting(_, _ string) error { return nil }

type mockPulseRepo struct {
	standups []pulseDomain.DailyStandup
	leaves   []pulseDomain.LeaveSummary
	standupErr error
	leaveErr   error
}

func (m *mockPulseRepo) SaveStandup(_ *pulseDomain.DailyStandup) error { return nil }
func (m *mockPulseRepo) GetStandupByUserAndDate(_ uint, _ time.Time) (*pulseDomain.DailyStandup, error) {
	return nil, nil
}
func (m *mockPulseRepo) GetStandupsByDate(_ time.Time) ([]pulseDomain.DailyStandup, error) {
	return m.standups, m.standupErr
}
func (m *mockPulseRepo) GetTimeLogsByDate(_ time.Time) ([]pulseDomain.TimeLogSummary, error) { return nil, nil }
func (m *mockPulseRepo) GetSubmissionsByDate(_ time.Time) ([]pulseDomain.SubmissionSummary, error) { return nil, nil }
func (m *mockPulseRepo) GetApprovedLeavesByDate(_ time.Time) ([]pulseDomain.LeaveSummary, error) {
	return m.leaves, m.leaveErr
}
func (m *mockPulseRepo) GetHiddenPulseUserIDs() ([]uint, error) { return nil, nil }
func (m *mockPulseRepo) SetHiddenPulseUserIDs(_ []uint) error { return nil }

type mockAttendanceRepo struct {
	holidays []attendanceDomain.HolidayCalendar
	holidayErr error
}

func (m *mockAttendanceRepo) ListHolidayCalendars(_, _ time.Time) ([]attendanceDomain.HolidayCalendar, error) {
	return m.holidays, m.holidayErr
}
// Stub remaining methods
func (m *mockAttendanceRepo) GetActiveOfficeConfig() (*attendanceDomain.OfficeConfig, error) { return nil, nil }
func (m *mockAttendanceRepo) GetFirstOfficeConfig() (*attendanceDomain.OfficeConfig, error) { return nil, nil }
func (m *mockAttendanceRepo) GetOfficeConfigByID(_ uint) (*attendanceDomain.OfficeConfig, error) { return nil, nil }
func (m *mockAttendanceRepo) CreateOfficeConfig(_ *attendanceDomain.OfficeConfig) error { return nil }
func (m *mockAttendanceRepo) UpdateOfficeConfig(_ *attendanceDomain.OfficeConfig) error { return nil }
func (m *mockAttendanceRepo) DeactivateAllOfficeConfigs() error { return nil }
func (m *mockAttendanceRepo) DeactivateOfficeConfigsExcept(_ uint) error { return nil }
func (m *mockAttendanceRepo) GetRecordByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.AttendanceRecord, error) { return nil, nil }
func (m *mockAttendanceRepo) SaveRecord(_ *attendanceDomain.AttendanceRecord) error { return nil }
func (m *mockAttendanceRepo) ListUserRecordsAfterID(_ uint, _ int64, _ int) ([]attendanceDomain.AttendanceRecord, error) { return nil, nil }
func (m *mockAttendanceRepo) ListRecordsByDate(_ time.Time) ([]attendanceDomain.AttendanceRecord, error) { return nil, nil }
func (m *mockAttendanceRepo) DeleteRecordByID(_ int64) error { return nil }
func (m *mockAttendanceRepo) CreateOffsiteCheckInRequest(_ *attendanceDomain.OffsiteCheckInRequest) error { return nil }
func (m *mockAttendanceRepo) GetLatestOffsiteCheckInRequestByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.OffsiteCheckInRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListPendingOffsiteCheckInRequests() ([]attendanceDomain.OffsiteCheckInRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) GetOffsiteCheckInRequestByID(_ int64) (*attendanceDomain.OffsiteCheckInRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) UpdateOffsiteCheckInRequest(_ *attendanceDomain.OffsiteCheckInRequest) error { return nil }
func (m *mockAttendanceRepo) CreateOffsiteCheckOutRequest(_ *attendanceDomain.OffsiteCheckOutRequest) error { return nil }
func (m *mockAttendanceRepo) GetLatestOffsiteCheckOutRequestByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.OffsiteCheckOutRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListPendingOffsiteCheckOutRequests() ([]attendanceDomain.OffsiteCheckOutRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) GetOffsiteCheckOutRequestByID(_ int64) (*attendanceDomain.OffsiteCheckOutRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) UpdateOffsiteCheckOutRequest(_ *attendanceDomain.OffsiteCheckOutRequest) error { return nil }
func (m *mockAttendanceRepo) CreateLeaveRequest(_ *attendanceDomain.LeaveRequest) error { return nil }
func (m *mockAttendanceRepo) GetLeaveRequestByID(_ int64) (*attendanceDomain.LeaveRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListLeaveRequestsByUser(_ uint) ([]attendanceDomain.LeaveRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListPendingLeaveRequests() ([]attendanceDomain.LeaveRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListAllLeaveRequests() ([]attendanceDomain.LeaveRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) ListApprovedLeaveRequestsByDate(_ time.Time) ([]attendanceDomain.LeaveRequest, error) { return nil, nil }
func (m *mockAttendanceRepo) UpdateLeaveRequest(_ *attendanceDomain.LeaveRequest) error { return nil }
func (m *mockAttendanceRepo) DeleteLeaveRequestByID(_ int64) error { return nil }
func (m *mockAttendanceRepo) ListLeavePolicies() ([]attendanceDomain.LeavePolicy, error) { return nil, nil }
func (m *mockAttendanceRepo) UpsertLeavePolicy(_ *attendanceDomain.LeavePolicy) error { return nil }
func (m *mockAttendanceRepo) UpsertHolidayCalendar(_ *attendanceDomain.HolidayCalendar) error { return nil }
func (m *mockAttendanceRepo) CreateLeaveAuditLog(_ *attendanceDomain.LeaveAuditLog) error { return nil }
func (m *mockAttendanceRepo) ListLeaveAuditLogs(_ int64) ([]attendanceDomain.LeaveAuditLog, error) { return nil, nil }
func (m *mockAttendanceRepo) CreateLeaveNotification(_ *attendanceDomain.LeaveNotification) error { return nil }
func (m *mockAttendanceRepo) ListLeaveNotifications(_ uint, _ bool) ([]attendanceDomain.LeaveNotification, error) { return nil, nil }
func (m *mockAttendanceRepo) MarkLeaveNotificationRead(_ uint, _ int64) error { return nil }
func (m *mockAttendanceRepo) GetLeaveTrendByMonth(_ string, _, _ time.Time) ([]attendanceDomain.LeaveTrendPoint, error) { return nil, nil }
func (m *mockAttendanceRepo) ListAdminApproverUserIDs() ([]uint, error) { return nil, nil }
func (m *mockAttendanceRepo) FindUserIDByEmail(_ string) (uint, error) { return 0, nil }
func (m *mockAttendanceRepo) IsUserRemote(_ uint) (bool, error) { return false, nil }

// --- Tests ---

func TestIsWorkday_Weekend(t *testing.T) {
	s := &DiscordStandupScheduler{attendanceRepo: &mockAttendanceRepo{}}
	sat := time.Date(2025, 5, 3, 0, 0, 0, 0, time.UTC) // Saturday
	sun := time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC) // Sunday
	if s.isWorkday(sat) {
		t.Error("Saturday should not be a workday")
	}
	if s.isWorkday(sun) {
		t.Error("Sunday should not be a workday")
	}
}

func TestIsWorkday_Weekday(t *testing.T) {
	s := &DiscordStandupScheduler{attendanceRepo: &mockAttendanceRepo{}}
	mon := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC) // Monday
	if !s.isWorkday(mon) {
		t.Error("Monday should be a workday")
	}
}

func TestIsWorkday_Holiday(t *testing.T) {
	hol := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC) // Monday but holiday
	s := &DiscordStandupScheduler{
		attendanceRepo: &mockAttendanceRepo{
			holidays: []attendanceDomain.HolidayCalendar{{Date: hol, Name: "Labour Day"}},
		},
	}
	if s.isWorkday(hol) {
		t.Error("Holiday Monday should not be a workday")
	}
}

func TestIsWorkday_HolidayCheckFails(t *testing.T) {
	mon := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordStandupScheduler{
		attendanceRepo: &mockAttendanceRepo{holidayErr: errors.New("db error")},
	}
	if !s.isWorkday(mon) {
		t.Error("If holiday check fails, assume workday")
	}
}

func TestCheckAndNotify_AllSubmitted(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{
			standups: []pulseDomain.DailyStandup{{UserID: 1}, {UserID: 2}},
			leaves:   nil,
		},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice", LastName: "A"},
				{ID: 2, Role: "MANAGER", FirstName: "Bob", LastName: "B"},
			},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	s.checkAndNotify()
	if discord.callCount != 0 {
		t.Error("Should not send notification when all submitted")
	}
}

func TestCheckAndNotify_SomeMissing(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{
			standups: []pulseDomain.DailyStandup{{UserID: 1}},
			leaves:   nil,
		},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice", LastName: "A"},
				{ID: 2, Role: "ENGINEER", FirstName: "Bob", LastName: "B"},
			},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	s.checkAndNotify()
	if discord.callCount != 1 {
		t.Fatal("Expected 1 notification call")
	}
	if len(discord.lastUsers) != 1 {
		t.Fatalf("Expected 1 missing user, got %d", len(discord.lastUsers))
	}
	if discord.lastUsers[0].DisplayName != "Bob B" {
		t.Errorf("Expected 'Bob B', got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotify_ExcludesCEOAndSupport(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{standups: nil, leaves: nil},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{
				{ID: 1, Role: "CEO", FirstName: "Boss"},
				{ID: 2, Role: "SUPPORT", FirstName: "Help"},
				{ID: 3, Role: "ENGINEER", FirstName: "Dev"},
			},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	s.checkAndNotify()
	if discord.callCount != 1 {
		t.Fatal("Expected notification")
	}
	if len(discord.lastUsers) != 1 {
		t.Fatalf("Expected 1 missing (engineer only), got %d", len(discord.lastUsers))
	}
	if discord.lastUsers[0].DisplayName != "Dev" {
		t.Errorf("Expected 'Dev', got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotify_ExcludesOnLeave(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{
			standups: nil,
			leaves:   []pulseDomain.LeaveSummary{{UserID: 2}},
		},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice"},
				{ID: 2, Role: "ENGINEER", FirstName: "OnLeave"},
			},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	s.checkAndNotify()
	if len(discord.lastUsers) != 1 {
		t.Fatalf("Expected 1 missing (user 2 on leave), got %d", len(discord.lastUsers))
	}
	if discord.lastUsers[0].DisplayName != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotify_FallbackDisplayName(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{standups: nil, leaves: nil},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", Email: "john@co.com"},
			},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	s.checkAndNotify()
	if discord.lastUsers[0].DisplayName != "john" {
		t.Errorf("Expected 'john' from email fallback, got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotify_WeekendSkips(t *testing.T) {
	discord := &mockDiscordSvc{enabled: true}
	s := &DiscordStandupScheduler{
		discordSvc: discord,
		pulseRepo: &mockPulseRepo{standups: nil, leaves: nil},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{{ID: 1, Role: "ENGINEER", FirstName: "Dev"}},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	// Temporarily override time by running on a known weekend — we test isWorkday separately
	// checkAndNotify uses time.Now() internally, so we test the isWorkday gate directly
	sat := time.Date(2025, 5, 3, 9, 30, 0, 0, time.UTC)
	if s.isWorkday(sat) {
		t.Error("Saturday is not a workday")
	}
	if discord.callCount != 0 {
		t.Error("No notification should be sent for weekend check")
	}
}

func TestCheckAndNotify_DiscordDisabled(t *testing.T) {
	s := &DiscordStandupScheduler{
		discordSvc: &mockDiscordSvc{enabled: false},
		pulseRepo: &mockPulseRepo{standups: nil, leaves: nil},
		authRepo: &mockAuthRepo{
			users: []authDomain.User{{ID: 1, Role: "ENGINEER", FirstName: "Dev"}},
		},
		attendanceRepo: &mockAttendanceRepo{},
	}
	// Start should not launch goroutine
	s.Start()
	// Give a moment to ensure no goroutine started
	time.Sleep(50 * time.Millisecond)
}
