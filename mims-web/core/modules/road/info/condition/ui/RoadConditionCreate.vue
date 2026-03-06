<script setup lang="ts">
import { useForm } from "vee-validate"
import { useConditionCreateStore, useConditionStore } from "../store"

defineEmits(["graph_width"])

const route = useRoute()
const id = Number(route.params.roadId)

const store = useConditionCreateStore()
const conditionStore = useConditionStore()

useStoreLifecycle(store)

const upLoadIRI = ref()
const upLoadImage = ref()

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)

	store.getLaneList(id)

	// เคลียไฟล์
	upLoadIRI.value.clearFile()
	upLoadImage.value.clearFile()

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
	useDownloadFile(`ดาวน์โหลด .CSV TEMPLATE`, `roads/${id}/condition_template`)
}

const { handleSubmit, handleReset } = useForm()

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.postConditions(id)
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await conditionStore.callBackUpdateData(id, "create")
			},
		})
		// Dismiss modal
		store.$reset()
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
					<h3 class="modal-title fw-semibold">เพิ่มข้อมูลสภาพทาง</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 col-lg-6 mb-5">
								<VDatePicker v-model="store.surveyedDate" name="surveyed_date" label="วันที่สำรวจ" :required="true" />
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VSelect
									v-model="store.laneNo"
									:options="store.getLaneListOptions"
									label="ช่องจราจร"
									name="lane_no"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadIRI"
									v-model="store.iriFile"
									label="ไฟล์ข้อมูลสำรวจ IRI/RUT/MPD/IFI"
									:required="true"
									total-file-size="5MB"
									name="iri_filename"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadImage"
									v-model="store.imageFile"
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
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VTextarea v-model="store.remarks" label="หมายเหตุ" name="remarks" min-height="100px" />
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
							<BtnSubmit :disabled="store.loading" :loading="store.loading" label="เพิ่ม" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
