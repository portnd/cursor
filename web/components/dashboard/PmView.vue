<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Sticky Page Header ──────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-blue-500/15 border border-blue-500/30 flex items-center justify-center">
            <svg class="w-4 h-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold text-white">Squad Command Center</h1>
            <p class="text-xs text-gray-500">
              <template v-if="myTeam">
                {{ myTeam.name }} ·
              </template>
              Subsidiary Profit Center
            </p>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <!-- Team badge -->
          <div v-if="myTeam" class="hidden sm:flex items-center gap-3 rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-2">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Squad</span>
            <span class="text-sm font-bold text-blue-300">{{ myTeam.name }}</span>
            <span class="h-3.5 w-px bg-gray-700"/>
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Members</span>
            <span class="text-sm font-bold text-white tabular-nums">{{ myTeamMembers.length }}</span>
          </div>
          <button
            @click="bootstrap"
            :disabled="isBootstrapping"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50"
          >
            <svg class="h-3.5 w-3.5" :class="isBootstrapping ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- ── Bootstrapping ───────────────────────────────────────────────────────── -->
    <div v-if="isBootstrapping" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-blue-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading squad data…</p>
    </div>

    <!-- ── Error ─────────────────────────────────────────────────────────────── -->
    <div v-else-if="bootstrapError" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Failed to load squad data</p>
          <p class="text-xs text-red-300 mt-0.5">{{ bootstrapError }}</p>
          <button @click="bootstrap" class="mt-2 text-xs underline hover:text-red-200 transition-colors">Try again</button>
        </div>
      </div>
    </div>

    <!-- ── No Team Warning ────────────────────────────────────────────────────── -->
    <div v-else-if="!myTeam" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-amber-500/30 bg-amber-900/20 px-5 py-4 text-amber-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Not assigned to a team</p>
          <p class="text-xs text-amber-300 mt-0.5">Ask your CEO to assign you to a squad team to enable the full dashboard.</p>
        </div>
      </div>
    </div>

    <!-- ── Main Dashboard ─────────────────────────────────────────────────────── -->
    <main v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">

      <!-- Continuous UAT Queue (highest priority: sub-task test approvals) -->
      <ContinuousUATQueue />

      <!-- Handover Approval Queue -->
      <ApprovalQueueBoard />

      <!-- Row 1: Squad P&L (includes Total Capital) -->
      <TeamFinancialWidget
        :team-member-ids="myTeamMemberIds"
        :completed-task-ids="[]"
        :total-capital="totalCapital"
        :project-capital-count="Object.keys(projectCapitals).length"
        :loaded-monthly-burn-rate="loadedTeamMonthlyCost"
        :initial-tasks="allTasks"
      />

      <!-- Row 2: B2B Board + Capacity Radar (2-col) -->
      <div class="grid gap-6 lg:grid-cols-2">
        <InternalB2BBoard
          :team-member-ids="myTeamMemberIds"
          :team-project-ids="myTeamProjectIds"
          :initial-tasks="allTasks"
        />
        <TeamCapacityRadar
          :team-members="myTeamMembers"
          :initial-tasks="allTasks"
        />
      </div>

      <!-- Row 3: Active Sprints (full width) -->
      <ActiveSprintsOverview />

      <!-- Row 4: Feature Roadmap Board (PM/CEO only) -->
      <section>
        <FeatureRoadmapBoard />
      </section>

      <!-- Row 5: Project Capital (Internal VC) -->
      <section v-if="teamProjects.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-xs font-semibold uppercase tracking-widest text-gray-500">Project Capital</h2>
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <span>Total:</span>
            <span class="font-bold text-emerald-400">{{ formatMoney(totalCapital) }}</span>
          </div>
        </div>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <NuxtLink
            v-for="project in teamProjects"
            :key="project.id"
            :to="`/projects/${project.code || project.id}?tab=capital`"
            class="group rounded-2xl border border-gray-700/60 bg-gray-800/50 p-4 hover:border-indigo-500/40 hover:bg-gray-800/80 transition-all"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-bold text-white truncate group-hover:text-indigo-300 transition-colors">{{ project.name }}</p>
                <p class="text-xs text-gray-500 font-mono">{{ project.code }}</p>
              </div>
              <svg class="w-4 h-4 text-gray-600 group-hover:text-indigo-400 transition-colors ml-2 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </div>

            <template v-if="projectCapitals[project.id]">
              <!-- Capital Balance -->
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-500">Capital Balance</span>
                <span class="text-sm font-bold text-emerald-400">{{ formatMoney(projectCapitals[project.id].capital_balance) }}</span>
              </div>
              <!-- Burn Rate -->
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-500">Monthly Burn</span>
                <span class="text-xs text-gray-300">{{ formatMoney(projectCapitals[project.id].team_monthly_cost) }}/mo</span>
              </div>
              <!-- Gross Margin -->
              <div class="flex items-center justify-between mb-3">
                <span class="text-xs text-gray-500">Gross Margin</span>
                <span
                  class="text-sm font-bold"
                  :class="(projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) >= 0 ? 'text-emerald-400' : 'text-red-400'"
                >
                  {{ (projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) >= 0 ? '+' : '' }}{{ formatMoney(projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) }}
                </span>
              </div>
              <!-- Runway bar -->
              <div>
                <div class="flex items-center justify-between text-xs mb-1">
                  <span class="text-gray-500">Runway</span>
                  <span class="font-bold" :class="runwayColor(projectCapitals[project.id].runway_months)">
                    {{ projectCapitals[project.id].runway_months > 0 ? projectCapitals[project.id].runway_months.toFixed(1) + ' mo' : '—' }}
                  </span>
                </div>
                <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="runwayBarColor(projectCapitals[project.id].runway_months)"
                    :style="{ width: Math.min((projectCapitals[project.id].runway_months / 3) * 100, 100) + '%' }"
                  ></div>
                </div>
              </div>
            </template>
            <template v-else>
              <div class="text-xs text-gray-600 italic">No capital injected</div>
            </template>
          </NuxtLink>
        </div>
      </section>

      <!-- Row 6: Estimation Delta Chart (full width) -->
      <EstimationDeltaChart
        :team-project-ids="myTeamProjectIds"
        :initial-tasks="allTasks"
      />

    </main>

  </div>
