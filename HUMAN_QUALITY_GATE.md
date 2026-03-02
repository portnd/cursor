# 🚦 Human Quality Gate - Implementation Complete

**Date:** 2026-01-26  
**Status:** ✅ DEPLOYED  
**Priority:** 🔴 CRITICAL WORKFLOW CHANGE  

---

## 🎯 What Changed

### **Critical Behavior Change:**

**BEFORE:**
```
AI Review: PASS → Task Status: COMPLETED ✅ (automatic)
```

**AFTER:**
```
AI Review: PASS → Task Status: REVIEW_PENDING 🚦 → Human Approval → COMPLETED ✅
```

---

## 🧠 Why This Matters

### **The Problem:**
- AI could mark tasks as COMPLETED automatically
- No human verification
- Risky for production code

### **The Solution:**
- **Human Quality Gate** - PM/CEO must approve before completion
- AI verdict is advisory, not final
- Ensures human oversight on all completed work

---

## 📦 Implementation Details

### **1. Updated SubmitWork Flow**

**File:** `api/internal/modules/sentinel/usecase/sentinel_usecase.go`

**Before (Lines 148-171):**
```go
// ⏰ Time Tracking: Mark task as COMPLETED if AI verdict is PASS
if sub.AIVerdict == "PASS" {
    task.Status = "COMPLETED"
    now := time.Now()
    task.CompletedAt = &now
    // ...
}
```

**After:**
```go
// 🚦 Human Quality Gate: If AI PASS, move to REVIEW_PENDING (not COMPLETED)
if sub.AIVerdict == "PASS" {
    task.Status = "REVIEW_PENDING" // 🔒 Requires PM/CEO approval
    // Do NOT set CompletedAt yet - human approval required
    // ...
}
```

---

### **2. New ApproveTask Feature**

#### **A. Repository Layer**

**File:** `api/internal/modules/sentinel/repository/postgres_repository.go`

**New Method:**
```go
func (r *postgresRepository) ApproveTask(id uuid.UUID) error {
	// Use SQL UPDATE with NOW() to ensure atomic operation
	result := r.db.Exec(`
		UPDATE tasks 
		SET status = 'COMPLETED', 
		    completed_at = NOW(),
		    updated_at = NOW()
		WHERE id = ?
	`, id)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	
	return nil
}
```

**What it does:**
- Updates task status to `COMPLETED`
- Sets `completed_at` timestamp atomically
- Returns error if task not found

---

#### **B. Usecase Layer**

**File:** `api/internal/modules/sentinel/usecase/sentinel_usecase.go`

**New Method:**
```go
func (u *sentinelUsecase) ApproveTask(taskID uuid.UUID, approverID uint, approverRole string) error {
	// 🔒 ROLE VALIDATION: Only PM or CEO can approve tasks
	if approverRole != "CEO" && approverRole != "PM" {
		return fmt.Errorf("access denied: only PM or CEO can approve tasks")
	}

	// 1️⃣ Get the task
	task, err := u.repo.GetTaskByID(taskID)
	// ...

	// 2️⃣ Verify task is in REVIEW_PENDING status
	if task.Status != "REVIEW_PENDING" {
		return fmt.Errorf("task is not pending review (current status: %s)", task.Status)
	}

	// 3️⃣ Approve the task
	if err := u.repo.ApproveTask(taskID); err != nil {
		return fmt.Errorf("failed to approve task: %w", err)
	}

	// 4️⃣ Log success with time metrics
	// ...
}
```

**What it does:**
- **Permission Check:** Only PM or CEO can approve
- **Status Validation:** Task must be in `REVIEW_PENDING` status
- **Approval:** Marks task as `COMPLETED` and sets timestamp
- **Logging:** Tracks who approved and calculates actual vs estimated time

---

#### **C. Handler Layer**

**File:** `api/internal/modules/sentinel/delivery/http/sentinel_handler.go`

**New Handler:**
```go
func (h *SentinelHandler) ApproveTask(c *gin.Context) {
	// 1. Parse Task ID
	taskID, err := uuid.Parse(c.Param("id"))
	// ...

	// 2. Get approver info from context
	approverID := getUserIDFromContext(c)
	approverRole := getUserRoleFromContext(c)
	// ...

	// 3. Call usecase to approve task
	if err := h.usecase.ApproveTask(taskID, approverID, approverRole); err != nil {
		// Handle errors (403, 404, 400, 500)
		// ...
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task approved and marked as COMPLETED",
	})
}
```

**HTTP Response Codes:**
- `200 OK` - Task approved successfully
- `400 Bad Request` - Task not in REVIEW_PENDING status
- `401 Unauthorized` - User not authenticated
- `403 Forbidden` - User is not PM or CEO
- `404 Not Found` - Task doesn't exist
- `500 Internal Server Error` - Database error

