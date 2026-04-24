<template>
  <!-- Loading skeleton -->
  <section v-if="isLoading" class="rounded-2xl border border-gray-700 bg-gray-800/60 p-6">
    <div class="flex items-center gap-3 mb-4">
      <div class="h-4 w-48 animate-pulse rounded bg-gray-700"/>
      <div class="h-5 w-8 animate-pulse rounded-full bg-gray-700 ml-auto"/>
    </div>
    <div class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-28 animate-pulse rounded-xl bg-gray-700/60"/>
    </div>
  </section>

  <!-- Empty state -->
  <section v-else-if="queue.length === 0">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-amber-400 mb-0.5">Final approval queue</h2>
        <p class="text-xs text-gray-500">Tasks awaiting your final approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-gray-100 dark:bg-gray-700/60 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-500">
        0 pending
      </span>
    </div>
    <div class="rounded-2xl border border-gray-700/50 bg-gray-800/30 px-6 py-10 text-center">
      <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-amber-500/10 border border-amber-500/20">
        <svg class="h-5 w-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
      </div>
      <p class="text-sm font-semibold text-gray-400">All clear</p>
      <p class="text-xs text-gray-600 mt-1">No tasks awaiting final approval</p>
    </div>
  </section>

  <!-- Queue list -->
  <section v-else>
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xs font-semibold uppercase tracking-widest text-amber-400 mb-0.5">Final approval queue</h2>
        <p class="text-xs text-gray-500">Product Owner has tested these tasks — review evidence and give final approval</p>
      </div>
      <span class="text-xs font-bold px-3 py-1 rounded-full bg-amber-500/10 border border-amber-500/25 text-amber-400">
        {{ queue.length }} awaiting approval
      </span>
    </div>

    <div class="space-y-3">
      <div
        v-for="task in queue"
        :key="task.id"
        class="rounded-2xl border border-amber-500/30 bg-gray-800/70 p-4 shadow-md hover:border-amber-400/50 transition-colors"
      >
        <div class="flex items-start justify-between gap-3">
          <!-- Task info (click through to task detail) -->
          <NuxtLink
            :to="`/task/${task.id}`"
            class="flex-1 min-w-0 rounded-xl -m-1 p-1 outline-none transition-colors cursor-pointer hover:bg-amber-500/5 focus-visible:ring-2 focus-visible:ring-amber-500/40 focus-visible:ring-offset-2 focus-visible:ring-offset-gray-900"
          >
            <!-- Project pill + task type -->
            <div class="flex items-center gap-2 mb-1.5">
              <span
                class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-semibold truncate max-w-[140px]"
                :style="{
                  borderColor: task.project_color || '#f59e0b',
                  color: task.project_color || '#f59e0b',
                  backgroundColor: (task.project_color || '#f59e0b') + '18'
                }"
              >
                <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :style="{ backgroundColor: task.project_color || '#f59e0b' }"/>
                {{ task.project_name || 'Unknown Project' }}
              </span>
              <span
                class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border"
                :class="task.task_type === 'BUG'
                  ? 'border-red-500/40 text-red-400 bg-red-500/10'
                  : 'border-indigo-500/40 text-indigo-400 bg-indigo-500/10'"
              >{{ task.task_type || 'TASK' }}</span>
              <span class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border border-amber-500/40 text-amber-400 bg-amber-500/10">
                READY FOR UAT
              </span>
            </div>

            <!-- Title -->
            <p class="text-sm font-semibold text-white line-clamp-2 leading-snug mb-1">{{ task.title }}</p>

            <!-- Code + assignee -->
            <div class="flex items-center gap-3 text-[10px] text-gray-500 mb-2">
              <span v-if="task.code" class="font-mono">{{ task.code }}</span>
              <span v-if="task.assigned_to_display_name || task.assigned_to_email" class="flex items-center gap-1">
                <svg class="w-3 h-3 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                {{ task.assigned_to_display_name || task.assigned_to_email }}
              </span>
            </div>

            <!-- Product Owner test evidence -->
            <template v-if="parsedPayload(task)">
              <div class="mt-2 p-3 rounded-xl bg-gray-900/60 border border-amber-600/20 space-y-2">
                <p class="text-[10px] font-semibold text-amber-400 uppercase tracking-wider mb-1">Product Owner test evidence</p>
                <div v-if="parsedPayload(task)?.test_url" class="flex items-center gap-2">
                  <span class="text-[10px] text-gray-500 w-16 shrink-0">URL</span>
                  <a
                    :href="parsedPayload(task)!.test_url"
                    target="_blank"
                    rel="noopener"
                    class="text-xs text-blue-400 hover:text-blue-300 underline truncate max-w-xs"
                    @click.stop
                  >{{ parsedPayload(task)!.test_url }}</a>
                </div>
                <div v-if="parsedPayload(task)?.test_steps" class="flex items-start gap-2">
                  <span class="text-[10px] text-gray-500 w-16 shrink-0 mt-0.5">Steps</span>
                  <pre class="text-xs text-gray-300 whitespace-pre-wrap break-words max-w-xs line-clamp-3">{{ parsedPayload(task)!.test_steps }}</pre>
                </div>
              </div>
            </template>
          </NuxtLink>

          <!-- Action buttons -->
          <div class="flex flex-col gap-2 flex-shrink-0">
            <button
              :disabled="actioningId === task.id || approveSubmitting"
              class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-semibold hover:bg-emerald-500/20 hover:border-emerald-400/50 transition-colors disabled:opacity-50"
              @click="openApproveModal(task)"
            >
              <span>👑</span>
              APPROVE
            </button>
            <button
              :disabled="actioningId === task.id"
              class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-xs font-semibold hover:bg-red-500/20 hover:border-red-400/50 transition-colors disabled:opacity-50"
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

  <!-- CEO Final Approve confirm modal -->
  <Teleport to="body">
    <div
      v-if="approveModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
      @click.self="closeApproveModal"
    >
      <div class="w-full max-w-md rounded-2xl border border-amber-500/30 bg-gray-900 shadow-2xl p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center flex-shrink-0">
            <span class="text-lg" aria-hidden="true">👑</span>
          </div>
          <div>
            <h3 class="text-sm font-bold text-white">Final Approval — Mark as Done?</h3>
            <p class="text-xs text-gray-500 truncate max-w-[280px]">{{ approveTarget?.title }}</p>
          </div>
        </div>
        <p class="text-xs text-gray-400 mb-4">
          This marks the task as <span class="text-emerald-400 font-medium">COMPLETED</span>. Continue only if you have verified the Product Owner's test evidence.
        </p>
        <div class="flex items-center justify-end gap-3">
          <button
            type="button"
            :disabled="approveSubmitting"
            class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-400 hover:border-gray-600 hover:text-gray-200 transition-colors disabled:opacity-50"
            @click="closeApproveModal"
          >Cancel</button>
          <button
            type="button"
            :disabled="approveSubmitting"
            class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-amber-500/15 border border-amber-500/40 text-amber-400 text-xs font-bold hover:bg-amber-500/25 transition-colors disabled:opacity-50"
            @click="submitApprove"
          >
            <svg v-if="approveSubmitting" class="w-3.5 h-3.5 animate-spin shrink-0" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <span v-if="approveSubmitting">Approving…</span>
            <span v-else>👑 Yes, Final Approve</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Reject Modal -->
  <Teleport to="body">
    <div
      v-if="rejectModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
      @click.self="closeRejectModal"
    >
      <div class="w-full max-w-md rounded-2xl border border-red-500/30 bg-gray-900 shadow-2xl p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-8 h-8 rounded-lg bg-red-500/15 border border-red-500/30 flex items-center justify-center flex-shrink-0">
            <svg class="w-4 h-4 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
          <div>
            <h3 class="text-sm font-bold text-white">Reject — Return to Product Owner</h3>
            <p class="text-xs text-gray-500 truncate max-w-[280px]">{{ rejectTarget?.title }}</p>
          </div>
        </div>
        <p class="text-xs text-gray-400 mb-3">Explain what failed or needs improvement. The task will be returned to IN_PROGRESS.</p>
        <textarea
          ref="rejectTextareaRef"
          v-model="rejectReason"
          rows="4"
          placeholder="Describe what the Product Owner needs to fix (min. 10 characters)…"
          class="w-full rounded-xl border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 focus:border-red-500/50 focus:outline-none focus:ring-1 focus:ring-red-500/30 resize-none"
        />
        <p v-if="rejectReason.length > 0 && rejectReason.length < 10" class="text-[11px] text-red-400 mt-1">
          At least {{ 10 - rejectReason.length }} more character(s) required
        </p>
        <div class="flex items-center justify-end gap-3 mt-4">
          <button class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-400 hover:border-gray-600 hover:text-gray-200 transition-colors" @click="closeRejectModal">Cancel</button>
          <button
            :disabled="rejectReason.length < 10 || rejectSubmitting"
            class="px-4 py-2 rounded-lg bg-red-500/15 border border-red-500/40 text-red-400 text-xs font-bold hover:bg-red-500/25 transition-colors disabled:opacity-50"
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
import { nextTick } from 'vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { GlobalActiveTask } from '~/core/modules/tasks/infrastructure/tasks-api'

