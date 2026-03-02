# 🧪 Appeal Review System - Testing Guide

## Test Data Setup

### 1. Check Existing Appeals
```bash
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  a.id as appeal_id,
  a.status as appeal_status,
  s.id as submission_id,
  s.ai_verdict,
  s.is_overridden,
  t.id as task_id,
  t.title as task_title,
  t.status as task_status
FROM appeals a
JOIN submissions s ON s.id = a.submission_id
JOIN tasks t ON t.id = s.task_id
WHERE a.status = 'PENDING'
LIMIT 5;
"
```

### 2. Create Test Appeal (If Needed)
```bash
# First, create a FAIL submission
TASK_ID="<existing_task_id>"
SUBMISSION_ID=$(uuidgen | tr '[:upper:]' '[:lower:]')

docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
INSERT INTO submissions (id, task_id, dev_id, commit_hash, ai_verdict, ai_score, ai_feedback, created_at)
VALUES (
  '$SUBMISSION_ID',
  '$TASK_ID',
  3,
  'test-commit-hash-123',
  'FAIL',
  25,
  '{\"feedback\": \"Test failure for appeal review\"}',
  NOW()
);

-- Create pending appeal
INSERT INTO appeals (id, submission_id, developer_id, reason, status, created_at, updated_at)
VALUES (
  uuid_generate_v4(),
  '$SUBMISSION_ID',
  3,
  'This is a test appeal. The AI misunderstood the context. I believe the code is actually correct.',
  'PENDING',
  NOW(),
  NOW()
)
RETURNING id, submission_id;
"
```

---

## Backend API Testing

### Test 1: Resolve Appeal as CEO (Should Succeed)
```bash
# Login as CEO
CEO_TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ceo@sentinel.com","password":"password123"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['token'])")

# Get a pending appeal ID
APPEAL_ID="<copy from test data setup>"

# Approve Appeal
curl -s -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "http://localhost:8080/api/v1/sentinel/appeals/$APPEAL_ID/resolve" \
  -d '{
    "status": "APPROVED",
    "note": "After careful review, the developer is correct. The AI missed important context."
  }' | python3 -m json.tool

# Expected: {"message": "Appeal resolved successfully"}
```

### Test 2: Verify Database Changes
```bash
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  'APPEAL' as type,
  a.status,
  a.resolver_id,
  a.resolver_note,
  u.email as resolver_email
FROM appeals a
LEFT JOIN users u ON u.id = a.resolver_id
WHERE a.id = '$APPEAL_ID'

UNION ALL

SELECT 
  'SUBMISSION',
  s.ai_verdict,
  NULL,
  CASE WHEN s.is_overridden THEN 'OVERRIDDEN' ELSE 'ORIGINAL' END,
  NULL
FROM submissions s
WHERE s.id = (SELECT submission_id FROM appeals WHERE id = '$APPEAL_ID')

UNION ALL

SELECT 
  'TASK',
  t.status,
  NULL,
  to_char(t.completed_at, 'YYYY-MM-DD HH24:MI:SS'),
  NULL
FROM tasks t
WHERE t.id = (
  SELECT task_id FROM submissions WHERE id = (
    SELECT submission_id FROM appeals WHERE id = '$APPEAL_ID'
  )
);
"

# Expected:
# Appeal: APPROVED, resolver_id=1, resolver_note="...", resolver_email="ceo@sentinel.com"
# Submission: PASS, OVERRIDDEN
# Task: COMPLETED, completed_at=<timestamp>
```

### Test 3: Try to Resolve as DEV (Should Fail with 403)
```bash
# Login as DEV
DEV_TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@sentinel.com","password":"password123"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['token'])")

# Try to resolve appeal
curl -s -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "http://localhost:8080/api/v1/sentinel/appeals/$APPEAL_ID/resolve" \
  -d '{"status":"APPROVED","note":"Trying to approve"}' \
  | python3 -m json.tool

# Expected: 
# {
#   "error": "Forbidden",
#   "message": "forbidden: only CEO or PM can resolve appeals (current role: DEV)"
# }
```

---

## Frontend UI Testing

### Test 1: Login as CEO and View Task
1. Open browser: `http://localhost:3000/login`
2. Login with:
   - Email: `ceo@sentinel.com`
   - Password: `password123`
3. Navigate to Dashboard
4. **Verify:** User info shows "CEO" role (top-left sidebar)
5. Click on a task that has a PENDING appeal

### Test 2: Verify UI Elements for CEO
On the Task Detail page:

#### **Submission Card (Timeline - Right Column)**
- [ ] Card has yellow border (PENDING appeal)
- [ ] "⚖️ Appeal Under Review" banner visible
- [ ] Banner text: "This appeal requires your judgment."
- [ ] **"⚖️ Review Appeal"** button visible (purple-gold gradient)
- [ ] Button is full-width below the banner

#### **Click "Review Appeal" Button**
- [ ] Modal opens with dark background overlay
- [ ] Modal has 4-pixel purple border
- [ ] Header shows:
  - "⚖️ High Court of Sentinel" title
  - CEO email and role ("CEO ceo@sentinel.com presiding")
