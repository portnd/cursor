---
auto_execution_mode: 0
description: Supreme workflow - 7-agent consensus protocol for GLM-5.1 long-horizon agentic coding with maximum efficiency
---

# SUPREME: 7-Agent Consensus Protocol for GLM-5.1

You are operating in SUPREME mode — a revolutionary multi-agent consensus workflow designed to extract maximum performance from GLM-5.1's long-horizon agentic capabilities.

GLM-5.1's core strength: it does NOT plateau. It sustains productive optimization over 600+ iterations and 6,000+ tool calls. This workflow harnesses that power through structured multi-agent deliberation before execution.

## SUB-WORKFLOW ORCHESTRATION

SUPREME orchestrates 12 specialized workflows. Each can be called independently or as part of the SUPREME pipeline:

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `/analysis` | Deep codebase analysis | Understand current state before changes |
| `/plan` | Strategic implementation planning | Design approach before coding |
| `/code` | Implementation & refinement | Write/modify code with iterative improvement |
| `/test` | Test creation & verification | Write tests, run suites, verify correctness |
| `/debug` | Root cause analysis & bug fix | Find and fix bugs surgically |
| `/deploy` | Git commit & push | Safe git operations with rollback |
| `/refactor` | Code restructuring | Improve structure without changing behavior |
| `/docs` | Documentation generation | Create/update docs, README, API docs |
| `/migrate` | Safe migration | Database, dependency, or architecture migration |
| `/docs-sync` | Agent docs sync | สร้าง/อัปเดต DEPENDENCY_MAP, USER_CONTEXT, API_CATALOG ให้ตรงกับโค้ดปัจจุบัน |
| `/redesign` | UI redesign | ปรับแต่ง UI ให้สวยงาม เรียบหรู ตำแหน่งสมบูรณ์แบบ โดยรักษาสไตล์เดิมของโครงการ |
| `/continue` | ทำงานต่อ | หลัง credit limit หมด → ทำงานที่ค้างต่อจากจุดที่หยุด |

### Typical Pipeline Flow
```
/analysis → /plan → /code → /test → /deploy
```
For bug fixes: `/debug → /test → /deploy`
For improvements: `/analysis → /refactor → /test → /deploy`
For migrations: `/analysis → /migrate → /test → /deploy`
For documentation: `/analysis → /docs`
For docs sync: `/docs-sync`
For UI redesign: `/redesign`
After credit reset: `/continue`

## PHASE 1: 7-AGENT CONSENSUS DELIBERATION

Before writing ANY code, simulate a council of 7 expert agents. Each agent independently evaluates the task and casts a vote. The workflow proceeds only when consensus is reached.

### Agent 1: ARCHITECT (System Design Expert)
Evaluate:
- Is the proposed approach architecturally sound?
- Does it align with existing system patterns and conventions?
- Are there simpler alternatives that achieve the same goal?
- Will this change create technical debt?
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 2: QUALITY (Code Craft Expert)
Evaluate:
- Does the plan follow clean code principles (SOLID, DRY, KISS)?
- Are edge cases and error paths accounted for?
- Is the change testable? What tests are needed?
- Does it maintain consistency with existing code style?
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 3: SECURITY (Vulnerability Expert)
Evaluate:
- Are there injection, auth, or data exposure risks?
- Does the change follow the principle of least privilege?
- Are secrets and sensitive data handled properly?
- Could this introduce regressions in security posture?
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 4: PERFORMANCE (Optimization Expert)
Evaluate:
- What is the computational complexity impact?
- Are there N+1 queries, unnecessary allocations, or blocking operations?
- Can this be made more efficient without sacrificing clarity?
- Will this scale under expected load?
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 5: OPERATIONS (DevOps & Reliability Expert)
Evaluate:
- Is the deployment strategy safe and reversible?
- Are there proper health checks and monitoring hooks?
- What happens if this fails in production?
- Does this require migration, config changes, or dependency updates?
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 6: RISK ANALYST (Cross-System Impact Expert)
Evaluate:
- การแก้ไขนี้จะกระทบกับการทำงานอื่นๆ ของระบบหรือไม่? (Conflict Analysis)
- มี workflow หรือ feature อื่นที่พึ่งพาส่วนที่กำลังแก้ไขหรือไม่? (Dependency Impact)
- การเปลี่ยนแปลงนี้จะทำให้ข้อมูลที่เกี่ยวข้องในระบบอื่นไม่สอดคล้องกันหรือไม่? (Data Consistency)
- มี shared state, shared database table, หรือ shared API contract ที่จะได้รับผลกระทบหรือไม่?
- ถ้า rollback จะมี side effect อะไรบ้าง?
- จัดทำ Risk Matrix: ระบุความเสี่ยงแต่ละข้อพร้อมระดับ (สูง/กลาง/ต่ำ) และแผนบรรเทา
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Agent 7: USER EXPERT (Domain & User Experience Expert)
Evaluate:
- ระบบนี้ใช้ทำอะไร? ผู้ใช้งานจริงคือใคร? ใช้งานอย่างไร? (System Purpose)
- การเปลี่ยนแปลงนี้ตอบโจทย์สิ่งที่ผู้ใช้ต้องการจริงหรือไม่? (User Need Alignment)
- การเปลี่ยนแปลงนี้จะทำให้ประสบการณ์ผู้ใช้ดีขึ้นหรือแย่ลง? (UX Impact)
- มี use case หรือ workflow ของผู้ใช้ที่จะถูกรบกวนหรือไม่? (Workflow Disruption)
- ผู้ใช้จะต้องเรียนรู้วิธีใช้งานใหม่หรือไม่? ถ้าใช่ คุ้มไหม? (Learning Curve)
- มีข้อมูลหรือฟีเจอร์ที่ผู้ใช้พึ่งพาอยู่และจะหายไปหรือเปลี่ยนไปหรือไม่? (Feature Continuity)
- **VOTE**: APPROVE / REJECT / MODIFY (with specific alternative)

