# 🧪 Task Edit & Delete - Browser Testing Guide

## 🎯 Testing Objective
Verify that Edit and Delete functionality works correctly with proper permission controls.

---

## ⏱️ Estimated Time: 10 Minutes

---

## 🔧 Prerequisites

### **1. Services Running**
```bash
cd /Users/tnd/Documents/the-sentinel-core
docker compose up -d

# Verify services
docker compose ps
# ✅ api: Up, port 8080
# ✅ web: Up, port 3000
# ✅ postgres: Up
```

### **2. Test Accounts**
| Email | Password | Role | User ID |
|-------|----------|------|---------|
| `ceo@sentinel.com` | `password123` | CEO | 1 |
| `pm@sentinel.com` | `password123` | PM | 2 |
| `dev@sentinel.com` | `password123` | DEV | 3 |

### **3. Test Data**
- At least 1 task created by CEO (created_by = 1)
- At least 1 task created by DEV (created_by = 3)

---

## 🧪 Test Suite

### **Test 1: CEO Can Edit Any Task ✅**

**Objective:** Verify CEO can edit any task and AI re-estimates

**Steps:**
1. Open: `http://localhost:3000/login`
2. Login as **CEO**: `ceo@sentinel.com` / `password123`
3. Click "Dashboard" or "Tasks Inbox"
4. Click on any task to open detail page
5. **Verify:** See "✏️ Edit" and "🗑️ Abort" buttons in header
6. Click **"✏️ Edit"** button
7. **Verify Edit Modal Opens:**
   - Title field is pre-filled
   - Description is pre-filled
   - Deadline is pre-filled (if exists)
   - Yellow warning banner shows: "AI Re-estimation Alert"
8. Change **Title** to: `"Test: Updated Secure Database System"`
9. Change **Description** to: `"Add advanced security with encryption, input validation, and audit logging"`
10. Click **"💾 Update Mission"**
11. **Verify:**
    - Loading spinner appears
    - Alert shows: "✅ Mission Updated! 🤖 AI is recalculating time estimate..."
    - Modal closes
    - Page refreshes
12. Check API logs:
    ```bash
    docker compose logs api --tail 20 | grep -E "AI|Task|Updated"
    ```
13. **Expected Log Output:**
    ```
    🔄 Task content changed. Triggering AI re-estimation...
       Old: [Original Title] Original description
       New: [Test: Updated...] Add advanced security...
    🧠 AI Estimation Request: Test: Updated...
    📡 Calling Gemini API (model: gemini-2.5-flash)
    ✅ AI Re-estimation Complete: 240 minutes (4.0 hours)
    ✅ Task Updated: uuid by CEO (User ID: 1)
    ```
14. **Verify on Page:**
    - Task title updated to new value
    - AI estimated time changed (e.g., 30min → 240min)
    - If negotiation was PENDING, now shows "NONE"

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 2: Creator Can Edit Own Task ✅**

**Objective:** Verify task creator can edit their own task

**Steps:**
1. Login as **DEV**: `dev@sentinel.com` / `password123`
2. Go to "Create" page
3. Create a new task:
   - Title: "Test Task for Edit"
   - Description: "Original description"
   - Deadline: Tomorrow
4. Click "Initialize Task"
5. After creation, click on the task to open detail page
6. **Verify:** See "✏️ Edit" and "🗑️ Abort" buttons (creator = you)
7. Click **"✏️ Edit"**
8. Change **Description** to: `"Updated description by creator"`
9. Click **"Update Mission"**
10. **Verify:**
    - Update succeeds
    - Page refreshes
    - Description updated

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 3: Non-Creator Cannot Edit Other's Task ❌**

**Objective:** Verify permission enforcement (non-creator cannot edit)

**Steps:**
1. Login as **DEV**: `dev@sentinel.com` / `password123`
2. Navigate to a task created by CEO (created_by = 1, not DEV)
   - Get CEO's task ID from dashboard or API
3. Open: `http://localhost:3000/task/{ceo-task-id}`
4. **Verify:**
   - "✏️ Edit" button **NOT VISIBLE**
   - "🗑️ Abort" button **NOT VISIBLE**
   - Only "Back" button shows
5. Try direct API call (to verify backend enforcement):
   ```bash
   DEV_TOKEN="<get from browser DevTools>"
   TASK_ID="<ceo-task-id>"
   
   curl -X PATCH "http://localhost:8080/api/v1/sentinel/tasks/$TASK_ID" \
     -H "Authorization: Bearer $DEV_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"title": "Hacked"}' \
     -s
   ```
