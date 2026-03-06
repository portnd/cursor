<script setup lang="ts">
import { useStrategicCreateStore } from "~/core/modules/analysis/strategic/store"
import { StrategicSelectRoad } from "~/core/modules/analysis/strategic/ui"

const store = useStrategicCreateStore()

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.manage_myself_maintenance_analysis)
		},
	],
})

const title = computed(() => {
	return store.step === 1 ? "บำรุงรักษาเชิงกลยุทธิ์ (เลือกสายทาง)" : "บำรุงรักษาเชิงกลยุทธิ์ (กำหนดเงื่อนไข)"
})

useHead({
	title: title.value,
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' }]" title="บำรุงรักษาเชิงกลยุทธ์" />
	<StrategicSelectRoad />
</template>

<style scoped></style>
