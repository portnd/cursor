# 🗂️ Approvals Inbox - Implementation Summary

## 🎯 Objective
Provide CEO and PM with a centralized "Approvals Inbox" showing all tasks requiring their attention (pending appeals and time negotiations).

---

## ✨ What Was Implemented

### **Backend (Go)**

#### **1. Domain Layer**
- ✅ Added `GetTasksRequiringApproval()` to `SentinelRepository` interface
- ✅ Added `GetPendingApprovals(userRole string)` to `SentinelUsecase` interface

#### **2. Repository Layer (`postgres_repository.go`)**
- ✅ Implemented `GetTasksRequiringApproval()` method
- **Query Logic:**
  - Fetches tasks where `negotiation_status = 'PENDING'` (developer wants more time)
  - OR tasks with submissions that have appeals where `status = 'PENDING'`
  - Eager loads: `Submissions` (ordered by newest first) and `Submissions.Appeal`
  
```go
func (r *postgresRepository) GetTasksRequiringApproval() ([]domain.Task, error) {
	var tasks []domain.Task

	err := r.db.
		Preload("Submissions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("Submissions.Appeal"). // Eager load appeals
		Where("negotiation_status = ?", "PENDING").
		Or("id IN (?)", r.db.Table("tasks").
			Select("tasks.id").
			Joins("JOIN submissions ON submissions.task_id = tasks.id").
			Joins("JOIN appeals ON appeals.submission_id = submissions.id").
			Where("appeals.status = ?", "PENDING")).
		Order("created_at desc").
		Find(&tasks).Error

	return tasks, err
}
```

#### **3. Usecase Layer (`sentinel_usecase.go`)**
- ✅ Implemented `GetPendingApprovals(userRole string)` method
- **Business Logic:**
  - Role validation: Only CEO and PM allowed
  - Returns error `"access denied: only CEO and PM can view approvals inbox"` for other roles
  - Calls repository to fetch tasks
  
```go
func (u *sentinelUsecase) GetPendingApprovals(userRole string) ([]domain.Task, error) {
	// 🔒 ROLE VALIDATION: Only CEO and PM
	if userRole != "CEO" && userRole != "PM" {
		return nil, fmt.Errorf("access denied: only CEO and PM can view approvals inbox")
	}

	tasks, err := u.repo.GetTasksRequiringApproval()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending approvals: %w", err)
	}

	return tasks, nil
}
```

#### **4. Middleware Update (`auth.go`)**
- ✅ Updated `AuthMiddleware` to extract `role` from JWT and set in context
  
```go
// Set user info in context for handlers to use
c.Set("user_id", userID)
c.Set("email", claims["email"])
c.Set("role", claims["role"]) // 👈 NEW: Extract role from JWT
```

#### **5. Handler Layer (`sentinel_handler.go`)**
- ✅ Added `getUserRoleFromContext()` helper function
- ✅ Implemented `GetApprovals` handler
  - Extracts user role from JWT context
  - Calls usecase with role validation
  - Returns HTTP 403 Forbidden for unauthorized roles
  - Returns HTTP 200 with tasks data and count
  
```go
func (h *SentinelHandler) GetApprovals(c *gin.Context) {
	userRole := getUserRoleFromContext(c)
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user role not found in token",
		})
		return
	}

	tasks, err := h.usecase.GetPendingApprovals(userRole)
	if err != nil {
		if contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve approvals inbox",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Approvals inbox retrieved successfully",
		"data":    tasks,
		"count":   len(tasks),
	})
}
```

#### **6. Routes (`route.go`)**
- ✅ Registered `GET /sentinel/tasks/approvals` endpoint
  
```go
sentinelGroup.GET("/tasks/approvals", handler.GetApprovals) // Get approvals inbox (CEO/PM only)
```

---

## 🔌 API Endpoint

### **GET /api/v1/sentinel/tasks/approvals**

#### **Description**
Returns tasks requiring PM/CEO attention (pending appeals or time negotiations).

