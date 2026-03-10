<template>
  <div class="rounded-lg border border-gray-700 bg-gray-800 overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-700">
      <h3 class="text-lg font-bold text-white">Team Leaderboard</h3>
      <p class="text-xs text-gray-400">Ranked by composite KPI score</p>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-700 text-left text-gray-400">
            <th class="px-4 py-2 font-medium">#</th>
            <th class="px-4 py-2 font-medium">Member</th>
            <th class="px-4 py-2 font-medium">Role</th>
            <th class="px-4 py-2 font-medium text-right">Delivery</th>
            <th class="px-4 py-2 font-medium text-right">Quality</th>
            <th class="px-4 py-2 font-medium text-right">Rework</th>
            <th class="px-4 py-2 font-medium text-right">Score</th>
            <th v-if="isCEO" class="px-4 py-2 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(m, i) in members"
            :key="m.user_id"
            class="border-b border-gray-700/50 hover:bg-gray-700/30 transition-colors"
          >
            <td class="px-4 py-2 text-gray-500">{{ i + 1 }}</td>
            <td class="px-4 py-2 text-white font-medium">{{ m.email }}</td>
            <td class="px-4 py-2">
              <span
                class="rounded px-2 py-0.5 text-xs font-medium"
                :class="roleBadgeClass(m.role)"
              >
                {{ m.role }}
              </span>
            </td>
            <td class="px-4 py-2 text-right">
              <span :class="pctColor(m.delivery_rate_pct)">{{ m.delivery_rate_pct.toFixed(1) }}%</span>
            </td>
            <td class="px-4 py-2 text-right text-gray-300">{{ m.code_quality_index.toFixed(0) }}</td>
            <td class="px-4 py-2 text-right">
              <span :class="reworkColor(m.rework_rate_pct)">{{ m.rework_rate_pct.toFixed(1) }}%</span>
            </td>
            <td class="px-4 py-2 text-right font-bold" :class="scoreColor(m.composite_score)">
              {{ m.composite_score.toFixed(1) }}
            </td>

            <!-- CEO-only action column -->
            <td v-if="isCEO" class="px-4 py-2 text-right">
              <div v-if="confirmingUserId === m.user_id" class="flex items-center justify-end gap-2">
                <span class="text-xs text-amber-400 font-medium">Confirm reset?</span>
                <button
                  @click="confirmReset(m.user_id)"
                  :disabled="resettingUserId === m.user_id"
                  class="rounded px-2 py-0.5 text-xs font-bold bg-red-600 hover:bg-red-500 text-white disabled:opacity-50 transition-colors"
                >
                  {{ resettingUserId === m.user_id ? '…' : 'Yes' }}
                </button>
                <button
                  @click="confirmingUserId = null"
                  class="rounded px-2 py-0.5 text-xs font-semibold text-gray-400 hover:text-white transition-colors"
                >
                  No
                </button>
              </div>
              <button
                v-else
                @click="confirmingUserId = m.user_id"
                class="inline-flex items-center gap-1 rounded px-2 py-1 text-[10px] font-semibold border border-gray-600/60 text-gray-500 hover:border-amber-500/40 hover:text-amber-400 hover:bg-amber-500/10 transition-all"
                title="Reset Rework Rate for this developer"
              >
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                </svg>
                Reset Rework
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="members.length === 0" class="px-4 py-8 text-center text-gray-500">
      No team members to display.
    </div>

    <!-- Success toast -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div
        v-if="toastVisible"
        class="mx-4 mb-3 flex items-center gap-2 rounded-lg border border-emerald-500/30 bg-emerald-950/30 px-3 py-2 text-xs text-emerald-400"
      >
        <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        Rework Rate reset — counter now starts from this moment.
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import type { TeamMemberKPI } from '~/core/modules/performance/performance-api'
import { usePerformanceApi } from '~/core/modules/performance/performance-api'
import { useAuth } from '~/composables/useAuth'

defineProps<{
  members: TeamMemberKPI[]
}>()

const emit = defineEmits<{
  (e: 'refresh'): void
}>()

const { resetReworkRate } = usePerformanceApi()
const { currentUser } = useAuth()

const isCEO = computed(() => currentUser.value?.role === 'CEO')

const confirmingUserId = ref<number | null>(null)
const resettingUserId = ref<number | null>(null)
const toastVisible = ref(false)
let toastTimer: ReturnType<typeof setTimeout> | null = null

async function confirmReset(userId: number) {
  resettingUserId.value = userId
  try {
    await resetReworkRate(userId)
    confirmingUserId.value = null
    toastVisible.value = true
    if (toastTimer) clearTimeout(toastTimer)
    toastTimer = setTimeout(() => { toastVisible.value = false }, 3500)
    emit('refresh')
  } catch {
    confirmingUserId.value = null
  } finally {
    resettingUserId.value = null
  }
}

function roleBadgeClass(role: string) {
  if (role === 'CEO') return 'bg-purple-600/30 text-purple-300'
  if (role === 'PM') return 'bg-blue-600/30 text-blue-300'
  return 'bg-gray-600/30 text-gray-300'
}

function pctColor(pct: number) {
  if (pct >= 85) return 'text-emerald-400'
  if (pct >= 70) return 'text-amber-400'
  return 'text-red-400'
}

function reworkColor(pct: number) {
  if (pct <= 15) return 'text-emerald-400'
  if (pct <= 25) return 'text-amber-400'
  return 'text-red-400'
}

function scoreColor(score: number) {
  if (score >= 80) return 'text-emerald-400'
  if (score >= 60) return 'text-amber-400'
  return 'text-gray-400'
}
</script>
