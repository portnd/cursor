# Project Detail Page – Comprehensive Context for Gemini

**URL Pattern:** `http://localhost:3000/projects/:id` (e.g. `/projects/hdmap`)  
**Source File:** `web/pages/projects/[id].vue`  
**Auth:** Required (JWT middleware)  
**Route param `:id`:** รับได้ทั้ง UUID ของ project หรือ project code (slug) เช่น `hdmap`

---

## 0. Data Loaded on Page Mount

เมื่อเปิดหน้า ระบบ call API พร้อมกัน 4 รายการ:

| API | Endpoint | เก็บใน |
|-----|----------|--------|
| Get Project | `GET /api/v1/sentinel/projects/:id` | `project` |
| Get Tasks | `GET /api/v1/sentinel/tasks?project_id=...` | `allTasks[]` |
| Get Sprints | `GET /api/v1/sentinel/sprints?project_id=...` | `sprints[]` |
| Get Milestones | `GET /api/v1/sentinel/milestones?project_id=...` | `milestones[]` |
| Get Epics | `GET /api/v1/sentinel/epics?project_id=...` | `epics[]` |

**หมายเหตุ:** Tab `Analytics` และ `Timeline` โหลดข้อมูลเพิ่มเติมแบบ lazy (โหลดเฉพาะเมื่อเปิด tab นั้น)

---

## 1. Header ที่แสดงตลอดทุก Tab

Header ติดอยู่ด้านบน (sticky) แสดง:
- ลิงก์ `← Projects` กลับไปหน้า project list
- ชื่อ project (`project.name`)
- ปุ่ม ✏️ → เปิด **Edit Project Modal**
- Badge สถานะ project: `ACTIVE` (เขียว), `COMPLETED` (น้ำเงิน), `ON_HOLD` (เหลือง)
- Project code (slug) เช่น `hdmap` — แสดงใน medium screens ขึ้นไป
- ปุ่ม **Refresh** (แสดงเฉพาะ tab Timeline)
- แถบ Tabs: Overview · Board · Timeline · Backlog · Sprints · Analytics · Costing

---

## 2. Tab List

| Tab ID | Icon | Label |
|--------|------|-------|
| `overview` | 📊 | Overview |
| `board` | 🗂 | Board |
| `timeline` | 📅 | Timeline |
| `backlog` | 📋 | Backlog |
| `sprints` | 🏃 | Sprints |
| `analytics` | 📈 | Analytics |
| `costing` | 💰 | Costing |

URL query `?tab=xxx` จะเปลี่ยนตาม tab ที่เลือก เช่น `?tab=backlog`

---

## 3. Tab: Overview (📊)

**วัตถุประสงค์:** สรุปภาพรวมของ project ในหน้าเดียว

### 3.1 Key Metrics (4 การ์ดบนสุด)

| Metric | แสดงค่า | สีข้อความ | Logic การคำนวณ |
|--------|---------|-----------|----------------|
| **Active Sprint** | ชื่อ sprint ที่ status = `ACTIVE` หรือ "No sprint" | Purple หรือ Gray | `sprints.find(s => s.status === 'ACTIVE')` |
| **Tasks Done** | `completedCount / totalTasks` | Green | `completedCount = allTasks.filter(t => t.status === 'COMPLETED').length` |
| **In Progress** | จำนวน task ที่ status = `IN_PROGRESS` | Yellow | `allTasks.filter(t => t.status === 'IN_PROGRESS').length` |
| **Overdue** | จำนวน task ที่ overdue | Red (ถ้า > 0) หรือ Green | Task ที่ `status !== 'COMPLETED'` AND `due_at < now` |

**Progress Bar** ใต้ Tasks Done: `completionPct = round(completedCount / totalTasks * 100)`%

### 3.2 Current Sprint Card

แสดงข้อมูลของ active sprint:
- ชื่อ sprint + goal
- 3 ตัวเลข:
  - **Total**: tasks ทั้งหมดที่ `sprint_id === activeSprint.id`
  - **Done**: tasks ที่ `status === 'COMPLETED'` ใน sprint นั้น
  - **Story Pts**: ผลรวม `story_points` ของ tasks ใน sprint
- **Sprint Progress Bar**: `round(done / total * 100)`%
- ปุ่ม **Complete Sprint** (เปลี่ยน status เป็น COMPLETED)
- ปุ่ม **+ Sprint** → เปิด Create Sprint Modal

