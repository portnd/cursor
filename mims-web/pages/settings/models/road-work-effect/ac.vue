<script setup lang="ts">
import { RoadWorkEffectAC } from "~~/core/modules/setting/model/road-work-effect/ui"
useHead({
	title: "ตั้งค่าแบบจำลอง (Road Work Effect)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }, { name: 'Road Work Effect' }]"
		title="ลาดยาง (AC)"
	/>

	<RoadWorkEffectAC />
</template>

<style scoped></style>
