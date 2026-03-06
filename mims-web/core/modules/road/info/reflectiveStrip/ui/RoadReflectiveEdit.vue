<script setup lang="ts">
import { useForm } from "vee-validate"
import { useReflectiveStore, useReflectiveEditStore } from "../store"

const route = useRoute()
const roadId = Number(route.params.roadId)
const editStore = useReflectiveEditStore()
const reflectStore = useReflectiveStore()

const upLoadCSV = ref()

useStoreLifecycle(editStore)

const downloadFile = () => {
	useDownloadFile(`ดาวน์โหลด .CSV TEMPLATE`, `roads/${roadId}/retro_reflectivity/template`)
}

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = async (idParent: number) => {
	// เคลียไฟล์
	upLoadCSV.value.clearFile()
	// upLoadImage.value.clearFile()

	editStore.getLineList(roadId)
	await editStore.getDefaultData(idParent)
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

	const res = await editStore.updateReflectData(roadId)
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				reflectStore.params.id_parent = res.data?.id_parent
				await reflectStore.callBackUpdateData(roadId, "update")
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
					<h3 class="modal-title fw-semibold">ปรับปรุงข้อมูลแถบการสะท้อนแสง</h3>
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
								<VDatePicker
									v-model="editStore.surveyedDate"
									name="surveyed_date"
									label="วันที่สำรวจ"
									:required="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VNumberInput
									v-model="editStore.lineNo"
									:options="editStore.getLineOptions"
									label="เส้นจราจร"
									name="line_no"
									align="start"
									:min="0"
									:required="true"
									:disabled="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="upLoadCSV"
									v-model="editStore.csvFile"
									:files="editStore.csv_file_path"
									label="ไฟล์ข้อมูลสำรวจ G7"
									:required="true"
									total-file-size="20MB"
									name="csv_file"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VTextarea v-model="editStore.remarks" label="หมายเหตุ" name="remarks" min-height="100px" />
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
							<BtnSubmit label="บันทึก" :loading="editStore.loading" :disabled="editStore.loading" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
