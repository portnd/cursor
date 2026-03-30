<script setup lang="ts">
/**
 * Feature roadmap lives on each project: /projects/:id?tab=feature-roadmap
 */
definePageMeta({ middleware: 'auth' })
await navigateTo('/projects', { replace: true })
</script>

<template>
  <div class="min-h-screen bg-gray-900" />
</template>
