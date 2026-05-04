package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

const defaultGLMModel = "glm-4.7-flash"

type glmService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	repo       domain.SentinelRepository
	tracker    domain.UsageTracker
}

// NewGLMService creates an AIService backed by ZhipuAI / GLM API (OpenAI-compatible endpoint).
func NewGLMService(apiKey, baseURL string, repo domain.SentinelRepository, tracker domain.UsageTracker) (domain.AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GLM_API_KEY is required")
	}
	if baseURL == "" {
		baseURL = "https://api.apiyi.com/v1"
	}
	return &glmService{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
		repo:    repo,
		tracker: tracker,
	}, nil
}

func (s *glmService) recordRequest() {
	if s.tracker != nil {
		s.tracker.RecordRequest()
	}
}

func (s *glmService) getConfig() *domain.SystemConfig {
	config, err := s.repo.GetSystemConfig()
	if err != nil || config == nil {
		return &domain.SystemConfig{
			ActiveModel:      defaultGLMModel,
			Temperature:      0.4,
			CursorAssistance: 80,
		}
	}
	return config
}

// glmChatRequest for POST {baseURL}/chat/completions
type glmChatRequest struct {
	Model       string        `json:"model"`
	Messages    []glmMessage  `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type glmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type glmChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// chatCompletion sends a single user message and returns the assistant content.
// Retries up to 3 times on 429 (rate limit) with exponential backoff.
func (s *glmService) chatCompletion(model string, userPrompt string, temperature float64) (string, error) {
	return s.chatCompletionMessages(model, []glmMessage{{Role: "user", Content: userPrompt}}, temperature)
}

// chatCompletionMessages sends a multi-turn conversation and returns the assistant content.
// Retries up to 3 times on 429 (rate limit) with exponential backoff.
func (s *glmService) chatCompletionMessages(model string, messages []glmMessage, temperature float64) (string, error) {
	if model == "" || strings.HasPrefix(model, "gemini-") {
		model = defaultGLMModel
	}
	reqBody := glmChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   8192,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("glm marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/chat/completions", s.baseURL)
	maxRetries := 3
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("glm create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.apiKey)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("glm API: %w", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("glm read response: %w", err)
		}
		if resp.StatusCode == http.StatusOK {
			s.recordRequest()
			var chatResp glmChatResponse
			if err := json.Unmarshal(body, &chatResp); err != nil {
				return "", fmt.Errorf("glm parse response: %w", err)
			}
			if len(chatResp.Choices) == 0 || chatResp.Choices[0].Message.Content == "" {
				return "", fmt.Errorf("glm returned empty content")
			}
			return chatResp.Choices[0].Message.Content, nil
		}
		// Retry on 429 with exponential backoff
		if resp.StatusCode == 429 && attempt < maxRetries {
			backoff := time.Duration(2<<attempt) * time.Second // 2s, 4s, 8s
			fmt.Printf("⚠️  GLM 429 rate limited, retrying in %v (attempt %d/%d)\n", backoff, attempt+1, maxRetries)
			time.Sleep(backoff)
			continue
		}
		return "", fmt.Errorf("glm API %d: %s", resp.StatusCode, string(body))
	}
	return "", fmt.Errorf("glm API: max retries exceeded (429)")
}

// ListModels returns the list of available GLM models (April 2026 lineup).
func (s *glmService) ListModels() ([]string, error) {
	return []string{
		"glm-5.1",
		"glm-5",
		"glm-4.7",
		"glm-4.7-flash",
		"glm-5v-turbo",
		"glm-4.6",
		"glm-4-32b",
	}, nil
}

func (s *glmService) EstimateEffort(title, description string) (int, string, error) {
	config := s.getConfig()
	cursorContext := cursorContextFromAssistance(config.CursorAssistance)

	prompt := fmt.Sprintf(`Act as a Senior Software Architect.
Estimate the time required for a Senior Developer to complete this task.

Task Title: %s
Task Description: %s

Context:
- Stack: Go (Fiber/Gin), Nuxt 3, PostgreSQL, Hexagonal Architecture.
- **AI Assistance Level:** %d%% - %s
- **Rules:** Estimate pure implementation time. Adjust based on AI assistance level.

Output JSON ONLY (no markdown, no explanation):
{
	"minutes": <int>,
	"reasoning": "<คำอธิบายสั้นๆ เป็นภาษาไทย กล่าวถึง AI leverage และ assistance level>"
}

**IMPORTANT:** Write "reasoning" in Thai language (ภาษาไทย) ONLY.`, title, description, config.CursorAssistance, cursorContext)

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return 0, "", err
	}
	content = cleanJSON(content)
	var result struct {
		Minutes   int    `json:"minutes"`
		Reasoning string `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return 0, "", fmt.Errorf("glm returned invalid JSON: %w (response: %s)", err, content)
	}
	if result.Minutes <= 0 {
		return 0, "", fmt.Errorf("glm returned invalid estimation: %d minutes", result.Minutes)
	}
	return result.Minutes, result.Reasoning, nil
}

