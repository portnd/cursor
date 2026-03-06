<script setup lang="ts">
import { useRoadListStore, useRoadCreateStore } from "../store"
import SelectMainRoad from "./SelectMainRoad.vue"
import SelectModal from "./SelectModal.vue"
import { IRoadList, IRoadListRoads } from "~~/core/modules/road/roadList/infrastructure"

import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const modalCreate: Ref = ref()
const store = useRoadListStore()
const create = useRoadCreateStore()
const createItem = (data: any, id: number) => {
	create.roadInitParams.roadId = data?.road_section_id.toString()
	create.roadInitParams.level = data?.road_level.toString()
	create.getRoadInit(id.toString(), data?.road_level.toString(), "")
	modalCreate.value.showModal()
}

defineProps({
	roads: {
		type: Object as PropType<IRoadList[]>,
		required: true,
	},
	loading: {
		type: Boolean,
		required: true,
	},
})

const getCanEdit = (refDepotCode: number) => {
	return (
		initUserStore.accessPermissions[IUserRolesAccess.manage_all_road] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road],
			refDepotCode
		)
	)
}
</script>

<template>
	<div id="road_accordion" class="accordion border border-0">
		<template v-if="loading">
			<div v-for="n in 6" :key="n" class="mb-5">
				<VSkeletonLoader>
					<div class="card d-flex flex-row py-1 align-items-center justify-content-between w-100 px-4">
						<div class="w-50 h-50">
							<p class="mb-0">ชื่อสายทาง</p>
						</div>
						<div class="accordion-button collapsed fs-4 py-4 w-25 h-50"></div>
					</div>
				</VSkeletonLoader>
			</div>
		</template>
		<template v-else-if="roads.length > 0">
			<template v-for="(parent, index) of roads" :key="parent.id">
				<div class="accordion-item card py-1 mb-5">
					<h2 :id="`#road_accordion_main_header${parent.id}`" class="accordion-header">
						<button
							class="accordion-button fs-5 fw-semibold pt-4 py-0 text-gray-800 bg-white"
							:class="{ collapsed: index !== 0 }"
							type="button"
							data-bs-toggle="collapse"
							:data-bs-target="`#road_accordion_main_body${parent.id}`"
							:aria-controls="`road_accordion_main_body${parent.id}`"
						>
							<VPopover
								:is-hover="true"
								:hover="true"
								placement="bottom"
								class="p-0"
								style="height: fit-content"
								type="button"
							>
								<div class="fw-normal text-start" style="font-size: 16px">
									ทางหลวงพิเศษหมายเลข {{ parent.number }} {{ store.extractText(parent.short_name, true) }}
									<span v-show="store.extractText(parent.short_name) !== null" class="text-gray-600"
										>{{ " " + store.extractText(parent.short_name) }}
									</span>
								</div>
								<template #content>
									<span>ทางหลวงพิเศษหมายเลข {{ parent.number }} {{ parent.name }}</span>
								</template>
							</VPopover>
						</button>
						<div class="fw-normal sub-road text-gray-600 px-6 mt-2 mb-5">
							(กม. ที่ {{ convertMeterToKm(parent.km_start) }} - {{ convertMeterToKm(parent.km_end) }}) ระยะทาง
							{{ toNumber(parent.distance) }} กม.
						</div>
					</h2>
					<div
						:id="`road_accordion_main_body${parent.id}`"
						data-bs-parent="#road_accordion_main"
						:aria-labelledby="`#road_accordion_main_header${parent.id}`"
						class="accordion-collapse collapse"
						:class="{ show: index <= 0 }"
					>
						<div class="px-5 pt-1">
							<template v-for="(road, key) in parent.sections" :key="key">
								<div class="accordion-item card py-1 mb-5">
									<h2 :id="`#road_accordion_sub_header${road.id}`" class="accordion-header">
										<div
											class="d-flex justify-content-between align-items-center"
											:style="road.roads.length === 0 ? `visibility: hidden;` : ``"
										>
											<button
												class="accordion-button fs-5 fw-semibold py-4 text-gray-800 bg-white"
												:class="{ collapsed: key > 0 }"
												type="button"
												data-bs-toggle="collapse"
												:data-bs-target="`#road_accordion_sub_body${road.id}`"
												:aria-controls="`road_accordion_sub_body${road.id}`"
											>
												<span class="road-name">
													ตอนควบคุม {{ road.number }}: {{ road.name_origin_th }}
													{{ road.name_destination_th === "" ? "" : "- " + road.name_destination_th }}
												</span>
											</button>

											<div class="position-absolute end-0 me-16">
												<a
													v-show="getCanEdit(road.ref_depot.id)"
													class="btn btn-primary btn-add lh-xxl"
													@click="
														createItem(
															road.roads.find((item) => item.road_section_id === road.id) as IRoadListRoads,
															road.id
														)
													"
												>
													เพิ่มข้อมูล
												</a>
												<!--  -->
											</div>
										</div>
									</h2>
									<div
										:id="`road_accordion_sub_body${road.id}`"
										data-bs-parent="#road_accordion_sub"
										:aria-labelledby="`#road_accordion_sub_header${road.id}`"
										class="accordion-collapse collapse"
										:class="{ show: key <= 0 }"
									>
										<!-- <div class="m-4"> -->
										<SelectMainRoad :roads="road.roads" :can-edit="getCanEdit(road.ref_depot.id)" />
										<!-- </div> -->
									</div>
								</div>
							</template>
						</div>
					</div>
				</div>
			</template>
		</template>

		<template v-else>
			<VNotFound height="80vh" />
		</template>
	</div>
	<SelectModal ref="modalCreate" title="modal-road" />
</template>

<style scoped lang="scss">
.road-name {
	width: 85%;
	font-size: 16px;
	font-weight: 400;
	@media only screen and (max-width: 576px) {
		width: 80%;
	}
}
.sub-road {
	font-size: 12px;
}
.accordion,
.accordion-button {
	border-radius: 0.75rem !important;
	z-index: 0 !important;
	box-shadow: none !important;
}
.accordion-item,
.accordion-item .accordion-button.collapsed {
	border-top-right-radius: 16px !important;
	border-top-left-radius: 16px !important;
	border-bottom-left-radius: 16px !important;
	border-bottom-right-radius: 16px !important;
}
.accordion-item .accordion-item {
	border: 1px solid var(--kt-gray-300);
	box-shadow: none;
	.road-name {
		width: 75%;
		@media only screen and (max-width: 576px) {
			width: 60%;
		}
	}
}

.btn-add {
	height: 28px;
	padding: 0 1.5rem !important;
	font-size: 12px;
}
</style>
