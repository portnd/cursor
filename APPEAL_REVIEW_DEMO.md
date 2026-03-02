# 🎖️ Appeal Review System - Live Demo

## 📋 Test Data

### Task Information
- **Task ID:** `a517e15d-f9aa-4a19-931b-ecf52d967ebf`
- **Task Title:** "Implement secure database query"
- **Task URL:** http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf

### Appeal Information
- **Appeal ID:** `55cd64c4-2168-4eec-9422-7c70af32161b`
- **Submission ID:** `06f92f92-8b5a-4c4c-ae62-ddc9385ae661`
- **Developer:** Dev #1
- **Status:** PENDING
- **Reason:** 
  > "I believe the AI has made an incorrect assessment. The code uses prepared statements which prevent SQL injection. The AI may have misinterpreted the syntax. This is a false positive and should be reviewed by a human expert."

---

## 🧪 Quick Test Steps

### Step 1: View as DEV (No Review Access)
```bash
# 1. Open browser: http://localhost:3000/login
# 2. Login: dev@sentinel.com / password123
# 3. Go to: http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf
# 4. Scroll to Timeline (Right Column)
# 5. Find the FAIL submission (Submission ID: 06f92f92...)

EXPECTED:
✅ Yellow banner: "⚖️ Appeal Under Review"
✅ Banner text: "Your appeal is being reviewed by management. Decision pending."
❌ NO "Review Appeal" button (DEV users cannot review)
```

### Step 2: View as CEO (Full Review Access)
```bash
# 1. Logout (click Logout in sidebar)
# 2. Login: ceo@sentinel.com / password123
# 3. Go to: http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf
# 4. Scroll to Timeline (Right Column)
# 5. Find the FAIL submission with PENDING appeal

EXPECTED:
✅ Yellow banner: "⚖️ Appeal Under Review"
✅ Banner text: "This appeal requires your judgment."
✅ Purple-gold gradient button: "⚖️ Review Appeal"
```

### Step 3: Open Adjudication Modal
```bash
# Click the "⚖️ Review Appeal" button

EXPECTED MODAL CONTENTS:
┌──────────────────────────────────────────────────┐
│ 🎖️ HIGH COURT OF SENTINEL                       │
│ CEO ceo@sentinel.com presiding                   │
├──────────────────────────────────────────────────┤
│ CASE INFORMATION                                 │
│ • Case #: 55cd64c4-2168-4eec-9422-7c70af32161b  │
│ • Submission: 06f92f92...                        │
│ • Appellant: Dev #1 (with avatar)                │
├──────────────────────────────────────────────────┤
│ 📜 THE PLEA                                      │
│ "I believe the AI has made an incorrect          │
│  assessment. The code uses prepared statements   │
│  which prevent SQL injection. The AI may have    │
│  misinterpreted the syntax. This is a false      │
│  positive and should be reviewed by a human      │
│  expert."                                        │
├──────────────────────────────────────────────────┤
│ 📝 Your Verdict Statement (Optional)             │
│ [Textarea - empty]                               │
├──────────────────────────────────────────────────┤
│ [✅ Sustain Appeal]  [❌ Dismiss Appeal]         │
│ (Green button)        (Red button)               │
│                                                   │
│ [Cancel Review]                                  │
└──────────────────────────────────────────────────┘
```

### Step 4: Approve the Appeal
```bash
# 1. (Optional) Enter resolver note:
#    "After reviewing the code, the developer is correct. 
#     The prepared statements are properly implemented. 
#     This is a false positive from the AI."
#
# 2. Click "✅ Sustain Appeal (Set to PASS)"

EXPECTED BEHAVIOR:
1. Button shows spinner (⚙️)
2. Modal closes after ~1-2 seconds
3. Page refreshes automatically
4. Submission card NOW has:
   ✅ GOLD border (instead of red)
   ✅ "👑 OVERRIDDEN BY HUMAN" badge at top
   ✅ Header: "👑 VERDICT OVERRIDDEN" (amber text)
   ✅ Score shows in amber color
5. Task status badge (top-right) changes to:
   ✅ "COMPLETED" (green badge)
```

---

## 🔍 Verification Commands

### Verify Appeal Resolution
```bash
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  a.status as appeal_status,
  a.resolver_id,
  u.email as resolver_email,
  u.role as resolver_role,
  a.resolver_note
FROM appeals a
LEFT JOIN users u ON u.id = a.resolver_id
WHERE a.id = '55cd64c4-2168-4eec-9422-7c70af32161b';
"
```

**Expected Output:**
```
appeal_status | resolver_id | resolver_email   | resolver_role | resolver_note
--------------+-------------+------------------+---------------+----------------------------------
APPROVED      | 1           | ceo@sentinel.com | CEO           | After reviewing the code, the...
```

### Verify Submission Override
```bash
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  id,
  ai_verdict,
  is_overridden,
  created_at
FROM submissions
WHERE id = '06f92f92-8b5a-4c4c-ae62-ddc9385ae661';
"
```

