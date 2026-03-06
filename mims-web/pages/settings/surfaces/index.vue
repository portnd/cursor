<script lang="ts" setup>
import { SurfaceList } from "~~/core/modules/setting/surface/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.setting_surface_type_access)
		},
	],
})

useHead({
	title: "ตั้งค่าชนิดผิวทาง",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'ตั้งค่าระบบ' }]" title="ชนิดผิวทาง" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<SurfaceList />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