**All Sprints List** (ด้านล่างการ์ด):
- รายชื่อ sprint ทั้งหมด เรียงตาม `sort_order` แล้ว `created_at`
- สถานะแต่ละ sprint: `ACTIVE` (purple), `COMPLETED` (gray), `PLANNING` (yellow)
- **ลากเพื่อเรียงลำดับ** (drag-and-drop → PATCH sort_order ทีละ sprint)
- Actions:
  - `Start` (เฉพาะ PLANNING, disabled ถ้ามี active sprint อยู่แล้ว) → POST `/sentinel/sprints/:id/start`
  - `Reopen` (เฉพาะ COMPLETED) → POST `/sentinel/sprints/:id/reopen`
  - `+ Tasks` → Add Tasks to Sprint Modal
  - `Edit` → Edit Sprint Modal
  - `Delete` → Delete Sprint Modal

### 3.3 Milestone Tracker Card

Component: `MilestoneTimeline.vue`  
แสดง milestones เป็น timeline bar แนวนอน:
- แต่ละ milestone มีชื่อ, due_date, สถานะ `PENDING` / `REACHED` / `MISSED`
- คลิกที่ milestone → เปิด Edit Milestone Modal
- ปุ่ม `+ Add Milestone` → เปิด Create Milestone Modal

### 3.4 Recent Activity

แสดง 10 tasks ที่อัปเดตล่าสุด (`updated_at` ใหม่สุดก่อน):
- แต่ละแถวแสดง: code, title, priority badge, status badge
- คลิก → navigate ไปหน้า task detail `/task/:id`

---

## 4. Tab: Board (🗂)

**วัตถุประสงค์:** Kanban board แสดง tasks แยกตามสถานะ

Component: `KanbanBoard.vue`

### 4.1 ข้อมูลที่แสดง

- Tasks แบ่งเป็น column ตามสถานะ: `PENDING` · `IN_PROGRESS` · `REVIEW_PENDING` · `BLOCKED` · `COMPLETED`
- Card แต่ละใบแสดง: code (เช่น 001), title, priority badge
- Filter by Sprint: สามารถ filter tasks ตาม sprint ที่ active

### 4.2 Actions บน Board

| Action | ผลลัพท์ |
|--------|---------|
| **ลาก task ข้าม column** | เปลี่ยน status → PATCH `/api/v1/sentinel/tasks/bulk-status` |
| **คลิก task card** | Navigate ไปหน้า task detail |

### 4.3 Logic Status Change

เมื่อลาก task ไปวางที่ column ใหม่:
1. อัปเดต `allTasks[idx].status` ทันที (optimistic update)
2. Call `tasksApi.bulkUpdateStatus([taskId], newStatus)`
3. ถ้า error → reload ข้อมูลทั้งหมด

---

## 5. Tab: Timeline (📅)

**วัตถุประสงค์:** Gantt chart แสดง timeline ของ Epics หรือ Sprints พร้อม tasks

### 5.1 Toolbar

#### Mode Toggle
| Mode | สี | Endpoint |
|------|-----|---------|
| **Epic Roadmap** | Purple | `GET /api/v1/sentinel/projects/:id/timeline/epic-view` |
| **Sprint Execution** | Emerald/Green | `GET /api/v1/sentinel/projects/:id/timeline/sprint-view` |

#### View Toggle (granularity ของ column)
| View | 1 column = |
|------|-----------|
| Day | 1 วัน |
| Week (7d) | 7 วัน (จันทร์–อาทิตย์) |
| Month | 1 เดือน |

#### Toolbar Actions
| ปุ่ม | ผลลัพท์ |
|------|---------|
| **Expand all** | กางทุก row (แสดง tasks ภายใต้ Epic/Sprint) |
| **Collapse all** | ย่อทุก row |
| **Fullscreen** | ขยาย Gantt chart เต็มหน้าจอ (modal overlay) |
| **Export PDF** | Generate PDF timeline → เปิดในแท็บใหม่ (POST `..timeline/export-pdf`) |
| **Today** | Scroll timeline ให้วันปัจจุบันอยู่ในมุมมอง (offset 18% จากซ้าย) |
| **Refresh** (header) | โหลดข้อมูล timeline ใหม่ |

### 5.2 Gantt Chart

