<script setup lang="ts">
import { RoadWorkEffectConcrete } from "~~/core/modules/setting/model/road-work-effect/ui"
useHead({
	title: "ตั้งค่าแบบจำลอง (Road Work Effect)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }, { name: 'Road Work Effect' }]"
		title="คอนกรีต (Concrete)"
	/>

	<RoadWorkEffectConcrete />
</template>

<style scoped></style>
