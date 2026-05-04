---
auto_execution_mode: 0
description: Complete usage guide for all Windsurf workflows with step-by-step examples
---

# Windsurf Workflow คู่มือการใช้งาน

> คู่มือสำหรับใช้งานร่วมกับ GLM-5.1 ใน Windsurf IDE
> เรียกผ่าน Cascade panel: `/ชื่อ-workflow` + คำอธิบายงาน

---

## ภาพรวมระบบ

### Workflow เดี่ยว (ใช้เฉพาะที่ต้องการ)
| คำสั่ง | หน้าที่ | ใช้เมื่อ |
|--------|---------|----------|
| `/analysis` | วิเคราะห์ codebase | อยากรู้สถานะปัจจุบัน, เพิ่งเข้าโปรเจกต์ใหม่ |
| `/plan` | วางแผน | งานซับซ้อนแตะหลายไฟล์, มีความเสี่ยง |
| `/code` | เขียน/แก้โค้ด | มีแผนชัดเจนแล้ว, ต้องการ refinement |
| `/test` | เขียนเทส | หลัง /code เสมอ, ต้องการ coverage |
| `/debug` | หาแก้ bug | โค้ดทำงานผิด, มี error |
| `/deploy` | commit & push | โค้ดเสร็จแล้ว, ต้องการ push |
| `/refactor` | ปรับโครงสร้าง | โค้ดทำงานถูกแต่โครงสร้างไม่ดี |
| `/docs` | เขียนเอกสาร | ต้องการ README, API docs, docstrings |
| `/migrate` | ย้ายระบบ | อัปเกรด dependency, เปลี่ยน DB schema |
| `/docs-sync` | ซิงค์เอกสาร Agent | สร้าง/อัปเดต DEPENDENCY_MAP, USER_CONTEXT, API_CATALOG ให้ตรงกับโค้ดปัจจุบัน |
| `/redesign` | ปรับแต่ง UI | สวยงาม เรียบหรู ตำแหน่งสมบูรณ์แบบ โดยรักษาสไตล์เดิมของโครงการ |
| `/continue` | ทำงานต่อ | หลัง credit limit หมด → ทำงานที่ค้างต่อจากจุดที่หยุด |

### Supreme Pipeline (ทุกอย่างในคำสั่งเดียว)
| คำสั่ง | หน้าที่ |
|--------|---------|
| `/supreme` | 7-agent consensus + analysis + plan + code + test + deploy |

### Pipeline แนะนำ
```
งานใหม่:    /analysis → /plan → /code → /test → /deploy
แก้บั๊ก:     /debug → /test → /deploy
ปรับปรุง:    /analysis → /refactor → /test → /deploy
ย้ายระบบ:   /analysis → /migrate → /test → /deploy
เอกสาร:     /analysis → /docs
ซิงค์เอกสาร: /docs-sync
ปรับแต่ง UI: /redesign
ทำงานต่อ:    /continue
ครบจบ:      /supreme
```

---

## วิธีเรียกใช้

**วิธีที่ 1:** พิมพ์ใน Cascade (`Cmd+L` / `Ctrl+L`)
```
/[workflow] [คำอธิบายงาน]
```

**วิธีที่ 2:** Customizations → Workflows → คลิกเลือก

---

## /analysis — วิเคราะห์ Code

**ใช้เมื่อ:** อยากรู้ว่าโค้ดเป็นอย่างไร, มีปัญหาอะไร, เพิ่งเข้าโปรเจกต์ใหม่

**คำสั่ง:**
```
/analysis วิเคราะห์โมดูล authentication ทั้งหมด
```

**ผลลัพธ์:** Analysis Report พร้อม tech stack, dependency map, จุดเสี่ยง (🔴🟡🟢), คำแนะนำ

**ตัวอย่างผลลัพธ์:**
```
## Analysis Report: exat-model-service/Routes/Api.py
| Category | Finding | Severity | Location |
|----------|---------|----------|----------|
| Security | ไม่มี input validation | 🔴 | Api.py:45 |
| Quality | ฟังก์ชันยาว 80+ บรรทัด | 🟡 | Api.py:120 |
| Performance | query ซ้ำใน loop | 🟡 | Api.py:67 |
```

