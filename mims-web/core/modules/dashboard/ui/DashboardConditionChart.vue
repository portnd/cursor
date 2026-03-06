<script setup lang="ts">
import { useDashboardStore } from "../store"
import DashboardConditionIRI from "./DashboardConditionIRI.vue"

defineProps({
	collapsed: {
		type: Boolean,
		default: false,
	},
})

const pieChart = ref()
const barChart = ref()
const histChart = ref()
const store = useDashboardStore()
const reflective = [1, 2, 3, 4]
const reflectiveList = [
	{
		label: "เส้นทึบ",
		color: "",
		value: 1,
	},
	{
		label: "เส้นปะ",
		color: "",
		value: 2,
	},
	{
		label: "เส้นทึบ - ปะ",
		color: "",
		value: 3,
	},
	{
		label: "เส้นทึบคู่",
		color: "",
		value: 4,
	},
]

const pieCondition = reactive({
	series: store.dataArray,
	chartOptions: {
		plotOptions: {
			pie: {
				dataLabels: {
					offset: 40,
				},
			},
		},
		legend: {
			show: false,
		},
		title: {
			text: "ข้อมูลสภาพทาง",
			align: "center",
			style: {
				fontSize: "18px",
			},
		},
		dataLabels: {
			enabled: true,
			style: {
				fontSize: "12px",
				fontWeight: 400,
				colors: ["#3F4254"],
			},
			dropShadow: {
				enabled: false,
			},
			formatter: (_: any, opts: any) =>
				((opts.w.config.series[opts.seriesIndex] / store.totalData) * 100).toFixed(1) + " %",
		},
		grid: {
			padding: {
				top: 20,
				left: 0,
				right: 0,
				bottom: 20,
			},
		},
		stroke: {
			show: false,
		},
		colors: ["#A4FCA5", "#42D235", "#F77A14", "#FF290A", "#973131"],
		labels: ["ดีมาก", "ดี", "ปานกลาง", "แย่", "แย่มาก"],
		tooltip: {
			enabled: true,
			fillSeriesColor: false,
			theme: "light",
			y: {
				show: true,
				formatter: (value: number) => {
					return Number(value.toFixed(2)).toLocaleString() + " กม."
				},
			},
		},
	},
})