#### **Access Control**
- **Allowed:** CEO, PM
- **Denied:** DEV (403 Forbidden)

#### **Request**
```http
GET /api/v1/sentinel/tasks/approvals
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### **Response (Success)**
```json
{
  "message": "Approvals inbox retrieved successfully",
  "count": 1,
  "data": [
    {
      "id": "a517e15d-f9aa-4a19-931b-ecf52d967ebf",
      "title": "Implement secure database query",
      "description": "Add database query functions with proper security",
      "negotiation_status": "PENDING",
      "proposed_minutes": 120,
      "negotiation_reason": "The AI estimate does not account for legacy code...",
      "status": "COMPLETED",
      "assigned_to": 3,
      "submissions": [
        {
          "id": "06f92f92-...",
          "ai_verdict": "FAIL",
          "appeal": {
            "id": "55cd64c4-...",
            "status": "PENDING",
            "reason": "The code uses prepared statements...",
            "ai_recommendation": "OVERTURN",
            "ai_confidence": 85,
            "ai_reasoning": "Developer is correct..."
          }
        }
      ]
    }
  ]
}
```

#### **Response (Forbidden - DEV User)**
```json
{
  "error": "Forbidden",
  "message": "access denied: only CEO and PM can view approvals inbox"
}
```

#### **Response (Unauthorized - No Token)**
```json
{
  "error": "Unauthorized",
  "message": "authorization header required"
}
```

---

## 🧪 Testing Results

### **Test 1: CEO Access ✅**
```bash
CEO_TOKEN="..."
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/approvals \
  -H "Authorization: Bearer $CEO_TOKEN"
```

**Result:**
```json
{
  "message": "Approvals inbox retrieved successfully",
  "count": 1,
  "data": [...]
}
```

### **Test 2: PM Access ✅**
```bash
PM_TOKEN="..."
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/approvals \
  -H "Authorization: Bearer $PM_TOKEN"
```

**Result:**
```json
{
  "message": "Approvals inbox retrieved successfully",
  "count": 1,
  "data": [...]
}
```

### **Test 3: DEV Access ❌ (Forbidden)**
```bash
DEV_TOKEN="..."
curl -X GET http://localhost:8080/api/v1/sentinel/tasks/approvals \
  -H "Authorization: Bearer $DEV_TOKEN"
