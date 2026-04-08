<template>
  <!-- Backdrop -->
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
        @click.self="onBackdropClick"
      >
        <!-- Modal panel -->
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="scale-95 opacity-0"
          enter-to-class="scale-100 opacity-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="scale-100 opacity-100"
          leave-to-class="scale-95 opacity-0"
        >
          <div
            v-if="isOpen"
            class="w-full max-w-lg rounded-2xl border border-gray-700 bg-gray-900 shadow-2xl"
          >
            <!-- Header -->
            <div class="flex items-center justify-between border-b border-gray-700 px-6 py-4">
              <div>
                <h3 class="text-lg font-bold text-white">Daily Check-in</h3>
                <p class="text-xs text-gray-400 mt-0.5">{{ todayFormatted }} — Let the team know what you're up to</p>
              </div>
              <button
                v-if="!forced"
                class="rounded-lg p-1.5 text-gray-400 hover:bg-gray-800 hover:text-white transition"
                @click="close"
              >
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Body -->
            <form class="flex flex-col gap-5 px-6 py-5" @submit.prevent="handleSubmit">
              <!-- Yesterday summary -->
              <div class="flex flex-col gap-1.5">
                <label class="text-sm font-semibold text-gray-200" for="yesterday">
                  What did you accomplish yesterday?
                  <span class="ml-1 text-red-400">*</span>
                </label>
                <div class="relative">
                  <textarea
                    id="yesterday"
                    v-model="form.yesterday_summary"
                    rows="3"
                    placeholder="e.g. Finished the authentication flow, reviewed 2 PRs…"
                    class="w-full resize-none rounded-lg border border-gray-700 bg-gray-800 px-3 py-2 pr-10 text-sm text-gray-200 placeholder-gray-600 focus:border-indigo-500 focus:outline-none"
                    :class="{ 'border-red-600': errors.yesterday_summary }"
                  />
                  <button
                    v-if="voiceSupported"
                    type="button"
                    :title="yesterdayListening ? 'หยุดฟัง' : 'พูดแทนพิมพ์'"
                    class="absolute right-2 top-2 flex h-6 w-6 items-center justify-center rounded-md transition-colors"
                    :class="yesterdayListening ? 'bg-red-500/20 text-red-400 ring-1 ring-red-500/50 animate-pulse' : 'text-gray-500 hover:bg-gray-700 hover:text-gray-300'"
                    @click="toggleYesterday"
                  >
                    <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="2" width="6" height="12" rx="3"/>
                      <path stroke-linecap="round" d="M5 10a7 7 0 0014 0"/>
                      <line x1="12" y1="17" x2="12" y2="22"/>
                      <line x1="9" y1="22" x2="15" y2="22"/>
                    </svg>
                  </button>
                </div>
                <p v-if="errors.yesterday_summary" class="text-xs text-red-400">{{ errors.yesterday_summary }}</p>
                <p v-if="voiceErrorYesterday" class="text-xs text-amber-400 flex items-center gap-1">
                  <span>⚠</span> {{ voiceErrorYesterday }}
                </p>
                <p v-if="yesterdayListening" class="text-xs text-indigo-400 animate-pulse">🎙 กำลังฟัง… พูดได้เลย</p>
              </div>

              <!-- Today tasks -->
              <div class="flex flex-col gap-1.5">
                <label class="text-sm font-semibold text-gray-200">
                  Which task codes are you working on today?
                  <span class="text-xs font-normal text-gray-500 ml-1">(optional, press Enter to add)</span>
                </label>
                <!-- Tag input -->
                <div
                  class="flex min-h-[2.5rem] flex-wrap gap-1.5 rounded-lg border border-gray-700 bg-gray-800 px-3 py-2 focus-within:border-indigo-500"
                >
                  <span
                    v-for="(tid, i) in form.today_task_ids"
                    :key="i"
                    class="flex items-center gap-1 rounded bg-indigo-900/70 px-2 py-0.5 text-xs font-mono text-indigo-300"
                  >
                    {{ tid }}
                    <button
                      type="button"
                      class="text-indigo-400 hover:text-white"
                      @click="removeTaskId(i)"
                    >×</button>
                  </span>
                  <input
                    v-model="taskInput"
                    type="text"
                    placeholder="e.g. SENT-042"
                    class="min-w-[8rem] flex-1 bg-transparent text-sm text-gray-200 placeholder-gray-600 outline-none"
                    @keydown.enter.prevent="addTaskId"
                    @keydown.backspace="onBackspace"
                    @keydown.comma.prevent="addTaskId"
                  />
                </div>
              </div>

              <!-- Blocker -->
              <div class="flex flex-col gap-1.5">
                <label class="text-sm font-semibold text-gray-200" for="blocker">
                  Any blockers?
                  <span class="text-xs font-normal text-gray-500 ml-1">(leave empty if none)</span>
                </label>
                <div class="relative">
                  <textarea
                    id="blocker"
                    v-model="form.blocker"
                    rows="2"
                    placeholder="e.g. Waiting on design assets for the dashboard…"
                    class="w-full resize-none rounded-lg border border-gray-700 bg-gray-800 px-3 py-2 pr-10 text-sm text-gray-200 placeholder-gray-600 focus:border-indigo-500 focus:outline-none"
                    :class="form.blocker ? 'border-red-700 bg-red-900/10' : ''"
                  />
                  <button
                    v-if="voiceSupported"
                    type="button"
                    :title="blockerListening ? 'หยุดฟัง' : 'พูดแทนพิมพ์'"
                    class="absolute right-2 top-2 flex h-6 w-6 items-center justify-center rounded-md transition-colors"
                    :class="blockerListening ? 'bg-red-500/20 text-red-400 ring-1 ring-red-500/50 animate-pulse' : 'text-gray-500 hover:bg-gray-700 hover:text-gray-300'"
                    @click="toggleBlocker"
                  >
                    <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="2" width="6" height="12" rx="3"/>
                      <path stroke-linecap="round" d="M5 10a7 7 0 0014 0"/>
                      <line x1="12" y1="17" x2="12" y2="22"/>
                      <line x1="9" y1="22" x2="15" y2="22"/>
                    </svg>
                  </button>
                </div>
                <p v-if="form.blocker" class="text-xs text-red-400 font-medium">
                  Blocker will be highlighted for the team lead.
                </p>
                <p v-if="voiceErrorBlocker" class="text-xs text-amber-400 flex items-center gap-1">
                  <span>⚠</span> {{ voiceErrorBlocker }}
                </p>
                <p v-if="blockerListening" class="text-xs text-indigo-400 animate-pulse">🎙 กำลังฟัง… พูดได้เลย</p>
              </div>

              <!-- API error -->
              <div
                v-if="store.error"
                class="rounded-lg border border-red-700 bg-red-900/40 px-3 py-2 text-sm text-red-300"
              >
                {{ store.error }}
              </div>

              <!-- Actions -->
              <div class="flex items-center justify-end gap-3 border-t border-gray-700 pt-4">
                <button
                  v-if="!forced"
                  type="button"
                  class="rounded-lg px-4 py-2 text-sm text-gray-400 hover:bg-gray-800 hover:text-white transition"
                  @click="close"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="store.submitting"
                  class="flex items-center gap-2 rounded-lg bg-indigo-600 px-5 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500 disabled:opacity-50"
                >
                  <svg v-if="store.submitting" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                  </svg>
                  {{ store.submitting ? 'Submitting…' : 'Submit Check-in' }}
                </button>
              </div>
            </form>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { usePulseStore } from '../store/pulse-store'
