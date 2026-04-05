<template>
  <div class="kanban-board w-full min-w-0">

    <!-- ── Active Sprint Blinder (DEV role only) ──────────────────────────── -->
    <div v-if="activeSprint && isDev" class="relative overflow-hidden rounded-2xl border border-purple-500/40 bg-gradient-to-r from-purple-950/60 via-gray-900/80 to-gray-900/60 px-6 py-5 mb-4 shadow-lg shadow-purple-500/10">
      <div class="pointer-events-none absolute -top-10 -left-10 h-40 w-40 rounded-full bg-purple-600/10 blur-3xl" />
      <div class="relative flex flex-col sm:flex-row sm:items-center gap-4">
        <div class="flex items-center gap-3 min-w-0 flex-1">
          <div class="shrink-0 w-9 h-9 rounded-xl bg-purple-500/20 border border-purple-500/40 flex items-center justify-center">
            <svg class="w-4 h-4 text-purple-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
          </div>
          <div class="min-w-0">
            <div class="flex items-center gap-2 mb-0.5">
              <span class="text-[10px] font-bold uppercase tracking-widest text-purple-400">Active Sprint</span>
              <span class="inline-flex h-1.5 w-1.5 rounded-full bg-green-400 animate-pulse" />
            </div>
            <h3 class="text-sm font-bold text-white truncate">{{ activeSprint.name }}</h3>
            <p v-if="activeSprint.goal" class="text-xs text-gray-400 mt-0.5 line-clamp-1">{{ activeSprint.goal }}</p>
          </div>
        </div>
        <div v-if="activeSprint.end_date" class="shrink-0 text-right">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500">Sprint Ends</p>
          <p class="text-sm font-bold text-amber-400 mt-0.5">{{ formatEndDate(activeSprint.end_date) }}</p>
          <p class="text-xs text-gray-500 mt-0.5">{{ getCountdown(activeSprint.end_date) }}</p>
        </div>
      </div>
    </div>

    <!-- Filter Bar (responsive: stack on narrow, wrap on wide) -->
    <!-- Sprint selector is hidden for DEV role — they are locked to their active sprint -->
    <div class="flex flex-col sm:flex-row sm:flex-wrap items-stretch sm:items-center gap-3 mb-4 sm:mb-5 p-3 bg-gray-800/50 rounded-xl border border-gray-700/50">
      <div v-if="!isDev" class="flex flex-wrap items-center gap-2 sm:gap-2">
        <label class="text-xs text-gray-400 uppercase tracking-wide font-medium w-14 sm:w-auto shrink-0">Sprint</label>
        <select v-model="filterSprint" class="input-select text-sm min-w-0 flex-1 sm:min-w-[160px] sm:flex-none max-w-full">
          <option value="">All Sprints</option>
          <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
        <span v-if="sprints.length === 0" class="text-xs text-gray-500 w-full sm:w-auto">— No sprints yet</span>
        <span v-else-if="filterSprint" class="text-xs text-purple-400 truncate max-w-full">{{ sprintNameById(filterSprint) || 'Sprint' }}</span>
      </div>
      <div v-if="isDev && activeSprint" class="flex items-center gap-2">
        <span class="text-xs text-gray-400 uppercase tracking-wide font-medium w-14 sm:w-auto shrink-0">Sprint</span>
        <span class="text-sm font-semibold text-purple-300 px-3 py-1.5 rounded-lg bg-purple-500/10 border border-purple-500/30">{{ activeSprint.name }}</span>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <label class="text-xs text-gray-400 uppercase tracking-wide w-14 sm:w-auto shrink-0">Priority</label>
        <select v-model="filterPriority" class="input-select text-sm min-w-0 flex-1 sm:flex-none">
          <option value="">All</option>
          <option value="CRITICAL">Critical</option>
          <option value="HIGH">High</option>
          <option value="MEDIUM">Medium</option>
          <option value="LOW">Low</option>
        </select>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <label class="text-xs text-gray-400 uppercase tracking-wide w-20 sm:w-auto shrink-0">Swim Lane</label>
        <select v-model="swimLane" class="input-select text-sm min-w-0 flex-1 sm:flex-none">
          <option value="none">None</option>
          <option value="priority">By Priority</option>
          <option value="assignee">By Assignee</option>
        </select>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <label class="text-xs text-gray-400 uppercase tracking-wide w-8 sm:w-auto shrink-0">Type</label>
        <select v-model="filterType" class="input-select text-sm min-w-0 flex-1 sm:flex-none">
          <option value="">All</option>
          <option value="FEATURE">★ Feature</option>
          <option value="TASK">📋 Task</option>
          <option value="BUG">⚠ Bug</option>
        </select>
      </div>
      <div class="text-xs text-gray-500 sm:ml-auto shrink-0">
        {{ filteredTasks.length }} tasks
      </div>
    </div>

    <!-- WIP Limit Warning -->
    <div v-if="wipWarning" class="mb-4 flex items-center gap-2 px-3 sm:px-4 py-2 bg-yellow-500/10 border border-yellow-500/30 rounded-lg text-yellow-400 text-xs sm:text-sm">
      <span>⚠️</span>
      <span class="min-w-0">WIP limit exceeded in <strong>{{ wipWarning }}</strong></span>
    </div>

    <!-- Kanban Columns: horizontal scroll on small screens, grid on large -->
    <div class="flex lg:grid lg:grid-cols-6 gap-3 overflow-x-auto lg:overflow-visible pb-2 -mx-1 px-1 lg:mx-0 lg:px-0 snap-x snap-mandatory scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-transparent">
      <div
        v-for="col in columns"
        :key="col.status"
        class="kanban-col min-w-[200px] sm:min-w-[220px] w-[200px] sm:w-[220px] lg:w-auto shrink-0 lg:shrink snap-start"
        @dragover.prevent="col.droppable ? onDragOver($event, col.status) : undefined"
        @dragleave="onDragLeave"
        @drop="col.droppable ? onDrop($event, col.status) : undefined"
        :class="{ 'drop-target': dragOverCol === col.status && col.droppable, 'col-locked': !col.droppable }"
      >
        <!-- Column Header -->
        <div class="flex items-center justify-between mb-3 px-1">
          <div class="flex items-center gap-2">
            <span class="text-lg">{{ col.icon }}</span>
            <span class="text-sm font-semibold" :class="col.headerClass">{{ col.label }}</span>
            <!-- Lock badge for non-droppable columns -->
            <span v-if="!col.droppable" title="Cannot drag tasks here" class="text-[10px] text-gray-500 px-1.5 py-0.5 rounded bg-gray-800 border border-gray-700">🔒</span>
          </div>
          <div class="flex items-center gap-2">
            <span
              class="badge-count"
              :class="isWipExceeded(col.status) ? 'bg-red-500/20 text-red-400 border-red-500/30' : 'bg-gray-700 text-gray-400'"
              :title="col.wipLimit ? `${tasksByCol[col.status]?.length || 0} of ${col.wipLimit} max (WIP limit)` : `${tasksByCol[col.status]?.length || 0} tasks`"
            >
              {{ tasksByCol[col.status]?.length || 0 }}
              <span v-if="col.wipLimit" class="text-gray-500">/{{ col.wipLimit }}</span>
            </span>
          </div>
        </div>

        <!-- Task Cards (inline to avoid runtime template compiler) -->
        <div class="flex flex-col gap-2 min-h-[120px]">
          <template v-if="swimLane === 'none'">
            <div
              v-for="task in tasksByCol[col.status]"
              :key="task.id"
              class="kanban-card group"
              draggable="true"
              @dragstart="onDragStart($event, task)"
              @click="$emit('task-click', task)"
            >
              <div class="flex items-start justify-between gap-1 mb-2">
                <div class="flex items-center gap-1.5 min-w-0">
                  <span class="text-xs text-gray-500 font-mono shrink-0">{{ taskDisplayCode(task) }}</span>
                  <span class="shrink-0 text-[11px] font-semibold px-1.5 py-0.5 rounded-full border" :class="taskTypeCls(task.task_type)">
                    {{ taskTypeIcon(task.task_type) }} {{ task.task_type || 'TASK' }}
                  </span>
                </div>
                <span class="text-xs border rounded px-1 shrink-0" :class="priorityCls(task.priority)">{{ task.priority }}</span>
              </div>
              <p v-if="task.sprint_id && sprintNameById(task.sprint_id)" class="text-[10px] text-purple-400 mb-1 truncate" :title="sprintNameById(task.sprint_id)">📌 {{ sprintNameById(task.sprint_id) }}</p>
              <p class="text-sm text-gray-200 font-medium leading-snug mb-2 line-clamp-2">{{ task.title }}</p>
              <!-- WAIT_FOR_DEPLOY: no deployment request yet -->
              <div
                v-if="col.status === 'WAIT_FOR_DEPLOY' && !(props.deployedTaskIds ?? []).includes(task.id)"
                class="flex items-center gap-1.5 mb-2 px-2 py-1 rounded-md bg-orange-500/10 border border-orange-500/25 text-[10px] text-orange-400 font-medium"
              >
                <span>⚠️</span> No deployment request
              </div>
              <div class="flex items-center justify-between text-xs text-gray-500">
                <span v-if="task.story_points" class="flex items-center gap-1"><span class="text-purple-400">◆</span> {{ task.story_points }} SP</span>
                <span v-if="assigneeLabel(task)" class="flex items-center gap-1" :title="assigneeLabel(task)"><span class="w-4 h-4 rounded-full bg-purple-600 flex items-center justify-center text-white text-[10px] font-bold">{{ assigneeInitial(task) }}</span></span>
                <span v-if="daysUntilDue(task) !== null" :class="daysUntilDue(task)! < 0 ? 'text-red-400' : daysUntilDue(task)! <= 2 ? 'text-yellow-400' : 'text-gray-500'">
                  {{ daysUntilDue(task)! < 0 ? Math.abs(daysUntilDue(task)!) + 'd overdue' : daysUntilDue(task) + 'd left' }}
                </span>
              </div>
              <div v-if="task.progress > 0" class="mt-2">
                <div class="h-1 bg-gray-700 rounded-full overflow-hidden">
                  <div class="h-full bg-purple-500 rounded-full transition-all" :style="{ width: task.progress + '%' }"></div>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <div v-for="lane in getLanes(col.status)" :key="lane.key" class="mb-2">
              <div class="text-xs text-gray-500 px-1 mb-1 font-medium">{{ lane.label }}</div>
              <div
                v-for="task in lane.tasks"
                :key="task.id"
                class="kanban-card group"
                draggable="true"
                @dragstart="onDragStart($event, task)"
                @click="$emit('task-click', task)"
              >
                <div class="flex items-start justify-between gap-1 mb-2">
                  <div class="flex items-center gap-1.5 min-w-0">
                    <span class="text-xs text-gray-500 font-mono shrink-0">{{ taskDisplayCode(task) }}</span>
                    <span class="shrink-0 text-[11px] font-semibold px-1.5 py-0.5 rounded-full border" :class="taskTypeCls(task.task_type)">
                      {{ taskTypeIcon(task.task_type) }} {{ task.task_type || 'TASK' }}
                    </span>
                  </div>
                  <span class="text-xs border rounded px-1 shrink-0" :class="priorityCls(task.priority)">{{ task.priority }}</span>
                </div>
                <p v-if="task.sprint_id && sprintNameById(task.sprint_id)" class="text-[10px] text-purple-400 mb-1 truncate">📌 {{ sprintNameById(task.sprint_id) }}</p>
                <p class="text-sm text-gray-200 font-medium leading-snug mb-2 line-clamp-2">{{ task.title }}</p>
                <!-- WAIT_FOR_DEPLOY: no deployment request yet -->
                <div
                  v-if="col.status === 'WAIT_FOR_DEPLOY' && !(props.deployedTaskIds ?? []).includes(task.id)"
                  class="flex items-center gap-1.5 mb-2 px-2 py-1 rounded-md bg-orange-500/10 border border-orange-500/25 text-[10px] text-orange-400 font-medium"
                >
                  <span>⚠️</span> No deployment request
                </div>
                <div class="flex items-center justify-between text-xs text-gray-500">
                  <span v-if="task.story_points" class="flex items-center gap-1"><span class="text-purple-400">◆</span> {{ task.story_points }} SP</span>
                  <span v-if="assigneeLabel(task)" class="flex items-center gap-1" :title="assigneeLabel(task)"><span class="w-4 h-4 rounded-full bg-purple-600 flex items-center justify-center text-white text-[10px] font-bold">{{ assigneeInitial(task) }}</span></span>
                  <span v-if="daysUntilDue(task) !== null" :class="daysUntilDue(task)! < 0 ? 'text-red-400' : daysUntilDue(task)! <= 2 ? 'text-yellow-400' : 'text-gray-500'">
                    {{ daysUntilDue(task)! < 0 ? Math.abs(daysUntilDue(task)!) + 'd overdue' : daysUntilDue(task) + 'd left' }}
                  </span>
                </div>
                <div v-if="task.progress > 0" class="mt-2">
                  <div class="h-1 bg-gray-700 rounded-full overflow-hidden">
                    <div class="h-full bg-purple-500 rounded-full transition-all" :style="{ width: task.progress + '%' }"></div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Empty State -->
          <div
            v-if="!tasksByCol[col.status]?.length"
            class="flex flex-col items-center justify-center h-20 border-2 border-dashed border-gray-700 rounded-lg text-gray-600 text-xs"
          >
            <template v-if="col.droppable">Drop tasks here</template>
            <template v-else>—</template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { isEngineerLikeRole } from '~/utils/roles'
