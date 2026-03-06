<script setup lang="ts">
import { useForm } from "vee-validate"
import { useConditionStore } from "../store"

const store = useConditionStore()

const validations = computed(() => {
	const validate = handleFieldValidation(store.conditionInput, store.params.condition_type)
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
	handleFieldCondition(store.conditionInput, store.params.condition_type)

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

onUnmounted(() => {
	// Do NOT dispose the shared ConditionStore here.
	// RoadCondition.vue (parent) handles store.$reset() on unmount.
})
</script>

<template>
	<div id="modal-compare-graph" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-xl modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0 bg-primary">
					<h4 class="fw-semibold mb-0" style="color: #444444">
						กำหนดเกณฑ์การจำแนกสภาพทาง {{ store.params.condition_type }}
					</h4>
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
									:options="store.getSurveyRangeOptions"
									:can-clear="false"
									:can-deselect="false"
									label="เลือกเกณฑ์การจำแนกสภาพทาง"
									name="ref_condition_range_id"
									@update:model-value="() => store.onUpdateOwner()"
								/>
							</div>
							<div class="row px-5 pt-5">
								<div class="col-lg-6 col-12 mb-lg-0">
									<div class="card survey-rule mb-5">
										<div class="card-header bg-light-primary text-center ps-5">
											<h3 class="card-title">เกณฑ์การจำแนกสภาพทาง ผิวทางลาดยาง</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ store.checkConditionInput("ac")?.left_unit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ store.checkConditionInput("ac")?.right_unit }}</div>
											</div>
											<template
												v-for="(ac, acIndex) of store.checkConditionInput('ac')?.conditionList"
												:key="`${store.params.condition_type}${acIndex}`"
											>
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="ac.left_value"
															:name="ac.left_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ ac.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="{
																			background: store.checkCriteria(ac.grade_id)?.color,
																		}"
																	></div>
																	<div>
																		<label>{{ store.checkCriteria(ac.grade_id)?.name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ ac.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="ac.right_value"
															:name="ac.right_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
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
											<h3 class="card-title">เกณฑ์การจำแนกสภาพทาง ผิวทางคอนกรีต</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ store.checkConditionInput("cc")?.left_unit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ store.checkConditionInput("cc")?.right_unit }}</div>
											</div>
											<template
												v-for="(cc, ccIndex) of store.checkConditionInput('cc')?.conditionList"
												:key="`${store.params.condition_type}${ccIndex}`"
											>
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="cc.left_value"
															:name="cc.left_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ cc.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="{ background: store.checkCriteria(cc.grade_id)?.color }"
																	></div>
																	<div>
																		<label>{{ store.checkCriteria(cc.grade_id)?.name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ cc.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="cc.right_value"
															:name="cc.right_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
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

							<BtnSubmit label="ยืนยัน" @click="onSubmit" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
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
.survey-rule {
	overflow: hidden;
}
.card.survey-rule {
	border: 1px solid var(--kt-gray-300);
	box-shadow: none !important;
	.card-header {
		min-height: 48px;
		border-bottom: 1px solid var(--kt-gray-300);
		.card-title {
			color: var(--kt-gray-800);
			font-size: 15px;
			font-weight: 400;
		}
		.card-header {
			font-size: 1rem;
		}
	}
}
</style>
