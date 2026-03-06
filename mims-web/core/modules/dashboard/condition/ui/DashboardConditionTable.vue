<script setup lang="ts">
import { useDashboardConditionStore } from "../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const store = useDashboardConditionStore()

enum EConditionType {
	IRI = 1,
	MPD = 2,
	RUT = 3,
	IFI = 4,
	Reflect = 5,
}

const condition = computed(() => {
	if (store.data.chart) {
		return store.data.chart.lable
	} else {
		return []
	}
})

const handleHeader = computed(() => {
	const headers: THeader[] = [
		{ text: "เลน", value: "no" },
		{ text: "ระยะทาง (กม.)", value: "total_km" },
		{ text: "ค่าเฉลี่ย (ม./กม.)", value: "avg_value" },
	]

	if (store.data.chart) {
		store.data.chart.lable?.forEach((item: string, index: string) => {
			headers.push({ text: item, value: "value_" + index })
		})
	}
	return headers
})

const handleData = computed(() => {
	const items = ref<Array<any>>([])
	store.data.table?.forEach((item: any) => {
		const dataRow = ref<any>({})
		for (const key in item) {
			if (key.includes("lane_no")) {
				dataRow.value.lane_no = item.lane_no
			} else if (key.includes("detail_km")) {
				item[key].forEach((detail: any, index: string) => {
					dataRow.value[`value_${index}`] = detail.value.toFixed(2)
				})
			} else {
				dataRow.value[key] = item[key].toFixed(2)
			}
		}
		items.value.push(dataRow.value)
	})

	return items.value
})

const toggleType = computed(() => {
	if (store.conditionType === EConditionType.IRI) {
		return "(ม./กม.)"
	} else if (store.conditionType === EConditionType.RUT || store.conditionType === EConditionType.MPD) {
		return "(มม.)"
	} else {
		return ""
	}
})
</script>
<template>
	<VDatatable :headers="handleHeader" :items="handleData">
		<template #customize-headers>
			<thead>
				<tr>
					<th rowspan="2" class="text-center">เลน</th>
					<th rowspan="2" class="text-center">ระยะทาง (กม.)</th>
					<th rowspan="2" class="text-center">ค่าเฉลี่ย {{ toggleType }}</th>
					<th :colspan="condition?.length" class="text-center border-bottom border-1">ระยะทางในแต่ละช่วง (กม.)</th>
				</tr>
				<tr>
					<th v-for="item in condition" :key="item.id" class="text-center">{{ item }}</th>
				</tr>
			</thead>
		</template>
		<template v-for="(header, index) in handleHeader" #[`item-${header.value}`]="{ item }" :key="index">
			<div>
				{{ item[header.value] }}
			</div>
		</template>
	</VDatatable>
</template>

<style scoped></style>
