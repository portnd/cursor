<template>
  <section class="bg-gray-800/50 border border-gray-700/80 rounded-xl p-5">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Sub-tasks</h2>
        <span class="text-xs bg-gray-700 text-gray-400 rounded-full px-2 py-0.5">{{ subtasks.length }}</span>
      </div>
      <button
        v-if="canEdit && !showAddForm"
        type="button"
        @click="openAddForm"
        class="flex items-center gap-1.5 text-xs px-2.5 py-1.5 bg-blue-600/20 hover:bg-blue-600/40 text-blue-400 hover:text-blue-300 border border-blue-600/40 rounded-lg transition-colors"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Add Sub-task
      </button>
    </div>

    <!-- Roll-up summary bar -->
    <div v-if="subtasks.length > 0" class="mb-4 p-3 bg-gray-900/60 rounded-lg border border-gray-700/50 space-y-2">
      <div class="flex items-center justify-between text-xs">
        <span class="text-gray-500">Total Estimated Effort</span>
        <span class="font-semibold text-white">
          {{ totalEstimatedMinutes }} min
          <span class="text-gray-400 font-normal">({{ (totalEstimatedMinutes / 60).toFixed(1) }}h)</span>
        </span>
      </div>
      <div class="flex items-center justify-between text-xs">
        <span class="text-gray-500">Aggregate Progress</span>
        <span class="font-semibold" :class="aggregateProgress >= 100 ? 'text-green-400' : 'text-blue-400'">{{ aggregateProgress }}%</span>
      </div>
      <div class="w-full bg-gray-700 rounded-full h-1.5 overflow-hidden">
        <div class="h-1.5 rounded-full transition-all duration-500" :class="aggregateProgress >= 100 ? 'bg-green-500' : 'bg-blue-500'" :style="{ width: `${aggregateProgress}%` }" />
      </div>
      <div class="text-xs text-gray-500">{{ completedCount }} / {{ subtasks.length }} completed</div>
    </div>

    <!-- Empty state -->
    <div v-if="subtasks.length === 0 && !showAddForm" class="text-center py-6 text-gray-600 text-sm">
      No sub-tasks yet. Break this task down into smaller pieces.
    </div>

    <!-- Sub-task list -->
    <ul v-if="subtasks.length > 0" class="space-y-2 mb-3">
      <li
        v-for="sub in subtasks"
        :key="sub.id"
        class="flex items-center gap-3 p-3 bg-gray-900/50 rounded-lg border border-gray-700/40 hover:border-gray-600/60 group transition-colors cursor-pointer"
        @click.self="navigateTo(`/task/${sub.id}`)"
      >
        <!-- Status dot -->
        <div class="w-2 h-2 rounded-full flex-shrink-0 mt-0.5" :class="statusDotClass(sub.status)" :title="sub.status" />

        <!-- Title + assignee -->
        <div class="flex-1 min-w-0" @click="navigateTo(`/task/${sub.id}`)">
          <p
            class="text-sm text-gray-200 group-hover:text-white truncate transition-colors"
            :class="{ 'line-through text-gray-500': sub.status === 'COMPLETED' }"
          >
            {{ sub.title }}
          </p>
          <div class="flex items-center gap-2 mt-0.5 flex-wrap">
            <span v-if="sub.assigned_to" class="text-xs text-gray-500">
              {{ sub.assigned_to_display_name || sub.assigned_to_email || `Dev #${sub.assigned_to}` }}
            </span>
            <span v-else class="text-xs text-gray-600 italic">Unassigned</span>
            <span v-if="sub.estimated_minutes" class="text-xs text-gray-600">· {{ sub.estimated_minutes }}min</span>
          </div>
        </div>

        <!-- Status badge -->
        <span class="text-xs px-2 py-0.5 rounded-full font-medium flex-shrink-0" :class="statusBadgeClass(sub.status)">
          {{ statusLabel(sub.status) }}
        </span>

        <!-- Split button (PM/CEO only, visible on hover) -->
        <button
          v-if="canEdit"
          type="button"
          @click.stop="openSplitModal(sub)"
          title="Duplicate & Split this sub-task"
          class="opacity-0 group-hover:opacity-100 flex items-center gap-1 text-xs px-2 py-1 rounded-lg bg-amber-900/30 hover:bg-amber-900/50 text-amber-400 hover:text-amber-300 border border-amber-700/40 hover:border-amber-600/60 transition-all shrink-0"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/>
          </svg>
          Split
        </button>
      </li>
    </ul>

    <!-- Add Sub-task form (inline) -->
    <div v-if="showAddForm" class="mt-3 p-4 bg-gray-900/70 rounded-xl border border-blue-600/30">
      <div class="text-xs font-semibold text-blue-400 uppercase tracking-wider mb-3">New Sub-task</div>
      <div class="space-y-3">
        <input
          ref="titleInputRef"
          v-model="newSubtask.title"
          type="text"
          placeholder="Sub-task title..."
          class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
          @keydown.enter="submitAddSubtask"
          @keydown.esc="cancelAddForm"
        />
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block text-xs text-gray-500 mb-1">Assignee</label>
            <select v-model.number="newSubtask.assigned_to" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 focus:ring-2 focus:ring-blue-500 outline-none">
              <option :value="null">— Unassigned —</option>
              <option v-for="u in assigneeOptions" :key="u.id" :value="u.id">{{ u.display_name || u.email }} ({{ u.role }})</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">Estimated Minutes</label>
            <input v-model.number="newSubtask.estimated_minutes" type="number" min="0" step="15" placeholder="e.g. 60" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-blue-500 outline-none" />
          </div>
        </div>
        <p v-if="addError" class="text-xs text-red-400">{{ addError }}</p>
        <div class="flex items-center gap-2 pt-1">
          <button type="button" :disabled="isAdding || !newSubtask.title.trim()" @click="submitAddSubtask" class="px-4 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors flex items-center gap-1.5">
            <svg v-if="isAdding" class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
            {{ isAdding ? 'Adding...' : 'Add Sub-task' }}
          </button>
          <button type="button" @click="cancelAddForm" class="px-4 py-1.5 text-sm text-gray-400 hover:text-gray-200 rounded-lg transition-colors">Cancel</button>
        </div>
      </div>
    </div>
  </section>

  <!-- ══ SPLIT MODAL ══ -->
  <Teleport to="body">
    <div v-if="showSplitModal && splitTarget" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto" @click.self="closeSplitModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-2xl my-auto flex flex-col max-h-[90vh]">

        <!-- Header -->
        <div class="flex items-center justify-between px-6 pt-5 pb-4 border-b border-gray-700/60 shrink-0">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <div class="w-6 h-6 rounded-md bg-amber-900/40 border border-amber-700/50 flex items-center justify-center shrink-0">
                <svg class="w-3.5 h-3.5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/></svg>
              </div>
              <h2 class="text-base font-bold text-white">Duplicate & Split Sub-task</h2>
            </div>
            <p class="text-xs text-gray-500 truncate max-w-md">
              Original: <span class="text-gray-300">{{ splitTarget.title }}</span>
              <span v-if="splitTarget.estimated_minutes" class="ml-2 text-gray-500">({{ splitTarget.estimated_minutes }}min)</span>
            </p>
          </div>
          <button @click="closeSplitModal" class="text-gray-500 hover:text-white transition-colors ml-4 shrink-0">✕</button>
        </div>

        <!-- Body -->
        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">

          <!-- Info banner -->
          <div class="flex items-start gap-3 p-3 bg-amber-900/20 border border-amber-700/40 rounded-xl text-xs text-amber-300">
            <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            <span>Task ต้นฉบับจะ<strong> ถูกลบ</strong> และแทนที่ด้วย {{ splitItems.length }} sub-tasks ใหม่ที่มี parent เดิม — ไม่สามารถ undo ได้</span>
          </div>

          <!-- Split items -->
          <div class="space-y-3">
            <div
              v-for="(item, idx) in splitItems"
              :key="idx"
              class="p-4 bg-gray-900/60 border border-gray-700/60 rounded-xl space-y-3"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Sub-task {{ idx + 1 }}</span>
                <button
                  v-if="splitItems.length > 2"
                  type="button"
                  @click="removeSplitItem(idx)"
                  class="text-xs text-red-400 hover:text-red-300 px-2 py-0.5 rounded hover:bg-red-900/20 transition-colors"
                >
                  Remove
                </button>
              </div>
              <!-- Title -->
              <input
                v-model="item.title"
                type="text"
                :placeholder="`e.g. ${splitTarget?.title} - Part ${idx + 1}`"
                class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-amber-500/60 focus:border-amber-500/60 outline-none"
                :class="splitSubmitted && !item.title.trim() ? 'border-red-500/60' : ''"
              />
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <!-- Assignee -->
                <div class="sm:col-span-1">
                  <label class="block text-xs text-gray-500 mb-1">Assignee</label>
                  <select v-model="item.assignee_id" class="w-full px-2.5 py-2 bg-gray-800 border border-gray-600 rounded-xl text-xs text-gray-100 focus:ring-2 focus:ring-amber-500/60 outline-none">
                    <option :value="null">— Unassigned —</option>
                    <option v-for="u in assigneeOptions" :key="u.id" :value="u.id">{{ u.display_name || u.email }}</option>
                  </select>
                </div>
                <!-- Est. Minutes -->
                <div>
                  <label class="block text-xs text-gray-500 mb-1">Est. Minutes</label>
                  <input v-model.number="item.estimated_minutes" type="number" min="0" step="15" placeholder="0" class="w-full px-2.5 py-2 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-amber-500/60 outline-none" />
                </div>
                <!-- Priority -->
                <div>
                  <label class="block text-xs text-gray-500 mb-1">Priority</label>
                  <select v-model="item.priority" class="w-full px-2.5 py-2 bg-gray-800 border border-gray-600 rounded-xl text-xs text-gray-100 focus:ring-2 focus:ring-amber-500/60 outline-none">
                    <option value="">Inherit</option>
                    <option value="CRITICAL">🔴 Critical</option>
                    <option value="HIGH">🟠 High</option>
                    <option value="MEDIUM">🟡 Medium</option>
                    <option value="LOW">🟢 Low</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          <!-- Add another split -->
          <button
            type="button"
            @click="addSplitItem"
            class="w-full flex items-center justify-center gap-2 py-2.5 border border-dashed border-gray-600 hover:border-amber-600/60 text-gray-500 hover:text-amber-400 rounded-xl text-sm transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
            Add another piece
          </button>

          <!-- Total minutes summary -->
          <div class="flex items-center justify-between text-xs px-1">
            <span class="text-gray-500">Total estimated effort after split</span>
            <span class="font-semibold" :class="splitTotalMinutes > (splitTarget?.estimated_minutes || 0) ? 'text-amber-400' : 'text-gray-300'">
              {{ splitTotalMinutes }} min
              <span v-if="splitTarget?.estimated_minutes" class="text-gray-500 font-normal ml-1">(was {{ splitTarget.estimated_minutes }} min)</span>
            </span>
          </div>

          <p v-if="splitError" class="p-3 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-sm">{{ splitError }}</p>
        </div>

        <!-- Footer -->
        <div class="flex gap-3 px-6 py-4 border-t border-gray-700/60 shrink-0">
          <button
            @click="submitSplit"
            :disabled="isSplitting || splitItems.some(i => !i.title.trim())"
            class="flex-1 flex items-center justify-center gap-2 py-2.5 bg-gradient-to-r from-amber-600 to-orange-600 hover:from-amber-700 hover:to-orange-700 disabled:opacity-40 text-white text-sm font-semibold rounded-xl transition-colors"
          >
            <svg v-if="isSplitting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/></svg>
            {{ isSplitting ? 'Splitting...' : `Confirm Split into ${splitItems.length} Sub-tasks` }}
          </button>
          <button @click="closeSplitModal" :disabled="isSplitting" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 text-gray-300 text-sm rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from 'vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useAuth } from '~/composables/useAuth'

