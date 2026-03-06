<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useMaintenanceEditStore } from "../store"
import { ISurfacesRuleItem } from "../infrastructure"
import { ModalStandard } from "./index"

interface IValidate {
	[key: string]: string
}

// const category = [
// 	{
// 		type: { name: "ลาดยาง (AC)", value: "asphalt" },
// 		categoryList: [
// 			{ name: "OL-Overlay", value: "ol_overlay" },
// 			{ name: "SS-Slurry Seal", value: "ss_slurry_seal" },
// 			{ name: "M&OL-Mill&Overlay", value: "mol_mill_overlay" },
// 			{ name: "RCL-Recycling", value: "rcl_recycling" },
// 			{ name: "Rc-Reconstruction", value: "rc_reconstruction" },
// 		],
// 	},
// 	{
// 		type: { name: "คอนกรีต (Concrete)", value: "concrete" },
// 		categoryList: [
// 			{ name: "FDR", value: "fdr" },
// 			{ name: "BCO", value: "bco" },
// 			{ name: "M&OL", value: "mol" },
// 			{ name: "Seal", value: "seal" },
// 		],
// 	},
// ]

const linkOption = [
	{ label: "AND", value: "AND" },
	{ label: "OR", value: "OR" },
]

const operationOption = [
	{ label: "<", value: "<" },
	{ label: "<=", value: "<=" },
]

const modalStandard: Ref = ref()
const store = useMaintenanceEditStore()
useStoreLifecycle(store)
const summaryData = ref()

watch(
	store,
	() => {
		summaryData.value = store.generateSummary
	},
	{ deep: true }
)

watch([() => store.methodId, () => store.type, () => store.activeStandardIndex, () => store.data], () => {
	store.updateInterventionCriteriasSelected()
	store.updateStandartsSelected()
})

defineRule("isM&OL", (value: any) => {
	if ([3, 8].includes(store.methodId)) {
		if (isNaN(Number(value)) || value === null) {
			return "โปรดระบุ"
		} else {
			return true
		}
	}
	return true
})

const handleValidate = computed(() => {
	const validations: IValidate = {}
	validations.standard_name = "required"
	validations.surfaceType = "required"
	validations.thinkness = "required"
	validations.scraping = "isM&OL"
	validations.perUnit = "required"
	return validations
})

const { handleSubmit, validate, handleReset, resetField, isSubmitting, errors } = useForm({
	validationSchema: handleValidate,
})

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.edit()
	if (res.status === true) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.get()
				summaryData.value = store.generateSummary
			},
		})
	}
})

const handleType = async (category: ISurfacesRuleItem, type: any) => {
	if ((store.getInterventionCriterias.length ?? 0) > 0) {
		const { valid } = await validate()
		if (valid) {
			store.type = type
			store.category = category.name
			store.methodId = category.id

			store.activeStandardIndex = 0
		} else {
			showAlert({
				html: `โปรดกรอกข้อมูลให้ครบถ้วน`,
				type: "warning",
				callBack: async () => {},
			})
		}
	} else {
		store.type = type
		store.category = category.name
		store.methodId = category.id
	}
}

const addStandard = async () => {
	if (store.getInterventionCriterias.length > 0) {
		const { valid } = await validate()
		if (valid) {
			store.addStandard()
			resetField("surfaceType", { value: undefined })
			resetField("thinkness", { value: undefined })
			resetField("scraping", { value: undefined })
			resetField("perUnit", { value: undefined })
		} else {
			showAlert({
				html: `โปรดกรอกข้อมูลให้ครบถ้วน`,
				type: "warning",
				callBack: async () => {},
			})
		}
	} else {
		store.addStandard()
		handleReset()
	}
}

const deleteStandard = (index: any) => {
	if (store.getInterventionCriterias.length > 0) {
		showAlert({
			title: `ลบรายการ`,
			html: `คุณต้องการ <b class="fw-semibold text-danger">"ลบ"</b> รายการมาตรฐาน ${store.getInterventionCriterias[index].maintenance_standard_name}  ใช่หรือไม่`,
			type: "question",
			callBack: async () => {
				if (!store.getInterventionCriterias[index].is_new) {
					await store.delete(store.getInterventionCriterias[index].id)
				}
				await store.deleteStandard(index)
			},
		})
	} else {
		showAlert({
			title: `ลบรายการ`,
			html: `ไม่มีรายการมาตรฐานที่จะทำการ<b class="fw-semibold text-danger">"ลบ"</b><br/> โปรดกดเพิ่มมาตรฐาน`,
			type: "warning",
			callBack: async () => {},
		})
	}
}

