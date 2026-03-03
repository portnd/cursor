<template>
  <div
    v-if="totalWidth > 0 && sortedMilestones.length > 0"
    class="gantt-milestone-row flex items-center border-b border-slate-600/60 bg-slate-800/90 h-11 shrink-0 relative"
    :style="{ width: totalWidth + 'px', minWidth: totalWidth + 'px' }"
  >
    <div
      v-for="m in sortedMilestones"
      :key="m.id"
      class="absolute flex flex-col items-center -translate-x-1/2 z-10 cursor-pointer group"
      :style="{ left: leftPx(m) + 'px' }"
      @click="$emit('milestone-click', m)"
    >
      <div
        class="milestone-diamond w-4 h-4 rotate-45 border-2 shadow-sm transition-all flex-shrink-0 group-hover:scale-110"
        :class="markerClass(m)"
        :title="m.title + (m.due_date ? ' — ' + formatDate(m.due_date) : '')"
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

defineEmits<{
  (e: 'milestone-click', m: Milestone): void
}>()

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

function leftPx(m: Milestone): number {
  if (!m.due_date) return 0
  const start = toLocalMidnight(props.dateRangeStart)
  const end = toLocalMidnight(props.dateRangeEnd)
  const date = toLocalMidnight(m.due_date)
  if (end <= start) return 0
  const pct = Math.max(0, Math.min(1, (date - start) / (end - start)))
  const offset = props.gridOffset ?? 0
  return offset + pct * props.gridWidth
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function markerClass(m: Milestone) {
  if (m.status === 'REACHED') return 'bg-emerald-500 border-emerald-400 shadow-md'
  if (m.status === 'MISSED') return 'bg-rose-500 border-rose-400'
  const isOverdue = m.due_date && new Date(m.due_date) < new Date()
  if (isOverdue) return 'bg-rose-500/20 border-rose-400'
  return 'bg-slate-800 border-indigo-400/90 group-hover:border-indigo-300 group-hover:bg-indigo-500/60'
}
</script>

<style scoped>
.milestone-diamond {
  border-radius: 2px;
}
</style>
