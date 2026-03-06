<script setup lang="ts">
import HistoryProjectTab from "../HistoryProjectTab.vue"
import { useMaintenanceHistoryPlanStore } from "../../../store/MaintenanceHistoryPlanStore"
import HistoryProjectPlanTable from "./HistoryProjectPlanTable.vue"

const route = useRoute()
const id = Number(route.params.id)
const store = useMaintenanceHistoryPlanStore()
useStoreLifecycle(store)
const refPlanGraph: Ref = ref()
const refDisbursementGraph: Ref = ref()

onMounted(async () => {
	await store.getPlanList(id)
	await store.setGraph(id, store.params.planId, refPlanGraph, refDisbursementGraph)
})

const planChart = reactive({
	chartOptions: {
		chart: {
			type: "line",
			toolbar: {
				show: false,
			},
			zoom: {
				enabled: false,
			},
		},
		dataLabels: {
			enabled: false,
			background: {
				borderRadius: 8,
			},
		},
		xaxis: {
			categories: [],
			tickAmount: 10,
		},
		yaxis: {
			min: 0,
			max: 100,
			tickAmount: 4,
			title: {
				text: "ความก้าวหน้า (%)",
			},
			labels: {
				formatter: function (value: number) {
					return value.toFixed(0)
				},
			},
		},
		title: {
			text: "แผนดำเนินงาน",
			align: "center",
		},
		legend: {
			position: "bottom",
			offsetY: 10,
			horizontalAlign: "center",
			itemMargin: {
				horizontal: 15,
			},
		},
		grid: {
			show: true,
			borderColor: "#e7e7e7",
			strokeDashArray: 5,
			xaxis: {
				lines: {
					show: false,
				},
			},
			yaxis: {
				lines: {
					show: true,
				},
			},
		},
		colors: ["#e7e7e7"],
		markers: {
			size: 5,
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ series, dataPointIndex, w }: any) {
				const labels = w.globals.categoryLabels
				const colors = w.globals.colors

				let html = `<div class="apexcharts-result">`
				html += `<label class="fw-bold">${labels[dataPointIndex]}</label>`
				w.globals.seriesNames.forEach((name: any, key: number) => {
					if (series[key][dataPointIndex] !== null) {
						html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${colors[key]};"></span> ${name}: ${series[key][dataPointIndex]}%</div>`
					}
				})
				html += `</div>`

				return html
			},
		},
	},
})

const disbursementChart = reactive({
	chartOptions: {
		chart: {
			type: "line",
			toolbar: {
				show: false,
			},
			zoom: {
				enabled: false,
			},
		},
		dataLabels: {
			enabled: false,
			background: {
				borderRadius: 8,
			},
		},
		yaxis: {
			title: {
				offsetX: 6,
				offsetY: 2,
				text: "การเบิกจ่าย (บาท)",
			},
			labels: {
				formatter: function (value: number) {
					if (value !== null) {
						return formatNumberWithSuffix(value)
					}
				},
			},
		},
		xaxis: {
			categories: [],
			tickAmount: 10,
		},
		title: {
			text: "การเบิกจ่ายเงิน",
			align: "center",
		},
		plotOptions: {
			series: {
				turboThreshold: 10000,
				marker: {
					enabled: true,
				},
			},
		},
		legend: {
			position: "bottom",
			offsetY: 10,
			horizontalAlign: "center",
			itemMargin: {
				horizontal: 15,
			},
		},
		grid: {
			show: true,
			borderColor: "#e7e7e7",
			strokeDashArray: 5,
			xaxis: {
				lines: {
					show: false,
				},
			},
			yaxis: {
				lines: {
					show: true,
				},
			},
		},
		colors: ["#e7e7e7"],
		markers: {
			size: 5,
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ series, dataPointIndex, w }: any) {
				const labels = w.globals.categoryLabels
				const colors = w.globals.colors
				let html = `<div class="apexcharts-result">`
				html += `<label class="fw-bold">${labels[dataPointIndex]}</label>`
				w.globals.seriesNames.forEach((name: any, key: number) => {
					if (series[key][dataPointIndex] !== null) {
						html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
							colors[key]
						};"></span> ${name}: ${toNumber(series[key][dataPointIndex])} บาท</div>`
					}
				})
				html += `</div>`

				return html
			},
		},
	},
})

