# 🛡️ Sentinel Module

**Version:** 1.0.0  
**Purpose:** Task Management, AI Evaluation, and Developer Performance Tracking

---

## 📁 Module Structure

```
internal/modules/sentinel/
├── domain/
│   ├── entities.go       # Core domain models (Task, Submission, Appeal, AuditLog)
│   └── repository.go     # Repository interface (PORT in Hexagonal Architecture)
├── repository/           # [TODO] Database implementations (PostgreSQL via GORM)
├── usecase/              # [TODO] Business logic layer
└── delivery/
    └── http/             # [TODO] HTTP handlers (REST API)
```

---

## 🏗️ Domain Entities

### 1. Task

**Purpose:** Represents a development task assigned to developers

**Fields:**
```go
type Task struct {
    ID                 uuid.UUID       // Primary key
    Title              string          // Task title (max 255 chars)
    Description        string          // Full description
    ResourceURLs       datatypes.JSON  // Figma, docs, API specs
    AIEstimatedMinutes *int            // AI-calculated effort
    DueAt              *time.Time      // Deadline
    StartedAt          *time.Time      // When developer started
    CompletedAt        *time.Time      // When completed
    Status             string          // PENDING, IN_PROGRESS, SUBMITTED, etc.
    AssignedTo         *uint           // Developer assigned (FK to users)
    CreatedBy          *uint           // CEO/PM who created it (FK to users)
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

**Status Flow:**
```
PENDING → IN_PROGRESS → SUBMITTED → COMPLETED
                           ↓
                        REJECTED → (Can appeal)
