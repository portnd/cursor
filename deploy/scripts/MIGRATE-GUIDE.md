# คู่มือ Migrate ฐานข้อมูล (Local + Server)

เอกสารนี้อธิบายการรัน migration บน **local** และการ **migrate ไปที่ server** โดยใช้โค้ดและเครื่องมือที่มีอยู่ในโปรเจกต์แล้ว

---

## 1. โครงสร้างที่เกี่ยวข้องในโปรเจกต์

### 1.1 โปรแกรม migrate

- **ที่อยู่:** `api/cmd/migrate/main.go`
- **หน้าที่:** อ่าน config (ที่อยู่ของ Postgres + user/password/db) จาก **environment variables** แล้วต่อ Postgres แล้วรันไฟล์ `.up.sql` ในโฟลเดอร์ migrations ตามลำดับ และบันทึกว่า migration ไหนรันไปแล้วในตาราง `schema_migrations`

### 1.2 Config ที่ migrate ใช้

- **ที่อยู่:** `api/internal/core/config/config.go`
- **การโหลดค่า:** ใช้ `godotenv.Load()` (อ่านไฟล์ `.env` ใน **current working directory**) และอ่าน env vars ชุดนี้:

| ตัวแปร | ความหมาย | ค่าเริ่มต้น (ถ้าไม่ตั้ง) |
|--------|----------|--------------------------|
| `POSTGRES_HOST` | โฮสต์ของ Postgres | `localhost` |
| `POSTGRES_PORT` | พอร์ต | `5432` |
| `POSTGRES_USER` | User | `komgrip` |
| `POSTGRES_PASSWORD` | รหัสผ่าน | `komgrip_secret` |
| `POSTGRES_DB` | ชื่อ database | `komgrip_db` |

- **หมายเหตุ:** ถ้าตั้งค่าใน **environment** (เช่น `export POSTGRES_HOST=...`) ค่าจะ override ค่าในไฟล์ `.env`

### 1.3 ไฟล์ migration

- **โฟลเดอร์:** `api/databases/migrations/`
- **รูปแบบชื่อ:** `YYYYMMDDHHMMSS_ชื่อ.up.sql` และ `YYYYMMDDHHMMSS_ชื่อ.down.sql`
- **การรัน:** โปรแกรม migrate จะรันจาก **current working directory** โดยคาดว่าโฟลเดอร์ `databases/migrations` อยู่ใต้ directory ปัจจุบัน ดังนั้นต้องรันคำสั่ง migrate จาก **โฟลเดอร์ `api`** (หรือจาก root โดยให้ CWD มี path ไปที่ `api` ได้)

### 1.4 Local (docker-compose)

- **ที่อยู่:** `docker-compose.yml` (ที่ root โปรเจกต์)
- **บริการที่เกี่ยวข้อง:**
  - **postgres:** container ชื่อ `komgrip_db` (หรือตาม `PROJECT_NAME`), เปิดพอร์ต `5432` ออกมาที่เครื่องคุณ
  - **api:** mount โฟลเดอร์ `./api` ไปที่ `/app` ใน container และตั้ง env ให้ `POSTGRES_HOST=postgres` (ชื่อ service ใน docker-compose)
- **การรัน migrate บน local ตอนนี้:** ใช้คำสั่ง `make migrate-up` ซึ่งจะรัน **ภายใน container api** ดังนั้นจึงเห็น Postgres ที่ชื่อ `postgres` และอ่านไฟล์ migration จาก `/app/databases/migrations`

### 1.5 Server (deploy)

- **ที่อยู่:** โฟลเดอร์ `deploy/` — มี `docker-compose.yml` และ `.env.staging` (หรือ `.env` บน server ที่ copy ไปไว้ที่ `~/kg`)
- **บน server:** รัน `cd ~/kg && docker compose up -d` → มี container **postgres** (sentinel_db) เปิดพอร์ต 5432 ที่ server
- **ภาพรวม:** เครื่องคุณ (local) จะ **ไม่** ต่อ Postgres บน server โดยตรงได้ถ้า server ไม่ได้ expose พอร์ต 5432 ออกอินเทอร์เน็ต ดังนั้นจะใช้ **SSH tunnel** ให้เครื่องคุณส่งการเชื่อมต่อไปที่ Postgres บน server ผ่าน SSH