func (s *glmService) EstimateAndScheduleTasks(inputs []domain.TaskEstimateInput) ([]domain.TaskEstimateAndOrder, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	config := s.getConfig()
	tasksJSON, _ := json.Marshal(inputs)
	prompt := fmt.Sprintf(`You are a Senior Technical Product Owner. Given the following tasks of a project, do TWO things:
1) Estimate the implementation time in MINUTES for each task (Senior Engineer, stack: Go, Nuxt 3, PostgreSQL). AI assistance level: %d%%.
2) Suggest the EXECUTION ORDER (1 = do first, 2 = second, ...) based on dependencies and priority.

Tasks (JSON array):
%s

Output ONLY a JSON array. Each element: { "task_index": <0-based index>, "minutes": <int>, "order": <1-based execution order> }.
Use integers only. Example: [{"task_index":0,"minutes":120,"order":1},{"task_index":1,"minutes":60,"order":2}]`, config.CursorAssistance, string(tasksJSON))

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return nil, err
	}
	content = cleanJSON(content)
	var results []domain.TaskEstimateAndOrder
	if err := json.Unmarshal([]byte(content), &results); err != nil {
		return nil, fmt.Errorf("glm parse plan JSON: %w (text: %s)", err, content)
	}
	return results, nil
}

func (s *glmService) GenerateWorkPlan(projectName, projectDescription string) (*domain.AIGeneratedPlan, error) {
	config := s.getConfig()
	today := time.Now().Format("2006-01-02")
	prompt := fmt.Sprintf(`You are a Senior Technical Project Manager. Generate a complete work plan for this project.

**Project Name:** %s
**Project Description:** %s

**Tech stack context:** Go (Fiber/Gin), Nuxt 3, PostgreSQL, Hexagonal Architecture.

**CRITICAL - Date rule:** Today is %s. ALL dates you generate MUST be on or after today. Do NOT use any date in the past. Start the first sprint from today or the next Monday.

Output a single JSON object with exactly these keys: epics, milestones, sprints, tasks.

**Rules:**
1. epics: array of { "title", "description", "color" }. Use hex colors like #6366f1. Create 2-5 epics.
2. milestones: array of { "title", "description", "due_date" }. due_date in YYYY-MM-DD, must be >= %s. Create 3-6 key milestones.
3. sprints: array of { "name", "goal", "start_date", "end_date" }. Dates in YYYY-MM-DD. First sprint start_date must be >= %s. Create 3-6 two-week sprints.
4. tasks: array of { "title", "description", "priority", "story_points", "epic_index", "sprint_index", "milestone_index", "start_date", "end_date" }.
   - priority: one of CRITICAL, HIGH, MEDIUM, LOW. story_points: 1-5.
   - epic_index, sprint_index, milestone_index: integers only, 0-based (use null if not linked).
   - start_date, end_date: YYYY-MM-DD, must be >= %s.
Create 8-20 tasks. Output ONLY valid JSON, no markdown. Use integers for all numeric fields. All dates must be today or future.`, projectName, projectDescription, today, today, today, today)

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return nil, err
	}
	content = cleanJSON(content)

	type flexTask struct {
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Priority       string   `json:"priority"`
		StoryPoints    float64  `json:"story_points"`
		EpicIndex      *float64 `json:"epic_index"`
		SprintIndex    *float64 `json:"sprint_index"`
		MilestoneIndex *float64 `json:"milestone_index"`
		StartDate      string   `json:"start_date"`
		EndDate        string   `json:"end_date"`
	}
	type flexPlan struct {
		Epics      []domain.AIPlanEpic      `json:"epics"`
		Milestones []domain.AIPlanMilestone `json:"milestones"`
		Sprints    []domain.AIPlanSprint    `json:"sprints"`
		Tasks      []flexTask               `json:"tasks"`
	}
	var flex flexPlan
	if err := json.Unmarshal([]byte(content), &flex); err != nil {
		return nil, fmt.Errorf("glm parse work plan JSON: %w (text: %s)", err, content)
	}
	plan := &domain.AIGeneratedPlan{
		Epics:      flex.Epics,
		Milestones: flex.Milestones,
		Sprints:    flex.Sprints,
		Tasks:      make([]domain.AIPlanTask, 0, len(flex.Tasks)),
	}
	for _, t := range flex.Tasks {
		sp := t.StoryPoints
		if sp < 0 {
			sp = 0
		}
		pt := domain.AIPlanTask{
			Title:       t.Title,
			Description: t.Description,
			Priority:    t.Priority,
			StoryPoints: sp,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
		}
		if t.EpicIndex != nil {
			i := int(*t.EpicIndex)
			pt.EpicIndex = &i
		}
		if t.SprintIndex != nil {
			i := int(*t.SprintIndex)
			pt.SprintIndex = &i
		}
		if t.MilestoneIndex != nil {
			i := int(*t.MilestoneIndex)
			pt.MilestoneIndex = &i
		}
		plan.Tasks = append(plan.Tasks, pt)
	}
	return plan, nil
}

