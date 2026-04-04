<template>
  <div class="min-h-full w-full bg-gray-900 pb-8">
    <div class="w-full max-w-7xl mx-auto px-4 sm:px-8 lg:px-12 py-6 md:py-8">
      <!-- Header -->
      <div class="mb-6 md:mb-8">
        <h1 class="text-3xl md:text-4xl font-bold text-white tracking-tight">Create Task</h1>
        <p class="text-base text-gray-400 mt-2">Add a new task to your project</p>
      </div>

      <!-- Card — single scroll via layout <main>; avoids nested overflow trapping wheel in Chrome -->
      <div class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl flex flex-col">

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
                :class="form.task_type === 'FEATURE' ? 'border-purple-500 bg-purple-500/20 text-purple-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-purple-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">★</span> Feature
              </button>
              <button
                type="button"
                @click="form.task_type = 'TASK'"
                :class="form.task_type === 'TASK' ? 'border-blue-500 bg-blue-500/20 text-blue-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-blue-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">📋</span> Task
              </button>
              <button
                type="button"
                @click="form.task_type = 'BUG'"
                :class="form.task_type === 'BUG' ? 'border-red-500 bg-red-500/20 text-red-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-red-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">⚠</span> Bug
              </button>
            </div>
            <!-- PM Rule hint -->
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
            <textarea v-model="form.description" rows="6" class="input-field w-full resize-y min-h-[10rem]" placeholder="Describe the task objectives and requirements..."></textarea>
          </div>

          <!-- Project -->
          <div>
            <label class="label">Project</label>
            <select v-model="form.project_id" @change="onProjectChange" class="input-field w-full">
              <option value="">— No project —</option>
              <option v-for="p in projects" :key="p.id" :value="p.id">{{ p.name }}</option>
            </select>
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
              <input v-model.number="form.story_points" type="number" min="0" class="input-field w-full" placeholder="0" />
            </div>
          </div>

          <!-- Sprint (only when project is selected) -->
          <div v-if="form.project_id && sprints.length">
            <label class="label">Sprint</label>
            <select v-model="form.sprint_id" class="input-field w-full">
              <option value="">Backlog</option>
              <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </div>

          <!-- Epic (only when project is selected) -->
          <div v-if="form.project_id && epics.length">
            <label class="label">Epic</label>
            <select v-model="form.epic_id" class="input-field w-full">
              <option value="">No Epic</option>
              <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
            </select>
          </div>

          <!-- Due Date -->
          <div>
            <label class="label">Due Date</label>
            <input v-model="form.due_date" type="datetime-local" class="input-field w-full" />
          </div>

          <!-- Start / End Dates -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Start Date</label>
              <input v-model="form.start_date" type="datetime-local" class="input-field w-full" />
            </div>
            <div>
              <label class="label">End Date</label>
              <input v-model="form.end_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>

        </div>

        <!-- Footer actions -->
        <div class="flex flex-col-reverse sm:flex-row gap-3 sm:gap-4 p-6 md:p-8 pt-4 border-t border-gray-700">
          <button
            @click="handleSubmit"
            :disabled="isSubmitting || !form.title.trim()"
            class="flex-1 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 disabled:opacity-40 text-white font-semibold rounded-xl py-4 md:py-4 text-base md:text-lg transition-colors flex items-center justify-center gap-2 min-h-[3.25rem]"
          >
            <svg v-if="isSubmitting" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
            <span>{{ isSubmitting ? 'Creating...' : 'Create Task' }}</span>
          </button>
          <NuxtLink to="/dashboard" class="sm:shrink-0 px-6 py-4 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-xl transition-colors text-base font-medium text-center min-h-[3.25rem] flex items-center justify-center">Cancel</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { effortHoursToMinutes } from '~/utils/effortHours'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const router = useRouter()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()
const { showSuccess, showError } = useNotification()

interface ProjectItem { id: string; name: string }
interface SprintItem  { id: string; name: string }
interface EpicItem   { id: string; title: string }

const projects = ref<ProjectItem[]>([])
const sprints  = ref<SprintItem[]>([])
const epics    = ref<EpicItem[]>([])

const isSubmitting   = ref(false)
const showSuccessMsg = ref(false)
const errorMessage   = ref('')

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
})

onMounted(async () => {
  try {
    const list = await projectsApi.getProjects()
    projects.value = list.map((p: any) => ({ id: p.id, name: p.name }))
  } catch {
    // non-critical
  }
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
  isSubmitting.value  = true
  errorMessage.value  = ''
  showSuccessMsg.value = false

  try {
    const payload: any = {
      title:             form.value.title,
      description:       form.value.description,
      task_type:         form.value.task_type || 'TASK',
      priority:          form.value.priority,
      story_points:      form.value.story_points,
      estimated_minutes: form.value.task_type === 'FEATURE' ? 0 : effortHoursToMinutes(Number(form.value.estimated_hours) || 0),
    }
    if (form.value.project_id)  payload.project_id  = form.value.project_id
    if (form.value.sprint_id)   payload.sprint_id   = form.value.sprint_id
    if (form.value.epic_id)     payload.epic_id     = form.value.epic_id
    if (form.value.due_date)    payload.due_date    = new Date(form.value.due_date).toISOString()
    if (form.value.start_date)  payload.start_date  = new Date(form.value.start_date).toISOString()
    if (form.value.end_date)    payload.end_date    = new Date(form.value.end_date).toISOString()

    await tasksApi.createTask(payload)
    showSuccessMsg.value = true
    showSuccess('Task created successfully!', 'Done')
    setTimeout(() => {
      if (form.value.project_id) {
        router.push(`/projects/${form.value.project_id}?tab=backlog`)
      } else {
        router.push('/dashboard')
      }
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
.label {
  @apply block text-sm sm:text-base text-gray-300 mb-2 font-medium;
}
.input-field {
  @apply bg-gray-700 border border-gray-500 rounded-xl px-4 py-3.5 text-base text-gray-100 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
</style>
