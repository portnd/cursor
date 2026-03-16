# 🛡️ KOMGRIP API (Go Backend)

Production-ready Go API built with Hexagonal Architecture (Modular Monolith).

---

## 📁 Project Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── core/                       # Shared infrastructure
│   │   ├── config/                 # Configuration management
│   │   │   └── config.go           # Environment loader
│   │   └── database/               # Database connections
│   │       └── database.go         # Postgres, Redis init
│   └── modules/                    # Business modules (Hexagonal)
│       └── health/                 # Health check module
│           └── delivery/
│               └── http/
│                   ├── health_handler.go  # Handler logic
│                   └── route.go           # Route registration
├── docker/
│   └── Dockerfile.dev              # Development container with hot reload
├── .air.toml                       # Air hot reload config
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── Dockerfile                      # Production build
└── .env.example                    # Environment variables template
```

---

## 🏗️ Architecture: Hexagonal (Ports & Adapters)

Each module follows this structure:

```
internal/modules/{feature}/
├── domain/           # Business entities & interfaces (Pure Go, no dependencies)
├── usecase/          # Business logic implementation
├── repository/       # Data persistence adapters (Postgres, Redis)
└── delivery/         # HTTP handlers, gRPC, CLI, etc.
```

### Dependency Flow
```
Delivery → Usecase → Repository → Database
   ↓          ↓           ↓
  Gin      Business    GORM
          Logic
```

---

## 🚀 Getting Started

### 1. Start Services (From Root Directory)

```bash
# From repository root
make up

# View API logs
make logs

# Access API container shell
make shell-api
```

### 2. Health Check

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "UP",
  "timestamp": "2026-01-22T10:30:00+07:00",
  "services": {
    "postgres": "UP",
    "redis": "UP"
  }
}
```

If any service is down, status will be `DEGRADED` with HTTP 503.

---

## 🛠️ Development

### Hot Reload (Air)

The development container uses [Air](https://github.com/cosmtrek/air) for automatic reloading:

- **Watch:** `.go` files in all directories
- **Exclude:** `tmp/`, `vendor/`, `*_test.go`
- **Build:** Compiles to `./tmp/main` on file change
- **Restart:** Automatically restarts the server

Edit any `.go` file and see changes instantly!

### Environment Variables

Create `.env` file (or run `make init` from root):

```bash
cp .env.example .env
```

**Key Variables:**
- `APP_ENV`: `development` or `production`
- `APP_PORT`: API server port (default: `8080`)
- `POSTGRES_*`: PostgreSQL connection details
- `REDIS_*`: Redis connection details
- `JWT_SECRET`: Secret for JWT signing (change in production!)

---

## 📦 Dependencies

| Package | Purpose | Documentation |
| :--- | :--- | :--- |
| `gin-gonic/gin` | HTTP framework | https://gin-gonic.com |
| `gorm.io/gorm` | ORM for PostgreSQL | https://gorm.io |
| `go-redis/v9` | Redis client | https://redis.uptrace.dev |
| `gin-contrib/cors` | CORS middleware | https://github.com/gin-contrib/cors |
| `godotenv` | .env loader | https://github.com/joho/godotenv |

### Install/Update Dependencies

```bash
# Inside container
make shell-api

# Download dependencies
go mod download

# Add new dependency
go get github.com/example/package

# Update all dependencies
go get -u ./...

# Tidy up
go mod tidy
```

---

## 🧪 Testing

```bash
# Run all tests
make test-api

# Or inside container
make shell-api
go test -v ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Test specific module
go test -v ./internal/modules/health/...
```

---

## 🗄️ Database Management

### PostgreSQL (Primary DB)

```bash
# Access PostgreSQL shell
make shell-db

# Inside psql
\dt          # List tables
\d users     # Describe table
SELECT * FROM users LIMIT 10;
```

### Redis (Cache)

```bash
# Access Redis CLI
make shell-redis

# Inside redis-cli
INFO memory
KEYS *
GET user:123
TTL session:abc
```

---

## 📝 Adding a New Module

Example: Adding a `user` module

### 1. Create Directory Structure

```bash
mkdir -p internal/modules/user/{domain,usecase,repository,delivery/http}
```

### 2. Define Domain (Pure Go)

**`internal/modules/user/domain/user.go`**
```go
package domain

import "time"

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Name      string    `json:"name" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
    Create(user *User) error
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
}

