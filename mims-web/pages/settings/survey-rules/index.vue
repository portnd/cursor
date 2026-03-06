<script lang="ts" setup>
import { SurveyRuleList } from "~~/core/modules/setting/surveyRule/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_condition_criteria_access)
		},
	],
})

useHead({
	title: "ตั้งค่าเกณฑ์การจำแนกสภาพทาง",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="เกณฑ์การจำแนกสภาพทาง" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<SurveyRuleList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
