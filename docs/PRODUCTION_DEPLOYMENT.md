# 🚀 KOMGRIP - PRODUCTION DEPLOYMENT GUIDE

## 📋 Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Pre-Deployment Checklist](#pre-deployment-checklist)
4. [Build Production Images](#build-production-images)
5. [Deploy to Production](#deploy-to-production)
6. [Post-Deployment Verification](#post-deployment-verification)
7. [Monitoring & Maintenance](#monitoring--maintenance)
8. [Rollback Procedures](#rollback-procedures)

---

## 🎯 Overview

This guide covers deploying the Komgrip Starter Kit to production using **Production-Optimized Dockerfiles**.

**Key Features:**
- ✅ **Multi-stage builds** (minimal image sizes)
- ✅ **Non-root users** (enhanced security)
- ✅ **Health checks** (automatic recovery)
- ✅ **Resource limits** (prevent resource exhaustion)
- ✅ **Alpine-based** images (smaller attack surface)
- ✅ **Node 20 LTS** (frontend - latest stable)
- ✅ **Go 1.23** (backend - latest stable)
- ✅ **npm ci** (100% reproducible builds)

---

## 📦 Prerequisites

### Required Software
- **Docker** 20.10+ & **Docker Compose** v2.0+
- **Git** (for version control)
- **SSL Certificates** (Let's Encrypt recommended)
- **Reverse Proxy** (Nginx or Caddy)

### Server Requirements (Minimum)
- **CPU:** 2 cores (4 cores recommended)
- **RAM:** 4 GB (8 GB recommended)
- **Storage:** 20 GB SSD
- **OS:** Ubuntu 20.04+ or Debian 11+

### Firewall Rules
```bash
# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow SSH (change port if custom)
sudo ufw allow 22/tcp

# Enable firewall
sudo ufw enable
```

---

## ✅ Pre-Deployment Checklist

### 1. **Security Configuration**

```bash
# Generate strong JWT secret (64 characters)
openssl rand -base64 64

# Generate database passwords (32 characters)
openssl rand -base64 32
```

### 2. **Create Production Environment File**

Create `.env.prod` with the following structure:

```bash
# Project
PROJECT_NAME=komgrip
VERSION=1.0.0

# API
API_PORT=8080
JWT_SECRET=<64-char-random-string>
CORS_ORIGINS=https://yourdomain.com

# Frontend
WEB_PORT=3000
NUXT_PUBLIC_API_BASE=https://api.yourdomain.com

# PostgreSQL
POSTGRES_USER=komgrip_prod
POSTGRES_PASSWORD=<32-char-random-string>
POSTGRES_DB=komgrip_prod_db

# MongoDB
MONGO_USER=komgrip_admin
MONGO_PASSWORD=<32-char-random-string>
MONGO_DB=komgrip_logs_prod

# Redis
REDIS_PASSWORD=<32-char-random-string>
```

### 3. **Update Application Configuration**

```bash
# Update CORS origins in production
# Edit api/.env or docker-compose.prod.yml
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

---

## 🏗️ Build Production Images

### **Step 1: Navigate to Project Root**

```bash
cd /path/to/komgrip-starter
```

### **Step 2: Build Backend Image**

```bash
# Build Go API image
docker build -f api/docker/Dockerfile.prod -t komgrip-api:1.0.0 ./api

# Tag as latest
docker tag komgrip-api:1.0.0 komgrip-api:latest

# Verify image size (should be ~20-30 MB)
docker images | grep komgrip-api
```

**Expected Output:**
```
komgrip-api    1.0.0    <image-id>    2 minutes ago    25.4MB
komgrip-api    latest   <image-id>    2 minutes ago    25.4MB
```

### **Step 3: Build Frontend Image**

```bash
# Build Nuxt web image
docker build -f web/docker/Dockerfile.prod -t komgrip-web:1.0.0 ./web

# Tag as latest
docker tag komgrip-web:1.0.0 komgrip-web:latest

# Verify image size (should be ~150-200 MB)
docker images | grep komgrip-web
```

**Expected Output:**
```
komgrip-web    1.0.0    <image-id>    5 minutes ago    178MB
komgrip-web    latest   <image-id>    5 minutes ago    178MB
```

### **Step 4: Test Images Locally (Optional)**

```bash
# Test API
docker run -d -p 8080:8080 \
  -e APP_ENV=production \
  -e JWT_SECRET=test-secret \
  --name test-api \
  komgrip-api:latest

# Check health
curl http://localhost:8080/health

# Clean up
docker stop test-api && docker rm test-api
```

---

## 🚀 Deploy to Production

### **Option 1: Docker Compose (Recommended for Single Server)**

```bash
# 1. Copy production compose file
cp docker-compose.prod.yml docker-compose.yml

# 2. Build all images
docker-compose build

# 3. Start services
docker-compose --env-file .env.prod up -d

# 4. Check status
docker-compose ps

# 5. View logs
docker-compose logs -f
```

### **Option 2: Docker Swarm (For Cluster)**

```bash
# 1. Initialize Swarm
docker swarm init

# 2. Create secrets
echo "your-jwt-secret" | docker secret create jwt_secret -
echo "your-db-password" | docker secret create db_password -

# 3. Deploy stack
docker stack deploy -c docker-compose.prod.yml komgrip

# 4. Check services
docker service ls
```

### **Option 3: Kubernetes (For Large Scale)**

See `k8s/` directory for Kubernetes manifests (to be created).

---

## ✅ Post-Deployment Verification

### **1. Health Checks**

```bash
# Check API health
curl https://api.yourdomain.com/health

# Expected response:
# {
#   "status": "UP",
#   "timestamp": "2026-01-23T15:00:00Z",
#   "services": {
#     "postgres": "UP",
#     "mongodb": "UP",
#     "redis": "UP"
#   }
# }

# Check Web health
curl https://yourdomain.com

# Should return HTML page
```

### **2. Database Connectivity**

```bash
# Test PostgreSQL
docker exec komgrip_db_prod psql -U komgrip_prod -d komgrip_prod_db -c "SELECT 1;"

# Test MongoDB
docker exec komgrip_mongo_prod mongosh --eval "db.runCommand({ping: 1})"

# Test Redis
docker exec komgrip_redis_prod redis-cli -a YOUR_REDIS_PASSWORD ping
```

### **3. Load Testing (Optional)**

```bash
# Install Apache Bench
sudo apt install apache2-utils

# Test API endpoint
ab -n 1000 -c 10 https://api.yourdomain.com/health

# Test Web endpoint
ab -n 1000 -c 10 https://yourdomain.com/
```

---

## 📊 Monitoring & Maintenance

### **1. View Logs**

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api

# Last 100 lines
docker-compose logs --tail=100 api
```

### **2. Resource Usage**

```bash
# Container stats
docker stats

# Disk usage
docker system df

# Clean up unused images
docker system prune -a
```

### **3. Database Backups**

```bash
# Backup PostgreSQL
docker exec komgrip_db_prod pg_dump -U komgrip_prod komgrip_prod_db > backup_$(date +%Y%m%d).sql

# Backup MongoDB
docker exec komgrip_mongo_prod mongodump --out=/backup

# Automate with cron
0 2 * * * /path/to/backup-script.sh
```

### **4. SSL Certificate Renewal (Let's Encrypt)**

```bash
# Install Certbot
sudo apt install certbot

# Renew certificates
sudo certbot renew

# Automate renewal
0 3 * * * certbot renew --quiet
```

---

## 🔄 Rollback Procedures

### **Quick Rollback to Previous Version**

```bash
# 1. Stop current services
docker-compose down

# 2. Revert to previous image version
docker tag komgrip-api:1.0.0 komgrip-api:latest
docker tag komgrip-web:1.0.0 komgrip-web:latest

# 3. Restart services
docker-compose up -d
```

### **Database Rollback**

```bash
# Restore PostgreSQL backup
docker exec -i komgrip_db_prod psql -U komgrip_prod komgrip_prod_db < backup_20260123.sql
```

---

## 🛡️ Security Best Practices

1. **Always use HTTPS** (with valid SSL certificates)
2. **Change default passwords** for all services
3. **Limit container resources** (prevent DoS attacks)
4. **Run containers as non-root** (already configured)
5. **Keep images updated** (patch vulnerabilities)
6. **Enable audit logging** (track all changes)
7. **Use secrets management** (HashiCorp Vault, AWS Secrets Manager)
8. **Implement rate limiting** (prevent brute force)
9. **Regular security audits** (penetration testing)
10. **Backup everything** (databases, configs, SSL certs)

---

## 📞 Support & Troubleshooting

### **Common Issues**

**Issue:** Container keeps restarting
```bash
# Check logs
docker logs komgrip_api_prod

# Check health status
docker inspect komgrip_api_prod | grep -A 10 Health
```

**Issue:** Database connection failed
```bash
# Check if database is running
docker ps | grep postgres

# Check database logs
docker logs komgrip_db_prod
```

**Issue:** Out of disk space
```bash
# Clean up old images
docker image prune -a

# Clean up volumes
docker volume prune
```

---

## 📝 Production Checklist Summary

- [ ] All default passwords changed
- [ ] JWT secret generated (64+ characters)
- [ ] CORS origins configured
- [ ] SSL certificates installed
- [ ] Reverse proxy configured
- [ ] Health checks passing
- [ ] Database backups automated
- [ ] Monitoring set up
- [ ] Log aggregation configured
- [ ] Firewall rules applied
- [ ] Resource limits set
- [ ] Load testing completed
- [ ] Rollback procedure tested
- [ ] Documentation updated

---

**Built with 💪 by the Komgrip Team for Production-Ready National-Scale Applications.**

© 2026 Komgrip. All rights reserved.
