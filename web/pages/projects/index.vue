<template>
  <div class="projects-index-enterprise min-h-screen text-white p-6">
    <!-- Header -->
    <div class="mb-8 flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white">
          Project Command Center
        </h1>
        <p class="text-sm text-gray-400 mt-1">{{ projects.length }} active workspaces</p>
      </div>
      <button
        v-if="canManageProjects"
        type="button"
        class="flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-xl shadow-lg shadow-purple-500/20 transition-all"
        @click="openCreateModal"
      >
        <span>+</span>
        <span>New Project</span>
      </button>
    </div>

    <!-- Squad banner (Product Owner only, when squads are enabled) -->
    <div
      v-if="!isCEO && (currentUser?.role === 'PRODUCT_OWNER' || currentUser?.role === 'PM') && teamsStore.teamsFeatureEnabled && squadName"
      class="mb-6 flex items-center gap-3 px-4 py-3 bg-purple-900/30 border border-purple-500/30 rounded-xl"
    >
      <div class="flex h-8 w-8 items-center justify-center rounded-full bg-purple-600/20 text-purple-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
      </div>
      <div>
        <span class="text-xs text-purple-400 uppercase tracking-widest font-semibold">Your Squad</span>
        <p class="text-white font-bold leading-tight">{{ squadName }}</p>
      </div>
      <span class="ml-auto text-xs text-gray-500">Showing projects for your team only</span>
    </div>
    <!-- No-squads mode: Product Owner sees only CEO-assigned projects -->
    <div
      v-else-if="!isCEO && (currentUser?.role === 'PRODUCT_OWNER' || currentUser?.role === 'PM') && !teamsStore.teamsFeatureEnabled"
      class="mb-6 flex items-center gap-3 px-4 py-3 bg-amber-900/20 border border-amber-600/40 rounded-xl text-amber-100"
    >
      <div class="flex h-8 w-8 items-center justify-center rounded-full bg-amber-600/20 text-amber-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
      </div>
      <div class="min-w-0">
        <span class="text-xs text-amber-400 uppercase tracking-widest font-semibold">Product Owner access</span>
        <p class="text-sm text-amber-100/95 leading-snug mt-0.5">Squads are off — you only see projects the CEO assigned to you as Product Owner. Ask the CEO to add you on the project card if something is missing.</p>
      </div>
    </div>

    <!-- System Metrics Row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8 projects-metrics-grid">
      <div class="metric-card metric-card-enterprise">
        <div class="text-2xl font-bold text-purple-400">{{ totalActive }}</div>
        <div class="metric-label">Active Projects</div>
      </div>
      <div class="metric-card metric-card-enterprise">
        <div class="text-2xl font-bold text-green-400">{{ totalCompleted }}</div>
        <div class="metric-label">Completed</div>
      </div>
      <div class="metric-card metric-card-enterprise">
        <div class="text-2xl font-bold text-yellow-400">{{ onHold }}</div>
        <div class="metric-label">On Hold</div>
      </div>
      <div class="metric-card metric-card-enterprise">
        <div class="text-2xl font-bold text-purple-400">{{ totalProjects }}</div>
        <div class="metric-label">Total Projects</div>
      </div>
    </div>

    <!-- Filter & Search -->
    <div class="projects-filter-enterprise flex flex-wrap items-center gap-3 mb-6">
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
            ? 'bg-gradient-to-r from-purple-600 to-pink-600 text-white'
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
      <button v-if="!search && canManageProjects" @click="openCreateModal" class="btn-primary px-6 py-2.5">
        Create Project
      </button>
    </div>

    <!-- Projects Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
      <div
        v-for="project in filteredProjects"
        :key="project.id"
        class="project-card project-card-enterprise group block cursor-pointer"
        @click="navigateTo(`/projects/${project.id}`)"
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
            <h3 class="text-lg font-bold text-white truncate group-hover:text-purple-300 transition-colors">
              {{ project.name }}
            </h3>
            <p class="text-xs text-gray-500 font-mono mt-0.5">{{ project.code }}</p>
          </div>
          <!-- Delete Button -->
          <button
            v-if="canDeleteProjects"
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

        <!-- Team Badge + CEO Assign Dropdown (hidden when teams feature disabled) -->
        <div v-if="teamsStore.teamsFeatureEnabled" class="mb-4">
          <!-- Current Team Badge (all roles) -->
          <div v-if="project.team_id" class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 mb-2">
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v1h8v-1zM6 8a2 2 0 11-4 0 2 2 0 014 0zM16 18v-1a5.972 5.972 0 00-.75-2.906A3.005 3.005 0 0119 15v1h-3zM4.75 12.094A5.973 5.973 0 004 15v1H1v-1a3 3 0 013.75-2.906z"/></svg>
            {{ teamsStore.teamNameById(project.team_id) }}
          </div>
          <div v-else class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-gray-700/50 text-gray-500 mb-2">
            No team assigned
          </div>
          <!-- CEO-only: Assign Team Dropdown -->
          <div v-if="isCEO" class="mt-1">
            <select
              :value="project.team_id ?? ''"
              @change="assignTeam(project, ($event.target as HTMLSelectElement).value ? Number(($event.target as HTMLSelectElement).value) : null)"
              class="w-full bg-gray-700/60 border border-gray-600/50 rounded-lg px-3 py-1.5 text-xs text-gray-300 focus:outline-none focus:border-purple-500 transition-colors cursor-pointer"
              @click.stop
            >
              <option value="">— Unassign team —</option>
              <option v-for="team in teamsStore.teams" :key="team.id" :value="team.id">
                {{ team.name }}
              </option>
            </select>
          </div>
        </div>

        <!-- CEO/MANAGER: assign Product Owner users when squads are disabled (multi-select) -->
        <div
          v-if="!teamsStore.teamsFeatureEnabled && isCEO"
          class="mb-4 space-y-2"
          @click.stop
        >
          <div class="text-xs text-gray-500 font-medium uppercase tracking-wide">Product Owner users</div>
          <p class="text-[11px] text-gray-500 leading-snug">
            คลิกชื่อเพื่อเลือกหรือยกเลิก — เลือกได้หลายคนพร้อมกัน (คลิกชื่อที่เลือกอยู่แล้วจะถอนการเลือก)
          </p>
          <div
            v-if="pmOwnersAssignError"
            class="flex items-start gap-2 rounded-lg border border-red-500/45 bg-red-950/40 px-3 py-2 text-xs text-red-100"
            role="alert"
          >
            <span class="flex-1 leading-snug">{{ pmOwnersAssignError }}</span>
            <button
              type="button"
              class="shrink-0 rounded p-0.5 text-red-400 hover:bg-red-900/50 hover:text-red-200"
              aria-label="Dismiss error"
              @click.stop="pmOwnersAssignError = ''"
            >
              ✕
            </button>
          </div>
          <div
            v-if="pmUserOptions.length === 0"
            class="text-xs text-amber-400/90"
          >
            No Product Owner users in the system.
          </div>
          <div
            v-else
            class="max-h-36 overflow-y-auto rounded-lg border border-gray-600/50 bg-gray-700/40 px-1 py-1.5 space-y-0.5"
          >
            <label
              v-for="u in pmUserOptions"
              :key="u.id"
              class="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-xs text-gray-200 hover:bg-gray-600/35 cursor-pointer select-none"
              @click.stop
            >
              <input
                type="checkbox"
                class="h-3.5 w-3.5 shrink-0 rounded border-gray-500 bg-gray-800 text-purple-600 focus:ring-1 focus:ring-purple-500/50"
                :checked="isPmOwner(project, u.id)"
                @change="togglePmOwner(project, u.id, $event)"
              />
              <span class="truncate">{{ u.display_name || u.email }}</span>
            </label>
          </div>
        </div>

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

        <!-- Mode Toggle -->
        <div class="flex gap-2 mb-5 p-1 bg-gray-900/60 rounded-xl">
          <button
            @click="importMode = false"
            class="flex-1 py-2 text-xs font-semibold rounded-lg transition-colors"
            :class="!importMode ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-gray-200'"
          >
            ✨ New Project
          </button>
          <button
            @click="importMode = true"
            class="flex-1 py-2 text-xs font-semibold rounded-lg transition-colors"
            :class="importMode ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-gray-200'"
          >
            📥 Import from Backup
          </button>
        </div>

        <!-- New Project Form -->
        <div v-if="!importMode" class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Project Name <span class="text-red-400">*</span></label>
            <input
              v-model="createForm.name"
              type="text"
              placeholder="e.g. MIMS HDMap Main หรือ โปรเจกต์หลัก"
              class="input-field w-full"
              @keyup.enter="createProject"
              :class="{ 'border-red-500': createError }"
            />
            <p class="text-xs text-gray-500 mt-1">Use letters, numbers, spaces, hyphens, underscores (รองรับภาษาไทย)</p>
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

        <!-- Import from Backup Form -->
        <div v-else class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Project Name <span class="text-red-400">*</span></label>
            <input
              v-model="importForm.name"
              type="text"
              placeholder="ชื่อโครงการใหม่ (รองรับภาษาไทย)"
              class="input-field w-full"
              :class="{ 'border-red-500': createError }"
            />
            <p class="text-xs text-gray-500 mt-1">โครงการจะได้รับ code ใหม่โดยอัตโนมัติ</p>
          </div>
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">ไฟล์ Backup (.sentinel.json) <span class="text-red-400">*</span></label>
            <label
              class="flex flex-col items-center justify-center w-full h-28 border-2 border-dashed rounded-xl cursor-pointer transition-colors"
              :class="importFile ? 'border-purple-500/60 bg-purple-500/5' : 'border-gray-600 bg-gray-700/40 hover:border-gray-500'"
            >
              <template v-if="!importFile">
                <svg class="w-8 h-8 text-gray-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <p class="text-xs text-gray-400">คลิกเพื่อเลือกไฟล์ หรือลาก & วาง</p>
                <p class="text-xs text-gray-600 mt-1">.sentinel.json</p>
              </template>
              <template v-else>
                <svg class="w-6 h-6 text-purple-400 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-xs text-purple-300 font-semibold">{{ importFile.name }}</p>
                <p class="text-xs text-gray-500 mt-0.5">{{ formatFileSize(importFile.size) }}</p>
              </template>
              <input type="file" accept=".json,.sentinel.json" class="hidden" @change="onImportFileChange" />
            </label>
          </div>

          <!-- Preview -->
          <div v-if="importPayloadPreview" class="bg-gray-700/40 border border-gray-600/50 rounded-xl p-3 space-y-1">
            <p class="text-xs font-semibold text-gray-300 mb-2">ข้อมูลที่จะ import:</p>
            <div class="grid grid-cols-2 gap-1 text-xs">
              <span class="text-gray-500">โครงการต้นทาง:</span>
              <span class="text-gray-200 font-medium">{{ importPayloadPreview.project?.name || '-' }}</span>
              <span class="text-gray-500">Tasks:</span>
              <span class="text-gray-200">{{ importPayloadPreview.tasks?.length || 0 }} รายการ</span>
              <span class="text-gray-500">Sprints:</span>
              <span class="text-gray-200">{{ importPayloadPreview.sprints?.length || 0 }} รายการ</span>
              <span class="text-gray-500">Milestones:</span>
              <span class="text-gray-200">{{ importPayloadPreview.milestones?.length || 0 }} รายการ</span>
              <span class="text-gray-500">Epics:</span>
              <span class="text-gray-200">{{ importPayloadPreview.epics?.length || 0 }} รายการ</span>
            </div>
          </div>

          <div v-if="createError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ createError }}
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button
            v-if="!importMode"
            @click="createProject"
            :disabled="isCreating || !createForm.name.trim()"
            class="flex-1 btn-primary py-2.5 disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <div v-if="isCreating" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            {{ isCreating ? 'Creating...' : 'Create Project' }}
          </button>
          <button
            v-else
            @click="importProjectFromBackup"
            :disabled="isCreating || !importForm.name.trim() || !importPayloadPreview"
            class="flex-1 btn-primary py-2.5 disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <div v-if="isCreating" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            {{ isCreating ? 'Importing...' : '📥 Import & Create' }}
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
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'

