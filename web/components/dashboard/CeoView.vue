<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Sticky Page Header ─────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center">
            <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold text-white">CEO Command Center</h1>
            <p class="text-xs text-gray-500">Strategic metrics & company health at a glance</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs text-gray-600">Updated {{ lastUpdated }}</span>
          <button
            @click="refresh"
            :disabled="isLoading"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-gray-900 dark:text-white transition-colors disabled:opacity-50"
          >
            <svg class="h-3.5 w-3.5" :class="isLoading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- ── Loading ───────────────────────────────────────────────────────────── -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-amber-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading strategic data…</p>
    </div>

    <!-- ── Error ─────────────────────────────────────────────────────────────── -->
    <div v-else-if="error" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Failed to load data</p>
          <p class="text-xs text-red-300 mt-0.5">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- ── Main Content ───────────────────────────────────────────────────────── -->
    <main v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">

      <!-- ── CEO UAT Approval Queue (sub-tasks Product Owner has tested, awaiting CEO final approval) ── -->
      <CeoUATApprovalQueue />

      <!-- ── Section label helper ──────────────────────────────────────────────── -->

      <!-- ── Five headline KPIs ──────────────────────────────────────────────── -->
      <section>
        <h2 class="section-label">Five numbers that matter</h2>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">

          <!-- Cash Runway -->
          <div
            class="rounded-2xl border p-5 shadow-lg"
            :class="financeSummary
              ? (financeSummary.runway_months >= 12
                  ? 'border-emerald-500/30 bg-gradient-to-br from-emerald-900/30 to-emerald-950/20'
                  : financeSummary.runway_months >= 6
                    ? 'border-amber-500/30 bg-gradient-to-br from-amber-900/30 to-amber-950/20'
                    : 'border-red-500/30 bg-gradient-to-br from-red-900/30 to-red-950/20')
              : 'border-gray-700 bg-gray-800/60'"
          >
            <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1">Cash Runway</p>
            <p class="text-2xl font-black text-white tabular-nums">
              {{ financeSummary && financeSummary.runway_months > 0 ? financeSummary.runway_months.toFixed(1) + ' mo' : '—' }}
            </p>
            <p class="text-xs text-gray-500 mt-1.5">
              {{ financeSummary ? 'From accounting' : 'กรอกที่หน้าบัญชี' }}
            </p>
          </div>

          <!-- Engineering Health -->
          <div class="rounded-2xl border p-5 shadow-lg" :class="kpiCardClass(overview?.engineering_health_index ?? 0, 70, 50)">
            <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1">Engineering Health</p>
            <p class="text-2xl font-black text-white tabular-nums">{{ overview ? overview.engineering_health_index.toFixed(1) : '—' }}</p>
            <p class="text-xs text-gray-500 mt-1.5">0–100 composite</p>
          </div>

          <!-- Projects On Track -->
          <div class="rounded-2xl border p-5 shadow-lg" :class="kpiCardClass(overview?.project_on_track_rate_pct ?? 0, 80, 60)">
            <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1">Projects On Track</p>
            <p class="text-2xl font-black text-white tabular-nums">{{ overview ? overview.project_on_track_rate_pct.toFixed(1) : '—' }}%</p>
            <p class="text-xs text-gray-500 mt-1.5">Delivery reliability</p>
          </div>

          <!-- Velocity Trend -->
          <div
            class="rounded-2xl border p-5 shadow-lg"
            :class="overview
              ? (overview.team_velocity_trend_pct >= 0 ? 'border-emerald-500/30 bg-gray-800/60' : 'border-amber-500/30 bg-gray-800/60')
              : 'border-gray-700 bg-gray-800/60'"
          >
            <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1">Velocity Trend</p>
            <p class="text-2xl font-black text-white tabular-nums">
              {{ overview ? (overview.team_velocity_trend_pct >= 0 ? '+' : '') + overview.team_velocity_trend_pct.toFixed(1) : '—' }}%
            </p>
            <p class="text-xs text-gray-500 mt-1.5">Sprint-over-sprint</p>
          </div>

          <!-- Sprint Success -->
          <div class="rounded-2xl border p-5 shadow-lg" :class="kpiCardClass(overview?.sprint_success_rate_pct ?? 0, 80, 60)">
            <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1">Sprint Success</p>
            <p class="text-2xl font-black text-white tabular-nums">{{ overview ? overview.sprint_success_rate_pct.toFixed(1) : '—' }}%</p>
            <p class="text-xs text-gray-500 mt-1.5">Commitments met</p>
          </div>
        </div>
      </section>

      <!-- ── Company Pulse (Finance) ─────────────────────────────────────────── -->
      <section>
        <h2 class="section-label">Company pulse</h2>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="rounded-2xl border border-gray-700 bg-gray-800/60 px-5 py-4">
            <p class="text-xs font-semibold uppercase tracking-widest text-cyan-300 mb-1">MRR / Revenue</p>
            <p class="text-xl font-extrabold text-white tabular-nums">
              {{ financeSummary ? formatFinance(financeSummary.last_month_mrr) + ' ' + financeSummary.currency : '—' }}
            </p>
            <p class="text-xs text-gray-500 mt-1">รายได้เดือนล่าสุด</p>
          </div>
          <div class="rounded-2xl border border-gray-700 bg-gray-800/60 px-5 py-4">
            <p class="text-xs font-semibold uppercase tracking-widest text-red-300 mb-1">Burn Rate</p>
            <p class="text-xl font-extrabold text-white tabular-nums">
              {{ financeSummary ? formatFinance(financeSummary.burn_rate) + ' ' + financeSummary.currency : '—' }}
            </p>
            <p class="text-xs text-gray-500 mt-1">เฉลี่ย/เดือน (12 เดือน)</p>
          </div>
          <div class="rounded-2xl border border-gray-700 bg-gray-800/60 px-5 py-4">
            <p class="text-xs font-semibold uppercase tracking-widest text-emerald-300 mb-1">Net New ARR</p>
            <p class="text-xl font-extrabold text-white tabular-nums">
              {{ financeSummary ? formatFinance(financeSummary.net_new_arr) + ' ' + financeSummary.currency : '—' }}
            </p>
            <p class="text-xs text-gray-500 mt-1">เทียบเดือนก่อน</p>
          </div>
        </div>
      </section>

      <!-- ── Project Capital Overview (only when teams feature is enabled) ───── -->
      <section v-if="teamsStore.teamsFeatureEnabled">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label" style="margin-bottom:0">Project Capital</h2>
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <span>Total deployed:</span>
            <span class="font-bold text-emerald-400">{{ formatMoney(totalProjectCapital) }}</span>
          </div>
        </div>
        <div v-if="isLoadingCapitals" class="text-center py-8 text-gray-600 text-sm">Loading capital data…</div>
        <div v-else-if="projectCapitals.length === 0" class="text-center py-8 text-gray-600 text-sm italic">No projects with capital yet.</div>
        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <NuxtLink
            v-for="c in projectCapitals"
            :key="c.project_id"
            :to="`/projects/${c.project_id}?tab=capital`"
            class="group rounded-2xl border border-gray-700/60 bg-gray-800/50 p-4 hover:border-amber-500/40 hover:bg-gray-800/80 transition-all"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-bold text-white truncate group-hover:text-amber-300 transition-colors">{{ c.project_name }}</p>
              </div>
              <svg class="w-4 h-4 text-gray-600 group-hover:text-amber-400 transition-colors ml-2 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </div>
            <div class="flex items-center justify-between mb-2">
              <span class="text-xs text-gray-500">Capital</span>
              <span class="text-sm font-bold text-emerald-400">{{ formatMoney(c.capital_balance) }}</span>
            </div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs text-gray-500">Monthly Burn</span>
              <span class="text-xs text-gray-300">{{ formatMoney(c.team_monthly_cost) }}/mo</span>
            </div>
            <div>
              <div class="flex items-center justify-between text-xs mb-1">
                <span class="text-gray-500">Runway</span>
                <span class="font-bold" :class="runwayColor(c.runway_months)">
                  {{ c.runway_months > 0 ? c.runway_months.toFixed(1) + ' mo' : '—' }}
                </span>
              </div>
              <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all"
                  :class="runwayBarColor(c.runway_months)"
                  :style="{ width: Math.min((c.runway_months / 3) * 100, 100) + '%' }"
                ></div>
              </div>
            </div>
          </NuxtLink>
        </div>
      </section>

      <!-- ── Quick Links ────────────────────────────────────────────────────── -->
      <section>
        <h2 class="section-label">Quick access</h2>
        <div class="grid grid-cols-2 sm:grid-cols-5 gap-3">
          <NuxtLink
            v-for="link in quickLinks"
            :key="link.to"
            :to="link.to"
            class="group rounded-2xl border border-gray-700/60 bg-gray-800/40 p-4 hover:border-gray-600 hover:bg-gray-800/70 transition-all"
          >
            <div class="w-8 h-8 rounded-lg flex items-center justify-center mb-3" :class="link.iconBg">
              <svg class="w-4 h-4" :class="link.iconColor" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="link.icon"/>
              </svg>
            </div>
            <p class="text-sm font-semibold text-white">{{ link.label }}</p>
            <p class="text-xs text-gray-500 mt-0.5">{{ link.desc }}</p>
          </NuxtLink>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import type { OverviewKPIs } from '~/core/modules/performance/performance-api'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import { useFinanceApi } from '~/core/modules/finance/finance-api'
