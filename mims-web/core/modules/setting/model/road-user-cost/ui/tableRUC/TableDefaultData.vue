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
					<table class="table customize-basic-table mb-0 text-truncate">
						<thead>
							<tr>
								<th style="width: 160px" class="text-center fw-semibold">ชื่อ</th>
								<th style="width: 220px" class="text-center fw-semibold">
									vehicle_name <span class="required"></span>
								</th>
								<th style="width: 130px" class="text-center fw-semibold">Fuel_UCost <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">OIL_UCost <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">TYRE_UCost <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">Veh_UCost <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">M <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">m <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">wheels <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">NUM_PASS <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">
									T_UCost (บาท/ชม.) <span class="required"></span>
								</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="(value, name, index) in rucStore.rucData" :key="`${name}-${index}`">
								<td class="text-center">
									{{ rucStore.generateName(name) }}
									<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
										<i class="fi fi-sr-interrogation fs-5"></i>
										<template #content>
											<div>{{ CarTypeModel[name] }}</div>
										</template>
									</VPopover>
								</td>
								<td><VTextInput v-model="value.vehicle_name" :name="`vehicle_name_${index}`" /></td>
								<td><VNumberInput v-model="value.fuel_u_cost" :name="`fuel_u_cost_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.oil_u_cost" :name="`oil_u_cost_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.type_u_cost" :name="`type_u_cost_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.veh_u_cost" :name="`veh_u_cost_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.m_upper" :name="`m_upper_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.m_lower" :name="`m_lower_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.wheels" :name="`wheels_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.num_pass" :name="`num_pass_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="value.t_u_cost" :name="`t_u_cost_${index}`" :precision="4" /></td>
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

.customize-basic-table {
	width: max-content;
	@media (min-width: 1750px) {
		width: 100%;
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
