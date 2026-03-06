<script setup lang="ts">
import { BudgetCreate } from "~/core/modules/setting/budget/ui"

useHead({
	title: "เพิ่มประเภทงบประมาณ",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'ประเภทงบประมาณ', to: '/settings/budgets' }]"
		title="เพิ่มข้อมูล"
	/>
	<BudgetCreate />
</template>

<style scoped></style>
