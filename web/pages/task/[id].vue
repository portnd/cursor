<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลด task...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-4xl mx-auto">
      <div class="bg-red-900/20 border border-red-500 rounded p-6 text-red-400">
        <h2 class="text-xl font-bold mb-2">Failed to load task</h2>
        <p>{{ error }}</p>
        <button 
          @click="goToDashboard"
          class="mt-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded transition-colors"
        >
          ← Back to Dashboard
        </button>
      </div>
    </div>

    <!-- Content: Enterprise Task Detail -->
    <div v-else-if="task">
      <!-- Header: minimal, professional -->
      <div class="border-b border-gray-700/80 pb-5 mb-6">
        <!-- Row 1: nav (left) + actions (right) -->
        <div class="flex flex-wrap items-center justify-between gap-4">
          <nav class="flex items-center gap-2 text-sm text-gray-500">
            <NuxtLink :to="backTarget" class="hover:text-gray-300 transition-colors">← Back</NuxtLink>
            <span class="text-gray-600">/</span>
            <span class="text-gray-400 truncate">Task</span>
          </nav>
          <div class="flex flex-wrap items-center gap-2 shrink-0">
            <span
              :class="getStatusClass(task.status)"
              class="px-3 py-1.5 text-xs font-medium rounded-full uppercase tracking-wide"
            >
              {{ getStatusLabel(task.status) }}
            </span>
            <template v-if="canEditOrDelete">
              <button
                @click="openEditModal"
                class="px-3 py-1.5 text-sm font-medium text-gray-200 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                Edit
              </button>
              <button
                @click="openDeleteConfirmation"
                class="px-3 py-1.5 text-sm font-medium text-red-300 bg-red-900/30 hover:bg-red-900/50 border border-red-800 rounded-lg transition-colors"
              >
                Delete
              </button>
            </template>
            <NuxtLink
              :to="backTarget"
              class="px-3 py-1.5 text-sm font-medium text-gray-400 hover:text-white rounded-lg transition-colors"
            >
              Back
            </NuxtLink>
          </div>
        </div>
        <!-- Row 2: title + code + date + by (bottom-left) | Prev/Next (bottom-right) -->
        <div class="mt-4 flex flex-wrap items-end justify-between gap-4">
          <div class="min-w-0">
            <h1 class="text-xl sm:text-2xl font-semibold text-white tracking-tight">{{ task.title }}</h1>
            <div class="mt-2 flex flex-wrap items-center gap-3 text-xs text-gray-500">
              <code class="font-mono text-gray-400">{{ task.code || task.id }}</code>
              <span>Created {{ formatDate(task.created_at) }}</span>
              <span v-if="creatorLabel" class="text-gray-500">by {{ creatorLabel }}</span>
            </div>
          </div>
          <div v-if="inSprintContext" class="flex flex-wrap items-center gap-2 shrink-0">
            <button
              type="button"
              aria-label="Task ก่อนหน้า"
              :class="prevTaskLink
                ? 'px-3 py-1.5 text-sm font-medium text-gray-400 hover:text-white rounded-lg transition-colors border border-gray-600 hover:border-gray-500 cursor-pointer'
                : 'px-3 py-1.5 text-sm text-gray-600 rounded-lg border border-gray-700 cursor-pointer hover:border-gray-600'"
              @click="goToPrevTask"
            >
              ← Task ก่อนหน้า
            </button>
            <button
              type="button"
              aria-label="Task ถัดไป"
              :class="nextTaskLink
                ? 'px-3 py-1.5 text-sm font-medium text-gray-400 hover:text-white rounded-lg transition-colors border border-gray-600 hover:border-gray-500 cursor-pointer'
                : 'px-3 py-1.5 text-sm text-gray-600 rounded-lg border border-gray-700 cursor-pointer hover:border-gray-600'"
              @click="goToNextTask"
            >
              Task ถัดไป →
            </button>
            <span v-if="sprintNavMessage" class="text-xs text-amber-400 ml-1">{{ sprintNavMessage }}</span>
          </div>
        </div>
      </div>

      <!-- Main: two columns -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left: Description + Details -->
        <div class="lg:col-span-2 space-y-6">
          <section class="bg-gray-800/50 border border-gray-700/80 rounded-xl p-5">
            <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
              <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Description</h2>
              <div class="flex items-center gap-2">
                <a
                  v-if="slideOpenInSlidesURL"
                  :href="slideOpenInSlidesURL"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-xs text-indigo-400 hover:text-indigo-300 flex items-center gap-1"
                >
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                  </svg>
                  Open in Slides
                </a>
                <button
                  v-if="canEditOrDelete && !isEditingDescription"
                  @click="startInlineEdit"
                  class="text-xs text-gray-500 hover:text-blue-400 transition-colors flex items-center gap-1"
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                  </svg>
                  Edit
                </button>
                <div v-if="isEditingDescription" class="flex items-center gap-2">
                  <button
                    @click="saveInlineDescription"
                    :disabled="isSavingDescription"
                    class="text-xs px-2.5 py-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded transition-colors"
                  >
                    {{ isSavingDescription ? 'Saving...' : 'Save' }}
                  </button>
                  <button
                    @click="cancelInlineEdit"
                    class="text-xs px-2.5 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded transition-colors"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>

            <!-- Inline rich editor (edit mode) -->
            <RichTextEditor
              v-if="isEditingDescription"
              v-model="inlineDescriptionHtml"
              placeholder="Describe what needs to be done... (paste images with ⌘V)"
            />

            <!-- Read-only rich text view -->
            <div v-else>
              <RichTextEditor
                v-if="task.description && task.description.trim()"
                :model-value="task.description"
                :readonly="true"
              />
              <p v-else class="text-gray-500 text-sm italic">No description. Click Edit to add one.</p>
            </div>
          </section>
        </div>

        <!-- Right: Key details card -->
        <div class="lg:col-span-1">
          <section class="bg-gray-800/50 border border-gray-700/80 rounded-xl p-5 sticky top-4">
            <h2 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">Details</h2>
            <dl class="space-y-4">
              <div>
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Assignee</dt>
                <dd v-if="task.assigned_to" class="text-sm font-medium text-white">Dev #{{ task.assigned_to }}</dd>
                <dd v-else class="text-sm text-gray-500">Unassigned</dd>
              </div>
              <div v-if="task.priority">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Priority</dt>
                <dd class="text-sm font-medium text-white">{{ task.priority }}</dd>
              </div>
              <div v-if="task.story_points != null">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Story points</dt>
                <dd class="text-sm font-medium text-white">{{ task.story_points }}</dd>
              </div>
              <div v-if="task.due_at">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Due date</dt>
                <dd class="text-sm font-medium text-white">{{ formatDateTime(task.due_at) }}</dd>
                <div
                  v-if="task.status !== 'COMPLETED' && task.due_at"
                  :class="[
                    'text-xs mt-0.5',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400' : 'text-gray-500'
                  ]"
                >
                  {{ getDeadlineCountdown(task.due_at) }}
                </div>
              </div>
              <div v-if="task.start_date">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Start date</dt>
                <dd class="text-sm text-gray-300">{{ formatDateTime(task.start_date) }}</dd>
              </div>
              <div v-if="task.end_date">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">End date</dt>
                <dd class="text-sm text-gray-300">{{ formatDateTime(task.end_date) }}</dd>
              </div>
              <div v-if="task.completed_at">
                <dt class="text-xs text-gray-500 uppercase tracking-wider mb-1">Completed</dt>
                <dd class="text-sm text-green-400">{{ formatDateTime(task.completed_at) }}</dd>
                <dd v-if="task.started_at" class="text-xs text-gray-500 mt-0.5">
                  Duration: {{ calculateDuration(task.started_at, task.completed_at) }}
                </dd>
              </div>
            </dl>
          </section>
        </div>
      </div>

      <!-- BOTTOM SECTION: Discussion & Time Tracking -->
      <div class="mt-6 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Discussion / Comments -->
        <div class="bg-gray-800 border border-gray-700 rounded-xl p-6">
          <h2 class="text-sm font-bold text-gray-400 uppercase mb-4 flex items-center gap-2">
            <span>💬</span> Discussion
            <span class="ml-auto text-xs font-normal text-gray-500">{{ comments.length }} comments</span>
          </h2>
          <TaskComments
            :comments="comments"
            :loading="commentsLoading"
            @add-comment="handleAddComment"
          />
        </div>

        <!-- Time Tracking -->
        <div class="bg-gray-800 border border-gray-700 rounded-xl p-6">
          <h2 class="text-sm font-bold text-gray-400 uppercase mb-4 flex items-center gap-2">
            <span>⏱️</span> Time Tracking
          </h2>
          <TimeLogger
            :time-logs="timeLogs"
            :estimated-minutes="task.ai_estimated_minutes || 0"
            :loading="timeLogsLoading"
            @log-time="handleLogTime"
          />
        </div>
      </div>
    </div>

    <!-- Edit Task Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeEditModal"
    >
      <div class="bg-gray-800 border-2 border-blue-600 rounded-lg p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6 pb-4 border-b-2 border-gray-700">
          <div>
            <h2 class="text-2xl font-bold text-white flex items-center gap-3">
              <span class="text-3xl">✏️</span>
              <span>Edit Task</span>
            </h2>
            <div class="text-sm text-gray-400 mt-1">Update task details</div>
          </div>
          <button
            @click="closeEditModal"
            class="text-gray-400 hover:text-white transition-colors text-2xl"
            :disabled="isUpdatingTask"
          >
            ✕
          </button>
        </div>

        <!-- Error Message -->
        <div v-if="editError" class="mb-4 p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
          {{ editError }}
        </div>

        <!-- Form -->
        <div class="space-y-5">
          <!-- Title -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📋</span>
                <span>Title</span>
              </span>
            </label>
            <input
              v-model="editForm.title"
              type="text"
              placeholder="Enter task title..."
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            />
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📝</span>
                <span>Description</span>
              </span>
            </label>
            <RichTextEditor
              v-model="editForm.description"
              placeholder="Describe the task objectives and requirements... (paste images with ⌘V)"
            />
          </div>

          <!-- Priority & Story Points (same as Create Task) -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">Priority</label>
              <select
                v-model="editForm.priority"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :disabled="isUpdatingTask"
              >
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">Story Points</label>
              <input
                v-model.number="editForm.story_points"
                type="number"
                min="0"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :disabled="isUpdatingTask"
              />
            </div>
          </div>

          <!-- Sprint (same as Create Task) -->
          <div v-if="editSprints.length > 0">
            <label class="block text-sm font-medium text-gray-300 mb-2">Sprint</label>
            <select
              v-model="editForm.sprint_id"
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            >
              <option value="">Backlog</option>
              <option v-for="s in editSprints" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </div>

          <!-- Due date / Start / End (same as Create Task) -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2"><span>📅</span>Due date</span>
            </label>
            <input
              v-model="editForm.deadline"
              type="datetime-local"
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">Start date</label>
              <input
                v-model="editForm.start_date"
                type="datetime-local"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :disabled="isUpdatingTask"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">End date</label>
              <input
                v-model="editForm.end_date"
                type="datetime-local"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :disabled="isUpdatingTask"
              />
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4 mt-6 pt-4 border-t-2 border-gray-700">
          <button
            @click="submitEdit"
            :disabled="isUpdatingTask || !editForm.title.trim()"
            class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
          >
            <span v-if="isUpdatingTask" class="animate-spin">⚙️</span>
            <span v-else>💾</span>
            <span>{{ isUpdatingTask ? 'Updating...' : 'Update task' }}</span>
          </button>
          
          <button
            @click="closeEditModal"
            :disabled="isUpdatingTask"
            class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded-lg transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeDeleteModal"
    >
      <div class="bg-gray-800 border-2 border-red-600 rounded-lg p-6 max-w-lg w-full relative">
        <!-- Close Button (X) -->
        <button
          @click="closeDeleteModal"
          class="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors"
          :disabled="isDeletingTask"
          title="Close"
        >
          <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>

        <!-- Header -->
        <div class="flex items-center gap-3 mb-6">
          <span class="text-4xl">🗑️</span>
          <div>
            <h2 class="text-2xl font-bold text-white">Delete task?</h2>
            <div class="text-sm text-gray-400 mt-1">This action cannot be undone</div>
          </div>
        </div>

        <!-- Warning -->
        <div class="mb-6 p-4 bg-red-900/30 border-2 border-red-600/50 rounded-lg">
          <div class="text-sm text-red-300 leading-relaxed">
            <p class="font-bold mb-2">⚠️ Critical Operation</p>
            <p>Are you sure you want to <strong>permanently delete</strong> this task?</p>
            <p class="mt-2">This will remove:</p>
            <ul class="list-disc list-inside ml-2 mt-1">
              <li>Task data</li>
              <li>Comments and time logs</li>
              <li>All appeals</li>
              <li>Complete audit trail</li>
            </ul>
          </div>
        </div>

        <!-- Task Info -->
        <div class="mb-6 p-4 bg-gray-900 border border-gray-700 rounded">
          <div class="text-xs text-gray-500 uppercase mb-1">Task to delete:</div>
          <div class="font-bold text-white">{{ task?.title }}</div>
          <div class="text-xs text-gray-400 mt-1">ID: {{ task?.id }}</div>
        </div>

        <!-- Error Message -->
        <div v-if="deleteError" class="mb-4 p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
          {{ deleteError }}
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4">
          <button
            @click="confirmDelete"
            :disabled="isDeletingTask"
            class="flex-1 px-6 py-3 bg-red-600 hover:bg-red-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
          >
            <span v-if="isDeletingTask" class="animate-spin">⚙️</span>
            <span v-else>💥</span>
            <span>{{ isDeletingTask ? 'Deleting...' : 'Yes, Delete Forever' }}</span>
          </button>
          
          <button
            @click="closeDeleteModal"
            :disabled="isDeletingTask"
            class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded-lg transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'
