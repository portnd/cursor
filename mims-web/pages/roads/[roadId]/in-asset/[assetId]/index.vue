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
const router = useRouter()

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
	title: assetName,
})

useStoreLifecycle(roadTitleStore, { resetOnEnter: false })

const previousPage = () => {
	router.push({ path: `/roads/${roadId}/in-asset` })
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
						to: `/roads/${roadId}`,
					},
					{ name: 'สินทรัพย์', to: `/roads/${roadId}/in-asset` },
				]"
				:title="assetName"
			/>
		</div>
	</div>
	<RoadAssetList asset-type="in" />
</template>

<style scoped></style>