---

#### **D. Route Registration**

**File:** `api/internal/modules/sentinel/delivery/http/route.go`

**New Endpoint:**
```go
sentinelGroup.POST("/tasks/:id/approve", handler.ApproveTask) // Approve task after review (PM/CEO only)
```

**Full URL:** `POST /api/v1/sentinel/tasks/:id/approve`

---

## 🔄 Complete Workflow

### **Developer Submits Code:**

```
1. Developer: POST /api/v1/sentinel/tasks/{id}/submit
   {
     "commit_hash": "abc123",
     "diff": "..."
   }

2. AI Reviews Code:
   • Runs security audit
   • Checks for SQL injection, secrets, etc.
   • Returns verdict: PASS or FAIL

3. IF AI verdict = FAIL:
   → Task Status: PENDING
   → Developer can appeal

4. IF AI verdict = PASS:
   → Task Status: REVIEW_PENDING 🚦 (NEW!)
   → Awaiting PM/CEO approval
```

---

### **PM/CEO Approves Task:**

```
5. PM/CEO: GET /api/v1/sentinel/tasks/approvals
   → Sees tasks with Status = REVIEW_PENDING

6. PM/CEO Reviews:
   • Checks code diff
   • Verifies AI feedback
   • Makes final decision

7. PM/CEO: POST /api/v1/sentinel/tasks/{id}/approve
   → Task Status: COMPLETED ✅
   → completed_at = NOW()
   → Time metrics calculated

8. System Logs:
   ✅ Task {id} APPROVED by CEO (ID: 1)
   🎉 Task marked COMPLETED at 2026-01-26T10:46:12Z
   📊 Actual Time: 3.5 hours (AI Estimated: 4.0 hours)
```

---

## 🎨 Task Status Flow

```
PENDING
   ↓ (Assign to developer)
ASSIGNED
   ↓ (Developer submits work)
UNDER_REVIEW (AI Review in progress)
   ↓
   ├─ AI: FAIL → PENDING (can appeal)
   │
   └─ AI: PASS → REVIEW_PENDING 🚦 (NEW STATUS!)
                     ↓
                  PM/CEO Approves
                     ↓
                  COMPLETED ✅
```

---

## 🧪 Testing Guide

### **Test Case 1: Submit Work with AI PASS**

```bash
# 1. Submit work
curl -X POST http://localhost:8080/api/v1/sentinel/tasks/{task_id}/submit \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commit_hash": "abc123",
    "diff": "db.Where(\"email = ?\", email).First(&user)"
  }'

# Expected Response:
{
  "id": "...",
  "ai_verdict": "PASS",
  "ai_score": 90,
  "ai_feedback": {
    "feedback": "✅ ใช้ Parameterized Query ถูกต้อง..."
  }
}

# 2. Check task status
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/{task_id} \
  -H "Authorization: Bearer $DEV_TOKEN"

# Expected:
{
  "status": "REVIEW_PENDING",  ← NOT "COMPLETED"!
  "completed_at": null           ← Still null
}
```

---

### **Test Case 2: PM/CEO Approves Task**

```bash
# 1. PM/CEO checks approvals inbox
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/approvals \
  -H "Authorization: Bearer $CEO_TOKEN"

# Expected: List of tasks with status = REVIEW_PENDING

# 2. PM/CEO approves task
curl -X POST http://localhost:8080/api/v1/sentinel/tasks/{task_id}/approve \
  -H "Authorization: Bearer $CEO_TOKEN"

# Expected Response (200 OK):
{
  "message": "Task approved and marked as COMPLETED"
}

# 3. Verify task is completed
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/{task_id} \
  -H "Authorization: Bearer $CEO_TOKEN"

# Expected:
{
  "status": "COMPLETED",
  "completed_at": "2026-01-26T10:46:12Z"  ← Now set!
}
```

---

### **Test Case 3: Developer Tries to Approve (Should Fail)**

```bash
curl -X POST http://localhost:8080/api/v1/sentinel/tasks/{task_id}/approve \
  -H "Authorization: Bearer $DEV_TOKEN"

# Expected Response (403 Forbidden):
{
  "error": "Forbidden",
  "message": "access denied: only PM or CEO can approve tasks (your role: DEV)"
}
```

---

### **Test Case 4: Approve Task Not in REVIEW_PENDING**

```bash
# Try to approve a task that's already COMPLETED
curl -X POST http://localhost:8080/api/v1/sentinel/tasks/{task_id}/approve \
  -H "Authorization: Bearer $CEO_TOKEN"

# Expected Response (400 Bad Request):
{
  "error": "Invalid Status",
  "message": "task is not pending review (current status: COMPLETED)"
}
```

---

## 📊 Database Changes

### **No Schema Changes Required! ✅**