import TaskComments from '~/components/tasks/TaskComments.vue'
import TimeLogger from '~/components/tasks/TimeLogger.vue'
import RichTextEditor from '~/components/editor/RichTextEditor.vue'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { TaskComment, TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

interface Appeal {
  id: string
  submission_id: string
  developer_id: number
  reason: string
  status: string // PENDING, APPROVED, REJECTED
  
  // AI Advisory System
  ai_recommendation: string // OVERTURN or UPHOLD
  ai_confidence: number     // 0-100
  ai_reasoning: string      // Advice for CEO/PM
  
  resolver_id: number | null
  resolver_note: string
  created_at: string
  updated_at: string
}

interface Submission {
  id: string
  task_id: string
  dev_id: number
  commit_hash: string
  ai_verdict: string
  ai_score: number
  ai_feedback: any
  is_overridden: boolean
  appeal?: Appeal
  created_at: string
}

interface Task {
  id: string
  code?: string // e.g. mims-hdmap-main-001
  project_id?: string | null // UUID of project (for Back to project when task has no code)
  title: string
  description: string
  resource_urls: any
  ai_estimated_minutes: number
  
  // Time Negotiation
  negotiation_status: string // NONE, PENDING, APPROVED, REJECTED
  proposed_minutes: number
  negotiation_reason: string
  
  due_at: string | null
  start_date: string | null
  end_date: string | null
  priority?: string
  story_points?: number
  sprint_id?: string | null
  started_at: string | null
  completed_at: string | null
  status: string
  assigned_to: number | null
  created_by: number
  created_by_role?: string
  created_by_email?: string
  created_at: string
  updated_at: string
  submissions?: Submission[]
}

const route = useRoute()
const { fetchWithAuth, currentUser: authCurrentUser } = useAuth()
const authStore = useAuthStore()
const projectsApi = useProjectsApi()

// State
const task = ref<Task | null>(null)
const isLoading = ref(true)
const error = ref('')

// Sprint context: task IDs in order (for Prev/Next when from_sprint + from_project)
const sprintTaskIds = ref<string[]>([])
const sprintNavMessage = ref('')

// Edit Task State
const showEditModal = ref(false)
const editForm = ref({
  title: '',
  description: '',
  deadline: '',
  priority: 'MEDIUM',
  story_points: 0,
  sprint_id: '',
  start_date: '',
  end_date: ''
})
const editSprints = ref<{ id: string; name: string }[]>([])
const isUpdatingTask = ref(false)
const editError = ref('')

// Delete Task State
const showDeleteModal = ref(false)
const isDeletingTask = ref(false)
const deleteError = ref('')

// Inline Description Edit State
const isEditingDescription = ref(false)
const isSavingDescription = ref(false)
const inlineDescriptionHtml = ref('')

// Comments & Time Logs State
const comments = ref<TaskComment[]>([])
const timeLogs = ref<TimeLog[]>([])
const commentsLoading = ref(false)
const timeLogsLoading = ref(false)

// Parsed slide resource URLs (for Google Slides imported tasks; images are now in description, only metadata for "Open in Slides")
const slideResourceURLs = computed(() => {
  if (!task.value?.resource_urls) return null
  const ru = typeof task.value.resource_urls === 'string'
    ? JSON.parse(task.value.resource_urls)
    : task.value.resource_urls
  if (ru?.source !== 'google_slides') return null
  return ru as {
    thumbnail_url: string
    images: string[]
    slide_url: string
    source: string
    slide_index: number
    presentation_id: string
    comments: Array<{ content: string; author: string; resolved: boolean }>
  }
})

// Open-in-Slides URL: use stored URL if it has #slide= fragment, else build with slide_index so we open to the right slide
const slideOpenInSlidesURL = computed(() => {
  const ru = slideResourceURLs.value
  if (!ru?.slide_url) return ''
  if (ru.slide_url.includes('#slide=')) return ru.slide_url
  if (ru.presentation_id && ru.slide_index != null)
    return `https://docs.google.com/presentation/d/${ru.presentation_id}/edit#slide=${ru.slide_index}`
  return ru.slide_url
})

// Single source: store first, then JWT payload (so Edit works right after login)
const effectiveUser = computed(() => {
  if (authStore.user) return authStore.user
  const payload = authCurrentUser.value
  if (!payload) return null
  const id = payload.user_id ?? (payload as any).userId
  return id != null || payload.role ? { id: Number(id) || 0, role: payload.role || '', email: payload.email || '' } : null
})

const canEditOrDelete = computed(() => {
  if (!task.value || !effectiveUser.value) return false
  const user = effectiveUser.value
  const role = (user.role || '').trim().toUpperCase()
  if (role === 'CEO' || role === 'PM') return true
  const creatorId = Number(task.value.created_by)
  const userId = Number(user.id ?? authStore.userId ?? 0)
  return creatorId === userId && !Number.isNaN(userId)
})

const creatorLabel = computed(() => {
  if (!task.value) return ''
  const role = task.value.created_by_role
  const email = task.value.created_by_email
  if (role === 'CEO') return email ? `CEO (${email})` : 'CEO'
  if (role === 'PM') return email ? `PM (${email})` : 'PM'
  return `Dev #${task.value.created_by}`
})

// Methods
const fetchTask = async () => {
  const taskId = (route.params.id as string)?.trim?.() || ''

  if (!taskId || taskId === 'undefined' || taskId === 'null') {
    error.value = 'Invalid or missing task ID. Check the URL.'
    isLoading.value = false
    return
  }

  try {
    isLoading.value = true
    error.value = ''

    const response = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${taskId}`)
    task.value = response.data

    // Load sprint task order for Prev/Next when opened from sprint page
    // Use task.project_id (UUID) for API — from_project in URL may be project code (e.g. mims-hdmap), API requires UUID
    const fromSprint = route.query.from_sprint as string
    if (fromSprint && task.value?.project_id) {
      const tasksApi = useTasksApi()
      try {
        const projectTasks = await tasksApi.getTasksByProject(task.value.project_id)
        const inSprint = projectTasks.filter((t) => t.sprint_id === fromSprint)
        inSprint.sort((a, b) => (a.code || a.id).localeCompare(b.code || b.id))
        sprintTaskIds.value = inSprint.map((t) => t.id)
      } catch {
        sprintTaskIds.value = []
      }
    } else {
      sprintTaskIds.value = []
    }
  } catch (err: any) {
    console.error('Failed to fetch task:', err)
    const apiMsg = err?.data?.message ?? err?.data?.error
    const status = err?.statusCode ?? err?.status
    error.value = apiMsg || (status === 404 ? 'Task not found.' : err?.message || 'Failed to load task.')
  } finally {
    isLoading.value = false
  }
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    'COMPLETED': 'bg-green-700 text-green-100 border border-green-500',
    'IN_PROGRESS': 'bg-blue-700 text-blue-100 border border-blue-500',
    'PENDING': 'bg-yellow-700 text-yellow-100 border border-yellow-500',
    'BLOCKED': 'bg-red-700 text-red-100 border border-red-500',
    'REVIEW_PENDING': 'bg-indigo-900 text-indigo-200 border border-indigo-600' // 🚦 Quality Gate
  }
  return classes[status] || 'bg-gray-700 text-gray-100 border border-gray-500'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'COMPLETED': '✅ COMPLETED',
    'IN_PROGRESS': '🔄 IN PROGRESS',
    'PENDING': '⏳ PENDING',
    'BLOCKED': '🚫 BLOCKED',
    'REVIEW_PENDING': '⏳ WAITING FOR APPROVAL', // 🚦 Quality Gate
    'ASSIGNED': '📌 ASSIGNED'
  }
  return labels[status] || status
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Deadline Helpers
const getDeadlineUrgency = (task: Task) => {
  if (!task.due_at || task.status === 'COMPLETED') return 'none'
  
  const now = new Date().getTime()
  const dueDate = new Date(task.due_at).getTime()
  const hoursUntilDue = (dueDate - now) / (1000 * 60 * 60)
  
  if (hoursUntilDue < 0) return 'overdue'
  if (hoursUntilDue < 24) return 'urgent'
  return 'normal'
}

const getDeadlineCountdown = (dueAt: string) => {
  const now = new Date().getTime()
  const due = new Date(dueAt).getTime()
  const diff = due - now
  
  if (diff < 0) {
    // Overdue
    const hours = Math.abs(Math.floor(diff / (1000 * 60 * 60)))
    const days = Math.floor(hours / 24)
    if (days > 0) return `Overdue by ${days} days`
    return `Overdue by ${hours} hours`
  }
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `${days} days left`
  if (hours > 0) return `${hours} hours left`
  return 'Due very soon!'
}

const calculateDuration = (startAt: string, completedAt: string) => {
  const start = new Date(startAt).getTime()
  const end = new Date(completedAt).getTime()
  const diff = end - start
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const days = Math.floor(hours / 24)
  
  if (days > 0) {
    const remainingHours = hours % 24
    return `${days}d ${remainingHours}h`
  }
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes}m`
}

const toDatetimeLocal = (iso: string | null | undefined) => {
  if (!iso) return ''
  const d = new Date(iso)
  return isNaN(d.getTime()) ? '' : d.toISOString().slice(0, 16)
}

// Edit Task Methods
const openEditModal = async () => {
  if (!task.value) return
  
  // Pre-fill form with current values (same fields as Create Task)
  editForm.value.title = task.value.title
  editForm.value.description = task.value.description || ''
  editForm.value.deadline = toDatetimeLocal(task.value.due_at)
  editForm.value.priority = task.value.priority || 'MEDIUM'
  editForm.value.story_points = task.value.story_points ?? 0
  editForm.value.sprint_id = task.value.sprint_id ?? ''
  editForm.value.start_date = toDatetimeLocal(task.value.start_date)
  editForm.value.end_date = toDatetimeLocal(task.value.end_date)
  
  // Load sprints for Sprint dropdown when task has project
  editSprints.value = []
  if (task.value.project_id) {
    try {
      const list = await projectsApi.getSprints(task.value.project_id)
      editSprints.value = list.map((s) => ({ id: s.id, name: s.name }))
    } catch {
      // ignore
    }
  }
  
  editError.value = ''
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.value = {
    title: '',
    description: '',
    deadline: '',
    priority: 'MEDIUM',
    story_points: 0,
    sprint_id: '',
    start_date: '',
    end_date: ''
  }
  editError.value = ''
}

const submitEdit = async () => {
  if (!task.value) return
  
  // Validation
  if (!editForm.value.title.trim()) {
    editError.value = 'Title is required'
    return
  }
  
  try {
    isUpdatingTask.value = true
    editError.value = ''
    
    const taskId = route.params.id as string
    
    // Prepare request body (only send changed fields) — same fields as Create Task
    const body: Record<string, string | number> = {}
    
    if (editForm.value.title !== task.value.title) {
      body.title = editForm.value.title
    }
    if (editForm.value.description !== (task.value.description || '')) {
      body.description = editForm.value.description
    }
    if (editForm.value.priority && editForm.value.priority !== (task.value.priority || 'MEDIUM')) {
      body.priority = editForm.value.priority
    }
    const currentSp = task.value.story_points ?? 0
    if (Number(editForm.value.story_points) !== currentSp) {
      body.story_points = Number(editForm.value.story_points) || 0
    }
    const currentSprint = task.value.sprint_id ?? ''
    if (editForm.value.sprint_id !== currentSprint) {
      body.sprint_id = editForm.value.sprint_id || ''
    }
    const currentStart = toDatetimeLocal(task.value.start_date)
    if (editForm.value.start_date !== currentStart && editForm.value.start_date) {
      body.start_date = new Date(editForm.value.start_date).toISOString()
    }
    const currentEnd = toDatetimeLocal(task.value.end_date)
    const newEnd = editForm.value.end_date || editForm.value.deadline
    const newEndStr = newEnd ? new Date(newEnd).toISOString() : ''
    if (newEnd && toDatetimeLocal(newEndStr) !== currentEnd) {
      body.end_date = newEndStr
    }
    
    if (Object.keys(body).length === 0) {
      editError.value = 'No changes detected. Please modify at least one field.'
      isUpdatingTask.value = false
      return
    }
    
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'PATCH',
      body
    })
    
    // Success
    const hasContentChange = body.title !== undefined || body.description !== undefined
    alert(hasContentChange ? 'Task updated.' : 'Task updated.')
    
    closeEditModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to update task:', err)
    editError.value = err.data?.message || err.message || 'Failed to update task'
  } finally {
    isUpdatingTask.value = false
  }
}

