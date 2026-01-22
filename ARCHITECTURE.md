# 🏗️ KOMGRIP ARCHITECTURE OVERVIEW

Complete architectural documentation for the God-Tier Starter Kit.

---

## 📊 System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                         │
│                    (Browser / Mobile App)                    │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      │ HTTP/HTTPS
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      FRONTEND (Nuxt 3)                       │
│                  Feature-Sliced Design (FSD)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Pages      │  │   Modules    │  │   Shared     │      │
│  │  (Routing)   │→ │   (Logic)    │→ │  (Utils)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                   │                   │            │
│         └───────────────────┴───────────────────┘            │
│                             │                                │
│                      Pinia Store                             │
│                             │                                │
└─────────────────────────────┼────────────────────────────────┘
                              │
                              │ REST API
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    BACKEND (Go + Gin)                        │
│              Hexagonal Architecture (Modular)                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Delivery Layer (HTTP Handlers)                      │   │
│  │         ↓                                             │   │
│  │  Usecase Layer (Business Logic)                      │   │
│  │         ↓                                             │   │
│  │  Repository Layer (Data Access)                      │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────┬────────────────────────────────┘
                              │
                 ┌────────────┼────────────┐
                 │            │            │
                 ▼            ▼            ▼
         ┌──────────┐  ┌──────────┐  ┌──────────┐
         │PostgreSQL│  │ MongoDB  │  │  Redis   │
         │   (15)   │  │   (6)    │  │   (7)    │
         │  PRIMARY │  │LOGS/AUDIT│  │  CACHE   │
         │   ACID   │  │ SCHEMA-  │  │  128MB   │
         │  TXNS    │  │   LESS   │  │   LRU    │
         └──────────┘  └──────────┘  └──────────┘
```

---

## 🎯 Architecture Patterns

### Backend: Hexagonal Architecture (Ports & Adapters)

```
internal/modules/{feature}/
├── domain/              # Core business logic (Pure Go)
│   ├── entities.go      # Domain models
│   └── interfaces.go    # Port definitions (Repository, Usecase)
│
├── usecase/             # Business logic implementation
│   └── feature_usecase.go   # Implements domain.Usecase
│
├── repository/          # Adapter: Database
│   ├── postgres_repository.go   # Implements domain.Repository
│   └── mongo_repository.go      # Alternative implementation
│
└── delivery/            # Adapter: External interfaces
    ├── http/            # HTTP/REST adapter
    │   ├── handler.go   # Gin handlers
    │   └── route.go     # Route registration
    └── grpc/            # gRPC adapter (future)
```

**Benefits:**
- ✅ Business logic independent of frameworks
- ✅ Easy to test (mock interfaces)
- ✅ Swap implementations without changing core logic
- ✅ Clear dependency flow (inward only)

---

### Frontend: Feature-Sliced Design (FSD)

```
core/modules/{feature}/
├── infrastructure/      # External dependencies
│   ├── api.ts          # HTTP API calls
│   └── storage.ts      # localStorage, IndexedDB
│
├── store/              # State management
│   └── feature-store.ts    # Pinia store
│
└── ui/                 # UI components
    ├── FeatureCard.vue
    └── FeatureForm.vue

pages/                  # File-based routing (Nuxt)
└── feature.vue         # Dumb component (no logic)
```

**Benefits:**
- ✅ Features are isolated and portable
- ✅ Clear separation of concerns
- ✅ Easy to onboard new developers
- ✅ Scales to large teams

---

## 🗄️ Database Strategy

| Database | Use Case | Data Examples | Constraints |
| :--- | :--- | :--- | :--- |
| **PostgreSQL** | Core business data requiring ACID | Users, Wallets, Transactions, Orders, Payments | Must use transactions, enforce foreign keys |
| **MongoDB** | Logs, audits, analytics | System logs, Audit trails, User activities, Analytics events | No foreign keys, TTL indexes recommended |
| **Redis** | High-speed cache | Sessions, Rate limiting, OTP codes, Temporary tokens | **128MB limit**, LRU eviction, Always set TTL |

### Decision Tree

```
New data entity?
    │
    ├─ Need ACID guarantees? ────────────────────→ PostgreSQL
    │
    ├─ Frequent writes, flexible schema? ─────────→ MongoDB
    │
    └─ Temporary data (< 1 hour)? ────────────────→ Redis
```

---

## 🔄 Request Flow

### Example: User Login

```
1. Browser
   ↓ POST /api/auth/login
   
2. Nuxt Frontend (pages/login.vue)
   ↓ Calls LoginForm component
   
3. Feature Module (core/modules/auth/)
   ├─ UI: LoginForm.vue
   ├─ Store: authStore.login()
   └─ Infrastructure: authApi.login()
       ↓ useHttp().post('/api/auth/login', {...})
       
4. Go Backend (api/cmd/server/main.go)
   ↓ Gin Router receives request
   
5. Auth Module (internal/modules/auth/)
   ├─ Delivery: AuthHandler.Login(c *gin.Context)
   ├─ Usecase: authUsecase.Login(email, password)
   │   ├─ Validates credentials
   │   ├─ Checks password hash
   │   └─ Generates JWT token
   └─ Repository: userRepo.FindByEmail(email)
       ↓ GORM query to PostgreSQL
       
6. PostgreSQL
   ↓ Returns user record
   
7. Redis
   ↓ Stores session token (TTL: 24h)
   
8. Response flows back through layers
   
9. Frontend receives JWT
   ↓ Stores in cookie
   ↓ Updates authStore.user
   ↓ Redirects to dashboard
