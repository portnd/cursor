<script lang="ts" setup>
import { ModelDeteriorationConcrete } from "~~/core/modules/setting/model/deterioration/ui"

useHead({
	title: "ตั้งค่าแบบจำลอง (Deterioration)",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'แบบจำลอง' }, { name: 'Deterioration' }]"
		title="คอนกรีต (Concrete)"
	/>
	<ModelDeteriorationConcrete />
</template>

<style scoped></style>
