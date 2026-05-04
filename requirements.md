# Requirements: The Sentinel

## Overview

The Sentinel is an **AI-driven Development OS** — a project management and engineering quality platform that uses AI to estimate task effort, audit code, and enforce delivery standards. It is built for software teams and is operated by three roles: CEO, PM, and DEV.

---

## Problem Statement

Engineering teams lack objective, data-driven tooling to:
1. Estimate task effort accurately before work begins.
2. Audit code quality consistently without human reviewer bias.
3. Enforce accountability through structured appeal and negotiation workflows.
4. Visualise project health, velocity, and cost in one place.

The Sentinel solves all four problems in a single, integrated system.

---

## Goals

- Provide AI-generated time estimates for every task before development starts.
- Automate code quality review via AI (Groq or Gemini) on every submission.
- Give developers a structured path to appeal AI verdicts or negotiate estimates.
- Give PMs and CEOs full visibility into project status, team performance, and cost.
- Run entirely in Docker; the web app must be reachable at `http://localhost:3000`.

## Non-Goals

- Not a general-purpose issue tracker (no GitHub/Jira sync in current scope).
- Not a real-time code editor or IDE plugin.
- Not a billing or invoicing system (finance module is internal accounting only).
- No mobile-native app (web-only).

---

## Users / Personas

| Role | Description | Key Capabilities |
|------|-------------|-----------------|
| **CEO** | Full system access | User/team management, AI config, finance, leaderboard, all project data |
| **PM** | Project and team lead | Create/manage projects, tasks, sprints, epics, milestones; assign DEVs; approve tasks; PM-scoped leaderboard |
| **DEV** | Individual contributor | View/work on assigned tasks, submit code, log time, appeal verdicts, negotiate estimates |

---

## Scope

### In Scope (Current)

- Multi-project workspace with tasks, epics, sprints, milestones, and task dependencies.
- AI time estimation per task (single and batch/schedule).
- AI code audit: submit diff → AI returns PASS/FAIL verdict, score (0–100), and structured feedback.
- Appeal system: DEV appeals failed submission → AI advises → PM/CEO resolves.
- Time negotiation: DEV proposes minutes → AI advises → PM/CEO approves/rejects.
- Human quality gate: PM/CEO can approve a task regardless of AI verdict.
- Google Slides import: parse a presentation URL into tasks (with preview + selective import).
- Timeline views: Epic Roadmap and Sprint Execution (Gantt), drag/resize bars, milestone overlay.
- PDF export: timeline and quotation PDF via headless Chrome (chromedp).
- Kanban board with drag-and-drop status changes.
- Product Backlog management: sort, group by epic, inline editing.
- Sprint management: create, start, complete, reopen; one active sprint per project.
- Project analytics: burndown chart, velocity chart, team capacity table, cycle time.
- Costing tab: fully-loaded cost calculation + PDF quotation export.
- Finance module: monthly accounting entries (revenue, expenses, cash balance), CEO summary.
- Performance module: personal KPIs (DEV), team leaderboard (PM/CEO), overview KPIs (CEO).
- AI configuration: active model, temperature, Cursor assistance %; usage stats (RPM/RPD).
- Role-based access control (CEO / PM / DEV) enforced at API level.
- Dark mode UI by default.
- Real-time dashboard updates via WebSockets.

### Out of Scope (Current)

- Third-party issue tracker integration (GitHub, Jira, Linear).
- Email / SMS notifications.
- GraphQL API.
- gRPC inter-service communication.
- Multi-tenancy (single tenant currently).
- Mobile app.

---

## User Stories / Use Cases

### Authentication
- As any user, I can register and log in with email/password.
- As any user, I can update my profile and change my password.
- As CEO, I can manage all user accounts and assign roles.

### Projects & Tasks
- As PM/CEO, I can create, edit, and delete projects.
- As PM/CEO, I can create tasks with title, description, priority, story points, epic, sprint, and dates.
- As PM, I can assign tasks to DEVs.
- As DEV, I can view only my assigned tasks.
- As any user, I can view a task's full detail page (comments, time logs, submission history).
- As any user, I can add comments to a task.
- As DEV, I can log time against a task.
- As PM/CEO, I can do bulk status updates on tasks (kanban drag).
- As any user, I can import tasks from a Google Slides presentation.

### AI Time Estimation
- As PM/CEO, I can trigger AI estimation on a single task.
- As PM/CEO, I can trigger a batch AI schedule for all tasks in a project.
- Every task must store `ai_estimated_minutes`; it is required by the system.

### AI Code Audit
- As DEV, I can submit a commit hash and diff for AI review on a task.
- The system must return: `verdict` (PASS/FAIL/PENDING), `score` (0–100), and `ai_feedback` (JSONB).
- Audits are stored in the `submissions` table.