- [ ] Case Information card shows:
  - Case # (Appeal ID)
  - Submission ID (truncated)
  - Appellant name with avatar (Dev #3)
- [ ] "The Plea" section shows:
  - Amber-themed blockquote
  - Developer's appeal reason in quotes
- [ ] Resolver Note textarea visible (optional)
- [ ] Two verdict buttons side-by-side:
  - Left: Green "✅ Sustain Appeal (Set to PASS)"
  - Right: Red "❌ Dismiss Appeal (Keep as FAIL)"
- [ ] Cancel button at bottom

### Test 3: Approve Appeal (Full Flow)
1. In the modal, enter resolver note: "Developer is correct. AI false positive."
2. Click "✅ Sustain Appeal"
3. **Verify:**
   - [ ] Button shows spinner (⚙️ rotating)
   - [ ] Modal closes after 1-2 seconds
   - [ ] Page refreshes automatically
4. **After Refresh:**
   - [ ] Submission card border is NOW GOLD (not red)
   - [ ] "👑 OVERRIDDEN BY HUMAN" badge appears at top of card
   - [ ] Verdict header shows "👑 VERDICT OVERRIDDEN" (amber text)
   - [ ] Score number is amber-colored
   - [ ] Task status badge (top-right) shows "COMPLETED" (green)

### Test 4: View as DEV User
1. Logout
2. Login as `dev@sentinel.com`
3. Navigate to the SAME task
4. **Verify:**
   - [ ] Submission card shows GOLD border (overridden)
   - [ ] "👑 OVERRIDDEN BY HUMAN" badge visible
   - [ ] NO "Review Appeal" button (DEV users can't review)
   - [ ] Appeal status shows as resolved (not pending)

### Test 5: Reject Appeal (Create Another Test Case)
1. Login as CEO
2. Navigate to a different task with PENDING appeal
3. Click "⚖️ Review Appeal"
4. Enter resolver note: "AI was correct. Security vulnerability is real."
5. Click "❌ Dismiss Appeal"
6. **Verify:**
   - [ ] Modal closes
   - [ ] Page refreshes
   - [ ] Submission card remains RED (FAIL verdict)
   - [ ] NO "OVERRIDDEN" badge
   - [ ] Appeal status shows "REJECTED"
   - [ ] If you expand appeal info, resolver note is visible

---

## Visual Regression Checklist

### Submission Card States
1. **FAIL + No Appeal:** Red border, "Appeal Verdict" button visible (for dev)
2. **FAIL + PENDING Appeal (DEV view):** Yellow border, "Appeal Under Review" banner
3. **FAIL + PENDING Appeal (CEO/PM view):** Yellow border, banner + "Review Appeal" button
4. **FAIL + REJECTED Appeal:** Red border, rejection banner with resolver note
5. **PASS + APPROVED Appeal:** GOLD border, "OVERRIDDEN BY HUMAN" badge, amber text
6. **PASS + No Appeal:** Green border, "MISSION ACCOMPLISHED"

### Modal States
1. **Initial Open:** All fields empty, buttons enabled
2. **During Submission:** Buttons show spinner, disabled
3. **Error State:** Red error banner appears above form

---

## Performance Testing

### Loading States
- [ ] Modal opens instantly (<100ms)
- [ ] API call completes in <2 seconds (typical)
- [ ] Page refresh completes in <3 seconds
- [ ] No console errors during flow

### Error Handling
Test error scenarios:
1. **Network Error:** Disconnect internet, try to resolve
   - [ ] Error banner shows "Failed to resolve appeal"
2. **Invalid Appeal ID:** Use fake UUID
   - [ ] Error shows "Appeal not found"
3. **Already Resolved:** Try to resolve same appeal twice
   - [ ] Error shows "Appeal already resolved"

---

## Accessibility Testing

- [ ] Modal can be closed with ESC key
- [ ] Modal can be closed by clicking overlay
- [ ] Tab order is logical (Resolver Note → Approve → Reject → Cancel)
- [ ] Buttons have clear labels (not just icons)
- [ ] Error messages are clear and actionable

---

## Cross-Browser Testing

Test in:
- [ ] Chrome/Edge (Chromium)
- [ ] Firefox
- [ ] Safari

Verify:
- [ ] Gradient backgrounds render correctly
- [ ] Modal overlay is dark and blocks content
- [ ] Buttons are clickable and responsive
- [ ] Text is readable (contrast)

---

## Final Verification

After all tests pass:

```bash
# Check API logs for errors
docker compose logs api --tail 50 | grep -i "error\|fail\|panic"

# Check web logs for console errors
# (Use browser DevTools Console tab)

# Verify database consistency
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  COUNT(*) FILTER (WHERE a.status = 'PENDING') as pending_appeals,
  COUNT(*) FILTER (WHERE a.status = 'APPROVED') as approved_appeals,
  COUNT(*) FILTER (WHERE a.status = 'REJECTED') as rejected_appeals,
  COUNT(*) FILTER (WHERE s.is_overridden = true) as overridden_submissions,
  COUNT(*) FILTER (WHERE t.status = 'COMPLETED' AND s.ai_verdict = 'PASS' AND s.is_overridden = true) as auto_completed_tasks
FROM appeals a
LEFT JOIN submissions s ON s.id = a.submission_id
LEFT JOIN tasks t ON t.id = s.task_id;
"
```

Expected output:
- `pending_appeals`: 0 (all resolved in test)
- `approved_appeals`: >= 1
- `rejected_appeals`: >= 1
- `overridden_submissions`: Same as `approved_appeals`
- `auto_completed_tasks`: Same as `approved_appeals`

---

**🎯 All Tests Passed? The Appeal Review System is Production Ready!**
