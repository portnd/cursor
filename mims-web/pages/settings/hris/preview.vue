<script setup lang="ts">
import { HRISPreview } from "~~/core/modules/setting/hris/ui"

useHead({
	title: "รายการสายทางที่เชื่อมกับระบบ HRIS",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'รายการสายทางที่เชื่อมกับระบบ HRIS', to: '/settings/hris' }]"
		title="Preview"
	/>
	<HRISPreview />
</template>

<style scoped></style>
