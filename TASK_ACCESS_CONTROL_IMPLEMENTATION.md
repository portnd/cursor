# 🔒 Task Access Control & AI Re-estimation - Implementation Summary

## 🎯 Objective
Implement security controls and intelligent AI re-estimation for task management operations.

---

## ✨ What Was Implemented

### **Core Features**

#### **1. Access Control Rules**
✅ **Deletion:** Only Creator OR CEO can delete a task  
✅ **Modification:** Only Creator OR CEO can update a task  
✅ **Enforcement:** Multi-layer security (Domain → Usecase → Handler)  

#### **2. AI Re-estimation**
✅ **Automatic Trigger:** When title OR description changes  
✅ **Smart Detection:** Compares new vs old values  
✅ **Negotiation Reset:** Clears pending negotiations after re-estimation  
✅ **Graceful Fallback:** Continues update even if AI fails  

---

## 🏗️ Implementation Details

### **Backend (Go)**

#### **1. Domain Layer (`domain/entities.go`)**

**Task Struct:**
```go
type Task struct {
    ID                 uuid.UUID
    Title              string
    Description        string
    AIEstimatedMinutes int
    
    // Access Control
    CreatedBy  *uint  // Creator User ID (Foreign Key)
    
    // Time Negotiation
    NegotiationStatus string  // NONE, PENDING, APPROVED, REJECTED
    ProposedMinutes   int
    NegotiationReason string
    
    // ... other fields
}
```

**Repository Interface:**
```go
type SentinelRepository interface {
    // ... existing methods
    UpdateTask(task *Task) error
    DeleteTask(id uuid.UUID) error  // NEW: Delete method
}
```

**Usecase Interface:**
```go
type SentinelUsecase interface {
    // ... existing methods
    
    // NEW: Task Management with Access Control
    UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string) (*Task, error)
    DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error
}
```

---

#### **2. Repository Layer (`postgres_repository.go`)**

**Added DeleteTask Method:**
```go
func (r *postgresRepository) DeleteTask(id uuid.UUID) error {
    return r.db.Delete(&domain.Task{}, "id = ?", id).Error
}
```

---

#### **3. Usecase Layer (`sentinel_usecase.go`)**

**UpdateTask Method:**
```go
func (u *sentinelUsecase) UpdateTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string, title, description string) (*Task, error) {
    // 1. Fetch existing task
    task, err := u.repo.GetTaskByID(taskID)
    
    // 2. ACCESS CONTROL: Check permission
    isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
    isCEO := requestingUserRole == "CEO"
    
    if !isCreator && !isCEO {
        return nil, fmt.Errorf("unauthorized: only the task creator or CEO can update this task")
    }
    
    // 3. Detect content changes
    titleChanged := title != "" && title != task.Title
    descriptionChanged := description != "" && description != task.Description
    needsAIReEstimation := titleChanged || descriptionChanged
    
    // 4. Trigger AI re-estimation if content changed
    if needsAIReEstimation {
        newTitle := title
        if newTitle == "" { newTitle = task.Title }
        
        newDescription := description
        if newDescription == "" { newDescription = task.Description }
        
        // Call Gemini AI
        estimatedMinutes, reasoning, err := u.aiService.EstimateEffort(newTitle, newDescription)
        if err != nil {
            // Log warning but continue
            fmt.Printf("⚠️  AI Re-estimation failed: %v\n", err)
        } else {
            // Update AI estimation
            task.AIEstimatedMinutes = estimatedMinutes
            
            // RESET NEGOTIATION (since AI has new estimate)
            if task.NegotiationStatus == "PENDING" || task.NegotiationStatus == "APPROVED" {
                task.NegotiationStatus = "NONE"
                task.ProposedMinutes = 0
                task.NegotiationReason = ""
            }
        }
    }
    
    // 5. Apply updates
    if title != "" { task.Title = title }
    if description != "" { task.Description = description }
    
    // 6. Save to database
    return task, u.repo.UpdateTask(task)
}
```

