package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type geminiService struct {
	apiKey     string
	httpClient *http.Client
	repo       domain.SentinelRepository
	tracker    domain.UsageTracker // optional: record API calls for quota display
}

// NewGeminiService สร้าง Instance ใหม่ (ใช้ REST API แทน SDK). tracker may be nil.
func NewGeminiService(apiKey string, repo domain.SentinelRepository, tracker domain.UsageTracker) (domain.AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	return &geminiService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		repo:    repo,
		tracker: tracker,
	}, nil
}

// listModelsResponse matches Gemini API GET /v1beta/models response
type listModelsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
	NextPageToken string `json:"nextPageToken"`
}

func (s *geminiService) recordRequest() {
	if s.tracker != nil {
		s.tracker.RecordRequest()
	}
}

// ListModels calls Gemini List Models API and returns model IDs (e.g. gemini-2.5-flash-lite).
func (s *geminiService) ListModels() ([]string, error) {
	var all []string
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s&pageSize=100", s.apiKey)
	for {
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("list models: %w", err)
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("list models API %d: %s", resp.StatusCode, string(bodyBytes))
		}
		s.recordRequest()
		var data listModelsResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("parse list models: %w", err)
		}
		for _, m := range data.Models {
			name := strings.TrimPrefix(m.Name, "models/")
			if name != "" {
				all = append(all, name)
			}
		}
		if data.NextPageToken == "" {
			break
		}
		apiURL = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s&pageSize=100&pageToken=%s", s.apiKey, url.QueryEscape(data.NextPageToken))
	}
	return all, nil
}

// Gemini REST API Request/Response structures
type geminiRequest struct {
	Contents         []geminiContent   `json:"contents"`
	GenerationConfig *generationConfig `json:"generationConfig,omitempty"`
}

type generationConfig struct {
	Temperature float64 `json:"temperature"` // 0.0 = deterministic, 2.0 = creative
	TopK        int     `json:"topK"`        // Top-K sampling (1 = most deterministic)
	TopP        float64 `json:"topP"`        // Nucleus sampling
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// โครงสร้างสำหรับรับ Response จาก AI
type aiEstimationResponse struct {
	Minutes   int    `json:"minutes"`
	Reasoning string `json:"reasoning"`
}

// โครงสร้างสำหรับรับ Code Review Response จาก AI
type aiCodeReviewResponse struct {
	Verdict  string          `json:"verdict"`  // "PASS" or "FAIL"
	Score    int             `json:"score"`    // 0-100
	Feedback json.RawMessage `json:"feedback"` // Can be string or array
}

func (s *geminiService) EstimateEffort(title, description string) (int, string, error) {
	fmt.Printf("🧠 AI Estimation Request: %s\n", title)
	
	// 🔥 DYNAMIC CONFIG: Fetch current system configuration
	config, err := s.repo.GetSystemConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to fetch system config, using defaults: %v\n", err)
		// Fallback to defaults if config fetch fails
		config = &domain.SystemConfig{
			ActiveModel:      "gemini-2.5-flash-lite",
			Temperature:      0.4,
			CursorAssistance: 80,
		}
	}
	
	fmt.Printf("⚙️  AI Config: Model=%s, Temp=%.2f, Cursor=%d%%\n", 
		config.ActiveModel, config.Temperature, config.CursorAssistance)
	
	// 🎯 Dynamic Cursor Assistance Context
	cursorContext := ""
	if config.CursorAssistance <= 20 {
		cursorContext = "The developer is coding mostly manually with minimal AI assistance (20% or less). Expect slower implementation times similar to traditional development."
	} else if config.CursorAssistance <= 50 {
		cursorContext = "The developer uses AI moderately for code suggestions and debugging (~50%). Estimate time with moderate AI speedup."
	} else if config.CursorAssistance <= 80 {
		cursorContext = "The developer heavily relies on AI tools (Cursor/Windsurf/Copilot) for boilerplate, refactoring, and debugging (~80%). Expect significant time savings."
	} else {
		cursorContext = "The developer works in an AI-first workflow with near-full assistance (90%+). Expect very aggressive time estimates - AI handles most implementation."
	}
	
	// God-Tier Prompt Engineering with Dynamic Context
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

	// Build request with DYNAMIC generation config
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &generationConfig{
			Temperature: float64(config.Temperature), // 🔥 Dynamic temperature
			TopK:        1,
			TopP:        0.95,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 60, "JSON marshal error", err
	}

	// Call Gemini API with DYNAMIC MODEL 🔥
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", 
		config.ActiveModel, s.apiKey)
	fmt.Printf("📡 Calling Gemini API (model: %s, temp: %.2f, cursor: %d%%)\n", 
		config.ActiveModel, config.Temperature, config.CursorAssistance)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 60, "HTTP request error", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Printf("❌ Gemini API call failed: %v\n", err)
		return 0, "", fmt.Errorf("gemini API network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("❌ Gemini API error %d: %s\n", resp.StatusCode, string(bodyBytes))
		return 0, "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	s.recordRequest()
	fmt.Printf("✅ Gemini API responded: Status %d\n", resp.StatusCode)

	// Parse response
	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return 0, "", fmt.Errorf("failed to decode gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return 0, "", fmt.Errorf("gemini returned empty response (no candidates)")
	}

	// Extract text
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Clean Markdown if any (Gemini sometimes wraps in ```json ... ```)
	responseText = strings.ReplaceAll(responseText, "```json", "")
	responseText = strings.ReplaceAll(responseText, "```", "")
	responseText = strings.TrimSpace(responseText)

	// Parse AI estimation
	var result aiEstimationResponse
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		fmt.Printf("❌ AI JSON parse error: %v (text: %s)\n", err, responseText)
		return 0, "", fmt.Errorf("gemini returned invalid JSON: %w (response: %s)", err, responseText)
	}

	// Validate (กัน AI มั่วเลขติดลบหรือ 0)
	if result.Minutes <= 0 {
		return 0, "", fmt.Errorf("gemini returned invalid estimation: %d minutes", result.Minutes)
	}

	fmt.Printf("✅ AI Estimated: %d minutes (Reasoning: %s)\n", result.Minutes, result.Reasoning)
	return result.Minutes, result.Reasoning, nil
}

