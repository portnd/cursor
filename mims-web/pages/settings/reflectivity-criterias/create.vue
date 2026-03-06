<script setup lang="ts">
import { ReflectivityCreate } from "~~/core/modules/setting/reflectivityCriteria/ui"

useHead({
	title: "เพิ่มตั้งค่าเกณฑ์ค่าการสะท้อนแสง",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'เกณฑ์ค่าการสะท้อนแสง', to: '/settings/reflectivity-criterias' }]"
		title="เพิ่มข้อมูล"
	/>
	<ReflectivityCreate />
</template>

<style scoped></style>
