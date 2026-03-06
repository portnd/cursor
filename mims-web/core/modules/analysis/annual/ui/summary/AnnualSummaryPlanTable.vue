<script setup lang="ts">
import { useAnnualSummaryDashboardStore } from "../../store"

const store = useAnnualSummaryDashboardStore()
</script>

<template>
	<div class="table-responsive">
		<table class="table customize-basic-table mb-0">
			<tbody>
				<tr>
					<td width="150" class="text-start">งบประมาณซ่อมบำรุง</td>
					<td width="250" class="text-end">{{ toNumber(store.getTable1?.budget, 2) || "-" }}</td>
					<td width="80" class="text-end">บาท</td>
				</tr>
				<tr class="text-end border-top">
					<td width="150" class="text-start">IRI (ก่อนวิเคราะห์)</td>
					<td width="250" class="text-end">{{ toNumber(store.getTable1?.iri_before, 2) || "-" }}</td>
					<td width="80" class="text-end">ม./กม.</td>
				</tr>
				<tr class="text-end border-top">
					<td width="150" class="text-start">IRI (หลังวิเคราะห์)</td>
					<td width="250" class="text-end">{{ toNumber(store.getTable1?.iri_after, 2) || "-" }}</td>
					<td width="80" class="text-end">ม./กม.</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped lang="scss">
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

tbody tr td:nth-of-type(2) {
	@include border;
	/* border-left: 1px solid $border-color; */
}

tr:nth-child(2) th:first-of-type {
	@include border;
	border-left: 1px solid $border-color !important;
}
th.border-bottom {
	@include border;
	border-bottom: 1px solid $border-color !important;
}

dd,
dd p:last-of-type {
	margin-bottom: 0px;
}

.border-top {
	border-top: 1px solid var(--kt-gray-300) !important;
}

tbody tr:hover td {
	background-color: #f4f4f4 !important;
}

// tr:hover td:first-child:hover,
// tr:hover td:first-child:hover ~ td {
// 	background-color: transparent !important;
// }
</style>
