<script lang="ts" setup>
import { HRISList } from "~~/core/modules/setting/hris/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_hris_access)
		},
	],
})

useHead({
	title: "รายการสายทางที่เชื่อมกับระบบ HRIS",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="รายการสายทางที่เชื่อมกับระบบ HRIS" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5"><HRISList /></div>
		</div>
	</div>
</template>

<style scoped></style>