**Best Practice:**
- ✅ ระบุขอบเขตชัดเจน: ทั้งโปรเจกต์ / โมดูล / ไฟล์
- ✅ รันก่อนเริ่มงานทุกครั้งถ้าไม่คุ้นเคยโค้ด
- ❌ ไม่แก้ไขโค้ด — วิเคราะห์อย่างเดียว

---

## /plan — วางแผน

**ใช้เมื่อ:** งานซับซ้อน, แตะหลายไฟล์, มีความเสี่ยง, อยากเปรียบเทียบหลายวิธี

**คำสั่ง:**
```
/plan เพิ่มระบบ CRUD conditions API พร้อม pagination โดยใช้ pattern เดียวกับ inspection
```

**ผลลัพธ์:** หลาย Approach เปรียบเทียบ, Task breakdown, Acceptance criteria, Risk assessment

**ตัวอย่างผลลัพธ์:**
```
### Approach A: ใช้ pattern เดิม — Risk: Low, Effort: M ✅ Recommended
### Approach B: สร้าง generic handler — Risk: Medium, Effort: L

### Task Breakdown
- [ ] Task 1: สร้าง model Conditions ใน Models/
- [ ] Task 2: สร้าง handler ใน Routes/Api.py
- [ ] Task 3: เพิ่ม pagination logic
- [ ] Task 4: เพิ่ม input validation

### Acceptance Criteria
1. GET /api/v1/conditions คืนค่าพร้อม pagination
2. POST/PUT/DELETE ทำงานถูกต้อง
3. ไม่ทำลาย API เดิม
```

**Best Practice:**
- ✅ ให้ context มาก: "ใช้ pattern เดียวกับ X", "ห้ามทำลาย Y"
- ✅ ระบุ constraints: "ห้ามเพิ่ม DB field", "ต้อง backward compatible"
- ✅ รัน /analysis ก่อน /plan เสมอถ้ายังไม่เข้าใจโค้ด
- ❌ อย่าข้ามไป /code เลยถ้างานแตะ 2+ ไฟล์

---

## /code — เขียนโค้ด

**ใช้เมื่อ:** มีแผนชัดเจน, ต้องการเขียน/แก้โค้ด, ต้องการ GLM-5.1 refinement

**คำสั่ง:**
```
/code ตามแผน — เพิ่ม CRUD conditions API ตาม pattern ของ inspection
```

**Agent ทำอะไร:**
1. อ่านไฟล์เป้าหมายทั้งหมดก่อน
2. ระบุ assumptions ก่อนเขียน
3. เขียนโค้ดทีละ task, ตรวจทุก task ก่อนไปต่อ
4. Refinement loop 4 รอบ: Correctness → Edge Cases → Integration → Optimization
5. Self-Review Checklist ก่อนจบ

**Best Practice:**
- ✅ ระบุ pattern: "ตาม pattern ของ X"
- ✅ ระบุ constraints: "ห้ามแก้ไฟล์ X", "ต้องใช้ library Y"
- ❌ อย่าสั่งแบบกว้าง: "ทำระบบใหม่ทั้งหมด" → แบ่งทีละส่วน
- ❌ อย่าหยุดแค่ "มันทำงานได้" — ให้ refinement loop ทำจบ

---

## /test — เขียนเทส

**ใช้เมื่อ:** หลัง /code เสมอ, ต้องการ coverage, มี bug ต้องเขียน regression test

**คำสั่ง:**
```
/test เขียนเทสสำหรับ conditions API — happy path, edge cases, error paths
```

**Agent ทำอะไร:**
1. ค้นหา test infrastructure และ pattern
2. เขียนเทสตาม priority: Happy Path → Edge Cases → Error Paths → Integration
3. รันเทสใหม่ — ต้องผ่านทั้งหมด
4. รันเทสทั้งหมด — ต้องไม่มี regression
5. วัด coverage

**Best Practice:**
- ✅ รัน /test ทุกครั้งหลัง /code — ไม่มีข้อยกเว้น
- ✅ ถ้าโค้ดเดิมไม่มีเทส → เขียน happy path ขั้นต่ำก่อน
- ❌ อย่าเขียนเทสแค่เพื่อ coverage — เทสต้องตรวจพฤติกรรมจริง

---

## /debug — หาและแก้ Bug

**ใช้เมื่อ:** มี bug, โค้ดทำงานผิด, เทส fail ไม่รู้ทำไม่

**คำสั่ง (ให้ข้อมูลให้มาก):**
```
/debug API /api/v1/conditions คืน 500 เมื่อ page=-1
error: "Offset must be non-negative"
เกิดเฉพาะค่า page ติดลบ
```

