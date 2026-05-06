<template>
  <div class="create-enterprise-page min-h-full w-full pb-8">
    <div class="w-full max-w-7xl mx-auto px-4 sm:px-8 lg:px-12 py-6 md:py-8">
      <!-- Header -->
      <div class="mb-6 md:mb-8">
        <h1 class="text-3xl md:text-4xl font-bold text-white tracking-tight">Create Task</h1>
        <p class="text-base text-gray-400 mt-2">Add a new task to your project</p>
      </div>

      <!-- Card — single scroll via layout <main>; avoids nested overflow trapping wheel in Chrome -->
      <div class="create-enterprise-card rounded-2xl flex flex-col">

        <div class="p-6 md:p-10 space-y-6 md:space-y-7">

          <!-- Success -->
          <div v-if="showSuccessMsg" class="flex items-center gap-3 p-4 md:p-5 bg-green-900/30 border border-green-600 rounded-xl text-green-400 text-base">
            <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            Task created! Redirecting...
          </div>

          <!-- Error -->
          <div v-if="errorMessage" class="p-4 md:p-5 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-base">{{ errorMessage }}</div>

          <!-- Task Type Selector -->
          <div>
            <label class="label">Type *</label>
            <div class="grid grid-cols-3 gap-3 sm:gap-4">
              <button
                type="button"
                @click="form.task_type = 'FEATURE'"
                :class="form.task_type === 'FEATURE' ? 'border-purple-300 dark:border-purple-500 bg-purple-100 dark:bg-purple-500/20 text-purple-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-purple-300 dark:border-purple-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">★</span> Feature
              </button>
              <button
                type="button"
                @click="form.task_type = 'TASK'"
                :class="form.task_type === 'TASK' ? 'border-blue-300 dark:border-blue-500 bg-blue-100 dark:bg-blue-500/20 text-blue-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-blue-300 dark:border-blue-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">📋</span> Task
              </button>
              <button
                type="button"
                @click="form.task_type = 'BUG'"
                :class="form.task_type === 'BUG' ? 'border-red-300 dark:border-red-500 bg-red-100 dark:bg-red-500/20 text-red-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-red-300 dark:border-red-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">⚠</span> Bug
              </button>
            </div>
            <!-- Product Owner rule hint -->
            <div v-if="form.task_type === 'FEATURE'" class="mt-3 flex items-start gap-3 p-4 bg-purple-900/20 border border-purple-500/30 rounded-xl text-sm sm:text-base text-purple-300 leading-relaxed">
              <span class="shrink-0 mt-0.5">★</span>
              <span><strong>Feature mode:</strong> Acts as a parent container. Estimated effort is disabled — add sub-tasks of type Task/Bug to assign work.</span>
            </div>
          </div>

          <!-- Title -->
          <div>
            <label class="label">Title *</label>
            <input v-model="form.title" type="text" class="input-field w-full" placeholder="e.g. Implement user authentication system" />
          </div>

          <!-- Description -->
          <div>
            <label class="label">Description</label>
            <RichTextEditor v-model="form.description" placeholder="Describe the task objectives and requirements… (paste images with ⌘V)" />
          </div>

          <!-- Assignee -->
          <div>
            <label class="label" :class="form.task_type === 'FEATURE' ? 'text-gray-500' : ''">
              Assignee
              <span v-if="form.task_type === 'FEATURE'" class="text-gray-600 font-normal">(disabled for Features)</span>
            </label>
            <select
              v-model="form.assigned_to"
              class="input-field w-full transition-opacity"
              :class="form.task_type === 'FEATURE' ? 'opacity-40 cursor-not-allowed' : ''"
              :disabled="form.task_type === 'FEATURE' || assigneeLoading"
            >
              <option value="">— Unassigned —</option>
              <option v-for="u in visibleAssigneeUsers" :key="u.id" :value="u.id">{{ assigneeLabel(u) }}</option>
            </select>
          </div>

          <!-- Project -->
          <div>
            <label class="label">Project *</label>
            <select
              v-model="form.project_id"
              @change="onProjectChange"
              class="input-field w-full"
              :class="projectFieldError ? 'border-amber-600 ring-1 ring-amber-600/40' : ''"
              required
            >
              <option value="" disabled>Select a project or Komgrip</option>
              <option value="__komgrip__">📋 Komgrip (ไม่อยู่ในโครงการ)</option>
              <option v-for="p in projects" :key="p.id" :value="p.id">{{ p.name }}</option>
            </select>
            <p v-if="form.project_id === '__komgrip__'" class="text-sm text-violet-400 mt-2">งานนี้จะถูกสร้างใน Komgrip — ไม่ผูกกับโครงการใด สถานะมีแค่ Pending / Done</p>
            <p v-if="projectFieldError" class="text-sm text-amber-400 mt-2">{{ projectFieldError }}</p>
          </div>

          <!-- Estimated Effort -->
          <div>
            <label class="label" :class="form.task_type === 'FEATURE' ? 'text-gray-500' : ''">
              Estimated Effort (hours)
              <span v-if="form.task_type === 'FEATURE'" class="text-gray-600 font-normal">(disabled for Features)</span>
            </label>
            <input
              v-model.number="form.estimated_hours"
              type="number"
              min="0"
              step="0.1"
              class="input-field w-full transition-opacity"
              :class="form.task_type === 'FEATURE' ? 'opacity-40 cursor-not-allowed' : ''"
              :disabled="form.task_type === 'FEATURE'"
              placeholder="e.g. 1.5"
            />
            <p v-if="form.task_type !== 'FEATURE'" class="text-sm text-gray-500 mt-2">Hours, up to 1 decimal place (e.g. 1.5). Used for Manday and Quotation (Costing Engine).</p>
          </div>

          <!-- Priority & Story Points -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Priority</label>
              <select v-model="form.priority" class="input-field w-full">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points</label>
              <input v-model.number="form.story_points" type="number" min="0" step="0.5" class="input-field w-full" placeholder="0" />
            </div>
          </div>

          <!-- Sprint (only when project is selected and not Komgrip) -->
          <div v-if="form.project_id && form.project_id !== '__komgrip__' && sprints.length">
            <label class="label">Sprint</label>
            <select v-model="form.sprint_id" class="input-field w-full">
              <option value="">Backlog</option>
              <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </div>

          <!-- Epic (only when project is selected and not Komgrip) -->
          <div v-if="form.project_id && form.project_id !== '__komgrip__' && epics.length">
            <label class="label">Epic</label>
            <select v-model="form.epic_id" class="input-field w-full">
              <option value="">No Epic</option>
              <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
            </select>
          </div>

          <!-- Due Date -->
          <div>
            <label class="label">Due Date</label>
            <UiDatePicker v-model="form.due_date" placeholder="Select due date…" />
          </div>

          <!-- Start / End Dates -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Start Date</label>
              <UiDatePicker v-model="form.start_date" placeholder="Select start date…" />
            </div>
            <div>
              <label class="label">End Date</label>
              <UiDatePicker v-model="form.end_date" placeholder="Select end date…" />
            </div>
          </div>

        </div>

        <!-- Footer actions -->
        <div class="flex flex-col-reverse sm:flex-row gap-3 sm:gap-4 p-6 md:p-8 pt-4 border-t border-gray-700">
          <button
            @click="handleSubmit"
            :disabled="isSubmitting || !form.title.trim() || !canSubmitProject"
            class="btn-primary-enterprise flex-1 disabled:opacity-40 text-white font-semibold rounded-xl py-4 md:py-4 text-base md:text-lg transition-all flex items-center justify-center gap-2 min-h-[3.25rem]"
          >
            <svg v-if="isSubmitting" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
            <span>{{ isSubmitting ? 'Creating...' : 'Create Task' }}</span>
          </button>
          <NuxtLink to="/dashboard" class="btn-ghost-enterprise sm:shrink-0 px-6 py-4 text-base font-medium text-center min-h-[3.25rem] flex items-center justify-center">Cancel</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { effortHoursToMinutes } from '~/utils/effortHours'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { isTaskAssigneeRole } from '~/utils/roles'
