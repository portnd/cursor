<template>
  <div class="flex flex-wrap items-center gap-2" :class="dense ? '' : 'mb-2'">
    <span
      class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-semibold max-w-[200px] truncate"
      :style="pillStyle"
    >{{ task.project_name || 'Project' }}</span>
    <span v-if="task.sprint_name" class="text-[10px] font-medium text-purple-400/90 truncate max-w-[180px]">{{ task.sprint_name }}</span>
    <span v-if="task.task_type" class="text-[10px] font-bold uppercase px-1.5 py-0.5 rounded border border-gray-600 text-gray-400">{{ task.task_type }}</span>
    <span v-if="task.priority" class="text-[10px] border border-gray-600 rounded px-1 text-gray-400">{{ task.priority }}</span>
    <span class="text-[10px] font-semibold px-2 py-0.5 rounded bg-gray-700/80 text-gray-300">{{ statusLabel }}</span>
    <span v-if="detailed && task.code" class="text-[10px] font-mono text-gray-600">{{ task.code }}</span>
  </div>
  <h3 v-if="detailed" class="text-lg font-bold text-white">{{ task.title }}</h3>
</template>

<script setup lang="ts">
interface BoardTask {
  title: string
  code?: string
  project_name?: string
  project_color?: string
  sprint_name?: string
  task_type?: string
  priority?: string
}

const props = withDefaults(
  defineProps<{
    task: BoardTask
    statusLabel: string
    detailed?: boolean
    /** Tighter spacing (e.g. table cells) */
    dense?: boolean
  }>(),
  { detailed: false, dense: false }
)

const pillStyle = computed(() => {
  const c = props.task.project_color || '#6366f1'
  return {
    borderColor: c,
    color: c,
    backgroundColor: `${c}22`,
  }
})
</script>
