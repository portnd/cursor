# 🚀 KOMGRIP QUICKSTART GUIDE

Get your God-Tier application running in 3 minutes.

> **📚 Need detailed instructions?**  
> - **🇹🇭 ภาษาไทย:** [SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md) - คู่มือละเอียดฉบับสมบูรณ์  
> - **🇬🇧 English:** [SETUP_GUIDE_EN.md](./SETUP_GUIDE_EN.md) - Complete step-by-step guide  
> - **📖 All Documentation:** [DOCUMENTATION_INDEX.md](./DOCUMENTATION_INDEX.md)

---

## ⚡ Quick Setup (First Time)

```bash
# 1. Initialize project (rename modules, create .env)
./init_project.sh
# or
make init

# 2. Start all services (Postgres, Mongo, Redis, API, Web)
make up

# 3. Wait 30 seconds for services to initialize, then test
curl http://localhost:8080/health
```

**Expected Output:**
```json
{
  "status": "UP",
  "timestamp": "2026-01-22T10:30:00+07:00",
  "services": {
    "postgres": "UP",
    "mongodb": "UP",
    "redis": "UP"
  }
}
```

**4. Access the beautiful landing page**
```bash
# Open in browser
open http://localhost:3000
```

You'll see a stunning God-Tier landing page with live system status!

---

## 🎯 Access Points

| Service | URL | Description |
| :--- | :--- | :--- |
| **Web App** | http://localhost:3000 | Nuxt 3 Frontend |
| **API** | http://localhost:8080 | Go Backend |
| **Health Check** | http://localhost:8080/health | System status |

---

## 🛠️ Common Commands

```bash
# View logs (Ctrl+C to exit)
make logs

# Stop all services
make down

# Restart services
make restart

# Access database shells
make shell-db        # PostgreSQL
make shell-mongo     # MongoDB
make shell-redis     # Redis

# Access container shells
make shell-api       # Go API container
make shell-web       # Nuxt Web container

# Clean everything (DESTRUCTIVE - deletes all data!)
make clean
```

---

## 🔥 Hot Reload (Development)

Both API and Web have hot reload enabled:

### Backend (Go + Air)
- Edit any `.go` file in `api/`
- Air automatically rebuilds and restarts
- Check logs: `make logs`

### Frontend (Nuxt)
- Edit any `.vue`, `.ts`, `.css` file in `web/`
- Nuxt HMR automatically updates
- Check browser console for updates

---

## 📁 Project Structure Overview

```
komgrip-starter/
├── api/                        # Go Backend (Hexagonal Architecture)
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/
│   │   ├── core/               # Shared utilities
│   │   │   ├── config/         # Environment config
│   │   │   └── database/       # DB connections
│   │   └── modules/            # Feature modules
│   │       └── health/         # Health check module
│   └── docker/Dockerfile.dev   # Dev container with hot reload
│
├── web/                        # Nuxt 3 Frontend (Feature-Sliced Design)
│   ├── core/modules/           # Feature modules
│   └── pages/                  # Nuxt pages (routing)
│
├── docker-compose.yml          # Multi-service orchestration
├── Makefile                    # Development commands
├── init_project.sh             # Project initialization
└── README.md                   # Full documentation
```

---

## 🗄️ Database Info

### PostgreSQL (Primary)
- **Host:** localhost:5432
- **User:** komgrip
- **Password:** komgrip_secret (change in production!)
- **Database:** komgrip_db
- **Use for:** Users, transactions, core business data

### MongoDB (Secondary)
- **Host:** localhost:27017
- **User:** komgrip
- **Password:** komgrip_secret
- **Database:** komgrip_logs
- **Use for:** Audit logs, system logs, analytics

### Redis (Cache)
- **Host:** localhost:6379
- **Password:** komgrip_secret
- **Memory Limit:** 128MB (with LRU eviction)
- **Use for:** Sessions, rate limiting, temporary cache

---

## 🧪 Testing Your Setup

```bash
# 1. Check all services are running
make ps

# 2. Test database connections
curl http://localhost:8080/health

# 3. Test PostgreSQL
make shell-db
# Inside psql: \l (list databases)

# 4. Test MongoDB
make shell-mongo
# Inside mongosh: show dbs

# 5. Test Redis
make shell-redis
# Inside redis-cli: PING
```

---

## 🚨 Troubleshooting

### Services won't start
```bash
# Check if ports are already in use
lsof -i :8080    # API port
lsof -i :3000    # Web port
lsof -i :5432    # Postgres
lsof -i :27017   # MongoDB
lsof -i :6379    # Redis

# Force cleanup and restart
make down
docker system prune -f
make up
```

### Health check shows "DEGRADED"
```bash
# Check individual service logs
docker-compose logs postgres
docker-compose logs mongo
docker-compose logs redis

# Restart specific service
docker-compose restart postgres
```

### Hot reload not working
```bash
# Backend (Air)
make shell-api
cat build-errors.log    # Check for build errors

# Frontend (Nuxt)
make shell-web
npm run dev             # Run manually to see errors
```

### Permission denied on init_project.sh
```bash
chmod +x init_project.sh
./init_project.sh
```

---

## 📖 Next Steps

1. **Read Documentation:**
   - Root: [README.md](./README.md) - Complete guide
   - API: [api/README.md](./api/README.md) - Backend architecture
   - Web: [web/README.md](./web/README.md) - Frontend guide (coming soon)

2. **Add Your First Module:**
   - See `api/README.md` section "Adding a New Module"
   - Example: User registration, wallet management

3. **Configure Production:**
   - Update all passwords in `.env`
   - Generate strong JWT secrets
   - Enable HTTPS
   - Set up monitoring

4. **Deploy:**
   - See main README.md "Deployment" section
   - Production Dockerfile included
   - Ready for Docker Swarm / Kubernetes

---

## 💡 Tips

- **Use Makefile:** All common tasks have shortcuts (`make up`, `make logs`, etc.)
- **Check Logs First:** 90% of issues are visible in logs (`make logs`)
- **Hot Reload is Magic:** No need to restart containers during development
- **Shell Access:** Use `make shell-*` commands to debug inside containers
- **Keep .env Secret:** Never commit `.env` to git (already in `.gitignore`)

---

## 📞 Support

- **CEO:** Thanandorn
- **Architecture:** Hexagonal Modular Monolith + Feature-Sliced Design
- **Issue Tracker:** [GitHub Issues](#)

---

**Happy Coding! 🛡️**

*Built for national-scale applications with zero compromises.*
