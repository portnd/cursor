# 🛡️ คู่มือการติดตั้งและเริ่มใช้งาน Komgrip Starter Kit

**คู่มือฉบับภาษาไทย** สำหรับการเริ่มโปรเจกต์ใหม่ตั้งแต่ต้น  
**สำหรับ:** Developer ที่จะ clone โปรเจกต์ครั้งแรก  
**เวลาที่ใช้:** ประมาณ 5-10 นาที

---

## 📋 สารบัญ

1. [สิ่งที่ต้องเตรียมก่อนเริ่ม](#1-สิ่งที่ต้องเตรียมก่อนเริ่ม)
2. [Clone โปรเจกต์จาก GitHub](#2-clone-โปรเจกต์จาก-github)
3. [ตั้งค่าโปรเจกต์ครั้งแรก](#3-ตั้งค่าโปรเจกต์ครั้งแรก)
4. [เริ่มใช้งาน Services](#4-เริ่มใช้งาน-services)
5. [ตรวจสอบการทำงาน](#5-ตรวจสอบการทำงาน)
6. [การแก้ปัญหาที่พบบ่อย](#6-การแก้ปัญหาที่พบบ่อย)
7. [คำสั่งที่ใช้บ่อย](#7-คำสั่งที่ใช้บ่อย)

---

## 1. สิ่งที่ต้องเตรียมก่อนเริ่ม

### 1.1 ตรวจสอบโปรแกรมที่จำเป็น

คุณต้องติดตั้งโปรแกรมเหล่านี้ก่อน:

#### **Docker Desktop** (จำเป็นมาก!)
```bash
# ตรวจสอบว่าติดตั้งแล้วหรือยัง
docker --version
docker-compose --version
```

**ถ้ายังไม่มี:**
- **Mac:** ดาวน์โหลดจาก https://www.docker.com/products/docker-desktop
- **Windows:** ดาวน์โหลดจาก https://www.docker.com/products/docker-desktop
- ติดตั้งและเปิดโปรแกรม Docker Desktop
- ตรวจสอบว่า Docker กำลังทำงาน (มีไอคอน 🐳 ที่ menu bar)

#### **Git**
```bash
# ตรวจสอบ
git --version
```

**ถ้ายังไม่มี:**
- **Mac:** `brew install git` หรือดาวน์โหลดจาก https://git-scm.com/
- **Windows:** ดาวน์โหลดจาก https://git-scm.com/

#### **Make** (ไม่บังคับ แต่แนะนำ)
```bash
# ตรวจสอบ
make --version
```

**ถ้ายังไม่มี:**
- **Mac:** มีติดตั้งมาแล้ว (Xcode Command Line Tools)
- **Windows:** ติดตั้ง `choco install make` หรือใช้ Git Bash

---

### 1.2 ตรวจสอบว่า Port ไม่ถูกใช้งาน

โปรเจกต์ใช้ ports เหล่านี้:
- **3000** - Nuxt (Frontend)
- **8080** - Go API (Backend)
- **5432** - PostgreSQL
- **6379** - Redis

```bash
# ตรวจสอบว่า port ว่างหรือไม่ (Mac/Linux)
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :6379
```

ถ้ามีโปรแกรมใช้ port เหล่านี้อยู่ ต้องปิดก่อน หรือเปลี่ยน port ใน `.env` (ดูขั้นตอนที่ 3.3)

---

## 2. Clone โปรเจกต์จาก GitHub

### 2.1 Clone Repository

```bash
# ไปที่ folder ที่ต้องการเก็บโปรเจกต์
cd ~/Documents/dev_project/

# Clone (เปลี่ยน URL ตามของคุณ)
git clone https://github.com/portnd/komgrip-starter.git

# เข้าไปใน folder
cd komgrip-starter
```

หรือถ้าใช้ SSH:
```bash
git clone git@github.com:portnd/komgrip-starter.git
cd komgrip-starter
```

### 2.2 ตรวจสอบโครงสร้างโปรเจกต์

```bash
# ดู structure
ls -la
```

คุณควรเห็นไฟล์เหล่านี้:
- ✅ `docker-compose.yml`
- ✅ `Makefile`
- ✅ `init_project.sh`
- ✅ `README.md`
- ✅ `api/` (โฟลเดอร์ Backend)
- ✅ `web/` (โฟลเดอร์ Frontend)

---

## 3. ตั้งค่าโปรเจกต์ครั้งแรก

### 3.1 รัน Script ตั้งค่าอัตโนมัติ

```bash
# ทำให้ script สามารถรันได้
chmod +x init_project.sh

# รัน script
./init_project.sh
```

**Script จะถามคำถาม 2 ข้อ:**

#### **คำถามที่ 1: ชื่อโปรเจกต์**
```
Enter your project name (e.g., my-awesome-project):
```

ตอบ: `komgrip-starter` (หรือชื่อที่ต้องการ)

#### **คำถามที่ 2: Go module path**
```
Enter your Go module path (e.g., github.com/yourname/komgrip-starter):
```

ตอบ: `github.com/portnd/komgrip-starter` (เปลี่ยนตาม GitHub repo ของคุณ)

#### **ยืนยัน**
```
Proceed with initialization? [Y/n]:
```

กด: `Y` หรือ `Enter`

### 3.2 สิ่งที่ Script ทำให้อัตโนมัติ

Script จะทำการ:
1. ✅ เปลี่ยน Go module path ทั้งหมดในโปรเจกต์
2. ✅ สร้างไฟล์ `.env` (root, api/, web/)
3. ✅ สร้าง JWT secret แบบสุ่ม
4. ✅ ตั้งค่า database passwords
5. ✅ ลบ `api/go.sum` (จะ regenerate ใหม่)

### 3.3 ตรวจสอบไฟล์ .env (ไม่บังคับ)

```bash
# ดูไฟล์ .env หลัก
cat .env

# ดูไฟล์ .env ของ API
cat api/.env

# ดูไฟล์ .env ของ Web
cat web/.env
```

**ถ้าต้องการเปลี่ยน Port หรือ Password:**
```bash
# แก้ไฟล์ .env (ใช้ editor ที่ชอบ)
nano .env
# หรือ
code .env
```

**ตัวอย่างการเปลี่ยน:**
```env
# ถ้า port 8080 ถูกใช้แล้ว เปลี่ยนเป็น 8081
API_PORT=8081

# ถ้าต้องการเปลี่ยน password
POSTGRES_PASSWORD=your_strong_password_here
```

---

## 4. เริ่มใช้งาน Services

### 4.1 ตรวจสอบว่า Docker กำลังทำงาน

```bash
# ตรวจสอบ Docker
docker ps
```

ถ้า error แสดงว่า Docker ยังไม่ทำงาน:
- เปิด **Docker Desktop** แล้วรอให้ขึ้น
- ดูที่ menu bar ต้องมีไอคอน 🐳

### 4.2 เริ่มทุก Services (ครั้งแรกจะช้าหน่อย)

```bash
# เริ่ม services ทั้งหมด
make up
```

หรือถ้าไม่มี Make:
```bash
docker-compose up -d
```

**สิ่งที่จะเกิดขึ้น (ครั้งแรกใช้เวลา 3-5 นาที):**

1. **ดาวน์โหลด Docker Images:**
   - ⏳ PostgreSQL 15
   - ⏳ Redis 7
   - ⏳ Node 18 (สำหรับ Web)
   - ⏳ Go 1.23 (สำหรับ API)

2. **Build Docker Containers:**
   - ⏳ Build API (Go)
   - ⏳ Build Web (Nuxt) - **ใช้เวลานานที่สุด (npm install)**

3. **Start Services:**
   - ✅ PostgreSQL
   - ✅ Redis
   - ✅ API (Go)
   - ✅ Web (Nuxt)

### 4.3 ดู Logs (ตรวจสอบว่าทุกอย่างเริ่มทำงาน)

```bash
# ดู logs ทุก services
make logs

# หรือ
docker-compose logs -f
```

**กด `Ctrl + C` เพื่อออกจาก logs**

---

## 5. ตรวจสอบการทำงาน

### 5.1 ตรวจสอบว่าทุก Service ทำงาน

```bash
# ดูสถานะ containers
make ps

# หรือ
docker-compose ps
```

**ผลลัพธ์ที่ต้องการเห็น:**
```
NAME            STATUS
komgrip_api     Up X minutes
komgrip_db      Up X minutes (healthy)
komgrip_redis   Up X minutes (healthy)
komgrip_web     Up X minutes
```

ทุก service ต้อง **Up** และ database ต้อง **healthy**

### 5.2 ทดสอบ API Backend

```bash
# ทดสอบ health check endpoint
curl http://localhost:8080/health
```

**ผลลัพธ์ที่ต้องการ:**
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

### 5.3 ทดสอบ Web Frontend

**เปิดเบราว์เซอร์:**
```bash
# Mac
open http://localhost:3000

# หรือเปิด browser เองและไปที่
# http://localhost:3000
```

**คุณควรเห็น:**
- ✅ หน้าจอสวยงามพื้นหลัง gradient (ม่วง-ดำ)
- ✅ ข้อความ "🛡️ KOMGRIP"
- ✅ "God-Tier Starter Kit"
- ✅ กล่องแสดงสถานะ database (Postgres, Redis) เป็นสีเขียว "UP"
- ✅ Feature cards 6 อัน
- ✅ Tech stack badges ด้านล่าง

### 5.4 ทดสอบการเชื่อมต่อ Database

```bash
# เข้าไปใน PostgreSQL
make shell-db
# ใน psql shell พิมพ์:
\l
\q

# เข้าไปใน Redis
make shell-redis
# ใน redis-cli พิมพ์:
PING
exit
```

---

## 6. การแก้ปัญหาที่พบบ่อย

### 6.1 ปัญหา: Port ถูกใช้งานแล้ว

**Error:**
```
Error: Ports are not available: port is already allocated
```

**วิธีแก้:**
```bash
# 1. หาโปรแกรมที่ใช้ port (เช่น 8080)
lsof -i :8080

# 2. ปิดโปรแกรมนั้น หรือเปลี่ยน port ใน .env
nano .env
# เปลี่ยน API_PORT=8081

# 3. Restart
make down
make up
```

### 6.2 ปัญหา: Docker ไม่ทำงาน

**Error:**
```
Cannot connect to the Docker daemon
```

**วิธีแก้:**
1. เปิด **Docker Desktop**
2. รอให้ Docker ทำงาน (มีไอคอน 🐳 ที่ menu bar)
3. รันคำสั่งใหม่

### 6.3 ปัญหา: Services ไม่ Healthy

**Error:**
```
komgrip_db | [Warning] Health check failed
```

**วิธีแก้:**
```bash
# 1. ดู logs ของ service นั้น
docker-compose logs postgres
# หรือ
docker-compose logs redis

# 2. Restart service นั้น
docker-compose restart postgres

# 3. ถ้ายังไม่หาย ลบและสร้างใหม่
make down
docker volume prune -f
make up
```

### 6.4 ปัญหา: API ไม่เชื่อมต่อ Database

**Error ใน logs:**
```
Failed to connect to postgres
```

**วิธีแก้:**
```bash
# 1. ตรวจสอบว่า databases healthy
make ps

# 2. รอให้ databases พร้อม (ใช้เวลา 10-30 วินาที)
sleep 30

# 3. Restart API
docker-compose restart api

# 4. ดู logs
docker-compose logs api
```

### 6.5 ปัญหา: Web ไม่แสดงผล

**Error ใน browser:**
```
This site can't be reached
```

**วิธีแก้:**
```bash
# 1. ตรวจสอบว่า web container ทำงาน
docker-compose ps web

# 2. ดู logs
docker-compose logs web

# 3. ถ้า npm install ล้มเหลว rebuild
docker-compose stop web
docker-compose rm -f web
docker-compose up -d web

# 4. รอ 1-2 นาที แล้วลองใหม่
```

### 6.6 ปัญหา: Cache ไฟล์เก่า (Browser)

**อาการ:** หน้าเว็บแสดงแต่มี error ใน Console

**วิธีแก้:**
```bash
# ใน Browser กด
Mac: Cmd + Shift + R
Windows: Ctrl + Shift + R

# หรือเปิดใน Incognito
Mac: Cmd + Shift + N
Windows: Ctrl + Shift + N
```

### 6.7 ปัญหา: Memory เต็ม

**Error:**
```
Docker error: out of memory
```

**วิธีแก้:**
1. เปิด **Docker Desktop**
2. ไปที่ **Settings** → **Resources**
3. เพิ่ม Memory เป็น **4GB** ขึ้นไป
4. คลิก **Apply & Restart**

---

## 7. คำสั่งที่ใช้บ่อย

### 7.1 การจัดการ Services

```bash
# เริ่ม services ทั้งหมด
make up

# หยุด services ทั้งหมด
make down

# Restart services
make restart

# ดู logs
make logs

# ดูสถานะ
make ps
```

### 7.2 เข้าถึง Database

```bash
# PostgreSQL
make shell-db

# Redis
make shell-redis
```

### 7.3 เข้าถึง Container

```bash
# เข้าไปใน API container
make shell-api

# เข้าไปใน Web container
make shell-web
```

### 7.4 ดู Logs แยก Service

```bash
# Logs ของ API เท่านั้น
docker-compose logs -f api

# Logs ของ Web เท่านั้น
docker-compose logs -f web

# Logs ของ Database
docker-compose logs -f postgres
docker-compose logs -f redis
```

### 7.5 Clean Up (ลบทุกอย่างและเริ่มใหม่)

```bash
# หยุดและลบ containers + volumes
make clean

# หรือแบบละเอียด
docker-compose down -v --remove-orphans
docker system prune -a -f

# เริ่มใหม่
make up
```

---

## 8. ขั้นตอนถัดไป (หลังติดตั้งเสร็จ)

### 8.1 อ่านเอกสารเพิ่มเติม

- **README.md** - ภาพรวมโปรเจกต์
- **QUICKSTART.md** - คู่มือเริ่มต้นฉบับย่อ
- **ARCHITECTURE.md** - สถาปัตยกรรมระบบ
- **api/README.md** - คู่มือ Backend
- **web/README.md** - คู่มือ Frontend

### 8.2 เริ่มพัฒนา

```bash
# สร้าง branch ใหม่
git checkout -b feature/your-feature-name

# เริ่มเขียนโค้ด
# แก้ไขไฟล์ใน api/ หรือ web/

# Hot reload จะทำงานอัตโนมัติ
# API: Air จะ rebuild
# Web: Nuxt HMR จะ reload
```

### 8.3 ทดสอบ

```bash
# ทดสอบ API
make test-api

# ทดสอบ Web
make test-web
```

---

## 9. Checklist การติดตั้ง

ใช้ checklist นี้เพื่อตรวจสอบว่าทุกอย่างพร้อมแล้ว:

### เตรียมความพร้อม
- [ ] ติดตั้ง Docker Desktop แล้ว
- [ ] ติดตั้ง Git แล้ว
- [ ] ติดตั้ง Make แล้ว (optional)
- [ ] Port 3000, 8080, 5432, 6379 ว่าง

### Clone และตั้งค่า
- [ ] Clone repository สำเร็จ
- [ ] รัน `./init_project.sh` สำเร็จ
- [ ] มีไฟล์ `.env` ครบทั้ง 3 ที่ (root, api/, web/)

### เริ่ม Services
- [ ] `make up` ทำงานสำเร็จ
- [ ] ทุก containers สถานะ "Up"
- [ ] Database ทุกตัวสถานะ "healthy"

### ทดสอบการทำงาน
- [ ] `curl http://localhost:8080/health` ได้ผลลัพธ์ {"status":"UP"}
- [ ] เปิด http://localhost:3000 เห็นหน้าจอ Landing Page
- [ ] หน้าจอแสดงสถานะ database เป็นสีเขียว
- [ ] ไม่มี error ใน browser console

### พร้อมพัฒนา
- [ ] Hot reload ทำงาน (แก้ไฟล์แล้วเห็นการเปลี่ยนแปลง)
- [ ] อ่านเอกสารแล้ว
- [ ] เข้าใจโครงสร้างโปรเจกต์

---

## 10. ติดต่อและขอความช่วยเหลือ

### หากพบปัญหา:

1. **ตรวจสอบ Logs:**
   ```bash
   make logs
   ```

2. **ดู Documentation:**
   - README.md
   - DEPLOYMENT_STATUS.md

3. **Restart Services:**
   ```bash
   make down
   make up
   ```

4. **Clean และเริ่มใหม่:**
   ```bash
   make clean
   make up
   ```

5. **ถาม Team:**
   - GitHub Issues: (URL ของโปรเจกต์)
   - Email: thanandorn14@gmail.com

---

## สรุป

คุณได้ทำตามขั้นตอนเหล่านี้แล้ว:

1. ✅ ติดตั้ง Docker, Git, Make
2. ✅ Clone โปรเจกต์
3. ✅ รัน init_project.sh
4. ✅ Start services ด้วย make up
5. ✅ ทดสอบว่าทุกอย่างทำงาน
6. ✅ พร้อมเริ่มพัฒนา!

**ยินดีต้อนรับสู่ Komgrip God-Tier Starter Kit! 🛡️🚀**

---

**เอกสารนี้เขียนโดย:** Thanandorn (Komgrip CEO)  
**อัพเดทล่าสุด:** 2026-01-22  
**เวอร์ชัน:** 1.0.0
