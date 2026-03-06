<script setup lang="ts">
import { toRaw } from "vue"
import type { Options as HighchartsOptions, Chart } from "highcharts"

const props = defineProps<{
	options: HighchartsOptions
	constructorType?: "chart" | "stockChart" | "mapChart" | "ganttChart"
}>()

// Use the Highcharts instance provided by plugins/highcharts.ts
// This avoids repeated dynamic imports that break in esbuild pre-bundled environments
const { $highcharts } = useNuxtApp() as any

const container = ref<HTMLElement | null>(null)
let chart: Chart | null = null

const initChart = () => {
	if (!container.value) {
		console.warn("[Highchart] container.value is null, skipping")
		return
	}
	if (!$highcharts) {
		console.error("[Highchart] $highcharts is not provided by plugin")
		return
	}

	const HC = $highcharts as any
	const type = props.constructorType ?? "chart"
	const rawOptions = toRaw(props.options)

	console.log("[Highchart] initChart - type:", type, "HC.chart:", typeof HC[type], "HC.Chart:", typeof HC.Chart)

	if (typeof HC[type] === "function") {
		chart = HC[type](container.value, rawOptions)
		console.log("[Highchart] chart created via HC[type]:", chart)
	} else if (typeof HC.Chart === "function") {
		chart = new HC.Chart(container.value, rawOptions)
		console.log("[Highchart] chart created via new HC.Chart:", chart)
	} else {
		console.error("[Highchart] Cannot create chart - HC[type] and HC.Chart are both not functions", { HC, type })
	}
}

onMounted(() => {
	initChart()
})

watch(
	() => props.options,
	(newOptions) => {
		if (chart) {
			chart.update(toRaw(newOptions), true, true, true)
		} else {
			initChart()
		}
	},
	{ deep: true }
)

onBeforeUnmount(() => {
	if (chart) {
		chart.destroy()
		chart = null
	}
})
</script>

<template>
	<div ref="container"></div>
</template>
