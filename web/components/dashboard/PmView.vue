<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <!-- Enterprise Header -->
    <header class="sticky top-0 z-10 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 py-6 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
              Resource Control
            </h1>
            <p class="mt-1 text-sm text-gray-400">
              Task allocation and workforce management
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-3">
            <!-- Inline stats -->
            <div class="flex items-center gap-4 rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-2.5">
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Total</span>
              <span class="text-lg font-bold text-white">{{ tasks.length }}</span>
              <span class="h-4 w-px bg-gray-600" />
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Unassigned</span>
              <span class="text-lg font-bold text-amber-400">{{ unassignedTasks.length }}</span>
              <span class="h-4 w-px bg-gray-600" />
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">In review</span>
              <span class="text-lg font-bold text-purple-400">{{ reviewPendingTasks.length }}</span>
            </div>
            <button
              type="button"
              @click="fetchTasks"
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
      <div class="h-10 w-10 animate-spin rounded-full border-2 border-purple-500 border-t-transparent" />
      <p class="mt-4 text-sm text-gray-500">Loading resource data...</p>
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

    <!-- Main Content -->
    <main v-else class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <!-- KPI Strip -->
      <section class="mb-8">
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Accountability KPIs
        </h2>
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <PerformanceKpiScoreCard
            label="Team Delivery Rate"
            :value="avgTeamDeliveryRate"
            format="pct"
            :status="deliveryStatus(avgTeamDeliveryRate)"
          />
          <PerformanceKpiScoreCard
            label="Team Code Quality"
            :value="avgTeamCodeQuality"
            :status="qualityStatus(avgTeamCodeQuality)"
          />
          <PerformanceKpiScoreCard
            label="Team Rework Rate"
            :value="avgTeamReworkRate"
            format="pct"
            :status="reworkStatus(avgTeamReworkRate)"
          />
          <PerformanceKpiScoreCard
            label="Backlog Health"
            :value="backlogHealthPct"
            format="pct"
            :status="deliveryStatus(backlogHealthPct)"
          />
        </div>
      </section>

      <!-- Team Leaderboard -->
      <section v-if="performanceStore.team.length > 0" class="mb-8">
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Team Performance
        </h2>
        <PerformanceTeamLeaderboard :members="performanceStore.team" />
      </section>

      <!-- Quality Gate -->
      <section v-if="reviewPendingTasks.length > 0" class="mb-8">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Quality Gate — Ready for approval
          </h2>
          <span class="rounded-full bg-purple-500/20 px-3 py-1 text-xs font-medium text-purple-300">
            {{ reviewPendingTasks.length }} awaiting
          </span>
        </div>
        <div class="overflow-hidden rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/80 text-left text-xs font-medium uppercase tracking-wider text-gray-400">
                  <th class="px-5 py-4">Task</th>
                  <th class="px-5 py-4">Developer</th>
                  <th class="px-5 py-4 text-center">AI Score</th>
                  <th class="px-5 py-4 text-right">Submitted</th>
                  <th class="px-5 py-4 text-center">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/80">
                <tr
                  v-for="task in reviewPendingTasks"
                  :key="task.id"
                  class="transition-colors hover:bg-gray-700/30"
                >
                  <td class="px-5 py-4">
                    <div class="font-medium text-white">{{ task.title }}</div>
                    <div class="mt-0.5 line-clamp-1 text-xs text-gray-500">
                      {{ task.description?.substring(0, 60) }}...
                    </div>
                  </td>
                  <td class="px-5 py-4">
                    <div class="flex items-center gap-2.5">
                      <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-purple-600 text-xs font-bold text-white">
                        D{{ task.assigned_to || task.created_by }}
                      </div>
                      <span class="text-gray-300">Dev #{{ task.assigned_to || task.created_by }}</span>
                    </div>
                  </td>
                  <td class="px-5 py-4 text-center">
                    <span class="inline-flex items-center gap-1 rounded-md bg-emerald-900/50 px-2 py-1 text-xs font-semibold text-emerald-300 ring-1 ring-emerald-500/30">
                      ✓ PASS
                    </span>
                  </td>
                  <td class="px-5 py-4 text-right text-xs text-gray-500">
                    {{ formatTimeAgo(task.updated_at || task.created_at) }}
                  </td>
                  <td class="px-5 py-4 text-center">
                    <button
                      type="button"
                      @click="goToTask(task)"
                      class="inline-flex items-center gap-1.5 rounded-lg bg-emerald-600 px-3 py-2 text-xs font-semibold text-white transition-colors hover:bg-emerald-500"
                    >
                      Review
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Resource Control: Unassigned + Active Development -->
      <section>
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          Resource Control
        </h2>
        <div class="grid gap-6 lg:grid-cols-2">
          <!-- Unassigned Queue -->
          <div class="flex flex-col rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
            <div class="flex items-center justify-between border-b border-gray-700/80 px-5 py-4">
              <h3 class="font-semibold text-white">Unassigned Queue</h3>
              <span class="rounded-full bg-amber-500/15 px-2.5 py-0.5 text-xs font-medium text-amber-400">
                {{ unassignedTasks.length }} tasks
              </span>
            </div>
            <div class="flex-1 overflow-auto p-4">
              <div v-if="unassignedTasks.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
                <div class="rounded-full bg-gray-700/80 p-4 text-gray-500">
                  <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                </div>
                <p class="mt-3 text-sm font-medium text-gray-400">No unassigned tasks</p>
                <p class="mt-1 text-xs text-gray-500">New tasks will appear here</p>
              </div>
              <ul v-else class="space-y-2">
                <li
                  v-for="task in unassignedTasks"
                  :key="task.id"
                  class="group flex items-center justify-between gap-3 rounded-lg border border-amber-500/20 bg-gray-900/50 p-3 transition-all hover:border-amber-500/40 hover:bg-gray-800/80"
                >
                  <button
                    type="button"
                    class="min-w-0 flex-1 text-left"
                    @click="goToTask(task)"
                  >
                    <span class="block truncate font-medium text-white group-hover:text-purple-300">
                      {{ task.title }}
                    </span>
                    <span class="mt-0.5 block text-xs text-gray-500">
                      Est: {{ (task.ai_estimated_minutes / 60).toFixed(1) }}h · {{ formatDate(task.created_at) }}
                    </span>
                  </button>
                  <button
                    type="button"
                    @click.stop="openAssignModal(task)"
                    class="shrink-0 rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 px-3 py-1.5 text-xs font-semibold text-white shadow-md transition-all hover:from-purple-500 hover:to-pink-500 hover:shadow-purple-500/25"
                  >
                    Assign
                  </button>
                </li>
              </ul>
            </div>
          </div>

          <!-- Active Development -->
          <div class="flex flex-col rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
            <div class="flex items-center justify-between border-b border-gray-700/80 px-5 py-4">
              <h3 class="font-semibold text-white">Active Development</h3>
              <span class="rounded-full bg-purple-500/15 px-2.5 py-0.5 text-xs font-medium text-purple-300">
                {{ Object.keys(tasksByDeveloper).length }} developers
              </span>
            </div>
            <div class="flex-1 overflow-auto p-4">
              <div v-if="Object.keys(tasksByDeveloper).length === 0" class="flex flex-col items-center justify-center py-12 text-center">
                <div class="rounded-full bg-gray-700/80 p-4 text-gray-500">
                  <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                </div>
                <p class="mt-3 text-sm font-medium text-gray-400">No assigned tasks</p>
                <p class="mt-1 text-xs text-gray-500">Assign tasks from the queue</p>
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="(devTasks, devId) in tasksByDeveloper"
                  :key="devId"
                  class="rounded-lg border p-4 transition-colors"
                  :class="
                    devTasks.length > 3
                      ? 'border-red-500/40 bg-red-950/20'
                      : 'border-gray-700/80 bg-gray-900/40'
                  "
                >
                  <div class="mb-3 flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <span class="font-semibold text-white">Developer #{{ devId }}</span>
                      <span class="text-xs text-gray-500">{{ devTasks.length }} tasks</span>
                      <span
                        v-if="devTasks.length > 3"
                        class="rounded bg-red-500/20 px-2 py-0.5 text-xs font-medium text-red-400"
                      >
                        Overloaded
                      </span>
                    </div>
                    <span class="text-xs font-medium text-gray-500">{{ getDevWorkload(devTasks) }}h</span>
                  </div>
                  <ul class="space-y-1.5">
                    <li
                      v-for="task in devTasks"
                      :key="task.id"
                      class="flex items-center justify-between gap-2 rounded-md bg-gray-800/80 py-2 pl-3 pr-2 transition-colors hover:bg-gray-700/60"
                    >
                      <button
                        type="button"
                        class="flex min-w-0 flex-1 items-center gap-2 text-left"
                        @click="goToTask(task)"
                      >
                        <span
                          :class="getStatusDot(task.status)"
                          class="h-2 w-2 shrink-0 rounded-full"
                        />
                        <span class="truncate text-sm text-gray-200">{{ task.title }}</span>
                        <span class="shrink-0 text-xs text-gray-500">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}h</span>
                      </button>
                      <button
                        type="button"
                        @click.stop="openReassignModal(task)"
                        class="shrink-0 rounded px-2 py-1 text-xs font-medium text-gray-400 transition-colors hover:bg-gray-600 hover:text-white"
                      >
                        Re-assign
                      </button>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Assign / Re-assign Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showAssignModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="assign-modal-title"
          @keydown.escape="closeAssignModal"
        >
          <div
            class="fixed inset-0 bg-black/70 backdrop-blur-sm"
            @click="closeAssignModal"
          />
          <div
            class="relative w-full max-w-md rounded-2xl border border-gray-700 bg-gray-800 shadow-2xl"
            @click.stop
          >
            <div class="border-b border-gray-700 px-6 py-4">
              <h2 id="assign-modal-title" class="text-lg font-semibold text-white">
                {{ isReassign ? 'Re-assign task' : 'Assign task' }}
              </h2>
              <p class="mt-1 line-clamp-2 text-sm text-gray-400">{{ selectedTask?.title }}</p>
            </div>
            <form class="p-6" @submit.prevent="assignTask">
              <label for="assign-dev-id" class="mb-2 block text-sm font-medium text-gray-300">
                Developer ID
              </label>
              <input
                id="assign-dev-id"
                v-model.number="assignDevId"
                type="number"
                min="1"
                placeholder="e.g. 1, 2, 3"
                class="w-full rounded-xl border border-gray-600 bg-gray-900 px-4 py-3 text-white placeholder-gray-500 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/30 focus:outline-none"
              />
              <div class="mt-6 flex gap-3">
                <button
                  type="button"
                  @click="closeAssignModal"
                  class="flex-1 rounded-xl border border-gray-600 bg-gray-700 px-4 py-2.5 text-sm font-medium text-gray-200 transition-colors hover:bg-gray-600"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="isAssigning || !assignDevId"
                  class="flex-1 rounded-xl bg-gradient-to-r from-purple-600 to-pink-600 px-4 py-2.5 text-sm font-semibold text-white shadow-lg transition-all hover:from-purple-500 hover:to-pink-500 disabled:pointer-events-none disabled:opacity-50"
                >
                  {{ isAssigning ? 'Assigning…' : 'Assign' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Success Toast -->
    <Teleport to="body">
      <Transition name="toast">
        <div
          v-if="showSuccess"
          class="fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-xl border border-emerald-500/40 bg-gray-800 px-5 py-4 shadow-xl shadow-black/30"
        >
          <span class="flex h-9 w-9 items-center justify-center rounded-full bg-emerald-600 text-white">✓</span>
          <p class="text-sm font-medium text-white">{{ successMessage }}</p>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import PerformanceKpiScoreCard from '~/components/performance/KpiScoreCard.vue'
import PerformanceTeamLeaderboard from '~/components/performance/TeamLeaderboard.vue'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'

interface Task {
  id: string
  title: string
  description: string
  status: string
  ai_estimated_minutes: number
  assigned_to?: number
  created_at: string
  updated_at?: string
  due_at?: string | null
}

const { fetchWithAuth } = useAuth()
const performanceStore = usePerformanceStore()

const tasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')

function deliveryStatus(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct >= 85) return 'good'
  if (pct >= 70) return 'warn'
  return 'bad'
}

function qualityStatus(q: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (q >= 75) return 'good'
  if (q >= 60) return 'warn'
  return 'bad'
}

function reworkStatus(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct <= 15) return 'good'
  if (pct <= 25) return 'warn'
  return 'bad'
}

