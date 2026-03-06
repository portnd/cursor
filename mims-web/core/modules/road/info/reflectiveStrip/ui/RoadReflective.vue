<script setup lang="ts">
import { useReflectiveStore } from "../store"
// import RoadReason from "../../reason/ui"
import { RoadReflectiveIRI, RoadReflectiveChart, RoadReflectiveCreate, RoadReflectiveEdit } from "./index"
import RoadMenu from "~~/core/modules/road/info/menu/ui"

import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()
const roadTitleStore = useRoadTitleStore()

const store = useReflectiveStore()

useStoreLifecycle(store, { resetOnEnter: false })

const route = useRoute()
const id = Number(route.params.roadId)

onBeforeMount(async () => {
	// getReflectivityList จะเรียก setDefault() ภายในอัตโนมัติหลัง reflectList ถูก set
	await store.getReflectivityList(id)

	await store
		.getReflectivityDetails()
		.then(() => store.createCriteriaCheckbox())
		.then(() => store.createDataTable())

	await store.getLineList(id)
})

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

const modalCreate: Ref = ref()
const modalEdit: Ref = ref()
// const modalReason = ref()

// ความสูงของ div
const divRoadMenu = ref()
const divRoadMenuHeight = ref(0)

const divCondition = ref()
const divConditionHeight = ref(0)

// เพิ่มข้อมูล
const createItem = () => {
	modalCreate.value.showModal()
}

// แก้ไขข้อมูล
const editItem = () => {
	modalEdit.value.showModal(store.params.id_parent)
}

// const showReason = () => {
// 	modalReason.value.showModal()
// }

const onDelete = () => {
	const data = store.reflectList
	const year = data.find((item) => item.year === store.params.year)
	const dateSurveyed = year?.items.find((item) => item.id_parent === store.params.id_parent)

	if (dateSurveyed?.surveyed_date) {
		useDeleteItem({
			name: `${dateSurveyed.line_no} (สำรวจ: ${buddhistFormatDate(dateSurveyed.surveyed_date, "dd mmm yy")})`,
			url: `roads/retro_reflectivity/${store.params.id_parent}`,
			callBack: async () => {
				await store.callBackUpdateData(id, "delete")
			},
		})
	}
}

const notFoundHeight = ref("")
const handleDivHeight = () => {
	divRoadMenuHeight.value = divRoadMenu.value?.offsetHeight
	divConditionHeight.value = divCondition.value?.offsetHeight

	// กรณี ซ่อนแผนที่
	if (mapShow.value.collapsed) {
		if (divRoadMenuHeight.value + divConditionHeight.value === 0) {
			notFoundHeight.value = `59%`
		} else {
			notFoundHeight.value = `calc(96% - ${divRoadMenuHeight.value + divConditionHeight.value + 4}px)`
		}
	} else {
		notFoundHeight.value = `60vh`
	}
}

const canEdit = ref<boolean>()

watchEffect(() => {
	setTimeout(() => {
		handleDivHeight()
	}, 1500)
})

onMounted(() => {
	canEdit.value =
		initUserStore.accessPermissions[IUserRolesAccess.manage_road_retro] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_retro],
			roadTitleStore.data.ref_depot.id
		)
	handleDivHeight()

	setTimeout(() => {
		handleDivHeight()
	}, 1500)
})

