# 🎖️ Appeal Review System Implementation

## Overview
CEO/PM can now review and resolve developer appeals directly from the Task Detail page.

---

## ✨ Features Implemented

### 1. **Role-Based Access Control**
- Only users with role `'CEO'` or `'PM'` can see the "Review Appeal" button
- Regular developers see the standard "Appeal Under Review" banner

### 2. **Review Appeal Button**
- **Location:** Timeline (Right Column) → Inside Submission Card
- **Condition:** Visible when:
  - User is CEO or PM (`isCeoOrPm`)
  - Submission has an appeal with status `'PENDING'`
- **Style:** Purple-to-Gold gradient button
- **Action:** Opens "Adjudication Modal"

### 3. **The Adjudication Modal**
High-contrast, dramatic UI styled as a "High Court":

#### **Header**
- Title: "⚖️ High Court of Sentinel"
- Subtitle: Shows presiding officer (CEO/PM email and role)
- Purple/Gold gradient background

#### **Case Information**
- **Case #:** Appeal ID
- **Submission:** Submission ID (truncated)
- **Appellant:** Developer avatar + name

#### **The Plea**
- Displays the developer's `appeal.reason` in a styled blockquote
- Amber-themed background to emphasize importance

#### **Resolver Note (Optional)**
- Textarea for CEO/PM to explain their decision
- Placeholder suggests both approval and rejection examples

#### **Verdict Actions**
Two large gradient buttons side-by-side:

1. **✅ Sustain Appeal (APPROVE)**
   - Green gradient
   - Sets submission verdict to `PASS`
   - Auto-completes task (backend logic)
   - Overrides AI verdict (`is_overridden = true`)

2. **❌ Dismiss Appeal (REJECT)**
   - Red gradient
   - Keeps submission as `FAIL`
   - Records resolver's decision in appeal

---

## 🔌 Backend Integration

### API Endpoint
```
POST /api/v1/sentinel/appeals/:id/resolve
```

### Request Body
```json
{
  "status": "APPROVED" | "REJECTED",
  "note": "Optional explanation for the decision"
}
```

### Backend Logic (When APPROVED)
1. Update appeal status to `'APPROVED'`
2. Set `resolver_id` and `resolver_note`
3. Override submission: `ai_verdict = 'PASS'`, `is_overridden = true`
4. Auto-complete task: `status = 'COMPLETED'`, set `completed_at`
5. Calculate and log actual vs estimated time

---

## 🎨 UI/UX Design Principles

### Visual Hierarchy
- **Dramatic & Authoritative:** Purple/gold color scheme conveys gravitas
- **High Contrast:** Clear distinction between Approve (green) and Reject (red)
- **Responsive Feedback:** Loading spinners during API calls
- **Error Handling:** Red error banner for API failures

### User Flow
```
CEO/PM views Task Detail Page
  → Sees submission with PENDING appeal
  → Clicks "⚖️ Review Appeal" button
  → Modal opens with full case details
  → Reads developer's plea
  → (Optional) Writes resolver note
  → Clicks "Sustain" or "Dismiss"
  → API call + loading state
  → Modal closes
  → Page refreshes automatically
  → Submission card updates (Green if approved, shows OVERRIDDEN badge)
```

---

## 📦 Files Modified

### Frontend
1. **`web/pages/task/[id].vue`**
   - Added `authStore` import
   - Added `isCeoOrPm` computed property
   - Added `canReviewAppeal()` method
   - Added adjudication modal state management
   - Added `openAdjudicationModal()`, `closeAdjudicationModal()`, `resolveAppeal()` methods
   - Added "Review Appeal" button in submission card
   - Added complete Adjudication Modal UI

2. **`web/core/modules/auth/store/auth-store.ts`**
   - Updated `User` interface to include `role: string`

3. **`web/core/modules/auth/infrastructure/auth-api.ts`**
   - Updated `User` interface to include `role: string`

---

