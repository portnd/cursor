<script setup lang="ts">
import { useForm } from "vee-validate"
import AnnualStep from "../AnnualStep.vue"
import { useAnnualCreateStore } from "../../../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useAnnualCreateStore()

onMounted(() => {
	store.wasStep2 = true
})

const onCancel = () => {
	store.step = 1
}

const validate = computed(() => {
	const validations: IValidate = {}
	const keys = Object.keys(store.step2Data)

	const conditionKeys = {
		1: ["condition_id", "target"],
		2: ["condition_id", "target", "budget"],
		3: ["condition_id", "target", "iri"],
	}

	keys.forEach((key) => {
		const conditionId = store.step2Data.condition_id

		if (conditionId in conditionKeys && conditionKeys[conditionId as keyof typeof conditionKeys].includes(key)) {
			validations[key] = "required"
		}
	})

	return validations
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.analyseData()
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/analyses")
				handleReset()
			},
		})
	}
})
</script>

<template>
	<AnnualStep />
	<div class="row">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="store.step2Data.surface_type" label="ชนิดผิวทาง" :disabled="true" name="surface_type" />
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.total_km"
				label="ระยะทาง"
				:disabled="true"
				text-end="กม."
				align="start"
				name="total_km"
				:precision="2"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.iri_avg"
				label="IRI เฉลี่ย"
				:disabled="true"
				text-end="ม/กม."
				align="start"
				name="iri_avg"
				:precision="2"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.ifi_avg"
				label="IFI เฉลี่ย"
				:precision="2"
				:disabled="true"
				align="start"
				name="gn_avg"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.condition_id"
				:options="store.getConditionOptions"
				name="condition_id"
				placeholder="เลือก"
				label="เงื่อนไข"
				:required="true"
				:can-clear="false"
				:can-deselect="false"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.target"
				:options="store.getTargetOptions"
				label="เป้าหมาย"
				placeholder="เลือก"
				name="target"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div v-if="store.step2Data.condition_id === 2" class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.budget"
				label="งบประมาณ"
				:required="true"
				text-end="ล้านบาท"
				:precision="2"
				align="start"
				name="budget"
				:min="0"
			/>
		</div>
		<div v-else-if="store.step2Data.condition_id === 3" class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.iri"
				label="ค่า IRI"
				:precision="2"
				:required="true"
				text-end="ม/กม."
				align="start"
				:min="0"
				name="iri"
			/>
		</div>
		<!-- <div class="row"> -->
		<div class="col-md-3 col-12 mb-2">
			<VTextarea v-model="store.step2Data.comment" name="comment" label="ความเห็น" />
		</div>
		<!-- </div> -->
	</div>
	<div class="row mt-5">
		<div class="col-12 d-flex justify-content-between text-end">
			<div>
				<BtnCancel label="ย้อนกลับ" @click="onCancel" />
			</div>
			<div>
				<button type="button" class="btn btn-outline-primary rounded-4 me-5 fw-semibold text-gray-800">ยกเลิก</button>
				<BtnSubmit :disabled="store.loading" :loading="store.loading" label="เริ่มวิเคราะห์" @click="onSubmit" />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
