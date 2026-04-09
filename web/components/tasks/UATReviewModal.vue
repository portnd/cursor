<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click.self="closeModal"
      >
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 translate-y-2"
        >
          <div v-if="modelValue" class="w-full max-w-xl bg-gray-900 border border-gray-700/80 rounded-2xl shadow-2xl overflow-hidden">
            <!-- Header -->
            <div class="flex items-center justify-between px-6 py-4 border-b border-gray-700/60 bg-gray-800/40">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-violet-500/15 border border-violet-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-sm font-bold text-white">UAT Review</h3>
                  <p class="text-xs text-gray-500 mt-0.5 truncate max-w-xs">{{ feature?.title }}</p>
                </div>
              </div>
              <button
                @click="closeModal"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-500 hover:text-gray-900 dark:text-white hover:bg-gray-700/60 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <!-- UAT Payload (if available) -->
            <div v-if="uatPayload" class="px-6 py-5 space-y-4">
              <!-- Staging URL -->
              <div class="rounded-xl border border-indigo-500/20 bg-indigo-950/20 px-4 py-4">
                <p class="text-[10px] font-bold uppercase tracking-wider text-indigo-400 mb-2">Staging Environment</p>
                <div class="flex items-center gap-3">
                  <p class="flex-1 text-sm font-mono text-indigo-200 truncate">{{ uatPayload.staging_url }}</p>
                  <a
                    :href="uatPayload.staging_url"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-1.5 shrink-0 px-3 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold transition-colors"
                  >
                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                    </svg>
                    Test in Browser
                  </a>
                </div>
              </div>

              <!-- Test Credentials -->
              <div v-if="uatPayload.test_credentials" class="rounded-xl border border-gray-700/60 bg-gray-800/30 px-4 py-3">
                <p class="text-[10px] font-bold uppercase tracking-wider text-gray-500 mb-2">Testing Instructions / Credentials</p>
                <p class="text-sm text-gray-300 whitespace-pre-wrap leading-relaxed">{{ uatPayload.test_credentials }}</p>
              </div>

              <!-- Release Notes -->
              <div v-if="uatPayload.release_notes" class="rounded-xl border border-gray-700/60 bg-gray-800/30 px-4 py-3">
                <p class="text-[10px] font-bold uppercase tracking-wider text-gray-500 mb-2">Release Notes</p>
                <p class="text-sm text-gray-300 whitespace-pre-wrap leading-relaxed">{{ uatPayload.release_notes }}</p>
              </div>
            </div>

            <!-- No UAT payload yet -->
            <div v-else class="px-6 py-8 text-center">
              <div class="w-12 h-12 rounded-2xl bg-gray-700/40 flex items-center justify-center mx-auto mb-3">
                <svg class="w-6 h-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <p class="text-sm font-semibold text-gray-300">No UAT Details Yet</p>
              <p class="text-xs text-gray-500 mt-1">The developer has not submitted UAT details for this feature.</p>
            </div>

            <!-- Reject reason input (shown when rejecting) -->
            <div v-if="showRejectForm" class="px-6 pb-4">
              <div class="rounded-xl border border-red-500/30 bg-red-950/20 px-4 py-3">
                <label class="block text-xs font-semibold text-red-300 mb-2">Bug / Rejection Reason <span class="text-red-400">*</span></label>
                <textarea
                  v-model="rejectReason"
                  rows="3"
                  placeholder="Describe what failed or what needs to be fixed…"
                  class="w-full rounded-lg border border-red-700/40 bg-gray-900/60 px-3 py-2 text-sm text-white placeholder-gray-600 outline-none focus:border-red-500/60 transition-colors resize-none"
                />
              </div>
            </div>

            <!-- Error -->
            <div v-if="actionError" class="px-6 pb-3">
              <p class="text-xs text-red-400 bg-red-900/20 border border-red-500/30 rounded-lg px-3 py-2">{{ actionError }}</p>
            </div>

            <!-- Action Buttons -->
            <div class="px-6 py-4 border-t border-gray-700/60 space-y-3">
              <!-- Approve -->
              <button
                v-if="!showRejectForm"
                @click="handleApprove"
                :disabled="actioning"
                class="w-full flex items-center justify-center gap-3 py-3.5 rounded-xl font-black text-base tracking-wide bg-emerald-100 dark:bg-emerald-600 hover:bg-emerald-100 dark:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed text-gray-900 dark:text-white transition-all shadow-lg shadow-emerald-900/30 hover:shadow-emerald-800/40"
              >
                <svg v-if="actioning && approving" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                </svg>
                {{ approving ? 'Approving…' : 'APPROVE FEATURE' }}
              </button>

              <!-- Reject row -->
              <div v-if="!showRejectForm" class="flex items-center gap-3">
                <button
                  @click="showRejectForm = true"
                  :disabled="actioning"
                  class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl font-bold text-sm border border-red-300 dark:border-red-500/40 text-red-400 bg-red-100 dark:bg-red-500/10 hover:bg-red-100 dark:bg-red-500/20 hover:border-red-300 dark:border-red-500/60 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                  </svg>
                  REJECT — Bug Found
                </button>
                <button
                  @click="closeModal"
                  class="px-4 py-2.5 rounded-xl text-sm font-semibold text-gray-500 hover:text-gray-900 dark:text-white hover:bg-gray-700/60 transition-colors"
                >
                  Close
                </button>
              </div>

              <!-- Confirm reject -->
              <div v-if="showRejectForm" class="flex items-center gap-3">
                <button
                  @click="handleReject"
                  :disabled="actioning || !rejectReason.trim()"
                  class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl font-bold text-sm bg-red-100 dark:bg-red-600 hover:bg-red-100 dark:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed text-gray-900 dark:text-white transition-all"
                >
                  <svg v-if="actioning && !approving" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                  {{ !approving && actioning ? 'Rejecting…' : 'Confirm Rejection' }}
                </button>
                <button
                  @click="showRejectForm = false; rejectReason = ''"
                  :disabled="actioning"
                  class="px-4 py-2.5 rounded-xl text-sm font-semibold text-gray-500 hover:text-gray-900 dark:text-white hover:bg-gray-700/60 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { FeatureRoadmapItem, UATPayload } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{
  modelValue: boolean
  feature: FeatureRoadmapItem | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
  (e: 'success'): void
}>()

