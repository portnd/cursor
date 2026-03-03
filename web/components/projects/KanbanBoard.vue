<template>
  <div class="kanban-board w-full min-w-0">
    <!-- Filter Bar (responsive: stack on narrow, wrap on wide) -->
    <div class="flex flex-col sm:flex-row sm:flex-wrap items-stretch sm:items-center gap-3 mb-4 sm:mb-5 p-3 bg-gray-800/50 rounded-xl border border-gray-700/50">
      <div class="flex flex-wrap items-center gap-2 sm:gap-2">
        <label class="text-xs text-gray-400 uppercase tracking-wide font-medium w-14 sm:w-auto shrink-0">Sprint</label>
        <select v-model="filterSprint" class="input-select text-sm min-w-0 flex-1 sm:min-w-[160px] sm:flex-none max-w-full">
          <option value="">All Sprints</option>
          <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
        <span v-if="sprints.length === 0" class="text-xs text-gray-500 w-full sm:w-auto">— No sprints yet</span>
        <span v-else-if="filterSprint" class="text-xs text-indigo-400 truncate max-w-full">{{ sprintNameById(filterSprint) || 'Sprint' }}</span>
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
    <div class="flex lg:grid lg:grid-cols-5 gap-3 overflow-x-auto lg:overflow-visible pb-2 -mx-1 px-1 lg:mx-0 lg:px-0 snap-x snap-mandatory scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-transparent">
      <div
        v-for="col in columns"
        :key="col.status"
        class="kanban-col min-w-[200px] sm:min-w-[220px] w-[200px] sm:w-[220px] lg:w-auto shrink-0 lg:shrink snap-start"
        @dragover.prevent="onDragOver($event, col.status)"
        @dragleave="onDragLeave"
        @drop="onDrop($event, col.status)"
        :class="{ 'drop-target': dragOverCol === col.status }"
      >
        <!-- Column Header -->
        <div class="flex items-center justify-between mb-3 px-1">
          <div class="flex items-center gap-2">
            <span class="text-lg">{{ col.icon }}</span>
            <span class="text-sm font-semibold" :class="col.headerClass">{{ col.label }}</span>
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
                <span class="text-xs text-gray-500 font-mono">{{ task.code }}</span>
                <span class="text-xs border rounded px-1" :class="priorityCls(task.priority)">{{ task.priority }}</span>
              </div>
              <p v-if="task.sprint_id && sprintNameById(task.sprint_id)" class="text-[10px] text-indigo-400 mb-1 truncate" :title="sprintNameById(task.sprint_id)">📌 {{ sprintNameById(task.sprint_id) }}</p>
              <p class="text-sm text-gray-200 font-medium leading-snug mb-2 line-clamp-2">{{ task.title }}</p>
              <div class="flex items-center justify-between text-xs text-gray-500">
                <span v-if="task.story_points" class="flex items-center gap-1"><span class="text-purple-400">◆</span> {{ task.story_points }} SP</span>
                <span v-if="task.assigned_to" class="flex items-center gap-1"><span class="w-4 h-4 rounded-full bg-indigo-600 flex items-center justify-center text-white text-[10px] font-bold">{{ task.assigned_to }}</span></span>
                <span v-if="daysUntilDue(task) !== null" :class="daysUntilDue(task)! < 0 ? 'text-red-400' : daysUntilDue(task)! <= 2 ? 'text-yellow-400' : 'text-gray-500'">
                  {{ daysUntilDue(task)! < 0 ? Math.abs(daysUntilDue(task)!) + 'd overdue' : daysUntilDue(task) + 'd left' }}
                </span>
              </div>
              <div v-if="task.progress > 0" class="mt-2">
                <div class="h-1 bg-gray-700 rounded-full overflow-hidden">
                  <div class="h-full bg-indigo-500 rounded-full transition-all" :style="{ width: task.progress + '%' }"></div>
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
                  <span class="text-xs text-gray-500 font-mono">{{ task.code }}</span>
                  <span class="text-xs border rounded px-1" :class="priorityCls(task.priority)">{{ task.priority }}</span>
                </div>
                <p v-if="task.sprint_id && sprintNameById(task.sprint_id)" class="text-[10px] text-indigo-400 mb-1 truncate">📌 {{ sprintNameById(task.sprint_id) }}</p>
                <p class="text-sm text-gray-200 font-medium leading-snug mb-2 line-clamp-2">{{ task.title }}</p>
                <div class="flex items-center justify-between text-xs text-gray-500">
                  <span v-if="task.story_points" class="flex items-center gap-1"><span class="text-purple-400">◆</span> {{ task.story_points }} SP</span>
                  <span v-if="task.assigned_to" class="flex items-center gap-1"><span class="w-4 h-4 rounded-full bg-indigo-600 flex items-center justify-center text-white text-[10px] font-bold">{{ task.assigned_to }}</span></span>
                  <span v-if="daysUntilDue(task) !== null" :class="daysUntilDue(task)! < 0 ? 'text-red-400' : daysUntilDue(task)! <= 2 ? 'text-yellow-400' : 'text-gray-500'">
                    {{ daysUntilDue(task)! < 0 ? Math.abs(daysUntilDue(task)!) + 'd overdue' : daysUntilDue(task) + 'd left' }}
                  </span>
                </div>
                <div v-if="task.progress > 0" class="mt-2">
                  <div class="h-1 bg-gray-700 rounded-full overflow-hidden">
                    <div class="h-full bg-indigo-500 rounded-full transition-all" :style="{ width: task.progress + '%' }"></div>
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
            Drop tasks here
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Task } from '~/core/modules/projects/infrastructure/projects-api'
import type { Sprint } from '~/core/modules/projects/infrastructure/projects-api'

