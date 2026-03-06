<script setup lang="ts">
import { ReflectivityEdit } from "~~/core/modules/setting/reflectivityCriteria/ui"

useHead({
	title: "แก้ไขเกณฑ์ค่าการสะท้อนแสง",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'เกณฑ์ค่าการสะท้อนแสง', to: '/settings/reflectivity-criterias' }]"
		title="แก้ไขข้อมูล"
	/>
	<ReflectivityEdit />
</template>

<style scoped></style>
