<script setup lang="ts">
import { useRoadTitleStore } from "../store"

const props = defineProps({
	roadId: {
		type: Number,
		default: 0,
	},
	isLeftSumPosition: {
		type: Boolean,
		default: false,
	},
})

const store = useRoadTitleStore()

useStoreLifecycle(store, { resetOnEnter: false })

onBeforeMount(async () => {
	if (props.roadId > 0) {
		await store.getData(Number(props.roadId))
	}
})

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div
			v-show="store.data"
			class="row align-items-center gx-0 py-2"
			:class="isLeftSumPosition ? '' : 'justify-content-between'"
		>
			<div class="col-lg-2 col-md-2 col-3 text-center" style="max-width: 135px">
				<RoadIcon
					:id="store.data.id"
					:icon="store.data.road_info?.ref_road_type?.icon"
					:color="`${store.data.road_info?.road_color_code}`"
					width="30px"
					heigth="30px"
					custom-class="p-4"
				/>
			</div>
			<div class="col-lg-6 col-md-6 col-6">
				<h3 class="fw-normal mt-2 mb-0">
					<span v-if="store.data.road_level === 1" class="lh-lg me-2 text-gray-800" style="font-size: 16px"
						>ชื่อถนน: {{ store.data.road_section_name_th }}</span
					>
					<span v-else class="lh-lg me-2 text-gray-800" style="font-size: 16px">{{ store.data.road_info?.name }}</span>
				</h3>
				<p v-if="store.data.road_level === 1" class="form-text fs-6 mt-1 mb-2">
					<span>จาก - ถึง: {{ store.data.origin_to_destination }} </span>
				</p>
				<p class="form-text fs-7 mt-1 mb-2">
					<span class="text-gray-600"
						>{{
							`${store.data.road_code} | ${store.data.responsible_code} | กม.ที่ ${convertMeterToKm(
								store.data.road_info?.km_start
							)} - ${convertMeterToKm(store.data.road_info?.km_end)}`
						}}
					</span>
				</p>
			</div>
			<div v-show="!isLeftSumPosition" class="col-lg-4 col-md-4 col-3 text-md-end text-end">
				<div class="row align-items-center">
					<div class="col-md-11 col-12 text-center text-md-end">
						<span class="fw-normal fs-5 mt-3"
							>{{ calculateDistance(store.data.road_info?.km_start, store.data.road_info?.km_end) }} กม.</span
						>
					</div>
					<div class="col-md-11 col-12 text-center text-md-end pt-3">
						<template v-for="(surface, index) in store.data.road_surface_icon" :key="index">
							<span
								:style="`background-color: ${surface.color_code}`"
								class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal"
								:class="store.data.road_surface_icon.length - 1 === index ? '' : 'me-2'"
							>
								{{ surface.name }}
							</span>
						</template>
					</div>
				</div>
			</div>
		</div>
	</VSkeletonLoader>
</template>

<style scoped></style>