const props = defineProps<{
  tasks: Task[]
  sprints: Sprint[]
}>()

const emit = defineEmits<{
  (e: 'task-click', task: Task): void
  (e: 'status-change', taskId: string, status: string): void
}>()

const filterSprint = ref('')
const filterPriority = ref('')
const swimLane = ref<'none' | 'priority' | 'assignee'>('none')
const dragTask = ref<Task | null>(null)
const dragOverCol = ref('')

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
function daysUntilDue(task: Task): number | null {
  if (!task.due_at) return null
  const diff = new Date(task.due_at).getTime() - Date.now()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

const columns = [
  { status: 'PENDING', label: 'Backlog', icon: '📋', headerClass: 'text-gray-300', wipLimit: 0 },
  { status: 'IN_PROGRESS', label: 'In Progress', icon: '⚡', headerClass: 'text-blue-400', wipLimit: 5 },
  { status: 'REVIEW_PENDING', label: 'In Review', icon: '🔍', headerClass: 'text-yellow-400', wipLimit: 3 },
  { status: 'COMPLETED', label: 'Done', icon: '✅', headerClass: 'text-green-400', wipLimit: 0 },
  { status: 'BLOCKED', label: 'Blocked', icon: '🚫', headerClass: 'text-red-400', wipLimit: 0 },
]

const filteredTasks = computed(() =>
  props.tasks.filter((t) => {
    if (filterSprint.value && t.sprint_id !== filterSprint.value) return false
    if (filterPriority.value && t.priority !== filterPriority.value) return false
    return true
  })
)

const COLUMN_STATUSES = ['PENDING', 'IN_PROGRESS', 'REVIEW_PENDING', 'COMPLETED', 'BLOCKED'] as const

const tasksByCol = computed(() => {
  const map: Record<string, Task[]> = {}
  for (const col of columns) map[col.status] = []
  for (const t of filteredTasks.value) {
    const status = t.status && COLUMN_STATUSES.includes(t.status as any) ? t.status : 'PENDING'
    if (map[status]) map[status].push(t)
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
    const assignees = [...new Set(tasks.map((t) => t.assigned_to ?? 'Unassigned'))]
    return assignees.map((a) => ({
      key: String(a),
      label: a === 'Unassigned' ? 'Unassigned' : `Dev #${a}`,
      tasks: tasks.filter((t) => (t.assigned_to ?? 'Unassigned') === a),
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

function onDrop(e: DragEvent, newStatus: string) {
  dragOverCol.value = ''
  const task = dragTask.value
  if (!task || task.status === newStatus) return
  emit('status-change', task.id, newStatus)
  dragTask.value = null
}
</script>

<style scoped>
.kanban-col {
  @apply bg-gray-800/50 rounded-xl p-3 flex flex-col border border-gray-700/50 transition-all;
}

.drop-target {
  @apply border-indigo-500/60 bg-indigo-500/5;
}

.kanban-card {
  @apply bg-gray-900 border border-gray-700 rounded-lg p-3 cursor-grab active:cursor-grabbing hover:border-indigo-500/50 hover:shadow-lg hover:shadow-indigo-500/5 transition-all select-none;
}

.badge-count {
  @apply text-xs px-2 py-0.5 rounded-full border font-mono;
}

.input-select {
  @apply bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-gray-200 focus:outline-none focus:border-indigo-500 transition-colors;
}
</style>
