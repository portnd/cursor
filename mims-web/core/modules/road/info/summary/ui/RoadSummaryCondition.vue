<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import { ILaneList } from "../infrastructure"

const store = useRoadSummaryStore()

useStoreLifecycle(store)

const lanes = computed(() => {
	return store.laneList.map((e: ILaneList) => {
		return { value: e.lane_no.toString(), label: e.lane_no.toString() }
	})
})

const getPercentage = (index: number, items: Array<number>, times: Array<string>) => {
	let message = ""
	if (index > 0) {
		let diff = 0
		const result = []
		for (let i = 0; i < index; i++) {
			diff = Math.abs(((items[i] - items[index]) / items[i]) * 100)
			result.push(
				`${items[index] > items[i] ? ` - เพิ่มขึ้น ${diff.toFixed(1)}% ` : ` - ลดลง ${diff.toFixed(1)}% `}  (ปี ${
					times[i]
				})`
			)
		}
		message = result.join("\n")
	}

	return message
}

const chartOptions = (id: string, title: string, yAxisTitle: string) => {
	return {
		chart: {
			id,
			fontFamily: "Prompt, sans-serif",
			height: 350,
			type: "line",
			dropShadow: {
				enabled: false,
				color: "#000",
				top: 18,
				left: 7,
				blur: 10,
				opacity: 0.2,
			},
			toolbar: {
				show: false,
			},
			zoom: {
				enabled: false,
			},
		},
		dataLabels: {
			enabled: true,
		},
		stroke: {
			curve: "smooth",
		},
		title: {
			text: title,
			align: "center",
			style: {
				fontSize: 16,
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
		markers: {
			size: 1,
		},
		yaxis: {
			title: {
				text: yAxisTitle,
			},
		},
		xaxis: {
			categories: store.dataCategories,
		},
		fill: {
			opacity: 1,
		},
		tooltip: {
			Html: true,
			enabled: true,
			custom: function ({ series, seriesIndex, dataPointIndex, w }: any) {
				const seriesName = w.globals.seriesNames[0]
				const value = series[seriesIndex][dataPointIndex]

				const years = w.globals.categoryLabels
				const percentage = getPercentage(dataPointIndex, w.globals.series[seriesIndex], years)

				return `<div class="apexcharts-result">
				    <label class="fw-bold">ปี ${years[dataPointIndex]}</label>
				    <div class="mt-2 fs-7">${seriesName} ${seriesName === "ค่า IRI" ? "(ม./กม.) " : ""} : ${value}</div>
				    <div class="fs-7">${percentage.split("\n").join("<br>")}</div>
				  </div>
				`
			},
		},
	}
}

const refIRI = ref()
const refGN = ref()

// ค้นหาข้อมูล
const onSearch = async () => {
	await store.getConditionCompareAverge()
}

const IRIseries = computed(() => {
	return [
		{
			name: "ค่า IRI",
			data: store.dataIRIChart,
		},
	]
})
const IRIchartOptions = computed(() => {
	return chartOptions("iri", "ค่าเฉลี่ย IRI (ม./กม.)", "ค่าเฉลี่ย IRI (ม./กม.)")
})

const GNseries = computed(() => {
	return [
		{
			name: "ค่า GN",
			data: store.dataGNChart,
		},
	]
})
const GNchartOptions = computed(() => {
	return chartOptions("gn", "ค่าเฉลี่ย GN", "ค่าเฉลี่ย GN")
})

onMounted(() => {
	onSearch()
})

</script>

<template>
	<div class="row mb-3 px-3 pt-0">
		<div class="col-5 col-md-5 mb-2">
			<VLabel label="ช่องจราจร" class="mt-0" />
			<VSelect
				v-model="store.conditionLaneId"
				:options="lanes"
				placeholder="เลือก"
				name="lane"
				:close-on-select="true"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="onSearch"
			/>
		</div>
		<div class="col-7 col-md-7">
			<VLoading :loading="store.loading" float="end" />
		</div>
	</div>
	<div class="row mt-3">
		<div class="col-12">
			<ClientOnly>
				<apexchart ref="refIRI" type="area" :height="325" :options="IRIchartOptions" :series="IRIseries" />
			</ClientOnly>
		</div>
	</div>
	<div class="row">
		<div class="col-12">
			<ClientOnly>
				<apexchart ref="refGN" type="area" :height="325" :options="GNchartOptions" :series="GNseries" />
			</ClientOnly>
		</div>
	</div>
</template>

<style></style>
