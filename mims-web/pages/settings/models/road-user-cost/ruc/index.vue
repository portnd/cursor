<script setup lang="ts">
import { ModelRoadUserCostRUC } from "~~/core/modules/setting/model/road-user-cost/ui"

onMounted(() => {
	navigateTo("/settings/models/road-user-cost/ruc/default-data")
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง', to: 'road-user-cost' }]"
		title="RoadUserCost"
	/>
	<ModelRoadUserCostRUC type="ruc" />
</template>

<style scoped></style>
