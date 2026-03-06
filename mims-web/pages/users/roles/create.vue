<script lang="ts" setup>
import { RolesCreate } from "~~/core/modules/users/roles/ui"

useHead({
	title: "เพิ่มสิทธิ์ผู้ใช้งาน",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ผู้ใช้งาน' }, { name: 'ตั้งค่าสิทธิ์ผู้ใช้งาน', to: '/users/roles' }]"
		title="เพิ่มรายการ"
	/>
	<RolesCreate />
</template>

<style scoped></style>