definePageMeta({ layout: 'default', middleware: 'auth' })

const { fetchWithAuth, currentUser } = useAuth()
const teamsStore = useTeamsStore()
const teamsApi = useTeamsApi()
const projectsApi = useProjectsApi()

const isCEO = computed(() => currentUser.value?.role === 'CEO' || currentUser.value?.role === 'MANAGER')
const canManageProjects = computed(() =>
  ['CEO', 'MANAGER', 'PRODUCT_OWNER', 'PM'].includes(currentUser.value?.role ?? ''))
const canDeleteProjects = computed(() => ['CEO'].includes(currentUser.value?.role ?? ''))
const squadName = computed(() => {
  const tid = currentUser.value?.team_id
  if (!tid) return null
  return teamsStore.teamNameById(tid)
})

const projects = ref<Project[]>([])
const isLoading = ref(true)
const error = ref('')
const search = ref('')
const statusFilter = ref('ALL')

const showCreateModal = ref(false)
const createForm = ref({ name: '', description: '', status: 'ACTIVE' })
const isCreating = ref(false)
const createError = ref('')

// Import-from-backup mode
const importMode = ref(false)
const importForm = ref({ name: '' })
const importFile = ref<File | null>(null)
const importPayloadPreview = ref<Record<string, any> | null>(null)

