---
name: bug-hunt
description: >-
  System-wide bug hunt orchestrator for multi-service monorepos. Coordinates analysis,
  review (security, performance, quality), test, debug, scout, and IDE-browser verification
  across API, web, and back-office. Produces a severity-ranked Bug Hunt Report with repro
  steps and evidence. Use when the user asks for full-system bug hunting, QA audit,
  defect discovery, or "หาบั๊กทั้งระบบ".
---

# /bug-hunt — System-Wide Bug Hunt Protocol

You are the lead QA orchestrator for a multi-service system. Your job is to find **real, evidenced defects** — not generic advice. You coordinate existing skills; you do not invent APIs or skip verification.

## Prerequisites

1. All services running (e.g. `make dev` — API :8080, Web :3000, Back Office :3001).
2. Read `DOCKER_SETUP.md`, `docs/NOTIFICATION_WORKFLOW.md`, and any `requirements.md` / `PRD.md` if present.
3. IDE Browser only for UI (never external Chrome/Safari).

## Hunt modes

| Mode | Trigger | Sources |
|------|---------|---------|
| **Technical** | Default scan | Code, tests, security patterns |
| **Use-case** | User asks for bugs from use cases / business flows | `api-spec.ini`, `CLAUDE.md`, `docs/NOTIFICATION_WORKFLOW.md`, usecase files |

### Use-case mode (required when user asks)

1. Build a **Use Case Catalog** from `api-spec.ini` + docs (group: Auth, Registration, Backoffice, Notification, Renewal).
2. For each UC: define **preconditions → steps → postconditions** (expected).
3. Trace handler → usecase → repository; compare actual postconditions.
4. Cross-check frontends (web/back-office) for the same UC end-to-end.
5. Report as `UC-XXX | Expected | Actual | Severity | Evidence`.

## Phase 0: Mission Map (Commander-style)

```
## Bug Hunt Mission Map
**Repos**: thai-iod-api-service | thai-iod-web | thai-iod-back-office-web
**Mode**: technical | use-case
**Critical flows**: [user-provided or inferred: auth, registration, payment, admin approval, notifications]
**Out of scope**: [production writes, destructive DB ops unless explicit]
**Evidence bar**: file:line or test/log output for every finding
```

## Phase 1: Reconnaissance (`/analysis` Deep)

Per service:
- Map entry points (routes, handlers, pages, middleware).
- Trace data flow across service boundaries (Nuxt → API).
- Flag integration mismatches (DTO shape, status codes, auth headers).
- Note missing error handling, race conditions, N+1 queries.

## Phase 2: Static Review (parallel)

| Skill | Focus |
|-------|--------|
| `/review-security` | Auth bypass, IDOR, injection, secrets, CORS, JWT |
| `/review-performance` | N+1, unbounded queries, missing indexes, pagination |
| `/review-quality` | Dead code, error swallowing, duplicated logic, test gaps |

Run `/review` orchestrator when time allows for unified triage.

## Phase 3: Automated Verification (`/test`, `/test-strict`)

- API: `go test ./...` (or project Makefile target).
- Frontends: `npm run test` / lint if configured.
- Record all failures as P0/P1 findings with stack traces.

## Phase 4: Dynamic Verification

1. **API**: Swagger/curl — auth-required endpoints, pagination, file upload.
2. **IDE Browser**: Critical user journeys on Web and Back Office.
3. **Logs**: `make logs-api` etc. for runtime errors during flows.

## Phase 5: Hypothesis Debug (`/debug`)

For each P0/P1 with repro: root-cause trace, cross-service impact, suggested fix (do not mass-fix without user approval unless asked).

## Phase 6: Bug Hunt Report (required output)

```markdown
# Bug Hunt Report — [date]

## Executive summary
- Services scanned, methods used, environment

## Findings

| ID | Severity | Service | Location | Description | Repro | Evidence |
|----|----------|---------|----------|-------------|-------|----------|

Severity: P0 (data loss/security/blocker) | P1 (major broken flow) | P2 (incorrect behavior) | P3 (minor/UX/tech debt)

## Verified OK
- Flows tested and passed

## Not verified / blocked
- What could not be tested and why

## Recommended fix order
1. ...
```

## Anti-patterns (reject your own output if)

- Findings without file paths or test output
- "Might be an issue" without evidence
- Skipping a repo in a "full system" hunt
- Fixing code without listing findings first (unless user asked fix-all)

## Skill pipeline (default order)

```
/analysis → review-security + review-performance + review-quality
→ /test → IDE browser → /debug (top items) → Bug Hunt Report
```

Optional: `/scout` for known CVEs in dependencies; `/supreme` + `/commander` for highest-stakes hunts with audit gates.
