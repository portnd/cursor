<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadUserCostAccStore } from "../../store/RoadUserCostAccStore"
import { IValidate } from "~/core/shared/types/Validate"

const store = useRoadUserCostAccStore()

const handleValidate = computed(() => {
	const validation: IValidate = {}
	for (const key in store.dataLoss) {
		validation[key] = "required"
	}

	return validation
})

const { handleSubmit } = useForm({ validationSchema: handleValidate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.postAccidentLossValue()
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getLossAccidentData()
			},
		})
	}
})

const onCancel = async () => {
	await store.getLossAccidentData()
}

// onUnmounted(() => store.$dispose())
// onUnmounted(() => {
// 	accLossStore.$dispose()
// })

// onBeforeRouteLeave(() => {
// 	accLossStore.$dispose()
// })
</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="col-12 mb-5">
			<div class="table-responsive">
				<table class="table customize-basic-table mb-0">
					<thead>
						<tr>
							<th class="text-center fw-semibold">ประเภท</th>
							<th class="text-center fw-semibold">ราคา (บาท/ครั้ง) <span class="required"></span></th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>มูลค่าความสูญเสียจากอุบัติเหตุที่มีผู้เสียชีวิต</td>
							<td>
								<VNumberInput
									v-model="store.dataLoss.value_of_fatal_accidents"
									name="value_of_fatal_accidents"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>มูลค่าความสูญเสียจากอุบัติเหตุที่มีผู้บาดเจ็บสาหัส</td>
							<td>
								<VNumberInput
									v-model="store.dataLoss.value_of_accidents_with_serious_injuries"
									name="value_of_accidents_with_serious_injuries"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>มูลค่าความสูญเสียจากอุบัติเหตุที่มีผู้บาดเจ็บเล็กน้อย</td>
							<td>
								<VNumberInput
									v-model="store.dataLoss.value_of_accidents_with_minor_injuries"
									name="value_of_accidents_with_minor_injuries"
									:precision="2"
								/>
							</td>
						</tr>
						<tr>
							<td>มูลค่าความสูญเสียจากอุบัติเหตุที่มีทรัพย์สินเสียหาย</td>
							<td>
								<VNumberInput
									v-model="store.dataLoss.value_of_accidents_with_property_damaged"
									name="value_of_accidents_with_property_damaged"
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
