<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @mousedown.self="$emit('update:modelValue', false)"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" />

        <!-- Panel -->
        <div class="relative w-full max-w-lg bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl shadow-black/50 flex flex-col max-h-[90vh]">
          <!-- Header -->
          <div class="flex items-center justify-between px-5 py-4 border-b border-gray-700/60">
            <div class="flex items-center gap-2.5">
              <span class="text-purple-400 text-lg">⏱️</span>
              <div>
                <h2 class="text-sm font-bold text-gray-100">Quick Log Time</h2>
                <p class="text-xs text-gray-500">Log time against any active task in your team</p>
              </div>
            </div>
            <button type="button" @click="$emit('update:modelValue', false)"
              class="p-1.5 rounded-lg text-gray-500 hover:text-gray-300 hover:bg-gray-700 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-y-auto px-5 py-4 space-y-4">

            <!-- Timer fill notice (shown when modal opened from timer stop) -->
            <div v-if="fromTimer" class="flex items-center gap-2 p-2.5 bg-cyan-900/20 border border-cyan-700/30 rounded-lg">
              <span class="text-cyan-400 text-sm">⏱</span>
              <p class="text-xs text-cyan-300/80">Pre-filled from timer session — adjust if needed.</p>
            </div>

            <!-- Step 1: Task Search -->
            <div>
              <label class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2 block">Search Task</label>
              <div class="relative">
                <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search by task ID or title..."
                  class="w-full pl-9 pr-4 py-2.5 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
                  @focus="showDropdown = true"
                  @input="showDropdown = true"
                  @blur="onSearchBlur"
                />
              </div>

              <!-- Task Dropdown -->
              <Transition name="dropdown">
                <div
                  v-if="showDropdown && (filteredTasks.length || loadingTasks || teamTasks.length === 0)"
                  class="mt-1.5 bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-xl max-h-52 overflow-y-auto"
                >
                  <div v-if="loadingTasks" class="flex items-center justify-center py-6 text-gray-500 text-sm gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    Loading tasks...
                  </div>
                  <template v-else-if="filteredTasks.length">
                    <button
                      v-for="task in filteredTasks" :key="task.id" type="button"
                      class="w-full flex items-start gap-3 px-4 py-3 hover:bg-gray-700/60 transition-colors text-left border-b border-gray-700/40 last:border-0"
                      @click="selectTask(task)"
                    >
                      <span class="mt-0.5 shrink-0 text-sm">{{ taskTypeIcon(task.task_type) }}</span>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 flex-wrap">
                          <span class="text-xs font-mono text-purple-400 shrink-0">{{ task.code }}</span>
                          <span class="text-sm text-gray-200 truncate">{{ task.title }}</span>
                        </div>
                        <div class="flex items-center gap-2 mt-1">
                          <span v-if="task.project_name" class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium text-gray-300 bg-gray-700/80">
                            <span v-if="task.project_color" class="w-1.5 h-1.5 rounded-full shrink-0" :style="{ backgroundColor: task.project_color }" />
                            {{ task.project_name }}
                          </span>
                          <span v-if="assigneeLabel(task)" class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-medium bg-indigo-900/60 text-indigo-300 border border-indigo-700/40">
                            <svg class="w-2.5 h-2.5" fill="currentColor" viewBox="0 0 20 20"><path d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" /></svg>
                            {{ assigneeLabel(task) }}
                          </span>
                          <span v-else class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-gray-700/60 text-gray-500">Unassigned</span>
                        </div>
                      </div>
                    </button>
                  </template>
                  <div v-else class="py-6 text-center text-gray-500 text-sm">
                    <p v-if="searchQuery">No tasks match "{{ searchQuery }}"</p>
                    <p v-else>No active sprint tasks found in your team</p>
                  </div>
                </div>
              </Transition>

              <!-- Selected Task Card -->
              <div v-if="selectedTask" class="mt-2 flex items-center gap-3 p-3 bg-purple-900/20 border border-purple-700/40 rounded-xl">
                <span class="text-lg">{{ taskTypeIcon(selectedTask.task_type) }}</span>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-mono text-purple-400">{{ selectedTask.code }}</span>
                    <span class="text-sm font-medium text-gray-200 truncate">{{ selectedTask.title }}</span>
                  </div>
                  <div class="flex items-center gap-2 mt-1">
                    <span v-if="selectedTask.project_name" class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] text-gray-400 bg-gray-700/60">
                      <span v-if="selectedTask.project_color" class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: selectedTask.project_color }" />
                      {{ selectedTask.project_name }}
                    </span>
                    <span v-if="assigneeLabel(selectedTask)" class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] bg-indigo-900/60 text-indigo-300 border border-indigo-700/40">
                      Assigned to: {{ assigneeLabel(selectedTask) }}
                    </span>
                  </div>
                </div>
                <button type="button" @click="clearTask" class="shrink-0 p-1 rounded-lg text-gray-500 hover:text-red-400 hover:bg-red-900/20 transition-colors">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Step 2: Time entry (shown after task selected) -->
            <Transition name="slide-down">
              <div v-if="selectedTask" class="space-y-3">

                <!-- Work Type -->
                <div>
                  <label class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2 block">Work Type</label>
                  <div class="flex gap-1.5 flex-wrap">
                    <button
                      v-for="wt in WORK_TYPES" :key="wt.value" type="button"
                      @click="workType = wt.value"
                      :class="['px-2.5 py-1 rounded-lg text-xs font-medium border transition-all',
                        workType === wt.value ? 'bg-purple-600/30 border-purple-500/60 text-purple-300' : 'bg-gray-800 border-gray-600/40 text-gray-400 hover:border-gray-500']"
                    >{{ wt.emoji }} {{ wt.label }}</button>
                  </div>
                </div>

                <!-- Time + Presets -->
                <div>
                  <label class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2 block">Time Spent</label>
                  <div class="flex items-center gap-3 mb-2">
                    <div class="flex items-center gap-2">
                      <input v-model.number="logHours" type="number" min="0" max="16" placeholder="0" class="time-input" />
                      <span class="text-xs text-gray-500">h</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <input v-model.number="logMins" type="number" min="0" max="59" placeholder="0" class="time-input" />
                      <span class="text-xs text-gray-500">m</span>
                    </div>
                    <div v-if="totalMinutes > 0" class="ml-auto text-xs text-purple-400 font-medium">
                      {{ formatMinutes(totalMinutes) }}
                    </div>
                  </div>
                  <!-- Presets -->
                  <div class="flex gap-1.5 flex-wrap">
                    <button v-for="p in presets" :key="p.min" type="button" @click="applyPreset(p.min)" class="preset-btn">
                      {{ p.label }}
                    </button>
                  </div>
                </div>

                <!-- Description -->
                <div>
                  <label class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2 block">
                    Description <span class="normal-case text-gray-600">(optional)</span>
                  </label>
                  <input v-model="logDescription" type="text" placeholder="What did you work on?"
                    class="w-full px-3 py-2.5 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
                  />
                </div>

                <!-- Date (backfill up to 7 days) -->
                <div>
                  <label class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2 block">Date</label>
                  <select v-model="loggedDate"
                    class="w-full px-3 py-2.5 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-300 focus:outline-none focus:border-purple-500 transition-colors">
                    <option v-for="d in dateOptions" :key="d.value" :value="d.value">{{ d.label }}</option>
                  </select>
                </div>

                <!-- Swarm notice -->
                <div v-if="selectedTask && assigneeLabel(selectedTask)" class="flex items-start gap-2 p-3 bg-blue-900/20 border border-blue-700/30 rounded-lg">
                  <svg class="w-4 h-4 text-blue-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p class="text-xs text-blue-300/80 leading-relaxed">
                    Time will be credited to <strong class="text-blue-300">you</strong> — not the task owner.
                    Ownership remains with {{ assigneeLabel(selectedTask) }}.
                  </p>
                </div>
              </div>
            </Transition>
          </div>

          <!-- Footer -->
          <div class="px-5 py-4 border-t border-gray-700/60 flex items-center justify-between gap-3">
            <button type="button" @click="$emit('update:modelValue', false)"
              class="px-4 py-2 text-sm text-gray-400 hover:text-gray-200 transition-colors">Cancel</button>
            <button
              type="button"
              :disabled="!canSubmit || submitting"
              @click="submitLog"
              class="px-5 py-2 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 disabled:opacity-40 disabled:cursor-not-allowed text-white text-sm font-medium rounded-xl transition-all flex items-center gap-2"
            >
              <svg v-if="submitting" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {{ submitting ? 'Logging...' : 'Log Time' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'
import { WORK_TYPES, useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTimer } from '~/composables/useTimer'
import { localDateStr } from '~/composables/useLocalDate'

const props = defineProps<{
  modelValue: boolean
  /** Pre-selected task ID — used when opening from Timer stop */
  preselectedTaskId?: string
  /** Pre-filled minutes — used when opening from Timer stop */
  prefilledMinutes?: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
  (e: 'logged'): void
}>()

const tasksApi = useTasksApi()
const { timerState } = useTimer()

const teamTasks    = ref<GlobalActiveTask[]>([])
const loadingTasks = ref(false)
const searchQuery  = ref('')
const showDropdown = ref(false)
const selectedTask = ref<GlobalActiveTask | null>(null)

const logHours       = ref(0)
const logMins        = ref(0)
const logDescription = ref('')
const workType       = ref('DEV')
const loggedDate     = ref(todayValue())
const submitting     = ref(false)
const fromTimer      = ref(false)

const totalMinutes = computed(() => logHours.value * 60 + logMins.value)
const canSubmit    = computed(() => !!selectedTask.value && totalMinutes.value > 0)

const presets = [
  { label: '+15m', min: 15 },
  { label: '+30m', min: 30 },
  { label: '+1h',  min: 60 },
  { label: '+2h',  min: 120 },
]

function applyPreset(min: number) {
  const total = totalMinutes.value + min
  logHours.value = Math.floor(total / 60)
  logMins.value  = total % 60
}

const dateOptions = computed(() => {
  const opts: { label: string; value: string }[] = []
  for (let i = 0; i < 7; i++) {
    const d = new Date()
    d.setDate(d.getDate() - i)
    const v = localDateStr(d)
    opts.push({ label: i === 0 ? `Today (${v})` : i === 1 ? `Yesterday (${v})` : v, value: v })
  }
  return opts
})

function todayValue() { return localDateStr() }

const filteredTasks = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return teamTasks.value.slice(0, 20)
  return teamTasks.value.filter(t =>
    t.code?.toLowerCase().includes(q)
    || t.title?.toLowerCase().includes(q)
    || t.assigned_to_display_name?.toLowerCase().includes(q)
    || t.assigned_to_email?.toLowerCase().includes(q)
    || t.project_name?.toLowerCase().includes(q),
  ).slice(0, 20)
})

