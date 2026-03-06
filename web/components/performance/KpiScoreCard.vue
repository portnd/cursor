<template>
  <div
    class="rounded-lg border p-4 transition-colors"
    :class="statusClass"
  >
    <div class="text-xs uppercase tracking-wide text-gray-400 mb-1">
      {{ label }}
    </div>
    <div class="text-2xl font-bold text-white">
      {{ displayValue }}
    </div>
    <div v-if="sublabel" class="text-xs text-gray-500 mt-1">
      {{ sublabel }}
    </div>
    <div v-if="trend" class="mt-2 flex items-center gap-1 text-xs">
      <span v-if="trend === 'up'" class="text-emerald-400">↑ {{ trendLabel }}</span>
      <span v-else-if="trend === 'down'" class="text-red-400">↓ {{ trendLabel }}</span>
      <span v-else class="text-gray-500">→ {{ trendLabel }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    value: number | string
    format?: 'pct' | 'number' | 'raw'
    sublabel?: string
    trend?: 'up' | 'down' | 'stable'
    trendLabel?: string
    status?: 'good' | 'warn' | 'bad' | 'neutral'
  }>(),
  { format: 'raw', trendLabel: '', status: 'neutral' }
)

const displayValue = computed(() => {
  if (props.format === 'pct') {
    return `${Number(props.value).toFixed(1)}%`
  }
  if (props.format === 'number') {
    return Number(props.value).toFixed(1)
  }
  return props.value
})

const statusClass = computed(() => {
  const base = 'bg-gray-800 border-gray-700'
  if (props.status === 'good') return `${base} border-emerald-500/50`
  if (props.status === 'warn') return `${base} border-amber-500/50`
  if (props.status === 'bad') return `${base} border-red-500/50`
  return base
})
</script>
