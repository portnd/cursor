<template>
  <!-- Loading skeleton -->
  <section v-if="isLoading" class="rounded-2xl border border-gray-700 bg-gray-800/60 p-6">
    <div class="flex items-center gap-3 mb-4">
      <div class="h-4 w-40 animate-pulse rounded bg-gray-700"/>
      <div class="h-5 w-8 animate-pulse rounded-full bg-gray-700 ml-auto"/>
    </div>
    <div class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-24 animate-pulse rounded-xl bg-gray-700/60"/>
    </div>
  </section>

  <!-- Empty state -->
  <section v-else-if="queue.length === 0">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-cyan-400 mb-0.5">Continuous UAT Queue</h2>
        <p class="text-xs text-gray-500">Tasks awaiting your test approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-gray-700/60 border border-gray-600 text-gray-500">
        0 awaiting
      </span>
    </div>
    <div class="rounded-2xl border border-gray-700/50 bg-gray-800/30 px-6 py-10 text-center">
      <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-cyan-500/10 border border-cyan-500/20">
        <svg class="h-5 w-5 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
      </div>
      <p class="text-sm font-semibold text-gray-400">All clear</p>
      <p class="text-xs text-gray-600 mt-1">No tasks pending test review</p>
    </div>
  </section>

  <!-- Queue list -->
  <section v-else>
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-cyan-400 mb-0.5">Continuous UAT Queue</h2>
        <p class="text-xs text-gray-500">Tasks awaiting your test approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-cyan-500/10 border border-cyan-500/25 text-cyan-400">
        {{ queue.length }} awaiting test
      </span>
    </div>

    <div class="space-y-3">
      <div
        v-for="task in queue"
        :key="task.id"
        class="rounded-2xl border border-cyan-500/30 bg-gray-800/70 p-4 shadow-md hover:border-cyan-400/50 transition-colors"
      >
        <div class="flex items-start justify-between gap-3">
          <!-- Task info (click through to task detail) -->
          <NuxtLink
            :to="`/task/${task.id}`"
            class="flex-1 min-w-0 rounded-xl -m-1 p-1 outline-none transition-colors cursor-pointer hover:bg-cyan-500/5 focus-visible:ring-2 focus-visible:ring-cyan-500/40 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900"
          >
            <!-- Project pill + task type -->
            <div class="flex items-center gap-2 mb-1.5">
              <span
                class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-semibold truncate max-w-[140px]"
                :style="{
                  borderColor: task.project_color || '#6366f1',
                  color: task.project_color || '#6366f1',
                  backgroundColor: (task.project_color || '#6366f1') + '18'
                }"
              >
                <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :style="{ backgroundColor: task.project_color || '#6366f1' }"/>
                {{ task.project_name || 'Unknown Project' }}
              </span>
              <span
                class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border"
                :class="task.task_type === 'BUG'
                  ? 'border-red-500/40 text-red-400 bg-red-500/10'
                  : 'border-indigo-500/40 text-indigo-400 bg-indigo-500/10'"
              >{{ task.task_type || 'TASK' }}</span>
            </div>

            <!-- Title -->
            <p class="text-sm font-semibold text-white line-clamp-2 leading-snug mb-1">{{ task.title }}</p>

            <!-- Code + assignee -->
            <div class="flex items-center gap-3 text-[10px] text-gray-500">
              <span v-if="task.code" class="font-mono">{{ task.code }}</span>
              <span v-if="task.assigned_to_display_name || task.assigned_to_email" class="flex items-center gap-1">
                <svg class="w-3 h-3 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                {{ task.assigned_to_display_name || task.assigned_to_email }}
              </span>
            </div>
          </NuxtLink>

          <!-- Action buttons -->
          <div class="flex flex-col gap-2 flex-shrink-0">
            <button
              :disabled="actioningId === task.id || approveSubmitting"
              class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-100 dark:bg-emerald-500/10 border border-emerald-300 dark:border-emerald-500/30 text-emerald-400 text-xs font-semibold hover:bg-emerald-100 dark:bg-emerald-500/20 hover:border-emerald-400/50 transition-colors disabled:opacity-50"
              @click="openApproveModal(task)"
            >
              <span>✅</span>
              APPROVE
            </button>
            <button
              :disabled="actioningId === task.id"
              class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-100 dark:bg-red-500/10 border border-red-300 dark:border-red-500/30 text-red-400 text-xs font-semibold hover:bg-red-100 dark:bg-red-500/20 hover:border-red-400/50 transition-colors disabled:opacity-50"
              @click="openRejectModal(task)"
            >
              <span>❌</span>
              REJECT
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Product Owner approve: Test Evidence Form Modal -->
  <Teleport to="body">
    <div
      v-if="approveModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
      @click.self="!approveSubmitting && closeApproveModal()"
    >
      <div class="w-full max-w-lg rounded-2xl border border-emerald-500/30 bg-gray-900 shadow-2xl p-6">
        <!-- Header -->
        <div class="flex items-start gap-3 mb-5">
          <div class="w-8 h-8 rounded-lg bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center flex-shrink-0 mt-0.5">
            <span class="text-lg" aria-hidden="true">🧪</span>
          </div>
          <div class="min-w-0">
            <h3 class="text-sm font-bold text-white">Submit Test Evidence to CEO</h3>
            <p class="text-xs text-gray-500 truncate max-w-[320px]">{{ approveTarget?.title }}</p>
            <p class="text-xs text-amber-400/80 mt-1">Task will be forwarded to CEO for final approval — not marked as Done yet.</p>
          </div>
        </div>

        <!-- Test URL -->
        <div class="mb-4">
          <label class="block text-[11px] font-semibold text-gray-400 uppercase tracking-wider mb-1.5">
            Test / Staging URL <span class="text-red-400">*</span>
          </label>
          <input
            ref="approveUrlRef"
            v-model="approveTestUrl"
            type="url"
            placeholder="https://staging.example.com/feature-xyz"
            :disabled="approveSubmitting"
            class="w-full rounded-xl border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 focus:border-emerald-500/50 focus:outline-none focus:ring-1 focus:ring-emerald-500/30 disabled:opacity-50"
          />
          <p v-if="approveTestUrl.length > 0 && !approveTestUrl.startsWith('http')" class="text-[11px] text-red-400 mt-1">
            URL must start with http:// or https://
          </p>
        </div>

        <!-- Test Steps -->
        <div class="mb-5">
          <label class="block text-[11px] font-semibold text-gray-400 uppercase tracking-wider mb-1.5">
            Test Steps for CEO <span class="text-red-400">*</span>
          </label>
          <textarea
            v-model="approveTestSteps"
            rows="6"
            placeholder="Describe step-by-step how the CEO should test this feature:&#10;1. Navigate to...&#10;2. Click on...&#10;3. Verify that..."
            :disabled="approveSubmitting"
            class="w-full rounded-xl border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 focus:border-emerald-500/50 focus:outline-none focus:ring-1 focus:ring-emerald-500/30 resize-none disabled:opacity-50"
          />
          <div class="flex items-center justify-between mt-1">
            <p v-if="approveTestSteps.length > 0 && approveTestSteps.length < 20" class="text-[11px] text-red-400">
              At least {{ 20 - approveTestSteps.length }} more character(s) required
            </p>
            <span class="text-[11px] text-gray-600 ml-auto">{{ approveTestSteps.length }} chars</span>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3">
          <button
            type="button"
            :disabled="approveSubmitting"
            class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-400 hover:border-gray-600 hover:text-gray-200 transition-colors disabled:opacity-50"
            @click="closeApproveModal"
          >Cancel</button>
          <button
            type="button"
            :disabled="approveSubmitting || !isApproveFormValid"
            class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-emerald-100 dark:bg-emerald-500/15 border border-emerald-300 dark:border-emerald-500/40 text-emerald-400 text-xs font-bold hover:bg-emerald-100 dark:bg-emerald-500/25 transition-colors disabled:opacity-50"
            @click="submitApprove"
          >
            <svg v-if="approveSubmitting" class="w-3.5 h-3.5 animate-spin shrink-0" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <span v-if="approveSubmitting">Submitting…</span>
            <span v-else>✅ Submit to CEO</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Reject Modal (Teleport) -->
  <Teleport to="body">
    <div
      v-if="rejectModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
      @click.self="closeRejectModal"
    >
      <div class="w-full max-w-md rounded-2xl border border-red-500/30 bg-gray-900 shadow-2xl p-6">
        <!-- Header -->
        <div class="flex items-center gap-3 mb-4">
          <div class="w-8 h-8 rounded-lg bg-red-500/15 border border-red-500/30 flex items-center justify-center flex-shrink-0">
            <svg class="w-4 h-4 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
          <div>
            <h3 class="text-sm font-bold text-white">Reject Task</h3>
            <p class="text-xs text-gray-500 truncate max-w-[280px]">{{ rejectTarget?.title }}</p>
          </div>
        </div>

        <p class="text-xs text-gray-400 mb-3">Explain what failed so the developer can fix it. This will be logged as a comment.</p>

        <textarea
          ref="rejectTextareaRef"
          v-model="rejectReason"
          rows="4"
          placeholder="Describe the issue (min. 10 characters)…"
          class="w-full rounded-xl border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 focus:border-red-500/50 focus:outline-none focus:ring-1 focus:ring-red-500/30 resize-none"
        />

        <p v-if="rejectReason.length > 0 && rejectReason.length < 10" class="text-[11px] text-red-400 mt-1">
          At least {{ 10 - rejectReason.length }} more character(s) required
        </p>

        <div class="flex items-center justify-end gap-3 mt-4">
          <button
            class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-400 hover:border-gray-600 hover:text-gray-200 transition-colors"
            @click="closeRejectModal"
          >Cancel</button>
          <button
            :disabled="rejectReason.length < 10 || rejectSubmitting"
            class="px-4 py-2 rounded-lg bg-red-100 dark:bg-red-500/15 border border-red-300 dark:border-red-500/40 text-red-400 text-xs font-bold hover:bg-red-100 dark:bg-red-500/25 transition-colors disabled:opacity-50"
            @click="submitReject"
          >
            <span v-if="rejectSubmitting">Rejecting…</span>
            <span v-else>Confirm Reject</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { nextTick, computed } from 'vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'