const emit = defineEmits<{ (e: 'refresh'): void }>()

const { getTasksReadyForCEOApproval, approveSubTask, rejectSubTask } = useTasksApi()

const queue = ref<GlobalActiveTask[]>([])
const isLoading = ref(false)
const actioningId = ref<string | null>(null)

const approveModalOpen = ref(false)
const approveTarget = ref<GlobalActiveTask | null>(null)
const approveSubmitting = ref(false)

const rejectModalOpen = ref(false)
const rejectTarget = ref<GlobalActiveTask | null>(null)
const rejectReason = ref('')
const rejectSubmitting = ref(false)
const rejectTextareaRef = ref<HTMLTextAreaElement | null>(null)

function parsedPayload(task: GlobalActiveTask): { test_url?: string; test_steps?: string } | null {
  const raw = (task as any).uat_payload
  if (!raw) return null
  try {
    return typeof raw === 'string' ? JSON.parse(raw) : raw
  } catch {
    return null
  }
}

async function load() {
  isLoading.value = true
  try {
    queue.value = await getTasksReadyForCEOApproval()
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
  approveModalOpen.value = true
}

function closeApproveModal() {
  if (approveSubmitting.value) return
  approveModalOpen.value = false
  approveTarget.value = null
}

async function submitApprove() {
  const t = approveTarget.value
  if (!t) return
  approveSubmitting.value = true
  actioningId.value = t.id
  try {
    await approveSubTask(t.id)
    queue.value = queue.value.filter(x => x.id !== t.id)
    closeApproveModal()
    emit('refresh')
  } catch (e: any) {
    alert(e?.data?.message || e?.message || 'Failed to approve task')
  } finally {
    approveSubmitting.value = false
    actioningId.value = null
  }
}

function openRejectModal(task: GlobalActiveTask) {
  if (approveModalOpen.value) closeApproveModal()
  rejectTarget.value = task
  rejectReason.value = ''
  rejectModalOpen.value = true
  nextTick(() => rejectTextareaRef.value?.focus())
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
  }
}

onMounted(() => load())
</script>
