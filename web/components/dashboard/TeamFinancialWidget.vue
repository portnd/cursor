<template>
  <section>
    <h2 class="section-label">Squad P&amp;L</h2>

    <!-- Total Capital hero card -->
    <div
      v-if="totalCapital != null"
      class="mt-3 rounded-2xl border border-emerald-500/40 bg-gradient-to-br from-emerald-950/40 via-emerald-900/20 to-gray-900/60 p-6 mb-4 shadow-xl overflow-hidden relative"
    >
      <div class="absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(16,185,129,0.15),transparent)] pointer-events-none"/>
      <div class="relative flex flex-wrap items-end justify-between gap-4">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-emerald-500/20 border border-emerald-500/40 flex items-center justify-center flex-shrink-0">
            <svg class="w-6 h-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <div>
            <p class="text-xs font-semibold uppercase tracking-widest text-emerald-400/80 mb-0.5">Total Capital</p>
            <p class="text-3xl font-black text-white tabular-nums tracking-tight">
              {{ formatCurrency(totalCapital) }}
            </p>
            <p class="text-sm text-gray-400 mt-1">
              <span class="font-medium text-gray-300">{{ projectCapitalCount }}</span> project{{ projectCapitalCount !== 1 ? 's' : '' }} with capital
            </p>
          </div>
        </div>
        <NuxtLink
          to="/projects"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-emerald-500/20 border border-emerald-500/40 text-emerald-300 text-sm font-medium hover:bg-emerald-500/30 hover:border-emerald-400/50 transition-all"
        >
          View projects
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
        </NuxtLink>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">

      <!-- Monthly Burn Rate -->
      <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-5 shadow-lg">
        <div class="flex items-center gap-2 mb-3">
          <div class="w-7 h-7 rounded-lg bg-red-500/15 border border-red-500/30 flex items-center justify-center flex-shrink-0">
            <svg class="w-3.5 h-3.5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"/>
            </svg>
          </div>
          <p class="text-xs font-semibold uppercase tracking-widest text-gray-400">Monthly Burn Rate</p>
        </div>
        <div v-if="isLoading" class="h-8 w-24 animate-pulse rounded-lg bg-gray-700"/>
        <p v-else class="text-2xl font-extrabold tabular-nums text-white">
          {{ formatCurrency(displayBurnRate) }}
        </p>
        <p class="text-xs text-gray-500 mt-1">
          {{ loadedMonthlyBurnRate != null && loadedMonthlyBurnRate > 0 ? 'Total loaded cost across all projects' : 'Total team salary cost' }}
        </p>
      </div>

      <!-- Runway -->
      <div
        class="rounded-2xl border p-5 shadow-lg transition-colors"
        :class="runwayCardBorderClass"
      >
        <div class="flex items-center gap-2 mb-3">
          <div
            class="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0"
            :class="runwayIconBgClass"
          >
            <svg class="w-3.5 h-3.5" :class="runwayTextClass" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <p class="text-xs font-semibold uppercase tracking-widest text-gray-400">Runway</p>
        </div>
        <div v-if="isLoading" class="h-8 w-24 animate-pulse rounded-lg bg-gray-700"/>
        <template v-else>
          <p class="text-2xl font-extrabold tabular-nums" :class="runwayTextClass">
            {{ runwayMonths > 0 ? runwayMonths.toFixed(1) + ' mo' : '—' }}
          </p>
          <p class="text-xs text-gray-500 mt-1">Capital ÷ monthly burn</p>
          <div v-if="runwayMonths > 0" class="mt-2 h-1.5 bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="runwayBarClass"
              :style="{ width: Math.min((runwayMonths / 3) * 100, 100) + '%' }"
            />
          </div>
        </template>
      </div>

      <!-- Gross Margin -->
      <div class="rounded-2xl border bg-gray-800/60 p-5 shadow-lg transition-colors"
           :class="netMargin >= 0 ? 'border-emerald-500/40' : 'border-red-500/40'">
        <div class="flex items-center gap-2 mb-3">
          <div class="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0"
               :class="netMargin >= 0 ? 'bg-emerald-500/15 border border-emerald-500/30' : 'bg-red-500/15 border border-red-500/30'">
            <svg class="w-3.5 h-3.5" :class="netMargin >= 0 ? 'text-emerald-400' : 'text-red-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path v-if="netMargin >= 0" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"/>
            </svg>
          </div>
          <p class="text-xs font-semibold uppercase tracking-widest text-gray-400">Gross Margin</p>
        </div>
        <div v-if="isLoading" class="h-8 w-24 animate-pulse rounded-lg bg-gray-700"/>
        <template v-else>
          <p class="text-2xl font-extrabold tabular-nums" :class="netMargin >= 0 ? 'text-emerald-400' : 'text-red-400'">
            {{ netMargin >= 0 ? '+' : '' }}{{ formatCurrency(netMargin) }}
          </p>
          <p class="text-xs mt-1" :class="netMargin >= 0 ? 'text-emerald-600' : 'text-red-600'">
            {{ netMargin >= 0 ? 'Squad is profitable' : 'Squad is under water' }}
          </p>
        </template>
      </div>

    </div>
  </section>
