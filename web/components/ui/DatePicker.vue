<template>
  <div ref="rootEl" class="dp-root">
    <!-- ── Trigger ── -->
    <button
      type="button"
      @click="toggle"
      class="dp-trigger"
      :class="[modelValue ? 'dp-trigger--filled' : 'dp-trigger--empty', isOpen && 'dp-trigger--open']"
    >
      <svg class="dp-trigger-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <rect x="3" y="4" width="18" height="18" rx="2" stroke-width="1.75"/>
        <path stroke-linecap="round" stroke-width="1.75" d="M16 2v4M8 2v4M3 10h18"/>
      </svg>
      <span class="dp-trigger-value">{{ displayValue }}</span>
      <button
        v-if="modelValue"
        type="button"
        @click.stop="clearValue"
        class="dp-clear-btn"
        tabindex="-1"
        aria-label="Clear date"
      >
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-3 h-3">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
      <svg
        class="dp-chevron"
        :class="isOpen && 'rotate-180'"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
      </svg>
    </button>

    <!-- ── Portal Panel ── -->
    <Teleport to="body">
      <Transition name="dp-fade">
        <div
          v-if="isOpen"
          ref="panelEl"
          class="dp-panel"
          :class="{ 'dp-panel--light': !isDark }"
          :style="panelStyle"
          @mousedown.prevent
          role="dialog"
          aria-label="Date picker"
        >
          <!-- Month / Year Navigation -->
          <div class="dp-nav">
            <button type="button" @click="addYear(-1)" class="dp-nav-btn" title="Previous year">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-3.5 h-3.5">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7M18 19l-7-7 7-7"/>
              </svg>
            </button>
            <button type="button" @click="addMonth(-1)" class="dp-nav-btn" title="Previous month">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <button type="button" @click="cycleYearInput" class="dp-nav-label">
              {{ MONTHS[viewMonth] }} {{ viewYear }}
            </button>
            <button type="button" @click="addMonth(1)" class="dp-nav-btn" title="Next month">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>
            <button type="button" @click="addYear(1)" class="dp-nav-btn" title="Next year">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-3.5 h-3.5">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M6 5l7 7-7 7"/>
              </svg>
            </button>
          </div>

          <!-- Quick Shortcuts -->
          <div class="dp-shortcuts">
            <button type="button" @click="pickOffset(0)" class="dp-shortcut-btn">Today</button>
            <button type="button" @click="pickOffset(1)" class="dp-shortcut-btn">Tomorrow</button>
            <button type="button" @click="pickOffset(7)" class="dp-shortcut-btn">+7 days</button>
            <button type="button" @click="pickOffset(30)" class="dp-shortcut-btn">+30 days</button>
          </div>

          <!-- Weekday Headers -->
          <div class="dp-weekday-row">
            <span v-for="w in WEEKDAYS" :key="w" class="dp-weekday-cell">{{ w }}</span>
          </div>

          <!-- Day Grid -->
          <div class="dp-days-grid">
            <button
              v-for="d in calendarDays"
              :key="d.key"
              type="button"
              @click="pickDay(d)"
              :class="dayClass(d)"
              :tabindex="d.isCurrentMonth ? 0 : -1"
            >
              <span class="dp-day-num">{{ d.d }}</span>
              <span v-if="d.isToday && !d.isSelected" class="dp-today-dot" />
            </button>
          </div>

          <!-- Footer: selected value display -->
          <div v-if="modelValue" class="dp-footer">
            <svg class="w-3.5 h-3.5 text-violet-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            <span class="dp-footer-label">{{ displayValue }}</span>
            <button type="button" @click="clearValue" class="dp-footer-clear">Clear</button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
const { isDark } = useTheme()

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]
const WEEKDAYS = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa']

interface DayCell {
  key: string
  d: number
  ymd: string
  isCurrentMonth: boolean
  isToday: boolean
  isSelected: boolean
  isWeekend: boolean
  isDisabled: boolean
}

const props = withDefaults(defineProps<{
  modelValue?: string
  placeholder?: string
  disabled?: boolean
  min?: string
  max?: string
}>(), {
  modelValue: '',
  placeholder: 'Select date…',
  disabled: false,
  min: '',
  max: '',
})

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
}>()

const rootEl = ref<HTMLElement | null>(null)
const panelEl = ref<HTMLElement | null>(null)
const isOpen = ref(false)

const today = new Date()
today.setHours(0, 0, 0, 0)

function parsedOrToday(): Date {
  if (props.modelValue) {
    const [y, m, d] = props.modelValue.split('-').map(Number)
    if (y && m && d) return new Date(y, m - 1, d)
  }
  return new Date(today)
}

