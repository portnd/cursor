<template>
  <div class="w-full max-w-md mx-auto">
    <!-- Card Container -->
    <div class="bg-gray-800/80 backdrop-blur-xl border border-gray-700 rounded-2xl shadow-2xl p-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h2 class="text-3xl font-bold text-white mb-2">
          Create Account
        </h2>
        <p class="text-gray-400">
          Join The Sentinel and start your journey
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Email Field -->
        <div>
          <label for="email" class="block text-sm font-medium text-gray-300 mb-2">
            Email Address
          </label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            required
            placeholder="you@example.com"
            class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200 outline-none"
            :class="{ 'border-red-500': errors.email }"
            :disabled="isLoading"
          />
          <p v-if="errors.email" class="mt-1 text-sm text-red-400">
            {{ errors.email }}
          </p>
        </div>

        <!-- Password Field -->
        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-2">
            Password
          </label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200 outline-none"
            :class="{ 'border-red-500': errors.password }"
            :disabled="isLoading"
          />
          <p v-if="errors.password" class="mt-1 text-sm text-red-400">
            {{ errors.password }}
          </p>
          <p class="mt-1 text-xs text-gray-500">
            Must be at least 8 characters
          </p>
        </div>

        <!-- Confirm Password Field -->
        <div>
          <label for="confirm-password" class="block text-sm font-medium text-gray-300 mb-2">
            Confirm Password
          </label>
          <input
            id="confirm-password"
            v-model="formData.confirmPassword"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-3 bg-gray-900/50 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200 outline-none"
            :class="{ 'border-red-500': errors.confirmPassword }"
            :disabled="isLoading"
          />
          <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-400">
            {{ errors.confirmPassword }}
          </p>
        </div>

        <!-- Error Message -->
        <div v-if="authStore.error" class="bg-red-900/20 border border-red-500/50 text-red-400 px-4 py-3 rounded-lg">
          <p class="text-sm">{{ authStore.error }}</p>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full py-3 px-4 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 hover:from-purple-200 dark:from-purple-700 hover:to-pink-200 dark:to-pink-700 text-gray-900 dark:text-white font-semibold rounded-lg shadow-lg hover:shadow-purple-500/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          <span v-if="!isLoading">Create Account</span>
          <span v-else class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          </span>
        </button>
      </form>

      <!-- Footer -->
      <div class="mt-6 text-center">
        <p class="text-sm text-gray-400">
          Already have an account?
          <NuxtLink
            to="/login"
            class="font-semibold text-purple-400 hover:text-purple-300 transition-colors"
          >
            Sign in
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
  confirmPassword: '',
})

const errors = ref({
  email: '',
  password: '',
  confirmPassword: '',
})

// Loading state
const isLoading = computed(() => authStore.isLoading)

/**
 * Validate form inputs
 */
const validateForm = (): boolean => {
  errors.value = { email: '', password: '', confirmPassword: '' }
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

  // Confirm password validation
  if (!formData.value.confirmPassword) {
    errors.value.confirmPassword = 'Please confirm your password'
    isValid = false
  } else if (formData.value.password !== formData.value.confirmPassword) {
    errors.value.confirmPassword = 'Passwords do not match'
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

  // Call register action
  const result = await authStore.register(
    formData.value.email,
    formData.value.password,
    formData.value.confirmPassword
  )

  // Redirect on success
  if (result?.success) {
    await router.push('/')
  }
}
</script>
