<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-white">Team Performance</h1>
      <p class="text-sm text-gray-400 mt-1">Company-wide evaluation and KPIs</p>
    </div>

    <div v-if="store.loading && !store.overview" class="text-center py-20 text-gray-400">
      Loading performance data...
    </div>
    <div v-else-if="store.error" class="rounded-lg border border-red-500/50 bg-red-900/20 p-4 text-red-400">
      {{ store.error }}
    </div>

    <!-- CEO only: Engineering Health + KPI cards + Team Leaderboard -->
    <template v-else>
      <div v-if="store.overview" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        <div class="rounded-lg border-2 border-purple-500/50 bg-gray-800 p-4 col-span-full md:col-span-1">
          <div class="text-xs uppercase tracking-wide text-gray-400 mb-1">Engineering Health Index</div>
          <div class="text-4xl font-bold text-white">
            {{ store.overview.engineering_health_index.toFixed(1) }}
          </div>
          <div class="text-xs text-gray-500 mt-1">Weighted composite (0–100)</div>
        </div>
        <PerformanceKpiScoreCard
          label="Sprint Success Rate"
          :value="store.overview.sprint_success_rate_pct"
          format="pct"
          :status="statusForDelivery(store.overview.sprint_success_rate_pct)"
        />
        <PerformanceKpiScoreCard
          label="Project On-Track"
          :value="store.overview.project_on_track_rate_pct"
          format="pct"
          :status="statusForDelivery(store.overview.project_on_track_rate_pct)"
        />
        <PerformanceKpiScoreCard
          label="Milestone Hit Rate"
          :value="store.overview.milestone_hit_rate_pct"
          format="pct"
          :status="statusForDelivery(store.overview.milestone_hit_rate_pct)"
        />
        <PerformanceKpiScoreCard
          label="Cursor Adoption"
          :value="store.overview.cursor_adoption_score"
          sublabel="AI assistance level"
          :status="store.overview.cursor_adoption_score >= 80 ? 'good' : store.overview.cursor_adoption_score >= 50 ? 'warn' : 'bad'"
        />
        <PerformanceKpiScoreCard
          label="Velocity Trend"
          :value="store.overview.team_velocity_trend_pct"
          format="pct"
          sublabel="Sprint-over-sprint growth"
          :status="store.overview.team_velocity_trend_pct >= 0 ? 'good' : 'warn'"
        />
      </div>
      <PerformanceTeamLeaderboard :members="store.team" @refresh="store.fetchTeam()" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { usePerformanceStore } from '~/core/modules/performance/performance-store'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const { currentUser } = useAuth()
const store = usePerformanceStore()

const role = computed(() => (currentUser.value?.role as string)?.toUpperCase() || 'DEV')

onMounted(() => {
  if (role.value !== 'CEO') {
    navigateTo('/dashboard')
    return
  }
  store.fetchAll('CEO')
})

function statusForDelivery(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct >= 85) return 'good'
  if (pct >= 70) return 'warn'
  return 'bad'
}

function statusForQuality(q: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (q >= 75) return 'good'
  if (q >= 60) return 'warn'
  return 'bad'
}

function statusForRework(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct <= 15) return 'good'
  if (pct <= 25) return 'warn'
  return 'bad'
}

function statusForAccuracy(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct >= 70) return 'good'
  if (pct >= 50) return 'warn'
  return 'bad'
}

function statusForHealth(h: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (h >= 80) return 'good'
  if (h >= 60) return 'warn'
  return 'bad'
}
</script>
