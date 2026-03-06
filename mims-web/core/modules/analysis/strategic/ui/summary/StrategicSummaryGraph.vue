<script setup lang="ts">
import { useStrategicAnalysisDashboardStore } from "../../store"

const store = useStrategicAnalysisDashboardStore()

const graph1 = ref()
const bar1 = ref()
const bar2 = ref()

watch(
	() => store.data,
	() => {
		if (Object.keys(store.data).length > 0) {
			graph1?.value?.updateOptions({
				color: store.data.graph1?.color,
				xaxis: {
					categories: store.getGraphCategories1,
					title: {
						text: "ปี พ.ศ.",
					},

					// tickAmount: 5,
				},
				legend: {
					position: "bottom",
					offsetY: 10,
					horizontalAlign: "center",
					showForSingleSeries: true,
					customLegendItems: store.getGraph1LegendLabel,
					itemMargin: {
						horizontal: 15,
					},
					onItemClick: {
						toggleDataSeries: true,
					},
					onItemHover: {
						highlightDataSeries: true,
					},
				},
			})

			bar1?.value?.updateOptions({
				colors: store.getBar1ColorList,
				xaxis: {
					categories: store.getBarChartCategories1,
					title: {
						text: "ปี พ.ศ.",
					},
					// tickAmount: 5,
				},
				legend: {
					position: "bottom",
					offsetY: 10,
					horizontalAlign: "center",
					showForSingleSeries: true,
					customLegendItems: store.getBar1LegendLabel,
					itemMargin: {
						horizontal: 15,
					},
					onItemClick: {
						toggleDataSeries: true,
					},
					onItemHover: {
						highlightDataSeries: true,
					},
				},
			})

			bar2?.value?.updateOptions({
				// colors: store.data.bar2?.name,
				xaxis: {
					categories: store.getBarChartCategories2,
					title: {
						text: "ปี พ.ศ.",
					},
					legend: {
						itemMargin: {
							horizontal: 15,
						},
						showForSingleSeries: true,
						customLegendItems: store.getLegendBar2,
						markers: {
							radius: 12,
							// fillColors: ["#FDB833", "#66abe6"],
						},
						onItemClick: {
							toggleDataSeries: true,
						},
						onItemHover: {
							highlightDataSeries: true,
						},
					},
					// tickAmount: 5,
				},
			})
		}
		// store.updateGraph(graph1, bar1, bar2)
	},
	{ deep: true }
)

watch(
	() => store.data,
	() => {
		if (!Object.keys(store.data).length) {
			graph1?.value?.updateOptions({
				legend: {
					show: false,
				},
			})
			bar1?.value?.updateOptions({
				legend: {
					show: false,
				},
			})
			bar2?.value?.updateOptions({
				legend: {
					show: false,
				},
			})
		}
	},
	{ deep: true }
)

// watch(
// 	() => store.getLegendBar2,
// 	() => {
// 		if (store.getLegendBar2.length === 0) {
// 			bar2.value.updateOptions({
// 				legend: {
// 					show: false,
// 				},
// 			})
// 		}
// 	}
// )

const lineAnnualInfo = reactive({
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
		stroke: {
			width: 4,
		},
		colors: [],
		xaxis: {
			categories: [2567, 2568, 2569, 2570, 2571],
			// tickAmount: 5,
			title: {
				text: "ปี พ.ศ.",
			},
			axisBorder: {
				show: false,
			},
			axisTicks: {
				show: false,
			},
		},
		yaxis: {
			title: {
				text: "ค่า IRI (ม./กม.)",
			},
			labels: {
				formatter: function (value: number) {
					if (value !== null) {
						return value.toFixed(2)
					}
				},
			},
		},
		title: {
			text: "IRI หลังวิเคราะห์รายปี",
			align: "center",
			offsetY: 10,
			style: {
				fontSize: "14px",
			},
		},
		legend: {
			show: true,
		},
		markers: {
			size: 4,
		},
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
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ series, dataPointIndex, w }: any) {
				const labels = w.globals.categoryLabels
				const colors = w.globals.colors
				// const filterPlan = store.data?.bar2?.datasets
				// 	?.filter((item) => item.plan === store.plan)
				// 	.flatMap((item) => {
				// 		return item.data.flatMap((child) => {
				// 			return child.lable
				// 		})
				// 	})
				let html = `<div class="apexcharts-result">`
				html += `<label class="fw-bold">ปี ${labels[dataPointIndex]}</label>`
				w.globals.seriesNames.forEach((name: any, key: number) => {
					if (series[key][dataPointIndex] !== null) {
						html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
							colors[key]
						};"></span> ${name}: ${series[key][dataPointIndex].toFixed(2)}</div>`
					}
				})
				html += `</div>`

				return html
			},
		},
	},
})

