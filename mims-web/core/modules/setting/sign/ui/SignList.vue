<script setup lang="ts">
import { ISign } from "../infrastructure"
import { SignCreate, SignEdit } from "./index"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import { useInitDataStore } from "~/core/modules/initData/store"

// ตั้งค่าตาราง
const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 125 },
	{ text: "ชื่อ", value: "name" },
	{ text: "ตัวอักษรย่อ", value: "abbr", width: 175 },
	{ text: "รูปภาพ", value: "sign_image_filepath", width: 175 },
	{ text: "จัดการ", value: "operation", width: 175 },
]
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

const modalCreateKey: Ref = ref(0)
const modalEditKey: Ref = ref(0)

// เพิ่มข้อมูล
const createItem = () => {
	modalCreate.value.showModal()
}

// แก้ไขข้อมูล
const editItem = (item: ISign) => {
	modalEdit.value.showModal(item.id)
}

// ลบข้อมูล
const dataTable = ref()
const deleteItem = (item: ISign) => {
	useDeleteItem({
		name: item.name,
		url: `/settings/signs/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
			const initDataStore = useInitDataStore()
			initDataStore.initData()
		},
	})
}

// ค้นหาข้อมูล
interface ISignsSearch {
	search_name: string
}

const search: ISignsSearch = reactive({
	search_name: "",
})

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.search_name,
	})
}

const resetSearch = async () => {
	await dataTable.value.searchData({
		name: "",
	})
	search.search_name = ""
}
</script>

<template>
	<div class="row mb-3">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.search_name" label="ชื่อเครื่องหมายจราจร" name="name" @keyup.enter="onSearch" />
		</div>
		<div class="col-12 col-md mb-2 d-flex justify-content-between align-items-end">
			<div class="d-flex mb-2 align-self-end">
				<BtnSearch @click="onSearch" />
				<button
					type="button"
					class="btn rounded-4 ms-5 mt-md-0 mt-sm-2 mt-2 fw-semibold text-gray-700"
					@click="resetSearch()"
				>
					รีเซ็ต
				</button>
			</div>
			<div class="mb-2 text-end align-self-end order-first order-md-last">
				<BtnCreate @click="createItem" />
			</div>
		</div>
	</div>

	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/signs">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-name="{ item }">
					<div class="text-center">{{ item.name }}</div>
				</template>
				<template #item-abbr="{ item }">
					<div class="text-center">{{ item.abbr }}</div>
				</template>
				<template #item-sign_image_filepath="{ item }">
					<div v-viewer class="symbol symbol-40px cursor-pointer">
						<img :src="item.sign_image_filepath" class="rounded-1" />
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
	<SignCreate
		ref="modalCreate"
		:data-table="dataTable"
		@on-finish="
			() => {
				modalCreateKey++
			}
		"
	/>
	<SignEdit
		ref="modalEdit"
		:data-table="dataTable"
		@on-finish="
			() => {
				modalEditKey++
			}
		"
	/>
</template>

<style scoped></style>
