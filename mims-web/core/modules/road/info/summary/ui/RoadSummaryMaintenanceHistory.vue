<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import RoadSummaryMaintenanceHistoryTable from "./RoadSummaryMaintenanceHistoryTable.vue"

const store = useRoadSummaryStore()

const route = useRoute()
const id = Number(route.params.roadId)

onMounted(() => store.getMaintenanceYears(id))
</script>

<template>
	<div class="row mb-3 pt-0">
		<div class="col-4 mb-2">
			<VLabel label="ปีงบประมาณ" class="mt-0" />
			<VSelect
				v-if="store.getYearsOptions.length > 0"
				v-model="store.yearParams"
				:can-clear="false"
				:can-deselect="false"
				name="year"
				:options="store.getYearsOptions"
				placeholder="เลือก"
				@update:model-value="(e: any) => store.getMaintenanceProjects(id, e)"
			/>
			<p v-else class="text-gray-600 mb-0 mt-1 small">
				ไม่มีปีงบประมาณสำหรับสายทางนี้
			</p>
		</div>
		<template v-if="store.maintenanceProjects.length > 0">
			<div
				v-for="(item, index) in store.maintenanceProjects"
				v-show="store.maintenanceProjects.length > 0"
				:key="`project_accordion${index}`"
				class="accordion"
			>
				<div class="accordion-item card shadow-none pb-4 mt-4">
					<h3 :id="`#project_accordion${index}`" class="accordion-header">
						<button
							class="accordion-button pt-4 pb-0 fw-normal text-gray-900 bg-white fs-6 shadow-none"
							:class="{ collapsed: index > 0 }"
							type="button"
							data-bs-toggle="collapse"
							:data-bs-target="`#project_body${index}`"
							:aria-controls="`project_body${index}`"
						>
							ชื่อโครงการ
						</button>
					</h3>
					<div class="row accordion-body py-0">
						<div class="col-12">
							<p class="text-gray-600">
								{{ item.name }}
							</p>
						</div>
						<div class="col-sm-4 col-12">
							<p class="text-gray-900 mb-0">ปีงบประมาณ</p>
							<p class="text-gray-600">
								{{ item.budget_year + 543 }}
							</p>
						</div>
						<!-- <div class="col-sm-4 col-12">
                <p class="text-gray-900 mb-0">งบประมาณการซ่อมบำรุง</p>
                <p class="text-gray-600">{{ toNumber(item.budget_maintenance) }} บาท</p>
              </div> -->
						<div class="col-sm-4 col-12">
							<p class="text-gray-900 mb-0">ประเภทงบประมาณ</p>
							<p class="text-gray-600">{{ item.budget?.name }}</p>
						</div>
					</div>
					<div class="row accordion-body py-0">
						<div class="col-12">
							<p class="text-gray-900 mb-0">ประเภทการซ่อมบำรุง</p>
							<p class="text-gray-600">{{ item.budget_method?.method_name }}</p>
						</div>
						<div class="col-sm-4 col-12">
							<p class="text-gray-900 mb-0">วันที่ตรวจรับงานงวดสุดท้าย</p>
							<p class="mb-0 text-gray-600">
								{{ buddhistFormatDate(item.project_end_date, "dd mmm yyyy") }}
							</p>
						</div>
						<div class="col-sm-4 col-12">
							<p class="text-gray-900 mb-0">วันที่หมดการค้ำประกัน</p>
							<p class="mb-0" :style="{ color: `${item.color} !important` }">
								{{ buddhistFormatDate(item.guarantee_expiration_date, "dd mmm yyyy") }}
							</p>
						</div>
					</div>
					<div
						:id="`project_body${index}`"
						:class="{ show: index <= 0 }"
						class="accordion-collapse collapse"
						:aria-labelledby="`#project_accordion${index}`"
						data-bs-parent="#project_accordion"
					>
						<div class="accordion-body">
							<RoadSummaryMaintenanceHistoryTable
								:maintenance-data="item.roads"
								:maintenance-history-data="item.road_histories"
							/>
						</div>
					</div>
				</div>
			</div>
		</template>
		<VNotFound v-else class="mt-5" :is-not-shadow="true" height="50dvh" />
	</div>
</template>

<style scoped>
p {
	font-size: 1.025rem;
	margin-top: 0.25rem;
}
.accordion-item {
	border: 1px solid var(--kt-gray-300);
}
</style>
