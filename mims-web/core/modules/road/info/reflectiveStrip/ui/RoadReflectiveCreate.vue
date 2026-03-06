<script setup lang="ts">
import { useForm } from "vee-validate"
import { useReflectiveCreateStore } from "../store/ReflectiveCreateStore"
import { useReflectiveStore } from "../store"
import VNumberInput from "~/components/inputs/VNumberInput.vue"

defineEmits(["graph_width"])

const route = useRoute()
const id = Number(route.params.roadId)

const createStore = useReflectiveCreateStore()
const reflectiveStore = useReflectiveStore()

const upLoadG7 = ref()
// const upLoadImage = ref()

useStoreLifecycle(createStore)

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	createStore.getLineList(id)

	// เคลียไฟล์
	upLoadG7.value.clearFile()
	// upLoadImage.value.clearFile()

	handleReset()

	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

const downloadFile = () => {
	useDownloadFile(`ดาวน์โหลด .CSV TEMPLATE`, `roads/${id}/retro_reflectivity/template`)
}

const { handleSubmit, handleReset } = useForm()

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await createStore.createReflectData(id)
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await reflectiveStore.callBackUpdateData(id, "create")
			},
		})
		// Dismiss modal
		hideModal()
	}
})

defineExpose({
	showModal,
	hideModal,
})

</script>

<template>
	<div id="modal-condition-create" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">เพิ่มข้อมูลแถบสะท้อนแสง</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit.prevent="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 col-lg-6 mb-5">
								<VDatePicker
									v-model="createStore.surveyedDate"
									name="surveyed_date"
									label="วันที่สำรวจ"
									:required="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VNumberInput
									v-model="createStore.lineNo"
									:options="createStore.getLineOptions"
									label="เส้นจราจร"
									align="start"
									name="line_no"
									:min="0"
									:required="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadG7"
									v-model="createStore.csvFile"
									label="ไฟล์ข้อมูลสำรวจ G7"
									:required="true"
									total-file-size="5MB"
									name="csv_file"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
								/>
							</div>
							<!-- <div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadImage"
									v-model="createStore.imageFile"
									label="ไฟล์ภาพกล้องหน้า"
									:required="false"
									total-file-size="1024MB"
									name="image_filename"
									aspect-ratio="0.225"
									:accepted-file-types="[
										'application/zip',
										'application/x-rar-compressed',
										'application/x-rar',
										'application/x-zip-compressed',
									]"
								/>
							</div> -->
							<div class="col-12 col-lg-6 mb-5">
								<VTextarea v-model="createStore.remarks" label="หมายเหตุ" name="remarks" min-height="100px" />
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<button
							type="button"
							class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 mb-3 fw-semibold fs-7 float-start"
							@click="downloadFile"
						>
							ดาวน์โหลด .CSV TEMPLATE
						</button>

						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit :disabled="createStore.loading" :loading="createStore.loading" label="เพิ่ม" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
