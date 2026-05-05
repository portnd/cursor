package repository

import (
	"errors"
	"sync"
	"testing"
	"time"

	attendanceDomain "github.com/portnd/the-sentinel-core/internal/modules/attendance/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// --- Mocks for DiscordLeaveScheduler ---

type mockDiscordSvcForLeave struct {
	mu        sync.Mutex
	lastLeaves []domain.LeaveEntryForDiscord
	lastDate   string
	enabled    bool
	sendErr    error
	callCount  int
}

func (m *mockDiscordSvcForLeave) IsEnabled() bool { return m.enabled }
func (m *mockDiscordSvcForLeave) SendTimeLogNotification(_ []domain.TimeLogEntryForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvcForLeave) SendMissingLogNotification(_ []domain.UserWithoutLogForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvcForLeave) SendLeaveNotification(leaves []domain.LeaveEntryForDiscord, date string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	m.lastLeaves = leaves
	m.lastDate = date
	return m.sendErr
}
func (m *mockDiscordSvcForLeave) SendMissingStandupNotification(_ []domain.UserWithoutStandupForDiscord, _ string) error {
	return nil
}

type mockAttendanceRepoForLeave struct {
	leaves []attendanceDomain.LeaveRequest
	err    error
}

func (m *mockAttendanceRepoForLeave) ListApprovedLeaveRequestsByDate(_ time.Time) ([]attendanceDomain.LeaveRequest, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.leaves, nil
}

// Stub remaining AttendanceRepository methods
func (m *mockAttendanceRepoForLeave) GetActiveOfficeConfig() (*attendanceDomain.OfficeConfig, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) GetFirstOfficeConfig() (*attendanceDomain.OfficeConfig, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) GetOfficeConfigByID(_ uint) (*attendanceDomain.OfficeConfig, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) CreateOfficeConfig(_ *attendanceDomain.OfficeConfig) error { return nil }
func (m *mockAttendanceRepoForLeave) UpdateOfficeConfig(_ *attendanceDomain.OfficeConfig) error { return nil }
func (m *mockAttendanceRepoForLeave) DeactivateAllOfficeConfigs() error                         { return nil }
func (m *mockAttendanceRepoForLeave) DeactivateOfficeConfigsExcept(_ uint) error                { return nil }
func (m *mockAttendanceRepoForLeave) GetRecordByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.AttendanceRecord, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) SaveRecord(_ *attendanceDomain.AttendanceRecord) error { return nil }
func (m *mockAttendanceRepoForLeave) ListUserRecordsAfterID(_ uint, _ int64, _ int) ([]attendanceDomain.AttendanceRecord, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListRecordsByDate(_ time.Time) ([]attendanceDomain.AttendanceRecord, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) DeleteRecordByID(_ int64) error { return nil }
func (m *mockAttendanceRepoForLeave) CreateOffsiteCheckInRequest(_ *attendanceDomain.OffsiteCheckInRequest) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) GetLatestOffsiteCheckInRequestByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.OffsiteCheckInRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListPendingOffsiteCheckInRequests() ([]attendanceDomain.OffsiteCheckInRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) GetOffsiteCheckInRequestByID(_ int64) (*attendanceDomain.OffsiteCheckInRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) UpdateOffsiteCheckInRequest(_ *attendanceDomain.OffsiteCheckInRequest) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) CreateOffsiteCheckOutRequest(_ *attendanceDomain.OffsiteCheckOutRequest) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) GetLatestOffsiteCheckOutRequestByUserAndDate(_ uint, _ time.Time) (*attendanceDomain.OffsiteCheckOutRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListPendingOffsiteCheckOutRequests() ([]attendanceDomain.OffsiteCheckOutRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) GetOffsiteCheckOutRequestByID(_ int64) (*attendanceDomain.OffsiteCheckOutRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) UpdateOffsiteCheckOutRequest(_ *attendanceDomain.OffsiteCheckOutRequest) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) CreateLeaveRequest(_ *attendanceDomain.LeaveRequest) error { return nil }
func (m *mockAttendanceRepoForLeave) GetLeaveRequestByID(_ int64) (*attendanceDomain.LeaveRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListLeaveRequestsByUser(_ uint) ([]attendanceDomain.LeaveRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListPendingLeaveRequests() ([]attendanceDomain.LeaveRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListAllLeaveRequests() ([]attendanceDomain.LeaveRequest, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) UpdateLeaveRequest(_ *attendanceDomain.LeaveRequest) error { return nil }
func (m *mockAttendanceRepoForLeave) DeleteLeaveRequestByID(_ int64) error                      { return nil }
func (m *mockAttendanceRepoForLeave) ListLeavePolicies() ([]attendanceDomain.LeavePolicy, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) UpsertLeavePolicy(_ *attendanceDomain.LeavePolicy) error { return nil }
func (m *mockAttendanceRepoForLeave) ListHolidayCalendars(_, _ time.Time) ([]attendanceDomain.HolidayCalendar, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) UpsertHolidayCalendar(_ *attendanceDomain.HolidayCalendar) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) CreateLeaveAuditLog(_ *attendanceDomain.LeaveAuditLog) error { return nil }
func (m *mockAttendanceRepoForLeave) ListLeaveAuditLogs(_ int64) ([]attendanceDomain.LeaveAuditLog, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) CreateLeaveNotification(_ *attendanceDomain.LeaveNotification) error {
	return nil
}
func (m *mockAttendanceRepoForLeave) ListLeaveNotifications(_ uint, _ bool) ([]attendanceDomain.LeaveNotification, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) MarkLeaveNotificationRead(_ uint, _ int64) error { return nil }
func (m *mockAttendanceRepoForLeave) GetLeaveTrendByMonth(_ string, _, _ time.Time) ([]attendanceDomain.LeaveTrendPoint, error) {
	return nil, nil
}
func (m *mockAttendanceRepoForLeave) ListAdminApproverUserIDs() ([]uint, error) { return nil, nil }
func (m *mockAttendanceRepoForLeave) FindUserIDByEmail(_ string) (uint, error)  { return 0, nil }
func (m *mockAttendanceRepoForLeave) IsUserRemote(_ uint) (bool, error)          { return false, nil }

// --- Tests ---

func TestCheckAndNotifyLeaves_NoLeaves(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	s := &DiscordLeaveScheduler{
		discordSvc:     discord,
		attendanceRepo: &mockAttendanceRepoForLeave{leaves: nil},
	}
	s.checkAndNotifyLeaves()
	if discord.callCount != 0 {
		t.Error("Should not send notification when no leaves")
	}
}

func TestCheckAndNotifyLeaves_WithLeaves(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	today := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			leaves: []attendanceDomain.LeaveRequest{
				{
					UserID:          1,
					UserDisplayName: "Alice",
					LeaveType:       "ANNUAL",
					StartDate:       today,
					EndDate:         today,
					IsHalfDay:       false,
				},
				{
					UserID:    2,
					UserEmail: "bob@co.com",
					LeaveType: "SICK",
					StartDate: today,
					EndDate:   today,
					IsHalfDay: true,
				},
			},
		},
	}
	s.checkAndNotifyLeaves()
	if discord.callCount != 1 {
		t.Fatal("Expected 1 notification call")
	}
	if len(discord.lastLeaves) != 2 {
		t.Fatalf("Expected 2 leave entries, got %d", len(discord.lastLeaves))
	}
	if discord.lastLeaves[0].DisplayName != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", discord.lastLeaves[0].DisplayName)
	}
	if discord.lastLeaves[1].DisplayName != "bob@co.com" {
		t.Errorf("Expected 'bob@co.com' fallback, got '%s'", discord.lastLeaves[1].DisplayName)
	}
	if !discord.lastLeaves[1].IsHalfDay {
		t.Error("Expected IsHalfDay=true for Bob")
	}
}

