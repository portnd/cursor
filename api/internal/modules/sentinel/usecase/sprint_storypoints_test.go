package usecase

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// mockRepoForSprintSP implements domain.SentinelRepository for StartSprint/ReopenSprint tests.
type mockRepoForSprintSP struct {
	sprints       map[uuid.UUID]*domain.Sprint
	tasksBySprint map[uuid.UUID][]domain.Task
	activeSprint  *domain.Sprint
}

func (m *mockRepoForSprintSP) GetSprintByID(id uuid.UUID) (*domain.Sprint, error) {
	s, ok := m.sprints[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}
func (m *mockRepoForSprintSP) GetActiveSprintByProjectID(projectID uuid.UUID) (*domain.Sprint, error) {
	return m.activeSprint, nil
}
func (m *mockRepoForSprintSP) GetTasksBySprintID(sprintID uuid.UUID) ([]domain.Task, error) {
	return m.tasksBySprint[sprintID], nil
}
func (m *mockRepoForSprintSP) UpdateSprint(s *domain.Sprint) error { return nil }

// --- Stub implementations for all other SentinelRepository methods ---

func (m *mockRepoForSprintSP) CreateProject(p *domain.Project) error { return nil }
func (m *mockRepoForSprintSP) GetAllProjects(ctx domain.CallerContext) ([]domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetProjectByID(id uuid.UUID, ctx domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetProjectByCode(code string, ctx domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTasksByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTasksByProjectIDForProjectPage(projectID uuid.UUID, limit int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTasksByProjectIDForProjectPageCursor(projectID uuid.UUID, limit int, cursorCreatedAt *time.Time, cursorID *uuid.UUID, offset int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateProject(p *domain.Project) error { return nil }
func (m *mockRepoForSprintSP) DeleteProject(id uuid.UUID) error      { return nil }
func (m *mockRepoForSprintSP) DeleteProjectPlan(projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForSprintSP) AssignProjectTeam(projectID uuid.UUID, teamID *uint) error {
	return nil
}
func (m *mockRepoForSprintSP) ReplaceProjectPmAssignments(projectID uuid.UUID, userIDs []uint) error {
	return nil
}
func (m *mockRepoForSprintSP) CreateTask(task *domain.Task) error { return nil }
func (m *mockRepoForSprintSP) GetTaskByID(id uuid.UUID) (*domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTaskByCode(code string) (*domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) CountTasksForCode(projectID *uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockRepoForSprintSP) GetMaxTaskCodeSuffix(prefix string) (int, error) {
	return 0, nil
}
func (m *mockRepoForSprintSP) GetTasksByAssignee(userID uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetActiveSprintTasksByAssignee(userID uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetActiveSprintsForUser(userID uint) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetGlobalActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTeamActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetActiveFeatures(teamID uint, projectID *uuid.UUID) ([]domain.FeatureRoadmapItem, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetUnassignedTasks() ([]domain.Task, error)    { return nil, nil }
func (m *mockRepoForSprintSP) GetAllTasks() ([]domain.Task, error)           { return nil, nil }
func (m *mockRepoForSprintSP) GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTasksRequiringApproval() ([]domain.Task, error) { return nil, nil }
func (m *mockRepoForSprintSP) UpdateTask(task *domain.Task) error                { return nil }
func (m *mockRepoForSprintSP) DeleteTask(id uuid.UUID) error                    { return nil }
func (m *mockRepoForSprintSP) CreateTaskActivity(e *domain.TaskActivityEvent) error {
	return nil
}
func (m *mockRepoForSprintSP) ListTaskActivitiesByTaskID(taskID uuid.UUID) ([]domain.TaskActivityEvent, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetAllTaskDependencies() ([]domain.TaskDependency, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) CreateTaskDependency(dep *domain.TaskDependency) error { return nil }
func (m *mockRepoForSprintSP) DeleteTaskDependency(id uuid.UUID) error             { return nil }
func (m *mockRepoForSprintSP) CreateSubmission(sub *domain.Submission) error       { return nil }
func (m *mockRepoForSprintSP) GetSubmissionByID(id uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateSubmission(sub *domain.Submission) error { return nil }
func (m *mockRepoForSprintSP) GetLatestSubmission(taskID uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) CreateAppeal(appeal *domain.Appeal) error { return nil }
func (m *mockRepoForSprintSP) GetAppealBySubmissionID(subID uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetAppealByID(id uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateAppeal(appeal *domain.Appeal) error { return nil }
func (m *mockRepoForSprintSP) ApproveTask(id uuid.UUID) error           { return nil }
func (m *mockRepoForSprintSP) RejectTask(taskID uuid.UUID, rejectorID uint, reason string) error {
	return nil
}
func (m *mockRepoForSprintSP) GetTasksReadyForTest(teamID uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) SetTaskReadyForUAT(taskID uuid.UUID, uatPayload []byte) error {
	return nil
}
func (m *mockRepoForSprintSP) SetTaskWaitForDeploy(taskID uuid.UUID, uatPayload []byte) error {
	return nil
}
func (m *mockRepoForSprintSP) AdvanceTaskToReadyForUAT(taskID uuid.UUID) error { return nil }
func (m *mockRepoForSprintSP) GetTasksReadyForCEOApproval(teamID uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) CreateSprint(sprint *domain.Sprint) error { return nil }
func (m *mockRepoForSprintSP) GetSprintsByProjectID(projectID uuid.UUID) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) DeleteSprint(id uuid.UUID) error { return nil }
func (m *mockRepoForSprintSP) CreateMilestone(mil *domain.Milestone) error { return nil }
func (m *mockRepoForSprintSP) GetMilestoneByID(id uuid.UUID) (*domain.Milestone, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetMilestonesByProjectID(projectID uuid.UUID) ([]domain.Milestone, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateMilestone(mil *domain.Milestone) error { return nil }
func (m *mockRepoForSprintSP) DeleteMilestone(id uuid.UUID) error          { return nil }
func (m *mockRepoForSprintSP) CreateTaskComment(c *domain.TaskComment) error { return nil }
func (m *mockRepoForSprintSP) GetCommentsByTaskID(taskID uuid.UUID) ([]domain.TaskComment, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTaskCommentByID(commentID uuid.UUID) (*domain.TaskComment, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateTaskComment(c *domain.TaskComment) error { return nil }
func (m *mockRepoForSprintSP) DeleteTaskComment(commentID uuid.UUID) error  { return nil }
func (m *mockRepoForSprintSP) CreateTimeLog(t *domain.TimeLog) error { return nil }
func (m *mockRepoForSprintSP) GetTimeLogsByTaskID(taskID uuid.UUID) ([]domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetTimeLogByID(logID uuid.UUID) (*domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateTimeLog(t *domain.TimeLog) error { return nil }
func (m *mockRepoForSprintSP) DeleteTimeLog(logID uuid.UUID) error  { return nil }
func (m *mockRepoForSprintSP) GetTimeLogsByUserAndDate(userID uint, date time.Time) ([]domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) BulkCreateTimeLogs(logs []domain.TimeLog) error { return nil }
func (m *mockRepoForSprintSP) GetTotalLoggedMinutes(taskID uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockRepoForSprintSP) CountChildTasks(parentID uuid.UUID) (int, error) { return 0, nil }
func (m *mockRepoForSprintSP) GetChildTasksByParentID(parentID uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetProjectAnalytics(projectID uuid.UUID) (*domain.ProjectAnalytics, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error {
	return nil
}
func (m *mockRepoForSprintSP) GetSystemConfig() (*domain.SystemConfig, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateSystemConfig(config *domain.SystemConfig) error { return nil }
func (m *mockRepoForSprintSP) CreateEpic(epic *domain.Epic) error                  { return nil }
func (m *mockRepoForSprintSP) GetEpicByID(id uuid.UUID) (*domain.Epic, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetEpicsByProjectID(projectID uuid.UUID) ([]domain.Epic, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateEpic(epic *domain.Epic) error { return nil }
func (m *mockRepoForSprintSP) DeleteEpic(id uuid.UUID) error      { return nil }
func (m *mockRepoForSprintSP) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetSprintTimelineData(projectID uuid.UUID) (*domain.SprintTimelineData, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetImportedSlideIndicesByPresentationID(presentationID string) ([]int, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateProjectCapital(projectID uuid.UUID, newBalance float64, bonusPct *float64) error {
	return nil
}
func (m *mockRepoForSprintSP) CreateProjectTransaction(tx *domain.ProjectTransaction) error {
	return nil
}
func (m *mockRepoForSprintSP) GetProjectTransactions(projectID uuid.UUID) ([]domain.ProjectTransaction, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) DeleteProjectTransaction(txID int64, projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForSprintSP) CreateB2BRequest(req *domain.B2BRequest) error { return nil }
func (m *mockRepoForSprintSP) GetB2BRequests(teamID uint, direction string) ([]domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetB2BRequestByID(id uuid.UUID) (*domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) UpdateB2BRequest(req *domain.B2BRequest) error { return nil }
func (m *mockRepoForSprintSP) GetFirstProjectByTeamID(teamID uint) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) CreateProjectBackup(backup *domain.ProjectBackup) error {
	return nil
}
func (m *mockRepoForSprintSP) GetProjectBackups(projectID uuid.UUID) ([]domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) GetProjectBackupByID(id uuid.UUID) (*domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockRepoForSprintSP) DeleteProjectBackup(id uuid.UUID, projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForSprintSP) GetKomgripTasks(userID uint) ([]domain.Task, error) {
	return nil, nil
}

// --- StartSprint tests ---

func TestStartSprint_BlockedWhenTasksMissingStoryPoints(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()
	task2ID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {
				{ID: uuid.New(), Code: "A001", Title: "Task with SP", StoryPoints: 3},
				{ID: task2ID, Code: "A002", Title: "Task without SP", StoryPoints: 0},
			},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.StartSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when tasks missing story points, got nil")
	}
	if !domain.IsSprintStoryPointsMissing(err) {
		t.Fatalf("expected ErrSprintStoryPointsMissing, got %T: %v", err, err)
	}

	var spErr *domain.ErrSprintStoryPointsMissing
	if !errors.As(err, &spErr) {
		t.Fatal("errors.As failed to extract ErrSprintStoryPointsMissing")
	}
	if len(spErr.MissingTasks) != 1 {
		t.Fatalf("expected 1 missing task, got %d", len(spErr.MissingTasks))
	}
	if spErr.MissingTasks[0].ID != task2ID {
		t.Errorf("expected missing task ID %v, got %v", task2ID, spErr.MissingTasks[0].ID)
	}
	if spErr.MissingTasks[0].Code != "A002" {
		t.Errorf("expected missing task code A002, got %s", spErr.MissingTasks[0].Code)
	}
	if spErr.MissingTasks[0].Title != "Task without SP" {
		t.Errorf("expected missing task title 'Task without SP', got %s", spErr.MissingTasks[0].Title)
	}
}

func TestStartSprint_AllowedWhenAllTasksHaveStoryPoints(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {
				{ID: uuid.New(), Code: "A001", Title: "Task 1", StoryPoints: 2},
				{ID: uuid.New(), Code: "A002", Title: "Task 2", StoryPoints: 5},
			},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	sprint, err := uc.StartSprint(sprintID)

	if err != nil {
		t.Fatalf("expected no error when all tasks have SP, got %v", err)
	}
	if sprint.Status != "ACTIVE" {
		t.Errorf("expected sprint status ACTIVE, got %s", sprint.Status)
	}
}

func TestStartSprint_AllowedWhenSprintHasNoTasks(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{},
	}

	uc := &sentinelUsecase{repo: mock}
	sprint, err := uc.StartSprint(sprintID)

	if err != nil {
		t.Fatalf("expected no error for empty sprint, got %v", err)
	}
	if sprint.Status != "ACTIVE" {
		t.Errorf("expected sprint status ACTIVE, got %s", sprint.Status)
	}
}

func TestStartSprint_AllTasksMissingStoryPoints(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {
				{ID: uuid.New(), Code: "A001", Title: "Task 1", StoryPoints: 0},
				{ID: uuid.New(), Code: "A002", Title: "Task 2", StoryPoints: 0},
				{ID: uuid.New(), Code: "A003", Title: "Task 3", StoryPoints: 0},
			},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.StartSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when all tasks missing story points")
	}

	var spErr *domain.ErrSprintStoryPointsMissing
	errors.As(err, &spErr)
	if len(spErr.MissingTasks) != 3 {
		t.Fatalf("expected 3 missing tasks, got %d", len(spErr.MissingTasks))
	}
}

func TestStartSprint_BlockedWhenAlreadyActiveSprint(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()
	activeSprintID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 2"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {{ID: uuid.New(), Code: "A001", Title: "Task 1", StoryPoints: 3}},
		},
		activeSprint: &domain.Sprint{ID: activeSprintID, ProjectID: projectID, Status: "ACTIVE", Name: "Sprint 1"},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.StartSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when active sprint already exists")
	}
	if domain.IsSprintStoryPointsMissing(err) {
		t.Fatal("expected conflict error, not story points error")
	}
}

func TestStartSprint_SetsStartDateWhenNil(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {{ID: uuid.New(), Code: "A001", Title: "Task 1", StoryPoints: 2}},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	sprint, err := uc.StartSprint(sprintID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sprint.StartDate == nil {
		t.Error("expected StartDate to be set")
	}
}

// --- ReopenSprint tests ---

func TestReopenSprint_BlockedWhenTasksMissingStoryPoints(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()
	taskID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "COMPLETED", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {
				{ID: taskID, Code: "B001", Title: "Reopened task no SP", StoryPoints: 0},
			},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.ReopenSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when reopening sprint with tasks missing SP")
	}
	if !domain.IsSprintStoryPointsMissing(err) {
		t.Fatalf("expected ErrSprintStoryPointsMissing, got %T: %v", err, err)
	}

	var spErr *domain.ErrSprintStoryPointsMissing
	errors.As(err, &spErr)
	if len(spErr.MissingTasks) != 1 {
		t.Fatalf("expected 1 missing task, got %d", len(spErr.MissingTasks))
	}
	if spErr.MissingTasks[0].Code != "B001" {
		t.Errorf("expected missing task code B001, got %s", spErr.MissingTasks[0].Code)
	}
}

func TestReopenSprint_AllowedWhenAllTasksHaveStoryPoints(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "COMPLETED", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {
				{ID: uuid.New(), Code: "B001", Title: "Task 1", StoryPoints: 5},
			},
		},
	}

	uc := &sentinelUsecase{repo: mock}
	sprint, err := uc.ReopenSprint(sprintID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sprint.Status != "ACTIVE" {
		t.Errorf("expected sprint status ACTIVE, got %s", sprint.Status)
	}
}

func TestReopenSprint_BlockedWhenNotCompleted(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "PLANNING", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.ReopenSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when reopening non-COMPLETED sprint")
	}
	if domain.IsSprintStoryPointsMissing(err) {
		t.Fatal("expected status error, not story points error")
	}
}

func TestReopenSprint_BlockedWhenActiveSprintExists(t *testing.T) {
	sprintID := uuid.New()
	projectID := uuid.New()
	activeSprintID := uuid.New()

	mock := &mockRepoForSprintSP{
		sprints: map[uuid.UUID]*domain.Sprint{
			sprintID: {ID: sprintID, ProjectID: projectID, Status: "COMPLETED", Name: "Sprint 1"},
		},
		tasksBySprint: map[uuid.UUID][]domain.Task{
			sprintID: {{ID: uuid.New(), Code: "B001", Title: "Task 1", StoryPoints: 3}},
		},
		activeSprint: &domain.Sprint{ID: activeSprintID, ProjectID: projectID, Status: "ACTIVE", Name: "Sprint 2"},
	}

	uc := &sentinelUsecase{repo: mock}
	_, err := uc.ReopenSprint(sprintID)

	if err == nil {
		t.Fatal("expected error when active sprint already exists")
	}
	if domain.IsSprintStoryPointsMissing(err) {
		t.Fatal("expected conflict error, not story points error")
	}
}

// --- ErrSprintStoryPointsMissing error type tests ---

func TestErrSprintStoryPointsMissing_Error(t *testing.T) {
	err := &domain.ErrSprintStoryPointsMissing{
		MissingTasks: []domain.StoryPointMissingTask{
			{ID: uuid.New(), Code: "A001", Title: "First task"},
			{ID: uuid.New(), Code: "A002", Title: "Second task"},
		},
	}

	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
	if !strings.Contains(msg, "A001") {
		t.Errorf("expected error message to contain A001, got: %s", msg)
	}
	if !strings.Contains(msg, "A002") {
		t.Errorf("expected error message to contain A002, got: %s", msg)
	}
}

func TestIsSprintStoryPointsMissing_Nil(t *testing.T) {
	if domain.IsSprintStoryPointsMissing(nil) {
		t.Error("expected false for nil error")
	}
}

func TestIsSprintStoryPointsMissing_OtherError(t *testing.T) {
	if domain.IsSprintStoryPointsMissing(errors.New("some other error")) {
		t.Error("expected false for non-SP error")
	}
}

func TestIsSprintStoryPointsMissing_WrappedError(t *testing.T) {
	inner := &domain.ErrSprintStoryPointsMissing{
		MissingTasks: []domain.StoryPointMissingTask{
			{ID: uuid.New(), Code: "X001", Title: "Test"},
		},
	}
	wrapped := fmt.Errorf("wrapped: %w", inner)

	if !domain.IsSprintStoryPointsMissing(wrapped) {
		t.Error("expected true for wrapped ErrSprintStoryPointsMissing")
	}
}
