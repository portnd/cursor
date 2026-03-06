<script setup lang="ts">
import { MaintenanceHistoryCreate } from "~/core/modules/maintenance/history/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.manage_all_maint_history,
				IUserRolesAccess.manage_owner_maint_history,
			])
		},
	],
})

useHead({
	title: "เพิ่มโครงการ",
})

const router = useRouter()

const previousPage = () => {
	router.push({ path: `/maintenances/history` })
}
</script>

<template>
	<div class="row mb-3">
		<div class="col-auto">
			<button
				type="button"
				class="btn btn-primary rounded-4 px-4 fw-semibold"
				style="font-size: 12px; height: fit-content"
				@click="previousPage"
			>
				&lt; ย้อนกลับ
			</button>
		</div>
		<div class="col-auto d-flex align-self-end" style="margin-top: 0.7em">
			<TheBreadcrumb :breadcrumbs="[{ name: 'ประวัติการซ่อมบำรุง' }]" title="เพิ่มประวัติการซ่อมบำรุง" />
		</div>
	</div>

	<MaintenanceHistoryCreate />
</template>

<style scoped></style>
