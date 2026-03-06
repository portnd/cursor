<script setup lang="ts">
import { useMaintenanceHistoryPlanStore } from "~/core/modules/maintenance/history/store"
import { HistoryProjectPlan } from "~~/core/modules/maintenance/history/ui"

useHead({
	title: "แผนการดำเนินงาน",
})

const store = useMaintenanceHistoryPlanStore()
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'ติดตามการซ่อมบำรุง', to: '/maintenances/follow-up' },
			{ name: 'ประวัติการซ่อมบำรุง', to: '/maintenances/history' },
		]"
		:title="store.projectInfo.name"
	/>
	<HistoryProjectPlan />
</template>

<style scoped></style>
