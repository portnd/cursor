# Reverse Engineering & Analysis — Example Outputs

## Example 1: Full System Analysis + Reverse Engineering

```
Input: "Analyze and reverse engineer the submission module"

## Forensic Analysis Report: Submission Module

### Executive Summary
The submission module handles code submission intake, AI-powered estimation, and audit
storage. The original author was a mid-level Go developer with a REST background, who
built this iteratively over ~2 weeks. The core is well-structured (Hexagonal), but
the AI integration was bolted on later and bypasses the usecase layer.

---

### Part I: System Analysis

#### Overview
- **Tech Stack**: Go (Gin), PostgreSQL (GORM), Redis (cache), OpenAI API
- **Architecture**: Hexagonal (handlers → usecases → repositories), 70% adherence
- **Scale**: 12 files, ~1,800 lines, 3 layers

#### Architecture Map
```
HTTP Request
  → Router (routes/submission.go)
    → Auth Middleware (checks JWT)
      → Handler (handlers/submission.go)
        → Usecase (usecases/submission.go)
          → Repository (repositories/submission.go)
            → PostgreSQL (submissions table)
          → AI Service (services/ai.go) ← BYPASS: called directly from handler
          → Cache (repositories/cache.go)
            → Redis
```

#### Dependency Map
- **Internal**: handler → usecase → repo → DB (clean), handler → AI service (VIOLATION)
- **External**: gin v1.9, gorm v1.25, go-redis v9, openai-go v1.0
- **API Contracts**: POST /api/v1/submissions, GET /api/v1/submissions/:id, GET /api/v1/submissions (list)

#### Data Flow Map
```
1. Client sends POST with {task_id, code, language}
2. Handler validates input, checks rate limit
3. Usecase creates submission record (status: pending)
4. AI service estimates time (called from handler — WRONG LAYER)
5. Usecase updates submission with ai_estimated_minutes
6. Usecase runs audit analysis
7. Usecase stores audit result as JSONB in submissions table
8. Return submission with estimation + audit

List endpoint: cursor-based pagination, HMAC-signed cursor
```

#### Quality Findings
| Category | Finding | Severity | Location |
|----------|---------|----------|----------|
| Architecture | AI service called from handler, bypasses usecase | 🔴 | handlers/submission.go:87-112 |
| Performance | No timeout on AI service call | 🔴 | handlers/submission.go:95 |
| Performance | N+1 query when listing submissions with tasks | 🟡 | repositories/submission.go:67 |
| Security | Rate limit only checks by user_id, not by IP | 🟡 | handlers/submission.go:52 |
| Quality | Magic number 30 (days expiry) should be constant | 🟢 | repositories/submission.go:23 |
| Quality | Missing error wrapping in repo layer | 🟢 | repositories/submission.go:45 |

#### Patterns & Conventions
- Error wrapping with `fmt.Errorf("functionName: %w", err)` — consistent in usecase/repo, inconsistent in handler
- Cursor-based pagination with HMAC — implemented correctly per project standards
- Response format: `{data: ..., meta: {next_cursor: ...}}`
- Validation done in handler before reaching usecase

---

### Part II: Author Intent & Reverse Engineering

#### Architecture Intent
- **Original Design Goal**: Clean submission pipeline with AI-powered code analysis
- **Chosen Pattern**: Hexagonal Architecture (per project standard)
- **Deviations**: AI service bypasses usecase layer — was likely added as a quick feature after initial build

#### Author Profile
- **Experience Level**: Mid — Evidence: Understands interfaces and layering but breaks own rules under pressure
- **Primary Language**: Go (native thinker) — Evidence: Proper context usage, idiomatic error wrapping
- **Development Style**: Iterative with rushed final additions — Evidence: Clean core (first 5 commits), messy AI integration (last 3 commits)
- **Blind Spots**: Concurrency (no mutex), Timeouts (no context deadline on external calls), Pagination edge cases

#### Discovered Business Rules
| Rule | Location | Confidence | Evidence |
|------|----------|------------|----------|
| Max 10 submissions per user per day | handler.go:52 | High | Explicit rate limit check |
| AI estimation is optional (graceful degradation) | handler.go:87 | High | nil check on AI service with fallback |
| Audit results expire after 30 days | repo.go:23 | Medium | Hardcoded constant, unclear if enforced |
| Code must be > 0 and < 1MB | handler.go:45 | High | Size validation check |
| Only task owner can submit | usecase.go:34 | High | Ownership check before creation |

#### Assumption Map
| Assumption | Evidence | Valid? | Impact if Wrong |
|------------|----------|--------|-----------------|
| "AI service responds in < 5s" | No context timeout | ❌ | Request hangs indefinitely |
| "Redis is always available" | No fallback on cache miss | ❌ | Submission list fails entirely |
| "User context is always set" | Direct context.Value without check | ❌ | Panic if middleware skipped |
| "Submissions list is small" | No cursor cap enforcement | 🟡 | Could return huge responses |

#### Hidden Contracts
- **Cache dependency**: handler.go:78 assumes Redis is warm; no cache-aside pattern for new keys
- **Auth middleware required**: handler expects `user_id` in Gin context; will panic with nil pointer if middleware is misconfigured
- **Task must exist before submission**: usecase.go:30 checks task existence but doesn't lock it — race condition possible

#### Decision Timeline
1. FIRST: Core CRUD with Hexagonal pattern (clean, 5 commits over 3 days)
2. THEN: Added cursor-based pagination (well-implemented, follows project standard)
3. LATER: Added AI estimation — bypassed usecase layer (2 commits, rushed)
4. LAST: Added rate limiting (single commit, clean implementation)

---

### Part III: Risk & Recommendations

#### Risk Map
| Area | Fragility | Why | Safe Change Approach |
|------|-----------|-----|---------------------|
| AI integration | 🔴 | Bypasses usecase, no timeout, no circuit breaker | Move to usecase, add context deadline + retry |
| Auth context | 🔴 | Panics if middleware skips | Add nil check before accessing context values |
| Cache layer | 🟡 | No fallback if Redis down | Add cache-aside pattern with DB fallback |
| Rate limiter | 🟢 | Isolated, clean | Safe to modify independently |

#### Complexity Hotspots
| Location | Type | Effort to Understand | Key Insight |
|----------|------|---------------------|-------------|
| handler.go:87-112 | Layer violation + missing timeout | 2h | The AI call was a rush job — it needs to move to usecase |
| usecase.go:45-78 | Nested business rules | 1h | Clean but dense — each step is a separate concern |
| repo.go:55-70 | Cursor pagination + N+1 | 1h | Pagination is correct, but eager loading is missing |

#### Recommendations
1. 🔴 Move AI service call from handler to usecase layer; add context.WithTimeout (30s)
2. 🔴 Add nil check for user context in handler; return 401 instead of panicking
3. 🟡 Fix N+1 query in list — use Preload or JOIN for task data
4. 🟡 Add Redis fallback / circuit breaker for cache failures
5. 🟢 Extract magic numbers (expiry days, rate limits) to constants

#### Mental Model Reconstruction
"The author wanted a clean, Hexagonal submission pipeline. They started right —
handlers for HTTP concerns, usecases for business logic, repos for data. The CRUD
and pagination are textbook implementations.

Then AI estimation was requested as a 'quick addition.' Instead of threading it
through the usecase layer properly, they called the AI service directly from the
handler to save time. This is the main architectural debt — it works, but it
breaks the layering and makes the handler do too much.

The biggest operational risk is the unprotected AI call with no timeout. If the
OpenAI API is slow, every submission request will hang. This should be the first
thing to fix."

#### System Health Score
| Dimension | Score | Reason |
|-----------|-------|--------|
| Architecture Clarity | 🟡 | Core is clean, AI integration breaks layering |
| Code Quality | 🟢 | Good naming, consistent patterns in original code |
| Test Coverage | 🟡 | CRUD is tested, AI path is not |
| Security Posture | 🟡 | Auth works, but context access is unsafe |
| Performance | 🟡 | Pagination works, N+1 issue on list |
| Maintainability | 🟢 | Clear structure, easy to navigate |
```