Library: `vue-ganttastic` (dark theme)  
Label column width: 220px

**Epic Roadmap Mode:**
- Row Epic (สีม่วง): แท่งครอบ start_date → end_date ของ epic
- Row Task (สีม่วงอ่อน): แท่งตาม task.start_date → task.end_date

**Sprint Execution Mode:**
- Row Sprint (สีเขียว): แท่งครอบ start_date → end_date ของ sprint
- Row Task (สีม่วงอ่อน): แท่งตาม task.start_date → task.end_date

**ปัจจุบัน (Now):** เส้นแนวตั้งสีน้ำเงิน (`border-t-blue-400`)

### 5.3 Milestone Lines

เส้นแนวตั้งสี purple (`bg-purple-500/50`) แสดงตำแหน่งของแต่ละ milestone บน Gantt:
- คลิก diamond icon → Edit Milestone Modal
- ลาก milestone diamond → อัปเดต `due_date` แบบ realtime (PATCH `/sentinel/milestones/:id`)

### 5.4 Gantt Bar Interactions (Drag & Resize)

| Action | ผลลัพท์ |
|--------|---------|
| **ลากแท่ง Task** กลางแท่ง | เลื่อน start/end date ของ task |
| **ลากขอบ** ซ้าย/ขวาของแท่ง Task | Resize (ขยาย/ย่อ) start/end date |
| **ลากแท่ง Epic** | Scale tasks ทั้งหมดภายใต้ epic ตาม ratio |
| **คลิกแท่ง** | Navigate ไปหน้า task/epic/sprint detail |

เมื่อ drag/resize เสร็จ: PATCH `tasksApi.updateTask(id, { start_date, end_date })` + sync Epic dates

**Epic Date Sync Logic (`syncEpicDatesFromTasks`):**  
เมื่อ task ใน epic ถูกเลื่อน ระบบคำนวณ epic.start_date = min(tasks.start_date), epic.end_date = max(tasks.end_date) แล้ว PATCH epic

### 5.5 Milestone Legend

แสดงรายชื่อ milestones ทั้งหมดด้านล่าง Gantt พร้อม due_date

---

## 6. Tab: Backlog (📋)

**วัตถุประสงค์:** จัดการ Product Backlog — สร้าง/จัดเรียง tasks, จัดกลุ่มตาม Epic, กำหนด Sprint

### 6.1 Epics Section (ด้านบน)

แสดง epics เป็น chip/badge:
- แต่ละ chip: dot สี epic, ชื่อ epic, status badge (ถ้าไม่ใช่ PLANNING)
- **ลากเพื่อเรียงลำดับ** (drag-and-drop → PATCH `epic.sort_order`)
- Hover → ปุ่ม ✎ (Edit) / ✕ (Delete)
- ปุ่ม **+ Epic** → Create Epic Modal

### 6.2 Product Backlog Table

Header controls:
- Task count ทั้งหมด (top-level tasks เท่านั้น ไม่นับ sub-tasks)
- **Expand all / Collapse all** — กาง/ย่อทุก Epic group และ sub-tasks
- **Import Slides** → Import from Google Slides Modal
- **+ Task** → Create Task Modal

#### Table Structure (per Epic Group)

1. **Epic Group Header** (clickable toggle, draggable):
   - dot สี, ชื่อ epic, จำนวน tasks
   - hover → ปุ่ม `+ Task` ใน epic นั้น
   - ลากเพื่อเรียง epic

2. **Table Header Row**: (Drag handle) | ID | Task | SP | Priority | Epic | Sprint | Status | (Actions)

3. **Task Rows** (ภายใต้ epic):
   | Column | ข้อมูล | Editable |
   |--------|--------|---------|
   | Drag handle | ⋮⋮ ลากเพื่อเรียงลำดับ | — |
   | ID | ลำดับ 001, 002... (backlog order) | — |
   | Task | ชื่อ task (คลิก → task detail), ✎ (edit title), ⎘ (duplicate) | via modal |
   | SP | Story Points (คลิก → browser prompt แก้ไข) | inline prompt |
   | Priority | Dropdown: 🔴 CRITICAL / 🟠 HIGH / 🟡 MEDIUM / 🟢 LOW | inline select |
   | Epic | Dropdown ย้าย task ไป epic อื่น | inline select |
   | Sprint | Dropdown กำหนด sprint (หรือ Backlog) | inline select |
   | Status | Badge สถานะ (ย่อ 6 ตัวอักษร) | — |
   | + Sub | ปุ่ม add sub-task | — |

