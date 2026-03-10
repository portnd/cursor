# ย้ายข้อมูลจาก Local ขึ้น Server (PostgreSQL + MongoDB)

ใช้เมื่อ server รันแล้วแต่ DB ว่าง ต้องการ copy ข้อมูลจากเครื่อง local ขึ้นไป

---

## สรุปขั้นตอน

1. **Local:** dump PostgreSQL + MongoDB
2. **Copy ไฟล์ขึ้น server**
3. **Server:** หยุด backend ชั่วคราว → restore ทั้งสอง DB → เริ่ม backend ใหม่

---

## 1. บนเครื่อง Local (ในโฟลเดอร์โปรเจกต์)

### 1.1 Dump PostgreSQL

ถ้ารันด้วย Docker (docker-compose ชื่อ container เป็น `komgrip_db` หรือดูด้วย `docker ps`):

```bash
# ดูชื่อ container: docker ps | grep postgres
docker exec komgrip_db pg_dump -U komgrip komgrip_db --no-owner --no-acl > backup_postgres.sql
```

ถ้า Postgres ติดตั้งบนเครื่องโดยตรง:

```bash
pg_dump -h localhost -U komgrip komgrip_db --no-owner --no-acl > backup_postgres.sql
```

### 1.2 Dump MongoDB

ถ้ารันด้วย Docker (container ชื่อ `komgrip_mongo`):

```bash
# ใส่รหัส MONGO_PASSWORD ตาม .env ของคุณ
docker exec komgrip_mongo mongodump \
  --username komgrip \
  --password "komgrip_secret" \
  --authenticationDatabase admin \
  --db komgrip_logs \
  --out /tmp/mongo_dump

docker cp komgrip_mongo:/tmp/mongo_dump ./mongo_dump
```

ถ้า Mongo ติดตั้งบนเครื่องโดยตรง:

```bash
mongodump --host localhost --port 27017 \
  -u komgrip -p "komgrip_secret" \
  --authenticationDatabase admin \
  --db komgrip_logs \
  --out ./mongo_dump
```

จะได้โฟลเดอร์ `backup_postgres.sql` และ `mongo_dump/` (หรือ `mongo_dump/komgrip_logs/`)

---

## 2. Copy ไฟล์ขึ้น Server

```bash
scp backup_postgres.sql USER@SERVER_IP:~/
scp -r mongo_dump USER@SERVER_IP:~/
```

แทน `USER` และ `SERVER_IP` ด้วย user กับ IP ของ server

---

## 3. บน Server

### 3.1 หยุด backend (และ frontend ถ้าต้องการให้ไม่มีการเขียน DB)

```bash
cd ~/kg
docker compose stop backend frontend
```

### 3.2 Restore PostgreSQL

```bash
# ส่งไฟล์เข้า container แล้ว restore (แทน komgrip/komgrip_db ตาม .env บน server)
docker exec -i sentinel_db psql -U komgrip -d komgrip_db < ~/backup_postgres.sql
```

ถ้าใช้ user/db อื่น ให้แก้ `-U` และ `-d` ตามค่าใน ~/kg/.env (POSTGRES_USER, POSTGRES_DB)

### 3.3 Restore MongoDB

```bash
# copy โฟลเดอร์ dump เข้า container
docker cp ~/mongo_dump sentinel_mongo:/tmp/mongo_dump

# restore (แทน komgrip / รหัส / komgrip_logs ตาม .env บน server)
docker exec sentinel_mongo mongorestore \
  --username komgrip \
  --password "YOUR_MONGO_PASSWORD" \
  --authenticationDatabase admin \
  --db komgrip_logs \
  --drop \
  /tmp/mongo_dump/komgrip_logs
```

`--drop` = ลบ collection เดิมก่อน restore (ถ้าไม่ใช้ ข้อมูลเก่าจะยังอยู่ อาจซ้ำ key)

### 3.4 เริ่ม backend + frontend ใหม่

```bash
cd ~/kg
docker compose up -d backend frontend
```

### 3.5 ลบไฟล์ backup บน server (ถ้าไม่ต้องการเก็บ)

```bash
rm ~/backup_postgres.sql
rm -rf ~/mongo_dump
docker exec sentinel_mongo rm -rf /tmp/mongo_dump
```

---

## หมายเหตุ

- **รหัสผ่าน:** ใช้ค่าจาก `.env` ของ local ตอน dump และจาก `~/kg/.env` บน server ตอน restore
- **User/DB name:** ถ้า server ใช้ user หรือชื่อ DB ไม่ตรงกับ local ให้แก้ในคำสั่ง `psql` / `mongorestore` ให้ตรงกับ server
- **Redis:** ไม่จำเป็นต้อง migrate (มักเป็น cache/session) ถ้าต้องการให้ user login ใหม่หลังย้ายก็พอ