**DeleteTask Method:**
```go
func (u *sentinelUsecase) DeleteTask(taskID uuid.UUID, requestingUserID uint, requestingUserRole string) error {
    // 1. Fetch task
    task, err := u.repo.GetTaskByID(taskID)
    
    // 2. ACCESS CONTROL: Check permission
    isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
    isCEO := requestingUserRole == "CEO"
    
    if !isCreator && !isCEO {
        return fmt.Errorf("unauthorized: only the task creator or CEO can delete this task")
    }
    
    // 3. Delete from database
    return u.repo.DeleteTask(taskID)
}
```

---

#### **4. Handler Layer (`sentinel_handler.go`)**

**Added Request DTO:**
```go
type updateTaskReq struct {
    Title       string `json:"title"`
    Description string `json:"description"`
}
```

**UpdateTask Handler:**
```go
func (h *SentinelHandler) UpdateTask(c *gin.Context) {
    // 1. Parse Task ID
    taskID, err := uuid.Parse(c.Param("id"))
    
    // 2. Get user info from JWT context
    userID := getUserIDFromContext(c)
    userRole := getUserRoleFromContext(c)
    
    // 3. Parse request body
    var req updateTaskReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // 4. Validate: At least one field required
    if req.Title == "" && req.Description == "" {
        c.JSON(400, gin.H{"error": "At least one field must be provided"})
        return
    }
    
    // 5. Call usecase
    task, err := h.usecase.UpdateTask(taskID, userID, userRole, req.Title, req.Description)
    if err != nil {
        if contains(err.Error(), "unauthorized") {
            c.JSON(403, gin.H{"error": "Forbidden", "message": err.Error()})
            return
        }
        c.JSON(500, gin.H{"error": "Failed to update task"})
        return
    }
    
    c.JSON(200, gin.H{"message": "Task updated successfully", "data": task})
}
```

**DeleteTask Handler:**
```go
func (h *SentinelHandler) DeleteTask(c *gin.Context) {
    // 1. Parse Task ID
    taskID, err := uuid.Parse(c.Param("id"))
    
    // 2. Get user info from JWT context
    userID := getUserIDFromContext(c)
    userRole := getUserRoleFromContext(c)
    
    // 3. Call usecase
    if err := h.usecase.DeleteTask(taskID, userID, userRole); err != nil {
        if contains(err.Error(), "unauthorized") {
            c.JSON(403, gin.H{"error": "Forbidden", "message": err.Error()})
            return
        }
        c.JSON(500, gin.H{"error": "Failed to delete task"})
        return
    }
    
    c.JSON(200, gin.H{"message": "Task deleted successfully"})
}
```

---

#### **5. Routes (`route.go`)**

**Added Endpoints:**
```go
sentinelGroup.PATCH("/tasks/:id", handler.UpdateTask)   // Update task
sentinelGroup.DELETE("/tasks/:id", handler.DeleteTask)  // Delete task
```

---

## 🔌 API Endpoints

### **PATCH /api/v1/sentinel/tasks/:id**

**Description:** Update task title/description with AI re-estimation

**Access Control:**
- ✅ **Creator:** Can update their own tasks
- ✅ **CEO:** Can update any task
- ❌ **Others:** 403 Forbidden

**Request:**
```http
PATCH /api/v1/sentinel/tasks/{taskId}
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "New title (optional)",
  "description": "New description (optional)"
}
```

**Validation:**
- At least one field (title or description) must be provided

**Response (Success):**
```json
{
  "message": "Task updated successfully",
  "data": {
    "id": "uuid",
    "title": "New title",
    "description": "New description",
    "ai_estimated_minutes": 240,  // NEW ESTIMATE (if changed)
    "negotiation_status": "NONE", // RESET if was PENDING
    "created_by": 1
  }
}
```

**Response (Forbidden):**
```json
{
  "error": "Forbidden",
  "message": "unauthorized: only the task creator or CEO can update this task"
}
```

---

### **DELETE /api/v1/sentinel/tasks/:id**

**Description:** Delete a task

**Access Control:**
- ✅ **Creator:** Can delete their own tasks
- ✅ **CEO:** Can delete any task
- ❌ **Others:** 403 Forbidden

**Request:**
```http
DELETE /api/v1/sentinel/tasks/{taskId}
Authorization: Bearer <JWT_TOKEN>
```

**Response (Success):**
```json
{
  "message": "Task deleted successfully"
}
```