4. **Sub-task Rows** (expandable):
   - แสดง indented ด้วย ↳
   - Epic column แสดง "Inherits", Sprint column แสดง "Inherits"
   - ลากเพื่อเรียงใน epic เดียวกันไม่ได้ (parent handles ordering)

5. **Unassigned Group**: tasks ที่ไม่มี epic_id

#### Logic การเรียงลำดับ Tasks (backlogSprintOrderIndex)

```
Sort order:
1. Sprint order: Backlog first (sprint_id = null → index 0), จากนั้นเรียงตาม sprint.sort_order
2. task.sort_order
3. task.created_at
```

#### Task ID Display (taskDisplayCode)

```
allTasksInBacklogOrder: [epics tasks+subs in order] then [unassigned tasks+subs in order]
→ แต่ละ task ได้ index 001, 002, 003...
→ แสดงแทน task.code ที่เป็น string เช่น "hdmap-001"
```

### 6.3 Actions ใน Backlog

| Action | Modal/Method | ผลลัพท์ |
|--------|-------------|---------|
| **+ Task** | Create Task Modal | POST `/sentinel/tasks` สร้าง task ใหม่ |
| **+ Sub** | Create Task Modal (parent_id set) | POST `/sentinel/tasks` สร้าง sub-task |
| **✎ title** | Edit Task Title Modal | PATCH `/sentinel/tasks/:id` update title |
| **⎘ duplicate** | — (inline) | POST `/sentinel/tasks` copy task title+(copy), same epic/sprint |
| **Priority dropdown** | Inline | PATCH `/sentinel/tasks/:id` update priority |
| **Epic dropdown** | Inline | PATCH `/sentinel/tasks/:id` update epic_id + reload Epic timeline |
| **Sprint dropdown** | Inline | PATCH `/sentinel/tasks/:id` update sprint_id + clamp dates + reload Sprint timeline |
| **SP click** | browser `prompt()` | PATCH `/sentinel/tasks/:id` update story_points |
| **Drag to reorder** | Drop handler | PATCH sort_order ของ tasks ที่เปลี่ยน |
| **Import Slides** | Google Slides Modal | POST `/sentinel/import/google-slides/preview` + POST `/sentinel/import/google-slides` |

---

## 7. Tab: Sprints (🏃)

**วัตถุประสงค์:** จัดการ Sprints ทั้งหมดของ project

### 7.1 ข้อมูลที่แสดง

**Active Sprint Hero Card** (ถ้ามี active sprint):
- ชื่อ sprint + สถานะ badge "Active"
- Goal (ถ้ามี)
- Start date / End date
- 3 ตัวเลข: Tasks total, Done, Story Pts
- Progress bar: `round(done / total * 100)`%
- ลิงก์ `Open Sprint →` → ไปหน้า `/projects/sprint/:sprintId`

**All Sprints List**:
- เรียงตาม sort_order
- แต่ละ sprint แสดง: ชื่อ, status badge, goal, date range, stats (`done/total tasks`, SP)
- **ลากเพื่อเรียงลำดับ** (drag-and-drop)

### 7.2 Sprint Stats Logic (getSprintStats)

```typescript
tasks = allTasks.filter(t => !t.parent_id && t.sprint_id === sprintId)
return {
  total: tasks.length,
  done: tasks.filter(t => t.status === 'COMPLETED').length,
  sp: tasks.reduce((s, t) => s + (t.story_points || 0), 0)
}
```

(นับเฉพาะ top-level tasks ไม่นับ sub-tasks)

### 7.3 Actions ต่อ Sprint

