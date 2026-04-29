---
name: docs-sync
description: >-
  Agent documentation sync — create and update DEPENDENCY_MAP, USER_CONTEXT, API_CATALOG
  to match current codebase. Use when synchronizing agent reference docs, updating documentation
  maps, or when the user mentions docs-sync, sync docs, dependency map, or API catalog.
---

# /docs-sync: Agent Documentation Synchronization

Create or update reference documentation to match the current codebase state.

## Managed Documents

| Document | Purpose |
|----------|---------|
| `DEPENDENCY_MAP.md` | Cross-system dependency map |
| `USER_CONTEXT.md` | User context and personas |
| `API_CATALOG.md` | API endpoint catalog |

## Step 1: Check Existing Documents

Check whether all 3 document sets exist:
- If a doc exists in project root → read current content → mark as "needs update"
- If not found → mark as "needs creation"

## Step 2: Scan Current Code

Scan code to collect latest information:

### 2.1 API Endpoints
- Search for route definitions in all relevant route files
- Check handler files for any new handlers not yet in routes

### 2.2 Frontend Pages & API Usage
- List all pages in the frontend pages directory
- Search for service classes that call APIs
- Check the base HTTP client configuration

### 2.3 Database Models
- Search for all model/entity definitions
- Check repository layer for DB queries

### 2.4 Inter-service Communication
- Search for outgoing HTTP calls between services
- Search for database access points (MongoDB, PostgreSQL, etc.)

### 2.5 Permission Middleware
- Search for permission middleware definitions
- Search for access control structures

## Step 3: Compare and Update Documents

### 3.1 DEPENDENCY_MAP.md

**Verify:**
- [ ] Service Architecture Overview — any services added/removed or port changes?
- [ ] Frontend ↔ Backend API Contract Map — new modules? new pages? endpoint changes?
- [ ] Shared Database Tables — new tables? new shared tables?
- [ ] Inter-service Integration — endpoint changes? auth changes?
- [ ] Permission Middleware — new middleware?
- [ ] Frontend Shared State — new stores?
- [ ] High-Risk Change Zones — new risk zones?

**If changes found → update that section of the document**

### 3.2 USER_CONTEXT.md

**Verify:**
- [ ] User Personas — new user roles?
- [ ] Key Workflows — new workflows? (check for new pages)
- [ ] Feature Criticality — new daily-use features?
- [ ] Pain Points — new issues?
- [ ] UX Decision Guide — are the iron rules still correct?

**If changes found → update that section of the document**

### 3.3 API_CATALOG.md

**Verify:**
- [ ] API Base URLs — env var changes? port changes?
- [ ] Every endpoint group — added/removed/changed endpoints?
- [ ] Permission middleware — changed permissions?
- [ ] Response Format — format changes?

**If changes found → update that section of the document**

## Step 4: Create New Documents (if missing)

If any document doesn't exist, create it using data from Step 2.

## Critical Rules

1. **Don't write docs from memory**: All information must come from scanning actual code
2. **Don't delete existing data unnecessarily**: Only add/edit, unless data is definitively wrong
3. **Preserve existing structure**: If a doc already exists, update only the changed sections — don't rewrite the entire file
4. **Scan thoroughly**: Don't skip any modules, especially newly created ones
5. **Always check permissions**: Every time an endpoint is added, specify the correct permission
6. **Frontend mapping must be accurate**: If adding an endpoint, find the frontend page that uses it