const viewYear  = ref(parsedOrToday().getFullYear())
const viewMonth = ref(parsedOrToday().getMonth())

watch(() => props.modelValue, (v) => {
  if (v) {
    const [y, m] = v.split('-').map(Number)
    if (y && m) { viewYear.value = y; viewMonth.value = m - 1 }
  }
})

function toYMD(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${dd}`
}

const displayValue = computed(() => {
  if (!props.modelValue) return props.placeholder
  const [y, m, d] = props.modelValue.split('-').map(Number)
  if (!y || !m || !d) return props.placeholder
  const date = new Date(y, m - 1, d)
  return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
})

const calendarDays = computed((): DayCell[] => {
  const firstOfMonth = new Date(viewYear.value, viewMonth.value, 1)
  const lastOfMonth  = new Date(viewYear.value, viewMonth.value + 1, 0)
  const startDow = firstOfMonth.getDay()
  const todayYMD = toYMD(today)
  const days: DayCell[] = []

  const isDisabledYMD = (ymd: string) =>
    (props.min ? ymd < props.min : false) || (props.max ? ymd > props.max : false)

  for (let i = startDow - 1; i >= 0; i--) {
    const d = new Date(viewYear.value, viewMonth.value, -i)
    const ymd = toYMD(d)
    days.push({ key: ymd, d: d.getDate(), ymd, isCurrentMonth: false, isToday: ymd === todayYMD, isSelected: props.modelValue === ymd, isWeekend: d.getDay() === 0 || d.getDay() === 6, isDisabled: isDisabledYMD(ymd) })
  }

  for (let i = 1; i <= lastOfMonth.getDate(); i++) {
    const d = new Date(viewYear.value, viewMonth.value, i)
    const ymd = toYMD(d)
    days.push({ key: ymd, d: i, ymd, isCurrentMonth: true, isToday: ymd === todayYMD, isSelected: props.modelValue === ymd, isWeekend: d.getDay() === 0 || d.getDay() === 6, isDisabled: isDisabledYMD(ymd) })
  }

  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const d = new Date(viewYear.value, viewMonth.value + 1, i)
    const ymd = toYMD(d)
    days.push({ key: ymd, d: d.getDate(), ymd, isCurrentMonth: false, isToday: ymd === todayYMD, isSelected: props.modelValue === ymd, isWeekend: d.getDay() === 0 || d.getDay() === 6, isDisabled: isDisabledYMD(ymd) })
  }

  return days
})

function dayClass(d: DayCell): string[] {
  const base = ['dp-day-btn']
  if (d.isDisabled) { base.push('dp-day--disabled'); return base }
  if (d.isSelected) { base.push('dp-day--selected'); return base }
  if (!d.isCurrentMonth) { base.push('dp-day--other'); return base }
  base.push('dp-day--current')
  if (d.isToday) base.push('dp-day--today')
  if (d.isWeekend) base.push('dp-day--weekend')
  return base
}

function addMonth(n: number) {
  const d = new Date(viewYear.value, viewMonth.value + n, 1)
  viewYear.value  = d.getFullYear()
  viewMonth.value = d.getMonth()
}

function addYear(n: number) {
  viewYear.value += n
}

function cycleYearInput() {
  const y = parseInt(window.prompt('Jump to year:', String(viewYear.value)) ?? '')
  if (!isNaN(y) && y > 1900 && y < 2200) viewYear.value = y
}

function pickOffset(days: number) {
  const d = new Date(today)
  d.setDate(d.getDate() + days)
  selectYMD(toYMD(d))
}

function pickDay(d: DayCell) {
  if (d.isDisabled) return
  if (!d.isCurrentMonth) {
    const [y, m] = d.ymd.split('-').map(Number)
    viewYear.value  = y
    viewMonth.value = m - 1
  }
  selectYMD(d.ymd)
}

function selectYMD(ymd: string) {
  emit('update:modelValue', ymd)
  isOpen.value = false
}

function clearValue() {
  emit('update:modelValue', '')
  isOpen.value = false
}

function toggle() {
  if (props.disabled) return
  if (!isOpen.value) {
    const base = parsedOrToday()
    viewYear.value  = base.getFullYear()
    viewMonth.value = base.getMonth()
  }
  isOpen.value = !isOpen.value
}

const panelStyle = computed(() => {
  if (!isOpen.value || !rootEl.value) return {}
  const rect = rootEl.value.getBoundingClientRect()
  const panelW = Math.max(rect.width, 300)
  const IDEAL_H  = 390
  const spaceBelow = window.innerHeight - rect.bottom - 8
  const spaceAbove = rect.top - 8
  // Open upward when below doesn't have enough room AND above has more (or equal) space
  const openUp = spaceBelow < IDEAL_H && spaceAbove >= spaceBelow
  const left = Math.max(8, Math.min(rect.left, window.innerWidth - panelW - 8))

  if (openUp) {
    // Anchor panel BOTTOM flush to trigger TOP — works regardless of actual panel height
    return {
      position: 'fixed' as const,
      bottom: `${window.innerHeight - rect.top + 4}px`,
      left: `${left}px`,
      width: `${panelW}px`,
      maxHeight: `${Math.max(spaceAbove, 260)}px`,
      zIndex: 9999,
    }
  }
  return {
    position: 'fixed' as const,
    top: `${rect.bottom + 4}px`,
    left: `${left}px`,
    width: `${panelW}px`,
    maxHeight: `${Math.max(spaceBelow, 260)}px`,
    zIndex: 9999,
  }
})

function onOutsideClick(e: MouseEvent) {
  const t = e.target as Node
  if (rootEl.value?.contains(t) || panelEl.value?.contains(t)) return
  isOpen.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') isOpen.value = false
}

onMounted(() => {
  document.addEventListener('mousedown', onOutsideClick)
  document.addEventListener('keydown', onKeydown)
})
onUnmounted(() => {
  document.removeEventListener('mousedown', onOutsideClick)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
/* ── Root ── */
.dp-root {
  position: relative;
  display: block;
}

/* ── Trigger (dark default) ── */
.dp-trigger {
  @apply w-full flex items-center gap-2 px-3 py-2
    bg-gray-800/80 border border-gray-700 rounded-lg
    text-sm transition-all duration-150 cursor-pointer select-none
    focus:outline-none;
}
.dp-trigger:hover { @apply border-gray-600; }
.dp-trigger--open  { @apply border-violet-500 ring-2 ring-violet-500/25; }
.dp-trigger--filled .dp-trigger-value { @apply text-white; }
.dp-trigger--empty  .dp-trigger-value { @apply text-gray-500; }

.dp-trigger-icon {
  @apply w-3.5 h-3.5 shrink-0 text-gray-500;
}
.dp-trigger-value {
  @apply flex-1 text-left text-sm truncate;
}
.dp-chevron {
  @apply w-4 h-4 shrink-0 text-slate-500 transition-transform duration-200;
}

.dp-clear-btn {
  @apply w-5 h-5 rounded-full flex items-center justify-center shrink-0
    text-slate-500 hover:text-slate-200 hover:bg-slate-600/70 transition-colors;
}

/* ── Panel (dark default) ── */
.dp-panel {
  @apply rounded-2xl overflow-y-auto;
  background: rgba(10, 14, 26, 0.97);
  border: 1px solid rgba(148, 163, 184, 0.12);
  box-shadow:
    0 32px 64px rgba(0, 0, 0, 0.72),
    0 8px 24px rgba(0, 0, 0, 0.48),
    0 0 0 1px rgba(255, 255, 255, 0.04) inset;
  backdrop-filter: blur(20px) saturate(150%);
}

/* ── Panel (light mode) — Arctic Violet ── */
.dp-panel--light {
  background: #FFFFFF !important;
  border: 1px solid #E2E8F6 !important;
  box-shadow:
    0 24px 64px rgba(14,17,36,.12),
    0 8px 24px rgba(14,17,36,.08),
    0 0 0 1px rgba(124,58,237,.05) inset !important;
  backdrop-filter: none !important;
}

/* ── Navigation (dark default) ── */
.dp-nav {
  @apply flex items-center gap-1 px-3 pt-4 pb-3;
}
.dp-nav-btn {
  @apply w-8 h-8 rounded-lg flex items-center justify-center shrink-0
    text-slate-400 hover:text-white hover:bg-slate-700/60
    transition-colors duration-150;
}
.dp-nav-label {
  @apply flex-1 text-center text-sm font-semibold text-white tracking-wide
    hover:text-violet-300 transition-colors duration-150 cursor-pointer;
}

/* Navigation light overrides */
.dp-panel--light .dp-nav-btn {
  color: #6B7A9A !important;
}
.dp-panel--light .dp-nav-btn:hover {
  color: #0E1117 !important;
  background-color: #F7F9FE !important;
}
.dp-panel--light .dp-nav-label {
  color: #0E1117 !important;
  font-weight: 600 !important;
}
.dp-panel--light .dp-nav-label:hover {
  color: #7C3AED !important;
}

/* ── Shortcuts (dark default) ── */
.dp-shortcuts {
  @apply flex gap-1.5 px-3 pb-3;
  border-bottom: 1px solid rgba(148, 163, 184, 0.1);
}
.dp-shortcut-btn {
  @apply text-xs font-medium px-2.5 py-1 rounded-md
    text-slate-400 border border-slate-700/60
    hover:text-violet-300 hover:border-violet-500/50 hover:bg-violet-500/10
    transition-all duration-150 whitespace-nowrap;
}

/* Shortcuts light overrides */
.dp-panel--light .dp-shortcuts {
  border-bottom: 1px solid #E2E8F6 !important;
}
.dp-panel--light .dp-shortcut-btn {
  color: #6B7A9A !important;
  border-color: #E2E8F6 !important;
  background: transparent !important;
}
.dp-panel--light .dp-shortcut-btn:hover {
  color: #7C3AED !important;
  border-color: #C4B5FD !important;
  background: #F5F3FF !important;
}

/* ── Weekday Row ── */
.dp-weekday-row {
  @apply grid grid-cols-7 px-3 pt-3 pb-1;
}
.dp-weekday-cell {
  @apply text-center text-xs font-medium tracking-wide text-slate-600 py-1;
}
.dp-panel--light .dp-weekday-cell {
  color: #9DA6BD !important;
  font-weight: 600 !important;
  letter-spacing: 0.04em !important;
}

/* ── Day Grid ── */
.dp-days-grid {
  @apply grid grid-cols-7 px-3 pb-3 gap-y-0.5;
}

.dp-day-btn {
  @apply relative flex flex-col items-center justify-center
    h-9 w-full rounded-lg text-sm font-medium
    transition-all duration-100 leading-none;
}
.dp-day-num { @apply leading-none; }

/* Day states — dark */
.dp-day--current      { @apply text-slate-200 hover:bg-slate-700/60 hover:text-white; }
.dp-day--current.dp-day--weekend { @apply text-slate-400; }
.dp-day--today        { @apply text-white font-semibold; }
.dp-today-dot {
  @apply absolute bottom-1 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full bg-violet-400;
}
.dp-day--other        { @apply text-slate-700 hover:text-slate-500 hover:bg-slate-800/40; }
.dp-day--disabled     { @apply text-slate-700 cursor-not-allowed line-through; }
.dp-day--selected {
  @apply text-white font-semibold;
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
  box-shadow: 0 0 0 1px rgba(139,92,246,.4), 0 4px 12px rgba(109,40,217,.45);
}
.dp-day--selected:hover { background: linear-gradient(135deg, #8b5cf6, #7c3aed); }

/* Day states — light (Arctic Violet) */
.dp-panel--light .dp-day--current {
  color: #374261 !important;
}
.dp-panel--light .dp-day--current:hover {
  background-color: #F2F4F9 !important;
  color: #0E1117 !important;
}
.dp-panel--light .dp-day--current.dp-day--weekend {
  color: #6B7A9A !important;
}
.dp-panel--light .dp-day--today {
  color: #7C3AED !important;
  font-weight: 700 !important;
}
.dp-panel--light .dp-today-dot {
  background-color: #7C3AED !important;
}
.dp-panel--light .dp-day--other {
  color: #C8CDD9 !important;
}
.dp-panel--light .dp-day--other:hover {
  color: #9DA6BD !important;
  background-color: #F7F9FE !important;
}
.dp-panel--light .dp-day--disabled {
  color: #D2D8EC !important;
}

/* ── Footer (dark default) ── */
.dp-footer {
  @apply flex items-center gap-2 px-4 py-2.5;
  border-top: 1px solid rgba(148, 163, 184, 0.08);
  background: rgba(255, 255, 255, 0.02);
}
.dp-footer-label { @apply flex-1 text-xs text-slate-400 font-medium; }
.dp-footer-clear  { @apply text-xs text-slate-600 hover:text-red-400 transition-colors font-medium; }

/* Footer light overrides */
.dp-panel--light .dp-footer {
  border-top: 1px solid #EBF0FA !important;
  background: #F7F9FE !important;
}
.dp-panel--light .dp-footer-label { color: #6B7A9A !important; }
.dp-panel--light .dp-footer-clear {
  color: #9DA6BD !important;
}
.dp-panel--light .dp-footer-clear:hover { color: #DC2626 !important; }

/* ── Transition ── */
.dp-fade-enter-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.dp-fade-leave-active { transition: opacity 0.1s ease, transform 0.1s ease; }
.dp-fade-enter-from  { opacity: 0; transform: translateY(-6px) scale(0.97); }
.dp-fade-leave-to    { opacity: 0; transform: translateY(-4px) scale(0.98); }
</style>