const showDeleteModal = ref(false)
const projectToDelete = ref<Project | null>(null)
const isDeleting = ref(false)
const deleteError = ref('')

const pmUserOptions = ref<{ id: number; email: string; display_name?: string }[]>([])
const pmOwnersAssignError = ref('')

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
  // Prefer counts from list API (populated by backend)
  if (project.task_total !== undefined) {
    if (type === 'total') return project.task_total
    if (type === 'done') return project.task_completed ?? 0
    if (type === 'overdue') return project.task_overdue ?? 0
  }
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
  if (pct >= 40) return 'bg-purple-500'
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

async function loadProjects(opts?: { quiet?: boolean }) {
  const quiet = !!opts?.quiet
  if (!quiet) {
    isLoading.value = true
    error.value = ''
  }
  try {
    const data = await fetchWithAuth<{ data: Project[] }>('/sentinel/projects')
    projects.value = data.data || []
  } catch (e: any) {
    if (!quiet) {
      error.value = e.message || 'Failed to load projects'
    }
  } finally {
    if (!quiet) {
      isLoading.value = false
    }
  }
}

function openCreateModal() {
  createForm.value = { name: '', description: '', status: 'ACTIVE' }
  createError.value = ''
  importMode.value = false
  importForm.value = { name: '' }
  importFile.value = null
  importPayloadPreview.value = null
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  importFile.value = null
  importPayloadPreview.value = null
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

function onImportFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  importFile.value = file
  importPayloadPreview.value = null
  createError.value = ''
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const json = JSON.parse(e.target?.result as string)
      // Support both raw payload and wrapped { payload: ... } formats
      const payload = json.payload ?? json
      if (!payload.project) {
        createError.value = 'ไฟล์ไม่ถูกต้อง: ไม่พบข้อมูล project ในไฟล์ backup'
        importFile.value = null
        return
      }
      importPayloadPreview.value = payload
      // Pre-fill name from backup if empty
      if (!importForm.value.name && payload.project?.name) {
        importForm.value.name = payload.project.name + ' (Imported)'
      }
    } catch {
      createError.value = 'ไม่สามารถอ่านไฟล์ได้ กรุณาตรวจสอบว่าเป็นไฟล์ .sentinel.json ที่ถูกต้อง'
      importFile.value = null
    }
  }
  reader.readAsText(file)
}

