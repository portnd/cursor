<script setup lang="ts">
import { useForm } from "vee-validate"
import { useReflectiveStore } from "../store"

const store = useReflectiveStore()

useStoreLifecycle(store, { resetOnEnter: false })

const validations = computed(() => {
	const validate = handleFieldReflectivityValidation(store.reflectivityInput, store.params.owner_name)
	validate.name = ""
	return validate
})

const { handleSubmit } = useForm({
	validationSchema: validations,
})

const onSubmit = handleSubmit((_, actions) => {
	useAction(actions)
	store.submitRule()
	useHandlerSuccess(0, { fn: () => {} })

	hideModal()
})

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = new $bootstrap.Modal(modalElement)

	handleFieldReflectivityCondition(store.reflectivityInput, store.params.owner_name)

	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	// @ts-ignore
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})

</script>

<template>
	<div id="modal-compare-graph" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0 bg-primary">
					<h4 class="fw-semibold mb-0" style="color: #444444">กำหนดเกณฑ์ค่าการสะท้อนแสง</h4>
					<button
						type="button"
						class="btn-close"
						data-bs-dismiss="modal"
						aria-label="Close"
						@click="store.cancelRule"
					></button>
				</div>
				<form @submit.prevent="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-lg-4 px-5">
								<VSelect
									v-model="store.params.owner_id"
									:options="store.getOwnerOptions"
									:can-clear="false"
									:can-deselect="false"
									label="เลือกเกณฑ์ค่าการสะท้อนแสง"
									name="ref_condition_range_id"
									@update:model-value="() => store.onUpdateOwner()"
								/>
							</div>
							<div class="row px-5 pt-5">
								<div class="col-lg-6 col-12 mb-lg-0">
									<div class="card survey-rule mb-5">
										<div class="card-header bg-light-primary text-center ps-5">
											<h3 class="card-title">กำหนดเกณฑ์ค่าการสะท้อนแสง - เส้นสีขาว</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ store.checkReflectivitiyInput("white")?.left_unit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ store.checkReflectivitiyInput("white")?.right_unit }}</div>
											</div>
											<template
												v-for="(white, whiteIndex) of store.checkReflectivitiyInput('white')?.reflectivity_list"
												:key="`white${whiteIndex}`"
											>
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="white.left_value"
															:name="white.left_name"
															align="center"
															:precision="2"
															:max="999"
															:min="0"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ white.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="{
																			background: store.checkCriteria(white.grade_id)?.color,
																		}"
																	></div>
																	<div>
																		<label>{{ store.checkCriteria(white.grade_id)?.name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ white.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="white.right_value"
															:name="white.right_name"
															align="center"
															:max="999"
															:min="0"
															:precision="2"
															:max-length="true"
														/>
													</div>
												</div>
											</template>
										</div>
									</div>
								</div>
								<div class="col-lg-6 col-12 mb-lg-0">
									<div class="card survey-rule mb-5">
										<div class="card-header bg-light-primary text-center ps-5">
											<h3 class="card-title">กำหนดเกณฑ์ค่าการสะท้อนแสง - เส้นสีเหลือง</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ store.checkReflectivitiyInput("yellow")?.left_unit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ store.checkReflectivitiyInput("yellow")?.right_unit }}</div>
											</div>
											<template
												v-for="(yellow, yellowIndex) of store.checkReflectivitiyInput('yellow')?.reflectivity_list"
												:key="`yellow${yellowIndex}`"
											>
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="yellow.left_value"
															:name="yellow.left_name"
															align="center"
															:max="999"
															:min="0"
															:precision="2"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ yellow.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="{ background: store.checkCriteria(yellow.grade_id)?.color }"
																	></div>
																	<div>
																		<label>{{ store.checkCriteria(yellow.grade_id)?.name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ yellow.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="yellow.right_value"
															:name="yellow.right_name"
															align="center"
															:max="999"
															:min="0"
															:precision="2"
															:max-length="true"
														/>
													</div>
												</div>
											</template>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
					<div class="modal-footer px-8">
						<button
							type="button"
							class="btn btn-outline btn-outline-primary rounded-4 px-3 py-2 me-5 fw-semibold fs-6"
							@click="store.defaultInput"
						>
							คืนค่าเริ่มต้น
						</button>
						<div>
							<BtnCancel data-bs-dismiss="modal" @click="store.cancelRule" />

							<BtnSubmit label="ยืนยัน" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-header {
	padding: 1.25rem 1.75rem;
}
.modal-body {
	padding: 0.5rem 1.5rem;
}
.btn-close {
	background: transparent url("/images/icons/svg/close-white.svg") center/1rem auto no-repeat;
	opacity: 0.75;
}
</style>
