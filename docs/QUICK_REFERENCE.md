# 🚀 KOMGRIP STARTER - QUICK REFERENCE CARD

> พิมพ์หรือบันทึกไฟล์นี้ไว้ใช้งานบ่อยๆ

---

## 📖 เอกสารหลัก

| เอกสาร | ใช้เมื่อ |
|:-------|:---------|
| **[SETUP_GUIDE_TH.md](./SETUP_GUIDE_TH.md)** | 🆕 ติดตั้งครั้งแรก (คู่มือละเอียด) |
| **[SETUP_GUIDE_EN.md](./SETUP_GUIDE_EN.md)** | 🆕 First-time setup (English) |
| **[QUICKSTART.md](./QUICKSTART.md)** | ⚡ Quick 3-min start |
| **[DOCUMENTATION_INDEX.md](./DOCUMENTATION_INDEX.md)** | 📚 หาเอกสารอื่นๆ |

---

## ⚡ คำสั่งที่ใช้บ่อย

### การจัดการ Services

```bash
make up          # เริ่ม services ทั้งหมด
make down        # หยุด services ทั้งหมด
make restart     # Restart services
make logs        # ดู logs ทั้งหมด
make ps          # ดูสถานะ containers
make clean       # ลบทุกอย่าง (containers + volumes)
```

### เข้าถึง Database

```bash
make shell-db     # PostgreSQL shell (psql)
make shell-redis  # Redis CLI
```

### เข้าถึง Container

```bash
make shell-api    # เข้าไปใน API container
make shell-web    # เข้าไปใน Web container
```

### ดู Logs แยก Service

```bash
docker-compose logs -f api       # API logs
docker-compose logs -f web       # Web logs
docker-compose logs -f postgres  # Database logs
```

---

## 🌐 URLs

| Service | URL | Purpose |
|:--------|:----|:--------|
| **Frontend** | http://localhost:3000 | Nuxt web app |
| **API** | http://localhost:8080 | Go REST API |
| **Health** | http://localhost:8080/health | Status check |

---

## 🐛 แก้ปัญหาเร่งด่วน

### 1. Services ไม่ทำงาน
```bash
make down
docker system prune -f
make up
```

### 2. Port ถูกใช้แล้ว
```bash
# ดูว่าโปรแกรมไหนใช้ port 8080
lsof -i :8080
# ปิดโปรแกรมนั้น หรือเปลี่ยน port ใน .env
```

### 3. Database ไม่ healthy
```bash
docker-compose restart postgres
# หรือ
docker-compose restart redis
```

### 4. API ไม่เชื่อมต่อ DB
```bash
# รอ 30 วินาทีให้ DB พร้อม
sleep 30
docker-compose restart api
```

### 5. Web ไม่แสดงผล
```bash
docker-compose stop web
docker-compose rm -f web
docker-compose up -d web
# รอ 1-2 นาที
```

### 6. ไฟล์ Cache เก่า (Browser)
```bash
# Hard refresh
Mac: Cmd + Shift + R
Windows: Ctrl + Shift + R

# หรือเปิด Incognito
Mac: Cmd + Shift + N
Windows: Ctrl + Shift + N
```

### 7. Docker Memory เต็ม
```
Docker Desktop → Settings → Resources
เพิ่ม Memory เป็น 4GB+
Apply & Restart
```

---

## 🔍 ตรวจสอบระบบ

### Health Check
```bash
curl http://localhost:8080/health
```

**ผลลัพธ์ที่ต้องการ:**
```json
{
  "status": "UP",
  "services": {
    "postgres": "UP",
    "redis": "UP"
  }
}
```

### Container Status
```bash
make ps
```

**ต้องเห็น:**
- ✅ `komgrip_api` - Up
- ✅ `komgrip_db` - Up (healthy)
- ✅ `komgrip_redis` - Up (healthy)
- ✅ `komgrip_web` - Up

---

## 📂 โครงสร้างโปรเจกต์

```
komgrip-starter/
├── api/                 # Go Backend
│   ├── cmd/server/      # Main entry point
│   ├── internal/
│   │   ├── core/        # Config, DB, Middleware
│   │   └── modules/     # Feature modules
│   └── go.mod
├── web/                 # Nuxt Frontend
│   ├── pages/           # File-based routing
│   ├── core/modules/    # Feature-Sliced Design
│   └── package.json
├── docker-compose.yml   # Services orchestration
├── Makefile             # Dev commands
└── .env                 # Environment variables
```

---

## 🔑 Environment Variables

### Root `.env`
```env
PROJECT_NAME=komgrip-starter
API_PORT=8080
WEB_PORT=3000
```

### `api/.env`
```env
DB_HOST=postgres
DB_USER=komgrip
DB_PASS=<generated>
REDIS_ADDR=redis:6379
JWT_SECRET=<generated>
```

### `web/.env`
```env
NUXT_PUBLIC_API_BASE=http://localhost:8080
```

---

## 🎯 ขั้นตอนพัฒนา

### เริ่มต้น
```bash
git checkout -b feature/your-feature
# เขียนโค้ด (hot reload ทำงานอัตโนมัติ)
```

### ก่อน Commit
```bash
make test          # Run tests
make logs          # ตรวจสอบ errors
git add .
git commit -m "feat: your feature"
```

### Push
```bash
git push origin feature/your-feature
# สร้าง Pull Request บน GitHub
```

---

## 📞 ติดต่อ

- **Email:** thanandorn14@gmail.com
- **GitHub Issues:** (URL ของโปรเจกต์)
- **Documentation:** [DOCUMENTATION_INDEX.md](./DOCUMENTATION_INDEX.md)

---

## 🧠 เคล็ดลับ

1. **ใช้ `make logs` บ่อยๆ** - จะเห็น errors ทันที
2. **ตรวจสอบ `make ps` เสมอ** - ดูว่า services healthy
3. **Hard refresh browser** - Cmd/Ctrl + Shift + R
4. **Clean start ถ้าสงสัย** - `make clean && make up`
5. **อ่าน error message** - มักบอกสาเหตุชัดเจน
6. **รอให้ DB พร้อม** - ใช้เวลา 10-30 วินาทีหลัง `make up`
7. **ใช้ Incognito** - ทดสอบโดยไม่มี cache

---

**Last Updated:** 2026-01-22  
**Version:** 1.0.0

---

> 💡 **บันทึกไฟล์นี้ไว้ใน Desktop หรือพิมพ์ออกมาเพื่อดูง่ายๆ!**