// EstimateAndScheduleTasks ประเมินเวลาและลำดับการทำของแต่ละ task จากข้อมูลที่มี
func (s *geminiService) EstimateAndScheduleTasks(inputs []domain.TaskEstimateInput) ([]domain.TaskEstimateAndOrder, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	fmt.Printf("🤖 AI Estimate & Schedule: %d tasks\n", len(inputs))

	config, err := s.repo.GetSystemConfig()
	if err != nil || config == nil {
		config = &domain.SystemConfig{
			ActiveModel: "gemini-2.5-flash-lite",
			Temperature: 0.4,
			CursorAssistance: 80,
		}
	}
	tasksJSON, _ := json.Marshal(inputs)
	prompt := fmt.Sprintf(`You are a Senior Technical PM. Given the following tasks of a project, do TWO things:
1) Estimate the implementation time in MINUTES for each task (Senior Dev, stack: Go, Nuxt 3, PostgreSQL). AI assistance level: %d%%.
2) Suggest the EXECUTION ORDER (1 = do first, 2 = second, ...) based on dependencies and priority. Higher priority and blocking work should have lower order number.

Tasks (JSON array):
%s

Output ONLY a JSON array. Each element: { "task_index": <0-based index>, "minutes": <int>, "order": <1-based execution order> }.
Use integers only. Example: [{"task_index":0,"minutes":120,"order":1},{"task_index":1,"minutes":60,"order":2}]
`, config.CursorAssistance, string(tasksJSON))

	reqBody := geminiRequest{
		Contents:         []geminiContent{{Parts: []geminiPart{{Text: prompt}}}},
		GenerationConfig: &generationConfig{Temperature: float64(config.Temperature), TopK: 1, TopP: 0.95},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", config.ActiveModel, s.apiKey)
	client := &http.Client{Timeout: 90 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini API: %w", err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API %d: %s", resp.StatusCode, string(bodyBytes))
	}
	s.recordRequest()
	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty gemini response")
	}
	rawText := geminiResp.Candidates[0].Content.Parts[0].Text
	cleaned := strings.TrimSpace(rawText)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)
	var results []domain.TaskEstimateAndOrder
	if err := json.Unmarshal([]byte(cleaned), &results); err != nil {
		return nil, fmt.Errorf("parse plan JSON: %w (text: %s)", err, cleaned)
	}
	fmt.Printf("✅ AI Estimate & Order: %d results\n", len(results))
	return results, nil
}

