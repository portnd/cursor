package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// projectBackupUsecase implements the backup/restore methods on the main sentinel usecase.
// We attach these methods to sentinelUsecase (defined in sentinel_usecase.go) so the
// existing SentinelUsecase interface is satisfied by one struct.

// CreateProjectBackup snapshots the full project (epics, sprints, milestones, tasks) into JSONB.
func (u *sentinelUsecase) CreateProjectBackup(projectID uuid.UUID, label string, createdBy *uint) (*sentinelDomain.ProjectBackup, error) {
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}

	project, err := u.repo.GetProjectByID(projectID, ctx)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	epics, err := u.repo.GetEpicsByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load epics: %w", err)
	}

	sprints, err := u.repo.GetSprintsByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load sprints: %w", err)
	}

	milestones, err := u.repo.GetMilestonesByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load milestones: %w", err)
	}

	tasks, err := u.repo.GetTasksByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	// Strip non-serialisable GORM associations to keep payload clean.
	cleanProject := *project
	cleanProject.Tasks = nil

	payload := sentinelDomain.ProjectBackupPayload{
		Project:    cleanProject,
		Epics:      epics,
		Sprints:    sprints,
		Milestones: milestones,
		Tasks:      tasks,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise backup: %w", err)
	}

	if label == "" {
		label = fmt.Sprintf("Backup %s", time.Now().Format("2006-01-02 15:04"))
	}

	backup := &sentinelDomain.ProjectBackup{
		ProjectID: projectID,
		Label:     label,
		Payload:   datatypes.JSON(payloadBytes),
		CreatedBy: createdBy,
	}

	if err := u.repo.CreateProjectBackup(backup); err != nil {
		return nil, fmt.Errorf("failed to save backup: %w", err)
	}

	return backup, nil
}

// GetProjectBackups returns all backup records for a project (payload excluded for list view).
func (u *sentinelUsecase) GetProjectBackups(projectID uuid.UUID) ([]sentinelDomain.ProjectBackup, error) {
	return u.repo.GetProjectBackups(projectID)
}

// RestoreProjectBackup replaces the current project plan with the snapshot stored in the backup.
// It deletes all existing epics/sprints/milestones/tasks then re-creates them from the payload.
// The project metadata (name, description, status, color, capital, bonus) is also restored.
func (u *sentinelUsecase) RestoreProjectBackup(backupID uuid.UUID, projectID uuid.UUID) error {
	backup, err := u.repo.GetProjectBackupByID(backupID)
	if err != nil {
		return fmt.Errorf("backup not found: %w", err)
	}
	if backup.ProjectID != projectID {
		return fmt.Errorf("backup does not belong to this project")
	}

	var payload sentinelDomain.ProjectBackupPayload
	if err := json.Unmarshal(backup.Payload, &payload); err != nil {
		return fmt.Errorf("failed to parse backup payload: %w", err)
	}

	// 1. Update project metadata.
	ctx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	project, err := u.repo.GetProjectByID(projectID, ctx)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}
	project.Name = payload.Project.Name
	project.Description = payload.Project.Description
	project.Status = payload.Project.Status
	project.Color = payload.Project.Color
	project.CapitalBalance = payload.Project.CapitalBalance
	project.BonusPercentage = payload.Project.BonusPercentage
	if err := u.repo.UpdateProject(project); err != nil {
		return fmt.Errorf("failed to restore project metadata: %w", err)
	}

	// 2. Wipe existing plan (tasks, sprints, milestones, epics).
	if err := u.repo.DeleteProjectPlan(projectID); err != nil {
		return fmt.Errorf("failed to clear existing plan: %w", err)
	}

	// 3. Re-create epics (preserve original IDs so task.epic_id references remain valid).
	epicIDMap := map[uuid.UUID]uuid.UUID{} // old → new (same in our case, we keep IDs)
	for i := range payload.Epics {
		e := payload.Epics[i]
		e.ProjectID = projectID
		e.Tasks = nil
		_ = epicIDMap
		if err := u.repo.CreateEpic(&e); err != nil {
			return fmt.Errorf("failed to restore epic %s: %w", e.ID, err)
		}
	}

	// 4. Re-create sprints.
	for i := range payload.Sprints {
		s := payload.Sprints[i]
		s.ProjectID = projectID
		s.Tasks = nil
		if err := u.repo.CreateSprint(&s); err != nil {
			return fmt.Errorf("failed to restore sprint %s: %w", s.ID, err)
		}
	}

	// 5. Re-create milestones.
	for i := range payload.Milestones {
		m := payload.Milestones[i]
		m.ProjectID = projectID
		if err := u.repo.CreateMilestone(&m); err != nil {
			return fmt.Errorf("failed to restore milestone %s: %w", m.ID, err)
		}
	}

	// 6. Re-create tasks (strip nested sub-task slices to avoid duplicate inserts).
	for i := range payload.Tasks {
		t := payload.Tasks[i]
		t.ProjectID = &projectID
		t.SubTasks = nil
		t.Submissions = nil
		t.Epic = nil
		t.ParentTask = nil
		if err := u.repo.CreateTask(&t); err != nil {
			return fmt.Errorf("failed to restore task %s: %w", t.ID, err)
		}
	}

	return nil
}

// DeleteProjectBackup removes a single backup record.
func (u *sentinelUsecase) DeleteProjectBackup(backupID uuid.UUID, projectID uuid.UUID) error {
	return u.repo.DeleteProjectBackup(backupID, projectID)
}

