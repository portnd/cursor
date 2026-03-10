<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-white">YOUR ASSIGNMENTS & APPROVALS</h1>
          <p class="text-sm text-gray-400 mt-1">Centralized mission control & decision center</p>
        </div>
        <div v-if="authStore.user" class="text-right">
          <div class="text-xs text-gray-500 uppercase">Operator</div>
          <div class="text-sm text-white font-medium">{{ authStore.user.email }}</div>
          <div 
            :class="[
              'text-xs font-bold px-2 py-0.5 rounded mt-1 inline-block',
              authStore.user.role === 'CEO' ? 'bg-gray-700 text-gray-200' :
              authStore.user.role === 'PM' ? 'bg-blue-700 text-blue-100' :
              'bg-green-700 text-green-100'
            ]"
          >
            {{ authStore.user.role }}
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลด inbox...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-900/20 border-2 border-red-500 rounded-lg p-6 text-red-400">
      <div class="flex items-center gap-3 mb-2">
        <span class="text-2xl">❌</span>
        <span class="font-bold">Error Loading Inbox</span>
      </div>
      <p class="text-sm">{{ error }}</p>
      <button 
        @click="fetchData"
        class="mt-4 px-4 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded transition-colors"
      >
        Retry
      </button>
    </div>

    <!-- Content -->
    <div v-else class="space-y-8">
      <!-- SECTION 1: ⏱️ PENDING TIME NEGOTIATIONS (CEO/PM ONLY) -->
      <div 
        v-if="showApprovals && timeNegotiations.length > 0"
        class="bg-gray-800 border border-gray-700 rounded p-6 mb-8"
      >
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-lg font-bold text-white uppercase flex items-center gap-2">
              <span>⏱️</span>
              <span>Pending Time Negotiations</span>
            </h2>
            <p class="text-sm text-gray-400 mt-1">Developers requesting more time</p>
          </div>
          <div class="px-3 py-1 bg-yellow-700 text-yellow-100 font-bold rounded text-sm">
            {{ timeNegotiations.length }}
          </div>
        </div>

        <div class="space-y-3">
          <div
            v-for="task in timeNegotiations"
            :key="task.id"
            class="bg-gray-900 border border-gray-700 hover:border-gray-600 rounded p-4 transition-all cursor-pointer"
            @click="goToTask(task)"
          >
            <div class="flex items-start justify-between gap-4 mb-3">
              <div class="flex-1">
                <h3 class="text-base font-bold text-white mb-1">{{ task.title }}</h3>
                <p class="text-xs text-gray-500">ID: {{ taskCodeDisplay(task) }}</p>
              </div>
              
              <!-- Time Negotiation Badge -->
              <div class="px-2 py-1 bg-yellow-700 text-yellow-100 text-xs font-bold rounded flex items-center gap-1">
                <span>⏱️</span>
                <span>TIME</span>
              </div>
            </div>

            <!-- Task Info -->
            <div class="mb-4">
              <p class="text-sm text-gray-400 mb-2">{{ task.description || 'No description' }}</p>
              
              <!-- Time Negotiation Details -->
              <div class="bg-gray-950 border border-gray-700 rounded p-3">
                <div class="space-y-2">
                  <div class="flex items-center gap-4 text-xs">
                    <div>
                      <span class="text-gray-500">Estimate:</span>
                      <span class="text-yellow-400 font-bold ml-1">{{ task.estimated_minutes }} min</span>
                    </div>
                    <div>
                      <span class="text-gray-500">→ Dev Proposes:</span>
                      <span class="text-yellow-400 font-bold ml-1">{{ task.proposed_minutes }} min</span>
                    </div>
                  </div>
                  <p class="text-xs text-gray-400">
                    <span class="text-gray-500">Reason:</span> {{ task.negotiation_reason }}
                  </p>
                  
                  <!-- AI Advisory -->
                  <div class="pt-2 mt-2 border-t border-gray-700 space-y-2">
                    <div class="flex items-center gap-4 text-xs">
                      <div>
                        <span class="text-gray-500">AI Suggests:</span>
                        <span 
                          v-if="task.negotiation_ai_recommendation"
                          :class="[
                            'font-bold ml-1',
                            task.negotiation_ai_recommendation === 'APPROVE' ? 'text-green-400' : 'text-red-400'
                          ]"
                        >
                          {{ task.negotiation_ai_recommendation }} ({{ task.negotiation_ai_confidence }}%)
                        </span>
                        <span v-else class="text-gray-500 ml-1 italic">
                          (ไม่มีข้อมูล)
                        </span>
                      </div>
                    </div>
                    
                    <!-- AI Reasoning -->
                    <div class="bg-gray-900 border border-gray-700 rounded p-2">
                      <p class="text-xs"
                         :class="task.negotiation_ai_reasoning && task.negotiation_ai_reasoning.trim()
                           ? 'text-gray-300'
                           : 'text-gray-500 italic'">
                        <span class="font-bold">💡 ความเห็น AI:</span> 
                        <span v-if="task.negotiation_ai_reasoning && task.negotiation_ai_reasoning.trim()">
                          "{{ task.negotiation_ai_reasoning }}"
                        </span>
                        <span v-else>
                          ⚠️ ไม่มีข้อมูล - กรุณาใช้ดุลยพินิจในการตัดสิน
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex items-center gap-3">
              <button
                @click.stop="approveNegotiation(task.code || task.id, task.proposed_minutes || 0)"
                class="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-bold rounded transition-all flex items-center justify-center gap-2"
              >
                <span>✅</span>
                <span>Approve Time</span>
              </button>
              <button
                @click.stop="rejectNegotiation(task.code || task.id)"
                class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-bold rounded transition-all flex items-center justify-center gap-2"
              >
                <span>❌</span>
                <span>Reject</span>
              </button>
              <NuxtLink 
                :to="`/task/${task.code || task.id}`"
                @click.stop
                class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm font-medium rounded transition-all"
              >
                View Details →
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 3: ⚡ MY ACTIVE MISSIONS (FOR EVERYONE) -->
      <div v-if="myTasks.length > 0" class="bg-gray-800 border border-gray-700 rounded p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-lg font-bold text-white uppercase flex items-center gap-2">
              <span>⚡</span>
              <span>My Active Missions</span>
            </h2>
            <p class="text-sm text-gray-400 mt-1">Tasks currently assigned to you</p>
          </div>
          <div class="px-3 py-1 bg-blue-700 text-blue-100 font-bold rounded text-sm">
            {{ myTasks.length }}
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="task in myTasks"
            :key="task.id"
            :class="[
              'bg-gray-900 border rounded p-4 transition-all relative cursor-pointer hover:bg-gray-800',
              getDeadlineBorderClass(task)
            ]"
            @click="goToTask(task)"
          >
            <!-- Urgency Badge -->
            <div 
              v-if="getDeadlineUrgency(task) === 'overdue'"
              class="absolute -top-2 -right-2 px-2 py-1 bg-red-700 text-red-100 text-xs font-bold rounded"
            >
              🚨 OVERDUE
            </div>
            <div 
              v-else-if="getDeadlineUrgency(task) === 'urgent'"
              class="absolute -top-2 -right-2 px-2 py-1 bg-yellow-700 text-yellow-100 text-xs font-bold rounded"
            >
              ⚠️ URGENT
            </div>

            <!-- Task Header -->
            <div class="flex items-start justify-between mb-3">
              <h3 class="text-lg font-bold text-white flex-1">{{ task.title }}</h3>
              <span 
                :class="[
                  'px-2 py-1 text-xs font-bold rounded',
                  task.status === 'COMPLETED' ? 'bg-green-600 text-white' :
                  task.status === 'IN_PROGRESS' ? 'bg-blue-600 text-white' :
                  'bg-gray-600 text-white'
                ]"
              >
                {{ task.status }}
              </span>
            </div>

            <!-- Task Description -->
            <p class="text-sm text-gray-400 mb-4">{{ task.description || 'No description' }}</p>

            <!-- Task Metrics -->
            <div class="space-y-2 mb-4 text-xs">
              <div class="flex items-center gap-2">
                <span class="text-gray-500">Estimated:</span>
                <span class="text-blue-400 font-bold">
                  {{ task.estimated_minutes }} min ({{ (task.estimated_minutes / 60).toFixed(1) }}h)
                </span>
              </div>
              <div v-if="task.due_at" class="flex items-center gap-2">
                <span 
                  :class="[
                    'font-bold',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-400' :
                    'text-gray-400'
                  ]"
                >
                  ⏰ Due: {{ formatDeadline(task.due_at) }}
                </span>
                <span 
                  :class="[
                    'px-2 py-0.5 font-bold rounded',
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-900/50 text-red-300' :
                    getDeadlineUrgency(task) === 'urgent' ? 'bg-yellow-900/50 text-yellow-300' :
                    'bg-gray-700 text-gray-300'
                  ]"
                >
                  {{ getDeadlineCountdown(task.due_at) }}
                </span>
              </div>
            </div>

            <!-- Action Button -->
            <NuxtLink 
              :to="`/task/${task.code || task.id}`"
              @click.stop
              class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-bold rounded transition-colors"
            >
              <span>Execute</span>
              <span>→</span>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- EMPTY STATE -->
      <div 
        v-if="!isLoading && !error && timeNegotiations.length === 0 && myTasks.length === 0"
        class="text-center py-20"
      >
        <div class="text-6xl mb-4">✨</div>
        <h2 class="text-2xl font-bold text-gray-400 mb-2">All Systems Clear</h2>
        <p class="text-gray-500">No pending actions. You're all caught up!</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

