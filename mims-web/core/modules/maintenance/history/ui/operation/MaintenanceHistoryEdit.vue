<script setup lang="ts">
import { useForm } from "vee-validate"
import { useMaintenanceHistoryEditStore } from "../../store"

const route = useRoute()
const id = Number(route.params.id)
const store = useMaintenanceHistoryEditStore()
useStoreLifecycle(store)

onMounted(() => {
	store
		.getMaintenanceBudget()
		.then(() => store.getDivisiontOptions())
		.then(() => store.getDefault(id))
})

const { handleReset, errors, isSubmitting, handleSubmit } = useForm()

const onSubmit = handleSubmit(async (_, action) => {
	useAction(action)

	const res = await store.updateMaintenanceData(id)

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo(`/maintenances/history/${id}/info`)
			},
		})
	}
})

const onDelete = () => {
	useDeleteItem({
		name: store.defaultData.name,
		url: `maintenance/${id}`,
		callBack: () => {
			navigateTo(`/maintenances/history`)
		},
	})
}

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onCancel = () => {
	navigateTo(`/maintenances/history/${id}/info`)
}

const checkPermission = computed(() => {
	const initUser = useInitUser()
	const check = initUser?.access_control.some((item) => item.access_key === "manage_owner_maint_history")

	return check
})

onUnmounted(() => {
	handleReset()
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<!-- begin::Form -->
				<form @submit="onSubmit">
					<div class="row">
						<div class="col-md-12">
							<h4 class="fw-semibold">ข้อมูลโครงการ</h4>
							<div class="row mt-3">
								<div class="col-md-4 col-12 mb-3">
									<VTextInput v-model="store.defaultData.name" label="ชื่อโครงการ" :required="true" name="name" />
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VTree
										v-model="store.defaultData.owner_code"
										label="หน่วยงาน"
										name="owner_code"
										:mode="checkPermission ? 'LEAF_PRIORITY' : 'ALL'"
										:disable-branch-nodes="checkPermission"
										:options="store.getDivisionOption()"
										:required="true"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VTextInput v-model="store.defaultData.contract_number" label="เลขที่สัญญา" name="contract_number" />
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VDatePicker
										v-model="store.defaultData.budget_year"
										label="ปีงบประมาณ"
										:year-picker="true"
										:required="true"
										name="budget_year"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VSelect
										v-model="store.budgetId"
										:options="store.getBudgetOptions"
										label="ประเภทงบประมาณ"
										placeholder="เลือก"
										:required="true"
										name="budget_id"
										:close-on-select="true"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VSelect
										v-model="store.budgetMethodId"
										:options="store.getMaintenanceCriteriaOptions"
										label="ประเภทการซ่อมบำรุง"
										placeholder="เลือก"
										:required="true"
										name="budget_method_id"
										:close-on-select="true"
									/>
								</div>

								<div class="col-md-4 col-12 mb-3">
									<VNumberInput
										v-model="store.defaultData.budget_maintenance"
										label="วงเงินงบประมาณ (ไม่รวม vat)"
										name="budget_maintenance"
										:precision="3"
										text-end="บาท"
										align="start"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VNumberInput
										v-model="store.defaultData.middle_price"
										label="ราคากลาง (รวม vat)"
										name="middle_price"
										:precision="3"
										text-end="บาท"
										align="start"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VNumberInput
										v-model="store.defaultData.contract_work_value"
										label="มูลค่างานตามสัญญา (รวม vat)"
										name="contract_work_value"
										:precision="3"
										text-end="บาท"
										align="start"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VNumberInput
										v-model="store.defaultData.budget_procurement"
										label="ราคาที่จัดซื้อจัดจ้าง"
										name="budget_procurement"
										:precision="3"
										text-end="บาท"
										align="start"
										:required="true"
									/>
								</div>

								<div class="col-md-4 col-12 mb-3">
									<VTextInput
										v-model="store.defaultData.contractor_name"
										label="บริษัทที่ปรึกษาโครงการ"
										name="contractor_name"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VTextInput v-model="store.defaultData.advisor_name" label="บริษัทผู้รับจ้าง" name="advisor_name" />
								</div>

								<div class="col-md-4 col-12 mb-3">
									<VTextInput
										v-model="store.defaultData.project_secretary_name"
										label="ชื่อ-นามสกุลเจ้าหน้าที่เลขาโครงการ"
										name="project_secretary_name"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VDatePicker
										v-model="store.defaultData.project_end_date"
										label="วันที่สิ้นสุดโครงการ"
										name="project_end_date"
										:max-date="null"
										:required="true"
									/>
								</div>
								<div class="col-md-4 col-12 mb-3">
									<VDatePicker
										v-model="store.defaultData.guarantee_expiration_date"
										label="วันที่หมดค้ำประกัน"
										name="guarantee_expiration_date"
										:max-date="null"
										:required="true"
									/>
								</div>

								<div class="col-md-4 col-12 mb-3">
									<VTextarea
										v-model="store.defaultData.project_details"
										label="รายละเอียดโครงการ"
										name="project_details"
										min-height="100px"
									/>
								</div>
								<div class="col-md col-12 mb-3">
									<VUploadFile
										ref="upLoadFile"
										v-model="store.files"
										:files="store.filesPath"
										label="เอกสารขอการเบิกจ่าย (รองรับ 20 ไฟล์)"
										max-file-size="10MB"
										name="attachments"
										aspect-ratio="0.4"
										:multiple="true"
										:max-files="20"
										:accepted-file-types="[
											'image/png',
											'image/jpg',
											'image/jpeg',
											'application/pdf',
											'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
											'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
											'.dwg',
										]"
									/>
								</div>
							</div>
						</div>
					</div>

					<div class="d-flex justify-content-between mt-5">
						<span>
							<button type="button" class="btn btn-outline-danger btn-delete" @click="onDelete">ลบโครงการ</button>
						</span>
						<span>
							<BtnCancel @click="onCancel" />
							<BtnSubmit label="บันทึก" :loading="store.submitLoading" :disabled="store.submitLoading" />
						</span>
					</div>
				</form>
				<!-- end::Form -->
			</div>
		</div>
	</div>
</template>

<style scoped>
.btn-delete {
	background-color: transparent;
	border: 1px solid #f1416c !important;
	color: #f1416c !important;
}

.btn-delete:hover {
	color: #fff !important;
	background-color: #f1416c !important;
}
</style>
