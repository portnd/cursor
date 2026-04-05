<template>
  <section class="bg-gray-800/50 border border-gray-700/80 rounded-xl p-5">
    <!-- Header -->
    <div class="flex flex-wrap items-center justify-between gap-2 mb-4">
      <div class="flex items-center gap-2">
        <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Sub-tasks</h2>
        <span class="text-xs bg-gray-700 text-gray-400 rounded-full px-2 py-0.5">{{ subtasks.length }}</span>
      </div>
      <div v-if="canEdit && !isMaxDepth" class="flex flex-wrap items-center gap-2 justify-end">
        <button
          v-if="!showAddForm"
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
      <span
        v-else-if="isMaxDepth"
        class="text-xs text-gray-600 italic w-full text-right sm:text-left"
      >Max depth (Level C)</span>
    </div>

    <!-- Roll-up summary bar -->
    <div v-if="subtasks.length > 0" class="mb-4 p-3 bg-gray-900/60 rounded-lg border border-gray-700/50 space-y-2">
      <div class="flex items-center justify-between text-xs">
        <span class="text-gray-500">Total Estimated Effort</span>
        <span class="font-semibold text-white">
          {{ formatMinutesAsHours(totalEstimatedMinutes) }}
          <span class="text-gray-400 font-normal">h</span>
          <span class="text-gray-500 font-normal text-xs ml-1">({{ totalEstimatedMinutes }} min)</span>
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
            <span v-if="sub.estimated_minutes" class="text-xs text-gray-600">· {{ formatMinutesAsHours(sub.estimated_minutes) }}h</span>
          </div>
        </div>

        <!-- Status badge -->
        <span class="text-xs px-2 py-0.5 rounded-full font-medium flex-shrink-0" :class="statusBadgeClass(sub.status)">
          {{ statusLabel(sub.status) }}
        </span>

        <!-- Split button (Product Owner / CEO only, visible on hover — hidden at max depth because split creates children) -->
        <button
          v-if="canEdit && !isMaxDepth"
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
      <div class="flex flex-wrap items-center justify-between gap-2 mb-3">
        <div class="text-xs font-semibold text-blue-400 uppercase tracking-wider">New Sub-task</div>
        <button
          v-if="parentTask"
          type="button"
          class="flex items-center gap-1.5 text-xs px-2.5 py-1.5 bg-emerald-900/30 hover:bg-emerald-800/40 text-emerald-300 border border-emerald-700/50 rounded-lg transition-colors"
          @click="fillFromParent"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
          </svg>
          Duplicate parent
        </button>
      </div>
      <div class="space-y-3">
        <input
          ref="titleInputRef"
          v-model="newSubtask.title"
          type="text"
          placeholder="Sub-task title..."
          class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
          @keydown.esc="cancelAddForm"
        />
        <div>
          <label class="block text-xs text-gray-500 mb-1">Description</label>
          <RichTextEditor v-model="newSubtask.description" placeholder="Optional — paste from parent or write details…" />
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block text-xs text-gray-500 mb-1">Assignee</label>
            <select v-model.number="newSubtask.assigned_to" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 focus:ring-2 focus:ring-blue-500 outline-none">
              <option :value="null">— Unassigned —</option>
              <option v-for="u in assigneeOptions" :key="u.id" :value="u.id">{{ u.display_name || u.email }} ({{ u.role }})</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">Estimated Effort (hours)</label>
            <input v-model.number="newSubtask.estimated_hours" type="number" min="0" step="0.1" placeholder="e.g. 1.5" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-blue-500 outline-none" />
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label class="block text-xs text-gray-500 mb-1">Type</label>
            <select v-model="newSubtask.task_type" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 focus:ring-2 focus:ring-blue-500 outline-none">
              <option value="TASK">TASK</option>
              <option value="FEATURE">FEATURE</option>
              <option value="BUG">BUG</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">Priority</label>
            <select v-model="newSubtask.priority" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 focus:ring-2 focus:ring-blue-500 outline-none">
              <option value="">— Default —</option>
              <option value="CRITICAL">CRITICAL</option>
              <option value="HIGH">HIGH</option>
              <option value="MEDIUM">MEDIUM</option>
              <option value="LOW">LOW</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">Story points</label>
            <input v-model.number="newSubtask.story_points" type="number" min="0" step="1" class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-sm text-gray-100 focus:ring-2 focus:ring-blue-500 outline-none" />
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
              <span v-if="splitTarget.estimated_minutes" class="ml-2 text-gray-500">({{ formatMinutesAsHours(splitTarget.estimated_minutes) }}h)</span>
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
                <!-- Est. hours -->
                <div>
                  <label class="block text-xs text-gray-500 mb-1">Est. (hours)</label>
                  <input v-model.number="item.estimated_hours" type="number" min="0" step="0.1" placeholder="0" class="w-full px-2.5 py-2 bg-gray-800 border border-gray-600 rounded-xl text-sm text-gray-100 placeholder-gray-500 focus:ring-2 focus:ring-amber-500/60 outline-none" />
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
              {{ formatMinutesAsHours(splitTotalMinutes) }} h
              <span class="text-gray-500 font-normal ml-1">({{ splitTotalMinutes }} min)</span>
              <span v-if="splitTarget?.estimated_minutes" class="text-gray-500 font-normal ml-1">· was {{ formatMinutesAsHours(splitTarget.estimated_minutes) }} h</span>
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
import RichTextEditor from '~/components/editor/RichTextEditor.vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { useAuth } from '~/composables/useAuth'
import { minutesToEffortHours, effortHoursToMinutes, formatMinutesAsHours } from '~/utils/effortHours'
import { isTaskAssigneeRole } from '~/utils/roles'

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