import { useAuth } from '~/composables/useAuth'
import RichTextEditor from '~/components/editor/RichTextEditor.vue'

const KOMGRIP_VALUE = '__komgrip__'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const router = useRouter()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()
const { showSuccess, showError } = useNotification()
const { currentUser, fetchWithAuth } = useAuth()
const teamsStore = useTeamsStore()

type AssigneeUser = { id: number; email: string; display_name: string; first_name?: string; last_name?: string; role: string }
const assigneeUsers = ref<AssigneeUser[]>([])
const assigneeLoading = ref(false)

const visibleAssigneeUsers = computed(() => {
  const isCeo = currentUser.value?.role?.toUpperCase() === 'CEO'
  return assigneeUsers.value.filter((u) => isCeo || u.role?.toUpperCase() !== 'CEO')
})

async function loadAssignees() {
  if (assigneeUsers.value.length > 0 || assigneeLoading.value) return
  assigneeLoading.value = true
  try {
    await teamsStore.fetchTeamsFeatureEnabled()
    if (teamsStore.teamsFeatureEnabled) {
      await teamsStore.fetchTeams()
      const allUsers = teamsStore.teams.flatMap((t: any) => t.users ?? [])
      const seen = new Set<number>()
      assigneeUsers.value = allUsers.filter((u: AssigneeUser) => {
        if (seen.has(u.id)) return false
        seen.add(u.id)
        return isTaskAssigneeRole(u.role)
      })
    } else {
      const res = await fetchWithAuth<{ data: AssigneeUser[] }>('/auth/users')
      assigneeUsers.value = (res.data ?? []).filter((u: AssigneeUser) => isTaskAssigneeRole(u.role))
    }
  } catch {
    // non-critical
  } finally {
    assigneeLoading.value = false
  }
}

