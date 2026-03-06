<script setup lang="ts">
import { IRoadListRoads } from "../infrastructure/RoadListModel"
import { useRoadCreateStore, useRoadListStore } from "../store"
import SelectRampRoad from "./SelectRampRoad.vue"
import SelectModal from "./SelectModal.vue"

const modalCreate: Ref = ref()
const store = useRoadListStore()
const create = useRoadCreateStore()

const createItem = (data: any) => {
	create.roadInitParams.roadId = data.id.toString()
	create.roadInitParams.level = "2"
	create.getRoadInit(create.roadInitParams.roadId, create.roadInitParams.level, "")
	modalCreate.value.showModal()
}

defineProps({
	roads: {
		type: Object as PropType<IRoadListRoads[]>,
		required: true,
	},
	canEdit: {
		type: Boolean,
		required: true,
	},
})
</script>

<template>
	<template v-for="road in roads" :key="road.road_id">
		<div
			class="card overflow-hidden mb-4 border border-1 border-gray-300 shadow-none mx-4"
			:class="road.id === 2 ? 'bg-light-danger' : ''"
		>
			<div class="row align-items-center gx-0 pt-2 pb-4 bg-hover-gray-100 position-relative">
				<div class="col-lg-2 col-md-2 col-3 text-center cursor-pointer" @click="navigateTo(`/roads/${road.id}`)">
					<NuxtLink :to="`/roads/${road.id}`">
						<RoadIcon
							:id="road.id"
							:icon="road.road_info?.ref_road_type?.icon"
							:color="`${road.road_info?.road_color_code}`"
							width="30px"
							heigth="30px"
							custom-class="p-4"
						/>
					</NuxtLink>
				</div>
				<div class="col-lg-6 col-md-6 col-6 cursor-pointer" @click="navigateTo(`/roads/${road.id}`)">
					<h3 class="fw-normal align-items-start mt-2 mb-0" style="font-size: 16px">
						<NuxtLink :to="`/roads/${road.id}`" class="lh-lg me-2 text-gray-800"
							>ชื่อถนน: {{ road.road_info?.name }} ({{ road.road_info?.direction?.name }})</NuxtLink
						>
						<br />
					</h3>
					<p class="form-text fs-7 mt-2 mb-2">
						<NuxtLink :to="`/roads/${road.id}`">
							<span class="text-gray-800 fs-6">จาก - ถึง: {{ road.road_info.origin_to_destination }} </span>
						</NuxtLink>
						<br />
						<NuxtLink :to="`/roads/${road.id}`" class="text-gray-600">
							<span
								>{{ road.road_info?.road_code }} | {{ road.road_info?.responsible_code }} | กม. ที่
								{{ convertMeterToKm(road.road_info?.km_start) }} - {{ convertMeterToKm(road.road_info?.km_end) }}</span
							>
						</NuxtLink>
						<NuxtLink class="cursor-pointer" @click.stop="store.setLocation(road)">
							<i class="fi fi-sr-marker align-middle ms-4 fs-2" :style="{ color: road.road_info?.road_color_code }"></i>
						</NuxtLink>
						<br />
					</p>
					<div v-if="road.survey_status" class="mb-4 d-flex">
						<img class="me-2" src="/images/icons/svg/survey-alert.svg" alt="survey-alert" />
						<span class="text-danger">ต้องดำเนินการสำรวจอีกครั้งเนื่องจาก มีการซ่อมบำรุงผิวทาง</span>
					</div>
					<div class="d-flex gap-2">
						<div v-if="road.road_surface_icon.length > 0" class="d-flex gap-2">
							<template v-for="(surface, index) in road.road_surface_icon" :key="index">
								<span
									:style="`background-color: ${surface.color_code}`"
									class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
								>
									{{ surface.name }}
								</span>
							</template>
							<span
								v-show="road.condition_status"
								:style="`background-color: ${road.condition_status_color}`"
								class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
							>
								R. Condition
							</span>

							<span
								v-show="road.retro_status"
								:style="`background-color: ${road.retro_status_color}`"
								class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
							>
								R. Pavement Marking
							</span>
							<span
								v-show="road.damage_status"
								:style="`background-color: ${road.damage_status_color}`"
								class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
							>
								R. Damage
							</span>
						</div>

						<div v-else class="d-flex gap-2" style="height: 2em"></div>
						<div class="position-absolute" style="bottom: 12px; right: 12px">
							<a v-show="canEdit" class="btn btn-primary btn-add lh-xxl" @click.stop="createItem(road)">เพิ่มข้อมูล</a>
						</div>
					</div>
				</div>
				<div class="col-lg-4 col-md-4 col-3 align-self-start text-md-end text-start">
					<div class="d-flex justify-content-end pe-3">
						<h3 class="fw-normal fs-5 mt-3">
							{{ calculateDistance(road.road_info?.km_start, road.road_info?.km_end) }} กม.
						</h3>
						<div :style="road.child_roads.length === 0 ? 'visibility: hidden;' : ''">
							<button
								class="btn btn-collapsed rounded-3 collapsed"
								data-bs-toggle="collapse"
								:data-bs-target="`#collapseRamps${road.id}`"
								:aria-controls="`collapseRamps${road.id}`"
								type="button"
							>
								<i class="fi fi-rr-caret-up fs-1 text-gray-600 pe-0"></i>
							</button>
						</div>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-12">
					<SelectRampRoad :road-id="road.id" :ramps="road.child_roads" />
				</div>
			</div>
		</div>
	</template>
	<SelectModal ref="modalCreate" title="modal-main-road" />
</template>

<style lang="scss" scoped>
.collapsed i {
	transition: transform 0.3s ease-in-out, padding-right 0.3s ease-in-out;
	padding-right: 0rem !important;
}

i {
	transition: transform 0.3s ease-in-out;
}

/* Define the rotated state */
.collapsed i {
	transform: rotate(180deg) !important;
}

.btn-collapsed {
	height: auto !important;
	padding: 7px !important;
	line-height: 1.25;
}

.btn-add {
	height: 28px;
	padding: 0 1.5rem !important;
	font-size: 12px;
}

.card {
	border-radius: 16px !important;
}

p.form-text {
	width: 162%;
	@media screen and (max-width: 768px) {
		width: 145%;
	}
}
</style>
