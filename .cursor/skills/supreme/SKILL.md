---
name: supreme
description: >-
  Supreme workflow — 8-agent expert consultation protocol (Scout + 7 domain experts) fused with
  the Commander discipline backbone for maximum quality end-to-end implementation. Commander
  acts as the active orchestrator across every phase: dispatching sub-agents, auditing outputs,
  rejecting laziness, enforcing project rules, and gating final delivery. Use for important tasks
  requiring highest quality, complex features, or when the user mentions supreme, full pipeline,
  end-to-end, maximum quality, commander+supreme, or zero-defect execution.
---

# SUPREME: 8-Agent Consultation × Commander Discipline

You are operating in SUPREME mode — a multi-agent expert consultation workflow designed to extract maximum quality through structured expert feedback before execution, with the **Commander** protocol acting as the **enforcement spine** through every phase.

## COMMANDER AS THE ACTIVE SPINE (READ FIRST)

`/commander` is not a later-stage handoff — it is the **active orchestrator** across ALL phases of `/supreme`. At every phase you:

1. **Load Commander's rule manifest** (R1–R12 from `.cursor/skills/commander/SKILL.md`) and keep it active.
2. **Dispatch via Commander's Task Packet format** whenever a sub-agent is invoked — no free-form delegation.
3. **Audit every agent output** using Commander's Laziness Signal Matrix before proceeding.
4. **Reject & re-dispatch** any deliverable that fails audit, up to 3 attempts, then escalate.
5. **Gate each phase transition** with a Commander checkpoint — no phase advances while evidence is missing.
6. **Produce Commander's After-Action Report** as the final artifact (in addition to the Supreme Delivery Report).

If Commander's rule manifest or Task Packet protocol is unclear, **read `.cursor/skills/commander/SKILL.md` in full before continuing**.

## SUB-WORKFLOW ORCHESTRATION

SUPREME orchestrates these specialized workflows (all dispatched through Commander):

|| Command | Purpose | When to Use |
||---------|---------|-------------|
|| `/scout` | Technology research & best practices | Research current best practices before making tech decisions |
|| `/analysis` | Deep codebase analysis | Understand current state before changes |
|| `/plan` | Strategic implementation planning | Design approach before coding |
|| `/code` | Implementation & refinement | Write/modify code with iterative improvement |
|| `/test` | Test creation & verification | Write tests, run suites, verify correctness |
|| `/debug` | Root cause analysis & bug fix | Find and fix bugs surgically |
|| `/deploy` | Git commit & push | Safe git operations with rollback |
|| `/refactor` | Code restructuring | Improve structure without changing behavior |
|| `/docs` | Documentation generation | Create/update docs, README, API docs |
|| `/migrate` | Safe migration | Database, dependency, or architecture migration |
|| `/docs-sync` | Agent docs sync | Update reference docs to match codebase |
|| `/redesign` | UI redesign | Pixel-perfect UI polish within existing style |
|| `/continue` | Resume work | Continue from where work stopped |
|| `/review` | Code review | Review changes for bugs and improvements |
|| `/commander` | **Dispatch + discipline + audit (spine of every phase)** | **Always active throughout supreme** |

### Typical Pipeline Flow (Commander wraps every arrow)
```
Commander ─▶ /scout ─▶ Commander audit ─▶ /analysis ─▶ Commander audit ─▶ /plan ─▶
Commander audit ─▶ /code ─▶ Commander audit ─▶ /test ─▶ Commander audit ─▶ /deploy
                                                              │
                                           every ─▶ is gated by Commander's Task Packet
                                           verification + Laziness Signal Matrix audit
```
For bug fixes: Commander wraps `/debug → /test → /deploy`
For improvements: Commander wraps `/scout → /analysis → /refactor → /test → /deploy`
For migrations: Commander wraps `/scout → /analysis → /migrate → /test → /deploy`
For documentation: Commander wraps `/analysis → /docs`
For docs sync: Commander wraps `/docs-sync`
For UI redesign: Commander wraps `/redesign`
After interruption: Commander wraps `/continue`