const duplicateStandard = async (index: any) => {
	if (store.getInterventionCriterias.length > 0) {
		const { valid } = await validate()
		if (valid) {
			showAlert({
				title: `คัดลอกรายการ`,
				html: `คุณต้องการ <b class="fw-semibold text-danger">"คัดลอก"</b> รายการมาตรฐาน ${store.getInterventionCriterias[index].maintenance_standard_name}  ใช่หรือไม่`,
				type: "question",
				callBack: async () => {
					await store.duplicateStandard(index)
				},
			})
		}
	} else {
		showAlert({
			title: `คัดลอกรายการ`,
			html: `ไม่มีรายการมาตรฐานที่จะทำการ<b class="fw-semibold text-danger">"คัดลอก"</b><br/> โปรดกดเพิ่มมาตรฐาน`,
			type: "warning",
			callBack: async () => {},
		})
	}
}

const switchStandard = async (index: number) => {
	const { valid } = await validate()
	if (valid) {
		store.activeStandardIndex = index
	} else {
		showAlert({
			html: `โปรดกรอกข้อมูลให้ครบถ้วน`,
			type: "warning",
			callBack: async () => {},
		})
	}
}

const criteriaOptions = computed(() => {
	const criteria = useInitData().refCriteriaType()

	const options = criteria?.map((item) => ({ label: item.name, value: item.name }))
	return options ?? []
})

const handleRestore = async () => {
	const lastIndex = store.getInterventionCriterias.length - 1
	if (store.getInterventionCriterias[lastIndex].is_new ?? false) {
		resetField("surfaceType", { value: undefined })
		resetField("thinkness", { value: undefined })
		resetField("scraping", { value: undefined })
		resetField("perUnit", { value: undefined })
		resetField("description", { value: undefined })
	}
	await store.restore()
	summaryData.value = store.generateSummary
}

const getSurfaceName = (key: string) => {
	switch (key) {
		case "asphalt":
			return "ลาดยาง (AC)"

		case "concrete":
			return "คอนกรีต (Concrete)"

		default:
			return key
	}
}

onMounted(async () => {
	await store.getData()
	summaryData.value = store.generateSummary
})

</script>

