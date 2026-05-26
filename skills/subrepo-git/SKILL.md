---
name: subrepo-git
description: >-
  Connects nested git sub-repositories in the workspace to their origin remotes,
  fetches latest refs, reports status, and stops until the user specifies which
  branch to use per repo. Use when the user asks to connect sub-repos to origin,
  sync git for thai-iod-api-service, thai-iod-web, or thai-iod-back-office-web,
  เชื่อมต่อ git repo ย่อย, fetch origin, or switch branch in nested repos.
---

# Sub-Repository Git Connect & Branch Gate

The workspace root (`iod`) is **not** a git repository. Application code lives in **nested repos** with their own `.git` directories. This skill connects them to `origin`, refreshes remote refs, then **waits for an explicit branch choice** before any checkout.

## Known sub-repositories (this workspace)

| Directory | Origin (default) |
|-----------|------------------|
| `thai-iod-api-service/` | `https://gitlab.komgrip.co.th/komgrip/thai-iod-api-service.git` |
| `thai-iod-web/` | `https://gitlab.komgrip.co.th/komgrip/thai-iod-web.git` |
| `thai-iod-back-office-web/` | `https://gitlab.komgrip.co.th/komgrip/thai-iod-back-office-web.git` |

For discovery rules and exclusions, see [reference.md](reference.md).

## Phase 1 — Discover & connect (auto-run)

Run for **each** target sub-repo (all three unless the user names specific dirs):

```bash
REPO_ROOT="/Users/tnd/Documents/iod"   # or workspace root
SUBREPO="thai-iod-api-service"         # repeat per repo

cd "$REPO_ROOT/$SUBREPO"
git rev-parse --is-inside-work-tree
git remote -v
git branch -vv
git status -sb
```

### 1.1 Verify or set `origin`

- If `origin` exists → record URL; do **not** change it unless the user asks.
- If `origin` is missing → **ask the user** for the remote URL, then:
  ```bash
  git remote add origin <URL>
  ```
- If `origin` URL is wrong → show current vs expected; **ask** before `git remote set-url origin <URL>`.

### 1.2 Connect to remote (safe, no branch change)

```bash
git fetch origin --prune
```

Optional (only if user requested full sync metadata):

```bash
git remote show origin
```

### 1.3 Collect status per repo

Record for the report:

- Current branch and upstream (`git status -sb`)
- Whether local is ahead/behind origin
- Uncommitted changes (modified / untracked counts)
- Whether the requested branch exists locally or only on `remotes/origin/*`

**Do not run** `git checkout`, `git switch`, `git pull`, `git merge`, `git rebase`, or `git reset` in Phase 1.

## Phase 2 — Branch gate (mandatory stop)

After Phase 1, **stop and ask the user** which branch to use. Do not guess `develop` or any default.

Present a summary table, then ask clearly:

```markdown
## Sub-repo git status (connected to origin)

| Repo | Origin | Current branch | vs origin | Local changes |
|------|--------|----------------|-----------|-----------------|
| thai-iod-api-service | ✓ | refactoring/unit_test | even/ahead/behind | M: N, ?: M |
| thai-iod-web | ✓ | ... | ... | ... |
| thai-iod-back-office-web | ✓ | ... | ... | ... |

**รอคำสั่งจากคุณ:** ต้องการให้แต่ละ repo อยู่ branch อะไร?
- ระบุทีละ repo (เช่น `api: develop`, `web: develop`, `backoffice: feature/foo`)
- หรือ branch เดียวกันทุก repo (เช่น `ทุก repo ใช้ develop`)
- หรือ `คง branch ปัจจุบัน` สำหรับ repo ที่ไม่ระบุ

จะ checkout / pull ก็ต่อเมื่อคุณยืนยัน branch แล้วเท่านั้น
```

If the user only asked to "connect origin" without mentioning branches, Phase 2 is the **final step** until they reply.

## Phase 3 — Switch branch (only after explicit user command)

Proceed only when the user names branch(es). Map aliases:

| User says | Directory |
|-----------|-----------|
| api, backend, service | `thai-iod-api-service` |
| web, frontend | `thai-iod-web` |
| backoffice, back-office, bo | `thai-iod-back-office-web` |

### 3.1 Pre-checkout safety

For each repo to switch:

1. If **uncommitted changes** exist → warn and ask: stash, commit, discard, or skip this repo.
2. Confirm branch name exists: `git branch -a | grep -E 'branch-name$'`
3. **Ask once more** if switching away from a branch with unpushed commits or dirty tree.

### 3.2 Checkout workflow

```bash
cd "$REPO_ROOT/$SUBREPO"
git fetch origin --prune

# Prefer local branch if it exists and tracks origin
git switch <branch> 2>/dev/null || git switch -c <branch> --track origin/<branch>
```

- Use `git switch -c <branch> --track origin/<branch>` when the branch exists only on remote.
- **Do not** `git pull` unless the user asked to pull/update working tree after checkout.
- **Never** `git push --force` unless explicitly requested.

### 3.3 Post-switch report

```markdown
## Branch switch result

| Repo | Previous | Now | Tracking | Pull needed? |
|------|----------|-----|----------|--------------|
| thai-iod-api-service | ... | ... | origin/... | yes/no |

Notes: [any stash, skipped repo, or error]
```

## Critical rules

1. **Workspace root is not git** — run all git commands inside sub-repo directories, not `/Users/tnd/Documents/iod` (unless that path becomes a repo later).
2. **Exclude tooling repos** — never touch `.cursor/.git`, `.windsurf/**/.git`, or other nested `.git` under config folders unless the user explicitly includes them.
3. **Connect ≠ checkout** — `git fetch` is allowed in Phase 1; branch changes require Phase 3 and user-specified branch names.
4. **No silent branch picks** — do not switch to `develop`, `main`, or remote `HEAD` default without user instruction.
5. **Ask before destructive ops** — `push`, `reset --hard`, `clean -fd`, `rebase`, force push.
6. **Parallel status checks** — run `git remote`, `git fetch`, `git status` for multiple sub-repos in parallel when possible.
7. **Respect user git rules** — no `git config` changes; no `--no-verify` unless user asks.

## When user intent is ambiguous

| User request | Action |
|--------------|--------|
| "เชื่อม origin" / "connect git" | Phase 1 + Phase 2 only |
| "ไป branch develop" | Phase 1 (if not done) → Phase 3 for named branch |
| "fetch อย่างเดียว" | Phase 1 only, no Phase 2 question about checkout unless they also want status |
| "pull latest" | Ask which repo + branch first, then fetch + pull on current or specified branch |

## Related skills

- **deploy** — commit/push after code is ready; still ask before push.
- **continue** — resume work; read `git status` per sub-repo being edited.