**Why:**
- We're using existing `status` field (already supports custom values)
- `completed_at` field already exists as `*time.Time` (nullable)
- New status value: `REVIEW_PENDING` (no migration needed)

**Existing Schema:**
```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    status VARCHAR(50) DEFAULT 'PENDING',  ← Can store "REVIEW_PENDING"
    completed_at TIMESTAMP,                ← Already nullable
    -- ...
);
```

---

## 🔐 Security & Permissions

### **Role-Based Access Control:**

| Action | DEV | PM | CEO |
|--------|-----|----|----|
| Submit Work | ✅ | ✅ | ✅ |
| Approve Task | ❌ | ✅ | ✅ |
| View Approvals Inbox | ❌ | ✅ | ✅ |

**Implementation:**
```go
if approverRole != "CEO" && approverRole != "PM" {
    return fmt.Errorf("access denied: only PM or CEO can approve tasks")
}
```

---

## 📈 Metrics & Logging

### **Success Log Example:**
```
🚦 Task abc-123 moved to REVIEW_PENDING - awaiting PM/CEO approval
📋 AI Review: PASS (90/100) - Ready for human verification

✅ Task abc-123 APPROVED by CEO (ID: 1)
🎉 Task marked COMPLETED at 2026-01-26T10:46:12Z
📊 Actual Time: 3.5 hours (AI Estimated: 4.0 hours)
```

### **Error Log Example:**
```
⚠️  Approval failed: access denied: only PM or CEO can approve tasks (your role: DEV)
⚠️  Approval failed: task is not pending review (current status: COMPLETED)
```

---

## 🎓 Key Benefits

### **1. Human Oversight**
- ✅ Every completed task verified by PM/CEO
- ✅ Catches AI false positives
- ✅ Final decision with human judgment

### **2. Quality Assurance**
- ✅ Double-check on critical code
- ✅ Ensures production readiness
- ✅ Reduces deployment risks

### **3. Accountability**
- ✅ Tracks who approved each task
- ✅ Audit trail for completed work
- ✅ Clear responsibility chain

### **4. Process Control**
- ✅ CEO/PM gate before completion
- ✅ Prevents premature completion
- ✅ Maintains quality standards

---

## 📚 API Documentation

### **POST /api/v1/sentinel/tasks/:id/approve**

**Description:** Approve a task in REVIEW_PENDING status (PM/CEO only)

**Authentication:** Required (JWT token)

**Authorization:** PM or CEO role

**Request:**
```bash
POST /api/v1/sentinel/tasks/abc-123-456/approve
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK):**
```json
{
  "message": "Task approved and marked as COMPLETED"
}
```

**Error Responses:**

**401 Unauthorized:**
```json
{
  "error": "Unauthorized",
  "message": "user not authenticated"
}
```

**403 Forbidden:**
```json
{
  "error": "Forbidden",
  "message": "access denied: only PM or CEO can approve tasks (your role: DEV)"
}
```

**404 Not Found:**
```json
{
  "error": "Not Found",
  "message": "task not found"
}
```

**400 Bad Request:**
```json
{
  "error": "Invalid Status",
  "message": "task is not pending review (current status: COMPLETED)"
}
```

---

## 🚀 Deployment Status

```bash
✅ Domain interfaces updated
✅ Repository implementation added
✅ Usecase logic implemented
✅ HTTP handler created
✅ Route registered
✅ API restarted successfully
✅ Endpoint: POST /api/v1/sentinel/tasks/:id/approve
✅ Linter errors: 0
✅ Status: LIVE & READY
```

---

## ✅ Checklist

- [x] Updated SubmitWork to set REVIEW_PENDING instead of COMPLETED
- [x] Added ApproveTask to repository interface
- [x] Implemented ApproveTask in postgres_repository.go
- [x] Added ApproveTask to usecase interface
- [x] Implemented ApproveTask in sentinel_usecase.go with permission checks
- [x] Created ApproveTask handler in sentinel_handler.go
- [x] Registered POST /tasks/:id/approve route
- [x] No linter errors
- [x] API restarted and verified
- [x] Documentation created

---

## 🎉 Summary

**Feature:** Human Quality Gate  
**Status:** ✅ COMPLETE & DEPLOYED  
**Impact:** HIGH - Changes core workflow  

**Key Changes:**
1. ✅ AI PASS → REVIEW_PENDING (not COMPLETED)
2. ✅ PM/CEO approval required
3. ✅ New endpoint: POST /tasks/:id/approve
4. ✅ Role-based access control
5. ✅ Time metrics tracking

**Next Steps:**
- Test the new workflow
- Update frontend to show REVIEW_PENDING status
- Add approval button in UI for PM/CEO

---

**The Human Quality Gate is now active! All tasks require PM/CEO approval before completion! 🚦**