## PHASE 0: SCOUT RECONNAISSANCE (Technology Intelligence)

**Commander's role in Phase 0**:
- Dispatch Scout via a Task Packet (TP-SCOUT) with explicit scope + deliverables (Scout Intelligence Brief).
- After Scout reports, audit the brief: was WebSearch actually used? Are findings dated 2026? Are versions/sources concrete? Any vague "probably" language → REJECT and re-dispatch (up to 3 attempts).
- Load Commander's rule manifest R1–R12 into working memory before the brief is consumed by Phase 1 agents.

Before the expert consultation, deploy the **Scout agent** to gather real-time technology intelligence. This ensures all subsequent expert agents base their recommendations on **current best practices**, not outdated assumptions.

### Agent 0: SCOUT (Technology Research & Trend Analyst)
**Research and provide:**
- Current best practices for the technologies involved (web search required)
- Latest stable versions, breaking changes, and deprecations
- Production case studies from top engineering teams
- Community consensus on recommended approaches
- Security advisories and known issues
- Alternative tools/libraries worth considering

**Output**: Scout Intelligence Brief — a concise research report that all other agents reference during consultation. Use WebSearch to ensure all findings are current (include year 2026 in queries).

**When to activate Scout:**
- ALWAYS before Phase 1 for complex features or technology decisions
- When any agent needs library/version/pattern recommendations
- When the user asks "is there a better/newer way to do this?"

---

## PHASE 1: 8-AGENT EXPERT CONSULTATION

**Commander's role in Phase 1**:
- Dispatch each of the 8 experts via a numbered Task Packet (TP-EXPERT-01 through TP-EXPERT-08) with a single, specific question set per packet.
- After each expert reports, audit for laziness signals: "looks fine", "probably works", generic platitudes without codebase evidence → REJECT and re-dispatch.
- Cross-check every agent's recommendation against the Scout Intelligence Brief and Commander's rule manifest R1–R12. Any recommendation that violates a rule must be flagged and revised before Phase 2.
- Produce an evidence ledger: for each expert, record the codebase files/snippets they referenced. No file references → evidence of a lazy pass → REJECT.

Before writing ANY code, simulate a council of 8 expert agents (including Scout from Phase 0). Each agent provides **expert feedback** from their domain perspective, then proposes **their recommended approach**. All agents should reference the **Scout Intelligence Brief** when evaluating technologies.

### Agent 0: SCOUT (Technology Research & Trend Analyst)
**Provide feedback on:**
- What are the current best practices for the technologies involved?
- Are there newer libraries, patterns, or approaches that are better?
- What breaking changes or deprecations should we be aware of?
- What are production-proven solutions used by top engineering teams?
- Are there security advisories or known issues with current choices?

**Propose approach**: Based on real-time research, propose the most modern and well-supported approach.

### Agent 1: ARCHITECT (System Design Expert)
**Provide feedback on:**
- Is the proposed approach architecturally sound?
- Does it align with existing system patterns and conventions?
- Are there simpler alternatives that achieve the same goal?
- Will this change create technical debt?

**Propose approach**: Based on architectural analysis, propose the most architecturally sound approach.

### Agent 2: QUALITY (Code Craft Expert)
**Provide feedback on:**
- Does the plan follow clean code principles (SOLID, DRY, KISS)?
- Are edge cases and error paths accounted for?
- Is the change testable? What tests are needed?
- Does it maintain consistency with existing code style?

**Propose approach**: Based on quality analysis, propose the cleanest and most maintainable approach.

### Agent 3: SECURITY (Vulnerability Expert)
**Provide feedback on:**
- Are there injection, auth, or data exposure risks?
- Does the change follow the principle of least privilege?
- Are secrets and sensitive data handled properly?
- Could this introduce regressions in security posture?

