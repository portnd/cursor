<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-white">DEV EXECUTION DASHBOARD</h1>
      <p class="text-sm text-gray-400 mt-1">Deep work & task pipeline</p>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="text-center py-20">
      <div class="text-gray-400">Loading your tasks...</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded p-4 text-red-400 mb-6">
      {{ error }}
    </div>

    <!-- Content -->
    <div v-else>
      <!-- TOP: Current Focus (IN_PROGRESS) -->
      <div v-if="currentFocus.length > 0" class="mb-8">
        <h2 class="text-lg font-bold text-blue-400 uppercase mb-3">⚡ Current Focus</h2>
        
        <div class="grid grid-cols-1 gap-4">
          <div
            v-for="task in currentFocus"
            :key="task.id"
            :class="[
              'bg-blue-900/20 border-2 rounded-lg p-6 relative cursor-pointer hover:scale-[1.01] transition-all',
              getDeadlineBorderClass(task)
            ]"
            @click="goToTask(task)"
          >
            <!-- Urgency Badge (Top Right) -->
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

            <div class="flex items-start justify-between mb-3">
              <h3 class="text-xl font-bold text-white flex-1">{{ task.title }}</h3>
              <span class="px-3 py-1 bg-blue-600 text-blue-100 text-xs font-bold rounded">
                IN PROGRESS
              </span>
            </div>
            
            <p class="text-gray-400 text-sm mb-4">{{ task.description || 'No description' }}</p>
            
            <div class="flex items-center gap-6 text-sm flex-wrap">
              <div class="flex items-center gap-2">
                <span class="text-gray-500">Estimated:</span>
                <span class="text-blue-400 font-bold">
                  {{ task.ai_estimated_minutes }} min ({{ (task.ai_estimated_minutes / 60).toFixed(1) }}h)
                </span>
              </div>
              <div v-if="task.started_at" class="flex items-center gap-2">
                <span class="text-gray-500">Started:</span>
                <span class="text-gray-400">{{ formatDate(task.started_at) }}</span>
              </div>
              <div v-if="task.due_at" class="flex items-center gap-2">
                <span 
                  :class="[
                    'text-xs font-bold',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-400' :
                    'text-gray-400'
                  ]"
                >
                  ⏰ Due: {{ formatDeadline(task.due_at) }}
                </span>
                <span 
                  :class="[
                    'px-2 py-0.5 text-xs font-bold rounded',
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-900/50 text-red-300' :
                    getDeadlineUrgency(task) === 'urgent' ? 'bg-yellow-900/50 text-yellow-300' :
                    'bg-gray-700 text-gray-300'
                  ]"
                >
                  {{ getDeadlineCountdown(task.due_at) }}
                </span>
              </div>
            </div>

            <div class="mt-4 flex gap-2">
              <button 
                @click.stop="openSubmitModal(task)"
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded transition-colors"
              >
                Submit Work
              </button>
              <NuxtLink 
                :to="`/task/${task.code || task.id}`"
                @click.stop
                class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 text-sm font-medium rounded transition-colors"
              >
                View Details
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- MIDDLE: My Backlog (PENDING) -->
      <div v-if="myBacklog.length > 0" class="mb-8">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-yellow-400 uppercase">📋 My Backlog</h2>
          <span class="text-sm text-gray-500">{{ myBacklog.length }} pending</span>
        </div>

        <div class="bg-gray-800 border border-gray-700 rounded overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-gray-900 text-xs text-gray-400 uppercase">
              <tr>
                <th class="text-left p-3">Task</th>
                <th class="text-right p-3">Est. Hours</th>
                <th class="text-right p-3">Deadline</th>
                <th class="text-right p-3">Action</th>
              </tr>
            </thead>
            <tbody class="text-gray-300">
              <tr 
                v-for="task in myBacklog" 
                :key="task.id"
                :class="[
                  'border-t border-gray-700 hover:bg-gray-700/50 cursor-pointer',
                  getDeadlineUrgency(task) === 'overdue' ? 'bg-red-900/20' :
                  getDeadlineUrgency(task) === 'urgent' ? 'bg-yellow-900/20' : ''
                ]"
                @click="goToTask(task)"
              >
                <td class="p-3 font-medium">
                  {{ task.title }}
                  <span 
                    v-if="getDeadlineUrgency(task) === 'overdue'"
                    class="ml-2 px-2 py-0.5 bg-red-600 text-white text-xs font-bold rounded"
                  >
                    OVERDUE
                  </span>
                  <span 
                    v-else-if="getDeadlineUrgency(task) === 'urgent'"
                    class="ml-2 px-2 py-0.5 bg-yellow-500 text-black text-xs font-bold rounded"
                  >
                    URGENT
                  </span>
                </td>
                <td class="p-3 text-right">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}</td>
                <td 
                  :class="[
                    'p-3 text-right text-xs',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-bold' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-400 font-bold' :
                    'text-gray-500'
                  ]"
                >
                  <div v-if="task.due_at">
                    {{ formatDeadline(task.due_at) }}
                    <div class="text-xs mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</div>
                  </div>
                  <div v-else class="text-gray-600">No deadline</div>
                </td>
                <td class="p-3 text-right">
                  <button 
                    @click.stop="goToTask(task)"
                    class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium rounded transition-colors"
                  >
                    Start
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- BOTTOM: Available Missions (Unassigned) -->
      <div v-if="availableMissions.length > 0" class="mb-8">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-green-400 uppercase">🎯 Available Missions</h2>
          <span class="text-sm text-gray-500">{{ availableMissions.length }} unassigned</span>
        </div>

        <div class="bg-gray-800 border border-gray-700 rounded overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-gray-900 text-xs text-gray-400 uppercase">
              <tr>
                <th class="text-left p-3">Task</th>
                <th class="text-left p-3">Description</th>
                <th class="text-right p-3">Est. Hours</th>
                <th class="text-right p-3">Deadline</th>
                <th class="text-right p-3">Action</th>
              </tr>
            </thead>
            <tbody class="text-gray-300">
              <tr 
                v-for="task in availableMissions" 
                :key="task.id"
                :class="[
                  'border-t border-gray-700 hover:bg-gray-700/50 cursor-pointer',
                  getDeadlineUrgency(task) === 'overdue' ? 'bg-red-900/20' :
                  getDeadlineUrgency(task) === 'urgent' ? 'bg-yellow-900/20' : ''
                ]"
                @click="goToTask(task)"
              >
                <td class="p-3 font-medium">
                  {{ task.title }}
                  <span 
                    v-if="getDeadlineUrgency(task) === 'overdue'"
                    class="ml-2 px-2 py-0.5 bg-red-600 text-white text-xs font-bold rounded"
                  >
                    OVERDUE
                  </span>
                  <span 
                    v-else-if="getDeadlineUrgency(task) === 'urgent'"
                    class="ml-2 px-2 py-0.5 bg-yellow-500 text-black text-xs font-bold rounded"
                  >
                    URGENT
                  </span>
                </td>
                <td class="p-3 text-gray-500 text-xs truncate max-w-xs">
                  {{ task.description || 'No description' }}
                </td>
                <td class="p-3 text-right">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}</td>
                <td 
                  :class="[
                    'p-3 text-right text-xs',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-bold' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-400 font-bold' :
                    'text-gray-500'
                  ]"
                >
                  <div v-if="task.due_at">
                    {{ formatDeadline(task.due_at) }}
                    <div class="text-xs mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</div>
                  </div>
                  <div v-else class="text-gray-600">No deadline</div>
                </td>
                <td class="p-3 text-right">
                  <button
                    @click.stop="claimTask(task.id)"
                    :disabled="claiming === task.id"
                    class="px-3 py-1 bg-green-600 hover:bg-green-700 disabled:bg-gray-700 text-white text-xs font-bold rounded transition-colors disabled:cursor-not-allowed"
                  >
                    {{ claiming === task.id ? 'Claiming...' : '✋ Claim' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="currentFocus.length === 0 && myBacklog.length === 0 && availableMissions.length === 0" class="bg-gray-800 border border-gray-700 rounded p-12 text-center">
        <div class="text-gray-500 mb-4">No tasks available</div>
        <NuxtLink
          to="/create"
          class="inline-block px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded transition-colors"
        >
          Create New Task
        </NuxtLink>
      </div>
    </div>

    <!-- Success Notification -->
    <div
      v-if="showSuccess"
      class="fixed top-4 right-4 bg-green-700 border border-green-600 text-white px-4 py-3 rounded shadow-lg z-50"
    >
      {{ successMessage }}
    </div>

    <!-- Submit Work Modal -->
    <div
      v-if="showSubmitModal"
      class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
      @click.self="closeSubmitModal"
    >
      <div class="bg-gray-800 border-2 border-blue-500 rounded-lg p-6 max-w-3xl w-full max-h-[90vh] overflow-y-auto relative">
        <!-- Close Button (X) -->
        <button
          @click="closeSubmitModal"
          class="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors"
          title="Close"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>

        <h3 class="text-2xl font-bold text-white mb-2">
          🚀 Submit Mission
        </h3>
        <p class="text-gray-400 text-sm mb-6 line-clamp-2">
          {{ selectedTask?.title }}
        </p>

        <!-- Error Message -->
        <div v-if="error" class="mb-4 p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
          {{ error }}
        </div>

        <!-- Form -->
        <div class="space-y-4 mb-6">
          <!-- Commit Hash -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Commit Hash <span class="text-red-400">*</span>
            </label>
            <input
              v-model="submitForm.commitHash"
              type="text"
              placeholder="e.g., a3f5d9c7b2e8..."
              class="w-full px-4 py-3 bg-gray-900 border border-gray-600 rounded text-white text-sm font-mono focus:outline-none focus:border-blue-500 transition-colors"
            />
          </div>

          <!-- Code Diff -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Code Diff / Changes <span class="text-red-400">*</span>
            </label>
            <textarea
              v-model="submitForm.diff"
              placeholder="Paste your code changes here...

Example:
+ function newFeature() {
+   return 'implemented';
+ }
- // old code removed"
              rows="12"
              class="w-full px-4 py-3 bg-gray-900 border border-gray-600 rounded text-white text-sm font-mono focus:outline-none focus:border-blue-500 transition-colors resize-none"
              style="min-height: 250px"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">
              💡 Tip: Include context lines for better AI analysis
            </p>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button
            @click="submitWork"
            :disabled="isSubmitting || !submitForm.commitHash || !submitForm.diff"
            class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-bold rounded transition-colors disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="isSubmitting">⚡ Analyzing...</span>
            <span v-else>🚀 Submit to AI Auditor</span>
          </button>
          <button
            @click="closeSubmitModal"
            :disabled="isSubmitting"
            class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 text-gray-300 rounded transition-colors disabled:cursor-not-allowed"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Result Modal (The Verdict) -->
    <div
      v-if="showResultModal"
      class="fixed inset-0 bg-black/90 flex items-center justify-center z-50 p-4"
      @click.self="closeResultModal"
    >
      <div 
        :class="[
          'rounded p-8 max-w-2xl w-full border max-h-[90vh] overflow-y-auto relative',
          submissionResult?.ai_verdict === 'PASS' 
            ? 'bg-gray-800 border-green-700' 
            : 'bg-gray-800 border-red-700'
        ]"
      >
        <!-- Close Button (X) -->
        <button
          @click="closeResultModal"
          :class="[
            'absolute top-4 right-4 transition-colors',
            submissionResult?.ai_verdict === 'PASS' 
              ? 'text-green-300 hover:text-white' 
              : 'text-red-300 hover:text-white'
          ]"
          title="Close"
        >
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>

        <!-- Verdict Header -->
        <div class="text-center mb-8">
          <div 
            :class="[
              'text-6xl font-black mb-4',
              submissionResult?.ai_verdict === 'PASS' ? 'text-green-400' : 'text-red-400'
            ]"
          >
            {{ submissionResult?.ai_verdict === 'PASS' ? '✅ MISSION' : '❌ MISSION' }}
          </div>
          <div 
            :class="[
              'text-5xl font-black',
              submissionResult?.ai_verdict === 'PASS' ? 'text-green-100' : 'text-red-100'
            ]"
          >
            {{ submissionResult?.ai_verdict === 'PASS' ? 'ACCOMPLISHED' : 'FAILED' }}
          </div>
        </div>

        <!-- Score -->
        <div class="text-center mb-8">
          <div class="text-gray-300 text-sm uppercase tracking-wide mb-2">AI Auditor Score</div>
          <div 
            :class="[
              'text-7xl font-black',
              submissionResult?.ai_verdict === 'PASS' ? 'text-green-200' : 'text-red-200'
            ]"
          >
            {{ submissionResult?.ai_score }}<span class="text-4xl">/100</span>
          </div>
        </div>

        <!-- Feedback -->
        <div 
          :class="[
            'rounded-lg p-6 mb-8 border-2',
            submissionResult?.ai_verdict === 'PASS' 
              ? 'bg-green-950/50 border-green-600' 
              : 'bg-red-950/50 border-red-600'
          ]"
        >
          <div class="text-white font-bold mb-3 flex items-center gap-2">
            <span>🤖</span>
            <span>AI Auditor Feedback:</span>
          </div>
          <div class="text-gray-200 text-sm whitespace-pre-wrap leading-relaxed">
            {{ getFeedbackText() }}
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="space-y-3">
          <!-- Primary Action Button -->
          <button
            @click="closeResultModal"
            :class="[
              'w-full px-6 py-4 font-bold rounded-lg transition-all text-white text-lg',
              submissionResult?.ai_verdict === 'PASS'
                ? 'bg-green-600 hover:bg-green-700'
                : 'bg-red-600 hover:bg-red-700'
            ]"
          >
            {{ submissionResult?.ai_verdict === 'PASS' ? '🎉 Continue Working' : '📢 Go to Task & Appeal' }}
          </button>
          
          <!-- Additional Info for FAIL/PENDING -->
          <div v-if="submissionResult?.ai_verdict !== 'PASS'" class="text-center text-gray-400 text-sm">
            You will be taken to the task page where you can:
            <div class="mt-1 space-y-1">
              <div>• 📢 <span class="text-amber-400">Appeal the AI verdict</span></div>
              <div>• 🔍 Review detailed feedback</div>
              <div>• 🔧 Submit revised code</div>
            </div>
          </div>
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
  due_at?: string
  started_at?: string
  completed_at?: string
}