## Example 2: Quick Function-Level Reverse Engineering

```
Input: "What is this CalculateAI Estimate function doing?"

## Reverse Engineering Report: CalculateEstimation

### Purpose Reconstruction
This function estimates how long a coding task will take, using a combination of
code complexity metrics and historical data. It's the core algorithm of the
"AI estimation" feature.

### Control Flow (simplified)
```
Input: {code, language, task_description}
  ↓
1. Tokenize code → count lines, functions, complexity score
  ↓
2. Match language to base multiplier (Go=1.2x, Python=0.9x, etc.)
  ↓
3. Look up similar historical tasks (by complexity range)
  ↓
4. IF historical data exists → weighted average (70% historical, 30% heuristic)
   IF no historical data → pure heuristic estimate
  ↓
5. Apply task priority modifier (urgent = 1.3x, normal = 1.0x)
  ↓
6. Clamp to [15 min, 480 min] (15min to 8hrs)
  ↓
Output: {estimated_minutes, confidence_score, method_used}
```

### Hidden Business Rules
- Line 45: If code has > 500 lines, confidence drops to "low" — this is an undocumented threshold
- Line 62: Historical matching uses ±20% complexity range — this window is hardcoded and arbitrary
- Line 78: The "urgent" modifier (1.3x) means urgent tasks are estimated to take LONGER, not shorter — this is intentional (less available focus time)

### Author's Mental State
"The author built this incrementally:
1. First, a simple heuristic (line count × language multiplier)
2. Then added historical matching (the `if historicalData` branch)
3. The priority modifier was added last (different indentation style)

The 70/30 weighting was likely a business decision, not a data-driven one.
The clamping range (15-480 min) suggests the product only handles tasks
that fit within a single work day."
```

## Example 3: Architecture Heat Map Output

```
Input: "Give me a heat map of the entire API module"

## Architecture Heat Map

🔴 HOT (Fragile / High Risk):
- handlers/submission.go: AI integration bypasses usecase, no timeout, rate limiter incomplete
- middleware/auth.go: JWT validation has a time window vulnerability (accepts expired tokens within 5min grace)
- repositories/submission.go: N+1 query on list, raw SQL mixed with GORM

🟡 WARM (Needs Attention):
- usecases/submission.go: Business rules are dense but correct — needs test coverage
- routes/routes.go: Route registration is manual and error-prone (no auto-discovery)
- services/ai.go: Retry logic exists but no exponential backoff

🟢 COOL (Healthy):
- models/submission.go: Clean GORM model, proper tags, well-documented
- repositories/cache.go: Solid cache-aside pattern, proper TTL handling
- middleware/cors.go: Standard CORS setup, properly configured
- utils/cursor.go: HMAC cursor implementation is correct and well-tested

❄️ COLD (Dead / Unused):
- handlers/legacy_handler.go: Old v1 API endpoints, not called by frontend anymore
- repositories/search.go: ElasticSearch integration started but never completed
- services/notification.go: Notification service stub, all functions return nil
```