| Action | Condition | Endpoint | ผลลัพท์ |
|--------|-----------|----------|---------|
| **Create Sprint** (header) | — | POST `/sentinel/sprints` | สร้าง sprint ใหม่ status = PLANNING |
| **Start** | sprint.status = PLANNING AND ไม่มี active sprint | POST `/sentinel/sprints/:id/start` | เปลี่ยน status → ACTIVE |
| **Complete Sprint** | sprint.status = ACTIVE | POST `/sentinel/sprints/:id/complete` | เปลี่ยน status → COMPLETED |
| **Reopen** | sprint.status = COMPLETED | POST `/sentinel/sprints/:id/reopen` | เปลี่ยน status → ACTIVE |
| **+ Add Tasks** | — | POST `/sentinel/sprints/:id/tasks` | เพิ่ม tasks เข้า sprint (batch) |
| **Edit** | — | PATCH `/sentinel/sprints/:id` | แก้ไข name/goal/dates |
| **Delete** | — | DELETE `/sentinel/sprints/:id` | ลบ sprint (tasks ย้ายกลับ Backlog) |
| **Open Sprint →** | ACTIVE sprint hero | navigate | ไปหน้า Sprint detail |

**Rule:** หนึ่งโปรเจกต์มีได้แค่ 1 sprint ที่ ACTIVE พร้อมกัน

---

## 8. Tab: Analytics (📈)

**วัตถุประสงค์:** ดูสถิติและ KPIs ของ project

API: `GET /api/v1/sentinel/projects/:id/analytics` → โหลด lazy เมื่อเข้า tab

Component: `ProjectAnalytics.vue`

### 8.1 Key Metrics (4 การ์ด)

| Metric | Logic | หน่วย |
|--------|-------|------|
| **Tasks Completed** | `completed_tasks / total_tasks` | ตัวเลข + progress bar |
| **Story Points Done** | `completed_story_points / total_story_points` | ตัวเลข + progress bar |
| **Hours Logged** | `total_logged_minutes / 60` | ชั่วโมง (1 decimal) |
| **Avg Cycle Time** | `avg_cycle_time_days` (จาก backend) | วัน หรือชั่วโมง (ถ้า < 1 วัน) |

**Completion %:** `round(completed_tasks / total_tasks * 100)`  
**SP %:** `round(completed_story_points / total_story_points * 100)`  
**Hours:** `(total_logged_minutes / 60).toFixed(1)`  
**Avg Cycle Time:** ถ้า `d < 1` แสดง `{d*24}h` ถ้า `d >= 1` แสดง `{d.toFixed(1)}d`

### 8.2 Sprint Burndown Chart (Chart.js Line Chart)

- X-axis: วันที่ใน sprint
- Y-axis: งานที่เหลือ (remaining tasks/SP)
- 2 เส้น:
  - **Ideal** (เส้นประ, สีม่วงจาง): เส้นตรงจาก total → 0 ตลอด sprint
  - **Actual** (เส้นทึบ, สีม่วง): `burndown[i].remaining` จาก backend
- ข้อมูลมาจาก `analytics.burndown[]` → `[{ day, ideal, remaining }]`

### 8.3 Team Velocity Chart (Chart.js Bar Chart)

- X-axis: ชื่อ sprint ที่ผ่านมา
- Y-axis: Story Points
- 2 ชุด bars:
  - **Planned SP** (สีม่วงอ่อน): `velocity[i].planned_sp`
  - **Completed SP** (สีเขียว): `velocity[i].completed_sp`
- ข้อมูลมาจาก `analytics.velocity[]` → `[{ sprint_name, planned_sp, completed_sp }]`

### 8.4 Team Capacity Table

แสดงรายนักพัฒนาที่ assigned กับ project:

| Column | ข้อมูล | หมายเหตุ |
|--------|--------|---------|
| Developer | email หรือ `Dev #id` | พร้อม avatar (ตัวอักษรแรก) |
| Tasks | `assigned_tasks` | จำนวน tasks ที่ assign |
| Est. Hours | `estimated_hours.toFixed(1)h` | จาก `ai_estimated_minutes / 60` รวม |
| Logged Hours | `logged_hours.toFixed(1)h` | จาก time_logs รวม |
| Utilization | `utilization_pct.toFixed(0)%` | `logged_hours / estimated_hours * 100` |
| Bar | Progress bar | visual |

**Utilization Color Logic:**
- `> 120%` → Red (overloaded)
- `> 90%` → Green (optimal)
- `> 50%` → Yellow (moderate)
- `≤ 50%` → Gray (underutilized)

---

## 9. Tab: Costing (💰)

**วัตถุประสงค์:** คำนวณราคาโปรเจกต์แบบ Fully Loaded Cost + สร้าง Quotation PDF

Component: `QuotationBuilder.vue`  
Store: `useCostingStore`

### 9.1 Cost Parameters Section