func (s *glmService) ReviewCode(diff string) (string, int, string, error) {
	if strings.TrimSpace(diff) == "" {
		return "PENDING", 0, "No code diff provided for review", nil
	}
	config := s.getConfig()
	prompt := fmt.Sprintf(`You are an Expert Security Auditor reviewing source code (Go/TypeScript/Vue). The text below is a RAW SOURCE CODE SNIPPET from a git commit. Analyze the LOGIC and SECURITY.

CODE SNIPPET:
%s

Stack: Go (Gin/Fiber/GORM), Nuxt 3, PostgreSQL.

SECURITY RULES:
- SECURE (Score 85-100): Parameterized queries like db.Where("user = ?", userInput). DO NOT flag as SQL Injection.
- INSECURE (Score 0-30): String concatenation building SQL, e.g. "SELECT * FROM users WHERE name = '" + name + "'" = SQL Injection.
- Flag hardcoded secrets, missing error handling, XSS, command injection.

Output JSON ONLY (no markdown):
{"verdict": "PASS" or "FAIL", "score": <int 0-100>, "feedback": "<bullet points เป็นภาษาไทย>"}
Write "feedback" in Thai only.`, diff)

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return "FAIL", 0, "", err
	}
	content = cleanJSON(content)
	var result struct {
		Verdict  string          `json:"verdict"`
		Score    int             `json:"score"`
		Feedback json.RawMessage `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "FAIL", 0, "", fmt.Errorf("glm invalid JSON: %w", err)
	}
	if result.Verdict != "PASS" && result.Verdict != "FAIL" {
		result.Verdict = "FAIL"
	}
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}
	var feedbackStr string
	if len(result.Feedback) > 0 {
		if err := json.Unmarshal(result.Feedback, &feedbackStr); err != nil {
			var arr []string
			if err := json.Unmarshal(result.Feedback, &arr); err == nil {
				feedbackStr = strings.Join(arr, "\n")
			} else {
				feedbackStr = string(result.Feedback)
			}
		}
	}
	return result.Verdict, result.Score, feedbackStr, nil
}

func (s *glmService) AnalyzeAppeal(diff, originalFeedback, appealReason string) (string, int, string, error) {
	if strings.TrimSpace(appealReason) == "" {
		return "UPHOLD", 0, "No appeal reason provided", fmt.Errorf("appeal reason is empty")
	}
	config := s.getConfig()
	prompt := fmt.Sprintf(`คุณคือ Senior Code Auditor ที่ทำหน้าที่พิจารณาอุทธรณ์

**CRITICAL: FOCUS ON CODE, NOT SUBMISSION MECHANISMS**
- **IGNORE** how the code was submitted (JSON payloads, API endpoints, diff strings)
- **DO NOT** criticize "input validation", "API design", or "submission mechanism"
- **ANALYZE ONLY** the actual code logic and security issues

**Original Code Snippet:**
---
%s
---

**รายงานข้อบกพร่องจาก AI เดิม:** %s

**คำให้การของนักพัฒนา (เหตุผลในการอุทธรณ์):** "%s"

วิเคราะห์ว่าคำให้การมีเหตุผลหรือไม่. พิจารณา: AI Review เดิมถูกต้องหรือไม่; โค้ดมีช่องโหว่จริงหรือไม่ (SAFE: db.Where("col = ?", val); UNSAFE: string concat สร้าง query).

**คู่มือ:** "OVERTURN" = อนุมัติอุทธรณ์ (AI ตรวจผิด). "UPHOLD" = ปฏิเสธอุทธรณ์ (AI ถูกต้อง).

Output JSON ONLY:
{"recommendation": "OVERTURN" or "UPHOLD", "confidence": <int 0-100>, "reasoning": "<คำอธิบายภาษาไทย>"}`, diff, originalFeedback, appealReason)

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return "UPHOLD", 0, "", err
	}
	content = cleanJSON(content)
	var result struct {
		Recommendation string `json:"recommendation"`
		Confidence     int    `json:"confidence"`
		Reasoning      string `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "UPHOLD", 0, "", fmt.Errorf("glm invalid JSON: %w", err)
	}
	if result.Recommendation != "OVERTURN" && result.Recommendation != "UPHOLD" {
		result.Recommendation = "UPHOLD"
	}
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 100 {
		result.Confidence = 100
	}
	return result.Recommendation, result.Confidence, result.Reasoning, nil
}