/** Fields read when duplicating the parent task into a new sub-task (API task shape is a superset). */
export interface ParentTaskCopySource {
  title?: string
  description?: string
  task_type?: string
  assigned_to?: number | null
  estimated_minutes?: number
  priority?: string
  story_points?: number
  epic_id?: string | null
  sprint_id?: string | null
  milestone_id?: string | null
  start_date?: string | null
  end_date?: string | null
  due_at?: string | null
}

interface SplitItem {
  title: string
  estimated_hours: number
  assignee_id: number | null
  priority: string
}

const props = defineProps<{
  parentTaskId: string
  projectId?: string | null
  /** Current page task — used to duplicate creatable fields into a new sub-task */
  parentTask?: ParentTaskCopySource | null
  subtasks: SubTask[]
  canEdit: boolean
  /** true when this task is already at level C (parent itself has a parent) — blocks adding more sub-tasks */
  isMaxDepth?: boolean
}>()

const emit = defineEmits<{
  (e: 'subtask-added', task: SubTask): void
  (e: 'refresh'): void
}>()

const { fetchWithAuth, currentUser } = useAuth()
const tasksApi = useTasksApi()
const { getTeams } = useTeamsApi()
const teamsStore = useTeamsStore()

// ── Add form ──────────────────────────────────────────────
const showAddForm = ref(false)
const isAdding = ref(false)
const addError = ref('')
const titleInputRef = ref<HTMLInputElement | null>(null)

interface NewSubtaskDraft {
  title: string
  description: string
  task_type: 'FEATURE' | 'TASK' | 'BUG'
  assigned_to: number | null
  estimated_hours: number
  priority: string
  story_points: number
  epic_id: string | null
  sprint_id: string | null
  milestone_id: string | null
  start_date: string | null
  end_date: string | null
  due_date: string | null
}

function emptyDraft(): NewSubtaskDraft {
  return {
    title: '',
    description: '',
    task_type: 'TASK',
    assigned_to: null,
    estimated_hours: 0,
    priority: '',
    story_points: 0,
    epic_id: null,
    sprint_id: null,
    milestone_id: null,
    start_date: null,
    end_date: null,
    due_date: null,
  }
}

const newSubtask = ref<NewSubtaskDraft>(emptyDraft())

// ── Shared assignees ──────────────────────────────────────
const assigneeOptions = ref<AssigneeOption[]>([])

// ── Split modal ───────────────────────────────────────────
const showSplitModal = ref(false)
const splitTarget = ref<SubTask | null>(null)
const splitItems = ref<SplitItem[]>([])
const isSplitting = ref(false)
const splitError = ref('')
const splitSubmitted = ref(false)

const splitTotalMinutes = computed(() =>
  splitItems.value.reduce((s, i) => s + effortHoursToMinutes(Number(i.estimated_hours) || 0), 0)
)

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
    if (role === 'PRODUCT_OWNER' || role === 'PM') {
      await teamsStore.fetchTeamsFeatureEnabled()
      if (teamsStore.teamsFeatureEnabled) {
        const userId = currentUser.value?.user_id
        const teams = await getTeams()
        const myTeam = teams.find((t) => t.users?.some((u) => u.id === userId))
        assigneeOptions.value = (myTeam?.users ?? [])
          .filter((u) => isTaskAssigneeRole(u.role))
          .map((u) => ({ id: u.id, email: u.email, display_name: u.display_name, role: u.role }))
      } else {
        const res = await fetchWithAuth<{ data: AssigneeOption[] }>('/auth/users')
        assigneeOptions.value = (res.data ?? []).filter((u) => isTaskAssigneeRole(u.role))
      }
    } else {
      const res = await fetchWithAuth<{ data: AssigneeOption[] }>('/auth/users')
      assigneeOptions.value = (res.data ?? []).filter((u) => isTaskAssigneeRole(u.role))
    }
  } catch { /* non-critical */ }
}

