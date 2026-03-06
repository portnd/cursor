<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useGrowthRateStore } from "../../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useGrowthRateStore()
useStoreLifecycle(store)

onMounted(() => {
	store.get()
})

const handleValidate = computed(() => {
	const validations: IValidate = {}
	store.data.forEach((_, index: number) => {
		validations[`name${index}`] = `required|not_over_1`
	})
	return validations
})

defineRule("not_over_1", (value: any) => {
	const input = value
	let message = ""

	if (input > 1) {
		message = "ค่าต้องไม่เกิน 1"
	} else {
		message = ""
	}

	return message
})

const { handleSubmit, isSubmitting, errors } = useForm({ validationSchema: handleValidate })

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.post()
	if (res) {
		useHandlerSuccess(res.code, {
			showAlert: true,
		})
	}
})

const onCancel = () => {
	store.get()
}

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<form @submit="onSubmit">
			<div class="row mt-3">
				<div class="col-12">
					<div class="table-responsive">
						<table class="table customize-basic-table mb-0">
							<thead>
								<tr>
									<th class="text-center fw-semibold">ลำดับ</th>
									<th class="text-center fw-semibold">รหัสสายทาง</th>
									<th class="text-center fw-semibold">ชื่อสายทาง</th>
									<th class="text-center fw-semibold">อัตราการเติบโตของปริมาณจราจร <span class="required"></span></th>
								</tr>
							</thead>
							<tbody>
								<tr v-for="(item, index) in store.data" :key="index">
									<td class="text-center">{{ item.road_group_id }}</td>
									<td class="text-center">{{ item.code }}</td>
									<td>{{ item.road_group_name }}</td>
									<td><VNumberInput v-model="item.r" :min="0" :max="100" :name="`name${index}`" :precision="2" /></td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
			<div class="row mt-5">
				<div class="col-12 text-end">
					<VLoading :loading="store.loading" />
					<div>
						<BtnCancel @click="onCancel()" />
						<BtnSubmit label="บันทึก" :disabled="store.loading" />
					</div>
				</div>
			</div>
		</form>
	</VSkeletonLoader>
</template>

<style scoped>
thead tr th {
	padding: 18px 0px;
	border: 0;
}

thead tr th:first-of-type {
	border-top-left-radius: 7px;
}

thead tr th:last-of-type {
	border-top-right-radius: 7px;
}

tbody tr td {
	padding-left: 1.5rem !important;
	padding-right: 1.5rem !important;
	padding-bottom: 0.5rem;
	padding-top: 0.5rem;
	vertical-align: middle;
}

tbody tr td[colspan="2"] {
	background-color: var(--kt-gray-300);
}

.customize-basic-table tr td {
	border-right: 1px solid var(--kt-gray-300) !important;
}

.customize-basic-table tr td:last-of-type {
	border: 0 !important;
}

.customize-basic-table tr:last-of-type td:last-of-type,
.customize-basic-table tr:hover:last-of-type td:last-of-type {
	border-radius: 0px !important;
}
</style>
