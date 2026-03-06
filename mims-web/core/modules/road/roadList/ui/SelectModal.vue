<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadCreateStore, useRoadListStore } from "../store"
import { IValidate } from "~/core/shared/types/Validate"
import { IFile } from "~/core/shared/types/File"

const store = useRoadCreateStore()
const roadListStore = useRoadListStore()
const refLineFiles = ref()
const refLaneFiles = ref()
const lineFile = ref<IFile>()
const laneFile = ref<IFile>()
const distance = ref("")
const lengthKmStartInput = ref(0)
const lengthKmEndInput = ref(0)
// const router = useRouter()

// Modal
const { $bootstrap }: any = useNuxtApp()

const modal = ref()

const showModal = async () => {
	setTimeout(() => {
		setErrors({})
		const inputStart = document.querySelectorAll(".input-kmStart")
		const inputEnd = document.querySelectorAll(".input-kmEnd")
		inputStart.forEach((elementStart) => {
			elementStart.classList.remove("input-km--error")
		})
		inputEnd.forEach((elementEnd) => {
			elementEnd.classList.remove("input-km--error")
		})
		const separator = document.querySelectorAll(".separator")
		separator.forEach((element) => {
			if (element) {
				element.classList.remove("input-separator--error")
				element.classList.add("input-separator")
			}
		})
	}, 300)
	const modalElement = modal.value
	refLineFiles.value?.onParentActive()
	refLaneFiles.value?.onParentActive()
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	handleReset()
	store.params.road_color_code = "#f57e20"
	await bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// เคลียไฟล์
	refLineFiles.value?.clearFile()
	refLaneFiles.value?.clearFile()
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineProps({
	title: {
		type: String,
		default: "",
	},
})

const refRoadType = computed(() => {
	if (store.roadInitParams.level === "1") {
		return toOptions(useInitData().refRoadTypeLevelFirst())
	} else {
		return toOptions(useInitData().refRoadTypeLevelSecond())
	}
})

watch(lineFile, () => {
	store.params.center_line_shape_file = lineFile.value?.data?.file
})

watch(laneFile, () => {
	store.params.center_lane_shape_file = laneFile.value?.data?.file
})

const validate = computed(() => {
	const validations: IValidate = {}
	validations.refRoadType = "required"
	if (store.roadInitParams.level === "2") {
		validations.ramp_id = "required"
		validations.name = "required"
	}
	validations.km_start = "required|km"
	validations.km_end = "required|km"
	validations.center_line_shape_file = "required"
	validations.center_lane_shape_file = "required"
	validations.year_construction_completed = "required"
	return validations
})

const { handleSubmit, setErrors, handleReset, errors, submitCount, resetField } = useForm({
	validationSchema: validate,
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.createRoad()
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await roadListStore.getData()
				handleReset()
			},
		})
		onCancel()
		hideModal()
		navigateTo(`/roads/${store.roadId}/summary`)
	}
})

const onCancel = () => {
	refLineFiles.value.clearFile()
	refLaneFiles.value.clearFile()
	resetField("center_lane_shape_file", { value: undefined })
	resetField("center_line_shape_file", { value: undefined })
	store.clearFile()
	handleReset()
}

const addClassError = (type: number) => {
	if (type === 1) {
		const inputStart = document.querySelectorAll(".input-kmStart")
		inputStart.forEach((element) => {
			if (element) {
				if (errors.value.km_start === "โปรดระบุ" || errors.value.km_start === "รูปแบบ กม. ไม่ถูกต้อง") {
					element.classList.add("input-km--error")
				} else {
					element.classList.remove("input-km--error")
				}
			}
		})
	} else {
		const inputEnd = document.querySelectorAll(".input-kmEnd")
		inputEnd.forEach((element) => {
			if (inputEnd) {
				if (errors.value.km_end === "โปรดระบุ" || errors.value.km_end === "รูปแบบ กม. ไม่ถูกต้อง") {
					element.classList.add("input-km--error")
				} else {
					element.classList.remove("input-km--error")
				}
			}
		})
	}
}

watch(
	() => store.params.km_start,
	(_, oldValue) => {
		if (oldValue) {
			lengthKmStartInput.value = oldValue.toString().length
		}
	}
)

watch(
	() => store.params.km_end,
	(_, oldValue) => {
		if (oldValue) {
			lengthKmEndInput.value = oldValue.toString().length
		}
	}
)

watch(
	() => errors.value,
	() => {
		if (submitCount.value === 0) {
			if (store.params.km_start !== 0) {
				addClassError(1)
			} else if (lengthKmStartInput.value > 0) {
				addClassError(1)
			}
			if (store.params.km_end !== 0) {
				addClassError(2)
			} else if (lengthKmEndInput.value > 0) {
				addClassError(2)
			}
		} else if (submitCount.value > 0) {
			if (Object.keys(errors.value).includes("km_start")) {
				addClassError(1)
			}
			if (store.params.km_start !== 0) {
				addClassError(1)
			}
			if (Object.keys(errors.value).includes("km_end")) {
				addClassError(2)
			}
			if (store.params.km_end !== 0) {
				addClassError(2)
			}
		}
	},
	{ deep: true }
)

