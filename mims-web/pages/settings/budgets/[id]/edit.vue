<script setup lang="ts">
import { BudgetEdit } from "~~/core/modules/setting/budget/ui"
useHead({
	title: "แก้ไขประเภทงบประมาณ",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'ประเภทงบประมาณ', to: '/settings/budgets' }]"
		title="แก้ไขข้อมูล"
	/>
	<BudgetEdit />
</template>
<style scoped></style>
