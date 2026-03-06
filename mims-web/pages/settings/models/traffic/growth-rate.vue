<script setup lang="ts">
import { ModelTrafficGrowthRate } from "~~/core/modules/setting/model/traffic/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Traffic)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }, { name: 'Traffic' }]"
		title="Growth Rate"
	/>
	<ModelTrafficGrowthRate />
</template>

<style scoped></style>