watch(
	() => [store.params.km_end, store.params.km_start],
	() => {
		store.updateDistance(distance)
	}
)

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h1 class="modal-title fw-semibold">เพิ่มข้อมูล</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" @click="onCancel"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.roadInit.road_code"
									name="road_code"
									label="หมายเลขตอนควบคุม 8 หลัก"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.roadInit.road_section_name_th"
									name="name_th"
									label="ชื่อตอนควบคุม (ไทย)"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.roadInit.road_section_name_en"
									name="name_en"
									label="ชื่อตอนควบคุม (อังกฤษ)"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput v-model="store.roadInit.province" label="จังหวัด" name="province" :disabled="true" />
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="store.roadInit.district"
										class="w-100 w-lg-auto text-nowrap"
										label="หน่วยงานที่รับผิดชอบ"
										name="district"
										:disabled="true"
									/>
									<div class="align-self-end mb-4">-</div>
									<VTextInput
										v-model="store.roadInit.depot"
										class="w-100 w-lg-auto"
										label=""
										name="depot"
										:disabled="true"
									/>
								</div>
							</div>
							<div v-if="store.roadInitParams.level === '1'" class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="store.roadInit.origin"
										class="w-100 w-lg-auto"
										label="จาก - ถึง"
										name="origin"
										:disabled="true"
									/>
									<div class="align-self-end mb-4">-</div>
									<VTextInput
										v-model="store.roadInit.destination"
										class="w-100 w-lg-auto"
										label=""
										name="destination"
										:disabled="true"
									/>
								</div>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VSelect
									v-model="store.params.ref_road_type_id"
									name="refRoadType"
									:options="refRoadType"
									label="ประเภทของถนน"
									:required="true"
									:can-clear="false"
									:can-deselect="false"
									@update:model-value="(e: any) => store.setRoadType(e)"
								/>
							</div>

							<div v-if="store.roadInitParams.level === '2'" class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.params.ramp_id"
									name="ramp_id"
									label="รหัส Ramp"
									:required="true"
									:validate-number="true"
								/>
							</div>
							<div v-if="store.roadInitParams.level === '2'" class="col-lg-4 col-12 mb-2">
								<VTextInput v-model="store.params.name" name="name" label="ชื่อ Ramp" :required="true" />
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="store.params.km_start"
										name="km_start"
										class="w-100 w-lg-auto input-kmStart"
										label="ช่วง กม."
										:required="true"
									/>
									<div class="align-self-end separator">-</div>

									<VTextInput
										v-model="store.params.km_end"
										name="km_end"
										class="w-100 w-lg-auto input-kmEnd"
										label=""
										:required="true"
									/>
								</div>
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VTextInput v-model="distance" name="distance" label="ระยะทาง (กม.)" :disabled="true" />
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VDatePicker
									v-model="store.params.year_construction_completed"
									label="ปีที่อนุมัติลงทะเบียน"
									name="year_construction_completed"
									:year-picker="true"
									:required="true"
									:max-date="null"
								/>
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VColorPicker
									v-model="store.params.road_color_code"
									name="color"
									label="เลือกสีสายทาง"
									:required="true"
								/>
							</div>
						</div>
						<div class="row">
							<div class="col-lg-6 col-12 mb-2">
								<VUploadFile
									ref="refLineFiles"
									v-model="lineFile"
									label="แนบไฟล์ Centerline"
									:required="true"
									total-file-size="20MB"
									name="center_line_shape_file"
									aspect-ratio="0.225"
									:accepted-file-types="[
										'application/zip',
										'application/x-rar-compressed',
										'application/x-rar',
										'application/x-zip-compressed',
									]"
								/>
							</div>
							<div class="col-lg-6 col-12 mb-2">
								<VUploadFile
									ref="refLaneFiles"
									v-model="laneFile"
									label="แนบไฟล์ Centerlane"
									:required="true"
									total-file-size="20MB"
									name="center_lane_shape_file"
									aspect-ratio="0.225"
									:accepted-file-types="[
										'application/zip',
										'application/x-rar-compressed',
										'application/x-rar',
										'application/x-zip-compressed',
									]"
								/>
							</div>
							<div class="col-lg-6 col-12 mb-2">
								<VTextarea v-model="store.params.remark" label="หมายเหตุ" name="remark" min-height="140px" />
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<div class="row w-100">
							<div class="col-12 col-lg-7 text-start mb-5 mb-lg-0">
								<NuxtLink
									class="btn btn-outline btn-outline-primary rounded-2 fw-semibold float-start me-4"
									to="/files/pdf/20240110_MIMS_คู่มือการนำเข้าข้อมูลสายทาง.pdf"
									target="_blank"
									download
								>
									ดาวน์โหลดคู่มือ Centerline, Centerlane
								</NuxtLink>
							</div>
							<div class="col-12 col-lg-5 text-end">
								<BtnCancel data-bs-dismiss="modal" @click="onCancel" />
								<BtnSubmit :loading="store.loading" :disabled="store.loading" label="บันทึก" />
							</div>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped>
.input-km--error {
	margin-bottom: -15px !important;
}
.input-separator--error {
	margin-bottom: 3.5rem !important;
}
.input-separator {
	margin-bottom: 2.5rem !important;
}
</style>
