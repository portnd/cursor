<template>
  <div class="w-full space-y-6">
    <!-- ── Header ──────────────────────────────────────────────────────────── -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900 dark:text-white tracking-tight">Team Pulse Board</h2>
        <p class="mt-0.5 text-sm text-gray-500 dark:text-gray-400">Async daily check-in for {{ displayDate }}</p>
      </div>
      <div class="flex items-center gap-3">
        <UiDatePicker
          v-model="selectedDate"
          placeholder="Select date…"
          @update:modelValue="onDateChange"
        />
        <button
          class="rounded-lg bg-indigo-100 dark:bg-indigo-600 px-4 py-1.5 text-sm font-medium text-gray-900 dark:text-white transition hover:bg-indigo-100 dark:bg-indigo-500 disabled:opacity-50 flex items-center gap-1.5"
          :disabled="store.loading"
          @click="store.fetchDailyPulse(selectedDate)"
        >
          <svg v-if="store.loading" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
          </svg>
          {{ store.loading ? 'Loading…' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- ── Summary Stats ───────────────────────────────────────────────────── -->
    <div v-if="store.pulse" class="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-5">
      <div class="flex flex-col gap-1 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3">
        <span class="text-lg">👥</span>
        <span class="text-2xl font-bold text-gray-900 dark:text-white">{{ visibleTeamSize }}</span>
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Team Size</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-orange-200 dark:border-orange-700 bg-orange-50/60 dark:bg-orange-900/20 px-4 py-3">
        <span class="text-lg">🏖️</span>
        <span class="text-2xl font-bold text-orange-700 dark:text-orange-300">{{ visibleOnLeaveCount }}</span>
        <span class="text-xs text-orange-700/80 dark:text-orange-300/80 uppercase tracking-wide">On Leave</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3">
        <span class="text-lg">✅</span>
        <span class="text-2xl font-bold text-gray-900 dark:text-white">{{ visibleCheckedIn }} / {{ activeMembers }}</span>
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Checked In (Active)</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3">
        <span class="text-lg">📊</span>
        <span class="text-2xl font-bold text-gray-900 dark:text-white">{{ visibleCheckinRate }}%</span>
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Check-in Rate</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3">
        <span class="text-lg">⏱️</span>
        <span class="text-2xl font-bold text-gray-900 dark:text-white">{{ (visibleTotalMinutesLogged / 60).toFixed(1) }}h</span>
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Total Hours Logged</span>
      </div>
    </div>

    <!-- ── Error ───────────────────────────────────────────────────────────── -->
    <div
      v-if="store.error"
      class="rounded-lg border border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-900/40 px-4 py-3 text-sm text-red-600 dark:text-red-300"
    >
      {{ store.error }}
    </div>

    <!-- ── Skeleton ────────────────────────────────────────────────────────── -->
    <div v-if="store.loading && !store.pulse" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 6" :key="i" class="h-52 animate-pulse rounded-xl bg-gray-200 dark:bg-gray-800" />
    </div>

    <!-- ── Visibility Manager (CEO) ────────────────────────────────────────── -->
    <div
      v-if="store.pulse && canManageVisibility"
      class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3"
    >
      <div class="flex items-center justify-between gap-3">
        <p class="text-sm font-semibold text-gray-900 dark:text-white">Team size visibility (CEO)</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">Checked users are hidden from Team Size</p>
      </div>
      <div class="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2 xl:grid-cols-3">
        <label
          v-for="member in store.pulse.members"
          :key="`vis-${member.user_id}`"
          class="flex items-center gap-2 rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2 text-sm text-gray-700 dark:text-gray-200"
        >
          <input
            type="checkbox"
            class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
            :checked="isHidden(member.user_id)"
            @change="onToggleVisibility(member.user_id, ($event.target as HTMLInputElement).checked)"
          >
          <span class="truncate">{{ member.user_display_name || member.user_email }}</span>
        </label>
      </div>
    </div>

    <!-- ── Member Grid ─────────────────────────────────────────────────────── -->
    <div v-if="store.pulse" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <PulseMemberCard
        v-for="member in visibleMembers"
        :key="member.user_id"
        :member="member"
      />
    </div>

    <!-- ── Empty state ─────────────────────────────────────────────────────── -->
    <div
      v-if="store.pulse && store.pulse.members.length === 0"
      class="rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 py-16 text-center text-gray-500 dark:text-gray-400"
    >
      <p class="text-3xl">📭</p>
      <p class="mt-2 text-sm">No team members found.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { usePulseStore } from '../store/pulse-store'
import PulseMemberCard from './PulseMemberCard.vue'
import { localDateStr } from '~/composables/useLocalDate'

const props = withDefaults(defineProps<{
  canManageVisibility?: boolean
  hiddenUserIds?: number[]
}>(), {
  canManageVisibility: false,
  hiddenUserIds: () => [],
})

const emit = defineEmits<{
  (e: 'toggle-user-visibility', payload: { userId: number; hidden: boolean }): void
}>()

const store = usePulseStore()

const today = localDateStr()
const selectedDate = ref(today)
const pollMs = 10000 // Polling so other users' check-ins appear without manual refresh
let pollTimer: ReturnType<typeof setInterval> | null = null

const displayDate = computed(() => {
  const d = new Date(selectedDate.value + 'T00:00:00')
  return d.toLocaleDateString('en-GB', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
})

const hiddenSet = computed(() => new Set(props.hiddenUserIds))

const visibleMembers = computed(() => {
  const members = store.pulse?.members ?? []
  return members.filter((m) => !hiddenSet.value.has(m.user_id))
})

const visibleTeamSize = computed(() => visibleMembers.value.length)

const visibleOnLeaveCount = computed(() =>
  visibleMembers.value.filter((m) => m.is_on_leave).length,
)

const activeMembers = computed(() => Math.max(visibleTeamSize.value - visibleOnLeaveCount.value, 0))

const visibleCheckedIn = computed(() =>
  visibleMembers.value.filter((m) => !m.is_on_leave && m.standup !== null).length,
)

const visibleCheckinRate = computed(() => {
  if (activeMembers.value === 0) return 0
  return Math.round((visibleCheckedIn.value / activeMembers.value) * 100)
})

const visibleTotalMinutesLogged = computed(() =>
  visibleMembers.value.reduce((sum, m) => sum + (m.total_logged_minutes ?? 0), 0),
)

function isHidden(userId: number) {
  return hiddenSet.value.has(userId)
}

function onToggleVisibility(userId: number, checked: boolean) {
  emit('toggle-user-visibility', { userId, hidden: checked })
}

function onDateChange() {
  store.fetchDailyPulse(selectedDate.value)
}

function startPolling() {
  stopPolling()
  pollTimer = setInterval(async () => {
    // Avoid background polling when tab is hidden (saves both CPU + API load).
    if (typeof document !== 'undefined' && document.visibilityState !== 'visible') return
    // Keep at most one in-flight request.
    if (store.loading) return
    // Only auto-refresh for today's pulse; historical dates can be refreshed manually.
    if (selectedDate.value !== today) return
    await store.fetchDailyPulse(selectedDate.value)
  }, pollMs)
}

function stopPolling() {
  if (!pollTimer) return
  clearInterval(pollTimer)
  pollTimer = null
}

onMounted(() => {
  if (!store.pulse || store.lastFetchedDate !== selectedDate.value) {
    store.fetchDailyPulse(selectedDate.value)
  }

  startPolling()
})

watch(selectedDate, (d) => {
  if (d === today) startPolling()
  else stopPolling()
})

onBeforeUnmount(() => {
  stopPolling()
})
</script>
