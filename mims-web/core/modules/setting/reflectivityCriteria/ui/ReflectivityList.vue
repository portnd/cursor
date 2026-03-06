<script setup lang="ts">
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import { useInitDataStore } from "~/core/modules/initData/store"

export interface ISurveyRuleList {
	id: number
	name: string
}

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 125 },
	{ text: "ชื่อเกณฑ์ค่าการสะท้อนแสง", value: "name", width: 125 },
	{ text: "ประเภท", value: "type", width: 125 },
	{ text: "จัดการ", value: "operation", width: 125 },
]

const initDataStore = useInitDataStore()

// เพิ่มข้อมูล
const createItem = () => {
	return navigateTo(`/settings/reflectivity-criterias/create`)
}

// แก้ไขข้อมูล
const editItem = (item: ISurveyRuleList) => {
	return navigateTo(`/settings/reflectivity-criterias/${item.id}/edit`)
}

// ค้นหาข้อมูล
interface ISurveyRulesSearch {
	name: string
	ref_reflectivity_range_id?: number
}

const search: ISurveyRulesSearch = reactive({
	name: "",
	ref_reflectivity_range_id: undefined,
})

const onSearch = async () => {
	await dataTable.value.searchData({
		name: search.name,
		ref_reflectivity_range_id: search.ref_reflectivity_range_id,
	})
}

// watch(
// 	() => search.ref_reflectivity_range_id,
// 	() => {
// 		onSearch()
// 	}
// )

const resetSearch = async () => {
	await dataTable.value.searchData({
		name: "",
		ref_reflectivity_range_id: undefined,
	})
	search.name = ""
	search.ref_reflectivity_range_id = 0
}

const optionGenerator = () => {
	return toOptions(initDataStore.data?.ref_reflectivity_range)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: ISurveyRuleList) => {
	useDeleteItem({
		name: item.name,
		url: `/settings/owners_road_line/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()

			initDataStore.initData()
		},
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
</script>
<template>
	<form class="row mb-3" @submit.prevent="onSearch" @keyup.enter.prevent="onSearch">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="search.name" label="ชื่อเกณฑ์ค่าการสะท้อนแสง" name="name" @keyup.enter="onSearch" />
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="search.ref_reflectivity_range_id"
				:options="optionGenerator()"
				label="ประเภท"
				:can-clear="true"
				name="type"
				placeholder="ทั้งหมด"
			/>
		</div>
		<div class="col-12 col-md d-flex justify-content-between align-items-end">
			<div class="d-flex mb-2">
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
				<BtnCreate class="btn-outline btn-outline-primary px-8" label="เพิ่ม" @click="createItem" />
			</div>
		</div>
	</form>
	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/owners_road_line">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-name="{ item }">
					<div class="text-center">{{ item.name }}</div>
				</template>
				<template #item-type="{ item }">
					<div class="text-center">{{ item.ref_reflectivity_range?.name }}</div>
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