แสดงพารามิเตอร์จาก **Admin Cost Config** (`GET /api/v1/pricing/cost-config` + `GET /api/v1/pricing/salaries`):

| พารามิเตอร์ | Logic | ค่า default |
|------------|-------|------------|
| **Cost / Manday** | `fullyLoadedCost / billableDays` | — |
| **Cost / Hour** | `costPerManday / workingHoursPerDay` | — |
| **Billable Days** | `workingDaysPerMonth / overheadMultiplier` | 22 / 1.3 ≈ 16.9 |
| **Utilisation** | `1 / overheadMultiplier` | 1/1.3 ≈ 77% |
| **Risk Buffer** | `risk_margin_pct` | 10% |
| **Profit Margin** | `profit_margin_pct` | 25% |

**Derived calculations:**
```
overheadPerDev = (exec_salaries + overhead_amount + pm_salaries) / devCount
costPerDev = (totalDevSalary + totalDevSS) / devCount
fullyLoadedCost = overheadPerDev + costPerDev
billableDays = workingDaysPerMonth / overheadMultiplier
costPerManday = fullyLoadedCost / billableDays
costPerHour = costPerManday / workingHoursPerDay
```

ลิงก์ ⚙️ Edit → `/admin/cost-config`

### 9.2 Actions

| ปุ่ม | เงื่อนไข | ผลลัพท์ |
|------|---------|---------|
| **Calculate Cost** | มี DEV ใน cost config | เปิด Task Selection Modal |
| **Export PDF Quotation** | มีผลการคำนวณแล้ว | POST `/api/v1/sentinel/projects/:id/quotation/export` → เปิด PDF ในแท็บใหม่ |

### 9.3 Task Selection Modal

1. โหลด epics + tasks ของ project
2. แสดง tasks จัดกลุ่มตาม epic (ไม่รวม sub-tasks)
3. สามารถ select/deselect แต่ละ task หรือทั้ง epic
4. ปุ่ม **Select All** / **Deselect All**
5. แสดง: task code, title, date range, priority
6. Tasks ที่ไม่มี dates จะแสดง "No dates set" (สีแดงจาง)
7. กด **Calculate** → POST `/api/v1/sentinel/projects/:id/quotation/calculate`

### 9.4 ผลการคำนวณ (Results Section)

**Model Metrics** (4 การ์ด):
- Cost / Manday
- Total Mandays (ผลรวม `ai_estimated_minutes / (workingHoursPerDay * 60)`)
- Tasks Costed (จำนวน tasks ที่เลือก)

**Task Cost Breakdown Table**:
| Column | ข้อมูล |
|--------|--------|
| Epic | epic_title หรือ "—" |
| Task | ชื่อ task |
| Mandays | `ai_estimated_minutes / (workingHoursPerDay * 60)` |
| Cost (THB) | `mandays * cost_per_manday` |

**Quotation Summary** (Financial):
| รายการ | Logic |
|--------|-------|
| Labor Subtotal | ผลรวม cost ของทุก task ที่เลือก |
| Risk Buffer | `subtotal * risk_margin_pct` |
| Profit Margin | `(subtotal + risk) * profit_margin_pct` |
| Total (before VAT) | `subtotal + risk + profit` |
| VAT (7%) | `total_before_vat * 0.07` |
| **Grand Total** | `total_before_vat + vat` |

สกุลเงิน: **THB** (แสดงด้วย `Intl.NumberFormat 'th-TH'`)

---

## 10. Modals ทั้งหมดในหน้านี้

