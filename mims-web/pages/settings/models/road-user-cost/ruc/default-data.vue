<script setup lang="ts">
import { ModelRoadUserCostRUC } from "~~/core/modules/setting/model/road-user-cost/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Road User Cost)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'ตั้งค่าระบบ' },
			{ name: 'แบบจำลอง' },
			{ name: 'Road User Cost' },
			{ name: 'Road User Cost (RUC)' },
		]"
		title="ข้อมูลตั้งต้น"
	/>
	<ModelRoadUserCostRUC type="ruc" />
</template>

<style scoped></style>
