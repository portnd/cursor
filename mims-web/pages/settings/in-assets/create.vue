<script lang="ts" setup>
import { InAssetCreate } from "~~/core/modules/setting/assetTable/ui"

useHead({
	title: "เพิ่มสินทรัพย์",
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'สินทรัพย์', to: '/settings/in-assets' }]"
		title="เพิ่มรายการ"
	/>

	<InAssetCreate />
</template>

<style scoped></style>