interface SubTask {
  id: string
  title: string
  status: string
  assigned_to: number | null
  assigned_to_display_name?: string
  assigned_to_email?: string
  estimated_minutes: number
  progress: number
}

interface AssigneeOption {
  id: number
  email: string
  display_name: string
  role: string
}

interface SplitItem {
  title: string
  estimated_minutes: number
  assignee_id: number | null
  priority: string
}

const props = defineProps<{
  parentTaskId: string
  projectId?: string | null
  subtasks: SubTask[]
  canEdit: boolean
}>()

const emit = defineEmits<{
  (e: 'subtask-added', task: SubTask): void
  (e: 'refresh'): void
}>()

const { fetchWithAuth, currentUser } = useAuth()
const tasksApi = useTasksApi()
const { getTeams } = useTeamsApi()

// ── Add form ──────────────────────────────────────────────
const showAddForm = ref(false)
const isAdding = ref(false)
const addError = ref('')
const titleInputRef = ref<HTMLInputElement | null>(null)

const newSubtask = ref<{ title: string; assigned_to: number | null; estimated_minutes: number }>({
  title: '', assigned_to: null, estimated_minutes: 0,
})

// ── Shared assignees ──────────────────────────────────────
const assigneeOptions = ref<AssigneeOption[]>([])

