<template>
  <div class="min-h-screen bg-gray-900 text-gray-100" style="background-color: #111827;">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiAppModal />
  </div>
</template>

<script setup lang="ts">
useHead({
  htmlAttrs: {
    lang: 'en',
    class: 'antialiased bg-gray-900'
  },
  bodyAttrs: {
    class: 'bg-gray-900 text-gray-100'
  }
})
</script>