**Agent ทำอะไร:**
1. Reproduce bug
2. สร้าง hypotheses เรียงความน่าจะเป็น
3. ทดสอบ hypothesis ทีละอัน
4. ระบุ root cause พร้อมตำแหน่งไฟล์:บรรทัด
5. แก้ที่ต้นเหตุ (ไม่ใช่อาการ)
6. เขียน regression test บังคับ

**Best Practice:**
- ✅ ให้ข้อมูลมาก: error message, stack trace, เงื่อนไข, สิ่งที่เปลี่ยนล่าสุด
- ❌ อย่าบอกแค่ "มันไม่ทำงาน" — บอกว่าผิดยังไง
- ❌ ทุก bug fix ต้องมี regression test — ไม่มีข้อยกเว้น

---

## /deploy — Commit & Push

**ใช้เมื่อ:** โค้ดเสร็จ, ต้องการ push, ต้องการตรวจก่อน push

**คำสั่ง:**
```
/deploy push ขึ้น feature/conditions-api
```

**Agent ทำอะไร:**
1. ตรวจ `git status` + `git diff`
2. สแกนไม่มี secret ใน diff
3. สร้าง commit message ตาม convention
4. **ขออนุมัติก่อน commit**
5. Stage + commit + push (ขออนุมัติก่อน push)
6. ตรวจ CI/CD status

**Best Practice:**
- ✅ รัน /test ก่อน /deploy เสมอ
- ✅ ตรวจ commit message ก่อนอนุมัติ
- ❌ อย่า push ไป main โดยตรง (ถ้ามี branch rules)
- ❌ อย่า force push ยกเว้นได้รับอนุมัติ

---

## /refactor — ปรับโครงสร้าง

**ใช้เมื่อ:** โค้ดทำงานถูกแต่โครงสร้างไม่ดี, ซ้ำ, ยาว, coupling สูง

**คำสั่ง:**
```
/refactor ฟังก์ชัน process_inspection_data ใน Models/Inspection.py ยาวเกิน แบ่งเป็น sub-functions
```

**Agent ทำอะไร:**
1. ตรวจมีเทสไหม (ไม่มี → เขียนก่อน)
2. เลือก technique: Extract Function, Extract Utility, Guard Clauses ฯลฯ
3. ปรับทีละขั้นเล็กๆ
4. รันเทสหลังทุกขั้น — fail = revert ทันที

**Best Practice:**
- ✅ ต้องมีเทสก่อน refactor — ไม่มีเทส = ไม่ refactor
- ✅ ทำทีละขั้นเล็กๆ รันเทสทุกขั้น
- ✅ แยก refactor ออกจาก feature change
- ❌ อย่าเปลี่ยนพฤติกรรม — refactor = output เดิม
- ❌ อย่าเปลี่ยนเทส — ถ้าต้องเปลี่ยนแสดงว่าเปลี่ยนพฤติกรรม

---

## /docs — เขียนเอกสาร

**ใช้เมื่อ:** ต้องการ README, API docs, docstrings, architecture docs

**คำสั่ง:**
```
/docs สร้าง API documentation สำหรับ /api/v1/conditions endpoints
```

**Best Practice:**
- ✅ ระบุประเภท: API docs, README, inline docs, runbook
- ✅ ระบุ audience: developer ใหม่, API consumer, ops
- ❌ อย่าเขียนขัดแย้งกับโค้ดจริง — agent ต้องอ่านโค้ดก่อนเขียน

---

## /migrate — ย้ายระบบ

**ใช้เมื่อ:** อัปเกรด dependency, เปลี่ยน DB schema, เปลี่ยน framework

**คำสั่ง:**
```
/migrate อัปเกรด Flask จาก 2.x เป็น 3.x ใน exat-model-service
```

**Agent ทำอะไร:**
1. จำแนกประเภทและประเมินความเสี่ยง
2. ตรวจมีเทส (ไม่มี → เขียนก่อน)
3. สร้าง plan 3 phase: Preparation → Transition → Cutover
4. ดำเนินการทีละ phase, commit แยกทุก phase
5. แต่ละ phase มี rollback plan ของตัวเอง

