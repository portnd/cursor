<script setup lang="ts">
import { useForm } from "vee-validate"
import { useConditionEditStore, useConditionStore } from "../store"

const route = useRoute()
const id = Number(route.params.roadId)
const store = useConditionEditStore()
const conditionStore = useConditionStore()
useStoreLifecycle(store)

const upLoadIRI = ref()
const upLoadImage = ref()

const downloadFile = () => {
	useDownloadFile(`ดาวน์โหลด .CSV TEMPLATE`, `roads/${id}/condition_template`)
}

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = async (idParent: number) => {
	// เคลียไฟล์
	upLoadIRI.value.clearFile()
	upLoadImage.value.clearFile()

	await store.getDefaultData(idParent)
	await store.getLaneList(id)
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
}

const { handleSubmit, handleReset } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.updateCondition(id)
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				conditionStore.params.id_parent = res.data?.id_parent
				await conditionStore.callBackUpdateData(id, "update")
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
					<h3 class="modal-title fw-semibold">ปรับปรุงข้อมูลสภาพทาง</h3>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
						@click="handleReset()"
					></button>
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
									:disabled="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadIRI"
									v-model="store.iriFile"
									:files="store.iri_file_path"
									label="ไฟล์ข้อมูลสำรวจ IRI/RUT/MPD/IFI"
									:required="true"
									total-file-size="20MB"
									name="iri_filename"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadImage"
									v-model="store.imageFile"
									:files="store.image_file_path"
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
							<BtnSubmit label="บันทึก" :loading="store.loading" :disabled="store.loading" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
