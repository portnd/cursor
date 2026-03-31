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
        <NuxtLink v-if="projectId" :to="`/projects/${projectId}?tab=sprints`" class="mt-4 inline-block text-sm text-gray-400 hover:text-white">← Back to Project</NuxtLink>
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
                :class="sprint.status === 'ACTIVE' ? 'bg-purple-500/20 text-purple-400' : sprint.status === 'COMPLETED' ? 'bg-gray-600 text-gray-400' : 'bg-yellow-500/20 text-yellow-400'"
              >
                {{ sprint.status }}
              </span>
              <span v-if="sprint.start_date || sprint.end_date" class="text-xs text-gray-500">
                {{ formatDate(sprint.start_date) }} – {{ formatDate(sprint.end_date) }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <NuxtLink
              v-if="prevSprint"
              :to="`/projects/sprint/${prevSprint.id}?project=${projectId}`"
              class="btn-ghost-sm inline-flex items-center gap-1"
            >
              ← ก่อนหน้า
            </NuxtLink>
            <NuxtLink
              v-if="nextSprint"
              :to="`/projects/sprint/${nextSprint.id}?project=${projectId}`"
              class="btn-ghost-sm inline-flex items-center gap-1"
            >
              ถัดไป →
            </NuxtLink>
            <button
              type="button"
              @click="openBacklogImportModal"
              class="inline-flex items-center gap-1.5 rounded-lg border border-emerald-500/40 bg-emerald-500/10 px-3 py-1.5 text-xs font-medium text-emerald-200 hover:bg-emerald-500/20 transition-colors"
              title="Add existing tasks from the backlog (or another sprint)"
            >
              <svg class="w-3.5 h-3.5 shrink-0 opacity-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
              </svg>
              From backlog
            </button>
            <button @click="openCreateTaskModal()" class="btn-primary-sm">+ New task</button>
            <NuxtLink :to="`/projects/${projectId}?tab=sprints`" class="btn-ghost-sm inline-flex">← Back</NuxtLink>
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
          <div class="text-2xl font-bold text-purple-400">{{ totalSp }}</div>
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
          <div class="flex flex-wrap items-center justify-end gap-2 sm:gap-3">
            <button
              v-if="backlogTaskCount > 0 || otherSprintTaskCount > 0"
              type="button"
              @click="openBacklogImportModal"
              class="inline-flex items-center gap-1.5 rounded-lg border border-emerald-500/35 bg-emerald-500/10 px-2.5 py-1 text-[11px] font-medium text-emerald-300 hover:bg-emerald-500/20 transition-colors"
            >
              + From backlog
              <span v-if="backlogTaskCount" class="text-emerald-400/90">({{ backlogTaskCount }})</span>
            </button>
            <template v-if="sprintTasks.length">
              <label class="flex items-center gap-2 cursor-pointer select-none text-sm text-gray-400 hover:text-gray-200">
                <input
                  ref="checkAllInputRef"
                  type="checkbox"
                  :checked="isCheckAllChecked"
                  class="w-4 h-4 rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
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
            </template>
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
                class="w-4 h-4 rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                @change="toggleTaskSelection(t.id)"
              />
            </label>
            <div
              class="flex flex-1 items-center justify-between min-w-0 cursor-pointer"
              @click="navigateToTask(t.id)"
            >
              <div class="flex items-center gap-3 min-w-0">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(t.code) }}</span>
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
        <div v-else class="text-center py-12 px-4 text-gray-500 text-sm">
          <p class="mb-1 text-gray-400 font-medium">No tasks in this sprint yet</p>
          <p class="mb-6 text-xs text-gray-600 max-w-md mx-auto leading-relaxed">
            Pull work that already exists from the <span class="text-gray-500">backlog</span> (no sprint), or create a new task.
          </p>
          <div class="flex flex-wrap items-center justify-center gap-2">
            <button
              v-if="backlogTaskCount > 0 || otherSprintTaskCount > 0"
              type="button"
              @click="openBacklogImportModal"
              class="inline-flex items-center gap-2 rounded-xl border border-emerald-500/40 bg-emerald-500/15 px-4 py-2.5 text-sm font-medium text-emerald-200 hover:bg-emerald-500/25 transition-colors"
            >
              Add from backlog
              <span v-if="backlogTaskCount" class="text-xs font-normal text-emerald-400/90">({{ backlogTaskCount }} ready)</span>
            </button>
            <button type="button" @click="openCreateTaskModal()" class="btn-primary-sm px-4 py-2.5">+ New task</button>
          </div>
          <NuxtLink
            :to="`/projects/${projectId}?tab=backlog`"
            class="mt-5 inline-block text-xs text-purple-400/90 hover:text-purple-300"
          >
            Open full Backlog tab → assign or create tasks there
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Add existing tasks (backlog / other sprints) → this sprint -->
    <div
      v-if="showBacklogImportModal && sprint"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeBacklogImportModal"
    >
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-5 sm:p-6 max-w-xl w-full shadow-2xl max-h-[88vh] flex flex-col">
        <div class="flex items-start justify-between gap-3 mb-4">
          <div>
            <h2 class="text-lg font-bold text-white">Add tasks to this sprint</h2>
            <p class="text-sm text-gray-400 mt-1">
              เลือกงานที่มีอยู่แล้วในโปรเจกต์ แล้วดึงเข้า
              <span class="text-gray-200 font-medium">“{{ sprint.name }}”</span>
            </p>
          </div>
          <button type="button" class="text-gray-500 hover:text-white shrink-0" aria-label="Close" @click="closeBacklogImportModal">✕</button>
        </div>

        <div class="rounded-xl border border-emerald-500/20 bg-emerald-500/5 px-3 py-2.5 mb-4 text-xs text-gray-400 leading-relaxed">
          <span class="text-emerald-300/95 font-medium">Backlog</span>
          = งานที่ยังไม่ได้อยู่ใน sprint ใด
          <span class="text-gray-600 mx-1">·</span>
          <span class="text-amber-300/90 font-medium">Other sprint</span>
          = ย้ายจาก sprint อื่นมาที่นี่ได้ (โหมดขยายด้านล่าง)
        </div>

        <div class="mb-3">
          <label class="sr-only" for="backlog-import-search">Search tasks</label>
          <input
            id="backlog-import-search"
            v-model="backlogImportQuery"
            type="search"
            class="input-field w-full"
            placeholder="ค้นหาชื่อหรือรหัส task…"
            autocomplete="off"
          />
        </div>

        <div class="flex rounded-xl border border-gray-600 p-1 bg-gray-900/40 mb-3">
          <button
            type="button"
            class="flex-1 rounded-lg px-3 py-2 text-xs font-medium transition-colors"
            :class="backlogImportScope === 'backlog' ? 'bg-emerald-600/30 text-emerald-200' : 'text-gray-500 hover:text-gray-300'"
            @click="setBacklogImportScope('backlog')"
          >
            Backlog only
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg px-3 py-2 text-xs font-medium transition-colors"
            :class="backlogImportScope === 'anywhere' ? 'bg-amber-600/25 text-amber-200' : 'text-gray-500 hover:text-gray-300'"
            @click="setBacklogImportScope('anywhere')"
          >
            + Other sprints
          </button>
        </div>

        <p class="text-[11px] text-gray-500 mb-2">
          แสดง {{ backlogImportFiltered.length }} รายการ
          <span v-if="backlogImportScope === 'backlog'">(ใน backlog {{ backlogTaskCount }} งาน)</span>
          <span v-else>(ยังไม่อยู่ใน sprint นี้ {{ tasksNotInThisSprintCount }} งาน)</span>
          · เลือกแล้ว {{ selectedBacklogTaskIds.length }}
        </p>

        <div class="flex flex-wrap gap-2 mb-3">
          <button type="button" class="text-xs text-purple-400 hover:text-purple-300" @click="selectAllVisibleBacklogImport">
            เลือกทั้งหมดที่แสดง
          </button>
          <button type="button" class="text-xs text-gray-500 hover:text-gray-400" @click="selectedBacklogTaskIds = []">
            ล้างการเลือก
          </button>
        </div>

        <div class="flex-1 min-h-[12rem] overflow-y-auto rounded-xl border border-gray-700 bg-gray-900/30 p-2">
          <div v-if="backlogImportFiltered.length === 0" class="text-center py-10 px-3 text-sm text-gray-500">
            <template v-if="backlogImportScope === 'backlog'">
              ไม่มีงานใน backlog
              <span v-if="otherSprintTaskCount > 0" class="block mt-2 text-xs text-gray-600">
                มีงานใน sprint อื่น {{ otherSprintTaskCount }} รายการ — กดแท็บ “+ Other sprints”
              </span>
            </template>
            <template v-else>ไม่มีงานที่จะย้ายเข้า sprint นี้ (หรือลองค้นหาอย่างอื่น)</template>
          </div>
          <label
            v-for="t in backlogImportFiltered"
            :key="t.id"
            class="flex items-center gap-3 py-2.5 px-2 rounded-lg hover:bg-gray-700/45 cursor-pointer"
          >
            <input
              v-model="selectedBacklogTaskIds"
              type="checkbox"
              :value="t.id"
              class="rounded border-gray-600 bg-gray-700 text-emerald-500 focus:ring-emerald-500 shrink-0"
            />
            <span class="text-xs font-mono text-gray-500 shrink-0 w-14 truncate" :title="t.code ?? ''">{{ taskCodeSuffix(t.code) }}</span>
            <span class="text-sm text-gray-200 truncate flex-1 min-w-0">{{ t.title }}</span>
            <span
              v-if="!t.sprint_id"
              class="text-[10px] px-1.5 py-0.5 rounded shrink-0 bg-gray-600/80 text-gray-300"
            >Backlog</span>
            <span
              v-else
              class="text-[10px] px-1.5 py-0.5 rounded shrink-0 max-w-[7rem] truncate bg-amber-500/15 text-amber-400"
              :title="sprintNameById(t.sprint_id) || ''"
            >{{ sprintNameById(t.sprint_id) || 'Sprint' }}</span>
          </label>
        </div>

        <div v-if="backlogImportError" class="mt-3 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
          {{ backlogImportError }}
        </div>

        <div class="flex flex-col sm:flex-row gap-3 mt-4">
          <button
            type="button"
            class="flex-1 btn-primary py-2.5 disabled:opacity-45"
            :disabled="isBacklogImporting || selectedBacklogTaskIds.length === 0"
            @click="confirmBacklogImport"
          >
            {{ isBacklogImporting ? 'กำลังเพิ่ม…' : `เพิ่ม ${selectedBacklogTaskIds.length} งานเข้า sprint` }}
          </button>
          <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="closeBacklogImportModal">
            ยกเลิก
          </button>
        </div>
        <NuxtLink
          :to="`/projects/${projectId}?tab=backlog`"
          class="mt-3 block text-center text-xs text-gray-500 hover:text-purple-400"
          @click="closeBacklogImportModal"
        >
          เปิดแท็บ Backlog เต็มหน้าจอ (สร้าง / แก้ / import slides)
        </NuxtLink>
      </div>
    </div>

    <!-- Create Task Modal (same structure as project page) -->
    <div v-if="showCreateTaskModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-start justify-center z-50 p-3 sm:p-6 overflow-y-auto" @click.self="closeCreateTaskModal">
      <div class="create-task-modal bg-gray-800 border border-gray-700 rounded-2xl w-full max-w-7xl shadow-2xl my-4 sm:my-8 flex flex-col max-h-[calc(100dvh-2rem)] min-h-0">
        <div class="flex items-center justify-between px-6 sm:px-8 pt-6 sm:pt-8 pb-4 shrink-0 border-b border-gray-700/80">
          <h2 class="text-2xl sm:text-3xl font-bold text-white tracking-tight">Add Task</h2>
          <button type="button" @click="closeCreateTaskModal" class="shrink-0 w-11 h-11 flex items-center justify-center rounded-xl text-gray-400 hover:text-white hover:bg-gray-700 text-xl leading-none" aria-label="Close">✕</button>
        </div>
        <div class="px-6 sm:px-8 py-6 sm:py-8 space-y-6 sm:space-y-7 flex-1 overflow-y-auto overscroll-contain min-h-0">
          <div>
            <label class="label">Title *</label>
            <input v-model="createTaskForm.title" type="text" class="input-field w-full" placeholder="Task title..." />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="createTaskForm.description" rows="6" class="input-field w-full resize-y min-h-[10rem]" placeholder="Describe the task..."></textarea>
          </div>
          <div>
            <label class="label">Estimated Effort (Minutes/Hours) *</label>
            <input
              v-model.number="createTaskForm.estimated_minutes"
              type="number"
              min="0"
              step="1"
              class="input-field w-full"
              placeholder="e.g. 60 (minutes)"
              required
            />
            <p class="text-sm text-gray-500 mt-2">Minutes. Used for Manday and Quotation (Costing Engine).</p>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
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
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
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
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Start Date</label>
              <input v-model="createTaskForm.start_date" type="datetime-local" class="input-field w-full" />
            </div>
            <div>
              <label class="label">End Date</label>
              <input v-model="createTaskForm.end_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <div v-if="createTaskError" class="p-4 md:p-5 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-base">{{ createTaskError }}</div>
        </div>
        <div class="flex flex-col-reverse sm:flex-row gap-3 sm:gap-4 px-6 sm:px-8 py-5 sm:py-6 border-t border-gray-700 shrink-0">
          <button @click="submitCreateTask" :disabled="isCreatingTask || !createTaskForm.title.trim() || (Number(createTaskForm.estimated_minutes) ?? 0) < 0" class="flex-1 btn-primary py-4 text-base sm:text-lg font-semibold rounded-xl disabled:opacity-40 min-h-[3.25rem]">
            {{ isCreatingTask ? 'Creating...' : 'Create Task' }}
          </button>
          <button type="button" @click="closeCreateTaskModal" class="sm:shrink-0 px-6 py-4 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-xl transition-colors text-base font-medium min-h-[3.25rem]">Cancel</button>
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
const { showSuccess } = useNotification()

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
const sprints = ref<Sprint[]>([])
const allTasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')