</template>

<script setup lang="ts">
import TeamFinancialWidget from '~/components/dashboard/TeamFinancialWidget.vue'
import InternalB2BBoard from '~/components/dashboard/InternalB2BBoard.vue'
import TeamCapacityRadar from '~/components/dashboard/TeamCapacityRadar.vue'
import ActiveSprintsOverview from '~/components/dashboard/ActiveSprintsOverview.vue'
import EstimationDeltaChart from '~/components/dashboard/EstimationDeltaChart.vue'
import FeatureRoadmapBoard from '~/components/tasks/FeatureRoadmapBoard.vue'
import ApprovalQueueBoard from '~/components/dashboard/ApprovalQueueBoard.vue'
import ContinuousUATQueue from '~/components/dashboard/ContinuousUATQueue.vue'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { Team, TeamUser } from '~/core/modules/teams/infrastructure/teams-api'
import type { Project, ProjectCapitalResponse } from '~/core/modules/projects/infrastructure/projects-api'

const { currentUser } = useAuth()
const { getTeams } = useTeamsApi()
const { getProjects, getProjectCapital } = useProjectsApi()
const { fetchWithAuth } = useAuth()

const isBootstrapping = ref(true)
const bootstrapError = ref('')
const myTeam = ref<Team | null>(null)
const teamProjects = ref<Project[]>([])
const projectCapitals = ref<Record<string, ProjectCapitalResponse>>({})
const allTasks = ref<any[]>([])
const myTeamMembers = computed<TeamUser[]>(() => myTeam.value?.users ?? [])
const myTeamMemberIds = computed<number[]>(() => myTeamMembers.value.map(u => u.id))
const myTeamProjectIds = computed<string[]>(() => teamProjects.value.map(p => p.id))

const totalCapital = computed(() =>
  Object.values(projectCapitals.value).reduce((s, c) => s + (c.capital_balance ?? 0), 0)
)

/** Loaded monthly burn rate for this squad — sum across all projects (team_monthly_cost × number of projects). */
const loadedTeamMonthlyCost = computed(() => {
  const vals = Object.values(projectCapitals.value)
  if (!vals.length) return 0
  return vals.reduce((sum, c) => sum + c.team_monthly_cost, 0)
})

function runwayColor(months: number) {
  if (months > 2) return 'text-emerald-400'
  if (months > 1) return 'text-yellow-400'
  return 'text-red-400'
}
function runwayBarColor(months: number) {
  if (months > 2) return 'bg-emerald-500'
  if (months > 1) return 'bg-yellow-500'
  return 'bg-red-500'
}
function formatMoney(v: number) {
  return '฿' + Math.round(v).toLocaleString('th-TH')
}

const bootstrap = async () => {
  isBootstrapping.value = true
  bootstrapError.value = ''
  try {
    const [teams, projects, tasksRes] = await Promise.all([
      getTeams(),
      getProjects(),
      fetchWithAuth<{ data: any[] }>('/sentinel/tasks').catch(() => ({ data: [] })),
    ])

    allTasks.value = (tasksRes as any)?.data ?? []

    const userId = currentUser.value?.user_id
    const found = teams.find(t => t.users?.some(u => u.id === userId))
    myTeam.value = found ?? null

    teamProjects.value = projects

    // Fetch capital for each project in parallel
    const capitals = await Promise.allSettled(
      projects.map(p => getProjectCapital(p.id).then(c => ({ id: p.id, capital: c })))
    )
    const map: Record<string, ProjectCapitalResponse> = {}
    for (const r of capitals) {
      if (r.status === 'fulfilled') map[r.value.id] = r.value.capital
    }
    projectCapitals.value = map
  } catch (err: any) {
    bootstrapError.value = err?.data?.message || err?.message || 'Failed to load data'
  } finally {
    isBootstrapping.value = false
  }
}

onMounted(() => bootstrap())
</script>