import { useVoiceInput } from '~/composables/useVoiceInput'

// ── Props & Emits ─────────────────────────────────────────────────────────────

interface Props {
  /**
   * When true the modal cannot be dismissed until the standup is submitted.
   * Intended for the auto-trigger on first page load of the day.
   */
  forced?: boolean
}

const props = withDefaults(defineProps<Props>(), { forced: false })
const emit = defineEmits<{ (e: 'close'): void; (e: 'submitted'): void }>()

// ── State ─────────────────────────────────────────────────────────────────────

const store = usePulseStore()
const isOpen = ref(false)

const today = new Date().toISOString().slice(0, 10) // YYYY-MM-DD

const todayFormatted = computed(() => {
  const d = new Date(today + 'T00:00:00')
  return d.toLocaleDateString('en-GB', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
})

const form = reactive({
  yesterday_summary: '',
  today_task_ids: [] as string[],
  blocker: '',
})

// ── Voice input ───────────────────────────────────────────────────────────────
const { isListening: yesterdayListening, isSupported: voiceSupported, error: voiceErrorYesterday, toggle: toggleYesterday } =
  useVoiceInput((text) => {
    const sep = form.yesterday_summary && !form.yesterday_summary.endsWith(' ') ? ' ' : ''
    form.yesterday_summary += sep + text
  })

const { isListening: blockerListening, error: voiceErrorBlocker, toggle: toggleBlocker } =
  useVoiceInput((text) => {
    const sep = form.blocker && !form.blocker.endsWith(' ') ? ' ' : ''
    form.blocker += sep + text
  })

const taskInput = ref('')
const errors = reactive({ yesterday_summary: '' })

// ── Task tag helpers ──────────────────────────────────────────────────────────

function addTaskId() {
  const val = taskInput.value.trim().replace(/,/g, '')
  if (val && !form.today_task_ids.includes(val)) {
    form.today_task_ids.push(val)
  }
  taskInput.value = ''
}

function removeTaskId(index: number) {
  form.today_task_ids.splice(index, 1)
}

function onBackspace() {
  if (!taskInput.value && form.today_task_ids.length) {
    form.today_task_ids.pop()
  }
}

// ── Validation ────────────────────────────────────────────────────────────────

function validate(): boolean {
  errors.yesterday_summary = ''
  if (!form.yesterday_summary.trim()) {
    errors.yesterday_summary = 'Please describe what you did yesterday.'
    return false
  }
  return true
}

// ── Submit ────────────────────────────────────────────────────────────────────

async function handleSubmit() {
  if (!validate()) return
  // Flush any pending tag input
  if (taskInput.value.trim()) addTaskId()

  const ok = await store.submitStandup({
    date: today,
    yesterday_summary: form.yesterday_summary.trim(),
    today_task_ids: [...form.today_task_ids],
    blocker: form.blocker.trim(),
  })

  if (ok) {
    emit('submitted')
    forceClose()
  }
}

// ── Open / Close ──────────────────────────────────────────────────────────────

function open() {
  store.error = null
  isOpen.value = true
}

/** Force-close regardless of forced prop — used after a successful submit */
function forceClose() {
  isOpen.value = false
  emit('close')
}

function close() {
  if (props.forced) return
  forceClose()
}

function onBackdropClick() {
  if (!props.forced) close()
}

// ── Auto-trigger for forced mode ──────────────────────────────────────────────

onMounted(async () => {
  if (props.forced) {
    // Load today's pulse to check if the current user already checked in
    await store.fetchDailyPulse(today)
    if (!store.hasCheckedInToday) {
      open()
    }
  }
})

// ── Expose open/close for parent components ───────────────────────────────────
defineExpose({ open, close })
</script>