const sprintsOrdered = computed(() =>
  [...sprints.value].sort(
    (a, b) =>
      (a.sort_order ?? 0) - (b.sort_order ?? 0) ||
      new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  )
)
const currentSprintIndex = computed(() => {
  if (!sprint.value) return -1
  return sprintsOrdered.value.findIndex((s) => s.id === sprint.value!.id)
})
const prevSprint = computed(() => {
  const i = currentSprintIndex.value
  return i > 0 ? sprintsOrdered.value[i - 1] ?? null : null
})
const nextSprint = computed(() => {
  const i = currentSprintIndex.value
  return i >= 0 && i < sprintsOrdered.value.length - 1 ? sprintsOrdered.value[i + 1] ?? null : null
})

const sprintTasks = computed(() => allTasks.value.filter((t) => t.sprint_id === sprintId.value))
const doneCount = computed(() => sprintTasks.value.filter((t) => t.status === 'COMPLETED').length)
const totalSp = computed(() => sprintTasks.value.reduce((s, t) => s + (t.story_points || 0), 0))
const overdueCount = computed(() => {
  const now = Date.now()
  return sprintTasks.value.filter((t) => t.status !== 'COMPLETED' && t.due_at && new Date(t.due_at).getTime() < now).length
})

const backlogTaskCount = computed(() => allTasks.value.filter((t) => !t.sprint_id).length)
const otherSprintTaskCount = computed(() =>
  allTasks.value.filter((t) => !!t.sprint_id && t.sprint_id !== sprintId.value).length
)
const tasksNotInThisSprintCount = computed(() => allTasks.value.filter((t) => t.sprint_id !== sprintId.value).length)

