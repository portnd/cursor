<template>
  <div class="project-analytics">
    <!-- Key Metrics Row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <div class="metric-card">
        <div class="metric-value text-purple-400">{{ analytics.completed_tasks }}/{{ analytics.total_tasks }}</div>
        <div class="metric-label">Tasks Completed</div>
        <div class="mt-2 h-1.5 bg-gray-700 rounded-full overflow-hidden">
          <div class="h-full bg-purple-500 rounded-full" :style="{ width: completionPct + '%' }"></div>
        </div>
      </div>
      <div class="metric-card">
        <div class="metric-value text-purple-400">{{ analytics.completed_story_points }}/{{ analytics.total_story_points }}</div>
        <div class="metric-label">Story Points Done</div>
        <div class="mt-2 h-1.5 bg-gray-700 rounded-full overflow-hidden">
          <div class="h-full bg-purple-500 rounded-full" :style="{ width: spPct + '%' }"></div>
        </div>
      </div>
      <div class="metric-card">
        <div class="metric-value text-green-400">{{ totalLoggedHours }}h</div>
        <div class="metric-label">Hours Logged</div>
      </div>
      <div class="metric-card">
        <div class="metric-value text-yellow-400">{{ avgCycleTimeFormatted }}</div>
        <div class="metric-label">Avg Cycle Time</div>
      </div>
    </div>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <!-- Burndown Chart -->
      <div class="chart-card">
        <h3 class="chart-title">Sprint Burndown</h3>
        <div v-if="analytics.burndown?.length" class="h-56">
          <canvas ref="burndownCanvas"></canvas>
        </div>
        <div v-else class="flex items-center justify-center h-56 text-gray-500 text-sm">
          No active sprint or no burndown data available.
        </div>
      </div>

      <!-- Velocity Chart -->
      <div class="chart-card">
        <h3 class="chart-title">Team Velocity</h3>
        <div v-if="analytics.velocity?.length" class="h-56">
          <canvas ref="velocityCanvas"></canvas>
        </div>
        <div v-else class="flex items-center justify-center h-56 text-gray-500 text-sm">
          No sprint history available yet.
        </div>
      </div>
    </div>

    <!-- Team Capacity Table -->
    <div class="chart-card">
      <h3 class="chart-title mb-4">Team Capacity</h3>
      <div v-if="analytics.team_capacity?.length" class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-700">
              <th class="text-left py-2 px-3 text-gray-400 font-medium">Developer</th>
              <th class="text-right py-2 px-3 text-gray-400 font-medium">Tasks</th>
              <th class="text-right py-2 px-3 text-gray-400 font-medium">Est. Hours</th>
              <th class="text-right py-2 px-3 text-gray-400 font-medium">Logged Hours</th>
              <th class="text-right py-2 px-3 text-gray-400 font-medium">Utilization</th>
              <th class="py-2 px-3 text-gray-400 font-medium"></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in analytics.team_capacity"
              :key="row.user_id"
              class="border-b border-gray-800 hover:bg-gray-800/50"
            >
              <td class="py-3 px-3">
                <div class="flex items-center gap-2">
                  <div class="w-7 h-7 rounded-full bg-purple-600 flex items-center justify-center text-white text-xs font-bold">
                    {{ capacityInitial(row) }}
                  </div>
                  <div class="min-w-0">
                    <span class="text-gray-300 block truncate">{{ capacityLabel(row) }}</span>
                    <span
                      v-if="row.user_display_name?.trim() && row.user_email"
                      class="text-xs text-gray-500 truncate block"
                    >{{ row.user_email }}</span>
                  </div>
                </div>
              </td>
              <td class="text-right py-3 px-3 text-gray-300">{{ row.assigned_tasks }}</td>
              <td class="text-right py-3 px-3 text-gray-300">{{ row.estimated_hours.toFixed(1) }}h</td>
              <td class="text-right py-3 px-3 text-gray-300">{{ row.logged_hours.toFixed(1) }}h</td>
              <td class="text-right py-3 px-3">
                <span :class="utilizationClass(row.utilization_pct)">
                  {{ row.utilization_pct.toFixed(0) }}%
                </span>
              </td>
              <td class="py-3 px-3">
                <div class="w-24 h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="utilizationBarClass(row.utilization_pct)"
                    :style="{ width: Math.min(row.utilization_pct, 100) + '%' }"
                  ></div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-center py-8 text-gray-500 text-sm">
        No team assignment data yet.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch, nextTick } from 'vue'
import type { ProjectAnalytics } from '~/core/modules/projects/infrastructure/projects-api'

const props = defineProps<{
  analytics: ProjectAnalytics
}>()

const burndownCanvas = ref<HTMLCanvasElement | null>(null)
const velocityCanvas = ref<HTMLCanvasElement | null>(null)

let burndownChart: any = null
let velocityChart: any = null

