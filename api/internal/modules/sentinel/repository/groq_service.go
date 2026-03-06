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

const defaultGroqModel = "llama-3.3-70b-versatile"
const groqChatURL = "https://api.groq.com/openai/v1/chat/completions"
const groqModelsURL = "https://api.groq.com/openai/v1/models"

type groqService struct {
	apiKey     string
	httpClient *http.Client
	repo       domain.SentinelRepository
	tracker    domain.UsageTracker
}

// NewGroqService creates an AIService that uses Groq API. tracker may be nil.
func NewGroqService(apiKey string, repo domain.SentinelRepository, tracker domain.UsageTracker) (domain.AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is required")
	}
	return &groqService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
		repo:    repo,
		tracker: tracker,
	}, nil
}

func (s *groqService) recordRequest() {
	if s.tracker != nil {
		s.tracker.RecordRequest()
	}
}

// groqModelsResponse matches GET /openai/v1/models
type groqModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (s *groqService) ListModels() ([]string, error) {
	req, err := http.NewRequest("GET", groqModelsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("groq list models: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("groq list models API %d: %s", resp.StatusCode, string(body))
	}
	s.recordRequest()
	var data groqModelsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse groq models: %w", err)
	}
	var ids []string
	for _, m := range data.Data {
		if m.ID != "" {
			ids = append(ids, m.ID)
		}
	}
	return ids, nil
}

// groqChatRequest for POST /openai/v1/chat/completions
type groqChatRequest struct {
	Model       string          `json:"model"`
	Messages    []groqMessage   `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// chatCompletion sends a single user message and returns the assistant content.
func (s *groqService) chatCompletion(model string, userPrompt string, temperature float64) (string, error) {
	if model == "" || strings.HasPrefix(model, "gemini-") {
		model = defaultGroqModel
	}
	reqBody := groqChatRequest{
		Model:       model,
		Messages:    []groqMessage{{Role: "user", Content: userPrompt}},
		Temperature: temperature,
		MaxTokens:   4096,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", groqChatURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("groq API: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API %d: %s", resp.StatusCode, string(body))
	}
	s.recordRequest()
	var chatResp groqChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("parse groq response: %w", err)
	}
	if len(chatResp.Choices) == 0 || chatResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("groq returned empty content")
	}
	return chatResp.Choices[0].Message.Content, nil
}

func (s *groqService) getConfig() *domain.SystemConfig {
	config, err := s.repo.GetSystemConfig()
	if err != nil || config == nil {
		return &domain.SystemConfig{
			ActiveModel:      defaultGroqModel,
			Temperature:      0.4,
			CursorAssistance: 80,
		}
	}
	return config
}

func cleanJSON(s string) string {
	s = strings.ReplaceAll(s, "```json", "")
	s = strings.ReplaceAll(s, "```", "")
	return strings.TrimSpace(s)
}

// EstimateEffort implements AIService.
func (s *groqService) EstimateEffort(title, description string) (int, string, error) {
	config := s.getConfig()
	cursorContext := "The developer uses AI moderately (~50%)."
	if config.CursorAssistance <= 20 {
		cursorContext = "The developer codes mostly manually. Estimate without AI speedup."
	} else if config.CursorAssistance <= 50 {
		cursorContext = "The developer uses AI moderately for suggestions and debugging (~50%)."
	} else if config.CursorAssistance <= 80 {
		cursorContext = "The developer heavily relies on AI tools (~80%). Expect significant time savings."
	} else {
		cursorContext = "The developer works in an AI-first workflow (90%+). Expect aggressive time estimates."
	}
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
		return 0, "", fmt.Errorf("groq returned invalid JSON: %w (response: %s)", err, content)
	}
	if result.Minutes <= 0 {
		return 0, "", fmt.Errorf("groq returned invalid estimation: %d minutes", result.Minutes)
	}
	return result.Minutes, result.Reasoning, nil
}

// EstimateAndScheduleTasks implements AIService.
func (s *groqService) EstimateAndScheduleTasks(inputs []domain.TaskEstimateInput) ([]domain.TaskEstimateAndOrder, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	config := s.getConfig()
	tasksJSON, _ := json.Marshal(inputs)
	prompt := fmt.Sprintf(`You are a Senior Technical PM. Given the following tasks of a project, do TWO things:
1) Estimate the implementation time in MINUTES for each task (Senior Dev, stack: Go, Nuxt 3, PostgreSQL). AI assistance level: %d%%.
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
		return nil, fmt.Errorf("parse plan JSON: %w (text: %s)", err, content)
	}
	return results, nil
}

// GenerateWorkPlan implements AIService.
func (s *groqService) GenerateWorkPlan(projectName, projectDescription string) (*domain.AIGeneratedPlan, error) {
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
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Priority        string   `json:"priority"`
		StoryPoints     float64  `json:"story_points"`
		EpicIndex       *float64 `json:"epic_index"`
		SprintIndex     *float64 `json:"sprint_index"`
		MilestoneIndex  *float64 `json:"milestone_index"`
		StartDate       string   `json:"start_date"`
		EndDate         string   `json:"end_date"`
	}
	type flexPlan struct {
		Epics      []domain.AIPlanEpic      `json:"epics"`
		Milestones []domain.AIPlanMilestone `json:"milestones"`
		Sprints    []domain.AIPlanSprint    `json:"sprints"`
		Tasks      []flexTask               `json:"tasks"`
	}
	var flex flexPlan
	if err := json.Unmarshal([]byte(content), &flex); err != nil {
		return nil, fmt.Errorf("parse work plan JSON: %w (text: %s)", err, content)
	}
	plan := &domain.AIGeneratedPlan{
		Epics:      flex.Epics,
		Milestones: flex.Milestones,
		Sprints:    flex.Sprints,
		Tasks:      make([]domain.AIPlanTask, 0, len(flex.Tasks)),
	}
	for _, t := range flex.Tasks {
		sp := int(t.StoryPoints)
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

// ReviewCode implements AIService (simplified prompt; same JSON contract).
func (s *groqService) ReviewCode(diff string) (string, int, string, error) {
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
		return "FAIL", 0, "", fmt.Errorf("invalid JSON: %w", err)
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

// AnalyzeAppeal implements AIService.
func (s *groqService) AnalyzeAppeal(diff, originalFeedback, appealReason string) (string, int, string, error) {
	if strings.TrimSpace(appealReason) == "" {
		return "UPHOLD", 0, "No appeal reason provided", fmt.Errorf("appeal reason is empty")
	}
	config := s.getConfig()
	prompt := fmt.Sprintf(`คุณคือ Senior Code Auditor ที่ทำหน้าที่พิจารณาอุทธรณ์

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
		return "UPHOLD", 0, "", fmt.Errorf("invalid JSON: %w", err)
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

// AnalyzeTimeNegotiation implements AIService.
func (s *groqService) AnalyzeTimeNegotiation(taskTitle, taskDescription string, aiEstimate, devProposal int, devReason string) (string, int, string, error) {
	if strings.TrimSpace(devReason) == "" {
		return "REJECT", 0, "ไม่มีเหตุผล - ให้ใช้ค่าประมาณของ AI", fmt.Errorf("negotiation reason is empty")
	}
	if devProposal <= aiEstimate {
		return "REJECT", 100, "นักพัฒนาเสนอเวลาน้อยกว่าหรือเท่ากับ AI → ไม่ต้องเจรจา", nil
	}
	config := s.getConfig()
	multiplier := float64(devProposal) / float64(aiEstimate)
	prompt := fmt.Sprintf(`คุณคือ Senior PM ตรวจสอบคำขอเจรจาเวลาจากนักพัฒนา

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