</template>

<script setup lang="ts">
import { usePricingApi } from '~/core/modules/pricing/infrastructure/pricing-api'

const props = withDefaults(
  defineProps<{
    teamMemberIds: number[]
    completedTaskIds: string[]
    totalCapital?: number | null
    projectCapitalCount?: number
    /** Loaded monthly burn rate (salary + SS + overhead). When provided, used for Monthly Burn Rate card instead of raw salary sum. */
    loadedMonthlyBurnRate?: number | null
    initialTasks?: any[]
  }>(),
  { totalCapital: null, projectCapitalCount: 0, loadedMonthlyBurnRate: null }
)

const { listSalaries, getCostConfig } = usePricingApi()
const { fetchWithAuth } = useAuth()

const isLoading = ref(true)
const burnRate = ref(0)
const earnedValue = ref(0)

/** Display burn rate: prefer loaded (salary+SS+overhead) when provided, else raw salary sum. */
const displayBurnRate = computed(() => {
  const loaded = props.loadedMonthlyBurnRate
  if (loaded != null && loaded > 0) return loaded
  return burnRate.value
})

/** Runway in months = total capital ÷ monthly burn. */
const runwayMonths = computed(() => {
  const cap = props.totalCapital ?? 0
  const burn = displayBurnRate.value
  if (burn <= 0) return 0
  return cap / burn
})

function runwayClass(mo: number) {
  if (mo > 2) return { border: 'border-emerald-500/40', icon: 'bg-emerald-500/15 border-emerald-500/30', text: 'text-emerald-400', bar: 'bg-emerald-500' }
  if (mo > 1) return { border: 'border-amber-500/40', icon: 'bg-amber-500/15 border-amber-500/30', text: 'text-amber-400', bar: 'bg-amber-500' }
  return { border: 'border-red-500/40', icon: 'bg-red-500/15 border-red-500/30', text: 'text-red-400', bar: 'bg-red-500' }
}
const runwayCardBorderClass = computed(() => {
  const r = runwayMonths.value > 0 ? runwayClass(runwayMonths.value) : { border: 'border-gray-700' }
  return r.border + ' bg-gray-800/60'
})
const runwayIconBgClass = computed(() => {
  const r = runwayMonths.value > 0 ? runwayClass(runwayMonths.value) : { icon: 'bg-gray-500/15 border-gray-500/30' }
  return r.icon + ' border'
})
const runwayTextClass = computed(() => {
  const r = runwayMonths.value > 0 ? runwayClass(runwayMonths.value) : { text: 'text-gray-400' }
  return r.text
})
const runwayBarClass = computed(() => {
  const r = runwayMonths.value > 0 ? runwayClass(runwayMonths.value) : { bar: 'bg-gray-500' }
  return r.bar
})

/** Gross Margin = Total Capital − Monthly Burn Rate (Internal VC view). */
const netMargin = computed(() => {
  const cap = props.totalCapital ?? 0
  return cap - displayBurnRate.value
})

interface Task {
  id: string
  status: string
  estimated_minutes: number
  assigned_to?: number
  completed_at?: string | null
}

const formatCurrency = (val: number) => {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    maximumFractionDigits: 0,
  }).format(val)
}

const fetchFinancials = async () => {
  isLoading.value = true
  try {
    const [salaries, config, tasksRes] = await Promise.all([
      listSalaries(),
      getCostConfig(),
      props.initialTasks && props.initialTasks.length > 0
        ? Promise.resolve({ data: props.initialTasks })
        : fetchWithAuth<{ data: Task[] }>('/sentinel/tasks'),
    ])

    const memberIds = new Set(props.teamMemberIds)

    // Burn Rate = sum of monthly salaries for team members
    const teamSalaries = salaries.filter(s => memberIds.has(s.user_id))
    burnRate.value = teamSalaries.reduce((sum, s) => sum + s.monthly_salary, 0)

    // Earned Value = completed tasks this month × cost_per_minute × estimated_minutes
    const tasks: Task[] = tasksRes?.data ?? []
    const now = new Date()
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)

    const completedThisMonth = tasks.filter(t => {
      if (t.status !== 'COMPLETED') return false
      if (!t.assigned_to || !memberIds.has(t.assigned_to)) return false
      if (t.completed_at) {
        return new Date(t.completed_at) >= monthStart
      }
      return true
    })

    const salaryMap = new Map(teamSalaries.map(s => [s.user_id, s.cost_per_minute]))
    earnedValue.value = completedThisMonth.reduce((sum, t) => {
      const cpm = salaryMap.get(t.assigned_to!) ?? 0
      return sum + cpm * (t.estimated_minutes || 0)
    }, 0)
  } catch (e) {
    // silent — parent handles errors
  } finally {
    isLoading.value = false
  }
}

watch(() => props.teamMemberIds, (newVal) => {
  if (newVal && newVal.length > 0) fetchFinancials()
}, { deep: true, immediate: true })
</script>