```

---

## 🔐 Security Layers

### Frontend
- ✅ HTTP-only cookies for tokens
- ✅ CSRF protection (Nuxt built-in)
- ✅ XSS protection (Vue sanitization)
- ✅ Input validation (TypeScript + Zod)

### Backend
- ✅ JWT authentication
- ✅ Password hashing (bcrypt)
- ✅ SQL injection protection (GORM prepared statements)
- ✅ CORS configuration
- ✅ Rate limiting (Redis)
- ✅ Request validation

### Infrastructure
- ✅ HTTPS in production
- ✅ Database encryption at rest
- ✅ Secrets management (.env)
- ✅ Container isolation (Docker)

---

## 📈 Scalability Strategy

### Horizontal Scaling

```
Load Balancer (Nginx/Caddy)
    ↓
┌───────┬───────┬───────┐
│ API 1 │ API 2 │ API 3 │  ← Scale Go instances
└───┬───┴───┬───┴───┬───┘
    │       │       │
    └───────┼───────┘
            ↓
    ┌──────────────┐
    │  PostgreSQL  │  ← Read replicas
    │   Master +   │
    │  2 Replicas  │
    └──────────────┘
```

### Vertical Scaling

| Service | CPU | Memory | Notes |
| :--- | :--- | :--- | :--- |
| Go API | 2-4 cores | 2-4 GB | Efficient, low memory |
| Nuxt SSR | 2-4 cores | 2-4 GB | Node.js memory |
| PostgreSQL | 4-8 cores | 8-16 GB | CPU for queries, RAM for cache |
| MongoDB | 2-4 cores | 4-8 GB | Document processing |
| Redis | 1-2 cores | 256 MB | Memory-focused, **128MB limit** |

---

## 🚀 Deployment Architecture

### Development
```
Docker Compose (Local)
├── PostgreSQL (port 5432)
├── MongoDB (port 27017)
├── Redis (port 6379)
├── API (port 8080) ← Hot reload with Air
└── Web (port 3000) ← Hot reload with Nuxt HMR
```

### Staging/Production
```
Cloud Provider (AWS/GCP/Azure)
├── Container Orchestration
│   ├── Kubernetes / Docker Swarm
│   └── Auto-scaling groups
├── Managed Databases
│   ├── RDS (PostgreSQL)
│   ├── DocumentDB / Atlas (MongoDB)
│   └── ElastiCache (Redis)
├── CDN (CloudFront/CloudFlare)
├── Load Balancer (ALB/NLB)
└── Monitoring (Prometheus + Grafana)
```

---

## 📊 Module Communication

### Backend Modules (Go)

```go
// Modules are independent but can communicate via:

// 1. Direct usecase injection
type WalletUsecase struct {
    walletRepo  domain.WalletRepository
    userUsecase user.UserUsecase  // ← Dependency injection
}

// 2. Event bus (future)
eventBus.Publish("user.created", userID)
```

### Frontend Modules (Nuxt)

```typescript
// Modules are independent but can communicate via:

// 1. Pinia stores
const authStore = useAuthStore()
const walletStore = useWalletStore()
walletStore.loadForUser(authStore.user.id)

// 2. Event bus (mitt)
eventBus.emit('wallet:updated', balance)
```

---

## 🔍 Monitoring & Observability

### Logs
- **Go:** Structured logging (Zap/Zerolog)
- **Nuxt:** Console + Sentry
- **Aggregation:** ELK Stack or Loki

### Metrics
- **Backend:** Prometheus + Go metrics
- **Frontend:** Web Vitals
- **Databases:** Native monitoring tools

### Tracing
- **Distributed:** OpenTelemetry
- **APM:** DataDog / New Relic

### Health Checks
- **Backend:** `GET /health` (checks all DBs)
- **Frontend:** Uptime monitoring
- **Databases:** Native health endpoints

---

## 🧪 Testing Strategy

### Backend (Go)
```
Unit Tests        → Test business logic (usecases)
Integration Tests → Test with real databases
E2E Tests        → Test full API flows
```

### Frontend (Nuxt)
```
Unit Tests        → Test composables, stores
Component Tests   → Test Vue components (Vitest)
E2E Tests        → Test user flows (Playwright)
```

---

## 📦 CI/CD Pipeline

```
Developer Push
    ↓
GitHub Actions / GitLab CI
    ├─ Lint (ESLint, golangci-lint)
    ├─ Unit Tests
    ├─ Integration Tests
    ├─ Build Docker Images
    ├─ Security Scan (Trivy)
    └─ Deploy to Staging
        ↓ Manual Approval
        Deploy to Production
```

---

## 🎓 Key Principles

### Backend (Go)
1. **Dependency Inversion:** Depend on interfaces, not implementations
2. **Single Responsibility:** Each module handles one domain
3. **Error Handling:** Always return errors, never panic
4. **Testability:** All business logic is unit-testable

### Frontend (Nuxt)
1. **Feature Independence:** Each feature is self-contained
2. **Dumb Pages:** Pages only handle routing
3. **Smart Modules:** Business logic lives in modules
4. **Type Safety:** Strict TypeScript mode

### Database
1. **Right Tool for the Job:** Use PostgreSQL for ACID, Mongo for logs, Redis for cache
2. **Transactions:** Always use for multi-step operations
3. **Migrations:** Never manually alter schema
4. **Indexes:** Add for frequently queried columns

---

## 🚧 Future Enhancements

- [ ] GraphQL API layer (optional)
- [ ] gRPC for inter-service communication
- [ ] WebSocket support for real-time features
- [ ] Background job processing (Redis + Bull)
- [ ] File upload service (S3 integration)
- [ ] Email service (SendGrid/SES)
- [ ] SMS service (Twilio)
- [ ] Multi-tenancy support
- [ ] Role-based access control (RBAC)
- [ ] API versioning strategy

---

**This architecture is designed to scale from prototype to national-level production.**

Built with 💪 by Thanandorn (Komgrip CEO)