---

## 2. การ Migrate บน Local (แบบละเอียด)

### 2.1 สิ่งที่ต้องมี

- Docker / Docker Compose ติดตั้งแล้ว
- โปรเจกต์อยู่ที่ path เช่น `~/Documents/sentinel-core` และมีโฟลเดอร์ `api/` พร้อม `databases/migrations/`

### 2.2 ขั้นตอน

**ขั้นที่ 1:** เปิด terminal ที่ **root โปรเจกต์** (ที่มี `docker-compose.yml`)

```bash
cd /Users/komgrip/Documents/sentinel-core
```

**ขั้นที่ 2:** ให้บริการ postgres (และ api ถ้าต้องการ) รันอยู่

```bash
make up
# หรือรันแค่ postgres: docker compose up -d postgres
```

รอจน postgres สุขภาพดี (ดูได้จาก `docker compose ps`)

**ขั้นที่ 3:** รัน migration (รัน **ใน container api** เพื่อให้ใช้ network เดียวกับ postgres และเห็นโฟลเดอร์ migrations ที่ mount ไว้)

```bash
make migrate-up
```

คำสั่งนี้เทียบเท่ากับ:

```bash
docker compose exec api go run cmd/migrate/main.go up
```

- รันจาก **ภายใน container api**
- Working directory ใน container คือ `/app` (= โฟลเดอร์ `api/` ของคุณ)
- โปรแกรมจะอ่าน env ใน container เช่น `POSTGRES_HOST=postgres`, `POSTGRES_PORT=5432`, … แล้วต่อ Postgres แล้วรันทุกไฟล์ `.up.sql` ที่ยังไม่เคยรัน (ตาม `schema_migrations`)

**ขั้นที่ 4:** ดูผลลัพธ์

- ถ้าสำเร็จจะเห็นข้อความประมาณ: `✅ All migrations completed successfully!` หรือ `✅ Applied N migration(s)`
- ถ้ามี migration ใหม่จะเห็น `🚀 Applying migration: 20260312120000_add_app_settings` และ `✅ Applied: ...`

**การ rollback (ถ้าต้องการย้อน migration ล่าสุด):**

```bash
make migrate-down
# เทียบเท่า: docker compose exec api go run cmd/migrate/main.go down
```

---

## 3. การ Migrate ไปที่ Server (แบบละเอียด)

แนวคิด: **รันโปรแกรม migrate จากเครื่องคุณ (local)** แต่ให้มันต่อกับ **Postgres บน server** ผ่าน **SSH tunnel** ดังนั้นไม่ต้องเปิดพอร์ต 5432 ของ server ออกอินเทอร์เน็ต และไม่ต้อง copy โค้ดขึ้น server แค่ใช้ไฟล์ migration ที่มีอยู่บนเครื่องคุณ

### 3.1 สิ่งที่ต้องมี

- SSH เข้า server ได้ (มี key หรือ password)
- รู้ค่า **user** และ **host** (หรือ IP) ของ server (เช่นใน GitHub Actions ใช้ `secrets.USERNAME`, `secrets.HOST`)
- รู้ค่า Postgres บน server: จากไฟล์ `deploy/.env.staging` (หรือ `.env` ที่ใช้บน server ที่ `~/kg`) คือ `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` (พอร์ตบน server เป็น 5432 ใน docker)

### 3.2 ขั้นตอน

**ขั้นที่ 1:** เปิด tunnel จากเครื่องคุณไปยัง Postgres บน server

ใน terminal หนึ่ง (ทิ้งไว้เปิดตลอดที่ต้องการ migrate):

