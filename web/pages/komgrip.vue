<template>
  <div class="komgrip-page min-h-full w-full pb-8">
    <div class="w-full max-w-5xl mx-auto px-4 sm:px-8 py-6 md:py-8">

      <!-- Header -->
      <div class="mb-6 md:mb-8 flex items-start justify-between gap-4 flex-wrap">
        <div>
          <div class="flex items-center gap-3 mb-1">
            <div class="flex items-center justify-center w-10 h-10 rounded-xl bg-gradient-to-br from-violet-600 to-purple-700 shadow-lg shadow-violet-500/30">
              <svg class="w-5 h-5 !text-white" fill="none" stroke="white" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 8h6m-6 4h4" />
              </svg>
            </div>
            <h1 class="text-3xl md:text-4xl font-bold text-white tracking-tight">Komgrip</h1>
          </div>
          <p class="text-sm text-gray-400 ml-13">งานที่ไม่ได้อยู่ในโครงการ · เข้าถึงได้ทุกคน · ค้นหาได้ใน Work Log</p>
        </div>
        <button
          @click="openCreateModal"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-violet-600 to-purple-600 hover:from-violet-500 hover:to-purple-500 text-white font-semibold text-sm shadow-lg shadow-violet-500/25 transition-all shrink-0"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          สร้างงาน
        </button>
      </div>

      <!-- Stats row -->
      <div class="grid grid-cols-3 gap-3 mb-6">
        <div class="stat-card">
          <span class="text-2xl font-bold text-white">{{ tasks.length }}</span>
          <span class="text-xs text-gray-400 mt-0.5">ทั้งหมด</span>
        </div>
        <div class="stat-card">
          <span class="text-2xl font-bold text-amber-400">{{ pendingTasks.length }}</span>
          <span class="text-xs text-gray-400 mt-0.5">Pending</span>
        </div>
        <div class="stat-card">
          <span class="text-2xl font-bold text-green-400">{{ doneTasks.length }}</span>
          <span class="text-xs text-gray-400 mt-0.5">Done</span>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex items-center gap-1 mb-5 bg-white dark:bg-slate-800/60 border border-gray-200 dark:border-white/10 rounded-xl p-1 w-fit">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="activeTab = tab.value"
          :class="activeTab === tab.value
            ? 'bg-gradient-to-r from-violet-600 to-purple-600 text-white shadow-md'
            : 'text-gray-400 hover:text-white'"
          class="px-4 py-1.5 rounded-lg text-sm font-semibold transition-all"
        >
          {{ tab.label }}
          <span class="ml-1.5 text-xs opacity-70">{{ tab.count }}</span>
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-16 text-gray-500">
        <svg class="animate-spin w-7 h-7 mr-3" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
        </svg>
        กำลังโหลด...
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredTasks.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
        <div class="w-16 h-16 rounded-2xl bg-violet-900/30 border border-violet-700/30 flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <p class="text-gray-400 font-medium">ยังไม่มีงาน{{ activeTab !== 'all' ? 'ใน ' + (activeTab === 'pending' ? 'Pending' : 'Done') : '' }}</p>
        <p class="text-gray-600 text-sm mt-1">กดปุ่ม "สร้างงาน" เพื่อเพิ่มงานใหม่</p>
      </div>

      <!-- Task list -->
      <div v-else class="space-y-2">
        <TransitionGroup name="task-list">
          <div
            v-for="task in filteredTasks"
            :key="task.id"
            class="task-card group"
          >
            <!-- Status toggle -->
            <button
              @click="toggleStatus(task)"
              :disabled="updatingId === task.id"
              class="flex-shrink-0 mt-0.5 transition-all"
              :title="task.status === 'COMPLETED' ? 'Mark as Pending' : 'Mark as Done'"
            >
              <div v-if="updatingId === task.id" class="w-5 h-5 rounded-full border-2 border-violet-500 flex items-center justify-center">
                <svg class="animate-spin w-3 h-3 text-violet-400" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
              </div>
              <div
                v-else
                :class="task.status === 'COMPLETED'
                  ? 'bg-green-500/20 border-green-500 text-green-400 hover:bg-green-500/10'
                  : 'border-gray-600 text-transparent hover:border-violet-500 hover:bg-violet-500/10'"
                class="w-5 h-5 rounded-full border-2 flex items-center justify-center transition-all"
              >
                <svg v-if="task.status === 'COMPLETED'" class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </button>

            <!-- Task info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start gap-2 flex-wrap">
                <span class="text-[10px] font-mono text-violet-400 shrink-0 mt-0.5 bg-violet-500/10 px-1.5 py-0.5 rounded">{{ task.code }}</span>
                <span
                  :class="task.status === 'COMPLETED' ? 'line-through text-gray-500' : 'text-gray-100'"
                  class="text-sm font-medium flex-1 min-w-0"
                >{{ task.title }}</span>
              </div>
              <div v-if="task.description" class="text-xs text-gray-500 mt-1 truncate">{{ task.description }}</div>
              <div class="flex items-center gap-3 mt-1.5 flex-wrap">
                <span :class="priorityBadge(task.priority)" class="text-[10px] font-bold px-1.5 py-0.5 rounded uppercase">{{ task.priority }}</span>
                <span v-if="task.estimated_minutes" class="text-[10px] text-gray-500 flex items-center gap-0.5">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                  {{ formatMinutes(task.estimated_minutes) }}
                </span>
                <span class="text-[10px] text-gray-600">{{ formatDate(task.created_at) }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
              <NuxtLink
                :to="`/task/${task.id}`"
                class="p-1.5 rounded-lg text-gray-500 hover:text-violet-400 hover:bg-violet-500/10 transition-colors"
                title="Open task detail"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </NuxtLink>
              <button
                @click="confirmDelete(task)"
                class="p-1.5 rounded-lg text-gray-500 hover:text-red-400 hover:bg-red-500/10 transition-colors"
                title="Delete task"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </TransitionGroup>
      </div>
    </div>

    <!-- Create Task Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="showModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
          @mousedown.self="closeModal"
        >
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" />
          <div class="relative w-full max-w-lg bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl shadow-black/50">
            <!-- Modal header -->
            <div class="flex items-center justify-between px-6 py-4 border-b border-gray-700/60">
              <div class="flex items-center gap-2.5">
                <span class="text-violet-400 text-lg">📋</span>
                <h2 class="text-base font-bold text-gray-100">สร้างงาน Komgrip</h2>
              </div>
              <button type="button" @click="closeModal" class="p-1.5 rounded-lg text-gray-500 hover:text-gray-300 hover:bg-gray-700 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>

            <!-- Modal body -->
            <div class="px-6 py-5 space-y-4">
              <!-- Title -->
              <div>
                <label class="modal-label">ชื่องาน *</label>
                <input
                  v-model="newTask.title"
                  type="text"
                  class="modal-input w-full"
                  placeholder="ระบุชื่องาน..."
                  @keydown.enter="submitCreate"
                  autofocus
                />
              </div>

              <!-- Description -->
              <div>
                <label class="modal-label">รายละเอียด</label>
                <textarea
                  v-model="newTask.description"
                  rows="3"
                  class="modal-input w-full resize-none"
                  placeholder="รายละเอียดเพิ่มเติม (ไม่บังคับ)"
                />
              </div>

              <!-- Priority & Estimated Hours -->
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="modal-label">Priority</label>
                  <select v-model="newTask.priority" class="modal-input w-full">
                    <option value="CRITICAL">🔴 Critical</option>
                    <option value="HIGH">🟠 High</option>
                    <option value="MEDIUM">🟡 Medium</option>
                    <option value="LOW">🟢 Low</option>
                  </select>
                </div>
                <div>
                  <label class="modal-label">ชั่วโมงโดยประมาณ</label>
                  <input
                    v-model.number="newTask.estimated_hours"
                    type="number"
                    min="0"
                    step="0.5"
                    class="modal-input w-full"
                    placeholder="0"
                  />
                </div>
              </div>

              <!-- Error -->
              <p v-if="createError" class="text-sm text-red-400 bg-red-900/20 border border-red-700/30 rounded-lg px-3 py-2">{{ createError }}</p>
            </div>

            <!-- Modal footer -->
            <div class="flex gap-3 px-6 pb-5">
              <button
                @click="submitCreate"
                :disabled="creating || !newTask.title.trim()"
                class="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-violet-600 to-purple-600 hover:from-violet-500 hover:to-purple-500 text-white font-semibold text-sm transition-all disabled:opacity-40 flex items-center justify-center gap-2"
              >
                <svg v-if="creating" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
                {{ creating ? 'กำลังสร้าง...' : 'สร้างงาน' }}
              </button>
              <button @click="closeModal" class="px-5 py-2.5 rounded-xl bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm font-medium transition-colors border border-gray-700">
                ยกเลิก
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { Task } from '~/core/modules/projects/infrastructure/projects-api'
import { effortHoursToMinutes } from '~/utils/effortHours'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const { showSuccess, showError, confirm } = useNotification()
const tasksApi = useTasksApi()

const tasks = ref<Task[]>([])
const loading = ref(false)
const updatingId = ref<string | null>(null)

const activeTab = ref<'all' | 'pending' | 'done'>('all')

const tabs = computed(() => [
  { label: 'ทั้งหมด', value: 'all' as const, count: tasks.value.length },
  { label: 'Pending', value: 'pending' as const, count: pendingTasks.value.length },
  { label: 'Done', value: 'done' as const, count: doneTasks.value.length },
])

const pendingTasks = computed(() => tasks.value.filter(t => t.status !== 'COMPLETED'))
const doneTasks = computed(() => tasks.value.filter(t => t.status === 'COMPLETED'))
const filteredTasks = computed(() => {
  if (activeTab.value === 'pending') return pendingTasks.value
  if (activeTab.value === 'done') return doneTasks.value
  return tasks.value
})

async function loadTasks() {
  loading.value = true
  try {
    tasks.value = await tasksApi.getKomgripTasks()
  } catch {
    showError('โหลดข้อมูลไม่สำเร็จ')
  } finally {
    loading.value = false
  }
}

async function toggleStatus(task: Task) {
  updatingId.value = task.id
  const newStatus = task.status === 'COMPLETED' ? 'PENDING' : 'COMPLETED'
  try {
    const updated = await tasksApi.updateKomgripTaskStatus(task.id, newStatus)
    const idx = tasks.value.findIndex(t => t.id === task.id)
    if (idx !== -1) tasks.value[idx] = updated
  } catch {
    showError('อัปเดต status ไม่สำเร็จ')
  } finally {
    updatingId.value = null
  }
}

async function confirmDelete(task: Task) {
  const ok = await confirm({
    title: 'ลบงาน',
    message: `ต้องการลบ "${task.title}" ใช่หรือไม่?`,
    confirmLabel: 'ลบ',
    cancelLabel: 'ยกเลิก',
    variant: 'danger',
  })
  if (!ok) return
  try {
    await tasksApi.deleteKomgripTask(task.id)
    tasks.value = tasks.value.filter(t => t.id !== task.id)
    showSuccess('ลบงานสำเร็จ')
  } catch {
    showError('ลบงานไม่สำเร็จ')
  }
}

// Create modal
const showModal = ref(false)
const creating = ref(false)
const createError = ref('')
const newTask = ref({
  title: '',
  description: '',
  priority: 'MEDIUM',
  estimated_hours: 0,
})

function openCreateModal() {
  newTask.value = { title: '', description: '', priority: 'MEDIUM', estimated_hours: 0 }
  createError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function submitCreate() {
  if (!newTask.value.title.trim()) return
  creating.value = true
  createError.value = ''
  try {
    const created = await tasksApi.createKomgripTask({
      title: newTask.value.title.trim(),
      description: newTask.value.description,
      priority: newTask.value.priority,
      estimated_minutes: effortHoursToMinutes(newTask.value.estimated_hours || 0),
    })
    tasks.value.unshift(created)
    showSuccess('สร้างงานสำเร็จ')
    closeModal()
  } catch (err: any) {
    createError.value = err?.data?.message ?? err?.message ?? 'สร้างงานไม่สำเร็จ'
  } finally {
    creating.value = false
  }
}

function priorityBadge(priority: string) {
  const map: Record<string, string> = {
    CRITICAL: 'bg-red-900/40 text-red-400 border border-red-700/30',
    HIGH: 'bg-orange-900/40 text-orange-400 border border-orange-700/30',
    MEDIUM: 'bg-yellow-900/40 text-yellow-400 border border-yellow-700/30',
    LOW: 'bg-green-900/40 text-green-400 border border-green-700/30',
  }
  return map[priority] || map.MEDIUM
}

function formatMinutes(mins: number): string {
  if (!mins) return ''
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function formatDate(iso: string): string {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: '2-digit' })
  } catch { return '' }
}

onMounted(() => loadTasks())
</script>

<style scoped>
.komgrip-page {
  background:
    radial-gradient(1200px 620px at 84% -16%, rgba(139, 92, 246, 0.18), transparent 60%),
    radial-gradient(960px 520px at -8% 0%, rgba(59, 130, 246, 0.16), transparent 56%),
    linear-gradient(180deg, #070b17 0%, #0b1220 54%, #090f1a 100%);
}

.stat-card {
  @apply flex flex-col items-center justify-center py-4 rounded-xl border border-white/10 bg-slate-900/60 backdrop-blur-sm;
}

.task-card {
  @apply flex items-start gap-3 px-4 py-3.5 rounded-xl border border-white/10 bg-slate-900/60 hover:bg-slate-800/70 hover:border-violet-500/20 transition-all;
}

.modal-label {
  @apply block text-xs font-semibold text-slate-300 mb-1.5 tracking-wide uppercase;
}

.modal-input {
  @apply bg-slate-800/90 border border-slate-600/75 rounded-xl px-3.5 py-2.5 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-violet-400 focus:ring-2 focus:ring-violet-500/30 transition-all;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.task-list-enter-active,
.task-list-leave-active {
  transition: all 0.2s ease;
}
.task-list-enter-from,
.task-list-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
