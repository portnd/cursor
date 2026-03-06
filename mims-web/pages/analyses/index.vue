<script setup lang="ts">
import { AnalysisList } from "~/core/modules/analysis/list/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_owner_maintenance_analysis,
				IUserRolesAccess.view_all_maintenance_analysis,
				IUserRolesAccess.manage_myself_maintenance_analysis,
			])
		},
	],
})

useHead({
	title: "วิเคราะห์การซ่อมบำรุง",
})
</script>

<template>
	<TheBreadcrumb title="วิเคราะห์การซ่อมบำรุง" />
	<AnalysisList />
</template>

<style scoped></style>
