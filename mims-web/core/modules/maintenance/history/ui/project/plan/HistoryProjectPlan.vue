<script setup lang="ts">
import { useMaintenanceHistoryPlanStore } from "../../../store/MaintenanceHistoryPlanStore"
import HistoryProjectPlanGraph from "./HistoryProjectPlanGraph.vue"

const store = useMaintenanceHistoryPlanStore()
useStoreLifecycle(store)

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<HistoryProjectPlanGraph />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" height="calc(97vh)" :is-sticky="true" :loading="false" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