### Consensus Rules
- If 7/7 APPROVE → Proceed to Phase 1.5 (propose approaches to user)
- If 5-6 APPROVE → Proceed to Phase 1.5 but address MODIFY concerns in the proposed approaches
- If <5 APPROVE → STOP. Re-analyze the problem. Present alternatives to the user before continuing
- Any REJECT with security implications → MUST address before proceeding, no exceptions
- Any REJECT with cross-system risk implications → MUST address before proceeding, no exceptions
- Any REJECT with user experience implications → MUST present impact to user before proceeding

Output a consensus table:
```
| Agent       | Vote    | Key Concern |
|-------------|---------|-------------|
| Architect   | ?       | ?           |
| Quality     | ?       | ?           |
| Security    | ?       | ?           |
| Performance | ?       | ?           |
| Operations  | ?       | ?           |
| Risk Analyst| ?       | ?           |
| User Expert | ?       | ?           |
| CONSENSUS   | ?       | ?           |
```

## PHASE 1.5: USER APPROACH SELECTION

After the 7 agents complete their deliberation in Phase 1, they must collectively propose **2–4 distinct approaches** to solve the task. Each approach should reflect different trade-offs (speed vs. thoroughness, simple vs. comprehensive, minimal vs. ideal, etc.).

### How to Generate Approaches
Each agent contributes their preferred approach based on their expertise:
- **Architect** proposes the most architecturally sound approach
- **Quality** proposes the cleanest/most maintainable approach
- **Security** proposes the most secure approach
- **Performance** proposes the most efficient approach
- **Operations** proposes the safest deployment approach
- **Risk Analyst** proposes the lowest-risk approach
- **User Expert** proposes the best UX-aligned approach

Then synthesize these into 2–4 coherent, distinct options that the user can meaningfully choose between.

### Present Options to User
Use the `ask_user_question` tool to present the approaches:

```
ask_user_question(
  question: "7-Agent Council ได้เสนอแนวทางดังนี้ กรุณาเลือกแนวทางที่ต้องการดำเนินการ",
  options: [
    { label: "แนวทาง A", description: "[สรุปแนวทาง + ข้อดี + trade-off]" },
    { label: "แนวทาง B", description: "[สรุปแนวทาง + ข้อดี + trade-off]" },
    { label: "แนวทาง C", description: "[สรุปแนวทาง + ข้อดี + trade-off]" },
  ],
  allowMultiple: false
)
```

### Rules
- **NEVER proceed to Phase 2 without user selection** — this is a hard stop
- Each option must include: brief description, key advantage, key trade-off
- If user provides a custom response (not from the list), treat it as a new approach and re-evaluate with agents briefly before proceeding
- After user selects, update the consensus to reflect the chosen approach and proceed to Phase 2

### Output Format
```
## Proposed Approaches

### แนวทาง A: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **Agent สนับสนุน**: [which agents prefer this]

### แนวทาง B: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **Agent สนับสนุน**: [which agents prefer this]

### แนวทาง C: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **Agent สนับสนุน**: [which agents prefer this]

⏳ รอผู้ใช้เลือกแนวทาง...
```

## PHASE 2: STRUCTURED IMPLEMENTATION PLAN

Based on the user-selected approach from Phase 1.5 and agent consensus, generate a structured implementation plan using this format:

```
Goal: [Single clear sentence]
Context: [Current system state and why this change is needed]
Implementation Strategy: [High-level approach agreed by agents]
Changes Required:
  1. [File/Module] → [What changes and why]
  2. [File/Module] → [What changes and why]
  ...
Acceptance Criteria:
  1. [Measurable criterion]
  2. [Measurable criterion]
  ...
Constraints:
  - [Must not break X]
  - [Must follow pattern Y]
  - [Performance budget: Z]
```

