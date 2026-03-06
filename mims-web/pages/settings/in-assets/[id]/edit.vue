<script lang="ts" setup>
import { InAssetEdit } from "~~/core/modules/setting/assetTable/ui"

useHead({
	title: "แก้ไขสินทรัพย์",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'สินทรัพย์', to: '/settings/in-assets' }]"
		title="แก้ไขข้อมูล"
	/>

	<InAssetEdit />
</template>

<style scoped></style>
