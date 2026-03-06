<script setup lang="ts">
import { useForm } from "vee-validate"
import { useGraphCompareStore, useConditionStore } from "../store"

const store = useGraphCompareStore()
const conditionStore = useConditionStore()
useStoreLifecycle(store)
const route = useRoute()
const id = Number(route.params.roadId)

onMounted(async () => {
	await store.getLaneList(id)
})

const toggleType = (name: string) => {
	const typeName = name.toUpperCase()

	if (typeName === "IRI" || typeName === "IFI") {
		return "ค่า " + typeName + " (ม./กม.)"
	} else if (typeName === "RUT" || typeName === "MPD") {
		return "ค่า " + typeName + " (มม.)"
	} else {
		return "ค่า " + typeName
	}
}

const { setFieldError } = useForm()

const laneChartData = computed(() => {
	const toggleName = toggleType(store.conditionType)

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
							text: toggleName,
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
			if (!laneMap.has(item.lane_no)) {
				laneMap.set(item.lane_no, [])
			}

			laneMap.get(item.lane_no).push(item.value)
		})
	})

	return Array.from(yearLaneMap.entries()).map(([year, laneMap]: [number, Map<number, number[]>]) => {
		const series = Array.from(laneMap.entries()).map(([laneNo, values]) => ({
			name: `ช่องจราจร ${laneNo}`,
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
					text: toggleName,
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
						text: toggleType(store.conditionType),
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
	const toggleName = toggleType(store.conditionType)

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
						text: `ช่องจราจร `,
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
							text: toggleName,
						},
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
	const yearData = store.year.sort((a, b) => a.lane - b.lane)

	yearData.forEach(({ lane, items }) => {
		if (!laneYearMap.has(lane)) {
			laneYearMap.set(lane, new Map())
		}

		const laneMap = laneYearMap.get(lane)

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
				text: `ช่องจราจร ${lane}`,
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
					text: toggleName,
				},
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

const showModal = (dataType: string) => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	store.conditionType = dataType
	store.getLaneList(id)
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
	setFieldError("lane_lanes", "")
	setFieldError("lane_years", "")
	setFieldError("year_lanes", "")
	setFieldError("year_years", "")
	hideModal()
}
const onLaneSubmit = async () => {
	let errors = false
	if (!store.laneParams.laneInput.length) {
		setFieldError("lane_lanes", "โปรดระบุ")
		errors = true
	}

	if (!store.laneParams.yearsInput.length) {
		setFieldError("lane_years", "โปรดระบุ")
		errors = true
	}

	if (!errors) {
		store.laneParams.lanes = store.laneParams.laneInput
		store.laneParams.years = store.laneParams.yearsInput
		await store.getCompareLane(id)
	}
}

const onYearSubmit = async () => {
	let errors = false
	if (!store.yearParams.laneInput.length) {
		setFieldError("year_lanes", "โปรดระบุ")
		errors = true
	}

	if (!store.yearParams.yearsInput.length) {
		setFieldError("year_years", "โปรดระบุ")
		errors = true
	}

	if (!errors) {
		store.yearParams.lanes = store.yearParams.laneInput
		store.yearParams.years = store.yearParams.yearsInput
		await store.getCompareYear(id)
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
									<h4 class="fw-semibold text-gray-600">เปรียบเทียบรายช่องจราจร</h4>
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
										:options="conditionStore.getYearOptions"
										mode="tags"
										label="ปี"
										name="lane_years"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-6 px-2">
									<VSelect
										v-model="store.laneParams.laneInput"
										mode="tags"
										:options="store.getLaneListOptions"
										label="ช่องจราจร"
										name="lane_lanes"
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
										v-model="store.yearParams.laneInput"
										:options="store.getLaneListOptions"
										mode="tags"
										label="ช่องจราจร"
										name="year_lanes"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
								<div class="col-6 px-2">
									<VSelect
										v-model="store.yearParams.yearsInput"
										:options="conditionStore.getYearOptions"
										mode="tags"
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
