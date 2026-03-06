<script setup lang="ts">
import { useMaintenanceHistoryAttachmentStore } from "../../../store"
import ProjectFile from "./ProjectFile.vue"

const store = useMaintenanceHistoryAttachmentStore()
useStoreLifecycle(store)

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<div class="row">
				<ProjectFile />
			</div>
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="false" height="calc(97vh)" :is-sticky="true" @map="store.setMap" />
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