const emit = defineEmits<{ (e: 'refresh'): void }>()

const { getTasksReadyForTest, pmApproveSubTask, rejectSubTask } = useTasksApi()

const queue = ref<GlobalActiveTask[]>([])
const isLoading = ref(false)
const actioningId = ref<string | null>(null)
const actionType = ref<'approve' | 'reject' | null>(null)

// Product Owner approve: test evidence form modal
const approveModalOpen = ref(false)
const approveTarget = ref<GlobalActiveTask | null>(null)
const approveSubmitting = ref(false)
const approveTestUrl = ref('')
const approveTestSteps = ref('')
const approveUrlRef = ref<HTMLInputElement | null>(null)

const isApproveFormValid = computed(() =>
  approveTestUrl.value.startsWith('http') && approveTestSteps.value.trim().length >= 20
)

// Reject modal state
const rejectModalOpen = ref(false)
const rejectTarget = ref<GlobalActiveTask | null>(null)
const rejectReason = ref('')
const rejectSubmitting = ref(false)
const rejectTextareaRef = ref<HTMLTextAreaElement | null>(null)

async function load() {
  isLoading.value = true
  try {
    queue.value = await getTasksReadyForTest()
  } catch {
    queue.value = []
  } finally {
    isLoading.value = false
  }
}