// GetProjectBackupPayload returns the full payload of a single backup (for download/export).
func (u *sentinelUsecase) GetProjectBackupPayload(projectID uuid.UUID, backupID uuid.UUID) (*sentinelDomain.ProjectBackupPayload, error) {
	backup, err := u.repo.GetProjectBackupByID(backupID)
	if err != nil {
		return nil, fmt.Errorf("backup not found: %w", err)
	}
	if backup.ProjectID != projectID {
		return nil, fmt.Errorf("backup does not belong to this project")
	}
	var payload sentinelDomain.ProjectBackupPayload
	if err := json.Unmarshal(backup.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse backup payload: %w", err)
	}
	return &payload, nil
}

// ImportProjectFromBackup creates a brand-new project and populates it from a backup payload.
// The new project gets a fresh UUID and code; all child records also get fresh UUIDs to avoid
// primary-key conflicts with any existing data from the source project.
func (u *sentinelUsecase) ImportProjectFromBackup(newName string, payload *sentinelDomain.ProjectBackupPayload, createdBy *uint) (*sentinelDomain.Project, error) {
	// 1. Determine description/status from payload but use the caller-supplied name.
	description := payload.Project.Description
	status := payload.Project.Status
	if status == "" {
		status = "ACTIVE"
	}

	// CreateProject generates a unique code from the name.
	callerCtx := sentinelDomain.CallerContext{Role: sentinelDomain.RoleCEO}
	newProject, err := u.CreateProject(newName, description, status, callerCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// 2. Build UUID remapping tables so IDs don't collide.
	epicMap := map[uuid.UUID]uuid.UUID{}
	sprintMap := map[uuid.UUID]uuid.UUID{}
	milestoneMap := map[uuid.UUID]uuid.UUID{}
	taskMap := map[uuid.UUID]uuid.UUID{}

	for _, e := range payload.Epics {
		epicMap[e.ID] = uuid.New()
	}
	for _, s := range payload.Sprints {
		sprintMap[s.ID] = uuid.New()
	}
	for _, m := range payload.Milestones {
		milestoneMap[m.ID] = uuid.New()
	}
	for _, t := range payload.Tasks {
		taskMap[t.ID] = uuid.New()
	}

	// 3. Re-create epics with new IDs.
	for _, e := range payload.Epics {
		newEpic := e
		newEpic.ID = epicMap[e.ID]
		newEpic.ProjectID = newProject.ID
		newEpic.Tasks = nil
		if err := u.repo.CreateEpic(&newEpic); err != nil {
			return nil, fmt.Errorf("failed to import epic: %w", err)
		}
	}

	// 4. Re-create sprints with new IDs.
	for _, s := range payload.Sprints {
		newSprint := s
		newSprint.ID = sprintMap[s.ID]
		newSprint.ProjectID = newProject.ID
		newSprint.Tasks = nil
		if err := u.repo.CreateSprint(&newSprint); err != nil {
			return nil, fmt.Errorf("failed to import sprint: %w", err)
		}
	}

	// 5. Re-create milestones with new IDs.
	for _, m := range payload.Milestones {
		newMilestone := m
		newMilestone.ID = milestoneMap[m.ID]
		newMilestone.ProjectID = newProject.ID
		if err := u.repo.CreateMilestone(&newMilestone); err != nil {
			return nil, fmt.Errorf("failed to import milestone: %w", err)
		}
	}

	// 6. Re-create tasks — two passes: first create all, then patch parent references.
	// Pass 1: create with remapped IDs (clear parent/epic/sprint/milestone for now).
	for _, t := range payload.Tasks {
		newTask := t
		newTask.ID = taskMap[t.ID]
		newTask.ProjectID = &newProject.ID
		newTask.SubTasks = nil
		newTask.Submissions = nil
		newTask.Epic = nil
		newTask.ParentTask = nil
		// Remap foreign keys.
		if t.EpicID != nil {
			if newID, ok := epicMap[*t.EpicID]; ok {
				newTask.EpicID = &newID
			} else {
				newTask.EpicID = nil
			}
		}
		if t.SprintID != nil {
			if newID, ok := sprintMap[*t.SprintID]; ok {
				newTask.SprintID = &newID
			} else {
				newTask.SprintID = nil
			}
		}
		if t.MilestoneID != nil {
			if newID, ok := milestoneMap[*t.MilestoneID]; ok {
				newTask.MilestoneID = &newID
			} else {
				newTask.MilestoneID = nil
			}
		}
		// Defer parent_id to pass 2 to avoid FK violation.
		newTask.ParentID = nil
		// Generate a unique code for this imported task so it doesn't conflict with
		// any existing task codes. Use the new task UUID (8-char prefix) as suffix.
		newTask.Code = fmt.Sprintf("%s-%s", newProject.Code, newTask.ID.String()[:8])
		if err := u.repo.CreateTask(&newTask); err != nil {
			return nil, fmt.Errorf("failed to import task: %w", err)
		}
	}

	// Pass 2: patch parent_id references.
	for _, t := range payload.Tasks {
		if t.ParentID == nil {
			continue
		}
		newTaskID := taskMap[t.ID]
		newParentID, ok := taskMap[*t.ParentID]
		if !ok {
			continue
		}
		existing, err := u.repo.GetTaskByID(newTaskID)
		if err != nil {
			continue
		}
		existing.ParentID = &newParentID
		_ = u.repo.UpdateTask(existing)
	}

	return newProject, nil
}
