# The Sentinel – System Context for AI (Gemini / Cursor)

This document describes what the system does, its structure, key files, and data schema. Use it to give context to AI assistants (e.g. Gemini, Cursor).

---

## 1. What This System Does

**The Sentinel** is an **AI-driven Development OS**: it helps teams **estimate task time** and **audit code** using AI (Groq or Gemini).

### Core capabilities

| Area | Description |
|------|-------------|
| **Projects & tasks** | Multi-project; tasks have `ai_estimated_minutes`, `due_at`, story points, priority, status (PENDING → IN_PROGRESS → COMPLETED). |
| **AI time estimation** | AI estimates effort (minutes) from task title/description; optional batch estimate + schedule for a project. |
| **Code review (audit)** | Dev submits code (diff); AI returns verdict (PASS/FAIL), score (0–100), and feedback. Audits stored in `submissions` table as JSONB. |
| **Appeal system** | If AI fails a submission, dev can appeal; AI advises (UPHOLD/OVERTURN); PM/CEO resolves. |
| **Time negotiation** | Dev can propose different minutes; AI advises (APPROVE/REJECT); PM/CEO approves. |
| **Human quality gate** | PM/CEO can approve a task (mark COMPLETED) regardless of AI. |
| **Roles** | **CEO** (full access, user/team/AI config), **PM** (projects, tasks, assign, approve, PM-scoped leaderboard), **DEV** (my tasks, submit, log time, appeal, negotiate). |
| **Sprints & epics** | Sprints (time-boxed), Epics (feature groups), Milestones; Gantt view with task dependencies. |
| **Google Slides import** | Import slides as tasks (preview + import); optional sprint/epic. |
| **Timeline export** | Epic/Sprint timeline views; export to PDF (chromedp). |
| **Finance** | Monthly accounting entries (revenue, expenses, cash balance); CEO summary (runway, burn rate). |
| **Performance** | Personal KPIs (DEV); team leaderboard (PM/CEO); overview KPIs (CEO). |
| **AI config** | CEO can set active model, temperature, Cursor assistance %; view AI usage (RPM/RPD). |

### Tech stack

- **Backend:** Go 1.23, Gin, GORM, PostgreSQL (primary), MongoDB (logs/audit), Redis (cache/sessions).
- **Frontend:** Nuxt 3 (Vue 3, TypeScript), TailwindCSS, dark mode by default, Pinia.
- **Run:** Docker; web on **port 3000**, API port from env (e.g. 8080).

---

## 2. Project Structure (High Level)

```
sentinel-core/
├── api/                    # Go backend (Hexagonal, modular monolith)
│   ├── cmd/
│   │   ├── server/         # main.go – wire modules, routes, DB
│   │   └── migrate/       # DB migration entry
│   └── internal/
│       ├── core/           # config, database, middleware, pdf (chromedp)
│       └── modules/       # auth, finance, health, performance, sentinel
│
├── web/                    # Nuxt 3 frontend (FSD)
│   ├── core/modules/       # auth, finance, performance, projects, tasks
│   ├── pages/              # file-based routes
│   ├── components/         # shared UI (dashboard, tasks, projects, performance, editor)
│   ├── layouts/
│   ├── plugins/
│   └── composables/
│
├── docs/                   # ARCHITECTURE.md, setup guides
├── docker-compose.yml
├── Makefile
└── .cursorrules            # AI coding standards (Hexagonal, FSD, no panic)
```

---

## 3. Backend (API) – File Layout

### 3.1 Entry and core

| Path | Purpose |
|------|--------|
| `api/cmd/server/main.go` | Load config, init Postgres/Mongo/Redis, AutoMigrate entities, wire all modules, register routes, run server. |
| `api/cmd/migrate/main.go` | Migration runner. |
| `api/internal/core/config/config.go` | Env/config (DB, JWT, Groq/Gemini API keys, Google API key, etc.). |
| `api/internal/core/database/database.go` | Postgres, Mongo, Redis init. |
| `api/internal/core/middleware/auth.go` | JWT auth middleware. |
| `api/internal/core/pdf/chromepdf.go` | Timeline PDF export (chromedp). |

### 3.2 Modules (Hexagonal: domain → usecase → repository → delivery/http)

Each module follows: **handlers (delivery)** → **usecases** → **repositories**.

