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
								<th style="width: 130px" class="text-center fw-semibold">IDLE_FUEL <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">DFUEL <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">ZETAB <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">EHP <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">EDT <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">PRAT <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">KPEA <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">PACCS_A0 <span class="required"></span></th>
								<th style="width: 130px" class="text-center fw-semibold">PCTPENG <span class="required"></span></th>
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
								<td><VNumberInput v-model="item.idle_fuel" :name="`idle_fuel_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.df_fuel" :name="`df_fuel_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.ze_tab" :name="`ze_tab_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.ehp" :name="`ehp_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.edt" :name="`edt_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.prat" :name="`prat_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.kpea" :name="`kpea_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.paccs_a0" :name="`paccs_a0_${index}`" :precision="4" /></td>
								<td><VNumberInput v-model="item.pct_peng" :name="`pct_peng_${index}`" :precision="4" /></td>
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
	@media (min-width: 1400px) {
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
