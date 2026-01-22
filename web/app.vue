<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
useHead({
  htmlAttrs: {
    lang: 'en',
    class: 'antialiased'
  }
})
</script>