### Appeal System
- As DEV, I can appeal a FAIL verdict with a written reason.
- The AI must advise UPHOLD or OVERTURN with confidence and reasoning.
- As PM/CEO, I can resolve the appeal (approve or reject).

### Time Negotiation
- As DEV, I can propose different estimated minutes for a task with a reason.
- The AI must advise APPROVE or REJECT with confidence and reasoning.
- As PM/CEO, I can approve or reject the negotiation.

### Sprints & Epics
- A project can have many sprints; only one may be ACTIVE at a time.
- Sprints follow the flow: PLANNING → ACTIVE → COMPLETED (with reopen allowed).
- Epics group tasks; epics have color, sort order, and date range.
- Task dependencies (FS, SS types) can be defined between tasks.

### Timeline
- PM/CEO can view a Gantt chart in Epic Roadmap or Sprint Execution mode.
- Gantt bars can be dragged to shift dates; edges can be resized.
- Milestones appear as vertical overlay lines; can be dragged to change due_date.
- Moving a task within an epic auto-recalculates epic start/end dates.
- Timeline can be exported to PDF.

### Analytics
- As PM/CEO, I can view per-project analytics: task completion, story points, cycle time, burndown, velocity, team capacity/utilisation.

### Costing
- As CEO/PM, I can select tasks and calculate a fully-loaded cost estimate.
- The system uses: dev salaries, overhead, exec salaries, PM salaries, risk buffer, profit margin, and VAT (7% THB).
- A PDF quotation can be exported.

### Finance
- As CEO, I can enter monthly revenue, expenses, and cash balance.
- As CEO, I can see runway and burn rate summary.

### Performance
- As DEV, I can view my personal KPIs.
- As PM, I can view a team leaderboard scoped to my assigned DEVs.
- As CEO, I can view system-wide performance overview.

### AI Configuration (CEO)
- As CEO, I can set the active AI model (default: `gemini-2.5-flash-lite`), temperature (default: 0.4), and Cursor assistance % (default: 80).
- As CEO, I can view AI usage stats (RPM/RPD).

---

## Functional Requirements

### Backend (Go + Gin, Hexagonal Architecture)

| # | Requirement |
|---|-------------|
| B-01 | All modules follow `handlers → usecases → repositories` layering. No business logic in handlers. |
| B-02 | Every task record must have `ai_estimated_minutes` (int) and `due_at` (time, nullable). |
| B-03 | Code audit results must be stored in `submissions.ai_feedback` as JSONB. |
| B-04 | No use of `panic`. Errors must always be returned and wrapped. |
| B-05 | No use of `interface{}` unless strictly justified with a comment. |
| B-06 | Auth uses JWT. Sessions/OTPs stored in Redis (TTL enforced). |
| B-07 | PostgreSQL for all core ACID data; Redis for sessions, rate limiting, and temporary tokens only. |
| B-08 | GORM AutoMigrate runs on startup for all registered entities. |
| B-09 | API is prefixed `/api/v1`; health check is at `GET /health` (no prefix). |
| B-10 | All routes must be protected by JWT middleware except `/auth/register`, `/auth/login`, `/health`. |
| B-11 | Role-based access enforced in middleware or usecase layer (not only in handlers). |
| B-12 | AI services (Gemini, Groq) are injected as interfaces; a `noop_ai_service` must exist for offline/test. |
| B-13 | AI usage tracked per request (RPM/RPD) in `ai_usage_tracker`. |
| B-14 | Timeline PDF export uses `chromedp` (headless Chrome). |
| B-15 | One ACTIVE sprint per project at a time; enforced at the usecase layer. |

### Frontend (Nuxt 3 + TypeScript, FSD)

| # | Requirement |
|---|-------------|
| F-01 | Feature modules live under `web/core/modules/{feature}/` with `infrastructure/`, `store/`, `ui/` sub-folders. |
| F-02 | Pages under `web/pages/` are dumb components — no business logic; all logic lives in modules. |
| F-03 | Dark mode is the default theme. |
| F-04 | All HTTP calls go through `web/core/shared/api/http.ts`. |
| F-05 | Pinia stores manage all client state. No direct API calls from pages. |
| F-06 | TypeScript strict mode enabled; no `any` types. |
| F-07 | TailwindCSS for all styling; no inline styles or external CSS frameworks. |
| F-08 | Dashboard uses WebSockets for real-time updates. |
| F-09 | Gantt chart uses `vue-ganttastic` (dark theme, label column 220px). |
| F-10 | Charts use Chart.js (burndown as Line, velocity as Bar). |
| F-11 | Drag-and-drop (sprint reorder, epic reorder, task reorder, kanban, Gantt) must persist via PATCH API calls after drop. |
| F-12 | Session storage persists backlog expand/collapse state per project. |
| F-13 | All currency displayed in THB using `Intl.NumberFormat 'th-TH'`. |

### Infrastructure / DevOps

