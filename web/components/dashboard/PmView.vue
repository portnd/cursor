<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-white">PM RESOURCE CONTROL</h1>
        <p class="text-sm text-gray-400 mt-1">Task allocation and workforce management</p>
      </div>
      <button
        @click="fetchTasks"
        class="px-4 py-2 bg-gray-800 hover:bg-gray-700 border border-gray-600 text-gray-300 text-sm font-medium rounded transition-colors"
      >
        ⟳ Refresh
      </button>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="text-center py-20">
      <div class="text-gray-400">Loading resource data...</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded p-4 text-red-400">
      {{ error }}
    </div>

    <!-- 🚦 QUALITY GATE: Tasks Ready for Approval -->
    <div v-if="!isLoading && !error && reviewPendingTasks.length > 0" class="mb-6">
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
                  {{ formatTimeAgo(task.updated_at || task.created_at) }}
                </td>
                <td class="p-3 text-center">
                  <button
                    @click="goToTask(task.id)"
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

    <!-- Split View Layout -->
    <div v-else class="grid grid-cols-2 gap-6">
      <!-- LEFT: Unassigned Queue -->
      <div>
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-lg font-bold text-orange-400 uppercase">Unassigned Queue</h2>
          <span class="text-sm text-gray-500">{{ unassignedTasks.length }} tasks</span>
        </div>

        <div v-if="unassignedTasks.length === 0" class="bg-gray-800 border border-gray-700 rounded p-8 text-center text-gray-500">
          No unassigned tasks
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="task in unassignedTasks"
            :key="task.id"
            class="bg-gray-800 border border-orange-500/30 hover:border-orange-500/60 rounded p-3 transition-colors cursor-pointer hover:scale-[1.01]"
            @click="goToTask(task.id)"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="flex-1 min-w-0">
                <div class="font-medium text-white text-sm mb-1 truncate">
                  {{ task.title }}
                </div>
                <div class="text-xs text-gray-500">
                  Est: {{ (task.ai_estimated_minutes / 60).toFixed(1) }}h
                  | Created: {{ formatDate(task.created_at) }}
                </div>
              </div>
              <button
                @click.stop="openAssignModal(task)"
                class="shrink-0 px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded transition-colors"
              >
                Assign
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT: Active Development (Grouped by Developer) -->
      <div>
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-lg font-bold text-blue-400 uppercase">Active Development</h2>
          <span class="text-sm text-gray-500">{{ Object.keys(tasksByDeveloper).length }} developers</span>
        </div>

        <div v-if="Object.keys(tasksByDeveloper).length === 0" class="bg-gray-800 border border-gray-700 rounded p-8 text-center text-gray-500">
          No assigned tasks
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="(devTasks, devId) in tasksByDeveloper"
            :key="devId"
            class="bg-gray-800 border rounded p-4"
            :class="[
              devTasks.length > 3 
                ? 'border-red-500/50 bg-red-900/10' 
                : 'border-gray-700'
            ]"
          >
            <!-- Developer Header -->
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-2">
                <span class="font-bold text-white">Developer #{{ devId }}</span>
                <span class="text-xs text-gray-500">{{ devTasks.length }} tasks</span>
                <span 
                  v-if="devTasks.length > 3"
                  class="text-xs font-bold text-red-400 bg-red-900/30 px-2 py-0.5 rounded"
                >
                  ⚠️ OVERLOADED
                </span>
              </div>
              <div class="text-xs text-gray-500">
                {{ getDevWorkload(devTasks) }}h workload
              </div>
            </div>

            <!-- Developer's Tasks -->
            <div class="space-y-2">
              <div
                v-for="task in devTasks"
                :key="task.id"
                class="bg-gray-900/50 rounded p-2 flex items-center justify-between gap-2 cursor-pointer hover:bg-gray-800/70 transition-colors"
                @click="goToTask(task.id)"
              >
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span 
                      :class="getStatusDot(task.status)"
                      class="w-2 h-2 rounded-full shrink-0"
                    ></span>
                    <span class="text-sm text-gray-300 truncate">{{ task.title }}</span>
                  </div>
                  <div class="text-xs text-gray-600 ml-4">
                    {{ (task.ai_estimated_minutes / 60).toFixed(1) }}h
                  </div>
                </div>
                <button
                  @click.stop="openReassignModal(task)"
                  class="shrink-0 px-2 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 text-xs rounded transition-colors"
                >
                  Re-assign
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Assignment Modal -->
    <div
      v-if="showAssignModal"
      class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4"
      @click.self="closeAssignModal"
    >
      <div class="bg-gray-800 border border-gray-600 rounded-lg p-6 max-w-md w-full">
        <h3 class="text-lg font-bold text-white mb-4">
          {{ isReassign ? 'Re-assign Task' : 'Assign Task' }}
        </h3>
        
        <p class="text-sm text-gray-400 mb-4 line-clamp-2">{{ selectedTask?.title }}</p>

        <div class="mb-6">
          <label class="block text-sm text-gray-400 mb-2">Developer ID</label>
          <input
            v-model.number="assignDevId"
            type="number"
            placeholder="Enter developer ID (e.g., 1, 2, 3)"
            class="w-full px-3 py-2 bg-gray-900 border border-gray-600 rounded text-white text-sm focus:outline-none focus:border-blue-500"
          />
        </div>

        <div class="flex gap-3">
          <button
            @click="assignTask"
            :disabled="isAssigning || !assignDevId"
            class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 text-white font-bold rounded transition-colors disabled:cursor-not-allowed"
          >
            {{ isAssigning ? 'Assigning...' : 'Assign' }}
          </button>
          <button
            @click="closeAssignModal"
            class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Success Notification -->
    <div
      v-if="showSuccess"
      class="fixed top-4 right-4 bg-green-700 border border-green-600 text-white px-4 py-3 rounded shadow-lg z-50"
    >
      {{ successMessage }}
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

