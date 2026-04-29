# Master Planner — Advanced Reference

## Planning Frameworks

### 1. The CRAFT Framework (for feature design)

When designing a new feature, work through CRAFT:

- **Context**: What is the current state? Why is this change needed NOW?
- **Requirements**: What must be true for this to be successful? (Not "nice to have" — MUST)
- **Approach**: What is the highest-ROI implementation strategy?
- **Failure Modes**: What are the top 5 ways this could go wrong?
- **Tests of Success**: How will we KNOW it worked? (Measurable criteria)

### 2. The IMPACT Matrix (for prioritizing tasks within a plan)

Rate each task on two axes:

```
         High Impact
              |
    Do First  |  Do Second
              |
Low Effort ───+─── High Effort
              |
    Do Last   |  Skip / Defer
              |
         Low Impact
```

- **Do First**: High impact, low effort → immediate wins
- **Do Second**: High impact, high effort → core value
- **Do Last**: Low impact, low effort → polish
- **Skip/Defer**: Low impact, high effort → not worth it now

### 3. The REVERSIBILITY Test

Before deciding rigor level, classify each decision:

| Type | Example | Rigor Level |
|------|---------|-------------|
| **Easily Reversible** | UI text, config values, feature flags | Quick decision, minimal analysis |
| **Moderately Reversible** | API design, DB schema, library choice | Standard analysis, document trade-offs |
| **Nearly Irreversible** | Architecture pattern, data migration, auth system | Deep analysis, get stakeholder buy-in |

Apply analysis effort proportional to irreversibility.

### 4. The PRE-MORTEM Technique

Before finalizing any plan, run a pre-mortem:

```
## Pre-Mortem Exercise

Imagine it's 3 months from now and this project FAILED.
Write 5 specific reasons why it failed:

1. [Specific failure mode]
2. [Specific failure mode]
3. [Specific failure mode]
4. [Specific failure mode]
5. [Specific failure mode]

For each failure, what can we do NOW to prevent it?
- Prevention for #1: [specific action in the plan]
- Prevention for #2: [specific action in the plan]
...
```

### 5. The MIGRATION SAFETY Protocol

When a plan involves data/schema changes:

```
## Migration Safety Checklist

- [ ] Migration is reversible (down migration exists and is tested)
- [ ] Migration handles existing data (no data loss)
- [ ] Migration is idempotent (safe to run multiple times)
- [ ] New code works with BOTH old and new schema (during rollout)
- [ ] Old code fails gracefully with new schema (during rollback)
- [ ] Large data changes are batched (no table locks)
- [ ] Foreign key constraints are respected
- [ ] Index creation doesn't block writes (use CONCURRENTLY)
```

---

## Advanced Planning Patterns

### Pattern: Strangler Fig

For replacing legacy systems incrementally:

```
Phase 1: Route NEW traffic to new system (behind feature flag)
Phase 2: Migrate one use case at a time to new system
Phase 3: Once all use cases migrated, remove old system
```

Key: Both systems coexist during transition. Each phase is independently valuable and reversible.

### Pattern: Feature Flag Driven Delivery

For risky or gradual rollouts:

```
Phase 1: Deploy code behind feature flag (OFF) — no user impact
Phase 2: Enable for internal users — test in production
Phase 3: Enable for X% of users — gradual rollout
Phase 4: Enable for all users — full release
Phase 5: Remove feature flag — cleanup
```

### Pattern: Blue-Green Deployment

For zero-downtime deployments:

```
Phase 1: Deploy new version to "green" environment
Phase 2: Run smoke tests against green
Phase 3: Switch traffic from "blue" to "green"
Phase 4: Monitor for issues
Phase 5: If issues → switch back to "blue" (instant rollback)
Phase 6: If stable → "green" becomes new "blue"
```

### Pattern: Event Sourcing for Complex State

When the feature involves complex state transitions:

```
Instead of: UPDATE tasks SET status = 'done' WHERE id = 1
Use: INSERT INTO task_events (task_id, event_type, payload) VALUES (1, 'status_changed', '{"from":"in_progress","to":"done"}')

Benefits: Full audit trail, replay capability, temporal queries
Trade-off: Higher complexity, eventual consistency
```

