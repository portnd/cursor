<script lang="ts" setup>
import { InAssetList } from "~~/core/modules/setting/assetTable/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_in_asset_access)
		},
	],
})

useHead({
	title: "ตั้งค่าสินทรัพย์",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="สินทรัพย์" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<InAssetList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
