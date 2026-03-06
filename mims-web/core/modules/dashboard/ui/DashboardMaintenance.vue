<script setup lang="ts">
import ServerSideDataTable from "../../common/datatable/ui"
import { useDashboardStore } from "../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const emits = defineEmits(["dataTable"])
const store = useDashboardStore()
const dataTable = ref()
const refQuantityChart = ref()
const refMaintenanceBudgetChart = ref()
const refBarChart = ref()

const countChart = reactive({
	chartOptions: {
		title: {
			text: "จำนวนโครงการในช่วงรับประกันผลงาน",
			align: "center",
		},
		plotOptions: {
			pie: {
				donut: {
					size: "70%",
					labels: {
						show: true,
						value: {
							show: true,
							fontSize: "16px",
						},
						total: {
							show: true,
							showAlways: true,
							label: 0,
							fontSize: "22px",
							fontWeight: 500,
							color: "#FDB833",
							formatter: function () {
								return "โครงการ"
							},
						},
					},
				},
			},
		},
		stroke: {
			show: false,
		},
		labels: {
			formatter: function (value: number) {
				if (value !== null) {
					return formatNumberWithSuffix(value)
				}
			},
		},
		// labels: ["ทางหลวงพิเศษหมายเลข 7", "ทางหลวงพิเศษหมายเลข 9", "โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง"],
		dataLabels: {
			enabled: false,
		},
		grid: {
			padding: {
				top: 20,
				left: 0,
				right: 0,
				bottom: 20,
			},
		},
		// colors: ["#FDB833", "#FFDA6A", "#A86F00"],
		legend: {
			show: false,
		},
		tooltip: {
			Html: true,
			enabled: true,
			theme: "light",
			custom: function ({ series, seriesIndex, w }: any) {
				const total = series.reduce((a: number, b: number) => a + b, 0)
				let html = `<div class="apexcharts-result justify-content-start">`
				html += `<label class="fs-7 mb-2" style="border-bottom: 1px solid #c9c9c9">${w.config.labels[seriesIndex]}</label>`
				html += `
          <div class="d-flex w-100 fs-7">
            <span style="background-color: ${w.globals.colors[seriesIndex]};" class="me-2 mt-1 dot"></span>
            ${series[seriesIndex]} โครงการ: ${toNumber((series[seriesIndex] / total) * 100, 2)} %
          </div>
          `
				html += `</div>`
				return html
			},
		},
	},
})
const budgetChart = reactive({
	chartOptions: {
		title: {
			text: "งบประมาณซ่อมบำรุง",
			align: "center",
		},
		plotOptions: {
			pie: {
				donut: {
					size: "70%",
					labels: {
						show: true,
						value: {
							show: true,
							fontSize: "16px",
						},
						total: {
							show: true,
							showAlways: true,
							fontSize: "22px",
							fontWeight: 500,
							color: "#FDB833",
							formatter: function () {
								return "บาท"
							},
						},
					},
				},
			},
		},
		stroke: {
			show: false,
		},
		dataLabels: {
			enabled: false,
		},
		grid: {
			padding: {
				top: 20,
				left: 0,
				right: 0,
				bottom: 20,
			},
		},
		colors: ["#FDB833", "#FFDA6A", "#A86F00"],
		legend: {
			show: false,
		},
		tooltip: {
			Html: true,
			enabled: true,
			theme: "light",
			custom: function ({ series, seriesIndex, w }: any) {
				const total = series.reduce((a: number, b: number) => a + b, 0)
				let html = `<div class="apexcharts-result justify-content-start">`
				html += `<label class="fs-7 mb-2" style="border-bottom: 1px solid #c9c9c9">${w.config.labels[seriesIndex]}</label>`
				html += `
          <div class="d-flex w-100 fs-7">
            <span style="background-color: ${w.globals.colors[seriesIndex]};" class="me-2 mt-1 dot"></span>
            ${toNumber(series[seriesIndex])} บาท: ${toNumber((series[seriesIndex] / total) * 100, 2)} %
          </div>
          `
				html += `</div>`
				return html
			},
		},
	},
})

