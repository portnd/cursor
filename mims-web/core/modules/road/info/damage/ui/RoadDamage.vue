<script setup lang="ts">
// import { useRoadListStore } from "../../../roadList/store"
import { useRoadDamageStore } from "../store/RoadDamageStore"
import { RoadDamageSearch, RoadDamageTable, RoadDamagePhotoViewer, RoadDamageSummary } from "./index"
import RoadMenu from "~~/core/modules/road/info/menu/ui"

const store = useRoadDamageStore()

const route = useRoute()
const roadId = Number(route.params.roadId)

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

onBeforeMount(async () => {
	store.$reset()
	await store.getDamageList(roadId).then(() => store.setDefaultparams())
	await store.getLaneList(roadId)

	if (store.damageList?.length > 0) {
		await store.getRoadDamageDetail(roadId)
	}
})

onUnmounted(() => {
	store.$reset()
})

onBeforeMount(() => {
	window.addEventListener("beforeunload", () => {
		localStorage.removeItem("search")
	})
})
</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<RoadMenu @color="store.setColor" @geom="store.setMainGeom" />
			<VSkeletonLoader :loading="store.loading">
				<RoadDamageSearch :date="store.damageList" />
				<RoadDamageSummary v-show="store.damageList" />

				<RoadDamageTable v-show="store.damageList" :data="store.getRoadDamageRange" />
			</VSkeletonLoader>
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="false" height="calc(60vh)" :is-sticky="true" @map="store.setMap" />
				</KeepAlive>

				<RoadDamagePhotoViewer :map-show="mapShow.collapsed" :image-data="store.image.path" />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
