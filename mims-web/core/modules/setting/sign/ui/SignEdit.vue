<script setup lang="ts">
import { useForm } from "vee-validate"
import { useSignEditStore } from "../store"
import { useInitDataStore } from "~/core/modules/initData/store"

const props = defineProps({
	dataTable: {
		type: null,
		required: true,
	},
})

const emit = defineEmits(["onFinish"])

const store = useSignEditStore()
useStoreLifecycle(store)

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()
const imageRef = ref()

const showModal = async (id: number) => {
	// เคลียไฟล์
	imageRef.value.clearFile()

	await store.get(id)
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	emit("onFinish")
	bootstrapModal?.hide()
}

const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.edit()
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				// Datatable reload ข้อมูล
				const dataTable = props.dataTable
				dataTable.loadData()
				const initDataStore = useInitDataStore()
				initDataStore.initData()
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
	<div id="modal-sign-create" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">แก้ไขข้อมูล</h3>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
						@click="
							() => {
								emit('onFinish')
							}
						"
					></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 mb-5">
								<VTextInput v-model="store.data.name" name="name" label="ชื่อเครื่องหมาย" :required="true" />
							</div>
						</div>
						<div class="row">
							<div class="col-12 mb-5">
								<VTextInput v-model="store.data.abbr" name="abbr" label="ชื่อย่อ" :required="true" />
							</div>
						</div>
						<div class="row">
							<div class="col-12 mb-5">
								<VUploadFile
									ref="imageRef"
									v-model="store.image"
									:files="store.data.sign_image_filepath"
									total-file-size="1MB"
									:image-size="300"
									name="image"
									:multiple="false"
									label="รูปภาพ"
									:required="true"
									aspect-ratio="0.225"
									:accepted-file-types="['image/png', 'image/jpg', 'image/jpeg']"
								/>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
