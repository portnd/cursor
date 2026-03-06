<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useMaintenanceHistoryGuaranteeEditStore } from "../../../store/MaintenanceHistoryGuaranteeEditStore"
import { useMaintenanceHistoryDetailsStore } from "../../../store/MaintenanceHistoryDetailsStore"
import { IValidate } from "~/core/shared/types/Validate"

const route = useRoute()
const id = Number(route.params.id)
const store = useMaintenanceHistoryGuaranteeEditStore()
const detailsStore = useMaintenanceHistoryDetailsStore()
useStoreLifecycle(store)

// Modal
const { $bootstrap }: any = useNuxtApp()
const modal = ref()
const upLoadFile = ref()

const validate = computed(() => {
	const validation: IValidate = {}
	const paramsKeys = Object.keys(store.params)
	paramsKeys.forEach((key) => {
		if (key !== "attachments" && key !== "id" && key !== "road_group_id") {
			switch (key) {
				case "intervention_criteria_id":
					validation[key] = store.isShowMethod ? "required" : ""
					break
				case "lane":
					validation[key] = "required"
					break
				case "road_group_id":
					validation[key] = "required"
					break
				case "road_id":
					validation[key] = "required"
					break
				case "km_end":
					validation[key] = `required|km|kmEndInRange|kmEndParams`
					break
				case "km_start":
					validation[key] = `required|km|kmStartInRange|kmStartParams`
					break
			}
		}
	})
	return validation
})

defineRule("kmStartInRange", (value: any) => {
	const inputValue = convertStringToKm(value)

	const roadChildGroup = store.roadList.id === Number(store.params.road_group_id) ? store.roadList.roads : []

	const roads = roadChildGroup.find((item) => item.id === Number(store.params.road_id))
	const lane = roads?.lanes.find((l) => l.lane === Number(store.params.lane))
	const kmStart = lane?.km_start || 0
	const kmEnd = lane?.km_end || 0

	let message = ""
	if (kmStart < kmEnd) {
		if (inputValue < kmStart || inputValue > kmEnd) {
			message = "กม. เริ่มต้นไม่ได้อยู่ในช่วงของสายทาง"
		} else {
			return true
		}
	} else if (inputValue > kmStart || inputValue < kmEnd) {
		message = "กม. เริ่มต้นไม่ได้อยู่ในช่วงของสายทาง"
	} else {
		return true
	}

	if (message) {
		return message
	} else {
		return true
	}
})

defineRule("kmEndInRange", () => {
	const inputValue = convertStringToKm(store.params.km_end)
	const roadChildGroup = store.roadList.id === Number(store.params.road_group_id) ? store.roadList.roads : []

	const roads = roadChildGroup.find((item) => item.id === Number(store.params.road_id))
	const lane = roads?.lanes.find((l) => l.lane === Number(store.params.lane))
	const kmStart = lane?.km_start || 0
	const kmEnd = lane?.km_end || 0

	let message = ""
	if (kmStart < kmEnd) {
		if (inputValue > kmEnd || inputValue < kmStart) {
			message = "กม. สิ้นสุดไม่ได้อยู่ในช่วงของสายทาง"
		} else {
			return true
		}
	} else if (inputValue < kmEnd || inputValue > kmStart) {
		message = "กม. สิ้นสุดไม่ได้อยู่ในช่วงของสายทาง"
	} else {
		return true
	}

	if (message) {
		return message
	} else {
		return true
	}
})

defineRule("kmStartParams", () => {
	const inputValue = convertStringToKm(store.params.km_start)

	const roadChildGroup = store.roadList.id === Number(store.params.road_group_id) ? store.roadList.roads : []

	const roads = roadChildGroup.find((item) => item.id === Number(store.params.road_id))
	const lane = roads?.lanes.find((l) => l.lane === Number(store.params.lane))
	const kmStart = lane?.km_start || 0
	const kmEnd = lane?.km_end || 0

	const kmStartParams = store.params.km_start
	const kmEndParams = store.params.km_end
	const inputKmEnd = kmStartParams ? convertStringToKm(kmEndParams!) : 0

	let message = ""
	if (kmStart < kmEnd) {
		if (inputValue > inputKmEnd) {
			message = "กม. เริ่มต้น ต้องมีค่าน้อยกว่า กม. สิ้นสุด"
		} else if (inputValue === inputKmEnd) {
			message = "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		} else {
			return true
		}
	} else if (inputValue < inputKmEnd) {
		message = "ค่าของกม.เริ่มต้น ต้องมากกว่า ค่าของกม.สิ้นสุด"
	} else if (inputValue === inputKmEnd) {
		message = "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
	} else {
		return true
	}

	if (message) {
		return message
	} else {
		return true
	}
})

