---
name: continue
description: >-
  Resume work from where it stopped after credit limit reset or interruption.
  Use when continuing incomplete work, resuming after a break, or when the user mentions
  continue, resume, pick up where we left off, or keep going.
---

# /continue: Resume Work from Stopped Point

Use when the agent stopped mid-task because of credit limit or interruption. After resetting credits, use this to resume work.

## Step 1: Read Latest State

Read these to know what was being done and how far it got:

1. **Read `progress.txt`** (if it exists) — summary of progress written by the agent
2. **Read todo list** — check for pending todos (status: `in_progress` or `pending`)
3. **Read recently modified files** — `git diff --stat` to see which files were already modified

```
1. Read progress.txt (if exists)
2. Check current todo list
3. Run: git diff --stat
4. Run: git status
```

## Step 2: Identify Stop Point

From Step 1 data, identify:

- **Task in progress**: [what]
- **Completed**: [what's done]
- **Pending**: [what's not yet done]
- **File being edited**: [which file]

## Step 3: Resume Work

Continue from the stopped point **without starting over**:

- If editing a file → read that file, then continue editing
- If running a command → run it again
- If in the middle of a workflow → continue from the step that was stopped
- If there are pending todos → continue from the next todo

## Step 4: Update Status

After completing each resumed portion:
- Update todo list (mark completed)
- If work is still substantial → write `progress.txt` summarizing current state (in case credit runs out again)

---

## Critical Rules

1. **Never start over**: Read state first, then continue from the stop point
2. **Never repeat work**: If git diff shows a change was already made → skip it
3. **Write progress.txt**: If work is not done and may be long → summarize for the next continue cycle
4. **Be concise**: Don't re-explain, just continue
