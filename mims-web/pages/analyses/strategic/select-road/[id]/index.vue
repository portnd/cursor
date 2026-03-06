<script setup lang="ts">
import { useStrategicEditStore } from "~/core/modules/analysis/strategic/store"
import EditStrategicSelectRoad from "~/core/modules/analysis/strategic/ui/edit/EditStrategicSelectRoad.vue"

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
		:title="store.isCopy ? 'บำรุงรักษาเชิงกลยุทธ์' : 'แก้ไขข้อมูล'"
	/>
	<EditStrategicSelectRoad />
</template>

<style scoped></style>
