<script setup lang="ts">
import { useAnnualSummaryDashboardStore } from "../../store"

const store = useAnnualSummaryDashboardStore()

const items = [
	{
		name: "SS : Fibro Seal",
		workload: 3675,
		budget: 0,
		range: 0.73,
	},
	{
		name: "SS : Para Slurry Seal",
		workload: 55527.5,
		budget: 9439675,
		range: 0.15,
	},
	{
		name: "M&OL : PMA",
		workload: 92890,
		budget: 102179000,
		range: 5.18,
	},
	{
		name: "M&OL : AC60/70",
		workload: 22365,
		budget: 4473000,
		range: 6.39,
	},
]
</script>

<template>
	<div class="table-responsive">
		<table class="table customize-basic-table mb-0">
			<thead>
				<tr>
					<th rowspan="2" class="text-center">วิธีการซ่อม</th>
					<th colspan="3" class="text-center border-bottom border-1">ปี 2567</th>
				</tr>
				<tr>
					<th class="text-center">ปริมาณงาน (ตร.ม.)</th>
					<th class="text-center">ค่าซ่อมบำรุง (บาท)</th>
					<th class="text-center">ระยะทาง (กม.)</th>
				</tr>
			</thead>
			<tbody v-if="store.getTable2?.length">
				<tr v-for="(item, index) of store.getTable2" :key="index">
					<td class="text-center">{{ item.name }}</td>
					<td class="text-end">{{ toNumber(item.aera, 2) }}</td>
					<td class="text-end">{{ toNumber(item.iri_after, 2) }}</td>
					<td class="text-end">{{ toNumber(item.range, 2) }}</td>
				</tr>
				<tr v-show="items.length" class="text-end border-top data-row">
					<td class="text-center">รวม</td>
					<td class="text-end">
						{{
							toNumber(
								store.getTable2?.reduce((acc: number, item) => acc + item.aera, 0),
								2
							)
						}}
					</td>
					<td class="text-end">
						{{
							toNumber(
								store.getTable2?.reduce((acc: number, item) => acc + item.iri_after, 0),
								2
							)
						}}
					</td>
					<td class="text-end">
						{{
							toNumber(
								store.getTable2?.reduce((acc: number, item) => acc + item.range, 0),
								2
							)
						}}
					</td>
				</tr>
			</tbody>
			<tbody v-else>
				<tr>
					<td colspan="4">ไม่พบข้อมูล</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped lang="scss">
// @media (max-width: 767px) {
// 	.customize-basic-table {
// 		width: max-content;
// 	}
// }
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
	border-left: solid $border-color;
}

tr > td:first-child {
	// border-right: 1px solid $border-color !important;
	border-left: transparent !important;
}

dd,
dd p:last-of-type {
	margin-bottom: 0px;
}

.border-top {
	border-top: 1px solid var(--kt-gray-300) !important;
}

tr:hover td {
	background-color: #f4f4f4 !important;
}

// tr:hover td:first-child:hover,
// tr:hover td:first-child:hover ~ td {
// 	background-color: transparent;
// }

tbody tr td:nth-child(2) {
	border-left: 1px solid var(--kt-gray-300) !important;
}
</style>
