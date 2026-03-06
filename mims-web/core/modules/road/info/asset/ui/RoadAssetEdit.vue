<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useRoadAssetEditStore } from "../store"
import { IOption } from "~/core/shared/types/Option"
import { useInitDataStore } from "~/core/modules/initData/store"
import { IValidate } from "~/core/shared/types/Validate"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"

const route = useRoute()
const roadId = Number(route.params.roadId)

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

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()
const title = ref()
const uploadFileKm: Ref = ref()

const initDataStore = useInitDataStore()
const store = useRoadAssetEditStore()
const roadStore = useRoadTitleStore()
useStoreLifecycle([store, initDataStore])

const showModal = (item: any, template: any) => {
	handleReset()
	store.id = Number(roadId)
	store.idParentAsset = item.raw_data.id_parent_asset
	store.assetType = props.assetType
	store.template = template
	store.setRoadGeom(roadStore.data.road_info?.the_geom)
	const titles = useInitData()
		.refAssetTable()
		?.map((e) => {
			return { tableLabel: e.table_label, id: e.id }
		})
	if (titles) {
		for (let index = 0; index < titles.length; index++) {
			if (item.refAssetTableId === titles[index].id) {
				title.value = titles[index].tableLabel
			}
		}
	}
	store.pushDataToTemplate(item.raw_data)
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	store.$reset()
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineRule("kmRoadEdit", (value: any) => {
	const km = convertStringToKm(value)
	if (roadStore.data.road_info.ref_direction_id === 1) {
		if (km < roadStore.data.road_info?.km_start || km > roadStore.data.road_info?.km_end) {
			return "กม. ที่ท่านกรอกไม่ได้อยู่ในสายทาง"
		} else {
			return true
		}
	} else if (km > roadStore.data.road_info?.km_start || km < roadStore.data.road_info?.km_end) {
		return "กม. ที่ท่านกรอกไม่ได้อยู่ในสายทาง"
	} else {
		return true
	}
})

defineRule("requiredImageEdit", (_: any, index: any) => {
	const i = Number(index[0])
	if (store.template.length !== 0) {
		if (store.template[i].value === "") {
			return "โปรดระบุ"
		} else if (store.template[i].value === null) {
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
	return store.template?.length
})

const { handleSubmit, handleReset } = useForm({
	validationSchema: handleValidate,
})
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.putAssetRoad(roadId)

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
	store.resetTemplate()
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
					<h3 class="modal-title fw-semibold">แก้ไขข้อมูล {{ title }}</h3>
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
										(item.component_type === 'text-number' && item.column_name.toLowerCase().includes('latitude')) ||
										item.column_name.toLowerCase().includes('longitude')
									"
									class="col-12 mb-5"
									:class="templateLength > 4 ? 'col-md-6' : ''"
								>
									<VLongDecimalInput
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
											:searchable="true"
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
										:max-date="null"
										:year-picker="true"
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
											:files="store.imageUrlList.find((image) => image.name === item.column_name)?.path"
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
							<BtnSubmit label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
