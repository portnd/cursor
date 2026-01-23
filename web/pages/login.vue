<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-600 via-indigo-600 to-blue-700 p-4">
    <!-- Background Pattern -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-1/2 -left-1/2 w-full h-full bg-purple-500/10 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-1/2 -right-1/2 w-full h-full bg-indigo-500/10 rounded-full blur-3xl"></div>
    </div>

    <!-- Login Form Container -->
    <div class="relative z-10 w-full">
      <!-- Logo/Branding -->
      <div class="text-center mb-8">
        <h1 class="text-5xl font-bold text-white mb-2">
          🛡️ KOMGRIP
        </h1>
        <p class="text-white/80 text-lg">
          God-Tier Starter Kit
        </p>
      </div>

      <!-- Login Form Component -->
      <LoginForm />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'
import LoginForm from '~/core/modules/auth/ui/LoginForm.vue'

// Redirect if already logged in
const authStore = useAuthStore()
const router = useRouter()

onMounted(() => {
  if (authStore.isLoggedIn) {
    router.push('/')
  }
})

// Meta tags
definePageMeta({
  layout: false, // No layout for auth pages
})

useHead({
  title: 'Login - Komgrip Starter Kit',
  meta: [
    {
      name: 'description',
      content: 'Sign in to your Komgrip account',
    },
  ],
})
</script>
