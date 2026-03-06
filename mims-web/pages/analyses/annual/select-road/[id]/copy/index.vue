<script setup lang="ts">
import { useAnnualAnalyseCopyStore } from "~/core/modules/analysis/annual/store"
import { CopyAnnualSelectRoad } from "~/core/modules/analysis/annual/ui/copy"

const store = useAnnualAnalyseCopyStore()

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
	<CopyAnnualSelectRoad />
</template>

<style scoped></style>
