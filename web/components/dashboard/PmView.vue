<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Sticky Page Header ──────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-3">
          <!-- Icon: blue for Squad mode, violet for Portfolio mode -->
          <div
            class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors"
            :class="teamsEnabled ? 'bg-blue-500/15 border border-blue-500/30' : 'bg-violet-500/15 border border-violet-500/30'"
          >
            <!-- Squad icon -->
            <svg v-if="teamsEnabled" class="w-4 h-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
            <!-- Portfolio / project icon -->
            <svg v-else class="w-4 h-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold text-white">
              {{ teamsEnabled ? 'Squad Command Center' : 'Project Command Center' }}
            </h1>
            <p class="text-xs text-gray-500">
              <template v-if="teamsEnabled && myTeam">{{ myTeam.name }} · </template>
              {{ teamsEnabled ? 'Subsidiary Profit Center' : 'Portfolio Management' }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <!-- Squad badge (teams enabled + team found) -->
          <div v-if="teamsEnabled && myTeam" class="hidden sm:flex items-center gap-3 rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-2">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Squad</span>
            <span class="text-sm font-bold text-blue-300">{{ myTeam.name }}</span>
            <span class="h-3.5 w-px bg-gray-700"/>
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Members</span>
            <span class="text-sm font-bold text-white tabular-nums">{{ myTeamMembers.length }}</span>
          </div>
          <!-- Portfolio mode badge (teams disabled) -->
          <div v-if="!teamsEnabled" class="hidden sm:flex items-center gap-2 rounded-xl border border-violet-500/30 bg-violet-900/20 px-3 py-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-violet-400"/>
            <span class="text-xs font-medium text-violet-300">Individual Mode</span>
            <span class="text-xs text-violet-500">· {{ teamProjects.length }} projects</span>
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

    <!-- ── Teams-Disabled Banner ───────────────────────────────────────────────── -->
    <div v-if="!teamsEnabled && !isBootstrapping" class="max-w-screen-xl mx-auto px-6 pt-5">
      <div class="flex items-center gap-3 rounded-xl border border-violet-500/20 bg-violet-900/10 px-5 py-3">
        <svg class="h-4 w-4 flex-shrink-0 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p class="text-xs text-violet-300">
          <span class="font-semibold">Running in Individual Project Mode</span>
          <span class="text-violet-400/70 ml-1.5">— Squad P&amp;L, capacity radar and B2B boards are hidden. All project-level features remain active.</span>
        </p>
        <NuxtLink to="/admin/teams" class="ml-auto text-xs font-medium text-violet-400 hover:text-violet-300 transition-colors whitespace-nowrap underline underline-offset-2">
          Manage Teams
        </NuxtLink>
      </div>
    </div>

    <!-- ── Bootstrapping ───────────────────────────────────────────────────────── -->
    <div v-if="isBootstrapping" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-blue-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading project data…</p>
    </div>

    <!-- ── Error ─────────────────────────────────────────────────────────────── -->
    <div v-else-if="bootstrapError" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Failed to load data</p>
          <p class="text-xs text-red-300 mt-0.5">{{ bootstrapError }}</p>
          <button @click="bootstrap" class="mt-2 text-xs underline hover:text-red-200 transition-colors">Try again</button>
        </div>
      </div>
    </div>

    <!-- ── No Team Warning (teams enabled but not assigned) ───────────────────── -->
    <div v-else-if="teamsEnabled && !myTeam" class="max-w-screen-xl mx-auto px-6 py-8">
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

      <PmPerformanceSection :projects="teamProjects" />

      <!-- ════════════════════════════════════════════════════════════════════════
           TEAM MODE — Squad financial & capacity widgets
           ════════════════════════════════════════════════════════════════════════ -->
      <template v-if="teamsEnabled && myTeam">
        <!-- Squad P&L (includes Total Capital) -->
        <TeamFinancialWidget
          :team-member-ids="myTeamMemberIds"
          :completed-task-ids="[]"
          :total-capital="totalCapital"
          :project-capital-count="Object.keys(projectCapitals).length"
          :loaded-monthly-burn-rate="loadedTeamMonthlyCost"
          :initial-tasks="allTasks"
        />

        <!-- B2B Board + Capacity Radar (2-col) -->
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
      </template>

      <!-- ════════════════════════════════════════════════════════════════════════
           INDIVIDUAL MODE — Portfolio overview (teams feature disabled)
           ════════════════════════════════════════════════════════════════════════ -->
      <template v-else-if="!teamsEnabled">

        <!-- Portfolio Summary Stats (no aggregate Capital when feature team is off — squad capital is N/A) -->
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3">
          <div class="rounded-2xl border border-gray-700/60 bg-gray-800/50 p-5">
            <p class="text-xs font-medium uppercase tracking-widest text-gray-500 mb-2">Projects</p>
            <p class="text-3xl font-bold text-white tabular-nums">{{ teamProjects.length }}</p>
            <p class="text-xs text-gray-600 mt-1">total portfolio</p>
          </div>
          <div class="rounded-2xl border border-gray-700/60 bg-gray-800/50 p-5">
            <p class="text-xs font-medium uppercase tracking-widest text-gray-500 mb-2">Active</p>
            <p class="text-3xl font-bold text-emerald-400 tabular-nums">{{ portfolioStats.activeCount }}</p>
            <p class="text-xs text-gray-600 mt-1">running now</p>
          </div>
          <div class="rounded-2xl border border-gray-700/60 bg-gray-800/50 p-5">
            <p class="text-xs font-medium uppercase tracking-widest text-gray-500 mb-2">Overdue</p>
            <p class="text-3xl font-bold tabular-nums" :class="portfolioStats.totalOverdue > 0 ? 'text-red-400' : 'text-gray-500'">
              {{ portfolioStats.totalOverdue }}
            </p>
            <p class="text-xs text-gray-600 mt-1">tasks past due</p>
          </div>
        </div>

        <!-- Project Portfolio Health Grid -->
        <section>
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xs font-semibold uppercase tracking-widest text-gray-500">Project Portfolio</h2>
            <NuxtLink to="/projects" class="text-xs text-gray-500 hover:text-gray-300 transition-colors flex items-center gap-1">
              View all
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </NuxtLink>
          </div>
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <NuxtLink
              v-for="project in teamProjects"
              :key="project.id"
              :to="`/projects/${project.code || project.id}`"
              class="group rounded-2xl border border-gray-700/60 bg-gray-800/50 p-5 hover:border-violet-500/40 hover:bg-gray-800/80 transition-all"
            >
              <!-- Project header -->
              <div class="flex items-start justify-between mb-4 gap-2">
                <div class="flex items-center gap-2 min-w-0 flex-1">
                  <span
                    v-if="project.color"
                    class="w-2.5 h-2.5 rounded-full flex-shrink-0"
                    :style="{ backgroundColor: project.color }"
                  />
                  <div class="min-w-0">
                    <p class="text-sm font-bold text-white truncate group-hover:text-violet-300 transition-colors">{{ project.name }}</p>
                    <p class="text-xs font-mono text-gray-500">{{ project.code }}</p>
                  </div>
                </div>
                <span class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-semibold" :class="projectStatusClass(project.status)">
                  {{ project.status }}
                </span>
              </div>

              <!-- Task completion progress -->
              <div v-if="(project.task_total ?? 0) > 0" class="mb-4">
                <div class="flex items-center justify-between text-xs mb-1.5">
                  <span class="text-gray-500">Tasks</span>
                  <span class="font-semibold tabular-nums text-gray-300">
                    {{ project.task_completed ?? 0 }}/{{ project.task_total }}
                    <span class="text-gray-500 font-normal ml-1">({{ taskPct(project) }}%)</span>
                  </span>
                </div>
                <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="taskPct(project) >= 80 ? 'bg-emerald-500' : taskPct(project) >= 40 ? 'bg-violet-500' : 'bg-amber-500'"
                    :style="{ width: taskPct(project) + '%' }"
                  />
                </div>
              </div>
              <div v-else class="mb-4">
                <p class="text-xs text-gray-600 italic">No tasks yet</p>
              </div>

              <!-- Overdue only (no capital / runway when feature team is disabled) -->
              <div class="flex items-center justify-end pt-3 border-t border-gray-700/50">
                <div class="flex items-center gap-1.5 text-xs">
                  <template v-if="(project.task_overdue ?? 0) > 0">
                    <span class="w-1.5 h-1.5 rounded-full bg-red-500"/>
                    <span class="font-semibold text-red-400">{{ project.task_overdue }} overdue</span>
                  </template>
                  <template v-else>
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500/50"/>
                    <span class="text-gray-600">On track</span>
                  </template>
                </div>
              </div>
            </NuxtLink>
          </div>
        </section>

      </template>

      <!-- Active Sprints (full width — both modes) — reuses project list from bootstrap (no duplicate getProjects) -->
      <ActiveSprintsOverview :projects="teamProjects" />

      <!-- ── Project Capital Deep-dive (Squad mode only — teams have dedicated capital) -->
      <section v-if="teamsEnabled && teamProjects.length > 0">
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
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-500">Capital Balance</span>
                <span class="text-sm font-bold text-emerald-400">{{ formatMoney(projectCapitals[project.id].capital_balance) }}</span>
              </div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-500">Monthly Burn</span>
                <span class="text-xs text-gray-300">{{ formatMoney(projectCapitals[project.id].team_monthly_cost) }}/mo</span>
              </div>
              <div class="flex items-center justify-between mb-3">
                <span class="text-xs text-gray-500">Gross Margin</span>
                <span
                  class="text-sm font-bold"
                  :class="(projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) >= 0 ? 'text-emerald-400' : 'text-red-400'"
                >
                  {{ (projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) >= 0 ? '+' : '' }}{{ formatMoney(projectCapitals[project.id].capital_balance - projectCapitals[project.id].team_monthly_cost) }}
                </span>
              </div>
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
                  />
                </div>
              </div>
            </template>
            <template v-else>
              <div class="text-xs text-gray-600 italic">No capital injected</div>
            </template>
          </NuxtLink>
        </div>
      </section>

      <!-- Estimation Delta Chart: squad projects when teams on; PM-owned projects when teams off -->
      <EstimationDeltaChart
        :team-project-ids="estimationDeltaProjectIds"
        :initial-tasks="allTasks"
        :section-title="estimationDeltaSectionTitle"
        :scope-description="estimationDeltaScopeDescription"
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
import ContinuousUATQueue from '~/components/dashboard/ContinuousUATQueue.vue'
import PmPerformanceSection from '~/components/dashboard/PmPerformanceSection.vue'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { Team, TeamUser } from '~/core/modules/teams/infrastructure/teams-api'
import type { Project, ProjectCapitalResponse } from '~/core/modules/projects/infrastructure/projects-api'

const { currentUser } = useAuth()
const performanceStore = usePerformanceStore()
const { getTeams } = useTeamsApi()
const { getProjects, getProjectCapital } = useProjectsApi()
const { fetchWithAuth } = useAuth()
const teamsStore = useTeamsStore()

const isBootstrapping = ref(true)
const bootstrapError = ref('')
const myTeam = ref<Team | null>(null)
const teamProjects = ref<Project[]>([])
const projectCapitals = ref<Record<string, ProjectCapitalResponse>>({})
const allTasks = ref<any[]>([])

const teamsEnabled = computed(() => teamsStore.teamsFeatureEnabled)
const myTeamMembers = computed<TeamUser[]>(() => myTeam.value?.users ?? [])
const myTeamMemberIds = computed<number[]>(() => myTeamMembers.value.map(u => u.id))
const myTeamProjectIds = computed<string[]>(() => teamProjects.value.map(p => p.id))

/** Squad mode: all projects on the PM’s team. Individual (no squads): only projects where this PM is in pm_owners (fallback: full list from API if owners not loaded). */
const estimationDeltaProjectIds = computed<string[]>(() => {
  if (teamsEnabled.value) {
    return teamProjects.value.map(p => p.id)
  }
  const uid = (currentUser.value as { user_id?: number; id?: number } | null)?.user_id
    ?? (currentUser.value as { user_id?: number; id?: number } | null)?.id
  if (uid == null) return teamProjects.value.map(p => p.id)
  const owned = teamProjects.value.filter((p) => {
    const owners = p.pm_owners
    if (!owners?.length) return true
    return owners.some(o => o.user_id === uid)
  })
  return owned.map(p => p.id)
})

const estimationDeltaSectionTitle = computed(() =>
  teamsEnabled.value
    ? 'Estimation Accuracy — Delta Analysis'
    : 'Estimation Accuracy — Portfolio delta analysis'
)

const estimationDeltaScopeDescription = computed(() =>
  teamsEnabled.value
    ? 'Tasks in projects linked to your squad.'
    : 'Only tasks in projects where you are the assigned PM (portfolio mode — not squad-scoped).'
)

const totalCapital = computed(() =>
  Object.values(projectCapitals.value).reduce((s, c) => s + (c.capital_balance ?? 0), 0)
)

const loadedTeamMonthlyCost = computed(() => {
  const vals = Object.values(projectCapitals.value)
  if (!vals.length) return 0
  return vals.reduce((sum, c) => sum + c.team_monthly_cost, 0)
})

const portfolioStats = computed(() => ({
  activeCount: teamProjects.value.filter(p => p.status === 'ACTIVE').length,
  totalOverdue: teamProjects.value.reduce((s, p) => s + (p.task_overdue ?? 0), 0),
}))

function taskPct(project: Project): number {
  if (!project.task_total || project.task_total === 0) return 0
  return Math.round(((project.task_completed ?? 0) / project.task_total) * 100)
}

function projectStatusClass(status: string) {
  if (status === 'ACTIVE') return 'bg-emerald-500/15 border border-emerald-500/30 text-emerald-400'
  if (status === 'COMPLETED') return 'bg-blue-500/15 border border-blue-500/30 text-blue-400'
  return 'bg-amber-500/15 border border-amber-500/30 text-amber-400'
}

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
    // Feature flag + project list + tasks in parallel (was serial: flag → then projects/tasks)
    const [, projects, tasksRes] = await Promise.all([
      teamsStore.fetchTeamsFeatureEnabled(),
      getProjects(),
      fetchWithAuth<{ data: any[] }>('/sentinel/tasks').catch(() => ({ data: [] })),
    ])

    allTasks.value = (tasksRes as any)?.data ?? []
    teamProjects.value = projects

    // Squad mode: resolve team + per-project capital in parallel (was: teams then N capitals)
    if (teamsStore.teamsFeatureEnabled) {
      const userId = currentUser.value?.user_id
      const [teams, capitals] = await Promise.all([
        getTeams(),
        Promise.allSettled(
          projects.map(p => getProjectCapital(p.id).then(c => ({ id: p.id, capital: c })))
        ),
      ])
      myTeam.value = teams.find(t => t.users?.some(u => u.id === userId)) ?? null
      const map: Record<string, ProjectCapitalResponse> = {}
      for (const r of capitals) {
        if (r.status === 'fulfilled') map[r.value.id] = r.value.capital
      }
      projectCapitals.value = map
    } else {
      myTeam.value = null
      projectCapitals.value = {}
    }
  } catch (err: any) {
    bootstrapError.value = err?.data?.message || err?.message || 'Failed to load data'
  } finally {
    isBootstrapping.value = false
  }
}

onMounted(() => {
  // KPIs load alongside bootstrap so “Performance” is ready when the shell appears
  performanceStore.fetchAll('PM')
  bootstrap()
})
</script>
