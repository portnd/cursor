<script setup lang="ts">
import { useAnnualSummaryDashboardStore } from "../../store"

const store = useAnnualSummaryDashboardStore()

const pieChart = ref()
const treeChart = ref()
watch(
	() => [store.dashboard, store.pie_series_index, store.tree_series_index],
	() => {
		if (store.isInit) {
			setTimeout(() => {
				store.updateGraph(treeChart, pieChart)
			}, 100)
		} else {
			store.updateGraph(treeChart, pieChart)
		}
	},
	{ deep: true }
)

const pieCondition = reactive({
	chartOptions: {
		plotOptions: {
			pie: {
				dataLabels: {
					offset: -5,
				},
			},
		},
		labels: [
			"OL-Overlay",
			"M&OL-Mill&Overlay",
			"RCL-Recycling",
			"Rc-Reconstruction",
			"SS-SlurrySeal",
			"FDR",
			"BCO",
			"M&OL",
			"Seal",
		],
		title: {
			text: "ปริมาณงาน",
			align: "center",
			offsetY: 10,
			style: {
				fontSize: "14.95px",
			},
		},
		dataLabels: {
			style: {
				fontSize: "10px",
				colors: ["#fff"],
			},
		},
		legend: {
			show: false,
		},
		stroke: {
			show: false,
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ series, seriesIndex, w }: any) {
				const labels = w.globals.labels
				const colors = w.globals.colors
				let html = ""
				html += `<div class="apexcharts-result " style="background-color: ${colors[seriesIndex]};">`
				if (series[seriesIndex] !== null) {
					html += `${labels[seriesIndex]}: ${toNumber(store.getPieChartOptions?.area[seriesIndex])}  ตร.ม.</div>`
				}
				html += `</div>`

				return html
			},
		},
	},
})

const treeAnnualRepair = reactive({
	chartOptions: {
		chart: {
			height: 300,
			type: "treemap",
			toolbar: {
				show: false,
			},
		},
		grid: {
			show: false,
			padding: {
				left: 15,
			},
		},
		legend: {
			show: true,
			position: "bottom",
			offsetY: 10,
			horizontalAlign: "center",
			showForSingleSeries: true,
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
		// colors: ["#FF7A8D", "#FFB800", "#AF85FF", "#B22727", "#82E0AA", "#418FFF", "#FF69B4", "#7FFFD4", "#FF7F33"],
		title: {
			text: "ค่าซ่อมบำรุง",
			align: "center",
			offsetY: 10,
			style: {
				fontSize: "14.95px",
			},
		},
		dataLabels: {
			enabled: true,
			style: {
				fontSize: "12px",
				background: {
					enabled: true,
					dropShadow: {
						enabled: false,
					},
				},
			},
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ _, dataPointIndex, w }: any) {
				const labels = w.globals.categoryLabels
				const colors = w.globals.colors
				const budget = store.getTreeChartOptions?.budgets[dataPointIndex] || 0
				let html = `<div class="apexcharts-result">`
				html += `<div class="mt-1 fs-7"><span class="dot" style="background-color:${colors[dataPointIndex]};"></span> ${
					labels[dataPointIndex]
				} : ${toNumber(budget)}  บาท</div>`
				html += `</div>`

				return html
			},
		},
		plotOptions: {
			treemap: {
				distributed: true,
				enableShades: false,
			},
		},
	},
})
</script>

<template>
	<div class="col-12 mb-6">
		<div class="card p-1 shadow-none border border-1 border-gray-200 h-100">
			<dl class="row p-5 mb-0">
				<dt class="col-12 h5 fw-semibold">สายทาง</dt>
				<dd class="col-12 p-1 card shadow-none border border-ra border-1 border-gray-200">
					<ul v-if="!store.loading && store.getRoadsList.length" class="scroll-y-custom mb-0">
						<li v-for="(item, index) in store.getRoadsList" :key="index" class="fs-6 mb-2">
							<div class="row">
								<div class="col">
									<span class="text-gray-600"> {{ item }}</span>
								</div>
							</div>
						</li>
					</ul>
					<span v-else class="ps-4">-</span>
				</dd>
				<dt class="col-12 mt-3 h5 fw-semibold">ตัวกรอง</dt>
				<dd class="col-12">
					<ul v-show="!store.loading">
						<li class="fs-6">
							<span class="text-gray-600">{{ store.getDashboardFilter }}</span>
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
						<p v-show="store.getComment.length === 1">{{ store.getComment[0] === "-" ? "-" : store.getComment[0] }}</p>
					</ul>
				</dd>
			</dl>
		</div>
	</div>
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 px-5 mb-5">
		<div class="card shadow-none border border-1 border-gray-200 h-100">
			<ClientOnly>
				<apexchart
					v-show="!store.loading"
					ref="pieChart"
					type="pie"
					height="300"
					:options="pieCondition.chartOptions"
					:series="store.getPieChartOptions?.series"
				/>
			</ClientOnly>
			<div class="row justify-content-center mt-7">
				<template v-for="(item, key) in store.getPieCategories?.label" :key="key">
					<div
						class="col-auto d-flex align-items-center cursor-pointer"
						:class="store.pie_series_index?.includes(key) ? '' : 'selected'"
						@click="store.handleCheckbox(key, 'pie-chart-index')"
					>
						<div class="square my-2 me-2" :style="`background: ${store.getPieCategories?.color[key]}`"></div>
						<span>{{ item }}</span>
					</div>
				</template>
			</div>
		</div>
	</div>
	<div class="col-xl-6 col-lg-6 col-md-6 col-12 px-5 mb-5">
		<div class="card shadow-none border border-1 border-gray-200 h-100">
			<ClientOnly>
				<apexchart
					v-if="!store.loading"
					ref="treeChart"
					type="treemap"
					height="300"
					:options="treeAnnualRepair.chartOptions"
					:series="store.getTreeChartOptions?.series"
				/>
				<div v-else class="h-100 w-100"></div>
			</ClientOnly>

			<div class="row justify-content-center mt-2">
				<template v-for="(item, key) in store.getTreeLegends?.label" :key="key">
					<div
						class="col-auto d-flex align-items-center cursor-pointer"
						:class="store.tree_series_index.includes(key) ? '' : 'selected'"
						@click="store.handleCheckbox(key, 'tree-map-chart')"
					>
						<div class="square my-2 me-2" :style="`background: ${store.getTreeLegends?.color[key]}`"></div>
						<span>{{ item }}</span>
					</div>
				</template>
			</div>
		</div>
	</div>
	<!-- <div>
  <BudgetLimitTable />
</div> -->
</template>

<style scoped lang="scss">
li::marker {
	color: var(--kt-text-gray-500) !important;
}

.scroll-content {
	overflow-y: scroll;
	overflow-x: hidden;
	max-height: 450px;
}

.scroll-y-custom {
	height: 160px;
	overflow-y: scroll;
	overflow-x: hidden;
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
