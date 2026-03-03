<template>
  <div class="min-h-screen bg-gray-900 text-white p-6">
    <!-- Header -->
    <div class="mb-8 flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white">
          Project Command Center
        </h1>
        <p class="text-sm text-gray-400 mt-1">{{ projects.length }} active workspaces</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl shadow-lg shadow-indigo-500/20 transition-all"
        @click="openCreateModal"
      >
        <span>+</span>
        <span>New Project</span>
      </button>
    </div>

    <!-- System Metrics Row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <div class="metric-card">
        <div class="text-2xl font-bold text-indigo-400">{{ totalActive }}</div>
        <div class="metric-label">Active Projects</div>
      </div>
      <div class="metric-card">
        <div class="text-2xl font-bold text-green-400">{{ totalCompleted }}</div>
        <div class="metric-label">Completed</div>
      </div>
      <div class="metric-card">
        <div class="text-2xl font-bold text-yellow-400">{{ onHold }}</div>
        <div class="metric-label">On Hold</div>
      </div>
      <div class="metric-card">
        <div class="text-2xl font-bold text-purple-400">{{ totalProjects }}</div>
        <div class="metric-label">Total Projects</div>
      </div>
    </div>

    <!-- Filter & Search -->
    <div class="flex flex-wrap items-center gap-3 mb-6">
      <input
        v-model="search"
        type="text"
        placeholder="Search projects..."
        class="input-field flex-1 min-w-[200px] max-w-xs"
      />
      <div class="flex gap-2">
        <button
          v-for="s in ['ALL', 'ACTIVE', 'ON_HOLD', 'COMPLETED']"
          :key="s"
          @click="statusFilter = s"
          class="px-3 py-1.5 text-sm rounded-lg transition-colors font-medium"
          :class="statusFilter === s
            ? 'bg-indigo-600 text-white'
            : 'bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-gray-200'"
        >
          {{ s === 'ALL' ? 'All' : s.replace('_', ' ') }}
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลดโปรเจกต์...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded-xl p-6 text-red-400">
      <h2 class="font-bold mb-1">Failed to load projects</h2>
      <p class="text-sm">{{ error }}</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredProjects.length === 0" class="text-center py-20">
      <div class="text-6xl mb-4">📭</div>
      <h2 class="text-xl font-semibold text-gray-300 mb-2">No projects found</h2>
      <p class="text-gray-500 mb-6">{{ search ? 'Try a different search term.' : 'Create your first project to get started.' }}</p>
      <button v-if="!search" @click="openCreateModal" class="btn-primary px-6 py-2.5">
        Create Project
      </button>
    </div>

    <!-- Projects Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
      <div
        v-for="project in filteredProjects"
        :key="project.id"
        class="project-card group"
      >
        <!-- Card Header -->
        <div class="flex items-start justify-between mb-4">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold border"
                :class="statusClass(project.status)"
              >
                {{ project.status.replace('_', ' ') }}
              </span>
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="healthClass(project)"
              >
                {{ healthLabel(project) }}
              </span>
            </div>
            <h3 class="text-lg font-bold text-white truncate group-hover:text-indigo-300 transition-colors">
              {{ project.name }}
            </h3>
            <p class="text-xs text-gray-500 font-mono mt-0.5">{{ project.code }}</p>
          </div>
          <!-- Delete Button -->
          <button
            @click.stop="confirmDelete(project)"
            class="opacity-0 group-hover:opacity-100 p-1.5 text-gray-600 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-all ml-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>

        <!-- Description -->
        <p class="text-sm text-gray-400 line-clamp-2 mb-4 min-h-[2.5rem]">
          {{ project.description || 'No description provided.' }}
        </p>

        <!-- Task Stats -->
        <div class="grid grid-cols-3 gap-2 mb-4">
          <div class="text-center p-2 bg-gray-800/60 rounded-lg">
            <div class="text-sm font-bold text-gray-200">{{ getTaskCount(project, 'total') }}</div>
            <div class="text-[10px] text-gray-500 uppercase">Tasks</div>
          </div>
          <div class="text-center p-2 bg-gray-800/60 rounded-lg">
            <div class="text-sm font-bold text-green-400">{{ getTaskCount(project, 'done') }}</div>
            <div class="text-[10px] text-gray-500 uppercase">Done</div>
          </div>
          <div class="text-center p-2 bg-gray-800/60 rounded-lg">
            <div class="text-sm font-bold text-red-400">{{ getTaskCount(project, 'overdue') }}</div>
            <div class="text-[10px] text-gray-500 uppercase">Overdue</div>
          </div>
        </div>

        <!-- Progress Bar -->
        <div class="mb-4" v-if="getTaskCount(project, 'total') > 0">
          <div class="flex justify-between text-xs text-gray-500 mb-1">
            <span>Progress</span>
            <span>{{ progressPct(project) }}%</span>
          </div>
          <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="progressBarClass(project)"
              :style="{ width: progressPct(project) + '%' }"
            ></div>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between pt-3 border-t border-gray-700/50">
          <span class="text-xs text-gray-600">
            {{ formatDate(project.created_at) }}
          </span>
          <div class="flex gap-2">
            <NuxtLink
              :to="`/projects/${project.code || project.id}?tab=gantt`"
              class="btn-ghost-sm"
              @click.stop
            >
              Gantt
            </NuxtLink>
            <NuxtLink
              :to="`/projects/${project.code || project.id}`"
              class="btn-primary-sm"
            >
              Open Hub →
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Project Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeCreateModal"
    >
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-bold text-white">Create New Project</h2>
          <button @click="closeCreateModal" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Project Name <span class="text-red-400">*</span></label>
            <input
              v-model="createForm.name"
              type="text"
              placeholder="e.g. MIMS HDMap Main (English only)"
              class="input-field w-full"
              @keyup.enter="createProject"
              :class="{ 'border-red-500': createError }"
            />
            <p class="text-xs text-gray-500 mt-1">Use English letters, numbers, spaces, hyphens only.</p>
          </div>
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Description</label>
            <textarea
              v-model="createForm.description"
              rows="3"
              placeholder="What is this project about?"
              class="input-field w-full resize-none"
            ></textarea>
          </div>
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Status</label>
            <select v-model="createForm.status" class="input-field w-full">
              <option value="ACTIVE">Active</option>
              <option value="ON_HOLD">On Hold</option>
              <option value="COMPLETED">Completed</option>
            </select>
          </div>

          <div v-if="createError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ createError }}
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button
            @click="createProject"
            :disabled="isCreating || !createForm.name.trim()"
            class="flex-1 btn-primary py-2.5 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {{ isCreating ? 'Creating...' : 'Create Project' }}
          </button>
          <button @click="closeCreateModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
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
      <div class="bg-gray-800 border border-red-600 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-xl font-bold text-white mb-2">Delete Project</h2>
        <p class="text-gray-400 text-sm mb-6">
          Are you sure you want to delete <strong class="text-white">{{ projectToDelete?.name }}</strong>? This will permanently delete all tasks, submissions, and history. This cannot be undone.
        </p>
        <div v-if="deleteError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
          {{ deleteError }}
        </div>
        <div class="flex gap-3">
          <button
            @click="deleteProject"
            :disabled="isDeleting"
            class="flex-1 py-2.5 bg-red-600 hover:bg-red-700 disabled:bg-gray-600 text-white font-semibold rounded-xl transition-colors disabled:cursor-not-allowed"
          >
            {{ isDeleting ? 'Deleting...' : 'Yes, Delete' }}
          </button>
          <button @click="closeDeleteModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import type { Project } from '~/core/modules/projects/infrastructure/projects-api'