const { approveTask, rejectTask } = useTasksApi()

const showRejectForm = ref(false)
const rejectReason = ref('')
const actioning = ref(false)
const approving = ref(false)
const actionError = ref('')

const uatPayload = computed<UATPayload | null>(() => {
  const raw = props.feature?.uat_payload
  if (!raw) return null
  if (typeof raw === 'object') return raw as UATPayload
  try { return JSON.parse(raw as unknown as string) as UATPayload } catch { return null }
})

watch(() => props.modelValue, (open) => {
  if (open) {
    showRejectForm.value = false
    rejectReason.value = ''
    actionError.value = ''
    approving.value = false
  }
})

function closeModal() {
  if (!actioning.value) emit('update:modelValue', false)
}

async function handleApprove() {
  if (!props.feature) return
  actioning.value = true
  approving.value = true
  actionError.value = ''
  try {
    await approveTask(props.feature.id)
    emit('update:modelValue', false)
    emit('success')
  } catch (e: any) {
    actionError.value = e?.data?.message || e?.message || 'Failed to approve feature'
  } finally {
    actioning.value = false
    approving.value = false
  }
}

async function handleReject() {
  if (!props.feature || !rejectReason.value.trim()) return
  actioning.value = true
  approving.value = false
  actionError.value = ''
  try {
    await rejectTask(props.feature.id, rejectReason.value.trim())
    emit('update:modelValue', false)
    emit('success')
  } catch (e: any) {
    actionError.value = e?.data?.message || e?.message || 'Failed to reject feature'
  } finally {
    actioning.value = false
  }
}
</script>
