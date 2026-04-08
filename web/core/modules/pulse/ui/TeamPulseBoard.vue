<template>
  <div class="w-full space-y-6">
    <!-- ── Header ──────────────────────────────────────────────────────────── -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Team Pulse Board</h2>
        <p class="mt-0.5 text-sm text-gray-400">Async daily check-in for {{ displayDate }}</p>
      </div>
      <div class="flex items-center gap-3">
        <UiDatePicker
          v-model="selectedDate"
          placeholder="Select date…"
          @update:modelValue="onDateChange"
        />
        <button
          class="rounded-lg bg-indigo-600 px-4 py-1.5 text-sm font-medium text-white transition hover:bg-indigo-500 disabled:opacity-50 flex items-center gap-1.5"
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
    <div v-if="store.pulse" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div class="flex flex-col gap-1 rounded-xl border border-gray-700 bg-gray-800 px-4 py-3">
        <span class="text-lg">👥</span>
        <span class="text-2xl font-bold text-white">{{ store.pulse.total_members }}</span>
        <span class="text-xs text-gray-400 uppercase tracking-wide">Team Size</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-700 bg-gray-800 px-4 py-3">
        <span class="text-lg">✅</span>
        <span class="text-2xl font-bold text-white">{{ store.pulse.checked_in }} / {{ store.pulse.total_members }}</span>
        <span class="text-xs text-gray-400 uppercase tracking-wide">Checked In</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-700 bg-gray-800 px-4 py-3">
        <span class="text-lg">📊</span>
        <span class="text-2xl font-bold text-white">{{ store.checkinRate }}%</span>
        <span class="text-xs text-gray-400 uppercase tracking-wide">Check-in Rate</span>
      </div>
      <div class="flex flex-col gap-1 rounded-xl border border-gray-700 bg-gray-800 px-4 py-3">
        <span class="text-lg">⏱️</span>
        <span class="text-2xl font-bold text-white">{{ (store.pulse.total_minutes_logged / 60).toFixed(1) }}h</span>
        <span class="text-xs text-gray-400 uppercase tracking-wide">Total Hours Logged</span>
      </div>
    </div>

    <!-- ── Error ───────────────────────────────────────────────────────────── -->
    <div
      v-if="store.error"
      class="rounded-lg border border-red-700 bg-red-900/40 px-4 py-3 text-sm text-red-300"
    >
      {{ store.error }}
    </div>

    <!-- ── Skeleton ────────────────────────────────────────────────────────── -->
    <div v-if="store.loading && !store.pulse" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 6" :key="i" class="h-52 animate-pulse rounded-xl bg-gray-800" />
    </div>

    <!-- ── Member Grid ─────────────────────────────────────────────────────── -->
    <div v-if="store.pulse" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <PulseMemberCard
        v-for="member in store.pulse.members"
        :key="member.user_id"
        :member="member"
      />
    </div>

    <!-- ── Empty state ─────────────────────────────────────────────────────── -->
    <div
      v-if="store.pulse && store.pulse.members.length === 0"
      class="rounded-xl border border-gray-700 bg-gray-800/50 py-16 text-center text-gray-400"
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

const store = usePulseStore()

const today = localDateStr()
const selectedDate = ref(today)
const pollMs = 10000 // Polling so other users' check-ins appear without manual refresh
let pollTimer: ReturnType<typeof setInterval> | null = null

const displayDate = computed(() => {
  const d = new Date(selectedDate.value + 'T00:00:00')
  return d.toLocaleDateString('en-GB', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
})

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