| Module | Domain | Usecase | Repository | Delivery (HTTP) |
|--------|--------|---------|------------|------------------|
| **auth** | `domain/entity.go` | `usecase/auth_usecase.go` | `repository/postgres_repo.go` | `delivery/http/auth_handler.go`, `route.go` |
| **sentinel** | `domain/entities.go` | `usecase/sentinel_usecase.go` | `repository/postgres_repository.go` + `gemini_service.go` / `groq_service.go` / `noop_ai_service.go`, `ai_usage_tracker.go` | `delivery/http/sentinel_handler.go`, `route.go` |
| **finance** | `domain/entity.go` | `usecase/finance_usecase.go` | `repository/postgres_repo.go` | `delivery/http/finance_handler.go`, `route.go` |
| **performance** | `domain/entities.go` | `usecase/performance_usecase.go` | `repository/postgres_repo.go` | `delivery/http/performance_handler.go`, `route.go` |
| **health** | – | – | – | `delivery/http/health_handler.go`, `route.go` (checks Postgres, Mongo, Redis) |

### 3.3 API route prefixes (all under `/api/v1` unless noted)

- **Health:** `GET /health` (no prefix).
- **Auth:** `POST /auth/register`, `POST /auth/login`; protected: `GET/PATCH /auth/me`, `PATCH /auth/me/password`, `GET/POST/DELETE/PATCH /auth/users`, `PATCH /auth/users/:id/role`, `PATCH /auth/users/:id/password`.
- **Sentinel:**  
  - Projects: `POST/GET/PATCH/DELETE /sentinel/projects`, `GET/PATCH/DELETE /sentinel/projects/:id`, `POST /sentinel/projects/:id/ai-schedule`, `POST /sentinel/projects/:id/ai-plan`, `POST /sentinel/projects/:id/clear-plan`, `GET /sentinel/projects/:id/analytics`.  
  - Tasks: `POST/GET /sentinel/tasks`, `GET /sentinel/tasks/my`, `GET /sentinel/tasks/unassigned`, `GET /sentinel/tasks/approvals`, `GET /sentinel/tasks/gantt`, `GET/PATCH/DELETE /sentinel/tasks/:id`, `POST /sentinel/tasks/:id/assign`, `POST /sentinel/tasks/:id/submit`, `POST /sentinel/tasks/:id/estimate`, `POST /sentinel/tasks/:id/negotiate`, `POST /sentinel/tasks/:id/approve`, `POST /sentinel/tasks/:id/comments`, `GET /sentinel/tasks/:id/comments`, `POST/GET /sentinel/tasks/:id/time-logs`, `PATCH /sentinel/tasks/bulk-status`.  
  - Dependencies: `POST /sentinel/tasks/dependencies`, `DELETE /sentinel/tasks/dependencies/:id`.  
  - Appeals: `POST /sentinel/submissions/:id/appeal`, `POST /sentinel/appeals/:id/resolve`.  
  - Sprints: `POST/GET/PATCH/DELETE /sentinel/sprints`, `POST /sentinel/sprints/:id/start|complete|reopen`, `POST /sentinel/sprints/:id/tasks`.  
  - Milestones: `POST/GET/PATCH/DELETE /sentinel/milestones`.  
  - Epics: `POST/GET/PATCH/DELETE /sentinel/epics`.  
  - Timeline: `GET /sentinel/projects/:id/timeline/epic-view`, `sprint-view`, `export-pdf`.  
  - Import: `POST /sentinel/import/google-slides/preview`, `POST /sentinel/import/google-slides`.
- **Admin:** `GET/PUT /admin/config`, `GET /admin/models`, `GET /admin/ai-usage`.
- **Finance:** (under same group, typically `/api/v1/finance` or similar) entries + summary.
- **Performance:** (under same group) `/performance/me`, `/performance/team`, `/performance/overview`.

---

## 4. Frontend (Web) – File Layout

### 4.1 Pages (Nuxt file-based routing)

