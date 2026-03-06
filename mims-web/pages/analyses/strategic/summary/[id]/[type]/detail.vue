<script setup lang="ts">
import StrategicSummaryDetail from "~/core/modules/analysis/strategic/ui/summary/StrategicSummaryDetail.vue"
useHead({
	title: "รายละเอียด",
})

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.manage_myself_maintenance_analysis)
		},
	],
})

const route = useRoute()
</script>
<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' },
			{ name: 'บำรุงรักษาเชิงกลยุทธ์', to: `/analyses/strategic/summary/${route.params.id}/${route.params.type}` },
		]"
		title="รายละเอียด"
	/>
	<StrategicSummaryDetail />
</template>

<style scoped></style>
;
