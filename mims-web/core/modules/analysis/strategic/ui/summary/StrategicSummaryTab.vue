<script setup lang="ts">
import { useStrategicAnalysisDashboardStore } from "../../store"
import StrategicSummaryTable from "./StrategicSummaryTable.vue"
import StrategicSummaryPlanTable from "./StrategicSummaryPlanTable.vue"
const store = useStrategicAnalysisDashboardStore()
</script>

<template>
	<div class="mt-0">
		<ul class="nav nav-tabs nav-line-tabs mb-5">
			<li
				class="nav-item active"
				data-bs-toggle="tab"
				data-bs-target="#detail-summary"
				role="tab"
				aria-selected="true"
				@click="() => store.togglePlanTable('สรุปรวม')"
			>
				<span class="nav-link cursor-pointer">สรุปรวม</span>
				<span class="line"></span>
			</li>
			<li
				v-for="i of store.data?.number_plan"
				:key="i"
				class="nav-item"
				data-bs-toggle="tab"
				:data-bs-target="`#detail-plan-${i}`"
				role="tab"
				aria-selected="false"
				@click="() => store.togglePlanTable(`แผนที่ ${i}`)"
			>
				<span class="nav-link cursor-pointer">แผนที่ {{ i }}</span>
				<span class="line"></span>
			</li>
			<li
				v-show="store.data?.table?.unlimited_plan.length > 0"
				class="nav-item"
				data-bs-toggle="tab"
				data-bs-target="#detail-plan-unlimited"
				role="tab"
				aria-selected="false"
				@click="() => store.togglePlanTable('ไม่จำกัดงบประมาณ')"
			>
				<span class="nav-link cursor-pointer">ไม่จำกัดงบประมาณ</span>
				<span class="line"></span>
			</li>
		</ul>
	</div>

	<div class="col-12 order-2">
		<!-- begin::Content -->
		<div class="tab-content p-2">
			<div v-show="!store.loading" id="detail-summary" class="tab-pane fade active show" role="tabpanel">
				<StrategicSummaryTable />
			</div>
			<div id="detail-plan-1" class="tab-pane fade" role="tabpanel">
				<StrategicSummaryPlanTable />
			</div>
			<div id="detail-plan-2" class="tab-pane fade" role="tabpanel">
				<StrategicSummaryPlanTable />
			</div>
			<div id="detail-plan-3" class="tab-pane fade" role="tabpanel">
				<StrategicSummaryPlanTable />
			</div>
			<div id="detail-plan-unlimited" class="tab-pane fade" role="tabpanel">
				<StrategicSummaryPlanTable />
			</div>
		</div>
		<!-- end::Content -->
	</div>
</template>

<style scoped lang="scss">
#nav-tab .active {
	background-color: #fff0d9 !important;
	color: #fdb833 !important;
}

#nav-tab .nav-link {
	background: var(--kt-gray-100);
}

.nav-line-tabs .nav-item .nav-link {
	padding: 10px 22px;
}
</style>
