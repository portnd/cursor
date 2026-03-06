<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import { RoadDetailModal } from "./index"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const store = useRoadSummaryStore()
const route = useRoute()
const modalEdit: Ref = ref()
const roadLevel = ref(1)

onMounted(() => {
	store.getRoadDetail(Number(route.params.roadId))
})

const editItem = (roadId: number) => {
	modalEdit.value?.showModal(store.road)
	roadLevel.value = roadId
}
</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div id="row" class="row mb-5">
			<div class="col-12 col-md-9 text-start align-self-center ps-4 order-last order-md-first">
				<label class="text-gray-900 fs-6 me-1">ปรับปรุงข้อมูลโดย</label>
				<VUser
					:label="store.road.road_info?.user.firstname + ' ' + store.road.road_info?.user.lastname"
					:name="store.road.road_info?.user.firstname + ' ' + store.road.road_info?.user.lastname"
					:role="'-'"
				/>
				<label class="text-gray-900 fs-6 ms-1">
					เมื่อวันที่
					{{ buddhistFormatDate(new Date(store.road.road_info?.updated_at), "dd mmm yyyy เวลา HH:ii น.") }}
				</label>
			</div>
			<div class="col-12 col-md-3 text-end">
				<NuxtLink
					v-if="
						initUserStore.accessPermissions[IUserRolesAccess.manage_road_summary] ||
						initUserStore.getIsOwnerManagePermission(
							initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_summary],
							store.road.ref_depot.id
						)
					"
					class="btn btn-outline btn-outline-primary rounded-2 px-5 fw-semibold fs-6"
					@click="editItem(store.road.road_level)"
				>
					ปรับปรุงข้อมูล
				</NuxtLink>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">หมายเลขตอนควบคุม 8 หลัก</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_code }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ชื่อตอนควบคุม (ไทย)</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_section_name_th }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ชื่อตอนควบคุม (อังกฤษ)</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_section_name_en }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">จังหวัด</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.province }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">หน่วยงานที่รับผิดชอบ</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.responsible_code }}</span>
			</div>
		</div>
		<div v-if="store.road.road_level === 1" class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">จาก - ถึง</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.origin_to_destination }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ประเภทของถนน</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_info?.ref_road_type.name }}</span>
			</div>
		</div>
		<div v-if="store.road.road_level === 2" class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">รหัส Ramp</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_info?.ramp_id }}</span>
			</div>
		</div>
		<div v-if="store.road.road_level === 2" class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ชื่อ Ramp</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700">{{ store.road.road_info?.name }}</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ช่วง กม.</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700"
					>{{ convertMeterToKm(store.road.road_info?.km_start) }} -
					{{ convertMeterToKm(store.road.road_info?.km_end) }}</span
				>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">ระยะทาง</span>
			</div>
			<div class="col-md-7 col-12">
				<span class="fs-6 text-gray-700"
					>{{ calculateDistance(store.road.road_info?.km_start, store.road.road_info?.km_end) }} กม.</span
				>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">Shape File</span>
			</div>
			<div class="col-md-7 col-12">
				<ul
					v-if="store.road.road_info?.center_line_shape_file_path && store.road.road_info?.center_lane_shape_file_path"
					class="row ps-0 mt-2 mb-0"
				>
					<li class="col col-file text-primary d-flex flex-column align-items-center text-center p-0">
						<NuxtLink
							class="align-self-center"
							style="width: fit-content"
							:href="store.road.road_info?.center_line_shape_file_path"
							target="_blank"
							download
						>
							<img src="/images/files/zip.png" width="50" />
						</NuxtLink>
						<p class="filename">center-line</p>
					</li>
					<li class="col col-file text-primary d-flex flex-column align-items-center text-center p-0">
						<NuxtLink
							class="align-self-center"
							style="width: fit-content"
							:href="store.road.road_info?.center_lane_shape_file_path"
							target="_blank"
							download
						>
							<img src="/images/files/zip.png" width="50" />
						</NuxtLink>
						<p class="filename">center-lane</p>
					</li>
				</ul>
				<span v-else class="text-gray-600">ไม่พบข้อมูล</span>
			</div>
		</div>
		<div class="row mb-5">
			<div class="col-md-5 col-12">
				<span class="fs-6 fw-semibold">หมายเหตุ</span>
			</div>
			<div class="col-md-7 col-12">
				<span
					class="fs-6 text-gray-700 text-preline"
					v-html="store.road.road_info?.remark === '' ? '-' : store.road.road_info?.remark"
				></span>
			</div>
		</div>
	</VSkeletonLoader>
	<RoadDetailModal ref="modalEdit" :road-id="roadLevel" />
</template>

<style scoped lang="scss">
ul {
	list-style-type: none;
}

.filename {
	margin-top: 1px;
	font-size: 0.925rem;
	font-weight: 400;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 66px;
	display: inline-block;
}
.col-file {
	max-width: 6em;
}
</style>
