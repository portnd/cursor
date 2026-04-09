<template>
  <section>
    <div class="flex items-center justify-between mb-4 gap-4">
      <div class="min-w-0">
        <h2 class="section-label mb-0">{{ sectionTitle }}</h2>
        <p v-if="scopeDescription" class="text-xs text-gray-500 mt-1 max-w-xl leading-relaxed">{{ scopeDescription }}</p>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-xs text-gray-500">{{ deltaRows.length }} tasks tracked</span>
        <button
          @click="fetchData"
          :disabled="isLoading"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-gray-900 dark:text-white transition-colors disabled:opacity-50"
        >
          <svg class="h-3 w-3" :class="isLoading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div v-if="!isLoading && deltaRows.length > 0" class="grid grid-cols-2 gap-4 sm:grid-cols-4 mb-5">
      <div class="rounded-xl border border-gray-700 bg-gray-800/60 p-4">
        <p class="text-xs text-gray-500 mb-1">Avg Variance</p>
        <p class="text-lg font-bold tabular-nums" :class="summary.avgVariancePct >= 0 ? 'text-red-400' : 'text-emerald-400'">
          {{ summary.avgVariancePct >= 0 ? '+' : '' }}{{ summary.avgVariancePct.toFixed(1) }}%
        </p>
        <p class="text-xs text-gray-600 mt-0.5">vs estimated</p>
      </div>
      <div class="rounded-xl border border-gray-700 bg-gray-800/60 p-4">
        <p class="text-xs text-gray-500 mb-1">On / Under Budget</p>
        <p class="text-lg font-bold tabular-nums text-emerald-400">{{ summary.onBudget }}</p>
        <p class="text-xs text-gray-600 mt-0.5">tasks within estimate</p>
      </div>
      <div class="rounded-xl border border-gray-700 bg-gray-800/60 p-4">
        <p class="text-xs text-gray-500 mb-1">Over Budget</p>
        <p class="text-lg font-bold tabular-nums text-red-400">{{ summary.overBudget }}</p>
        <p class="text-xs text-gray-600 mt-0.5">tasks exceeded estimate</p>
      </div>
      <div class="rounded-xl border border-gray-700 bg-gray-800/60 p-4">
        <p class="text-xs text-gray-500 mb-1">Total Overrun</p>
        <p class="text-lg font-bold tabular-nums text-amber-400">{{ formatMinutes(summary.totalOverrunMins) }}</p>
        <p class="text-xs text-gray-600 mt-0.5">unplanned time</p>
      </div>
    </div>

    <!-- Table -->
    <div class="rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden shadow-lg">

      <!-- Loading -->
      <div v-if="isLoading" class="flex items-center justify-center py-16">
        <svg class="h-6 w-6 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
      </div>

      <!-- Empty -->
      <div v-else-if="deltaRows.length === 0" class="flex flex-col items-center justify-center py-14 text-center">
        <svg class="h-9 w-9 text-gray-600 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>
        <p class="text-sm text-gray-400 font-medium">No time logs found</p>
        <p class="text-xs text-gray-600 mt-0.5">Developers must log time on tasks for delta analysis</p>
      </div>

      <!-- Data Table -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-700 bg-gray-900/60 text-xs uppercase tracking-wider text-gray-500">
              <th class="px-5 py-3.5 text-left">Task</th>
              <th class="px-5 py-3.5 text-right">Estimated</th>
              <th class="px-5 py-3.5 text-right">Actual Logged</th>
              <th class="px-5 py-3.5 text-right">Delta</th>
              <th class="px-5 py-3.5 text-right">Variance</th>
              <th class="px-5 py-3.5 text-center">Accuracy</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700/60">
            <tr
              v-for="row in paginatedRows"
              :key="row.taskId"
              class="hover:bg-gray-700/15 transition-colors"
            >
              <td class="px-5 py-3.5 max-w-xs">
                <button @click="navigateTo(`/task/${row.taskCode || row.taskId}`)" class="text-left hover:text-blue-300 transition-colors">
                  <p class="font-medium text-white truncate">{{ row.taskTitle }}</p>
                  <p v-if="row.taskCode" class="text-xs text-gray-600 font-mono mt-0.5">{{ row.taskCode }}</p>
                </button>
              </td>
              <td class="px-5 py-3.5 text-right text-gray-300 tabular-nums">
                {{ formatMinutes(row.estimatedMins) }}
              </td>
              <td class="px-5 py-3.5 text-right tabular-nums" :class="row.actualMins > row.estimatedMins ? 'text-red-300' : 'text-gray-300'">
                {{ formatMinutes(row.actualMins) }}
              </td>
              <td class="px-5 py-3.5 text-right tabular-nums font-semibold" :class="row.deltaMins <= 0 ? 'text-emerald-400' : 'text-red-400'">
                {{ row.deltaMins > 0 ? '+' : '' }}{{ formatMinutes(row.deltaMins) }}
              </td>
              <td class="px-5 py-3.5 text-right tabular-nums font-semibold" :class="row.variancePct <= 0 ? 'text-emerald-400' : 'text-red-400'">
                {{ row.variancePct > 0 ? '+' : '' }}{{ row.variancePct.toFixed(1) }}%
              </td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center">
                  <!-- Mini bar: fill = min(actual/estimated, 1) -->
                  <div class="relative h-1.5 w-20 rounded-full bg-gray-700 overflow-hidden">
                    <div
                      class="absolute left-0 top-0 h-full rounded-full"
                      :class="row.actualMins <= row.estimatedMins ? 'bg-emerald-500' : 'bg-red-500'"
                      :style="{ width: `${Math.min((row.estimatedMins > 0 ? (row.actualMins / row.estimatedMins) * 100 : 100), 100)}%` }"
                    />
                    <!-- Overflow indicator -->
                    <div
                      v-if="row.actualMins > row.estimatedMins"
                      class="absolute right-0 top-0 h-full w-1 rounded-full bg-red-400 animate-pulse"
                    />
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-700 px-5 py-3">
          <span class="text-xs text-gray-500">
            Showing {{ (currentPage - 1) * pageSize + 1 }}–{{ Math.min(currentPage * pageSize, deltaRows.length) }} of {{ deltaRows.length }}
          </span>
          <div class="flex items-center gap-1">
            <button
              @click="currentPage--"
              :disabled="currentPage === 1"
              class="px-2 py-1 rounded text-xs text-gray-400 hover:text-gray-900 dark:text-white hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >Prev</button>
            <span class="px-2 text-xs text-gray-500">{{ currentPage }}/{{ totalPages }}</span>
            <button
              @click="currentPage++"
              :disabled="currentPage === totalPages"
              class="px-2 py-1 rounded text-xs text-gray-400 hover:text-gray-900 dark:text-white hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >Next</button>
          </div>
        </div>
      </div>

    </div>
  </section>
