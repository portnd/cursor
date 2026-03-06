<script setup lang="ts">
import { useDashboardStore } from "../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const store = useDashboardStore()

const handleSurface = (id: number) => {
	const index = store.data.surface_array.indexOf(id)
	if (index === -1) {
		store.data.surface_array.push(id)
	} else {
		store.data.surface_array.splice(index, 1)
	}
	store.data.surface_array.sort((a, b) => a - b)
	store.colors = store.data.surface_colors.flatMap((item, key) => {
		return store.data.surface_array
			.map((id) => {
				if (key === id) {
					return item
				}
				return undefined
			})
			.filter((item: any) => item !== undefined)
	})

	// store.data.surface_map.forEach((item) => {
	// 	store.createLine(item.the_geom?.coordinates, item.color)
	// })
}

const handleHeader = computed(() => {
	const headers: THeader[] = [
		{ text: "", value: "no" },
		{ text: "", value: "surface" },
	]
	const data = store.data.surface?.surface_dashboard_table?.find((item: any) => item)

	if (data) {
		const lane = Object.keys(data?.surface_lane_type).reverse()
		lane.forEach((item) => {
			headers.push({ text: "", value: item })
		})
	}

	// headers.push({ text: "", value: "summary_right" })
	return headers
})

const handleData = computed(() => {
	const items = ref<Array<any>>([])
	store.data.surface?.surface_dashboard_table?.forEach((item: any) => {
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

	if (items.value.length) {
		items.value.push({
			surface: "รวม",
			more_than_four: toNumber(totalMoreThanFour, 2),
			four_lane: toNumber(totalFourLane, 2),
			three_lane: toNumber(totalThreeLane, 2),
			two_lane: toNumber(totalTwoLane, 2),
			one_lane: toNumber(totalOneLane, 2),
			summary_right: toNumber(totalSummaryRight, 2),
			isTotal: true,
		})
	}

	return items.value
})

const onSync = () => {
	store.syncDataMart()
}
</script>
<template>
	<div class="row">
		<div class="col-12">
			<div class="row">
				<div class="col-9"></div>
				<div class="col-3">
					<button
						class="col-12 btn btn-outline btn-outline-primary btn-disabled rounded-2 mt-md-0 mt-sm-2 mt-3 p-0"
						:disabled="store.syncing"
						@click="onSync()"
					>
						<div v-if="store.syncing">
							<span class="spinner-border spinner-border-sm align-middle me-5"></span>
							<span>{{ store.dataMart?.percent.toFixed(2) }}%</span>
						</div>
						<div v-else>
							<i class="fi fi-br-refresh"></i>
							<span>ซิงค์ข้อมูล</span>
						</div>
					</button>
				</div>
				<div class="fw-normal fs-5 text-end mb-2">
					<template v-if="store?.dataMart">
						<label class="text-gray-600 me-2 fs-8 mt-3">ข้อมูลอัปเดตข้อมูลโดย {{ store.dataMart.updated_by }}</label>
						<label class="text-gray-600 fs-8 mt-3">
							เมื่อวันที่
							{{ buddhistFormatDate(store.dataMart.updated_at, "dd mmm yyyy เวลา HH:ii น.") }}
						</label>
					</template>
				</div>
			</div>
		</div>
		<div class="col-12 pt-5">
			<div class="card-chart h-100 text-center p-4 pt-4">
				<ClientOnly>
					<apexchart ref="pieChart" type="pie" height="275" :options="store.barOptions()" :series="store.barSeries()" />
				</ClientOnly>
				<div class="row justify-content-center mt-8">
					<template v-for="(item, key) in store.data?.surface?.summary" :key="key">
						<div
							class="col-auto d-flex align-items-center cursor-pointer"
							:class="store.data.surface_array.includes(key) ? '' : 'selected'"
							@click="handleSurface(key)"
						>
							<div class="square my-2 me-2" :style="`background: ${item.surface?.color_code}`"></div>
							<span>{{ item.surface?.name }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
		<div class="col-12 pt-5">
			<VDatatable :headers="handleHeader" :items="handleData">
				<template #customize-headers>
					<thead>
						<tr>
							<th rowspan="2" class="text-center">ลำดับ</th>
							<th rowspan="2" class="text-center">ประเภทผิวทาง</th>
							<th colspan="5" class="text-center border-bottom border-1">ระยะทาง (กม.)</th>
						</tr>
						<tr>
							<th class="text-center">มากกว่า 4 ช่องจราจร</th>
							<th class="text-center">4 ช่องจราจร</th>
							<th class="text-center">3 ช่องจราจร</th>
							<th class="text-center">2 ช่องจราจร</th>
							<th class="text-center">1 ช่องจราจร</th>
						</tr>
					</thead>
				</template>
				<template v-for="(header, index) in handleHeader" #[`item-${header.value}`]="{ item }" :key="index">
					<div v-if="header.value === 'no'">
						{{ item.isTotal ? "" : item[header.value] }}
					</div>
					<div v-else>
						{{ item[header.value] }}
					</div>
				</template>
			</VDatatable>
		</div>
	</div>
</template>

<style lang="scss" scoped>
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
