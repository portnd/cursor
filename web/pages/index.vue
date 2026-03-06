<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 flex items-center justify-center" style="background-color: #111827;">
    <div class="text-center">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-gray-400">Redirecting...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated } = useAuth()

onMounted(() => {
  if (isAuthenticated.value) {
    navigateTo('/dashboard')
  } else {
    navigateTo('/login')
  }
})
</script>
