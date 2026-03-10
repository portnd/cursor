<template>
  <div class="w-full">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-indigo-500/15 border border-indigo-500/30 flex items-center justify-center">
          <svg class="w-4 h-4 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/>
          </svg>
        </div>
        <div>
          <h2 class="text-sm font-bold text-white">Global Active Board</h2>
          <p class="text-xs text-gray-500">All pending & in-progress TASK and BUG across projects</p>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <!-- Project legend -->
        <div class="hidden md:flex items-center gap-2 flex-wrap max-w-md">
          <span
            v-for="proj in projectLegend"
            :key="proj.name"
            class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium"
            :style="{ borderColor: proj.color, color: proj.color, backgroundColor: proj.color + '18' }"
          >
            <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :style="{ backgroundColor: proj.color }"/>
            {{ proj.name }}
          </span>
        </div>

        <button
          @click="load"
          :disabled="loading"
          class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50"
        >
          <svg class="h-3.5 w-3.5" :class="loading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-16">
      <svg class="h-7 w-7 animate-spin text-indigo-400 mr-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <span class="text-sm text-gray-500">Loading global board…</span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
      <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
      </svg>
      <div>
        <p class="text-sm font-semibold">Failed to load global board</p>
        <p class="text-xs text-red-300 mt-0.5">{{ error }}</p>
      </div>
    </div>

    <!-- Empty state -->
    <div
      v-else-if="tasks.length === 0"
      class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-700 bg-gray-800/30 py-16 text-center"
    >
      <div class="w-12 h-12 rounded-2xl bg-gray-700/50 flex items-center justify-center mb-3">
        <svg class="h-6 w-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
        </svg>
      </div>
      <p class="text-sm font-semibold text-gray-300">No active tasks across projects</p>
      <p class="text-xs text-gray-500 mt-1">Create TASK or BUG items and assign them to a project</p>
    </div>

    <!-- Kanban Board -->
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4 items-start">
      <div
        v-for="col in columns"
        :key="col.status"
        class="rounded-2xl border bg-gray-800/40 overflow-hidden"
        :class="col.borderClass"
      >
        <!-- Column header -->
        <div class="flex items-center justify-between px-4 py-3 border-b" :class="col.headerBorderClass">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full" :class="col.dotClass"/>
            <span class="text-xs font-bold uppercase tracking-widest" :class="col.labelClass">{{ col.label }}</span>
            <!-- Tooltip badge for READY_FOR_TEST -->
            <span
              v-if="col.status === 'READY_FOR_TEST'"
              class="text-[9px] font-semibold uppercase tracking-wide px-1.5 py-0.5 rounded border border-cyan-500/30 text-cyan-400/70 bg-cyan-500/10"
            >Awaiting PM review</span>
          </div>
          <span
            class="text-xs font-bold tabular-nums px-2 py-0.5 rounded-full"
            :class="col.countClass"
          >{{ tasksByStatus[col.status].length }}</span>
        </div>

        <!-- Cards -->
        <div class="p-3 space-y-3 min-h-[200px]">
          <div
            v-for="task in tasksByStatus[col.status]"
            :key="task.id"
            class="group relative rounded-xl border bg-gray-900/60 p-4 cursor-pointer transition-all shadow-sm hover:shadow-md"
            :class="task.status === 'READY_FOR_TEST'
              ? 'border-cyan-500/50 hover:border-cyan-400/70 hover:bg-gray-800/80 ring-1 ring-cyan-500/20 animate-pulse-border'
              : 'border-gray-700/60 hover:border-gray-600 hover:bg-gray-800/80'"
            @click="goToTask(task)"
          >
            <!-- Project color pill -->
            <div class="flex items-center justify-between mb-2.5">
              <span
                class="inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-[10px] font-semibold leading-none max-w-[160px] truncate"
                :style="{
                  borderColor: task.project_color || '#6366f1',
                  color: task.project_color || '#6366f1',
                  backgroundColor: (task.project_color || '#6366f1') + '18'
                }"
              >
                <span
                  class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                  :style="{ backgroundColor: task.project_color || '#6366f1' }"
                />
                {{ task.project_name || 'Unknown Project' }}
              </span>

              <!-- Priority badge -->
              <span
                class="text-[10px] font-bold uppercase tracking-wide px-1.5 py-0.5 rounded"
                :class="priorityClass(task.priority)"
              >{{ task.priority }}</span>
            </div>

            <!-- Title -->
            <p class="text-sm font-semibold text-white leading-snug line-clamp-2 mb-2">
              {{ task.title }}
            </p>

            <!-- Task type badge + Code badge -->
            <div class="flex items-center gap-2 mb-2.5">
              <span
                class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border"
                :class="task.task_type === 'BUG'
                  ? 'border-red-500/40 text-red-400 bg-red-500/10'
                  : 'border-indigo-500/40 text-indigo-400 bg-indigo-500/10'"
              >{{ task.task_type || 'TASK' }}</span>
              <p v-if="task.code" class="text-[10px] font-mono text-gray-600">{{ task.code }}</p>
            </div>

            <!-- Meta row -->
            <div class="flex items-center justify-between gap-2 text-[10px] text-gray-500">
              <span class="flex items-center gap-1">
                <svg class="w-3 h-3 text-indigo-400/70" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span class="text-indigo-300/80 font-medium">{{ estimatedHours(task.estimated_minutes) }}</span>
              </span>

              <span
                v-if="task.due_at"
                class="font-medium"
                :class="deadlineClass(task.due_at)"
              >{{ formatDeadline(task.due_at) }}</span>
              <span v-else class="text-gray-700">No deadline</span>
            </div>

            <!-- Story points -->
            <div v-if="task.story_points" class="mt-2 flex items-center gap-1">
              <span class="text-[10px] text-gray-600">SP:</span>
              <span class="text-[10px] font-bold text-gray-400">{{ task.story_points }}</span>
            </div>

            <!-- Quick action: Mark Ready for Test (IN_PROGRESS cards only) -->
            <div
              v-if="task.status === 'IN_PROGRESS'"
              class="mt-3 pt-3 border-t border-gray-700/50"
              @click.stop
            >
              <button
                :disabled="markingReadyId === task.id"
                class="w-full flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-lg bg-cyan-500/10 border border-cyan-500/30 text-cyan-400 text-[11px] font-semibold hover:bg-cyan-500/20 hover:border-cyan-400/50 transition-colors disabled:opacity-50"
                @click="handleMarkReadyForTest(task)"
              >
                <svg v-if="markingReadyId === task.id" class="w-3 h-3 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                <svg v-else class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m0 10l2 2 4-4"/>
                </svg>
                {{ markingReadyId === task.id ? 'Submitting…' : 'Mark Ready for Test' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'

type KanbanStatus = 'PENDING' | 'IN_PROGRESS' | 'READY_FOR_TEST'

const columns = [
  {
    status: 'PENDING' as KanbanStatus,
    label: 'Pending',
    borderClass: 'border-gray-700/60',
    headerBorderClass: 'border-gray-700/60',
    dotClass: 'bg-gray-400',
    labelClass: 'text-gray-400',
    countClass: 'bg-gray-700/60 text-gray-300',
  },
  {
    status: 'IN_PROGRESS' as KanbanStatus,
    label: 'In Progress',
    borderClass: 'border-purple-500/30',
    headerBorderClass: 'border-purple-500/20',
    dotClass: 'bg-purple-400 animate-pulse',
    labelClass: 'text-purple-400',
    countClass: 'bg-purple-500/15 text-purple-300',
  },
  {
    status: 'READY_FOR_TEST' as KanbanStatus,
    label: 'Ready for Test',
    borderClass: 'border-cyan-500/30',
    headerBorderClass: 'border-cyan-500/20',
    dotClass: 'bg-cyan-400 animate-pulse',
    labelClass: 'text-cyan-400',
    countClass: 'bg-cyan-500/15 text-cyan-300',
  },
]

const { getGlobalActiveTasks, markReadyForTest } = useTasksApi()

const tasks = ref<GlobalActiveTask[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const markingReadyId = ref<string | null>(null)

const tasksByStatus = computed<Record<KanbanStatus, GlobalActiveTask[]>>(() => {
  const map: Record<KanbanStatus, GlobalActiveTask[]> = { PENDING: [], IN_PROGRESS: [], READY_FOR_TEST: [] }
  for (const t of tasks.value) {
    const s = t.status as KanbanStatus
    if (map[s] !== undefined) map[s].push(t)
  }
  return map
})

const projectLegend = computed(() => {
  const seen = new Map<string, string>()
  for (const t of tasks.value) {
    if (t.project_name && !seen.has(t.project_name)) {
      seen.set(t.project_name, t.project_color || '#6366f1')
    }
  }
  return Array.from(seen.entries()).map(([name, color]) => ({ name, color }))
})

async function load() {
  loading.value = true
  error.value = null
  try {
    tasks.value = await getGlobalActiveTasks()
  } catch (e: any) {
    error.value = e?.data?.message || e?.message || 'Failed to load tasks'
  } finally {
    loading.value = false
  }
}

async function handleMarkReadyForTest(task: GlobalActiveTask) {
  markingReadyId.value = task.id
  try {
    await markReadyForTest(task.id)
    // Optimistically update status in place
    const idx = tasks.value.findIndex(t => t.id === task.id)
    if (idx !== -1) tasks.value[idx] = { ...tasks.value[idx], status: 'READY_FOR_TEST' }
  } catch (e: any) {
    alert(e?.data?.message || e?.message || 'Failed to mark task as ready for test')
  } finally {
    markingReadyId.value = null
  }
}

function goToTask(task: GlobalActiveTask) {
  navigateTo(`/task/${task.code || task.id}?from=dashboard&from_tab=board`)
}

function estimatedHours(minutes: number): string {
  if (!minutes) return '—'
  const h = (minutes / 60).toFixed(1)
  return `${h}h est.`
}

function deadlineClass(dueAt: string): string {
  const diff = new Date(dueAt).getTime() - Date.now()
  if (diff < 0) return 'text-red-400 font-semibold'
  if (diff < 24 * 3600 * 1000) return 'text-amber-400 font-semibold'
  return 'text-gray-500'
}

function formatDeadline(dueAt: string): string {
  const diff = new Date(dueAt).getTime() - Date.now()
  if (diff < 0) {
    const d = Math.floor(Math.abs(diff) / (86400 * 1000))
    return d > 0 ? `Overdue ${d}d` : 'Overdue'
  }
  const d = Math.floor(diff / (86400 * 1000))
  if (d > 0) return `${d}d left`
  const h = Math.floor(diff / 3600000)
  return h > 0 ? `${h}h left` : 'Due soon'
}

function priorityClass(priority: string): string {
  switch (priority) {
    case 'CRITICAL': return 'bg-red-500/20 text-red-400'
    case 'HIGH': return 'bg-orange-500/20 text-orange-400'
    case 'MEDIUM': return 'bg-amber-500/20 text-amber-400'
    case 'LOW': return 'bg-gray-500/20 text-gray-400'
    default: return 'bg-gray-500/20 text-gray-400'
  }
}

onMounted(() => {
  load()
})
</script>