const barCondition = reactive({
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
					const total = store.dataArray.reduce((acc, item) => acc + item)
					return ((value / total) * 100).toFixed(1) + " %"
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
			text: "ข้อมูลสรุปสภาพทาง",
			align: "center",
			style: {
				fontSize: "16px",
			},
		},
		colors: ["#A4FCA5", "#42D235", "#F77A14", "#FF290A", "#973131"],
		plotOptions: {
			bar: {
				horizontal: false,
				columnWidth: "40",
				distributed: true,
				endingShape: "rounded",
				borderRadius: 4,
				dataLabels: {
					position: "top",
				},
			},
		},
		xaxis: {
			categories: ["ดีมาก", "ดี", "ปานกลาง", "แย่", "แย่มาก"],
			title: {
				text: "เกณฑ์การประเมิน",
			},
			labels: {
				// show: false,
				style: {
					colors: "#304758",
					fontSize: "12px",
				},
			},
			axisTicks: {
				show: false,
			},
		},
		yaxis: [
			{
				title: {
					text: "ระยะทาง (กม.)",
				},
				labels: {
					formatter: (value: number) => {
						if (value) {
							return Number(value.toFixed(0)).toLocaleString()
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

const histCondition = reactive({
	chartOptions: {
		chart: {
			toolbar: {
				show: false,
			},
		},
		dataLabels: {
			enabled: true,
			distributed: true,
			offsetY: -30,
			style: { position: "absolute", fontSize: "12px", colors: ["#000000"], fontWeight: "regular" },
			formatter: function (value: any, _: any) {
				const total = store.dataArray.reduce((a: number, b: number) => a + b, 0)
				return `${value} กม. (${((value / total) * 100).toFixed(1)}%)`
			},
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ dataPointIndex }: any) {
				const total = store.dataArray.reduce((a: number, b: number) => a + b, 0)
				let html = `<div class="apexcharts-result">`
				html += `<div><span class="dot" style="background-color: ${store.colors[dataPointIndex]};"></span>${
					store.labelsArr[dataPointIndex]
				}: ${((store.dataArray[dataPointIndex] / total) * 100).toFixed(1) + " %"}  </div>`
				html += "</div>"
				return html
			},
		},
		stroke: {
			show: true,
			width: 2,
			colors: ["transparent"],
			curve: "smooth",
		},
		legend: { show: false },
		colors: ["#A4FCA5", "#42D235", "#F77A14", "#FF290A", "#973131"],
		plotOptions: {
			bar: {
				horizontal: false,
				columnWidth: "50",
				distributed: true,
				endingShape: "rounded",
				borderRadius: 4,
				dataLabels: {
					position: "top",
				},
			},
		},
		xaxis: {
			categories: ["ดีมาก", "ดี", "ปานกลาง", "แย่", "แย่มาก"],
			title: { text: "IRI (ม./กม.)", style: { color: "#3F4254", fontSize: "12px", fontWeight: 400 } },
			labels: {
				// show: false,
				style: {
					colors: "#304758",
					fontSize: "12px",
				},
			},
			axisTicks: {
				show: false,
			},
		},
		yaxis: { title: { text: "ระยะทาง (กม.)" } },
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

watch(
	() => store.conditionType,
	() => {
		store.getCheckBox()
		store.updateConditionChart(pieChart, barChart)
		if (store.toggle.hist) {
			setTimeout(() => {
				store.updatehisChart(histChart)
			}, 100)
		}
	}
)

watch(
	() => store.toggle.hist,
	() => {
		if (store.toggle.hist) {
			setTimeout(() => {
				store.updatehisChart(histChart)
			}, 100)
		}
	}
)

watch(
	() => store.menu,
	() => {
		if (store.menu === "condition") {
			store.getCheckBox()
			store.updateConditionChart(pieChart, barChart)
		}
	}
)
</script>
<template>
	<div v-if="store.roads.length === 1" class="row">
		<div class="col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<div class="row">
					<div class="col-12" :class="!collapsed ? 'col-md-12 col-lg-12 col-xl-8' : 'col-md-8'">
						<div class="d-block d-sm-flex p-0">
							<ul class="nav nav-tabs nav-line-tabs mb-5 me-8 justify-content-center justify-content-md-start">
								<li class="nav-item">
									<span
										class="nav-link cursor-pointer"
										:class="{ active: !store.toggle.hist }"
										@click="store.toggleGraph"
										>GRAPH</span
									>
									<span class="line"></span>
								</li>
								<li class="nav-item">
									<span
										class="nav-link cursor-pointer"
										:class="{ active: store.toggle.hist }"
										@click="store.toggleGraph"
										>HISTOGRAM</span
									>
									<span class="line"></span>
								</li>
							</ul>
							<ul
								v-if="store.conditionType === 5"
								class="nav nav-tabs nav-line-tabs mb-5 justify-content-center justify-content-md-start"
							>
								<li class="nav-item">
									<span
										class="nav-link cursor-pointer"
										:class="{ active: !store.toggle.whiteLine }"
										@click="store.toggleReflect"
										>เส้นสีขาว</span
									>
									<span class="line"></span>
								</li>
								<li class="nav-item">
									<span
										class="nav-link cursor-pointer"
										:class="{ active: store.toggle.whiteLine }"
										@click="store.toggleReflect"
										>เส้นสีเหลือง</span
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
						<button
							v-if="store.conditionType === 5"
							type="button"
							class="btn btn-outline btn-outline-primary px-3 py-2 me-3 mb-2 fw-semibold fs-6"
						>
							กำหนดเกณฑ์
						</button>
						<button type="button" class="btn btn-outline btn-outline-primary px-3 py-2 mb-2 fw-semibold fs-6">
							กราฟเปรียบเทียบ
						</button>
					</div>
				</div>
				<div ref="containerGraph" class="row container-chart">
					<div class="col-12">
						<div class="w-100">
							<ClientOnly>
								<highchart v-if="!store.toggle.hist" ref="highcharts" :options="store.dataGraph" class="mt-5" />
								<apexchart
									v-else
									ref="histChart"
									height="295"
									class="hist-chart"
									:options="histCondition.chartOptions"
									:type="'bar'"
									:series="store.histSeries()"
								/>
							</ClientOnly>
						</div>
						<div class="row">
							<div
								class="mt-3 text-left d-flex"
								:class="
									store.conditionType === 5
										? 'col-12 col-md-4 col-xl-3 justify-content-md-start justify-content-center'
										: 'col-12 col-lg-12'
								"
							>
								<VCheckbox v-model="store.graphSelectFilter" :options="store.conditionList" name="criteria" />
							</div>
							<div
								v-if="store.conditionType === 5"
								class="col-12 col-md-8 col-xl-9 mt-3 d-flex justify-content-md-end justify-content-center"
							>
								<VCheckbox v-model="reflective" :options="reflectiveList" name="reflective" />
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<DashboardConditionIRI />
			</div>
		</div>
	</div>
	<div class="row">
		<div class="col-xl-6 col-lg-12 col-md-6 col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<ClientOnly>
					<apexchart
						ref="pieChart"
						type="pie"
						height="275"
						:options="pieCondition.chartOptions"
						:series="pieCondition.series"
					/>
				</ClientOnly>
				<div class="row justify-content-center mt-8">
					<template v-for="(item, key) in store.conditionList" :key="key">
						<div
							class="col-auto d-flex align-items-center cursor-pointer"
							:class="store.conditionArray.includes(key + 1) ? '' : 'selected'"
						>
							<div class="square my-2 me-2" :style="`background: ${item.color}`"></div>
							<span>{{ item.label }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
		<div class="col-xl-6 col-lg-12 col-md-6 col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<ClientOnly>
					<apexchart
						ref="barChart"
						type="bar"
						height="315"
						:options="barCondition.chartOptions"
						:series="store.barSeries()"
					/>
				</ClientOnly>
				<div class="row justify-content-center">
					<template v-for="(item, key) in store.conditionList" :key="key">
						<div
							class="col-auto d-flex align-items-center cursor-pointer"
							:class="store.conditionArray.includes(key + 1) ? '' : 'selected'"
						>
							<div class="square my-2 me-2" :style="`background: ${item.color}`"></div>
							<span>{{ item.label }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}
.card-chart {
	border: solid 1px var(--kt-gray-300);
	border-radius: 1rem;
}
.square {
	width: 20px;
	height: 20px;
	border-radius: 5px;
	display: inline-block;
}

.selected {
	opacity: 0.2;
}
</style>
