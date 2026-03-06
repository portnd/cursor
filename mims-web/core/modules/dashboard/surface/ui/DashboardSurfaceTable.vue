<script setup lang="ts">
import { useDashboardSurfaceStore } from "../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const store = useDashboardSurfaceStore()

const surface = computed(() => {
	const label = ref<Array<string>>([])
	const data = store.data.surface_dashboard_table?.find((item: any) => item)
	if (data) {
		const lane = Object.keys(data?.surface_lane_type)
		lane.forEach((item) => {
			if (item.includes("more")) {
				label.value.unshift("มากกว่า 4 เลน")
			} else if (item.includes("four")) {
				label.value.unshift("4 เลน")
			} else if (item.includes("three")) {
				label.value.unshift("3 เลน")
			} else if (item.includes("two")) {
				label.value.unshift("2 เลน")
			} else {
				label.value.unshift("1 เลน")
			}
		})
	}
	label.value.push("รวม")
	return label.value
})

const handleHeader = computed(() => {
	const headers: THeader[] = [{ text: "", value: "surface" }]
	const data = store.data.surface_dashboard_table?.find((item: any) => item)
	if (data) {
		const lane = Object.keys(data?.surface_lane_type).reverse()
		lane.forEach((item) => {
			headers.push({ text: "", value: item })
		})
	}

	headers.push({ text: "", value: "summary_right" })
	return headers
})

const handleData = computed(() => {
	const items = ref<Array<any>>([])
	store.data.surface_dashboard_table?.forEach((item: any) => {
		const dataRow = ref<any>({})
		for (const key in item) {
			if (key.includes("name")) {
				dataRow.value.surface = item[key]
			} else if (key.includes("lane")) {
				const lane = Object.keys(item[key]).reverse()
				const value = Object.values(item[key]).reverse()
				for (let i = 0; i < lane.length; i++) {
					dataRow.value[lane[i]] = Number(value[i]).toLocaleString("th-TH", {
						minimumFractionDigits: 2,
						maximumFractionDigits: 2,
					})
				}
			}
		}

		dataRow.value.summary_right = toNumber(
			Number(dataRow.value.more_than_four) +
				Number(dataRow.value.three_lane) +
				Number(dataRow.value.two_lane) +
				Number(dataRow.value.one_lane),
			2
		)

		items.value.push(dataRow.value)
	})

	const totalMoreThanFour = items.value.reduce(
		(acc, curr) => acc + parseFloat(curr.more_than_four.replaceAll(",", "")),
		0
	)
	const totalFourLane = items.value.reduce((acc, curr) => acc + parseFloat(curr.four_lane.replaceAll(",", "")), 0)
	const totalThreeLane = items.value.reduce((acc, curr) => acc + parseFloat(curr.three_lane.replaceAll(",", "")), 0)
	const totalTwoLane = items.value.reduce((acc, curr) => acc + parseFloat(curr.two_lane.replaceAll(",", "")), 0)
	const totalOneLane = items.value.reduce((acc, curr) => acc + parseFloat(curr.one_lane.replaceAll(",", "")), 0)
	const totalSummaryRight = items.value.reduce(
		(acc, curr) => acc + parseFloat(curr.summary_right.replaceAll(",", "")),
		0
	)

	items.value.push({
		surface: "รวม",
		more_than_four: toNumber(totalMoreThanFour, 2),
		four_lane: toNumber(totalFourLane, 2),
		three_lane: toNumber(totalThreeLane, 2),
		two_lane: toNumber(totalTwoLane, 2),
		one_lane: toNumber(totalOneLane, 2),
		summary_right: toNumber(totalSummaryRight, 2),
	})

	return items.value
})
</script>

<template>
	<div class="row">
		<div class="col-12">
			<div class="card p-5">
				<VDatatable :headers="handleHeader" :items="handleData">
					<template #customize-headers>
						<thead>
							<tr>
								<th rowspan="2" class="text-center">ประเภทผิวทาง</th>
								<th colspan="6" class="text-center border-bottom border-1">ระยะทาง (กม.)</th>
							</tr>
							<tr>
								<th v-for="(item, key) in surface" :key="key" class="text-center">{{ item }}</th>
							</tr>
						</thead>
					</template>
					<template v-for="(header, index) in handleHeader" #[`item-${header.value}`]="{ item }" :key="index">
						<div>
							{{ item[header.value] }}
						</div>
					</template>
				</VDatatable>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
