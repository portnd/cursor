<script lang="ts" setup>
definePageMeta({
	middleware: [
		"landing-page",
		// function () {
		// 	const permission = usePermission()
		// 	const hasDashPermission = permission.hasPermission("dashboard_road_condition_access")
		// 	const hasRoadPermissionAll = permission.hasPermission("view_all_road_in_asset")
		// 	const hasRoadPermissionDepartment = permission.hasPermission("view_department_road_in_asset")

		// 	if (hasDashPermission) {
		// 		return navigateTo("/dashboards")
		// 	}

		// 	if (hasRoadPermissionAll || hasRoadPermissionDepartment) {
		// 		return navigateTo("/roads")
		// 	}

		// 	if (!hasDashPermission) {
		// 		return navigateTo("/profile")
		// 	}
		// },
	],
})
</script>

<template>
	<div></div>
</template>

<style scoped></style>
