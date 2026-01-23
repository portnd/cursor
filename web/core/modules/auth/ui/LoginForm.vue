<template>
  <div class="w-full max-w-md mx-auto">
    <!-- Card Container -->
    <div class="bg-white rounded-2xl shadow-xl p-8 border border-slate-200">
      <!-- Header -->
      <div class="text-center mb-8">
        <h2 class="text-3xl font-bold text-slate-900 mb-2">
          Welcome Back
        </h2>
        <p class="text-slate-600">
          Sign in to your account to continue
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Email Field -->
        <div>
          <label for="email" class="block text-sm font-medium text-slate-700 mb-2">
            Email Address
          </label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            required
            placeholder="you@example.com"
            class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200 outline-none"
            :class="{ 'border-red-500': errors.email }"
            :disabled="isLoading"
          />
          <p v-if="errors.email" class="mt-1 text-sm text-red-600">
            {{ errors.email }}
          </p>
        </div>

        <!-- Password Field -->
        <div>
          <label for="password" class="block text-sm font-medium text-slate-700 mb-2">
            Password
          </label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200 outline-none"
            :class="{ 'border-red-500': errors.password }"
            :disabled="isLoading"
          />
          <p v-if="errors.password" class="mt-1 text-sm text-red-600">
            {{ errors.password }}
          </p>
        </div>

        <!-- Error Message -->
        <div v-if="authStore.error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          <p class="text-sm">{{ authStore.error }}</p>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full bg-gradient-to-r from-purple-600 to-indigo-600 text-white font-semibold py-3 px-6 rounded-lg hover:from-purple-700 hover:to-indigo-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="!isLoading">Sign In</span>
          <span v-else class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </span>
        </button>
      </form>

      <!-- Footer -->
      <div class="mt-6 text-center">
        <p class="text-sm text-slate-600">
          Don't have an account?
          <NuxtLink
            to="/register"
            class="font-semibold text-purple-600 hover:text-purple-700 transition-colors"
          >
            Sign up
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '../store/auth-store'

const authStore = useAuthStore()
const router = useRouter()

// Form state
const formData = ref({
  email: '',
  password: '',
})

const errors = ref({
  email: '',
  password: '',
})

// Loading state
const isLoading = computed(() => authStore.isLoading)

/**
 * Validate form inputs
 */
const validateForm = (): boolean => {
  errors.value = { email: '', password: '' }
  let isValid = true

  // Email validation
  if (!formData.value.email) {
    errors.value.email = 'Email is required'
    isValid = false
  } else if (!/\S+@\S+\.\S+/.test(formData.value.email)) {
    errors.value.email = 'Email is invalid'
    isValid = false
  }

  // Password validation
  if (!formData.value.password) {
    errors.value.password = 'Password is required'
    isValid = false
  } else if (formData.value.password.length < 8) {
    errors.value.password = 'Password must be at least 8 characters'
    isValid = false
  }

  return isValid
}

/**
 * Handle form submission
 */
const handleSubmit = async () => {
  // Clear previous errors
  authStore.clearError()

  // Validate form
  if (!validateForm()) {
    return
  }

  // Call login action
  const result = await authStore.login(
    formData.value.email,
    formData.value.password
  )

  // Redirect on success
  if (result?.success) {
    await router.push('/')
  }
}
</script>
