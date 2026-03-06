<script setup lang="ts">
import { useForm } from "vee-validate"
import { useReflectivityRuleCreateStore } from "../store"
import { useInitDataStore } from "~/core/modules/initData/store"

const store = useReflectivityRuleCreateStore()
const initDataStore = useInitDataStore()
useStoreLifecycle(store)

const handleValidate = () => {
	return handleFieldReflectivityValidation(store.rule)
}

const { handleSubmit, isSubmitting, errors, handleReset } = useForm({
	validationSchema: handleValidate(),
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.create()
	if (res?.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res!.code, {
			showAlert: true,
			fn: () => {
				initDataStore.initData()
				handleReset()
				navigateTo(`/settings/reflectivity-criterias`)
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
watch(store.rule, () => {
	console.log("watch handleFieldCondition")
	// handleFieldReflectivityCondition(store.rule)
})

const onCancel = () => {
	handleReset()
	return navigateTo("/settings/reflectivity-criterias")
}

onMounted(() => {
	handleFieldReflectivityCondition(store.rule)
})

onBeforeRouteLeave(() => {
	handleReset()
})
</script>
<template>
	<div class="row">
		<div class="col-xl-12">
			<form @submit.prevent="onSubmit">
				<div class="card p-5 pb-1 mb-5">
					<div class="row mb-8">
						<div class="col-md-3 col-12 mb-2">
							<VTextInput v-model="store.name" label="ชื่อเกณฑ์ค่าการสะท้อนแสง" name="name" :required="true" />
						</div>
						<div class="col-md-3 col-12 mb-2">
							<VSelect
								v-model="store.reflectivityRangeId"
								:options="store.optionGenerator"
								label="ประเภท"
								name="type"
								:can-clear="false"
								:can-deselect="false"
								:required="true"
							/>
						</div>
					</div>

					<div class="row justify-content-center gx-10">
						<div class="row">
							<div class="col-md-6 mb-lg-0">
								<div class="card survey-rule mb-5">
									<div class="card-header bg-primary text-center ps-5">
										<h3 class="card-title">กำหนดเกณฑ์ค่าการสะท้อนแสง - เส้นสีขาว</h3>
									</div>
									<div class="card-body py-5">
										<div class="row text-center mb-5">
											<div class="col-3 text-gray-800">ค่าต่ำสุด {{ store.rule[0].leftUnit }}</div>
											<div class="col-6"></div>
											<div class="col-3 text-gray-800">ค่าสูงสุด {{ store.rule[0].rightUnit }}</div>
										</div>
										<template v-for="(whiteItem, key) in store.rule[0].white.reflectivity_list" :key="key">
											<div class="row align-items-center text-center mb-5">
												<div class="col-3">
													<VNumberInput
														v-model="whiteItem.left_value"
														:name="whiteItem.left_name"
														align="center"
														:precision="2"
														:max="999"
														:min="0"
														:max-length="true"
													/>
												</div>
												<div class="col-6">
													<div class="row space-between">
														<div class="col m-auto">{{ whiteItem.left_symbol }}</div>
														<div class="col-6">
															<div class="justify-content-center">
																<div class="square" :style="`background: ${getGrade(whiteItem.grade_id).color}`"></div>
																<div>
																	<label class="text-gray-800">{{ getGrade(whiteItem.grade_id).name }}</label>
																</div>
															</div>
														</div>
														<div class="col m-auto">{{ whiteItem.right_symbol }}</div>
													</div>
												</div>
												<div class="col-3">
													<VNumberInput
														v-model="whiteItem.right_value"
														:name="whiteItem.right_name"
														align="center"
														:precision="2"
														:max="999"
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
									<div class="card-header bg-primary text-center ps-5">
										<h3 class="card-title text-gray-800">กำหนดเกณฑ์ค่าการสะท้อนแสง - เส้นสีเหลือง</h3>
									</div>
									<div class="card-body py-5">
										<div class="row text-center mb-5">
											<div class="col-3 text-gray-800">ค่าต่ำสุด {{ store.rule[0].leftUnit }}</div>
											<div class="col-6"></div>
											<div class="col-3 text-gray-800">ค่าสูงสุด {{ store.rule[0].rightUnit }}</div>
										</div>
										<template v-for="(yellowItem, key) in store.rule[0].yellow.reflectivity_list" :key="key">
											<div class="row align-items-center text-center mb-5">
												<div class="col-3">
													<VNumberInput
														v-model="yellowItem.left_value"
														:name="yellowItem.left_name"
														align="center"
														:precision="2"
														:max="999"
														:min="0"
														:max-length="true"
													/>
												</div>
												<div class="col-6">
													<div class="row space-between">
														<div class="col m-auto">{{ yellowItem.left_symbol }}</div>
														<div class="col-6">
															<div class="justify-content-center">
																<div class="square" :style="`background: ${getGrade(yellowItem.grade_id).color}`"></div>
																<div>
																	<label class="text-gray-800">{{ getGrade(yellowItem.grade_id).name }}</label>
																</div>
															</div>
														</div>
														<div class="col m-auto">{{ yellowItem.right_symbol }}</div>
													</div>
												</div>
												<div class="col-3">
													<VNumberInput
														v-model="yellowItem.right_value"
														:name="yellowItem.right_name"
														align="center"
														:precision="2"
														:max="999"
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