```bash
ssh -L 5433:localhost:5432 USER@SERVER_IP
```

- แทน `USER` ด้วย user ที่ใช้ SSH (เช่น `deploy`, `ubuntu`)
- แทน `SERVER_IP` ด้วย IP หรือ hostname ของ server (เช่น `kg.how.co.th`)
- ความหมาย: การเชื่อมต่อไปที่ **localhost:5433 บนเครื่องคุณ** จะถูกส่งผ่าน SSH ไปที่ **localhost:5432 บน server** (ซึ่งคือ Postgres ใน Docker)

ถ้าใช้ key:

```bash
ssh -i ~/.ssh/your_key -L 5433:localhost:5432 USER@SERVER_IP
```

หลังรันแล้ว **ไม่ต้องรันคำสั่งอื่นใน session นี้** แค่ปล่อยให้ tunnel ทำงาน (หรือถ้าต้องการให้ tunnel อยู่หลังปิด terminal ได้ใช้ `-f -N` ตามคู่มือ SSH)

**ขั้นที่ 2:** เปิด terminal อีกอันที่เครื่องคุณ แล้วไปที่โฟลเดอร์ `api`

```bash
cd /Users/komgrip/Documents/sentinel-core/api
```

**ขั้นที่ 3:** ตั้งค่า environment ให้ชี้ไปที่ Postgres ผ่าน tunnel

ใช้ค่าจาก `deploy/.env.staging` (หรือจาก `.env` บน server):

```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5433
export POSTGRES_USER=komgrip
export POSTGRES_PASSWORD=staging_pg_secret_kg_2024
export POSTGRES_DB=komgrip_db
```

- **สำคัญ:** `POSTGRES_PORT=5433` เพราะเรา forward พอร์ต 5433 บนเครื่องคุณไปยัง 5432 บน server
- `POSTGRES_PASSWORD` ต้องตรงกับที่ตั้งบน server

**ขั้นที่ 4:** รัน migration (จากเครื่องคุณ โดยใช้โฟลเดอร์ `api` ปัจจุบัน = มี `databases/migrations`)

```bash
go run cmd/migrate/main.go up
```

- โปรแกรมจะโหลด config (จาก env ด้านบน และอาจจาก `api/.env` ด้วย แต่ env ที่ export ไว้จะ override)
- ต่อ Postgres ที่ localhost:5433 → ผ่าน tunnel ไปที่ Postgres บน server
- อ่านไฟล์ใน `api/databases/migrations/` แล้วรันเฉพาะ migration ที่ยังไม่เคยรันบน server

**ขั้นที่ 5:** ดูผลลัพธ์

- เหมือนกับ local: จะเห็น `✅ All migrations completed successfully!` หรือรายการ migration ที่ถูก apply
- ถ้า error เช็ค: tunnel ยังเปิดอยู่ไหม, user/password/db ตรงกับบน server ไหม, server เปิด Postgres ที่ localhost:5432 ไหม

**สรุปคำสั่งรวม (ใน terminal ที่ 2):**

```bash
cd /Users/komgrip/Documents/sentinel-core/api

export POSTGRES_HOST=localhost
export POSTGRES_PORT=5433
export POSTGRES_USER=komgrip
export POSTGRES_PASSWORD=staging_pg_secret_kg_2024
export POSTGRES_DB=komgrip_db

go run cmd/migrate/main.go up
```

(ถ้าไม่มี Go บนเครื่อง ให้ใช้วิธีในหัวข้อ 4 แทน หรือรันผ่าน Docker ด้านล่าง)

### 3.3 ถ้าเครื่องคุณไม่มี Go — รัน migrate ผ่าน Docker (ยังใช้ tunnel เหมือนเดิม)

จาก **root โปรเจกต์** (มีโฟลเดอร์ `api/`):

