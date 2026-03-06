<script setup lang="ts">
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import { useInitDataStore } from "~/core/modules/initData/store"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 50 },
	{ text: "ชื่อสินทรัพย์", value: "table_label", width: 300 },
	{ text: "กลุ่มสินทรัพย์", value: "asset_group", width: 300 },
	{ text: "ชื่อตารางข้อมูล", value: "table_name", width: 300 },
	{ text: "จัดการ", value: "operation", width: 50 },
]

// เพิ่มข้อมูล
const createItem = () => {
	return navigateTo(`/settings/in-assets/create`)
}

// แก้ไขข้อมูล
const editItem = (item: any) => {
	return navigateTo(`/settings/in-assets/${item.id}/edit`)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: any) => {
	useDeleteItem({
		name: item.table_name,
		url: `/settings/asset_tables/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
			const initDataStore = useInitDataStore()
			initDataStore.initData()
		},
	})
}

// สร้าง options VSelect
const groupType = ref("")
const generateOptionTable = () => {
	const data: { id: number; name: string }[] = []
	useInitData()
		.refAsset()
		?.map((e) => data.push({ id: e.id, name: e.name }))
	return data
}

// ค้นหาข้อมูล
interface IAssetTableSearch {
	search_name: string
	group_type: string
	// asset_type: string
}

const search: IAssetTableSearch = reactive({
	search_name: "",
	group_type: "",
	// asset_type: "in",
})

const resetSearch = async () => {
	groupType.value = ""
	search.search_name = ""
	search.group_type = ""

	await onSearch()
}

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.search_name,
		group_id: groupType.value,
		// asset_type: search.asset_type,
	})
}

const handleKeyDown = (event: any) => {
	if (event.key === "Enter") {
		onSearch()
	}
}

onMounted(() => {
	window.addEventListener("keydown", handleKeyDown)
})

onUnmounted(() => {
	window.removeEventListener("keydown", handleKeyDown)
})

// watch(
// 	() => groupType.value,
// 	() => {
// 		onSearch()
// 	}
// )
</script>

<template>
	<form class="row mb-3" @submit.prevent="onSearch" @keyup.enter.prevent="onSearch">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.search_name" label="ชื่อสินทรัพย์" name="name" />
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="groupType"
				:options="toOptions(generateOptionTable())"
				label="กลุ่มสินทรัพย์"
				name="group"
				:close-on-select="true"
				@keyup.enter="onSearch"
			/>
		</div>
		<div class="col-12 col-md-6 d-flex justify-content-between align-items-end">
			<div class="mb-2">
				<BtnSearch @click="onSearch" />
				<!-- <button type="submit" class="btn btn-primary rounded-4" @submit.prevent="onSearch">ค้นหา</button> -->
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
	</form>

	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable
				ref="dataTable"
				:headers="headers"
				url="/settings/asset_tables"
				:additional-params="{ asset_type: 'in' }"
			>
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-table_label="{ item }">
					<div class="text-start">{{ item.table_label }}</div>
				</template>
				<template #item-asset_group="{ item }">
					<div class="text-center">{{ item.asset_group }}</div>
				</template>
				<template #item-table_name="{ item }">
					<div class="text-start">{{ item.table_name }}</div>
				</template>
				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item)" />
					<BtnDelete v-show="item.can_delete" @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>
</template>

<style scoped></style>