function assigneeLabel(task: GlobalActiveTask): string {
  if (task.assigned_to_display_name) return task.assigned_to_display_name
  if (task.assigned_to_email) return task.assigned_to_email.split('@')[0]
  return ''
}

function taskTypeIcon(type: string): string {
  const icons: Record<string, string> = { FEATURE: '⚡', TASK: '✅', BUG: '🐛' }
  return icons[type] || '📋'
}

function formatMinutes(mins: number): string {
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function selectTask(task: GlobalActiveTask) {
  selectedTask.value = task
  searchQuery.value  = ''
  showDropdown.value = false
}

function clearTask() {
  selectedTask.value = null
  logHours.value     = 0
  logMins.value      = 0
  logDescription.value = ''
  workType.value     = 'DEV'
}

async function loadTeamTasks() {
  loadingTasks.value = true
  try { teamTasks.value = await tasksApi.getTeamActiveTasks() }
  catch { teamTasks.value = [] }
  finally { loadingTasks.value = false }
}

async function submitLog() {
  if (!canSubmit.value || !selectedTask.value) return
  submitting.value = true
  try {
    await tasksApi.logTime(
      selectedTask.value.id,
      totalMinutes.value,
      logDescription.value.trim(),
      workType.value,
      loggedDate.value,
      fromTimer.value,
    )
    emit('logged')
    emit('update:modelValue', false)
    clearTask()
  }
  catch (e) { console.error('Failed to log time:', e) }
  finally { submitting.value = false }
}

watch(() => props.modelValue, (open) => {
  if (!open) return
  fromTimer.value = false
  searchQuery.value  = ''
  showDropdown.value = false
  workType.value     = 'DEV'
  loggedDate.value   = todayValue()
  logDescription.value = ''

  // Pre-fill from timer if provided
  if (props.prefilledMinutes && props.prefilledMinutes > 0) {
    logHours.value = Math.floor(props.prefilledMinutes / 60)
    logMins.value  = props.prefilledMinutes % 60
    fromTimer.value = true
  }
  else {
    logHours.value = 0
    logMins.value  = 0
  }

  loadTeamTasks().then(() => {
    // Auto-select task if pre-selected (e.g. from timer)
    if (props.preselectedTaskId) {
      const found = teamTasks.value.find(t => t.id === props.preselectedTaskId)
      if (found) selectedTask.value = found
    }
    else {
      selectedTask.value = null
    }
  })
})

function onSearchBlur() {
  setTimeout(() => { showDropdown.value = false }, 150)
}
</script>

<style scoped>
.time-input {
  @apply w-16 text-center bg-gray-800 border border-gray-600 rounded-xl px-3 py-2.5 text-sm text-gray-200 focus:outline-none focus:border-purple-500 transition-colors;
}
.preset-btn {
  @apply bg-gray-800 border border-gray-600/40 text-purple-400 hover:bg-purple-900/20 hover:border-purple-500/40 rounded-lg px-2.5 py-1 text-xs font-medium transition-colors;
}

.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.2s ease; }
.modal-fade-enter-active .relative, .modal-fade-leave-active .relative { transition: transform 0.2s ease, opacity 0.2s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-fade-enter-from .relative { transform: scale(0.96) translateY(-8px); opacity: 0; }

.dropdown-enter-active, .dropdown-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-4px); }

.slide-down-enter-active, .slide-down-leave-active { transition: opacity 0.2s ease, transform 0.2s ease; }
.slide-down-enter-from, .slide-down-leave-to { opacity: 0; transform: translateY(-6px); }
</style>
