<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900/30 to-gray-900 flex items-center justify-center p-4">
    <!-- Background Decorations -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl animate-pulse"></div>
      <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-pink-500/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
    </div>

    <!-- Register Form Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Branding -->
      <div class="text-center mb-8">
        <h1 class="text-5xl font-bold mb-2">
          <span class="bg-gradient-to-r from-purple-400 via-pink-500 to-purple-600 bg-clip-text text-transparent">
            🛡️ The Sentinel
          </span>
        </h1>
        <p class="text-gray-400">Create your account — AI-Powered Task Management</p>
      </div>

      <!-- Register Form Component -->
      <RegisterForm />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'
import RegisterForm from '~/core/modules/auth/ui/RegisterForm.vue'

// Redirect if already logged in
const authStore = useAuthStore()
const router = useRouter()

onMounted(() => {
  if (authStore.isLoggedIn) {
    router.push('/')
  }
})

definePageMeta({
  layout: false,
})

useHead({
  title: 'Register - The Sentinel',
  meta: [
    {
      name: 'description',
      content: 'Create your Sentinel account',
    },
  ],
})
</script>
