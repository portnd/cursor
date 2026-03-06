<script lang="ts" setup>
import { SignList } from "~~/core/modules/setting/sign/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_traffic_sign_access)
		},
	],
})

useHead({
	title: "ตั้งค่าเครื่องหมายจราจร",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="เครื่องหมายจราจร" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<SignList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