// ── Split modal ───────────────────────────────────────────
const showSplitModal = ref(false)
const splitTarget = ref<SubTask | null>(null)
const splitItems = ref<SplitItem[]>([])
const isSplitting = ref(false)
const splitError = ref('')
const splitSubmitted = ref(false)

const splitTotalMinutes = computed(() => splitItems.value.reduce((s, i) => s + (i.estimated_minutes || 0), 0))

// ── Roll-up ───────────────────────────────────────────────
const totalEstimatedMinutes = computed(() =>
  props.subtasks.reduce((sum, s) => sum + (s.estimated_minutes || 0), 0)
)
const completedCount = computed(() => props.subtasks.filter((s) => s.status === 'COMPLETED').length)
const aggregateProgress = computed(() => {
  if (props.subtasks.length === 0) return 0
  const total = props.subtasks.reduce((sum, s) => {
    if (s.status === 'COMPLETED') return sum + 100
    return sum + (s.progress || 0)
  }, 0)
  return Math.round(total / props.subtasks.length)
})

// ── Load assignees ────────────────────────────────────────
async function loadAssignees() {
  if (assigneeOptions.value.length > 0) return
  try {
    const role = (currentUser.value?.role || '').toUpperCase()
    if (role === 'PM') {
      const userId = currentUser.value?.user_id
      const teams = await getTeams()
      const myTeam = teams.find((t) => t.users?.some((u) => u.id === userId))
      assigneeOptions.value = (myTeam?.users ?? [])
        .filter((u) => ['DEV', 'PM', 'MANAGER', 'SUPPORT'].includes(u.role))
        .map((u) => ({ id: u.id, email: u.email, display_name: u.display_name, role: u.role }))
    } else {
      const res = await fetchWithAuth<{ data: AssigneeOption[] }>('/auth/users')
      assigneeOptions.value = (res.data ?? []).filter((u) =>
        ['DEV', 'PM', 'MANAGER', 'SUPPORT'].includes(u.role)
      )
    }
  } catch { /* non-critical */ }
}

