<script setup lang="ts">
// import { api as viewerApi } from "v-viewer"
import { useMaintenanceHistoryDetailsStore } from "../../../store/MaintenanceHistoryDetailsStore"
// import { IMaintenanceHistoryDetailAttrachment, IMaintenanceHistoryDetailMaintenanceRoadHistory } from '../../../infrastructure/MaintenanceHistoryModel';
import { HistoryProjectInfoModalRoads, HistoryProjectInfoModalEdit } from "../.."
import {
	IMaintenanceHistoryDetailMaintenanceRoad,
	IMaintenanceHistoryDetailMaintenanceRoadHistory,
} from "../../../infrastructure"
import BtnCreate from "~/components/buttons/BtnCreate.vue"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const route = useRoute()
const id = Number(route.params.id)
const store = useMaintenanceHistoryDetailsStore()
const canEdit = ref<boolean>()

// const data = [
// 	{
// 		id: 1,
// 		road_info: {
// 			name: "-",
// 		},
// 		lane: "1",
// 		intervention_criteria: {
// 			id: 1,
// 			maintenance_standard_name: "-",
// 		},
// 		km_start: 0,
// 		km_end: 1000,
// 		the_geom: "",
// 		attachments: [],
// 	},
// ]

onMounted(async () => {
	await store.getMaintenanceHistoryDetail(id)
	await store.getMaintenanceBudgets().then(() => store.checkIsShowMethod(store.data.budget_method.id))

	canEdit.value =
		initUserStore.accessPermissions[IUserRolesAccess.manage_all_maint_history] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_maint_history],
			store.data.ref_depot.id
		)
})

// Modal
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

// เพิ่มข้อมูล
const createItem = (
	item: IMaintenanceHistoryDetailMaintenanceRoadHistory | IMaintenanceHistoryDetailMaintenanceRoad
) => {
	modalCreate.value.showModal(item)
}

const editInfroItem = (itemId: number) => {
	navigateTo(`/maintenances/history/${id}/maintenance-roads/${itemId}/edit`)
	// modalEdit.value.showModal(itemId, store.getRoadGroupId, store.data.budget_method.id)
}

const editHistItem = (itemId: number) => {
	navigateTo(`/maintenances/history/${id}/warranty/${itemId}/edit`)
}

// const onDelete = () => {
// 	const data = store.conditionList
// 	const year = data.find((item) => item.year === store.params.year)
// 	const dateSurveyed = year?.items.find((item) => item.id_parent === store.params.id_parent)

// 	useDeleteItem({
// 		name: `${dateSurveyed?.lane_no} (สำรวจ: ${buddhistFormatDate(dateSurveyed?.surveyed_date, "dd mmm yy")})`,
// 		url: `roads/condition/${store.params.id_parent}`,
// 		callBack: async () => {
// 			await store.callBackUpdateData(id, "delete")
// 		},
// 	})
// }

const onDeleteRoads = (item: IMaintenanceHistoryDetailMaintenanceRoad) => {
	useDeleteItem({
		name: item.road_name,
		url: `maintenance/${id}/road/${item.id}`,
		callBack: async () => {
			await store.getMaintenanceHistoryDetail(id)
			await store.getMaintenanceBudgets().then(() => store.checkIsShowMethod(store.data.budget_method.id))
		},
	})
}
const onDeleteWarranty = (item: IMaintenanceHistoryDetailMaintenanceRoadHistory) => {
	useDeleteItem({
		name: item.road_name,
		url: `maintenance/${id}/road_history/${item.id}`,
		callBack: async () => {
			await store.getMaintenanceHistoryDetail(id)
			await store.getMaintenanceBudgets().then(() => store.checkIsShowMethod(store.data.budget_method.id))
		},
	})
}