// Delete Task Methods
const openDeleteConfirmation = () => {
  deleteError.value = ''
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deleteError.value = ''
}

const confirmDelete = async () => {
  if (!task.value) return
  
  try {
    isDeletingTask.value = true
    deleteError.value = ''
    
    const taskId = route.params.id as string
    
    // Call DELETE API
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'DELETE'
    })
    
    // Success! Show notification and redirect to back target (e.g. project backlog)
    alert('Task deleted.')
    goToDashboard()
  } catch (err: any) {
    console.error('Failed to delete task:', err)
    deleteError.value = err.data?.message || err.message || 'Failed to delete task'
  } finally {
    isDeletingTask.value = false
  }
}

// Inline Description Edit
const startInlineEdit = () => {
  inlineDescriptionHtml.value = task.value?.description || ''
  isEditingDescription.value = true
}

const cancelInlineEdit = () => {
  isEditingDescription.value = false
  inlineDescriptionHtml.value = ''
}

const saveInlineDescription = async () => {
  if (!task.value) return
  isSavingDescription.value = true
  try {
    const taskId = route.params.id as string
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'PATCH',
      body: { description: inlineDescriptionHtml.value }
    })
    task.value.description = inlineDescriptionHtml.value
    isEditingDescription.value = false
    inlineDescriptionHtml.value = ''
  } catch (err: any) {
    console.error('Failed to save description:', err)
  } finally {
    isSavingDescription.value = false
  }
}