function fillNullWithPrevious(arr: Array<number | null>) {
	let lastNonNullValue: number | null = null
	return arr.map((value) => {
		if (value !== null) {
			lastNonNullValue = value
		}
		return value !== null ? value : lastNonNullValue
	})
}

const planSeries = computed(() => {
	if (store.graphData.length === 0) {
		return [
			{
				name: "งวดงาน",
				data: [],
			},
			{
				name: "การเบิกจ่าย (บาท)",
				data: [],
			},
		]
	}

	const series = store.graphData
		.flatMap((parent) => {
			if (parent.name === "แผนการดำเนินงาน") {
				return parent.data.map((child) => {
					const filledData = fillNullWithPrevious(child.data)
					return { name: child.name, data: filledData.map((value) => value) }
				})
			}

			return undefined
		})
		.filter((item) => item !== undefined)

	return series
})

const disbursmentSeries = computed(() => {
	if (store.graphData.length === 0) {
		return [
			{
				name: "งวดงาน",
				data: [],
			},
			{
				name: "การเบิกจ่าย (บาท)",
				data: [],
			},
		]
	}

	const series = store.graphData
		.flatMap((parent) => {
			if (parent.name === "การเบิกจ่ายเงิน") {
				return parent.data.map((child) => {
					const filledData = fillNullWithPrevious(child.data)
					return { name: child.name, data: filledData.map((value) => value) }
				})
			}

			return undefined
		})
		.filter((item) => item !== undefined)
	return series
})

const onSearch = async (id: number, planId: number[]) => {
	await store.setGraph(id, planId, refPlanGraph, refDisbursementGraph)
	store.getSumDataTable()
}

const downloadFile = () => {
	const params = store.params.planId.length > 0 ? `?plan_id=${encodeURIComponent(store.params.planId.join(","))}` : ""
	useDownloadFile(`ดาวน์โหลดไฟล์`, `maintenance/${id}/plan_progress_export_report${params}`)
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-12">
			<div class="card p-5">
				<HistoryProjectTab />
				<div class="row px-5 pt-3">
					<div class="col-12 text-end p-0">
						<NuxtLink
							class="btn btn-outline btn-outline-primary rounded-4 px-4 py-3 mb-3 mb-sm-0 fw-semibold"
							@click="downloadFile"
						>
							ดาวน์โหลด PDF
						</NuxtLink>
					</div>
					<div class="col-12 col-md-5 mb-2">
						<VSelect
							v-model="store.params.planId"
							:options="store.getPlanListOptions"
							mode="multiple"
							:close-on-select="false"
							label="แผน"
							name="plan"
							placeholder="เลือก"
							limit="1"
							@update:model-value="(planId: any) => onSearch(id, planId)"
						/>
					</div>
					<div class="col-12 col-md-7 text-end">
						<VLoading :loading="store.loading" float="end" class="mt-md-10" />
					</div>

					<VSkeletonLoader :loading="store.loading">
						<div class="col-12 mb-5 container-chart mt-3 mt-md-0">
							<ClientOnly>
								<apexchart
									ref="refPlanGraph"
									type="line"
									height="340"
									:options="planChart.chartOptions"
									:series="planSeries"
								/>
							</ClientOnly>
						</div>
						<div class="col-12 container-chart">
							<ClientOnly>
								<apexchart
									ref="refDisbursementGraph"
									type="line"
									height="340"
									:options="disbursementChart.chartOptions"
									:series="disbursmentSeries"
								/>
							</ClientOnly>
						</div>
					</VSkeletonLoader>
				</div>
				<VSkeletonLoader :loading="store.loading">
					<HistoryProjectPlanTable />
				</VSkeletonLoader>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