---

## Decision-Making Techniques

### Technique: A/B/C Speed Decisions

When the team needs to make a decision quickly:

- **A-type** (2 min): Reversible, low impact → Decider picks, informs team
- **B-type** (15 min): Reversible, medium impact → Decider proposes, 1 objection blocks
- **C-type** (1 hour): Irreversible, high impact → Team discusses, consensus required

### Technique: The 10/10/10 Test

For evaluating whether a choice matters:

- How will we feel about this decision in **10 minutes**?
- How will we feel about this decision in **10 months**?
- How will we feel about this decision in **10 years**?

If only the 10-minute view matters → it's a low-stakes decision. Move fast.

### Technique: Happiness Metric

For user-facing features, measure by:

```
Before: "How do users accomplish X today?" → Pain level: 1-10
After:  "How will users accomplish X with this?" → Pain level: 1-10

If delta < 3 → reconsider the approach. Is this worth building?
If delta >= 5 → this is a high-value feature. Prioritize quality.
```

---

## Common Anti-Patterns to Avoid

### Anti-Pattern: The Everything Plan
**Symptom**: One massive phase with 50 tasks
**Fix**: Break into phases where each delivers user value independently

### Anti-Pattern: The Assumption Cascade
**Symptom**: Task 5 depends on Task 4 depends on Task 3... all assumptions
**Fix**: Verify critical assumptions before planning dependent tasks. Add verification tasks first.

### Anti-Pattern: The Perfect Plan
**Symptom**: Planning takes longer than implementation
**Fix**: Time-box planning. A good plan today beats a perfect plan next week.

### Anti-Pattern: The Generic Plan
**Symptom**: Plan could apply to any project, not specific to THIS codebase
**Fix**: Reference specific files, functions, patterns from the actual codebase. "Update handler" → "Add validation to `api/internal/handler/task.go:CreateTask()`"

### Anti-Pattern: The No-Rollback Plan
**Symptom**: No way to undo changes if things go wrong
**Fix**: Every phase must have a rollback strategy. Feature flags, migrations with down scripts, or config-based toggles.

---

## Plan Communication Templates

### Template: Quick Plan (for small features)

```
## Plan: [Feature Name]

**Goal**: [1 sentence]
**User Impact**: [1 sentence about what changes for users]

**Approach**: [chosen approach in 2-3 sentences]

**Tasks**:
1. [File] → [Change] → [Why]
2. [File] → [Change] → [Why]
3. [File] → [Change] → [Why]

**Acceptance**:
- [ ] [measurable criterion 1]
- [ ] [measurable criterion 2]

**Risk**: [top risk + mitigation in 1 sentence]
```

### Template: Full Plan (for major features)

```
## Master Plan: [Feature Name]

### Executive Summary
[2-3 sentences: what, why, approach]

### Key Insights
1. [Most important finding from analysis]
2. [Second most important]
3. [Third most important]

### Recommended Approach
[Name] — [why in 1-2 sentences]

### Implementation Phases
[Phase 1: Foundation] → [Phase 2: Core] → [Phase 3: Polish]

### Acceptance Criteria
[Numbered list of measurable criteria]

### Top Risks
[Risk matrix with mitigations]

### Effort Estimate
[Total: X phases, estimated Y-Z days]
```

---

## Estimation Techniques

### Technique: Three-Point Estimate

For each task:
- **Optimistic** (O): Everything goes perfectly
- **Most Likely** (M): Normal complications
- **Pessimistic** (P): Things go wrong

Expected = (O + 4M + P) / 6

### Technique: Complexity Multiplier

Base estimate × complexity factor:

| Complexity | Multiplier | Example |
|-----------|-----------|---------|
| Simple CRUD | ×1.0 | Add a field to existing entity |
| Standard Feature | ×1.5 | New entity with relationships |
| Complex Feature | ×2.0 | Multi-entity workflow |
| Cross-Cutting | ×3.0 | Auth, logging, shared state |
| Migration | ×3.5 | Schema change with data transform |
