<script setup lang="ts">
import { useForm } from "vee-validate"
import { useSurfaceEditStore } from "../store"

const props = defineProps({
	dataTable: {
		type: null,
		required: true,
	},
})

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = (id: number) => {
	store.get(id)
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()

	// Reset form validation
	handleReset()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

const store = useSurfaceEditStore()
useStoreLifecycle(store)
const { handleSubmit, handleReset } = useForm()
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
			},
		})

		// Dismiss modal
		hideModal()
	}
})

const generateOptionTable = () => {
	let data: { id: number; name: string }[] = []
	if (data.length > 0) {
		data = []
	}
	useInitData()
		.refSurfaceType()
		?.map((e) => data.push({ id: e.id, name: e.name }))
	return data
}

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-surface-create" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-lg modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">แก้ไขข้อมูล</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 col-md-4 col-lg-4 mb-5">
								<VSelect
									v-model="store.data.surfaceType"
									:options="toOptions(generateOptionTable())"
									name="type"
									label="รูปแบบวัสดุ"
									:required="true"
									placeholder="เลือก"
									:close-on-select="true"
									@update:model-value="($e:any) => store.checkValidate($e)"
								/>
							</div>
						</div>
						<div class="row">
							<div class="col-12 col-md-4 col-lg-4 mb-5">
								<VTextInput v-model="store.data.surfaceName" name="name" label="ชนิดวัสดุผิวทาง" :required="true" />
							</div>
							<div class="col-6 col-md-4 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.surfaceDrainage"
									:required="true"
									label="Drainage Coefficient"
									name="drainage"
									:precision="2"
									:min="0"
									:max="100"
								/>
							</div>
							<div class="col-6 col-md-4 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.surfaceLayer"
									:required="true"
									label="Layer Coefficient"
									name="layer_coefficient"
									:precision="2"
									:min="0"
									:max="100"
								/>
							</div>
							<div class="col-12 col-md-12 col-lg-12">
								<label class="fw-semibold">Grip Number Coefficient</label>
							</div>
							<div class="col-6 col-md-3 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.a"
									:required="true"
									label="a"
									name="a"
									:precision="6"
									:min="0"
									:max="100"
								/>
							</div>
							<div class="col-6 col-md-3 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.b"
									:required="true"
									label="b"
									name="b"
									:precision="6"
									:min="0"
									:max="100"
								/>
							</div>
							<div class="col-6 col-md-3 col-lg-2 mb-5">
								<VNumberInput v-model="store.data.c" :required="true" label="c" name="c1" :min="1" :max="10" />
							</div>
							<div class="col-6 col-md-3 col-lg-2 mb-5">
								<VNumberInput
									v-model="store.data.subC"
									label="ยกกำลัง"
									:required="true"
									name="c2"
									:allow-minus="true"
									:max="100"
								/>
							</div>
							<div class="col-6 col-md-4 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.crt"
									:required="store.validates"
									label="CRT"
									name="crt"
									:precision="2"
									:min="0"
									:max="100"
								/>
							</div>
							<div class="col-6 col-md-4 col-lg-4 mb-5">
								<VNumberInput
									v-model="store.data.rrf"
									:required="store.validates"
									label="RRF"
									name="rrf"
									:precision="2"
									:min="0"
									:max="100"
								/>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit :disabled="store.loading" label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
