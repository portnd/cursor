<script setup lang="ts">
import { IBudget } from "../infrastructure"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import { useInitDataStore } from "~/core/modules/initData/store"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 125 },
	{ text: "ชื่องบประมาณ", value: "name" },
	{ text: "จัดการ", value: "operation", width: 125 },
]

const router = useRouter()

// เพิ่มข้อมูล
const createItem = () => {
	router.push({ path: `/settings/budgets/create` })
}

// แก้ไขข้อมูล
const editItem = (item: IBudget) => {
	router.push({ path: `/settings/budgets/${item.id}/edit` })
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: IBudget) => {
	useDeleteItem({
		name: item.name,
		url: `/settings/budget/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
			const initDataStore = useInitDataStore()
			initDataStore.initData()
		},
	})
}

// ค้นหาข้อมูล
interface IBudgetSearch {
	searchName: string
}

const search: IBudgetSearch = reactive({
	searchName: "",
})

const resetSearch = () => {
	search.searchName = ""
	onSearch()
}

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.searchName,
	})
}
</script>

<template>
	<div class="row mb-3">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.searchName" label="ชื่องบประมาณ" name="name" @keyup.enter="onSearch" />
		</div>
		<div class="col-12 col-md mb-2 d-flex justify-content-between align-items-end">
			<div class="d-flex align-items-end">
				<BtnSearch @click="onSearch" />
				<button
					type="button"
					class="btn btn-outline-primary rounded-4 ms-5 fw-semibold text-gray-700"
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
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/budget">
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
</template>

<style scoped></style>