**Propose approach**: Based on security analysis, propose the most secure approach.

### Agent 4: PERFORMANCE (Optimization Expert)
**Provide feedback on:**
- What is the computational complexity impact?
- Are there N+1 queries, unnecessary allocations, or blocking operations?
- Can this be made more efficient without sacrificing clarity?
- Will this scale under expected load?

**Propose approach**: Based on performance analysis, propose the most efficient approach.

### Agent 5: OPERATIONS (DevOps & Reliability Expert)
**Provide feedback on:**
- Is the deployment strategy safe and reversible?
- Are there proper health checks and monitoring hooks?
- What happens if this fails in production?
- Does this require migration, config changes, or dependency updates?

**Propose approach**: Based on operational analysis, propose the safest deployment approach.

### Agent 6: RISK ANALYST (Cross-System Impact Expert)
**Provide feedback on:**
- Will this change impact other system operations? (Conflict Analysis)
- Are there workflows or features that depend on the part being changed? (Dependency Impact)
- Will this change make related data in other systems inconsistent? (Data Consistency)
- Are there shared states, shared DB tables, or shared API contracts affected?
- What are the rollback side effects?
- Build Risk Matrix with severity levels and mitigation plans

**Propose approach**: Based on risk analysis, propose the lowest-risk approach.

### Agent 7: USER EXPERT (Domain & User Experience Expert)
**Provide feedback on:**
- What is this system for? Who are the real users? How do they use it?
- Does this change address what users actually need?
- Will this change improve or worsen user experience?
- Are there user workflows that will be disrupted?
- Will users need to learn new patterns? Is it worth it?
- Are there features or data users depend on that will change or disappear?

**Propose approach**: Based on user experience analysis, propose the best UX-aligned approach.

### Expert Consultation Output
```
## 8-Agent Expert Consultation

|| Agent | Feedback Summary | Proposed Approach |
||-------|-----------------|-------------------|
|| Scout | [current tech intelligence + best practices] | [modern-approach-backed recommendation] |
|| Architect | [key architectural insights] | [their recommended approach] |
|| Quality | [code quality insights] | [their recommended approach] |
|| Security | [security insights] | [their recommended approach] |
|| Performance | [performance insights] | [their recommended approach] |
|| Operations | [operational insights] | [their recommended approach] |
|| Risk Analyst | [risk insights + risk matrix] | [their recommended approach] |
|| User Expert | [user experience insights] | [their recommended approach] |
```

---

## PHASE 2: CURATE 5 APPROACHES FOR USER SELECTION

**Commander's role in Phase 2**:
- Verify the 5 approaches are genuinely distinct (no slight variants). Similar approaches → merge and add a new distinct one.
- Confirm each approach is traceable to at least 2 expert proposals from Phase 1 (cite the agents in the "ข้อเสนอแนะจาก Agent" line).
- Confirm each approach explicitly lists which of R1–R12 could be violated if chosen poorly, plus mitigation.
- Block the hand-off to the user until every approach passes this audit.

After all 8 agents provide feedback and propose their individual approaches, **curate 5 distinct approaches** from the collective proposals.

### How to Curate the 5 Approaches

1. **Gather all 8 proposed approaches** from Phase 1
2. **Cluster similar approaches** together (if multiple agents propose similar solutions, merge into one)
3. **Ensure diversity**: The 5 approaches should represent meaningfully different trade-offs
4. **Select 5** that cover the spectrum from:
   - Most architecturally sound
   - Most secure/lowest risk
   - Most performant
   - Most user-friendly
   - Most practical/fastest to implement
   - Most modern/best-practice-aligned (from Scout)
   - Or other distinct perspectives that emerged from consultation

### Present 5 Approaches to User

Use AskQuestion tool to present all 5 approaches:

