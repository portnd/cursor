---
description: Redesign - ปรับแต่ง UI ให้สวยงาม เรียบหรู ตำแหน่งสมบูรณ์แบบ โดยรักษาสไตล์เดิมของโครงการ
---

# /redesign: Pixel-Perfect UI Redesign

ปรับแต่ง UI ให้สวยงามและเรียบหรู โดย **รักษาสไตล์เดิม** (Keenthemes Metronic + Bootstrap 5) และ **จัดตำแหน่งให้สมบูรณ์แบบเหมือนมนุษย์ทำ**

> ห้ามเปลี่ยน framework, ห้ามเปลี่ยน theme system, ห้ามเพิ่ม library CSS ใหม่ — ทำงานภายใต้ระบบที่มีอยู่

---

## เอกลักษณ์สไตล์เดิมของโครงการ (ห้ามทำลาย)

### ระบบสี (Keenthemes Metronic)
- **Primary**: `#0e4285` (น้ำเงินเข้ม) — ปุ่มหลัก, link, active state
- **Success**: `#50cd89` (เขียว) — สถานะสำเร็จ, เกรด A
- **Warning**: `#f6c000` (เหลือง) — เตือน, เกรด C
- **Danger**: `#f1416c` (แดง) — ผิดพลาด, เกรด E
- **Info**: `#7239ea` (ม่วง) — ข้อมูล, analysis
- **Dark**: `#181c32` (เข้ม) — aside, header
- **Page bg**: `#f9f9f9` (เทาอ่อน) — พื้นหลังหน้า

### ระบบ CSS Variables
- ทุกสีใช้ผ่าน `--kt-*` CSS variables (เช่น `--kt-primary`, `--kt-gray-300`)
- รองรับ dark mode ผ่าน `[data-theme="dark"]`
- อยู่ใน `assets/themes/sass/mains/components/_root.scss`

### Layout
- **Aside**: 125px desktop, 80px mobile, พื้นหลัง `#181c32`
- **Content border-radius**: `1.5rem`
- **Root font-size**: `13px`
- **Input height**: `46px`
- **Card**: มี box-shadow, border-radius ตาม Bootstrap 5

### Component Style
- ใช้ Bootstrap 5 classes เป็นหลัก
- เติม Keenthemes custom classes (`.btn-light-primary`, `.badge-light-success`, etc.)
- Form input: solid style (bg `#f4f4f4`, focus bg white)
- Table: striped, hover, standard Bootstrap

---

## Step 1: วิเคราะห์หน้าที่ต้อง redesign

ระบุหน้า/ส่วนที่ต้องปรับ:

1. อ่านไฟล์ Vue ของหน้านั้น
2. อ่าน component ที่เกี่ยวข้อง
3. จับภาพ screenshot (ถ้ามี) หรือเปิด browser preview เพื่อดูสภาพปัจจุบัน
4. ระบุปัญหา: อะไรที่ดูไม่เรียบ? ตำแหน่งไม่ตรง? ไม่สมดุล?

### Checklist การวิเคราะห์
- [ ] Alignment: ข้อความ/ปุ่ม/ตาราง จัดชิด/กึ่งกลาง สม่ำเสมอ?
- [ ] Spacing: padding/margin สม่ำเสมอ? ไม่แคบเกิน/ห่างเกิน?
- [ ] Typography: font-size, font-weight สมดุล? heading เด่นชัด?
- [ ] Visual Hierarchy: อะไรสำคัญที่สุดดูเด่นที่สุด?
- [ ] Consistency: component เดียวกันใช้ style เดียวกันทุกหน้า?
- [ ] Responsive: มือถือ/แท็บเล็ต ดูดีไหม?
- [ ] Whitespace: มีพื้นที่ว่างพอให้ตาพัก?

---

## Step 2: ออกแบบการปรับปรุง (ก่อนแก้โค้ด)

ก่อนแก้ไฟล์ใด ให้อธิบายแผนก่อน:

