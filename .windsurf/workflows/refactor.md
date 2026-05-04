---
auto_execution_mode: 0
description: Refactor - restructure existing code to improve quality without changing behavior
---

# /refactor — Safe Code Refactoring Protocol

You are a senior software engineer performing a refactoring. Your goal is to improve code structure, readability, and maintainability WITHOUT changing external behavior. Every refactoring must be behavior-preserving.

## Step 1: Identify Refactoring Target

Clarify what needs refactoring:
- If the user specifies a target → analyze it
- If not → ask: "Which module/function/file should I refactor?"
- If /analysis was run → use its quality findings as refactoring targets

Common refactoring triggers:
- Function/method too long (>50 lines)
- High cyclomatic complexity (>10)
- Code duplication across files
- Tightly coupled modules
- Poor naming that obscures intent
- Missing abstraction layers

## Step 2: Baseline Capture

Before touching any code, establish a safety net:

1. **Run existing tests**: Confirm all pass. This is your baseline.
2. **Identify test coverage**: If coverage is low on the target code → **STOP and write tests first** (call /test)
3. **Document current behavior**: Briefly note what the code does externally (API contracts, return values, side effects)

**CRITICAL**: Never refactor code without tests. If no tests exist, write them first.

## Step 3: Refactoring Strategy

Choose the appropriate refactoring technique:

| Smell | Technique | Description |
|-------|-----------|-------------|
| Long function | Extract Function | Break into smaller named functions |
| Duplication | Extract Shared Utility | Move common logic to shared module |
| Complex conditionals | Replace Conditional with Polymorphism/Strategy | Use pattern matching or strategy pattern |
| Large class | Extract Class | Split responsibilities |
| Tight coupling | Introduce Interface/Dependency Injection | Decouple via abstraction |
| Poor naming | Rename Variable/Function | Use descriptive names |
| Nested conditionals | Guard Clauses / Early Return | Flatten control flow |
| Magic numbers | Replace with Named Constants | Extract meaningful constants |
| Feature envy | Move Method | Move method to the class it belongs to |
| Dead code | Remove Dead Code | Delete unused functions, imports, variables |

## Step 4: Execute Refactoring (Small Steps)

Refactor in tiny, verifiable steps. Each step must leave the code in a working state:

### For each refactoring step:
1. **Make ONE small change** (e.g., extract one function, rename one variable)
2. **Run tests immediately** — all must pass
3. **If tests fail** → revert the change, analyze why, try again
4. **If tests pass** → commit mentally (or suggest a git checkpoint) and proceed

### Refactoring Order:
1. **Rename first**: Fix naming — this makes subsequent refactoring easier to reason about
2. **Extract next**: Pull out functions, classes, utilities
3. **Restructure last**: Change relationships between modules (highest risk)

### Rules During Refactoring:
- **Never mix refactoring with feature changes**: If you find a bug during refactoring, note it but fix it separately
- **Never change public API**: External contracts must remain identical
- **Never change test expectations**: Tests verify behavior. If you need to change a test, you're changing behavior — that's not refactoring
- **Keep each step small**: If a step touches more than 3 files, it's too big

## Step 5: Post-Refactoring Verification

```
## Refactoring Verification

- [ ] All existing tests pass (same as baseline)
- [ ] No new test failures introduced
- [ ] Public API contracts unchanged
- [ ] No behavior changes (same inputs → same outputs)
- [ ] Code is more readable than before
- [ ] Code is less complex than before (fewer lines, lower cyclomatic complexity)
- [ ] No new dependencies introduced (unless abstraction requires it)
- [ ] No performance regression
```

## Step 6: Refactoring Report

```
## Refactoring Report

**Target**: [module/function refactored]
**Technique(s) Used**: [list techniques applied]

### Changes Made
| Step | Change | Files Modified | Tests Status |
|------|--------|---------------|-------------|
| 1 | [what was done] | [files] | ✅ Pass |

### Metrics
| Metric | Before | After |
|--------|--------|-------|
| Lines of code | [n] | [n] |
| Functions | [n] | [n] |
| Cyclomatic complexity | [n] | [n] |
| Duplication instances | [n] | [n] |

### Behavior Verification
- All existing tests: ✅ PASS (same as baseline)
- No new behavior introduced
- No behavior removed
```

## Critical Rules

1. **No tests = no refactor**: If the target code lacks tests, write them first.
2. **Behavior preservation**: If external behavior changes, it's not refactoring — it's a feature change.
3. **Small steps**: Each refactoring step should be independently verifiable.
4. **Run tests after every step**: Don't accumulate multiple refactorings before testing.
5. **No feature creep**: Don't add "improvements" that change behavior. Only restructure.
6. **Revert on failure**: If a refactoring step breaks tests, revert immediately and try a different approach.
