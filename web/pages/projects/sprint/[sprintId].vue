<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลด sprint...</p>
    </div>

    <div v-else-if="error" class="p-8 max-w-2xl mx-auto">
      <div class="bg-red-900/20 border border-red-500 rounded-xl p-6 text-red-400">
        <h2 class="font-bold text-lg mb-1">Failed to load sprint</h2>
        <p class="text-sm">{{ error }}</p>
        <NuxtLink v-if="projectId" :to="`/projects/${projectId}`" class="mt-4 inline-block text-sm text-gray-400 hover:text-white">← Back to Project</NuxtLink>
        <NuxtLink v-else to="/projects" class="mt-4 inline-block text-sm text-gray-400 hover:text-white">← Back to Projects</NuxtLink>
      </div>
    </div>

    <div v-else-if="project && sprint" class="p-3 sm:p-6">
      <!-- Breadcrumb & header -->
      <div class="border-b border-gray-800 pb-4 mb-6">
        <div class="flex flex-wrap items-center gap-2 text-sm text-gray-400 mb-2">
          <NuxtLink to="/projects" class="hover:text-white transition-colors">Projects</NuxtLink>
          <span>/</span>
          <NuxtLink :to="`/projects/${projectId}`" class="hover:text-white transition-colors truncate">{{ project.name }}</NuxtLink>
          <span>/</span>
          <span class="text-gray-200 font-medium truncate">{{ sprint.name }}</span>
        </div>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-bold text-white">{{ sprint.name }}</h1>
            <p v-if="sprint.goal" class="text-sm text-gray-400 mt-1">{{ sprint.goal }}</p>
            <div class="flex items-center gap-3 mt-2">
              <span
                class="px-2 py-0.5 text-xs font-semibold rounded-full"
                :class="sprint.status === 'ACTIVE' ? 'bg-indigo-500/20 text-indigo-400' : sprint.status === 'COMPLETED' ? 'bg-gray-600 text-gray-400' : 'bg-yellow-500/20 text-yellow-400'"
              >
                {{ sprint.status }}
              </span>
              <span v-if="sprint.start_date || sprint.end_date" class="text-xs text-gray-500">
                {{ formatDate(sprint.start_date) }} – {{ formatDate(sprint.end_date) }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button @click="openCreateTaskModal()" class="btn-primary-sm">+ Add Task</button>
            <button @click="openImportModal()" class="btn-import-sm">
              <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
              Import Slides
            </button>
            <NuxtLink :to="`/projects/${projectId}`" class="btn-ghost-sm">← Project</NuxtLink>
          </div>
        </div>
      </div>

      <!-- Sprint stats -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
        <div class="metric-card">
          <div class="text-2xl font-bold text-gray-200">{{ sprintTasks.length }}</div>
          <div class="metric-label">Tasks</div>
        </div>
        <div class="metric-card">
          <div class="text-2xl font-bold text-green-400">{{ doneCount }}</div>
          <div class="metric-label">Done</div>
        </div>
        <div class="metric-card">
          <div class="text-2xl font-bold text-indigo-400">{{ totalSp }}</div>
          <div class="metric-label">Story points</div>
        </div>
        <div class="metric-card">
          <div class="text-2xl font-bold" :class="overdueCount > 0 ? 'text-red-400' : 'text-gray-400'">{{ overdueCount }}</div>
          <div class="metric-label">Overdue</div>
        </div>
      </div>

      <!-- Task list -->
      <div class="card">
        <div class="flex flex-wrap items-center justify-between gap-3 mb-4">
          <h3 class="section-title mb-0">Tasks in this sprint</h3>
          <div v-if="sprintTasks.length" class="flex items-center gap-3">
            <label class="flex items-center gap-2 cursor-pointer select-none text-sm text-gray-400 hover:text-gray-200">
              <input
                ref="checkAllInputRef"
                type="checkbox"
                :checked="isCheckAllChecked"
                class="w-4 h-4 rounded border-gray-500 bg-gray-700 text-indigo-500 focus:ring-indigo-500"
                @change="toggleCheckAll"
              />
              <span>Check all</span>
            </label>
            <button
              v-if="selectedTaskIds.length > 0"
              type="button"
              class="px-3 py-1.5 text-sm font-medium text-red-300 bg-red-900/30 hover:bg-red-900/50 border border-red-800 rounded-lg transition-colors flex items-center gap-2"
              :disabled="isDeletingTasks"
              @click="confirmDeleteSelected"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
              Delete ({{ selectedTaskIds.length }})
            </button>
          </div>
        </div>
        <div v-if="sprintTasks.length" class="space-y-2">
          <div
            v-for="t in sprintTasks"
            :key="t.id"
            class="flex items-center gap-3 py-3 px-4 rounded-lg hover:bg-gray-700/40 transition-colors border-b border-gray-700/50 last:border-0"
          >
            <label class="flex items-center shrink-0 cursor-pointer" @click.stop>
              <input
                type="checkbox"
                :checked="selectedTaskIds.includes(t.id)"
                class="w-4 h-4 rounded border-gray-500 bg-gray-700 text-indigo-500 focus:ring-indigo-500"
                @change="toggleTaskSelection(t.id)"
              />
            </label>
            <div
              class="flex flex-1 items-center justify-between min-w-0 cursor-pointer"
              @click="navigateToTask(t.id)"
            >
              <div class="flex items-center gap-3 min-w-0">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ t.code }}</span>
                <span class="text-sm text-gray-200 truncate">{{ t.title }}</span>
                <span class="px-1.5 py-0.5 text-[10px] rounded font-medium shrink-0" :class="priorityBadge(t.priority)">{{ t.priority }}</span>
              </div>
              <div class="flex items-center gap-3 shrink-0">
                <span v-if="t.story_points" class="text-xs text-purple-400">{{ t.story_points }} SP</span>
                <span class="text-xs px-2 py-0.5 rounded-full" :class="taskStatusBadge(t.status)">{{ t.status.replace('_', ' ') }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-12 text-gray-500 text-sm">
          <p class="mb-4">No tasks in this sprint yet.</p>
          <button @click="openCreateTaskModal()" class="btn-primary-sm">+ Add first task</button>
        </div>
      </div>
    </div>

    <!-- Import from Google Slides Modal -->
    <div v-if="showImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeImportModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-xl w-full shadow-2xl">

        <!-- Header -->
        <div class="flex items-center justify-between mb-5">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center">
              <svg class="w-4 h-4 text-indigo-400" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import from Google Slides</h2>
              <p class="text-xs text-gray-400">สร้าง task อัตโนมัติจากแต่ละ slide</p>
            </div>
          </div>
          <button @click="closeImportModal" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>

        <!-- Result state -->
        <div v-if="importStep === 'result' && importResult" class="space-y-4">
          <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
            <div class="flex items-center gap-2 mb-2">
              <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
              <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
            </div>
            <p class="text-gray-300 text-sm font-medium mb-1">{{ importResult.presentation_title }}</p>
            <p class="text-gray-400 text-xs">สร้าง {{ importResult.created_count }} tasks จาก {{ importResult.slide_count }} slides</p>
          </div>
          <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
            <div
              v-for="task in importResult.tasks"
              :key="task.id"
              class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm"
            >
              <span class="text-xs font-mono text-gray-500 shrink-0">{{ task.code }}</span>
              <span class="text-gray-200 truncate">{{ task.title }}</span>
            </div>
          </div>
          <button @click="closeImportModal" class="w-full btn-primary py-2.5">Done</button>
        </div>

        <!-- Step 2: Select slides to import -->
        <div v-else-if="importStep === 'select' && importPreview" class="space-y-4">
          <div class="p-3 bg-gray-700/40 rounded-xl">
            <p class="text-sm font-medium text-white">{{ importPreview.presentation_title }}</p>
            <p class="text-xs text-gray-500 mt-0.5">
            เลือก slide ที่จะ import ({{ importSelectedIndices.length }} / {{ importPreview.slides.length }})
            <span v-if="(importPreview.already_imported_slide_indices?.length ?? 0) > 0">— หน้าที่นำเข้าแล้วจะถูก uncheck ไว้ เลือกเฉพาะหน้าที่ต้องการเพิ่มได้</span>
          </p>
          </div>
          <div class="flex items-center gap-2 flex-wrap">
            <button type="button" @click="importSelectAll" class="btn-ghost-sm">เลือกทั้งหมด</button>
            <button type="button" @click="importDeselectAll" class="btn-ghost-sm">ยกเลิกทั้งหมด</button>
            <button
              type="button"
              @click="importSelectOnlyNew"
              class="btn-ghost-sm text-indigo-400"
            >
              เลือกเฉพาะที่ยังไม่เคยนำเข้า
            </button>
          </div>
          <div class="max-h-56 overflow-y-auto space-y-1.5 pr-1 border border-gray-700/60 rounded-xl p-2">
            <label
              v-for="s in importPreview.slides"
              :key="s.index"
              class="flex items-center gap-3 py-2 px-2 rounded-lg hover:bg-gray-700/40 cursor-pointer"
              :class="{ 'opacity-70': s.hidden }"
            >
              <input
                v-model="importSelectedIndices"
                type="checkbox"
                :value="s.index"
                class="rounded border-gray-500 bg-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-xs text-gray-400 w-8 shrink-0">#{{ s.index }}</span>
              <span class="text-sm text-gray-200 truncate flex-1">{{ s.title || '(ไม่มีชื่อ)' }}</span>
              <span v-if="s.hidden" class="text-xs text-amber-400/90 shrink-0">ซ่อน</span>
              <span v-else-if="(importPreview.already_imported_slide_indices || []).includes(s.index)" class="text-xs text-gray-500 shrink-0">นำเข้าแล้ว</span>
            </label>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Priority ของทุก task</label>
              <select v-model="importForm.priority" class="input-field w-full" :disabled="isImporting">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points (ต่อ task)</label>
              <input v-model.number="importForm.story_points" type="number" min="0" class="input-field w-full" placeholder="1" :disabled="isImporting" />
            </div>
          </div>
          <div v-if="importError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ importError }}</div>
          <div class="flex gap-3">
            <button
              @click="submitImport"
              :disabled="isImporting || importSelectedIndices.length === 0"
              class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
            >
              <svg v-if="isImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
              {{ isImporting ? 'กำลัง import...' : `Import ${importSelectedIndices.length} Slides` }}
            </button>
            <button type="button" @click="importStep = 'form'" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">กลับ</button>
          </div>
        </div>

        <!-- Step 1: Form (URL + Load slides) -->
        <div v-else class="space-y-4">
          <div>
            <label class="label">Google Slides URL *</label>
            <input
              v-model="importForm.presentation_url"
              type="url"
              class="input-field w-full"
              placeholder="https://docs.google.com/presentation/d/..."
              :disabled="isLoadingPreview"
            />
            <p class="text-xs text-gray-500 mt-1">ต้องเปิดสิทธิ์ "Anyone with the link can view"</p>
          </div>

          <div v-if="importError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ importError }}</div>

          <div class="flex gap-3 mt-1">
            <button
              type="button"
              @click="loadImportPreview"
              :disabled="isLoadingPreview || !importForm.presentation_url.trim()"
              class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
            >
              <svg v-if="isLoadingPreview" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
              {{ isLoadingPreview ? 'กำลังโหลด...' : 'โหลดรายการ slide' }}
            </button>
            <button @click="closeImportModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
          </div>
        </div>

      </div>
    </div>

    <!-- Create Task Modal (same structure as project page) -->
    <div v-if="showCreateTaskModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeCreateTaskModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">Add Task</h2>
          <button @click="closeCreateTaskModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="label">Title *</label>
            <input v-model="createTaskForm.title" type="text" class="input-field w-full" placeholder="Task title..." />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="createTaskForm.description" rows="3" class="input-field w-full resize-none" placeholder="Describe the task..."></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Priority</label>
              <select v-model="createTaskForm.priority" class="input-field w-full">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points</label>
              <input v-model.number="createTaskForm.story_points" type="number" min="0" class="input-field w-full" placeholder="0" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Sprint</label>
              <select v-model="createTaskForm.sprint_id" class="input-field w-full" disabled>
                <option v-if="sprint" :value="sprint.id">{{ sprint.name }}</option>
              </select>
            </div>
            <div>
              <label class="label">Due Date</label>
              <input v-model="createTaskForm.due_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Start Date</label>
              <input v-model="createTaskForm.start_date" type="datetime-local" class="input-field w-full" />
            </div>
            <div>
              <label class="label">End Date</label>
              <input v-model="createTaskForm.end_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <div v-if="createTaskError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ createTaskError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="submitCreateTask" :disabled="isCreatingTask || !createTaskForm.title.trim()" class="flex-1 btn-primary py-2.5 disabled:opacity-40">
            {{ isCreatingTask ? 'Creating...' : 'Create Task' }}
          </button>
          <button @click="closeCreateTaskModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete selected tasks confirmation -->
    <div v-if="showDeleteConfirmModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeDeleteConfirmModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h3 class="text-lg font-bold text-white mb-2">Delete {{ selectedTaskIds.length }} task(s)?</h3>
        <p class="text-sm text-gray-400 mb-4">This cannot be undone.</p>
        <div v-if="deleteError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm mb-4">{{ deleteError }}</div>
        <div class="flex gap-3">
          <button
            type="button"
            class="flex-1 px-4 py-2.5 bg-red-600 hover:bg-red-700 text-white font-medium rounded-xl transition-colors disabled:opacity-50"
            :disabled="isDeletingTasks"
            @click="deleteSelectedTasks"
          >
            {{ isDeletingTasks ? 'Deleting...' : 'Delete' }}
          </button>
          <button
            type="button"
            class="px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors"
            :disabled="isDeletingTasks"
            @click="closeDeleteConfirmModal"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { Project, Sprint, Task } from '~/core/modules/projects/infrastructure/projects-api'

definePageMeta({ layout: 'default', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()

// Route: /projects/sprint/:sprintId?project=:projectIdOrCode (project required for loading)
const projectId = computed(() => (route.query.project as string) || '')
const sprintId = computed(() => route.params.sprintId as string)

// Ensure browser Back from sprint goes to project (fix duplicate history or direct-open)
const SPRINT_BACK_STATE = 'sprint-back-to-project'
onMounted(() => {
  if (!projectId.value || typeof window === 'undefined') return
  const fullPath = route.fullPath
  const projectPath = `/projects/${projectId.value}`

  const onPopState = () => {
    const state = window.history.state as { key?: string; projectId?: string } | null
    if (state?.key === SPRINT_BACK_STATE && state?.projectId) {
      navigateTo(`/projects/${state.projectId}`)
      return
    }
    // If back landed on sprint URL again (duplicate history), go to project
    if (typeof window !== 'undefined' && window.location.pathname.includes('/projects/sprint/') && projectId.value) {
      navigateTo(`/projects/${projectId.value}`)
    }
  }
  window.addEventListener('popstate', onPopState)
  onUnmounted(() => window.removeEventListener('popstate', onPopState))

  // If user opened this sprint directly (e.g. new tab / bookmark), history has only this entry.
  // Push project then current URL so Back goes to project.
  if (window.history.length === 1) {
    window.history.pushState(
      { key: SPRINT_BACK_STATE, projectId: projectId.value },
      '',
      projectPath
    )
    window.history.pushState({}, '', fullPath)
  }
})

const project = ref<Project | null>(null)
const sprint = ref<Sprint | null>(null)
const allTasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')

const sprintTasks = computed(() => allTasks.value.filter((t) => t.sprint_id === sprintId.value))
const doneCount = computed(() => sprintTasks.value.filter((t) => t.status === 'COMPLETED').length)
const totalSp = computed(() => sprintTasks.value.reduce((s, t) => s + (t.story_points || 0), 0))
const overdueCount = computed(() => {
  const now = Date.now()
  return sprintTasks.value.filter((t) => t.status !== 'COMPLETED' && t.due_at && new Date(t.due_at).getTime() < now).length
})

function formatDate(d: string | null) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function priorityBadge(p: string) {
  if (p === 'CRITICAL') return 'bg-red-500/20 text-red-400'
  if (p === 'HIGH') return 'bg-orange-500/20 text-orange-400'
  if (p === 'MEDIUM') return 'bg-yellow-500/20 text-yellow-400'
  return 'bg-green-500/20 text-green-400'
}

function taskStatusBadge(status: string) {
  if (status === 'COMPLETED') return 'bg-green-500/20 text-green-400'
  if (status === 'IN_PROGRESS') return 'bg-blue-500/20 text-blue-400'
  if (status === 'REVIEW_PENDING') return 'bg-yellow-500/20 text-yellow-400'
  if (status === 'BLOCKED') return 'bg-red-500/20 text-red-400'
  return 'bg-gray-700 text-gray-400'
}

function navigateToTask(id: string) {
  router.push({
    path: `/task/${id}`,
    query: { from_sprint: sprintId.value, from_project: projectId.value }
  })
}

// Bulk delete: selection state
const selectedTaskIds = ref<string[]>([])
const checkAllInputRef = ref<HTMLInputElement | null>(null)
const isCheckAllChecked = computed(() => sprintTasks.value.length > 0 && selectedTaskIds.value.length === sprintTasks.value.length)
const isCheckAllIndeterminate = computed(() => {
  const n = selectedTaskIds.value.length
  return n > 0 && n < sprintTasks.value.length
})
watch([selectedTaskIds, sprintTasks], () => {
  const el = checkAllInputRef.value
  if (el) (el as HTMLInputElement).indeterminate = isCheckAllIndeterminate.value
}, { immediate: true })
function toggleTaskSelection(id: string) {
  const idx = selectedTaskIds.value.indexOf(id)
  if (idx === -1) selectedTaskIds.value = [...selectedTaskIds.value, id]
  else selectedTaskIds.value = selectedTaskIds.value.filter((x) => x !== id)
}
function toggleCheckAll() {
  if (isCheckAllChecked.value) selectedTaskIds.value = []
  else selectedTaskIds.value = sprintTasks.value.map((t) => t.id)
}

// Delete selected: confirmation modal + API
const showDeleteConfirmModal = ref(false)
const isDeletingTasks = ref(false)
const deleteError = ref('')
function confirmDeleteSelected() {
  if (selectedTaskIds.value.length === 0) return
  deleteError.value = ''
  showDeleteConfirmModal.value = true
}
function closeDeleteConfirmModal() {
  showDeleteConfirmModal.value = false
  deleteError.value = ''
}
async function deleteSelectedTasks() {
  if (selectedTaskIds.value.length === 0) return
  isDeletingTasks.value = true
  deleteError.value = ''
  try {
    for (const id of selectedTaskIds.value) {
      await tasksApi.deleteTask(id)
    }
    allTasks.value = allTasks.value.filter((t) => !selectedTaskIds.value.includes(t.id))
    selectedTaskIds.value = []
    closeDeleteConfirmModal()
  } catch (e: any) {
    deleteError.value = e?.message ?? 'Failed to delete tasks'
  } finally {
    isDeletingTasks.value = false
  }
}

async function loadAll() {
  if (!projectId.value) {
    error.value = 'Project not specified. Open this sprint from a project page.'
    isLoading.value = false
    return
  }
  isLoading.value = true
  error.value = ''
  try {
    const p = await projectsApi.getProject(projectId.value)
    project.value = p
    const [sprints, tasks] = await Promise.all([
      projectsApi.getSprints(p.id),
      tasksApi.getTasksByProject(p.id),
    ])
    const s = sprints.find((x) => x.id === sprintId.value)
    if (!s) {
      error.value = 'Sprint not found'
      return
    }
    sprint.value = s
    allTasks.value = tasks
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load sprint'
  } finally {
    isLoading.value = false
  }
}

// Create Task (sprint pre-selected; same form shape as project page)
const showCreateTaskModal = ref(false)
const createTaskForm = ref({
  title: '',
  description: '',
  priority: 'MEDIUM' as const,
  story_points: 0,
  sprint_id: '',
  due_date: '',
  start_date: '',
  end_date: '',
})
const isCreatingTask = ref(false)
const createTaskError = ref('')

function openCreateTaskModal() {
  createTaskForm.value = {
    title: '',
    description: '',
    priority: 'MEDIUM',
    story_points: 0,
    sprint_id: sprint.value?.id ?? '',
    due_date: '',
    start_date: '',
    end_date: '',
  }
  createTaskError.value = ''
  showCreateTaskModal.value = true
}

function closeCreateTaskModal() {
  showCreateTaskModal.value = false
}

async function submitCreateTask() {
  if (!project.value || !sprint.value) return
  isCreatingTask.value = true
  createTaskError.value = ''
  try {
    const payload: any = {
      title: createTaskForm.value.title,
      description: createTaskForm.value.description,
      priority: createTaskForm.value.priority,
      story_points: createTaskForm.value.story_points,
      project_id: project.value.id,
      sprint_id: sprint.value.id,
    }
    if (createTaskForm.value.due_date) payload.due_date = new Date(createTaskForm.value.due_date).toISOString()
    if (createTaskForm.value.start_date) payload.start_date = new Date(createTaskForm.value.start_date).toISOString()
    if (createTaskForm.value.end_date) payload.end_date = new Date(createTaskForm.value.end_date).toISOString()
    const task = await tasksApi.createTask(payload)
    allTasks.value.unshift(task)
    closeCreateTaskModal()
  } catch (e: any) {
    createTaskError.value = e?.message ?? 'Failed to create task'
  } finally {
    isCreatingTask.value = false
  }
}

// Import from Google Slides
const showImportModal = ref(false)
const importStep = ref<'form' | 'select' | 'result'>('form')
const isImporting = ref(false)
const isLoadingPreview = ref(false)
const importError = ref('')
const importResult = ref<{ created_count: number; slide_count: number; presentation_title: string; tasks: any[] } | null>(null)
const importPreview = ref<{
  presentation_title: string
  slides: { index: number; title: string; hidden?: boolean }[]
  import_mode?: string
  api_key_status?: string
  api_key_error?: string
} | null>(null)
const importSelectedIndices = ref<number[]>([])
const importForm = ref({
  presentation_url: '',
  api_key: '',
  priority: 'MEDIUM' as const,
  story_points: 1,
})

function openImportModal() {
  importForm.value = { presentation_url: '', api_key: '', priority: 'MEDIUM', story_points: 1 }
  importStep.value = 'form'
  importError.value = ''
  importResult.value = null
  importPreview.value = null
  importSelectedIndices.value = []
  showImportModal.value = true
}

function closeImportModal() {
  showImportModal.value = false
  if (importResult.value) {
    loadAll()
  }
}

async function loadImportPreview() {
  if (!importForm.value.presentation_url.trim()) return
  isLoadingPreview.value = true
  importError.value = ''
  try {
    const data = await tasksApi.previewGoogleSlides({
      presentation_url: importForm.value.presentation_url.trim(),
    })
    importPreview.value = data
    // By default select only NEW slides (not yet imported) and non-hidden; so re-importing same link can add only new pages
    const alreadySet = new Set(data.already_imported_slide_indices ?? [])
    importSelectedIndices.value = data.slides
      .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
      .map((s: { index: number }) => s.index)
    importStep.value = 'select'
  } catch (e: any) {
    importError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
  } finally {
    isLoadingPreview.value = false
  }
}

function importSelectAll() {
  if (importPreview.value) importSelectedIndices.value = importPreview.value.slides.map((s) => s.index)
}

function importDeselectAll() {
  importSelectedIndices.value = []
}

function importSelectOnlyNew() {
  if (!importPreview.value) return
  const alreadySet = new Set(importPreview.value.already_imported_slide_indices ?? [])
  importSelectedIndices.value = importPreview.value.slides
    .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
    .map((s: { index: number }) => s.index)
}

async function submitImport() {
  if (!project.value || !sprint.value) return
  isImporting.value = true
  importError.value = ''
  try {
    const payload: any = {
      presentation_url: importForm.value.presentation_url.trim(),
      sprint_id: sprint.value.id,
      project_id: project.value.id,
      priority: importForm.value.priority,
      story_points: importForm.value.story_points,
    }
    if (importForm.value.api_key.trim()) payload.api_key = importForm.value.api_key.trim()
    if (importSelectedIndices.value.length > 0) payload.slide_indices = importSelectedIndices.value
    importResult.value = await tasksApi.importGoogleSlides(payload)
    importStep.value = 'result'
  } catch (e: any) {
    importError.value = e?.data?.message ?? e?.message ?? 'Import failed'
  } finally {
    isImporting.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.card {
  @apply bg-gray-800 border border-gray-700 rounded-xl p-5;
}
.metric-card {
  @apply bg-gray-800/60 border border-gray-700/50 rounded-xl p-4;
}
.metric-label {
  @apply text-xs text-gray-500 mt-1 uppercase tracking-wide;
}
.section-title {
  @apply text-sm font-semibold text-gray-300;
}
.label {
  @apply block text-xs text-gray-400 mb-1.5 font-medium;
}
.input-field {
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-colors;
}
.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg transition-colors;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
.btn-import-sm {
  @apply px-3 py-1.5 text-xs bg-indigo-900/50 hover:bg-indigo-800/60 border border-indigo-700/50 text-indigo-300 font-medium rounded-lg transition-colors flex items-center gap-1.5;
}
.btn-primary {
  @apply bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl transition-colors;
}
</style>
