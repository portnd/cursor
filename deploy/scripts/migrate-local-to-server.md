# ย้ายข้อมูลจาก Local ขึ้น Server (PostgreSQL)

ใช้เมื่อ server รันแล้วแต่ DB ว่าง ต้องการ copy ข้อมูลจากเครื่อง local ขึ้นไป

---

## สรุปขั้นตอน

1. **Local:** dump PostgreSQL
2. **Copy ไฟล์ขึ้น server**
3. **Server:** หยุด backend ชั่วคราว → restore PostgreSQL → เริ่ม backend ใหม่

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

จะได้ไฟล์ `backup_postgres.sql`

---

## 2. Copy ไฟล์ขึ้น Server

```bash
scp backup_postgres.sql USER@SERVER_IP:~/
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

### 3.3 เริ่ม backend + frontend ใหม่

```bash
cd ~/kg
docker compose up -d backend frontend
```

### 3.4 ลบไฟล์ backup บน server (ถ้าไม่ต้องการเก็บ)

```bash
rm ~/backup_postgres.sql
```

---

## หมายเหตุ

- **รหัสผ่าน:** ใช้ค่าจาก `.env` ของ local ตอน dump และจาก `~/kg/.env` บน server ตอน restore
- **User/DB name:** ถ้า server ใช้ user หรือชื่อ DB ไม่ตรงกับ local ให้แก้ในคำสั่ง `psql` ให้ตรงกับ server
- **Redis:** ไม่จำเป็นต้อง migrate (มักเป็น cache/session) ถ้าต้องการให้ user login ใหม่หลังย้ายก็พอ
