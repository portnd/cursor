<script lang="ts" setup>
import { RolesEdit } from "~~/core/modules/users/roles/ui"

useHead({
	title: "แก้ไขสิทธิ์ผู้ใช้งาน",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ผู้ใช้งาน' }, { name: 'ตั้งค่าสิทธิ์ผู้ใช้งาน', to: '/users/roles' }]"
		title="แก้ไขข้อมูล"
	/>
	<RolesEdit />
</template>

<style scoped></style>
