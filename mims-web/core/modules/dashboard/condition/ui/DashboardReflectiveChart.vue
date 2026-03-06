<script setup lang="ts">
import { useDashboardReflectiveStore } from "../store"
// import { useDashboardStore } from "../store"
import // RoadReflectiveSurveyRules,
// RoadReflectiveCompareGraph,
"~/core/modules/road/info/reflectiveStrip/ui/index"

import {
	// DashboardReflectiveSurveyRules,
	DashboardReflectiveCompareGraph,
} from "~/core/modules/dashboard/condition/reflective/ui/index"

// import { DashboardConditionCompareGraph } from "./index"

// import { identity } from "@vueuse/core";

const store = useDashboardReflectiveStore()

const props = defineProps({
	collapsed: {
		type: Boolean,
	},
	mapShow: {
		type: Boolean,
	},
	searchCall: {
		type: Boolean,
		default: false,
	},
	reloadData: {
		type: Boolean,
		default: false,
	},
})

const barChart = ref()
const highcharts = ref()

// begin::กราฟเปรียบเทียบ
const modalCompareGraph: Ref = ref()
const compareGraph = () => {
	modalCompareGraph.value.showModal()
}
// end::กราฟเปรียบเทียบ

// begin::กำหนดเกณฑ์
// const modalCompareSurveyRules: Ref = ref()
// const compareSurveyRules = () => {
// 	modalCompareSurveyRules.value.showModal()
// }
// end::กำหนดเกณฑ์

const graphWidth = ref(0)
const containerGraph = ref()
const handleResizeGraph = () => {
	graphWidth.value = containerGraph.value?.offsetWidth

	if (graphWidth.value > 0) {
		store.lineChart.chart.width = Number(graphWidth.value) - 20
	}
}

onMounted(() => {
	window.addEventListener("resize", handleResizeGraph)
	handleResizeGraph()

	setTimeout(() => {
		handleResizeGraph()
	}, 1500)
})

// watch(
// 	() => [props.searchCall, props.reloadData],
// 	async (_) => {
// 		await store.getRoads()
// 	}
// )

watch(
	() => [props.mapShow, store.lineChart, props.collapsed],
	(_) => {
		window.addEventListener("resize", handleResizeGraph)
		handleResizeGraph()

		setTimeout(() => {
			handleResizeGraph()
		}, 1500)
	}
)

onUnmounted(() => {
	window.removeEventListener("resize", handleResizeGraph)
})
</script>

<template>
	<div class="d-flex justify-content-between" style="min-height: 40px">
		<div class="d-block d-sm-flex p-0">
			<ul class="nav nav-tabs nav-line-tabs mb-5 justify-content-center justify-content-md-start">
				<li class="nav-item">
					<span
						class="nav-link cursor-pointer"
						:class="{ active: store.params.graph_type === 'graph', disabled: store.loading }"
						@click="store.toggleGraph"
						>GRAPH</span
					>
					<span class="line"></span>
				</li>
				<li class="nav-item">
					<span
						class="nav-link cursor-pointer"
						:class="{ active: store.params.graph_type === 'histogram', disabled: store.loading }"
						@click="store.toggleGraph"
						>HISTOGRAM</span
					>
					<span class="line"></span>
				</li>
				<li v-show="store.details.has_white_line" class="nav-item">
					<span
						class="nav-link cursor-pointer"
						:class="`${store.params.line_color === 'white' ? 'active' : ''}
			${store.loading ? 'disabled' : ''}`"
						@click="store.toggleData('white')"
						>เส้นสีขาว</span
					>
					<span class="line"></span>
				</li>
				<li v-show="store.details.has_white_line" class="nav-item">
					<span
						class="nav-link cursor-pointer"
						:class="`${store.params.line_color === 'yellow' ? 'active' : ''}
			${store.loading ? 'disabled' : ''}`"
						@click="store.toggleData('yellow')"
						>เส้นสีเหลือง</span
					>
					<span class="line"></span>
				</li>
			</ul>
		</div>
		<div class="">
			<!-- <div class="col-12 text-center text-sm-end mb-2"
			:class="!collapsed ? 'col-md-12 col-lg-12 col-xl-4 ps-0' : 'col-md-4'"> -->
			<!-- <a
				type="button"
				style="max-height: 38px"
				class="btn btn-outline btn-outline-primary rounded-4 px-2 py-1 me-3 mb-2 fw-semibold fs-6 lh-xxl"
				:disabled="store.loading"
				@click="compareSurveyRules()"
			>
				กำหนดเกณฑ์
			</a> -->
			<a
				type="button"
				style="max-height: 38px"
				class="btn btn-outline btn-outline-primary rounded-4 px-2 py-1 mb-2 fw-semibold fs-6 lh-xxl"
				:disabled="store.loading"
				@click="compareGraph()"
			>
				กราฟเปรียบเทียบ
			</a>
		</div>
	</div>
	<!-- <VSkeletonLoader :loading="store.loading"> -->
	<div ref="containerGraph" class="row container-chart">
		<div class="col-12">
			<div class="w-100">
				<ClientOnly>
					<highchart
						v-if="store.params.graph_type !== 'histogram'"
						ref="highcharts"
						:options="store.lineChart"
						class="mt-5"
					/>
					<apexchart
						v-else
						ref="barChart"
						height="295"
						class="hist-chart"
						:options="store.histGraph()"
						:type="'bar'"
						:series="store.histSeries()"
					/>
					<!-- <h1 v-if="store.params.graph_type !== 'histogram'" ref="highcharts"> --- highchart : {{
							store.params.graph_type }}</h1>
						<h1 v-else> --- apex : {{ store.params.graph_type }}</h1> -->
				</ClientOnly>
			</div>
			<div class="col-12 col-lg-12 mt-3">
				<div class="d-flex justify-content-between">
					<VCheckbox
						v-model="store.params.criteria_id"
						:options="store.criteriaOptions"
						name="criteria"
						@update:model-value="store.onUpdateCheckbox"
					/>
					<VCheckbox
						v-model="store.params.line_type_id"
						:options="store.getLineTypeOptions"
						name="line_type"
						@update:model-value="store.onUpdateCheckbox"
					/>
				</div>
			</div>
			<!-- <div class="col-12 col-lg-5 mt-1 text-end">
			<span class="text-gray-900 me-4">ระยะทางที่สนใจ 2.67 กม.</span>
			<span class="text-gray-900">ระยะทางทั้งหมด 2.67 กม.</span>
		</div> -->
		</div>
	</div>
	<!-- </VSkeletonLoader> -->
	<!-- Modal -->
	<DashboardReflectiveCompareGraph
		ref="modalCompareGraph"
		:year-options="store.getYearOptions"
		:lane-options="store.getLineListOptions"
		:line-type="store.params.line_color"
	/>
	<!-- <DashboardReflectiveSurveyRules ref="modalCompareSurveyRules" /> -->
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}
</style>
