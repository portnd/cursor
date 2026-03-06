<script lang="ts" setup>
import { UsersList } from "~~/core/modules/users/users/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_user_access)
		},
	],
})

useHead({
	title: "จัดการผู้ใช้งานภายในระบบ",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ผู้ใช้งาน' }]" title="จัดการผู้ใช้งานภายในระบบ" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<UsersList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