const barAnnualInfo = reactive({
	chartOptions: {
		chart: {
			type: "bar",
			toolbar: {
				show: false,
			},
			zoom: {
				enabled: false,
			},
		},
		title: {
			text: "งบประมาณซ่อมบำรุงรายปี",
			align: "center",
			offsetY: 10,
			style: {
				fontSize: "14px",
			},
		},
		labels: [],
		colors: [],
		plotOptions: {
			bar: {
				horizontal: false,
				columnWidth: "70%",
				endingShape: "rounded",
				// distributed: true,
				borderRadius: 4,
				dataLabels: {
					position: "top",
				},
			},
		},
		dataLabels: {
			enabled: false,
			offsetY: -20,
			style: {
				fontSize: "10px",
				colors: ["#304758"],
			},
		},
		stroke: {
			show: true,
			width: 1,
			colors: ["transparent"],
			curve: "smooth",
		},
		legend: {
			position: "bottom",
			offsetY: 5,
			horizontalAlign: "center",
			showForSingleSeries: true,
			itemMargin: {
				horizontal: 5,
				verticalal: 5,
			},
			onItemClick: {
				toggleDataSeries: true,
			},
			onItemHover: {
				highlightDataSeries: true,
			},
			markers: {
				width: 16,
				height: 16,
				radius: 4,
			},
		},
		xaxis: {
			categories: [],
			title: {
				text: "ปี พ.ศ.",
			},
			labels: {
				show: true,
			},
			axisBorder: {
				show: false,
			},
			axisTicks: {
				show: false,
			},
			// tickAmount: 0,
		},
		yaxis: {
			title: {
				text: "ค่าซ่อมบำรุง (ล้านบาท)",
			},
			labels: {
				formatter: function (value: number) {
					if (value !== null) {
						return formatNumberWithSuffix(value)
					}
				},
			},
		},

		tooltip: {
			shared: true,
			intersect: false,
			Html: true,
			enabled: true,
			custom: function ({ series, dataPointIndex, w }: any) {
				const labels = w.globals.labels
				const colors = w.globals.colors

				let html = `<div class="apexcharts-result">`
				html += `<label class="fw-bold">ปี ${labels[dataPointIndex]}</label>`
				w.globals.seriesNames.forEach((name: any, key: number) => {
					const color = store.getBar1Color(key) || colors[key]

					if (series[key][dataPointIndex] !== null) {
						html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${color};"></span> ${name}: ${toNumber(
							series[key][dataPointIndex]
						)}</div>`
					}
				})
				html += `</div>`

				return html
			},
		},
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

const barRepair = reactive({
	series: [
		{
			name: "OL-Overlay",
			data: [null, null, null, null, null, null],
		},
		{
			name: "M&OL-Mill&Overlay",
			data: [8.76, 8.76, 8.76, 8.76, 8.76, 8.76],
		},
		{
			name: "RCL-Recycling",
			data: [null, null, null, null, null, null],
		},
		{
			name: "Rc-Reconstruction",
			data: [78.18, 78.18, 78.18, 78.18, 78.18, 78.18],
		},
		{
			name: "SS-SlurrySeal",
			data: [null, null, null, null, null, null],
		},
		{
			name: "FDR",
			data: [9.57, 9.57, 9.57, 9.57, 9.57, 9.57],
		},
		{
			name: "BCO",
			data: [null, null, null, null, null, null],
		},
		{
			name: "M&OL",
			data: [null, null, null, null, null, null],
		},
		{
			name: "Seal",
			data: [3.49, 3.49, 3.49, 3.49, 3.49, 3.49],
		},
	],
	chartOptions: {
		chart: {
			type: "bar",
			toolbar: {
				show: false,
			},
			stacked: true,
			zoom: {
				enabled: false,
			},
		},
		colors: [],
		plotOptions: {
			bar: {
				horizontal: false,
				borderRadius: 4,
				columnWidth: "70%",
				dataLabels: {
					total: {
						enabled: false,
						style: {
							fontSize: "13px",
							fontWeight: 900,
						},
					},
				},
			},
		},
		dataLabels: {
			enabled: true,
			offsetY: 3.2,
			style: {
				fontSize: "8px",
				colors: ["#fff"],
			},
			formatter: function (val: any) {
				return val ? val + "%" : ""
			},
		},
		stroke: {
			show: true,
			// width: 4,
			colors: ["transparent"],
			curve: "smooth",
		},
		legend: {
			show: true,
			itemMargin: {
				horizontal: 5,
				verticalal: 5,
			},
			showForSingleSeries: true,
			markers: {
				width: 16,
				height: 16,
				radius: 4,
			},
			onItemClick: {
				toggleDataSeries: true,
			},
			onItemHover: {
				highlightDataSeries: true,
			},
		},
		xaxis: {
			categories: ["2567", "2568", "2569", "2570", "2571", "2572"],
			title: {
				text: "ปี พ.ศ.",
			},
			labels: {
				show: true,
			},
			axisBorder: {
				show: false,
			},
			axisTicks: {
				show: false,
			},
		},
		yaxis: {
			title: {
				text: "สัดส่วนวิธีการซ่อมบำรุง (%)",
			},
			max: 100,
			labels: {
				formatter: function (value: any) {
					if (value !== null) {
						return value
					}
				},
			},
		},
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
		tooltip: {
			shared: true,
			intersect: false,
			Html: true,
			enabled: true,
			custom: function ({ series, dataPointIndex, w }: any) {
				const labels = w.globals.labels
				const colors = w.globals.colors
				const budgets = [0, 64695400, 0, 577288250, 0, 70661010, 0, 0, 25757900]
				let html = `<div class="apexcharts-result">`
				html += `<label class="fw-bold">ปี ${labels[dataPointIndex]}</label>`
				w.globals.seriesNames.forEach((name: any, key: number) => {
					if (series[key][dataPointIndex]) {
						html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
							colors[key]
						};"></span> ${name}: ${toNumber(budgets[key])} (${toNumber(series[key][dataPointIndex], 2)}%)</div>`
					}
				})
				html += `</div>`

				return html
			},
		},
	},
})

