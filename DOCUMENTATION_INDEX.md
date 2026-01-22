# 📚 Komgrip Starter Kit - Documentation Index

Complete guide to all documentation files in this project.

---

## 🚀 Getting Started (Start Here!)

### For First-Time Setup

Choose your language preference:

| Document | Language | Purpose | Time Required |
|:---------|:---------|:--------|:--------------|
| **[SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md)** | 🇹🇭 Thai | คู่มือติดตั้งและเริ่มใช้งานฉบับสมบูรณ์ | 5-10 นาที |
| **[SETUP_GUIDE_EN.md](./SETUP_GUIDE_EN.md)** | 🇬🇧 English | Complete installation and setup guide | 5-10 minutes |
| **[QUICKSTART.md](./QUICKSTART.md)** | 🇬🇧 English | Quick 3-minute start for experienced developers | 3 นาที |

**Recommendation:** 
- **New to the project?** → Start with `SETUP_GUIDE_TH.md` or `SETUP_GUIDE_EN.md`
- **Experienced developer?** → Use `QUICKSTART.md`

---

## 📖 Core Documentation

### Overview & Architecture

| Document | Purpose | Audience |
|:---------|:--------|:---------|
| **[README.md](./README.md)** | Project overview, tech stack, basic usage | Everyone |
| **[ARCHITECTURE.md](./ARCHITECTURE.md)** | System architecture, design decisions, patterns | Architects, Senior Developers |
| **[DEPLOYMENT_STATUS.md](./DEPLOYMENT_STATUS.md)** | Current deployment status, running services | DevOps, Developers |

---

## 🔧 Technical Guides

### Backend (Go API)

| Document | Purpose | Audience |
|:---------|:--------|:---------|
| **[api/README.md](./api/README.md)** | Go backend architecture, how to add modules | Backend Developers |
| **[api/internal/modules/*/README.md](./api/internal/modules/)** | Module-specific documentation | Backend Developers |

**Key Topics Covered:**
- Hexagonal Architecture implementation
- How to create new modules
- Database integration patterns
- API endpoint conventions
- Error handling strategies

### Frontend (Nuxt 3)

| Document | Purpose | Audience |
|:---------|:--------|:---------|
| **[web/README.md](./web/README.md)** | Nuxt 3 frontend architecture, FSD structure | Frontend Developers |
| **[web/core/modules/*/README.md](./web/core/modules/)** | Module-specific documentation | Frontend Developers |

**Key Topics Covered:**
- Feature-Sliced Design (FSD) implementation
- How to create new features
- State management with Pinia
- API integration patterns
- Component design guidelines

---

## 🛠️ Development Guides

### Setup & Configuration

| File | Purpose |
|:-----|:--------|
| **[.cursorrules](./.cursorrules)** | AI assistant coding standards (for Cursor IDE) |
| **[docker-compose.yml](./docker-compose.yml)** | Local development environment setup |
| **[Makefile](./Makefile)** | Development workflow commands |
| **[init_project.sh](./init_project.sh)** | Project initialization script |

### Environment Variables

| File | Purpose |
|:-----|:--------|
| **[.env.example](./.env.example)** | Root environment variables template |
| **[api/.env.example](./api/.env.example)** | Backend environment variables template |
| **[web/.env.example](./web/.env.example)** | Frontend environment variables template |

---

## 🎯 Quick Reference

### Common Tasks

| Task | Documentation | Command |
|:-----|:--------------|:--------|
| **First-time setup** | [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) | `./init_project.sh && make up` |
| **Start services** | [QUICKSTART.md](./QUICKSTART.md) | `make up` |
| **Stop services** | [QUICKSTART.md](./QUICKSTART.md) | `make down` |
| **View logs** | [QUICKSTART.md](./QUICKSTART.md) | `make logs` |
| **Add backend module** | [api/README.md](./api/README.md) | Follow guide |
| **Add frontend feature** | [web/README.md](./web/README.md) | Follow guide |
| **Check health** | [DEPLOYMENT_STATUS.md](./DEPLOYMENT_STATUS.md) | `curl localhost:8080/health` |

---

## 📝 Cheat Sheets

### Makefile Commands (Quick Reference)

```bash
# Service Management
make up              # Start all services
make down            # Stop all services
make restart         # Restart all services
make logs            # View all logs
make ps              # Show service status

# Database Access
make shell-db        # PostgreSQL shell
make shell-mongo     # MongoDB shell
make shell-redis     # Redis CLI

# Container Access
make shell-api       # Enter API container
make shell-web       # Enter Web container

# Cleanup
make clean           # Stop and remove all containers + volumes
```

### Docker Compose Commands (Alternative)

```bash
# If you don't have Make installed
docker-compose up -d                    # Start
docker-compose down                     # Stop
docker-compose logs -f                  # Logs
docker-compose ps                       # Status
docker-compose restart api              # Restart specific service
docker-compose exec api sh              # Enter container
```

