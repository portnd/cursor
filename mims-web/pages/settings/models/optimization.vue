<script setup lang="ts">
import { ModelOptimization } from "~~/core/modules/setting/model/optimization/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Optimization)",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }]" title="Optimization" />
	<ModelOptimization />
</template>

<style scoped></style>
