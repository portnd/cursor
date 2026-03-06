<script lang="ts" setup>
import { BudgetList } from "~~/core/modules/setting/budget/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_ref_budget_access)
		},
	],
})

useHead({
	title: "ตั้งค่าประเภทงบประมาณ",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="ประเภทงบประมาณ" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<BudgetList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
