<template>
  <!-- ── RUNNING STATE ── -->
  <div v-if="isRunning" class="timer-running mx-1">
    <div class="flex items-center gap-2">
      <!-- Live pulse dot -->
      <span class="relative flex h-2 w-2 shrink-0">
        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-purple-400 opacity-75" />
        <span class="relative inline-flex rounded-full h-2 w-2 bg-purple-500" />
      </span>

      <div v-show="!sidebarCollapsed" class="flex-1 min-w-0">
        <div class="font-mono text-xs font-bold text-purple-300 leading-none tabular-nums">{{ elapsedDisplay }}</div>
        <div class="text-[10px] text-gray-500 truncate mt-0.5 max-w-[130px]">
          <span class="font-mono text-purple-500">{{ timerState?.taskCode }}</span>
          {{ ' ' + timerState?.taskTitle }}
        </div>
      </div>

      <!-- Collapsed: just show time -->
      <div v-show="sidebarCollapsed" class="font-mono text-[10px] text-purple-300 tabular-nums">
        {{ elapsedDisplay }}
      </div>

      <!-- Stop button -->
      <button
        type="button"
        @click="handleStop"
        class="shrink-0 px-2 py-0.5 bg-red-900/30 border border-red-300 dark:border-red-700/30 text-red-400 hover:bg-red-800/40 rounded text-[10px] font-bold tracking-wide transition-colors"
        title="Stop & log time"
      >
        ■
      </button>
    </div>
  </div>

  <!-- ── IDLE STATE ── -->
  <div v-else class="relative">
      <button
        type="button"
        @click="togglePicker"
        :title="sidebarCollapsed ? 'Start Timer' : undefined"
        class="nav-link-timer w-full"
        :class="showPicker
          ? 'border-purple-600/50 bg-purple-900/20 text-purple-300'
          : 'text-purple-400 hover:text-purple-300 hover:bg-purple-700/20 hover:border-purple-600/30'"
      >
        <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span v-show="!sidebarCollapsed" class="font-medium truncate">Start Timer</span>
      </button>

      <!-- Task picker popover -->
      <Transition name="picker-drop">
        <div
          v-if="showPicker && !sidebarCollapsed"
          class="picker-panel"
          @mousedown.prevent
          @click.stop
        >
          <div class="px-3 pt-2.5 pb-1.5 border-b border-gray-700/60">
            <p class="text-[10px] font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Pick a task to time</p>
            <input
              v-model="pickerSearch"
              type="text"
              placeholder="Search..."
              class="w-full px-2.5 py-1.5 bg-gray-700/60 border border-gray-600/50 rounded-lg text-xs text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
              ref="pickerInput"
            />
          </div>

          <div class="overflow-y-auto max-h-52">
            <div v-if="loadingTasks" class="py-4 text-center text-xs text-gray-500">Loading...</div>
            <template v-else-if="filteredPickerTasks.length">
              <button
                v-for="task in filteredPickerTasks"
                :key="task.id"
                type="button"
                @click="startTimerFor(task)"
                class="w-full flex items-start gap-2 px-3 py-2.5 hover:bg-gray-700/50 text-left border-b border-gray-700/30 last:border-0 transition-colors"
              >
                <span class="text-sm mt-0.5 shrink-0">{{ taskIcon(task.task_type) }}</span>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1.5">
                    <span class="font-mono text-[10px] text-purple-400 shrink-0">{{ task.code }}</span>
                    <span v-if="task.assigned_to_display_name || task.assigned_to_email" class="text-[10px] text-indigo-400/80 truncate">
                      · {{ task.assigned_to_display_name || task.assigned_to_email }}
                    </span>
                  </div>
                  <span class="text-xs text-gray-200 truncate block">{{ task.title }}</span>
                  <span v-if="task.project_name" class="text-[10px] text-gray-600">{{ task.project_name }}</span>
                </div>
              </button>
            </template>
            <div v-else class="py-4 text-center text-xs text-gray-500">
              {{ pickerSearch ? 'No matches' : 'No active tasks' }}
            </div>
          </div>
        </div>
      </Transition>
    </div>

  <!-- Log Time modal after stop -->
  <TasksQuickLogTimeModal
    v-model="showLog"
    :preselected-task-id="stoppedTaskId"
    :prefilled-minutes="stoppedMinutes"
    @logged="onLogged"
  />
</template>

<script setup lang="ts">
import type { GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTimer } from '~/composables/useTimer'

const props = defineProps<{ sidebarCollapsed?: boolean }>()

const { timerState, elapsedDisplay, isRunning, start, stop } = useTimer()
const tasksApi = useTasksApi()

// ── Stop & log ─────────────────────────────────────────────────
const showLog        = ref(false)
const stoppedTaskId  = ref<string | undefined>()
const stoppedMinutes = ref<number | undefined>()

function handleStop() {
  const result = stop()
  if (!result) return
  stoppedTaskId.value  = result.taskId
  stoppedMinutes.value = result.minutes
  showLog.value = true
}
function onLogged() {
  stoppedTaskId.value  = undefined
  stoppedMinutes.value = undefined
}

// ── Task picker ─────────────────────────────────────────────────
const showPicker   = ref(false)
const pickerSearch = ref('')
const loadingTasks = ref(false)
const teamTasks    = ref<GlobalActiveTask[]>([])
const pickerInput  = ref<HTMLInputElement | null>(null)

const filteredPickerTasks = computed(() => {
  const q = pickerSearch.value.trim().toLowerCase()
  const list = q
    ? teamTasks.value.filter(t =>
        t.code?.toLowerCase().includes(q) ||
        t.title?.toLowerCase().includes(q) ||
        t.project_name?.toLowerCase().includes(q) ||
        t.assigned_to_display_name?.toLowerCase().includes(q) ||
        t.assigned_to_email?.toLowerCase().includes(q),
      )
    : teamTasks.value
  return list.slice(0, 15)
})

async function togglePicker() {
  showPicker.value = !showPicker.value
  if (showPicker.value) {
    loadingTasks.value = true
    try { teamTasks.value = await tasksApi.getTeamActiveTasks() }
    catch { teamTasks.value = [] }
    finally { loadingTasks.value = false }
    await nextTick()
    pickerInput.value?.focus()
  }
}

function startTimerFor(task: GlobalActiveTask) {
  start(task.id, task.title, task.code)
  showPicker.value = false
  pickerSearch.value = ''
}

function taskIcon(type: string) {
  return ({ FEATURE: '⚡', TASK: '✅', BUG: '🐛' } as Record<string, string>)[type] ?? '📋'
}

// close picker when clicking outside
function onDocClick(e: MouseEvent) {
  if (!(e.target as HTMLElement).closest('.picker-panel') && !(e.target as HTMLElement).closest('.nav-link-timer')) {
    showPicker.value = false
  }
}
onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.timer-running {
  @apply px-3 py-2.5 bg-gray-800/80 border border-purple-700/30 rounded-lg transition-all;
}
.nav-link-timer {
  @apply flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-all hover:translate-x-0.5;
}

.picker-panel {
  @apply absolute left-0 right-0 z-50 mt-1 bg-gray-900 border border-gray-700 rounded-xl shadow-2xl shadow-black/60 overflow-hidden;
  top: 100%;
}

.picker-drop-enter-active, .picker-drop-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.picker-drop-enter-from, .picker-drop-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
