<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <!-- Page Header -->
    <header class="sticky top-0 z-10 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 py-5 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-600 to-indigo-600 shadow-lg text-xl">
              📡
            </div>
            <div>
              <h1 class="text-xl font-bold tracking-tight text-white sm:text-2xl">Daily Standup</h1>
              <p class="text-xs text-gray-400 mt-0.5">Async standup & team activity tracker</p>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <!-- My check-in status badge (CEO & SUPPORT are exempt) -->
            <div
              v-if="currentUser && !exemptFromPulse"
              :class="[
                'flex items-center gap-2 rounded-xl border px-3 py-1.5 text-xs font-semibold',
                hasCheckedIn
                  ? 'border-emerald-700 bg-emerald-900/40 text-emerald-300'
                  : 'border-amber-700 bg-amber-900/40 text-amber-300'
              ]"
            >
              <span>{{ hasCheckedIn ? '✅' : '⏳' }}</span>
              {{ hasCheckedIn ? 'Checked in today' : 'Not checked in' }}
            </div>


            <button
              v-if="!exemptFromPulse && !hasCheckedIn"
              @click="checkinModal?.open()"
              class="flex items-center gap-2 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-lg transition hover:from-violet-500 hover:to-indigo-500"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Check in now
            </button>
            <button
              v-else-if="!exemptFromPulse"
              @click="checkinModal?.open()"
              class="flex items-center gap-2 rounded-xl border border-gray-600 bg-gray-800 px-4 py-2 text-sm font-medium text-gray-300 transition hover:bg-gray-700"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Update standup
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main content -->
    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <TeamPulseBoard />
    </main>

    <!-- Check-in Modal -->
    <DailyCheckinModal
      ref="checkinModal"
      :forced="false"
      @submitted="onCheckinSubmitted"
    />
  </div>
</template>

<script setup lang="ts">
import TeamPulseBoard from '~/core/modules/pulse/ui/TeamPulseBoard.vue'
import DailyCheckinModal from '~/core/modules/pulse/ui/DailyCheckinModal.vue'
import { usePulseStore } from '~/core/modules/pulse/store/pulse-store'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const { currentUser } = useAuth()
const store = usePulseStore()
const checkinModal = ref<InstanceType<typeof DailyCheckinModal> | null>(null)

import { localDateStr } from '~/composables/useLocalDate'
const today = localDateStr()

const exemptFromPulse = computed(() => {
  const r = currentUser.value?.role?.toUpperCase()
  return r === 'CEO' || r === 'SUPPORT'
})

const hasCheckedIn = computed(() => {
  if (!store.pulse || !currentUser.value) return false
  const uid = currentUser.value.user_id
  return store.pulse.members.some((m) => m.user_id === uid && m.standup !== null)
})

function onCheckinSubmitted() {}

onMounted(() => {
  store.fetchDailyPulse(today)
})
</script>
