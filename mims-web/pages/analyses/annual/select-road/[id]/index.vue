<script setup lang="ts">
import { useAnnualAnalyseEditStore } from "~/core/modules/analysis/annual/store"
import { EditAnnualSelectRoad } from "~/core/modules/analysis/annual/ui/edit"

const store = useAnnualAnalyseEditStore()

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
	<EditAnnualSelectRoad />
</template>

<style scoped></style>
