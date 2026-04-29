---
name: redesign
description: >-
  Pixel-perfect UI redesign — adjust UI to be beautiful and polished while preserving the
  project's existing style. Use when polishing UI, fixing layout issues, adjusting alignment,
  or when the user mentions UI redesign, pixel perfect, layout fix, styling, or visual polish.
---

# /redesign: Pixel-Perfect UI Redesign

Adjust UI for beauty and polish while **preserving existing style** and **achieving pixel-perfect alignment**.

> Do not change frameworks, theme systems, or add new CSS libraries — work within the existing system.

---

## Step 1: Analyze the Target Page

Identify the page/section to adjust:

1. Read the Vue/React file of the page
2. Read related components
3. Capture screenshot or open browser preview to see current state
4. Identify issues: what looks unpolished? misaligned? unbalanced?

### Analysis Checklist
- [ ] Alignment: text/buttons/tables aligned consistently?
- [ ] Spacing: padding/margin consistent? not too tight/loose?
- [ ] Typography: font-size, font-weight balanced? headings prominent?
- [ ] Visual Hierarchy: what's most important looks most prominent?
- [ ] Consistency: same component uses same style across all pages?
- [ ] Responsive: mobile/tablet looks good?
- [ ] Whitespace: enough breathing room?

---

## Step 2: Design Improvements (before editing code)

Before editing any file, explain the plan:

### 2.1 Identify What Will Change
```
Page: [page name]
File: [path to .vue/.tsx file]

Changes:
1. [what] → [change to what] (reason: ...)
2. [what] → [change to what] (reason: ...)

Not changing (preserving):
- [what] (reason: still good / part of existing style)
```

### 2.2 Pixel-Perfect Alignment Rules

| Element | Rule | Standard Value |
|---------|------|---------------|
| **Card padding** | Equal on all sides | `p-5` (20px) or `p-8` (30px) |
| **Section gap** | Between card/section | `mb-5` or `mb-10` |
| **Label-input gap** | Between label and input | `mb-2` |
| **Button group gap** | Between buttons | `me-3` (12px) |
| **Table header** | Bold, darker than body | `fw-bold`, `text-gray-800` |
| **Table cell** | vertical-align middle | `align-middle` |
| **Icon+Text** | Gap between icon and text | `me-2` or `ms-2` |
| **Page title** | Prominent, with breadcrumb | `fs-2 fw-bold text-gray-900` |
| **Badge/Tag** | Not too big, has padding | `badge-light-*` + `px-3 py-2` |
| **Modal header** | Has border-bottom, title prominent | `border-bottom` + `fs-4 fw-bold` |
| **Form row** | label and input aligned | `row` + `col-lg-3` label + `col-lg-9` input |

---

## Step 3: Surgical Edits

**Do NOT rewrite entire files** — edit only the points that need adjustment:

### 3.1 Edit Order (large to small)
1. **Layout structure** — grid, row, col correct first
2. **Spacing** — margin, padding balanced
3. **Typography** — font-size, weight, color
4. **Alignment** — text-align, vertical-align, flex alignment
5. **Details** — border, shadow, radius, transition

### 3.2 How to Edit (use CSS framework classes primarily)

**If adding spacing:**
```html
<!-- Bad: inline style -->
<div style="margin-top: 15px">

<!-- Good: framework class -->
<div class="mt-5">
```

**If aligning:**
```html
<!-- Bad: inline flex -->
<div style="display: flex; align-items: center; justify-content: space-between">

<!-- Good: framework classes -->
<div class="d-flex align-items-center justify-content-between">
```

**If aligning table cells:**
```html
<!-- Bad: no alignment -->
<td>{{ value }}</td>

<!-- Good: explicit alignment -->
<td class="align-middle text-center">{{ value }}</td>
<td class="align-middle text-end">{{ number }}</td>
<td class="align-middle">{{ text }}</td>
```

### 3.3 Forbidden Actions
- No `style=""` inline — use framework classes
- No hardcoded colors — use theme variables or classes
- No `!important` — fix specificity instead
- No changing SCSS variables unless necessary — use class override
- No adding new CSS frameworks — work within existing system
- No removing dark mode support — must support light/dark
- No changing component library — use existing components

---

## Step 4: Pixel-Perfect Verification

After editing, check every point:

### 4.1 Visual Checklist
- [ ] All cards have equal padding
- [ ] All sections have equal gap
- [ ] All table cells have `align-middle`
- [ ] Numbers right-aligned (`text-end`)
- [ ] Text left-aligned (default)
- [ ] Buttons in same position on every page
- [ ] Form label-input on same line
- [ ] Icon and text have consistent gap everywhere
- [ ] Badge/Tag same size in all contexts
- [ ] Modal dialog centered, appropriate size

### 4.2 Responsive Checklist
- [ ] Desktop (>=1200px): looks complete
- [ ] Tablet (768-1199px): aside collapsed, content adjusted
- [ ] Mobile (<768px): stacked layout, buttons large enough to tap

### 4.3 Dark Mode Checklist
- [ ] Dark mode on → everything readable
- [ ] Colors contrast enough
- [ ] No white backgrounds that should be dark

---

## Step 5: Test in Browser

Open browser preview to verify:
1. Open the page in browser
2. Check the edited page
3. Toggle light/dark mode
4. Try responsive (resize window)
5. If not right → go back to Step 3

---

## Step 6: Repeat for Next Page

If there are multiple pages:
1. Do one page at a time (Step 1-5)
2. After each page → commit separately
3. If you find reusable patterns → create shared component or mixin
