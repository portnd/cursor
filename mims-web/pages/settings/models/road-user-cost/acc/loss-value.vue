<script setup lang="ts">
import { ModelRoadUserCostACC } from "~~/core/modules/setting/model/road-user-cost/ui"

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
			{ name: 'Accident (ACC)' },
		]"
		title="มูลค่าความสูญเสีย"
	/>
	<ModelRoadUserCostACC type="ruc" />
</template>

<style scoped></style>
