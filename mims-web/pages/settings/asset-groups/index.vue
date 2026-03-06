<script lang="ts" setup>
import { AssetGroupList } from "~~/core/modules/setting/assetGroup/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_asset_place_group_access)
		},
	],
})

useHead({
	title: "ตั้งค่ากลุ่มสินทรัพย์",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="กลุ่มสินทรัพย์" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<AssetGroupList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
