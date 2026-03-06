<script setup lang="ts">
import { useAnnualAnalyseCreateStore } from "~/core/modules/analysis/annual/store"
import { AnnualSelectRoad } from "~/core/modules/analysis/annual/ui/create"

const store = useAnnualAnalyseCreateStore()

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.manage_myself_maintenance_analysis)
		},
	],
})

const title = computed(() => {
	return store.step === 1 ? "บำรุงรักษาประจำปี (เลือกสายทาง)" : "บำรุงรักษาประจำปี (กำหนดเงื่อนไข)"
})

useHead({
	title: title.value,
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' }]" title="บำรุงรักษาประจำปี" />
	<AnnualSelectRoad />
</template>

<style scoped></style>
