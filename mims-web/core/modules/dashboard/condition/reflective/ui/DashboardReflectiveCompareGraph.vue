<script setup lang="ts">
import { useForm } from "vee-validate"
import { useDashboardReflectiveGraphStore } from "../store"
import { IOption } from "~/core/shared/types/Option"

const store = useDashboardReflectiveGraphStore()
useStoreLifecycle(store)
// const route = useRoute()
// const id = Number(route.params.roadId)

defineProps({
	laneOptions: {
		type: Array<IOption>,
		default: [],
	},
	yearOptions: {
		type: Array<IOption>,
		default: [],
	},
	// lineType: {
	// 	type: String,
	// },
	dataType: {
		type: String,
		default: "",
	},
})

const { setFieldError } = useForm()

const laneChartData = computed(() => {
	const baseGraphOptions = {
		chart: {
			stacked: false,
		},
		dataLabels: {
			enabled: true,
		},
		stroke: {
			curve: "smooth",
		},
		markers: {
			size: 0,
			strokeColors: "none",
		},
		legend: {
			show: true,
			showForSingleSeries: true,
		},
	}

	if (store.lane.length === 0) {
		return [
			{
				option: {
					...baseGraphOptions,
					title: {
						text: `ปี `,
						align: "center",
						style: {
							fontSize: "16px",
						},
					},
					legend: {
						show: false,
						showForSingleSeries: true,
					},
					chart: {
						...baseGraphOptions.chart,
						zoom: {
							enabled: true,
						},
						toolbar: {
							show: true,
						},
					},
					xaxis: {
						tickAmount: 10,
						title: {
							text: "กิโลเมตร (กม.)",
							offsetY: -25,
							style: {
								color: "#3F4254",
								fontSize: "12px",
								fontWeight: 400,
							},
						},
					},
					yaxis: {
						title: {
							text: "Retro Reflectivity Average",
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
				series: [],
			},
		]
	}

	const yearLaneMap = new Map()
	const laneData = store.lane.sort((a, b) => b.year - a.year)

	laneData.forEach(({ year, items }) => {
		if (!yearLaneMap.has(year)) {
			yearLaneMap.set(year, new Map())
		}

		const laneMap = yearLaneMap.get(year)

		items.forEach((item) => {
			if (!laneMap.has(item.line_no)) {
				laneMap.set(item.line_no, [])
			}

			laneMap.get(item.line_no).push(item.value)
		})
	})

	return Array.from(yearLaneMap.entries()).map(([year, laneMap]: [number, Map<number, number[]>]) => {
		const series = Array.from(laneMap.entries()).map(([laneNo, values]) => ({
			name: `เส้นจราจร ${laneNo}`,
			data: values,
		}))

		const labelXAxis = store.lane.flatMap((item) =>
			item.year === year ? item.items.map((e) => convertMeterToKm(e.km_start)) : []
		)

		const graph = {
			...baseGraphOptions,
			chart: {
				...baseGraphOptions.chart,
				type: "line",
				zoom: {
					enabled: true,
				},
				toolbar: {
					show: true,
					tools: {
						download: false,
						selection: false,
						zoom: true,
						zoomin: false,
						zoomout: false,
						pan: false,
						reset: false,
						customIcons: [
							{
								icon: '<button class="btn btn-outline btn-outline-primary rounded-3 apxchart-button-box">รีเซ็ต</button>',
								title: "รีเซ็ต",
								class: "apexcharts-custom-reset-icon",
								click: function (chart: any) {
									chart.zoomX(0)
								},
							},
						],
					},
				},
			},
			title: {
				text: `ปี ${year + 543}`,
				align: "center",
				style: {
					fontSize: "16px",
				},
			},
			xaxis: {
				categories: labelXAxis,
				tickAmount: 10,
				title: {
					text: "กิโลเมตร (กม.)",
					offsetY: -25,
					style: {
						color: "#3F4254",
						fontSize: "12px",
						fontWeight: 400,
					},
				},
			},
			yaxis: {
				// min: this.min,
				// max: this.max,
				tickAmount: 4,
				labels: {
					formatter: (val: number) => val?.toFixed(0),
				},
				title: {
					text: "Retro Reflectivity Average",
				},
			},
			dataLabels: {
				enabled: false,
			},
			stroke: {
				curve: "smooth",
				width: 2,
			},
			tooltip: {
				Html: true,
				enabled: true,
				custom: function ({ series, dataPointIndex, w }: any) {
					const labels = w.globals.categoryLabels
					const colors = w.globals.colors
					console.log(w.globals.seriesNames)

					let html = `<div class="apexcharts-result">`
					html += `<label class="fw-bold">${labels[dataPointIndex]}</label>`
					w.globals.seriesNames.forEach((name: any, key: number) => {
						if (series[key][dataPointIndex] || series[key][dataPointIndex] === 0) {
							html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
								colors[key]
							};"></span> ${name}: ${series[key][dataPointIndex]?.toFixed(2)}</div>`
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
					title: {
						text: "Retro Reflectivity Average",
					},
				},
			},
			fill: {
				opacity: 1,
			},
			legend: {
				show: true,
				showForSingleSeries: true,
			},
		}

		return { option: graph, series }
	})
})

const yearChartData = computed(() => {
	// สร้าง base graph
	const baseGraphOptions = {
		chart: {
			type: "line",
			stacked: false,
		},
		legend: {
			show: true,
			showForSingleSeries: true,
		},
	}

	if (store.year.length === 0) {
		return [
			{
				option: {
					...baseGraphOptions,
					title: {
						text: `เส้นจราจร `,
						align: "center",
						style: {
							fontSize: "16px",
						},
					},
					legend: {
						show: false,
						showForSingleSeries: true,
					},
					chart: {
						...baseGraphOptions.chart,
						zoom: {
							enabled: true,
						},
						toolbar: {
							show: true,
						},
					},

					xaxis: {
						tickAmount: 10,
						title: {
							text: "กิโลเมตร (กม.)",
							offsetY: -25,
							style: {
								color: "#3F4254",
								fontSize: "12px",
								fontWeight: 400,
							},
						},
					},
					yaxis: {
						// min: this.min,
						// max: this.max,
						tickAmount: 4,
						labels: {
							formatter: (val: number) => val?.toFixed(0),
						},
						title: {
							text: "Retro Reflectivity Average",
						},
					},
					tooltip: {
						Html: true,
						enabled: true,
						custom: function ({ series, dataPointIndex, w }: any) {
							const labels = w.globals.categoryLabels
							console.log(labels)
							const colors = w.globals.colors

							let html = `<div class="apexcharts-result">`
							html += `<label class="fw-bold">${labels[dataPointIndex]}</label>`
							w.globals.seriesNames.forEach((name: any, key: number) => {
								if (series[key][dataPointIndex] !== null) {
									html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
										colors[key]
									};"></span> ${name}: ${generateNumber(series[key][dataPointIndex])}</div>`
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
				series: [],
			},
		]
	}

	const laneYearMap = new Map()
	const yearData = store.year.sort((a, b) => a.line - b.line)

	yearData.forEach(({ line, items }) => {
		if (!laneYearMap.has(line)) {
			laneYearMap.set(line, new Map())
		}

		const laneMap = laneYearMap.get(line)

		items.forEach((item) => {
			if (!laneMap.has(item.year)) {
				laneMap.set(item.year, [])
			}

			laneMap.get(item.year).push(item.value)
		})
	})

	const result = Array.from(laneYearMap.entries()).map(([lane, laneMap]: [number, Map<number, number[]>]) => {
		const series = Array.from(laneMap.entries()).map(([year, values]) => ({
			name: `ปี ${year + 543}`,
			data: values,
		}))

		const maxItem = [...new Set(store.year.flatMap((item) => item.items.flatMap((child) => child.km_start)))]

		// const toggleName = toggleType(props.dataType)

		const labelXAxis = maxItem.map((e) => convertMeterToKm(e))

		// Merge specific options with base options
		const graph = {
			...baseGraphOptions,
			chart: {
				...baseGraphOptions.chart,
				zoom: {
					enabled: true,
				},
				toolbar: {
					show: true,
					tools: {
						download: false,
						selection: false,
						zoom: true,
						zoomin: false,
						zoomout: false,
						pan: false,
						reset: false,
						customIcons: [
							{
								icon: '<button class="btn btn-outline btn-outline-primary rounded-3 apxchart-button-box">รีเซ็ต</button>',
								title: "รีเซ็ต",
								class: "apexcharts-custom-reset-icon",
								click: function (chart: any) {
									chart.zoomX(0)
								},
							},
						],
					},
				},
			},
			legend: {
				show: true,
				showForSingleSeries: true,
			},
			dataLabels: {
				enabled: false,
			},
			stroke: {
				curve: "smooth",
				width: 2,
			},
			markers: {
				size: 0,
				strokeColors: "none",
			},
			title: {
				text: `เส้นจราจร ${lane}`,
				align: "center",
				style: {
					fontSize: "16px",
				},
			},
			xaxis: {
				categories: labelXAxis,
				tickAmount: 10,
				title: {
					text: "กิโลเมตร (กม.)",
					offsetY: -25,
					style: {
						color: "#3F4254",
						fontSize: "12px",
						fontWeight: 400,
					},
				},
			},
			yaxis: {
				// min: this.min,
				// max: this.max,
				tickAmount: 4,
				labels: {
					formatter: (val: number) => val?.toFixed(0),
				},
				title: {
					text: "Retro Reflectivity Average",
				},
			},
			tooltip: {
				Html: true,
				enabled: true,
				custom: function ({ series, dataPointIndex, w }: any) {
					const labels = w.globals.categoryLabels
					console.log(labels)
					const colors = w.globals.colors

					let html = `<div class="apexcharts-result">`
					html += `<label class="fw-bold">${labels[dataPointIndex]}</label>`
					w.globals.seriesNames.forEach((name: any, key: number) => {
						if (series[key][dataPointIndex] || series[key][dataPointIndex] === 0) {
							html += `<div class="mt-1 fs-7"><span class="dot" style="background-color: ${
								colors[key]
							};"></span> ${name}: ${generateNumber(series[key][dataPointIndex])}</div>`
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
		}

		return { option: graph, series }
	})

	return result
})

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)

	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})

const onClose = () => {
	store.$reset()
	setFieldError("lane_lines", "")
	setFieldError("lane_years", "")
	setFieldError("year_lines", "")
	setFieldError("year_years", "")
	hideModal()
}
const onLaneSubmit = async () => {
	let errors = false
	if (!store.laneParams.lineInput.length) {
		setFieldError("lane_lines", "โปรดระบุ")
		errors = true
	}

	if (!store.laneParams.yearsInput.length) {
		setFieldError("lane_years", "โปรดระบุ")
		errors = true
	}

	if (!errors) {
		store.laneParams.lines = store.laneParams.lineInput
		store.laneParams.years = store.laneParams.yearsInput
		await store.getCompareLane()
	}
}

const onYearSubmit = async () => {
	let errors = false
	if (!store.yearParams.lineInput.length) {
		setFieldError("year_lines", "โปรดระบุ")
		errors = true
	}

	if (!store.yearParams.yearsInput.length) {
		setFieldError("year_years", "โปรดระบุ")
		errors = true
	}

	if (!errors) {
		store.yearParams.lines = store.yearParams.lineInput
		store.yearParams.years = store.yearParams.yearsInput
		await store.getCompareYear()
	}
}

</script>

<template>
	<div id="modal-compare-graph" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xxl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">กราฟเปรียบเทียบ</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" @click="onClose"></button>
				</div>
				<div class="modal-body">
					<div class="row">
						<div class="col-lg-6">
							<div class="row">
								<div class="col-12 col-md-7">
									<h4 class="fw-semibold text-gray-600">เปรียบเทียบรายเส้นจราจร</h4>
								</div>
								<div class="col-12 col-md-5 d-block d-lg-none">
									<div class="ms-4">
										<VLoading :loading="store.loading" :margin-top="'0'" float="end" />
									</div>
								</div>
							</div>
							<div class="row px-5">
								<div class="col-6 px-2">
									<VSelect
										v-model="store.laneParams.yearsInput"
										:options="yearOptions"
										mode="multiple"
										label="ปี"
										name="lane_years"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-6 px-2">
									<VSelect
										v-model="store.laneParams.lineInput"
										mode="multiple"
										:options="laneOptions"
										label="เส้นจราจร"
										name="lane_lines"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-md-12 text-end mt-3">
									<BtnSearch :disabled="store.loading" @click="onLaneSubmit" />
								</div>
							</div>
							<div class="row min-height">
								<div v-for="(item, index) of laneChartData" :key="`chartlane-${index}`" class="col-12 mt-5">
									<ClientOnly>
										<apexchart
											:key="index"
											height="350"
											:options="item.option"
											:series="item.series"
											:name="`indexlane${index}`"
										/>
									</ClientOnly>
								</div>
							</div>
						</div>
						<!-- <hr class="divider d-none d-lg-block" /> -->
						<div class="col-lg-6 mt-5 mt-lg-0">
							<div class="row">
								<div class="col-12 col-md-7">
									<h4 class="fw-semibold text-gray-600">เปรียบเทียบรายปี</h4>
								</div>
								<div class="col-12 col-md-5">
									<div class="ms-4">
										<VLoading :loading="store.loading" :margin-top="'0'" float="end" />
									</div>
								</div>
							</div>
							<div class="row px-5">
								<div class="col-6 px-2">
									<VSelect
										v-model="store.yearParams.lineInput"
										:options="laneOptions"
										mode="multiple"
										label="เส้นจราจร"
										name="year_lines"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-6 px-2">
									<VSelect
										v-model="store.yearParams.yearsInput"
										:options="yearOptions"
										mode="multiple"
										label="ปี"
										name="year_years"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-12 text-end mt-3">
									<BtnSearch :disabled="store.loading" @click="onYearSubmit" />
								</div>
							</div>
							<div class="row min-height">
								<div v-for="(item, index) of yearChartData" :key="`chartyear-${index}`" class="col-12 mt-5">
									<ClientOnly>
										<apexchart height="350" :options="item.option" :series="item.series" :name="`indexyear${index}`" />
									</ClientOnly>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.min-height {
	min-height: 25rem;
	padding-bottom: 1.5rem;
}
.divider {
	align-self: stretch;
	border: 1px solid #ddd;
	height: inherit;
	min-height: 88%;
	max-width: 0;
	width: 0;
	vertical-align: text-bottom;
	padding: 0;
	margin: 0;
	position: absolute;
	left: 50%;
	margin-top: 10%;
}
</style>