</template>

<script setup lang="ts">
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { Task } from '~/core/modules/projects/infrastructure/projects-api'

interface DeltaRow {
  taskId: string
  taskCode?: string
  taskTitle: string
  estimatedMins: number
  actualMins: number
  deltaMins: number
  variancePct: number
}

const props = withDefaults(
  defineProps<{
    teamProjectIds: string[]
    initialTasks?: any[]
    /** Main heading (e.g. portfolio wording when squads are off). */
    sectionTitle?: string
    /** Shown under the title (e.g. squad vs Product Owner–owned scope). */
    scopeDescription?: string
  }>(),
  {
    sectionTitle: 'Estimation Accuracy — Delta Analysis',
    scopeDescription: '',
  }
)

const { getTasksByProject, getTimeLogs } = useTasksApi()

const isLoading = ref(true)
const deltaRows = ref<DeltaRow[]>([])
const currentPage = ref(1)
const pageSize = 15

const totalPages = computed(() => Math.ceil(deltaRows.value.length / pageSize))

const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return deltaRows.value.slice(start, start + pageSize)
})

const summary = computed(() => {
  const rows = deltaRows.value
  const onBudget = rows.filter(r => r.deltaMins <= 0).length
  const overBudget = rows.filter(r => r.deltaMins > 0).length
  const totalOverrunMins = rows.filter(r => r.deltaMins > 0).reduce((s, r) => s + r.deltaMins, 0)
  const avgVariancePct = rows.length > 0
    ? rows.reduce((s, r) => s + r.variancePct, 0) / rows.length
    : 0
  return { onBudget, overBudget, totalOverrunMins, avgVariancePct }
})

const formatMinutes = (mins: number) => {
  if (!mins && mins !== 0) return '—'
  const sign = mins < 0 ? '-' : ''
  const abs = Math.abs(mins)
  const h = Math.floor(abs / 60)
  const m = abs % 60
  if (h === 0) return `${sign}${m}m`
  if (m === 0) return `${sign}${h}h`
  return `${sign}${h}h ${m}m`
}

const fetchData = async () => {
  isLoading.value = true
  currentPage.value = 1
  try {
    if (props.teamProjectIds.length === 0) return

    // Use pre-fetched tasks from parent if available (avoids duplicate API call)
    let allTasks: any[]
    if (props.initialTasks && props.initialTasks.length > 0) {
      allTasks = props.initialTasks.filter((t: any) =>
        props.teamProjectIds.includes(t.project_id) && t.estimated_minutes > 0
      )
    } else {
      const taskArrays = await Promise.all(
        props.teamProjectIds.map(pid => getTasksByProject(pid).catch(() => [] as Task[]))
      )
      allTasks = taskArrays.flat().filter(t => t.estimated_minutes > 0)
    }

    // Batch fetch time logs (max 20 concurrent)
    const BATCH = 20
    const rows: DeltaRow[] = []

    for (let i = 0; i < allTasks.length; i += BATCH) {
      const batch = allTasks.slice(i, i + BATCH)
      const logBatch = await Promise.all(
        batch.map(t => getTimeLogs(t.id).catch(() => [] as TimeLog[]))
      )
      batch.forEach((task, idx) => {
        const logs = logBatch[idx]
        if (logs.length === 0) return
        const actualMins = logs.reduce((s, l) => s + l.minutes, 0)
        const estimatedMins = task.estimated_minutes
        const deltaMins = actualMins - estimatedMins
        const variancePct = estimatedMins > 0 ? (deltaMins / estimatedMins) * 100 : 0
        rows.push({
          taskId: task.id,
          taskCode: task.code,
          taskTitle: task.title,
          estimatedMins,
          actualMins,
          deltaMins,
          variancePct,
        })
      })
    }

    // Sort: worst over-runs first
    deltaRows.value = rows.sort((a, b) => b.variancePct - a.variancePct)
  } catch {
    // silent
  } finally {
    isLoading.value = false
  }
}

onMounted(() => fetchData())

watch(() => props.teamProjectIds, () => {
  if (props.teamProjectIds.length > 0) fetchData()
}, { deep: true })
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
</style>
