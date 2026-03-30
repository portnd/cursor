<template>
  <section class="space-y-6">
    <div>
      <h2 class="section-label text-blue-400/90">Performance</h2>
      <p class="text-xs text-gray-500 -mt-2 mb-4 max-w-3xl leading-relaxed">
        Portfolio metrics use projects on this dashboard. Developer KPIs are scoped to tasks
        <span class="text-gray-400">assigned by you</span> (same logic as the engineering leaderboard).
      </p>
    </div>

    <div v-if="performanceStore.error" class="rounded-xl border border-red-500/30 bg-red-900/15 px-4 py-3 text-sm text-red-300">
      {{ performanceStore.error }}
    </div>

    <div v-else-if="performanceStore.loading && !performanceStore.personal" class="text-sm text-gray-500 py-4">
      Loading performance…
    </div>

    <template v-else>
      <!-- Command metrics: portfolio + squad roll-up -->
      <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
        <div class="rounded-2xl border border-blue-500/20 bg-blue-950/20 p-4">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500 mb-1.5">Health pulse</p>
          <p
            v-if="performanceStore.personal"
            class="text-xl font-extrabold tabular-nums"
            :class="healthColor(performanceStore.personal.health_score)"
          >
            {{ performanceStore.personal.health_score.toFixed(0) }}
          </p>
          <p class="text-xs text-gray-500 mt-0.5">Org wellness index</p>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500 mb-1.5">Portfolio done</p>
          <p
            class="text-xl font-extrabold tabular-nums"
            :class="portfolioCompletionPct === null ? 'text-gray-500' : pctColor(portfolioCompletionPct)"
          >
            {{ portfolioCompletionPct === null ? '—' : portfolioCompletionPct.toFixed(1) + '%' }}
          </p>
          <p class="text-xs text-gray-500 mt-0.5">Tasks completed / total</p>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500 mb-1.5">Overdue</p>
          <p
            class="text-xl font-extrabold tabular-nums"
            :class="portfolioOverdue > 0 ? 'text-red-400' : 'text-gray-400'"
          >
            {{ portfolioOverdue }}
          </p>
          <p class="text-xs text-gray-500 mt-0.5">Across visible projects</p>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500 mb-1.5">Active projects</p>
          <p class="text-xl font-extrabold tabular-nums text-emerald-400">{{ activeProjectCount }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Status ACTIVE</p>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gradient-to-br from-blue-900/25 to-indigo-900/20 p-4">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500 mb-1.5">Squad delivery index</p>
          <p
            class="text-xl font-black tabular-nums"
            :class="squadIndex === null ? 'text-gray-500' : scoreColor(squadIndex)"
          >
            {{ squadIndex === null ? '—' : squadIndex.toFixed(1) }}
          </p>
          <p class="text-xs text-gray-500 mt-0.5">Avg composite · devs you assign</p>
        </div>
      </div>

      <PerformanceTeamLeaderboard
        :members="performanceStore.team"
        title="Developers you assign"
        description="Delivery, quality, rework & velocity — only work assigned by you"
        empty-message="No developers yet with tasks assigned by you. Assign work to start tracking squad KPIs."
        @refresh="performanceStore.fetchTeam()"
      />
    </template>
  </section>
</template>

<script setup lang="ts">
import type { Project } from '~/core/modules/projects/infrastructure/projects-api'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'

const props = defineProps<{
  projects: Project[]
}>()

const performanceStore = usePerformanceStore()

const portfolioCompletionPct = computed(() => {
  let completed = 0
  let total = 0
  for (const p of props.projects) {
    completed += p.task_completed ?? 0
    total += p.task_total ?? 0
  }
  if (total === 0) return null
  return (completed / total) * 100
})

const portfolioOverdue = computed(() =>
  props.projects.reduce((s, p) => s + (p.task_overdue ?? 0), 0),
)

const activeProjectCount = computed(() =>
  props.projects.filter(p => p.status === 'ACTIVE').length,
)

const squadIndex = computed(() => {
  const members = performanceStore.team
  if (!members.length) return null
  const sum = members.reduce((acc, m) => acc + m.composite_score, 0)
  return sum / members.length
})

function pctColor(pct: number): string {
  if (pct >= 85) return 'text-emerald-400'
  if (pct >= 70) return 'text-amber-400'
  return 'text-red-400'
}

function scoreColor(score: number): string {
  if (score >= 80) return 'text-emerald-400'
  if (score >= 60) return 'text-amber-400'
  return 'text-gray-400'
}

function healthColor(h: number): string {
  if (h >= 80) return 'text-emerald-400'
  if (h >= 60) return 'text-amber-400'
  return 'text-red-400'
}

onMounted(() => {
  performanceStore.fetchAll('PM')
})
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest mb-4;
}
</style>
