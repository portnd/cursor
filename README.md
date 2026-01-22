# 🛡️ KOMGRIP GOD-TIER STARTER KIT

> A production-ready, scalable monorepo starter kit built for national-level applications.
> 
> **Architecture:** Modular Monolith + Hexagonal + Feature-Sliced Design  
> **CEO:** Thanandorn

---

## 📋 Table of Contents

- [Architecture & Stack](#architecture--stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Database Strategy](#database-strategy)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

---

## 🏗️ Architecture & Stack

| Layer | Technology | Purpose | Constraints |
| :--- | :--- | :--- | :--- |
| **Backend** | Go 1.21+ (Gin Framework) | REST API, Business Logic | Modular Monolith + Hexagonal Architecture |
| **Frontend** | Nuxt 3 (Vue 3 + TypeScript) | SSR/SPA Web Application | Feature-Sliced Design (FSD) |
| **Primary DB** | PostgreSQL 15 | Core Business Data (ACID) | Users, Wallets, Transactions |
| **Secondary DB** | MongoDB 6 | Logs, Audits, Unstructured Data | System logs, audit trails |
| **Cache** | Redis 7 | Sessions, Rate Limiting, Cache | **128MB Memory Limit** |
| **Styling** | TailwindCSS 3 | Utility-first CSS | Mobile-first responsive design |
| **DevOps** | Docker Compose | Local development environment | Health checks enabled |

---

## 📁 Project Structure

```
komgrip-starter/
├── api/                                # Go Backend (Modular Monolith)
│   ├── cmd/
│   │   ├── server/                     # Main API server entry point
│   │   └── migrate/                    # Database migration tool
│   ├── internal/
│   │   ├── core/                       # Shared utilities & config
│   │   │   ├── config/                 # App configuration
│   │   │   ├── middleware/             # HTTP middleware (auth, cors, etc.)
│   │   │   ├── database/               # DB connection managers
│   │   │   └── utils/                  # Common helpers
│   │   └── modules/                    # Business modules (Hexagonal)
│   │       ├── user/
│   │       │   ├── domain/             # Entities & Interfaces (Pure Go)
│   │       │   ├── usecase/            # Business logic implementation
│   │       │   ├── repository/         # Data access layer (PostgreSQL)
│   │       │   └── delivery/           # HTTP handlers (Gin)
│   │       ├── wallet/
│   │       │   ├── domain/
│   │       │   ├── usecase/
│   │       │   ├── repository/
│   │       │   └── delivery/
│   │       └── audit/                  # Audit logs (MongoDB)
│   │           ├── domain/
│   │           ├── usecase/
│   │           ├── repository/
│   │           └── delivery/
│   ├── migrations/                     # SQL migration files
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .env
│
├── web/                                # Nuxt 3 Frontend (FSD)
│   ├── core/
│   │   ├── modules/                    # Feature modules (FSD)
│   │   │   ├── auth/
│   │   │   │   ├── infrastructure/     # API clients
│   │   │   │   ├── store/              # Pinia stores
│   │   │   │   └── ui/                 # Vue components
│   │   │   ├── wallet/
│   │   │   │   ├── infrastructure/
│   │   │   │   ├── store/
│   │   │   │   └── ui/
│   │   │   └── dashboard/
│   │   │       ├── infrastructure/
│   │   │       ├── store/
│   │   │       └── ui/
│   │   └── shared/                     # Shared UI/utils
│   │       ├── api/                    # Base HTTP service
│   │       ├── composables/            # Reusable composables
│   │       └── ui/                     # Common components
│   ├── pages/                          # Nuxt pages (routing)
│   ├── layouts/                        # Layout components
│   ├── plugins/                        # Nuxt plugins
│   ├── public/                         # Static assets
│   ├── nuxt.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   ├── Dockerfile
│   └── .env
│
├── docker-compose.yml                  # Multi-service orchestration
├── Makefile                            # Development commands
├── init_project.sh                     # Project initialization script
├── .cursorrules                        # AI coding standards
├── .gitignore
├── .env.example
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- **Docker** & **Docker Compose** (v2.0+)
- **Make** (optional, but recommended)
- **Go** 1.21+ (for local development)
- **Node.js** 18+ (for local development)

### 1. Initialize the Project

```bash
# Clone the repository
git clone <your-repo-url>
cd komgrip-starter

# Run the initialization script
chmod +x init_project.sh
./init_project.sh

# Or use Make
make init
```

The script will:
- Rename Go module paths to your project name
- Create `.env` files from examples
- Generate JWT secrets

### 2. Start All Services

```bash
# Start all services (PostgreSQL, MongoDB, Redis, API, Web)
make up

# View logs
make logs

# Check service status
make ps
```

### 3. Access the Application

- **Web App:** http://localhost:3000
- **API:** http://localhost:8080
- **API Health:** http://localhost:8080/health

---

## 🛠️ Development Workflow

### Makefile Commands

```bash
make help              # Show all available commands
make up                # Start all services (detached)
make down              # Stop all services
make restart           # Restart all services
make logs              # Show logs (Ctrl+C to exit)
make clean             # Stop and DELETE all volumes (DESTRUCTIVE)

# Shell access
make shell-db          # PostgreSQL shell (psql)
make shell-mongo       # MongoDB shell (mongosh)
make shell-redis       # Redis CLI
make shell-api         # API container shell
make shell-web         # Web container shell

# Database migrations
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations

# Testing
make test-api          # Run API tests with coverage
make test-web          # Run Web tests
```

### Backend Development (Go)

```bash
# Access API container
make shell-api

# Run tests
go test -v ./...

# Run with hot reload (using air)
air

# Format code
gofmt -w .
```

### Frontend Development (Nuxt)

```bash
# Access Web container
make shell-web

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Type check
npm run type-check
```

---

## 🗄️ Database Strategy

| Database | Use Case | Data Type | ACID Compliance |
| :--- | :--- | :--- | :--- |
| **PostgreSQL** | Core business data | Users, Wallets, Transactions, Orders | ✅ Required |
| **MongoDB** | System logs & audits | Audit trails, error logs, analytics | ❌ Not required |
| **Redis** | High-speed cache | Sessions, rate limiting, temp data | ❌ Not required |

### Redis Memory Management

- **Hard Limit:** 128MB
- **Eviction Policy:** `allkeys-lru` (Least Recently Used)
- **Optimization Tips:**
  - Use short key names
  - Set TTL on all keys
  - Prefer hashes over strings for structured data
  - Monitor memory usage: `make shell-redis` → `INFO memory`

### PostgreSQL Best Practices

- Always use **parameterized queries** (prevent SQL injection)
- Use **transactions** for multi-step operations
- Add **indexes** on frequently queried columns
- Use **migrations** for schema changes (never manual ALTER)

### MongoDB Best Practices

- Use for **write-heavy** operations (logs, audits)
- Schema-less but use **struct validation** in Go
- Add **TTL indexes** for auto-expiring documents
- Avoid complex joins (denormalize data)

---

## 🧪 Testing

### Backend Tests (Go)

```bash
# Run all tests with race detection
make test-api

# Run specific module tests
docker-compose exec api go test -v ./internal/modules/user/...

# Generate coverage report
docker-compose exec api go test -coverprofile=coverage.out ./...
docker-compose exec api go tool cover -html=coverage.out -o coverage.html
```

### Frontend Tests (Nuxt)

```bash
# Run all tests
make test-web

# Run Vitest in watch mode
docker-compose exec web npm run test:watch

# Run E2E tests (Playwright)
docker-compose exec web npm run test:e2e
```

---

## 🚀 Deployment

### Production Checklist

- [ ] Change all default passwords in `.env`
- [ ] Generate strong JWT secrets
- [ ] Enable HTTPS (reverse proxy: Nginx/Caddy)
- [ ] Set `APP_ENV=production`
- [ ] Enable database backups (PostgreSQL + MongoDB)
- [ ] Configure Redis persistence (AOF/RDB)
- [ ] Set resource limits in `docker-compose.yml`
- [ ] Enable monitoring (Prometheus + Grafana)
- [ ] Configure log aggregation (ELK/Loki)

### Docker Production Build

```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

---

## 📖 Contributing

### Code Standards

- **Go:** Follow `.cursorrules` (snake_case files, PascalCase structs)
- **TypeScript:** Strict mode enabled, no `any` types
- **Commits:** Use conventional commits (`feat:`, `fix:`, `refactor:`)
- **PRs:** Must pass tests and linting

### Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - New features
- `fix/*` - Bug fixes

---

## 📞 Support

For issues or questions:
- **CEO:** Thanandorn
- **Documentation:** [Internal Wiki](#)
- **Issue Tracker:** [GitHub Issues](#)

---

## 📄 License

Proprietary - © 2026 Komgrip. All rights reserved.

---

**Built with 💪 for national-scale applications.**