// Store lifecycle is managed by the parent StrategicSummaryList component — do NOT dispose here
</script>

<template>
	<!-- ใส่รายละเอียด -->
	<!-- เส้นกึ่งกลางกราฟ กับรายละเอียดข้อมูลให้เท่ากัน -->
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 mb-5">
		<div class="card p-1 shadow-none border border-1 border-gray-200 h-100">
			<dl class="row p-5 mb-0">
				<dt class="col-12 h5 fw-semibold">สายทาง</dt>
				<dd class="col-12 p-1 card shadow-none border border-ra border-1 border-gray-200">
					<ul v-if="!store.loading && store.getRoadsList.length" class="scroll-y-custom mb-0">
						<template v-for="(road, index) in store.getRoadsList" :key="index">
							<li class="fs-6 mb-2">
								<div class="row">
									<div class="col">
										<span class="text-gray-600">{{ road }}</span>
									</div>
								</div>
							</li>
						</template>
					</ul>
					<div v-else class="ps-3">
						<span>-</span>
					</div>
				</dd>
				<dt class="col-12 mt-3 h5 fw-semibold">ตัวกรอง</dt>
				<dd class="col-12">
					<ul v-show="!store.loading">
						<li class="fs-6">
							<span class="text-gray-600">{{ store.getFilter }}</span>
						</li>
						<li v-show="store.getFilterList.length" class="fs-6">
							<div v-for="(item, index) of store.getFilterList" :key="index" class="text-gray-600">{{ item }}</div>
						</li>
					</ul>
				</dd>
				<dt class="col-12 mt-3 h5 fw-semibold">เงื่อนไข</dt>
				<dd class="col-12">
					<ul v-show="!store.loading">
						<li class="fs-6">
							<span class="text-gray-600">{{ store.getDashboardConditions }}</span>
						</li>
					</ul>
				</dd>
				<dt class="col-12 mt-3 h5 fw-semibold">ความเห็น</dt>
				<dd class="col-12">
					<ul v-show="!store.loading">
						<li
							v-for="(item, index) of store.getComment"
							v-show="store.getComment.length > 1"
							:key="index"
							class="fs-6 text-gray-600"
						>
							{{ item }}
						</li>
						<p v-show="store.getComment.length === 1">
							{{ store.getComment[0] === "-" ? "-" : store.getComment[0] }}
						</p>
					</ul>
				</dd>
			</dl>
		</div>
	</div>
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 mb-5">
		<div class="card p-0 shadow-none border border-1 border-gray-200 h-100">
			<ClientOnly>
				<apexchart
					v-show="!store.loading"
					ref="graph1"
					type="line"
					height="430"
					:options="lineAnnualInfo.chartOptions"
					:series="store.getLineChartSeries"
				/>
			</ClientOnly>
		</div>
	</div>
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 mb-5">
		<div class="card p-0 shadow-none border border-1 border-gray-200 h-100">
			<!-- ---------- -->
			<ClientOnly>
				<apexchart
					v-show="!store.loading"
					ref="bar1"
					type="bar"
					height="500"
					:options="barAnnualInfo.chartOptions"
					:series="store.getBar1Series"
				/>
			</ClientOnly>
		</div>
	</div>
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 mb-5">
		<div class="card p-2 shadow-none border border-1 border-gray-200 h-100">
			<span class="text-center fw-semibold mt-1 text-gray-800" style="font-size: 14px">
				แผนงบประมาณวิธีการซ่อมรายปี</span
			>

			<ul class="nav nav-tabs nav-line-tabs mt-2" style="overflow-x: auto; overflow-y: hidden; flex-wrap: unset">
				<li
					v-for="i of store.data?.number_plan"
					:key="i"
					class="nav-item"
					:class="store.plan === `แผนที่ ${i}` ? 'active' : ''"
					data-bs-toggle="tab"
					data-bs-target="#detail-graph"
					role="tab"
					aria-selected="true"
					@click="() => store.togglePlan(`แผนที่ ${i}`)"
				>
					<span class="nav-link cursor-pointer px-4" style="font-size: 14px; width: max-content">แผนที่ {{ i }}</span>
					<span class="line"></span>
				</li>
				<li
					v-show="store.data?.bar2?.datasets?.some((item) => item.plan === 'ไม่จำกัดงบประมาณ')"
					class="nav-item"
					:class="store.plan === 'ไม่จำกัดงบประมาณ' ? 'active' : ''"
					data-bs-toggle="tab"
					data-bs-target="#detail-graph"
					role="tab"
					aria-selected="false"
					@click="() => store.togglePlan('ไม่จำกัดงบประมาณ')"
				>
					<span class="nav-link cursor-pointer px-3 nolimit">ไม่จำกัดงบประมาณ</span>
					<span class="line"></span>
				</li>
			</ul>

			<div v-show="!store.loading" class="col-12 order-2">
				<!-- begin::Content -->
				<div class="tab-content p-2">
					<div id="detail-graph" class="tab-pane fade active show" role="tabpanel">
						<ClientOnly>
							<apexchart
								ref="bar2"
								type="bar"
								height="400"
								:options="barRepair.chartOptions"
								:series="store.getBar2Series"
							/>
						</ClientOnly>
					</div>
				</div>
				<!-- end::Content -->
			</div>
		</div>
	</div>
</template>

<style scoped>
li::marker {
	color: var(--kt-text-gray-500) !important;
}

.scroll-content {
	min-height: 505px;
}

.scroll-y-custom {
	height: 160px;
	overflow-y: scroll;
	overflow-x: hidden;
}

.square {
	width: 16px;
	height: 16px;
	border-radius: 5px;
	display: inline-block;
}

.selected {
	opacity: 0.2;
}

.nolimit {
	text-overflow: ellipsis;
	overflow: hidden;
	min-width: 100px;
	width: max-content;
	white-space: nowrap;
}
</style>
