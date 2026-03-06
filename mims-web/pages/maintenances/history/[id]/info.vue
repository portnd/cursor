<script setup lang="ts">
import TheBreadcrumb from "~/components/layouts/headers/TheBreadcrumb.vue"
import { HistoryProjectInfo } from "~~/core/modules/maintenance/history/ui"

definePageMeta({
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.view_owner_maint_history,
				IUserRolesAccess.view_all_maint_history,
			])
		},
	],
})

const router = useRouter()

const previousPage = () => {
	router.push({ path: `/maintenances/history` })
}

useHead({
	title: "ข้อมูลโครงการ",
})
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
			<TheBreadcrumb :breadcrumbs="[{ name: 'ประวัติการซ่อมบำรุง' }]" title="รายละเอียด" />
		</div>
	</div>
	<HistoryProjectInfo />
</template>

<style scoped></style>