onUnmounted(() => {
	// Do NOT dispose/reset here: nav from list→info causes spurious unmount before first fetch
	// completes; details store is reset when entering list page instead.
})
</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="row">
			<div class="col-12">
				<div class="card p-5">
					<div class="row">
						<div class="col-5 col-md-8 ps-4 d-flex align-items-center">
							<h5 class="fw-semibold fs-5 mt-6 mt-md-2">ข้อมูลโครงการ</h5>
							<!-- <i>
                <label class="text-gray-600 me-2 fs-7 mt-2">ปรับปรุงข้อมูลโดย</label>
                <VUser :label="store.getUpdater.name" :name="store.getUpdater.name" :role="store.getUpdater.department" />
                <label class="text-gray-600 mx-0 mx-sm-2 fs-7"> {{ store.getUpdater.date }} </label>
              </i> -->
						</div>
						<div class="col-7 col-md-4 text-end">
							<NuxtLink
								v-if="canEdit"
								:to="`/maintenances/history/${id}/edit`"
								class="btn btn-outline btn-outline-primary rounded-2 me-3 fw-semibold"
								style="padding-top: 0.8em"
							>
								แก้ไขโครงการ
							</NuxtLink>
						</div>
					</div>

					<div class="row">
						<div class="col-12 col-md-7 ps-4 order-1 order-md-0">
							<dl class="row mt-4 mt-md-4">
								<dt class="col-sm-6">ชื่อโครงการ :</dt>
								<dd class="col-sm-6">
									{{ store.data.name }}
								</dd>

								<dt class="col-sm-6">หน่วยงาน :</dt>
								<dd class="col-sm-6">
									{{ store.data.owner_name }}
								</dd>

								<dt class="col-sm-6">เลขที่สัญญา :</dt>
								<dd class="col-sm-6">
									{{ store.data.contract_number }}
								</dd>

								<dt class="col-sm-6">ปีงบประมาณ :</dt>
								<dd class="col-sm-6">
									{{ store.data.budget_year + 543 }}
								</dd>

								<dt class="col-sm-6">ประเภทงบประมาณ :</dt>
								<dd class="col-sm-6">
									{{ store.getBudgetName }}
								</dd>

								<dt class="col-sm-6">วิธีการซ่อมบำรุง :</dt>
								<dd class="col-sm-6">
									{{ store.getBudgetMethodName }}
								</dd>

								<dt class="col-sm-6">วงเงินงบประมาณ (ไม่รวม VAT) :</dt>
								<dd class="col-sm-6">
									{{ store.data.budget_maintenance ? toNumber(store.data.budget_maintenance) + " บาท" : "-" }}
								</dd>

								<dt class="col-sm-6">ราคากลาง (รวม VAT) :</dt>
								<dd class="col-sm-6">{{ toNumber(store.data.middle_price) }} บาท</dd>

								<dt class="col-sm-6">มูลค่างานตามสัญญา (รวม VAT) :</dt>
								<dd class="col-sm-6">{{ toNumber(store.data.contract_work_value) }} บาท</dd>

								<dt class="col-sm-6">ราคาที่จัดซื้อจัดจ้าง :</dt>
								<dd class="col-sm-6">{{ toNumber(store.data.budget_procurement) }} บาท</dd>

								<dt class="col-sm-6">บริษัทที่ปรึกษาโครงการ :</dt>
								<dd class="col-sm-6">{{ store.data.advisor_name }}</dd>

								<dt class="col-sm-6">บริษัทผู้รับจ้าง :</dt>
								<dd class="col-sm-6">{{ store.data.contractor_name }}</dd>

								<dt class="col-sm-6">ชื่อ-นามสกุลเจ้าหน้าที่เลขาโครงการ :</dt>
								<dd class="col-sm-6">{{ store.data.project_secretary_name }}</dd>

								<dt class="col-sm-6 fw-normal">รายละเอียดโครงการ :</dt>
								<dd class="col-sm-6">{{ store.data.project_details }}</dd>

								<!-- <dt class="col-sm-6 fw-normal">วันที่ตรวจรับงานงวดสุดท้าย</dt>
                <dd class="col-sm-6 text-preline">{วันที่ตรวจรับงานงวดสุดท้าย}</dd>
                <dt class="col-sm-6 fw-normal">วันที่หมดการค้ำประกันโครงการ</dt>
                <dd class="col-sm-6 text-preline">{วันที่หมดการค้ำประกันโครงการ}</dd> -->
							</dl>
						</div>
						<!-- <div class="col-12 col-md-5 order-0 order-md-1">
                <div class="col-12">
                  <HistoryProjectInfoStatus />
                </div>
              </div> -->
					</div>
					<div class="row">
						<div class="row row-cols-12">
							<dt class="col-12 col-md-3 fw-normal me-md-10 mb-md-0 mb-4">เอกสารแนบ :</dt>
							<dd class="col-12 col-md-7 mt-md-0 mt-4">
								<ul v-if="store.data.attachments?.length" class="col-12 p-0 mb-0 d-flex flex-wrap ps-0 gap-10 gap-md-8">
									<li
										v-for="(item, index) in store.data.attachments"
										:key="index"
										style="list-style-type: none"
										class="text-primary cursor-pointer mb-3 ps-1 d-flex align-items-end"
									>
										<NuxtLink
											v-if="['png', 'jpg', 'jpeg'].includes(getFileExtension(item.file_name))"
											class="d-flex flex-column align-items-start"
										>
											<img v-viewer :src="item?.path" class="cursor-pointer" width="65" height="65" />
											<span class="filename cursor-pointer mb-3" :title="item?.file_name">{{
												textTruncate(item.file_name, 10)
											}}</span>
										</NuxtLink>
										<NuxtLink
											v-else-if="getFileExtension(item.file_name) === 'pdf'"
											class="d-flex flex-column align-items-center"
											:to="item.path"
											target="_blank"
										>
											<img :src="`/images/files/pdf.png`" width="50" />
											<span class="filename cursor-pointer mb-3" :title="item?.file_name">{{
												textTruncate(item.file_name, 10)
											}}</span>
										</NuxtLink>
										<NuxtLink
											v-else-if="getFileExtension(item.file_name) === 'docx'"
											class="d-flex flex-column align-items-center"
											:to="item.path"
											target="_blank"
											download
										>
											<img :src="`/images/files/docx.png`" width="50" />
											<span class="filename cursor-pointer mb-3" :title="item?.file_name">{{
												textTruncate(item.file_name, 10)
											}}</span>
										</NuxtLink>
										<NuxtLink
											v-else-if="getFileExtension(item.file_name) === 'xlsx'"
											class="d-flex flex-column align-items-center"
											:to="item.path"
											target="_blank"
											download
										>
											<img :src="`/images/files/xlsx.png`" width="50" />
											<span class="filename cursor-pointer mb-3" :title="item?.file_name">{{
												textTruncate(item.file_name, 10)
											}}</span>
										</NuxtLink>
									</li>
								</ul>
								<p v-else class="text-gray-600">ไม่พบข้อมูล</p>
							</dd>
						</div>
					</div>
					<div class="row">
						<div class="col-12 mt-3">
							<div class="row align-items-center">
								<div class="col-6">
									<h5 class="fw-semibold">ข้อมูลการซ่อมบำรุง</h5>
								</div>
								<div class="col-6 text-end">
									<BtnCreate
										v-if="canEdit"
										label="เพิ่ม"
										@click="navigateTo(`/maintenances/history/${id}/maintenance-roads/create`)"
									/>
								</div>
							</div>
							<div class="table-responsive mt-3">
								<table class="table customize-basic-table mb-0 text-truncate table-hover">
									<thead>
										<tr>
											<th class="text-center">ลำดับ</th>
											<th class="text-center">สายทาง</th>
											<th class="text-center">จาก - ถึง</th>
											<th class="text-center">ช่องจราจร</th>
											<th v-show="store.data?.is_show_method" class="text-center">วิธีการซ่อมบำรุง</th>
											<th class="text-center">ช่วง กม.</th>
											<th class="text-center">ระยะทาง (กม.)</th>
											<th class="text-center">จัดการ</th>
										</tr>
									</thead>
									<tbody>
										<tr
											v-for="(item, index) of store.getMaintenanceHistoryRoadsTable"
											:key="index"
											class="cursor-pointer cursor_hover"
										>
											<td class="text-center" @click="store.toLocation(item)">{{ index + 1 }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.road_group_name }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.road_name }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.lane_no }}</td>
											<td v-show="store.data?.is_show_method" class="text-center" @click="store.toLocation(item)">
												{{
													item.intervention_criteria.maintenance_standard_name
														? item.intervention_criteria.maintenance_standard_name
														: "-"
												}}
											</td>
											<td class="text-center" @click="store.toLocation(item)">
												{{ convertMeterToKm(item.km_start) }} - {{ convertMeterToKm(item.km_end) }}
											</td>
											<td class="text-center" @click="store.toLocation(item)">
												{{ +(Math.abs(item.km_end - item.km_start) / 1000).toFixed(2) }}
											</td>
											<td class="text-center">
												<a
													type="button"
													class="btn-icon align-middle text-secondary lh-0 me-2 border-0 bg-transparent"
													@click.stop="createItem(item)"
												>
													<i class="fi fi-ss-road fs-3"></i>
												</a>
												<BtnEdit @click.stop="editInfroItem(item.id)" />
												<BtnDelete @click.stop="onDeleteRoads(item)" />
											</td>
										</tr>
										<tr v-show="store.getMaintenanceHistoryRoadsTable.length === 0" class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
									<!-- <tfoot v-show="store.getMaintenanceHistoryRoadsTable.length !== 0">
                    <tr class="border">
                      <td :colspan="store.isShowMethod ? 4 : 3"></td>
                      <td class="text-center">ระยะทาง</td>
                      <td class="text-center">{{ store.getSumMaintenanceRoadsTable }}</td>
                    </tr>
                  </tfoot> -->
								</table>
							</div>
						</div>
					</div>
					<div class="row mt-5">
						<div class="col-12 mt-5">
							<div class="row align-items-center">
								<div class="col-6">
									<h5 class="fw-semibold">ประวัติการซ่อมบำรุง ในระยะเวลารับประกันผลงาน</h5>
								</div>
								<div class="col-6 text-end">
									<BtnCreate
										v-if="canEdit"
										label="เพิ่ม"
										@click="navigateTo(`/maintenances/history/${id}/warranty/create`)"
									/>
								</div>
							</div>
							<div class="table-responsive mt-3">
								<table class="table customize-basic-table mb-0 text-truncate table-hover">
									<thead>
										<tr>
											<th class="text-center">ลำดับ</th>
											<th class="text-center">สายทาง</th>
											<th class="text-center">จาก - ถึง</th>
											<th class="text-center">ช่องจราจร</th>
											<th v-show="store.data?.is_show_method" class="text-center">วิธีการซ่อมบำรุง</th>
											<th class="text-center">ช่วง กม.</th>
											<th class="text-center">ระยะทางรวม (กม.)</th>
											<th class="text-center">จัดการ</th>
										</tr>
									</thead>
									<tbody>
										<tr
											v-for="(item, index) of store.getMaintenanceHistoryDataTable"
											:key="index"
											class="align-middle cursor-pointer"
										>
											<td class="text-center" @click="store.toLocation(item)">{{ index + 1 }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.road_group_name }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.road_name }}</td>
											<td class="text-center" @click="store.toLocation(item)">{{ item.lane_no }}</td>
											<td v-show="store.data?.is_show_method" class="text-center" @click="store.toLocation(item)">
												{{
													item.intervention_criteria.maintenance_standard_name
														? item.intervention_criteria.maintenance_standard_name
														: "-"
												}}
											</td>
											<td class="text-center" @click="store.toLocation(item)">
												{{ convertMeterToKm(item.km_start) }} - {{ convertMeterToKm(item.km_end) }}
											</td>
											<td class="text-center" @click="store.toLocation(item)">
												{{ +(Math.abs(item.km_end - item.km_start) / 1000).toFixed(2) }}
											</td>
											<td class="text-center">
												<a
													type="button"
													class="btn-icon align-middle text-secondary lh-0 me-2 border-0 bg-transparent"
													@click.stop="createItem(item)"
												>
													<i class="fi fi-ss-road fs-3"></i>
												</a>
												<BtnEdit @click.stop="editHistItem(item.id)" />
												<BtnDelete @click.stop="onDeleteWarranty(item)" />
											</td>
										</tr>
										<tr v-show="store.getMaintenanceHistoryDataTable.length === 0" class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</VSkeletonLoader>
	<HistoryProjectInfoModalRoads ref="modalCreate" />
	<HistoryProjectInfoModalEdit ref="modalEdit" />
</template>

<style scoped></style>
