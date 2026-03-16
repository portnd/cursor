# 🛡️ Komgrip Starter Kit - Complete Setup Guide

**English Version** - Step-by-step guide for starting a new project  
**For:** Developers cloning the project for the first time  
**Time Required:** Approximately 5-10 minutes

---

## 📋 Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Clone Project from GitHub](#2-clone-project-from-github)
3. [First-Time Project Setup](#3-first-time-project-setup)
4. [Start Services](#4-start-services)
5. [Verify Everything Works](#5-verify-everything-works)
6. [Common Troubleshooting](#6-common-troubleshooting)
7. [Frequently Used Commands](#7-frequently-used-commands)

---

## 1. Prerequisites

### 1.1 Required Software

You must have these installed:

#### **Docker Desktop** (Critical!)
```bash
# Check if installed
docker --version
docker-compose --version
```

**If not installed:**
- **Mac:** Download from https://www.docker.com/products/docker-desktop
- **Windows:** Download from https://www.docker.com/products/docker-desktop
- Install and open Docker Desktop
- Verify Docker is running (🐳 icon in menu bar)

#### **Git**
```bash
# Check if installed
git --version
```

**If not installed:**
- **Mac:** `brew install git` or download from https://git-scm.com/
- **Windows:** Download from https://git-scm.com/

#### **Make** (Optional but recommended)
```bash
# Check if installed
make --version
```

**If not installed:**
- **Mac:** Pre-installed (Xcode Command Line Tools)
- **Windows:** Install `choco install make` or use Git Bash

---

### 1.2 Verify Ports Are Available

This project uses these ports:
- **3000** - Nuxt (Frontend)
- **8080** - Go API (Backend)
- **5432** - PostgreSQL
- **6379** - Redis

```bash
# Check if ports are available (Mac/Linux)
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :6379
```

If any programs are using these ports, stop them first or change ports in `.env` (see step 3.3)

---

## 2. Clone Project from GitHub

### 2.1 Clone Repository

```bash
# Navigate to your projects folder
cd ~/Documents/dev_project/

# Clone (replace URL with yours)
git clone https://github.com/portnd/komgrip-starter.git

# Enter folder
cd komgrip-starter
```

Or using SSH:
```bash
git clone git@github.com:portnd/komgrip-starter.git
cd komgrip-starter
```

### 2.2 Verify Project Structure

```bash
# View structure
ls -la
```

You should see these files:
- ✅ `docker-compose.yml`
- ✅ `Makefile`
- ✅ `init_project.sh`
- ✅ `README.md`
- ✅ `api/` (Backend folder)
- ✅ `web/` (Frontend folder)

---

## 3. First-Time Project Setup

### 3.1 Run Automated Setup Script

```bash
# Make script executable
chmod +x init_project.sh

# Run script
./init_project.sh
```

**The script will ask 2 questions:**

#### **Question 1: Project Name**
```
Enter your project name (e.g., my-awesome-project):
```

Answer: `komgrip-starter` (or your preferred name)

#### **Question 2: Go Module Path**
```
Enter your Go module path (e.g., github.com/yourname/komgrip-starter):
```

Answer: `github.com/portnd/komgrip-starter` (change to match your GitHub repo)

#### **Confirmation**
```
Proceed with initialization? [Y/n]:
```

Press: `Y` or `Enter`

### 3.2 What the Script Does Automatically

The script will:
1. ✅ Replace all Go module paths in the project
2. ✅ Create `.env` files (root, api/, web/)
3. ✅ Generate random JWT secret
4. ✅ Set database passwords
5. ✅ Remove `api/go.sum` (will regenerate)

### 3.3 Check .env Files (Optional)

```bash
# View main .env
cat .env

# View API .env
cat api/.env

# View Web .env
cat web/.env
```

**To change Ports or Passwords:**
```bash
# Edit .env file (use your preferred editor)
nano .env
# or
code .env
```

**Example changes:**
```env
# If port 8080 is in use, change to 8081
API_PORT=8081

# To change password
POSTGRES_PASSWORD=your_strong_password_here
```

---

## 4. Start Services

### 4.1 Verify Docker is Running

```bash
# Check Docker
docker ps
```

If error, Docker is not running:
- Open **Docker Desktop** and wait for it to start
- Check menu bar for 🐳 icon

### 4.2 Start All Services (First time will be slow)

```bash
# Start all services
make up
```

Or without Make:
```bash
docker-compose up -d
```

**What will happen (first time takes 3-5 minutes):**

1. **Download Docker Images:**
   - ⏳ PostgreSQL 15
   - ⏳ Redis 7
   - ⏳ Node 18 (for Web)
   - ⏳ Go 1.23 (for API)

2. **Build Docker Containers:**
   - ⏳ Build API (Go)
   - ⏳ Build Web (Nuxt) - **Takes longest (npm install)**

3. **Start Services:**
   - ✅ PostgreSQL
   - ✅ Redis
   - ✅ API (Go)
   - ✅ Web (Nuxt)

### 4.3 View Logs (Verify everything starts)

```bash
# View logs for all services
make logs

# Or
docker-compose logs -f
```

**Press `Ctrl + C` to exit logs**

---

## 5. Verify Everything Works

### 5.1 Check All Services Are Running

```bash
# View container status
make ps

# Or
docker-compose ps
```

**Expected output:**
```
NAME            STATUS
komgrip_api     Up X minutes
komgrip_db      Up X minutes (healthy)
komgrip_redis   Up X minutes (healthy)
komgrip_web     Up X minutes
```

All services must be **Up** and databases must be **healthy**

### 5.2 Test API Backend

```bash
# Test health check endpoint
curl http://localhost:8080/health
```

**Expected output:**
```json
{
  "status": "UP",
  "timestamp": "2026-01-22T10:00:00Z",
  "services": {
    "postgres": "UP",
    "redis": "UP"
  }
}
```

### 5.3 Test Web Frontend

**Open browser:**
```bash
# Mac
open http://localhost:3000

# Or manually open browser and go to
# http://localhost:3000
```

**You should see:**
- ✅ Beautiful page with gradient background (purple-black)
- ✅ Text "🛡️ KOMGRIP"
- ✅ "God-Tier Starter Kit"
- ✅ Status boxes showing databases (Postgres, Redis) in green "UP"
- ✅ 6 feature cards
- ✅ Tech stack badges at bottom

### 5.4 Test Database Connections

```bash
# Enter PostgreSQL
make shell-db
# In psql shell type:
\l
\q

# Enter Redis
make shell-redis
# In redis-cli type:
PING
exit
```

---

## 6. Common Troubleshooting

### 6.1 Problem: Port Already in Use

**Error:**
```
Error: Ports are not available: port is already allocated
```

**Solution:**
```bash
# 1. Find program using port (e.g., 8080)
lsof -i :8080

# 2. Stop that program or change port in .env
nano .env
# Change API_PORT=8081

# 3. Restart
make down
make up
```

### 6.2 Problem: Docker Not Running

**Error:**
```
Cannot connect to the Docker daemon
```

**Solution:**
1. Open **Docker Desktop**
2. Wait for Docker to start (🐳 icon in menu bar)
3. Run commands again

### 6.3 Problem: Services Not Healthy

**Error:**
```
komgrip_db | [Warning] Health check failed
```

**Solution:**
```bash
# 1. View logs for that service
docker-compose logs postgres
# or
docker-compose logs redis

# 2. Restart service
docker-compose restart postgres

# 3. If still not working, delete and recreate
make down
docker volume prune -f
make up
```

### 6.4 Problem: API Can't Connect to Database

**Error in logs:**
```
Failed to connect to postgres
```

**Solution:**
```bash
# 1. Verify databases are healthy
make ps

# 2. Wait for databases to be ready (takes 10-30 seconds)
sleep 30

# 3. Restart API
docker-compose restart api

# 4. View logs
docker-compose logs api
```

### 6.5 Problem: Web Not Displaying

**Error in browser:**
```
This site can't be reached
```

**Solution:**
```bash
# 1. Verify web container is running
docker-compose ps web

# 2. View logs
docker-compose logs web

# 3. If npm install failed, rebuild
docker-compose stop web
docker-compose rm -f web
docker-compose up -d web

# 4. Wait 1-2 minutes and try again
```

### 6.6 Problem: Old File Cache (Browser)

**Symptom:** Website displays but has errors in Console

**Solution:**
```bash
# In Browser press
Mac: Cmd + Shift + R
Windows: Ctrl + Shift + R

# Or open in Incognito
Mac: Cmd + Shift + N
Windows: Ctrl + Shift + N
```

### 6.7 Problem: Out of Memory

**Error:**
```
Docker error: out of memory
```

**Solution:**
1. Open **Docker Desktop**
2. Go to **Settings** → **Resources**
3. Increase Memory to **4GB** or more
4. Click **Apply & Restart**

---

## 7. Frequently Used Commands

### 7.1 Service Management

```bash
# Start all services
make up

# Stop all services
make down

# Restart services
make restart

# View logs
make logs

# View status
make ps
```

### 7.2 Database Access

```bash
# PostgreSQL
make shell-db

# Redis
make shell-redis
```

### 7.3 Container Access

```bash
# Enter API container
make shell-api

# Enter Web container
make shell-web
```

### 7.4 View Logs by Service

```bash
# API logs only
docker-compose logs -f api

# Web logs only
docker-compose logs -f web

# Database logs
docker-compose logs -f postgres
docker-compose logs -f redis
```

### 7.5 Clean Up (Delete everything and start fresh)

```bash
# Stop and delete containers + volumes
make clean

# Or detailed
docker-compose down -v --remove-orphans
docker system prune -a -f

# Start fresh
make up
```

---

## 8. Next Steps (After Installation)

### 8.1 Read More Documentation

- **README.md** - Project overview
- **QUICKSTART.md** - Quick start guide
- **ARCHITECTURE.md** - System architecture
- **api/README.md** - Backend guide
- **web/README.md** - Frontend guide

### 8.2 Start Developing

```bash
# Create new branch
git checkout -b feature/your-feature-name

# Start coding
# Edit files in api/ or web/

# Hot reload works automatically
# API: Air will rebuild
# Web: Nuxt HMR will reload
```

### 8.3 Testing

```bash
# Test API
make test-api

# Test Web
make test-web
```

---

## 9. Installation Checklist

Use this checklist to verify everything is ready:

### Preparation
- [ ] Docker Desktop installed
- [ ] Git installed
- [ ] Make installed (optional)
- [ ] Ports 3000, 8080, 5432, 6379 available

### Clone and Setup
- [ ] Repository cloned successfully
- [ ] `./init_project.sh` completed successfully
- [ ] `.env` files exist in all 3 locations (root, api/, web/)

### Start Services
- [ ] `make up` executed successfully
- [ ] All containers status "Up"
- [ ] All databases status "healthy"

### Verify Functionality
- [ ] `curl http://localhost:8080/health` returns {"status":"UP"}
- [ ] Open http://localhost:3000 shows Landing Page
- [ ] Page displays database status in green
- [ ] No errors in browser console

### Ready to Develop
- [ ] Hot reload works (edit files and see changes)
- [ ] Documentation read
- [ ] Project structure understood

---

## 10. Support and Help

### If you encounter problems:

1. **Check Logs:**
   ```bash
   make logs
   ```

2. **View Documentation:**
   - README.md
   - DEPLOYMENT_STATUS.md

3. **Restart Services:**
   ```bash
   make down
   make up
   ```

4. **Clean and Start Fresh:**
   ```bash
   make clean
   make up
   ```

5. **Ask Team:**
   - GitHub Issues: (Your project URL)
   - Email: thanandorn14@gmail.com

---

## Summary

You've completed these steps:

1. ✅ Installed Docker, Git, Make
2. ✅ Cloned project
3. ✅ Ran init_project.sh
4. ✅ Started services with make up
5. ✅ Verified everything works
6. ✅ Ready to develop!

**Welcome to Komgrip God-Tier Starter Kit! 🛡️🚀**

---

**Written by:** Thanandorn (Komgrip CEO)  
**Last Updated:** 2026-01-22  
**Version:** 1.0.0
