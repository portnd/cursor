<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadUserCostRucStore } from "../../store/RoadUserCostRucStore"
import { IRucParentParams } from "../../infrastructure/RoadUserCostRucRequest"
import { IValidate } from "~/core/shared/types/Validate"

const CarTypeModel = useCarTypeModel()
const rucStore = useRoadUserCostRucStore()

const validate = computed(() => {
	const validation: IValidate = {}

	Object.keys(rucStore.rucData).forEach((parentKey, index) => {
		for (const childKey in rucStore.rucData[parentKey as keyof IRucParentParams]) {
			validation[`${childKey}_${index}`] = "required"
		}
	})

	return validation
})

const { handleSubmit } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await rucStore.postRucParams()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				rucStore.getRucData(rucStore.rucListId)
			},
		})
	}
})

const onCancel = () => {
	rucStore.getRucData(rucStore.rucListId)
}
</script>

<template>
	<VSkeletonLoader :loading="rucStore.loading">
		<form @submit="onSubmit">
			<div class="col-12 mb-5">
				<div class="table-responsive">
					<table class="table customize-basic-table mb-0">
						<thead>
							<tr>
								<th style="width: 160px" class="text-center fw-semibold">ชื่อ</th>
								<th style="width: 140px" class="text-center fw-semibold">Kpc <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">AKMO <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">LIFE0 <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">Kp <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">a0 <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">a1 <span class="required"></span></th>
								<th style="width: 140px" class="text-center fw-semibold">CPCON <span class="required"></span></th>
							</tr>
						</thead>
						<tbody v-for="(item, name, index) in rucStore.rucData" :key="index">
							<tr>
								<td class="text-center">
									{{ rucStore.generateName(name) }}
									<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
										<i class="fi fi-sr-interrogation fs-5"></i>
										<template #content>
											<div>{{ CarTypeModel[name] }}</div>
										</template>
									</VPopover>
								</td>
								<td><VNumberInput v-model="item.kpc" :name="`kpc_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.akmo" :name="`akmo_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.life0" :name="`life0_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.kp" :name="`kp_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.a0" :name="`a0_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.a1" :name="`a1_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.cp_con" :name="`cp_con_${index}`" :precision="4" /></td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			<div class="col-12 text-end">
				<VLoading :loading="rucStore.loading" />
				<div>
					<BtnCancel @click="onCancel()" />
					<BtnSubmit label="บันทึก" :disabled="rucStore.loading" />
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
@media only screen and (max-width: 768px) {
	.customize-basic-table {
		width: max-content;
	}
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
