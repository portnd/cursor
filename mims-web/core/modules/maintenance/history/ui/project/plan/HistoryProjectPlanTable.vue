<script setup lang="ts">
import { useMaintenanceHistoryPlanStore } from "../../../store/MaintenanceHistoryPlanStore"

const store = useMaintenanceHistoryPlanStore()
const canEdit = ref<boolean>()

onMounted(() => {
	canEdit.value = usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.manage_all_maint_history)
})
</script>

<template>
	<div class="row px-3">
		<div v-for="(item, index) of store.getDataTable" :key="index" class="col-12 mt-5">
			<h6 class="fw-semibold">{{ item.plan_name }}</h6>
			<div class="table-responsive mt-3">
				<table class="table customize-basic-table mb-0 text-truncate table-hover">
					<thead>
						<tr>
							<th rowspan="2" class="text-center">เดือน/ปี</th>
							<th colspan="4" class="text-center border-bottom border-1">ความก้าวหน้า (%)</th>
							<th colspan="4" class="text-center border-bottom border-1">การเบิกจ่าย (บาท)</th>
						</tr>
						<tr>
							<th class="text-center">แผน</th>
							<th class="text-center">สะสม</th>
							<th class="text-center">ผล</th>
							<th class="text-center">สะสม</th>
							<th class="text-center">แผน</th>
							<th class="text-center">สะสม</th>
							<th class="text-center">ผล</th>
							<th class="text-center">สะสม</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="(child, i) in item.value" :key="i" class="text-end">
							<td class="text-center">{{ buddhistFormatDate(child.schedule, "mmm yy") }}</td>
							<td>{{ store.checkIsInteger(child.plan) }}</td>
							<td>{{ store.checkIsInteger(child.plan_total) }}</td>
							<td>{{ store.checkIsInteger(child.progress_plan) }}</td>
							<td>{{ store.checkIsInteger(child.progress_plan_total) }}</td>
							<td>{{ store.checkIsInteger(child.disbursement_plan) }}</td>
							<td>
								{{ store.checkIsInteger(child.disbursement_plan_total) }}
							</td>
							<td>
								{{ store.checkIsInteger(child.disbursement_progress) }}
							</td>
							<td>
								{{ store.checkIsInteger(child.disbursement_progress_total) }}
							</td>
						</tr>
						<tr class="text-end border-top">
							<td class="text-center">รวม</td>
							<td>
								{{ store.checkIsInteger(item.value.reduce((acc, curr) => acc + (curr.plan || 0), 0)) }}
							</td>
							<td>{{ store.checkIsInteger(item.value[item.value.length - 1].plan_total) }}</td>
							<td>
								{{ store.checkIsInteger(item.value.reduce((acc, curr) => acc + (curr.progress_plan || 0), 0)) }}
							</td>
							<td>{{ store.checkIsInteger(item.value[item.value.length - 1].progress_plan_total) }}</td>
							<td>
								{{ store.checkIsInteger(item.value.reduce((acc, curr) => acc + (curr.disbursement_plan || 0), 0)) }}
							</td>
							<td>{{ store.checkIsInteger(item.value[item.value.length - 1].disbursement_plan_total) }}</td>
							<td>
								{{ store.checkIsInteger(item.value.reduce((acc, curr) => acc + (curr.disbursement_progress || 0), 0)) }}
							</td>
							<td>{{ store.checkIsInteger(item.value[item.value.length - 1].disbursement_progress_total) }}</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<div class="col-12 text-end mt-6">
			<dl class="row text-start">
				<dt class="col-12 mb-3 h6 fw-semibold">ปัญหาและอุปสรรค</dt>
				<dd
					v-for="(prob, probIndex) of store.getProblemDetails"
					v-show="store.getProblemDetails.length > 0"
					:key="probIndex"
					class="col-12 text-preline"
				>
					<li>{{ prob === "" ? null : decodeHTML(prob) }}</li>
				</dd>
				<dd v-show="store.getProblemDetails?.length === 0">ไม่พบข้อมูล</dd>
			</dl>
		</div>
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
	border-left: 1px solid $border-color;
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
</style>
