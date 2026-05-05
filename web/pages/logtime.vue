<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
    <!-- Header -->
    <header class="sticky top-0 z-10 border-b border-gray-200 dark:border-gray-800 bg-white/95 dark:bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-5xl px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between gap-4 py-5">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-purple-100 to-indigo-100 dark:from-purple-600 dark:to-indigo-600 shadow-sm dark:shadow-lg text-xl">
              🕐
            </div>
            <div>
              <h1 class="text-xl font-bold tracking-tight text-gray-900 dark:text-white">Work Log</h1>
              <p class="text-xs text-gray-600 dark:text-gray-400 mt-0.5">บันทึกเวลาทำงาน · ดู · แก้ไข · ลบ</p>
            </div>
          </div>

          <!-- Action buttons -->
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="showBulkLog = true"
              class="flex items-center gap-2 px-4 py-2 rounded-xl bg-indigo-600 dark:bg-indigo-600 hover:bg-indigo-700 dark:hover:bg-indigo-500 text-white text-sm font-semibold transition-colors"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/>
              </svg>
              EOD Batch
              <kbd class="text-[9px] font-mono bg-indigo-700/60 border border-indigo-500/40 rounded px-1 py-0.5">⌘⇧L</kbd>
            </button>
            <button
              type="button"
              @click="showQuickLog = true"
              class="flex items-center gap-2 px-4 py-2 rounded-xl bg-purple-600 dark:bg-purple-600 hover:bg-purple-700 dark:hover:bg-purple-500 text-white text-sm font-semibold transition-colors"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              Quick Log
              <kbd class="text-[9px] font-mono bg-purple-700/60 border border-purple-500/40 rounded px-1 py-0.5">⌘L</kbd>
            </button>
          </div>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8 space-y-6">

      <!-- Timer + Today Summary row -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">

        <!-- Timer Card -->
        <div class="bg-white dark:bg-gray-800/60 border border-gray-200 dark:border-gray-700/60 rounded-2xl p-5">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-base">⏱</span>
            <h2 class="text-sm font-semibold text-gray-900 dark:text-white">Timer</h2>
          </div>

          <!-- Running -->
          <div v-if="isRunning" class="space-y-3">
            <div class="flex items-center gap-3">
              <span class="relative flex h-3 w-3 shrink-0">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-purple-400 opacity-75" />
                <span class="relative inline-flex rounded-full h-3 w-3 bg-purple-500" />
              </span>
              <span class="font-mono text-2xl font-bold text-purple-300 tabular-nums">{{ elapsedDisplay }}</span>
            </div>
            <div class="text-xs text-gray-400 truncate">
              <span class="font-mono text-purple-400">{{ timerState?.taskCode }}</span>
              {{ timerState?.taskTitle ? ' · ' + timerState.taskTitle : '' }}
            </div>
            <button
              type="button"
              @click="handleStop"
              class="w-full py-2 rounded-xl bg-red-900/30 border border-red-700/30 text-red-400 hover:bg-red-800/40 text-sm font-semibold transition-colors"
            >
              ■ Stop & Log
            </button>
          </div>

          <!-- Idle -->
          <div v-else class="space-y-3">
            <p class="text-sm text-gray-600 dark:text-gray-500">ไม่มี timer กำลังทำงาน</p>
            <div class="relative">
              <input
                v-model="timerSearch"
                type="text"
                placeholder="🔍 ค้นหา task เพื่อเริ่ม timer..."
                class="w-full bg-white dark:bg-gray-700/60 border border-gray-300 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm text-gray-900 dark:text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
                @focus="loadTimerTasks"
              />
              <div v-if="timerTasks.length" class="mt-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden max-h-48 overflow-y-auto">
                <div v-if="timerLoading" class="py-3 text-center text-xs text-gray-500">Loading...</div>
                <template v-else-if="filteredTimerTasks.length">
                  <button
                    v-for="task in filteredTimerTasks"
                    :key="task.id"
                    type="button"
                    @click="startTimerFor(task)"
                    class="w-full flex items-center gap-2.5 px-3 py-2.5 hover:bg-gray-100 dark:hover:bg-gray-700/60 text-left border-b border-gray-200 dark:border-gray-700/30 last:border-0 transition-colors"
                    :class="{ 'pl-7': task.parent_id }"
                  >
                    <span class="text-sm shrink-0">{{ task.parent_id ? '↳' : taskIcon(task.task_type) }}</span>
                    <span class="font-mono text-[10px] text-purple-400 shrink-0">{{ task.code }}</span>
                    <span class="text-xs text-gray-700 dark:text-gray-200 truncate flex-1">{{ task.title }}</span>
                    <span v-if="task.parent_task_code" class="font-mono text-[9px] text-indigo-400/60 shrink-0 hidden sm:block">↑{{ task.parent_task_code }}</span>
                    <span v-if="task.assigned_to_display_name || task.assigned_to_email" class="text-[10px] text-indigo-400/70 shrink-0 hidden sm:block">
                      {{ task.assigned_to_display_name || task.assigned_to_email }}
                    </span>
                  </button>
                  <div
                    v-if="!timerSearch && timerTasks.length > TIMER_TASK_LIMIT"
                    class="px-3 py-2 text-[11px] text-gray-500 border-t border-gray-700/40 bg-gray-800/60"
                  >
                    Showing first {{ TIMER_TASK_LIMIT }} of {{ timerTasks.length }} tasks. Use search to narrow down.
                  </div>
                </template>
                <div v-else class="py-3 text-center text-xs text-gray-500">
                  {{ timerSearch ? 'ไม่พบ task' : 'ไม่มี active tasks' }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Today Summary Card -->
        <div class="bg-white dark:bg-gray-800/60 border border-gray-200 dark:border-gray-700/60 rounded-2xl p-5">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span class="text-base">📊</span>
              <h2 class="text-sm font-semibold text-white">Today's Summary</h2>
            </div>
            <span class="text-xs text-gray-500">{{ today }}</span>
          </div>

          <div class="mb-3">
            <div class="flex items-end justify-between mb-1.5">
              <span class="text-3xl font-bold tabular-nums" :class="totalMinutes >= 480 ? 'text-green-400' : 'text-white'">
                {{ (totalMinutes / 60).toFixed(1) }}<span class="text-base font-normal text-gray-400 ml-1">h</span>
              </span>
              <span class="text-xs text-gray-500">/ 8h</span>
            </div>
            <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div class="h-full rounded-full transition-all duration-500"
                :class="totalMinutes >= 480 ? 'bg-green-500' : totalMinutes >= 240 ? 'bg-blue-500' : 'bg-blue-600'"
                :style="{ width: Math.min((totalMinutes / 480) * 100, 100) + '%' }" />
            </div>
          </div>

          <!-- Work type breakdown -->
          <div v-if="workTypeBreakdown.length" class="flex flex-wrap gap-1.5 mt-3">
            <div v-for="wt in workTypeBreakdown" :key="wt.type" class="flex items-center gap-1 text-[10px]">
              <span class="px-1.5 py-0.5 rounded font-semibold uppercase" :class="workTypeBadge(wt.type)">{{ wt.type }}</span>
              <span class="text-gray-500">{{ wt.minutes }}m</span>
            </div>
          </div>
          <p v-else class="text-xs text-gray-600 mt-2">ยังไม่มี log วันนี้</p>
        </div>
      </div>

      <!-- My Work Logs History -->
      <section class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700/60 rounded-2xl overflow-hidden">
        <!-- Section header -->
        <div class="flex items-center justify-between gap-3 px-5 py-4 border-b border-gray-200 dark:border-gray-700/60">
          <div class="flex items-center gap-2.5">
            <span class="text-base">📋</span>
            <div>
              <h2 class="text-sm font-semibold text-gray-900 dark:text-white">Log History</h2>
              <p class="text-xs text-gray-500 mt-0.5">
                {{ logsDate === today ? 'วันนี้' : logsDate }}
                <span v-if="!logsLoading" class="ml-1">
                  · <span :class="totalMinutes >= 480 ? 'text-green-400' : 'text-blue-400'">{{ formatMinutes(totalMinutes) }}</span>
                  <span class="text-gray-600"> / 8h</span>
                </span>
              </p>
            </div>
          </div>
          <!-- date nav -->
          <div class="flex items-center gap-1.5">
            <button type="button" @click="shiftDate(-1)"
              class="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <UiDatePicker v-model="logsDate" placeholder="Select date…" />
            <button type="button" @click="shiftDate(1)" :disabled="logsDate >= today"
              class="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors disabled:opacity-30 disabled:cursor-not-allowed">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>
            <button v-if="logsDate !== today" type="button" @click="logsDate = today"
              class="text-xs px-2.5 py-1 rounded-lg bg-indigo-700/30 text-indigo-300 hover:bg-indigo-700/50 transition-colors">
              Today
            </button>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="logsLoading" class="flex items-center justify-center py-10 gap-2 text-sm text-gray-500">
          <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
          </svg>
          กำลังโหลด...
        </div>

        <!-- Empty -->
        <div v-else-if="!logEntries.length" class="flex flex-col items-center justify-center py-12 gap-2">
          <span class="text-4xl">📭</span>
          <p class="text-sm text-gray-600">ไม่มี log สำหรับวันนี้</p>
          <button type="button" @click="showQuickLog = true"
            class="mt-2 text-xs px-3 py-1.5 rounded-lg bg-purple-700/30 text-purple-300 hover:bg-purple-700/50 transition-colors">
            + Quick Log
          </button>
        </div>

        <!-- Log list -->
        <div v-else class="divide-y divide-gray-200 dark:divide-gray-700/30">
          <div
            v-for="log in logEntries"
            :key="log.id"
            class="group flex items-center gap-4 px-5 py-3.5 hover:bg-gray-100 dark:hover:bg-gray-700/20 transition-colors"
          >
            <span class="text-lg shrink-0">{{ workTypeEmoji(log.work_type) }}</span>

            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{{ formatMinutes(log.minutes) }}</span>
                <span class="text-[10px] px-1.5 py-0.5 rounded font-semibold uppercase tracking-wide shrink-0" :class="workTypeBadge(log.work_type)">
                  {{ log.work_type || 'DEV' }}
                </span>
                <span v-if="log.is_timer_session" class="text-[10px] px-1.5 py-0.5 rounded bg-violet-900/50 text-violet-300 shrink-0">⏱ Timer</span>
              </div>
              <div class="flex items-center gap-1.5 mt-0.5 min-w-0">
                <span v-if="log.task_code" class="font-mono text-[10px] text-purple-400 shrink-0">{{ log.task_code }}</span>
                <span v-if="log.task_title" class="text-xs text-gray-400 truncate">{{ log.task_title }}</span>
              </div>
              <p v-if="log.description" class="text-xs text-gray-500 mt-0.5 truncate">{{ log.description }}</p>
            </div>

            <div class="shrink-0 text-right">
              <p class="text-[10px] text-gray-600">
                {{ log.logged_at ? new Date(log.logged_at).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) : '' }}
              </p>
            </div>

            <!-- Edit / Delete (hover) -->
            <div class="shrink-0 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <button type="button" @click="startEditLog(log)"
                class="p-1.5 rounded-lg text-gray-500 hover:text-purple-400 hover:bg-purple-900/20 transition-colors" title="แก้ไข">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                </svg>
              </button>
              <button type="button" @click="confirmDeleteLog(log)"
                class="p-1.5 rounded-lg text-gray-500 hover:text-red-400 hover:bg-red-900/20 transition-colors" title="ลบ">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </section>

    </main>

    <!-- Modals -->
    <TasksQuickLogTimeModal v-model="showQuickLog" @logged="loadLogs" />
    <TasksBulkEodLoggerModal :show="showBulkLog" @close="showBulkLog = false" @done="() => { showBulkLog = false; loadLogs() }" />
    <TasksQuickLogTimeModal
      v-model="showTimerLog"
      :preselected-task-id="stoppedTaskId"
      :prefilled-minutes="stoppedMinutes"
      @logged="onTimerLogged"
    />

    <!-- Edit Log Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="editingLog" class="fixed inset-0 z-50 flex items-center justify-center p-4" @keydown.escape="editingLog = null">
          <div class="fixed inset-0 bg-black/70 backdrop-blur-sm" @click="editingLog = null" />
          <div class="relative w-full max-w-md bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl p-5" @click.stop>
            <h3 class="text-sm font-semibold text-white mb-4">✏️ แก้ไข Time Log</h3>

            <!-- Task selector -->
            <div class="mb-3">
              <label class="text-xs text-gray-400 mb-1.5 block">Task</label>
              <!-- Current task chip -->
              <div class="flex items-center gap-2 p-2 bg-gray-800 border border-gray-700 rounded-xl mb-2">
                <span class="font-mono text-[10px] text-purple-400 shrink-0">{{ editTaskCode || '—' }}</span>
                <span class="text-xs text-gray-300 truncate flex-1">{{ editTaskTitle || 'ไม่มี task' }}</span>
                <button type="button" @click="clearEditTask" class="text-[10px] text-gray-500 hover:text-gray-300 shrink-0 transition-colors">รีเซ็ต</button>
              </div>
              <!-- Search -->
              <div class="relative">
                <input
                  v-model="editTaskSearch"
                  type="text"
                  placeholder="ค้นหา task เพื่อเปลี่ยน..."
                  class="w-full bg-gray-700/60 border border-gray-600 rounded-xl px-3 py-2 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
                  @focus="editTaskDropdown = true"
                  @blur="hideEditTaskDropdown()"
                />
                <div v-if="editTaskDropdown && (editTaskSearch || timerTasks.length)" class="absolute z-10 mt-1 w-full bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-xl max-h-44 overflow-y-auto">
                  <template v-if="filteredEditTasks.length">
                    <button
                      v-for="task in filteredEditTasks" :key="task.id" type="button"
                      @mousedown.prevent="selectEditTask(task)"
                      class="w-full flex items-center gap-2 px-3 py-2 hover:bg-gray-700/60 text-left border-b border-gray-700/30 last:border-0 transition-colors"
                      :class="{ 'pl-6': task.parent_id }"
                    >
                      <span class="text-xs shrink-0">{{ task.parent_id ? '↳' : taskIcon(task.task_type) }}</span>
                      <span class="font-mono text-[10px] text-purple-400 shrink-0">{{ task.code }}</span>
                      <span class="text-xs text-gray-200 truncate flex-1">{{ task.title }}</span>
                      <span v-if="task.parent_task_code" class="font-mono text-[9px] text-indigo-400/60 shrink-0">↑{{ task.parent_task_code }}</span>
                      <span v-if="task.project_name" class="text-[10px] text-gray-500 shrink-0 hidden sm:block">{{ task.project_name }}</span>
                    </button>
                  </template>
                  <div v-else class="py-3 text-center text-xs text-gray-500">ไม่พบ task</div>
                </div>
              </div>
            </div>

            <div class="mb-3">
              <label class="text-xs text-gray-400 mb-1 block">เวลา (นาที)</label>
              <div class="flex gap-1.5 mb-2">
                <button v-for="p in [15, 30, 60, 120]" :key="p" type="button" @click="editMinutes = p"
                  class="text-xs px-2.5 py-1 rounded-lg border transition-colors"
                  :class="editMinutes === p ? 'bg-indigo-600 border-indigo-500 text-white' : 'bg-gray-700 border-gray-600 text-gray-400 hover:text-white'">
                  +{{ p >= 60 ? (p / 60) + 'h' : p + 'm' }}
                </button>
              </div>
              <input v-model.number="editMinutes" type="number" min="1" max="960"
                class="w-full bg-gray-700/60 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div class="mb-3">
              <label class="text-xs text-gray-400 mb-1.5 block">ประเภทงาน</label>
              <div class="flex flex-wrap gap-1.5">
                <button v-for="wt in WORK_TYPES" :key="wt.value" type="button" @click="editWorkType = wt.value"
                  class="text-xs px-2.5 py-1 rounded-lg border transition-colors"
                  :class="editWorkType === wt.value ? workTypeActiveClass(wt.value) : 'bg-gray-700 border-gray-600 text-gray-400 hover:text-white'">
                  {{ wt.emoji }} {{ wt.label }}
                </button>
              </div>
            </div>
            <div class="mb-4">
              <label class="text-xs text-gray-400 mb-1 block">คำอธิบาย</label>
              <input v-model="editDescription" type="text" placeholder="Optional..."
                class="w-full bg-gray-700/60 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div class="flex gap-2">
              <button type="button" @click="editingLog = null"
                class="flex-1 py-2 text-sm text-gray-400 border border-gray-700 rounded-xl hover:bg-gray-800 transition-colors">
                ยกเลิก
              </button>
              <button type="button" @click="saveEdit" :disabled="editSaving || editMinutes < 1"
                class="flex-1 py-2 text-sm font-semibold bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl transition-colors disabled:opacity-50">
                {{ editSaving ? 'กำลังบันทึก...' : 'บันทึก' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useTasksApi, WORK_TYPES } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { TimeLog, GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTimer } from '~/composables/useTimer'
import { localDateStr } from '~/composables/useLocalDate'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const { getMyDailyTimeLogs, deleteTimeLog, editTimeLog, getTeamActiveTasks, getKomgripTasks } = useTasksApi()
const { timerState, elapsedDisplay, isRunning, start, stop } = useTimer()

const today = ref(localDateStr())

const showQuickLog = ref(false)
const showBulkLog = ref(false)

// ── Timer ────────────────────────────────────────────────────────
const showTimerLog = ref(false)
const stoppedTaskId = ref<string | undefined>()
const stoppedMinutes = ref<number | undefined>()
const timerSearch = ref('')
const timerLoading = ref(false)
const timerTasks = ref<GlobalActiveTask[]>([])

function handleStop() {
  const result = stop()
  if (!result) return
  stoppedTaskId.value = result.taskId
  stoppedMinutes.value = result.minutes
  showTimerLog.value = true
}

function onTimerLogged() {
  stoppedTaskId.value = undefined
  stoppedMinutes.value = undefined
  loadLogs()
}

async function loadTimerTasks() {
  if (timerTasks.value.length) return
  timerLoading.value = true
  try {
    const [active, komgrip] = await Promise.all([
      getTeamActiveTasks().catch(() => [] as GlobalActiveTask[]),
      getKomgripTasks().catch(() => [] as GlobalActiveTask[]),
    ])
    const activeKomgrip = komgrip.filter((t: any) => t?.status !== 'COMPLETED')
    const komgripAsGlobal = activeKomgrip.map((t: any) => ({
      ...t,
      project_name: 'Komgrip',
      project_color: '#8b5cf6',
    })) as GlobalActiveTask[]
    timerTasks.value = [...active, ...komgripAsGlobal]
  } catch {
    timerTasks.value = []
  } finally {
    timerLoading.value = false
  }
}

const TIMER_TASK_LIMIT = 100

const filteredTimerTasks = computed(() => {
  const q = timerSearch.value.trim().toLowerCase()
  const list = q
    ? timerTasks.value.filter(t =>
        t.code?.toLowerCase().includes(q) ||
        t.title?.toLowerCase().includes(q) ||
        t.assigned_to_display_name?.toLowerCase().includes(q) ||
        t.assigned_to_email?.toLowerCase().includes(q) ||
        t.project_name?.toLowerCase().includes(q) ||
        t.parent_task_code?.toLowerCase().includes(q) ||
        t.parent_task_title?.toLowerCase().includes(q),
      )
    : timerTasks.value
  return list.slice(0, TIMER_TASK_LIMIT)
})

function startTimerFor(task: GlobalActiveTask) {
  start(task.id, task.title, task.code)
  timerSearch.value = ''
  timerTasks.value = []
}

function taskIcon(type: string) {
  return ({ FEATURE: '⚡', TASK: '✅', BUG: '🐛' } as Record<string, string>)[type] ?? '📋'
}

// ── Work Log History ─────────────────────────────────────────────
const logsDate = ref(today.value)
const logsLoading = ref(false)
const logSummary = ref<{ date: string; total_minutes: number; entries: TimeLog[] } | null>(null)

// Sync logsDate to real local date after hydration
watch(today, (newToday) => {
  if (logsDate.value !== newToday) logsDate.value = newToday
})

const logEntries = computed(() => logSummary.value?.entries ?? [])
const totalMinutes = computed(() => logSummary.value?.total_minutes ?? 0)

const workTypeBreakdown = computed(() => {
  const map = new Map<string, number>()
  for (const log of logEntries.value) {
    const wt = log.work_type || 'DEV'
    map.set(wt, (map.get(wt) ?? 0) + log.minutes)
  }
  return Array.from(map.entries()).map(([type, minutes]) => ({ type, minutes })).sort((a, b) => b.minutes - a.minutes)
})

async function loadLogs() {
  logsLoading.value = true
  try { logSummary.value = await getMyDailyTimeLogs(logsDate.value) }
  catch { /* non-critical */ }
  finally { logsLoading.value = false }
}

function shiftDate(delta: number) {
  const d = new Date(logsDate.value)
  d.setDate(d.getDate() + delta)
  logsDate.value = localDateStr(d)
}

watch(logsDate, loadLogs, { immediate: true })

// ── Delete ───────────────────────────────────────────────────────
async function confirmDeleteLog(log: TimeLog) {
  if (!confirm(`ยืนยันลบ log ${formatMinutes(log.minutes)}?`)) return
  try {
    await deleteTimeLog(log.id)
    await loadLogs()
  } catch {
    alert('ลบไม่สำเร็จ')
  }
}

// ── Edit ─────────────────────────────────────────────────────────
const editingLog = ref<TimeLog | null>(null)
const editMinutes = ref(0)
const editWorkType = ref('DEV')
const editDescription = ref('')
const editSaving = ref(false)

// task search inside edit modal
const editTaskSearch = ref('')
const editTaskId = ref<string | undefined>()
const editTaskCode = ref('')
const editTaskTitle = ref('')
const editTaskDropdown = ref(false)

const filteredEditTasks = computed(() => {
  const q = editTaskSearch.value.trim().toLowerCase()
  const list = q
    ? timerTasks.value.filter(t =>
        t.code?.toLowerCase().includes(q) ||
        t.title?.toLowerCase().includes(q) ||
        t.project_name?.toLowerCase().includes(q) ||
        t.parent_task_code?.toLowerCase().includes(q) ||
        t.parent_task_title?.toLowerCase().includes(q),
      )
    : timerTasks.value
  return list.slice(0, 50)
})

function selectEditTask(task: GlobalActiveTask) {
  editTaskId.value = task.id
  editTaskCode.value = task.code ?? ''
  editTaskTitle.value = task.title
  editTaskSearch.value = ''
  editTaskDropdown.value = false
}

function hideEditTaskDropdown() {
  setTimeout(() => { editTaskDropdown.value = false }, 150)
}

function clearEditTask() {
  editTaskId.value = editingLog.value?.task_id
  editTaskCode.value = editingLog.value?.task_code ?? ''
  editTaskTitle.value = editingLog.value?.task_title ?? ''
  editTaskSearch.value = ''
  editTaskDropdown.value = false
}

async function startEditLog(log: TimeLog) {
  editingLog.value = log
  editMinutes.value = log.minutes
  editWorkType.value = log.work_type || 'DEV'
  editDescription.value = log.description || ''
  editTaskId.value = log.task_id
  editTaskCode.value = log.task_code ?? ''
  editTaskTitle.value = log.task_title ?? ''
  editTaskSearch.value = ''
  editTaskDropdown.value = false
  await loadTimerTasks()
}

async function saveEdit() {
  if (!editingLog.value) return
  editSaving.value = true
  try {
    await editTimeLog(editingLog.value.id, editMinutes.value, editDescription.value, editWorkType.value, editTaskId.value)
    editingLog.value = null
    await loadLogs()
  } catch {
    alert('บันทึกไม่สำเร็จ')
  } finally {
    editSaving.value = false
  }
}

// ── Helpers ──────────────────────────────────────────────────────
function formatMinutes(mins: number) {
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function workTypeEmoji(wt: string) {
  return WORK_TYPES.find(w => w.value === wt)?.emoji ?? '📋'
}

function workTypeBadge(wt: string): string {
  const map: Record<string, string> = {
    DEV: 'bg-blue-900/50 text-blue-300',
    REVIEW: 'bg-cyan-900/50 text-cyan-300',
    TESTING: 'bg-green-900/50 text-green-300',
    MEETING: 'bg-orange-900/50 text-orange-300',
    RESEARCH: 'bg-purple-900/50 text-purple-300',
    OTHER: 'bg-gray-700/80 text-gray-400',
  }
  return map[wt] || 'bg-gray-700/80 text-gray-400'
}

function workTypeActiveClass(wt: string): string {
  const map: Record<string, string> = {
    DEV: 'bg-blue-600 border-blue-500 text-white',
    REVIEW: 'bg-cyan-600 border-cyan-500 text-white',
    TESTING: 'bg-green-600 border-green-500 text-white',
    MEETING: 'bg-orange-600 border-orange-500 text-white',
    RESEARCH: 'bg-purple-600 border-purple-500 text-white',
    OTHER: 'bg-gray-600 border-gray-500 text-white',
  }
  return map[wt] || 'bg-indigo-600 border-indigo-500 text-white'
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
