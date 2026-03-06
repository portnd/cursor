<script setup lang="ts">
import { RoadSummaryEdit } from "~~/core/modules/road/info/summary/ui"
import { useRoadTitleStore } from "~~/core/modules/common/roadTitle/store"

const route = useRoute()
const roadId = route.params.roadId

definePageMeta({
	validate: (route) => {
		if (!isNumber(route.params.roadId)) {
			return false
		}
		return true
	},
	middleware: [
		function () {
			return usePermission().checkMultipleAccessMiddleware([
				IUserRolesAccess.manage_road_summary,
				IUserRolesAccess.manage_owner_road_summary,
			])
		},
	],
})

const roadTitleStore = useRoadTitleStore()
useStoreLifecycle(roadTitleStore, { resetOnEnter: false })
const router = useRouter()
const name = ref()
const tab = ref()

onMounted(async () => {
	if (!roadTitleStore.data.road_info?.id) {
		if (isNumber(roadId)) {
			await roadTitleStore.getData(Number(roadId))
		}
	}
	tab.value = route.query.tab
	if (tab.value === "surface") {
		name.value = "ข้อมูลผิวทาง"
	} else {
		name.value = "ข้อมูลหน้าตัดผิวทาง"
	}
	useHead({
		title: `ปรับปรุง${name.value}`,
	})
})

const previousPage = () => {
	router.push({ path: `/roads/${roadId}/summary`, query: { tab: tab.value } })
	//
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
						to: `/roads/${roadId}/summary?tab=${tab}`,
					},
				]"
				:title="`ปรับปรุง${name}`"
			/>
		</div>
	</div>
	<RoadSummaryEdit />
</template>

<style scoped></style>