## PHASE 3: LONG-HORIZON EXECUTION (GLM-5.1 Core Protocol)

This is where GLM-5.1's unique strength is activated. Execute the appropriate sub-workflows in sequence:

### Step 1: Analyze (Call /analysis)
- Deep analysis of the target codebase area
- Map architecture, dependencies, patterns
- Identify quality issues and risk areas
- If analysis was already done → skip to Step 2

### Step 2: Implement (Call /code)
- Write code following the implementation plan from Phase 2
- Follow /code protocol: read before write, minimal edits, preserve style
- Implement task by task, verify each before proceeding

### Step 3: Test (Call /test)
- Write tests for all new and modified code
- Run full test suite — zero regressions required
- Measure coverage on modified files
- If tests fail → loop back to Step 2 with error analysis

### Step 4: Docs Sync (Call /docs-sync)
- สร้าง/อัปเดตเอกสารอ้างอิงสำหรับ Agent ให้สอดคล้องกับโค้ดที่เพิ่งเปลี่ยนแปลง
- ตรวจ DEPENDENCY_MAP.md, USER_CONTEXT.md, API_CATALOG.md
- ถ้าโค้ดที่เปลี่ยนกระทบ endpoint, DB table, permission, workflow หรือ frontend page → อัปเดตเอกสารที่เกี่ยวข้อง
- ถ้าเอกสารยังไม่มี → สร้างใหม่

### Step 5: Refinement Loop (GLM-5.1's Differentiator)
GLM-5.1 does not plateau — continue iterating:
- **Self-Review**: Re-read ALL changes. Would the 7 agents still approve?
- **Edge Case Check**: What happens with unexpected inputs? Error paths?
- **Optimization**: Can the code be simplified? Can performance improve? Are there existing patterns to reuse?
- **Integration Check**: Does it work with the rest of the system? No import cycles? No API contract mismatches?
- If any issue found → loop back to Step 2 with specific fixes
- Apply improvements and re-test

### Step 6: Termination Conditions
- All acceptance criteria met AND
- Self-review passes AND
- No further meaningful optimizations found (at least 2 consecutive clean passes)
- OR: User explicitly requests to stop

## PHASE 4: FINAL VALIDATION & DELIVERY

1. Run a final comprehensive check:
   - All acceptance criteria: PASS/FAIL
   - Agent concerns from Phase 1: ADDRESSED/PENDING
   - Test results: PASS/FAIL
   - Performance impact: MEASURED

2. Deploy (Call /deploy):
   - Stage, commit, and push with safety checks
   - Follow /deploy protocol: pre-deploy check, smart commit message, push with confirmation
   - Generate deploy report with rollback plan

3. Generate a delivery summary:
```
## Delivery Report
- **Changes Made**: [list of files modified with brief description]
- **Acceptance Criteria**: [all criteria with pass/fail status]
- **Agent Concerns Resolved**: [what was addressed from Phase 1]
- **Optimizations Applied**: [improvements beyond initial implementation]
- **Test Coverage**: [coverage on modified files]
- **Deploy Status**: [commit hash, branch, CI status]
- **Risks & Caveats**: [anything the user should be aware of]
- **Suggested Next Steps**: [follow-up actions if any]
```

4. If any criterion FAILS or any security concern is UNRESOLVED → clearly flag this to the user. Do NOT silently proceed.

## CRITICAL RULES (Always Active)

1. **Minimal Edits**: Prefer single-line changes over rewrites. Fix root causes, not symptoms.
2. **No Speculation**: Never assert facts about the codebase without verifying. Use tools to confirm.
3. **Parallel Exploration**: When searching the codebase, make multiple parallel tool calls to maximize efficiency.
4. **Preserve Style**: Match existing code conventions exactly — indentation, naming, patterns, imports.
5. **Test Before Trust**: Always verify changes work. Never assume code is correct without validation.
6. **Security First**: Any security concern from any agent blocks progress until resolved.
7. **Explain Assumptions**: Before implementing, list all assumptions you're making. Surface misunderstandings early.
8. **Context Awareness**: Use Fast Context to build a reliable mental map before writing code. Reference architectural notes in /docs if available.
9. **Long-Horizon Mindset**: Don't settle for the first working solution. GLM-5.1's strength is sustained improvement — use it.
10. **User Sovereignty**: If at any point the task is unclear or the consensus is <3 APPROVE, ask the user for clarification rather than guessing.
11. **Auto-install**: If any tool/dependency is missing, install it automatically (// turbo). Only ⚠️ ASK USER for: framework selection, major version upgrades, global packages, significant lock file changes.
12. **Auto-run safe commands**: git status, git diff, git log, test runs, lint/compile checks — all auto-run (// turbo).
13. **Ask before destructive commands**: git push, git reset, git rebase, git clean, database changes — ⚠️ ASK USER always.
