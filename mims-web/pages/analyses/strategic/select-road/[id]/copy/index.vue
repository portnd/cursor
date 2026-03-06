<script setup lang="ts">
import { useStrategicEditStore } from "~/core/modules/analysis/strategic/store"
import { CopyStrategicSelectRoad } from "~/core/modules/analysis/strategic/ui/copy"

const store = useStrategicEditStore()

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
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' }]"
		title="บำรุงรักษาเชิงกลยุทธ์
"
	/>
	<CopyStrategicSelectRoad />
</template>

<style scoped></style>
