<script setup lang="ts">
import { useStrategicAnalysisDashboardStore } from "../../store"

const store = useStrategicAnalysisDashboardStore()
</script>

<template>
	<div class="table-responsive">
		<table class="table customize-basic-table mb-0">
			<thead>
				<tr>
					<th rowspan="2" class="text-center" width="80">วิธีการซ่อม</th>
					<th
						v-for="(year, index) of store.getSummaryPlanTable.years"
						:key="index"
						colspan="2"
						class="text-center border-bottom border-1"
					>
						ปี {{ year }}
					</th>
				</tr>
				<tr>
					<template v-for="(_, index) of store.getSummaryPlanTable.years" :key="index">
						<th class="text-center" width="90">ระยะทาง (กม.)</th>
						<th class="text-center" width="100">งบประมาณ (บาท)</th>
					</template>
				</tr>
			</thead>
			<tbody v-if="store.getSummaryPlanTable?.plan?.length">
				<tr v-for="(item, index) of store.getSummaryPlanTable.plan" :key="index" class="hover-effect-row">
					<td class="text-center">{{ item.method_name }}</td>
					<template v-for="(_, valueIndex) of item.data" :key="valueIndex">
						<td class="text-end">
							{{ toNumber(item.data[valueIndex]?.km, 3) }}
						</td>
						<td class="text-end">{{ toNumber(item.data[valueIndex]?.budget, 2) }}</td>
					</template>
				</tr>

				<tr v-show="store.getSummaryPlanTable.plan?.length" class="text-end border-top">
					<td class="text-center">รวม</td>
					<template v-for="(_, yearIndex) of store.getSummaryPlanTable.years" :key="'total-' + yearIndex">
						<td class="text-end">
							{{
								toNumber(
									store.getSummaryPlanTable.plan?.reduce((acc, curr) => acc + (curr.data[yearIndex]?.km || 0), 0),
									3
								)
							}}
						</td>
						<td class="text-end">
							{{
								toNumber(
									store.getSummaryPlanTable.plan?.reduce((acc, curr) => acc + (curr.data[yearIndex]?.budget || 0), 0),
									2
								)
							}}
						</td>
					</template>
				</tr>
			</tbody>
			<tbody v-else>
				<tr>
					<td class="text-center">ไม่พบข้อมูล</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped lang="scss">
@media (max-width: 767px) {
	.customize-basic-table {
		width: max-content;
	}
}
$border-color: var(--kt-gray-300);
@mixin border {
	font-weight: 400 !important;
	background-color: transparent !important;
}

th[rowspan] {
	vertical-align: middle;
}

tr:nth-child(1) th:first-of-type,
tr:nth-child(1) th:last-of-type {
	@include border;
	border-bottom: 1px solid $border-color !important;
}

tbody tr td:nth-of-type(1) {
	@include border;
	border-left: 1px solid $border-color;
}

tr:nth-child(2) th:first-of-type {
	@include border;
	border-left: 1px solid $border-color !important;
}
th.border-bottom {
	@include border;
	border: none;
	border-left: 1px solid $border-color;
	border-bottom: 1px solid $border-color !important;
}

dd,
dd p:last-of-type {
	margin-bottom: 0px;
}

.border-right {
	border-top: 1px solid var(--kt-gray-300) !important;
}

.border-top {
	border-top: 1px solid var(--kt-gray-300) !important;
}

// .customize-basic-table tbody tr:not(:last-child):hover td {
// 	background-color: #d9d9d9 !important;
// }

tbody tr td:nth-child(2) {
	border-left: 1px solid var(--kt-gray-300) !important;
}

tr > td:first-child {
	// border-right: 1px solid $border-color !important;
	border-left: transparent !important;
}

tr:hover td {
	background-color: #f4f4f4 !important;
}

// tr:hover td:first-child:hover,
// tr:hover td:first-child:hover ~ td {
// 	background-color: transparent;
// }
</style>