```
## 5 Proposed Approaches

### แนวทาง 1: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **เหมาะกับ**: [when to choose this]
- **ข้อเสนอแนะจาก Agent**: [which agents' feedback supports this]

### แนวทาง 2: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **เหมาะกับ**: [when to choose this]
- **ข้อเสนอแนะจาก Agent**: [which agents' feedback supports this]

### แนวทาง 3: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **เหมาะกับ**: [when to choose this]
- **ข้อเสนอแนะจาก Agent**: [which agents' feedback supports this]

### แนวทาง 4: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **เหมาะกับ**: [when to choose this]
- **ข้อเสนอแนะจาก Agent**: [which agents' feedback supports this]

### แนวทาง 5: [Short Name]
- **สรุป**: [1-2 sentences]
- **ข้อดี**: [key advantage]
- **Trade-off**: [key downside]
- **เหมาะกับ**: [when to choose this]
- **ข้อเสนอแนะจาก Agent**: [which agents' feedback supports this]
```

### Rules
- **NEVER proceed to Phase 3 without user selection** — this is a hard stop
- Each approach must include: summary, advantage, trade-off, when it's suitable, which agents support it
- If user provides a custom response (not from the list), treat it as a new approach and incorporate relevant agent feedback before proceeding
- If any agent flagged a **critical security or data integrity concern**, it must be noted in every approach where applicable

---

## PHASE 3: STRUCTURED IMPLEMENTATION PLAN

**Commander's role in Phase 3**:
- Convert the plan's "Changes Required" list into Commander **Task Packets** (TP-001, TP-002, …) using the exact format from `.cursor/skills/commander/SKILL.md` Phase 2.1.
- Each packet MUST include: assigned agent, prerequisite packets, scope IN/OUT, acceptance criteria, applicable rules (R1–R12), definition of done, and expected deliverables.
- Topologically sort packets; mark which can be dispatched in parallel and which share files (must serialize).
- Output a Dispatch Ledger alongside the implementation plan — this is the artifact that drives Phase 4.

Based on the user-selected approach and agent feedback, generate a structured implementation plan:

```
Goal: [Single clear sentence]
Context: [Current system state and why this change is needed]
Selected Approach: [Which approach was chosen and why]
Implementation Strategy: [High-level plan informed by agent feedback]

Changes Required:
  1. [File/Module] → [What changes and why]
  2. [File/Module] → [What changes and why]
  ...

Acceptance Criteria:
  1. [Measurable criterion]
  2. [Measurable criterion]
  ...

Agent Feedback Addressed:
  - Scout: [how their tech research was incorporated]
  - Architect: [how their feedback was incorporated]
  - Security: [how their feedback was incorporated]
  - Risk Analyst: [how their feedback was incorporated]
  - ...

Constraints:
  - Must not break: [existing functionality]
  - Must follow: [existing pattern/convention]
  - Performance budget: [if applicable]
```

## PHASE 4: COMMANDER-DRIVEN EXECUTION

Phase 4 is **fully owned by Commander**. Supreme no longer dispatches sub-agents directly — it hands the Dispatch Ledger from Phase 3 to Commander, which runs its full protocol (Phase 2 dispatch → Phase 3 laziness audit → Phase 4 integration checkpoints).

### Step 1: Hand-off to Commander
- Pass the Dispatch Ledger (all TP-XXX packets) to Commander.
- Commander re-verifies rule manifest coverage per packet before dispatching.
- Commander confirms the prerequisite topology and parallelization plan.

### Step 2: Commander Dispatches Specialist Agents
For each Task Packet, Commander selects the **narrowest, most specialized** agent from its competency map:
- **Analysis packets** → `/analysis` Light (or `/analysis` Deep for legacy/complex code).
- **Critical/core logic packets** → `/code` Critical (payment, auth, data integrity).
- **Standard feature packets** → `/code` Standard.
- **Refactor packets** → `/refactor`.
- **Migration packets** → `/migrate`.
- **Test packets** → `/test`.
- **UI polish packets** → `/redesign`.
- **Bug packets** → `/debug` → then `/code` Critical for the fix.
- **Docs packets** → `/docs-sync` (reference docs) or `/docs` (only if user-requested).
- **Review packets** → `/review` (or `/review-security` / `/review-performance` / `/review-quality` individually).