function assigneeLabel(u: AssigneeUser) {
  const parts = [u.first_name, u.last_name].map((p) => (p || '').trim()).filter(Boolean)
  const name = parts.length > 0 ? parts.join(' ') : (u.display_name?.trim() || u.email)
  return `${name} (${u.role})`
}

interface ProjectItem { id: string; name: string }
interface SprintItem  { id: string; name: string }
interface EpicItem   { id: string; title: string }

const projects = ref<ProjectItem[]>([])
const sprints  = ref<SprintItem[]>([])
const epics    = ref<EpicItem[]>([])

const isSubmitting   = ref(false)
const showSuccessMsg = ref(false)
const errorMessage   = ref('')
const projectFieldError = ref('')

const canSubmitProject = computed(() => {
  return Boolean(form.value.project_id?.trim())
})

const form = ref({
  task_type:         'TASK',
  title:             '',
  description:       '',
  project_id:        '',
  sprint_id:         '',
  epic_id:           '',
  priority:          'MEDIUM',
  story_points:      0,
  estimated_hours: 0,
  due_date:          '',
  start_date:        '',
  end_date:          '',
  assigned_to:       '' as string | number,
})

onMounted(async () => {
  try {
    const list = await projectsApi.getProjects()
    projects.value = list.map((p: any) => ({ id: p.id, name: p.name }))
  } catch {
    // non-critical
  }
  void loadAssignees()
})

async function onProjectChange() {
  sprints.value = []
  epics.value   = []
  form.value.sprint_id = ''
  form.value.epic_id   = ''
  if (!form.value.project_id) return
  try {
    const [sprintList, epicList] = await Promise.all([
      projectsApi.getSprints(form.value.project_id),
      projectsApi.getEpics(form.value.project_id),
    ])
    sprints.value = sprintList.map((s: any) => ({ id: s.id, name: s.name }))
    epics.value   = epicList.map((e: any) => ({ id: e.id, title: e.title }))
  } catch {
    // ignore
  }
}

