<script lang="ts" setup>
const DashboardMain = defineAsyncComponent(
	() => import("~/core/modules/dashboard/ui/DashboardMain.vue")
)

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_total_dashboard,
				IUserRolesAccess.view_all_asset_dashboard,
				IUserRolesAccess.view_owner_asset_dashboard,
				IUserRolesAccess.view_all_road_condition_dashboard,
				IUserRolesAccess.view_owner_road_condition_dashboard,
				IUserRolesAccess.view_all_surface_dashboard,
				IUserRolesAccess.view_owner_surface_dashboard,
				IUserRolesAccess.view_all_maint_history_dashboard,
				IUserRolesAccess.view_owner_maint_history_dashboard,
			])
		},
	],
})

useHead({
	title: "ข้อมูลสรุป",
})
</script>

<template>
	<TheBreadcrumb title="ข้อมูลสรุป" />
	<DashboardMain />
</template>

<style scoped></style>