### 2.1 ระบุสิ่งที่จะเปลี่ยน
```
หน้า: [ชื่อหน้า]
ไฟล์: [path ของ .vue file]

สิ่งที่จะเปลี่ยน:
1. [อะไร] → [เปลี่ยนเป็นอย่างไร] (เหตุผล: ...)
2. [อะไร] → [เปลี่ยนเป็นอย่างไร] (เหตุผล: ...)

สิ่งที่จะไม่เปลี่ยน (รักษาไว้):
- [อะไร] (เหตุผล: ยังดีอยู่ / เป็นสไตล์เดิม)
```

### 2.2 Pixel-Perfect Alignment Rules

กฎสำหรับจัดตำแหน่งให้สมบูรณ์แบบ:

| ธาตุ | กฎ | ค่ามาตรฐาน |
|------|----|-----------|
| **Card padding** | เท่ากันทุกด้าน | `p-5` (20px) หรือ `p-8` (30px) |
| **Section gap** | ระหว่าง card/section | `mb-5` หรือ `mb-10` |
| **Label-input gap** | ระหว่าง label กับ input | `mb-2` (Bootstrap form standard) |
| **Button group gap** | ระหว่างปุ่ม | `me-3` (12px) |
| **Table header** | ตัวหนา, สีเข้มกว่า body | `fw-bold`, `text-gray-800` |
| **Table cell** | vertical-align middle | `align-middle` |
| **Icon+Text** | ระยะห่าง icon กับข้อความ | `me-2` หรือ `ms-2` |
| **Page title** | เด่น, มี breadcrumb | `fs-2 fw-bold text-gray-900` |
| **Badge/Tag** | ไม่ใหญ่เกิน, มี padding | `badge-light-*` + `px-3 py-2` |
| **Modal header** | มี border-bottom, title เด่น | `border-bottom` + `fs-4 fw-bold` |
| **Form row** | label กับ input ตรงแนว | `row` + `col-lg-3` label + `col-lg-9` input |

---

## Step 3: แก้ไขทีละจุด (Surgical Edits)

**ห้ามเขียนไฟล์ใหม่ทั้งหมด** — แก้เฉพาะจุดที่ต้องปรับ:

### 3.1 ลำดับการแก้ (จากใหญ่ไปเล็ก)
1. **Layout structure** — grid, row, col ให้ถูกต้องก่อน
2. **Spacing** — margin, padding ให้สมดุล
3. **Typography** — font-size, weight, color
4. **Alignment** — text-align, vertical-align, flex alignment
5. **Details** — border, shadow, radius, transition

### 3.2 วิธีแก้ (ใช้ Bootstrap classes เป็นหลัก)

**ถ้าอยากเพิ่มระยะห่าง:**
```html
<!-- ไม่ดี: inline style -->
<div style="margin-top: 15px">

<!-- ดี: Bootstrap class -->
<div class="mt-5">
```

**ถ้าอยากจัดตำแหน่ง:**
```html
<!-- ไม่ดี: inline flex -->
<div style="display: flex; align-items: center; justify-content: space-between">

<!-- ดี: Bootstrap classes -->
<div class="d-flex align-items-center justify-content-between">
```

**ถ้าอยากให้ตารางตำแหน่งตรง:**
```html
<!-- ไม่ดี: ไม่มี alignment -->
<td>{{ value }}</td>

<!-- ดี: ระบุ alignment ชัดเจน -->
<td class="align-middle text-center">{{ value }}</td>
<td class="align-middle text-end">{{ number }}</td>
<td class="align-middle">{{ text }}</td>
```

### 3.3 ห้ามทำ
- ❌ ห้ามใช้ `style=""` inline — ใช้ Bootstrap class
- ❌ ห้าม hardcode สี (เช่น `color: #0e4285`) — ใช้ `text-primary` หรือ `--kt-primary`
- ❌ ห้ามใช้ `!important` — แก้ specificity ให้ถูกต้องแทน
- ❌ ห้ามเปลี่ยน SCSS variable โดยไม่จำเป็น — ใช้ class override
- ❌ ห้ามเพิ่ม CSS framework ใหม่ — ทำงานใน Bootstrap 5 + Keenthemes
- ❌ ห้ามลบ dark mode support — ต้องรองรับทั้ง light/dark
- ❌ ห้ามเปลี่ยน component library — ใช้ component เดิม

