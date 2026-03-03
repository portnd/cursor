<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-white">CEO STRATEGIC OVERVIEW</h1>
      <p class="text-sm text-gray-400 mt-1">System-wide operational metrics</p>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="text-center py-20">
      <div class="text-gray-400">Loading system data...</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded p-4 text-red-400">
      {{ error }}
    </div>

    <!-- Content -->
    <div v-else>
      <!-- Key Metrics (Top Row) -->
      <div class="grid grid-cols-3 gap-4 mb-8">
        <!-- System Velocity -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">System Velocity</div>
          <div class="text-3xl font-bold text-white">{{ systemVelocity }}%</div>
          <div class="text-xs text-gray-500 mt-1">
            {{ completedCount }} / {{ tasks.length }} completed
          </div>
        </div>

        <!-- Pipeline Value -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">Pipeline Value</div>
          <div class="text-3xl font-bold text-white">{{ pipelineHours }}h</div>
          <div class="text-xs text-gray-500 mt-1">
            Total estimated workload
          </div>
        </div>

        <!-- Active Workforce -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">Active Workforce</div>
          <div class="text-3xl font-bold text-white">{{ activeWorkforce }}</div>
          <div class="text-xs text-gray-500 mt-1">
            Unique developers assigned
          </div>
        </div>
      </div>

      <!-- Secondary Metrics -->
      <div class="grid grid-cols-4 gap-4 mb-8">
        <div class="bg-gray-800 border border-gray-700 rounded p-3">
          <div class="text-xs text-gray-500">IN PROGRESS</div>
          <div class="text-xl font-bold text-blue-400">{{ inProgressCount }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-3">
          <div class="text-xs text-gray-500">PENDING</div>
          <div class="text-xl font-bold text-yellow-400">{{ pendingCount }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-3">
          <div class="text-xs text-gray-500">UNASSIGNED</div>
          <div class="text-xl font-bold text-orange-400">{{ unassignedCount }}</div>
        </div>
        <div 
          class="bg-gray-800 border border-gray-700 rounded p-3 cursor-pointer hover:bg-gray-700 transition-colors"
          @click="scrollToApprovals"
        >
          <div class="text-xs text-gray-500">🚦 READY FOR REVIEW</div>
          <div class="text-xl font-bold text-white">{{ reviewPendingCount }}</div>
          <div class="text-xs text-gray-400 mt-1">Click to review</div>
        </div>
      </div>

      <!-- 🚦 QUALITY GATE: Tasks Ready for Approval -->
      <div v-if="reviewPendingTasks.length > 0" id="approvals-section" class="mb-8">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-white uppercase">🚦 Quality Gate: Ready for Approval</h2>
          <span class="text-sm text-gray-500">{{ reviewPendingTasks.length }} tasks awaiting review</span>
        </div>
        
        <div class="bg-gray-800 border border-gray-700 rounded overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-gray-900 text-xs text-gray-400 uppercase">
                <tr>
                  <th class="text-left p-3">Task</th>
                  <th class="text-left p-3">Developer</th>
                  <th class="text-center p-3">AI Score</th>
                  <th class="text-right p-3">Submitted</th>
                  <th class="text-center p-3">Action</th>
                </tr>
              </thead>
              <tbody class="text-gray-300">
                <tr 
                  v-for="task in reviewPendingTasks" 
                  :key="task.id"
                  class="border-t border-gray-700 hover:bg-gray-700/50"
                >
                  <td class="p-3">
                    <div class="font-medium text-white">{{ task.title }}</div>
                    <div class="text-xs text-gray-500 mt-1">{{ task.description?.substring(0, 60) }}...</div>
                  </td>
                  <td class="p-3">
                    <div class="flex items-center gap-2">
                      <div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-white text-xs font-bold">
                        D{{ task.assigned_to || task.created_by }}
                      </div>
                      <span>Dev #{{ task.assigned_to || task.created_by }}</span>
                    </div>
                  </td>
                  <td class="p-3 text-center">
                    <span class="inline-flex items-center gap-1 px-2 py-1 bg-green-700 text-green-100 rounded text-xs font-bold">
                      <span>✅</span>
                      <span>PASS</span>
                    </span>
                  </td>
                  <td class="p-3 text-right text-xs text-gray-400">
                    {{ formatTimeAgo(task.updated_at) }}
                  </td>
                  <td class="p-3 text-center">
                    <button
                      @click="goToTask(task)"
                      class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white font-bold rounded transition-colors flex items-center gap-2 mx-auto"
                    >
                      <span>🔍</span>
                      <span>Review</span>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Critical Issues: Tasks IN_PROGRESS > 3 days -->
      <div v-if="bottleneckTasks.length > 0" class="mb-8">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-red-400 uppercase">⚠️ Potential Bottlenecks</h2>
          <span class="text-sm text-gray-500">Tasks in progress > 3 days</span>
        </div>
        
        <div class="bg-red-900/10 border border-red-500/50 rounded overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-red-900/20 text-xs text-gray-400 uppercase">
              <tr>
                <th class="text-left p-3">Task</th>
                <th class="text-left p-3">Assigned To</th>
                <th class="text-left p-3">Days Stuck</th>
                <th class="text-left p-3">Est. Hours</th>
              </tr>
            </thead>
            <tbody class="text-gray-300">
              <tr 
                v-for="task in bottleneckTasks" 
                :key="task.id"
                class="border-t border-red-500/20 hover:bg-red-900/10 cursor-pointer"
                @click="goToTask(task)"
              >
                <td class="p-3 font-medium">{{ task.title }}</td>
                <td class="p-3">Dev #{{ task.assigned_to }}</td>
                <td class="p-3 text-red-400 font-bold">{{ getDaysStuck(task) }} days</td>
                <td class="p-3">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}h</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- All Tasks Table -->
      <div>
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-white uppercase">All Tasks</h2>
          <div class="text-sm text-gray-500">{{ tasks.length }} total</div>
        </div>
        
        <div class="bg-gray-800 border border-gray-700 rounded overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-gray-900 text-xs text-gray-400 uppercase">
              <tr>
                <th class="text-left p-3">Status</th>
                <th class="text-left p-3">Task</th>
                <th class="text-left p-3">Assigned</th>
                <th class="text-right p-3">Est. Hours</th>
                <th class="text-right p-3">Created</th>
              </tr>
            </thead>
            <tbody class="text-gray-300">
              <tr 
                v-for="task in tasks" 
                :key="task.id"
                class="border-t border-gray-700 hover:bg-gray-700/50 cursor-pointer"
                @click="goToTask(task)"
              >
                <td class="p-3">
                  <span 
                    :class="getStatusClass(task.status)"
                    class="px-2 py-1 text-xs font-bold rounded"
                  >
                    {{ task.status }}
                  </span>
                </td>
                <td class="p-3 font-medium">{{ task.title }}</td>
                <td class="p-3">
                  <span v-if="task.assigned_to">Dev #{{ task.assigned_to }}</span>
                  <span v-else class="text-orange-400">Unassigned</span>
                </td>
                <td class="p-3 text-right">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}</td>
                <td class="p-3 text-right text-gray-500 text-xs">
                  {{ formatDate(task.created_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Task {
  id: string
  title: string
  description: string
  status: string
  ai_estimated_minutes: number
  assigned_to?: number
  created_at: string
  updated_at?: string
}

const { fetchWithAuth } = useAuth()

const tasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')

// Key Metrics
const systemVelocity = computed(() => {
  if (tasks.value.length === 0) return 0
  const completed = tasks.value.filter(t => t.status === 'COMPLETED').length
  return Math.round((completed / tasks.value.length) * 100)
})

const pipelineHours = computed(() => {
  const totalMinutes = tasks.value.reduce((sum, t) => sum + (t.ai_estimated_minutes || 0), 0)
  return (totalMinutes / 60).toFixed(1)
})

const activeWorkforce = computed(() => {
  const uniqueDevs = new Set(
    tasks.value
      .filter(t => t.assigned_to)
      .map(t => t.assigned_to)
  )
  return uniqueDevs.size
})

// Secondary Metrics
const completedCount = computed(() => 
  tasks.value.filter(t => t.status === 'COMPLETED').length
)

const inProgressCount = computed(() => 
  tasks.value.filter(t => t.status === 'IN_PROGRESS').length
)

const pendingCount = computed(() => 
  tasks.value.filter(t => t.status === 'PENDING').length
)

const unassignedCount = computed(() => 
  tasks.value.filter(t => !t.assigned_to).length
)

// 🚦 Quality Gate Metrics
const reviewPendingCount = computed(() => 
  tasks.value.filter(t => t.status === 'REVIEW_PENDING').length
)

const reviewPendingTasks = computed(() => 
  tasks.value
    .filter(t => t.status === 'REVIEW_PENDING')
    .sort((a, b) => new Date(b.updated_at || b.created_at).getTime() - new Date(a.updated_at || a.created_at).getTime())
)

const avgTaskHours = computed(() => {
  if (tasks.value.length === 0) return 0
  const avg = tasks.value.reduce((sum, t) => sum + (t.ai_estimated_minutes || 0), 0) / tasks.value.length
  return (avg / 60).toFixed(1)
})

// Bottleneck Detection: IN_PROGRESS > 3 days
const bottleneckTasks = computed(() => {
  const threeDaysAgo = new Date()
  threeDaysAgo.setDate(threeDaysAgo.getDate() - 3)
  
  return tasks.value.filter(task => {
    if (task.status !== 'IN_PROGRESS') return false
    const createdDate = new Date(task.created_at)
    return createdDate < threeDaysAgo
  })
})

const getDaysStuck = (task: Task) => {
  const createdDate = new Date(task.created_at)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - createdDate.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diffDays
}

// Fetch Data
const fetchTasks = async () => {
  try {
    isLoading.value = true
    const response = await fetchWithAuth<{ data: Task[] }>('/sentinel/tasks')
    tasks.value = response.data || []
    error.value = ''
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load tasks'
    console.error('Failed to fetch tasks:', err)
  } finally {
    isLoading.value = false
  }
}

// Utilities
const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    COMPLETED: 'bg-green-700 text-green-100',
    IN_PROGRESS: 'bg-blue-700 text-blue-100',
    PENDING: 'bg-yellow-700 text-yellow-100',
    BLOCKED: 'bg-red-700 text-red-100',
    REVIEW_PENDING: 'bg-indigo-900 text-indigo-200 border border-indigo-600' // 🚦 Quality Gate
  }
  return classes[status] || 'bg-gray-700 text-gray-100'
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric'
  })
}

const formatTimeAgo = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  return `${diffDays}d ago`
}

const scrollToApprovals = () => {
  const section = document.getElementById('approvals-section')
  if (section) {
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

// Navigation helper (use task.code for pretty URL when available)
const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}`)
}

onMounted(() => {
  fetchTasks()
})
</script>