defineRule("kmEndParams", (value: any) => {
	const inputValue = convertStringToKm(value)

	const roadChildGroup = store.roadList.id === Number(store.params.road_group_id) ? store.roadList.roads : []

	const roads = roadChildGroup.find((item) => item.id === Number(store.params.road_id))
	const lane = roads?.lanes.find((l) => l.lane === Number(store.params.lane))
	const kmStart = lane?.km_start || 0
	const kmEnd = lane?.km_end || 0

	const params = store.params
	const kmStartParams = params.km_start
	const inputKmStart = kmStartParams ? convertStringToKm(kmStartParams) : 0

	let message = ""
	if (kmStart < kmEnd) {
		if (inputValue < inputKmStart) {
			message = "กม. สิ้นสุด ต้องมีค่ามากกว่า กม. เริ่มต้น"
		} else if (inputValue === inputKmStart) {
			message = "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		} else {
			return true
		}
	} else if (inputValue > inputKmStart) {
		message = "กม.สิ้นสุด ต้องน้อยกว่า กม. เริ่มต้น"
	} else if (inputValue === inputKmStart) {
		message = "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
	} else {
		return true
	}

	if (message) {
		return message
	} else {
		return true
	}
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const showModal = async (itemId: number, roadId: number, budgetMethodId: number) => {
	store.$reset()
	handleReset()

	// เคลียไฟล์
	store.filePaths = []
	upLoadFile.value.clearFile()

	await store.getHistoryDetail(id)
	await store.getRoadList(roadId)
	store.params.road_group_id = `${roadId}`
	store.createKmControlOptions(roadId)
	store.createInterventionOptions()
	store.checkIsShowMethod(budgetMethodId)
	store.setDefaultParams(itemId)

	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

watch(
	() => [store.params.km_start, store.params.km_end],
	() => {
		store.calculateSumDistance()
	}
)

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.putHistoryGuaranteeEdit(id)
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				handleReset()
				await detailsStore.getMaintenanceHistoryDetail(id)
			},
		})
	}

	hideModal()
	store.$reset()
})

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})

</script>

<template>
	<div id="modal-project-delete" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">ปรับปรุงประวัติการซ่อมบำรุง</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-md-4 col-12 mb-5">
								<VSelect
									v-model="store.params.road_id"
									:options="store.kmControlOptions"
									label="สายทาง"
									placeholder="เลือก"
									name="road_id"
									:required="true"
									@update:model-value="(e:any) => store.createLaneOptions(e)"
								/>
							</div>
							<div class="col-md-4 col-12 mb-5">
								<VSelect
									v-model="store.params.lane"
									:options="store.laneOptions"
									label="ช่องจราจร"
									placeholder="เลือก"
									name="lane"
									:required="true"
								/>
							</div>
							<div class="col-md-4 col-12 mb-5">
								<VTree
									v-if="store.isShowMethod"
									v-model="store.params.intervention_criteria_id"
									label="วิธีการซ่อมบำรุง"
									:options="store.interventionCriteriaOptions"
									placeholder="เลือกวิธีการซ่อมบำรุง"
									:required="true"
									:name="`intervention_criteria_id`"
									:disable-branch-nodes="true"
									:default-expand-level="1"
								/>
							</div>
							<div class="col-md-4 col-12 mb-5">
								<VTextInput
									v-model="store.params.km_start"
									:required="true"
									placeholder="0+000"
									label="ช่วง กม. เริ่มต้น"
									name="km_start"
								/>
							</div>
							<div class="col-md-4 col-12 mb-5">
								<VTextInput
									v-model="store.params.km_end"
									:required="true"
									placeholder="0+000"
									label="ช่วง กม. สิ้นสุด"
									name="km_end"
								/>
							</div>
							<div class="col-md-4 col-12 mb-5">
								<VTextInput v-model="store.sumInput" label="ระยะทาง (กม.)" name="distance" :disabled="true" />
							</div>
						</div>
						<div class="row">
							<div class="col-lg-8 col-12">
								<VUploadFile
									ref="upLoadFile"
									v-model="store.files"
									:files="store.filePaths"
									label="รูปภาพ (รองรับ 10 ไฟล์)"
									max-file-size="10MB"
									name="attrachment"
									:multiple="true"
									aspect-ratio="0.4"
									:max-files="10"
									:accepted-file-types="['image/png', 'image/jpg', 'image/jpeg']"
								/>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit :disabled="store.loading" label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