// ReviewCode ใช้ AI วิเคราะห์ code diff แบบ Senior Code Reviewer
func (s *geminiService) ReviewCode(diff string) (string, int, string, error) {
	fmt.Printf("🔍 AI Code Review Request (diff length: %d chars)\n", len(diff))

	// Handle empty diff
	if strings.TrimSpace(diff) == "" {
		return "PENDING", 0, "No code diff provided for review", nil
	}

	// 🔥 DYNAMIC CONFIG: Fetch current system configuration
	config, err := s.repo.GetSystemConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to fetch system config for code review, using defaults: %v\n", err)
		// Fallback to defaults if config fetch fails
		config = &domain.SystemConfig{
			ActiveModel:      "gemini-2.5-flash-lite",
			Temperature:      0.2,
			CursorAssistance: 80,
		}
	}
	
	fmt.Printf("⚙️  AI Code Review Config: Model=%s, Temp=%.2f\n", 
		config.ActiveModel, config.Temperature)

	// 🔒 ANTI-HALLUCINATION PROMPT: Force AI to understand input is SOURCE CODE, not SQL query
	prompt := fmt.Sprintf(`
╔════════════════════════════════════════════════════════════╗
║  ROLE: Expert Security Auditor reviewing Go source code   ║
╚════════════════════════════════════════════════════════════╝

🚨 INPUT CONTEXT - READ CAREFULLY:
The text below is a RAW SOURCE CODE SNIPPET from a git commit.
- It is NOT a database query string being executed.
- It is NOT user input being inserted into a database.
- It is the PROGRAM CODE ITSELF (Go/TypeScript/Vue).

🎯 YOUR MISSION:
Analyze the LOGIC and SECURITY of the code snippet below.

╔═══════════════════════════════════════════════════════════════════╗
║  CRITICAL RULES - ANTI-HALLUCINATION INSTRUCTIONS                  ║
╚═══════════════════════════════════════════════════════════════════╝

1. **DO NOT** treat the code as if it were a "string being inserted into SQL".
   The code IS the program logic. You are reviewing HOW it handles data.

2. **SQL INJECTION DETECTION RULES:**
   
   ✅ **SECURE (Score: 85-100)** - Parameterized Queries:
      • db.Where("user = ?", userInput)
      • db.Where("email = ?", email).First(&user)
      • db.Exec("UPDATE users SET name = $1 WHERE id = $2", name, id)
   
   → These use placeholders (? or $1). Database driver handles escaping.
   → This is the CORRECT way. DO NOT flag as SQL Injection.
   
   ❌ **INSECURE (Score: 0-30)** - String Concatenation:
      • query := "SELECT * FROM users WHERE name = '" + userName + "'"
      • query := fmt.Sprintf("DELETE FROM posts WHERE id = %%s", postID)
      • db.Raw("SELECT * FROM users WHERE email = '" + email + "'")
   
   → These build SQL strings dynamically with user input.
   → CRITICAL SQL INJECTION vulnerability.

3. **OTHER SECURITY CHECKS:**
   - Hardcoded secrets/API keys → FAIL (score < 40)
   - Missing error handling → Score < 70
   - XSS vulnerabilities (unescaped HTML) → FAIL
   - Command Injection (exec with user input) → FAIL

4. **CODE QUALITY:**
   - Clear variable names ✅
   - Proper error handling (if err != nil) ✅
   - Clean structure ✅
   - Spaghetti code → FAIL (score < 50)

5. **IGNORE:**
   - Missing imports (not shown in snippet)
   - Surrounding function context
   - JSON structure of how code was submitted

╔═══════════════════════════════════════════════════════════════════╗
║  CODE SNIPPET TO AUDIT                                             ║
╚═══════════════════════════════════════════════════════════════════╝

%s

╔═══════════════════════════════════════════════════════════════════╗
║  CONTEXT                                                           ║
╚═══════════════════════════════════════════════════════════════════╝
Stack: Go (Gin/Fiber/GORM), Nuxt 3, PostgreSQL, Hexagonal Architecture

╔═══════════════════════════════════════════════════════════════════╗
║  RESPONSE FORMAT                                                   ║
╚═══════════════════════════════════════════════════════════════════╝

Output JSON ONLY (no markdown, no extra text):
{
	"verdict": "PASS" or "FAIL",
	"score": <int 0-100>,
	"feedback": "<bullet points เป็นภาษาไทย อธิบาย verdict>"
}

**MANDATORY:**
- Write "feedback" in Thai language (ภาษาไทย) ONLY.
- If you see db.Where("col = ?", val), recognize it as SECURE parameterized query.
- Only flag SQL Injection if you see string concatenation building queries.
- Be strict on security, fair on quality.
`, diff)

	// Build request with DYNAMIC generation config
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &generationConfig{
			Temperature: float64(config.Temperature), // 🔥 Dynamic temperature
			TopK:        1,
			TopP:        0.95,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "FAIL", 0, "JSON marshal error", err
	}

	// Call Gemini API with DYNAMIC MODEL 🔥
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", 
		config.ActiveModel, s.apiKey)
	fmt.Printf("📡 Calling Gemini API for Code Review (model: %s, temp: %.2f)\n", 
		config.ActiveModel, config.Temperature)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "FAIL", 0, "HTTP request error", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Printf("❌ Gemini Code Review API call failed: %v\n", err)
		return "FAIL", 0, "", fmt.Errorf("gemini API network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Gemini Code Review API error %d: %s\n", resp.StatusCode, string(bodyBytes))
		return "FAIL", 0, "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	s.recordRequest()
	fmt.Printf("✅ Gemini Code Review API responded: Status %d\n", resp.StatusCode)

	// Parse response
	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "FAIL", 0, "", fmt.Errorf("failed to decode gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "FAIL", 0, "", fmt.Errorf("gemini returned empty response (no candidates)")
	}

	// Extract text
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Clean Markdown if any
	responseText = strings.ReplaceAll(responseText, "```json", "")
	responseText = strings.ReplaceAll(responseText, "```", "")
	responseText = strings.TrimSpace(responseText)

	// Parse AI code review
	var result aiCodeReviewResponse
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		fmt.Printf("❌ AI Code Review JSON parse error: %v (text: %s)\n", err, responseText)
		return "FAIL", 0, "", fmt.Errorf("gemini returned invalid JSON: %w (response: %s)", err, responseText)
	}

	// Validate verdict
	if result.Verdict != "PASS" && result.Verdict != "FAIL" {
		result.Verdict = "FAIL" // Default to FAIL if invalid
	}

	// Validate score (0-100)
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}

	// Convert feedback to string (handle both string and array from AI)
	var feedbackStr string
	if len(result.Feedback) > 0 {
		// Try to unmarshal as string first
		var strFeedback string
		if err := json.Unmarshal(result.Feedback, &strFeedback); err == nil {
			feedbackStr = strFeedback
		} else {
			// Try as array of strings
			var arrFeedback []string
			if err := json.Unmarshal(result.Feedback, &arrFeedback); err == nil {
				feedbackStr = strings.Join(arrFeedback, "\n")
			} else {
				// Fallback: use raw JSON
				feedbackStr = string(result.Feedback)
			}
		}
	}

	fmt.Printf("✅ AI Code Review Complete: %s (Score: %d/100)\n", result.Verdict, result.Score)
	fmt.Printf("📝 Feedback: %s\n", feedbackStr)

	return result.Verdict, result.Score, feedbackStr, nil
}

