<template>
  <div
    v-if="totalWidth > 0 && sortedMilestones.length > 0"
    ref="rowRef"
    class="gantt-milestone-row flex items-center border-b border-slate-600/60 bg-slate-800/90 h-11 shrink-0 relative"
    :style="{ width: totalWidth + 'px', minWidth: totalWidth + 'px' }"
  >
    <div
      v-for="m in sortedMilestones"
      :key="m.id"
      class="milestone-marker absolute flex flex-col items-center -translate-x-1/2 z-10 cursor-grab active:cursor-grabbing select-none group"
      :class="{ 'cursor-grabbing': draggingMilestone?.id === m.id }"
      :style="{ left: (draggingMilestone?.id === m.id ? dragLeftPx : leftPx(m)) + 'px' }"
      @mousedown.prevent="onMilestoneMouseDown($event, m)"
      @touchstart.passive="onMilestoneTouchStart($event, m)"
      @click.stop="onMilestoneClick(m)"
    >
      <div
        class="milestone-diamond w-4 h-4 rotate-45 border-2 shadow-sm transition-all flex-shrink-0 group-hover:scale-110"
        :class="markerClass(m)"
        :title="m.title + (m.due_date ? ' — ' + formatDate(m.due_date) : '') + ' — Drag to change date'"
      />
      <span class="text-[10px] font-medium text-slate-400 mt-1 whitespace-nowrap max-w-[90px] truncate group-hover:max-w-none group-hover:text-slate-200">
        {{ m.title }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Milestone } from '~/core/modules/projects/infrastructure/projects-api'

const props = defineProps<{
  milestones: Milestone[]
  dateRangeStart: string
  dateRangeEnd: string
  gridWidth: number
  gridOffset?: number
}>()

const totalWidth = computed(() => props.gridWidth + (props.gridOffset ?? 0))

const emit = defineEmits<{
  (e: 'milestone-click', m: Milestone): void
  (e: 'milestone-drag-move', payload: { milestoneId: string; leftPx: number }): void
  (e: 'milestone-drag-end', payload: { milestone: Milestone; newDueDate: string }): void
}>()

const rowRef = ref<HTMLElement | null>(null)
const draggingMilestone = ref<Milestone | null>(null)
const dragLeftPx = ref(0)
const rowRectLeft = ref(0)
const justFinishedDrag = ref(false)

const sortedMilestones = computed(() =>
  [...props.milestones]
    .filter((m) => m.due_date)
    .sort((a, b) => new Date(a.due_date!).getTime() - new Date(b.due_date!).getTime())
)

/** Parse YYYY-MM-DD (or ISO date part) as local midnight to match Vue-Ganttastic/dayjs. */
function toLocalMidnight(isoOrYmd: string): number {
  const s = isoOrYmd.split('T')[0]
  const [y, m, d] = s.split('-').map(Number)
  return new Date(y, m - 1, d).getTime()
}

const gridOffset = computed(() => props.gridOffset ?? 0)

function leftPx(m: Milestone): number {
  if (!m.due_date) return 0
  const start = toLocalMidnight(props.dateRangeStart)
  const end = toLocalMidnight(props.dateRangeEnd)
  const date = toLocalMidnight(m.due_date)
  if (end <= start) return 0
  const pct = Math.max(0, Math.min(1, (date - start) / (end - start)))
  return gridOffset.value + pct * props.gridWidth
}

/** Inverse of leftPx: pixel position -> YYYY-MM-DD string (noon UTC). */
function pxToDate(px: number): string {
  const start = toLocalMidnight(props.dateRangeStart)
  const end = toLocalMidnight(props.dateRangeEnd)
  if (end <= start) return props.dateRangeStart.split('T')[0]
  const clamped = Math.max(gridOffset.value, Math.min(gridOffset.value + props.gridWidth, px))
  const pct = (clamped - gridOffset.value) / props.gridWidth
  const ts = start + pct * (end - start)
  const d = new Date(ts)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function onMilestoneClick(m: Milestone) {
  if (draggingMilestone.value || justFinishedDrag.value) return
  emit('milestone-click', m)
}

function onMilestoneMouseDown(e: MouseEvent, m: Milestone) {
  if (!m.due_date || !rowRef.value) return
  const rect = rowRef.value.getBoundingClientRect()
  rowRectLeft.value = rect.left
  draggingMilestone.value = m
  dragLeftPx.value = leftPx(m)
  const onMove = (e2: MouseEvent) => {
    const left = e2.clientX - rowRectLeft.value
    dragLeftPx.value = Math.max(gridOffset.value, Math.min(gridOffset.value + props.gridWidth, left))
    emit('milestone-drag-move', { milestoneId: m.id, leftPx: dragLeftPx.value })
  }
  const onUp = () => {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    justFinishedDrag.value = true
    setTimeout(() => { justFinishedDrag.value = false }, 200)
    const newDueDate = pxToDate(dragLeftPx.value)
    const currentYmd = m.due_date?.split('T')[0]
    if (newDueDate !== currentYmd) emit('milestone-drag-end', { milestone: m, newDueDate })
    draggingMilestone.value = null
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

function onMilestoneTouchStart(e: TouchEvent, m: Milestone) {
  if (!m.due_date || !rowRef.value || e.touches.length !== 1) return
  const rect = rowRef.value.getBoundingClientRect()
  rowRectLeft.value = rect.left
  draggingMilestone.value = m
  dragLeftPx.value = leftPx(m)
  const onMove = (e2: TouchEvent) => {
    if (e2.touches.length !== 1) return
    e2.preventDefault()
    const left = e2.touches[0].clientX - rowRectLeft.value
    dragLeftPx.value = Math.max(gridOffset.value, Math.min(gridOffset.value + props.gridWidth, left))
    emit('milestone-drag-move', { milestoneId: m.id, leftPx: dragLeftPx.value })
  }
  const onEnd = () => {
    document.removeEventListener('touchmove', onMove, { capture: true })
    document.removeEventListener('touchend', onEnd)
    justFinishedDrag.value = true
    setTimeout(() => { justFinishedDrag.value = false }, 200)
    const newDueDate = pxToDate(dragLeftPx.value)
    const currentYmd = m.due_date?.split('T')[0]
    if (newDueDate !== currentYmd) emit('milestone-drag-end', { milestone: m, newDueDate })
    draggingMilestone.value = null
  }
  document.addEventListener('touchmove', onMove, { capture: true, passive: false })
  document.addEventListener('touchend', onEnd)
}

function markerClass(m: Milestone) {
  if (m.status === 'REACHED') return 'bg-emerald-500 border-emerald-400 shadow-md'
  if (m.status === 'MISSED') return 'bg-rose-500 border-rose-400'
  const isOverdue = m.due_date && new Date(m.due_date) < new Date()
  if (isOverdue) return 'bg-rose-500/20 border-rose-400'
  return 'bg-slate-800 border-purple-400/90 group-hover:border-purple-300 group-hover:bg-purple-500/60'
}
</script>

<style scoped>
.milestone-diamond {
  border-radius: 2px;
}
</style>
