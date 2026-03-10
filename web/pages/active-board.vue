<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <main class="max-w-screen-xl mx-auto px-6 py-8">
      <GlobalKanbanBoard />
    </main>
  </div>
</template>

<script setup lang="ts">
import GlobalKanbanBoard from '~/components/tasks/GlobalKanbanBoard.vue'

definePageMeta({ middleware: 'auth' })
useHead({ title: 'Global Active Board — The Sentinel' })
</script>
