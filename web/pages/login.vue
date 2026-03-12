<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
    <!-- Background Decorations -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl animate-pulse"></div>
      <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-pink-500/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
    </div>

    <!-- Login Card -->
    <div class="relative w-full max-w-md">
      <!-- Logo & Title -->
      <div class="text-center mb-8">
        <h1 class="text-5xl font-bold mb-2">
          <span class="bg-gradient-to-r from-purple-400 via-pink-500 to-purple-600 bg-clip-text text-transparent">
            🛡️ The Sentinel
          </span>
        </h1>
        <p class="text-gray-400">AI-Powered Task Management System</p>
      </div>

      <!-- Login Form Card -->
      <div class="bg-gray-800/50 backdrop-blur-xl border border-gray-700 rounded-2xl shadow-2xl p-8">
        <h2 class="text-2xl font-bold text-white mb-6">Welcome Back</h2>

        <!-- Session expired notice (shown after redirect from 401) -->
        <div v-if="sessionExpired" class="mb-4 p-4 bg-amber-500/10 border border-amber-500/50 rounded-lg text-amber-400 text-sm">
          Session expired or invalid. Please sign in again.
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="mb-4 p-4 bg-red-500/10 border border-red-500/50 rounded-lg text-red-400 text-sm">
          {{ errorMessage }}
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- Email Input -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-300 mb-2">
              Email Address
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              placeholder="mail@komgrip.com"
              class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
            />
          </div>

          <!-- Password Input -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-300 mb-2">
              Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              placeholder="••••••••"
              class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
            />
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-lg shadow-lg hover:shadow-purple-500/50 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="isLoading" class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-if="!isLoading">Sign In</span>
            <span v-else>Signing In...</span>
          </button>
        </form>
      </div>

      <!-- Footer -->
      <p class="text-center text-gray-500 text-sm mt-6">
        Powered by AI • Built with Nuxt 3 & Tailwind CSS
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

definePageMeta({
  layout: false // Use custom layout for login page
})

const { login } = useAuth()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

const sessionExpired = ref(false)
onMounted(() => {
  if (route.query.session === 'expired') {
    sessionExpired.value = true
    if (import.meta.client) router.replace({ path: '/login', query: {} })
  }
})

const handleLogin = async () => {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const result = await login(email.value, password.value)
    
    if (result.success) {
      // Sync token + user into auth store so role-based UI (Edit button, etc.) works
      const tokenCookie = useCookie('token')
      if (result.token) authStore.token = result.token
      else if (tokenCookie.value) authStore.token = tokenCookie.value
      if (result.user) {
        authStore.user = result.user
        authStore.isAuthenticated = true
        if (import.meta.client) {
          localStorage.setItem('user', JSON.stringify(result.user))
        }
      }
      await router.push('/dashboard')
    } else {
      errorMessage.value = result.error || 'Login failed. Please check your credentials.'
    }
  } catch (error: any) {
    errorMessage.value = error.message || 'An unexpected error occurred'
  } finally {
    isLoading.value = false
  }
}
</script>
