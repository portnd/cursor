<script setup lang="ts">
import { useForm } from "vee-validate"
import { useSurveyRuleEditStore } from "../store"
import useHandlerSuccess from "../../../../../composables/useHandlerSuccess"
import { useInitDataStore } from "~/core/modules/initData/store"

const route = useRoute()
const id = Number(route.params.id)

const store = useSurveyRuleEditStore()
useStoreLifecycle(store)
const initDataStore = useInitDataStore()
const emit = defineEmits(["onRequestReload"])

const handleValidate = computed(() => {
	return handleFieldValidation(store.survey)
})

const { handleSubmit, isSubmitting, errors, handleReset } = useForm({
	validationSchema: handleValidate,
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.edit()
	if (res?.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				initDataStore.initData()
				navigateTo(`/settings/survey-rules`)
			},
		})
	}
})

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

// ให้ค่าที่กรอกตามเงื่อนไข
watch(
	() => store.survey,
	() => {
		// handleFieldCondition(store.survey)
	},
	{
		deep: true,
	}
)

// watch(
// 	() => store.conditionRangeId,
// 	() => {
// 		store.updateRangeId()
// 		console.log("watch:", store.conditionRangeId)
// 		emit("onRequestReload")
// 	}
// )

const onCancel = () => {
	return navigateTo("/settings/survey-rules")
}

// This function updates the left_value and right_value properties of each condition in the survey data store.
// const handleValue = () => {
// 	// Loop through each survey item in the data store.
// 	store.survey.forEach((_: any, key: any) => {
// 		// Find the survey item that matches the current condition key.
// 		const result = store.survey.find((item: any) => item.name.toLowerCase() === Object.keys(store.condition)[key])
// 		// Get the conditions for the current condition key.
// 		const conditions = Object.values(store.condition)[key]

// 		// Loop through each condition in the current survey item's conditionList.
// 		result?.conditionList.forEach((condition: any, c: any) => {
// 			// Update the left_value and right_value properties of the current condition.
// 			condition.left_value = conditions[c].left_value
// 			condition.right_value = conditions[c].right_value
// 		})
// 	})
// }

onMounted(async () => {
	await store.get(id)
	handleFieldCondition(store.survey)
})

onUnmounted(() => {
	handleReset()
})

// onBeforeRouteLeave(() => {
// 	handleReset()
// })
</script>
<template>
	<div class="row">
		<div class="col-xl-12">
			<form @submit.prevent="onSubmit">
				<div class="card p-5 pb-1 mb-5">
					<div class="row mb-8">
						<div class="col-lg-4">
							<VTextInput v-model="store.name" label="ชื่อเกณฑ์การจำแนกสภาพทาง" name="name" :required="true" />
						</div>
						<div class="col-lg-4">
							<VSelect
								v-model="store.conditionRangeId"
								:options="toOptions(initDataStore.data?.ref_condition_range)"
								label="ประเภท"
								name="type"
								:can-clear="false"
								:can-deselect="false"
								:required="true"
								@update:model-value="
									() => {
										console.log('VSelect')
										store.updateSurveyRul()
										emit('onRequestReload')
									}
								"
							/>
						</div>
					</div>
					<!-- Begin::กำหนดเกณฑ์การจำแนกสภาพทาง IRI -->
					<template v-for="(survey, key) in store.survey" :key="key">
						<div class="row justify-content-center gx-10">
							<div class="col-12 mb-5">
								<div class="card bg-primary rounded-bottom-0 py-4 px-10">
									<h5 class="fw-normal mb-0">กำหนดเกณฑ์การจำแนกสภาพทาง {{ survey.name }}</h5>
								</div>
							</div>
							<div class="row px-lg-10">
								<div class="col-md-6 mb-lg-0">
									<div class="card survey-rule mb-5">
										<div class="card-header bg-light-primary text-center ps-5">
											<h3 class="card-title">เกณฑ์การจำแนกสภาพทาง ผิวทางลาดยาง</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ survey.leftUnit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ survey.rightUnit }}</div>
											</div>
											<template v-for="(acCondition, index) in survey.ac.conditionList" :key="index">
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="acCondition.left_value"
															:name="acCondition.left_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ acCondition.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="`background: ${getGrade(acCondition.grade_id).color}`"
																	></div>
																	<div>
																		<label>{{ getGrade(acCondition.grade_id).name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ acCondition.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="acCondition.right_value"
															:name="acCondition.right_name"
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
								<div class="col-md-6 mb-lg-0">
									<div class="card survey-rule mb-5">
										<div class="card-header bg-light-primary text-center ps-5">
											<h3 class="card-title">เกณฑ์การจำแนกสภาพทาง ผิวทางคอนกรีต</h3>
										</div>
										<div class="card-body py-5">
											<div class="row text-center mb-5">
												<div class="col-3">ค่าต่ำสุด {{ survey.leftUnit }}</div>
												<div class="col-6"></div>
												<div class="col-3">ค่าสูงสุด {{ survey.rightUnit }}</div>
											</div>
											<template v-for="(ccCondition, index) in survey.cc.conditionList" :key="index">
												<div class="row align-items-center text-center mb-5">
													<div class="col-3">
														<VNumberInput
															v-model="ccCondition.left_value"
															:name="ccCondition.left_name"
															align="center"
															:precision="2"
															:max="100"
															:min="0"
															:max-length="true"
														/>
													</div>
													<div class="col-6">
														<div class="row space-between">
															<div class="col m-auto">{{ ccCondition.left_symbol }}</div>
															<div class="col-6">
																<div class="justify-content-center">
																	<div
																		class="square"
																		:style="`background: ${getGrade(ccCondition.grade_id).color}`"
																	></div>
																	<div>
																		<label>{{ getGrade(ccCondition.grade_id).name }}</label>
																	</div>
																</div>
															</div>
															<div class="col m-auto">{{ ccCondition.right_symbol }}</div>
														</div>
													</div>
													<div class="col-3">
														<VNumberInput
															v-model="ccCondition.right_value"
															:name="ccCondition.right_name"
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
					</template>
					<!-- End::กำหนดเกณฑ์การจำแนกสภาพทาง IRI -->

					<div class="d-flex justify-content-end mt-0 my-5">
						<BtnCancel @click="onCancel" />
						<BtnSubmit :loading="store.loading" label="บันทึก" />
					</div>
				</div>
			</form>
		</div>
	</div>
</template>

<style lang="scss" scoped>
.card-body {
	min-height: 335px;
}

.card-title {
	font-size: 14px !important;
}
</style>
