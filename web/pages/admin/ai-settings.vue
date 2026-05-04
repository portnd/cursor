<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-white">AI CONFIGURATION</h1>
      <p class="text-sm text-gray-400 mt-1">GLM Model &amp; behavior settings</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลดการตั้งค่า...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded p-4 text-red-400">
      {{ error }}
    </div>

    <!-- Main Content -->
    <div v-else>
      <!-- Current Status -->
      <div class="mb-6 bg-gray-800 border border-gray-700 rounded p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs text-gray-500 uppercase mb-1">Active Configuration</p>
            <div class="flex items-center gap-2">
              <span class="text-lg font-bold text-white">{{ config.active_model }}</span>
              <span class="px-2 py-0.5 bg-green-700 text-green-100 text-xs font-bold rounded">
                LIVE
              </span>
              <span class="px-2 py-0.5 text-xs font-medium rounded" :class="getModelBadge(formData.active_model).class">
                {{ getModelBadge(formData.active_model).label }}
              </span>
            </div>
          </div>
          <div class="text-right">
            <p class="text-xs text-gray-500">Last Updated</p>
            <p class="text-sm text-gray-300">{{ formatDateTime(config.updated_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Model Cards (visual selector) -->
      <div class="mb-6">
        <h2 class="text-sm font-bold text-gray-400 uppercase mb-3">Select GLM Model</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
          <div
            v-for="m in modelCards"
            :key="m.id"
            @click="formData.active_model = m.id"
            class="cursor-pointer border rounded-lg p-4 transition-all"
            :class="formData.active_model === m.id
              ? 'border-blue-500 bg-blue-900/20 ring-1 ring-blue-500'
              : 'border-gray-700 bg-gray-800 hover:border-gray-500'"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-bold text-white">{{ m.name }}</span>
              <span v-if="formData.active_model === m.id" class="text-blue-400 text-lg">✓</span>
            </div>
            <p class="text-xs text-gray-400 mb-2">{{ m.description }}</p>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="tag in m.tags"
                :key="tag"
                class="px-1.5 py-0.5 text-[10px] font-medium rounded"
                :class="getTagClass(tag)"
              >
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Configuration Form -->
      <div class="space-y-6">

        <!-- Temperature Control -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase mb-3">Temperature</h2>

          <div class="flex items-center justify-between mb-2">
            <span class="text-sm text-gray-400">Value</span>
            <span class="text-xl font-bold text-white">{{ formData.temperature.toFixed(2) }}</span>
          </div>

          <input
            v-model.number="formData.temperature"
            type="range"
            min="0"
            max="1"
            step="0.05"
            class="w-full h-2 bg-gray-700 rounded appearance-none cursor-pointer"
          />

          <div class="flex justify-between text-xs text-gray-500 mt-2">
            <span>0.0 (Strict)</span>
            <span>0.5 (Balanced)</span>
            <span>1.0 (Creative)</span>
          </div>

          <div class="mt-3 text-sm text-gray-400">
            {{ getTemperatureDescription(formData.temperature) }}
          </div>
        </div>

        <!-- Cursor Assistance Factor -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase mb-3">Cursor AI Assistance Level</h2>

          <div class="flex items-center justify-between mb-2">
            <span class="text-sm text-gray-400">Level</span>
            <span class="text-xl font-bold text-white">{{ formData.cursor_assistance }}%</span>
          </div>

          <input
            v-model.number="formData.cursor_assistance"
            type="range"
            min="0"
            max="100"
            step="5"
            class="w-full h-2 bg-gray-700 rounded appearance-none cursor-pointer"
          />

          <div class="flex justify-between text-xs text-gray-500 mt-2">
            <span>0% (Manual)</span>
            <span>50% (Hybrid)</span>
            <span>100% (Full AI)</span>
          </div>

          <div class="mt-3 p-3 bg-gray-900 border border-gray-700 rounded">
            <p class="text-sm text-gray-300 mb-2">
              <strong>{{ getCursorTitle(formData.cursor_assistance) }}</strong>
            </p>
            <p class="text-xs text-gray-400">
              {{ getCursorDescription(formData.cursor_assistance) }}
            </p>
          </div>
        </div>

      </div>

      <!-- Action Buttons -->
      <div class="mt-6 flex items-center gap-3">
        <button
          @click="saveConfiguration"
          :disabled="isSaving"
          class="flex-1 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
        >
          <span v-if="isSaving" class="animate-spin">⚙️</span>
          <span v-else>💾</span>
          <span>{{ isSaving ? 'Saving...' : 'Save Configuration' }}</span>
        </button>

        <button
          @click="resetToDefaults"
          class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded transition-colors"
        >
          Reset
        </button>
      </div>

      <!-- Info Footer -->
      <div class="mt-6 p-4 bg-blue-900/10 border border-blue-500/30 rounded">
        <p class="text-sm text-blue-300">
          Changes take effect immediately for all AI operations. GLM models are served via ZhipuAI (OpenAI-compatible API).
        </p>
      </div>

      <!-- Discord Notification Testing -->
      <div class="mt-8 border-t border-gray-700 pt-6">
        <h2 class="text-lg font-bold text-white mb-4 flex items-center gap-2">
          <span>🔔</span>
          Discord Notification Testing
        </h2>
        <p class="text-sm text-gray-400 mb-4">Test Discord notifications manually (CEO only)</p>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Test Missing Log -->
          <div class="bg-gray-800 border border-gray-700 rounded p-4">
            <h3 class="text-sm font-bold text-gray-300 mb-2">Missing Time Log Notification</h3>
            <p class="text-xs text-gray-500 mb-3">Send notification for users who didn't log time yesterday</p>
            <button
              @click="testMissingLogNotification"
              :disabled="isTestingMissingLog"
              class="w-full px-4 py-2 bg-orange-600 hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded transition-colors flex items-center justify-center gap-2"
            >
              <span v-if="isTestingMissingLog" class="animate-spin">⏳</span>
              <span v-else>📤</span>
              <span>{{ isTestingMissingLog ? 'Sending...' : 'Test Missing Log' }}</span>
            </button>
            <div v-if="missingLogResult" class="mt-3 p-2 bg-gray-900 rounded text-xs">
              <p class="text-green-400">✅ {{ missingLogResult.message }}</p>
              <p class="text-gray-400">Date: {{ missingLogResult.date }}</p>
              <p class="text-gray-400">Users without logs: {{ missingLogResult.users_without_logs }}</p>
            </div>
          </div>

          <!-- Test Leave Notification -->
          <div class="bg-gray-800 border border-gray-700 rounded p-4">
            <h3 class="text-sm font-bold text-gray-300 mb-2">Leave Notification</h3>
            <p class="text-xs text-gray-500 mb-3">Send notification for users on leave today</p>
            <button
              @click="testLeaveNotification"
              :disabled="isTestingLeave"
              class="w-full px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded transition-colors flex items-center justify-center gap-2"
            >
              <span v-if="isTestingLeave" class="animate-spin">⏳</span>
              <span v-else>📤</span>
              <span>{{ isTestingLeave ? 'Sending...' : 'Test Leave' }}</span>
            </button>
            <div v-if="leaveResult" class="mt-3 p-2 bg-gray-900 rounded text-xs">
              <p class="text-green-400">✅ {{ leaveResult.message }}</p>
              <p class="text-gray-400">Date: {{ leaveResult.date }}</p>
              <p class="text-gray-400">Leaves: {{ leaveResult.leaves_count }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Toast -->
    <div
      v-if="showSuccessToast"
      class="fixed bottom-8 right-8 bg-gray-800 border-2 border-blue-500 text-white px-6 py-4 rounded-lg flex items-center gap-3 z-50 shadow-xl"
    >
      <span>✅</span>
      <div>
        <p class="font-bold">Configuration Saved!</p>
        <p class="text-sm">Changes are now active</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

interface SystemConfig {
  active_model: string
  temperature: number
  cursor_assistance: number
  updated_at: string
}

interface ModelCard {
  id: string
  name: string
  description: string
  tags: string[]
}

const { fetchWithAuth } = useAuth()
const { showError, confirm } = useNotification()

const modelCards: ModelCard[] = [
  {
    id: 'glm-5.1',
    name: 'GLM-5.1',
    description: 'Flagship — top-tier coding & agentic tasks. Best for project plan generation and complex code review.',
    tags: ['Recommended', '744B MoE', '200K ctx']
  },
  {
    id: 'glm-5',
    name: 'GLM-5',
    description: 'Complex systems engineering & long-horizon agentic tasks. Strong reasoning capabilities.',
    tags: ['Agentic', '744B MoE', '200K ctx']
  },
  {
    id: 'glm-4.7',
    name: 'GLM-4.7',
    description: 'Daily coding & agentic workflows. Good balance of quality and speed.',
    tags: ['Balanced', '200K ctx']
  },
  {
    id: 'glm-4.7-flash',
    name: 'GLM-4.7-Flash',
    description: 'High-efficiency, low-cost model. Best for task estimation, quick reviews, and volume work.',
    tags: ['Fast', 'Low Cost', 'Default']
  },
  {
    id: 'glm-5v-turbo',
    name: 'GLM-5V-Turbo',
    description: 'Multimodal — processes screenshots & UI mockups. For vision + code tasks.',
    tags: ['Vision', 'Multimodal']
  },
  {
    id: 'glm-4.6',
    name: 'GLM-4.6',
    description: 'Previous generation. Solid for general tasks with advanced reasoning.',
    tags: ['Stable', '128K ctx']
  }
]

// State
const config = ref<SystemConfig>({
  active_model: 'glm-4.7-flash',
  temperature: 0.4,
  cursor_assistance: 80,
  updated_at: new Date().toISOString()
})

const formData = ref({
  active_model: 'glm-4.7-flash',
  temperature: 0.4,
  cursor_assistance: 80
})

const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')
const showSuccessToast = ref(false)

// Discord testing state
const isTestingMissingLog = ref(false)
const isTestingLeave = ref(false)
const missingLogResult = ref<{ message: string; date: string; users_without_logs: number } | null>(null)
const leaveResult = ref<{ message: string; date: string; leaves_count: number } | null>(null)

// Fetch current configuration
const fetchConfig = async () => {
  try {
    isLoading.value = true
    const response = await fetchWithAuth<{ data: SystemConfig }>('/admin/config')
    config.value = response.data
    formData.value = { ...response.data }
    error.value = ''
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load configuration'
    console.error('Failed to fetch config:', err)
  } finally {
    isLoading.value = false
  }
}

// Save configuration
const saveConfiguration = async () => {
  try {
    isSaving.value = true
    const response = await fetchWithAuth<{ data: SystemConfig }>('/admin/config', {
      method: 'PUT',
      body: formData.value
    })

    config.value = response.data
    showSuccessToast.value = true

    setTimeout(() => {
      showSuccessToast.value = false
    }, 3000)
  } catch (err: any) {
    showError(err.data?.message || err.message || 'Failed to save configuration', 'Save failed')
    console.error('Failed to save config:', err)
  } finally {
    isSaving.value = false
  }
}

// Reset to defaults
const resetToDefaults = async () => {
  const ok = await confirm({
    title: 'Reset to defaults',
    message: 'Reset to default configuration? Model: glm-4.7-flash, Temperature: 0.4, Cursor Assistance: 80%',
    confirmLabel: 'Reset',
    cancelLabel: 'Cancel',
    variant: 'primary'
  })
  if (ok) {
    formData.value = {
      active_model: 'glm-4.7-flash',
      temperature: 0.4,
      cursor_assistance: 80
    }
  }
}

// Test missing log notification
const testMissingLogNotification = async () => {
  try {
    isTestingMissingLog.value = true
    missingLogResult.value = null
    const response = await fetchWithAuth<{ message: string; date: string; users_without_logs: number }>('/admin/discord/test-missing-log', {
      method: 'POST'
    })
    missingLogResult.value = response
  } catch (err: any) {
    showError(err.data?.message || err.message || 'Failed to send notification', 'Discord Error')
  } finally {
    isTestingMissingLog.value = false
  }
}

// Test leave notification
const testLeaveNotification = async () => {
  try {
    isTestingLeave.value = true
    leaveResult.value = null
    const response = await fetchWithAuth<{ message: string; date: string; leaves_count: number }>('/admin/discord/test-leave', {
      method: 'POST'
    })
    leaveResult.value = response
  } catch (err: any) {
    showError(err.data?.message || err.message || 'Failed to send notification', 'Discord Error')
  } finally {
    isTestingLeave.value = false
  }
}

// Model badge helper
const getModelBadge = (modelId: string): { label: string; class: string } => {
  if (modelId.startsWith('glm-5.1')) return { label: 'Flagship', class: 'bg-yellow-700 text-yellow-100' }
  if (modelId.startsWith('glm-5')) return { label: 'Agentic', class: 'bg-purple-700 text-purple-100' }
  if (modelId.includes('flash')) return { label: 'Fast', class: 'bg-green-700 text-green-100' }
  if (modelId.includes('5v')) return { label: 'Vision', class: 'bg-orange-700 text-orange-100' }
  return { label: 'Standard', class: 'bg-gray-700 text-gray-300' }
}

// Tag color helper
const getTagClass = (tag: string): string => {
  if (tag === 'Recommended') return 'bg-yellow-900/50 text-yellow-300'
  if (tag === 'Default') return 'bg-blue-900/50 text-blue-300'
  if (tag === 'Fast' || tag === 'Low Cost') return 'bg-green-900/50 text-green-300'
  if (tag === 'Agentic') return 'bg-purple-900/50 text-purple-300'
  if (tag === 'Vision' || tag === 'Multimodal') return 'bg-orange-900/50 text-orange-300'
  return 'bg-gray-700 text-gray-400'
}

// Temperature helpers
const getTemperatureDescription = (temp: number): string => {
  if (temp <= 0.2) return '🎯 Maximum Precision — Highly consistent estimates'
  if (temp <= 0.4) return '⚖️ Balanced — Stable with slight variation'
  if (temp <= 0.6) return '🎨 Creative — More diverse responses'
  if (temp <= 0.8) return '🌈 Highly Creative — Varied outputs'
  return '🚀 Experimental — Maximum creativity'
}

// Cursor Assistance Helpers
const getCursorTitle = (level: number): string => {
  if (level <= 20) return 'Traditional Development'
  if (level <= 50) return 'Hybrid AI Workflow'
  if (level <= 80) return 'AI-Powered Development'
  return 'Ultra AI-First Workflow'
}

const getCursorDescription = (level: number): string => {
  if (level <= 20) return 'Minimal AI assistance. Developers code mostly manually. Expect traditional development speeds.'
  if (level <= 50) return 'Moderate AI usage for suggestions and debugging. Balanced approach with moderate time savings.'
  if (level <= 80) return 'Heavy AI reliance for boilerplate, refactoring, and debugging. Significant productivity boost.'
  return 'AI-first workflow with near-complete AI assistance. Maximum velocity with aggressive estimates.'
}

// Date formatter
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Load data on mount
onMounted(async () => {
  await fetchConfig()
})
</script>
