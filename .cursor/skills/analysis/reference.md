# Reverse Engineering & Analysis — Advanced Techniques Reference

## Technique 1: Blame Archaeology

Use git history to reconstruct the evolution of code:

```bash
# See who wrote each line and when
git blame -C -C -M <file>

# See how a function evolved over time
git log -p --follow -S "function_name" -- <file>

# Find when a specific pattern was introduced
git log --all --oneline --grep="pattern"

# See the state of code at a specific commit
git show <commit>:<file>
```

### What Git History Reveals:
- **Commit messages**: Author's intent at the time of change
- **Diff patterns**: Rushed changes (large diffs, late night commits) vs careful ones
- **Branch patterns**: Feature branches (thoughtful) vs direct commits (hurried)
- **Timing**: Rapid successive commits (debugging session) vs spaced commits (iterative development)

## Technique 2: Dependency Inversion Tracing

When understanding a complex dependency graph, work backwards:

1. **Start from the effect**: What's the observable outcome?
2. **Find the producer**: What code produces this outcome?
3. **Trace the trigger**: What calls this producer?
4. **Follow the chain**: Keep going until you reach the root cause
5. **Map the full chain**: entry → middleware → handler → usecase → repo → DB

## Technique 3: Interface Contract Extraction

Reverse-engineer the contract from implementation:

```
Given a function, extract:
1. Input constraints (type, range, format, required/optional)
2. Output guarantees (type, range, null possibility)
3. Side effects (DB writes, API calls, state mutations)
4. Error conditions (what can fail, how it fails)
5. Idempotency (is it safe to call multiple times?)
6. Ordering requirements (must something happen first?)
7. Concurrency safety (thread-safe? needs locking?)
```

## Technique 4: Execution Simulation

For complex logic, mentally simulate execution:

```
## Simulation Worksheet

Given input: [specific input]
Step 1: [variable] = [value] — because [reasoning]
Step 2: [condition] evaluates to [true/false] — because [reasoning]
Step 3: enters [branch] — because [reasoning]
...
Final state: [output and side effects]
```

## Technique 5: Comparative Analysis

Compare similar code to find patterns:

- Compare multiple handlers to find the "template" the author used
- Compare error handling across modules to find inconsistency
- Compare naming conventions to detect multiple authors
- Compare test patterns to understand testing philosophy

## Technique 6: Dead Code Archaeology

Dead code reveals abandoned intent:

- Commented-out code → "I tried this approach and it didn't work"
- TODO/FIXME/HACK → "I know this is wrong but I had to ship"
- Unused imports → "I refactored but forgot to clean up"
- Unreachable branches → "I planned for this case but never implemented it"

## Technique 7: Meta-Pattern Detection

Identify patterns that span multiple files:

| Pattern | What It Reveals |
|---------|----------------|
| Consistent error wrapping | Author values error traceability |
| Inconsistent error handling | Multiple authors or rushed development |
| Extensive logging in one area | That area was previously buggy |
| Over-engineered abstractions | Author anticipated future needs (rightly or wrongly) |
| Copy-paste with minor changes | Author was unsure how to abstract |
| Defensive checks everywhere | Author was burned by edge cases before |
| Missing tests for complex logic | Logic was written before tests were valued |

## Technique 8: Architecture Heat Mapping

Visually map which areas of the codebase are healthy vs fragile:

```
## Architecture Heat Map

🔴 HOT (Fragile / High Risk):
- [module]: [why it's fragile]

🟡 WARM (Needs Attention):
- [module]: [what needs attention]

🟢 COOL (Healthy):
- [module]: [why it's healthy]

❄️ COLD (Dead / Unused):
- [module]: [evidence it's unused]
```

## Technique 9: Dependency Coupling Score

Measure how tightly coupled modules are:

```
## Coupling Analysis

| Module | Depends On | Depends On It | Coupling Score | Risk |
|--------|-----------|---------------|----------------|------|
| [A] | [B, C, D] | [E, F] | High (3 in, 2 out) | Changes ripple widely |
| [B] | [A] | [none] | Low (1 in, 0 out) | Safe to modify |
```

Scoring: 0-2 = Low, 3-4 = Medium, 5+ = High coupling

## Technique 10: Database Schema Archaeology

Reverse-engineer the domain model from the database:

```bash
# Find all migration files
find . -path "*/migrations/*" -o -path "*/migrate/*" | sort

# Find all model/entity definitions
grep -r "type.*struct" --include="*.go" -l
```

From schema, infer:
- Core domain entities (tables with most relationships)
- Event-sourced vs CRUD patterns (presence of event tables)
- Soft delete patterns (deleted_at columns)
- Audit trails (created_by, updated_by columns)
- Multi-tenancy (tenant_id/org_id columns)

## Technique 11: API Surface Analysis

Map the complete API surface from handlers/routes:

```bash
# Find all route registrations
grep -r "\.GET\|\.POST\|\.PUT\|\.DELETE\|\.PATCH" --include="*.go"

# Find all handler functions
grep -r "func.*Handler\|func.*handler" --include="*.go"
```

For each endpoint, extract:
- HTTP method + path
- Auth requirement (public / authenticated / admin)
- Request body schema
- Response body schema
- Error responses