// Modal State
const showAssignModal = ref(false)
const selectedTask = ref<Task | null>(null)
const assignDevId = ref<number | null>(null)
const isAssigning = ref(false)
const isReassign = ref(false)

// Success State
const showSuccess = ref(false)
const successMessage = ref('')

// Computed
const unassignedTasks = computed(() => 
  tasks.value.filter(t => !t.assigned_to)
)

const tasksByDeveloper = computed(() => {
  const grouped: Record<number, Task[]> = {}
  tasks.value
    .filter(t => t.assigned_to)
    .forEach(task => {
      const devId = task.assigned_to!
      if (!grouped[devId]) grouped[devId] = []
      grouped[devId].push(task)
    })
  return grouped
})

// 🚦 Quality Gate Metrics
const reviewPendingTasks = computed(() => 
  tasks.value
    .filter(t => t.status === 'REVIEW_PENDING')
    .sort((a, b) => new Date(b.updated_at || b.created_at).getTime() - new Date(a.updated_at || a.created_at).getTime())
)

// Methods
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

const openAssignModal = (task: Task) => {
  selectedTask.value = task
  assignDevId.value = null
  isReassign.value = false
  showAssignModal.value = true
}

const openReassignModal = (task: Task) => {
  selectedTask.value = task
  assignDevId.value = task.assigned_to || null
  isReassign.value = true
  showAssignModal.value = true
}

const closeAssignModal = () => {
  showAssignModal.value = false
  selectedTask.value = null
  assignDevId.value = null
  isReassign.value = false
}

const assignTask = async () => {
  if (!selectedTask.value || !assignDevId.value) return

  try {
    isAssigning.value = true
    
    await fetchWithAuth(`/sentinel/tasks/${selectedTask.value.id}/assign`, {
      method: 'POST',
      body: { dev_id: assignDevId.value }
    })

    // Update local state
    const taskIndex = tasks.value.findIndex(t => t.id === selectedTask.value?.id)
    if (taskIndex !== -1) {
      tasks.value[taskIndex].assigned_to = assignDevId.value
      tasks.value[taskIndex].status = 'IN_PROGRESS'
    }

    // Show success
    successMessage.value = `Task ${isReassign.value ? 're-assigned' : 'assigned'} to Dev #${assignDevId.value}`
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)

    closeAssignModal()
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to assign task'
    console.error('Failed to assign task:', err)
  } finally {
    isAssigning.value = false
  }
}

const getDevWorkload = (devTasks: Task[]) => {
  const totalMinutes = devTasks.reduce((sum, t) => sum + (t.ai_estimated_minutes || 0), 0)
  return (totalMinutes / 60).toFixed(1)
}

const getStatusDot = (status: string) => {
  const classes: Record<string, string> = {
    COMPLETED: 'bg-green-500',
    IN_PROGRESS: 'bg-blue-500',
    PENDING: 'bg-yellow-500',
    BLOCKED: 'bg-red-500',
    REVIEW_PENDING: 'bg-indigo-700' // 🚦 Quality Gate
  }
  return classes[status] || 'bg-gray-500'
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

// Navigation helper
const goToTask = (taskId: string) => {
  navigateTo(`/task/${taskId}`)
}

onMounted(() => {
  fetchTasks()
})
</script>