type UserUsecase interface {
    Register(email, name string) (*User, error)
    GetProfile(id uint) (*User, error)
}
```

### 3. Implement Repository

**`internal/modules/user/repository/postgres_repository.go`**
```go
package repository

import (
    "github.com/komgrip/starter-kit/internal/modules/user/domain"
    "gorm.io/gorm"
)

type postgresRepository struct {
    db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.UserRepository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *postgresRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *postgresRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

### 4. Implement Usecase

**`internal/modules/user/usecase/user_usecase.go`**
```go
package usecase

import (
    "errors"
    "github.com/komgrip/starter-kit/internal/modules/user/domain"
)

type userUsecase struct {
    repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
    return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(email, name string) (*domain.User, error) {
    // Check if user exists
    existing, _ := u.repo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }

    user := &domain.User{
        Email: email,
        Name:  name,
    }

    if err := u.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

func (u *userUsecase) GetProfile(id uint) (*domain.User, error) {
    return u.repo.FindByID(id)
}
```

### 5. Create HTTP Handler

**`internal/modules/user/delivery/http/user_handler.go`**
```go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/komgrip/starter-kit/internal/modules/user/domain"
)

type UserHandler struct {
    usecase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *UserHandler {
    return &UserHandler{usecase: usecase}
}

type RegisterRequest struct {
    Email string `json:"email" binding:"required,email"`
    Name  string `json:"name" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.usecase.Register(req.Email, req.Name)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := h.usecase.GetProfile(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
```

### 6. Register Routes

**`internal/modules/user/delivery/http/route.go`**
```go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/komgrip/starter-kit/internal/modules/user/domain"
)

func RegisterRoutes(router *gin.RouterGroup, usecase domain.UserUsecase) {
    handler := NewUserHandler(usecase)

    users := router.Group("/users")
    {
        users.POST("/register", handler.Register)
        users.GET("/:id", handler.GetProfile)
    }
}
```

### 7. Wire Up in `main.go`

```go
import (
    userRepo "github.com/komgrip/starter-kit/internal/modules/user/repository"
    userUsecase "github.com/komgrip/starter-kit/internal/modules/user/usecase"
    userHttp "github.com/komgrip/starter-kit/internal/modules/user/delivery/http"
)

// After database connections...
repo := userRepo.NewPostgresRepository(db)
uc := userUsecase.NewUserUsecase(repo)

// After router initialization...
api := router.Group("/api/v1")
userHttp.RegisterRoutes(api, uc)
```

---

## 🚀 Production Build

### Build Docker Image

```bash
docker build -t komgrip-api:latest -f Dockerfile .
```

### Run Production Container

```bash
docker run -d \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e POSTGRES_HOST=your-db-host \
  -e JWT_SECRET=your-secret \
  komgrip-api:latest
```

---

## 📖 Best Practices

### 1. Error Handling
- Always return errors, never panic in production code
- Use custom error types for business logic errors
- Log errors with context

### 2. Database Transactions
```go
tx := db.Begin()
if err := repo.Create(tx, user); err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

### 3. Context Propagation
- Always pass `context.Context` for cancellation
- Set timeouts for external calls
- Use context for request-scoped values

### 4. Security
- Never log passwords or secrets
- Use parameterized queries (GORM does this automatically)
- Validate all user inputs
- Implement rate limiting for public endpoints

---

## 📞 Support

For issues or questions:
- **CEO:** Thanandorn
- **Architecture:** Hexagonal Modular Monolith
- **Documentation:** See root README.md

---

**Built with 💪 for national-scale applications.**