interface Submission {
  id: string
  task_id: string
  dev_id: number
  commit_hash: string
  ai_verdict: string
  ai_score: number
  ai_feedback: any
  created_at: string
}

const { fetchWithAuth, currentUser } = useAuth()

const myTasks = ref<Task[]>([])
const unassignedTasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')
const claiming = ref<string | null>(null)

// Success State
const showSuccess = ref(false)
const successMessage = ref('')

// Submission State
const showSubmitModal = ref(false)
const showResultModal = ref(false)
const isSubmitting = ref(false)
const selectedTask = ref<Task | null>(null)
const submitForm = ref({
  commitHash: '',
  diff: ''
})
const submissionResult = ref<Submission | null>(null)

// Computed: Split My Tasks into Focus and Backlog (Sorted by due_date)
const currentFocus = computed(() => 
  myTasks.value
    .filter(t => t.status === 'IN_PROGRESS')
    .sort((a, b) => {
      // Sort by due_date (urgent first)
      if (!a.due_at) return 1
      if (!b.due_at) return -1
      return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
    })
)

const myBacklog = computed(() => 
  myTasks.value
    .filter(t => t.status === 'PENDING')
    .sort((a, b) => {
      // Sort by due_date (urgent first)
      if (!a.due_at) return 1
      if (!b.due_at) return -1
      return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
    })
)