definePageMeta({ layout: 'default', middleware: 'auth' })

const { fetchWithAuth } = useAuth()

const projects = ref<Project[]>([])
const isLoading = ref(true)
const error = ref('')
const search = ref('')
const statusFilter = ref('ALL')

const showCreateModal = ref(false)
const createForm = ref({ name: '', description: '', status: 'ACTIVE' })
const isCreating = ref(false)
const createError = ref('')

const showDeleteModal = ref(false)
const projectToDelete = ref<Project | null>(null)
const isDeleting = ref(false)
const deleteError = ref('')

const totalProjects = computed(() => projects.value.length)
const totalActive = computed(() => projects.value.filter((p) => p.status === 'ACTIVE').length)
const totalCompleted = computed(() => projects.value.filter((p) => p.status === 'COMPLETED').length)
const onHold = computed(() => projects.value.filter((p) => p.status === 'ON_HOLD').length)

const filteredProjects = computed(() =>
  projects.value.filter((p) => {
    const matchSearch = !search.value || p.name.toLowerCase().includes(search.value.toLowerCase()) || p.code.toLowerCase().includes(search.value.toLowerCase())
    const matchStatus = statusFilter.value === 'ALL' || p.status === statusFilter.value
    return matchSearch && matchStatus
  })
)

function statusClass(status: string) {
  if (status === 'ACTIVE') return 'bg-green-500/10 text-green-400 border-green-500/30'
  if (status === 'COMPLETED') return 'bg-blue-500/10 text-blue-400 border-blue-500/30'
  if (status === 'ON_HOLD') return 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30'
  return 'bg-gray-700 text-gray-400 border-gray-600'
}

