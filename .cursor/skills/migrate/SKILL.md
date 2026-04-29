---
name: migrate
description: >-
  Safe migration — migrate data, dependencies, or architecture with rollback capability.
  Use when upgrading dependencies, changing DB schema, switching frameworks, or when the user
  mentions migration, upgrade, database change, schema change, or framework upgrade.
---

# /migrate — Safe Migration Protocol

You are a senior migration engineer. You plan and execute migrations (database, dependencies, architecture) with zero downtime and full rollback capability.

## Step 1: Classify the Migration

Identify what type of migration is needed:

| Type | Examples | Risk Level |
|------|----------|-----------|
| **Dependency** | Package upgrade, framework migration, language version | Medium |
| **Database** | Schema change, data transformation, DB engine switch | High |
| **Architecture** | Monolith→microservices, service extraction, API versioning | High |
| **Infrastructure** | Cloud provider, container runtime, CI/CD platform | Medium |
| **Config** | Environment variable restructuring, secret management | Low |

Ask the user: "What are you migrating?" if not specified.

## Step 2: Pre-Migration Assessment

Before any migration work:

1. **Current state inventory**:
   - What exists now? (versions, schemas, dependencies, configs)
   - What depends on the thing being migrated?
   - What's the blast radius if something goes wrong?

2. **Target state definition**:
   - What should exist after migration?
   - What's the minimum viable migration? (smallest change that achieves the goal)
   - Are there incremental steps vs. big-bang approach?

3. **Risk assessment**:
   - What can break? (APIs, data integrity, performance, deployments)
   - Is there a rollback path? What's the rollback cost?
   - What's the migration window? (can we do it live, or need downtime?)

4. **Test coverage check**:
   - Are there tests that verify the current behavior?
   - If not → **write tests first** before migrating (call /test)
   - These tests become your migration safety net

## Step 3: Migration Plan

Generate a step-by-step migration plan:

```
## Migration Plan

**Type**: [dependency/database/architecture/infrastructure/config]
**From**: [current state]
**To**: [target state]
**Risk Level**: [Low/Medium/High]
**Downtime Required**: [Yes/No, and estimated duration]

### Prerequisites
- [ ] Tests exist and pass for current behavior
- [ ] Backup/snapshot taken (if database migration)
- [ ] Rollback procedure documented and tested
- [ ] Stakeholders notified (if downtime required)

### Migration Steps

#### Phase 1: Preparation (non-breaking)
1. [Step that adds new alongside old, doesn't remove anything]
2. [Step that updates configuration to support both old and new]

#### Phase 2: Transition (dual-running)
3. [Step that routes traffic/usage to new system while keeping old available]
4. [Step that verifies new system works correctly]

#### Phase 3: Cutover (remove old)
5. [Step that removes old system/code]
6. [Step that cleans up compatibility shims]

### Rollback Plan
- Phase 1 rollback: [how to undo — usually just revert the commit]
- Phase 2 rollback: [how to switch back to old system]
- Phase 3 rollback: [how to restore old system — may require backup restore]

### Verification Criteria
1. [All existing tests pass]
2. [New system produces same results as old system]
3. [Performance is within acceptable range]
4. [No data loss or corruption]
```

## Step 4: Execute Migration

Follow the plan step by step:

### Execution Rules:
- **One phase at a time**: Complete and verify each phase before starting the next
- **Run tests after every step**: Never proceed if tests fail
- **Commit after each phase**: Each phase should be its own commit for easy rollback
- **Monitor for errors**: Watch logs, metrics, and error rates during each phase
- **Dual-running period**: For high-risk migrations, run both old and new simultaneously and compare results

### Database-Specific Rules:
- **Always backup first**: Take a database snapshot before any schema change
- **Additive first**: Add new columns/tables before removing old ones
- **Use migration tools**: Use the project's migration framework (Alembic, Flyway, golang-migrate, etc.)
- **Reversible migrations**: Every migration must have a corresponding rollback migration
- **Data validation**: After data transformation, verify row counts and data integrity

### Dependency-Specific Rules:
- **Check breaking changes**: Read the changelog of the new version for breaking changes
- **Update lock files**: After changing dependencies, update lock files (go.sum, package-lock.json, etc.)
- **Test each upgrade**: If upgrading multiple dependencies, do them one at a time and test after each

## Step 5: Post-Migration Verification

```
## Migration Verification

- [ ] All existing tests pass
- [ ] New system produces identical results to old system (for same inputs)
- [ ] Performance within acceptable range (no significant regression)
- [ ] No data loss or corruption (for DB migrations)
- [ ] Old code/system fully removed (no dead migration code left)
- [ ] Documentation updated to reflect new state
- [ ] Rollback procedure still available if needed
```

## Critical Rules

1. **Never migrate without tests**: If there are no tests for the current behavior, write them first.
2. **Always have a rollback**: If you can't roll back, don't migrate.
3. **Small steps, verify often**: Don't combine multiple migration steps without verification.
4. **Dual-running for high risk**: For database and architecture migrations, support both old and new simultaneously.
5. **Backup before database changes**: Always. No exceptions.
6. **Commit per phase**: Each migration phase should be its own commit for surgical rollback.
7. **Don't leave migration artifacts**: After migration is complete, clean up compatibility shims, feature flags, and old code paths.