const devMembers = computed(() =>
  performanceStore.team.filter((m) => m.role === 'DEV')
)

const avgTeamDeliveryRate = computed(() => {
  if (devMembers.value.length === 0) return 0
  const sum = devMembers.value.reduce((s, m) => s + m.delivery_rate_pct, 0)
  return sum / devMembers.value.length
})

const avgTeamCodeQuality = computed(() => {
  if (devMembers.value.length === 0) return 0
  const sum = devMembers.value.reduce((s, m) => s + m.code_quality_index, 0)
  return sum / devMembers.value.length
})

const avgTeamReworkRate = computed(() => {
  if (devMembers.value.length === 0) return 0
  const sum = devMembers.value.reduce((s, m) => s + m.rework_rate_pct, 0)
  return sum / devMembers.value.length
})

const backlogHealthPct = computed(() => {
  const active = tasks.value.filter((t) => t.status !== 'COMPLETED')
  if (active.length === 0) return 100
  const withAssigneeAndDue = active.filter(
    (t) => t.assigned_to && t.due_at
  ).length
  return (withAssigneeAndDue / active.length) * 100
})

const showAssignModal = ref(false)
const selectedTask = ref<Task | null>(null)
const assignDevId = ref<number | null>(null)
const isAssigning = ref(false)
const isReassign = ref(false)

const showSuccess = ref(false)
const successMessage = ref('')

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

const reviewPendingTasks = computed(() =>
  tasks.value
    .filter(t => t.status === 'REVIEW_PENDING')
    .sort((a, b) => new Date(b.updated_at || b.created_at).getTime() - new Date(a.updated_at || a.created_at).getTime())
)

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

    const taskIndex = tasks.value.findIndex(t => t.id === selectedTask.value?.id)
    if (taskIndex !== -1) {
      tasks.value[taskIndex].assigned_to = assignDevId.value
      tasks.value[taskIndex].status = 'IN_PROGRESS'
    }

    successMessage.value = `Task ${isReassign.value ? 'Re-assigned' : 'Assigned'} to Dev #${assignDevId.value}`
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
    IN_PROGRESS: 'bg-purple-500',
    PENDING: 'bg-yellow-500',
    BLOCKED: 'bg-red-500',
    REVIEW_PENDING: 'bg-purple-700'
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

const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}`)
}

onMounted(() => {
  fetchTasks()
  performanceStore.fetchAll('PM')
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(0.5rem);
}
</style>