// ── Add form ──────────────────────────────────────────────
async function openAddForm() {
  showAddForm.value = true
  newSubtask.value = { title: '', assigned_to: null, estimated_minutes: 0 }
  addError.value = ''
  await loadAssignees()
  await nextTick()
  titleInputRef.value?.focus()
}

function cancelAddForm() {
  showAddForm.value = false
  addError.value = ''
}

async function submitAddSubtask() {
  if (!newSubtask.value.title.trim()) { addError.value = 'Title is required.'; return }
  isAdding.value = true
  addError.value = ''
  try {
    const payload: Record<string, unknown> = {
      title: newSubtask.value.title.trim(),
      parent_id: props.parentTaskId,
      estimated_minutes: newSubtask.value.estimated_minutes || 0,
    }
    if (props.projectId) payload.project_id = props.projectId
    if (newSubtask.value.assigned_to) payload.assigned_to = newSubtask.value.assigned_to

    const created = await tasksApi.createTask(payload as Parameters<typeof tasksApi.createTask>[0])
    if (newSubtask.value.assigned_to) {
      try { await tasksApi.assignTask(created.id, newSubtask.value.assigned_to) } catch { /* non-fatal */ }
    }
    emit('subtask-added', created as SubTask)
    newSubtask.value = { title: '', assigned_to: null, estimated_minutes: 0 }
    showAddForm.value = false
    emit('refresh')
  } catch (err: any) {
    addError.value = err?.data?.message ?? err?.message ?? 'Failed to create sub-task.'
  } finally {
    isAdding.value = false
  }
}