| Modal | เปิดโดย | Fields | API Call |
|-------|---------|--------|---------|
| **Create Task** | Backlog "+ Task", Board "+ Task" | Title*, Description, AI Estimated Minutes*, Priority, Story Points, Sprint, Epic, Due/Start/End Date | POST `/sentinel/tasks` |
| **Add Sub-task** | Backlog "+ Sub" | เหมือน Create Task แต่ไม่มี Epic/Dates (inherit จาก parent) | POST `/sentinel/tasks` (พร้อม parent_id) |
| **Edit Task Title** | Backlog ✎ | Title only | PATCH `/sentinel/tasks/:id` |
| **Create Epic** | Backlog "+ Epic" | Title*, Description, Color (color picker + presets), Start/End Date | POST `/sentinel/epics` |
| **Edit Epic** | Backlog hover ✎ | เหมือน Create | PATCH `/sentinel/epics/:id` |
| **Create Sprint** | Overview/Sprints "+ Sprint" | Name*, Goal, Start Date, End Date | POST `/sentinel/sprints` |
| **Edit Sprint** | Sprints "Edit" | เหมือน Create | PATCH `/sentinel/sprints/:id` |
| **Add Tasks to Sprint** | Sprints "+ Tasks" | Multi-select tasks (checkbox) | POST `/sentinel/sprints/:id/tasks` |
| **Delete Sprint** | Sprints "Delete" | Confirmation only | DELETE `/sentinel/sprints/:id` |
| **Complete Sprint** | Sprints "Complete Sprint" | Confirmation only | POST `/sentinel/sprints/:id/complete` |
| **Reopen Sprint** | Sprints "Reopen" | Confirmation only | POST `/sentinel/sprints/:id/reopen` |
| **Add Milestone** | Overview Milestone Tracker | Title*, Description, Due Date | POST `/sentinel/milestones` |
| **Edit Milestone** | คลิก milestone | Title*, Description, Due Date, Status (PENDING/REACHED/MISSED) | PATCH `/sentinel/milestones/:id`; DELETE ถ้ากด Delete |
| **Edit Project** | Header ✏️ | Name*, Description, Status, Update Code (checkbox) | PATCH `/sentinel/projects/:id` |
| **Import Google Slides** | Backlog "Import Slides" | Step 1: URL; Step 2: select slides + Priority/SP/Epic | POST `/sentinel/import/google-slides/preview` then POST `/sentinel/import/google-slides` |
| **Task Selection (Costing)** | Costing "Calculate Cost" | Select tasks by epic | POST `/pricing/quotation/calculate` |

---

## 11. Google Slides Import Flow

3 ขั้นตอน:

**Step 1 – Input URL**
- กรอก Google Slides URL (ต้องเปิดสิทธิ์ "Anyone with the link can view")
- กด "โหลดรายการ slide" → POST `/sentinel/import/google-slides/preview`
- Response: `{ presentation_title, slides: [{ index, title, hidden }], already_imported_slide_indices: [] }`

**Step 2 – Select Slides**
- แสดง slides ทั้งหมด พร้อม checkbox
- Slides ที่ "นำเข้าแล้ว" จะ uncheck ไว้โดย default
- ปุ่ม: "เลือกทั้งหมด" / "ยกเลิกทั้งหมด" / "เลือกเฉพาะที่ยังไม่เคยนำเข้า"
- กำหนด Priority (ทุก task ได้ priority เดียวกัน), Story Points, Epic (optional)
- กด "Import N Slides" → POST `/sentinel/import/google-slides`

**Step 3 – Result**
- แสดงจำนวน tasks ที่สร้าง และรายชื่อแต่ละ task
- กด "Done" → ปิด modal, reload backlog

---

## 12. Drag-and-Drop สรุป

| สิ่งที่ลากได้ | ลากเพื่อ | API |
|-------------|---------|-----|
| Sprint (overview/sprints tab) | เรียงลำดับ | PATCH sprint.sort_order ทีละ sprint |
| Epic chip (backlog epics section) | เรียงลำดับ | PATCH epic.sort_order ทีละ epic |
| Epic group header (backlog table) | เรียงลำดับ | PATCH epic.sort_order |
| Task row drag handle (backlog) | เรียงลำดับใน epic เดียวกัน | PATCH task.sort_order ทีละ task |
| Gantt bar กลาง (timeline) | เลื่อน start/end date | PATCH task.start_date + task.end_date |
| Gantt bar ขอบ (timeline) | resize start/end date | PATCH task.start_date + task.end_date |
| Gantt Epic bar (timeline) | scale tasks ทั้งหมด | PATCH ทุก task ใน epic + PATCH epic |
| Milestone diamond (timeline) | เลื่อน due_date | PATCH milestone.due_date |

---

## 13. Computed Values Summary

