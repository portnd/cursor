<script lang="ts" setup>
import { ReportMain } from "~~/core/modules/report/ui/index"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.download_all_asset_report,
				IUserRolesAccess.download_owner_asset_report,
				IUserRolesAccess.download_all_map_asset_report,
				IUserRolesAccess.download_owner_map_asset_report,
				IUserRolesAccess.download_all_result_asset_report,
				IUserRolesAccess.download_owner_result_asset_report,
				IUserRolesAccess.download_all_road_report,
				IUserRolesAccess.download_all_surface_report,
				IUserRolesAccess.download_owner_surface_report,
				IUserRolesAccess.download_all_road_condition_report,
				IUserRolesAccess.download_owner_road_condition_report,
				IUserRolesAccess.download_all_road_condition_result_report,
				IUserRolesAccess.download_owner_road_condition_result_report,
				IUserRolesAccess.download_all_road_retro_report,
				IUserRolesAccess.download_owner_road_retro_report,
				IUserRolesAccess.download_all_road_retro_result_report,
				IUserRolesAccess.download_owner_road_retro_result_report,
				IUserRolesAccess.download_all_road_damage_report,

				IUserRolesAccess.download_owner_road_damage_report,
				IUserRolesAccess.download_all_road_damage_result_report,

				IUserRolesAccess.download_owner_road_damage_result_report,
				IUserRolesAccess.download_all_maint_history_report,

				IUserRolesAccess.download_owner_maint_history_report,
				IUserRolesAccess.download_all_maint_kpi_report,
				IUserRolesAccess.download_owner_maint_kpi_report,
			])
		},
	],
})

useHead({
	title: "รายงาน",
})
</script>

<template>
	<TheBreadcrumb title="รายงาน" />
	<ReportMain />
</template>

<style scoped></style>