async function importProjectFromBackup() {
  if (!importForm.value.name.trim() || !importPayloadPreview.value) return
  isCreating.value = true
  createError.value = ''
  try {
    await projectsApi.importProjectFromBackup(importForm.value.name.trim(), importPayloadPreview.value)
    closeCreateModal()
    await loadProjects()
  } catch (e: any) {
    createError.value = e.message || 'Failed to import project from backup'
  } finally {
    isCreating.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
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

async function assignTeam(project: Project, teamId: number | null) {
  try {
    await teamsStore.assignProjectToTeam(project.id, teamId)
    project.team_id = teamId
    project.team_name = teamId ? teamsStore.teamNameById(teamId) : undefined
  } catch (_e) {
    // silent — user will see no change
  }
}

async function loadPmUserOptions() {
  if (!isCEO.value || teamsStore.teamsFeatureEnabled) {
    pmUserOptions.value = []
    return
  }
  try {
    const res = await fetchWithAuth<{ data: { id: number; email: string; display_name?: string; role: string }[] }>('/auth/users')
    pmUserOptions.value = (res.data || []).filter((u) => u.role === 'PRODUCT_OWNER' || u.role === 'PM')
  } catch {
    pmUserOptions.value = []
  }
}

function isPmOwner(project: Project, userId: number): boolean {
  return !!project.pm_owners?.some((o) => o.user_id === userId)
}

function pmOwnersAssignErrorMessage(e: unknown): string {
  const err = e as { data?: { message?: string; error?: string }; statusMessage?: string; message?: string }
  const fromData = err?.data && typeof err.data === 'object'
    ? (err.data.message || err.data.error)
    : undefined
  return (
    (typeof fromData === 'string' && fromData) ||
    err?.statusMessage ||
    (typeof err?.message === 'string' && err.message) ||
    'ไม่สามารถบันทึกการมอบหมาย Product Owner ได้'
  )
}

/** Toggle one Product Owner; click again on the same name clears that selection. Multiple Product Owners can stay selected. */
async function togglePmOwner(project: Project, userId: number, e: Event) {
  const checked = (e.target as HTMLInputElement).checked
  const set = new Set<number>((project.pm_owners ?? []).map((o) => o.user_id))
  if (checked) {
    set.add(userId)
  } else {
    set.delete(userId)
  }
  const ids = [...set].sort((a, b) => a - b)
  const idx = projects.value.findIndex((p) => p.id === project.id)
  if (idx === -1) return

  pmOwnersAssignError.value = ''

  const nextOwners = pmUserOptions.value
    .filter((u) => ids.includes(u.id))
    .map((u) => ({
      user_id: u.id,
      email: u.email,
      display_name: u.display_name,
    }))
  projects.value[idx] = { ...projects.value[idx], pm_owners: nextOwners }

  try {
    const updated = await teamsApi.assignProjectPmOwners(project.id, ids)
    projects.value[idx] = { ...projects.value[idx], pm_owners: updated.pm_owners }
  } catch (err) {
    pmOwnersAssignError.value = pmOwnersAssignErrorMessage(err)
    await loadProjects({ quiet: true })
  }
}

onMounted(async () => {
  await teamsStore.fetchTeamsFeatureEnabled()
  await loadPmUserOptions()
  await loadProjects()
  if (teamsStore.teamsFeatureEnabled) await teamsStore.fetchTeams()
})
</script>

<style scoped>
.projects-index-enterprise {
  background:
    radial-gradient(1200px 620px at 84% -16%, rgba(139, 92, 246, 0.18), transparent 60%),
    radial-gradient(960px 520px at -8% 0%, rgba(59, 130, 246, 0.16), transparent 56%),
    linear-gradient(180deg, #070b17 0%, #0b1220 54%, #090f1a 100%);
}

.projects-metrics-grid {
  @apply gap-4;
}

.metric-card {
  @apply rounded-2xl p-4;
}

.metric-card-enterprise {
  @apply bg-gradient-to-b from-slate-800/70 to-slate-900/80 border border-white/10 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)] backdrop-blur-sm;
}

.metric-label {
  @apply text-[11px] text-slate-400 mt-1 uppercase tracking-[0.08em] font-medium;
}

.projects-filter-enterprise {
  @apply rounded-2xl border border-white/10 bg-slate-900/70 px-4 py-3 shadow-[0_14px_34px_rgba(2,6,23,0.34)] backdrop-blur-sm;
}

.project-card {
  @apply rounded-2xl p-5 transition-all cursor-pointer;
}

.project-card-enterprise {
  @apply bg-slate-900/75 border border-white/10 shadow-[0_16px_38px_rgba(2,6,23,0.42)] hover:border-violet-400/45 hover:shadow-[0_18px_40px_rgba(124,58,237,0.2)] backdrop-blur-sm;
}

.input-field {
  @apply bg-slate-800/90 border border-slate-600/70 rounded-xl px-4 py-2.5 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-violet-400 focus:ring-2 focus:ring-violet-500/35 transition-all;
}

.btn-primary {
  @apply bg-gradient-to-r from-violet-600 via-fuchsia-600 to-indigo-600 hover:from-violet-500 hover:via-fuchsia-500 hover:to-indigo-500 text-white font-semibold rounded-xl transition-all shadow-[0_12px_25px_rgba(124,58,237,0.35)];
}

.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-gradient-to-r from-violet-600 via-fuchsia-600 to-indigo-600 hover:from-violet-500 hover:via-fuchsia-500 hover:to-indigo-500 text-white font-medium rounded-lg transition-all shadow-[0_8px_20px_rgba(124,58,237,0.35)];
}

.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-slate-800/90 hover:bg-slate-700 text-slate-200 font-medium rounded-lg border border-white/10 transition-colors;
}
</style>