<template>
	<VSkeletonLoader :loading="store.loading === true">
		<div class="row">
			<div class="col-xl-12">
				<form @submit.prevent="onSubmit">
					<div class="row">
						<div class="col-12 col-md-3 col-xl-2">
							<!-- < md -->
							<div class="row d-block d-md-none">
								<div class="col-12 align-self-end text-end mb-0">
									<button
										type="button"
										class="btn btn-primary rounded-4 px-6 py-3 fw-semibold"
										@click="modalStandard.showModal()"
									>
										จัดลำดับวิธีการซ่อมบำรุง
									</button>
								</div>
							</div>

							<div class="row">
								<div class="col-12">
									<label class="fw-semibold mt-1 mb-2 ms-1 fs-5">วิธีการซ่อมบำรุง</label>
									<table class="table customize-basic-table table-hover">
										<tbody>
											<template
												v-for="(surfaces, key) in {
													asphalt: store.methods.asphalt,
													concrete: store.methods.concrete,
												}"
												:key="key"
											>
												<tr>
													<th class="fw-semibold no-hover">{{ getSurfaceName(key) }}</th>
												</tr>
												<tr v-for="(method, index) in surfaces" :key="index">
													<td
														:id="method.name"
														:class="method.id === store.methodId ? 'active' : ''"
														class="cursor-pointer ps-6"
														@click="handleType(method, key)"
													>
														{{ method.name }}
													</td>
												</tr>
											</template>
										</tbody>
									</table>
								</div>
							</div>
							<div class="row">
								<div class="col-12">
									<div class="row">
										<div class="col-6">
											<label class="fw-semibold mt-3 mb-3 ms-1 fs-5">มาตรฐาน</label>
										</div>
										<div class="col-6 text-end align-self-center pt-1">
											<i
												class="fi-br-plus-small text-primary me-2 fs-2 cursor-pointer"
												title="เพิ่ม"
												@click="addStandard()"
											></i>
											<i
												class="fi-br-minus-small text-primary me-2 fs-2 cursor-pointer"
												title="ลบ"
												@click="deleteStandard(store.activeStandardIndex)"
											></i>
											<i
												class="fi-sr-duplicate text-primary fs-6 cursor-pointer"
												title="คัดลอก"
												@click="duplicateStandard(store.activeStandardIndex)"
											></i>
										</div>
									</div>
									<div class="standard">
										<table class="table customize-basic-table table-hover border-0 mb-0">
											<tbody>
												<tr v-for="(standard, index) in store.getInterventionCriterias" :key="index">
													<td
														:class="store.activeStandardIndex === index ? 'active' : ''"
														class="cursor-pointer"
														@click="() => switchStandard(index)"
													>
														{{ standard.maintenance_standard_name }}
													</td>
												</tr>
											</tbody>
										</table>
									</div>
								</div>
							</div>
						</div>
						<div class="col-12 col-md-9 col-xl-10">
							<!-- > md -->
							<div class="row d-none d-md-block">
								<div class="col-12 mt-md-0 mt-2 align-self-end text-end mb-6">
									<button
										type="button"
										class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6"
										@click="modalStandard.showModal()"
									>
										จัดลำดับวิธีการซ่อมบำรุง
									</button>
								</div>
							</div>

							<template v-if="store.getInterventionCriterias.length > 0">
								<div class="row">
									<div class="col-12 col-md-8">
										<div v-if="store.standardSelected" class="row">
											<div class="col-12 col-md-6 col-lg-4">
												<VTextInput
													v-model="store.standardSelected.maintenance_standard_name"
													label="ชื่อมาตรฐาน"
													name="standard_name"
													:required="true"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VTextInput
													v-model="store.category"
													label="วิธีการซ่อมบำรุง"
													name="category"
													:disabled="true"
													:required="true"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VNumberInput
													v-model="store.standardSelected.maintenance_sequence"
													label="ลำดับที่"
													name="sequence"
													:disabled="true"
													:required="true"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VSelect
													v-model="store.standardSelected.maintenance_surface_type_id"
													:options="toOptions(useInitData().refSurface())"
													label="ชนิดผิวทาง"
													:name="`surfaceType`"
													:close-on-select="true"
													:can-clear="false"
													:required="true"
													placeholder="เลือก"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VNumberInput
													v-model="store.standardSelected.maintenance_thickness"
													label="ความหนาของการซ่อม (ซม.)"
													:name="`thinkness`"
													:required="true"
													:precision="2"
												/>
											</div>
											<div v-show="[3, 8].includes(store.methodId)" class="col-12 col-md-6 col-lg-4">
												<VNumberInput
													v-model="store.standardSelected.maintenance_scraping"
													label="ความหนาของการขูด (ซม.)"
													:name="`scraping`"
													:required="true"
													:precision="2"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VNumberInput
													v-model="store.standardSelected.maintenance_cost_per_unit"
													label="ราคาต่อหน่วย (บาท/ตร.ม.)"
													:name="`perUnit`"
													:required="true"
													:precision="2"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VTextInput
													v-model="store.standardSelected.maintenance_description"
													label="คำอธิบาย"
													name="description"
												/>
											</div>
											<div class="col-12 col-md-6 col-lg-4">
												<VTextInput
													v-model="store.standardSelected.maintenance_description"
													label="คำอธิบาย"
													name="description"
												/>
											</div>
										</div>
									</div>
									<div class="col-12 col-md-4">
										<VLabel label="สรุปเงื่อนไขมาตรฐาน" />
										<div class="bg-gray-200 p-3 rounded-2 summary-data">
											<div v-html="summaryData"></div>
										</div>
									</div>
								</div>
								<div class="row mt-3">
									<div class="col-12">
										<hr class="border border-1 border-gray-300" />
										<label class="fw-semibold mt-1 mb-2 ms-1 form-text fs-5">เงื่อนไขการคำนวณ</label>
										<div class="table-responsive pb-3">
											<table class="table table-condition mb-0">
												<tbody>
													<tr>
														<td class="form-text text-center column-condition">เงื่อนไข</td>
														<td class="form-text text-center column-condition">ค่า</td>
														<td class="form-text text-center column-condition">เครื่องหมาย</td>
														<td class="form-text text-center column-condition">เกณฑ์</td>
														<td class="form-text text-center column-condition">เครื่องหมาย</td>
														<td class="form-text text-center column-condition">ค่า</td>
														<td class="form-text text-center column-condition"></td>
													</tr>
													<tr
														v-for="(condition, key) in store.standardSelected?.maintenance_condition ?? []"
														:key="key"
													>
														<td class="align-middle">
															<VSelect
																v-if="key !== 0"
																v-model="condition.condition_link"
																:options="linkOption"
																label=""
																:name="`link-${key}`"
																:close-on-select="true"
																:can-clear="false"
																placeholder="เลือก"
																:auto-height="true"
															/>
														</td>
														<td>
															<VNumberInput
																v-model="condition.condition_value_1"
																label=""
																:name="`leftValue-${key}`"
																:precision="2"
															/>
														</td>
														<td>
															<VSelect
																v-model="condition.condition_operation_1"
																:options="operationOption"
																label=""
																:name="`leftOperation-${key}`"
																:close-on-select="true"
																:can-clear="false"
																placeholder="เลือก"
																:auto-height="true"
															/>
														</td>
														<td class="col-2">
															<VSelect
																v-model="condition.condition_criterion"
																:options="criteriaOptions"
																label=""
																:name="`criteria-${key}`"
																:close-on-select="true"
																:can-clear="false"
																placeholder="เลือก"
																:auto-height="true"
															/>
														</td>
														<td>
															<VSelect
																v-model="condition.condition_operation_2"
																:options="operationOption"
																label=""
																:name="`rightOperation-${key}`"
																:close-on-select="true"
																:can-clear="false"
																placeholder="เลือก"
																:auto-height="true"
															/>
														</td>
														<td>
															<VNumberInput
																v-model="condition.condition_value_2"
																label=""
																:name="`rightValue-${key}`"
																:precision="2"
															/>
														</td>
														<td class="ps-10 align-middle">
															<a
																class="fw-semibold me-3 cursor-pointer"
																title="เลื่อนขึ้น"
																:class="key === 0 ? 'disabled' : ''"
																@click="store.switchCondition(key, key - 1, 'up')"
															>
																<i class="fi-br-angle-up ms-1"></i>
															</a>
															<a
																class="fw-semibold me-5 cursor-pointer"
																title="เลื่อนลง"
																:class="
																	key === (store.standardSelected?.maintenance_condition.length ?? 0) - 1
																		? 'disabled'
																		: ''
																"
																@click="store.switchCondition(key, key + 1, 'down')"
															>
																<i class="fi-br-angle-down ms-1"></i>
															</a>
															<a
																v-show="key !== 0"
																class="fw-semibold cursor-pointer"
																:class="
																	(store.standardSelected?.maintenance_condition.length ?? 0) - 1 === 0
																		? 'disabled'
																		: ''
																"
																@click="store.deleteCondition(key)"
															>
																<i class="fi-sr-trash ms-1 fs-5 text-danger"></i>
															</a>
														</td>
													</tr>
													<tr>
														<td colspan="6"></td>
														<td class="p-0 text-center">
															<button
																type="button"
																class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold"
																@click="store.addCondition()"
															>
																<i class="fi fi-rr-plus align-middle fs-8 d-inline-flex"></i>
																เพิ่ม
															</button>
														</td>
													</tr>
												</tbody>
											</table>
										</div>
									</div>
								</div>
							</template>
							<template v-else>
								<div class="mt-5">
									<VNotFound :is-not-shadow="true" message="ไม่พบข้อมูล โปรดเพิ่มมาตรฐาน" />
								</div>
							</template>
						</div>
					</div>
					<div v-if="store.getInterventionCriterias !== null" class="d-flex justify-content-between pt-3">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel label="คืนค่าเริ่มต้น" @click="handleRestore()" />
							<BtnSubmit :disabled="store.loading" label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</VSkeletonLoader>

	<!-- Modal -->
	<ModalStandard ref="modalStandard" />
