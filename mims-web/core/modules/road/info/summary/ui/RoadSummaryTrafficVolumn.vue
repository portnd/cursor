<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import RoadSummaryTrafficVolumnCreateModal from "./RoadSummaryTrafficVolumnCreateModal.vue"
import RoadSummaryTrafficVolumnEditModal from "./RoadSummaryTrafficVolumnEditModal.vue"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const store = useRoadSummaryStore()
const modalTrafficEdit: Ref = ref()
const modalTrafficCreate: Ref = ref()

const route = useRoute()
const roadId = Number(route.params.roadId)
const canEdit = ref<boolean>()

const editItem = () => {
	modalTrafficEdit.value.showModal()
}

const createItem = () => {
	modalTrafficCreate.value.showModal()
}

onBeforeMount(async () => {
	await store.getTrafficRevision(roadId).then(() => store.setDefaultOptions())
})

const onDelete = () => {
	const date = store.trafficDetail?.surveyed_date
		? buddhistFormatDate(store.trafficDetail.surveyed_date, "dd mmm yyyy")
		: ""
	useDeleteItem({
		name: "สำรวจ: " + date,
		url: `roads/${roadId}/volume_aadt/${store.toggle.aadtId}`,
		callBack: async () => {
			await store.getTrafficRevision(roadId).then(() => store.setDefaultOptions())
			await store.getTrafficDetail(store.toggle.aadtId)
		},
	})
}

watch(
	() => store.toggle.year,
	(newYear) => {
		if (newYear) {
			store.toggle.aadtId = store.getYearItemsOptions?.[0]?.value || 0
		}
	}
)

watch(
	() => store.toggle.aadtId,
	async (newAadtId) => {
		if (!newAadtId || !store.trafficRevision?.length) {
			return
		}

		store.trafficRevision.forEach((parent) => {
			if (parent.year === store.toggle.year && parent.items?.length) {
				parent.items.forEach(async (child) => {
					if (child.id === newAadtId) {
						store.toggle.parentID = child.id_parent
						store.toggle.aadtId = child.id
						await store.getTrafficDetail(store.toggle.aadtId)
					}
				})
			}
		})
	}
)

const checkWidth = () => {
	const section = document.querySelector(".editor") as HTMLElement
	const title = document.querySelector(".col-title") as HTMLElement
	const button = document.querySelector(".col-button") as HTMLElement
	if (section) {
		const offsetWidth = section.offsetWidth
		const innerWidth = window.innerWidth
		let titleWidth = "100%"
		let buttonWidth = "100%"
		if (offsetWidth > 550 && innerWidth > 1037) {
			titleWidth = "62%"
			buttonWidth = "38%"
		} else if (innerWidth > 767 && innerWidth < 992) {
			titleWidth = "65%"
			buttonWidth = "35%"
		} else if (offsetWidth < 712 && innerWidth < 768) {
			titleWidth = "100%"
			buttonWidth = "100%"
		}
		title.style.width = titleWidth

		if (button) {
			button.style.width = buttonWidth
		}
	}
}

const updateWidth = () => {
	setTimeout(() => {
		checkWidth()
	}, 300)
}

onMounted(() => {
	canEdit.value =
		initUserStore.accessPermissions[IUserRolesAccess.manage_road_traffic] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_traffic],
			store.road.ref_depot.id
		)

	checkWidth()
	window.addEventListener("resize", updateWidth)
})

onUnmounted(() => {
	window.removeEventListener("resize", updateWidth)
})

