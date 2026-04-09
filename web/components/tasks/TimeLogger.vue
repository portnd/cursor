<template>
  <div class="time-logger">
    <!-- Summary Bar -->
    <div class="flex items-center justify-between mb-4 p-3 bg-gray-50 dark:bg-gray-800/60 rounded-xl border border-gray-200 dark:border-gray-700/50">
      <div class="flex items-center gap-6">
        <div class="text-center">
          <div class="text-lg font-bold text-purple-400">{{ totalLoggedHours }}h</div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Logged</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-bold text-gray-400">{{ estimatedHours }}h</div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Estimated</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-bold" :class="varianceClass">{{ varianceLabel }}</div>
          <div class="text-xs text-gray-500 dark:text-gray-400">Variance</div>
        </div>
      </div>
      <!-- Progress bar -->
      <div class="flex-1 max-w-48 ml-6">
        <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mb-1">
          <span>Progress</span>
          <span>{{ progressPct }}%</span>
        </div>
        <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all"
            :class="progressPct > 100 ? 'bg-red-500' : progressPct > 80 ? 'bg-yellow-500' : 'bg-purple-500'"
            :style="{ width: Math.min(progressPct, 100) + '%' }"
          />
        </div>
      </div>
    </div>

    <!-- Daily Quota Bar -->
    <div v-if="dailySummary" class="mb-4 p-3 bg-cyan-50 dark:bg-cyan-900/10 border border-cyan-200 dark:border-cyan-700/20 rounded-xl">
      <div class="flex justify-between items-center mb-1.5">
        <span class="text-xs text-cyan-700 dark:text-cyan-400 font-medium">Today's Log</span>
        <span class="text-xs text-cyan-700 dark:text-cyan-300">
          {{ formatMinutes(dailySummary.total_minutes) }} / 8h
        </span>
      </div>
      <div class="h-2 bg-cyan-100 dark:bg-gray-700/60 rounded-full overflow-hidden">
        <div
          class="h-full rounded-full transition-all"
          :class="dailyPct > 100 ? 'bg-red-400' : 'bg-gradient-to-r from-cyan-500 to-purple-500'"
          :style="{ width: Math.min(dailyPct, 100) + '%' }"
        />
      </div>
    </div>

    <!-- Log Work Form -->
    <div class="bg-gray-50 dark:bg-gray-800/60 rounded-xl border border-gray-200 dark:border-gray-700/50 p-4 mb-4">
      <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">Log Work</h4>

      <!-- Work Type -->
      <div class="mb-3">
        <label class="text-xs text-gray-500 mb-1.5 block">Work Type</label>
        <div class="flex gap-1.5 flex-wrap">
          <button
            v-for="wt in WORK_TYPES"
            :key="wt.value"
            type="button"
            @click="workType = wt.value"
            :class="[
              'px-2.5 py-1 rounded-lg text-xs font-medium border transition-all',
              workType === wt.value
                ? 'bg-purple-100 dark:bg-purple-600/30 border-purple-300 dark:border-purple-500/60 text-purple-700 dark:text-purple-300'
                : 'bg-white dark:bg-gray-700/50 border-gray-300 dark:border-gray-600/40 text-gray-600 dark:text-gray-400 hover:border-purple-300 dark:hover:border-gray-500 hover:text-purple-700 dark:hover:text-gray-200'
            ]"
          >
            {{ wt.emoji }} {{ wt.label }}
          </button>
        </div>
      </div>

      <!-- Time + Presets -->
      <div class="mb-3">
        <label class="text-xs text-gray-500 mb-1.5 block">Time Spent</label>
        <div class="flex gap-2 items-center mb-2">
          <input
            v-model.number="logHours"
            type="number"
            min="0"
            max="16"
            placeholder="0"
            class="input-field w-14 text-center"
          />
          <span class="text-gray-500 text-xs">h</span>
          <input
            v-model.number="logMins"
            type="number"
            min="0"
            max="59"
            placeholder="0"
            class="input-field w-14 text-center"
          />
          <span class="text-gray-500 text-xs">m</span>
          <span v-if="totalMinutes > 0" class="ml-auto text-xs text-purple-400 font-medium">
            {{ formatMinutes(totalMinutes) }}
          </span>
        </div>
        <!-- Quick presets -->
        <div class="flex gap-1.5">
          <button
            v-for="p in presets"
            :key="p.min"
            type="button"
            @click="applyPreset(p.min)"
            class="preset-btn"
          >
            {{ p.label }}
          </button>
        </div>
      </div>

      <!-- Description -->
      <div class="mb-3">
        <label class="text-xs text-gray-500 mb-1 block">Description <span class="text-gray-600">(optional)</span></label>
        <input
          v-model="logDescription"
          type="text"
          placeholder="What did you work on?"
          class="input-field w-full"
        />
      </div>

      <!-- Date picker — backfill up to 7 days -->
      <div class="mb-3">
        <label class="text-xs text-gray-500 mb-1 block">Date</label>
        <select v-model="loggedDate" class="input-field w-full text-sm">
          <option v-for="d in dateOptions" :key="d.value" :value="d.value">{{ d.label }}</option>
        </select>
      </div>

      <button
        type="button"
        @click="submitLog"
        :disabled="!totalMinutes || loading"
        class="btn-primary w-full py-2 text-sm disabled:opacity-40"
      >
        {{ loading ? 'Logging...' : 'Log Time' }}
      </button>
    </div>

    <!-- Time Log History -->
    <div v-if="timeLogs.length" class="space-y-2">
      <h4 class="text-xs text-gray-500 uppercase tracking-wide mb-3">Work Log</h4>
      <div
        v-for="log in timeLogs"
        :key="log.id"
        class="group flex items-center justify-between py-2 px-3 bg-gray-800/40 rounded-lg border border-gray-700/30 hover:border-gray-700 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div class="w-6 h-6 rounded-full bg-purple-700 flex items-center justify-center text-white text-[10px] font-bold shrink-0">
            {{ (log.user_email || String(log.user_id)).charAt(0).toUpperCase() }}
          </div>
          <div>
            <div class="flex items-center gap-1.5">
              <span class="text-sm text-gray-300">{{ formatMinutes(log.minutes) }}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded bg-gray-700/60 text-gray-400 font-medium">
                {{ workTypeEmoji(log.work_type) }} {{ log.work_type || 'DEV' }}
              </span>
              <span v-if="log.is_timer_session" class="text-[10px] text-cyan-500" title="Timer session">⏱</span>
            </div>
            <span v-if="log.description" class="text-xs text-gray-500">{{ log.description }}</span>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs text-gray-500">{{ log.user_email || `User #${log.user_id}` }}</span>
          <span class="text-xs text-gray-600">{{ formatDate(log.logged_at) }}</span>
          <!-- Edit / Delete (own logs only) -->
          <div v-if="isOwn(log)" class="hidden group-hover:flex items-center gap-1">
            <button
              type="button"
              @click="startEdit(log)"
              class="p-1 rounded text-gray-500 hover:text-purple-400 hover:bg-purple-900/20 transition-colors"
              title="Edit"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
            <button
              type="button"
              @click="confirmDelete(log)"
              class="p-1 rounded text-gray-500 hover:text-red-400 hover:bg-red-900/20 transition-colors"
              title="Delete"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Log Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="editingLog" class="fixed inset-0 z-50 flex items-center justify-center p-4" @mousedown.self="editingLog = null">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" />
          <div class="relative w-full max-w-sm bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl p-5">
            <h3 class="text-sm font-bold text-gray-100 mb-4">Edit Time Log</h3>
            <div class="space-y-3">
              <div>
                <label class="text-xs text-gray-500 mb-1 block">Work Type</label>
                <div class="flex gap-1.5 flex-wrap">
                  <button
                    v-for="wt in WORK_TYPES" :key="wt.value" type="button"
                    @click="editWorkType = wt.value"
                    :class="['px-2.5 py-1 rounded-lg text-xs font-medium border transition-all',
                      editWorkType === wt.value
                        ? 'bg-purple-100 dark:bg-purple-600/30 border-purple-300 dark:border-purple-500/60 text-purple-700 dark:text-purple-300'
                        : 'bg-white dark:bg-gray-700/50 border-gray-300 dark:border-gray-600/40 text-gray-600 dark:text-gray-400 hover:border-purple-300 dark:hover:border-gray-500 hover:text-purple-700 dark:hover:text-gray-200']"
                  >{{ wt.emoji }} {{ wt.label }}</button>
                </div>
              </div>
              <div class="flex gap-3 items-end">
                <div>
                  <label class="text-xs text-gray-500 mb-1 block">Hours</label>
                  <input v-model.number="editHours" type="number" min="0" max="16" class="input-field w-14 text-center" />
                </div>
                <div>
                  <label class="text-xs text-gray-500 mb-1 block">Mins</label>
                  <input v-model.number="editMins" type="number" min="0" max="59" class="input-field w-14 text-center" />
                </div>
              </div>
              <div>
                <label class="text-xs text-gray-500 mb-1 block">Description</label>
                <input v-model="editDesc" type="text" class="input-field w-full" />
              </div>
            </div>
            <div class="flex gap-2 mt-4">
              <button type="button" @click="editingLog = null" class="flex-1 px-4 py-2 text-sm text-gray-400 hover:text-gray-200 transition-colors">
                Cancel
              </button>
              <button type="button" @click="saveEdit" :disabled="editTotalMinutes <= 0 || editSaving" class="flex-1 btn-primary py-2 text-sm disabled:opacity-40">
                {{ editSaving ? 'Saving...' : 'Save' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import type { TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'
import { WORK_TYPES, useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useAuth } from '~/composables/useAuth'

const props = defineProps<{
  timeLogs: TimeLog[]
  estimatedMinutes: number
  taskId: string
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'log-time', minutes: number, description: string, workType: string, loggedDate: string, isTimer: boolean): void
  (e: 'refresh'): void
}>()

const { currentUser } = useAuth()
const tasksApi = useTasksApi()

// ── Form state ──────────────────────────────────────────────
const logHours = ref(0)
const logMins = ref(0)
const logDescription = ref('')
const workType = ref('DEV')
const loggedDate = ref(todayValue())
const totalMinutes = computed(() => logHours.value * 60 + logMins.value)

const presets = [
  { label: '+15m', min: 15 },
  { label: '+30m', min: 30 },
  { label: '+1h',  min: 60 },
  { label: '+2h',  min: 120 },
  { label: '+4h',  min: 240 },
]

function applyPreset(min: number) {
  const total = totalMinutes.value + min
  logHours.value = Math.floor(total / 60)
  logMins.value = total % 60
}

// ── Date options (today + 6 past days) ──────────────────────
const dateOptions = computed(() => {
  const opts: { label: string; value: string }[] = []
  for (let i = 0; i < 7; i++) {
    const d = new Date()
    d.setDate(d.getDate() - i)
    const v = d.toISOString().slice(0, 10)
    const label = i === 0 ? `Today (${v})` : i === 1 ? `Yesterday (${v})` : v
    opts.push({ label, value: v })
  }
  return opts
})

function todayValue() {
  return new Date().toISOString().slice(0, 10)
}

// ── Stats ────────────────────────────────────────────────────
const totalLoggedMinutes = computed(() => props.timeLogs.reduce((s, l) => s + l.minutes, 0))
const totalLoggedHours   = computed(() => (totalLoggedMinutes.value / 60).toFixed(1))
const estimatedHours     = computed(() => (props.estimatedMinutes / 60).toFixed(1))
const progressPct        = computed(() => {
  if (!props.estimatedMinutes) return 0
  return Math.round((totalLoggedMinutes.value / props.estimatedMinutes) * 100)
})
const variance      = computed(() => totalLoggedMinutes.value - props.estimatedMinutes)
const varianceLabel = computed(() => {
  if (!props.estimatedMinutes) return '—'
  const h = Math.abs(variance.value) / 60
  return `${variance.value > 0 ? '+' : '-'}${h.toFixed(1)}h`
})
const varianceClass = computed(() => {
  if (!props.estimatedMinutes) return 'text-gray-400'
  return variance.value > 0 ? 'text-red-400' : 'text-green-400'
})

// ── Daily quota ──────────────────────────────────────────────
const dailySummary = ref<{ total_minutes: number } | null>(null)
const dailyPct     = computed(() => Math.round(((dailySummary.value?.total_minutes ?? 0) / 480) * 100))

async function loadDailySummary() {
  try {
    dailySummary.value = await tasksApi.getMyDailyTimeLogs(todayValue())
  }
  catch { /* non-critical */ }
}
onMounted(loadDailySummary)

// ── Submit ───────────────────────────────────────────────────
function submitLog() {
  if (!totalMinutes.value) return
  emit('log-time', totalMinutes.value, logDescription.value.trim(), workType.value, loggedDate.value, false)
  logHours.value = 0
  logMins.value  = 0
  logDescription.value = ''
  workType.value  = 'DEV'
  loggedDate.value = todayValue()
  loadDailySummary()
}

// ── Edit log ─────────────────────────────────────────────────
const editingLog   = ref<TimeLog | null>(null)
const editHours    = ref(0)
const editMins     = ref(0)
const editDesc     = ref('')
const editWorkType = ref('DEV')
const editSaving   = ref(false)
const editTotalMinutes = computed(() => editHours.value * 60 + editMins.value)

function startEdit(log: TimeLog) {
  editingLog.value   = log
  editHours.value    = Math.floor(log.minutes / 60)
  editMins.value     = log.minutes % 60
  editDesc.value     = log.description || ''
  editWorkType.value = log.work_type || 'DEV'
}

async function saveEdit() {
  if (!editingLog.value || editTotalMinutes.value <= 0) return
  editSaving.value = true
  try {
    await tasksApi.editTimeLog(editingLog.value.id, editTotalMinutes.value, editDesc.value.trim(), editWorkType.value)
    editingLog.value = null
    emit('refresh')
    loadDailySummary()
  }
  catch (e) {
    console.error('Failed to edit time log:', e)
  }
  finally {
    editSaving.value = false
  }
}

// ── Delete log ───────────────────────────────────────────────
async function confirmDelete(log: TimeLog) {
  if (!confirm(`Delete this ${formatMinutes(log.minutes)} log?`)) return
  try {
    await tasksApi.deleteTimeLog(log.id)
    emit('refresh')
    loadDailySummary()
  }
  catch (e) {
    console.error('Failed to delete time log:', e)
  }
}

// ── Helpers ──────────────────────────────────────────────────
function isOwn(log: TimeLog) {
  return currentUser.value?.id === log.user_id
}

function workTypeEmoji(wt: string) {
  return WORK_TYPES.find(w => w.value === wt)?.emoji ?? '📋'
}

function formatMinutes(mins: number) {
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.input-field {
  @apply bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:border-purple-500 transition-colors;
}
.btn-primary {
  @apply bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-lg font-medium transition-colors;
}
.preset-btn {
  @apply bg-white dark:bg-gray-700/60 border border-gray-300 dark:border-gray-600/40 text-purple-600 dark:text-purple-400 hover:bg-purple-50 dark:hover:bg-purple-900/20 hover:border-purple-400 dark:hover:border-purple-500/40 rounded-lg px-2.5 py-1 text-xs font-medium transition-colors;
}
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.2s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
</style>
