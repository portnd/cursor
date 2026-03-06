<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <!-- CEO Command Center Header -->
    <header class="sticky top-0 z-10 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 py-6 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
              CEO Command Center
            </h1>
            <p class="mt-1 text-sm text-gray-400">
              Strategic metrics at a glance — what founders need to see
            </p>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs text-gray-500">
              Updated {{ lastUpdated }}
            </span>
            <button
              type="button"
              @click="refresh"
              class="inline-flex items-center gap-2 rounded-xl border border-gray-600 bg-gray-800 px-4 py-2.5 text-sm font-medium text-gray-200 transition-colors hover:border-gray-500 hover:bg-gray-700 hover:text-white"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center py-24">
      <div class="h-10 w-10 animate-spin rounded-full border-2 border-amber-500 border-t-transparent" />
      <p class="mt-4 text-sm text-gray-500">Loading strategic data...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div class="rounded-xl border border-red-500/50 bg-red-950/20 px-5 py-4 text-red-400">
        <div class="flex items-start gap-3">
          <span class="text-xl">⚠️</span>
          <div>
            <p class="font-medium">Failed to load data</p>
            <p class="mt-1 text-sm text-red-300">{{ error }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content: Single-pane strategic view -->
    <main v-else class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <!-- Five headline numbers (founder-first KPIs) -->
      <section class="mb-8">
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Five numbers that matter
        </h2>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
          <div
            class="rounded-xl border p-5 shadow-lg"
            :class="financeSummary ? (financeSummary.runway_months >= 12 ? 'border-emerald-500/50 bg-gray-800/60' : financeSummary.runway_months >= 6 ? 'border-amber-500/50 bg-gray-800/60' : 'border-red-500/50 bg-gray-800/60') : 'border-gray-700/80 bg-gray-800/60'"
          >
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Cash Runway</div>
            <div class="mt-2 text-2xl font-bold text-white">
              {{ financeSummary && financeSummary.runway_months > 0 ? financeSummary.runway_months.toFixed(1) + ' mo' : '—' }}
            </div>
            <div class="mt-1 text-xs text-gray-500">
              {{ financeSummary ? 'From accounting' : 'กรอกที่หน้าบัญชี' }}
            </div>
          </div>
          <div
            class="rounded-xl border p-5 shadow-lg"
            :class="overviewStatus.engineering"
          >
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Engineering Health</div>
            <div class="mt-2 text-2xl font-bold text-white">
              {{ overview ? overview.engineering_health_index.toFixed(1) : '—' }}
            </div>
            <div class="mt-1 text-xs text-gray-500">0–100 composite</div>
          </div>
          <div
            class="rounded-xl border p-5 shadow-lg"
            :class="overviewStatus.projects"
          >
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Projects On Track</div>
            <div class="mt-2 text-2xl font-bold text-white">
              {{ overview ? overview.project_on_track_rate_pct.toFixed(1) : '—' }}%
            </div>
            <div class="mt-1 text-xs text-gray-500">Delivery reliability</div>
          </div>
          <div
            class="rounded-xl border p-5 shadow-lg"
            :class="overviewStatus.velocity"
          >
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Velocity Trend</div>
            <div class="mt-2 text-2xl font-bold text-white">
              {{ overview ? (overview.team_velocity_trend_pct >= 0 ? '+' : '') + overview.team_velocity_trend_pct.toFixed(1) : '—' }}%
            </div>
            <div class="mt-1 text-xs text-gray-500">Sprint-over-sprint</div>
          </div>
          <div
            class="rounded-xl border p-5 shadow-lg"
            :class="overviewStatus.sprint"
          >
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Sprint Success</div>
            <div class="mt-2 text-2xl font-bold text-white">
              {{ overview ? overview.sprint_success_rate_pct.toFixed(1) : '—' }}%
            </div>
            <div class="mt-1 text-xs text-gray-500">Commitments met</div>
          </div>
        </div>
      </section>

      <!-- Company pulse (from accounting) -->
      <section class="mb-8">
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Company pulse
        </h2>
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div class="rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-3">
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">MRR / Revenue</div>
            <div class="mt-1 text-lg font-semibold text-white">
              {{ financeSummary ? formatFinance(financeSummary.last_month_mrr) + ' ' + financeSummary.currency : '—' }}
            </div>
            <div class="mt-0.5 text-xs text-gray-500">รายได้เดือนล่าสุด</div>
          </div>
          <div class="rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-3">
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Burn rate</div>
            <div class="mt-1 text-lg font-semibold text-white">
              {{ financeSummary ? formatFinance(financeSummary.burn_rate) + ' ' + financeSummary.currency : '—' }}
            </div>
            <div class="mt-0.5 text-xs text-gray-500">เฉลี่ย/เดือน (12 เดือน)</div>
          </div>
          <div class="rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-3">
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500">Net new ARR</div>
            <div class="mt-1 text-lg font-semibold text-white">
              {{ financeSummary ? formatFinance(financeSummary.net_new_arr) + ' ' + financeSummary.currency : '—' }}
            </div>
            <div class="mt-0.5 text-xs text-gray-500">เทียบเดือนก่อน</div>
          </div>
        </div>
      </section>

      <!-- Delivery & product health -->
      <section class="mb-8">
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Delivery & product health
        </h2>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
          <KpiScoreCard
            label="Engineering health"
            :value="overview?.engineering_health_index ?? 0"
            format="number"
            sublabel="Composite 0–100"
            :status="rag(overview?.engineering_health_index ?? 0, 70, 50)"
          />
          <KpiScoreCard
            label="Sprint success rate"
            :value="overview?.sprint_success_rate_pct ?? 0"
            format="pct"
            sublabel="Commitments met"
            :status="rag(overview?.sprint_success_rate_pct ?? 0, 80, 60)"
          />
          <KpiScoreCard
            label="Milestone hit rate"
            :value="overview?.milestone_hit_rate_pct ?? 0"
            format="pct"
            sublabel="Milestones reached"
            :status="rag(overview?.milestone_hit_rate_pct ?? 0, 80, 60)"
          />
          <KpiScoreCard
            label="Cursor adoption"
            :value="overview?.cursor_adoption_score ?? 0"
            format="raw"
            sublabel="AI tool usage"
            status="neutral"
          />
          <KpiScoreCard
            label="Team velocity trend"
            :value="overview?.team_velocity_trend_pct ?? 0"
            format="pct"
            :trend="overview && overview.team_velocity_trend_pct > 0 ? 'up' : overview && overview.team_velocity_trend_pct < 0 ? 'down' : 'stable'"
            trend-label="vs last period"
            :status="overview && overview.team_velocity_trend_pct >= 0 ? 'good' : 'warn'"
          />
        </div>
      </section>

      <!-- Team: who drives results -->
      <section>
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Team — who drives results
          </h2>
          <span class="rounded-full bg-gray-600/30 px-3 py-1 text-xs font-medium text-gray-400">
            {{ team.length }} members
          </span>
        </div>
        <div class="overflow-hidden rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/80 text-left text-xs font-medium uppercase tracking-wider text-gray-400">
                  <th class="px-5 py-4">#</th>
                  <th class="px-5 py-4">Member</th>
                  <th class="px-5 py-4">Role</th>
                  <th class="px-5 py-4 text-right">Delivery</th>
                  <th class="px-5 py-4 text-right">Quality</th>
                  <th class="px-5 py-4 text-right">Rework</th>
                  <th class="px-5 py-4 text-right">Score</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/80">
                <tr
                  v-for="(m, i) in team"
                  :key="m.user_id"
                  class="transition-colors hover:bg-gray-700/30"
                >
                  <td class="px-5 py-4 text-gray-500">{{ i + 1 }}</td>
                  <td class="px-5 py-4 font-medium text-white">{{ m.email }}</td>
                  <td class="px-5 py-4">
                    <span
                      class="inline-flex rounded-md px-2 py-1 text-xs font-semibold"
                      :class="roleBadgeClass(m.role)"
                    >
                      {{ m.role }}
                    </span>
                  </td>
                  <td class="px-5 py-4 text-right" :class="pctColor(m.delivery_rate_pct)">
                    {{ m.delivery_rate_pct.toFixed(1) }}%
                  </td>
                  <td class="px-5 py-4 text-right text-gray-300">{{ m.code_quality_index.toFixed(0) }}</td>
                  <td class="px-5 py-4 text-right" :class="reworkColor(m.rework_rate_pct)">
                    {{ m.rework_rate_pct.toFixed(1) }}%
                  </td>
                  <td class="px-5 py-4 text-right font-bold" :class="scoreColor(m.composite_score)">
                    {{ m.composite_score.toFixed(1) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-if="team.length === 0" class="px-5 py-8 text-center text-gray-500">
            No team data. Ensure performance module is configured.
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import type { OverviewKPIs } from '~/core/modules/performance/performance-api'
import type { TeamMemberKPI } from '~/core/modules/performance/performance-api'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import { useFinanceApi } from '~/core/modules/finance/finance-api'
import type { FinanceSummary } from '~/core/modules/finance/finance-api'
import KpiScoreCard from '~/components/performance/KpiScoreCard.vue'

const performanceStore = usePerformanceStore()
const financeApi = useFinanceApi()
const { currentUser } = useAuth()

const isLoading = ref(true)
const error = ref('')
const lastUpdated = ref('—')
const financeSummary = ref<FinanceSummary | null>(null)

const overview = computed<OverviewKPIs | null>(() => performanceStore.overview)
const team = computed<TeamMemberKPI[]>(() => performanceStore.team)

function formatFinance(value: number): string {
  return new Intl.NumberFormat('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(value)
}

function rag(value: number, goodThreshold: number, warnThreshold: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (value >= goodThreshold) return 'good'
  if (value >= warnThreshold) return 'warn'
  return 'bad'
}

const overviewStatus = computed(() => {
  const o = overview.value
  const border = (status: 'good' | 'warn' | 'bad') => {
    if (status === 'good') return 'border-emerald-500/50 bg-gray-800/60'
    if (status === 'warn') return 'border-amber-500/50 bg-gray-800/60'
    return 'border-red-500/50 bg-gray-800/60'
  }
  return {
    engineering: o ? border(rag(o.engineering_health_index, 70, 50)) : 'border-gray-700/80 bg-gray-800/60',
    projects: o ? border(rag(o.project_on_track_rate_pct, 80, 60)) : 'border-gray-700/80 bg-gray-800/60',
    velocity: o
      ? (o.team_velocity_trend_pct >= 0 ? border('good') : border('warn'))
      : 'border-gray-700/80 bg-gray-800/60',
    sprint: o ? border(rag(o.sprint_success_rate_pct, 80, 60)) : 'border-gray-700/80 bg-gray-800/60',
  }
})

function roleBadgeClass(role: string): string {
  if (role === 'CEO') return 'bg-amber-500/20 text-amber-400 ring-1 ring-amber-500/40'
  if (role === 'PM') return 'bg-purple-500/20 text-purple-300 ring-1 ring-purple-500/40'
  return 'bg-emerald-500/20 text-emerald-300 ring-1 ring-emerald-500/40'
}

function pctColor(pct: number): string {
  if (pct >= 80) return 'text-emerald-400'
  if (pct >= 60) return 'text-amber-400'
  return 'text-red-400'
}

function reworkColor(pct: number): string {
  if (pct <= 10) return 'text-emerald-400'
  if (pct <= 20) return 'text-amber-400'
  return 'text-red-400'
}

function scoreColor(score: number): string {
  if (score >= 70) return 'text-emerald-400'
  if (score >= 50) return 'text-amber-400'
  return 'text-red-400'
}

async function refresh() {
  isLoading.value = true
  error.value = ''
  try {
    await performanceStore.fetchAll(currentUser.value?.role ?? 'CEO')
    financeSummary.value = await financeApi.getSummary()
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
