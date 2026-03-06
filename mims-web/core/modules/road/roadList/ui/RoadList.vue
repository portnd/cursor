<script setup lang="ts">
import { useRoadListStore } from "../store"
import SearchRoad from "./SearchRoad.vue"
import SelectRoad from "./SelectRoad.vue"

useHead({
	title: "สายทาง",
})

interface IStateParams {
	keyword: string
	road_group_id: string[]
	road_section_id: number[]
	road_id: string[]
	km_start: string
	km_end: string
	depot_code: number[]
	ref_surface_id: number[]
	is_iri_1000: boolean | null
	is_iri_100: boolean | null
	is_rut_100: boolean | null
	is_ifi_100: boolean | null
	is_g7_100: boolean | null
}

const store = useRoadListStore()
const router = useRouter()
const route = useRoute()

useStoreLifecycle(store)

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

onMounted(() => {
	if (typeof window !== "undefined") {
		const retrievedParams = window.localStorage.getItem("search")
		if (retrievedParams) {
			if (window.localStorage.getItem("visited")) {
				store.params = JSON.parse(retrievedParams)
			}
			window.addEventListener("beforeunload", () => {
				window.localStorage.removeItem("search")
				store.params = {} as IStateParams
			})
		}

		if (Object.keys(route.query).length) {
			store.setQuriesParams(route.query)
		}

		store.getData()
	}
})

onUnmounted(() => {
	localStorage.setItem("search", JSON.stringify(store.params))
	if (!router.hasRoute("roads")) {
		store.$reset()
	}
})
</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<SearchRoad />
			<SelectRoad :roads="store.roads" :loading="store.loading" />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="store.loading" height="97vh" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
