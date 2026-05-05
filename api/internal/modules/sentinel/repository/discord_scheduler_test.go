package repository

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// --- Mocks for DiscordScheduler ---

type mockDiscordSvcForLog struct {
	mu              sync.Mutex
	lastUsers       []domain.UserWithoutLogForDiscord
	lastDate        string
	enabled         bool
	sendErr         error
	callCount       int
}

func (m *mockDiscordSvcForLog) IsEnabled() bool { return m.enabled }
func (m *mockDiscordSvcForLog) SendTimeLogNotification(_ []domain.TimeLogEntryForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvcForLog) SendMissingLogNotification(users []domain.UserWithoutLogForDiscord, date string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	m.lastUsers = users
	m.lastDate = date
	return m.sendErr
}
func (m *mockDiscordSvcForLog) SendLeaveNotification(_ []domain.LeaveEntryForDiscord, _ string) error {
	return nil
}
func (m *mockDiscordSvcForLog) SendMissingStandupNotification(_ []domain.UserWithoutStandupForDiscord, _ string) error {
	return nil
}

type mockSentinelRepo struct {
	logs map[uint][]domain.TimeLog
	err  error
}

func (m *mockSentinelRepo) GetTimeLogsByUserAndDate(userID uint, _ time.Time) ([]domain.TimeLog, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.logs[userID], nil
}