| Path | Purpose |
|------|--------|
| `web/pages/index.vue` | Home/landing. |
| `web/pages/login.vue`, `web/pages/register.vue` | Auth. |
| `web/pages/dashboard.vue` | Role-based dashboard (CEO / PM / DEV). |
| `web/pages/projects/index.vue` | Project list. |
| `web/pages/projects/[id].vue` | Project detail (tasks, epics, sprints, milestones, analytics, timeline, import). |
| `web/pages/projects/gantt.vue` | Gantt view. |
| `web/pages/projects/sprint/[sprintId].vue` | Sprint view. |
| `web/pages/tasks/index.vue` | Task list. |
| `web/pages/task/[id].vue` | Task detail (submit, comments, time logs, appeal, negotiate). |
| `web/pages/create.vue` | Create project/task (if used). |
| `web/pages/accounting.vue` | Finance (CEO). |
| `web/pages/admin/team.vue` | User management (CEO). |
| `web/pages/admin/ai-settings.vue` | AI config & usage (CEO). |
| `web/pages/profile.vue` | Own profile. |
| `web/pages/performance.vue` | Performance KPIs / leaderboard. |

### 4.2 Core modules (FSD: infrastructure → store → ui)

| Module | Files | Purpose |
|--------|-------|--------|
| **auth** | `infrastructure/auth-api.ts`, `store/auth-store.ts`, `ui/LoginForm.vue`, `RegisterForm.vue` | Login, register, me, users (CEO). |
| **projects** | `infrastructure/projects-api.ts`, `store/projects-store.ts` | Projects CRUD, analytics, timeline. |
| **tasks** | `infrastructure/tasks-api.ts`, `store/tasks-store.ts` | Tasks CRUD, assign, submit, comments, time logs. |
| **finance** | `finance-api.ts` | Accounting entries, summary. |
| **performance** | `performance-api.ts`, `performance-store.ts` | KPIs, team leaderboard. |

### 4.3 Shared and components

- `web/core/shared/api/http.ts` – Base HTTP client (auth, base URL).
- `web/plugins/api.ts` – API plugin.
- `web/composables/useAuth.ts`, `useNotification.ts`.
- `web/components/dashboard/CeoView.vue`, `PmView.vue`, `DevView.vue`.
- `web/components/projects/KanbanBoard.vue`, `GanttMilestoneRow.vue`, `MilestoneTimeline.vue`, `ProjectAnalytics.vue`.
- `web/components/tasks/TaskComments.vue`, `TimeLogger.vue`.
- `web/components/performance/TeamLeaderboard.vue`, `KpiScoreCard.vue`, `SpaceRadar.vue`.
- `web/components/editor/RichTextEditor.vue`, `ImageAnnotator.vue`, etc.
- `web/utils/timelinePdfExport.ts` – Timeline PDF export helper.

---

## 5. Data Schema (PostgreSQL, GORM)

Entities are in `api/internal/modules/*/domain/`. Tables are created by GORM AutoMigrate in `main.go`.

### 5.1 Auth – `users`

| Column | Type | Note |
|--------|------|------|
| id | uint (PK) | |
| email | string, unique, not null | |
| password | string, not null | hashed, never JSON |
| role | varchar(20), default 'DEV' | CEO, PM, DEV |
| health_score | decimal(5,2), default 100 | |
| tech_stack | text[] | |
| display_name | varchar(100) | |
| created_at, updated_at | time | |

### 5.2 Sentinel – core tables

**projects**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK), default gen_random_uuid() | |
| code | varchar(64), unique | slug e.g. mims-hdmap-main |
| name, description | string, text | |
| status | default 'ACTIVE' | ACTIVE, COMPLETED, ON_HOLD |
| created_at, updated_at | time | |

**epics**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| project_id | uuid, not null | |
| title, description | string, text | |
| status | default 'PLANNING' | PLANNING, IN_PROGRESS, DONE |
| color | default '#6366f1' | |
| sort_order | int, default 0 | |
| start_date, end_date | date, nullable | |
| created_at, updated_at | time | |

**sprints**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| project_id | uuid, not null | |
| name, goal | string, text | |
| start_date, end_date | date, nullable | |
| status | default 'PLANNING' | PLANNING, ACTIVE, COMPLETED |
| sort_order | int | |
| created_at, updated_at | time | |

**milestones**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| project_id | uuid, not null | |
| title, description | string, text | |
| due_date | date, nullable | |
| status | default 'PENDING' | PENDING, REACHED, MISSED |
| created_at, updated_at | time | |