</template>

<style scoped lang="scss">
.table-hover {
	> tbody > tr:hover > th.no-hover {
		background-color: #fff0 !important;
	}

	.no-hover {
		box-shadow: inset 0 0 0 9999px #fff0;
	}
}

table:not(.table-condition) {
	tr:first-child td {
		border-top-left-radius: 0.5rem;
		border-top-right-radius: 0.5rem;
	}
	tr:last-child td {
		border-bottom-left-radius: 0.5rem;
		border-bottom-right-radius: 0.5rem;
	}
}
td.active {
	background-color: var(--kt-gray-300);
}
div.standard {
	overflow-y: auto;
	max-height: 200px;
	.table {
		width: 100%;
	}
}

.summary-data {
	min-height: 220px;
	line-height: 1.75;
	font-size: var(--kt-input-font-size);
}

a {
	&.disabled {
		color: var(--kt-text-gray-400);
		cursor: no-drop;
	}
}

.column-condition {
	min-width: 125px;
	width: 225px;
	font-size: var(--kt-input-font-size);
}

.standard {
	border-radius: 8px;
	border: 1px solid var(--kt-gray-300);
	height: 200px;
	.customize-basic-table tr:last-of-type td:last-of-type,
	.customize-basic-table tr:hover:last-of-type td:last-of-type {
		border-radius: 0px !important;
	}

	table:not(.table-condition) tr:first-child td {
		border-top-right-radius: 0 !important;
	}
}
</style>
