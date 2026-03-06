<script setup lang="ts">
import { useForm } from "vee-validate"
import {
	useTrafficParameterRoadGroupStore,
	useTrafficParameterStore,
	useTrafficParameterCreateStore,
} from "../../store"
import { IValidate } from "~/core/shared/types/Validate"

const aadtRoadGroup = useTrafficParameterRoadGroupStore()
const aadtParameterStore = useTrafficParameterStore()
const aadtCreateStore = useTrafficParameterCreateStore()
useStoreLifecycle([aadtRoadGroup, aadtParameterStore, aadtCreateStore])

const handleValidate = computed(() => {
	const validation: IValidate = {}
	const params = aadtCreateStore.params
	for (const key in params) {
		if (key !== "is_truck_factor") {
			validation[key] = "required"
		}
	}

	return validation
})

const { handleSubmit } = useForm({ validationSchema: handleValidate })

watch(
	() => aadtParameterStore.roadGroupID,
	async (newRoadID) => {
		if (newRoadID) {
			await aadtParameterStore.getAadtData()
			await aadtRoadGroup.getAadtVehicleType(newRoadID)
			aadtCreateStore.params = aadtParameterStore.data
			aadtCreateStore.updateParams(aadtRoadGroup.getCalculate)
		} else {
			aadtCreateStore.$reset()
			aadtParameterStore.$reset()
		}
	}
)

watch(
	() =>
		[
			aadtCreateStore.params.is_truck_factor,
			aadtCreateStore.params.six_wheel_axle_number_id,
			aadtCreateStore.params.ten_wheel_axle_number_id,
		] as const,
	([newIsTruckFactor, newSixAxleId, newTenAxleId]) => {
		// Update the isTruckFactor
		aadtRoadGroup.isTruckFactor = newIsTruckFactor

		// Update the sixAxleId
		aadtRoadGroup.sixAxleId = newSixAxleId

		// Update the tenAxleId
		aadtRoadGroup.tenAxleId = newTenAxleId

		// If isTruckFactor is false, execute calculateTruckFactor function
		if (!newIsTruckFactor) {
			aadtCreateStore.calculateTruckFactor(aadtRoadGroup.getLoadEquivalent)
		}
	}
)

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await aadtCreateStore.postAadtParameter()
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await aadtRoadGroup.getAadtRoadGroupData()
				await aadtParameterStore.getAadtData()
				await aadtRoadGroup.getAadtVehicleType(aadtParameterStore.roadGroupID)
			},
		})
	}
})

const onCancel = async () => {
	await aadtParameterStore.getAadtData()
	await aadtRoadGroup.getAadtVehicleType(aadtParameterStore.roadGroupID)
	aadtCreateStore.params = aadtParameterStore.data
	aadtCreateStore.updateParams(aadtRoadGroup.getCalculate)
}

</script>

