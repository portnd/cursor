<script lang="ts" setup>
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"
import { RoadAssetList } from "~~/core/modules/road/info/asset/ui"

definePageMeta({
	validate: (route) => {
		if (!isNumber(route.params.roadId) || !isNumber(route.params.assetId)) {
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

const route = useRoute()
const roadId = route.params.roadId
const assetId = route.params.assetId

const assetName = computed(() => {
	const name = useInitData()
		.refAssetTable()
		?.find((item) => item.id === Number(assetId))?.table_label
	return name ?? ""
})

useHead({
	title: `${assetName}`,
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[
			{ name: 'สายทาง', to: '/roads' },
			{ name: roadTitleStore.data.road_section_name_th, to: `/roads/${roadId}` },
			{ name: 'สินทรัพย์นอกเขตทาง', to: `/roads/${roadId}/out-asset` },
		]"
		:title="assetName"
	/>
	<RoadAssetList asset-type="out" />
</template>

<style scoped></style>