// ── Add form ──────────────────────────────────────────────
function normalizeTaskType(t?: string): 'FEATURE' | 'TASK' | 'BUG' {
  if (t === 'FEATURE' || t === 'BUG' || t === 'TASK') return t
  return 'TASK'
}

function fillFromParent() {
  const p = props.parentTask
  if (!p) return
  newSubtask.value = {
    title: p.title || '',
    description: p.description || '',
    task_type: normalizeTaskType(p.task_type),
    assigned_to: p.assigned_to ?? null,
    estimated_hours: minutesToEffortHours(p.estimated_minutes ?? 0),
    priority: p.priority || '',
    story_points: p.story_points ?? 0,
    epic_id: p.epic_id ?? null,
    sprint_id: p.sprint_id ?? null,
    milestone_id: p.milestone_id ?? null,
    start_date: p.start_date || null,
    end_date: p.end_date || null,
    due_date: p.due_at || null,
  }
}

async function openAddForm() {
  showAddForm.value = true
  newSubtask.value = emptyDraft()
  addError.value = ''
  await loadAssignees()
  await nextTick()
  titleInputRef.value?.focus()
}

function cancelAddForm() {
  showAddForm.value = false
  addError.value = ''
  newSubtask.value = emptyDraft()
}

async function submitAddSubtask() {
  if (!newSubtask.value.title.trim()) { addError.value = 'Title is required.'; return }
  isAdding.value = true
  addError.value = ''
  const d = newSubtask.value
  try {
    const payload: Record<string, unknown> = {
      title: d.title.trim(),
      parent_id: props.parentTaskId,
      estimated_minutes: effortHoursToMinutes(Number(d.estimated_hours) || 0),
      task_type: d.task_type || 'TASK',
      story_points: d.story_points ?? 0,
    }
    if (props.projectId) payload.project_id = props.projectId
    const desc = (d.description || '').trim()
    if (desc) payload.description = desc
    if (d.priority) payload.priority = d.priority
    if (d.epic_id) payload.epic_id = d.epic_id
    if (d.sprint_id) payload.sprint_id = d.sprint_id
    if (d.milestone_id) payload.milestone_id = d.milestone_id
    if (d.start_date) payload.start_date = d.start_date
    if (d.end_date) payload.end_date = d.end_date
    if (d.due_date) payload.due_date = d.due_date

    const created = await tasksApi.createTask(payload as Parameters<typeof tasksApi.createTask>[0])
    if (d.assigned_to) {
      try { await tasksApi.assignTask(created.id, d.assigned_to) } catch { /* non-fatal */ }
    }
    emit('subtask-added', created as SubTask)
    newSubtask.value = emptyDraft()
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
    estimated_hours: minutesToEffortHours(Math.round((src.estimated_minutes || 0) / 2)),
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
        estimated_minutes: effortHoursToMinutes(Number(i.estimated_hours) || 0),
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
    READY_FOR_TEST: 'bg-cyan-500',
    REVIEW_PENDING: 'bg-purple-500', BLOCKED: 'bg-red-500', PENDING: 'bg-yellow-500',
  }
  return map[status] ?? 'bg-gray-500'
}

function statusBadgeClass(status: string): string {
  const map: Record<string, string> = {
    COMPLETED: 'bg-green-900/50 text-green-300 border border-green-700',
    IN_PROGRESS: 'bg-blue-900/50 text-blue-300 border border-blue-700',
    READY_FOR_TEST: 'bg-cyan-900/50 text-cyan-300 border border-cyan-700',
    REVIEW_PENDING: 'bg-purple-900/50 text-purple-300 border border-purple-700',
    BLOCKED: 'bg-red-900/50 text-red-300 border border-red-700',
    PENDING: 'bg-yellow-900/30 text-yellow-400 border border-yellow-700/50',
  }
  return map[status] ?? 'bg-gray-700/50 text-gray-400 border border-gray-600'
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    COMPLETED: 'Done', IN_PROGRESS: 'In Progress',
    READY_FOR_TEST: 'Ready for Test',
    REVIEW_PENDING: 'Review', BLOCKED: 'Blocked', PENDING: 'Pending',
  }
  return map[status] ?? status
}

onMounted(() => {
  if (props.canEdit) loadAssignees()
})
</script>
