package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)

var canvaDesignIDRe = regexp.MustCompile(`(?i)canva\.com/design/([A-Za-z0-9_-]+)`)

const canvaAPIBase = "https://api.canva.com/rest/v1"

func extractCanvaDesignID(rawURL string) (string, error) {
	s := strings.TrimSpace(rawURL)
	if s == "" {
		return "", errors.New("empty Canva URL")
	}
	m := canvaDesignIDRe.FindStringSubmatch(s)
	if len(m) < 2 {
		return "", errors.New("could not find Canva design id in URL (expected ...canva.com/design/DESIGN_ID/...)")
	}
	return m[1], nil
}

type canvaDesignEnvelope struct {
	Design struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		PageCount *int   `json:"page_count"`
	} `json:"design"`
}

func fetchCanvaDesign(ctx context.Context, designID, accessToken string) (title string, pageCount int, err error) {
	token := strings.TrimSpace(accessToken)
	if token == "" {
		return "", 0, errors.New("Canva access token is empty")
	}
	url := fmt.Sprintf("%s/designs/%s", canvaAPIBase, designID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("Canva API request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Canva get design failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var env canvaDesignEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return "", 0, fmt.Errorf("invalid Canva design response: %w", err)
	}
	title = strings.TrimSpace(env.Design.Title)
	if title == "" {
		title = "Canva design"
	}
	if env.Design.PageCount != nil && *env.Design.PageCount > 0 {
		pageCount = *env.Design.PageCount
	} else {
		pageCount = 1
	}
	return title, pageCount, nil
}

type canvaExportCreateBody struct {
	DesignID string `json:"design_id"`
	Format   struct {
		Type string `json:"type"`
	} `json:"format"`
}

type canvaJobEnvelope struct {
	Job struct {
		ID     string   `json:"id"`
		Status string   `json:"status"`
		URLs   []string `json:"urls"`
		Error  *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"job"`
}

func createCanvaPPTXExport(ctx context.Context, designID, accessToken string) (jobID string, err error) {
	var body canvaExportCreateBody
	body.DesignID = designID
	body.Format.Type = "pptx"
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, canvaAPIBase+"/exports", bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Canva export create failed: %w", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Canva export create failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var env canvaJobEnvelope
	if err := json.Unmarshal(b, &env); err != nil {
		return "", fmt.Errorf("invalid Canva export create response: %w", err)
	}
	if env.Job.ID == "" {
		return "", errors.New("Canva export create returned no job id")
	}
	return env.Job.ID, nil
}

func getCanvaExportJob(ctx context.Context, jobID, accessToken string) (*canvaJobEnvelope, error) {
	url := fmt.Sprintf("%s/exports/%s", canvaAPIBase, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Canva export status failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var env canvaJobEnvelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

func exportCanvaAsPPTX(ctx context.Context, designID, accessToken string) ([]byte, error) {
	jobID, err := createCanvaPPTXExport(ctx, designID, accessToken)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		env, err := getCanvaExportJob(ctx, jobID, accessToken)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(strings.TrimSpace(env.Job.Status)) {
		case "success":
			if len(env.Job.URLs) == 0 {
				return nil, errors.New("Canva export succeeded but returned no download URLs")
			}
			return downloadBytes(ctx, env.Job.URLs[0])
		case "failed":
			if env.Job.Error != nil {
				return nil, fmt.Errorf("Canva export failed: %s — %s", env.Job.Error.Code, env.Job.Error.Message)
			}
			return nil, errors.New("Canva export failed")
		case "in_progress":
			time.Sleep(2 * time.Second)
		default:
			time.Sleep(2 * time.Second)
		}
	}
	return nil, errors.New("Canva export timed out after 5 minutes")
}

func downloadBytes(ctx context.Context, downloadURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download export file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slurp, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download export file: HTTP %d — %s", resp.StatusCode, strings.TrimSpace(string(slurp)))
	}
	return io.ReadAll(resp.Body)
}

// PreviewCanva loads design metadata from Canva Connect API and returns one preview row per page.
func (u *sentinelUsecase) PreviewCanva(req *domain.PreviewCanvaRequest, accessToken string) (*domain.PreviewCanvaResult, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("Canva is not configured: set CANVA_ACCESS_TOKEN on the server")
	}
	if req == nil || strings.TrimSpace(req.DesignURL) == "" {
		return nil, errors.New("design_url is required")
	}
	designID, err := extractCanvaDesignID(req.DesignURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	title, pageCount, err := fetchCanvaDesign(ctx, designID, accessToken)
	if err != nil {
		return nil, err
	}
	pages := make([]domain.PreviewSlideItem, 0, pageCount)
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, domain.PreviewSlideItem{
			Index:  i,
			Title:  fmt.Sprintf("Page %d", i),
			Hidden: false,
		})
	}
	return &domain.PreviewCanvaResult{
		DesignTitle: title,
		DesignID:    designID,
		Pages:       pages,
	}, nil
}

// ImportFromCanva exports the design as PPTX via Canva API, parses it, and creates one task per selected page.
func (u *sentinelUsecase) ImportFromCanva(req *domain.ImportCanvaRequest, accessToken string, creatorID uint) (*domain.ImportCanvaResult, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("Canva is not configured: set CANVA_ACCESS_TOKEN on the server")
	}
	if req == nil {
		return nil, errors.New("request is required")
	}
	if len(req.Pages) == 0 {
		return nil, errors.New("at least one page is required to import")
	}
	designID, err := extractCanvaDesignID(req.DesignURL)
	if err != nil {
		return nil, err
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
	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
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
			return nil, &domain.ErrBadRequest{Msg: "cannot attach Canva import under a nested sub-task"}
		}
		if parent.ProjectID == nil || *parent.ProjectID != projectUUID {
			return nil, &domain.ErrBadRequest{Msg: "parent task must belong to the same project"}
		}
	}

	triagedMap := make(map[int]domain.TriagedSlide)
	for _, p := range req.Pages {
		triagedMap[p.SlideIndex] = p
	}

	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if !validPriorities[priority] {
		priority = "MEDIUM"
	}
	storyPoints := req.StoryPoints
	if storyPoints < 0 {
		storyPoints = 0
	}

	slug := "task"
	if proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO}); err == nil && proj != nil {
		slug = slugify(proj.Name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()

	pptxBytes, err := exportCanvaAsPPTX(ctx, designID, accessToken)
	if err != nil {
		return nil, err
	}
	pres, err := parsePPTX(pptxBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse exported PPTX: %w", err)
	}

	var slides []slideInfo
	for _, s := range pres.Slides {
		title := s.Title
		if title == "" {
			title = fmt.Sprintf("Page %d", s.Index)
		}
		slides = append(slides, slideInfo{
			Index:  s.Index,
			Title:  title,
			Body:   s.Body,
			Notes:  s.Notes,
			Images: s.Images,
		})
	}

	allowed := make(map[int]bool)
	for idx := range triagedMap {
		allowed[idx] = true
	}
	filtered := slides[:0]
	for _, s := range slides {
		if allowed[s.Index] {
			filtered = append(filtered, s)
		}
	}
	if len(filtered) == 0 {
		return nil, errors.New("no matching pages in export — check that selected page numbers exist in the design")
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Index < filtered[j].Index })

	designTitle := pres.Title
	if strings.TrimSpace(designTitle) == "" {
		designTitle = "Canva design"
	}

	canvaEditURL := fmt.Sprintf("https://www.canva.com/design/%s/edit", designID)
	maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(slug)

	var createdTasks []*domain.Task
	for i, slide := range filtered {
		var htmlParts []string
		if slide.Body != "" {
			htmlParts = append(htmlParts, "<p>"+html.EscapeString(slide.Body)+"</p>")
		}
		for _, imgSrc := range slide.Images {
			htmlParts = append(htmlParts, "<img src=\""+html.EscapeString(imgSrc)+"\" class=\"editor-image\" alt=\"Canva page\" />")
		}
		if slide.Notes != "" {
			htmlParts = append(htmlParts, "<p><em>Notes:</em> "+html.EscapeString(slide.Notes)+"</p>")
		}
		description := strings.Join(htmlParts, "\n")
		if description == "" {
			description = "<p></p>"
		}

		resourceURLs := domain.SlideResourceURLs{
			ThumbnailURL:   "",
			Images:         nil,
			SlideURL:       canvaEditURL,
			Source:         "canva",
			SlideIndex:     slide.Index,
			PresentationID: designID,
			Comments:       nil,
		}
		resourceURLsJSON, err := json.Marshal(resourceURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource URLs for page %d: %w", slide.Index, err)
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
			return nil, fmt.Errorf("failed to create task for page %d: %w", slide.Index, err)
		}
		createdTasks = append(createdTasks, task)
	}

	return &domain.ImportCanvaResult{
		CreatedCount: len(createdTasks),
		PageCount:    len(filtered),
		DesignTitle:  designTitle,
		Tasks:        createdTasks,
	}, nil
}