// Back: if came from project page (from_project + from_tab) or sprint page (from_sprint + from_project), return there; else project or dashboard
const backTarget = computed(() => {
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  const fromTab = route.query.from_tab as string
  if (fromSprint && fromProject) {
    return `/projects/sprint/${fromSprint}?project=${encodeURIComponent(fromProject)}`
  }
  if (fromProject) {
    const tab = fromTab || 'backlog'
    return `/projects/${fromProject}?tab=${tab}`
  }
  const t = task.value
  if (!t) return '/dashboard'
  if (t.code) {
    const projectCode = t.code.replace(/-[0-9]+$/, '')
    if (projectCode !== t.code) return `/projects/${projectCode}`
  }
  if (t.project_id) return `/projects/${t.project_id}`
  return '/dashboard'
})

// Whether we're in sprint context (opened from sprint page) — show Prev/Next area
const inSprintContext = computed(() => {
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  return !!(fromSprint && fromProject)
})

// Prev/Next task links within same sprint (only when from_sprint + from_project)
// Normalize id comparison (UUID may differ in casing between getTask and getTasksByProject)
const currentSprintIndex = computed(() => {
  const t = task.value
  if (!t?.id || !sprintTaskIds.value.length) return -1
  const needle = String(t.id).toLowerCase()
  return sprintTaskIds.value.findIndex((id) => String(id).toLowerCase() === needle)
})
const prevTaskLink = computed(() => {
  if (currentSprintIndex.value <= 0) return null
  const id = sprintTaskIds.value[currentSprintIndex.value - 1]
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  if (!id || !fromSprint || !fromProject) return null
  return `/task/${id}?from_sprint=${encodeURIComponent(fromSprint)}&from_project=${encodeURIComponent(fromProject)}`
})
const nextTaskLink = computed(() => {
  const idx = currentSprintIndex.value
  if (idx < 0 || idx >= sprintTaskIds.value.length - 1) return null
  const id = sprintTaskIds.value[idx + 1]
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  if (!id || !fromSprint || !fromProject) return null
  return `/task/${id}?from_sprint=${encodeURIComponent(fromSprint)}&from_project=${encodeURIComponent(fromProject)}`
})

