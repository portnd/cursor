<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useRoadAssetCreateStore } from "../store/RoadAssetCreateStore"
import { IRoadAssetTableTemplate } from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"
import { useInitDataStore } from "~/core/modules/initData/store"
import { IValidate } from "~/core/shared/types/Validate"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import VNumberInput from "~/components/inputs/VNumberInput.vue"

const props = defineProps({
	propsStore: {
		type: null,
		required: true,
	},
	assetType: {
		type: String,
		default: "",
	},
	onCancel: {
		type: Function,
		require: true,
	},
})

const title = ref()
const initDataStore = useInitDataStore()
const store = useRoadAssetCreateStore()
const roadStore = useRoadTitleStore()

useStoreLifecycle([store, initDataStore])

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = (refAssetTableId: number, template: IRoadAssetTableTemplate[]) => {
	handleReset()
	store.template = template
	store.assetType = props.assetType
	const titles = useInitData()
		.refAssetTable()
		?.map((e) => {
			return { tableLabel: e.table_label, id: e.id }
		})
	if (titles) {
		for (let index = 0; index < titles.length; index++) {
			if (refAssetTableId === titles[index].id) {
				title.value = titles[index].tableLabel
			}
		}
	}
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	store.resetTemplate()
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineRule("kmRoadCreate", (value: any) => {
	const km = convertStringToKm(value)
	if (roadStore.data.road_info?.ref_direction_id === 1) {
		if (km < roadStore.data?.road_info?.km_start || km > roadStore.data?.road_info?.km_end) {
			return "กม. ที่ท่านกรอกไม่ได้อยู่ในสายทาง"
		} else {
			return true
		}
	} else if (km > roadStore.data?.road_info?.km_start || km < roadStore.data?.road_info?.km_end) {
		return "กม. ที่ท่านกรอกไม่ได้อยู่ในสายทาง"
	} else {
		return true
	}
})

defineRule("requiredImageCreate", (_: any, index: any) => {
	const i = Number(index[0])
	if (store.template.length !== 0) {
		if (store.template[i].value === "") {
			return "โปรดระบุ"
		} else if (store.template[i].value.data === null) {
			return "โปรดระบุ"
		}
	}
	return true
})

// Validate Input
const handleValidate = computed(() => {
	const validations: IValidate = {}
	for (let index = 0; index < store.template.length; index++) {
		if (store.template[index].is_required === true) {
			if (store.template[index].component_type === "text-km") {
				validations[`data${index}`] = "required|km|kmRoadEdit"
				switch (true) {
					case store.template[index].column_name.includes("altitude"):
						validations[`data${index}`] = "required"
						break
					case store.template[index].column_name.includes("latitude"):
						validations[`data${index}`] = "required"
						break
					case store.template[index].column_name.includes("longitude"):
						validations[`data${index}`] = "required"
						break
				}
			} else if (store.template[index].component_type === "image") {
				validations[`data${index}`] = `requiredImageEdit:${index}`
			} else {
				validations[`data${index}`] = "required"
			}
		} else if (store.template[index].component_type === "text-km") {
			const value = store.template[index].value

			if (value !== null && value !== undefined && value !== false && value !== "") {
				validations[`data${index}`] = "km|kmRoadEdit"
			}

			switch (true) {
				case store.template[index].column_name.includes("altitude"):
					validations[`data${index}`] = ""
					break
				case store.template[index].column_name.includes("latitude"):
					validations[`data${index}`] = ""
					break
				case store.template[index].column_name.includes("longitude"):
					validations[`data${index}`] = ""
					break
			}
		}
	}
	return validations
})

const templateLength = computed(() => {
	return store.template.length
})

const { handleSubmit, handleReset } = useForm({
	validationSchema: handleValidate,
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.postAssetRoad()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				const propsStore = props.propsStore
				const respropsStore = await propsStore.getRevision()
				await propsStore.getDataTable(respropsStore ?? 1)
			},
		})
		// Dismiss modal
		hideModal()
	}
})

const generateOptionTable = (item: any) => {
	const data: IOption[] = []
	const selectedForm: any = useInitData().selectTemplateForm(item as keyof typeof initDataStore.data)
	if (Array.isArray(selectedForm)) {
		selectedForm.map((e) => data.push({ value: e.id, label: e.name, image: e.sign_image_filepath }))
	}
	return data
}

