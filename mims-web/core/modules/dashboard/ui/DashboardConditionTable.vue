<script setup lang="ts">
import { useDashboardStore } from "../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const store = useDashboardStore()

const handleHeader = computed(() => {
	const headers: THeader[] = [
		{ text: "เลน", value: "no" },
		{ text: "ระยะทาง (กม.)", value: "total_km" },
		{ text: "ค่าเฉลี่ย (ม./กม.)", value: "avg_value" },
	]
	store.conditionList.forEach((item) => {
		headers.push({ text: item.label, value: "value_" + item.id })
	})

	return headers
})

const data = [
	{
		no: 1,
		total_km: 6,
		avg_value: 7,
		value_1: 11,
		value_2: 22,
		value_3: 33,
		value_4: 44,
	},
	{
		no: 2,
		total_km: 6,
		avg_value: 7,
		value_1: 11,
		value_2: 22,
		value_3: 33,
		value_4: 44,
	},
	{
		no: 3,
		total_km: 6,
		avg_value: 7,
		value_1: 11,
		value_2: 22,
		value_3: 33,
		value_4: 44,
	},
]
</script>
<template>
	<VDatatable :headers="handleHeader" :items="data">
		<template #customize-headers>
			<thead>
				<tr>
					<th rowspan="2" class="text-center">เลน</th>
					<th rowspan="2" class="text-center">ระยะทาง (กม.)</th>
					<th rowspan="2" class="text-center">ค่าเฉลี่ย (ม./กม.)</th>
					<th :colspan="store.labelsArr.length" class="text-center border-bottom border-1">ระยะทางในแต่ละช่วง (กม.)</th>
				</tr>
				<tr>
					<th v-for="item in store.conditionList" :key="item.id" class="text-center">{{ item.label }}</th>
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
