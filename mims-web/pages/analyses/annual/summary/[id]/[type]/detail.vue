<script setup lang="ts">
import AnnualSummaryDetail from "~/core/modules/analysis/annual/ui/summary/AnnualSummaryDetail.vue"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.manage_myself_maintenance_analysis)
		},
	],
})
useHead({
	title: "รายละเอียด",
})
const route = useRoute()
</script>
<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' },
			{ name: 'บำรุงรักษาประจำปี', to: `/analyses/annual/summary/${route.params.id}/${route.params.type}` },
		]"
		title="รายละเอียด"
	/>
	<AnnualSummaryDetail />
</template>

<style scoped></style>
;
