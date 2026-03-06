<script setup lang="ts">
import { useRoadSummaryStore } from "../store"
import { RoadSummaryTab } from "./index"
import RoadMenu from "~~/core/modules/road/info/menu/ui"

const store = useRoadSummaryStore()

useStoreLifecycle(store, { resetOnEnter: false })

const setMap = store.setMap

onMounted(async () => {
	const route = useRoute()
	await store.getRoadDetail(Number(route.params.roadId))
	// store.getLaneList()
	// await store.getConditionCompareAverge()
})

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<RoadMenu />
			<RoadSummaryTab />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" height="calc(97vh)" :is-sticky="true" @map="setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
