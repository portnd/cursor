<script setup lang="ts">
import { useForm } from "vee-validate"
import { useEditStore } from "../store/DamageEditStore"
import { useRoadDamageStore } from "../store/RoadDamageStore"

const store = useEditStore()
const damageStore = useRoadDamageStore()
useStoreLifecycle(store)
const route = useRoute()
const roadId = Number(route.params.roadId)

const uploadCSV: Ref = ref()
const uploadImage: Ref = ref()

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = async (parentId: number, id: number) => {
	// เคลียไฟล์
	uploadCSV.value.clearFile()
	uploadImage.value.clearFile()

	await store.getDamageDefault(roadId, parentId)
	store.id = id
	store.parentId = parentId

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

const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.updateDamageData(roadId)
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await damageStore.getDamageList(roadId).then(() => {
					damageStore.params.parentId = res.data.id_parent
					const filterIdParent = damageStore.damageList.find((parentItem) =>
						parentItem.items.some(
							(childItem) =>
								childItem.id_parent === damageStore.params.parentId && childItem.id === damageStore.params.id
						)
					)

					if (filterIdParent) {
						damageStore.params.year = filterIdParent.year
					}
				})

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

onUnmounted(() => {
	store.$reset()
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
					<h3 class="modal-title fw-semibold">ปรับปรุงข้อมูลความเสียหาย</h3>
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
									v-model="store.data.lane_no"
									:can-clear="false"
									:can-deselect="false"
									:options="damageStore.createLaneOptions"
									label="ช่องจราจร"
									name="lane_no"
									placeholder="เลือก"
									:disabled="true"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="uploadCSV"
									v-model="store.csv_data"
									label="ไฟล์ข้อมูลความเสียหาย"
									:files="store.data.damage_filename"
									:required="true"
									total-file-size="5MB"
									name="damage_filename"
									aspect-ratio="0.225"
									:accepted-file-types="['text/csv']"
								/>
							</div>
							<div class="col-12 col-lg-6 mb-5">
								<VUploadFile
									ref="uploadImage"
									v-model="store.image"
									label="ไฟล์ภาพกล้องหลัง"
									:files="store.data.img_filepath"
									total-file-size="30MB"
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