**tasks**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| code | varchar(64), unique | e.g. mims-hdmap-main-001 |
| title, description | string, text | |
| resource_urls | jsonb, default '{}' | e.g. slide images/URLs |
| ai_estimated_minutes | int | **Required by Sentinel** |
| project_id | uuid, nullable | |
| epic_id, sprint_id, milestone_id | uuid, nullable | |
| parent_id | uuid, nullable | WBS sub-tasks |
| sort_order | int | |
| start_date, end_date | date, nullable | |
| progress | int, default 0 | 0–100 |
| priority | default 'MEDIUM' | CRITICAL, HIGH, MEDIUM, LOW |
| story_points | int, default 0 | |
| negotiation_status | default 'NONE' | NONE, PENDING, APPROVED, REJECTED |
| proposed_minutes, negotiation_reason | int, text | |
| negotiation_ai_recommendation, negotiation_ai_confidence, negotiation_ai_reasoning | string, int, text | |
| due_at, started_at, completed_at | time, nullable | |
| status | default 'PENDING', index | |
| assigned_to | uint, nullable | dev |
| assigned_by_id | uint, nullable | PM/CEO (for PM-scoped leaderboard) |
| created_by | uint, nullable | |
| created_at, updated_at | time | |

**task_dependencies**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| predecessor_id, successor_id | uuid, not null | |
| type | varchar(2), default 'FS' | FS, SS, etc. |
| created_at | time | |

**task_comments**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| task_id | uuid, not null | |
| user_id | uint, not null | |
| content | text, not null | |
| created_at | time | |

**time_logs**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| task_id | user_id | uint, not null |
| minutes | int, not null | |
| description | text | |
| logged_at | time | |

**submissions** (code audit)

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| task_id | uuid, not null | |
| dev_id | uint, not null | |
| commit_hash | string, not null | |
| diff | text | |
| ai_verdict | string | PASS, FAIL, PENDING |
| ai_score | int | |
| ai_feedback | jsonb, default '{}' | **Audits stored here** |
| is_overridden | bool, default false | appeal approved |
| created_at | time | |

**appeals**

| Column | Type | Note |
|--------|------|------|
| id | uuid (PK) | |
| submission_id | uuid, not null | |
| developer_id | uint, not null | |
| reason | text, not null | |
| status | default 'PENDING' | PENDING, APPROVED, REJECTED |
| ai_recommendation, ai_confidence, ai_reasoning | text, int, text | |
| resolver_id | uint, nullable | |
| resolver_note | string | |
| created_at, updated_at | time | |

**system_configs** (singleton)

| Column | Type | Note |
|--------|------|------|
| id | uint (PK) | |
| active_model | default 'gemini-2.5-flash-lite' | |
| temperature | float32, default 0.4 | |
| cursor_assistance | int, default 80 | 0–100 |
| updated_at | time | |

### 5.3 Finance – `finance_monthly_entries`

| Column | Type | Note |
|--------|------|------|
| id | uint (PK) | |
| year, month | int, unique (year, month) | 1–12 |
| revenue, expenses, cash_balance | decimal(15,2) | |
| note | text | |
| created_at, updated_at | time | |
| deleted_at | soft delete | |

### 5.4 Performance

No dedicated table; aggregates from `users`, `tasks`, `submissions`, `sprints`, `milestones`, `system_configs`.

---

## 6. Conventions (for AI-generated code)

- **Backend:** Hexagonal: handlers → usecase → repository. No `panic`; avoid `interface{}`; wrap errors. Tasks must have `ai_estimated_minutes` and `due_date` where applicable. Audits in `submissions.ai_feedback` (JSONB).
- **Frontend:** FSD under `web/core/modules/{feature}` (infrastructure, store, ui). Tailwind; dark mode default. Real-time where needed (e.g. dashboard).
- **Vertical slice:** New feature = DB/entity → repo → usecase → handler → route (and frontend module + page if needed).
- **Port:** Web app is at **http://localhost:3000** when using Docker.

---

## 7. Quick reference – where to put what

| I want to… | Backend | Frontend |
|------------|--------|----------|
| Add an API endpoint | Add handler in module’s `delivery/http`, register in module’s `route.go` | Call from `core/modules/{feature}/infrastructure/*-api.ts` |
| Add a new entity/table | Add model in module’s `domain/`, add to `db.AutoMigrate()` in `api/cmd/server/main.go` | Types can mirror or use API response types |
| Add business logic | Module’s `usecase/*.go` | Store or composable |
| Add a new page | – | `web/pages/**/*.vue` (file-based route) |
| Add a new feature area | New folder under `api/internal/modules/` (domain, usecase, repository, delivery/http) | New folder under `web/core/modules/` (infrastructure, store, ui) |

---

*Last updated from codebase: Sentinel core (api + web). Excludes mims-api-service and mims-web.*
