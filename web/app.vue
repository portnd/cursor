<template>
  <div class="min-h-screen">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiAppModal />
  </div>
</template>

<script setup lang="ts">
const { initTheme } = useTheme()

useHead({
  htmlAttrs: {
    lang: 'en',
    class: 'antialiased'
  },
})

// Apply the theme class to <html> on client mount
// This runs after SSR hydration so the correct theme is in place
if (import.meta.client) {
  initTheme()
}
</script>