</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<div ref="divRoadMenu">
				<RoadMenu />
			</div>

			<div v-show="!store.reflectList.length" ref="divCondition" class="card card-rounded p-5 mt-5 pt-3">
				<div class="row">
					<div class="col-6 col-md-2">
						<VSelect :options="[]" label="ปี" name="year" placeholder="เลือก" />
					</div>
					<div class="col-6 col-md-4">
						<VSelect :options="[]" label="ช่องจราจร" name="lane" placeholder="เลือก" />
					</div>
					<div class="col-md-6 col-12 mt-md-0 mt-2 align-self-end text-end mb-1">
						<a
							v-show="canEdit"
							type="button"
							class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6 lh-xl"
							@click="createItem"
						>
							เพิ่มข้อมูล
						</a>
					</div>
				</div>
			</div>

			<!-- begin::กรองข้อมูล -->
			<VSkeletonLoader :loading="store.loading">
				<div v-show="store.reflectList.length > 0" class="card card-rounded p-5 mt-5 pt-3">
					<div class="row">
						<div class="col-6 col-md-2">
							<VSelect
								v-model="store.params.year"
								:options="store.getYearOptions"
								label="ปี"
								name="year"
								placeholder="เลือก"
								:can-clear="false"
								:can-deselect="false"
								@update:model-value="() => store.onUpdateYear()"
							/>
						</div>
						<div class="col-6 col-md-4">
							<VSelect
								v-model="store.params.id_parent"
								:options="store.getLineOptions"
								label="เส้นจราจร"
								name="lane"
								placeholder="เลือก"
								:can-clear="false"
								:can-deselect="false"
								@update:model-value="() => store.onUpdateIdParent()"
							/>
						</div>
						<div class="col-md col-12 mt-md-0 mt-2 align-self-end text-end mb-1">
							<a
								v-show="canEdit"
								type="button"
								class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6 lh-xl"
								@click="createItem"
							>
								เพิ่มข้อมูล
							</a>
						</div>
					</div>
					<div class="row mt-3">
						<div class="col-md-8 col-12">
							<label class="text-gray-900 mx-2 fs-7 mt-2">ปรับปรุงข้อมูลโดย</label>
							<VUser
								:label="store.getUserUpdated?.fullname"
								:name="store.getUserUpdated?.fullname"
								:role="store.getUserUpdated?.department"
								:image="store.getUserUpdated?.img_path"
							/>
							<label class="text-gray-900 ms-2 fs-7 mt-2">เมื่อวันที่ {{ store.getUserUpdated?.date }}</label>
						</div>
						<div class="col-md-4 col-12 mt-md-0 mt-2 text-end">
							<a
								v-show="canEdit"
								type="button"
								class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 mb-3 fw-semibold fs-6 lh-xxl"
								@click="editItem"
							>
								ปรับปรุงข้อมูล
							</a>
							<button
								v-show="canEdit"
								type="button"
								class="btn btn-light rounded-4 px-7 py-2 ms-3 mb-3 fw-semibold text-black fs-6"
								@click="onDelete"
							>
								ลบข้อมูล
							</button>
							<!-- <button
								type="button"
								class="btn btn-outline btn-outline-primary rounded-4 px-4 py-3 ms-3 fw-semibold fs-6"
								@click="showReason"
							>
								เหตุผลการส่งกลับแก้ไข
							</button> -->
						</div>
					</div>
				</div>
			</VSkeletonLoader>
			<!-- end::กรองข้อมูล -->

			<!-- begin::กราฟ -->

			<div v-show="store.reflectList.length > 0" class="card card-rounded p-5 mt-3">
				<RoadReflectiveChart :collapsed="mapShow.collapsed" :map-show="mapShow.collapsed" />
			</div>
			<!-- end::กราฟ -->

			<!-- begin::ค่า IRI -->
			<VSkeletonLoader :loading="store.loading">
				<template v-if="store.reflectList.length > 0">
					<RoadReflectiveIRI />
				</template>
			</VSkeletonLoader>
			<!-- end::ค่า IRI -->

			<!-- ไม่พบข้อมูล -->
			<VNotFound v-show="!store.reflectList.length" :height="notFoundHeight" class="mt-5" />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="store.loading" height="97vh" :is-sticky="true" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>

	<!-- Modal -->
	<RoadReflectiveCreate ref="modalCreate" />
	<RoadReflectiveEdit ref="modalEdit" />
	<!-- <RoadReason ref="modalReason" :message="graphStore.graphData.reject_reason" /> -->
</template>

<style scoped>
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}
</style>
