# Master Planner — Real-World Examples

## Example 1: Adding Task Time Estimation with AI

**Scenario**: User asks "Add AI time estimation to tasks"

### Phase 1: Intent Analysis
```
Surface Request: Add AI-powered time estimation to tasks
Underlying Need: Project managers need reliable time estimates without manual input
Success Metric: AI estimates are within 20% of actual time, 80% of the time
Failure Mode: Estimates are wildly off → users lose trust and ignore them
User Context: Project managers (non-technical), developers (technical)
```

### Phase 2: Key Perspective Insights
```
Architecture: Use existing task data as training context; don't over-engineer ML pipeline
UX: Estimate should appear automatically when task is created/updated, not require action
Data: Need to store estimated vs actual for future improvement (feedback loop)
Edge Cases: New project with no history → need reasonable defaults
Performance: AI call must be async — don't block task creation
Security: Don't send sensitive project data to external AI without consent
Evolution: Later: estimate accuracy dashboard, team velocity predictions
```

### Phase 3: Approaches

**Approach A: AI API Call (Recommended)**
- On task create/update → call OpenAI with task description + project context
- Store estimate in `ai_estimated_minutes` field
- Pros: Smart estimates, improves with context
- Cons: External dependency, cost per call
- Effort: M | Risk: Low | Reversibility: Easy

**Approach B: Historical Average**
- Calculate average task duration from completed tasks in same project
- Pros: No external dependency, free
- Cons: Less accurate for new types of work
- Effort: S | Risk: Low | Reversibility: Easy

**Approach C: Hybrid**
- Use historical average as baseline, enhance with AI for complex tasks
- Pros: Best of both worlds
- Cons: More complex implementation
- Effort: L | Risk: Medium | Reversibility: Moderate

### Phase 5: Implementation Blueprint (Approach A)
```
Phase 1: Data Layer
- Add `ai_estimated_minutes` column to tasks table
- Add `estimation_context` JSONB column for AI prompt context
- Deliverable: DB migration, model updated

Phase 2: Estimation Service
- Create `internal/service/estimator/ai.go`
- Implement prompt engineering with project context
- Add async estimation via goroutine (don't block task creation)
- Deliverable: Estimation service with AI integration

Phase 3: Integration
- Hook into task create/update use case
- Store estimate result in task
- Deliverable: Estimates auto-generated on task create

Phase 4: Frontend Display
- Show estimate in task card and detail view
- Add loading state while estimating
- Deliverable: Users see AI estimates in UI
```

---

## Example 2: Building a Real-Time Dashboard

**Scenario**: User asks "Build a real-time dashboard for project status"

### Phase 1: Intent Analysis
```
Surface Request: Real-time project dashboard
Underlying Need: Managers need instant visibility into project health without refreshing
Success Metric: Dashboard updates within 2 seconds of any change
Failure Mode: Dashboard is stale or shows incorrect data
User Context: Multiple users viewing simultaneously, high concurrency
```

### Phase 2: Key Perspective Insights
```
Architecture: WebSocket with pub/sub pattern — existing pattern in codebase?
UX: Cards/panels for: active tasks, completion rate, blockers, team velocity
Data: Need to track what each connected user is subscribed to
Edge Cases: What if 100 users connect? What if DB is under load?
Performance: Cache dashboard state, push only deltas
Security: Users must only see projects they have access to
Observability: Track: connection count, message latency, disconnection rate
Evolution: Later: customizable widgets, export, alerts
```

### Phase 3: Approaches

**Approach A: WebSocket with Redis Pub/Sub**
- Backend publishes events to Redis channels
- WebSocket servers subscribe and push to connected clients
- Pros: Scales horizontally, proven pattern
- Cons: Redis dependency for pub/sub
- Effort: M | Risk: Low

**Approach B: Server-Sent Events (SSE)**
- Simpler than WebSocket, unidirectional
- Client subscribes to event stream
- Pros: Simpler implementation, auto-reconnect
- Cons: Unidirectional only, connection limits per browser
- Effort: S | Risk: Low

**Approach C: WebSocket Direct**
- Direct WebSocket connections to Go server
- In-process fan-out to connected clients
- Pros: No Redis dependency, lowest latency
- Cons: Doesn't scale horizontally without sticky sessions
- Effort: M | Risk: Medium