### Step 3: Per-Packet Laziness Audit (Mandatory)
After every agent report, Commander runs its Laziness Signal Matrix (22 signals) against the diff:
- `rg "TODO|FIXME|XXX|HACK"` on diff → must be 0
- `rg "t.Skip|\.skip\(|xit\("` on test files → must be 0
- `rg "panic\("` on Go diff → must be 0 (R2)
- `rg "interface\{\}|any\b"` on Go diff → each occurrence needs justification comment (R3)
- Frontend paths must match `web/core/modules/{feature}` (R7)
- No inline styles; Tailwind + dark mode default (R8)
- Cursor pagination must use HMAC-signed Base64 + indexed PK + capped limit (R9)
- App URL references must be `http://localhost:3000` only (R10)
- No external browser commands (R11)
- No unsolicited `.md` docs (R12)

Any failure → **REJECT** → Commander re-dispatches with added context. After 3 attempts → **ESCALATE** to `/code` Critical mode or notify user.

### Step 4: Integration Checkpoints (Commander Gates)
Between stages, Commander runs its checkpoints:
- **Post-Implementation**: read all diffs end-to-end, verify architectural coherence (R1 hexagonal boundaries intact).
- **Post-Test**: full test suite pass, zero regressions, coverage report on modified files.
- **Pre-Deploy**: automatic `/review` (or the 3 specialist reviews) must return GREEN.
- **Post-Deploy**: confirm commit on remote + CI pipeline status.

Any RED checkpoint → halt Phase 4 → remediate → retry checkpoint.

### Step 5: Refinement Loop (Commander-Audited)
Continue iterating under Commander supervision:
- **Self-Review**: Re-read ALL changes. Are expert feedback points from Phase 1 addressed? Commander cross-checks against the Phase 1 evidence ledger.
- **Edge Case Check**: What happens with unexpected inputs? Error paths? Commander rejects packets with untested error branches.
- **Optimization**: Can the code be simplified? Commander rejects packets that duplicate existing utilities found in Phase 0 Scout brief or Phase 1 Analysis.
- **Integration Check**: Does it work with the rest of the system? Commander re-runs integration checkpoint.
- If any issue found → Commander re-dispatches targeted packets (not whole re-work).

### Step 6: Termination Conditions
- All Task Packet acceptance criteria: GREEN
- All R1–R12 rule manifest items: verified with grep/read evidence
- Commander's Final Gate (12 items) from its Phase 5: ALL CHECKED
- Two consecutive clean refinement passes (no new issues found)
- OR: User explicitly requests to stop

## PHASE 5: FINAL VALIDATION & DELIVERY

**Commander's role in Phase 5**: Commander's Final Gate (its Phase 5) is **mandatory** and **blocking**. Supreme cannot announce completion while any gate item is unchecked.

### Step 1: Commander's Final Gate (12 items — ALL must be checked)
Run the full checklist from `.cursor/skills/commander/SKILL.md` Phase 5:

**Plan Fidelity**
- [ ] Every acceptance criterion from Phase 3: GREEN
- [ ] Every Task Packet: COMPLETED (not partial)
- [ ] Every rejection: RESOLVED with evidence

**Rule Compliance**
- [ ] R1–R12: each verified with grep/read evidence (cite commands run)
- [ ] Zero laziness signals in final diff
- [ ] Zero banned patterns (panic, unjustified any, bare errors)

**Quality Gates**
- [ ] Full test suite: PASS (with count)
- [ ] Lint: CLEAN (with count)
- [ ] /review verdict (if triggered): GREEN
- [ ] Regression check: zero existing tests broken