func (s *glmService) AnalyzeTimeNegotiation(taskTitle, taskDescription string, aiEstimate, devProposal int, devReason string) (string, int, string, error) {
	if strings.TrimSpace(devReason) == "" {
		return "REJECT", 0, "ไม่มีเหตุผล - ให้ใช้ค่าประมาณของ AI", fmt.Errorf("negotiation reason is empty")
	}
	if devProposal <= aiEstimate {
		return "REJECT", 100, "นักพัฒนาเสนอเวลาน้อยกว่าหรือเท่ากับ AI → ไม่ต้องเจรจา", nil
	}
	config := s.getConfig()
	multiplier := float64(devProposal) / float64(aiEstimate)
	prompt := fmt.Sprintf(`คุณคือ Senior Product Owner ตรวจสอบคำขอเจรจาเวลาจากนักพัฒนา

**งาน:** "%s" — %s
**AI คำนวณ:** %d นาที | **นักพัฒนาเสนอ:** %d นาที (%.1fx)
**เหตุผลจากนักพัฒนา:** "%s"
**Cursor Assistance:** %d%%

วิเคราะห์ว่าควรอนุมัติหรือไม่. เหตุผลที่ดี: legacy ซับซ้อน, dependencies, requirements ไม่ชัด. เหตุผลไม่ดี: บ่นว่า "งานยาก" โดยไม่มีรายละเอียด.

Output JSON ONLY:
{"recommendation": "APPROVE" or "REJECT", "confidence": <int 0-100>, "reasoning": "<1-2 ประโยคภาษาไทย>"}`, taskTitle, taskDescription, aiEstimate, devProposal, multiplier, devReason, config.CursorAssistance)

	content, err := s.chatCompletion(config.ActiveModel, prompt, float64(config.Temperature))
	if err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", err
	}
	content = cleanJSON(content)
	var result struct {
		Recommendation string `json:"recommendation"`
		Confidence     int    `json:"confidence"`
		Reasoning      string `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "REJECT", 50, "AI ไม่แน่ใจ - ใช้ดุลยพินิจ", nil
	}
	if result.Recommendation != "APPROVE" && result.Recommendation != "REJECT" {
		result.Recommendation = "REJECT"
		result.Confidence = 50
	}
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 100 {
		result.Confidence = 100
	}
	if strings.TrimSpace(result.Reasoning) == "" {
		if result.Recommendation == "APPROVE" {
			result.Reasoning = "นักพัฒนามีเหตุผลที่สมควร - แนะนำให้อนุมัติคำขอ"
		} else {
			result.Reasoning = "เหตุผลไม่เพียงพอ - ให้ใช้ค่าประมาณของ AI"
		}
	}
	return result.Recommendation, result.Confidence, result.Reasoning, nil
}

func (s *glmService) EstimateStoryPoints(ctx domain.StoryPointTaskContext, fewShotExamples []domain.StoryPointExample) (float64, int, domain.StoryPointFactors, string, error) {
	config := s.getConfig()

	// Build multi-turn conversation to avoid "Prompt exceeds max length"
	// Turn 1 (system): Instructions
	// Turn 2 (user): Few-shot examples for calibration
	// Turn 3 (assistant): Acknowledgment
	// Turn 4 (user): The actual task to estimate

	messages := []glmMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(`You are a Senior Agile Coach estimating Story Points.
Scale: Fibonacci (0, 0.5, 1, 2, 3, 5, 8, 13, 21).
3-Factor Model: Work(1-10), Complexity(1-10), Risk(1-10).
Stack: Go/Nuxt3/PostgreSQL. AI Assistance: %d%% - %s.
Story Points = relative effort (work + complexity + risk), NOT hours.
Key context signals:
- TaskType: FEATURE (client-facing, usually larger), TASK (dev work, medium), BUG (fixes, usually smaller but risk varies)
- Priority: CRITICAL/HIGH = more urgency & scrutiny → may increase risk factor
- SubTaskCount: more subtasks = more work & coordination
- HasParent: sub-tasks are usually smaller scope than parent
- EstimatedMin: existing time estimate (if >0) helps calibrate scale
Output JSON ONLY: {"story_points":<float>,"confidence":<int 0-100>,"factors":{"work_amount":<int>,"complexity":<int>,"risk":<int>},"reasoning":"<ภาษาไทย สั้นๆ>"}`, config.CursorAssistance, cursorContextFromAssistance(config.CursorAssistance)),
		},
	}

	// Turn 2: Few-shot examples (up to 10 — no truncation needed, separate message)
	if len(fewShotExamples) > 0 {
		exJSON, _ := json.Marshal(fewShotExamples)
		messages = append(messages, glmMessage{
			Role:    "user",
			Content: fmt.Sprintf("REFERENCE TASKS from the same project (calibrate your scale):\n%s", string(exJSON)),
		})
		messages = append(messages, glmMessage{
			Role:    "assistant",
			Content: "Understood. I'll use these reference tasks to calibrate my estimation scale. Please provide the task to estimate.",
		})
	}

	// Turn 4: The actual task — full context, description cleaned of base64/HTML
	taskMsg := fmt.Sprintf(`Task Title: %s
Task Type: %s
Priority: %s
Sub-tasks: %d
Is Sub-task: %v
Existing Time Estimate: %d minutes
Project: %s
Description: %s`, ctx.Title, ctx.TaskType, ctx.Priority, ctx.SubTaskCount, ctx.HasParent, ctx.EstimatedMin, ctx.ProjectName, stripDescriptionForAI(ctx.Description, 8000))
	messages = append(messages, glmMessage{
		Role:    "user",
		Content: taskMsg,
	})

	content, err := s.chatCompletionMessages(config.ActiveModel, messages, float64(config.Temperature))
	if err != nil {
		return 0, 0, domain.StoryPointFactors{}, "", err
	}
	content = cleanJSON(content)

	var result struct {
		StoryPoints float64                 `json:"story_points"`
		Confidence  int                     `json:"confidence"`
		Factors     domain.StoryPointFactors `json:"factors"`
		Reasoning   string                  `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return 0, 0, domain.StoryPointFactors{}, "", fmt.Errorf("glm returned invalid JSON: %w (response: %s)", err, content)
	}

	if result.StoryPoints <= 0 {
		return 0, 0, domain.StoryPointFactors{}, "", fmt.Errorf("glm returned invalid story points: %.1f", result.StoryPoints)
	}
	result.Factors.WorkAmount = clampInt(result.Factors.WorkAmount, 1, 10)
	result.Factors.Complexity = clampInt(result.Factors.Complexity, 1, 10)
	result.Factors.Risk = clampInt(result.Factors.Risk, 1, 10)

	fmt.Printf("✅ GLM SP Estimate: %.1f SP (confidence: %d%%, factors: W=%d C=%d R=%d)\n",
		result.StoryPoints, result.Confidence, result.Factors.WorkAmount, result.Factors.Complexity, result.Factors.Risk)

	return result.StoryPoints, result.Confidence, result.Factors, result.Reasoning, nil
}

