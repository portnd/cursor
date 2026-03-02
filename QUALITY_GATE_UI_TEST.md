# 🧪 Quality Gate UI - Testing Guide

## ⚡ Quick Test Checklist

```
┌─────────────────────────────────────────────────────┐
│  QUALITY GATE UI TESTING                            │
├─────────────────────────────────────────────────────┤
│  Status: Ready for Manual Testing                   │
│  Files Modified: 3                                  │
│  Components Added: 6                                │
│  No Linter Errors: ✅                               │
└─────────────────────────────────────────────────────┘
```

---

## 🎯 Test Scenarios

### **Scenario 1: Developer Submits Code (AI PASS)**

**Objective:** Verify developer sees correct UI when task is pending approval

**Steps:**
```
1. Login as DEV
   URL: http://localhost:3000/login

2. Go to assigned task
   URL: http://localhost:3000/task/{task_id}

3. Submit work with secure code
   POST /sentinel/tasks/{id}/submit
   {
     "commit_hash": "test123",
     "diff": "db.Where(\"email = ?\", email).First(&user)"
   }

4. AI reviews: PASS
   → Task status: REVIEW_PENDING
```

**Expected UI:**

✅ **Status Badge:**
```
[⏳ WAITING FOR APPROVAL]
• Color: Purple/Indigo
• Border: Indigo-600
• Animation: Pulsing
```

✅ **Banner (Top of page):**
```
┌────────────────────────────────────────────────┐
│ 🎉  AI Security Checks Passed!      90/100    │
│                                                │
│ Your code passed all AI security audits.      │
│ Awaiting PM/CEO verification for              │
│ functionality review.                          │
└────────────────────────────────────────────────┘
• Background: Indigo gradient
• Border: 2px indigo
• Icon: 🎉 (pulsing)
• Score: Visible on right
```

✅ **Button Visibility:**
```
❌ "Approve & Complete" button: NOT VISIBLE
   (Developer cannot self-approve)
```

**Pass Criteria:**
- [ ] Badge is purple and pulsing
- [ ] Banner shows celebration message
- [ ] AI score is displayed (90/100)
- [ ] Approve button is NOT visible
- [ ] Message mentions "PM/CEO verification"

---

### **Scenario 2: CEO Views Dashboard**

**Objective:** Verify CEO sees quality gate section with pending tasks

**Steps:**
```
1. Login as CEO
   URL: http://localhost:3000/login

2. Navigate to dashboard
   URL: http://localhost:3000/dashboard
```

**Expected UI:**

✅ **Metric Card (Top Right):**
```
┌───────────────────────┐
│ 🚦 READY FOR REVIEW   │
│        3              │ ← Pulsing number
│ Click to review       │
└───────────────────────┘
• Background: Indigo gradient
• Border: 2px indigo
• Animation: Pulse
• Cursor: Pointer
• Hover: Scale 1.05x
```

✅ **Quality Gate Section:**
```
┌─────────────────────────────────────────────────────────┐
│ 🚦 Quality Gate: Ready for Approval                    │
│ 3 tasks passed AI review, awaiting your verification   │
├─────────────────────────────────────────────────────────┤
│ Task          │ Dev  │ Score │ Time │ Action           │
├─────────────────────────────────────────────────────────┤
│ Fix Login Bug │ Dev2 │✅PASS│ 2h   │[🔍 Review&Approve]│
│ Add API Docs  │ Dev1 │✅PASS│ 1h   │[🔍 Review&Approve]│
│ Update Tests  │ Dev3 │✅PASS│ 30m  │[🔍 Review&Approve]│
└─────────────────────────────────────────────────────────┘
• Location: High priority (top, after metrics)
• Background: Indigo/Purple gradient
• Border: 2px indigo
• Table: Full width, scrollable
```

✅ **Interactions:**
```
1. Click metric card → Auto-scroll to section ✅
2. Hover table rows → Background highlight ✅
3. Click "Review & Approve" → Navigate to task ✅
```

**Pass Criteria:**
- [ ] Metric card visible and pulsing
- [ ] Metric shows correct count (3)
- [ ] Click metric → Scrolls to section
- [ ] Quality Gate section at top of page
- [ ] Table shows all REVIEW_PENDING tasks
- [ ] Each row has "Review & Approve" button
- [ ] Buttons are green gradient with glow

---

### **Scenario 3: CEO Approves Task**

**Objective:** Verify approval flow works end-to-end

