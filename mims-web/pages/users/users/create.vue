<script lang="ts" setup>
import { UsersCreate } from "~~/core/modules/users/users/ui"

useHead({
	title: "เพิ่มข้อมูลผู้ใช้งาน",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ผู้ใช้งาน' }, { name: 'จัดการผู้ใช้งานภายในระบบ', to: '/users/users' }]"
		title="เพิ่มข้อมูล"
	/>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<UsersCreate />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