function showSprintNavFeedback(msg: string) {
  sprintNavMessage.value = msg
  setTimeout(() => { sprintNavMessage.value = '' }, 2500)
}

function goToPrevTask() {
  const link = prevTaskLink.value
  if (link) {
    navigateTo(link)
  } else {
    showSprintNavFeedback('ไม่มี task ก่อนหน้า')
  }
}

function goToNextTask() {
  const link = nextTaskLink.value
  if (link) {
    navigateTo(link)
  } else {
    showSprintNavFeedback('ไม่มี task ถัดไป')
  }
}

const goToDashboard = () => {
  navigateTo(backTarget.value)
}

async function fetchCommentsAndLogs() {
  const taskId = route.params.id as string
  if (!taskId) return
  const tasksApi = useTasksApi()
  commentsLoading.value = true
  timeLogsLoading.value = true
  try {
    const [c, l] = await Promise.all([
      tasksApi.getComments(taskId),
      tasksApi.getTimeLogs(taskId),
    ])
    comments.value = c
    timeLogs.value = l
  } catch {
    // non-critical
  } finally {
    commentsLoading.value = false
    timeLogsLoading.value = false
  }
}

async function handleAddComment(content: string) {
  const taskId = route.params.id as string
  const tasksApi = useTasksApi()
  commentsLoading.value = true
  try {
    const comment = await tasksApi.addComment(taskId, content)
    comments.value.push(comment)
  } catch (e: any) {
    console.error('Failed to add comment:', e)
  } finally {
    commentsLoading.value = false
  }
}

async function handleLogTime(minutes: number, description: string) {
  const taskId = route.params.id as string
  const tasksApi = useTasksApi()
  timeLogsLoading.value = true
  try {
    const log = await tasksApi.logTime(taskId, minutes, description)
    timeLogs.value.unshift(log)
  } catch (e: any) {
    console.error('Failed to log time:', e)
  } finally {
    timeLogsLoading.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchTask()
  fetchCommentsAndLogs()
})
</script>
