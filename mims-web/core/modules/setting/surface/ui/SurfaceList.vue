<script setup lang="ts">
import { ISurface } from "../infrastructure"
import { SurfaceCreate, SurfaceEdit } from "./index"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 75 },
	{ text: "รูปแบบวัสดุ", value: "type", width: 150 },
	{ text: "ชนิดวัสดุผิวทาง", value: "name", width: 150 },
	{ text: "Drainage Coefficient", value: "drainage", width: 150 },
	{ text: "Layer Coefficient", value: "layer_coefficient", width: 150 },
	{ text: "a", value: "a", width: 100 },
	{ text: "b", value: "b", width: 100 },
	{ text: "c", value: "c", width: 100 },
	{ text: "CRT", value: "crt", width: 100 },
	{ text: "RRF", value: "rrf", width: 100 },
	{ text: "จัดการ", value: "operation", width: 85 },
]

// Modal
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

// เพิ่มข้อมูล
// const createItem = () => {
// 	modalCreate.value.showModal()
// }

// แก้ไขข้อมูล
const editItem = (item: ISurface) => {
	modalEdit.value.showModal(item.id)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: ISurface) => {
	useDeleteItem({
		name: item.name,
		url: `/settings/ref/surface/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
		},
	})
}
const type = ref()
let data: { id: number; name: string }[] = []
const generateOptionTable = () => {
	if (data.length > 0) {
		data = []
	}
	useInitData()
		.refSurfaceType()
		?.map((e) => data.push({ id: e.id, name: e.name }))
	return data
}

// ค้นหาข้อมูล
interface ISurfaceSearch {
	name: string
	type: string
}

const search: ISurfaceSearch = reactive({
	name: "",
	type: "",
})

const onSearch = async () => {
	const foundItem = data.find((item) => item.id === type.value)
	if (foundItem) {
		search.type = foundItem.name
	} else {
		search.type = ""
	}
	await dataTable.value.searchData({
		name: search.name,
		type: search.type,
	})
}

const resetSearch = () => {
	type.value = ""
	search.name = ""
	search.type = ""
	onSearch()
}

const splitItemC = (c: string) => {
	if (c === "") {
		return ""
	} else {
		const splitted = c.split("^") // แยกตัวอักษรด้วยคอมม่า (,)
		if (splitted[1] !== undefined) {
			return `${toNumber(Number(splitted[0]))}x10<sup>${toNumber(Number(splitted[1]))}</sup>` // สร้างสตริงในรูปแบบของ HTML superscript
		} else {
			return `${toNumber(Number(splitted[0]))}` // สร้างสตริงในรูปแบบของ HTML superscript
		}
	}
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
			<VTextInput v-model="search.name" label="ชนิดวัสดุผิวทาง" placeholder="" name="name" @keyup.enter="onSearch" />
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="type"
				:options="toOptions(generateOptionTable())"
				label="รูปแบบวัสดุ"
				name="group"
				:close-on-select="true"
			/>
		</div>
		<div class="col-12 col-md-3 mb-2 align-self-end">
			<BtnSearch @click="onSearch" />
			<button type="button" class="btn rounded-4 ms-5 fw-semibold text-gray-700" @click="resetSearch()">รีเซ็ต</button>
		</div>
		<!-- <div class="col-12 col-md-3 mb-2 text-end align-self-end order-first order-md-last">
			<BtnCreate @click="createItem" />
		</div> -->
	</form>
	<div></div>
	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="/settings/ref/surface">
				<!-- begin::Items -->
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-type="{ item }">
					<div class="text-center">{{ item.type }}</div>
				</template>
				<template #item-name="{ item }">
					<div class="text-center">{{ item.name }}</div>
				</template>
				<template #item-drainage="{ item }">
					<div class="text-center">{{ generateNumber(item.drainage) }}</div>
				</template>
				<template #item-layer_coefficient="{ item }">
					<div class="text-center">{{ generateNumber(item.layer_coefficient) }}</div>
				</template>
				<template #item-a="{ item }">
					<div class="text-center">{{ generateNumber(item.a) }}</div>
				</template>
				<template #item-b="{ item }">
					<div class="text-center">{{ generateNumber(item.b) }}</div>
				</template>
				<template #item-c="{ item }">
					<div class="text-center" v-html="splitItemC(item.c)"></div>
				</template>
				<template #item-crt="{ item }">
					<div class="text-center">{{ generateNumber(item.crt) }}</div>
				</template>
				<template #item-rrf="{ item }">
					<div class="text-center">{{ generateNumber(item.rrf) }}</div>
				</template>
				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item)" />
					<BtnDelete v-show="item.can_delete" @click="deleteItem(item)" />
				</template>
				<!-- end::Items -->
			</ServerSideDataTable>
		</div>
	</div>

	<!-- Modal -->
	<SurfaceCreate ref="modalCreate" :data-table="dataTable" />
	<SurfaceEdit ref="modalEdit" :data-table="dataTable" />
</template>

<style scoped></style>
