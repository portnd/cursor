<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 p-6">
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent">
        THE ARCHITECT: GRAND TIMELINE (GANTT)
      </h1>
      <p class="text-sm text-gray-400 mt-1">Project planning with Vue-Ganttastic</p>
    </div>

    <div class="flex flex-wrap items-center gap-3 mb-4">
      <span class="text-sm text-gray-400">View:</span>
      <button
        type="button"
        :class="[ 'px-4 py-2 rounded-lg font-medium transition-all', viewMode === 'day' ? 'bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white shadow-lg shadow-purple-500/20' : 'bg-gray-700 text-gray-300 hover:bg-gray-600' ]"
        @click="viewMode = 'day'"
      >
        Day
      </button>
      <button
        type="button"
        :class="[ 'px-4 py-2 rounded-lg font-medium transition-all', viewMode === 'week' ? 'bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white shadow-lg shadow-purple-500/20' : 'bg-gray-700 text-gray-300 hover:bg-gray-600' ]"
        @click="viewMode = 'week'"
      >
        Week
      </button>
      <button
        type="button"
        :class="[ 'px-4 py-2 rounded-lg font-medium transition-all', viewMode === 'month' ? 'bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white shadow-lg shadow-purple-500/20' : 'bg-gray-700 text-gray-300 hover:bg-gray-600' ]"
        @click="viewMode = 'month'"
      >
        Month
      </button>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลด Gantt...</p>
    </div>

    <div
      v-else-if="error"
      class="bg-red-900/20 border-2 border-red-500 rounded-lg p-6 text-red-400"
    >
      <p class="font-bold">Failed to load Gantt</p>
      <p class="text-sm mt-1">{{ error }}</p>
      <button
        type="button"
        class="mt-4 px-4 py-2 bg-red-100 dark:bg-red-600 hover:bg-red-200 dark:bg-red-700 rounded text-gray-900 dark:text-white text-sm font-medium"
        @click="loadGanttData"
      >
        Retry
      </button>
    </div>

    <div
      v-else-if="ganttRows.length === 0"
      class="text-center py-20 text-gray-500"
    >
      <p class="text-lg">No tasks yet. Create tasks to see them on the timeline.</p>
    </div>

    <ClientOnly v-else>
      <div class="rounded-lg overflow-hidden border border-gray-700 bg-gray-800">
        <g-gantt-chart
          :chart-start="chartStart"
          :chart-end="chartEnd"
          :precision="precision"
          bar-start="barStart"
          bar-end="barEnd"
          date-format="YYYY-MM-DD"
          :width="chartWidth + 'px'"
          :row-height="36"
          :grid="true"
          :current-time="true"
          color-scheme="dark"
          label-column-title="Task"
          label-column-width="200px"
          @click-bar="onClickBar"
          @dragend-bar="onDragEnd"
        >
          <g-gantt-row
            v-for="row in ganttRows"
            :key="row.taskId"
            :label="row.label"
            :bars="row.bars"
          />
        </g-gantt-chart>
      </div>
      <template #fallback>
        <div class="min-h-[400px] flex items-center justify-center text-gray-500">Loading chart...</div>
      </template>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
  layout: 'default'
})

interface ApiTask {
  id: string
  title: string
  start_date: string | null
  end_date: string | null
  progress: number
}

interface ApiDependency {
  predecessor_id: string
  successor_id: string
}

const { fetchWithAuth } = useAuth()
const router = useRouter()

const isLoading = ref(true)
const error = ref('')
const tasks = ref<ApiTask[]>([])
const viewMode = ref<'day' | 'week' | 'month'>('week')

function toYMD(d: string | Date): string {
  if (typeof d === 'string') return d.split('T')[0]
  return d.toISOString().split('T')[0]
}

