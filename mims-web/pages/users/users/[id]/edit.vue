<script lang="ts" setup>
import { UsersEdit } from "~~/core/modules/users/users/ui"

useHead({
	title: "แก้ไขข้อมูลผู้ใช้งาน",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ผู้ใช้งาน' }, { name: 'จัดการผู้ใช้งานภายในระบบ', to: '/users/users' }]"
		title="แก้ไขข้อมูล"
	/>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<UsersEdit />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
