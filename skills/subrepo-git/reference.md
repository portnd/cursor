# Sub-Repository Git — Reference

## Dynamic discovery

When the known repo list may be stale, discover nested git repos from the workspace root:

```bash
find . -name ".git" -type d \
  ! -path "./.cursor/*" \
  ! -path "./.windsurf/*" \
  ! -path "./node_modules/*" \
  2>/dev/null
```

Parent directory of each `.git` folder is one sub-repo. Example: `./thai-iod-api-service/.git` → repo path `thai-iod-api-service`.

## Exclusion list (default)

| Path pattern | Reason |
|--------------|--------|
| `.cursor/**` | IDE/agent metadata |
| `.windsurf/**` | Editor workflows |
| `node_modules/**` | Dependencies (rare `.git`, skip if present) |

## Origin URLs (iod workspace, May 2026)

```
thai-iod-api-service     → https://gitlab.komgrip.co.th/komgrip/thai-iod-api-service.git
thai-iod-web             → https://gitlab.komgrip.co.th/komgrip/thai-iod-web.git
thai-iod-back-office-web → https://gitlab.komgrip.co.th/komgrip/thai-iod-back-office-web.git
```

## Common branch names (inform user, do not auto-select)

- `develop` — default on origin (`remotes/origin/HEAD` → `develop`)
- `refactoring/unit_test` — often used locally in this workspace
- Feature branches — `feature/*`, `bugfix/*`, etc. on origin

## Troubleshooting

| Symptom | Check | Fix (with user approval) |
|---------|-------|---------------------------|
| `fatal: not a git repository` | Wrong `cd` path | `cd` into sub-repo directory |
| `origin` missing | `git remote -v` empty | `git remote add origin <url>` |
| Auth failed on fetch | GitLab credentials / SSH | User fixes credential helper or SSH key |
| Branch not found | `git branch -a` | `git switch -c <b> --track origin/<b>` after fetch |
| Dirty tree blocks switch | `git status` | stash / commit / skip per user choice |

## One-shot status script (optional)

Agent may run this from workspace root for a quick Phase 1 report:

```bash
ROOT="$(pwd)"
for d in thai-iod-api-service thai-iod-web thai-iod-back-office-web; do
  [ -d "$ROOT/$d/.git" ] || continue
  echo "======== $d ========"
  git -C "$ROOT/$d" remote get-url origin 2>/dev/null || echo "origin: (none)"
  git -C "$ROOT/$d" fetch origin --prune 2>&1 | tail -1
  git -C "$ROOT/$d" status -sb
  echo
done
```