const showBacklogImportModal = ref(false)
const backlogImportScope = ref<'backlog' | 'anywhere'>('backlog')
const backlogImportQuery = ref('')
const selectedBacklogTaskIds = ref<string[]>([])
const backlogImportError = ref('')
const isBacklogImporting = ref(false)

function setBacklogImportScope(s: 'backlog' | 'anywhere') {
  backlogImportScope.value = s
}

const backlogImportCandidates = computed(() => {
  const sid = sprintId.value
  return allTasks.value
    .filter((t) => {
      if (t.sprint_id === sid) return false
      if (backlogImportScope.value === 'backlog') return !t.sprint_id
      return true
    })
    .sort((a, b) => (a.code ?? '').localeCompare(b.code ?? '', undefined, { numeric: true }))
})

const backlogImportFiltered = computed(() => {
  const q = backlogImportQuery.value.trim().toLowerCase()
  if (!q) return backlogImportCandidates.value
  return backlogImportCandidates.value.filter((t) => {
    const code = (t.code ?? '').toLowerCase()
    const title = (t.title ?? '').toLowerCase()
    return code.includes(q) || title.includes(q)
  })
})

function sprintNameById(id: string | null | undefined): string {
  if (!id) return ''
  return sprints.value.find((s) => s.id === id)?.name ?? ''
}