**Steps:**
```
1. From dashboard, click "🔍 Review & Approve"
   → Should navigate to task detail page

2. Verify CEO sees:
   • Status badge: Purple, pulsing
   • Banner: "Quality Gate: Awaiting Your Approval"
   • Button: "✅ Approve & Complete" (green)

3. Click "✅ Approve & Complete"

4. Verify loading state:
   • Button text: "Approving..."
   • Icon: ⚙️ (spinning)
   • Button: Disabled (gray)

5. Wait for API response

6. Verify success:
   • Alert: "🎉 Task Approved & Completed!"
   • Status badge: Changes to "✅ COMPLETED" (green)
   • Button: Disappears
   • Banner: Disappears
   • Page: Refreshes with new data
```

**Expected Behavior:**

✅ **Before Approval:**
```
Status: [⏳ WAITING FOR APPROVAL] (purple, pulsing)
Button: [✅ Approve & Complete] (green, visible)
Banner: "Quality Gate: Awaiting Your Approval"
```

✅ **During Approval:**
```
Button: [⚙️ Approving...] (gray, disabled)
Status: (unchanged)
```

✅ **After Approval:**
```
Alert: "🎉 Task Approved & Completed!"
Status: [✅ COMPLETED] (green, no pulse)
Button: (disappeared)
Banner: (disappeared)
Timestamp: completed_at now set
```

**Pass Criteria:**
- [ ] Navigation works from dashboard
- [ ] Banner visible for CEO/PM
- [ ] Approve button visible and enabled
- [ ] Click triggers loading state
- [ ] Success alert shown
- [ ] Status updates to COMPLETED
- [ ] Button and banner disappear
- [ ] Page refreshes with new data

---

### **Scenario 4: PM Views & Approves**

**Objective:** Verify PM has same capabilities as CEO

**Steps:**
```
1. Login as PM
2. View dashboard
3. Verify Quality Gate section visible
4. Click "Review & Approve"
5. Navigate to task detail
6. Click "✅ Approve & Complete"
7. Verify approval succeeds
```

**Expected:**
- ✅ PM sees all same UI elements as CEO
- ✅ PM can approve tasks
- ✅ No functional differences

**Pass Criteria:**
- [ ] PM dashboard shows Quality Gate section
- [ ] PM can navigate to tasks
- [ ] PM can see approve button
- [ ] PM can successfully approve tasks

---

### **Scenario 5: Developer Tries to Approve (Should Fail)**

**Objective:** Verify developers cannot approve

**Steps:**
```
1. Login as DEV
2. Navigate to REVIEW_PENDING task
3. Verify approve button is NOT visible
4. (Optional) Attempt direct API call
```

**Expected:**
- ❌ Approve button: NOT VISIBLE
- ❌ No approval UI elements for developers
- ❌ API call (if attempted): 403 Forbidden

**Pass Criteria:**
- [ ] Developer sees waiting message
- [ ] No approve button visible
- [ ] Cannot interact with approval system

---

### **Scenario 6: Multiple Pending Tasks**

**Objective:** Verify UI scales with multiple tasks

**Steps:**
```
1. Create 5 tasks with AI PASS (REVIEW_PENDING)
2. Login as CEO
3. View dashboard
```

**Expected:**
- ✅ Metric: "🚦 READY FOR REVIEW: 5"
- ✅ Table: Shows all 5 tasks
- ✅ Sorted: Most recent first
- ✅ No performance issues
- ✅ Scrollable if needed

**Pass Criteria:**
- [ ] Metric count is correct (5)
- [ ] Table shows all tasks
- [ ] Sorting correct (newest first)
- [ ] UI remains responsive
- [ ] No layout issues

---

### **Scenario 7: Status Transitions**

**Objective:** Verify status badge updates correctly

**Test Status Flow:**
```
PENDING → ASSIGNED → UNDER_REVIEW → REVIEW_PENDING → COMPLETED
```

**Verify Each Status:**
```
PENDING:
  • Color: Yellow
  • Label: "⏳ PENDING"
  • Animation: None

REVIEW_PENDING:
  • Color: Purple/Indigo
  • Label: "⏳ WAITING FOR APPROVAL"
  • Animation: Pulse

COMPLETED:
  • Color: Green
  • Label: "✅ COMPLETED"
  • Animation: None
```

**Pass Criteria:**
- [ ] Each status has correct color
- [ ] Labels are clear and actionable
- [ ] REVIEW_PENDING pulses
- [ ] Others do not pulse