const completionPct = computed(() =>
  props.analytics.total_tasks > 0
    ? Math.round((props.analytics.completed_tasks / props.analytics.total_tasks) * 100)
    : 0
)

const spPct = computed(() =>
  props.analytics.total_story_points > 0
    ? Math.round((props.analytics.completed_story_points / props.analytics.total_story_points) * 100)
    : 0
)

const totalLoggedHours = computed(() =>
  (props.analytics.total_logged_minutes / 60).toFixed(1)
)

const avgCycleTimeFormatted = computed(() => {
  const d = props.analytics.avg_cycle_time_days
  if (!d) return '—'
  if (d < 1) return `${Math.round(d * 24)}h`
  return `${d.toFixed(1)}d`
})

function utilizationClass(pct: number) {
  if (pct > 120) return 'text-red-400 font-semibold'
  if (pct > 90) return 'text-green-400 font-semibold'
  if (pct > 50) return 'text-yellow-400'
  return 'text-gray-400'
}

function utilizationBarClass(pct: number) {
  if (pct > 120) return 'bg-red-500'
  if (pct > 90) return 'bg-green-500'
  if (pct > 50) return 'bg-yellow-500'
  return 'bg-gray-500'
}

function capacityLabel(row: ProjectAnalytics['team_capacity'][number]) {
  const name = row.user_display_name?.trim()
  if (name) return name
  if (row.user_email?.trim()) return row.user_email.trim()
  return `Dev #${row.user_id}`
}

function capacityInitial(row: ProjectAnalytics['team_capacity'][number]) {
  const src = capacityLabel(row)
  return src.charAt(0).toUpperCase() || '?'
}

async function renderCharts() {
  if (!import.meta.client) return
  const { Chart, registerables } = await import('chart.js')
  Chart.register(...registerables)

  // Burndown Chart
  if (burndownCanvas.value && props.analytics.burndown?.length) {
    if (burndownChart) burndownChart.destroy()
    const ctx = burndownCanvas.value.getContext('2d')!
    burndownChart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: props.analytics.burndown.map((p) => p.day),
        datasets: [
          {
            label: 'Ideal',
            data: props.analytics.burndown.map((p) => p.ideal),
            borderColor: 'rgba(99, 102, 241, 0.5)',
            borderDash: [5, 5],
            pointRadius: 0,
            tension: 0.1,
          },
          {
            label: 'Actual',
            data: props.analytics.burndown.map((p) => p.remaining),
            borderColor: 'rgb(99, 102, 241)',
            backgroundColor: 'rgba(99, 102, 241, 0.1)',
            fill: true,
            tension: 0.3,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { labels: { color: '#9ca3af', font: { size: 11 } } } },
        scales: {
          x: { ticks: { color: '#6b7280', maxTicksLimit: 7 }, grid: { color: 'rgba(75,85,99,0.2)' } },
          y: { ticks: { color: '#6b7280' }, grid: { color: 'rgba(75,85,99,0.2)' }, min: 0 },
        },
      },
    })
  }

  // Velocity Chart
  if (velocityCanvas.value && props.analytics.velocity?.length) {
    if (velocityChart) velocityChart.destroy()
    const ctx = velocityCanvas.value.getContext('2d')!
    velocityChart = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: props.analytics.velocity.map((v) => v.sprint_name),
        datasets: [
          {
            label: 'Planned SP',
            data: props.analytics.velocity.map((v) => v.planned_sp),
            backgroundColor: 'rgba(99, 102, 241, 0.3)',
            borderColor: 'rgba(99, 102, 241, 0.6)',
            borderWidth: 1,
          },
          {
            label: 'Completed SP',
            data: props.analytics.velocity.map((v) => v.completed_sp),
            backgroundColor: 'rgba(34, 197, 94, 0.5)',
            borderColor: 'rgba(34, 197, 94, 0.8)',
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { labels: { color: '#9ca3af', font: { size: 11 } } } },
        scales: {
          x: { ticks: { color: '#6b7280' }, grid: { color: 'rgba(75,85,99,0.2)' } },
          y: { ticks: { color: '#6b7280' }, grid: { color: 'rgba(75,85,99,0.2)' }, min: 0 },
        },
      },
    })
  }
}

onMounted(async () => {
  await nextTick()
  renderCharts()
})

watch(() => props.analytics, async () => {
  await nextTick()
  renderCharts()
}, { deep: true })
</script>

<style scoped>
.metric-card {
  @apply bg-gray-800/60 rounded-xl p-4 border border-gray-700/50;
}
.metric-value {
  @apply text-2xl font-bold tabular-nums;
}
.metric-label {
  @apply text-xs text-gray-400 mt-1 uppercase tracking-wide;
}
.chart-card {
  @apply bg-gray-800/60 rounded-xl p-5 border border-gray-700/50;
}
.chart-title {
  @apply text-sm font-semibold text-gray-300 mb-4;
}
</style>
