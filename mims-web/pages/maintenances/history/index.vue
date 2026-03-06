<script setup lang="ts">
import { MaintenanceHistory } from "~~/core/modules/maintenance/history/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_all_maint_history,
				IUserRolesAccess.view_owner_maint_history,
			])
		},
	],
})

useHead({
	title: "ประวัติการซ่อมบำรุง",
})
</script>

<template>
	<TheBreadcrumb title="ประวัติการซ่อมบำรุง" />
	<MaintenanceHistory />
</template>

<style scoped></style>