const availableMissions = computed(() => 
  unassignedTasks.value.sort((a, b) => {
    // Sort by due_date (urgent first)
    if (!a.due_at) return 1
    if (!b.due_at) return -1
    return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
  })
)

// Methods
const fetchTasks = async () => {
  try {
    isLoading.value = true
    
    const [myTasksResponse, unassignedResponse] = await Promise.all([
      fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/my'),
      fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/unassigned')
    ])
    
    myTasks.value = myTasksResponse.data || []
    unassignedTasks.value = unassignedResponse.data || []
    error.value = ''
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load tasks'
    console.error('Failed to fetch tasks:', err)
  } finally {
    isLoading.value = false
  }
}

const claimTask = async (taskId: string) => {
  if (!currentUser.value?.user_id) {
    error.value = 'User ID not found. Please log in again.'
    return
  }

  try {
    claiming.value = taskId
    
    await fetchWithAuth(`/sentinel/tasks/${taskId}/assign`, {
      method: 'POST',
      body: { dev_id: currentUser.value.user_id }
    })
    
    // Optimistic Update
    const claimedTask = unassignedTasks.value.find(t => t.id === taskId)
    if (claimedTask) {
      claimedTask.status = 'IN_PROGRESS'
      claimedTask.assigned_to = currentUser.value.user_id
      
      unassignedTasks.value = unassignedTasks.value.filter(t => t.id !== taskId)
      myTasks.value.unshift(claimedTask)
    }
    
    successMessage.value = '✅ Task claimed successfully'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)
    
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to claim task'
    console.error('Failed to claim task:', err)
  } finally {
    claiming.value = null
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric'
  })
}

