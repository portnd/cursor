<script setup lang="ts">
import { useForm } from "vee-validate"
import EditStrategicStep from "../EditStrategicStep.vue"
import { useStrategicEditStore } from "../../../store"
import { IStrategicParamsPlan } from "../../../infrastructure"
import { IValidate } from "~/core/shared/types/Validate"
import { IOption } from "~/core/shared/types/Option"

const store = useStrategicEditStore()
useStoreLifecycle(store)
// const route = useRoute()

onBeforeMount(() => {
	handleReset()
	store.wasStep2 = true
})

const validate = computed(() => {
	const validations: IValidate = {}

	const keys = Object.keys(store.step2Data)
	const planKeys =
		store.step2Data.plans?.length === 0 ? [] : store.step2Data.plans?.flatMap((item) => Object.keys(item))
	keys.forEach((key) => {
		if (store.step2Data.condition_id === 1) {
			switch (key) {
				case "discount":
					validations[key] = "required"
					break
				case "condition":
					validations[key] = "required"
					break
				default:
					break
			}
		} else {
			switch (key) {
				case "discount":
					validations[key] = "required"
					break
				case "number_plan":
					validations[key] = "required"
					break
				case "year":
					validations[key] = "required"
					break
				case "target":
					validations[key] = "required"
					break
				case "condition":
					validations[key] = "required"
					break
				default:
					break
			}
		}
	})

	if (planKeys?.length > 0) {
		planKeys?.forEach((key) => {
			switch (key) {
				case "id":
					break
				case "plan_year":
					break

				default:
					store.step2Data.plans.forEach((_, index) => {
						if (!key.startsWith("isNew") && key.includes("plan")) {
							validations[`${key}${index}`] = "required"
						}
					})
					break
			}
		})
	}

	return validations
})

const onCancel = () => {
	store.step = 1
}

const yearOptions = computed(() => {
	const options: IOption[] = []
	for (let i = 1; i <= 10; i++) {
		options.push({ label: `${i} ปี`, value: i })
	}

	return options
})

const planOptions = computed(() => {
	const options: IOption[] = []
	for (let i = 1; i <= 3; i++) {
		options.push({ label: `${i}`, value: i })
	}

	return options
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.startAnalyze()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/analyses")
				handleReset()
				store.$reset()
			},
		})
	}
})

// const cancel = () => {
// 	const path = store.isCopy
// 		? "/analyses"
// 		: `/analyses/${store.getPath?.group}/summary/${route.params.id}/${store.getPath?.criteria}`
// 	const prevUrl = path
// 	navigateTo(prevUrl)
// }
</script>

<template>
	<div class="row">
		<EditStrategicStep />
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="store.step2Data.surface_type" label="ชนิดผิวทาง" :disabled="true" name="surface_type" />
		</div>

		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.total_km"
				label="ระยะทาง"
				:disabled="true"
				:precision="2"
				text-end="กม."
				align="start"
				name="total_km"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.iri_avg"
				label="IRI เฉลี่ย"
				:disabled="true"
				:precision="2"
				text-end="ม/กม."
				align="start"
				name="iri_avg"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.ifi_avg"
				label="IFI เฉลี่ย"
				:precision="2"
				:disabled="true"
				align="start"
				name="ifi_avg"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.condition_id"
				:options="store.getConditionOptions"
				name="condition"
				placeholder="เลือก"
				label="เงื่อนไข"
				:required="true"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="store.onUpdateConditionId"
			/>
		</div>
		<div v-show="store.step2Data.condition_id !== 1" class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.target"
				:options="store.getTargetOptions"
				label="เป้าหมาย"
				placeholder="เลือก"
				name="target"
				:required="true"
				:can-clear="false"
				:can-deselect="false"
			/>
		</div>
		<div v-show="store.step2Data.condition_id !== 1" class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.number_plan"
				:options="planOptions"
				label="จำนวนทางเลือก"
				placeholder="เลือก"
				name="number_plan"
				:required="true"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="store.createTable"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.year"
				:options="yearOptions"
				label="ระยะเวลา (ปี)"
				placeholder="เลือก"
				name="year"
				:required="true"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="store.createTable"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.discount"
				label="อัตราคิดลด (Discount Rate)"
				text-end="%"
				align="start"
				:required="true"
				name="discount"
				:precision="2"
			/>
		</div>
		<div v-if="store.step2Data.condition_id === 1" class="col-md-3 col-12 mb-2">
			<VTextarea v-model="store.step2Data.comment" name="comment" label="ความเห็น" />
		</div>
		<div v-else class="col-md-3 col-12 mb-2">
			<VTextarea v-model="store.step2Data.comment" name="comment" label="ความเห็น" />
		</div>
	</div>
	<div v-show="store.step2Data.condition_id !== 1" class="row mt-5">
		<div class="col-12">
			<div class="table-responsive">
				<table class="table customize-basic-table mb-0">
					<thead>
						<tr>
							<th class="text-center py-5">ปี</th>
							<th v-for="(plan, index) of store.step2Data.number_plan" :key="index" class="text-center pt-5">
								ทางเลือกที่ {{ plan
								}}<BtnDelete
									v-show="store.step2Data.number_plan! > 1"
									class="float-end"
									@click="store.removePlan(plan)"
								/>
							</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="(_, i) in store.step2Data.plans" :key="i">
							<td class="text-center align-middle">ปีที่ {{ i + 1 }}</td>
							<td v-for="(_, index) of store.step2Data.number_plan" :key="index" class="text-start">
								<VNumberInput
									v-model="store.step2Data.plans[i][`plan_${index + 1}` as keyof IStrategicParamsPlan]"
									:name="`plan_${index + 1}${i}`"
									:precision="2"
									:required="true"
								/>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
			<span v-if="store.step2Data.condition_id === 3" class="float-end fs-6 mt-5">หน่วย : ม./กม.</span>
			<span v-else class="float-end fs-6 mt-5">หน่วย : ล้านบาท</span>
		</div>
	</div>
	<div class="row mt-5">
		<div class="col-12 text-end">
			<BtnCancel label="ย้อนกลับ" @click="onCancel" />
			<BtnSubmit
				:disabled="store.submit_loading"
				:loading="store.submit_loading"
				label="เริ่มวิเคราะห์"
				@click="onSubmit"
			/>
		</div>
	</div>
</template>

<style scoped>
@media (max-width: 768px) {
	.table {
		width: max-content;
	}
}
</style>
