<template>
  <section>
    <h2 class="section-label">Team Capacity &amp; Pulse</h2>
    <div class="rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden shadow-lg">

      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-700 px-5 py-3.5">
        <div class="flex items-center gap-2.5">
          <div class="w-6 h-6 rounded-md bg-blue-500/15 border border-blue-500/30 flex items-center justify-center">
            <svg class="w-3 h-3 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </div>
          <h3 class="text-sm font-bold text-white">Member Workload</h3>
        </div>
        <span class="text-xs font-medium px-2.5 py-0.5 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400">
          {{ teamMembers.length }} members
        </span>
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="flex items-center justify-center py-16">
        <svg class="h-6 w-6 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
      </div>

      <!-- Empty -->
      <div v-else-if="teamMembers.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
        <svg class="h-8 w-8 text-gray-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
        <p class="text-sm text-gray-400 font-medium">No team members found</p>
      </div>

      <!-- Member List -->
      <ul v-else class="divide-y divide-gray-700/60">
        <li
          v-for="member in memberStats"
          :key="member.id"
          class="px-5 py-4 hover:bg-gray-700/10 transition-colors"
        >
          <div class="flex items-start gap-3">
            <!-- Avatar -->
            <div
              class="w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold text-white flex-shrink-0 shadow-md"
              :class="avatarGradient(member.id)"
            >
              {{ member.initials }}
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2 mb-1">
                <div>
                  <p class="text-sm font-semibold text-white truncate">{{ member.display_name || member.email }}</p>
                  <p class="text-xs text-gray-500">{{ member.role }}</p>
                </div>
                <div class="text-right flex-shrink-0">
                  <span class="text-xs font-bold tabular-nums" :class="workloadColor(member.workloadPct)">
                    {{ member.workloadPct.toFixed(0) }}%
                  </span>
                  <p class="text-xs text-gray-600">{{ formatMinutes(member.totalActiveMins) }} / {{ formatMinutes(member.capacityMins) }}</p>
                </div>
              </div>

              <!-- Workload Bar -->
              <div class="h-1.5 w-full rounded-full bg-gray-700 overflow-hidden mb-2">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="workloadBarColor(member.workloadPct)"
                  :style="{ width: `${Math.min(member.workloadPct, 100)}%` }"
                />
              </div>

              <!-- Active Task Count -->
              <div class="flex items-center gap-3 flex-wrap">
                <span class="text-xs text-gray-500">
                  <span class="font-medium text-gray-300">{{ member.activeTaskCount }}</span> active task{{ member.activeTaskCount !== 1 ? 's' : '' }}
                </span>

                <!-- Standup Status -->
                <template v-if="member.standup">
                  <span class="h-3 w-px bg-gray-700"/>
                  <span class="text-xs text-emerald-400 flex items-center gap-1">
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                    </svg>
                    Checked in
                  </span>
                  <template v-if="member.standup.blocker">
                    <span class="h-3 w-px bg-gray-700"/>
                    <span class="text-xs font-bold text-red-400 flex items-center gap-1">
                      <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                      </svg>
                      BLOCKER: {{ member.standup.blocker }}
                    </span>
                  </template>
                </template>
                <template v-else>
                  <span class="h-3 w-px bg-gray-700"/>
                  <span class="text-xs text-gray-600">No standup today</span>
                </template>
              </div>
            </div>
          </div>
        </li>
      </ul>

    </div>
  </section>
</template>

<script setup lang="ts">
import { usePricingApi } from '~/core/modules/pricing/infrastructure/pricing-api'
import type { UserPulse } from '~/core/modules/pulse/infrastructure/pulse-api'

interface TeamMember {
  id: number
  email: string
  display_name?: string
  role: string
  team_id?: number | null
}

interface MemberTask {
  id: string
  status: string
  estimated_minutes: number
  assigned_to?: number
}

const props = defineProps<{
  teamMembers: TeamMember[]
  initialTasks?: any[]
}>()

const { fetchWithAuth } = useAuth()
const { getCostConfig } = usePricingApi()

const isLoading = ref(true)
const tasks = ref<MemberTask[]>([])
const pulseMembers = ref<UserPulse[]>([])
const availableMinutesPerMonth = ref(9600) // default: 8h * 20d * 60min

const avatarColors = [
  'from-blue-600 to-cyan-600',
  'from-purple-600 to-pink-600',
  'from-emerald-600 to-teal-600',
  'from-amber-600 to-orange-600',
  'from-rose-600 to-red-600',
  'from-indigo-600 to-blue-600',
]

const avatarGradient = (id: number) => {
  return `bg-gradient-to-br ${avatarColors[id % avatarColors.length]}`
}

const workloadColor = (pct: number) => {
  if (pct >= 90) return 'text-red-400'
  if (pct >= 70) return 'text-amber-400'
  return 'text-emerald-400'
}

const workloadBarColor = (pct: number) => {
  if (pct >= 90) return 'bg-red-500'
  if (pct >= 70) return 'bg-amber-500'
  return 'bg-emerald-500'
}

const formatMinutes = (mins: number) => {
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h === 0) return `${m}m`
  if (m === 0) return `${h}h`
  return `${h}h ${m}m`
}

const memberStats = computed(() => {
  return props.teamMembers.map(member => {
    const activeTasks = tasks.value.filter(
      t => t.assigned_to === member.id && ['IN_PROGRESS', 'PENDING', 'REVIEW_PENDING'].includes(t.status)
    )
    const totalActiveMins = activeTasks.reduce((s, t) => s + (t.estimated_minutes || 0), 0)
    const workloadPct = availableMinutesPerMonth.value > 0
      ? (totalActiveMins / availableMinutesPerMonth.value) * 100
      : 0
    const pulse = pulseMembers.value.find(p => p.user_id === member.id)
    const nameParts = (member.display_name || member.email || '??').split(' ')
    const initials = nameParts.length >= 2
      ? `${nameParts[0][0]}${nameParts[1][0]}`.toUpperCase()
      : (member.display_name || member.email || '?').substring(0, 2).toUpperCase()

    return {
      ...member,
      initials,
      activeTaskCount: activeTasks.length,
      totalActiveMins,
      capacityMins: availableMinutesPerMonth.value,
      workloadPct,
      standup: pulse?.standup ?? null,
    }
  })
})

const fetchData = async () => {
  isLoading.value = true
  try {
    const today = new Date().toISOString().split('T')[0]
    const [tasksRes, pulseRes, config] = await Promise.all([
      props.initialTasks && props.initialTasks.length > 0
        ? Promise.resolve({ data: props.initialTasks })
        : fetchWithAuth<{ data: MemberTask[] }>('/sentinel/tasks'),
      fetchWithAuth<{ date: string; members: UserPulse[] }>(`/pulse/daily?date=${today}`).catch(() => null),
      getCostConfig().catch(() => null),
    ])
    tasks.value = tasksRes?.data ?? []
    pulseMembers.value = pulseRes?.members ?? []
    if (config) {
      availableMinutesPerMonth.value = config.working_hours_per_day * config.working_days_per_month * 60
    }
  } catch {
    // silent
  } finally {
    isLoading.value = false
  }
}

onMounted(() => fetchData())
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
</style>
