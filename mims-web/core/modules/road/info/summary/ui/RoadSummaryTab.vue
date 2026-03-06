<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import RoadReason from "../../reason/ui"
import RoadSummaryTrafficVolumnCreateModal from "./RoadSummaryTrafficVolumnCreateModal.vue"
import RoadSummaryTrafficVolumnEditModal from "./RoadSummaryTrafficVolumnEditModal.vue"
import {
	RoadSummaryRoadDetail,
	RoadSummaryTrafficVolumn,
	RoadSummaryPavement,
	RoadSummarySurface,
	RoadSummaryMaintenanceHistory,
	RoadDetailModal,
} from "./index"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import { useInitUserStore } from "~/core/modules/initUser/store"

const roadTitleStore = useRoadTitleStore()

const initUserStore = useInitUserStore()

const store = useRoadSummaryStore()

useStoreLifecycle(store, { resetOnEnter: false })

const route = useRoute()
const roadId = route.params.roadId

const isShow = ref(true)
const activeTab = ref("road")
const canEdit = ref<boolean>()
const trafficPage: Ref = ref()
const modalTrafficCreate: Ref = ref()
const modalTrafficEdit: Ref = ref()
const modalEdit: Ref = ref()
const roadLevel = ref(1)

const changeTab = async (tab: string) => {
	isShow.value = true
	activeTab.value = tab
	if (store.map) {
		if (tab === "condition") {
			isShow.value = false
		} else if (tab === "road") {
			store.getRoadDetail(Number(roadId))
			store.map.Overlays.clear()
			store.createLine()
		} else if (tab === "surface" || tab === "pavement") {
			store.map.Overlays.clear()
			store.createLine()
		} else if (tab === "maintenance-history") {
			isShow.value = false
			store.map.Overlays.clear()
			store.createMaintenanceHistoryLine()
		} else if (tab === "traffic") {
			await store.getTrafficRevision(Number(roadId)).then(() => store.setDefaultOptions())
			if (store.toggle.aadtId) {
				await store.getTrafficDetail(store.toggle.aadtId)
			}
			store.map.Overlays.clear()
			store.createTrafficLine()
		} else {
			store.map.Overlays.clear()
			store.createLine()
		}
	}
}

watch(
	() => store.toggle.aadtId,
	async (newAadtId) => {
		if (newAadtId) {
			await store.getTrafficDetail(newAadtId)
		}
	},
	{ flush: 'post' }
)

const stopWatcher = watch(
	() => activeTab.value,
	() => {
		// Hook for potential cleanup
	}
)

onBeforeMount(() => {
	store.$reset()
	trafficPage.value?.updateWidth()
	if (route.query.tab) {
		const currentPath = route.path
		const url = new URL(currentPath, window.location.origin)
		const searchParams = new URLSearchParams(url.search)
		searchParams.delete("tab")
		url.search = searchParams.toString()
		navigateTo(url.pathname + url.search, {
			replace: true,
		})
	}
})
onMounted(() => {
	canEdit.value =
		initUserStore.accessPermissions[IUserRolesAccess.manage_road_damage] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_damage],
			roadTitleStore.data.ref_depot.id
		)
	if (route.query.tab) {
		activeTab.value = route.query.tab.toString()
	}
})

onUnmounted(() => {
	if (stopWatcher) {
		stopWatcher()
	}
})
</script>

