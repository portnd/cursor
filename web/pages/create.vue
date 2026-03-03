<template>
  <div class="min-h-screen p-8">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-white mb-2">
        Create New Mission 🚀
      </h1>
      <p class="text-gray-400">Define your task and create the mission</p>
    </div>

    <!-- Creation Form Card -->
    <div class="max-w-3xl mx-auto">
      <div class="bg-gray-800/50 border border-gray-700 rounded-xl p-8 backdrop-blur">
        <!-- Success Message -->
        <div 
          v-if="showSuccess" 
          class="mb-6 p-4 bg-green-500/10 border border-green-500/50 rounded-lg text-green-400 flex items-center gap-3 animate-fade-in"
        >
          <svg class="w-6 h-6 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p class="font-semibold">Task Created Successfully!</p>
            <p class="text-sm text-green-300">Redirecting to dashboard...</p>
          </div>
        </div>

        <!-- Error Message -->
        <div 
          v-if="errorMessage" 
          class="mb-6 p-4 bg-red-500/10 border border-red-500/50 rounded-lg text-red-400 flex items-center gap-3"
        >
          <svg class="w-6 h-6 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p class="font-semibold">Failed to Create Task</p>
            <p class="text-sm text-red-300">{{ errorMessage }}</p>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Task Title -->
          <div>
            <label for="title" class="block text-sm font-medium text-gray-300 mb-2">
              Task Title <span class="text-red-400">*</span>
            </label>
            <input
              id="title"
              v-model="formData.title"
              type="text"
              required
              placeholder="e.g., Implement user authentication system"
              class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
              :disabled="isSubmitting"
            />
          </div>

          <!-- Description -->
          <div>
            <label for="description" class="block text-sm font-medium text-gray-300 mb-2">
              Description <span class="text-red-400">*</span>
            </label>
            <textarea
              id="description"
              v-model="formData.description"
              required
              rows="8"
              placeholder="Provide detailed requirements, technical specifications, and acceptance criteria..."
              class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition resize-none"
              :disabled="isSubmitting"
            ></textarea>
            <p class="mt-2 text-sm text-gray-500">
              💡 Tip: Add clear requirements and acceptance criteria
            </p>
          </div>

          <!-- Deadline -->
          <div>
            <label for="deadline" class="block text-sm font-medium text-gray-300 mb-2">
              ⏰ Deadline <span class="text-gray-500 text-xs">(Optional)</span>
            </label>
            <input
              id="deadline"
              v-model="formData.deadline"
              type="datetime-local"
              class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition"
              :disabled="isSubmitting"
              :min="minDeadline"
            />
            <p class="mt-2 text-sm text-red-400 flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Set a deadline to enforce urgency and track performance
            </p>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center gap-4 pt-4">
            <button
              type="submit"
              :disabled="isSubmitting || !formData.title || !formData.description"
              class="flex-1 py-4 px-6 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-lg shadow-lg hover:shadow-purple-500/50 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3"
            >
              <svg 
                v-if="isSubmitting" 
                class="animate-spin h-5 w-5" 
                xmlns="http://www.w3.org/2000/svg" 
                fill="none" 
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg 
                v-else 
                class="w-5 h-5" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <span v-if="isSubmitting">Analyzing with AI...</span>
              <span v-else>Initialize Task & Run AI Analysis</span>
            </button>

            <NuxtLink
              to="/dashboard"
              class="px-6 py-4 bg-gray-700/50 hover:bg-gray-700 border border-gray-600 hover:border-gray-500 text-gray-300 hover:text-white font-medium rounded-lg transition-all"
            >
              Cancel
            </NuxtLink>
          </div>
        </form>

        <!-- Additional Info -->
        <div class="mt-8 pt-6 border-t border-gray-700">
          <div class="flex items-center gap-2 text-sm text-gray-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            <span>Your data is secure and encrypted</span>
          </div>
        </div>
      </div>

      <!-- Pro Tips Card -->
      <div class="mt-6 bg-gradient-to-br from-purple-900/20 to-pink-900/20 border border-purple-500/30 rounded-xl p-6">
        <h3 class="text-lg font-semibold text-white mb-3">
          ✨ Pro Tips for Better AI Estimation:
        </h3>
        <ul class="space-y-2 text-sm text-gray-300">
          <li class="flex items-start gap-2">
            <span class="text-purple-400 shrink-0">•</span>
            <span>Include specific technical requirements (API endpoints, database schema, etc.)</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-purple-400 shrink-0">•</span>
            <span>Mention any third-party integrations or external dependencies</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-purple-400 shrink-0">•</span>
            <span>List acceptance criteria and expected deliverables</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-purple-400 shrink-0">•</span>
            <span>Note any performance or security requirements</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const { fetchWithAuth } = useAuth()
const router = useRouter()

// Form state
const formData = ref({
  title: '',
  description: '',
  deadline: ''
})

const isSubmitting = ref(false)
const errorMessage = ref('')
const showSuccess = ref(false)

// Minimum deadline is now
const minDeadline = computed(() => {
  const now = new Date()
  return now.toISOString().slice(0, 16) // Format: YYYY-MM-DDTHH:MM
})

// Handle form submission
const handleSubmit = async () => {
  isSubmitting.value = true
  errorMessage.value = ''
  showSuccess.value = false

  try {
    // Prepare request body
    const requestBody: any = {
      title: formData.value.title,
      description: formData.value.description
    }

    // Add deadline if provided (convert to ISO8601/RFC3339 format)
    if (formData.value.deadline) {
      requestBody.due_date = new Date(formData.value.deadline).toISOString()
    }

    const response = await fetchWithAuth('/sentinel/tasks', {
      method: 'POST',
      body: requestBody
    })

    // Show success message
    showSuccess.value = true

    // Wait a moment for user to see success message
    setTimeout(() => {
      router.push('/dashboard')
    }, 2000)

  } catch (error: any) {
    console.error('Failed to create task:', error)
    errorMessage.value = error.data?.message || error.message || 'Failed to create task. Please try again.'
    
    // Scroll to top to show error
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } finally {
    isSubmitting.value = false
  }
}

// Reset form when component unmounts
onUnmounted(() => {
  formData.value = { title: '', description: '', deadline: '' }
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}
</style>