function toYMD(d: string) {
  return d.split('T')[0]
}

function taskDatesInSprintRange(
  _task: { start_date?: string | null; end_date?: string | null; due_at?: string | null },
  sp: { start_date: string | null; end_date: string | null }
): { start_date: string; end_date: string } | null {
  if (!sp?.start_date) return null
  const addDays = (ymd: string, days: number) => {
    const dt = new Date(ymd + 'T12:00:00Z')
    dt.setUTCDate(dt.getUTCDate() + days)
    return toYMD(dt.toISOString())
  }
  const spStart = toYMD(sp.start_date)
  let spEnd = sp.end_date ? toYMD(sp.end_date) : addDays(spStart, 14)
  if (spEnd <= spStart) spEnd = addDays(spStart, 1)
  const toNoonUTC = (ymd: string) => ymd + 'T12:00:00.000Z'
  return { start_date: toNoonUTC(spStart), end_date: toNoonUTC(spEnd) }
}

function openBacklogImportModal() {
  backlogImportQuery.value = ''
  backlogImportScope.value = 'backlog'
  selectedBacklogTaskIds.value = []
  backlogImportError.value = ''
  showBacklogImportModal.value = true
}

function closeBacklogImportModal() {
  showBacklogImportModal.value = false
  backlogImportError.value = ''
}

