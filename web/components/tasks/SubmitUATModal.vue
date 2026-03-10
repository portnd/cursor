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
        @click.self="$emit('update:modelValue', false)"
      >
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 translate-y-2"
        >
          <div v-if="modelValue" class="w-full max-w-lg bg-gray-900 border border-gray-700/80 rounded-2xl shadow-2xl">
            <!-- Header -->
            <div class="flex items-center justify-between px-6 py-4 border-b border-gray-700/60">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center">
                  <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-sm font-bold text-white">Submit for UAT</h3>
                  <p class="text-xs text-gray-500 mt-0.5 truncate max-w-xs">{{ feature?.title }}</p>
                </div>
              </div>
              <button
                @click="$emit('update:modelValue', false)"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-500 hover:text-white hover:bg-gray-700/60 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <!-- Body -->
            <div class="px-6 py-5 space-y-4">
              <p class="text-xs text-gray-400 bg-amber-500/10 border border-amber-500/20 rounded-lg px-3 py-2">
                All child tasks are complete. Provide the staging environment details for PM/CEO to review and accept.
              </p>

              <!-- Staging URL -->
              <div>
                <label class="block text-xs font-semibold text-gray-300 mb-1.5">
                  Staging URL <span class="text-red-400">*</span>
                </label>
                <input
                  v-model="form.staging_url"
                  type="url"
                  placeholder="https://staging.example.com"
                  class="w-full rounded-lg border bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 outline-none transition-colors"
                  :class="urlError
                    ? 'border-red-500/60 focus:border-red-400'
                    : 'border-gray-700 focus:border-violet-500/60'"
                  @blur="validateURL"
                />
                <p v-if="urlError" class="mt-1 text-xs text-red-400">{{ urlError }}</p>
              </div>

              <!-- Test Credentials -->
              <div>
                <label class="block text-xs font-semibold text-gray-300 mb-1.5">Testing Instructions / Credentials</label>
                <textarea
                  v-model="form.test_credentials"
                  rows="3"
                  placeholder="e.g. admin@example.com / Test@1234 — test the checkout flow on the staging site"
                  class="w-full rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 outline-none focus:border-violet-500/60 transition-colors resize-none"
                />
              </div>

              <!-- Release Notes -->
              <div>
                <label class="block text-xs font-semibold text-gray-300 mb-1.5">Release Notes</label>
                <textarea
                  v-model="form.release_notes"
                  rows="4"
                  placeholder="What was built? What changed? Any known limitations?"
                  class="w-full rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2.5 text-sm text-white placeholder-gray-600 outline-none focus:border-violet-500/60 transition-colors resize-none"
                />
              </div>

              <!-- Error -->
              <p v-if="submitError" class="text-xs text-red-400 bg-red-900/20 border border-red-500/30 rounded-lg px-3 py-2">{{ submitError }}</p>
            </div>

            <!-- Footer -->
            <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-gray-700/60">
              <button
                @click="$emit('update:modelValue', false)"
                class="px-4 py-2 rounded-lg text-sm font-semibold text-gray-400 hover:text-white hover:bg-gray-700/60 transition-colors"
              >
                Cancel
              </button>
              <button
                @click="handleSubmit"
                :disabled="submitting || !form.staging_url"
                class="inline-flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-bold bg-amber-500 text-gray-900 hover:bg-amber-400 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                <svg v-if="submitting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                </svg>
                {{ submitting ? 'Submitting…' : 'Submit for UAT' }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { FeatureRoadmapItem } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{
  modelValue: boolean
  feature: FeatureRoadmapItem | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
  (e: 'success'): void
}>()

const { submitUAT } = useTasksApi()

const form = reactive({
  staging_url: '',
  test_credentials: '',
  release_notes: '',
})

const urlError = ref('')
const submitError = ref('')
const submitting = ref(false)

function validateURL() {
  if (!form.staging_url) {
    urlError.value = ''
    return
  }
  const valid = form.staging_url.startsWith('http://') || form.staging_url.startsWith('https://')
  urlError.value = valid ? '' : 'Must be a valid http:// or https:// URL'
}

watch(() => props.modelValue, (open) => {
  if (open) {
    form.staging_url = ''
    form.test_credentials = ''
    form.release_notes = ''
    urlError.value = ''
    submitError.value = ''
  }
})

async function handleSubmit() {
  validateURL()
  if (urlError.value || !form.staging_url) return
  if (!props.feature) return

  submitting.value = true
  submitError.value = ''
  try {
    await submitUAT(props.feature.id, {
      staging_url: form.staging_url,
      test_credentials: form.test_credentials,
      release_notes: form.release_notes,
    })
    emit('update:modelValue', false)
    emit('success')
  } catch (e: any) {
    submitError.value = e?.data?.message || e?.message || 'Failed to submit UAT'
  } finally {
    submitting.value = false
  }
}
</script>
