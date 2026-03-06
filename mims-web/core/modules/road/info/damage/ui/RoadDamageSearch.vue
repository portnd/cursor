<script setup lang="ts">
import { IDatum } from "../infrastructure/RoadDamageModel"
import { useRoadDamageStore } from "../store/RoadDamageStore"
// import RoadReason from "../../reason/ui"
import { RoadDamageCreate, RoadDamageEdit } from "./index"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store/RoadTitleStore"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"
const roadTitleStore = useRoadTitleStore()

const initUserStore = useInitUserStore()

defineProps({
	date: {
		type: Array as PropType<IDatum[]>,
		default: null,
	},
})

const store = useRoadDamageStore()
const route = useRoute()
const roadId = Number(route.params.roadId)

const modalCreate = ref()
const modalEdit = ref()
// const modalReason = ref()
const canEdit = ref<boolean>()

const createItem = () => {
	modalCreate.value.showModal()
}

const editItem = () => {
	modalEdit.value.showModal(store.params.parentId, store.params.id)
}

// const showReason = () => {
// 	modalReason.value.showModal()
// }

const onDelete = () => {
	const year = store.damageList.find((item) => item.year === store.params.year)
	const dateSurvey = year?.items.find((item) => item.id === store.params.id)

	if (dateSurvey?.surveyed_date) {
		useDeleteItem({
			name: `${dateSurvey.lane_no} (สำรวจ: ${buddhistFormatDate(dateSurvey.surveyed_date, "dd mmm yy")})`,
			url: `roads/${store.params.id}/damage_import/${store.params.parentId}`,
			showAlert: true,
			callBack: async () => {
				store.map.Overlays.clear()
				await store.getDamageList(roadId).then(() => store.setDefaultparams())

				if (store.damageList.length > 0) {
					await store.getRoadDamageDetail(roadId)
				}
			},
		})
	}
}

onMounted(
	() =>
		(canEdit.value =
			initUserStore.accessPermissions[IUserRolesAccess.manage_road_damage] ||
			initUserStore.getIsOwnerManagePermission(
				initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_damage],
				roadTitleStore.data.ref_depot.id
			))
)

</script>

<template>
	<div v-show="!date" class="card card-rounded p-5 mt-5 pt-3">
		<div class="row">
			<div class="col-6 col-md-2">
				<VSelect :options="[]" label="ปี" name="year" placeholder="เลือก" />
			</div>
			<div class="col-6 col-md-4">
				<VSelect :options="[]" label="ช่องจราจร" name="lane" placeholder="เลือก" />
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
	</div>

	<div v-if="!date" class="mt-5">
		<VNotFound height="55dvh" />
	</div>

	<div v-else class="card card-rounded p-5 mt-5 pt-3">
		<div class="row">
			<div class="col-6 col-md-2">
				<VSelect
					v-model="store.params.year"
					:options="store.getYears ? store.getYears : []"
					label="ปี"
					name="year"
					placeholder="เลือก"
					:can-clear="false"
					:can-deselect="false"
					@update:model-value="() => store.onUpdateYear(roadId)"
				/>
			</div>
			<div class="col-6 col-md-4">
				<VSelect
					v-model="store.params.id"
					:options="store.createOptionsDate()"
					label="ช่องจราจร"
					name="lane"
					placeholder="เลือก"
					:can-clear="false"
					:can-deselect="false"
					@update:model-value="() => store.onUpdateIdParent(roadId)"
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
			<div v-show="store.data" class="col-md-8 col-12">
				<span
					class="badge rounded-1 fs-7 fw-semibold px-4 py-2"
					:class="`badge-status-${store.data?.road_damage.road_damage_status.status_code}`"
					>{{ store.data?.road_damage.road_damage_status.name }}</span
				>
				<label class="text-gray-900 mx-2 fs-7 mt-2">ปรับปรุงข้อมูลโดย</label>
				<VUser
					:label="`${store.data?.updated_by.firstname} ${store.data?.updated_by.lastname}`"
					:name="`${store.data?.updated_by.firstname} ${store.data?.updated_by.lastname}`"
					:role="`${store.data?.updated_by.ref_user_owner?.email}`"
					:image="`${store.data?.updated_by.profile_img_path}`"
				/>
				<label class="text-gray-900 ms-2 fs-7 mt-2">{{
					`เมื่อวันที่ ${buddhistFormatDate(store.data?.updated_date, "dd mmm yyyy เวลา HH:ii น.")}`
				}}</label>
			</div>
			<div v-show="store.data && canEdit" class="col-md-4 col-12 mt-md-0 mt-2 text-end">
				<a
					v-show="store.data?.road_damage.status !== 'W'"
					type="button"
					class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 mb-3 fw-semibold fs-6 lh-xxl"
					@click="editItem"
				>
					ปรับปรุงข้อมูล
				</a>
				<button
					v-show="store.data?.road_damage.status !== 'W' && store.data"
					type="button"
					class="btn btn-light rounded-4 px-7 py-2 ms-3 mb-3 fw-semibold text-black fs-6"
					@click="onDelete()"
				>
					ลบข้อมูล
				</button>
				<!-- <button
					v-show="store.data?.road_damage.status === 'R'"
					type="button"
					class="btn btn-outline btn-outline-primary rounded-4 px-4 py-3 ms-3 fw-semibold fs-6"
					@click="showReason"
				>
					เหตุผลการส่งกลับแก้ไข
				</button> -->
			</div>
		</div>
	</div>

	<!-- Modal -->
	<RoadDamageCreate ref="modalCreate" />
	<RoadDamageEdit ref="modalEdit" />
	<!-- <RoadReason ref="modalReason" :message="store.data?.road_damage.reject_reason" /> -->
</template>

<style scoped></style>