## 🧪 Testing Instructions

### Prerequisites
```bash
# Ensure API and Web services are running
docker compose up -d

# Ensure you have test users with roles
# - ceo@sentinel.com (CEO)
# - pm@sentinel.com (PM)
# - dev@sentinel.com (DEV)

# Ensure you have a task with a FAIL submission and PENDING appeal
```

### Test Case 1: DEV User (Should NOT see Review button)
1. Login as `dev@sentinel.com`
2. Navigate to a task with a PENDING appeal
3. **Expected:** See "Appeal Under Review" banner ONLY
4. **Expected:** NO "Review Appeal" button visible

### Test Case 2: CEO User (Should see Review button)
1. Login as `ceo@sentinel.com`
2. Navigate to a task with a PENDING appeal
3. **Expected:** See "Appeal Under Review" banner with note "This appeal requires your judgment."
4. **Expected:** See purple-gold gradient "⚖️ Review Appeal" button
5. Click the button
6. **Expected:** Adjudication Modal opens with:
   - "High Court of Sentinel" header
   - CEO email shown as presiding officer
   - Case ID and Submission ID
   - Appellant name (Dev #X)
   - Developer's appeal reason in blockquote
   - Resolver note textarea
   - Two verdict buttons (Sustain/Dismiss)

### Test Case 3: Approve Appeal
1. In Adjudication Modal, optionally enter a resolver note
2. Click "✅ Sustain Appeal"
3. **Expected:** Loading spinner appears
4. **Expected:** Modal closes after success
5. **Expected:** Page refreshes
6. **Expected:** Submission card border changes to GOLD
7. **Expected:** "👑 OVERRIDDEN BY HUMAN" badge appears
8. **Expected:** Submission verdict shows as PASS
9. **Backend:** Task status should be "COMPLETED"

### Test Case 4: Reject Appeal
1. Open Adjudication Modal for a different PENDING appeal
2. Enter resolver note: "AI was correct."
3. Click "❌ Dismiss Appeal"
4. **Expected:** Loading spinner appears
5. **Expected:** Modal closes
6. **Expected:** Page refreshes
7. **Expected:** Submission card remains RED (FAIL)
8. **Expected:** Appeal status shows "REJECTED" with resolver note

### Test Case 5: PM User (Should have same access as CEO)
1. Login as `pm@sentinel.com`
2. Repeat Test Cases 2-4
3. **Expected:** Same behavior as CEO

---

## 🔒 Security Notes

- **Backend Validation:** Role validation happens in `usecase` layer
- **Frontend Guard:** `canReviewAppeal()` checks user role client-side (UX)
- **API Error Handling:** 403 Forbidden if non-CEO/PM tries to resolve
- **Token-Based Auth:** All requests include JWT with role claim

---

## 🎯 Key Design Decisions

### 1. Why Optional Resolver Note?
- Allows quick approval/rejection without mandatory explanation
- Provides flexibility for obvious cases
- Auto-generates basic note if left empty

### 2. Why Two Separate Buttons?
- Clear distinction between actions (no dropdown confusion)
- Reduces cognitive load
- Visually dramatic (matches "verdict" theme)

### 3. Why Auto-Refresh After Resolution?
- Ensures UI shows latest state
- Simpler than optimistic updates for complex nested data
- Provides confirmation that action succeeded

---

## 🚀 Production Readiness

- ✅ Role-based access control
- ✅ Comprehensive error handling
- ✅ Loading states for all async operations
- ✅ Auto-refresh after successful resolution
- ✅ Backend transaction safety
- ✅ Audit trail (resolver_id, resolver_note saved)
- ✅ Visual feedback for overridden verdicts
- ✅ TypeScript type safety

---

**STATUS: 🎖️ PRODUCTION READY!**

The Appeal Review System provides CEO/PM with a powerful, intuitive interface to adjudicate developer appeals with the gravitas and clarity befitting a "High Court of Sentinel" experience.
