<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click.self="close"
        @keydown.escape="close"
      >
        <div class="relative w-full max-w-lg rounded-2xl border border-blue-500/30 bg-gray-800 shadow-2xl" @click.stop>
          <!-- Header -->
          <div class="border-b border-gray-700 px-6 py-4">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-blue-500/15 border border-blue-500/30 flex items-center justify-center flex-shrink-0">
                <svg class="w-4 h-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 4H6a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-2m-4-1v8m0 0l3-3m-3 3L9 8m-5 5h2.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293h3.172a1 1 0 00.707-.293l2.414-2.414A1 1 0 0014.414 13H17"/>
                </svg>
              </div>
              <div>
                <h2 class="text-sm font-bold text-white">Outsource to Another Team</h2>
                <p class="text-xs text-gray-500 mt-0.5">Send a B2B work request to another team</p>
              </div>
              <button @click="close" class="ml-auto text-gray-500 hover:text-gray-300 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Form -->
          <form class="p-6 space-y-4" @submit.prevent="submit">
            <!-- Error -->
            <div v-if="error" class="flex items-start gap-2 p-3 bg-red-900/30 border border-red-500/40 rounded-xl text-xs text-red-300">
              <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ error }}
            </div>

            <!-- Target Team -->
            <div>
              <label class="block text-xs font-medium text-gray-400 mb-1.5">
                Target Team <span class="text-red-400">*</span>
              </label>
              <div v-if="isLoadingTeams" class="flex items-center gap-2 text-xs text-gray-500 py-2">
                <svg class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                Loading teams…
              </div>
              <select
                v-else
                v-model.number="form.target_team_id"
                required
                class="w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-colors"
              >
                <option value="" disabled>— Select a team —</option>
                <option v-for="team in otherTeams" :key="team.id" :value="team.id">
                  {{ team.name }}
                </option>
              </select>
              <p v-if="otherTeams.length === 0 && !isLoadingTeams" class="text-xs text-amber-400 mt-1">
                No other teams found. Ask your CEO to create teams.
              </p>
            </div>

            <!-- Title -->
            <div>
              <label class="block text-xs font-medium text-gray-400 mb-1.5">
                Work Title <span class="text-red-400">*</span>
              </label>
              <input
                v-model="form.title"
                type="text"
                required
                placeholder="e.g. Implement OAuth2 integration…"
                class="w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-colors"
              />
            </div>

            <!-- Description -->
            <div>
              <label class="block text-xs font-medium text-gray-400 mb-1.5">Description</label>
              <textarea
                v-model="form.description"
                rows="3"
                placeholder="Describe the scope, requirements, expected deliverables…"
                class="w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-colors resize-none"
              />
            </div>

            <!-- Estimated Minutes -->
            <div>
              <label class="block text-xs font-medium text-gray-400 mb-1.5">
                Estimated Duration (minutes) <span class="text-red-400">*</span>
              </label>
              <div class="flex items-center gap-3">
                <input
                  v-model.number="form.estimated_minutes"
                  type="number"
                  min="1"
                  required
                  placeholder="e.g. 480"
                  class="w-32 rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-colors"
                />
                <span class="text-xs text-gray-500">
                  {{ form.estimated_minutes > 0 ? formatMinutes(form.estimated_minutes) : '' }}
                </span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-3 pt-2">
              <button
                type="button"
                @click="close"
                class="flex-1 inline-flex items-center justify-center rounded-lg border border-gray-600 bg-transparent px-4 py-2.5 text-sm font-semibold text-gray-300 transition-colors hover:bg-gray-700 hover:text-gray-900 dark:text-white"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isSubmitting || !form.title || !form.target_team_id || !form.estimated_minutes"
                class="flex-1 inline-flex items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-blue-100 dark:from-blue-600 to-indigo-100 dark:to-indigo-600 px-4 py-2.5 text-sm font-semibold text-gray-900 dark:text-white shadow-lg transition-all hover:from-blue-200 dark:hover:from-blue-500 hover:to-indigo-200 dark:hover:to-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <svg v-if="isSubmitting" class="h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                Send Request
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useTeamsApi, type Team } from '~/core/modules/teams/infrastructure/teams-api'
import { useB2BApi } from '~/core/modules/b2b/infrastructure/b2b-api'

const props = defineProps<{
  modelValue: boolean
  prefillTitle?: string
  prefillDescription?: string
  prefillMinutes?: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
  (e: 'created'): void
}>()

const { getTeams } = useTeamsApi()
const { createRequest } = useB2BApi()
const { currentUser } = useAuth()

const isLoadingTeams = ref(true)
const allTeams = ref<Team[]>([])
const isSubmitting = ref(false)
const error = ref('')

const form = reactive({
  title: '',
  description: '',
  estimated_minutes: 0,
  target_team_id: 0 as number,
})

// Filter out the user's own team from the dropdown
const myTeamId = computed(() => {
  const u = currentUser.value as any
  return u?.team_id ?? null
})

const otherTeams = computed(() =>
  allTeams.value.filter(t => t.id !== myTeamId.value)
)

const formatMinutes = (mins: number) => {
  if (!mins || mins <= 0) return ''
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h === 0) return `${m}m`
  if (m === 0) return `${h}h`
  return `${h}h ${m}m`
}

const close = () => {
  emit('update:modelValue', false)
}

const submit = async () => {
  if (!form.title || !form.target_team_id || form.estimated_minutes <= 0) return
  error.value = ''
  isSubmitting.value = true
  try {
    await createRequest({
      title: form.title,
      description: form.description,
      estimated_minutes: form.estimated_minutes,
      target_team_id: form.target_team_id,
    })
    emit('created')
    close()
  } catch (err: any) {
    error.value = err?.data?.message ?? err?.message ?? 'Failed to send request'
  } finally {
    isSubmitting.value = false
  }
}

// Pre-fill form when props change (opened from task detail)
watch(() => props.modelValue, (open) => {
  if (open) {
    form.title = props.prefillTitle ?? ''
    form.description = props.prefillDescription ?? ''
    form.estimated_minutes = props.prefillMinutes ?? 0
    form.target_team_id = 0
    error.value = ''
  }
})

// Load teams once
onMounted(async () => {
  try {
    allTeams.value = await getTeams()
  } catch {
    // silent
  } finally {
    isLoadingTeams.value = false
  }
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