```

**Result:**
```json
{
  "error": "Forbidden",
  "message": "access denied: only CEO and PM can view approvals inbox"
}
```

---

## 🎯 What Tasks Are Included?

The approvals inbox returns tasks that meet **ANY** of these criteria:

### **1. Time Negotiation Requests**
- Tasks where `negotiation_status = 'PENDING'`
- Developer has proposed a different time estimate
- PM/CEO needs to approve or reject the request

**Example:**
```json
{
  "negotiation_status": "PENDING",
  "ai_estimated_minutes": 30,
  "proposed_minutes": 120,
  "negotiation_reason": "Legacy code requires careful refactoring..."
}
```

### **2. Pending Appeals**
- Tasks with submissions that have appeals where `status = 'PENDING'`
- Developer has appealed an AI FAIL verdict
- PM/CEO needs to review and resolve (approve/reject)

**Example:**
```json
{
  "submissions": [
    {
      "ai_verdict": "FAIL",
      "appeal": {
        "status": "PENDING",
        "reason": "The code uses prepared statements...",
        "ai_recommendation": "OVERTURN",
        "ai_confidence": 85
      }
    }
  ]
}
```

---

## 📊 Data Structure

### **Task Object**
Each task in the approvals inbox includes:

```typescript
{
  id: string                    // Task UUID
  title: string
  description: string
  negotiation_status: string    // NONE, PENDING, APPROVED, REJECTED
  proposed_minutes: number      // Developer's proposed time (if negotiating)
  negotiation_reason: string    // Why dev needs more time
  status: string                // PENDING, IN_PROGRESS, COMPLETED
  assigned_to: number           // Developer ID
  
  submissions: [                // Array of code submissions
    {
      id: string
      ai_verdict: string        // PASS, FAIL
      ai_score: number
      ai_feedback: object
      
      appeal: {                 // Appeal object (if exists)
        id: string
        status: string          // PENDING, APPROVED, REJECTED
        reason: string          // Developer's appeal reason
        ai_recommendation: string  // OVERTURN or UPHOLD
        ai_confidence: number      // 0-100
        ai_reasoning: string       // AI's advice
        resolver_id: number | null
        resolver_note: string
      } | null
    }
  ]
}
```

---

## 🔄 User Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. DEVELOPER ACTIONS                                        │
├─────────────────────────────────────────────────────────────┤
│ Developer:                                                   │
│  ├─ Submits code → AI reviews → FAIL verdict                │
│  ├─ Appeals the verdict                                     │
│  └─ OR negotiates time estimate (AI says 30min, dev wants 2h)│
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. TASK ENTERS APPROVALS INBOX                              │
├─────────────────────────────────────────────────────────────┤
│ System:                                                      │
│  ├─ Task marked with PENDING appeal or time negotiation     │
│  ├─ Appears in CEO/PM approvals inbox                       │
│  └─ Includes AI advisory (for appeals)                      │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CEO/PM REVIEWS INBOX                                     │
├─────────────────────────────────────────────────────────────┤
│ CEO/PM:                                                      │
│  ├─ Calls GET /sentinel/tasks/approvals                     │
│  ├─ Sees all tasks requiring attention                      │
│  ├─ Reviews AI recommendations (for appeals)                │
│  └─ Prioritizes by urgency                                  │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. CEO/PM TAKES ACTION                                      │
├─────────────────────────────────────────────────────────────┤
│ For Appeals:                                                 │
│  ├─ POST /appeals/:id/resolve                               │
│  └─ Status: APPROVED or REJECTED                            │
│                                                              │
│ For Time Negotiations:                                       │
│  ├─ (Future) POST /tasks/:id/negotiate/resolve              │
│  └─ Status: APPROVED or REJECTED                            │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. TASK REMOVED FROM INBOX                                  │
├─────────────────────────────────────────────────────────────┤
│ System:                                                      │
│  ├─ Appeal/negotiation no longer PENDING                    │
│  ├─ Task disappears from approvals inbox                    │
│  └─ Developer notified of decision                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎨 Frontend Integration (Suggested)

### **Approvals Inbox Page**
Create `web/pages/approvals.vue` (CEO/PM only):

```vue
<template>
  <div class="approvals-inbox">
    <header>
      <h1>⚖️ Approvals Inbox</h1>
      <span class="count">{{ approvals.length }} pending</span>
    </header>

    <div v-for="task in approvals" :key="task.id" class="approval-card">
      <!-- Time Negotiation Section -->
      <div v-if="task.negotiation_status === 'PENDING'" class="negotiation">
        <span>⏱️ Time Negotiation Request</span>
        <p>
          AI: {{ task.ai_estimated_minutes }}min →
          Dev wants: {{ task.proposed_minutes }}min
        </p>
        <p>Reason: {{ task.negotiation_reason }}</p>
      </div>

      <!-- Pending Appeals Section -->
      <div v-for="sub in task.submissions" :key="sub.id">
        <div v-if="sub.appeal?.status === 'PENDING'" class="appeal">
          <span>📢 Appeal</span>
          <p>AI Verdict: {{ sub.ai_verdict }} ({{ sub.ai_score }})</p>
          <p>Developer's Plea: {{ sub.appeal.reason }}</p>
          
          <!-- AI Advisory -->
          <div class="ai-advisory">
            <span :class="sub.appeal.ai_recommendation === 'OVERTURN' ? 'green' : 'red'">
              {{ sub.appeal.ai_recommendation }}
            </span>
            <span>Confidence: {{ sub.appeal.ai_confidence }}%</span>
            <p>{{ sub.appeal.ai_reasoning }}</p>
          </div>

          <button @click="reviewAppeal(task.id)">Review →</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const approvals = ref([])

