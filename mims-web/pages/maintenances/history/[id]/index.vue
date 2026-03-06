<script setup lang="ts">
definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_owner_maint_history,
				IUserRolesAccess.view_all_maint_history,
			])
		},
	],
})
const route = useRoute()
const id = Number(route.params.id)
onMounted(() => {
	console.log("route.params.id =", id)
	navigateTo(`/maintenances/history/${id}/info`)
})
</script>

<template>
	<div></div>
</template>

<style scoped></style>