---

## Step 4: ตรวจสอบ Pixel-Perfect

หลังแก้ไขแล้ว ตรวจสอบทุกจุด:

### 4.1 Visual Checklist
- [ ] ทุก card มี padding เท่ากัน
- [ ] ทุก section มี gap เท่ากัน
- [ ] ทุก table cell มี `align-middle`
- [ ] ตัวเลขจัดชิดขวา (`text-end`)
- [ ] ข้อความจัดชิดซ้าย (default)
- [ ] ปุ่มจัดตำแหน่งเดียวกันทุกหน้า (มุมขวาล่าง หรือมุมขวาบน)
- [ ] Form label-input ตรงแนวเดียวกัน
- [ ] Icon กับข้อความ ระยะห่างเท่ากันทุกที่
- [ ] Badge/Tag ขนาดเท่ากันในทุก context
- [ ] Modal dialog กึ่งกลางจอ ขนาดเหมาะสม

### 4.2 Responsive Checklist
- [ ] Desktop (≥1200px): ดูสมบูรณ์
- [ ] Tablet (768-1199px): aside ย่อ, content ปรับ
- [ ] Mobile (<768px): stack layout, ปุ่มใหญ่พอแตะ

### 4.3 Dark Mode Checklist
- [ ] เปิด dark mode → ทุกอย่างอ่านได้
- [ ] สีตัดกันพอ
- [ ] ไม่มีพื้นขาวที่ควรเป็นพื้นเข้ม

---

## Step 5: ทดสอบใน Browser

เปิด browser preview เพื่อตรวจสอบจริง:

```
// turbo
1. เปิด http://localhost:3000 ใน browser preview
2. ตรวจสอบหน้าที่แก้
3. สลับ light/dark mode
4. ลอง responsive (ย่อ/ขยายหน้าต่าง)
5. ถ้าไม่ตรง → กลับไป Step 3 แก้ไข
```

---

## Step 6: ทำซ้ำกับหน้าถัดไป

ถ้ามีหลายหน้า:
1. ทำทีละหน้า (Step 1-5)
2. หลังเสร็จแต่ละหน้า → commit แยก
3. ถ้าพบ pattern ที่ใช้ซ้ำ → สร้าง shared component หรือ SCSS mixin

---

## Quick Reference: Bootstrap 5 Alignment Classes

| ต้องการ | Class | ตัวอย่าง |
|---------|-------|---------|
| จัดกึ่งกลางแนวนอน | `text-center` | หัวข้อ, ข้อความสั้น |
| จัดชิดขวา | `text-end` | ตัวเลข, จำนวนเงิน |
| จัดกึ่งกลางแนวตั้ง | `align-middle` | table cell |
| Flex กึ่งกลางทั้ง 2 แกน | `d-flex align-items-center justify-content-center` | empty state, loading |
| Flex ซ้าย-ขวา | `d-flex justify-content-between align-items-center` | header, toolbar |
| Flex ชิดขวา | `d-flex justify-content-end` | ปุ่ม group |
| Grid 2 คอลัมน์ | `row` > `col-lg-6` + `col-lg-6` | form 2 คอลัมน์ |
| Grid 3:9 | `row` > `col-lg-3` + `col-lg-9` | label:input |
| Gap ระหว่าง flex | `gap-3` (12px) / `gap-5` (20px) | button group, card group |

## Quick Reference: Keenthemes Custom Classes

| ต้องการ | Class |
|---------|-------|
| ปุ่มสีอ่อน | `btn-light-primary`, `btn-light-success`, `btn-light-danger` |
| Badge สีอ่อน | `badge-light-primary`, `badge-light-success` |
| พื้นหลังสีอ่อน | `bg-light-primary`, `bg-light-success` |
| ข้อความสีเทา | `text-gray-500` (muted), `text-gray-800` (normal) |
| Card shadow | `card` (auto), `shadow-sm` |
| Symbol (avatar) | `symbol symbol-30px`, `symbol symbol-50px` |
| Bullet | `bullet bullet-dot` |
