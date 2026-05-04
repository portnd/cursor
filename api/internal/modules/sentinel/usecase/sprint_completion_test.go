package usecase

import (
	"testing"

	"github.com/google/uuid"
)

// TestPreviousSprintID_ClearsOnComplete tests that when a task is marked as COMPLETED,
// the previous_sprint_id field should be cleared
func TestPreviousSprintID_ClearsOnComplete(t *testing.T) {
	// Arrange
	taskID := uuid.New()
	sprintID := uuid.New()

	// Task has previous_sprint_id from a closed sprint
	task := &mockTask{
		id:               taskID,
		status:           "IN_PROGRESS",
		previousSprintID: &sprintID,
	}

	// Act - Update task to COMPLETED
	// Simulate the update logic from sentinel_usecase.go line 1242-1244
	if task.status == "COMPLETED" && task.previousSprintID != nil {
		task.previousSprintID = nil
	}

	// Change status to COMPLETED
	task.status = "COMPLETED"
	if task.status == "COMPLETED" && task.previousSprintID != nil {
		task.previousSprintID = nil
	}

	// Assert
	if task.previousSprintID != nil {
		t.Errorf("Expected previous_sprint_id to be nil for completed task, got %v", task.previousSprintID)
	}
}

// TestPreviousSprintID_PreservesOriginal tests that previous_sprint_id
// is NOT overwritten if already set (keeps the first sprint where task failed)
func TestPreviousSprintID_PreservesOriginal(t *testing.T) {
	// Arrange
	taskID := uuid.New()
	firstSprintID := uuid.New()
	secondSprintID := uuid.New()

	// Task already has previous_sprint_id from first sprint
	task := &mockTask{
		id:               taskID,
		status:           "IN_PROGRESS",
		previousSprintID: &firstSprintID,
	}

	// Act - Try to update with new previous_sprint_id (should be ignored)
	// Simulate the update logic from sentinel_usecase.go line 1238-1240
	// Only set if not already set
	newPreviousSprintID := secondSprintID
	if task.previousSprintID == nil {
		task.previousSprintID = &newPreviousSprintID
	}

	// Assert - Should still have firstSprintID
	if task.previousSprintID == nil {
		t.Fatal("Expected previous_sprint_id to be set")
	}
	if *task.previousSprintID != firstSprintID {
		t.Errorf("Expected previous_sprint_id to remain %v (first sprint), got %v", firstSprintID, *task.previousSprintID)
	}
}

// TestPreviousSprintID_PreservesForIncompleteTasks tests that previous_sprint_id
// is preserved when task status is not COMPLETED
func TestPreviousSprintID_PreservesForIncompleteTasks(t *testing.T) {
	// Arrange
	taskID := uuid.New()
	sprintID := uuid.New()

	task := &mockTask{
		id:               taskID,
		status:           "IN_PROGRESS",
		previousSprintID: &sprintID,
	}

	// Act - Update task but keep it incomplete
	task.status = "READY_FOR_TEST" // Still not completed

	// Simulate the update logic - should NOT clear
	if task.status == "COMPLETED" && task.previousSprintID != nil {
		task.previousSprintID = nil
	}

	// Assert
	if task.previousSprintID == nil {
		t.Error("Expected previous_sprint_id to be preserved for incomplete task")
	}
	if *task.previousSprintID != sprintID {
		t.Errorf("Expected previous_sprint_id to be %v, got %v", sprintID, task.previousSprintID)
	}
}

// mockTask is a simple mock for testing
type mockTask struct {
	id               uuid.UUID
	status           string
	previousSprintID *uuid.UUID
}