defineExpose({
	updateWidth,
})
</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="row mb-3">
			<template v-if="store.trafficRevision.length === 0">
				<div class="col-6 col-md-3">
					<VSelect :options="[]" label="ปี" name="year" placeholder="เลือก" />
				</div>
				<div class="col-6 col-md-5">
					<VSelect :options="[]" label="วันที่สำรวจ" name="lane" placeholder="เลือก" />
				</div>
				<div class="col-md-4 col-12 mt-md-0 mt-2 align-self-end text-end mb-1">
					<a
						v-if="canEdit"
						type="button"
						class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6 lh-xl"
						@click="createItem()"
					>
						เพิ่มข้อมูล
					</a>
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
				<div class="col-6 col-md-4">
					<VSelect
						v-model="store.toggle.aadtId"
						:options="store.getYearItemsOptions"
						:can-clear="false"
						:can-deselect="false"
						label="วันที่สำรวจ"
						name="lane"
						placeholder="เลือก"
					/>
				</div>
				<div class="col-md-5 col-12 mt-md-0 mt-2 align-self-end text-end mb-1 d-md-block d-none">
					<button
						v-if="canEdit"
						type="button"
						class="btn btn-primary rounded-2 px-6 py-3 fw-semibold fs-6"
						@click="createItem()"
					>
						เพิ่มข้อมูล
					</button>
				</div>
			</template>
		</div>
		<div v-if="store.trafficRevision.length > 0" class="row editor mt-3">
			<div class="col-title">
				<template v-if="store.getUpdateBy?.date !== ''">
					<label class="text-gray-900 fs-7 me-2">ปรับปรุงข้อมูลโดย</label>
					<VUser :label="`${store.getUpdateBy?.name ?? ''}`" :name="`${store.getUpdateBy?.name ?? ''}`" :role="'-'" />
					<span class="text-gray-900 ms-2 fs-7 mt-2" style="word-wrap: break-word"
						>เมื่อวันที่
						{{ store.trafficDetail?.updated_date ? buddhistFormatDate(new Date(store.trafficDetail.updated_date), "dd mmm yyyy เวลา HH:ii น.") : '' }}</span
					>
				</template>
			</div>
			<div v-if="canEdit" class="col-button mt-md-0 mt-2 text-end">
				<button
					type="button"
					class="btn btn-primary rounded-2 px-6 py-3 mb-3 fw-semibold me-3 fs-6 d-md-none"
					@click="createItem()"
				>
					เพิ่มข้อมูล
				</button>
				<button
					class="btn btn-outline btn-outline-primary rounded-2 px-3 mb-3 py-2 fw-semibold fs-6"
					@click="editItem()"
				>
					ปรับปรุงข้อมูล
				</button>
				<button
					type="button"
					class="btn btn-outline btn-outline-danger rounded-2 px-3 ms-3 mb-3 py-2 fw-semibold fs-6"
					@click="onDelete"
				>
					ลบข้อมูล
				</button>
			</div>
		</div>
		<VNotFound v-if="store.trafficRevision.length === 0" style="box-shadow: none" />
		<div v-else class="row">
			<div class="col-12">
				<div class="table-responsive">
					<table class="table customize-basic-table mb-0 text-truncate table-hover">
						<thead>
							<tr>
								<th class="text-center">ประเภท</th>
								<th class="text-center">ปริมาณ (คัน/วัน)</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>ปริมาณจราจร รถ 4 ล้อ</td>
								<td class="text-end pe-10">{{ toNumber(store.trafficDetail?.veh1) }}</td>
							</tr>
							<tr>
								<td>ปริมาณจราจร รถ 6 ล้อ</td>
								<td class="text-end pe-10">{{ toNumber(store.trafficDetail?.veh2) }}</td>
							</tr>
							<tr>
								<td>ปริมาณจราจร รถมากกว่า 6 ล้อ</td>
								<td class="text-end pe-10">{{ toNumber(store.trafficDetail?.veh3) }}</td>
							</tr>
							<tr class="border border-top">
								<td style="border-left: none">รวม</td>
								<td class="text-end pe-10" style="border-left: none">{{ toNumber(store.trafficDetail?.total) }}</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</VSkeletonLoader>
	<RoadSummaryTrafficVolumnCreateModal ref="modalTrafficCreate" />
	<RoadSummaryTrafficVolumnEditModal ref="modalTrafficEdit" />
</template>

<style scoped>
.editor {
	justify-content: space-between;
}
.col-title {
	width: 60%;
}
.col-button {
	width: 40%;
}
</style>