import type { FinanceSummary } from '~/core/modules/finance/finance-api'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { ProjectCapitalResponse } from '~/core/modules/projects/infrastructure/projects-api'
import CeoUATApprovalQueue from '~/components/dashboard/CeoUATApprovalQueue.vue'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'

const performanceStore = usePerformanceStore()
const teamsStore = useTeamsStore()
const financeApi = useFinanceApi()
const { currentUser } = useAuth()
const { getProjects, getProjectCapital } = useProjectsApi()

const isLoading = ref(true)
const error = ref('')
const lastUpdated = ref('—')
const financeSummary = ref<FinanceSummary | null>(null)
const projectCapitals = ref<ProjectCapitalResponse[]>([])
const isLoadingCapitals = ref(false)

const overview = computed<OverviewKPIs | null>(() => performanceStore.overview)

const totalProjectCapital = computed(() =>
  projectCapitals.value.reduce((s, c) => s + (c.capital_balance ?? 0), 0)
)

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

const quickLinks = [
  { to: '/projects', label: 'Projects', desc: 'All active projects', icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z', iconBg: 'bg-blue-500/15 border border-blue-500/30', iconColor: 'text-blue-400' },
  { to: '/admin/cost-config', label: 'Cost Config', desc: 'Salary & overhead', icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z', iconBg: 'bg-amber-500/15 border border-amber-500/30', iconColor: 'text-amber-400' },
  { to: '/create', label: 'New Task', desc: 'Create & estimate', icon: 'M12 4v16m8-8H4', iconBg: 'bg-purple-500/15 border border-purple-500/30', iconColor: 'text-purple-400' },
  { to: '/pulse', label: 'Team Pulse', desc: 'Daily standups', icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z', iconBg: 'bg-emerald-500/15 border border-emerald-500/30', iconColor: 'text-emerald-400' },
  { to: '/admin/teams', label: 'Squads', desc: 'Manage teams', icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z', iconBg: 'bg-indigo-500/15 border border-indigo-500/30', iconColor: 'text-indigo-400' },
]

function formatFinance(value: number): string {
  return new Intl.NumberFormat('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(value)
}

function kpiCardClass(value: number, good: number, warn: number): string {
  if (value >= good) return 'border-emerald-500/30 bg-gray-800/60'
  if (value >= warn) return 'border-amber-500/30 bg-gray-800/60'
  return 'border-red-500/30 bg-gray-800/60'
}

async function loadProjectCapitals() {
  isLoadingCapitals.value = true
  try {
    const projects = await getProjects()
    const results = await Promise.allSettled(projects.map(p => getProjectCapital(p.id)))
    projectCapitals.value = results
      .filter((r): r is PromiseFulfilledResult<ProjectCapitalResponse> => r.status === 'fulfilled')
      .map(r => r.value)
      .filter(c => c.capital_balance > 0)
      .sort((a, b) => b.capital_balance - a.capital_balance)
  } catch (_) {}
  finally { isLoadingCapitals.value = false }
}

async function refresh() {
  isLoading.value = true
  error.value = ''
  try {
    await Promise.all([
      performanceStore.fetchAll(currentUser.value?.role ?? 'CEO'),
      financeApi.getSummary().then(s => { financeSummary.value = s }),
      teamsStore.fetchTeamsFeatureEnabled(),
      loadProjectCapitals(),
    ])
    lastUpdated.value = new Date().toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  } catch (e: any) {
    error.value = e?.message ?? performanceStore.error ?? 'Failed to load'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  refresh()
})
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
</style>
