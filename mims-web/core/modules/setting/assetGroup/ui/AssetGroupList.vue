<script setup lang="ts">
import { IAssetGroup } from "../infrastructure"
import { AssetGroupCreate, AssetGroupEdit } from "./index"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import { useInitDataStore } from "~/core/modules/initData/store"

// ตั้งค่าตาราง
const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 125 },
	{ text: "ชื่อกลุ่มสินทรัพย์", value: "name" },
	{ text: "จัดการ", value: "operation", width: 125 },
]

const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

// เพิ่มข้อมูล
const createItem = () => {
	modalCreate.value.showModal()
}

// แก้ไขข้อมูล
const editItem = (item: IAssetGroup) => {
	modalEdit.value.showModal(item.id)
}

// ลบข้อมูล
const dataTable = ref()
const deleteItem = (item: IAssetGroup) => {
	useDeleteItem({
		name: item.name,
		url: `/settings/asset_groups/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
			const initDataStore = useInitDataStore()
			initDataStore.initData()
		},
	})
}

// ค้นหาข้อมูล
interface IAssetGroupSearch {
	searchName: string
}

const search: IAssetGroupSearch = reactive({
	searchName: "",
})

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.searchName,
	})
}

const resetSearch = () => {
	search.searchName = ""
	onSearch()
}
</script>

<template>
	<div class="row mb-3">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.searchName" label="ชื่อกลุ่มสินทรัพย์" name="name" @keyup.enter="onSearch" />
		</div>
		<div class="col-12 col-md-9 d-flex justify-content-between align-items-end">
			<div class="d-flex mb-2 align-items-end">
				<BtnSearch @click="onSearch" />
				<button type="button" class="btn rounded-4 ms-5 fw-semibold text-gray-700" @click="resetSearch()">
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
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/asset_groups">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-name="{ item }">
					<div class="text-center">{{ item.name }}</div>
				</template>
				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item)" />
					<BtnDelete v-show="item.can_delete" @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>
	<AssetGroupCreate ref="modalCreate" :data-table="dataTable" />
	<AssetGroupEdit ref="modalEdit" :data-table="dataTable" />
</template>

<style scoped></style>
