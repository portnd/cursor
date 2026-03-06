<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAssetGroupCreateStore } from "../store"
import { useInitDataStore } from "~/core/modules/initData/store"

const props = defineProps({
	dataTable: {
		type: null,
		required: true,
	},
})

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()

	handleReset()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

const store = useAssetGroupCreateStore()
useStoreLifecycle(store)
const { handleSubmit, handleReset } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.create()

	if (res?.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res!.code, {
			showAlert: true,
			fn: () => {
				// Datatable reload ข้อมูล
				const dataTable = props.dataTable
				dataTable.loadData()
				const initDataStore = useInitDataStore()
				initDataStore.initData()

				// Dismiss modal
				hideModal()
			},
		})
	}
})

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-asset-group-create" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">เพิ่มรายการ</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 mb-5">
								<VTextInput v-model="store.name" name="name" label="ชื่อกลุ่มสินทรัพย์" :required="true" />
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit label="เพิ่ม" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
