<script setup lang="ts">
import { ColumnModalCreate, ColumnModalEdit } from "../index"
import { IRequestColumnAssetTable } from "../../infrastructure"
import type { THeader } from "~~/core/shared/types/Datatable"
const emit = defineEmits(["update:column", "update:deleteColumn", "update:editColumn"])

defineProps({
	modelValue: {
		type: Array as PropType<IRequestColumnAssetTable[]>,
	},
	typeEdit: {
		type: Boolean,
		default: false,
	},
})
const headers: THeader[] = [
	{ text: "", value: "icon", width: 50 },
	{ text: "คอลัมน์", value: "column_name", align: "start", sortable: false },
	{ text: "คำอธิบาย", value: "component_title", align: "start" },
	{ text: "ประเภท", value: "component_type", align: "start" },
	{ text: "บังคับกรอก", value: "is_required", width: 120 },
	{ text: "เห็นเมื่อแก้ไข", value: "is_visible_edit", width: 120 },
	{ text: "เห็นเมื่อดูข้อมูล", value: "is_visible_view", width: 120 },
	{ text: "เห็นเมื่อดูรายงาน", value: "is_visible_report", width: 120 },
	{ text: "จัดการ", value: "operation", width: 100 },
]

// Modal
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()

// เพิ่มข้อมูล
const createItem = () => {
	modalCreate.value.showModal()
}

// แก้ไขข้อมูล
const editItem = (item: IRequestColumnAssetTable) => {
	modalEdit.value.showModal(item)
}

// ลบข้อมูล
const deleteItem = (item: IRequestColumnAssetTable) => {
	useDeleteItem({
		name: item.column_name,
		callBack: function () {
			emit("update:deleteColumn", item)
		},
	})
}
</script>

<template>
	<div class="d-flex justify-content-end mt-0">
		<button
			type="button"
			class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold"
			@click="createItem"
		>
			<i class="fi fi-rr-plus align-middle fs-8"></i>
			เพิ่มคอลัมน์
		</button>
	</div>
	<VDatatable :headers="headers" :items="modelValue" :no-border="false">
		<!-- begin::Items -->
		<template #item-icon>
			<div class="text-center">
				<inline-svg src="/images/icons/svg/align-right.svg" width="25" height="25"></inline-svg>
			</div>
		</template>
		<template #item-column_name="{ item }">
			<div class="text-start">{{ item.column_name }}</div>
		</template>
		<template #item-component_title="{ item }">
			<div class="text-start">{{ item.component_title }}</div>
		</template>
		<template #item-component_type="{ item }">
			<div class="text-start">{{ item.component_type }}</div>
		</template>
		<template #item-is_required="{ item }">
			<div class="text-center">
				<i v-if="item.is_required" class="fi fi-sr-checkbox align-middle text-primary fs-4 lh-0"></i>
				<i v-else class="fi fi-sr-square text-primary fs-4 lh-0"></i>
			</div>
		</template>
		<template #item-is_visible_edit="{ item }">
			<i v-if="item.is_visible_edit" class="fi fi-sr-checkbox align-middle text-primary fs-4 lh-0"></i>
			<i v-else class="fi fi-sr-square align-middle text-primary fs-4 lh-0"></i>
		</template>
		<template #item-is_visible_view="{ item }">
			<i v-if="item.is_visible_view" class="fi fi-sr-checkbox align-middle text-primary fs-4 lh-0"></i>
			<i v-else class="fi fi-sr-square align-middle text-primary fs-4 lh-0"></i>
		</template>
		<template #item-is_visible_report="{ item }">
			<i v-if="item.is_visible_report" class="fi fi-sr-checkbox align-middle text-primary fs-4 lh-0"></i>
			<i v-else class="fi fi-sr-square align-middle text-primary fs-4 lh-0"></i>
		</template>
		<template #item-operation="{ item }">
			<div v-if="item.is_mandatory">
				<div class="text-center">
					<BtnBan />
				</div>
			</div>
			<div v-else>
				<BtnEdit @click="editItem(item)" />
				<BtnDelete @click="deleteItem(item)" />
			</div>
		</template>
		<!-- end::Items -->
	</VDatatable>

	<ColumnModalCreate
		ref="modalCreate"
		:columns="modelValue"
		@update:column="($column: any)=>{emit('update:column', $column)}"
	/>
	<ColumnModalEdit
		ref="modalEdit"
		:columns="modelValue"
		:disabled="typeEdit"
		@update:column="($column: any)=>{emit('update:editColumn', $column)}"
	/>
</template>

<style scoped></style>
