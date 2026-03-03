<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center">
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
