<script setup lang="ts">
import RoadList from "~~/core/modules/road/roadList/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_all_road,
				IUserRolesAccess.view_owner_road,
			])
		},
	],
})

useHead({
	title: "สายทาง",
})
</script>

<template>
	<TheBreadcrumb title="สายทาง" />

	<RoadList />
</template>

<style scoped></style>
