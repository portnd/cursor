<script lang="ts" setup>
import { RolesList } from "~~/core/modules/users/roles/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_permission_access)
		},
	],
})

useHead({
	title: "ตั้งค่าสิทธิ์ผู้ใช้งาน",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ผู้ใช้งาน' }]" title="ตั้งค่าสิทธิ์ผู้ใช้งาน" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<RolesList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
