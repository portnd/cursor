<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadSummaryEditStore, useRoadSummaryStore } from "../../store"
import { IRoad } from "../../infrastructure"
import { IValidate } from "~/core/shared/types/Validate"
import { useRoadCreateStore } from "~/core/modules/road/roadList/store"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import VSelect from "~/components/inputs/VSelect.vue"

defineProps({
	roadId: {
		type: Number,
		default: 0,
	},
})
// Modal
const { $bootstrap }: any = useNuxtApp()

const summaryStore = useRoadSummaryStore()
const store = useRoadSummaryEditStore()
const createStore = useRoadCreateStore()
const titleStore = useRoadTitleStore()
const route = useRoute()

useStoreLifecycle([store, createStore, titleStore, summaryStore])
const modal = ref()
const district = ref()
const depot = ref()
const origin = ref()
const destination = ref()
const lineFiles = ref()
const laneFiles = ref()
const kmstart = ref<string>("")
const kmend = ref<string>("")
const remark = ref<string>("")
const color = ref<string>("")
const roadType = ref<number>(0)
const lengthKmStartInput = ref(0)
const lengthKmEndInput = ref(0)

const showModal = async (data: IRoad) => {
	await setTimeout(async () => {
		await setErrors({})
		const inputStart = document.querySelectorAll(".input-kmStart")
		const inputEnd = document.querySelectorAll(".input-kmEnd")
		const separator = document.querySelector("#separator") as HTMLElement
		inputStart.forEach((elementStart) => {
			elementStart.classList.remove("input-km--error")
		})
		inputEnd.forEach((elementEnd) => {
			elementEnd.classList.remove("input-km--error")
		})
		separator.style.setProperty("margin-bottom", "1rem", "important")
		// separator.forEach((elementEnd) => {
		// 	elementEnd.classList.remove("input-separator-error")
		// })
	}, 300)
	store.road = data
	const modalElement = modal.value
	lineFiles.value?.onParentActive()
	laneFiles.value?.onParentActive()
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	district.value = store.road.responsible_code.split(" - ")[0]
	depot.value = store.road.responsible_code.split(" - ")[1]
	origin.value = store.road.origin_to_destination.split(" - ")[0]
	destination.value = store.road.origin_to_destination.split(" - ")[1]
	remark.value = store.road.road_info?.remark ?? ""
	color.value =
		store.road.road_info?.road_color_code === "undefined" ? "#DDDDDD" : store.road.road_info?.road_color_code
	roadType.value = store.road.road_info?.ref_road_type_id ?? 0
	store.params.center_line_shape_filepath = store.road.road_info.center_line_shape_file_path
	store.params.center_lane_shape_filepath = store.road.road_info.center_lane_shape_file_path
	await bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// เคลียไฟล์
	lineFiles.value?.clearFile()
	laneFiles.value?.clearFile()
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

watch(
	() => store.road,
	(newDefault) => {
		if (Object.keys(newDefault).length) {
			store.converKmRef(kmstart, kmend)
			store.km_total = Math.abs(store.road.road_info?.km_end - store.road.road_info?.km_start) / 1000 ?? 0
		}
	}
)

const validate = computed(() => {
	const validations: IValidate = {}
	// validations.refRoadType = "required"
	if (store.road.road_level === 2) {
		validations.ramp_id = "required"
		validations.name = "required"
	}
	validations.km_start = "required|km"
	validations.km_end = "required|km"
	if (store.road.is_init === false) {
		validations.center_lane_shape_file = "required"
		validations.center_line_shape_file = "required"
	}
	validations.year_construction_completed = "required"
	return validations
})

const { handleSubmit, setErrors, errors, submitCount, handleReset, resetField } = useForm({
	validationSchema: validate,
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.updateRoads(Number(route.params.roadId))
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await titleStore.getData(Number(route.params.roadId))
				await summaryStore.getRoadDetail(Number(route.params.roadId))
			},
		})
		onCancel()
		hideModal()
	} else {
		store.loading = false
	}
})

const onCancel = async () => {
	lineFiles.value.clearFile()
	laneFiles.value.clearFile()
	resetField("lineFiles", { value: undefined })
	resetField("laneFiles", { value: undefined })
	store.clearFile()
	handleReset()
	store.$reset()
	await titleStore.getData(Number(route.params.roadId))
	await summaryStore.getRoadDetail(Number(route.params.roadId))
}