const chartStart = computed(() => {
  const list = tasks.value.filter((t) => t.start_date || t.end_date)
  const pad = 30
  if (!list.length) {
    const d = new Date()
    d.setDate(d.getDate() - pad)
    return toYMD(d.toISOString())
  }
  let min = Infinity
  for (const t of list) {
    const v = t.start_date ? new Date(t.start_date).getTime() : (t.end_date ? new Date(t.end_date).getTime() : null)
    if (v != null) min = Math.min(min, v)
  }
  const d = new Date(min === Infinity ? Date.now() : min)
  d.setDate(d.getDate() - pad)
  return toYMD(d.toISOString())
})

const chartEnd = computed(() => {
  const list = tasks.value.filter((t) => t.start_date || t.end_date)
  const pad = 30
  if (!list.length) {
    const d = new Date()
    d.setDate(d.getDate() + pad)
    return toYMD(d.toISOString())
  }
  let max = -Infinity
  for (const t of list) {
    const v = t.end_date ? new Date(t.end_date).getTime() : (t.start_date ? new Date(t.start_date).getTime() : null)
    if (v != null) max = Math.max(max, v)
  }
  const d = new Date(max === -Infinity ? Date.now() : max)
  d.setDate(d.getDate() + pad)
  return toYMD(d.toISOString())
})

const precision = computed(() => viewMode.value)

const chartWidth = computed(() => {
  const start = new Date(chartStart.value + 'T00:00:00Z').getTime()
  const end = new Date(chartEnd.value + 'T00:00:00Z').getTime()
  const days = Math.max(1, (end - start) / 86400000)
  const px = viewMode.value === 'month' ? 4 : viewMode.value === 'week' ? 24 : 40
  return Math.max(800, Math.min(5000, Math.round(days * px)))
})

const ganttRows = computed(() => {
  const today = toYMD(new Date().toISOString())
  const tomorrow = toYMD(new Date(Date.now() + 86400000).toISOString())
  const addDay = (ymd: string) => {
    const d = new Date(ymd + 'T12:00:00Z')
    d.setUTCDate(d.getUTCDate() + 1)
    return toYMD(d.toISOString())
  }
  return tasks.value.map((t) => {
    let start = t.start_date ? toYMD(t.start_date) : (t.end_date ? toYMD(t.end_date) : today)
    let end = t.end_date ? toYMD(t.end_date) : (t.start_date ? toYMD(t.start_date) : tomorrow)
    if (start === end) end = addDay(start)
    if (end < start) end = addDay(start)
    const label = t.title || t.id
    return {
      taskId: t.id,
      label: label.length > 50 ? label.slice(0, 47) + '...' : label,
      bars: [
        {
          barStart: start,
          barEnd: end,
          ganttBarConfig: { id: t.id, label: label.slice(0, 40), hasHandles: true }
        }
      ]
    }
  })
})

async function loadGanttData() {
  isLoading.value = true
  error.value = ''
  try {
    const res = await fetchWithAuth<{ data: { tasks: ApiTask[] } }>('/sentinel/tasks/gantt')
    tasks.value = res?.data?.tasks ?? []
  } catch (e: any) {
    error.value = e?.data?.message || e?.message || 'Failed to load Gantt data'
  } finally {
    isLoading.value = false
  }
}

function onClickBar(payload: { bar: { ganttBarConfig: { id: string } } }) {
  const id = payload?.bar?.ganttBarConfig?.id
  if (id) router.push(`/task/${id}`)
}

async function onDragEnd(payload: { movedBars?: Map<string, { start: string | Date; end: string | Date }> }) {
  const map = payload.movedBars
  if (!map || map.size === 0) return
  for (const [taskId, range] of map) {
    const start = typeof range.start === 'string' ? range.start : (range.start as Date).toISOString()
    const end = typeof range.end === 'string' ? range.end : (range.end as Date).toISOString()
    try {
      await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
        method: 'PATCH',
        body: { start_date: start, end_date: end }
      })
      const idx = tasks.value.findIndex((t) => t.id === taskId)
      if (idx !== -1) {
        tasks.value[idx] = { ...tasks.value[idx], start_date: start, end_date: end }
      }
    } catch (e) {
      console.error('Failed to update task dates:', e)
    }
  }
}

onMounted(loadGanttData)
</script>