// AnalyzeAppeal uses AI to analyze an appeal's validity and provide recommendation
func (s *geminiService) AnalyzeAppeal(diff string, originalFeedback string, appealReason string) (string, int, string, error) {
	fmt.Printf("⚖️  AI Appeal Analysis Request\n")
	fmt.Printf("   • Diff Length: %d chars\n", len(diff))
	fmt.Printf("   • Original Feedback Length: %d chars\n", len(originalFeedback))
	fmt.Printf("   • Appeal Reason Length: %d chars\n", len(appealReason))

	// Validate inputs
	if strings.TrimSpace(appealReason) == "" {
		return "UPHOLD", 0, "No appeal reason provided", fmt.Errorf("appeal reason is empty")
	}

	// 🔥 DYNAMIC CONFIG: Fetch current system configuration
	config, err := s.repo.GetSystemConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to fetch system config for appeal analysis, using defaults: %v\n", err)
		// Fallback to defaults if config fetch fails
		config = &domain.SystemConfig{
			ActiveModel:      "gemini-2.5-flash-lite",
			Temperature:      0.3,
			CursorAssistance: 80,
		}
	}
	
	fmt.Printf("⚙️  AI Appeal Analysis Config: Model=%s, Temp=%.2f\n", 
		config.ActiveModel, config.Temperature)

	// Enhanced Appeal Analysis Prompt - Focus on CODE, not submission mechanism
	prompt := fmt.Sprintf(`คุณคือ Senior Code Auditor ที่ทำหน้าที่เป็นผู้พิพากษาในการพิจารณาอุทธรณ์

**CRITICAL: FOCUS ON CODE, NOT SUBMISSION MECHANISMS**
- **IGNORE** how the code was submitted (JSON payloads, API endpoints, diff strings)
- **DO NOT** criticize "input validation", "API design", or "submission mechanism"
- **ANALYZE ONLY** the actual code logic and security issues

**บริบท:**
โค้ดที่ส่งมาตรวจไม่ผ่าน AI security review นักพัฒนากำลังยื่นอุทธรณ์

**Original Code Snippet:**
---
%s
---

**รายงานข้อบกพร่องจาก AI เดิม:**
%s

**คำให้การของนักพัฒนา (เหตุผลในการอุทธรณ์):**
"%s"

**งานของคุณ:**
วิเคราะห์ว่าคำให้การของนักพัฒนามีเหตุผลหรือไม่ โดยพิจารณา:

1. **AI Review เดิมถูกต้องหรือไม่?**
   - ถ้า AI บ่นเรื่อง "JSON structure", "API endpoint", "diff validation" → False Positive (ควร OVERTURN)
   - ถ้า AI บ่นเรื่อง SQL Injection, Hardcoded Secrets ในโค้ด → ตรวจว่าจริงหรือไม่

2. **โค้ดมีช่องโหว่ความปลอดภัยจริงหรือไม่?**
   - ✅ SAFE: db.Where("user = ?", input) → Parameterized query
   - ❌ UNSAFE: "SELECT * FROM users WHERE name = '" + input + "'" → SQL Injection
   
3. **คำอธิบายของนักพัฒนามีเหตุผลทางเทคนิคหรือไม่?**

4. **มีปัจจัยบรรเทาหรือ False Positive หรือไม่?**

**คู่มือการตัดสินใจ:**
- **"OVERTURN"** = อนุมัติอุทธรณ์ (AI ตรวจผิด, นักพัฒนาพูดถูก)
  - ใช้เมื่อ: AI บ่นเรื่อง submission mechanism แทนที่จะวิเคราะห์โค้ด
  - ใช้เมื่อ: โค้ดใช้ Parameterized Queries อย่างถูกต้อง แต่ AI บอกว่า SQL Injection
  
- **"UPHOLD"** = ปฏิเสธอุทธรณ์ (AI ตัดสินถูกต้อง)
  - ใช้เมื่อ: มี SQL Injection, Hardcoded Secrets, XSS จริงๆ ในโค้ด
  - ใช้เมื่อ: นักพัฒนาอ้างผิด

**ระดับความมั่นใจ:**
- 90-100: มั่นใจมากในคำแนะนำ
- 70-89: มั่นใจ แต่มีรายละเอียดบางอย่างที่ต้องพิจารณา
- 50-69: มั่นใจปานกลาง ต้องใช้ดุลยพินิจของมนุษย์
- 0-49: มั่นใจต่ำ ต้องมนุษย์ตรวจสอบแน่นอน

**ตอบเป็น JSON ONLY (ไม่ต้องใส่ markdown หรือข้อความอื่น):**
{
	"recommendation": "OVERTURN" or "UPHOLD",
	"confidence": <int 0-100>,
	"reasoning": "<1-2 ประโยคเป็นภาษาไทย แนะนำ CEO/PM ในการพิจารณาอุทธรณ์นี้>"
}

**CRITICAL:** 
- Write "reasoning" in Thai language (ภาษาไทย) ONLY.
- Must be 1-2 sentences max, clear and actionable for CEO/PM.`, diff, originalFeedback, appealReason)

	// Build request with DYNAMIC generation config
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &generationConfig{
			Temperature: float64(config.Temperature), // 🔥 Dynamic temperature
			TopK:        1,
			TopP:        0.95,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "UPHOLD", 0, "JSON marshal error", err
	}

	// Call Gemini API with DYNAMIC MODEL 🔥
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", 
		config.ActiveModel, s.apiKey)
	fmt.Printf("📡 Calling Gemini API for Appeal Analysis (model: %s, temp: %.2f)\n", 
		config.ActiveModel, config.Temperature)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "UPHOLD", 0, "HTTP request error", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "UPHOLD", 0, "API call failed", err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "UPHOLD", 0, "Failed to read response", err
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Gemini API Error (Status: %d): %s\n", resp.StatusCode, string(body))
		return "UPHOLD", 0, "API returned error", fmt.Errorf("gemini API status %d: %s", resp.StatusCode, string(body))
	}
	s.recordRequest()
	// Parse Gemini response
	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		fmt.Printf("❌ Failed to parse Gemini response: %v\n", err)
		return "UPHOLD", 0, "Response parse error", err
	}

	// Extract text
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "UPHOLD", 0, "Empty AI response", fmt.Errorf("gemini returned empty response")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	fmt.Printf("📄 Raw AI Appeal Analysis: %s\n", responseText)

	// Clean Markdown
	responseText = strings.ReplaceAll(responseText, "```json", "")
	responseText = strings.ReplaceAll(responseText, "```", "")
	responseText = strings.TrimSpace(responseText)

	// Parse Appeal Analysis Response
	type appealAnalysisResponse struct {
		Recommendation string `json:"recommendation"` // OVERTURN or UPHOLD
		Confidence     int    `json:"confidence"`     // 0-100
		Reasoning      string `json:"reasoning"`      // Advice for CEO
	}

	var result appealAnalysisResponse
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		fmt.Printf("❌ AI JSON parse error: %v (text: %s)\n", err, responseText)
		return "UPHOLD", 0, "Invalid JSON response", fmt.Errorf("gemini returned invalid JSON: %w", err)
	}

	// Validate recommendation
	if result.Recommendation != "OVERTURN" && result.Recommendation != "UPHOLD" {
		fmt.Printf("⚠️  AI returned invalid recommendation: %s, defaulting to UPHOLD\n", result.Recommendation)
		result.Recommendation = "UPHOLD"
	}

	// Validate confidence
	if result.Confidence < 0 || result.Confidence > 100 {
		fmt.Printf("⚠️  AI returned invalid confidence: %d, capping to 0-100\n", result.Confidence)
		if result.Confidence < 0 {
			result.Confidence = 0
		} else if result.Confidence > 100 {
			result.Confidence = 100
		}
	}

	fmt.Printf("⚖️  AI Appeal Analysis Complete:\n")
	fmt.Printf("   • Recommendation: %s\n", result.Recommendation)
	fmt.Printf("   • Confidence: %d%%\n", result.Confidence)
	fmt.Printf("   • Reasoning: %s\n", result.Reasoning)

	return result.Recommendation, result.Confidence, result.Reasoning, nil
}