function getTaskCount(project: Project, type: 'total' | 'done' | 'overdue') {
  const tasks = project.tasks || []
  if (type === 'total') return tasks.length
  if (type === 'done') return tasks.filter((t) => t.status === 'COMPLETED').length
  if (type === 'overdue') {
    const now = Date.now()
    return tasks.filter((t) => t.status !== 'COMPLETED' && t.due_at && new Date(t.due_at).getTime() < now).length
  }
  return 0
}

function progressPct(project: Project) {
  const total = getTaskCount(project, 'total')
  if (!total) return 0
  return Math.round((getTaskCount(project, 'done') / total) * 100)
}

function progressBarClass(project: Project) {
  const pct = progressPct(project)
  if (pct >= 80) return 'bg-green-500'
  if (pct >= 40) return 'bg-indigo-500'
  return 'bg-yellow-500'
}

function healthLabel(project: Project) {
  const overdue = getTaskCount(project, 'overdue')
  const total = getTaskCount(project, 'total')
  if (!total) return '— Empty'
  const ratio = overdue / total
  if (ratio === 0) return '● Healthy'
  if (ratio < 0.2) return '● At Risk'
  return '● Critical'
}

function healthClass(project: Project) {
  const overdue = getTaskCount(project, 'overdue')
  const total = getTaskCount(project, 'total')
  if (!total) return 'bg-gray-700/50 text-gray-500'
  const ratio = overdue / total
  if (ratio === 0) return 'bg-green-500/10 text-green-400'
  if (ratio < 0.2) return 'bg-yellow-500/10 text-yellow-400'
  return 'bg-red-500/10 text-red-400'
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

async function loadProjects() {
  isLoading.value = true
  error.value = ''
  try {
    const data = await fetchWithAuth<{ data: Project[] }>('/sentinel/projects')
    projects.value = data.data || []
  } catch (e: any) {
    error.value = e.message || 'Failed to load projects'
  } finally {
    isLoading.value = false
  }
}

function openCreateModal() {
  createForm.value = { name: '', description: '', status: 'ACTIVE' }
  createError.value = ''
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
}

async function createProject() {
  if (!createForm.value.name.trim()) return
  isCreating.value = true
  createError.value = ''
  try {
    await fetchWithAuth('/sentinel/projects', {
      method: 'POST',
      body: createForm.value,
    })
    closeCreateModal()
    await loadProjects()
  } catch (e: any) {
    createError.value = e.message || 'Failed to create project'
  } finally {
    isCreating.value = false
  }
}

function confirmDelete(project: Project) {
  projectToDelete.value = project
  deleteError.value = ''
  showDeleteModal.value = true
}

function closeDeleteModal() {
  showDeleteModal.value = false
  projectToDelete.value = null
}

async function deleteProject() {
  if (!projectToDelete.value) return
  isDeleting.value = true
  deleteError.value = ''
  try {
    await fetchWithAuth(`/sentinel/projects/${projectToDelete.value.id}`, { method: 'DELETE' })
    closeDeleteModal()
    await loadProjects()
  } catch (e: any) {
    deleteError.value = e.message || 'Failed to delete project'
  } finally {
    isDeleting.value = false
  }
}

onMounted(loadProjects)
</script>

<style scoped>
.metric-card {
  @apply bg-gray-800/60 rounded-xl p-4 border border-gray-700/50;
}
.metric-label {
  @apply text-xs text-gray-500 mt-1 uppercase tracking-wide;
}
.project-card {
  @apply bg-gray-800 border border-gray-700 rounded-2xl p-5 hover:border-indigo-500/40 hover:shadow-lg hover:shadow-indigo-500/5 transition-all cursor-default;
}
.input-field {
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-colors;
}
.btn-primary {
  @apply bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl transition-colors;
}
.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg transition-colors;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
</style>
