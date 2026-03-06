<script setup lang="ts">
const { $bootstrap } = useNuxtApp()

defineProps({
	message: {
		type: String,
		default: "",
	},
})
// Modal
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-reason" ref="modal" class="modal fade" data-bs-keyboard="false">
		<div class="modal-dialog modal- modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">เหตุผลการส่งกลับแก้ไข</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<div class="row">
						<div class="col-12 mb-5">
							<p class="text-preline">
								{{ message === "" ? "ไม่พบข้อมูล" : decodeHTML(message) }}
							</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
