<template>
  <div class="rounded-lg border border-gray-700 bg-gray-800 p-4">
    <div class="text-xs uppercase tracking-wide text-gray-400 mb-3">SPACE Dimensions</div>
    <div class="relative h-48 flex items-center justify-center">
      <svg viewBox="0 0 200 200" class="w-full max-w-[240px] h-48" xmlns="http://www.w3.org/2000/svg">
        <!-- Pentagon outline -->
        <polygon
          fill="none"
          stroke="currentColor"
          stroke-width="1"
          class="text-gray-600"
          :points="pentagonPoints"
        />
        <!-- Grid lines from center to each vertex -->
        <line
          v-for="(_, i) in 5"
          :key="i"
          :x1="100"
          :y1="100"
          :x2="100 + 80 * Math.cos(angle(i))"
          :y2="100 - 80 * Math.sin(angle(i))"
          stroke="currentColor"
          stroke-width="0.5"
          class="text-gray-600"
        />
        <!-- Data fill (scaled 0-100 to radius 0-80) -->
        <polygon
          v-if="dataPoints"
          :points="dataPoints"
          fill="rgba(139, 92, 246, 0.25)"
          stroke="rgb(139, 92, 246)"
          stroke-width="1.5"
        />
        <!-- Vertex labels -->
        <text
          v-for="(label, i) in labels"
          :key="'l' + i"
          :x="100 + 95 * Math.cos(angle(i))"
          :y="100 - 95 * Math.sin(angle(i))"
          text-anchor="middle"
          dominant-baseline="middle"
          class="fill-gray-400 text-[10px]"
        >
          {{ label }}
        </text>
      </svg>
    </div>
    <div class="mt-2 grid grid-cols-2 gap-2 text-xs text-gray-500">
      <span>S: {{ values.satisfaction.toFixed(0) }}</span>
      <span>P: {{ values.performance.toFixed(0) }}</span>
      <span>A: {{ values.activity.toFixed(0) }}</span>
      <span>C: {{ values.collaboration.toFixed(0) }}</span>
      <span class="col-span-2">E: {{ values.efficiency.toFixed(0) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  values: {
    satisfaction: number
    performance: number
    activity: number
    collaboration: number
    efficiency: number
  }
}>()

const labels = ['S', 'P', 'A', 'C', 'E']

function angle(i: number) {
  return (Math.PI * 2 * i) / 5 - Math.PI / 2
}

const pentagonPoints = computed(() => {
  return Array.from({ length: 5 }, (_, i) => {
    const a = angle(i)
    return `${100 + 80 * Math.cos(a)},${100 - 80 * Math.sin(a)}`
  }).join(' ')
})

const dataPoints = computed(() => {
  const v = props.values
  const vals = [v.satisfaction, v.performance, v.activity, v.collaboration, v.efficiency]
  const pts = vals.map((val, i) => {
    const r = Math.min(100, Math.max(0, val)) / 100 * 80
    const a = angle(i)
    return `${100 + r * Math.cos(a)},${100 - r * Math.sin(a)}`
  })
  return pts.join(' ')
})
</script>
