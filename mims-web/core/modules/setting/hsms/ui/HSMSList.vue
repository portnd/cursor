<script setup lang="ts">
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

// "asset_name": "string",
//         "depot_name": "string",
//         "id": 0,
//         "km": "string",
//         "km_range": "string",
//         "location_name": "string",
//         "location_type_name": "string",
//         "road_group_name": "string",
//         "road_name": "string",
//         "type": "string"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 75 },
	{ text: "รายการสินทรัพย์", value: "asset_name", width: 150 },
	{ text: "ทางหลวง", value: "depot_name", width: 150 },
	{ text: "ตอนควบคุม", value: "road_group_name", width: 150 },
	{ text: "จุดที่ตั้ง กม.", value: "km", width: 150 },
	{ text: "ช่วง กม. (จาก-ถึง)", value: "km_range", width: 150 },
	{ text: "ชื่อบริเวณที่ตั้ง", value: "location_type_name", width: 150 },
	{ text: "บริเวณที่ตั้ง", value: "location_name", width: 150 },
	{ text: "หมวด", value: "type", width: 150 },
	{ text: "จัดการ", value: "operation", width: 85 },
]

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: any) => {
	useDeleteItem({
		name: `รายการสินทรัพย์: ${item.asset_name}`,
		url: `/settings/hsms/${item.type}/table/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
		},
	})
}

onMounted(() => {})

onUnmounted(() => {})
</script>

<template>
	<div class="row mb-3">
		<div class="col-12 col-md-4 mb-2 align-self-end"></div>
		<div class="col-12 mt-8 text-end">
			<!-- <BtnCreate @click="createItem" /> -->
		</div>
	</div>

	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/hsms">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-asset_name="{ item }">
					<div class="text-center">{{ item.asset_name }}</div>
				</template>
				<template #item-depot_name="{ item }">
					<div class="text-center">{{ item.depot_name }}</div>
				</template>
				<template #item-road_group_name="{ item }">
					<div class="text-center">{{ item.road_group_name }}</div>
				</template>

				<template #item-km="{ item }">
					<div class="text-center">{{ item.km }}</div>
				</template>
				<template #item-km_range="{ item }">
					<div class="text-center">{{ item.km_range }}</div>
				</template>
				<template #item-location_type_name="{ item }">
					<div class="text-center">{{ item.location_type_name }}</div>
				</template>
				<template #item-location_name="{ item }">
					<div class="text-center">{{ item.location_name }}</div>
				</template>
				<template #item-type="{ item }">
					<div class="text-center">{{ item.type }}</div>
				</template>

				<template #item-operation="{ item }">
					<!-- <BtnEdit @click="editItem(item)" /> -->
					<BtnDelete @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>
</template>

<style scoped></style>
