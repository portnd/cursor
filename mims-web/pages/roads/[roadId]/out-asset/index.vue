<script lang="ts" setup>
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"
import { RoadAsset } from "~~/core/modules/road/info/asset/ui"

useHead({
	title: "สายทาง (สินทรัพย์นอกเขตทาง)",
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
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'สายทาง', to: '/roads' }, { name: roadTitleStore.data.road_section_name_th }]"
		title="สินทรัพย์นอกเขตทาง"
	/>
	<RoadAsset asset-type="assetout" />
</template>

<style scoped></style>
