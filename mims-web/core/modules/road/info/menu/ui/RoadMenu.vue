<script setup lang="ts">
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"
import RoadTitle from "~~/core/modules/common/roadTitle/ui"

import { useInitUserStore } from "~/core/modules/initUser/store"
const initUserStore = useInitUserStore()

const route = useRoute()
const roadId = Number(route.params.roadId)

const tab = route.path.split("/")[3]

const roadTitleStore = useRoadTitleStore()

useStoreLifecycle(roadTitleStore, { resetOnEnter: false })

const emit = defineEmits(["color", "geom"])

watch(
	() => roadTitleStore.data,
	(newData, oldData) => {
		if (newData !== oldData) {
			const roadColor = roadTitleStore.data.road_info?.road_color_code
			emit("color", roadColor)
			emit("geom", roadTitleStore.data.road_info?.the_geom)
		}
	}
)

</script>

<template>
	<div class="card card-rounded">
		<RoadTitle :road-id="roadId" />

		<div class="row text-center mb-2 gx-0">
			<!-- begin::Tab -->
			<ul class="nav nav-tabs nav-pills border-0 d-flex px-3 pb-3">
				<li v-if="initUserStore.accessPermissions[IUserRolesAccess.view_road_summary]" class="nav-item">
					<NuxtLink :to="`/roads/${roadId}/summary`" class="nav-link" :class="{ active: tab === 'summary' }"
						>ข้อมูลทั่วไป</NuxtLink
					>
				</li>
				<li v-if="initUserStore.accessPermissions[IUserRolesAccess.view_road_in_assets]" class="nav-item">
					<NuxtLink :to="`/roads/${roadId}/in-asset`" class="nav-link" :class="{ active: tab === 'in-asset' }"
						>สินทรัพย์</NuxtLink
					>
				</li>
				<li v-if="initUserStore.accessPermissions[IUserRolesAccess.view_road_condition]" class="nav-item">
					<NuxtLink :to="`/roads/${roadId}/condition`" class="nav-link" :class="{ active: tab === 'condition' }"
						>สภาพทาง</NuxtLink
					>
				</li>
				<li v-if="initUserStore.accessPermissions[IUserRolesAccess.view_road_retro]" class="nav-item">
					<NuxtLink
						:to="`/roads/${roadId}/reflective-strip`"
						class="nav-link"
						:class="{ active: tab === 'reflective-strip' }"
						>แถบสะท้อนแสง</NuxtLink
					>
				</li>
				<li v-if="initUserStore.accessPermissions[IUserRolesAccess.view_road_damage]" class="nav-item">
					<NuxtLink :to="`/roads/${roadId}/damage`" class="nav-link" :class="{ active: tab === 'damage' }"
						>ความเสียหาย</NuxtLink
					>
				</li>
			</ul>
			<!-- end::Tab -->
		</div>
	</div>
</template>

<style scoped lang="scss">
.nav-pills .nav-item {
	width: calc(100% / 5);
	margin-right: 0px;
	text-align: -webkit-center;
	margin-top: 5px;

	@media (max-width: 1360px) {
		width: calc(100% / 4);
	}
	@media (max-width: 1050px) {
		width: calc(100% / 3);
	}
	@media (max-width: 991px) {
		width: calc(100% / 5);
	}
	@media (max-width: 885px) {
		width: calc(100% / 4);
	}
	@media (max-width: 640px) {
		width: calc(100% / 3);
	}
	@media (max-width: 460px) {
		width: calc(100% / 2);
	}
}
.nav-pills .nav-link {
	width: 95%;
	height: 40px;
	color: var(--kt-gray-800);
	line-height: 1.5;
	border: 0;
}
</style>
