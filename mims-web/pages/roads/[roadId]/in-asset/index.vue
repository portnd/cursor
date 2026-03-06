<script lang="ts" setup>
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"
import { RoadAsset } from "~~/core/modules/road/info/asset/ui"

useHead({
	title: "สายทาง (สินทรัพย์)",
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
			return usePermission().checkMenuAccessMiddleware(IUserRolesAccess.view_road_in_assets)
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
				title="สินทรัพย์"
			/>
		</div>
	</div>

	<RoadAsset asset-type="assetin" />
</template>

<style scoped></style>
