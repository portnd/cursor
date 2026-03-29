package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)

// PreviewPPTX parses an uploaded PPTX file and returns a slide list for triage.
func (u *sentinelUsecase) PreviewPPTX(data []byte) (*domain.PreviewPPTXResult, error) {
	if len(data) == 0 {
		return nil, errors.New("PPTX file is empty")
	}
	pres, err := parsePPTX(data)
	if err != nil {
		return nil, fmt.Errorf("invalid PPTX file: %w", err)
	}
	if len(pres.Slides) == 0 {
		return nil, errors.New("no slides found in the PPTX file")
	}
	slides := make([]domain.PreviewSlideItem, 0, len(pres.Slides))
	for _, s := range pres.Slides {
		t := s.Title
		if t == "" {
			t = fmt.Sprintf("Slide %d", s.Index)
		}
		slides = append(slides, domain.PreviewSlideItem{
			Index:              s.Index,
			Title:              t,
			SuggestedTaskTitle: suggestedTaskTitleFromSlideText(s.Body, s.Index),
			Hidden:             s.Hidden,
		})
	}
	title := strings.TrimSpace(pres.Title)
	if title == "" {
		title = "Imported Presentation"
	}
	return &domain.PreviewPPTXResult{Title: title, Slides: slides}, nil
}

// ImportFromPPTX parses uploaded PPTX bytes and creates one task per selected slide.
func (u *sentinelUsecase) ImportFromPPTX(data []byte, req *domain.ImportPPTXRequest, creatorID uint) (*domain.ImportPPTXResult, error) {
	if len(data) == 0 {
		return nil, errors.New("PPTX file is empty")
	}
	if req == nil {
		return nil, errors.New("request payload is required")
	}
	if len(req.Pages) == 0 {
		return nil, errors.New("at least one page is required to import")
	}

	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}
	var sprintUUID *uuid.UUID
	if req.SprintID != "" {
		parsed, err := uuid.Parse(req.SprintID)
		if err != nil {
			return nil, fmt.Errorf("invalid sprint_id: %w", err)
		}
		sprintUUID = &parsed
	}
	var epicUUID *uuid.UUID
	if req.EpicID != "" {
		parsed, err := uuid.Parse(req.EpicID)
		if err != nil {
			return nil, fmt.Errorf("invalid epic_id: %w", err)
		}
		epicUUID = &parsed
	}
	var parentUUID *uuid.UUID
	if req.ParentID != "" {
		parsed, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentUUID = &parsed
	}

	if parentUUID != nil {
		parent, err := u.repo.GetTaskByID(*parentUUID)
		if err != nil || parent == nil {
			return nil, errors.New("parent task not found")
		}
		if parent.ParentID != nil {
			return nil, &domain.ErrBadRequest{Msg: "cannot attach PPTX import under a nested sub-task"}
		}
		if parent.ProjectID == nil || *parent.ProjectID != projectUUID {
			return nil, &domain.ErrBadRequest{Msg: "parent task must belong to the same project"}
		}
	}

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if !validPriorities[priority] {
		priority = "MEDIUM"
	}
	storyPoints := req.StoryPoints
	if storyPoints < 0 {
		storyPoints = 0
	}

	triagedMap := make(map[int]domain.TriagedSlide)
	for _, p := range req.Pages {
		triagedMap[p.SlideIndex] = p
	}

	slug := "task"
	if proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO}); err == nil && proj != nil {
		slug = slugify(proj.Name)
	}

	pres, err := parsePPTX(data)
	if err != nil {
		return nil, fmt.Errorf("invalid PPTX file: %w", err)
	}

	var allSlides []slideInfo
	for _, s := range pres.Slides {
		t := s.Title
		if t == "" {
			t = fmt.Sprintf("Slide %d", s.Index)
		}
		allSlides = append(allSlides, slideInfo{
			Index:  s.Index,
			Title:  t,
			Body:   s.Body,
			Notes:  s.Notes,
			Images: s.Images,
		})
	}

	allowed := make(map[int]bool)
	for idx := range triagedMap {
		allowed[idx] = true
	}
	var selected []slideInfo
	for _, s := range allSlides {
		if allowed[s.Index] {
			selected = append(selected, s)
		}
	}
	if len(selected) == 0 {
		return nil, errors.New("no matching slides found — check that selected slide numbers exist in the file")
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].Index < selected[j].Index })

	presTitle := strings.TrimSpace(pres.Title)
	if presTitle == "" {
		presTitle = "Imported Presentation"
	}

	maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(slug)

	var created []*domain.Task
	for i, slide := range selected {
		var htmlParts []string
		if slide.Body != "" {
			htmlParts = append(htmlParts, "<p>"+html.EscapeString(slide.Body)+"</p>")
		}
		for _, imgSrc := range slide.Images {
			htmlParts = append(htmlParts, "<img src=\""+html.EscapeString(imgSrc)+"\" class=\"editor-image\" alt=\"Slide\" />")
		}
		if slide.Notes != "" {
			htmlParts = append(htmlParts, "<p><em>Speaker Notes:</em> "+html.EscapeString(slide.Notes)+"</p>")
		}
		description := strings.Join(htmlParts, "\n")
		if description == "" {
			description = "<p></p>"
		}

		resourceURLs := domain.SlideResourceURLs{
			Source:     "pptx_upload",
			SlideIndex: slide.Index,
		}
		resourceURLsJSON, err := json.Marshal(resourceURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource URLs for slide %d: %w", slide.Index, err)
		}

		code := fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i)
		projectIDCopy := projectUUID

		taskTitle := suggestedTaskTitleFromSlideText(slide.Body, slide.Index)
		taskPriority := priority
		taskEstimatedMinutes := 0
		var taskAssigneeID *uint
		if triage, ok := triagedMap[slide.Index]; ok {
			if strings.TrimSpace(triage.Title) != "" {
				taskTitle = strings.TrimSpace(triage.Title)
			}
			p := strings.ToUpper(strings.TrimSpace(triage.Priority))
			if validPriorities[p] {
				taskPriority = p
			}
			if triage.EstimatedMinutes > 0 {
				taskEstimatedMinutes = triage.EstimatedMinutes
			}
			taskAssigneeID = triage.AssigneeID
		}

		task := &domain.Task{
			ID:               uuid.New(),
			Code:             code,
			Title:            taskTitle,
			Description:      description,
			TaskType:         string(domain.TaskTypeTask),
			CreatedBy:        &creatorID,
			Status:           "PENDING",
			Priority:         taskPriority,
			StoryPoints:      storyPoints,
			EstimatedMinutes: taskEstimatedMinutes,
			SprintID:         sprintUUID,
			EpicID:           epicUUID,
			ProjectID:        &projectIDCopy,
			ParentID:         parentUUID,
			ResourceURLs:     datatypes.JSON(resourceURLsJSON),
		}
		if taskAssigneeID != nil {
			task.AssignedTo = taskAssigneeID
		}

		if err := u.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("failed to create task for slide %d: %w", slide.Index, err)
		}
		created = append(created, task)
	}

	return &domain.ImportPPTXResult{
		CreatedCount: len(created),
		PageCount:    len(selected),
		Title:        presTitle,
		Tasks:        created,
	}, nil
}