6. **Expected API Response:**
   ```json
   {
     "error": "Forbidden",
     "message": "unauthorized: only the task creator or CEO can update this task"
   }
   ```

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 4: Edit with No Changes Shows Error ❌**

**Objective:** Verify validation (no empty updates)

**Steps:**
1. Login as **CEO**
2. Open any task detail page
3. Click **"✏️ Edit"**
4. **Don't change any field** (leave as is)
5. Click **"Update Mission"**
6. **Verify:**
   - Error message shows: "No changes detected. Please modify at least one field."
   - Modal stays open
   - No API call made (check browser DevTools Network tab)

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 5: Edit with Empty Title Shows Error ❌**

**Objective:** Verify title validation

**Steps:**
1. Login as **CEO**
2. Open any task detail page
3. Click **"✏️ Edit"**
4. **Clear the Title field** (delete all text)
5. Click **"Update Mission"**
6. **Verify:**
   - Error message shows: "Title is required"
   - Modal stays open
   - No API call made

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 6: CEO Can Delete Any Task ✅**

**Objective:** Verify CEO can delete any task

**Steps:**
1. Login as **CEO**
2. Navigate to a test task (not critical)
3. Click **"🗑️ Abort"** button
4. **Verify Delete Modal Opens:**
   - Red theme border
   - Header: "🗑️ Abort Mission?"
   - Warning: "This action cannot be undone"
   - Lists what will be removed (mission data, submissions, etc.)
   - Shows task title and ID
5. Click **"💥 Yes, Delete Forever"**
6. **Verify:**
   - Loading spinner appears
   - Alert shows: "💥 Mission Deleted Successfully! Returning to dashboard..."
   - Redirects to `/dashboard`
7. **Verify task is gone:**
   - Task not in dashboard list
   - Direct URL returns 404: `http://localhost:3000/task/{deleted-id}`
8. Check API logs:
   ```bash
   docker compose logs api --tail 10 | grep "Deleted"
   ```
9. **Expected Log:**
   ```
   🗑️ Task Deleted: uuid by CEO (User ID: 1)
   ```

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 7: Creator Can Delete Own Task ✅**

**Objective:** Verify creator can delete their own task

**Steps:**
1. Login as **DEV**
2. Create a new test task
3. Open the task detail page
4. Click **"🗑️ Abort"**
5. Confirm deletion
6. **Verify:**
   - Task deleted successfully
   - Redirects to dashboard

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 8: Non-Creator Cannot Delete Other's Task ❌**

**Objective:** Verify permission enforcement (non-creator cannot delete)

**Steps:**
1. Login as **DEV**
2. Navigate to a task created by CEO
3. **Verify:**
   - "🗑️ Abort" button **NOT VISIBLE**
4. Try direct API call:
   ```bash
   curl -X DELETE "http://localhost:8080/api/v1/sentinel/tasks/$CEO_TASK_ID" \
     -H "Authorization: Bearer $DEV_TOKEN" \
     -s
   ```
5. **Expected Response:**
   ```json
   {
     "error": "Forbidden",
     "message": "unauthorized: only the task creator or CEO can delete this task"
   }
   ```

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 9: Edit Triggers AI Re-estimation ✅**

**Objective:** Verify AI re-estimates when content changes

**Steps:**
1. Login as **CEO**
2. Open any task with current estimate (e.g., 30 minutes)
3. Note the current **AI Estimated Time**
4. Click **"✏️ Edit"**
5. Change **Title** to: `"Implement Complex Microservices Architecture with Kubernetes"`
6. Change **Description** to: `"Build a scalable microservices system with Docker, Kubernetes, service mesh, observability, and CI/CD pipeline"`
7. Click **"Update Mission"**
8. Wait for page refresh
9. **Verify:**
   - AI Estimated Time CHANGED (e.g., 30min → 480min)
   - New estimate reflects increased complexity
10. Check logs for AI activity:
    ```bash
    docker compose logs api --tail 30 | grep -A 5 "AI Re-estimation"
    ```

**Result:** ✅ PASS / ❌ FAIL

---

### **Test 10: Edit Resets Pending Negotiation ✅**

**Objective:** Verify negotiation status resets after AI re-estimation

**Steps:**
1. Login as **DEV**
2. Create a task or open existing one
3. Submit a **Time Negotiation** (click "⏱️✋ Dispute AI Estimate")
   - Propose: 480 minutes
   - Reason: "Complex legacy code"
4. **Verify** negotiation status is "PENDING"
5. Now login as **CEO**
6. Open the same task
7. Click **"✏️ Edit"**
8. Change the **Description** (to trigger AI re-estimation)
9. Click **"Update Mission"**
10. **Verify:**
    - Negotiation status changes from "PENDING" → "NONE"
    - Proposed minutes reset to 0
    - New AI estimate shown

