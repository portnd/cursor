<template>
  <section v-if="isLoading" class="rounded-2xl border border-gray-700 bg-gray-800/60 p-6">
    <div class="flex items-center gap-3 mb-4">
      <div class="h-4 w-32 animate-pulse rounded bg-gray-700"/>
      <div class="h-5 w-8 animate-pulse rounded-full bg-gray-700 ml-auto"/>
    </div>
    <div class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-28 animate-pulse rounded-xl bg-gray-700/60"/>
    </div>
  </section>

  <section v-else-if="handoverTasks.length === 0">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-purple-400 mb-0.5">Handover Queue</h2>
        <p class="text-xs text-gray-500">Tasks awaiting your approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-gray-700/60 border border-gray-600 text-gray-500">
        0 pending
      </span>
    </div>
    <div class="rounded-2xl border border-gray-700/50 bg-gray-800/30 px-6 py-10 text-center">
      <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500/10 border border-emerald-500/20">
        <svg class="h-5 w-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
      </div>
      <p class="text-sm font-semibold text-gray-400">All clear</p>
      <p class="text-xs text-gray-600 mt-1">No tasks pending review</p>
    </div>
  </section>

  <section v-else>
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-purple-400 mb-0.5">Handover Queue</h2>
        <p class="text-xs text-gray-500">Tasks awaiting your approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/25 text-purple-400">
        {{ handoverTasks.length }} pending
      </span>
    </div>

    <div class="space-y-4">
      <div
        v-for="task in handoverTasks"
        :key="task.id"
        class="rounded-2xl border-2 border-purple-500/30 bg-gray-800/70 p-5 shadow-lg"
      >
        <!-- Task Header -->
        <div class="flex items-start justify-between gap-3 mb-3">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <span v-if="task.code" class="shrink-0 text-xs font-mono font-bold px-2 py-0.5 rounded bg-gray-700 text-gray-400 border border-gray-600">
                {{ task.code }}
              </span>
              <span class="px-2 py-0.5 text-xs font-semibold rounded-full bg-purple-500/15 border border-purple-500/30 text-purple-300">
                REVIEW PENDING
              </span>
            </div>
            <h3 class="font-bold text-white text-sm leading-tight">{{ task.title }}</h3>
            <p v-if="task.assigned_to_display_name || task.assigned_to_email" class="text-xs text-gray-500 mt-0.5">
              by {{ task.assigned_to_display_name || task.assigned_to_email }}
            </p>
          </div>
          <div v-if="task.due_at" class="shrink-0 text-right text-xs">
            <p :class="getDeadlineClass(task.due_at)">{{ formatDate(task.due_at) }}</p>
            <p class="text-gray-600 mt-0.5">{{ getCountdown(task.due_at) }}</p>
          </div>
        </div>

        <!-- Submission Reference -->
        <template v-if="latestSubmission(task)">
          <div class="rounded-xl border border-gray-600/60 bg-gray-900/50 p-3 mb-4 space-y-2">
            <div class="flex items-center gap-2">
              <svg class="h-3.5 w-3.5 text-purple-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
              </svg>
              <a
                :href="latestSubmission(task)!.reference_url"
                target="_blank"
                rel="noopener noreferrer"
                class="text-xs text-purple-300 hover:text-purple-200 underline underline-offset-2 truncate transition-colors"
              >
                {{ latestSubmission(task)!.reference_url }}
              </a>
            </div>
            <p v-if="latestSubmission(task)!.note" class="text-xs text-gray-400 pl-5 italic">
              "{{ latestSubmission(task)!.note }}"
            </p>
            <p class="text-xs text-gray-600 pl-5">
              Submitted {{ formatRelative(latestSubmission(task)!.created_at) }}
            </p>
          </div>
        </template>
        <div v-else class="rounded-xl border border-gray-700 bg-gray-900/30 px-3 py-2 mb-4 text-xs text-gray-500">
          No submission reference attached
        </div>

        <!-- Reject reason textarea (shown when rejecting) -->
        <div v-if="rejectingTaskId === task.id" class="mb-3">
          <label class="block text-xs font-medium text-red-400 mb-1.5">
            Rejection reason <span class="text-red-400">*</span>
            <span class="text-gray-600 font-normal ml-1">(min 10 chars)</span>
          </label>
          <textarea
            v-model="rejectReason"
            rows="3"
            placeholder="Describe what the developer needs to fix or improve…"
            class="w-full rounded-lg border border-red-500/40 bg-gray-900/60 px-3 py-2.5 text-sm text-white placeholder-gray-500 transition-colors focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500/50 resize-none"
          />
        </div>

        <!-- Action Buttons -->
        <div v-if="rejectingTaskId !== task.id" class="flex gap-3">
          <button
            @click="approve(task)"
            :disabled="actionLoading === task.id"
            class="flex-1 inline-flex items-center justify-center gap-2 rounded-xl bg-emerald-100 dark:bg-emerald-600 hover:bg-emerald-100 dark:bg-emerald-500 px-4 py-3 text-sm font-bold text-gray-900 dark:text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-emerald-900/20"
          >
            <svg v-if="actionLoading === task.id" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <svg v-else class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
            </svg>
            APPROVE — Mark Completed
          </button>
          <button
            @click="startReject(task.id)"
            :disabled="actionLoading === task.id"
            class="flex-1 inline-flex items-center justify-center gap-2 rounded-xl bg-red-100 dark:bg-red-600/20 hover:bg-red-100 dark:bg-red-600 border border-red-300 dark:border-red-500/40 hover:border-red-300 dark:border-red-500 px-4 py-3 text-sm font-bold text-red-300 hover:text-gray-900 dark:text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            REJECT — Send Back
          </button>
        </div>
        <div v-else class="flex gap-3">
          <button
            @click="cancelReject"
            class="flex-1 rounded-xl border border-gray-600 px-4 py-3 text-sm font-semibold text-gray-400 hover:bg-gray-700 hover:text-gray-900 dark:text-white transition-colors"
          >
            Cancel
          </button>
          <button
            @click="confirmReject(task)"
            :disabled="rejectReason.trim().length < 10 || actionLoading === task.id"
            class="flex-1 inline-flex items-center justify-center gap-2 rounded-xl bg-red-100 dark:bg-red-600 hover:bg-red-100 dark:bg-red-500 px-4 py-3 text-sm font-bold text-gray-900 dark:text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-red-900/20"
          >
            <svg v-if="actionLoading === task.id" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            Confirm Reject
          </button>
        </div>

        <!-- Error message -->
        <p v-if="errorTaskId === task.id" class="mt-2 text-xs text-red-400">{{ errorMessage }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
interface Submission {
  id: string
  reference_url: string
  note: string
  created_at: string
}

interface Task {
  id: string
  code?: string
  title: string
  status: string
  assigned_to?: number
  assigned_to_display_name?: string
  assigned_to_email?: string
  due_at?: string
  submissions?: Submission[]
}

const { fetchWithAuth } = useAuth()

const isLoading = ref(true)
const tasks = ref<Task[]>([])
const actionLoading = ref<string | null>(null)
const rejectingTaskId = ref<string | null>(null)
const rejectReason = ref('')
const errorTaskId = ref<string | null>(null)
const errorMessage = ref('')

const handoverTasks = computed(() =>
  tasks.value.filter(t => t.status === 'REVIEW_PENDING')
)

const latestSubmission = (task: Task): Submission | undefined =>
  task.submissions?.[0]

const fetchTasks = async () => {
  isLoading.value = true
  try {
    const res = await fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/approvals')
    tasks.value = res.data || []
  } catch {
    tasks.value = []
  } finally {
    isLoading.value = false
  }
}

const approve = async (task: Task) => {
  actionLoading.value = task.id
  errorTaskId.value = null
  try {
    await fetchWithAuth(`/sentinel/tasks/${task.id}/approve`, { method: 'POST' })
    tasks.value = tasks.value.filter(t => t.id !== task.id)
  } catch (err: any) {
    errorTaskId.value = task.id
    errorMessage.value = err.data?.message || err.message || 'Failed to approve task'
  } finally {
    actionLoading.value = null
  }
}

const startReject = (taskId: string) => {
  rejectingTaskId.value = taskId
  rejectReason.value = ''
  errorTaskId.value = null
}

const cancelReject = () => {
  rejectingTaskId.value = null
  rejectReason.value = ''
}

const confirmReject = async (task: Task) => {
  if (rejectReason.value.trim().length < 10) return
  actionLoading.value = task.id
  errorTaskId.value = null
  try {
    await fetchWithAuth(`/sentinel/tasks/${task.id}/reject`, {
      method: 'POST',
      body: { reason: rejectReason.value.trim() },
    })
    tasks.value = tasks.value.filter(t => t.id !== task.id)
    rejectingTaskId.value = null
    rejectReason.value = ''
  } catch (err: any) {
    errorTaskId.value = task.id
    errorMessage.value = err.data?.message || err.message || 'Failed to reject task'
  } finally {
    actionLoading.value = null
  }
}

const formatDate = (dateStr: string): string =>
  new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })

const formatRelative = (dateStr: string): string => {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  return `${Math.floor(hours / 24)}d ago`
}

const getCountdown = (dateStr: string): string => {
  const diff = new Date(dateStr).getTime() - Date.now()
  if (diff < 0) {
    const h = Math.abs(Math.floor(diff / 3600000))
    return h >= 24 ? `${Math.floor(h / 24)}d overdue` : `${h}h overdue`
  }
  const h = Math.floor(diff / 3600000)
  return h >= 24 ? `${Math.floor(h / 24)}d left` : `${h}h left`
}

const getDeadlineClass = (dateStr: string): string => {
  const diff = new Date(dateStr).getTime() - Date.now()
  if (diff < 0) return 'text-red-400 font-semibold'
  if (diff < 86400000) return 'text-amber-400 font-semibold'
  return 'text-gray-500'
}

onMounted(fetchTasks)
</script>