// Page Meta
definePageMeta({
  middleware: 'auth',
  layout: 'default'
})

// Types
interface Task {
  id: string
  code?: string
  title: string
  description: string
  estimated_minutes: number
  proposed_minutes?: number
  negotiation_status: string
  negotiation_reason: string
  negotiation_ai_recommendation?: string
  negotiation_ai_confidence?: number
  negotiation_ai_reasoning?: string
  due_at: string | null
  started_at: string | null
  completed_at: string | null
  status: string
  assigned_to: number
  submissions?: Submission[]
}

interface Submission {
  id: string
  reference_url: string
  note: string
  created_at: string
}

// Composables
const authStore = useAuthStore()
const { fetchWithAuth, currentUser } = useAuth()
const { showError, confirm } = useNotification()

// State
const isLoading = ref(true)
const error = ref('')
const myTasks = ref<Task[]>([])
const approvals = ref<Task[]>([])

// Computed
const showApprovals = computed(() => {
  return currentUser.value?.role === 'CEO' || currentUser.value?.role === 'PM'
})

// Separate time negotiations from other approvals
const timeNegotiations = computed(() => {
  return approvals.value.filter(task => task.negotiation_status === 'PENDING')
})

// Methods
const fetchData = async () => {
  isLoading.value = true
  error.value = ''

  try {
    // 1. Fetch My Tasks (for everyone)
    const myTasksResponse = await fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/my')
    myTasks.value = myTasksResponse.data || []

    // 2. Fetch Approvals (only for CEO/PM)
    if (showApprovals.value) {
      try {
        const approvalsResponse = await fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/approvals')
        approvals.value = approvalsResponse.data || []
      } catch (approvalsError: any) {
        // If 403, it's expected for non-CEO/PM users, just skip
        if (approvalsError.statusCode !== 403) {
          console.error('Failed to fetch approvals:', approvalsError)
        }
      }
    }
  } catch (err: any) {
    console.error('Failed to fetch tasks:', err)
    error.value = err.data?.message || err.message || 'Failed to load inbox'
  } finally {
    isLoading.value = false
  }
}