```

**Constants:**
- `TaskStatusPending`
- `TaskStatusInProgress`
- `TaskStatusSubmitted`
- `TaskStatusCompleted`
- `TaskStatusRejected`
- `TaskStatusCancelled`

**Relations:**
- `AssignedUser` - Developer assigned to this task
- `CreatedUser` - CEO/PM who created this task

---

### 2. Submission

**Purpose:** Code submission for a task with AI evaluation

**Fields:**
```go
type Submission struct {
    ID         uuid.UUID       // Primary key
    TaskID     uuid.UUID       // Task being submitted (FK to tasks)
    DevID      uint            // Developer who submitted (FK to users)
    CommitHash string          // Git commit SHA (7-64 chars)
    AIVerdict  *string         // PASS, FAIL, PENDING
    AIScore    *int            // 0-100 score
    AIFeedback datatypes.JSON  // Structured feedback from AI
    CreatedAt  time.Time
}
```

**AI Verdict Values:**
- `VerdictPass` - Auto-accepted
- `VerdictFail` - Rejected (developer can appeal)
- `VerdictPending` - Under evaluation

**AI Feedback Structure (JSONB):**
```json
{
  "code_quality": 85,
  "test_coverage": 90,
  "security_score": 95,
  "performance_score": 80,
  "issues": [
    "Missing error handling in line 45",
    "Unused import detected"
  ],
  "strengths": [
    "Good test coverage",
    "Clean code structure"
  ]
}
```

**Relations:**
- `Task` - The task this submission belongs to
- `Developer` - The user who made this submission

---

### 3. Appeal

**Purpose:** Developer appeal for disputed AI verdicts

**Fields:**
```go
type Appeal struct {
    ID           uuid.UUID  // Primary key
    TaskID       uuid.UUID  // Task being appealed (FK to tasks)
    DevID        uint       // Developer filing appeal (FK to users)
    Reason       string     // Why AI was wrong (min 10 chars)
    Status       string     // PENDING, APPROVED, REJECTED
    ReviewedBy   *uint      // CEO/PM who reviewed (FK to users)
    AdminComment string     // Admin's decision explanation
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Status Values:**
- `AppealStatusPending` - Awaiting review
- `AppealStatusApproved` - Human override (task passes)
- `AppealStatusRejected` - AI was correct

**Relations:**
- `Task` - The task being appealed
- `Developer` - The user who filed the appeal
- `Reviewer` - The admin who reviewed it

---

### 4. AuditLog

**Purpose:** Immutable event trail for compliance and debugging

**Fields:**
```go
type AuditLog struct {
    ID        uuid.UUID       // Primary key
    EventType string          // Event classification
    Metadata  datatypes.JSON  // Flexible event data
    CreatedAt time.Time
}
```

**Event Types (Constants):**
- `EventTaskCreated`
- `EventTaskAssigned`
- `EventTaskStarted`
- `EventTaskCompleted`
- `EventSubmissionCreated`
- `EventSubmissionEvaluated`
- `EventAppealFiled`
- `EventAppealReviewed`
- `EventUserRoleChanged`
- `EventUserHealthUpdated`

**Metadata Example (JSONB):**
```json
{
  "task_id": "uuid-here",
  "assigned_to": 123,
  "assigned_by": 1,
  "task_title": "Implement authentication",
  "estimated_minutes": 240
}
```

---

### 5. User (Enhanced)

**Purpose:** Extended user model with Sentinel role management

**New Fields Added:**
```go
Role        string         // CEO, PM, or DEV
HealthScore float64        // Performance score (0-100)
TechStack   pq.StringArray // Array of technologies
```

**Role Constants:**
- `RoleCEO` - Can assign tasks, review appeals
- `RolePM` - Can assign tasks
- `RoleDEV` - Default role, can submit and appeal

---

## 🔌 Repository Interface

The `SentinelRepository` interface defines all data access operations:

### Task Operations
```go
CreateTask(task *Task) error
GetTaskByID(id uuid.UUID) (*Task, error)
GetTasksByAssignee(userID uint, status string) ([]*Task, error)
GetTasksByCreator(userID uint) ([]*Task, error)
UpdateTask(task *Task) error
DeleteTask(id uuid.UUID) error
GetOverdueTasks() ([]*Task, error)
GetTasksByStatus(status string) ([]*Task, error)
```

### Submission Operations
```go
CreateSubmission(submission *Submission) error
GetSubmissionByID(id uuid.UUID) (*Submission, error)
GetSubmissionsByTask(taskID uuid.UUID) ([]*Submission, error)
GetSubmissionsByDeveloper(devID uint) ([]*Submission, error)
UpdateSubmission(submission *Submission) error
GetSubmissionByTaskAndCommit(taskID uuid.UUID, commitHash string) (*Submission, error)
```

### Appeal Operations
```go
CreateAppeal(appeal *Appeal) error
GetAppealByID(id uuid.UUID) (*Appeal, error)
GetAppealsByTask(taskID uuid.UUID) ([]*Appeal, error)
GetAppealsByDeveloper(devID uint) ([]*Appeal, error)
GetPendingAppeals() ([]*Appeal, error)
UpdateAppeal(appeal *Appeal) error
GetAppealsByReviewer(reviewerID uint) ([]*Appeal, error)
```

### Audit Log Operations
```go
CreateAuditLog(log *AuditLog) error
GetAuditLogsByEventType(eventType string, limit int) ([]*AuditLog, error)
GetRecentAuditLogs(limit int) ([]*AuditLog, error)
```

---

## 📖 Usage Examples

### Creating a Task

```go
import (
    "github.com/google/uuid"
    "github.com/komgrip/starter-kit/internal/modules/sentinel/domain"
    "gorm.io/datatypes"
)

task := &domain.Task{
    Title:       "Implement User Authentication",
    Description: "Add JWT-based auth with refresh tokens",
    ResourceURLs: datatypes.JSON([]byte(`{
        "figma": "https://figma.com/file/...",
        "docs": ["https://docs.google.com/..."]
    }`)),
    Status:      domain.TaskStatusPending,
    AssignedTo:  &devUserID,
    CreatedBy:   &ceoUserID,
}

err := repo.CreateTask(task)
```

### Creating a Submission

```go
submission := &domain.Submission{
    TaskID:     taskUUID,
    DevID:      developerID,
    CommitHash: "a7f8d9e3c2b1",
    AIVerdict:  stringPtr(domain.VerdictPending),
}

err := repo.CreateSubmission(submission)
```

### Filing an Appeal

```go
appeal := &domain.Appeal{
    TaskID: taskUUID,
    DevID:  developerID,
    Reason: "The AI failed to recognize that I implemented the caching layer correctly. Please review the Redis integration code.",
    Status: domain.AppealStatusPending,
}

err := repo.CreateAppeal(appeal)
```

### Creating an Audit Log

```go
auditLog := &domain.AuditLog{
    EventType: domain.EventTaskAssigned,
    Metadata: datatypes.JSON([]byte(`{
        "task_id": "uuid-here",
        "assigned_to": 123,
        "assigned_by": 1,
        "task_title": "Implement authentication"
    }`)),
}

err := repo.CreateAuditLog(auditLog)
```

---

## 🚀 Next Steps

### 1. Implement Repository Layer

Create `repository/postgres_repository.go`:

```go
package repository

import (
    "github.com/komgrip/starter-kit/internal/modules/sentinel/domain"
    "gorm.io/gorm"
)

type postgresRepository struct {
    db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.SentinelRepository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateTask(task *domain.Task) error {
    return r.db.Create(task).Error
}

// ... implement all other methods
```

### 2. Implement Use Cases

Create `usecase/task_usecase.go`:

```go
package usecase

import (
    "github.com/komgrip/starter-kit/internal/modules/sentinel/domain"
)

type taskUsecase struct {
    repo domain.SentinelRepository
}

func NewTaskUsecase(repo domain.SentinelRepository) TaskUsecase {
    return &taskUsecase{repo: repo}
}

func (u *taskUsecase) AssignTask(taskID uuid.UUID, devID uint) error {
    task, err := u.repo.GetTaskByID(taskID)
    if err != nil {
        return err
    }
    
    task.AssignedTo = &devID
    task.Status = domain.TaskStatusInProgress
    
    return u.repo.UpdateTask(task)
}
```

### 3. Create HTTP Handlers

Create `delivery/http/task_handler.go`:

```go
package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/komgrip/starter-kit/internal/modules/sentinel/domain"
)

type TaskHandler struct {
    usecase TaskUsecase
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
    var req domain.Task
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.usecase.CreateTask(&req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, req)
}
```

---

## 🏛️ Architecture Notes

This module follows the **Hexagonal Architecture** (Ports & Adapters):

```
Delivery Layer (HTTP)
       ↓
Use Case Layer (Business Logic)
       ↓
Repository Layer (Data Access)
       ↓
PostgreSQL Database
```

**Key Principles:**
- Domain entities are pure Go structs with no framework dependencies
- Repository interface is the PORT (defined in domain)
- Repository implementation is the ADAPTER (in repository package)
- Use cases contain ALL business logic
- HTTP handlers are thin and delegate to use cases

---

## 📚 Related Documentation

- [Database Schema](../../../../docs/DATABASE_SCHEMA.md)
- [API Documentation](../../../../docs/API.md) (TODO)
- [Testing Guide](../../../../docs/TESTING.md) (TODO)

---

**Domain layer complete and ready for implementation! 🚀**
