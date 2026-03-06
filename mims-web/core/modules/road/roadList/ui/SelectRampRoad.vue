<script setup lang="ts">
import { IRoadListRoads } from "../infrastructure"
import { useRoadListStore } from "../store"

const store = useRoadListStore()

defineProps({
	ramps: {
		type: Object as PropType<IRoadListRoads[]>,
		default: null,
	},
	roadId: {
		type: Number,
		required: true,
	},
})
</script>

<template>
	<div :id="`collapseRamps${roadId}`" class="collapse ramps">
		<template v-for="ramp in ramps" :key="ramp.road_id">
			<div class="ramp-item">
				<div
					class="px-3 pt-1 pb-4 border border-start-0 border-top-1 border-end-0 border-bottom-0 border-gray-300 bg-hover-gray-100"
				>
					<div class="row align-items-center gx-0 ps-0 ps-md-5">
						<div
							class="col-lg-2 col-md-2 col-3 text-center cursor-pointer"
							@click="navigateTo(`/roads/${ramp.id}/summary`)"
						>
							<NuxtLink :to="`/roads/${ramp.id}/summary`">
								<RoadIcon
									:id="ramp.id"
									:icon="ramp.road_info?.ref_road_type?.icon"
									:color="`${ramp.road_info?.road_color_code}`"
									width="30px"
									heigth="30px"
									custom-class="p-4"
								/>
							</NuxtLink>
						</div>
						<div class="col-lg-7 col-md-7 col-7 cursor-pointer">
							<h3 class="fw-normal mt-2 mb-0" style="font-size: 16px">
								<NuxtLink :to="`/roads/${ramp.id}/summary`" class="lh-lg me-2 text-gray-800">
									{{ ramp.road_info?.name }}
								</NuxtLink>
								<br />
							</h3>
							<p class="form-text fs-7 mt-1 mb-2">
								<NuxtLink :to="`/roads/${ramp.id}/summary`" class="text-gray-600">
									{{ ramp.road_info?.road_code }} | {{ ramp.road_info?.responsible_code }} | กม. ที่
									{{ convertMeterToKm(ramp.road_info?.km_start) }} -
									{{ convertMeterToKm(ramp.road_info?.km_end) }}
								</NuxtLink>
								<NuxtLink class="cursor-pointer" @click.stop="store.setLocation(ramp)">
									<i
										class="fi fi-sr-marker align-middle ms-4 fs-2"
										:style="{ color: ramp.road_info?.road_color_code }"
									></i>
								</NuxtLink>
								<br />
								<!-- <span v-show="ramp.id === 10" class="text-danger">
									<div class="d-flex gap-1">
										<i class="fi fi-rr-exclamation mt-1 lh-0"></i>
										<span>ต้องดำเนินการสำรวจอีกครั้งเนื่องจาก มีการซ่อมบำรุงผิวทาง</span>
									</div>
								</span> -->
							</p>
							<div v-if="ramp.survey_status" class="mb-4 d-flex">
								<img class="me-2" src="/images/icons/svg/survey-alert.svg" alt="survey-alert" />
								<span class="text-danger">ต้องดำเนินการสำรวจอีกครั้งเนื่องจาก มีการซ่อมบำรุงผิวทาง</span>
							</div>
							<div v-if="ramp.road_surface_icon.length > 0" class="d-flex gap-2">
								<div class="d-flex gap-2">
									<template v-for="(surface, index) in ramp.road_surface_icon" :key="index">
										<div>
											<span
												:style="`background-color: ${surface.color_code}`"
												class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
											>
												{{ surface.name }}
											</span>
										</div>
									</template>
									<span
										v-if="ramp.condition_status"
										:style="`background-color: ${ramp.condition_status_color}`"
										class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
									>
										R. Condition
									</span>
									<span
										v-show="ramp.retro_status"
										:style="`background-color: ${ramp.retro_status_color}`"
										class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
									>
										R. Pavement Marking
									</span>
									<span
										v-show="ramp.damage_status"
										:style="`background-color: ${ramp.damage_status_color}`"
										class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
									>
										R. Damage
									</span>
								</div>
							</div>
						</div>
						<div class="col-lg-3 col-sm-3 col-2 text-md-end text-start align-self-start">
							<div class="row align-items-center">
								<div class="col-12 text-end">
									<h3 class="fw-normal fs-5 mt-3 ms-3">
										{{ calculateDistance(ramp.road_info?.km_start, ramp.road_info?.km_end) }} กม.
									</h3>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</template>
	</div>
</template>

<style lang="scss" scoped>
.ramps .ramp-item:last-child > div {
	border-radius: 0px 0px 8px 8px !important;
}

p.form-text {
	width: 145%;
}
</style>
