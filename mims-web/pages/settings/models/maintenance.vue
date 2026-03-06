<script setup lang="ts">
import { MainMaintenance } from "~~/core/modules/setting/model/maintenance/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Intervention Criteria)",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }]" title="Intervention Criteria" />
	<MainMaintenance />
</template>

<style scoped></style>