```bash
cd /Users/komgrip/Documents/sentinel-core

docker run --rm \
  -v "$(pwd)/api:/app" \
  -w /app \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_PORT=5433 \
  -e POSTGRES_USER=komgrip \
  -e POSTGRES_PASSWORD=staging_pg_secret_kg_2024 \
  -e POSTGRES_DB=komgrip_db \
  golang:1.23-alpine \
  sh -c "go run cmd/migrate/main.go up"
```

- บน Mac/Windows: `host.docker.internal` คือเครื่อง host ดังนั้นพอร์ต 5433 ที่คุณ forward ไว้จะอยู่ที่ `host.docker.internal:5433`
- บน Linux: บางทีต้องใช้ `--add-host=host.docker.internal:host-gateway` หรือใช้ IP ของ host แทน

---

## 4. ทางเลือก: รัน SQL เองบน Server (ไม่ใช้ tunnel)

ถ้าไม่สะดวกใช้ tunnel หรือต้องการรัน migration เดียว (เช่น `app_settings`) ด้วยมือ บน server ทำดังนี้

**ขั้นที่ 1:** SSH เข้า server

```bash
ssh USER@SERVER_IP
```

**ขั้นที่ 2:** ไปที่โฟลเดอร์ deploy และรัน psql ใน container postgres

```bash
cd ~/kg
docker compose exec postgres psql -U komgrip -d komgrip_db
```

(ถ้าใน `.env` ใช้ชื่อ user/db อื่น ให้เปลี่ยน `-U` และ `-d` ให้ตรง)

**ขั้นที่ 3:** ใน psql ให้รัน SQL ของ migration นั้น

ตัวอย่างสำหรับ migration `20260312120000_add_app_settings`:

```sql
CREATE TABLE IF NOT EXISTS app_settings (
    key VARCHAR(64) PRIMARY KEY,
    value TEXT NOT NULL
);
INSERT INTO app_settings (key, value) VALUES ('teams_feature_enabled', 'true')
ON CONFLICT (key) DO NOTHING;

INSERT INTO schema_migrations (version) VALUES ('20260312120000_add_app_settings')
ON CONFLICT (version) DO NOTHING;
```

(ตาราง `schema_migrations` ในโปรเจกต์ใช้คอลัมน์ `version` แบบ UNIQUE ดังนั้นใช้ ON CONFLICT ได้)

**ขั้นที่ 4:** ออกจาก psql

```sql
\q
```

วิธีนี้เหมาะกับ migration เดียวหรือกรณีฉุกเฉิน ถ้ามีหลายไฟล์ migration หลายไฟล์ แนะนำใช้วิธี tunnel + `go run cmd/migrate/main.go up` (หัวข้อ 3) จะได้ลำดับและประวัติครบ

---

## 5. สรุปสั้น ๆ

| งาน | คำสั่ง / วิธี |
|-----|----------------|
| **Migrate บน local** | `make up` แล้ว `make migrate-up` (จาก root โปรเจกต์) |
| **Migrate ไป server** | 1) เปิด SSH tunnel: `ssh -L 5433:localhost:5432 USER@SERVER` 2) ใน terminal อื่น: `cd api` แล้วตั้ง `POSTGRES_HOST=localhost`, `POSTGRES_PORT=5433` และค่าอื่นจาก server 3) `go run cmd/migrate/main.go up` |
| **Rollback บน local** | `make migrate-down` |
| **รัน SQL เดียวบน server** | `docker compose exec postgres psql -U ... -d ...` แล้ววาง SQL ของ migration นั้น + insert ใน `schema_migrations` |

ถ้าต้องการให้คราวหน้า deploy แล้วรัน migrate อัตโนมัติบน server ได้ สามารถเพิ่ม step ใน GitHub Actions (หรือ script บน server) ให้รัน migration หลัง `docker compose up -d` ได้ โดยต้องออกแบบให้ container หรือคำสั่งที่รัน migrate เห็นโฟลเดอร์ migrations และต่อ Postgres ของ server ได้
