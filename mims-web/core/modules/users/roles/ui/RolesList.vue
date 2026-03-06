<script setup lang="ts">
import { IRoles } from "../infrastructure"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 125 },
	{ text: "กลุ่มสิทธิ์การใช้งาน", value: "name" },
	{ text: "จัดการ", value: "operation", width: 125 },
]

// เพิ่มข้อมูล
const createItem = () => {
	return navigateTo(`/users/roles/create`)
}

// แก้ไขข้อมูล
const editItem = (item: IRoles) => {
	return navigateTo(`/users/roles/${item.id}/edit`)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: IRoles) => {
	useDeleteItem({
		name: item.name,
		url: `/roles/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
		},
	})
}

// ค้นหาข้อมูล
interface IDepartmentSearch {
	name: string
}

const search: IDepartmentSearch = reactive({
	name: "",
})

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.name,
	})
}
</script>

<template>
	<div class="row mb-3">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.name" label="ชื่อกลุ่มสิทธิ์การใช้งาน" name="name" />
		</div>
		<div class="col-12 col-md mb-2 align-self-end d-flex justify-content-between align-items-end">
			<BtnSearch @click="onSearch" />
			<BtnCreate @click="createItem" />
		</div>
		<!-- <div class="col-12 col-md-6 mb-2 text-end align-self-end order-first order-md-last">
			<BtnCreate @click="createItem" />
		</div> -->
	</div>

	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/roles">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-name="{ item }">
					<div class="text-center">{{ item.name }}</div>
				</template>

				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item)" />
					<BtnDelete @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>
</template>

<style scoped></style>