const barChart = reactive({
	chartOptions: {
		chart: {
			toolbar: {
				show: false,
			},
		},
		dataLabels: {
			enabled: false,
		},
		tooltip: {
			enabled: true,
			y: {
				show: true,
				formatter: (value: number) => {
					return toNumber(value, 2) + " บาท"
				},
			},
		},
		stroke: {
			show: true,
			width: 2,
			colors: ["transparent"],
			curve: "smooth",
		},
		legend: {
			show: false,
		},
		title: {
			text: "10 อันดับโครงการที่ได้รับงบประมาณการซ่อมมากที่สุด",
			align: "center",
			style: {
				fontSize: "16px",
			},
		},
		subtitle: {
			text: `ณ วันที่ ${buddhistFormatDate(new Date(), "dd/mm/yyyy")}`,
			align: "right",
		},
		colors: ["#FDB833"],
		plotOptions: {
			bar: {
				horizontal: false,
				columnWidth: "35",
				distributed: true,
				endingShape: "rounded",
				borderRadius: 4,
			},
		},
		// labels: {
		// 	formatter: function (value: number) {
		// 		if (value !== null) {
		// 			return formatNumberWithSuffix(value)
		// 		}
		// 	},
		// },
		xaxis: {
			categories: [],
		},
		yaxis: [
			{
				show: true,
				title: {
					text: "งบประมาณ (บาท)",
				},
				labels: {
					show: true,
					style: {
						fontWeight: 400,
					},
					formatter: (value: number) => {
						if (value) {
							return formatNumberWithSuffix(Number(value.toFixed(0)))
						} else {
							return 0
						}
					},
				},
			},
		],
		grid: {
			borderColor: "#EAEAEA",
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
		fill: {
			opacity: 1,
		},
	},
})

const handleHeader = computed(() => {
	const headers: THeader[] = [
		{ text: "ลำดับ", value: "no" },
		{ text: "เลขสัญญา", value: "contract_number" },
		{ text: "สายทาง", value: "road" },
		{ text: "ตอนควบคุม", value: "road_section" },
		{ text: "หน่วยงาน", value: "owner" },
		{ text: "งบประมาณ (ล้านบาท)", value: "budget" },
		{ text: "ระยะเวลาในช่วงรับประกันผลงาน", value: "work_guarantee" },
		{ text: "ระยะเวลาเหลือติดค้ำประกัน", value: "guarantee_date" },
	]
	return headers
})

watch(
	() => [store.data.maintenance, store.data.maintenance_quantity_index, store.data.maintenance_budget_index],
	() => {
		refQuantityChart.value.updateOptions({
			chart: {
				events: {
					dataPointSelection: (_: any, __: any, opts: any) => {
						const dataIndex = opts.dataPointIndex
						const roadGroupId = store.data?.maintenance?.number_maintenance_chart?.road_group_id ?? []

						if (roadGroupId.length) {
							navigateTo(`maintenances/history?road_group_id_dashboard=${roadGroupId[dataIndex]}`)
						}
					},
				},
			},
			plotOptions: {
				pie: {
					donut: {
						labels: {
							total: {
								show: true,
								showAlways: true,
								label: toNumber(store.getSumMaintenanceChart?.sum_quantity),
							},
						},
					},
				},
			},
			labels: store.getMaintnenanceQuantity.label,
			colors: store.getMaintnenanceQuantity.colors,
		})

		refMaintenanceBudgetChart.value.updateOptions({
			chart: {
				events: {
					dataPointSelection: (_: any, __: any, opts: any) => {
						const dataIndex = opts.dataPointIndex
						const roadGroupId = store.data?.maintenance?.maintenance_budget_chart?.road_group_id ?? []

						if (roadGroupId.length) {
							navigateTo(`maintenances/history?road_group_id=${roadGroupId[dataIndex]}`)
						}
					},
				},
			},
			plotOptions: {
				pie: {
					donut: {
						labels: {
							total: {
								show: true,
								showAlways: true,
								label: toNumber(store.getSumMaintenanceChart?.sum_budget),
							},
						},
					},
				},
			},
			labels: store.getMaintenanceBudget.label,
			colors: store.getMaintenanceBudget.colors,
		})

		refBarChart.value.updateOptions({
			chart: {
				events: {
					dataPointSelection: (_: any, __: any, opts: any) => {
						const dataIndex = opts.dataPointIndex
						const maintenanceId = store.data?.maintenance?.top_ten_maintenance_budget_chart?.maintenance_id ?? []
						if (maintenanceId.length) {
							navigateTo(`maintenances/history/${maintenanceId[dataIndex]}/info`)
						}
					},
				},
			},
			xaxis: {
				categories: store.getMaintenanceBarChart?.categories,
			},
			colors: store.getMaintenanceBarChart?.colors,
			subtitle: {
				text: `ณ วันที่ ${store.getUpdateMaintenanceDate}`,
				align: "right",
			},
		})
	},
	{ deep: true }
)

watch(
	() => dataTable.value,
	() => {
		emits("dataTable", dataTable.value)
	}
)
</script>
<template>
	<div class="row">
		<div class="col-6 pt-5">
			<div class="card-chart h-100 text-center p-4">
				<ClientOnly>
					<apexchart
						ref="refQuantityChart"
						type="donut"
						height="300"
						:options="countChart.chartOptions"
						:series="store.getMaintnenanceQuantity.series"
					/>
				</ClientOnly>
				<div class="row px-xl-1">
					<template v-for="(item, key) in store.getMaintenanceQuantityCheckbox" :key="key">
						<div
							class="col-md-6 col-12 d-flex align-items-center cursor-pointer mt-2"
							:class="store.data.maintenance_quantity_index?.includes(key) ? '' : 'selected'"
							@click="store.handleCheckbox(key, 'maintenance-quantity', refQuantityChart)"
						>
							<div class="square col-1" :style="`background: ${item.color}`"></div>
							<span class="ms-2 text-start">{{ item.name }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
		<div class="col-6 pt-5">
			<div class="card-chart h-100 text-center p-4">
				<ClientOnly>
					<apexchart
						ref="refMaintenanceBudgetChart"
						type="donut"
						height="300"
						:options="budgetChart.chartOptions"
						:series="store.getMaintenanceBudget.series"
					/>
				</ClientOnly>
				<div class="row px-xl-1">
					<template v-for="(item, key) in store.getMaintenanceBudgetCheckbox" :key="key">
						<div
							class="col-md-6 col-12 d-flex align-items-center cursor-pointer xxl-set-50 mt-2"
							:class="store.data.maintenance_budget_index?.includes(key) ? '' : 'selected'"
							@click="store.handleCheckbox(key, 'maintenance-budget', refMaintenanceBudgetChart)"
						>
							<div class="square col-1" :style="`background: ${item.color}`"></div>
							<span class="ms-2 fs-7 text-start">{{ item.name }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
		<div class="col-12 pt-5">
			<div class="card-chart h-100 text-center p-4">
				<ClientOnly>
					<apexchart
						ref="refBarChart"
						type="bar"
						height="315"
						:options="barChart.chartOptions"
						:series="store.getMaintenanceBarChart?.series"
					/>
				</ClientOnly>
			</div>
		</div>
		<div class="col-12 pt-5 px-2 px-4">
			<label class="fs-5">โครงการในช่วงรับประกันผลงาน</label>
			<ServerSideDataTable ref="dataTable" url="dashboard/maintenance_table" :headers="handleHeader">
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-contract_number="{ item }">
					<div class="text-center">{{ item.contract_number }}</div>
				</template>
				<template #item-road="{ item }">
					<div class="text-center">
						<template v-for="(text, i) in item.road_name" :key="i">
							<div>{{ text }}</div>
						</template>
					</div>
				</template>
				<template #item-road_section="{ item }">
					<div class="text-center">
						<template v-for="(text, i) in item.section_name" :key="i">
							<div>{{ text }}</div>
						</template>
					</div>
				</template>
				<template #item-owner="{ item }">
					<div class="text-center">
						<template v-for="(text, i) in item.ref_depot_name" :key="i">
							<div>{{ text }}</div>
						</template>
					</div>
				</template>
				<template #item-budget="{ item }">
					<div class="text-end">{{ toNumber(item.budget, 2) }}</div>
				</template>
				<template #item-work_guarantee="{ item }">
					<div class="text-center">{{ item.guarantee_expiration_date }}</div>
				</template>
				<template #item-guarantee_date="{ item }">
					<div class="text-center">{{ item.remain_date }}</div>
				</template>
			</ServerSideDataTable>
		</div>
	</div>
</template>

<style lang="scss" scoped>
.card-chart {
	border: solid 1px var(--kt-gray-300);
	border-radius: 1rem;
}
.square {
	width: 20px !important;
	height: 20px !important;
	border-radius: 5px;
	display: inline-block;
}

.selected {
	opacity: 0.2;
}

@media only screen and (min-width: 1877px) {
	.xxl-set-50 {
		width: 50%;
	}
}
</style>