**Best Practice:**
- ✅ ต้องมีเทสก่อน migrate — เทสคือ safety net
- ✅ Commit ทุก phase แยกกัน — สำหรับ rollback ง่าย
- ✅ สำรองข้อมูลก่อน DB migration เสมอ
- ❌ อย่าทำ big bang — ใช้ 3-phase approach
- ❌ อย่าลืมลบ compatibility shims หลัง migration เสร็จ

---

## /supreme — สุดยอด Workflow

**ใช้เมื่อ:** ต้องการ pipeline ครบจบในคำสั่งเดียว, งานสำคัญที่ต้องการ quality สูงสุด

**คำสั่ง:**
```
/supreme เพิ่มระบบ RBAC ให้ /users API โดยไม่ทำลาย authentication เดิม
```

**Agent ทำอะไร (4 Phase):**

**Phase 1: 7-Agent Consensus** — 7 ผู้เชี่ยวชาญโหวตก่อนเขียนโค้ด
- Architect: สถาปัตยกรรมถูกไหม?
- Quality: โค้ดสะอาดไหม?
- Security: มีช่องโหว่ไหม?
- Performance: ประสิทธิภาพดีไหม?
- Operations: deploy ปลอดภัยไหม?
- Risk Analyst: การแก้ไขจะกระทบการทำงานอื่นของระบบไหม?
- User Expert: ตอบโจทย์ผู้ใช้จริงไหม? UX ดีขึ้นหรือแย่ลง?

กฎโหวต: 7/7 = ไปเลย, 5-6 = ไปแต่แก้ข้อกังวล, <5 = หยุดวิเคราะห์ใหม่

**Phase 2: Implementation Plan** — สร้างแผนมีโครงสร้าง

**Phase 3: Long-Horizon Execution** — ใช้พลัง GLM-5.1 ที่ไม่ plateau
- /analysis → /code → /test → Refinement Loop (วนจนดีที่สุด)

**Phase 4: Final Validation & Delivery** — /deploy + Delivery Report

**Best Practice:**
- ✅ ใช้กับงานสำคัญที่ต้องการ quality สูงสุด
- ✅ ให้ context ชัดเจน — 5 agents ต้องมีข้อมูลพอประเมิน
- ❌ อย่าใช้กับงานเล็ก (แก้ typo, เปลี่ยนชื่อตัวแปร) — ใช้ /code ธรรมดาพอ

---

## Best Practices รวม

### 1. เลือก Workflow ให้ถูกงาน
- งานเล็ก (1 ไฟล์, เปลี่ยนน้อย): `/code` โดยตรง
- งานกลาง (2-5 ไฟล์): `/analysis` → `/plan` → `/code` → `/test`
- งานใหญ่ (5+ ไฟล์, มีความเสี่ยง): `/supreme`
- แก้บั๊ก: `/debug` → `/test` → `/deploy`

### 2. ให้ Context ให้มาก
```
❌ /code แก้ไข API
✅ /code เพิ่ม endpoint GET /api/v1/conditions พร้อม pagination
   ตาม pattern ของ /api/v1/inspections ใน Routes/Api.py
   ต้องรองรับ page, per_page parameters
```

### 3. รันตามลำดับ Pipeline
```
✅ /analysis → /plan → /code → /test → /deploy
❌ /code → /deploy (ข้ามเทส = เสี่ยง)
❌ /plan → /deploy (ข้ามเขียนโค้ด?)
```

### 4. ใช้ GLM-5.1 Long-Horizon ให้เต็มที่
- อย่าหยุดแค่ "มันทำงานได้" — ให้ refinement loop ทำงาน
- ระบุให้ agent หาโอกาสปรับปรุงเพิ่มเติม
- ใช้ /supreme สำหรับงานที่ต้องการ quality สูงสุด

### 5. ตรวจสอบผลลัพธ์เสมอ
- อ่าน Analysis Report ก่อนตัดสินใจ
- อ่าน Implementation Plan ก่อนอนุมัติให้ /code
- ตรวจ commit message ก่อนอนุมัติ /deploy
- อ่าน Test Report — ถ้ามี regression ต้องแก้ก่อน deploy

### 6. ถ้าไม่แน่ใจ — ถาม
- ถ้า agent ถาม — ตอบให้ชัดเจน
- ถ้า agent บอก consensus < 3 — อย่าฝืน ให้ข้อมูลเพิ่ม
- ถ้า agent บอกต้องเขียนเทสก่อน — ให้เขียนเทสก่อน