**Result:** ✅ PASS / ❌ FAIL

---

## 📊 Test Summary Template

```
Date: ___________
Tester: ___________

| Test # | Test Name | Result | Notes |
|--------|-----------|--------|-------|
| 1 | CEO Edit Any Task | ✅ / ❌ | |
| 2 | Creator Edit Own | ✅ / ❌ | |
| 3 | Non-Creator Denied | ✅ / ❌ | |
| 4 | No Changes Error | ✅ / ❌ | |
| 5 | Empty Title Error | ✅ / ❌ | |
| 6 | CEO Delete Any | ✅ / ❌ | |
| 7 | Creator Delete Own | ✅ / ❌ | |
| 8 | Non-Creator Delete Denied | ✅ / ❌ | |
| 9 | AI Re-estimation | ✅ / ❌ | |
| 10 | Negotiation Reset | ✅ / ❌ | |

Overall Status: ✅ ALL PASS / ❌ SOME FAILED

Issues Found:
1. ___________
2. ___________
```

---

## 🔍 Visual Verification Checklist

### **Edit Modal:**
- [ ] Modal has blue border (`border-blue-600`)
- [ ] Header shows "✏️ Edit Mission"
- [ ] Yellow warning banner visible
- [ ] Title field pre-filled
- [ ] Description field pre-filled
- [ ] Deadline field pre-filled (if exists)
- [ ] Update button enabled when form valid
- [ ] Cancel button works
- [ ] Loading spinner shows during update
- [ ] Error messages display correctly

### **Delete Modal:**
- [ ] Modal has red border (`border-red-600`)
- [ ] Header shows "🗑️ Abort Mission?"
- [ ] Red warning banner visible
- [ ] Task title and ID displayed
- [ ] Bullet list shows what will be removed
- [ ] Delete button has destructive styling
- [ ] Cancel button works
- [ ] Loading spinner shows during deletion
- [ ] Redirects after successful deletion

### **Permission Checks:**
- [ ] CEO sees Edit/Delete on ALL tasks
- [ ] Creator sees Edit/Delete on OWN tasks only
- [ ] Non-creator does NOT see Edit/Delete buttons
- [ ] Backend returns 403 for unauthorized API calls

---

## 🐛 Common Issues & Solutions

### **Issue 1: Buttons Not Showing**
**Symptom:** Edit/Delete buttons not visible even for CEO  
**Check:**
- Is user logged in?
- Is JWT token valid? (Check browser DevTools → Application → Cookies)
- Is `authStore.user` populated? (Console: `authStore.user`)
- Is `task.created_by` populated? (Console: check task data)

**Solution:**
```javascript
// In browser console
console.log('User:', authStore.user)
console.log('Task:', task)
console.log('Can Edit:', canEditOrDelete)
```

---

### **Issue 2: Modal Not Opening**
**Symptom:** Clicking Edit/Delete does nothing  
**Check:**
- Browser console for JavaScript errors
- Vue DevTools for component state

**Solution:**
- Refresh page
- Check if modals are blocked by z-index issues
- Verify `showEditModal` / `showDeleteModal` refs

---

### **Issue 3: API Call Fails**
**Symptom:** Error message in modal  
**Check:**
- API is running: `docker compose ps api`
- Backend logs: `docker compose logs api --tail 50`
- Network tab in DevTools (check status code)

**Solution:**
- If 403: Permission denied (expected for non-authorized)
- If 500: Check backend logs for error details
- If 404: Task not found (may have been deleted)

---

### **Issue 4: AI Not Re-estimating**
**Symptom:** AI estimate doesn't change after edit  
**Check:**
- Did you actually change title/description?
- Check API logs for "AI Re-estimation" message
- Gemini API key valid?

**Solution:**
```bash
# Check if AI estimation is being triggered
docker compose logs api | grep -A 10 "Task content changed"

# Check Gemini API response
docker compose logs api | grep "Gemini"
```

---

## 📞 Support Commands

### **Check Services**
```bash
docker compose ps
docker compose logs api --tail 50
docker compose logs web --tail 50
```

### **Reset Database (if needed)**
```bash
docker compose down
docker compose up -d postgres
# Wait 5 seconds
docker compose up -d api web
```

### **Get Fresh JWT Token**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ceo@sentinel.com","password":"password123"}' \
  -s | jq -r '.token'
```

---

## ✅ Sign-Off

**Tested By:** _________________  
**Date:** _________________  
**Status:** ✅ APPROVED / ❌ NEEDS FIXES  
**Signature:** _________________  

---

**🎉 Complete this testing guide to verify the feature is production-ready!** 🧪✅
