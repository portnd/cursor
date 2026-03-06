<script setup lang="ts">
import { IRoadSummaryItem } from "../infrastructure"
import { useRoadSummaryStore } from "../store"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"
import type { THeader } from "~~/core/shared/types/Datatable"
const initUserStore = useInitUserStore()

const store = useRoadSummaryStore()
const route = useRoute()

useStoreLifecycle(store)

const props = defineProps({
	show: {
		type: Boolean,
		default: false,
	},
})

const headers: THeader[] = [
	{ text: "เริ่มต้น", value: "km_start" },
	{ text: "สิ้นสุด", value: "km_end" },
	{ text: "ทิศทาง", value: "lane", width: 180 },
	{ text: "ผิวทาง", value: "width_surface" },
	{ text: "ไหล่ทางซ้าย", value: "width_shoulder_left" },
	{ text: "ไหล่ทางขวา", value: "width_shoulder_right" },
]

const activeRowNumber = ref(1)

const selectRow = (isLane: boolean, item: IRoadSummaryItem, index: number) => {
	store.setSurfaceLocation(isLane, item, index)

	// เพิ่ม class ใส่ใน row นั้น ๆ
	activeRowNumber.value = item.no ?? 0
}

onUpdated(() => {
	if (props.show) {
		const widthArray: any = []
		const left = document.querySelectorAll(".left")
		const right = document.querySelectorAll(".right")
		if (left) {
			left.forEach((item) => {
				const el = item as HTMLElement
				widthArray.push(el.offsetWidth)
			})
			const leftWidth = Math.max(...widthArray)
			left.forEach((item) => {
				const el = item.querySelector("div") as HTMLElement
				el.style.width = leftWidth.toString() + "px"
			})
		}
		if (right) {
			right.forEach((item) => {
				const el = item as HTMLElement
				widthArray.push(el.offsetWidth)
			})
			const rightWidth = Math.max(...widthArray)
			right.forEach((item) => {
				const el = item.querySelector("div") as HTMLElement
				el.style.width = rightWidth.toString() + "px"
			})
		}
	}
})

onBeforeMount(async () => {
	await setTimeout(async () => {
		await store.getRoadSurface(Number(route.params.roadId))
	}, 200)
	await store.getSurfaceIcon()
})

const editItem = () => {
	navigateTo(`/roads/${route.params.roadId}/summary/edit?tab=surface`)
}

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div id="row" class="row mb-3">
			<div class="col-12 col-md-9 text-start align-self-center ps-4 order-last order-md-first">
				<template v-if="store.update_date !== ''">
					<label class="text-gray-900 fs-6 me-1">ปรับปรุงข้อมูลโดย</label>
					<VUser :label="store.update_by.full_name" :name="store.update_by.full_name" :role="'-'" />
					<label class="text-gray-900 fs-6 ms-1">
						เมื่อวันที่
						{{ buddhistFormatDate(new Date(store.update_date), "dd mmm yyyy เวลา HH:ii น.") }}
					</label>
				</template>
			</div>
			<div class="col-12 col-md-3 text-end">
				<NuxtLink
					v-if="
						initUserStore.accessPermissions[IUserRolesAccess.manage_road_summary] ||
						initUserStore.getIsOwnerManagePermission(
							initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_summary],
							store.road.ref_depot.id
						)
					"
					class="btn btn-outline btn-outline-primary rounded-2 px-5 fw-semibold fs-6"
					@click="editItem()"
				>
					ปรับปรุงข้อมูล
				</NuxtLink>
			</div>
		</div>
		<VDatatable
			:headers="headers"
			:items="store.roadSurface"
			:no-border="true"
			active-item-class-name="active"
			:active-row-number="activeRowNumber"
		>
			<template #customize-headers>
				<thead>
					<tr>
						<th colspan="2" class="text-center border-bottom border-1">กม.</th>
						<th rowspan="2" class="text-center">ช่องจราจร</th>
						<th colspan="3" class="text-center border-bottom border-1">ความกว้าง (ม.)</th>
					</tr>
					<tr>
						<th class="text-center border-0">เริ่มต้น</th>
						<th class="text-center">สิ้นสุด</th>
						<th class="text-center">ผิวทาง</th>
						<th class="text-center">ไหล่ทางซ้าย</th>
						<th class="text-center">ไหล่ทางขวา</th>
					</tr>
				</thead>
			</template>

			<!-- begin::Items -->
			<template #item-km_start="{ item }">
				<NuxtLink class="cursor-pointer" @click="selectRow(false, item, 0)">
					<div class="text-center text-black">{{ item.km_start }}</div>
				</NuxtLink>
			</template>
			<template #item-km_end="{ item }">
				<NuxtLink class="cursor-pointer" @click="selectRow(false, item, 0)">
					<div class="text-center text-black">{{ item.km_end }}</div>
				</NuxtLink>
			</template>

			<template #item-lane="{ item }">
				<div class="d-flex justify-content-center">
					<div
						v-for="(lane, index) of item.lane"
						:key="index"
						:style="`background-color: ${lane.surface.color_code === undefined ? '' : lane.surface.color_code};`"
						class="fw-semibold px-5 py-3 mx-1 rounded-xs"
						:class="lane.surface.color_code === '' ? 'cursor-default' : 'cursor-pointer text-white'"
						@click="selectRow(true, item, index)"
					>
						{{ lane.lane_no }}
					</div>
				</div>
			</template>

			<template #item-width_surface="{ item }">
				<NuxtLink class="cursor-pointer" @click="selectRow(false, item, 0)">
					<div class="text-center text-black">{{ item.width_surface }}</div>
				</NuxtLink>
			</template>
			<template #item-width_shoulder_left="{ item }">
				<NuxtLink class="cursor-pointer" @click="selectRow(false, item, 0)">
					<div class="text-center text-black">
						{{ item.width_shoulder_left }}
					</div>
				</NuxtLink>
			</template>
			<template #item-width_shoulder_right="{ item }">
				<NuxtLink class="cursor-pointer" @click="selectRow(false, item, 0)">
					<div class="text-center text-black">
						{{ item.width_shoulder_right }}
					</div>
				</NuxtLink>
			</template>
			<!-- end::Items -->
		</VDatatable>

		<!-- begin::สัญลักษณ์ผิวทาง -->
		<div class="row mt-6">
			<div class="col-12 text-end">
				<small class="fw-normal fs-6">สัญลักษณ์ผิวทาง: </small>
				<span
					v-for="(item, index) in store.surfaceIcon"
					:key="index"
					:style="`background-color: ${item.color_code};`"
					class="badge badge-primary text-white px-5 py-2 mx-1 mb-2 rounded-1 fw-normal"
					>{{ item.name }}</span
				>
			</div>
		</div>
		<!-- end::สัญลักษณ์ผิวทาง -->
	</VSkeletonLoader>
</template>

<style scoped>
.transparent {
	color: transparent;
}
</style>