| Computed | Formula | ใช้ใน |
|----------|---------|-------|
| `activeSprint` | `sprints.find(s => s.status === 'ACTIVE')` | Overview, Board, Sprints |
| `totalTasks` | `allTasks.length` | Overview |
| `completedCount` | `allTasks.filter(t => t.status === 'COMPLETED').length` | Overview |
| `inProgressCount` | `allTasks.filter(t => t.status === 'IN_PROGRESS').length` | Overview |
| `completionPct` | `round(completed / total * 100)` | Overview |
| `overdueCount` | tasks ที่ `status !== 'COMPLETED' && due_at < now` | Overview |
| `recentTasks` | sort by `updated_at` DESC, slice(0, 10) | Overview |
| `sprintTaskCount('total')` | tasks ที่ `sprint_id === activeSprint.id` | Overview |
| `sprintTaskCount('done')` | + filter `status === 'COMPLETED'` | Overview |
| `sprintTaskCount('sp')` | ผลรวม `story_points` | Overview |
| `sprintsOrdered` | sort by `sort_order` then `created_at` | Overview, Sprints |
| `taskDisplayCodeMap` | index 001..N ใน backlog order | Backlog, Board |
| `tasksNotInSprint` | tasks ที่ `sprint_id !== sprintForAddTasks.id` | Add Tasks Modal |

---

## 14. Session Storage (Backlog State Persistence)

เมื่อ navigate จาก Backlog ไปหน้า task detail แล้วกลับมา:
- บันทึก: `expandedEpics`, `expandedEpicGroups`, scroll position
- Key: `sentinel-backlog-expanded-{project.id}`
- Restore เมื่อกลับมาที่ backlog tab ของ project เดิม

---

## 15. Export PDF (Timeline)

API: `GET /api/v1/sentinel/projects/:id/timeline/export-pdf`  
ใช้ `chromedp` (headless Chrome) render timeline แล้วส่งกลับเป็น PDF  
เปิดในแท็บใหม่ผ่าน `window.open`  
Helper: `web/utils/timelinePdfExport.ts`

---

## 16. สรุป API Endpoints ที่หน้านี้ใช้ทั้งหมด

| Method | Path | ใช้ใน |
|--------|------|-------|
| GET | `/sentinel/projects/:id` | load project |
| PATCH | `/sentinel/projects/:id` | edit project |
| GET | `/sentinel/tasks?project_id=...` | load tasks |
| POST | `/sentinel/tasks` | create task / duplicate |
| PATCH | `/sentinel/tasks/:id` | update task fields |
| PATCH | `/sentinel/tasks/bulk-status` | kanban drag status change |
| GET | `/sentinel/sprints?project_id=...` | load sprints |
| POST | `/sentinel/sprints` | create sprint |
| PATCH | `/sentinel/sprints/:id` | edit sprint |
| DELETE | `/sentinel/sprints/:id` | delete sprint |
| POST | `/sentinel/sprints/:id/start` | start sprint |
| POST | `/sentinel/sprints/:id/complete` | complete sprint |
| POST | `/sentinel/sprints/:id/reopen` | reopen sprint |
| POST | `/sentinel/sprints/:id/tasks` | add tasks to sprint |
| GET | `/sentinel/milestones?project_id=...` | load milestones |
| POST | `/sentinel/milestones` | create milestone |
| PATCH | `/sentinel/milestones/:id` | edit / drag milestone |
| DELETE | `/sentinel/milestones/:id` | delete milestone |
| GET | `/sentinel/epics?project_id=...` | load epics |
| POST | `/sentinel/epics` | create epic |
| PATCH | `/sentinel/epics/:id` | edit / reorder epic |
| DELETE | `/sentinel/epics/:id` | delete epic |
| GET | `/sentinel/projects/:id/analytics` | analytics tab |
| GET | `/sentinel/projects/:id/timeline/epic-view` | timeline (epic mode) |
| GET | `/sentinel/projects/:id/timeline/sprint-view` | timeline (sprint mode) |
| GET | `/sentinel/projects/:id/timeline/export-pdf` | export PDF timeline |
| POST | `/sentinel/import/google-slides/preview` | import slides step 1 |
| POST | `/sentinel/import/google-slides` | import slides step 2 |
| GET | `/pricing/cost-config` | costing tab config |
| GET | `/pricing/salaries` | costing tab salaries |
| POST | `/sentinel/projects/:id/quotation/calculate` | costing tab calculate |
| POST | `/sentinel/projects/:id/quotation/export` | costing tab PDF |

---

*Last updated: 2026-03-07 | Source: `web/pages/projects/[id].vue`, `web/components/projects/ProjectAnalytics.vue`, `web/core/modules/pricing/ui/QuotationBuilder.vue`*
