<script setup lang="ts">
import { IHRISItem } from "../infrastructure"
import { HRISCreate, HRISEdit } from "./index"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 75 },
	{ text: "รหัสสายทาง", value: "road_number", width: 150 },
	{ text: "รหัสตอนควบคุม", value: "section_road_number", width: 150 },
	{ text: "รหัสสำนักงาน", value: "office_of_highways_code", width: 150 },
	{ text: "สถานะ", value: "status", width: 150 },
	{ text: "จัดการ", value: "operation", width: 85 },
]

// Modal
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

// key
const createKey: Ref = ref(0)

const editKey: Ref = ref(0)

// เพิ่มข้อมูล
const createItem = () => {
	modalCreate.value.showModal()
}

// แก้ไขข้อมูล
const editItem = (item: IHRISItem) => {
	modalEdit.value.showModal(item.id)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: IHRISItem) => {
	useDeleteItem({
		name: `รหัสสายทาง ${item.road_number}`,
		url: `/settings/hris/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
		},
	})
}

const onCreateFinish = () => {
	createKey.value++
	setTimeout(() => {
		dataTable.value.loadData()
	}, 500)
}

const onUpdateFinish = () => {
	editKey.value++
	setTimeout(() => {
		dataTable.value.loadData()
	}, 500)
}

const preview = () => {
	navigateTo("/settings/hris/preview")
}

onMounted(() => {})

onUnmounted(() => {})
</script>

<template>
	<div class="row mb-3">
		<div class="col-12 col-md-4 mb-2 align-self-end"></div>
		<div class="col-12 mt-8 text-end">
			<button
				type="button"
				class="btn btn-outline btn-outline-primary rounded-2 px-5 mx-5 fw-semibold fs-6"
				@click="preview()"
			>
				Preview
			</button>
			<BtnCreate @click="createItem" />
		</div>
	</div>

	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/hris">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-road_number="{ item }">
					<div class="text-center">{{ item.road_number }}</div>
				</template>
				<template #item-section_road_number="{ item }">
					<div class="text-center">{{ item.section_road_number }}</div>
				</template>
				<template #item-office_of_highways_code="{ item }">
					<div class="text-center">{{ item.office_of_highways_code }}</div>
				</template>
				<template #item-status="{ item }">
					<div class="text-center">
						<span
							:style="`background-color: ${item.status ? '#E8FFF3' : '#FFF5F8'}`"
							class="badge badge-primary px-3 py-2 rounded-pill fw-normal"
						>
							{{ item.status ? "Active" : "Inactive" }}
						</span>
					</div>
				</template>

				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item)" />
					<BtnDelete @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>

	<!-- Modal -->
	<HRISCreate ref="modalCreate" :key="createKey" @on-finish="onCreateFinish()" />
	<HRISEdit ref="modalEdit" :key="editKey" @on-finish="onUpdateFinish()" />
</template>

<style scoped></style>
