package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/google/uuid"
	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)

func (u *sentinelUsecase) recordTaskActivity(taskID uuid.UUID, action string, actorID *uint, payload map[string]interface{}) {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		raw = []byte("{}")
	}
	ev := &domain.TaskActivityEvent{
		ID:      uuid.New(),
		TaskID:  taskID,
		Action:  action,
		ActorID: actorID,
		Payload: datatypes.JSON(raw),
	}
	if err := u.repo.CreateTaskActivity(ev); err != nil {
		log.Printf("task activity: failed task=%s action=%s: %v", taskID, action, err)
	}
}

func eventIndicatesInProgressTransition(ev domain.TaskActivityEvent) bool {
	var p struct {
		ToStatus string `json:"to_status"`
	}
	_ = json.Unmarshal(ev.Payload, &p)
	to := strings.ToUpper(strings.TrimSpace(p.ToStatus))
	if to != "IN_PROGRESS" {
		return false
	}
	switch ev.Action {
	case domain.TaskActivityStatusChanged, domain.TaskActivityWorkflowReject, domain.TaskActivityRejectedReview:
		return true
	default:
		return false
	}
}

func eventIndicatesCompletion(ev domain.TaskActivityEvent) bool {
	switch ev.Action {
	case domain.TaskActivityApprovedReview, domain.TaskActivityCEOFinalApproved, domain.TaskActivityAppealComplete:
		return true
	case domain.TaskActivityStatusChanged:
		var p struct {
			ToStatus string `json:"to_status"`
		}
		_ = json.Unmarshal(ev.Payload, &p)
		return strings.ToUpper(strings.TrimSpace(p.ToStatus)) == "COMPLETED"
	default:
		return false
	}
}

// GetTaskActivityTimeline returns persisted events plus inferred milestones from the task row when audit began mid-flight.
func (u *sentinelUsecase) GetTaskActivityTimeline(taskID uuid.UUID) ([]domain.TaskActivityItem, error) {
	task, err := u.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	events, err := u.repo.ListTaskActivitiesByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	items := make([]domain.TaskActivityItem, 0, len(events)+4)
	for _, e := range events {
		items = append(items, domain.TaskActivityItem{
			ID:          e.ID.String(),
			Action:      e.Action,
			At:          e.CreatedAt,
			ActorUserID: e.ActorID,
			Payload:     json.RawMessage(e.Payload),
			Inferred:    false,
		})
	}

	hasCreated := false
	hasInProgressHint := false
	hasCompletion := false
	for _, e := range events {
		if e.Action == domain.TaskActivityCreated {
			hasCreated = true
		}
		if eventIndicatesInProgressTransition(e) {
			hasInProgressHint = true
		}
		if eventIndicatesCompletion(e) {
			hasCompletion = true
		}
	}

	if !hasCreated && task.CreatedBy != nil {
		p, _ := json.Marshal(map[string]interface{}{
			"title": task.Title,
			"note":  "Created before activity logging — timestamp from task record",
		})
		items = append(items, domain.TaskActivityItem{
			ID:          "syn:created",
			Action:      domain.TaskActivityCreated,
			At:          task.CreatedAt,
			ActorUserID: task.CreatedBy,
			Payload:     p,
			Inferred:    true,
		})
	}

	if task.StartedAt != nil && !hasInProgressHint {
		p, _ := json.Marshal(map[string]interface{}{
			"to_status": "IN_PROGRESS",
			"note":      "Start time from started_at — actor was not recorded in legacy history",
		})
		items = append(items, domain.TaskActivityItem{
			ID:       "syn:started",
			Action:   domain.TaskActivityStatusChanged,
			At:       *task.StartedAt,
			Payload:  p,
			Inferred: true,
		})
	}

	if task.Status == "COMPLETED" && task.CompletedAt != nil && !hasCompletion {
		p, _ := json.Marshal(map[string]interface{}{
			"to_status": "COMPLETED",
			"note":      "Completion time from completed_at — approval trail may predate activity logging",
		})
		items = append(items, domain.TaskActivityItem{
			ID:       "syn:completed",
			Action:   domain.TaskActivityStatusChanged,
			At:       *task.CompletedAt,
			Payload:  p,
			Inferred: true,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].At.Equal(items[j].At) {
			return items[i].ID < items[j].ID
		}
		return items[i].At.Before(items[j].At)
	})

	u.enrichActivityActors(&items)
	u.enrichActivityAssigneeNames(&items)
	return items, nil
}

func activityActorLabel(uu *authDomain.User) string {
	if uu == nil {
		return ""
	}
	if strings.TrimSpace(uu.DisplayName) != "" {
		return strings.TrimSpace(uu.DisplayName)
	}
	return strings.TrimSpace(uu.Email)
}

func (u *sentinelUsecase) enrichActivityActors(items *[]domain.TaskActivityItem) {
	arr := *items
	seen := map[uint]struct{}{}
	var ids []uint
	for i := range arr {
		if arr[i].ActorUserID == nil {
			continue
		}
		id := *arr[i].ActorUserID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	cache := make(map[uint]*authDomain.User, len(ids))
	for _, id := range ids {
		uu, err := u.authRepo.FindByID(id)
		if err == nil && uu != nil {
			cache[id] = uu
		}
	}
	for i := range arr {
		if arr[i].ActorUserID == nil {
			continue
		}
		uu := cache[*arr[i].ActorUserID]
		if uu == nil {
			continue
		}
		arr[i].ActorEmail = uu.Email
		arr[i].ActorDisplayName = activityActorLabel(uu)
	}
}

func (u *sentinelUsecase) enrichActivityAssigneeNames(items *[]domain.TaskActivityItem) {
	arr := *items
	need := map[uint]struct{}{}
	for i := range arr {
		var payload map[string]interface{}
		if err := json.Unmarshal(arr[i].Payload, &payload); err != nil {
			continue
		}
		for _, key := range []string{"assignee_user_id", "assigned_to_user_id", "previous_assignee_user_id"} {
			if v, ok := payload[key]; ok {
				switch t := v.(type) {
				case float64:
					uid := uint(t)
					if uid != 0 {
						need[uid] = struct{}{}
					}
				}
			}
		}
	}
	cache := make(map[uint]string, len(need))
	for id := range need {
		uu, err := u.authRepo.FindByID(id)
		if err == nil && uu != nil {
			cache[id] = activityActorLabel(uu)
		}
	}
	for i := range arr {
		var payload map[string]interface{}
		if err := json.Unmarshal(arr[i].Payload, &payload); err != nil {
			continue
		}
		changed := false
		for _, pair := range []struct{ idKey, nameKey string }{
			{"assignee_user_id", "assignee_display_name"},
			{"assigned_to_user_id", "assigned_to_display_name"},
			{"previous_assignee_user_id", "previous_assignee_display_name"},
		} {
			if v, ok := payload[pair.idKey]; ok {
				if fv, ok2 := v.(float64); ok2 {
					uid := uint(fv)
					if label, ok3 := cache[uid]; ok3 && label != "" {
						payload[pair.nameKey] = label
						changed = true
					}
				}
			}
		}
		if !changed {
			continue
		}
		raw, err := json.Marshal(payload)
		if err != nil {
			continue
		}
		arr[i].Payload = raw
	}
}
