<script setup lang="ts">
import { useForm } from "vee-validate"
import { useCreateStore } from "../store/DamageCreateStore"
import { useRoadDamageStore } from "../store/RoadDamageStore"

const store = useCreateStore()
const damageStore = useRoadDamageStore()
const route = useRoute()
const roadId = Number(route.params.roadId)

useStoreLifecycle(store)

const uploadCSV: Ref = ref()
const uploadImage: Ref = ref()

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)

	// เคลียไฟล์
	uploadCSV.value.clearFile()
	uploadImage.value.clearFile()

	handleReset()

	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	store.$reset()
	bootstrapModal?.hide()
}

const { handleSubmit, handleReset } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.createDamageData(roadId)
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await damageStore.getDamageList(roadId).then(() => damageStore.setDefaultparams())
				await damageStore.getRoadDamageDetail(roadId)
			},
		})
		// Dismiss modal

		hideModal()
	}
})

const downloadFile = () => {
	useDownloadFile(`ดาวน์โหลด .CSV TEMPLATE`, `roads/${roadId}/damage_template`)
}

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
					<h3 class="modal-title fw-semibold">เพิ่มข้อมูลความเสียหาย</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 col-lg-6 mb-5">
								<VDatePicker v-model="store.date" name="surveyed_date" label="วันที่สำรวจ" :required="true" />
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VSelect
									v-model="store.lane"
									:options="damageStore.createLaneOptions"
									label="ช่องจราจร"
									name="lane_no"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="uploadCSV"
									v-model="store.itemData"
									label="ไฟล์ข้อมูลความเสียหาย"
									:required="true"
									total-file-size="5MB"
									name="damage_filename"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
									@update:model-value="(e: any) => (store.itemData ? (store.itemData = e) : undefined)"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="uploadImage"
									v-model="store.image"
									label="ไฟล์ภาพกล้องหลัง"
									total-file-size="50MB"
									name="image_filename"
									aspect-ratio="0.225"
									:accepted-file-types="[
										'application/zip',
										'application/x-rar-compressed',
										'application/x-rar',
										'application/x-zip-compressed',
									]"
									@update:model-value="(e:any) => store.image ? store.image = e : undefined"
								/>
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
							<div>
								<BtnCancel data-bs-dismiss="modal" />
								<BtnSubmit label="เพิ่ม" :loading="store.loading" :disabled="store.loading" />
							</div>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
