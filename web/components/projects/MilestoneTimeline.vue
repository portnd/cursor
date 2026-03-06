<template>
  <div class="milestone-timeline">
    <div v-if="!milestones.length" class="text-center py-8 text-gray-500 text-sm">
      No milestones defined. Add milestones to track key deliverables.
    </div>
    <div v-else class="relative">
      <!-- Timeline track -->
      <div class="absolute left-0 right-0 top-8 h-0.5 bg-gray-700"></div>

      <!-- Milestone markers -->
      <div class="flex items-start justify-between relative overflow-x-auto pb-4 min-w-0">
        <div
          v-for="(m, idx) in sortedMilestones"
          :key="m.id"
          class="flex flex-col items-center gap-2 relative px-3 group cursor-pointer"
          style="min-width: 100px"
          @click="$emit('milestone-click', m)"
        >
          <!-- Diamond marker -->
          <div
            class="w-5 h-5 rotate-45 border-2 transition-all mt-[22px] z-10"
            :class="markerClass(m)"
          ></div>

          <!-- Label below -->
          <div class="mt-2 text-center">
            <p class="text-xs font-semibold text-gray-300 leading-tight line-clamp-2 max-w-[90px]">{{ m.title }}</p>
            <p class="text-[10px] mt-1" :class="dueDateClass(m)">
              {{ formatDate(m.due_date) }}
            </p>
            <span
              class="text-[10px] px-1.5 py-0.5 rounded-full font-medium mt-1 inline-block"
              :class="statusBadgeClass(m.status)"
            >
              {{ m.status }}
            </span>
          </div>

          <!-- Tooltip on hover -->
          <div class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 hidden group-hover:block z-20 w-48 bg-gray-900 border border-gray-700 rounded-lg p-3 shadow-xl text-left">
            <p class="text-sm font-semibold text-white mb-1">{{ m.title }}</p>
            <p v-if="m.description" class="text-xs text-gray-400 mb-2">{{ m.description }}</p>
            <p class="text-xs text-gray-500">Due: {{ formatDate(m.due_date) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Milestone Button -->
    <div class="flex justify-end mt-2">
      <button @click="$emit('add-milestone')" class="btn-ghost text-xs">
        + Add Milestone
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Milestone } from '~/core/modules/projects/infrastructure/projects-api'

const props = defineProps<{
  milestones: Milestone[]
}>()

const emit = defineEmits<{
  (e: 'milestone-click', m: Milestone): void
  (e: 'add-milestone'): void
}>()

const sortedMilestones = computed(() =>
  [...props.milestones].sort((a, b) => {
    if (!a.due_date) return 1
    if (!b.due_date) return -1
    return new Date(a.due_date).getTime() - new Date(b.due_date).getTime()
  })
)

function formatDate(d: string | null) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function markerClass(m: Milestone) {
  if (m.status === 'REACHED') return 'bg-green-500 border-green-400'
  if (m.status === 'MISSED') return 'bg-red-500 border-red-400'
  const isOverdue = m.due_date && new Date(m.due_date) < new Date()
  if (isOverdue) return 'bg-red-500/20 border-red-400'
  return 'bg-gray-800 border-purple-400 group-hover:bg-purple-500'
}

function dueDateClass(m: Milestone) {
  if (m.status === 'REACHED') return 'text-green-400'
  if (m.status === 'MISSED') return 'text-red-400'
  if (m.due_date && new Date(m.due_date) < new Date()) return 'text-red-400'
  return 'text-gray-500'
}

function statusBadgeClass(status: string) {
  if (status === 'REACHED') return 'bg-green-500/20 text-green-400'
  if (status === 'MISSED') return 'bg-red-500/20 text-red-400'
  return 'bg-gray-700 text-gray-400'
}
</script>

<style scoped>
.btn-ghost {
  @apply px-3 py-1.5 text-gray-400 hover:text-gray-200 hover:bg-gray-700 rounded-lg transition-colors;
}
</style>
