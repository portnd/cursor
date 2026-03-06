<script setup lang="ts">
import { useMaintenanceHistoryAttachmentStore } from "~/core/modules/maintenance/history/store"
import { HistoryProjectFile } from "~~/core/modules/maintenance/history/ui"

useHead({
	title: "เอกสารและรูปภาพ",
})
const store = useMaintenanceHistoryAttachmentStore()
useStoreLifecycle(store)
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'ติดตามการซ่อมบำรุง', to: '/maintenances/follow-up' },
			{ name: 'ประวัติการซ่อมบำรุง', to: '/maintenances/history' },
		]"
		:title="store.projectDetais.name"
	/>
	<HistoryProjectFile />
</template>

<style scoped></style>