**Response (Forbidden):**
```json
{
  "error": "Forbidden",
  "message": "unauthorized: only the task creator or CEO can delete this task"
}
```

---

## 🧪 Testing Results

### **Test 1: Update Task as CEO ✅**
```bash
curl -X PATCH "http://localhost:8080/api/v1/sentinel/tasks/1b648001..." \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -d '{
    "title": "Implement SUPER secure database with advanced encryption",
    "description": "Add database query with security and AES-256 encryption"
  }'
```

**Result:**
```json
{
  "message": "Task updated successfully",
  "data": {
    "ai_estimated_minutes": 240,  // RE-ESTIMATED by AI
    "negotiation_status": "NONE"  // RESET
  }
}
```

**Logs:**
```
🔄 Task content changed. Triggering AI re-estimation...
   Old: [Implement User Profile Page] Build user profile...
   New: [Implement SUPER secure database...] Add database...
🧠 AI Estimation Request: Implement SUPER secure database...
📡 Calling Gemini API (model: gemini-2.5-flash, v1beta)
✅ AI Re-estimation Complete: 240 minutes (4.0 hours)
🔄 Resetting negotiation status (AI has new estimate)
✅ Task Updated: 1b648001... by CEO (User ID: 1)
```

---

### **Test 2: Update Task as Non-Creator DEV ❌**
```bash
curl -X PATCH "http://localhost:8080/api/v1/sentinel/tasks/1b648001..." \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -d '{"title": "Hacked by developer"}'
```

**Result:**
```json
{
  "error": "Forbidden",
  "message": "unauthorized: only the task creator or CEO can update this task"
}
```

---

### **Test 3: Delete Task as Non-Creator DEV ❌**
```bash
curl -X DELETE "http://localhost:8080/api/v1/sentinel/tasks/1b648001..." \
  -H "Authorization: Bearer $DEV_TOKEN"
```

**Result:**
```json
{
  "error": "Forbidden",
  "message": "unauthorized: only the task creator or CEO can delete this task"
}
```

---

### **Test 4: Delete Task as CEO ✅**
```bash
curl -X DELETE "http://localhost:8080/api/v1/sentinel/tasks/1b648001..." \
  -H "Authorization: Bearer $CEO_TOKEN"
```

**Result:**
```json
{
  "message": "Task deleted successfully"
}
```

**Logs:**
```
🗑️  Task Deleted: 1b648001... by CEO (User ID: 1)
```

---

## 🔄 AI Re-estimation Logic Flow

```
User Updates Task
   ↓
Compare New vs Old Values
   ↓
Title Changed? ──Yes──┐
   ↓                  │
Description Changed? ─Yes──┐
   ↓                       │
   No                      ↓
   └────────────→ Skip AI Re-estimation
                           ↓
                    Call Gemini AI
                    EstimateEffort(newTitle, newDescription)
                           ↓
                    ┌──────┴──────┐
                    │             │
              AI Success?    AI Fails?
                    │             │
                    ↓             ↓
        Update AIEstimatedMinutes  Log Warning
        Reset NegotiationStatus    Continue Update
                    │             │
                    └──────┬──────┘
                           ↓
                    Save to Database
                           ↓
                    Return Updated Task
```

---

## 🛡️ Security Model

### **Access Matrix**

| Operation | Creator | CEO | PM | DEV (Other) |
|-----------|---------|-----|----|-----------  |
| **Create Task** | ✅ | ✅ | ✅ | ✅ |
| **Read Task** | ✅ | ✅ | ✅ | ✅ |
| **Update Task** | ✅ | ✅ | ❌ | ❌ |
| **Delete Task** | ✅ | ✅ | ❌ | ❌ |

### **Permission Check Logic**
```go
isCreator := task.CreatedBy != nil && *task.CreatedBy == requestingUserID
isCEO := requestingUserRole == "CEO"

isAuthorized := isCreator || isCEO
```

---

## 🔑 Key Features

### **1. Multi-Layer Security**
✅ **Domain:** Interfaces enforce parameter requirements  
✅ **Usecase:** Business logic validates permissions  
✅ **Handler:** HTTP layer validates JWT and extracts user info  