function openApproveModal(task: GlobalActiveTask) {
  if (approveSubmitting.value) return
  if (rejectModalOpen.value) closeRejectModal()
  approveTarget.value = task
  approveTestUrl.value = ''
  approveTestSteps.value = ''
  approveModalOpen.value = true
  nextTick(() => approveUrlRef.value?.focus())
}

function closeApproveModal() {
  if (approveSubmitting.value) return
  approveModalOpen.value = false
  approveTarget.value = null
  approveTestUrl.value = ''
  approveTestSteps.value = ''
}

async function submitApprove() {
  const t = approveTarget.value
  if (!t || !isApproveFormValid.value) return
  approveSubmitting.value = true
  actioningId.value = t.id
  actionType.value = 'approve'
  try {
    await pmApproveSubTask(t.id, approveTestUrl.value.trim(), approveTestSteps.value.trim())
    queue.value = queue.value.filter(x => x.id !== t.id)
    approveSubmitting.value = false
    actioningId.value = null
    actionType.value = null
    closeApproveModal()
    emit('refresh')
  } catch (e: any) {
    approveSubmitting.value = false
    actioningId.value = null
    actionType.value = null
    alert(e?.data?.message || e?.message || 'Failed to submit for CEO approval')
  }
}

function openRejectModal(task: GlobalActiveTask) {
  if (approveModalOpen.value) closeApproveModal()
  rejectTarget.value = task
  rejectReason.value = ''
  rejectModalOpen.value = true
  nextTick(() => {
    rejectTextareaRef.value?.focus()
  })
}

function closeRejectModal() {
  rejectModalOpen.value = false
  rejectTarget.value = null
  rejectReason.value = ''
}

async function submitReject() {
  if (!rejectTarget.value || rejectReason.value.length < 10) return
  rejectSubmitting.value = true
  actioningId.value = rejectTarget.value.id
  actionType.value = 'reject'
  try {
    await rejectSubTask(rejectTarget.value.id, rejectReason.value)
    queue.value = queue.value.filter(t => t.id !== rejectTarget.value!.id)
    closeRejectModal()
    emit('refresh')
  } catch (e: any) {
    alert(e?.data?.message || e?.message || 'Failed to reject task')
  } finally {
    rejectSubmitting.value = false
    actioningId.value = null
    actionType.value = null
  }
}

onMounted(() => {
  load()
})
</script>
