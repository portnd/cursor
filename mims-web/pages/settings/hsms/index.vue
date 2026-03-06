<script lang="ts" setup>
import { HSMSList } from "~~/core/modules/setting/hsms/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_hsms_access)
		},
	],
})

useHead({
	title: "ข้อมูลสินทรัพย์ที่กรอกข้อมูลไม่ถูกหลักการในระบบ HSMS",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]"
		title="ข้อมูลสินทรัพย์ที่กรอกข้อมูลไม่ถูกหลักการในระบบ HSMS"
	/>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5"><HSMSList /></div>
		</div>
	</div>
</template>

<style scoped></style>