<template>
	<div class="row">
		<div class="col-12 mb-5">
			<div class="row">
				<div class="col-12 col-md-6 mb-10">
					<VSelect
						v-model="aadtParameterStore.roadGroupID"
						name="road_group_id"
						:options="aadtRoadGroup.getRoadGroupOptions"
						label="สายทาง"
						placeholder="เลือกสายทาง"
						:can-clear="false"
						:required="true"
						:searchable="true"
						:can-deselect="false"
					/>
				</div>
			</div>

			<VSkeletonLoader :loading="aadtRoadGroup.loading">
				<div class="row">
					<div class="col-md-6 col-12 mb-5">
						<div class="border border-1 border-gray-300 rounded px-8 pt-5 pb-3">
							<div class="row">
								<div class="col-12 text-center">
									<h2 class="text-start fw-semibold fs-5 mb-5">YAX</h2>
									<img src="/images/settings/models/yax.png" width="210" alt="" class="yax" />
								</div>
								<div class="col-12 mb-5 elane">
									<VNumberInput
										v-model="aadtCreateStore.params.elane"
										label="Elane"
										name="elane"
										:required="true"
										:precision="2"
									/>
								</div>
							</div>
						</div>
					</div>

					<div class="col-md-6 col-12">
						<div class="border border-1 border-gray-300 rounded px-8 pt-5 pb-3">
							<h2 class="text-start fw-semibold fs-5 mb-5">Speed</h2>
							<div class="row">
								<div class="col-md-6 col-12 mb-5">
									<VNumberInput
										v-model="aadtCreateStore.params.speed_average"
										label="Speed Average"
										name="speed_average"
										:required="true"
										:precision="2"
									/>
								</div>
								<div class="col-md-6 col-12 mb-5">
									<VNumberInput
										v-model="aadtCreateStore.params.speed_heavy_truck"
										label="Speed Heavy Truck"
										name="speed_heavy_truck"
										:required="true"
										:precision="2"
									/>
								</div>
								<div class="col-md-6 col-12 mb-5">
									<VNumberInput
										v-model="aadtCreateStore.params.lane_distribution_factor"
										label="Lane distribution factor"
										name="lane_distribution_factor"
										:required="true"
										:precision="2"
									/>
								</div>
								<div class="col-md-6 col-12 mb-5">
									<VNumberInput
										v-model="aadtCreateStore.params.directional_distribution_factor"
										label="Directional distribution factor"
										name="directional_distribution_factor"
										:required="true"
										:precision="2"
									/>
								</div>
							</div>
						</div>
					</div>
				</div>
			</VSkeletonLoader>
		</div>
	</div>

	<VSkeletonLoader :loading="aadtRoadGroup.loading">
		<div class="row">
			<div class="col-12">
				<h2 class="fw-semibold fs-5 mb-5">Truck Factor</h2>
				<div class="table-responsive">
					<table class="table customize-basic-table mb-0">
						<thead>
							<tr>
								<th class="text-center fw-semibold" width="10%">จำนวนล้อ</th>
								<th class="text-center fw-semibold" width="35%">จำนวนเพลา</th>
								<th class="text-center fw-semibold">ปริมาณรถ</th>
								<th class="text-center fw-semibold">% Truck</th>
								<th class="text-center fw-semibold">Load Equivalent (L)</th>
								<th class="text-end fw-semibold">Truck Factor</th>
								<th class="text-center fw-bold ps-2">
									<VSwitch
										v-model="aadtCreateStore.params.is_truck_factor"
										name="is_truck_factor"
										value=""
										:disabled="false"
									/>
								</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td class="text-center">รถ 4 ล้อ</td>
								<td>
									<VNumberInput v-model="aadtCreateStore.params.four_wheel_axle_number" name="four_wheel_axle_number" />
								</td>
								<td>
									<VNumberInput
										v-model="aadtCreateStore.params.four_wheel_vehicle_volume"
										name="four_wheel_vehicle_volume"
										:disabled="true"
									/>
								</td>
								<td></td>
								<td></td>
								<td colspan="2"></td>
							</tr>
							<tr>
								<td class="text-center">6 ล้อ</td>
								<td>
									<VSelect
										v-model="aadtCreateStore.params.six_wheel_axle_number_id"
										:options="aadtRoadGroup.getOptions.six"
										name="six_wheel_axle_number_id"
										placeholder="เลือก"
										:can-clear="false"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtCreateStore.params.six_wheel_vehicle_volume"
										:disabled="true"
										name="six_wheel_vehicle_volume"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtCreateStore.params.six_wheel_percentage_truck"
										name="six_wheel_percentage_truck"
										:disabled="true"
										:precision="3"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtRoadGroup.getLoadEquivalent.six"
										:disabled="true"
										:precision="3"
										name="six_load_equivalent"
									/>
								</td>
								<td colspan="2">
									<VNumberInput
										v-model="aadtCreateStore.params.six_wheel_factor_result"
										:disabled="!aadtCreateStore.params.is_truck_factor"
										name="six_wheel_factor_result"
										:precision="3"
									/>
								</td>
							</tr>
							<tr>
								<td colspan="7" class="text-center py-8">
									<img :src="aadtRoadGroup.getVehicleImagePath.six" width="800" alt="" />
								</td>
							</tr>
							<tr>
								<td class="text-center">> 6 ล้อ</td>
								<td>
									<VSelect
										v-model="aadtCreateStore.params.ten_wheel_axle_number_id"
										:options="aadtRoadGroup.getOptions.ten"
										name="ten_wheel_axle_number_id"
										placeholder="เลือก"
										:can-clear="false"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtCreateStore.params.ten_wheel_vehicle_volume"
										name="ten_wheel_vehicle_volume"
										:disabled="true"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtCreateStore.params.ten_wheel_percentage_truck"
										name="ten_wheel_percentage_truck"
										:disabled="true"
										:precision="3"
									/>
								</td>
								<td>
									<VNumberInput
										v-model="aadtRoadGroup.getLoadEquivalent.ten"
										:disabled="true"
										:precision="3"
										name="ten_load_equivalent"
									/>
								</td>
								<td colspan="2">
									<VNumberInput
										v-model="aadtCreateStore.params.ten_wheel_factor_result"
										:disabled="!aadtCreateStore.params.is_truck_factor"
										name="ten_wheel_factor_result"
										:precision="3"
									/>
								</td>
							</tr>
							<tr>
								<td colspan="7" class="text-center py-8">
									<img :src="aadtRoadGroup.getVehicleImagePath.ten" width="800" alt="" />
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
		<div class="row mt-5">
			<div class="col-12 text-end">
				<VLoading :loading="aadtCreateStore.loading || aadtRoadGroup.loading" />
				<div>
					<BtnCancel @click="onCancel()" />
					<BtnSubmit label="บันทึก" :disabled="aadtCreateStore.loading" @click="onSubmit" />
				</div>
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
	padding-bottom: 0.75rem;
	padding-top: 0.75rem;
	vertical-align: middle;
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

.customize-basic-table tr td[colspan="7"] {
	border-top: 1px solid var(--kt-gray-300) !important;
	border-bottom: 1px solid var(--kt-gray-300) !important;
}

.customize-basic-table tr:last-of-type td[colspan="7"] {
	border-bottom: 0 !important;
}

.border-bottom {
	border-bottom: 1px solid var(--kt-gray-300) !important;
}

.yax {
	margin-top: -25px;
}
.elane {
	margin-top: -23px;
}
@media only screen and (max-width: 460px) {
	.yax {
		margin-top: -10px;
	}
	.elane {
		margin-top: -10px;
	}
}
</style>