const onDelete = () => {
	useDeleteItem({
		name: `${store.road.origin_to_destination}`,
		url: `/roads/${Number(route.params.roadId)}`,
		callBack: () => {
			hideModal()
			navigateTo("/roads")
		},
	})
}

const refRoadType = computed(() => {
	if (store.road.road_level === 1) {
		return toOptions(useInitData().refRoadTypeLevelFirst())
	} else {
		return toOptions(useInitData().refRoadTypeLevelSecond())
	}
})

watch(
	() => store.params.ref_road_type_id,
	() => {
		createStore.getRoadInit(
			store.road.id.toString(),
			store.road.road_level.toString(),
			store.road.road_info.ref_road_type_id.toString()
		)
	}
)

watch(
	() => [store.road.road_info?.km_end, store.road.road_info?.km_start],
	() => {
		store.calculateDistance()
	}
)

watch(
	() => store.road.road_info?.ref_road_type_id,
	(newData, oldData) => {
		if (oldData === 1 || oldData === 3) {
			if (newData === 2 || newData === 4) {
				lineFiles.value.clearFile()
				laneFiles.value.clearFile()
				resetField("lineFiles", { value: undefined })
				resetField("laneFiles", { value: undefined })
				store.params.center_lane_shape_filepath = ""
				store.params.center_line_shape_filepath = ""
			}
		} else if (oldData === 2 || oldData === 4) {
			if (newData === 1 || newData === 3) {
				lineFiles.value.clearFile()
				laneFiles.value.clearFile()
				resetField("lineFiles", { value: undefined })
				resetField("laneFiles", { value: undefined })
				store.params.center_lane_shape_filepath = ""
				store.params.center_line_shape_filepath = ""
			}
		}
	}
)

const addClassError = () => {
	const separator = document.querySelector("#separator") as HTMLElement
	if (separator) {
		if (submitCount.value > 0) {
			if (Object.keys(errors.value).includes("km_start") && Object.keys(errors.value).includes("km_end")) {
				separator.style.setProperty("margin-bottom", "2rem", "important")
			} else {
				separator.style.setProperty("margin-bottom", "1rem", "important")
			}
		}
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
		if (Object.keys(errors.value).includes("km_start") || Object.keys(errors.value).includes("km_end")) {
			addClassError()
		}
	}
)

defineExpose({
	showModal,
	hideModal,
})

onUnmounted(() => {
	store.$reset()
	createStore.$reset()
	titleStore.$reset()
	summaryStore.$reset()
})
</script>

