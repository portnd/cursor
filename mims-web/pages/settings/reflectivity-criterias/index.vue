<script lang="ts" setup>
import { ReflectivityList } from "~~/core/modules/setting/reflectivityCriteria/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_retro_criteria_access)
		},
	],
})

useHead({
	title: "ตั้งค่าเกณฑ์ค่าการสะท้อนแสง",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="เกณฑ์ค่าการสะท้อนแสง" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<ReflectivityList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
