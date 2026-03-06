<script setup lang="ts">
definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_model_access)
		},
	],
})

onMounted(() => {
	navigateTo("/settings/models/traffic")
})
</script>

<template>
	<div></div>
</template>

<style scoped></style>
