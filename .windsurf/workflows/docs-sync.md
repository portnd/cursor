---
description: Docs Sync - สร้างและอัปเดตเอกสารอ้างอิงสำหรับ Agent ให้สอดคล้องกับโค้ดปัจจุบัน
---

# /docs-sync: Agent Documentation Synchronization

สร้างหรืออัปเดตเอกสารอ้างอิงที่จำเป็นสำหรับ Agent ทั้ง 7 ตัว ให้สอดคล้องกับสถานะโค้ดปัจจุบัน

## เอกสารที่จัดการ

| เอกสาร | วัตถุประสงค์ | Agent หลักที่ใช้ |
|---------|-------------|-----------------|
| `DEPENDENCY_MAP.md` | แผนที่พึ่งพาข้ามระบบ | Risk Analyst, Architect, Operations |
| `USER_CONTEXT.md` | บริบทผู้ใช้งาน | User Expert, Quality, Architect |
| `API_CATALOG.md` | แคตตาล็อก API | Risk Analyst, Architect, Quality, Security |

---

## Step 1: ตรวจสอบเอกสารที่มีอยู่

ตรวจสอบว่าเอกสารทั้ง 3 ชุดมีอยู่แล้วหรือไม่:

```
for each doc in [DEPENDENCY_MAP.md, USER_CONTEXT.md, API_CATALOG.md]:
  if doc exists in project root:
    read current content → mark as "needs update"
  else:
    mark as "needs creation"
```

---

## Step 2: Scan โค้ดปัจจุบัน

สแกนโค้ดเพื่อเก็บข้อมูลล่าสุด:

### 2.1 API Endpoints
// turbo
ใช้ `code_search` หา route definitions:
- `exat-api-service/routes/routes.go` — Go API routes ทั้งหมด
- `exat-model-service/Routes/Api.py` — Flask API routes
- ตรวจ `exat-api-service/src/*/handlers/*.go` — เช็ค handler ใหม่ที่อาจยังไม่มีใน routes.go

### 2.2 Frontend Pages & API Usage
// turbo
ใช้ `code_search` หา frontend API usage:
- `exat-web/pages/` — หน้าทั้งหมด (list_dir)
- `exat-web/core/modules/*/infrastructure/*Service.ts` — Service classes ที่เรียก API
- `exat-web/core/shared/http/HttpService.ts` — base HTTP client

### 2.3 Database Models
// turbo
ใช้ `code_search` หา model definitions:
- `exat-api-service/models/*.go` — GORM models ทั้งหมด
- `exat-api-service/src/*/repositories/*.go` — repository layer (ตรวจ DB queries)

### 2.4 Inter-service Communication
// turbo
ใช้ `grep_search` หา:
- `PYTHON` env var usage — API → Model service calls
- `http.NewRequest` — outgoing HTTP calls
- `MongoDb` usage — MongoDB access points

### 2.5 Permission Middleware
// turbo
ใช้ `grep_search` หา:
- `middlewares.*Permission()` — permission middleware definitions
- `role_access_controls` — access control structure

---

## Step 3: เปรียบเทียบและอัปเดตเอกสาร

### 3.1 DEPENDENCY_MAP.md

เปรียบเทียบข้อมูลที่ scan ได้กับเอกสารปัจจุบัน:

**ตรวจสอบ:**
- [ ] Service Architecture Overview — เพิ่ม/ลบ service หรือเปลี่ยน port?
- [ ] Frontend ↔ Backend API Contract Map — มี module ใหม่? หน้าใหม่? endpoint เปลี่ยน?
- [ ] Shared Database Tables — มีตารางใหม่? ตารางที่ใช้ร่วมกันเพิ่ม?
- [ ] API → Model Service Integration — endpoint เปลี่ยน? auth เปลี่ยน?
- [ ] Revision System — มีโมดูลใหม่ที่ใช้ revision pattern?
- [ ] Permission Middleware — มี middleware ใหม่?
- [ ] Frontend Shared State — มี Pinia store ใหม่?
- [ ] High-Risk Change Zones — จุดเสี่ยงเปลี่ยนแปลง?

**ถ้าพบการเปลี่ยนแปลง → อัปเดตส่วนนั้นของเอกสาร**

