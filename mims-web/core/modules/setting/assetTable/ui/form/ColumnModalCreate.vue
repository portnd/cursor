<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { IRequestColumnAssetTable } from "../../infrastructure"
import { IOption } from "~~/core/shared/types/Option"
import { IValidate } from "~/core/shared/types/Validate"

const props = defineProps({
	columns: {
		type: Array<IRequestColumnAssetTable>,
		default: [],
	},
})
const columnData = ref({
	column_id: 0,
	column_name: "",
	table_name_ref: "",
	component_title: "",
	component_type: "",
	is_required: false,
	is_visible_view: false,
	is_visible_edit: false,
	is_visible_report: false,
})

const options: IOption[] = [
	{
		value: "text",
		label: "text",
	},
	{
		value: "text-number",
		label: "text-number",
	},
	{
		value: "text-km",
		label: "text-km",
	},
	{
		value: "select",
		label: "select",
	},
	{
		value: "text-year",
		label: "text-year",
	},
	{
		value: "datepicker",
		label: "datepicker",
	},
	{
		value: "image",
		label: "image",
	},
]

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	handleReset()
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
	columnData.value = {
		column_id: 0,
		column_name: "",
		table_name_ref: "",
		component_title: "",
		component_type: "",
		is_required: false,
		is_visible_view: false,
		is_visible_edit: false,
		is_visible_report: false,
	}
}

defineRule("duplicate", (value: any) => {
	if (columnData.value.component_type === "image") {
		value = value + "_filepath"
	}
	for (let index = 0; index < props.columns.length; index++) {
		let name = ""
		name = props.columns[index].column_name
		if (value === name) {
			return "ข้อมูลนี้มีอยู่ในระบบแล้ว"
		}
	}
	return true
})

// defineRule("selected", (value: any) => {
// 	if (columnData.value.component_type === "select" || value === undefined) {
// 		return "โปรดระบุ"
// 	}
// 	return true
// })

watch(
	() => columnData.value.component_type,
	(_, oldValue) => {
		// if เพื่อเช็คค่าเก่าเท่ากับ select เพื่อเคลียค่าในกรณีที่ type ที่เลือกอยู่เป็น select แล้วเปลี่ยนเป็น type อื่น pm ต้องการให้เคลียชื่อ
		if (oldValue === "select") {
			columnData.value.column_name = ""
			setTimeout(() => {
				setErrors({})
			}, 20)
		}
	}
)

const validate = computed(() => {
	const validation = {} as IValidate
	validation.column_name = "duplicate|required"
	validation.component_title = "required"
	validation.component_type = "required"
	validation.table_name_ref = columnData.value.component_type === "select" ? "required" : ""
	return validation
})

const emit = defineEmits(["update:column"])
const { handleSubmit, handleReset, resetField, setErrors } = useForm({
	validationSchema: validate,
})
const onSubmit = handleSubmit(() => {
	if (columnData.value.component_type === "image") {
		columnData.value.column_name = columnData.value.column_name + "_filepath"
	}
	emit("update:column", columnData.value)
	hideModal()
	handleReset()
})

// แก้บัค validate errors เมื่อเปลี่ยน table_name_ref
const resetErrors = () => {
	nextTick(() => {
		resetField("table_name_ref", { value: null, errors: undefined })
		resetField("column_name", { value: null, errors: undefined })
		resetField("component_title", { value: null, errors: undefined })
	})
}

watch(
	() => columnData.value.component_type,
	() => {
		resetErrors()
	}
)

const generateOptionTable = () => {
	const data: IOption[] = []
	useInitData()
		.refTableList()
		?.map((e) => data.push({ value: e.ref_name, label: e.ref_desc }))
	return data
}

onUpdated(() => {
	if (columnData.value.component_type === "select" && columnData.value.table_name_ref !== "") {
		if (columnData.value.table_name_ref === null) {
			columnData.value.table_name_ref = ""
		} else {
			columnData.value.column_name = columnData.value.table_name_ref + "_id"
		}
	} else {
		columnData.value.table_name_ref = ""
	}
})

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-column" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">เพิ่มคอลัมน์</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 mb-5">
								<VSelect
									v-model="columnData.component_type"
									:options="options"
									label="ประเภท"
									name="component_type"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
							<div v-show="columnData.component_type === 'select'" class="col-12 mb-3">
								<VSelect
									v-model="columnData.table_name_ref"
									:options="generateOptionTable()"
									label="ตารางอ้างอิง"
									name="table_name_ref"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
							<div class="col-12 mb-5">
								<VLabel label="คอลัมน์" :required="true" />
								<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
									<i class="fi fi-sr-interrogation fs-5"></i>
									<template #content>
										<div>กรอกได้เฉพาะ a-z, 0-9 และเครื่องหมาย _</div>
									</template>
								</VPopover>
								<VTextInput
									v-model="columnData.column_name"
									:readonly="columnData.component_type === 'select' ? true : false"
									name="column_name"
									:text-end="columnData.component_type === 'image' ? '_filepath' : ''"
									:validate-english="true"
									:paste="columnData.component_type === 'select' ? false : true"
								/>
							</div>
							<div class="col-12 mb-5">
								<VTextInput
									v-model="columnData.component_title"
									name="component_title"
									label="คำอธิบาย"
									:required="true"
								/>
							</div>
							<div class="col-12 mb-5">
								<VLabel label="เงื่อนไขเพิ่มเติม" />
								<div class="d-block d-md-flex">
									<VCheckbox
										v-model="columnData.is_required"
										:option="{ label: 'บังคับกรอก' }"
										name="is_required"
										mode="single"
									/>
									<VCheckbox
										v-model="columnData.is_visible_edit"
										:option="{ label: 'เห็นเมื่อแก้ไข' }"
										name="is_visible_edit"
										mode="single"
									/>
									<VCheckbox
										v-model="columnData.is_visible_view"
										:option="{ label: 'เห็นเมื่อดูข้อมูล' }"
										name="is_visible_view"
										mode="single"
									/>
									<VCheckbox
										v-model="columnData.is_visible_report"
										:option="{ label: 'เห็นเมื่อดูรายงาน' }"
										name="is_visible_report"
										mode="single"
									/>
								</div>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="false" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit label="เพิ่ม" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
