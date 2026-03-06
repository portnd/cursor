<script setup lang="ts">
import { useDashboardSurfaceStore } from "../store"
import { DashboardSurfaceChart, DashboardSurfaceSearch, DashboardSurfaceTable } from "./index"

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

const store = useDashboardSurfaceStore()
useStoreLifecycle(store)

onMounted(async () => {
	await store.get()
	await store.getRoad()
	store.defaultLocation()
	store.createLine()
})

</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<DashboardSurfaceSearch />
			<div v-show="store.loading" class="row mb-3" style="margin-top: -10px">
				<div class="col-12 text-end">
					<VLoading :loading="store.loading" float="end" />
				</div>
			</div>
			<VSkeletonLoader :loading="store.loading">
				<DashboardSurfaceChart />
				<DashboardSurfaceTable />
			</VSkeletonLoader>
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="store.loading" :is-sticky="true" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