func TestCheckAndNotifyLeaves_FallbackDisplayName(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	today := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			leaves: []attendanceDomain.LeaveRequest{
				{
					UserID:    1,
					UserEmail: "john@co.com",
					LeaveType: "PERSONAL",
					StartDate: today,
					EndDate:   today,
				},
			},
		},
	}
	s.checkAndNotifyLeaves()
	if discord.lastLeaves[0].DisplayName != "john@co.com" {
		t.Errorf("Expected 'john@co.com' from email fallback, got '%s'", discord.lastLeaves[0].DisplayName)
	}
}

func TestCheckAndNotifyLeaves_LeaveTypes(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	today := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			leaves: []attendanceDomain.LeaveRequest{
				{UserID: 1, UserDisplayName: "A", LeaveType: "ANNUAL", StartDate: today, EndDate: today},
				{UserID: 2, UserDisplayName: "B", LeaveType: "SICK", StartDate: today, EndDate: today},
				{UserID: 3, UserDisplayName: "C", LeaveType: "PERSONAL", StartDate: today, EndDate: today},
				{UserID: 4, UserDisplayName: "D", LeaveType: "UNPAID", StartDate: today, EndDate: today},
			},
		},
	}
	s.checkAndNotifyLeaves()
	if len(discord.lastLeaves) != 4 {
		t.Fatalf("Expected 4 leave entries, got %d", len(discord.lastLeaves))
	}
	expectedTypes := []string{"ANNUAL", "SICK", "PERSONAL", "UNPAID"}
	for i, et := range expectedTypes {
		if discord.lastLeaves[i].LeaveType != et {
			t.Errorf("Entry %d: expected LeaveType '%s', got '%s'", i, et, discord.lastLeaves[i].LeaveType)
		}
	}
}

func TestCheckAndNotifyLeaves_GetLeavesFails(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			err: errors.New("db error"),
		},
	}
	s.checkAndNotifyLeaves()
	if discord.callCount != 0 {
		t.Error("Should not send notification when ListApprovedLeaveRequestsByDate fails")
	}
}

func TestCheckAndNotifyLeaves_SendFails(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true, sendErr: errors.New("webhook error")}
	today := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			leaves: []attendanceDomain.LeaveRequest{
				{UserID: 1, UserDisplayName: "Alice", LeaveType: "ANNUAL", StartDate: today, EndDate: today},
			},
		},
	}
	s.checkAndNotifyLeaves()
	if discord.callCount != 1 {
		t.Fatal("Expected notification attempt even if send fails")
	}
}

func TestCheckAndNotifyLeaves_DiscordDisabled(t *testing.T) {
	s := &DiscordLeaveScheduler{
		discordSvc:     &mockDiscordSvcForLeave{enabled: false},
		attendanceRepo: &mockAttendanceRepoForLeave{},
	}
	s.Start()
	time.Sleep(50 * time.Millisecond)
}

func TestCheckAndNotifyLeaves_HalfDay(t *testing.T) {
	discord := &mockDiscordSvcForLeave{enabled: true}
	today := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	s := &DiscordLeaveScheduler{
		discordSvc: discord,
		attendanceRepo: &mockAttendanceRepoForLeave{
			leaves: []attendanceDomain.LeaveRequest{
				{
					UserID:          1,
					UserDisplayName: "Alice",
					LeaveType:       "SICK",
					StartDate:       today,
					EndDate:         today,
					IsHalfDay:       true,
				},
			},
		},
	}
	s.checkAndNotifyLeaves()
	if !discord.lastLeaves[0].IsHalfDay {
		t.Error("Expected IsHalfDay=true")
	}
}