const windowWidth = ref(window.innerWidth)
const handleResize = () => {
	windowWidth.value = window.innerWidth
}

onMounted(() => {
	window.addEventListener("resize", handleResize)
})

onUnmounted(() => {
	window.removeEventListener("resize", handleResize)
})

const onDismiss = () => {
	store.resetTemplate()
	if (props.onCancel) {
		// added check for props.onCancel
		props.onCancel()
	}
}

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-asset-create" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered" :class="templateLength > 4 ? 'modal-lg' : ''">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">เพิ่มข้อมูล {{ title }}</h3>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
						@click="() => onDismiss()"
					></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<template v-for="(item, index) in store.template" :key="index">
								<div
									v-if="item.component_type === 'text-km'"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VTextInput
										v-model="item.value"
										:name="`data${index}`"
										:label="item.component_title"
										:required="item.is_required"
									/>
								</div>
								<div
									v-if="
										(item.component_type === 'text-number' && item.column_name.includes('latitude')) ||
										item.column_name.includes('longitude')
									"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VNumberInput
										v-if="item.column_name == 'price'"
										v-model="item.value"
										:required="item.is_required"
										:label="item.component_title"
										:name="`data${index}`"
										:precision="3"
									/>
									<VLongDecimalInput
										v-else
										v-model="item.value"
										:required="item.is_required"
										:label="item.component_title"
										:name="`data${index}`"
										:precision="16"
									/>
								</div>
								<div
									v-if="item.component_type === 'text-number'"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VNumberInput
										v-model="item.value"
										:required="item.is_required"
										:label="item.component_title"
										:name="`data${index}`"
										:precision="2"
									/>
								</div>

								<template v-if="item.component_type === 'select'">
									<div
										v-if="!(item.component_title === 'ป้ายจราจร' && store.hideSignImage == false)"
										class="col-12 mb-5"
										:class="templateLength > 4 ? 'col-md-6' : ''"
									>
										<VSelect
											v-model="item.value"
											:options="generateOptionTable(item.table_name_ref)"
											:label="item.component_title"
											:name="`data${index}`"
											placeholder="เลือก"
											:required="item.is_required"
											@update:model-value="store.checkSignImage(item)"
										/>
									</div>
								</template>
								<div
									v-if="item.component_type === 'datepicker'"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VDatePicker
										v-model="item.value"
										:required="item.is_required"
										:name="`data${index}`"
										:label="item.component_title"
										:teleport-center="windowWidth <= 767 ? true : false"
									/>
								</div>
								<div
									v-if="item.component_type === 'text-year'"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VDatePicker
										v-model="item.value"
										:required="item.is_required"
										:name="`data${index}`"
										:label="item.component_title"
										:teleport-center="windowWidth <= 767 ? true : false"
										:max-range="new Date().getFullYear()"
										:max-date="new Date()"
										:year-picker="true"
										:min-date="null"
									/>
								</div>

								<template v-if="item.component_type === 'image'">
									<div
										v-if="store.hideSignImage === false"
										class="col-12 mb-5"
										:class="templateLength > 4 ? 'col-lg-6' : ''"
									>
										<VUploadFile
											ref="uploadFileKm"
											v-model="item.value"
											:required="item.is_required"
											:label="item.component_title"
											total-file-size="1MB"
											:image-size="300"
											:name="`data${index}`"
											:multiple="false"
											aspect-ratio="0.3"
											:accepted-file-types="['image/png', 'image/jpg', 'image/jpeg']"
										/>
									</div>
								</template>
								<template v-if="item.component_type === 'text'">
									<div
										v-if="!(item.component_title === 'คำอธิบายภาพ' && store.hideSignImage === true)"
										class="col-12 mb-5"
										:class="templateLength > 4 ? 'col-md-6' : ''"
									>
										<VTextarea
											v-model="item.value"
											:required="item.is_required"
											:label="item.component_title"
											:name="`data${index}`"
										/>
									</div>
								</template>
							</template>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" @click="() => onDismiss()" />
							<BtnSubmit label="เพิ่ม" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