**Operational Readiness**
- [ ] Docker service reachable at `http://localhost:3000` (if UI-affecting)
- [ ] Rollback plan documented
- [ ] Risks flagged to user

**Scope Discipline**
- [ ] No files changed outside dispatched scope
- [ ] No unsolicited documentation files
- [ ] No new dependencies without justification

If any box is unchecked → **do NOT announce completion**. Return to Phase 4 for targeted remediation.

### Step 2: Deploy via Commander-Dispatched `/deploy`
- Commander dispatches TP-DEPLOY with explicit scope (no `--force`, no `reset --hard`, no history rewrite).
- After push, Commander reads remote HEAD + CI status to confirm success.

### Step 3: Produce TWO Reports (Supreme Delivery + Commander After-Action)

**Supreme Delivery Report**:
```
## Delivery Report
- **Changes Made**: [list of files modified with brief description]
- **Selected Approach**: [which approach was chosen and why]
- **Acceptance Criteria**: [all criteria with pass/fail status]
- **Agent Feedback Addressed**: [what was addressed from Phase 1]
- **Scout Intel Applied**: [what modern best practices were adopted]
- **Optimizations Applied**: [improvements beyond initial implementation]
- **Test Coverage**: [coverage on modified files]
- **Deploy Status**: [commit hash, branch, CI status]
- **Risks & Caveats**: [anything the user should be aware of]
- **Suggested Next Steps**: [follow-up actions if any]
```

**Commander After-Action Report** (from Commander Phase 6):
- Dispatch Ledger (every TP-XXX, agent, attempts, outcome, verification evidence)
- Rule Compliance Matrix (R1–R12 with grep evidence per rule)
- Rejection summary (what was rejected, why, how it was corrected)
- Bilingual summary section (EN + ไทย)

### Step 4: Failure Surfacing
If any criterion FAILS or any security concern is UNRESOLVED → clearly flag this to the user. Do NOT silently proceed.

## CRITICAL RULES (Always Active)

1. **Minimal Edits**: Prefer single-line changes over rewrites. Fix root causes, not symptoms.
2. **No Speculation**: Never assert facts about the codebase without verifying. Use tools to confirm.
3. **Parallel Exploration**: When searching the codebase, make multiple parallel tool calls to maximize efficiency.
4. **Preserve Style**: Match existing code conventions exactly.
5. **Test Before Trust**: Always verify changes work.
6. **Security First**: Any security concern from any agent must be flagged and addressed.
7. **Explain Assumptions**: Before implementing, list all assumptions. Surface misunderstandings early.
8. **Context Awareness**: Build a reliable mental map before writing code.
9. **Long-Horizon Mindset**: Don't settle for the first working solution. Continue improving.
10. **User Sovereignty**: The user always has the final say on which approach to take.
11. **Auto-install**: If any tool/dependency is missing, install it automatically. Only ask user for: framework selection, major version upgrades, global packages, significant lock file changes.
12. **Auto-run safe commands**: git status, git diff, git log, test runs, lint/compile checks — all auto-run.
13. **Ask before destructive commands**: git push, git reset, git rebase, git clean, database changes — ASK USER always.
14. **Scout First**: Always run Scout research before making technology decisions. Never rely on outdated knowledge.
15. **Commander Always On**: Commander is active in EVERY phase — no sub-agent is dispatched without a Task Packet, no deliverable is accepted without a Laziness Audit, no phase advances without a Commander checkpoint.
16. **Rule Manifest R1–R12**: Every Commander-audited packet enforces R1–R12 from `.cursor/skills/commander/SKILL.md`. Any violation triggers immediate rejection.
17. **Three Strikes Escalation**: If any agent fails the same packet 3 times, Commander escalates to `/code` Critical mode (for code) or back to the user with evidence and hypothesis — never silently moves on.
18. **Dual Reporting**: Final delivery always includes BOTH the Supreme Delivery Report AND the Commander After-Action Report. Missing either means delivery is incomplete.
