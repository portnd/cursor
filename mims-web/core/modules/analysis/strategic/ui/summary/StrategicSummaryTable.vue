<script setup lang="ts">
import { useStrategicAnalysisDashboardStore } from "../../store"

const store = useStrategicAnalysisDashboardStore()
</script>

<template>
	<div class="table-responsive">
		<table v-if="store.getSummaryTable?.length" class="table customize-basic-table mb-0" style="width: max-content">
			<thead>
				<tr>
					<th class="text-center" width="120"></th>
					<th class="text-center" width="150">แผนงบประมาณ</th>
					<th v-for="(year, index) of store.getBarChartCategories1" :key="index" class="text-center" width="300">
						ปี {{ year }}
					</th>
				</tr>
			</thead>
			<tbody>
				<template v-for="(section, secIndex) of store.getSummaryTable" :key="secIndex">
					<tr class="row-hover">
						<th
							:key="secIndex"
							:rowspan="store.getCalculateColspan"
							class="text-center align-middle border-1 border-bottom border-gray-300"
							:style="{ borderBottom: secIndex === store.getSummaryTable.length - 1 ? 'transparent !important' : '' }"
						>
							{{ section.name }}
						</th>
						<td class="text-start">{{ section?.plans?.name }}</td>
						<td v-for="(value, i) of section?.plans?.value" :key="i" class="text-end">
							{{ toNumber(value, 2) }}
						</td>
					</tr>
					<tr
						v-for="(_, i) of section.plans_2"
						:key="i"
						class="row-hover"
						:class="
							i === section.plans_2.length - 1 ? ' text-center align-middle border-1 border-bottom border-gray-300' : ''
						"
					>
						<td class="text-start">{{ section.plans_2[i].name }}</td>
						<td v-for="(value, valueIndex) of section.plans_2[i].value" :key="valueIndex" class="text-end">
							{{ toNumber(value, 2) }}
						</td>
					</tr>
				</template>
			</tbody>
		</table>
		<div v-else class="text-center">ไม่พบข้อมูล</div>
	</div>
</template>

<style scoped>
.cursor-hover:hover {
	background-color: #d9d9d9;
}

@media (max-width: 767px) {
	.customize-basic-table {
		width: max-content;
	}
}

tbody tr td:first-of-type {
	border-left: 1px solid var(--kt-gray-300) !important;
}

tr:hover td {
	background-color: #f4f4f4;
}

tr:hover th:hover + td,
tr:hover th:hover ~ td {
	background-color: transparent;
}

.customize-basic-table tr:last-of-type td:first-of-type,
.customize-basic-table tr:hover:last-of-type td:first-of-type {
	border-radius: 0px 0px 0px 0px !important;
}

tr td:last-of-type {
	border-radius: 0px;
}

tbody tr > td {
	border-bottom: 1px solid var(--kt-gray-300) !important;
}

tbody tr > th:last-child {
	border-bottom: transparent !important;
}
</style>