---

## 🌐 URLs (After Setup)

| Service | URL | Purpose |
|:--------|:----|:--------|
| **Frontend** | http://localhost:3000 | Nuxt 3 web application |
| **Backend API** | http://localhost:8080 | Go REST API |
| **Health Check** | http://localhost:8080/health | API health status |
| **PostgreSQL** | localhost:5432 | Primary database |
| **MongoDB** | localhost:27017 | Secondary database (logs) |
| **Redis** | localhost:6379 | Cache database |

---

## 🎓 Learning Path

### For New Team Members

**Day 1: Setup & Exploration**
1. Read [README.md](./README.md) - Project overview
2. Follow [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) - Complete setup
3. Test all services work
4. Browse the codebase

**Day 2: Architecture Understanding**
1. Read [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
2. Read [api/README.md](./api/README.md) - Backend architecture
3. Read [web/README.md](./web/README.md) - Frontend architecture
4. Understand module structure

**Day 3: First Feature**
1. Pick a simple task
2. Create new branch
3. Follow architecture guidelines
4. Submit PR for review

### For Experienced Developers

**Quick Start** (< 30 minutes)
1. Skim [README.md](./README.md)
2. Run [QUICKSTART.md](./QUICKSTART.md) steps
3. Check [ARCHITECTURE.md](./ARCHITECTURE.md) for patterns
4. Start coding!

---

## 🔍 Finding What You Need

### "I want to..."

| Goal | Documentation |
|:-----|:--------------|
| **Set up the project from scratch** | [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) or [SETUP_GUIDE_EN.md](./SETUP_GUIDE_EN.md) |
| **Start coding quickly** | [QUICKSTART.md](./QUICKSTART.md) |
| **Understand the architecture** | [ARCHITECTURE.md](./ARCHITECTURE.md) |
| **Add a new API endpoint** | [api/README.md](./api/README.md) |
| **Add a new page/feature** | [web/README.md](./web/README.md) |
| **Connect to databases** | [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) → Section 5.4 |
| **Troubleshoot issues** | [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) → Section 6 |
| **Learn Makefile commands** | [Makefile](./Makefile) or this document → Cheat Sheets |
| **Check service status** | [DEPLOYMENT_STATUS.md](./DEPLOYMENT_STATUS.md) |
| **Understand coding rules (for AI)** | [.cursorrules](./.cursorrules) |

---

## 📞 Support

### Getting Help

**Order of escalation:**
1. **Search documentation** using this index
2. **Check troubleshooting** in [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) → Section 6
3. **View logs** with `make logs`
4. **Ask team** in Slack/Discord/Teams
5. **Create GitHub Issue** with full error details
6. **Email** thanandorn14@gmail.com for urgent matters

### Reporting Issues

When reporting issues, include:
- [ ] What you were trying to do
- [ ] Steps to reproduce
- [ ] Error message (full output)
- [ ] Output of `make ps`
- [ ] Output of `docker-compose logs [service]`
- [ ] Your OS and Docker version

---

## 🔄 Documentation Updates

### When to Update Documentation

**You should update docs when:**
- Adding new features or modules
- Changing architecture patterns
- Adding new environment variables
- Changing deployment processes
- Finding common issues/solutions

### How to Update

1. Edit the relevant `.md` file
2. Follow markdown best practices
3. Update this index if adding new files
4. Submit PR with clear description

---

## ✅ Documentation Checklist

Use this to verify you've read the right docs:

### Initial Setup
- [ ] Read [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) or [SETUP_GUIDE_EN.md](./SETUP_GUIDE_EN.md)
- [ ] Completed all setup steps
- [ ] Verified all services work
- [ ] Understand the project structure

### Before Coding
- [ ] Read [ARCHITECTURE.md](./ARCHITECTURE.md)
- [ ] Read [api/README.md](./api/README.md) or [web/README.md](./web/README.md) (your focus area)
- [ ] Understand the architectural patterns
- [ ] Read [.cursorrules](./.cursorrules) if using Cursor IDE

### Ready to Code
- [ ] Know how to add modules/features
- [ ] Understand database strategy
- [ ] Know coding conventions
- [ ] Know how to run tests
- [ ] Know how to view logs

---

## 📊 Document Status

| Document | Status | Last Updated | Next Review |
|:---------|:-------|:-------------|:------------|
| README.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| SETUP_GUIDE_TH.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| SETUP_GUIDE_EN.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| QUICKSTART.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| ARCHITECTURE.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| DEPLOYMENT_STATUS.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| api/README.md | ✅ Current | 2026-01-22 | 2026-02-22 |
| web/README.md | ✅ Current | 2026-01-22 | 2026-02-22 |

---

**This index is maintained by:** Thanandorn (Komgrip CEO)  
**Last Updated:** 2026-01-22  
**Version:** 1.0.0

---

> 💡 **Pro Tip:** Bookmark this page! It's your map to all documentation in this project.