function selectAllVisibleBacklogImport() {
  selectedBacklogTaskIds.value = backlogImportFiltered.value.map((t) => t.id)
}

async function confirmBacklogImport() {
  if (!sprint.value || selectedBacklogTaskIds.value.length === 0) return
  isBacklogImporting.value = true
  backlogImportError.value = ''
  const sp = sprint.value
  const ids = [...selectedBacklogTaskIds.value]
  try {
    await projectsApi.addTasksToSprint(sp.id, ids)
    for (const id of ids) {
      const t = allTasks.value.find((x) => x.id === id)
      if (t) {
        t.sprint_id = sp.id
        const dates = taskDatesInSprintRange(t, sp)
        if (dates) {
          try {
            await tasksApi.updateTask(id, { start_date: dates.start_date, end_date: dates.end_date })
            t.start_date = dates.start_date
            t.end_date = dates.end_date
          } catch {
            // timeline dates optional
          }
        }
      }
    }
    showSuccess(`Added ${ids.length} task(s) to this sprint.`, 'Done')
    closeBacklogImportModal()
  } catch (e: any) {
    const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'เพิ่มงานไม่สำเร็จ'
    backlogImportError.value = typeof err === 'string' ? err : 'เพิ่มงานไม่สำเร็จ'
  } finally {
    isBacklogImporting.value = false
  }
}

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

function taskCodeSuffix(code: string | undefined): string {
  if (!code) return '–'
  const suffix = code.split('-').pop()
  return /^\d+$/.test(suffix || '') ? String(Number(suffix)).padStart(4, '0') : code
}

function taskStatusBadge(status: string) {
  if (status === 'COMPLETED') return 'bg-green-500/20 text-green-400'
  if (status === 'IN_PROGRESS') return 'bg-blue-500/20 text-blue-400'
  if (status === 'READY_FOR_TEST') return 'bg-cyan-500/20 text-cyan-400'
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
    const [sprintsList, tasks] = await Promise.all([
      projectsApi.getSprints(p.id),
      tasksApi.getTasksByProject(p.id),
    ])
    sprints.value = sprintsList
    const s = sprintsList.find((x) => x.id === sprintId.value)
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
  estimated_minutes: 0,
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
    estimated_minutes: 0,
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
      estimated_minutes: Number(createTaskForm.value.estimated_minutes) || 0,
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
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
.create-task-modal .label {
  @apply block text-sm sm:text-base text-gray-300 mb-2 font-medium;
}
.create-task-modal .input-field {
  @apply bg-gray-700 border border-gray-500 rounded-xl px-4 py-3.5 text-base text-gray-100 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-medium rounded-lg transition-colors;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
.btn-primary {
  @apply bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-xl transition-colors;
}
</style>
