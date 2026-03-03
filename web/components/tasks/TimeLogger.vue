<template>
  <div class="time-logger">
    <!-- Summary Bar -->
    <div class="flex items-center justify-between mb-4 p-3 bg-gray-800/60 rounded-xl border border-gray-700/50">
      <div class="flex items-center gap-6">
        <div class="text-center">
          <div class="text-lg font-bold text-indigo-400">{{ totalLoggedHours }}h</div>
          <div class="text-xs text-gray-500">Logged</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-bold text-gray-400">{{ estimatedHours }}h</div>
          <div class="text-xs text-gray-500">Estimated</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-bold" :class="varianceClass">{{ varianceLabel }}</div>
          <div class="text-xs text-gray-500">Variance</div>
        </div>
      </div>
      <!-- Progress bar -->
      <div class="flex-1 max-w-48 ml-6">
        <div class="flex justify-between text-xs text-gray-500 mb-1">
          <span>Progress</span>
          <span>{{ progressPct }}%</span>
        </div>
        <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all"
            :class="progressPct > 100 ? 'bg-red-500' : progressPct > 80 ? 'bg-yellow-500' : 'bg-indigo-500'"
            :style="{ width: Math.min(progressPct, 100) + '%' }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Log Work Form -->
    <div class="bg-gray-800/60 rounded-xl border border-gray-700/50 p-4 mb-4">
      <h4 class="text-sm font-semibold text-gray-300 mb-3">Log Work</h4>
      <div class="flex gap-3">
        <div class="flex-1">
          <div class="flex gap-3 mb-2">
            <div class="flex-1">
              <label class="text-xs text-gray-500 mb-1 block">Time Spent</label>
              <div class="flex gap-2">
                <input
                  v-model.number="logHours"
                  type="number"
                  min="0"
                  placeholder="0"
                  class="input-field w-16 text-center"
                />
                <span class="self-center text-gray-500 text-xs">h</span>
                <input
                  v-model.number="logMins"
                  type="number"
                  min="0"
                  max="59"
                  placeholder="0"
                  class="input-field w-16 text-center"
                />
                <span class="self-center text-gray-500 text-xs">m</span>
              </div>
            </div>
          </div>
          <div>
            <label class="text-xs text-gray-500 mb-1 block">Description (optional)</label>
            <input
              v-model="logDescription"
              type="text"
              placeholder="What did you work on?"
              class="input-field w-full"
            />
          </div>
        </div>
        <div class="flex items-end">
          <button
            @click="submitLog"
            :disabled="!totalMinutes || loading"
            class="btn-primary px-4 py-2 text-sm disabled:opacity-40"
          >
            {{ loading ? '...' : 'Log' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Time Log History -->
    <div v-if="timeLogs.length" class="space-y-2">
      <h4 class="text-xs text-gray-500 uppercase tracking-wide mb-3">Work Log</h4>
      <div
        v-for="log in timeLogs"
        :key="log.id"
        class="flex items-center justify-between py-2 px-3 bg-gray-800/40 rounded-lg border border-gray-700/30 hover:border-gray-700 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div class="w-6 h-6 rounded-full bg-indigo-700 flex items-center justify-center text-white text-[10px] font-bold">
            {{ (log.user_email || String(log.user_id)).charAt(0).toUpperCase() }}
          </div>
          <div>
            <span class="text-sm text-gray-300">{{ formatMinutes(log.minutes) }}</span>
            <span v-if="log.description" class="text-xs text-gray-500 ml-2">— {{ log.description }}</span>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs text-gray-500">{{ log.user_email || `User #${log.user_id}` }}</span>
          <span class="text-xs text-gray-600">{{ formatDate(log.logged_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{
  timeLogs: TimeLog[]
  estimatedMinutes: number
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'log-time', minutes: number, description: string): void
}>()

const logHours = ref(0)
const logMins = ref(0)
const logDescription = ref('')

const totalMinutes = computed(() => logHours.value * 60 + logMins.value)

const totalLoggedMinutes = computed(() => props.timeLogs.reduce((s, l) => s + l.minutes, 0))
const totalLoggedHours = computed(() => (totalLoggedMinutes.value / 60).toFixed(1))
const estimatedHours = computed(() => (props.estimatedMinutes / 60).toFixed(1))

const progressPct = computed(() => {
  if (!props.estimatedMinutes) return 0
  return Math.round((totalLoggedMinutes.value / props.estimatedMinutes) * 100)
})

const variance = computed(() => totalLoggedMinutes.value - props.estimatedMinutes)
const varianceLabel = computed(() => {
  if (!props.estimatedMinutes) return '—'
  const h = Math.abs(variance.value) / 60
  const sign = variance.value > 0 ? '+' : '-'
  return `${sign}${h.toFixed(1)}h`
})
const varianceClass = computed(() => {
  if (!props.estimatedMinutes) return 'text-gray-400'
  if (variance.value > 0) return 'text-red-400'
  return 'text-green-400'
})

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

function submitLog() {
  if (!totalMinutes.value) return
  emit('log-time', totalMinutes.value, logDescription.value.trim())
  logHours.value = 0
  logMins.value = 0
  logDescription.value = ''
}
</script>

<style scoped>
.input-field {
  @apply bg-gray-700 border border-gray-600 rounded-lg px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-indigo-500 transition-colors;
}
.btn-primary {
  @apply bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition-colors;
}
</style>