**Expected Output:**
```
id                                   | ai_verdict | is_overridden | created_at
-------------------------------------+------------+---------------+-------------------------
06f92f92-8b5a-4c4c-ae62-ddc9385ae661 | PASS       | t             | 2026-01-25 22:03:11.002
```

### Verify Task Auto-Completion
```bash
docker compose exec -T postgres psql -U komgrip -d komgrip_db -c "
SELECT 
  id,
  title,
  status,
  started_at,
  completed_at
FROM tasks
WHERE id = 'a517e15d-f9aa-4a19-931b-ecf52d967ebf';
"
```

**Expected Output:**
```
id                                   | title                           | status    | started_at          | completed_at
-------------------------------------+---------------------------------+-----------+---------------------+---------------------
a517e15d-f9aa-4a19-931b-ecf52d967ebf | Implement secure database query | COMPLETED | 2026-01-25 23:14:30 | 2026-01-25 23:19:XX
```

---

## 🎨 Visual Verification Checklist

### Before Approval
- [ ] Submission card has RED border (FAIL)
- [ ] "Appeal Under Review" yellow banner visible
- [ ] For CEO: "Review Appeal" button visible (purple-gold gradient)
- [ ] Task status: "IN_PROGRESS" (blue badge)

### After Approval
- [ ] Submission card has GOLD border
- [ ] "👑 OVERRIDDEN BY HUMAN" badge visible (amber, at top of card)
- [ ] Verdict header: "👑 VERDICT OVERRIDDEN" (amber text)
- [ ] Score number in amber color
- [ ] Task status: "COMPLETED" (green badge)
- [ ] NO "Review Appeal" button (appeal is resolved)

### Modal UI (Before Submission)
- [ ] Purple border (4px) around modal
- [ ] Purple-to-amber gradient header
- [ ] Case information card (gray-800 background)
- [ ] Appellant avatar shows (blue gradient circle)
- [ ] "The Plea" in amber blockquote style
- [ ] Resolver note textarea (optional, gray-800)
- [ ] Two large buttons side-by-side:
  - Left: Green gradient "✅ Sustain Appeal"
  - Right: Red gradient "❌ Dismiss Appeal"
- [ ] Cancel button at bottom (gray)

---

## 🚨 Error Scenarios to Test

### Test 1: Try Approving Same Appeal Twice
```bash
# After approving once, try to approve again via API:
CEO_TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ceo@sentinel.com","password":"password123"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['token'])")

curl -s -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "http://localhost:8080/api/v1/sentinel/appeals/55cd64c4-2168-4eec-9422-7c70af32161b/resolve" \
  -d '{"status":"APPROVED","note":"Try again"}' \
  | python3 -m json.tool
```

**Expected Error:**
```json
{
  "error": "Failed to resolve appeal",
  "message": "appeal already resolved with status: APPROVED"
}
```

### Test 2: DEV User Tries to Resolve via API
```bash
DEV_TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@sentinel.com","password":"password123"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['token'])")

curl -s -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "http://localhost:8080/api/v1/sentinel/appeals/55cd64c4-2168-4eec-9422-7c70af32161b/resolve" \
  -d '{"status":"APPROVED","note":"Hacking attempt"}' \
  | python3 -m json.tool
```

**Expected Error:**
```json
{
  "error": "Forbidden",
  "message": "forbidden: only CEO or PM can resolve appeals (current role: DEV)"
}
```

---

## 📸 Screenshot Locations

Take screenshots for documentation:

1. **DEV View - Pending Appeal:**
   - File: `pending_appeal_dev_view.png`
   - URL: http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf
   - User: dev@sentinel.com

2. **CEO View - Review Button:**
   - File: `review_button_ceo_view.png`
   - URL: http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf
   - User: ceo@sentinel.com

3. **Adjudication Modal:**
   - File: `adjudication_modal.png`
   - Action: Click "Review Appeal" button

4. **Overridden Submission:**
   - File: `overridden_submission.png`
   - URL: http://localhost:3000/task/a517e15d-f9aa-4a19-931b-ecf52d967ebf
   - User: ceo@sentinel.com (after approval)

---

## ✅ Success Criteria

All of the following must be true:

- [x] Backend appeal resolution logic works (verified via curl)
- [x] Database correctly updates (appeal, submission, task)
- [x] Task auto-completes on approval
- [ ] DEV users do NOT see "Review Appeal" button
- [ ] CEO/PM users see "Review Appeal" button
- [ ] Modal opens with correct data
- [ ] Approval changes submission to PASS + OVERRIDDEN
- [ ] Rejection keeps submission as FAIL
- [ ] Visual feedback (gold border, badge) appears after approval
- [ ] Page refreshes after resolution
- [ ] No console errors during entire flow

**If all checkboxes are checked: 🎉 FEATURE IS PRODUCTION READY! 🎉**