// Deadline Urgency Helpers
const getDeadlineUrgency = (task: Task) => {
  if (!task.due_at || task.status === 'COMPLETED') return 'none'
  
  const now = new Date().getTime()
  const dueDate = new Date(task.due_at).getTime()
  const hoursUntilDue = (dueDate - now) / (1000 * 60 * 60)
  
  if (hoursUntilDue < 0) return 'overdue'
  if (hoursUntilDue < 24) return 'urgent'
  return 'normal'
}

const getDeadlineBorderClass = (task: Task) => {
  const urgency = getDeadlineUrgency(task)
  if (urgency === 'overdue') return 'border-red-500'
  if (urgency === 'urgent') return 'border-yellow-500'
  return 'border-blue-500'
}

const formatDeadline = (dueAt: string) => {
  const date = new Date(dueAt)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getDeadlineCountdown = (dueAt: string) => {
  const now = new Date().getTime()
  const due = new Date(dueAt).getTime()
  const diff = due - now
  
  if (diff < 0) {
    // Overdue
    const hours = Math.abs(Math.floor(diff / (1000 * 60 * 60)))
    const days = Math.floor(hours / 24)
    if (days > 0) return `Overdue by ${days}d`
    return `Overdue by ${hours}h`
  }
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `${days}d left`
  if (hours > 0) return `${hours}h left`
  return 'Due soon!'
}

// Submission Functions
const openSubmitModal = (task: Task) => {
  selectedTask.value = task
  submitForm.value = {
    commitHash: '',
    diff: ''
  }
  showSubmitModal.value = true
}

const closeSubmitModal = () => {
  if (!isSubmitting.value) {
    showSubmitModal.value = false
    selectedTask.value = null
    submitForm.value = {
      commitHash: '',
      diff: ''
    }
  }
}

const submitWork = async () => {
  if (!selectedTask.value || !submitForm.value.commitHash || !submitForm.value.diff) {
    return
  }

  try {
    isSubmitting.value = true
    error.value = ''

    const response = await fetchWithAuth<{ data: Submission }>(`/sentinel/tasks/${selectedTask.value.id}/submit`, {
      method: 'POST',
      body: {
        commit_hash: submitForm.value.commitHash,
        diff: submitForm.value.diff
      }
    })

    submissionResult.value = response.data
    showSubmitModal.value = false
    showResultModal.value = true

    // Refresh tasks after successful submission
    await fetchTasks()

  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to submit work'
    console.error('Failed to submit work:', err)
  } finally {
    isSubmitting.value = false
  }
}

const closeResultModal = async () => {
  // If NOT PASS (FAIL or PENDING), navigate to task detail page for appeal/review
  if (submissionResult.value?.ai_verdict !== 'PASS' && selectedTask.value?.id) {
    const task = selectedTask.value
    showResultModal.value = false
    submissionResult.value = null
    selectedTask.value = null
    await navigateTo(`/task/${task?.code || task?.id}`)
    return // Don't fetch tasks if navigating away
  }
  
  showResultModal.value = false
  submissionResult.value = null
  selectedTask.value = null
  
  // Refresh tasks when closing result modal
  fetchTasks()
}

const getFeedbackText = () => {
  if (!submissionResult.value?.ai_feedback) return 'No feedback available'
  
  try {
    const feedback = submissionResult.value.ai_feedback
    
    // Handle string feedback
    if (typeof feedback === 'string') {
      return feedback
    }
    
    // Handle error case (AI service unavailable)
    if (feedback.error) {
      const errorMsg = feedback.error
      
      // Check if it's a Gemini quota error
      if (errorMsg.includes('429') || errorMsg.includes('quota') || errorMsg.includes('RESOURCE_EXHAUSTED')) {
        return `⏳ AI Review Service Temporarily Unavailable

The AI code reviewer has reached its quota limit. Your submission has been recorded as PENDING.

What happens next:
• Your code submission is safely saved
• You can appeal this submission to request manual review
• Or wait and resubmit with a new commit hash
• The system will automatically retry the AI review later

Tip: For immediate feedback, try submitting again in about 30 seconds.`
      }
      
      // Generic error
      return `⚠️ AI Review Service Error

The AI reviewer encountered an issue and couldn't analyze your code. Your submission has been recorded as PENDING.

You can:
• Appeal this submission for manual review by CEO/PM
• Resubmit your code with a new commit hash
• Check the task detail page for more options

Your work is saved and won't be lost.`
    }
    
    // Handle normal feedback
    if (feedback.feedback) {
      return feedback.feedback
    }
    
    // Fallback: stringify the object
    return JSON.stringify(feedback, null, 2)
  } catch {
    return 'Unable to parse feedback'
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