// AnalyzeTimeNegotiation uses AI to evaluate developer's time negotiation request
// Returns: (recommendation, confidence, reasoning, error)
// - recommendation: "APPROVE" (developer is right) or "REJECT" (AI estimate stands)
// - confidence: 0-100 (how confident AI is in recommendation)
// - reasoning: Thai explanation for CEO/PM
func (s *geminiService) AnalyzeTimeNegotiation(
	taskTitle string,
	taskDescription string,
	aiEstimate int,
	devProposal int,
	devReason string,
) (string, int, string, error) {
	fmt.Printf("⏱️  AI Time Negotiation Analysis Request\n")
	fmt.Printf("   • Task: %s\n", taskTitle)
	fmt.Printf("   • AI Estimate: %d min → Dev Proposal: %d min\n", aiEstimate, devProposal)
	
	// Validate inputs
	if strings.TrimSpace(devReason) == "" {
		return "REJECT", 0, "ไม่มีเหตุผล - ให้ใช้ค่าประมาณของ AI", fmt.Errorf("negotiation reason is empty")
	}

	if devProposal <= aiEstimate {
		return "REJECT", 100, "นักพัฒนาเสนอเวลาน้อยกว่าหรือเท่ากับ AI → ไม่ต้องเจรจา", nil
	}

	// 🔥 DYNAMIC CONFIG: Fetch current system configuration
	config, err := s.repo.GetSystemConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to fetch system config for time negotiation, using defaults: %v\n", err)
		config = &domain.SystemConfig{
			ActiveModel:      "gemini-2.5-flash-lite",
			Temperature:      0.3,
			CursorAssistance: 80,
		}
	}
	
	fmt.Printf("⚙️  AI Time Negotiation Config: Model=%s, Temp=%.2f, Cursor=%d%%\n", 
		config.ActiveModel, config.Temperature, config.CursorAssistance)

	// Calculate multiplier based on difference
	multiplier := float64(devProposal) / float64(aiEstimate)

	// Time Negotiation Analysis Prompt
	prompt := fmt.Sprintf(`คุณคือ Senior Project Manager ที่ทำหน้าที่ตรวจสอบคำขอเจรจาเวลาจากนักพัฒนา

**งาน (Task):**
ชื่อ: "%s"
รายละเอียด: "%s"

**ข้อมูลการประมาณเวลา:**
- 🤖 **AI คำนวณ:** %d นาที (%.1f ชั่วโมง)
- 👨‍💻 **นักพัฒนาเสนอ:** %d นาที (%.1f ชั่วโมง)
- 📊 **ส่วนต่าง:** %.1fx (นักพัฒนาขอเพิ่ม %.0f%%%%)

**เหตุผลจากนักพัฒนา:**
"%s"

**ระดับการใช้ AI ในทีม (Cursor Assistance): %d%%**
- 0-20%%: ทีมเขียนโค้ดเองเกือบทั้งหมด (Manual)
- 21-50%%: ใช้ AI ช่วยบางส่วน (Hybrid)
- 51-80%%: ใช้ AI ช่วยเยอะ (AI-Assisted)
- 81-100%%: ใช้ AI เป็นหลัก (AI-First)

**งานของคุณ:**
วิเคราะห์ว่าควรอนุมัติคำขอเพิ่มเวลาหรือไม่ โดยพิจารณา:

1. **เหตุผลของนักพัฒนามีหลักการหรือไม่?**
   ✅ เหตุผลที่ดี:
   - มี legacy code ซับซ้อนที่ต้องศึกษา
   - มี dependencies เยอะที่ต้อง integrate
   - Requirements ไม่ชัดเจน ต้องสื่อสารกับทีมอื่น
   - มีความเสี่ยงด้านความปลอดภัย ต้องทดสอบรอบคอบ
   - มี edge cases เยอะที่ต้องจัดการ
   
   ❌ เหตุผลที่ไม่ดี:
   - บ่นว่า "งานยาก" โดยไม่มีรายละเอียดเทคนิค
   - แค่บอกว่า "AI ประมาณต่ำไป" โดยไม่อธิบาย
   - ขอเวลาเพิ่มมากเกินไป (>3x) โดยไม่มีเหตุผลรองรับ

2. **ส่วนต่างเวลาเหมาะสมหรือไม่?**
   - 1.2-1.5x: เหมาะสมสำหรับงานที่มี unknown factors เล็กน้อย
   - 1.5-2.5x: เหมาะสมสำหรับงานซับซ้อนที่มี dependencies
   - >3x: ควรมีเหตุผลที่ชัดเจนมาก (เช่น ต้อง migrate architecture)

3. **ระดับ Cursor Assistance (%d%%) สมเหตุสมผลหรือไม่?**
   - ถ้าทีมใช้ AI เยอะ (>80%%) → AI estimate น่าจะถูกต้องอยู่แล้ว
   - ถ้าทีมเขียนเองเยอะ (<50%%) → อาจต้องเวลามากกว่า AI คิด

**คู่มือการตัดสินใจ:**
- **"APPROVE"** = อนุมัติคำขอ (นักพัฒนามีเหตุผล, งานซับซ้อนกว่า AI คิด)
  - ใช้เมื่อ: เหตุผลมีรายละเอียดเทคนิค + ส่วนต่างไม่มากเกินไป
  
- **"REJECT"** = ปฏิเสธคำขอ (ให้ใช้ค่าประมาณของ AI)
  - ใช้เมื่อ: เหตุผลไม่ชัดเจน หรือ ขอเวลาเพิ่มมากเกินไป

**ระดับความมั่นใจ:**
- 90-100: มั่นใจมากในคำแนะนำ
- 70-89: มั่นใจ แต่ควรพิจารณาบริบทเพิ่มเติม
- 50-69: มั่นใจปานกลาง ต้องใช้ดุลยพินิจของ PM/CEO
- 0-49: มั่นใจต่ำ ต้องมนุษย์ตรวจสอบแน่นอน

**ตอบเป็น JSON ONLY (ไม่ต้องใส่ markdown):**
{
	"recommendation": "APPROVE" or "REJECT",
	"confidence": <int 0-100>,
	"reasoning": "<1-2 ประโยคเป็นภาษาไทย แนะนำ PM/CEO>"
}

**CRITICAL:**
- Write "reasoning" in Thai language (ภาษาไทย) ONLY.
- Must be 1-2 sentences max, clear and actionable.
- Consider Cursor Assistance level (%d%%) in your analysis.`,
		taskTitle,
		taskDescription,
		aiEstimate,
		float64(aiEstimate)/60.0,
		devProposal,
		float64(devProposal)/60.0,
		multiplier,
		((float64(devProposal)-float64(aiEstimate))/float64(aiEstimate))*100,
		devReason,
		config.CursorAssistance,
		config.CursorAssistance,
		config.CursorAssistance,
	)

	// Build request with DYNAMIC generation config
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &generationConfig{
			Temperature: float64(config.Temperature), // 🔥 Dynamic temperature
			TopK:        1,
			TopP:        0.95,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 🔥 DYNAMIC MODEL: Use config.ActiveModel
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		config.ActiveModel, // 🔥 Dynamic model selection
		s.apiKey,
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("failed to call AI API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("AI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	s.recordRequest()
	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("failed to parse AI response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "REJECT", 0, "ระบบ AI ขัดข้อง", fmt.Errorf("empty AI response")
	}

	rawText := geminiResp.Candidates[0].Content.Parts[0].Text
	fmt.Printf("📄 Raw AI Response:\n%s\n", rawText)

	// Parse JSON response
	var result struct {
		Recommendation string `json:"recommendation"`
		Confidence     int    `json:"confidence"`
		Reasoning      string `json:"reasoning"`
	}

	// Clean markdown artifacts
	cleanedText := strings.TrimSpace(rawText)
	cleanedText = strings.TrimPrefix(cleanedText, "```json")
	cleanedText = strings.TrimPrefix(cleanedText, "```")
	cleanedText = strings.TrimSuffix(cleanedText, "```")
	cleanedText = strings.TrimSpace(cleanedText)

	if err := json.Unmarshal([]byte(cleanedText), &result); err != nil {
		return "REJECT", 50, "AI ไม่แน่ใจ - ใช้ดุลยพินิจในการตัดสิน", fmt.Errorf("failed to parse JSON: %w (raw: %s)", err, cleanedText)
	}

	// Validate recommendation
	if result.Recommendation != "APPROVE" && result.Recommendation != "REJECT" {
		result.Recommendation = "REJECT"
		result.Confidence = 50
	}

	// Ensure confidence is in range
	if result.Confidence < 0 {
		result.Confidence = 0
	} else if result.Confidence > 100 {
		result.Confidence = 100
	}

	// Default reasoning if empty
	if strings.TrimSpace(result.Reasoning) == "" {
		if result.Recommendation == "APPROVE" {
			result.Reasoning = "นักพัฒนามีเหตุผลที่สมควร - แนะนำให้อนุมัติคำขอ"
		} else {
			result.Reasoning = "เหตุผลไม่เพียงพอ - ให้ใช้ค่าประมาณของ AI"
		}
	}

	fmt.Printf("⏱️  AI Time Negotiation Analysis Complete:\n")
	fmt.Printf("   • Recommendation: %s\n", result.Recommendation)
	fmt.Printf("   • Confidence: %d%%\n", result.Confidence)
	fmt.Printf("   • Reasoning: %s\n", result.Reasoning)

	return result.Recommendation, result.Confidence, result.Reasoning, nil
}

// GenerateWorkPlan asks AI to generate a full work plan (epics, milestones, sprints, tasks) from project name and description.
func (s *geminiService) GenerateWorkPlan(projectName, projectDescription string) (*domain.AIGeneratedPlan, error) {
	fmt.Printf("🤖 AI Generate Work Plan: %s\n", projectName)

	config, err := s.repo.GetSystemConfig()
	if err != nil || config == nil {
		config = &domain.SystemConfig{
			ActiveModel: "gemini-2.5-flash-lite",
			Temperature: 0.5,
		}
	}

	today := time.Now().Format("2006-01-02")
	prompt := fmt.Sprintf(`You are a Senior Technical Project Manager. Generate a complete work plan for this project.

**Project Name:** %s
**Project Description:** %s

**Tech stack context:** Go (Fiber/Gin), Nuxt 3, PostgreSQL, Hexagonal Architecture.

**CRITICAL - Date rule:** Today is %s. ALL dates you generate MUST be on or after today. Do NOT use any date in the past (no 2024 or earlier, and no dates before %s). Start the first sprint from today or the next Monday. Use only the current year and future.

Output a single JSON object with exactly these keys: epics, milestones, sprints, tasks.

**Rules:**
1. epics: array of { "title", "description", "color" }. Use hex colors like #6366f1. Create 2-5 epics that group major features.
2. milestones: array of { "title", "description", "due_date" }. due_date in YYYY-MM-DD, must be >= %s. Create 3-6 key milestones.
3. sprints: array of { "name", "goal", "start_date", "end_date" }. Dates in YYYY-MM-DD. First sprint start_date must be >= %s. Create 3-6 two-week sprints (14 days apart).
4. tasks: array of { "title", "description", "priority", "story_points", "epic_index", "sprint_index", "milestone_index", "start_date", "end_date" }.
   - priority: one of CRITICAL, HIGH, MEDIUM, LOW.
   - story_points: 1-5.
   - epic_index, sprint_index, milestone_index: integers only, 0-based index (use null if not linked). Do not use decimals.
   - start_date, end_date: YYYY-MM-DD, must be >= %s.
   Create 8-20 tasks spread across epics and sprints. Make titles and descriptions concrete and technical.

Output ONLY valid JSON, no markdown or explanation. Use integers for all numeric fields (e.g. 0 not 0.0). All dates must be today or future, never past.
`, projectName, projectDescription, today, today, today, today, today)

	reqBody := geminiRequest{
		Contents: []geminiContent{{Parts: []geminiPart{{Text: prompt}}}},
		GenerationConfig: &generationConfig{Temperature: float64(config.Temperature), TopK: 1, TopP: 0.95},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", config.ActiveModel, s.apiKey)
	client := &http.Client{Timeout: 90 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini API: %w", err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API %d: %s", resp.StatusCode, string(bodyBytes))
	}
	s.recordRequest()
	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("empty gemini response (no candidates)")
	}
	parts := geminiResp.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty gemini response (no content parts)")
	}
	rawText := parts[0].Text
	cleaned := strings.TrimSpace(rawText)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	// Use flexible struct so JSON numbers can be int or float (Gemini sometimes returns 0.0)
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
	if err := json.Unmarshal([]byte(cleaned), &flex); err != nil {
		return nil, fmt.Errorf("parse plan JSON: %w (text: %s)", err, cleaned)
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
	fmt.Printf("✅ AI Plan: %d epics, %d milestones, %d sprints, %d tasks\n",
		len(plan.Epics), len(plan.Milestones), len(plan.Sprints), len(plan.Tasks))
	return plan, nil
}