### Phase 5: Implementation Blueprint (Approach A)
```
Phase 1: WebSocket Infrastructure
- Set up WebSocket hub in Go (upgrade HTTP → WS)
- Implement Redis pub/sub subscriber
- Deliverable: Clients can connect and receive messages

Phase 2: Event Publishing
- Publish events on: task create/update/delete, status change
- Include project ID in event for routing
- Deliverable: Changes trigger events

Phase 3: Dashboard Data API
- REST endpoint for initial dashboard state
- Client loads initial data, then subscribes to WS for updates
- Deliverable: Full dashboard data available

Phase 4: Frontend Dashboard
- Dashboard page with real-time updates
- Connection status indicator
- Reconnection handling
- Deliverable: Live dashboard in browser
```

---

## Example 3: Adding Code Audit Submission System

**Scenario**: User asks "Add code audit feature where developers submit code for review"

### Phase 1: Intent Analysis
```
Surface Request: Code audit submission system
Underlying Need: Track code quality over time, provide AI-powered audit feedback
Success Metric: Developers submit code and receive actionable audit results
Failure Mode: Audit results are too generic to be useful
User Context: Developers (technical), tech leads (review results)
```

### Phase 2: Key Perspective Insights
```
Architecture: New "submissions" entity; store audit results as JSONB (matches existing pattern)
UX: Drag-and-drop file upload → processing state → audit results display
Data: Submissions table with: code, language, audit_result (JSONB), score
Edge Cases: Large files → size limit; binary files → reject; rate limiting
Performance: AI audit is slow → must be async with progress tracking
Security: Don't store secrets accidentally pasted in code; sanitize inputs
Evolution: Later: audit history, comparison, team-level aggregation
```

### Phase 3: Approaches

**Approach A: Async Queue with Status Polling**
- Submit → queue job → process → update status → poll/WebSocket notify
- Pros: Handles long processing, scalable
- Cons: More infrastructure (queue)
- Effort: L | Risk: Medium

**Approach B: Simple Async Goroutine**
- Submit → spawn goroutine → process → update DB → WebSocket notify
- Pros: No queue infrastructure needed
- Cons: Lost jobs on server restart, doesn't scale to multiple instances
- Effort: M | Risk: Medium

**Approach C: Synchronous with Timeout**
- Submit → wait for result → return
- Pros: Simplest implementation
- Cons: Bad UX for long audits, timeout issues
- Effort: S | Risk: High (UX risk)

### Phase 5: Implementation Blueprint (Approach A)
```
Phase 1: Data Layer
- Create `submissions` table: id, user_id, code, language, status, audit_result (JSONB), created_at
- Create GORM model and repository
- Deliverable: Migration + model + repository

Phase 2: Submission Endpoint
- POST /api/submissions — accept code + language
- Validate input (size limit, language support)
- Set status = "pending"
- Deliverable: Working submission API

Phase 3: Audit Processing
- Create audit service that calls AI with code
- Parse AI response into structured audit result
- Store result as JSONB
- Deliverable: Audit results generated and stored

Phase 4: Status & Results API
- GET /api/submissions/:id — return status + results
- GET /api/submissions — list user's submissions (paginated)
- Deliverable: Full CRUD API

Phase 5: Frontend
- Upload form with drag-and-drop
- Processing state with progress indicator
- Audit results display with score, issues, suggestions
- Deliverable: Complete UI flow
```

---

## Example 4: Quick Plan (Small Feature)

**Scenario**: User asks "Add due date to tasks"

```
## Quick Plan: Task Due Dates

Goal: Add optional due date field to tasks with overdue indicators
User Impact: Users can set deadlines and see which tasks are overdue

Approach: Add `due_date` column (nullable DATE) to tasks, expose in API and UI

Tasks:
1. `api/internal/repository/task/model.go` → Add `DueDate *time.Time` field — storage
2. `api/internal/repository/task/migration.go` → Add migration for `due_date` column — DB change
3. `api/internal/handler/task/dto.go` → Add `due_date` to request/response DTOs — API contract
4. `web/core/modules/tasks/components/TaskForm.vue` → Add date picker — user input
5. `web/core/modules/tasks/components/TaskCard.vue` → Show overdue indicator — visual feedback

Acceptance:
- [ ] User can set due date when creating/editing a task
- [ ] Overdue tasks show red indicator in task list
- [ ] API accepts/returns due_date in ISO 8601 format
- [ ] Tasks without due date are unaffected

Risk: Low — additive change, no breaking changes to existing functionality
```
