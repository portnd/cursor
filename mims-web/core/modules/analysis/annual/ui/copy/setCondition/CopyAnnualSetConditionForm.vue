<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAnnualAnalyseCopyStore } from "../../../store"
import { CopyAnnualStep } from ".."
import { IValidate } from "~/core/shared/types/Validate"

const store = useAnnualAnalyseCopyStore()
useStoreLifecycle(store)
const route = useRoute()
const id = Number(route.params.id)

const onCancel = () => {
	store.step = 1
}

const validate = computed(() => {
	const validations: IValidate = {}
	const keys = Object.keys(store.step2Data)

	keys.forEach((key) => {
		if (store.step2Data.condition_id === 1) {
			switch (key) {
				case "condition_id":
					validations[key] = "required"
					break
				case "target":
					validations[key] = "required"
					break
				default:
					break
			}
		} else if (store.step2Data.condition_id === 2) {
			switch (key) {
				case "condition_id":
					validations[key] = "required"
					break
				case "target":
					validations[key] = "required"
					break
				case "budget":
					validations[key] = "required"
					break
				default:
					break
			}
		} else if (store.step2Data.condition_id === 3) {
			switch (key) {
				case "condition_id":
					validations[key] = "required"
					break
				case "target":
					validations[key] = "required"
					break
				case "iri":
					validations[key] = "required"
					break
				default:
					break
			}
		} else if (store.step2Data.condition_id === 4) {
			switch (key) {
				case "condition_id":
					validations[key] = "required"
					break
				case "target":
					validations[key] = "required"
					break
				case "gn":
					validations[key] = "required"
					break
				default:
					break
			}
		}
	})

	return validations
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.createAnalyse(id)
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

watch(
	() => store.step2Data.condition_id,
	() => {
		if (store.step2Data.condition_id) {
			store.createTargetOptions(store.step2Data.condition_id)
		}
	}
)
</script>

<template>
	<CopyAnnualStep />
	<div class="row">
		<div class="col-md-3 col-12 mb-2">
			<VTextInput v-model="store.step2Data.surface_type" label="ชนิดผิวทาง" :disabled="true" name="surface_type" />
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.total_km"
				:precision="2"
				label="ระยะทาง"
				:disabled="true"
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
				text-end="ม/กม."
				align="start"
				name="iri_avg"
				:precision="2"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.gn_avg"
				:precision="2"
				label="GN เฉลี่ย"
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
				@update:model-value="(e: any) => store.createTargetOptions(e)"
			/>
		</div>
		<div class="col-md-3 col-12 mb-2">
			<VSelect
				v-model="store.step2Data.target"
				:options="store.targetOptions"
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
				text-end="ล้าน"
				align="start"
				name="budget"
				:min="0"
			/>
		</div>
		<div v-else-if="store.step2Data.condition_id === 3" class="col-md-3 col-12 mb-2">
			<VNumberInput
				v-model="store.step2Data.iri"
				label="ค่า IRI"
				:required="true"
				text-end="ม/กม."
				align="start"
				name="iri"
				:min="0"
				:precision="2"
			/>
		</div>
		<div v-else-if="store.step2Data.condition_id === 4" class="col-md-3 col-12 mb-2">
			<VNumberInput v-model="store.step2Data.gn" label="ค่า GN" :required="true" name="gn" />
		</div>
		<div v-else class="col-md-3 col-12 mb-2">
			<VTextarea v-model="store.step2Data.comment" name="comment" label="ความเห็น" />
		</div>
		<div v-show="store.step2Data.condition_id !== 1" class="row">
			<div class="col-md-3 col-12 mb-2">
				<VTextarea v-model="store.step2Data.comment" name="comment" label="ความเห็น" />
			</div>
		</div>
	</div>
	<div class="row mt-5">
		<div class="col-12 text-end">
			<BtnCancel label="ย้อนกลับ" @click="onCancel" />
			<BtnSubmit :disabled="store.loading" :loading="store.loading" label="เริ่มวิเคราะห์" @click="onSubmit" />
		</div>
	</div>
</template>

<style scoped></style>
