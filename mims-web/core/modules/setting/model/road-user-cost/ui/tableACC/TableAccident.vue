<script setup lang="ts">
import { useForm } from "vee-validate"
import { IRucAccChanceParams } from "../../infrastructure/RoadUserCostAccRequest"
import { useRoadUserCostAccStore } from "../../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useRoadUserCostAccStore()

const validate = computed(() => {
	const validation: IValidate = {}
	for (const key in store.dataAcc) {
		if (key !== "road_group_id") {
			validation[key as keyof IRucAccChanceParams] = "required"
		}
	}
	return validation
})

const { handleSubmit } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.postAccidentChanceParams()
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getAccChanceData(store.roadGroupID)
			},
		})
	}
})

const onCancel = () => {
	store.getAccChanceData(store.roadGroupID)
}

// onUnmounted(() => store.$dispose())
</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="col-12 mb-5">
			<div class="table-responsive">
				<table class="table customize-basic-table mb-0">
					<thead>
						<tr>
							<th class="text-center align-middle fw-semibold">ประเภท</th>
							<th class="text-center align-middle fw-semibold">
								อัตราการเกิดอุบัติเหตุในแต่ละส่วนบนถนน <br />
								(จำนวนอุบัติเหตุต่อล้านคัน-กิโลเมตร) <span class="required"></span>
							</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>จำนวนอุบัติเหตุที่มีผู้เสียชีวิต</td>
							<td>
								<VNumberInput
									v-model="store.dataAcc.number_of_fatal_accidents"
									name="number_of_fatal_accidents"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>จำนวนอุบัติเหตุที่มีผู้บาดเจ็บสาหัส</td>
							<td>
								<VNumberInput
									v-model="store.dataAcc.number_of_accidents_with_serious_injuries"
									name="number_of_accidents_with_serious_injuries"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>จำนวนอุบัติเหตุที่มีบาดเจ็บเล็กน้อย</td>
							<td>
								<VNumberInput
									v-model="store.dataAcc.number_of_accidents_with_minor_injuries"
									name="number_of_accidents_with_minor_injuries"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>จำนวนอุบัติเหตุที่มีทรัพย์สินเสียหาย</td>
							<td>
								<VNumberInput
									v-model="store.dataAcc.number_of_accidents_with_property_damaged"
									name="number_of_accidents_with_property_damaged"
									:precision="2"
								/>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<div class="col-12 text-end">
			<VLoading :loading="store.loading" />
			<div>
				<BtnCancel @click="onCancel()" />
				<BtnSubmit label="บันทึก" :disabled="store.loading" @click="onSubmit" />
			</div>
		</div>
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
	width: 50%;
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
