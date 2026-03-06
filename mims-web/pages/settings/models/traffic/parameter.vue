<script setup lang="ts">
import { ModelTrafficParameter } from "~~/core/modules/setting/model/traffic/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Traffic)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }, { name: 'Traffic' }]"
		title="Parameter"
	/>
	<ModelTrafficParameter />
</template>

<style scoped></style>