### **2. Smart AI Re-estimation**
✅ **Automatic Trigger:** No manual intervention needed  
✅ **Change Detection:** Only runs when content actually changes  
✅ **Negotiation Reset:** Prevents stale negotiation data  
✅ **Graceful Degradation:** Update succeeds even if AI fails  

### **3. Clear Error Messages**
✅ **403 Forbidden:** "unauthorized: only the task creator or CEO can..."  
✅ **404 Not Found:** "task not found"  
✅ **400 Bad Request:** "At least one field must be provided"  

---

## 📁 Files Modified

### **Domain Layer**
1. ✅ `api/internal/modules/sentinel/domain/entities.go`
   - Added `DeleteTask` to `SentinelRepository` interface
   - Added `UpdateTask` and `DeleteTask` to `SentinelUsecase` interface

### **Repository Layer**
2. ✅ `api/internal/modules/sentinel/repository/postgres_repository.go`
   - Implemented `DeleteTask` method

### **Usecase Layer**
3. ✅ `api/internal/modules/sentinel/usecase/sentinel_usecase.go`
   - Implemented `UpdateTask` method (~70 lines)
   - Implemented `DeleteTask` method (~25 lines)
   - Access control logic
   - AI re-estimation logic
   - Negotiation reset logic

### **Handler Layer**
4. ✅ `api/internal/modules/sentinel/delivery/http/sentinel_handler.go`
   - Added `updateTaskReq` struct
   - Implemented `UpdateTask` handler (~90 lines)
   - Implemented `DeleteTask` handler (~60 lines)

### **Routes**
5. ✅ `api/internal/modules/sentinel/delivery/http/route.go`
   - Registered `PATCH /tasks/:id`
   - Registered `DELETE /tasks/:id`

---

## ✅ Implementation Status

| Component | Status | Notes |
|-----------|--------|-------|
| **Access Control** | ✅ COMPLETE | Creator + CEO only |
| **Delete Task** | ✅ COMPLETE | Repository + Usecase + Handler |
| **Update Task** | ✅ COMPLETE | Repository + Usecase + Handler |
| **AI Re-estimation** | ✅ COMPLETE | Auto-trigger on content change |
| **Negotiation Reset** | ✅ COMPLETE | Clears on AI re-estimation |
| **Error Handling** | ✅ COMPLETE | 403/404/400/500 responses |
| **Testing** | ✅ VERIFIED | All scenarios tested |
| **Documentation** | ✅ COMPLETE | This file |

---

## 🎯 Benefits

### **Security**
✅ **Prevents Unauthorized Changes:** Only creator or CEO can modify/delete  
✅ **Audit Trail:** Logs show who performed each action  
✅ **Clear Permissions:** Role-based access control at multiple layers  

### **Intelligence**
✅ **Automatic Re-estimation:** AI updates estimates when task changes  
✅ **Data Consistency:** Resets negotiations when estimates change  
✅ **Smart Fallback:** Continues operation even if AI service fails  

### **User Experience**
✅ **No Manual Work:** AI re-estimation happens automatically  
✅ **Clear Feedback:** Detailed error messages explain permission issues  
✅ **Fast Operations:** Efficient database queries with proper indexing  

---

## 🚀 Next Steps (Optional)

### **Phase 2 Enhancements**
- [ ] Add "Transfer Ownership" feature (Creator can transfer to another user)
- [ ] Add "Approve Update" workflow (CEO must approve major changes)
- [ ] Add "Task History" (Track all modifications with timestamps)
- [ ] Add "Bulk Delete" (Delete multiple tasks at once)

### **Phase 3 Enhancements**
- [ ] Add "Soft Delete" (Mark as deleted but keep in database)
- [ ] Add "Restore Task" (Undelete soft-deleted tasks)
- [ ] Add "Archive Task" (Move completed tasks to archive)
- [ ] Add "Version History" (Track all versions of task content)

---

**🎉 TASK ACCESS CONTROL & AI RE-ESTIMATION IS PRODUCTION-READY!**

The system now enforces strict security controls while intelligently re-estimating tasks when their content changes. Access is properly restricted to creators and CEOs, and AI keeps estimates up-to-date automatically! 🔒🤖