| # | Requirement |
|---|-------------|
| I-01 | System runs via Docker Compose. |
| I-02 | Web app served on port **3000**. API on port **8080** (or env-configured). |
| I-03 | PostgreSQL 15, Redis 7. |
| I-04 | Redis hard memory limit: **128MB**, eviction policy: `allkeys-lru`. All Redis keys must have TTL. |
| I-05 | `docker-compose.prod.yml` for production builds. |

---

## Non-Functional Requirements

| Category | Requirement |
|----------|-------------|
| **Performance** | Dashboard should load within 2s on local Docker. API responses for list endpoints < 500ms. |
| **Scalability** | Horizontal scaling via multiple Go API instances behind a load balancer. PostgreSQL read replicas for reporting queries. |
| **Security** | HTTP-only cookies for JWT. Bcrypt password hashing. GORM prepared statements (no raw SQL injection). CORS configured. Rate limiting via Redis. |
| **Reliability** | API errors must return structured JSON `{ "error": "..." }` responses. Never return stack traces to clients. |
| **Maintainability** | Each module is self-contained; changing one module must not break others. |
| **Testability** | All usecase layer logic must be unit-testable via interface mocks. |
| **Observability** | Structured logging (Zap/Zerolog). `GET /health` checks Postgres and Redis connectivity. |

---

## Acceptance Criteria

| Feature | Acceptance Criteria |
|---------|---------------------|
| Task creation | Task is saved with `ai_estimated_minutes` (can be 0 on creation); `due_at` is optional but present in schema. |
| AI estimation | Calling estimate on a task populates `ai_estimated_minutes` and returns the updated task. |
| Code audit | Submitting a diff returns a `submissions` record with `ai_verdict`, `ai_score`, and `ai_feedback` populated. |
| Appeal | An appeal changes `submissions.is_overridden = true` when resolved as OVERTURN. |
| Sprint constraint | Starting a sprint when another is ACTIVE returns a 400/409 error. |
| Role enforcement | A DEV calling a PM-only endpoint returns 403. |
| Timeline drag | Dragging a Gantt bar calls PATCH on task dates; epic dates auto-sync. |
| Costing | PDF quotation includes labor subtotal, risk buffer, profit margin, VAT, and grand total in THB. |
| Health check | `GET /health` returns 200 with Postgres and Redis status. |
| Docker | `docker-compose up` starts all four services (postgres, redis, api, web); web is accessible at `http://localhost:3000`. |

---

## Dependencies

| Dependency | Version / Notes |
|------------|----------------|
| Go | 1.23 |
| Gin | HTTP framework |
| GORM | ORM (PostgreSQL driver) |
| PostgreSQL | 15 |
| Redis | 7 |
| Nuxt | 3 (Vue 3, TypeScript) |
| TailwindCSS | 3 |
| Pinia | State management |
| vue-ganttastic | Gantt chart |
| Chart.js | Burndown + velocity charts |
| chromedp | PDF export (headless Chrome) |
| Gemini API / Groq API | AI services (model set via `system_configs`) |
| Google Slides API | Slide import (Google API key required) |
| Docker / Docker Compose | v2.0+ |

---

## Risks and Edge Cases

| Risk | Mitigation |
|------|-----------|
| AI service downtime | `noop_ai_service` fallback; system must not block task creation if AI is unavailable. |
| Large diffs causing slow AI response | Consider chunking or async processing for submissions > N lines. |
| Concurrent sprint start race condition | Enforce single ACTIVE sprint at DB level (partial unique index or usecase-level transaction). |
| Timeline drag producing invalid date ranges | Validate `start_date <= end_date` before PATCH; reject with 400 if invalid. |
| Redis hitting 128MB limit | LRU eviction; always set TTL; monitor with `INFO memory`. |
| Gantt epic scale with no task dates | Skip tasks without `start_date`/`end_date` in scale calculation. |
| Google Slides URL permissions | Show clear error if the presentation is not publicly accessible. |
| chromedp cold start latency | PDF export endpoint may be slow on first request; acceptable for non-critical path. |
| DEV submitting on another DEV's task | Enforce `task.assigned_to === currentUser.id` check in submit usecase. |

---

## Open Questions

| # | Question | Status |
|---|----------|--------|
| OQ-01 | Should AI estimation be triggered automatically on task creation, or always manual? | Open |
| OQ-02 | What is the maximum diff size accepted for code audit? | Open |
| OQ-03 | Should sprint capacity planning (available hours per DEV) be tracked? | Open |
| OQ-04 | Is multi-currency support needed for the costing module? | Open (THB assumed) |
| OQ-05 | Are WebSocket updates scoped per-project or system-wide? | Open |
| OQ-06 | Should time logs be editable after submission, or append-only? | Open |

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-04-29 | Initial requirements created from GEMINI_CONTEXT.md, ARCHITECTURE.md, PROJECT_DETAIL_PAGE_CONTEXT.md, .cursorrules, and README | Agent |