<template>
	<div id="modal" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h1 class="modal-title fw-semibold">ปรับปรุงข้อมูล</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" @click="onCancel"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.road.road_code"
									name="road_code"
									label="หมายเลขตอนควบคุม 8 หลัก"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.road.road_section_name_th"
									name="road_section_name_th"
									label="ชื่อตอนควบคุม (ไทย)"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.road.road_section_name_en"
									name="road_section_name_en"
									label="ชื่อตอนควบคุม (อังกฤษ)"
									:disabled="true"
								/>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VTextInput v-model="store.road.province" name="province" label="จังหวัด" :disabled="true" />
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="district"
										class="w-100 w-lg-auto text-nowrap"
										label="หน่วยงานที่รับผิดชอบ"
										name="district"
										:disabled="true"
									/>
									<div class="align-self-end mb-4">-</div>
									<VTextInput v-model="depot" class="w-100 w-lg-auto" label="" name="depot" :disabled="true" />
								</div>
							</div>
							<div v-if="store.road.road_level === 1" class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="origin"
										class="w-100 w-lg-auto"
										label="จาก - ถึง"
										name="origin"
										:disabled="true"
									/>
									<div class="align-self-end mb-4">-</div>
									<VTextInput
										v-model="destination"
										class="w-100 w-lg-auto"
										label=""
										name="destination"
										:disabled="true"
									/>
								</div>
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<VSelect
									v-model="roadType"
									:options="refRoadType"
									name="ref_road_type_id"
									label="ประเภทของถนน"
									:required="true"
									:can-clear="false"
									:can-deselect="false"
									@update:model-value="(e: string) => store.updateInfo(e, 'roadType')"
								/>
							</div>
							<div v-if="store.road.road_level === 2" class="col-lg-4 col-12 mb-2">
								<VTextInput
									v-model="store.road.road_info.ramp_id"
									name="ramp_id"
									label="รหัส Ramp"
									:validate-number="true"
									:required="true"
								/>
							</div>
							<div v-if="store.road.road_level === 2" class="col-lg-4 col-12 mb-2">
								<VTextInput v-model="store.road.road_info.name" name="name" label="ชื่อ Ramp" :required="true" />
							</div>
							<div class="col-lg-4 col-12 mb-2">
								<div class="d-flex align-items-end gap-2">
									<VTextInput
										v-model="kmstart"
										label="ช่วง กม."
										class="w-100 w-lg-auto input-kmStart"
										name="km_start"
										:required="true"
										@update:model-value="(e: string) => store.updateKm(e, 'km_start')"
									/>
									<div id="separator" class="align-self-end mb-4">-</div>
									<VTextInput
										v-model="kmend"
										label=""
										class="w-100 w-lg-auto input-kmEnd"
										name="km_end"
										:required="true"
										@update:model-value="(e: string) => store.updateKm(e, 'km_end')"
									/>
								</div>
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VTextInput v-model="store.km_total" name="total" label="ระยะทาง (กม.)" :disabled="true" />
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VDatePicker
									v-if="store.road.road_info && store.road.road_info.year_construction_completed !== undefined"
									v-model="store.road.road_info.year_construction_completed"
									label="ปีที่อนุมัติลงทะเบียน"
									name="year_construction_completed"
									:year-picker="true"
									:required="true"
									:max-date="null"
								/>
							</div>
							<div class="col-lg-4 col-6 mb-2">
								<VColorPicker
									v-model="color"
									name="color"
									label="เลือกสีสายทาง"
									:required="true"
									@update:model-value="(e: string) => store.updateInfo(e, 'color')"
								/>
							</div>
						</div>
						<div class="row">
							<div class="col-lg-6 col-12 mb-2">
								<VUploadFile
									ref="lineFiles"
									v-model="store.params.center_line_shape_file"
									:files="store.params.center_line_shape_filepath"
									label="แนบไฟล์ Centerline"
									:required="store.road.is_init === true ? false : true"
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
									ref="laneFiles"
									v-model="store.params.center_lane_shape_file"
									:files="store.params.center_lane_shape_filepath"
									label="แนบไฟล์ Centerlane"
									:required="store.road.is_init === true ? false : true"
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
								<VTextarea
									v-model="remark"
									label="หมายเหตุ"
									name="remark"
									min-height="140px"
									@update:model-value="(e: string) => store.updateInfo(e, 'remark')"
								/>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<div class="row w-100">
							<div
								class="col-12 col-lg-7 text-start mb-sm-5 mb-2 mb-lg-0 justify-content-sm-start justify-content-center"
							>
								<button
									type="button"
									class="d-lg-flex d-sm-none mb-sm-0 mb-2 d-flex btn btn-outline btn-outline-danger rounded-2 fw-semibold float-start me-4 order-sm-first order-last"
									@click="onDelete"
								>
									ลบข้อมูล
								</button>
								<NuxtLink
									class="btn btn-outline btn-outline-primary rounded-2 fw-semibold float-start me-4"
									to="/files/pdf/20240110_MIMS_คู่มือการนำเข้าข้อมูลสายทาง.pdf"
									target="_blank"
									download
								>
									ดาวน์โหลดคู่มือ Centerline, Centerlane
								</NuxtLink>
							</div>
							<div
								class="col-12 col-sm-6 d-lg-none d-sm-flex d-none justify-content-sm-start justify-content-center mb-sm-0 mb-2"
							>
								<button
									type="button"
									class="btn btn-outline btn-outline-danger rounded-2 fw-semibold float-start me-4"
									@click="onDelete"
								>
									ลบข้อมูล
								</button>
							</div>
							<div class="col-12 col-sm-6 col-lg-5 text-sm-end text-center">
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

<style scoped lang="scss">
.separator {
	margin-bottom: 1rem !important;
}
.input-separator-error {
	margin-bottom: 2rem !important;
}
</style>