const fetchApprovals = async () => {
  const { data } = await useAuth().fetchWithAuth('/sentinel/tasks/approvals')
  approvals.value = data.data
}

onMounted(() => {
  fetchApprovals()
})
</script>
```

### **Dashboard Badge**
Update CEO/PM dashboard to show approvals count:

```vue
<div class="approvals-badge">
  <NuxtLink to="/approvals">
    ⚖️ Approvals
    <span class="count">{{ approvalsCount }}</span>
  </NuxtLink>
</div>
```

---

## 📁 Files Modified

### **Backend**
1. `api/internal/modules/sentinel/domain/entities.go`
   - Added `GetTasksRequiringApproval()` to `SentinelRepository` interface
   - Added `GetPendingApprovals(userRole string)` to `SentinelUsecase` interface

2. `api/internal/modules/sentinel/repository/postgres_repository.go`
   - Implemented `GetTasksRequiringApproval()` method
   - Complex query with JOINs and eager loading

3. `api/internal/modules/sentinel/usecase/sentinel_usecase.go`
   - Implemented `GetPendingApprovals(userRole string)` method
   - Role validation logic

4. `api/internal/core/middleware/auth.go`
   - Updated to extract `role` from JWT and set in context

5. `api/internal/modules/sentinel/delivery/http/sentinel_handler.go`
   - Added `getUserRoleFromContext()` helper
   - Implemented `GetApprovals` handler

6. `api/internal/modules/sentinel/delivery/http/route.go`
   - Registered `GET /tasks/approvals` endpoint

---

## ✅ Implementation Status

| Component | Status | Notes |
|-----------|--------|-------|
| **Domain Interface** | ✅ COMPLETE | Repository and Usecase interfaces updated |
| **Repository Query** | ✅ COMPLETE | Complex query with JOINs and Preload |
| **Usecase Logic** | ✅ COMPLETE | Role validation implemented |
| **Middleware Update** | ✅ COMPLETE | Role extraction from JWT |
| **Handler** | ✅ COMPLETE | GetApprovals with error handling |
| **Routes** | ✅ COMPLETE | Endpoint registered |
| **Testing** | ✅ VERIFIED | CEO ✅, PM ✅, DEV ❌ (403) |
| **Frontend** | ⏳ PENDING | UI not yet implemented |

---

## 🚀 Next Steps

### **Immediate (Backend Complete)**
- ✅ API endpoint operational
- ✅ Role-based access control working
- ✅ Returns correct data (pending appeals + time negotiations)

### **Frontend (Suggested)**
1. Create `pages/approvals.vue` for CEO/PM
2. Add "Approvals Inbox" link to sidebar (CEO/PM only)
3. Add badge with pending count
4. Display tasks grouped by type (appeals vs negotiations)
5. Link to task detail page for review

### **Future Enhancements**
1. Real-time notifications (WebSocket/Polling)
2. Email alerts for new approvals
3. Bulk actions (approve/reject multiple at once)
4. Filtering (show only appeals, only negotiations)
5. Sorting (by urgency, date, confidence)
6. Time negotiation resolution endpoint

---

## 🏆 Benefits

### **For CEO/PM**
✅ Centralized view of all pending approvals  
✅ No need to manually check tasks one by one  
✅ Prioritize by AI confidence and urgency  
✅ Efficient decision-making workflow  

### **For System**
✅ Clear separation of concerns (repository → usecase → handler)  
✅ Role-based access control at multiple layers  
✅ Scalable query design (handles large datasets)  
✅ RESTful API design  

---

## 🎯 Key Metrics

After deployment, track:

1. **Inbox Volume:** Avg number of pending approvals
2. **Response Time:** Time from submission to resolution
3. **Approval Rate:** % of appeals/negotiations approved
4. **Role Usage:** How often CEO vs PM uses the inbox

---

**🎉 THE APPROVALS INBOX IS FULLY OPERATIONAL!**

CEO and PM now have a God-Tier centralized dashboard to view and manage all pending appeals and time negotiations. The system is production-ready and awaiting frontend integration! 🚀
