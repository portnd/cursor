<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-white">AI CONFIGURATION</h1>
      <p class="text-sm text-gray-400 mt-1">System-wide AI behavior settings</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-20">
      <div class="text-gray-400">Loading configuration...</div>
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
            </div>
          </div>
          <div class="text-right">
            <p class="text-xs text-gray-500">Last Updated</p>
            <p class="text-sm text-gray-300">{{ formatDateTime(config.updated_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Configuration Form -->
      <div class="space-y-6">
        
        <!-- 1. Model Selector -->
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase mb-3">AI Model Selection</h2>

          <label class="block text-sm text-gray-400 mb-2">Active Model</label>
          <select
            v-model="formData.active_model"
            class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white focus:border-blue-500 outline-none"
          >
            <option v-for="model in availableModels" :key="model" :value="model">
              {{ model }}
              <span v-if="model === config.active_model"> (Current)</span>
            </option>
          </select>
          <p class="text-xs text-gray-500 mt-2">
            Recommended: gemini-2.5-flash-lite for balanced performance
          </p>
        </div>

        <!-- 2. Temperature Control -->
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

        <!-- 3. Cursor Assistance Factor -->
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
          class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
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
          Changes take effect immediately for all AI operations.
        </p>
      </div>
    </div>

    <!-- Success Toast -->
    <div
      v-if="showSuccessToast"
      class="fixed bottom-8 right-8 bg-green-600 text-white px-6 py-4 rounded flex items-center gap-3 z-50"
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

const { fetchWithAuth } = useAuth()

// State
const config = ref<SystemConfig>({
  active_model: 'gemini-2.5-flash-lite',
  temperature: 0.4,
  cursor_assistance: 80,
  updated_at: new Date().toISOString()
})

const formData = ref({
  active_model: 'gemini-2.5-flash-lite',
  temperature: 0.4,
  cursor_assistance: 80
})

const availableModels = ref<string[]>([])
const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')
const showSuccessToast = ref(false)

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

// Fetch available models
const fetchModels = async () => {
  try {
    const response = await fetchWithAuth<{ data: string[] }>('/admin/models')
    availableModels.value = response.data || []
  } catch (err: any) {
    console.error('Failed to fetch models:', err)
    // Fallback to default models
    availableModels.value = [
      'gemini-1.5-flash',
      'gemini-1.5-pro',
      'gemini-2.0-flash-exp',
      'gemini-2.5-flash-lite',
      'gemini-exp-1206'
    ]
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
    
    // Hide toast after 3 seconds
    setTimeout(() => {
      showSuccessToast.value = false
    }, 3000)
  } catch (err: any) {
    alert(`Failed to save configuration: ${err.data?.message || err.message}`)
    console.error('Failed to save config:', err)
  } finally {
    isSaving.value = false
  }
}

// Reset to defaults
const resetToDefaults = () => {
  if (confirm('Reset to default configuration?\n\nModel: gemini-2.5-flash-lite\nTemperature: 0.4\nCursor Assistance: 80%')) {
    formData.value = {
      active_model: 'gemini-2.5-flash-lite',
      temperature: 0.4,
      cursor_assistance: 80
    }
  }
}

// Temperature helpers
const getTemperatureDescription = (temp: number): string => {
  if (temp <= 0.2) return '🎯 Maximum Precision - Highly consistent estimates'
  if (temp <= 0.4) return '⚖️ Balanced - Stable with slight variation'
  if (temp <= 0.6) return '🎨 Creative - More diverse responses'
  if (temp <= 0.8) return '🌈 Highly Creative - Varied outputs'
  return '🚀 Experimental - Maximum creativity'
}

// Cursor Assistance Helpers
const getCursorTitle = (level: number): string => {
  if (level <= 20) return 'Traditional Development'
  if (level <= 50) return 'Hybrid AI Workflow'
  if (level <= 80) return 'AI-Powered Development'
  return 'Ultra AI-First Workflow'
}

const getCursorDescription = (level: number): string => {
  if (level <= 20) {
    return 'Minimal AI assistance. Developers code mostly manually. Expect traditional development speeds.'
  }
  if (level <= 50) {
    return 'Moderate AI usage for suggestions and debugging. Balanced approach with moderate time savings.'
  }
  if (level <= 80) {
    return 'Heavy AI reliance for boilerplate, refactoring, and debugging. Significant productivity boost.'
  }
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
  await Promise.all([fetchConfig(), fetchModels()])
})
</script>