// ── Split modal ───────────────────────────────────────────
function makeSplitItem(src: SubTask, suffix: string): SplitItem {
  const baseTitle = src.title.replace(/\s*-\s*Part\s*\d+$/i, '').trim()
  return {
    title: `${baseTitle} - ${suffix}`,
    estimated_minutes: Math.round((src.estimated_minutes || 0) / 2),
    assignee_id: null,
    priority: '',
  }
}

async function openSplitModal(sub: SubTask) {
  splitTarget.value = sub
  splitItems.value = [
    makeSplitItem(sub, 'Part 1'),
    makeSplitItem(sub, 'Part 2'),
  ]
  splitError.value = ''
  splitSubmitted.value = false
  showSplitModal.value = true
  await loadAssignees()
}

function closeSplitModal() {
  showSplitModal.value = false
  splitTarget.value = null
  splitItems.value = []
  splitError.value = ''
  splitSubmitted.value = false
}

function addSplitItem() {
  if (!splitTarget.value) return
  splitItems.value.push(makeSplitItem(splitTarget.value, `Part ${splitItems.value.length + 1}`))
}

function removeSplitItem(idx: number) {
  if (splitItems.value.length <= 2) return
  splitItems.value.splice(idx, 1)
}

async function submitSplit() {
  splitSubmitted.value = true
  if (splitItems.value.some((i) => !i.title.trim())) {
    splitError.value = 'All sub-tasks must have a title.'
    return
  }
  if (!splitTarget.value) return
  isSplitting.value = true
  splitError.value = ''
  try {
    await tasksApi.splitTask(
      splitTarget.value.id,
      splitItems.value.map((i) => ({
        title: i.title.trim(),
        estimated_minutes: i.estimated_minutes || 0,
        assignee_id: i.assignee_id ?? null,
        priority: i.priority || undefined,
      }))
    )
    closeSplitModal()
    emit('refresh')
  } catch (err: any) {
    splitError.value = err?.data?.message ?? err?.message ?? 'Split failed.'
  } finally {
    isSplitting.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────
function statusDotClass(status: string): string {
  const map: Record<string, string> = {
    COMPLETED: 'bg-green-500', IN_PROGRESS: 'bg-blue-500',
    REVIEW_PENDING: 'bg-purple-500', BLOCKED: 'bg-red-500', PENDING: 'bg-yellow-500',
  }
  return map[status] ?? 'bg-gray-500'
}

function statusBadgeClass(status: string): string {
  const map: Record<string, string> = {
    COMPLETED: 'bg-green-900/50 text-green-300 border border-green-700',
    IN_PROGRESS: 'bg-blue-900/50 text-blue-300 border border-blue-700',
    REVIEW_PENDING: 'bg-purple-900/50 text-purple-300 border border-purple-700',
    BLOCKED: 'bg-red-900/50 text-red-300 border border-red-700',
    PENDING: 'bg-yellow-900/30 text-yellow-400 border border-yellow-700/50',
  }
  return map[status] ?? 'bg-gray-700/50 text-gray-400 border border-gray-600'
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    COMPLETED: 'Done', IN_PROGRESS: 'In Progress',
    REVIEW_PENDING: 'Review', BLOCKED: 'Blocked', PENDING: 'Pending',
  }
  return map[status] ?? status
}

onMounted(() => {
  if (props.canEdit) loadAssignees()
})
</script>