<template>
	<div class="row mt-4">
		<div class="col-12">
			<div class="card shadow">
				<div class="card-body p-5">
					<div class="mt-0">
						<ul class="nav nav-tabs nav-line-tabs mb-5">
							<li
								class="nav-item"
								:class="{ active: activeTab === 'road' }"
								data-bs-toggle="tab"
								data-bs-target="#detail-road"
								role="tab"
								aria-selected="true"
								@click="changeTab('road')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลสายทาง</span>
								<span class="line"></span>
							</li>
							<li
								class="nav-item"
								:class="{ active: activeTab === 'surface' }"
								data-bs-toggle="tab"
								data-bs-target="#detail-surface"
								role="tab"
								aria-selected="true"
								@click="changeTab('surface')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลผิวทาง</span>
								<span class="line"></span>
							</li>
							<li
								class="nav-item"
								:class="{ active: activeTab === 'pavement' }"
								data-bs-toggle="tab"
								data-bs-target="#detail-pavement"
								role="tab"
								aria-selected="false"
								@click="changeTab('pavement')"
							>
								<span class="nav-link cursor-pointer">หน้าตัดผิวทาง</span>
								<span class="line"></span>
							</li>
							<!-- <li
								class="nav-item"
								data-bs-toggle="tab"
								data-bs-target="#detail-condition"
								role="tab"
								aria-selected="false"
								@click="changeTab('condition')"
							>
								<span class="nav-link cursor-pointer">สภาพทาง</span>
								<span class="line"></span>
							</li> -->
							<li
								class="nav-item"
								:class="{ active: activeTab === 'maintenance-history' }"
								data-bs-toggle="tab"
								data-bs-target="#maintenance-history"
								role="tab"
								aria-selected="false"
								@click="changeTab('maintenance-history')"
							>
								<span class="nav-link cursor-pointer">ประวัติการซ่อมบำรุง</span>
								<span class="line"></span>
							</li>
							<li
								v-if="
									initUserStore.accessPermissions[IUserRolesAccess.view_road_traffic] ||
									initUserStore.getIsOwnerManagePermission(
										initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_traffic],
										roadTitleStore.data.ref_depot.id
									)
								"
								class="nav-item"
								:class="{ active: activeTab === 'traffic' }"
								data-bs-toggle="tab"
								data-bs-target="#traffic-volumn"
								role="tab"
								aria-selected="false"
								@click="changeTab('traffic')"
							>
								<span class="nav-link cursor-pointer">ปริมาณจราจร</span>
								<span class="line"></span>
							</li>
						</ul>
					</div>
					<!-- <VSkeletonLoader :loading="store.loading">
						<div v-show="activeTab === 'traffic'" class="row px-2 mb-3">
							<template v-if="store.trafficRevision.length === 0">
								<div class="col-6 col-md-3">
									<VSelect :options="[]" label="ปี" name="year" placeholder="เลือก" />
								</div>
								<div class="col-6 col-md-5">
									<VSelect :options="[]" label="ช่องจราจร" name="lane" placeholder="เลือก" />
								</div>
								<div v-show="canEdit" class="col-md-4 col-12 mt-md-0 mt-2 align-self-end text-end mb-1">
									<button
										type="button"
										class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6"
										@click="createItem()"
									>
										เพิ่มข้อมูล
									</button>
								</div>
							</template>
							<template v-else>
								<div class="col-6 col-md-3">
									<VSelect
										v-model="store.toggle.year"
										:options="store.getYearListOptions"
										:can-clear="false"
										:can-deselect="false"
										label="ปี"
										name="year"
										placeholder="เลือก"
									/>
								</div>
								<div class="col-6 col-md-5">
									{{ store.toggle.aadtId }}
									<VSelect
										v-model="store.toggle.aadtId"
										:options="store.getYearItemsOptions"
										:can-clear="false"
										:can-deselect="false"
										label="ช่องจราจร"
										name="lane"
										placeholder="เลือก"
									/>
								</div>
								<div v-show="canEdit" class="col-md-4 col-12 mt-md-0 mt-2 align-self-end text-end mb-1">
									<button
										type="button"
										class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6"
										@click="createItem()"
									>
										เพิ่มข้อมูล
									</button>
								</div>
							</template>
						</div>

						<div
							v-show="isShow"
							class="d-flex justify-content-between align-items-md-center flex-md-row flex-column gap-2 px-2 mb-3"
							:class="activeTab === 'maintenance-history' ? 'd-none' : ''"
						>
							<div class="d-flex align-items-center flex-wrap w-md-75 gap-2 fw-normal fs-5 mb-0">
								<div v-if="activeTab === 'surface' || activeTab === 'pavement'">
									<template v-if="store.update_date !== ''">
										<label class="text-gray-900 fs-6">ปรับปรุงข้อมูลโดย</label>
										<VUser :label="store.update_by.full_name" :name="store.update_by.full_name" :role="'-'" />
										<label class="text-gray-900 fs-6">
											เมื่อวันที่
											{{ buddhistFormatDate(new Date(store.update_date), "dd mmm yyyy เวลา HH:ii น.") }}
										</label>
									</template>
								</div>
								<div v-if="activeTab === 'traffic'">
									<template v-if="store.getUpdateBy?.date !== ''">
										<label class="text-gray-900 fs-6">ปรับปรุงข้อมูลโดย</label>
										<VUser :label="`${store.getUpdateBy?.name ?? ''}`" :name="`${store.getUpdateBy?.name ?? ''}`" :role="'-'" />
										<label class="text-gray-900 fs-6">
											เมื่อวันที่
											{{ store.trafficDetail?.updated_date ? buddhistFormatDate(new Date(store.trafficDetail.updated_date), "dd mmm yyyy เวลา HH:ii น.") : '' }}
										</label>
									</template>
								</div>
								<div v-if="activeTab === 'road'">
									<template v-if="store.road.road_info?.updated_at !== ''">
										<label class="text-gray-900 fs-6">ปรับปรุงข้อมูลโดย</label>
										<VUser
											:label="store.road.road_info?.user.firstname + ' ' + store.road.road_info?.user.lastname"
											:name="store.road.road_info?.user.firstname + ' ' + store.road.road_info?.user.lastname"
											:role="'-'"
										/>
										<label class="text-gray-900 fs-6">
											เมื่อวันที่
											{{ buddhistFormatDate(new Date(store.road.road_info?.updated_at), "dd mmm yyyy เวลา HH:ii น.") }}
										</label>
									</template>
								</div>
							</div>
							<div v-show="canEdit">
								<NuxtLink
									v-if="activeTab === 'road'"
									v-show="store.status_code !== 'W'"
									class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold fs-6 lh-xxl"
									@click="editItem(roadTitleStore.data.road_level)"
								>
									ปรับปรุงข้อมูล
								</NuxtLink>
								<NuxtLink
									v-else-if="activeTab === 'traffic'"
									v-show="store.status_code !== 'W' && store.trafficRevision.length > 0"
									class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold fs-6 lh-xxl"
									@click="editItem(roadTitleStore.data.road_level)"
								>
									ปรับปรุงข้อมูล
								</NuxtLink>
								<NuxtLink
									v-else
									v-show="store.status_code !== 'W'"
									class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold fs-6 lh-xxl"
									@click="editItem(roadTitleStore.data.road_level)"
								>
									ปรับปรุงข้อมูล
								</NuxtLink>
								<button
									v-show="store.status_code === 'R'"
									type="button"
									class="btn btn-outline btn-outline-primary rounded-4 px-3 py-3 ms-3 fw-semibold fs-6"
									@click="showReason"
								>
									เหตุผลการส่งกลับแก้ไข
								</button>
							</div>
						</div>
					</VSkeletonLoader> -->
					<div class="row">
						<div class="col-12 order-2">
							<!-- begin::Content -->
							<div class="tab-content p-2">
								<div
									id="detail-road"
									class="tab-pane fade"
									role="tabpanel"
									:class="{ active: activeTab === 'road', show: activeTab === 'road' }"
								>
									<RoadSummaryRoadDetail />
								</div>
								<div
									id="detail-surface"
									class="tab-pane fade"
									role="tabpanel"
									:class="{ active: activeTab === 'surface', show: activeTab === 'surface' }"
								>
									<RoadSummarySurface :show="activeTab === 'surface'" />
								</div>
								<div
									id="detail-pavement"
									class="tab-pane fade"
									role="tabpanel"
									:class="{ active: activeTab === 'pavement', show: activeTab === 'pavement' }"
								>
									<RoadSummaryPavement />
								</div>
								<!-- <div id="detail-condition" class="tab-pane fade" role="tabpanel">
									<RoadSummaryCondition />
								</div> -->
								<div
									id="maintenance-history"
									class="tab-pane fade"
									role="tabpanel"
									:class="{ active: activeTab === 'maintenance-history', show: activeTab === 'maintenance-history' }"
								>
									<RoadSummaryMaintenanceHistory />
								</div>
								<div
									v-if="
										initUserStore.accessPermissions[IUserRolesAccess.view_road_traffic] ||
										initUserStore.getIsOwnerManagePermission(
											initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_traffic],
											roadTitleStore.data.ref_depot.id
										)
									"
									id="traffic-volumn"
									class="tab-pane fade"
									role="tabpanel"
									:class="{ active: activeTab === 'traffic', show: activeTab === 'traffic' }"
								>
									<RoadSummaryTrafficVolumn ref="trafficPage" />
								</div>
							</div>
							<!-- end::Content -->
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Modal -->
	<RoadReason ref="modalReason" :message="store?.reject_reason" />
	<RoadSummaryTrafficVolumnCreateModal ref="modalTrafficCreate" />
	<RoadSummaryTrafficVolumnEditModal ref="modalTrafficEdit" />
	<RoadDetailModal ref="modalEdit" :road-id="roadLevel" />
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 22px !important;
}
</style>
