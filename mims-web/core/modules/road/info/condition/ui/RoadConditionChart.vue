<script setup lang="ts">
import { useConditionStore } from "../store"
import { RoadConditionSurveyRules, RoadConditionCompareGraph } from "./index"

const store = useConditionStore()

const props = defineProps({
	collapsed: {
		type: Boolean,
	},
	mapShow: {
		type: Boolean,
	},
})

const barChart = ref()
const highcharts = ref()

// begin::กราฟเปรียบเทียบ
const modalCompareGraph: Ref = ref()
const compareGraph = () => {
	modalCompareGraph.value.showModal(store.params.condition_type.toLowerCase())
}
// end::กราฟเปรียบเทียบ

// begin::กำหนดเกณฑ์
const modalCompareSurveyRules: Ref = ref()
const compareSurveyRules = () => {
	modalCompareSurveyRules.value.showModal()
}
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
	// handleResizeGraph()

	setTimeout(() => {
		handleResizeGraph()
	}, 1500)
})

watch(
	() => [props.mapShow, store.lineChart],
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
	<div class="row">
		<div class="col-12" :class="!collapsed ? 'col-md-12 col-lg-12 col-xl-8' : 'col-md-8'">
			<div class="d-block d-sm-flex p-0">
				<ul class="nav nav-tabs nav-line-tabs mb-5 me-8 justify-content-center justify-content-md-start">
					<li class="nav-item">
						<span
							class="nav-link cursor-pointer"
							:class="{ active: store.params.graph_type !== 'HISTOGRAM', disabled: store.loading }"
							@click="store.toggleGraph"
							>GRAPH</span
						>
						<span class="line"></span>
					</li>
					<li class="nav-item">
						<span
							class="nav-link cursor-pointer"
							:class="{ active: store.params.graph_type !== 'GRAPH', disabled: store.loading }"
							@click="store.toggleGraph"
							>HISTOGRAM</span
						>
						<span class="line"></span>
					</li>
				</ul>
				<ul class="nav nav-tabs nav-line-tabs mb-5 justify-content-center justify-content-md-start">
					<li v-for="(type, index) of ['IRI', 'MPD', 'RUT', 'IFI']" :key="index" class="nav-item">
						<span
							class="nav-link cursor-pointer"
							:class="`${store.params.condition_type === type ? 'active' : ''}
              ${store.loading ? 'disabled' : ''}`"
							@click="store.toggleData"
							>{{ type }}</span
						>
						<span class="line"></span>
					</li>
				</ul>
			</div>
		</div>
		<div
			class="col-12 text-center text-sm-end mb-2"
			:class="!collapsed ? 'col-md-12 col-lg-12 col-xl-4 ps-0' : 'col-md-4'"
		>
			<a
				type="button"
				class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 me-3 mb-2 fw-semibold fs-6 lh-xxl"
				:disabled="store.loading"
				@click="compareSurveyRules()"
			>
				กำหนดเกณฑ์
			</a>
			<a
				type="button"
				class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 mb-2 fw-semibold fs-6 lh-xxl"
				:disabled="store.loading"
				@click="compareGraph()"
			>
				กราฟเปรียบเทียบ
			</a>
		</div>
	</div>
	<VSkeletonLoader :loading="store.loading">
		<div ref="containerGraph" class="row container-chart">
			<div class="col-12">
				<div class="w-100">
					<ClientOnly>
						<highchart
							v-if="store.params.graph_type === 'GRAPH'"
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
					</ClientOnly>
				</div>
				<div class="col-12 col-lg-12 mt-3 d-flex justify-content-between align-items-center">
					<VCheckbox
						v-model="store.params.criteria_id"
						:options="store.criteriaOptions"
						name="criteria"
						@update:model-value="() => store.onUpdateCheckbox()"
					/>
					<span>ค่าเฉลี่ย {{ store?.graphAvg.toFixed(2) }}</span>
				</div>
				<!-- <div class="col-12 col-lg-5 mt-1 text-end">
			<span class="text-gray-900 me-4">ระยะทางที่สนใจ 2.67 กม.</span>
			<span class="text-gray-900">ระยะทางทั้งหมด 2.67 กม.</span>
		</div> -->
			</div>
		</div>
	</VSkeletonLoader>

	<!-- Modal -->
	<RoadConditionCompareGraph ref="modalCompareGraph" />
	<RoadConditionSurveyRules ref="modalCompareSurveyRules" />
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}
</style>