// Time Negotiation Actions
const approveNegotiation = async (taskId: string, minutes: number) => {
  const ok = await confirm({
    title: 'Approve time negotiation',
    message: `Approve time negotiation to ${minutes} minutes?`,
    confirmLabel: 'Approve',
    cancelLabel: 'Cancel',
    variant: 'primary'
  })
  if (!ok) return

  try {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/negotiate/resolve`, {
      method: 'POST',
      body: JSON.stringify({
        approved: true,
        approved_minutes: minutes
      })
    })
    
    await fetchData()
  } catch (err: any) {
    console.error('Failed to approve negotiation:', err)
    showError(err.data?.message || 'Failed to approve negotiation')
  }
}

const rejectNegotiation = async (taskId: string) => {
  const reason = prompt('Reason for rejection (optional):')
  if (reason === null) return // User cancelled

  try {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/negotiate/resolve`, {
      method: 'POST',
      body: JSON.stringify({
        approved: false,
        rejection_reason: reason || 'No reason provided'
      })
    })
    
    // Refresh data
    await fetchData()
  } catch (err: any) {
    console.error('Failed to reject negotiation:', err)
    showError(err.data?.message || 'Failed to reject negotiation')
  }
}

// Deadline Utilities
const getDeadlineUrgency = (task: Task): 'normal' | 'urgent' | 'overdue' => {
  if (!task.due_at || task.status === 'COMPLETED') return 'normal'
  
  const now = new Date()
  const deadline = new Date(task.due_at)
  const hoursLeft = (deadline.getTime() - now.getTime()) / (1000 * 60 * 60)
  
  if (hoursLeft < 0) return 'overdue'
  if (hoursLeft < 24) return 'urgent'
  return 'normal'
}

const getDeadlineBorderClass = (task: Task): string => {
  const urgency = getDeadlineUrgency(task)
  if (urgency === 'overdue') return 'border-red-700'
  if (urgency === 'urgent') return 'border-yellow-700'
  return 'border-gray-700'
}

const formatDeadline = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric', 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getDeadlineCountdown = (dateStr: string): string => {
  const now = new Date()
  const deadline = new Date(dateStr)
  const diff = deadline.getTime() - now.getTime()
  
  if (diff < 0) {
    const hoursOverdue = Math.abs(Math.floor(diff / (1000 * 60 * 60)))
    const daysOverdue = Math.floor(hoursOverdue / 24)
    
    if (daysOverdue > 0) {
      return `${daysOverdue}d overdue`
    }
    return `${hoursOverdue}h overdue`
  }
  
  const hoursLeft = Math.floor(diff / (1000 * 60 * 60))
  const daysLeft = Math.floor(hoursLeft / 24)
  
  if (daysLeft > 0) {
    return `${daysLeft}d left`
  }
  return `${hoursLeft}h left`
}

function taskCodeDisplay(task: { id: string; code?: string }): string {
  if (task.code) {
    const suffix = task.code.split('-').pop()
    if (suffix && /^\d+$/.test(suffix)) return String(Number(suffix)).padStart(4, '0')
    return task.code
  }
  return task.id.substring(0, 8) + '…'
}

// Navigation helper
const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}`)
}

// Lifecycle
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* Optional: Add any custom animations or styles */
</style>
