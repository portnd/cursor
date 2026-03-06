<script setup lang="ts">
import { useForm } from "vee-validate"
import { useHRISEditStore } from "../store"

const emit = defineEmits(["onFinish"])

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = (id: number) => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()

	// Reset form validation
	store.get(id)
	// handleReset()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
	emit("onFinish")
}

const store = useHRISEditStore()
useStoreLifecycle(store)
const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.edit()

	if (res?.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res!.code, {
			showAlert: true,
			fn: () => {
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
	<div ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">แก้ไขรายการ</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-6 mb-5">
								<VTextInput
									v-model="store.data.road_number"
									name="road_number"
									label="รหัสสายทาง"
									:validate-number="true"
								/>
							</div>

							<div class="col-6 mb-5">
								<VTextInput
									v-model="store.data.section_road_number"
									label="รหัสตอนควบคุม"
									name="section_road_number"
									:validate-number="true"
								/>
							</div>
						</div>
						<div class="row">
							<div class="col-6 mb-5">
								<VTextInput
									v-model="store.data.office_of_highways_code"
									label="รหัสสำนักงาน"
									name="office_of_highways_code"
									:validate-number="true"
								/>
							</div>
							<div class="col-6 mb-5">
								<VSelect v-model="store.status" :options="store.getStatusOption" label="สถานะ" name="status" />
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<!-- <BtnCancel data-bs-dismiss="modal" /> -->
							<button
								type="button"
								class="btn rounded-4 ms-5 mt-md-0 mt-sm-2 mt-2 fw-semibold text-gray-700"
								data-bs-dismiss="modal"
								@click="() => {}"
							>
								ยกเลิก
							</button>
							<BtnSubmit label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