---

### **Scenario 8: Error Handling**

**Objective:** Verify errors are handled gracefully

**Test Cases:**

**A. API Error (Network Failure):**
```
1. Disconnect from API
2. Try to approve task
3. Expected: Alert with error message
4. Expected: Button re-enables
```

**B. Permission Error (403):**
```
1. Login as DEV
2. Manually call approve API (dev tools)
3. Expected: 403 Forbidden
4. Expected: Error message displayed
```

**C. Invalid Status (400):**
```
1. Try to approve already COMPLETED task
2. Expected: 400 Bad Request
3. Expected: "Task is not pending review" message
```

**Pass Criteria:**
- [ ] Network errors: User-friendly message
- [ ] Permission errors: Clear denial message
- [ ] Status errors: Explains current state
- [ ] Buttons re-enable after error
- [ ] No app crashes

---

## 📊 Visual Verification

### **Color Accuracy:**

**REVIEW_PENDING Badge:**
```
Expected: bg-indigo-900 text-indigo-200 border-indigo-600
Verify: Opens Chrome DevTools → Inspect element
Check: Computed styles match expected values
```

**Approve Button:**
```
Expected: bg-gradient-to-r from-green-600 to-emerald-600
Verify: Hover and check shadow-green-500/50
Check: Glow effect visible on hover
```

---

### **Animation Verification:**

**Pulse Effect:**
```
Element: Status badge, metric card
Expected: Opacity fades in/out smoothly
Duration: ~1.5 seconds per cycle
Verify: Watch for 5 seconds, count cycles
```

**Spin Effect:**
```
Element: ⚙️ icon during loading
Expected: 360° rotation continuously
Duration: 1 second per rotation
Verify: Smooth rotation, no jank
```

---

## 🔍 Browser Testing

### **Chrome (Primary):**
- [ ] All features work
- [ ] Animations smooth
- [ ] Colors accurate
- [ ] No console errors

### **Firefox:**
- [ ] Gradient rendering correct
- [ ] Pulse animation works
- [ ] Table layout correct

### **Safari:**
- [ ] Border styles render
- [ ] Backdrop blur works
- [ ] Animation performance good

### **Mobile (Chrome/Safari):**
- [ ] Responsive layout
- [ ] Touch interactions
- [ ] Buttons sized appropriately

---

## 📝 Test Reporting

### **Test Results Template:**

```markdown
## Test Session

**Date:** 2026-01-26
**Tester:** [Name]
**Browser:** Chrome 121
**Viewport:** 1920x1080

### Scenario 1: Developer UI
- [x] Status badge: Purple, pulsing ✅
- [x] Banner: Celebration message ✅
- [x] Score: Displayed (90/100) ✅
- [x] Button: Not visible ✅

### Scenario 2: CEO Dashboard
- [x] Metric: Visible, pulsing ✅
- [x] Section: High priority ✅
- [x] Table: All tasks shown ✅
- [x] Buttons: Green, glowing ✅

### Scenario 3: Approval Flow
- [x] Navigation: Works ✅
- [x] Button: Visible ✅
- [x] Loading: Smooth ✅
- [x] Success: Alert shown ✅
- [x] Status: Updates ✅

### Issues Found:
- None

### Overall Result: ✅ PASS
```

---

## 🚀 Pre-Deployment Checklist

### **Code Quality:**
- [x] No linter errors
- [x] No TypeScript errors
- [x] All imports present
- [x] No unused variables
- [x] Consistent naming

### **Functionality:**
- [ ] Approve button works (CEO/PM)
- [ ] Approve button hidden (DEV)
- [ ] Status updates correctly
- [ ] Banners show/hide correctly
- [ ] Dashboard sections visible
- [ ] Navigation works

### **Visual Design:**
- [ ] Colors match spec (purple/indigo)
- [ ] Animations smooth (pulse, spin)
- [ ] Gradients render correctly
- [ ] Responsive design works
- [ ] Icons display properly

### **Performance:**
- [ ] No lag on status updates
- [ ] Smooth animations (60fps)
- [ ] Fast API responses
- [ ] No memory leaks

### **Accessibility:**
- [ ] Buttons keyboard accessible
- [ ] Tab order logical
- [ ] Tooltips present
- [ ] Contrast sufficient

---

## 🔧 Testing Tools

### **1. Browser DevTools**
```
• Chrome DevTools (F12)
• Inspect elements
• Check computed styles
• Monitor network requests
• View console errors
```