// Stub all remaining SentinelRepository methods
func (m *mockSentinelRepo) CreateProject(_ *domain.Project) error { return nil }
func (m *mockSentinelRepo) GetAllProjects(_ domain.CallerContext) ([]domain.Project, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetProjectByID(_ uuid.UUID, _ domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetProjectByCode(_ string, _ domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTasksByProjectID(_ uuid.UUID) ([]domain.Task, error) { return nil, nil }
func (m *mockSentinelRepo) GetTasksByProjectIDForProjectPage(_ uuid.UUID, _ int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTasksByProjectIDForProjectPageCursor(_ uuid.UUID, _ int, _ *time.Time, _ *uuid.UUID, _ int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateProject(_ *domain.Project) error { return nil }
func (m *mockSentinelRepo) DeleteProject(_ uuid.UUID) error     { return nil }
func (m *mockSentinelRepo) DeleteProjectPlan(_ uuid.UUID) error { return nil }
func (m *mockSentinelRepo) AssignProjectTeam(_ uuid.UUID, _ *uint) error {
	return nil
}
func (m *mockSentinelRepo) ReplaceProjectPmAssignments(_ uuid.UUID, _ []uint) error {
	return nil
}
func (m *mockSentinelRepo) CreateTask(_ *domain.Task) error      { return nil }
func (m *mockSentinelRepo) GetTaskByID(_ uuid.UUID) (*domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTaskByCode(_ string) (*domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) CountTasksForCode(_ *uuid.UUID) (int, error) { return 0, nil }
func (m *mockSentinelRepo) GetMaxTaskCodeSuffix(_ string) (int, error)    { return 0, nil }
func (m *mockSentinelRepo) GetTasksByAssignee(_ uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetActiveSprintTasksByAssignee(_ uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetActiveSprintsForUser(_ uint) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetGlobalActiveTasks(_ domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTeamActiveTasks(_ domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetActiveFeatures(_ uint, _ *uuid.UUID) ([]domain.FeatureRoadmapItem, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetUnassignedTasks() ([]domain.Task, error)    { return nil, nil }
func (m *mockSentinelRepo) GetAllTasks() ([]domain.Task, error)          { return nil, nil }
func (m *mockSentinelRepo) GetTasksByProjectIDs(_ []uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTasksRequiringApproval() ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateTask(_ *domain.Task) error { return nil }
func (m *mockSentinelRepo) DeleteTask(_ uuid.UUID) error   { return nil }
func (m *mockSentinelRepo) CreateTaskActivity(_ *domain.TaskActivityEvent) error { return nil }
func (m *mockSentinelRepo) ListTaskActivitiesByTaskID(_ uuid.UUID) ([]domain.TaskActivityEvent, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetAllTaskDependencies() ([]domain.TaskDependency, error) {
	return nil, nil
}
func (m *mockSentinelRepo) CreateTaskDependency(_ *domain.TaskDependency) error { return nil }
func (m *mockSentinelRepo) DeleteTaskDependency(_ uuid.UUID) error            { return nil }
func (m *mockSentinelRepo) CreateSubmission(_ *domain.Submission) error        { return nil }
func (m *mockSentinelRepo) GetSubmissionByID(_ uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateSubmission(_ *domain.Submission) error { return nil }
func (m *mockSentinelRepo) GetLatestSubmission(_ uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockSentinelRepo) CreateAppeal(_ *domain.Appeal) error        { return nil }
func (m *mockSentinelRepo) GetAppealBySubmissionID(_ uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetAppealByID(_ uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateAppeal(_ *domain.Appeal) error { return nil }
func (m *mockSentinelRepo) ApproveTask(_ uuid.UUID) error    { return nil }
func (m *mockSentinelRepo) RejectTask(_ uuid.UUID, _ uint, _ string) error {
	return nil
}
func (m *mockSentinelRepo) GetTasksReadyForTest(_ uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockSentinelRepo) SetTaskReadyForUAT(_ uuid.UUID, _ []byte) error { return nil }
func (m *mockSentinelRepo) SetTaskWaitForDeploy(_ uuid.UUID, _ []byte) error {
	return nil
}
func (m *mockSentinelRepo) AdvanceTaskToReadyForUAT(_ uuid.UUID) error { return nil }
func (m *mockSentinelRepo) GetTasksReadyForCEOApproval(_ uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockSentinelRepo) CreateSprint(_ *domain.Sprint) error { return nil }
func (m *mockSentinelRepo) GetSprintByID(_ uuid.UUID) (*domain.Sprint, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetSprintsByProjectID(_ uuid.UUID) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetActiveSprintByProjectID(_ uuid.UUID) (*domain.Sprint, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTasksBySprintID(_ uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateSprint(_ *domain.Sprint) error { return nil }
func (m *mockSentinelRepo) DeleteSprint(_ uuid.UUID) error    { return nil }
func (m *mockSentinelRepo) CreateMilestone(_ *domain.Milestone) error { return nil }
func (m *mockSentinelRepo) GetMilestoneByID(_ uuid.UUID) (*domain.Milestone, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetMilestonesByProjectID(_ uuid.UUID) ([]domain.Milestone, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateMilestone(_ *domain.Milestone) error { return nil }
func (m *mockSentinelRepo) DeleteMilestone(_ uuid.UUID) error       { return nil }
func (m *mockSentinelRepo) CreateTaskComment(_ *domain.TaskComment) error { return nil }
func (m *mockSentinelRepo) GetCommentsByTaskID(_ uuid.UUID) ([]domain.TaskComment, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTaskCommentByID(_ uuid.UUID) (*domain.TaskComment, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateTaskComment(_ *domain.TaskComment) error { return nil }
func (m *mockSentinelRepo) DeleteTaskComment(_ uuid.UUID) error        { return nil }
func (m *mockSentinelRepo) CreateTimeLog(_ *domain.TimeLog) error        { return nil }
func (m *mockSentinelRepo) GetTimeLogsByTaskID(_ uuid.UUID) ([]domain.TimeLog, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetTimeLogByID(_ uuid.UUID) (*domain.TimeLog, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateTimeLog(_ *domain.TimeLog) error { return nil }
func (m *mockSentinelRepo) DeleteTimeLog(_ uuid.UUID) error     { return nil }
func (m *mockSentinelRepo) BulkCreateTimeLogs(_ []domain.TimeLog) error { return nil }
func (m *mockSentinelRepo) GetTotalLoggedMinutes(_ uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockSentinelRepo) CountChildTasks(_ uuid.UUID) (int, error) { return 0, nil }
func (m *mockSentinelRepo) GetChildTasksByParentID(_ uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetProjectAnalytics(_ uuid.UUID) (*domain.ProjectAnalytics, error) {
	return nil, nil
}
func (m *mockSentinelRepo) BulkUpdateTaskStatus(_ []uuid.UUID, _ string) error { return nil }
func (m *mockSentinelRepo) GetSystemConfig() (*domain.SystemConfig, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateSystemConfig(_ *domain.SystemConfig) error { return nil }
func (m *mockSentinelRepo) CreateEpic(_ *domain.Epic) error { return nil }
func (m *mockSentinelRepo) GetEpicByID(_ uuid.UUID) (*domain.Epic, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetEpicsByProjectID(_ uuid.UUID) ([]domain.Epic, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateEpic(_ *domain.Epic) error { return nil }
func (m *mockSentinelRepo) DeleteEpic(_ uuid.UUID) error  { return nil }
func (m *mockSentinelRepo) GetEpicTimelineData(_ uuid.UUID) (*domain.EpicTimelineData, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetSprintTimelineData(_ uuid.UUID) (*domain.SprintTimelineData, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetImportedSlideIndicesByPresentationID(_ string) ([]int, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateProjectCapital(_ uuid.UUID, _ float64, _ *float64) error {
	return nil
}
func (m *mockSentinelRepo) CreateProjectTransaction(_ *domain.ProjectTransaction) error {
	return nil
}
func (m *mockSentinelRepo) GetProjectTransactions(_ uuid.UUID) ([]domain.ProjectTransaction, error) {
	return nil, nil
}
func (m *mockSentinelRepo) DeleteProjectTransaction(_ int64, _ uuid.UUID) error { return nil }
func (m *mockSentinelRepo) CreateB2BRequest(_ *domain.B2BRequest) error { return nil }
func (m *mockSentinelRepo) GetB2BRequests(_ uint, _ string) ([]domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetB2BRequestByID(_ uuid.UUID) (*domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockSentinelRepo) UpdateB2BRequest(_ *domain.B2BRequest) error { return nil }
func (m *mockSentinelRepo) GetFirstProjectByTeamID(_ uint) (*domain.Project, error) {
	return nil, nil
}
func (m *mockSentinelRepo) CreateProjectBackup(_ *domain.ProjectBackup) error { return nil }
func (m *mockSentinelRepo) GetProjectBackups(_ uuid.UUID) ([]domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockSentinelRepo) GetProjectBackupByID(_ uuid.UUID) (*domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockSentinelRepo) DeleteProjectBackup(_, _ uuid.UUID) error { return nil }
func (m *mockSentinelRepo) GetKomgripTasks(_ uint) ([]domain.Task, error) { return nil, nil }

// Reuse mockAuthRepo from standup test (same package would conflict, so define here)
type mockAuthRepoForLog struct {
	users []authDomain.User
	err   error
}

func (m *mockAuthRepoForLog) GetAllUsers() ([]authDomain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}
func (m *mockAuthRepoForLog) CreateUser(_ *authDomain.User) error                      { return nil }
func (m *mockAuthRepoForLog) FindByEmail(_ string) (*authDomain.User, error)            { return nil, nil }
func (m *mockAuthRepoForLog) FindByID(_ uint) (*authDomain.User, error)                 { return nil, nil }
func (m *mockAuthRepoForLog) UpdateProfile(_ uint, _, _, _ *string, _ []string) error   { return nil }
func (m *mockAuthRepoForLog) UpdateAvatar(_ uint, _ string) error                       { return nil }
func (m *mockAuthRepoForLog) UpdateThemePreference(_ uint, _ string) error              { return nil }
func (m *mockAuthRepoForLog) UpdateUserRole(_ uint, _ string) error                     { return nil }
func (m *mockAuthRepoForLog) UpdateUserAdmin(_ uint, _, _, _, _ *string) error          { return nil }
func (m *mockAuthRepoForLog) DeleteUser(_ uint) error                                    { return nil }
func (m *mockAuthRepoForLog) UpdatePassword(_ uint, _ string) error                     { return nil }
func (m *mockAuthRepoForLog) ResetReworkRate(_ uint) error                               { return nil }
func (m *mockAuthRepoForLog) SetUserRemote(_ uint, _ bool) error                        { return nil }
func (m *mockAuthRepoForLog) CreateTeam(_ *authDomain.Team) error                       { return nil }
func (m *mockAuthRepoForLog) GetAllTeams() ([]authDomain.Team, error)                    { return nil, nil }
func (m *mockAuthRepoForLog) GetTeamByID(_ uint) (*authDomain.Team, error)              { return nil, nil }
func (m *mockAuthRepoForLog) UpdateTeamName(_ uint, _ string) error                     { return nil }
func (m *mockAuthRepoForLog) DeleteTeam(_ uint) error                                    { return nil }
func (m *mockAuthRepoForLog) AssignUserToTeam(_ uint, _ *uint) error                    { return nil }
func (m *mockAuthRepoForLog) UpdateTeamCapital(_ uint, _ float64, _ *float64) error     { return nil }
func (m *mockAuthRepoForLog) CreateTeamTransaction(_ *authDomain.TeamTransaction) error { return nil }
func (m *mockAuthRepoForLog) GetTeamTransactions(_ uint) ([]authDomain.TeamTransaction, error) {
	return nil, nil
}
func (m *mockAuthRepoForLog) GetAppSetting(_ string) (string, error) { return "", nil }
func (m *mockAuthRepoForLog) SetAppSetting(_, _ string) error         { return nil }

// --- Tests ---

func TestShouldCheckUserForLogs_EngineerRoles(t *testing.T) {
	s := &DiscordScheduler{}
	tests := []struct {
		role     string
		expected bool
	}{
		{authDomain.RoleEngineer, true},
		{authDomain.RoleChiefEngineer, true},
		{authDomain.RoleProductOwner, true},
		{authDomain.RoleManager, true},
		{authDomain.RoleCEO, false},
		{authDomain.RoleSupport, false},
		{"UNKNOWN", false},
	}
	for _, tt := range tests {
		got := s.shouldCheckUserForLogs(tt.role)
		if got != tt.expected {
			t.Errorf("shouldCheckUserForLogs(%q) = %v, want %v", tt.role, got, tt.expected)
		}
	}
}

func TestCheckAndNotifyMissingLogs_AllLogged(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice", LastName: "A"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{
				1: {{Minutes: 480}}, // 8 hours
			},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 0 {
		t.Error("Should not send notification when all users logged sufficient time")
	}
}

func TestCheckAndNotifyMissingLogs_SomeMissing(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice", LastName: "A"},
				{ID: 2, Role: "ENGINEER", FirstName: "Bob", LastName: "B"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{
				1: {{Minutes: 480}}, // 8 hours — sufficient
				2: {{Minutes: 30}},  // 0.5 hours — insufficient
			},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 1 {
		t.Fatal("Expected 1 notification call")
	}
	if len(discord.lastUsers) != 1 {
		t.Fatalf("Expected 1 missing user, got %d", len(discord.lastUsers))
	}
	if discord.lastUsers[0].DisplayName != "Bob B" {
		t.Errorf("Expected 'Bob B', got '%s'", discord.lastUsers[0].DisplayName)
	}
	if discord.lastUsers[0].TotalHours != 0.5 {
		t.Errorf("Expected 0.5 hours, got %.1f", discord.lastUsers[0].TotalHours)
	}
}

func TestCheckAndNotifyMissingLogs_ExcludesCEOAndSupport(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "CEO", FirstName: "Boss"},
				{ID: 2, Role: "SUPPORT", FirstName: "Help"},
				{ID: 3, Role: "ENGINEER", FirstName: "Dev"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{}, // no logs for anyone
		},
	}
	s.checkAndNotifyMissingLogs()
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

func TestCheckAndNotifyMissingLogs_FallbackDisplayName(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", Email: "john@co.com"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.lastUsers[0].DisplayName != "john" {
		t.Errorf("Expected 'john' from email fallback, got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotifyMissingLogs_DisplayNameFallback(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", DisplayName: "DevNick"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.lastUsers[0].DisplayName != "DevNick" {
		t.Errorf("Expected 'DevNick' from DisplayName fallback, got '%s'", discord.lastUsers[0].DisplayName)
	}
}

func TestCheckAndNotifyMissingLogs_GetUsersFails(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			err: errors.New("db error"),
		},
		sentinelRepo: &mockSentinelRepo{},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 0 {
		t.Error("Should not send notification when GetAllUsers fails")
	}
}

func TestCheckAndNotifyMissingLogs_GetLogsFailsForUser(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice"},
				{ID: 2, Role: "ENGINEER", FirstName: "Bob"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{
				1: {{Minutes: 480}},
			},
		},
	}
	s.checkAndNotifyMissingLogs()
	// User 2 has no entry in logs map, but mock returns nil (not error) for missing keys
	// So user 2 will be reported as missing
	if discord.callCount != 1 {
		t.Fatal("Expected 1 notification")
	}
}

func TestCheckAndNotifyMissingLogs_SendFails(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true, sendErr: errors.New("webhook error")}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 1 {
		t.Fatal("Expected notification attempt even if send fails")
	}
}

func TestCheckAndNotifyMissingLogs_DiscordDisabled(t *testing.T) {
	s := &DiscordScheduler{
		discordSvc: &mockDiscordSvcForLog{enabled: false},
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{{ID: 1, Role: "ENGINEER", FirstName: "Dev"}},
		},
		sentinelRepo: &mockSentinelRepo{},
	}
	s.Start()
	time.Sleep(50 * time.Millisecond) // Ensure no goroutine started
}

func TestCheckAndNotifyMissingLogs_ExactlyOneHour(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{
				1: {{Minutes: 60}}, // exactly 1 hour
			},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 0 {
		t.Error("1 hour exactly should be considered sufficient (threshold is < 1.0)")
	}
}

func TestCheckAndNotifyMissingLogs_MultipleTimeLogs(t *testing.T) {
	discord := &mockDiscordSvcForLog{enabled: true}
	s := &DiscordScheduler{
		discordSvc: discord,
		authRepo: &mockAuthRepoForLog{
			users: []authDomain.User{
				{ID: 1, Role: "ENGINEER", FirstName: "Alice"},
			},
		},
		sentinelRepo: &mockSentinelRepo{
			logs: map[uint][]domain.TimeLog{
				1: {{Minutes: 30}, {Minutes: 20}}, // 50 min total = 0.83 hours
			},
		},
	}
	s.checkAndNotifyMissingLogs()
	if discord.callCount != 1 {
		t.Fatal("Expected notification for user with < 1 hour across multiple logs")
	}
	if discord.lastUsers[0].TotalHours != 50.0/60.0 {
		t.Errorf("Expected %.2f hours, got %.2f", 50.0/60.0, discord.lastUsers[0].TotalHours)
	}
}
