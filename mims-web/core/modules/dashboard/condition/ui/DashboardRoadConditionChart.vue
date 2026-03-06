<script setup lang="ts">
import { useDashboardRoadConditionStore } from "../store"
import { DashboardConditionCompareGraph } from "./index"


const store = useDashboardRoadConditionStore()

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
	roadId: {
		type: Array<Number>,
		required: true,
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
// 	() => props.searchCall,
// 	(_) => {
// 		if (props.roadId.length === 1 && !props.roadId.includes(5)) {
// 			store.getRoads()
// 		}
// 	}
// )

// watch(
// 	() => props.reloadData,
// 	(_) => {
// 		// store.getRoads()
// 		store.getConditionType()
// 	}
// )
// const loading = computed(() => store.loading)

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
	<div class="card card-rounded p-5 mt-5">
		<div class="d-flex">
			<div class="p-2 flex-grow-1">
				<div class="d-block d-sm-flex p-0">
					<ul class="w-100 nav nav-tabs nav-line-tabs justify-content-center justify-content-md-start">
						<li class="w-50 nav-item text-center">
							<span
								class="nav-link cursor-pointer fs-5"
								:class="{ active: store.params.graph_type !== 'HISTOGRAM', disabled: store.loading }"
								@click="store.toggleGraph"
								>GRAPH</span
							>
							<span class="line"></span>
						</li>
						<li class="w-50 nav-item text-center">
							<span
								class="nav-link cursor-pointer fs-5"
								:class="{ active: store.params.graph_type !== 'GRAPH', disabled: store.loading }"
								@click="store.toggleGraph"
								>HISTOGRAM</span
							>
							<span class="line"></span>
						</li>
					</ul>
				</div>
			</div>
			<div class="px-2">
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
		<!-- <VSkeletonLoader :loading="loading"> -->
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
				<div class="col-12 col-lg-12 mt-3">
					<VCheckbox
						v-model="store.params.criteria_id"
						:options="store.criteriaOptions"
						name="criteria_id"
						@update:model-value="() => store.onUpdateCheckbox()"
					/>
				</div>
				<!-- <div class="col-12 col-lg-5 mt-1 text-end">
				<span class="text-gray-900 me-4">ระยะทางที่สนใจ 2.67 กม.</span>
				<span class="text-gray-900">ระยะทางทั้งหมด 2.67 กม.</span>
			</div> -->
			</div>
		</div>
		<!-- </VSkeletonLoader> -->

		<!-- Modal -->
		<DashboardConditionCompareGraph
			ref="modalCompareGraph"
			:search-call="props.searchCall"
			:reload-data="props.reloadData"
		/>
		<!-- <RoadConditionSurveyRules ref="modalCompareSurveyRules" /> -->
	</div>
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}
</style>
