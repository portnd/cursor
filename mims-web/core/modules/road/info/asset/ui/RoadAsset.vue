<script setup lang="ts">
import { useRoadListStore } from "../../../roadList/store"
import { useRoadAssetItemStore } from "../store"
import { RoadAssetItem } from "./index"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import RoadMenu from "~~/core/modules/road/info/menu/ui"

const props = defineProps({
	assetType: {
		type: String,
	},
})

interface IStateParams {
	keyword: string
	road_group_id: string[]
	road_section_id: number[]
	km_start: string
	km_end: string
	depot_code: number[]
	ref_surface_id: number[]
}

const store = useRoadAssetItemStore()
const storeList = useRoadListStore()
const roadTitleStore = useRoadTitleStore()
const route = useRoute()

useStoreLifecycle(store)

onMounted(async () => {
	store.assetType = props.assetType ?? ""
	store.getAsset()
	// โหลดข้อมูลสายทางก่อน (หรือรอให้ RoadTitle โหลดเสร็จ) แล้วค่อยเซ็ต geom/color และวาดแผนที่ (แก้ refresh แล้วแผนที่เป็นสีฟ้า)
	const roadId = Number(route.params.roadId)
	if (!roadTitleStore.data?.road_info?.the_geom) {
		await roadTitleStore.getData(roadId)
	}
	store.color = roadTitleStore.data?.road_info?.road_color_code ?? ""
	store.geom = roadTitleStore.data?.road_info?.the_geom ?? ""
	if (store.map) {
		store.defaultLocation()
		store.createLine()
	}
})

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

onBeforeMount(() => {
	window.addEventListener("beforeunload", () => {
		localStorage.removeItem("search")
		storeList.params = {} as IStateParams
	})
})

</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<RoadMenu />
			<RoadAssetItem :data-list="store.menu" :asset-type="assetType" :loading="store.loading" />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="store.loading" height="calc(97vh)" :is-sticky="true" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