import type { Task } from '~/core/modules/projects/infrastructure/projects-api'
import type { Sprint } from '~/core/modules/projects/infrastructure/projects-api'

interface ActiveSprintInfo {
  id: string
  name: string
  goal: string
  end_date: string | null
  status: string
}

const props = defineProps<{
  tasks: Task[]
  sprints: Sprint[]
  /** Optional map taskId -> display code (e.g. 001, 002) from project backlog order */
  taskDisplayCodeMap?: Record<string, string>
  /** Current user role — ENGINEER activates blinder mode */
  userRole?: string
  /** Active sprint info for the current user (from DevView or project page) */
  activeSprint?: ActiveSprintInfo | null
  /** Set of task IDs that already have a deployment request; used to show warning on WAIT_FOR_DEPLOY cards */
  deployedTaskIds?: string[]
}>()

const isDev = computed(() => isEngineerLikeRole(props.userRole))

const emit = defineEmits<{
  (e: 'task-click', task: Task): void
  (e: 'status-change', taskId: string, status: string): void
}>()

const filterSprint = ref('')
const filterPriority = ref('')
const filterType = ref('')
const swimLane = ref<'none' | 'priority' | 'assignee'>('none')
const dragTask = ref<Task | null>(null)
const dragOverCol = ref('')