// EstimateEffortFromContext estimates implementation time (minutes) using SP + factors as calibration context.
// This is a SEPARATE focused AI call from EstimateStoryPoints for maximum accuracy.
func (s *glmService) EstimateEffortFromContext(ctx domain.StoryPointTaskContext, storyPoints float64, factors domain.StoryPointFactors) (int, string, error) {
	config := s.getConfig()
	cursorContext := cursorContextFromAssistance(config.CursorAssistance)

	messages := []glmMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(`You are a Senior Software Architect estimating implementation time.
Stack: Go (Fiber/Gin), Nuxt 3, PostgreSQL, Hexagonal Architecture.
AI Assistance Level: %d%% - %s

CRITICAL — AI Assistance Adjustment:
- 80%%+ AI assistance: developer is ~2-3x faster. Multiply base time by 0.4-0.5
- 50%% AI assistance: moderate speedup. Multiply base time by 0.7-0.8
- 20%% or less: nearly manual coding. Use base time as-is (1.0x)

Story Points to Hours Reference (BASE time BEFORE AI adjustment, senior developer):
- 1 SP ≈ 4h | 2 SP ≈ 8h | 3 SP ≈ 12h (1.5 days)
- 5 SP ≈ 20h (2.5 days) | 8 SP ≈ 32h (4 days) | 13 SP ≈ 52h (6.5 days)
- 21 SP ≈ 84h (10.5 days)

Rules:
- Start from SP-to-hours reference, then adjust based on factors and AI assistance level.
- Include: coding, testing, code review, documentation. NOT meetings or waiting.
- If task has sub-tasks, this estimate covers the PARENT task coordination + integration.
Output JSON ONLY: {"minutes":<int>,"reasoning":"<คำอธิบายสั้นๆ เป็นภาษาไทย>"}`, config.CursorAssistance, cursorContext),
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Estimate implementation time for this task.

Task Title: %s
Task Type: %s
Priority: %s
Sub-tasks: %d
Is Sub-task: %v
Project: %s

AI already estimated:
- Story Points: %.1f (Fibonacci scale)
- Work Amount: %d/10
- Complexity: %d/10
- Risk: %d/10

Description (summary): %s`, ctx.Title, ctx.TaskType, ctx.Priority, ctx.SubTaskCount, ctx.HasParent, ctx.ProjectName,
				storyPoints, factors.WorkAmount, factors.Complexity, factors.Risk,
				stripDescriptionForAI(ctx.Description, 2000)),
		},
	}

	content, err := s.chatCompletionMessages(config.ActiveModel, messages, float64(config.Temperature))
	if err != nil {
		return 0, "", err
	}
	content = cleanJSON(content)

	var result struct {
		Minutes  int    `json:"minutes"`
		Reasoning string `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return 0, "", fmt.Errorf("glm effort returned invalid JSON: %w (response: %s)", err, content)
	}
	if result.Minutes <= 0 {
		return 0, "", fmt.Errorf("glm effort returned invalid estimation: %d minutes", result.Minutes)
	}

	fmt.Printf("⏱️ GLM Effort Estimate: %d min for %.1f SP task\n", result.Minutes, storyPoints)
	return result.Minutes, result.Reasoning, nil
}

// cursorContextFromAssistance returns a human-readable context string for the AI prompt.
func cursorContextFromAssistance(level int) string {
	switch {
	case level <= 20:
		return "The developer is coding mostly manually with minimal AI assistance (20%% or less). Expect slower implementation times similar to traditional development."
	case level <= 50:
		return "The developer uses AI moderately for code suggestions and debugging (~50%%). Estimate time with moderate AI speedup."
	case level <= 80:
		return "The developer heavily relies on AI tools (Cursor/Windsurf/Copilot) for boilerplate, refactoring, and debugging (~80%%). Expect significant time savings."
	default:
		return "The developer works in an AI-first workflow with near-full assistance (90%%+). Expect very aggressive time estimates - AI handles most implementation."
	}
}