### 3.2 USER_CONTEXT.md

**ตรวจสอบ:**
- [ ] User Personas — มีบทบาทผู้ใช้ใหม่? (เช็คจาก roles + permission structure)
- [ ] Key Workflows — มี workflow ใหม่? (เช็คจากหน้าใหม่ใน pages/)
- [ ] Feature Criticality — มีฟีเจอร์ใหม่ที่ใช้ทุกวัน?
- [ ] Pain Points — มีปัญหาใหม่ที่ผู้ใช้รู้สึก? (เช็คจาก bug reports, TODO comments)
- [ ] UX Decision Guide — กฎเหล็กยังถูกต้อง?

**ถ้าพบการเปลี่ยนแปลง → อัปเดตส่วนนั้นของเอกสาร**

### 3.3 API_CATALOG.md

**ตรวจสอบ:**
- [ ] API Base URLs — env var เปลี่ยน? port เปลี่ยน?
- [ ] ทุก endpoint group — เพิ่ม/ลบ/เปลี่ยน endpoint?
- [ ] Permission middleware — เปลี่ยน permission ของ endpoint?
- [ ] Model Service API — endpoint เปลี่ยน? process เปลี่ยน?
- [ ] Response Format — format เปลี่ยน?

**ถ้าพบการเปลี่ยนแปลง → อัปเดตส่วนนั้นของเอกสาร**

---

## Step 4: สร้างเอกสารใหม่ (ถ้ายังไม่มี)

ถ้าเอกสารใดยังไม่มี ให้สร้างขึ้นใหม่โดยใช้ข้อมูลจาก Step 2:

### สร้าง DEPENDENCY_MAP.md
ใช้ template จากเอกสารปัจจุบัน (ถ้ามี) หรือสร้างใหม่ตามโครงสร้าง:
1. Service Architecture Overview (diagram)
2. Frontend ↔ Backend API Contract Map (table)
3. Shared Database Tables (table with risk levels)
4. API → Model Service Integration (detail)
5. Revision System cross-module pattern
6. Permission Middleware Chain
7. Frontend Shared State (Pinia stores)
8. High-Risk Change Zones

### สร้าง USER_CONTEXT.md
ใช้ template:
1. ระบบนี้คืออะไร
2. User Personas (5 personas)
3. Key User Workflows
4. Feature Criticality Matrix
5. User Pain Points & Priorities
6. UX Decision Guide + กฎเหล็ก

### สร้าง API_CATALOG.md
ใช้ template:
1. API Base URLs
2. Authentication endpoints
3. ทุก endpoint group (ตาม routes.go)
4. Model Service API
5. Static File Endpoints
6. Common Response Format

---

## Step 5: อัปเดต AGENTS.md

ตรวจสอบว่า `AGENTS.md` มี section "Agent Reference Documents" ที่อ้างอิงเอกสารทั้ง 3 ชุด:
- ถ้ามีอยู่แล้ว → ตรวจว่าข้อมูลถูกต้อง
- ถ้าไม่มี → เพิ่ม section ใหม่

---

## Critical Rules

1. **อย่าเขียนเอกสารจากความจำ**: ทุกข้อมูลต้องมาจากการ scan โค้ดจริง
2. **อย่าลบข้อมูลเดิมโดยไม่จำเป็น**: เพิ่ม/แก้ไขเท่านั้น ยกเว้นข้อมูลที่ผิดแน่นอน
3. **รักษาโครงสร้างเดิม**: ถ้าเอกสารมีอยู่แล้ว อัปเดตเฉพาะส่วนที่เปลี่ยน อย่าเขียนใหม่ทั้งไฟล์
4. **สแกนให้ครบ**: อย่าข้ามโมดูลใดๆ โดยเฉพาะโมดูลใหม่ที่อาจเพิ่งสร้าง
5. **เช็ค permission เสมอ**: ทุกครั้งที่เพิ่ม endpoint ต้องระบุ permission ที่ถูกต้อง
6. **Frontend mapping ต้องตรง**: ถ้าเพิ่ม endpoint ต้องหาหน้า frontend ที่ใช้ endpoint นั้น