function formatEndDate(endDate: string) {
  return new Date(endDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function getCountdown(endDate: string) {
  const diff = new Date(endDate).getTime() - Date.now()
  if (diff < 0) return 'Sprint overdue'
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  if (days === 0) return 'Ends today'
  if (days === 1) return '1 day left'
  return `${days} days left`
}

function sprintNameById(id: string) {
  return props.sprints.find((s) => s.id === id)?.name ?? ''
}

const priorityConfig: Record<string, string> = {
  CRITICAL: 'text-red-400 bg-red-500/10 border-red-500/30',
  HIGH: 'text-orange-400 bg-orange-500/10 border-orange-500/30',
  MEDIUM: 'text-yellow-400 bg-yellow-500/10 border-yellow-500/30',
  LOW: 'text-green-400 bg-green-500/10 border-green-500/30',
}
function priorityCls(p: string) {
  return priorityConfig[p] ?? 'text-gray-400 bg-gray-700 border-gray-600'
}

const taskTypeConfig: Record<string, string> = {
  FEATURE: 'text-purple-300 bg-purple-500/15 border-purple-500/40',
  TASK: 'text-blue-300 bg-blue-500/15 border-blue-500/30',
  BUG: 'text-red-300 bg-red-500/15 border-red-500/30',
}
function taskTypeCls(type: string) {
  return taskTypeConfig[type] ?? taskTypeConfig['TASK']
}
function taskTypeIcon(type: string) {
  if (type === 'FEATURE') return '★'
  if (type === 'BUG') return '⚠'
  return '📋'
}

function taskDisplayCode(task: Task): string {
  if (props.taskDisplayCodeMap?.[task.id]) return props.taskDisplayCodeMap[task.id]
  if (!task.code) return '–'
  const suffix = task.code.split('-').pop()
  return /^\d+$/.test(suffix || '') ? String(Number(suffix)).padStart(4, '0') : task.code
}
function daysUntilDue(task: Task): number | null {
  if (!task.due_at) return null
  const diff = new Date(task.due_at).getTime() - Date.now()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

/** Show assignee name: display_name, or email local part (e.g. sirin.s@komgrip.com → sirin), or Dev #id */
function assigneeLabel(task: Task): string {
  if (!task.assigned_to && !task.assigned_to_display_name && !task.assigned_to_email) return ''
  if (task.assigned_to_display_name) return task.assigned_to_display_name
  if (task.assigned_to_email) {
    const local = task.assigned_to_email.split('@')[0] || ''
    return local.split('.')[0] || local || task.assigned_to_email
  }
  return task.assigned_to != null ? `Dev #${task.assigned_to}` : ''
}
function assigneeInitial(task: Task): string {
  const label = assigneeLabel(task)
  return label ? label.charAt(0).toUpperCase() : '?'
}

const columns = [
  { status: 'PENDING',         label: 'To Do',           icon: '📋', headerClass: 'text-gray-300',   wipLimit: 0,  droppable: true },
  { status: 'IN_PROGRESS',     label: 'In Progress',     icon: '⚡', headerClass: 'text-blue-400',   wipLimit: 5,  droppable: true },
  { status: 'READY_FOR_TEST',  label: 'Ready for Test',  icon: '🧪', headerClass: 'text-cyan-400',   wipLimit: 3,  droppable: true },
  { status: 'WAIT_FOR_DEPLOY', label: 'Wait for Deploy', icon: '🚀', headerClass: 'text-orange-400', wipLimit: 0,  droppable: true },
  { status: 'READY_FOR_UAT',   label: 'Ready for UAT',   icon: '🔬', headerClass: 'text-amber-400',  wipLimit: 0,  droppable: false }, // set automatically on deployment
  { status: 'COMPLETED',       label: 'Done',            icon: '✅', headerClass: 'text-green-400',  wipLimit: 0,  droppable: false }, // CEO/MANAGER approve only — no drag
]

const filteredTasks = computed(() =>
  props.tasks.filter((t) => {
    // DEV role is locked to active sprint — only show tasks from that sprint
    if (isDev.value && props.activeSprint) {
      if (t.sprint_id !== props.activeSprint.id) return false
    } else if (filterSprint.value && t.sprint_id !== filterSprint.value) {
      return false
    }
    if (filterPriority.value && t.priority !== filterPriority.value) return false
    if (filterType.value && (t.task_type || 'TASK') !== filterType.value) return false
    return true
  })
)

const COLUMN_STATUSES = ['PENDING', 'IN_PROGRESS', 'READY_FOR_TEST', 'WAIT_FOR_DEPLOY', 'READY_FOR_UAT', 'COMPLETED'] as const

/** Maps a task to a kanban column key. */
function bucketForTask(t: Task): string {
  if (t.status === 'WAIT_FOR_DEPLOY') return 'WAIT_FOR_DEPLOY'
  if (t.status === 'READY_FOR_UAT') return 'READY_FOR_UAT'
  if (t.status === 'REVIEW_PENDING') {
    return t.task_type === 'FEATURE' ? 'READY_FOR_UAT' : 'READY_FOR_TEST'
  }
  if (t.status === 'READY_FOR_TEST') return 'READY_FOR_TEST'
  if (COLUMN_STATUSES.includes(t.status as (typeof COLUMN_STATUSES)[number])) return t.status
  return 'PENDING'
}

/** Drop target status for API.
 *  Returns null when the column is non-droppable (COMPLETED / READY_FOR_UAT).
 */
function resolvedStatusForColumn(colStatus: string, task: Task): string | null {
  const col = columns.find((c) => c.status === colStatus)
  if (col && !col.droppable) return null
  if (colStatus === 'READY_FOR_TEST') {
    return task.task_type === 'FEATURE' ? 'REVIEW_PENDING' : 'READY_FOR_TEST'
  }
  return colStatus
}

const tasksByCol = computed(() => {
  const map: Record<string, Task[]> = {}
  for (const col of columns) map[col.status] = []
  for (const t of filteredTasks.value) {
    const key = bucketForTask(t)
    if (map[key]) map[key].push(t)
  }
  return map
})

const wipWarning = computed(() => {
  for (const col of columns) {
    if (col.wipLimit > 0 && (tasksByCol.value[col.status]?.length || 0) > col.wipLimit) {
      return col.label
    }
  }
  return null
})

function isWipExceeded(status: string) {
  const col = columns.find((c) => c.status === status)
  if (!col || !col.wipLimit) return false
  return (tasksByCol.value[status]?.length || 0) > col.wipLimit
}

function getLanes(status: string) {
  const tasks = tasksByCol.value[status] || []
  if (swimLane.value === 'priority') {
    const priorities = ['CRITICAL', 'HIGH', 'MEDIUM', 'LOW']
    return priorities
      .map((p) => ({ key: p, label: p, tasks: tasks.filter((t) => t.priority === p) }))
      .filter((l) => l.tasks.length > 0)
  }
  if (swimLane.value === 'assignee') {
    const assigneeKeys = [...new Set(tasks.map((t) => assigneeLabel(t) || 'Unassigned'))]
    return assigneeKeys.map((key) => ({
      key,
      label: key,
      tasks: tasks.filter((t) => (assigneeLabel(t) || 'Unassigned') === key),
    }))
  }
  return []
}

function onDragStart(e: DragEvent, task: Task) {
  dragTask.value = task
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', task.id)
  }
}

function onDragOver(e: DragEvent, status: string) {
  dragOverCol.value = status
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
}

function onDragLeave() {
  dragOverCol.value = ''
}

function onDrop(e: DragEvent, colStatus: string) {
  dragOverCol.value = ''
  const task = dragTask.value
  if (!task) return
  const target = resolvedStatusForColumn(colStatus, task)
  // null means column is protected (COMPLETED / READY_FOR_UAT) — silently cancel
  if (target === null || task.status === target) {
    dragTask.value = null
    return
  }
  emit('status-change', task.id, target)
  dragTask.value = null
}

function isColumnDroppable(colStatus: string): boolean {
  return columns.find((c) => c.status === colStatus)?.droppable ?? true
}
</script>

<style scoped>
.kanban-col {
  @apply bg-gray-800/50 rounded-xl p-3 flex flex-col border border-gray-700/50 transition-all;
}

.drop-target {
  @apply border-purple-500/60 bg-purple-500/5;
}

.col-locked {
  @apply opacity-80 cursor-default;
}

.kanban-card {
  @apply bg-gray-900 border border-gray-700 rounded-lg p-3 cursor-grab active:cursor-grabbing hover:border-purple-500/50 hover:shadow-lg hover:shadow-purple-500/5 transition-all select-none;
}

.badge-count {
  @apply text-xs px-2 py-0.5 rounded-full border font-mono;
}

.input-select {
  @apply bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-gray-200 focus:outline-none focus:border-purple-500 transition-colors;
}
</style>
