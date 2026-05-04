package usecase

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// mockRepoForCycleDetection implements domain.SentinelRepository for testing wouldCreateCycle.
// Only GetTaskByID has real logic; all other methods are stubs.
type mockRepoForCycleDetection struct {
	tasks map[uuid.UUID]*domain.Task
}

func (m *mockRepoForCycleDetection) GetTaskByID(id uuid.UUID) (*domain.Task, error) {
	return m.tasks[id], nil
}

// --- Stub implementations for all other SentinelRepository methods ---

func (m *mockRepoForCycleDetection) CreateProject(p *domain.Project) error { return nil }
func (m *mockRepoForCycleDetection) GetAllProjects(ctx domain.CallerContext) ([]domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetProjectByID(id uuid.UUID, ctx domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetProjectByCode(code string, ctx domain.CallerContext) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTasksByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTasksByProjectIDForProjectPage(projectID uuid.UUID, limit int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTasksByProjectIDForProjectPageCursor(projectID uuid.UUID, limit int, cursorCreatedAt *time.Time, cursorID *uuid.UUID, offset int) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateProject(p *domain.Project) error { return nil }
func (m *mockRepoForCycleDetection) DeleteProject(id uuid.UUID) error      { return nil }
func (m *mockRepoForCycleDetection) DeleteProjectPlan(projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForCycleDetection) AssignProjectTeam(projectID uuid.UUID, teamID *uint) error {
	return nil
}
func (m *mockRepoForCycleDetection) ReplaceProjectPmAssignments(projectID uuid.UUID, userIDs []uint) error {
	return nil
}
func (m *mockRepoForCycleDetection) CreateTask(task *domain.Task) error { return nil }
func (m *mockRepoForCycleDetection) GetTaskByCode(code string) (*domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) CountTasksForCode(projectID *uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockRepoForCycleDetection) GetMaxTaskCodeSuffix(prefix string) (int, error) {
	return 0, nil
}
func (m *mockRepoForCycleDetection) GetTasksByAssignee(userID uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetActiveSprintTasksByAssignee(userID uint) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetActiveSprintsForUser(userID uint) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetGlobalActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTeamActiveTasks(ctx domain.CallerContext) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetActiveFeatures(teamID uint, projectID *uuid.UUID) ([]domain.FeatureRoadmapItem, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetUnassignedTasks() ([]domain.Task, error)    { return nil, nil }
func (m *mockRepoForCycleDetection) GetAllTasks() ([]domain.Task, error)           { return nil, nil }
func (m *mockRepoForCycleDetection) GetTasksByProjectIDs(projectIDs []uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTasksRequiringApproval() ([]domain.Task, error) { return nil, nil }
func (m *mockRepoForCycleDetection) UpdateTask(task *domain.Task) error                { return nil }
func (m *mockRepoForCycleDetection) DeleteTask(id uuid.UUID) error                    { return nil }
func (m *mockRepoForCycleDetection) CreateTaskActivity(e *domain.TaskActivityEvent) error {
	return nil
}
func (m *mockRepoForCycleDetection) ListTaskActivitiesByTaskID(taskID uuid.UUID) ([]domain.TaskActivityEvent, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetAllTaskDependencies() ([]domain.TaskDependency, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) CreateTaskDependency(dep *domain.TaskDependency) error { return nil }
func (m *mockRepoForCycleDetection) DeleteTaskDependency(id uuid.UUID) error             { return nil }
func (m *mockRepoForCycleDetection) CreateSubmission(sub *domain.Submission) error       { return nil }
func (m *mockRepoForCycleDetection) GetSubmissionByID(id uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateSubmission(sub *domain.Submission) error { return nil }
func (m *mockRepoForCycleDetection) GetLatestSubmission(taskID uuid.UUID) (*domain.Submission, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) CreateAppeal(appeal *domain.Appeal) error { return nil }
func (m *mockRepoForCycleDetection) GetAppealBySubmissionID(subID uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetAppealByID(id uuid.UUID) (*domain.Appeal, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateAppeal(appeal *domain.Appeal) error { return nil }
func (m *mockRepoForCycleDetection) ApproveTask(id uuid.UUID) error           { return nil }
func (m *mockRepoForCycleDetection) RejectTask(taskID uuid.UUID, rejectorID uint, reason string) error {
	return nil
}
func (m *mockRepoForCycleDetection) GetTasksReadyForTest(teamID uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) SetTaskReadyForUAT(taskID uuid.UUID, uatPayload []byte) error {
	return nil
}
func (m *mockRepoForCycleDetection) SetTaskWaitForDeploy(taskID uuid.UUID, uatPayload []byte) error {
	return nil
}
func (m *mockRepoForCycleDetection) AdvanceTaskToReadyForUAT(taskID uuid.UUID) error { return nil }
func (m *mockRepoForCycleDetection) GetTasksReadyForCEOApproval(teamID uint) ([]domain.GlobalActiveTask, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) CreateSprint(sprint *domain.Sprint) error { return nil }
func (m *mockRepoForCycleDetection) GetSprintByID(id uuid.UUID) (*domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetSprintsByProjectID(projectID uuid.UUID) ([]domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetActiveSprintByProjectID(projectID uuid.UUID) (*domain.Sprint, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTasksBySprintID(sprintID uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateSprint(s *domain.Sprint) error { return nil }
func (m *mockRepoForCycleDetection) DeleteSprint(id uuid.UUID) error    { return nil }
func (m *mockRepoForCycleDetection) CreateMilestone(mil *domain.Milestone) error { return nil }
func (m *mockRepoForCycleDetection) GetMilestoneByID(id uuid.UUID) (*domain.Milestone, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetMilestonesByProjectID(projectID uuid.UUID) ([]domain.Milestone, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateMilestone(mil *domain.Milestone) error { return nil }
func (m *mockRepoForCycleDetection) DeleteMilestone(id uuid.UUID) error          { return nil }
func (m *mockRepoForCycleDetection) CreateTaskComment(c *domain.TaskComment) error { return nil }
func (m *mockRepoForCycleDetection) GetCommentsByTaskID(taskID uuid.UUID) ([]domain.TaskComment, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTaskCommentByID(commentID uuid.UUID) (*domain.TaskComment, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateTaskComment(c *domain.TaskComment) error { return nil }
func (m *mockRepoForCycleDetection) DeleteTaskComment(commentID uuid.UUID) error  { return nil }
func (m *mockRepoForCycleDetection) CreateTimeLog(t *domain.TimeLog) error { return nil }
func (m *mockRepoForCycleDetection) GetTimeLogsByTaskID(taskID uuid.UUID) ([]domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetTimeLogByID(logID uuid.UUID) (*domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateTimeLog(t *domain.TimeLog) error { return nil }
func (m *mockRepoForCycleDetection) DeleteTimeLog(logID uuid.UUID) error  { return nil }
func (m *mockRepoForCycleDetection) GetTimeLogsByUserAndDate(userID uint, date time.Time) ([]domain.TimeLog, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) BulkCreateTimeLogs(logs []domain.TimeLog) error { return nil }
func (m *mockRepoForCycleDetection) GetTotalLoggedMinutes(taskID uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockRepoForCycleDetection) CountChildTasks(parentID uuid.UUID) (int, error) { return 0, nil }
func (m *mockRepoForCycleDetection) GetChildTasksByParentID(parentID uuid.UUID) ([]domain.Task, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetProjectAnalytics(projectID uuid.UUID) (*domain.ProjectAnalytics, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) BulkUpdateTaskStatus(taskIDs []uuid.UUID, status string) error {
	return nil
}
func (m *mockRepoForCycleDetection) GetSystemConfig() (*domain.SystemConfig, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateSystemConfig(config *domain.SystemConfig) error { return nil }
func (m *mockRepoForCycleDetection) CreateEpic(epic *domain.Epic) error                  { return nil }
func (m *mockRepoForCycleDetection) GetEpicByID(id uuid.UUID) (*domain.Epic, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetEpicsByProjectID(projectID uuid.UUID) ([]domain.Epic, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateEpic(epic *domain.Epic) error { return nil }
func (m *mockRepoForCycleDetection) DeleteEpic(id uuid.UUID) error      { return nil }
func (m *mockRepoForCycleDetection) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetSprintTimelineData(projectID uuid.UUID) (*domain.SprintTimelineData, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetImportedSlideIndicesByPresentationID(presentationID string) ([]int, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateProjectCapital(projectID uuid.UUID, newBalance float64, bonusPct *float64) error {
	return nil
}
func (m *mockRepoForCycleDetection) CreateProjectTransaction(tx *domain.ProjectTransaction) error {
	return nil
}
func (m *mockRepoForCycleDetection) GetProjectTransactions(projectID uuid.UUID) ([]domain.ProjectTransaction, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) DeleteProjectTransaction(txID int64, projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForCycleDetection) CreateB2BRequest(req *domain.B2BRequest) error { return nil }
func (m *mockRepoForCycleDetection) GetB2BRequests(teamID uint, direction string) ([]domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetB2BRequestByID(id uuid.UUID) (*domain.B2BRequest, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) UpdateB2BRequest(req *domain.B2BRequest) error { return nil }
func (m *mockRepoForCycleDetection) GetFirstProjectByTeamID(teamID uint) (*domain.Project, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) CreateProjectBackup(backup *domain.ProjectBackup) error {
	return nil
}
func (m *mockRepoForCycleDetection) GetProjectBackups(projectID uuid.UUID) ([]domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) GetProjectBackupByID(id uuid.UUID) (*domain.ProjectBackup, error) {
	return nil, nil
}
func (m *mockRepoForCycleDetection) DeleteProjectBackup(id uuid.UUID, projectID uuid.UUID) error {
	return nil
}
func (m *mockRepoForCycleDetection) GetKomgripTasks(userID uint) ([]domain.Task, error) {
	return nil, nil
}

func ptrUUID(id uuid.UUID) *uuid.UUID { return &id }

func TestWouldCreateCycle_SelfReference(t *testing.T) {
	taskID := uuid.New()
	uc := &sentinelUsecase{repo: &mockRepoForCycleDetection{}}

	if !uc.wouldCreateCycle(taskID, taskID) {
		t.Error("expected self-reference to create cycle")
	}
}

func TestWouldCreateCycle_NoCycle(t *testing.T) {
	taskA := uuid.New()
	taskB := uuid.New()

	mock := &mockRepoForCycleDetection{
		tasks: map[uuid.UUID]*domain.Task{
			taskB: {ID: taskB, ParentID: nil}, // B has no parent
		},
	}
	uc := &sentinelUsecase{repo: mock}

	if uc.wouldCreateCycle(taskA, taskB) {
		t.Error("expected no cycle when setting parent to task with no ancestors")
	}
}

func TestWouldCreateCycle_ParentIsDescendant(t *testing.T) {
	taskA := uuid.New()
	taskB := uuid.New()
	taskC := uuid.New()

	mock := &mockRepoForCycleDetection{
		tasks: map[uuid.UUID]*domain.Task{
			taskC: {ID: taskC, ParentID: ptrUUID(taskB)}, // C's parent is B
			taskB: {ID: taskB, ParentID: ptrUUID(taskA)}, // B's parent is A
		},
	}
	uc := &sentinelUsecase{repo: mock}

	// Setting A's parent to C would create cycle: A -> C -> B -> A
	if !uc.wouldCreateCycle(taskA, taskC) {
		t.Error("expected cycle when parent is a descendant")
	}
}

func TestWouldCreateCycle_ExistingCycleInData(t *testing.T) {
	taskA := uuid.New()
	taskB := uuid.New()

	mock := &mockRepoForCycleDetection{
		tasks: map[uuid.UUID]*domain.Task{
			taskA: {ID: taskA, ParentID: ptrUUID(taskB)}, // A's parent is B
			taskB: {ID: taskB, ParentID: ptrUUID(taskA)}, // B's parent is A (cycle!)
		},
	}
	uc := &sentinelUsecase{repo: mock}

	taskC := uuid.New()
	// Setting C's parent to A should detect the existing cycle
	if !uc.wouldCreateCycle(taskC, taskA) {
		t.Error("expected cycle detection when traversing existing cycle")
	}
}

func TestWouldCreateCycle_DeepChain(t *testing.T) {
	taskA := uuid.New()
	taskB := uuid.New()
	taskC := uuid.New()
	taskD := uuid.New()
	taskE := uuid.New()

	mock := &mockRepoForCycleDetection{
		tasks: map[uuid.UUID]*domain.Task{
			taskE: {ID: taskE, ParentID: ptrUUID(taskD)},
			taskD: {ID: taskD, ParentID: ptrUUID(taskC)},
			taskC: {ID: taskC, ParentID: ptrUUID(taskB)},
			taskB: {ID: taskB, ParentID: ptrUUID(taskA)},
			taskA: {ID: taskA, ParentID: nil},
		},
	}
	uc := &sentinelUsecase{repo: mock}

	// Setting A's parent to E would create cycle
	if !uc.wouldCreateCycle(taskA, taskE) {
		t.Error("expected cycle in deep chain")
	}

	// Setting F's parent to E should NOT create cycle (F is new, not in chain)
	taskF := uuid.New()
	if uc.wouldCreateCycle(taskF, taskE) {
		t.Error("expected no cycle when setting parent to leaf of deep chain")
	}
}
