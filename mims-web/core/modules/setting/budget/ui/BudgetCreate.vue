<script setup lang="ts">
import { useForm } from "vee-validate"
import { useBudgetCreateStore } from "../store"
import { useInitDataStore } from "~/core/modules/initData/store"
import { IValidate } from "~/core/shared/types/Validate"

const handleValidate = computed(() => {
	const validations: IValidate = {}
	validations.name = `required`
	store.budget.forEach((_, index: number) => {
		validations[`method_name${index}`] = `required`
		// validations[`cost_per_unit${index}`] = `required|min_value:0`
	})
	return validations
})

const store = useBudgetCreateStore()
useStoreLifecycle(store)
const { handleSubmit, handleReset, errors, isSubmitting } = useForm({
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
	const res = await store.create()
	if (res.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/settings/budgets")
				handleReset()
				const initDataStore = useInitDataStore()
				initDataStore.initData()
			},
		})
	}
})

const onCancel = () => {
	navigateTo("/settings/budgets")
	handleReset()
}

</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<form @submit.prevent="onSubmit">
				<div class="card p-5 pb-1 mb-5">
					<div class="row mb-3">
						<div class="col-md-6 mb-1">
							<VTextInput v-model="store.name" name="name" label="ชื่องบประมาณ" :required="true" />
						</div>
					</div>
					<div v-for="(item, index) in store.budget" :key="index" class="row">
						<div class="col-sm-6 col-12" style="min-height: 90px">
							<VTextInput
								v-model="item.methodName"
								:name="`method_name${index}`"
								label="วิธีการซ่อมบำรุง"
								:required="true"
							/>
						</div>
						<!-- <div class="col-md col-12 mb-3">
							<VNumberInput
								v-model="item.costPerUnit"
								:name="`cost_per_unit${index}`"
								label="ราคาต่อหน่วย (บาท)"
								:required="true"
								:precision="2"
							/>
						</div> -->
						<div class="col-md-auto col-12 align-self-center text-md-start text-end mt-4">
							<BtnDelete
								:style="store.budget.length === 1 ? `visibility: hidden;` : ``"
								@click="store.deleteItemBudget(item.id)"
							/>
						</div>
					</div>
					<div class="col-12 mt-1">
						<button
							type="button"
							class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold"
							@click="store.addItemBudget()"
						>
							<i class="fi fi-rr-plus align-middle fs-8"></i>
							เพิ่ม
						</button>
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

<style scoped></style>