### **2. Vue DevTools**
```
• Install Vue DevTools extension
• View component state
• Check reactive updates
• Debug computed properties
```

### **3. Network Tab**
```
• Monitor API calls
• Check request/response
• Verify status codes (200, 403, 400)
• Check timing
```

### **4. Console**
```
• Watch for errors
• Check API responses
• View debug logs
• Monitor warnings
```

---

## 📊 Expected API Calls

### **Page Load:**
```
GET /api/v1/sentinel/tasks/{id}
→ Response: Task data with status
```

### **Approve Click:**
```
POST /api/v1/sentinel/tasks/{id}/approve
→ Request: (empty body)
→ Response (200): { "message": "Task approved..." }
→ Then: GET /api/v1/sentinel/tasks/{id} (refresh)
```

### **Dashboard Load:**
```
GET /api/v1/sentinel/tasks
→ Response: Array of all tasks
→ Frontend filters: REVIEW_PENDING tasks
```

---

## 🎯 Success Criteria

### **Must Pass:**
- ✅ Status badge: Purple, pulsing, correct label
- ✅ Approve button: Visible for CEO/PM only
- ✅ Banners: Role-specific, correct content
- ✅ Dashboard: Quality Gate section at top
- ✅ Approval: Works and updates status
- ✅ Errors: Handled gracefully
- ✅ No console errors
- ✅ No visual glitches

### **Should Pass:**
- ✅ Animations smooth (60fps)
- ✅ Responsive design works
- ✅ Time ago formatting accurate
- ✅ Auto-scroll to approvals works

### **Nice to Have:**
- ✅ Confetti animation on approval
- ✅ Toast notifications (instead of alerts)
- ✅ Batch approval feature
- ✅ Keyboard shortcuts

---

## 🐛 Known Issues / Limitations

### **Current Implementation:**

1. **Alerts instead of Toasts**
   - Using native `alert()` for notifications
   - Consider: Vue Toastification library
   - Impact: Low (functional but not pretty)

2. **No Confetti Animation**
   - Could add canvas-confetti on approval
   - Impact: Low (nice-to-have)

3. **Manual Refresh**
   - Dashboard doesn't auto-refresh
   - User must click "Refresh" button
   - Consider: Auto-refresh every 30s
   - Impact: Medium (UX improvement)

4. **No Batch Operations**
   - Can only approve one task at a time
   - Consider: "Approve All" button
   - Impact: Low (rare use case)

---

## ✅ Test Sign-Off Template

```markdown
## Quality Gate UI - Test Sign-Off

**Tester:** [Name]
**Date:** [Date]
**Environment:** [Local/Staging/Production]

### Test Results:
- Scenario 1 (Dev UI): ✅ PASS / ❌ FAIL
- Scenario 2 (CEO Dashboard): ✅ PASS / ❌ FAIL
- Scenario 3 (Approval Flow): ✅ PASS / ❌ FAIL
- Scenario 4 (PM View): ✅ PASS / ❌ FAIL
- Scenario 5 (Dev Cannot Approve): ✅ PASS / ❌ FAIL
- Scenario 6 (Multiple Tasks): ✅ PASS / ❌ FAIL
- Scenario 7 (Status Transitions): ✅ PASS / ❌ FAIL
- Scenario 8 (Error Handling): ✅ PASS / ❌ FAIL

### Issues Found:
[List any bugs or issues]

### Overall Result:
✅ APPROVED FOR DEPLOYMENT
❌ NEEDS FIXES

### Signature:
[Name] - [Date]
```

---

## 🚀 Quick Start Testing

### **Fastest Way to Test:**

```bash
# 1. Ensure services running
docker-compose ps

# 2. Open browser
# 3. Login as CEO: http://localhost:3000/login
# 4. Go to dashboard: http://localhost:3000/dashboard
# 5. Check for "🚦 READY FOR REVIEW" metric
# 6. If count > 0, click it and test approval flow
# 7. If count = 0, create test data:
#    • Create task
#    • Assign to dev
#    • Submit with AI PASS
#    • Return to dashboard
```

---

## 📚 Documentation

- **Complete Guide:** `QUALITY_GATE_UI_COMPLETE.md`
- **Visual Guide:** `UI_VISUAL_GUIDE.md`
- **This Test Guide:** `QUALITY_GATE_UI_TEST.md`

---

**Status: ✅ READY FOR TESTING**

**Start testing the Quality Gate UI! 🧪**