async function handleSubmit() {
  if (!form.value.title.trim()) return
  projectFieldError.value = ''
  if (!form.value.project_id?.trim()) {
    projectFieldError.value = 'กรุณาเลือกโครงการหรือ Komgrip'
    return
  }
  isSubmitting.value  = true
  errorMessage.value  = ''
  showSuccessMsg.value = false

  try {
    if (form.value.project_id === KOMGRIP_VALUE) {
      await tasksApi.createKomgripTask({
        title:             form.value.title,
        description:       form.value.description,
        priority:          form.value.priority,
        estimated_minutes: effortHoursToMinutes(Number(form.value.estimated_hours) || 0),
      })
      showSuccessMsg.value = true
      showSuccess('Komgrip task created!', 'Done')
      setTimeout(() => router.push('/komgrip'), 1200)
      return
    }

    const dateOnlyToISO = (ymd: string) => new Date(`${ymd}T00:00:00`).toISOString()
    const payload: any = {
      title:             form.value.title,
      description:       form.value.description,
      task_type:         form.value.task_type || 'TASK',
      priority:          form.value.priority,
      story_points:      form.value.story_points,
      estimated_minutes: form.value.task_type === 'FEATURE' ? 0 : effortHoursToMinutes(Number(form.value.estimated_hours) || 0),
    }
    payload.project_id = form.value.project_id
    if (form.value.sprint_id)   payload.sprint_id   = form.value.sprint_id
    if (form.value.epic_id)     payload.epic_id     = form.value.epic_id
    if (form.value.due_date)    payload.due_date    = dateOnlyToISO(form.value.due_date)
    if (form.value.start_date)  payload.start_date  = dateOnlyToISO(form.value.start_date)
    if (form.value.end_date)    payload.end_date    = dateOnlyToISO(form.value.end_date)
    if (form.value.assigned_to) payload.assigned_to = Number(form.value.assigned_to)

    await tasksApi.createTask(payload)
    showSuccessMsg.value = true
    showSuccess('Task created successfully!', 'Done')
    setTimeout(() => {
      router.push(`/projects/${form.value.project_id}?tab=backlog`)
    }, 1200)
  } catch (err: any) {
    errorMessage.value = err?.data?.message ?? err?.message ?? 'Failed to create task.'
    showError(errorMessage.value)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.create-enterprise-page {
  background:
    radial-gradient(1200px 620px at 84% -16%, rgba(139, 92, 246, 0.18), transparent 60%),
    radial-gradient(960px 520px at -8% 0%, rgba(59, 130, 246, 0.16), transparent 56%),
    linear-gradient(180deg, #070b17 0%, #0b1220 54%, #090f1a 100%);
}

.create-enterprise-card {
  @apply border border-white/10 bg-slate-900/75 shadow-[0_18px_42px_rgba(2,6,23,0.46)] backdrop-blur-sm;
}

.label {
  @apply block text-sm sm:text-base text-slate-200 mb-2 font-semibold tracking-wide;
}

.input-field {
  @apply bg-slate-800/90 border border-slate-600/75 rounded-xl px-4 py-3.5 text-base text-slate-100 placeholder-slate-500 focus:outline-none focus:border-violet-400 focus:ring-2 focus:ring-violet-500/35 transition-all;
}

.btn-primary-enterprise {
  @apply bg-gradient-to-r from-violet-600 via-fuchsia-600 to-indigo-600 hover:from-violet-500 hover:via-fuchsia-500 hover:to-indigo-500 shadow-[0_12px_25px_rgba(124,58,237,0.35)];
}

.btn-ghost-enterprise {
  @apply bg-slate-800/90 hover:bg-slate-700 text-slate-200 rounded-xl border border-white/10 transition-colors;
}
</style>
