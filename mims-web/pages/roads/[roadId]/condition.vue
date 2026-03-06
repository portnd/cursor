<script lang="ts" setup>
import { RoadCondition } from "~~/core/modules/road/info/condition/ui"
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"

useHead({
	title: "สายทาง (สภาพทาง)",
})

definePageMeta({
	validate: (route) => {
		if (!isNumber(route.params.roadId)) {
			return false
		}
		return true
	},
	middleware: [
		function () {
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.view_road_condition)
		},
	],
})

const roadTitleStore = useRoadTitleStore()
useStoreLifecycle(roadTitleStore, { resetOnEnter: false })
const router = useRouter()

const previousPage = () => {
	router.push({ path: `/roads` })
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
			<TheBreadcrumb
				:breadcrumbs="[
					{ name: 'สายทาง', to: '/roads' },
					{
						name:
							roadTitleStore.data.road_level === 2
								? roadTitleStore.data.road_info?.name
								: roadTitleStore.data.road_section_name_th,
					},
				]"
				title="สภาพทาง"
			/>
		</div>
	</div>
	<RoadCondition />
</template>

<style scoped></style>
