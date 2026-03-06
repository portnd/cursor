<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useVehicleTypeStore } from "../../store"
import { IRoadData } from "~/core/modules/road/roadList/infrastructure"
import { useRoadListStore } from "~/core/modules/road/roadList/store"
import { IOption } from "~/core/shared/types/Option"
import { IValidate } from "~/core/shared/types/Validate"

const roadListStore = useRoadListStore()
const store = useVehicleTypeStore()
useStoreLifecycle([roadListStore, store])
const carTypeModel = useCarTypeModel()

onMounted(async () => {
	await roadListStore.getData()
	store.id = roadListStore.roads[0].id
	await store.get(store.id)
})

defineRule("checkSumWheel", (_: any, [wheel]: [string]) => {
	let sum
	switch (wheel) {
		case "4":
			sum =
				Number(store.data.four_wheel.car_less_than_equal_seven) +
				Number(store.data.four_wheel.car_over_than_seven) +
				Number(store.data.four_wheel.light_bus) +
				Number(store.data.four_wheel.light_truck)
			break
		case "6-10":
			sum = Number(store.data.six_to_ten_wheel.medium_bus) + Number(store.data.six_to_ten_wheel.medium_truck)
			break
		case "> 10":
			sum =
				Number(store.data.over_ten_wheel.full_trailor) +
				Number(store.data.over_ten_wheel.heavy_bus) +
				Number(store.data.over_ten_wheel.heavy_truck) +
				Number(store.data.over_ten_wheel.semi_trailor)

			break
	}
	if (sum !== 1) {
		return `ผลรวมประเภทรถ ${wheel} ล้อต้องเท่ากับ 1`
	} else {
		return true
	}
})

const handleValidate = computed(() => {
	const validations: IValidate = {}
	validations.car_less_than_equal_seven = `required|checkSumWheel:4`
	validations.car_over_than_seven = `required|checkSumWheel:4`
	validations.light_bus = `required|checkSumWheel:4`
	validations.light_truck = `required|checkSumWheel:4`
	validations.medium_bus = `required|checkSumWheel:6-10`
	validations.medium_truck = `required|checkSumWheel:6-10`
	validations.heavy_bus = `required|checkSumWheel:> 10`
	validations.heavy_truck = `required|checkSumWheel:> 10`
	validations.full_trailor = `required|checkSumWheel:> 10`
	validations.semi_trailor = `required|checkSumWheel:> 10`
	return validations
})

const { handleSubmit, isSubmitting, errors } = useForm({ validationSchema: handleValidate })
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.post()
	if (res) {
		useHandlerSuccess(res.code, {
			showAlert: true,
		})
	}
})

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onCancle = async () => {
	await roadListStore.getData()
	store.id = roadListStore.roads[0].id
	store.get(store.id)
}

const generateOptionTable = (item: IRoadData[]) => {
	const data: IOption[] = []
	const selectedForm = item
	if (Array.isArray(selectedForm)) {
		selectedForm.map((e) => data.push({ value: e.id, label: e.name }))
	}
	return data
}

</script>

<template>
	<form @submit="onSubmit">
		<div class="row">
			<div class="col-12 mb-5">
				<div class="row">
					<div class="col-md-6 col-12">
						<VSelect
							v-model="store.id"
							name="roads"
							:options="generateOptionTable(roadListStore.roads)"
							label="สายทาง"
							placeholder="เลือกสายทาง"
							:required="true"
							:searchable="true"
							:can-clear="false"
							:can-deselect="false"
							@update:model-value="($e:any) => store.get($e)"
						/>
					</div>
				</div>
			</div>
		</div>
		<VSkeletonLoader :loading="store.loading">
			<div class="row">
				<div class="col-md-6 col-12">
					<div class="table-responsive">
						<table class="table customize-basic-table bd-bottom-responsive mb-0">
							<thead>
								<tr>
									<th class="text-center fw-semibold">ตัวแปร</th>
									<th class="text-center fw-semibold">อัตราส่วน <span class="required"></span></th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td class="fw-semibold py-5">รถ 4 ล้อ</td>
									<td class="fw-semibold text-center py-5"></td>
								</tr>
								<tr>
									<td>
										Car &lt;= 7
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.car_less_than_equal_seven }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.four_wheel.car_less_than_equal_seven"
											:precision="2"
											:min="0"
											name="car_less_than_equal_seven"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Car > 7
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.car_over_than_seven }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.four_wheel.car_over_than_seven"
											:precision="2"
											:min="0"
											name="car_over_than_seven"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Light Bus
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.light_bus }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput v-model="store.data.four_wheel.light_bus" :min="0" :precision="2" name="light_bus" />
									</td>
								</tr>
								<tr>
									<td>
										Light Truck
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.light_truck }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.four_wheel.light_truck"
											:min="0"
											:precision="2"
											name="light_truck"
										/>
									</td>
								</tr>
								<tr>
									<td class="fw-semibold py-5">รถ 6 ล้อ</td>
									<td class="fw-semibold text-center py-5"></td>
								</tr>
								<tr>
									<td>
										Medium Bus
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.medium_bus }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.six_to_ten_wheel.medium_bus"
											:min="0"
											:precision="2"
											name="medium_bus"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Medium Truck
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.medium_truck }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.six_to_ten_wheel.medium_truck"
											:precision="2"
											:min="0"
											name="medium_truck"
										/>
									</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
				<div class="col-md-6 col-12">
					<div class="table-responsive">
						<table class="table customize-basic-table bd-top-responsive mb-0">
							<thead class="thead-hidden">
								<tr>
									<th class="text-center fw-semibold">ตัวแปร</th>
									<th class="text-center fw-semibold">อัตราส่วน <span class="required"></span></th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td class="fw-semibold py-5">รถ > 6 ล้อ</td>
									<td class="fw-semibold text-center py-5"></td>
								</tr>
								<tr>
									<td>
										Heavy Bus
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.heavy_bus }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.over_ten_wheel.heavy_bus"
											:min="0"
											:precision="2"
											name="heavy_bus"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Heavy Truck
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.heavy_truck }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.over_ten_wheel.heavy_truck"
											:min="0"
											:precision="2"
											name="heavy_truck"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Full Trailer
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.full_trailor }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.over_ten_wheel.full_trailor"
											:min="0"
											:precision="2"
											name="full_trailor"
										/>
									</td>
								</tr>
								<tr>
									<td>
										Semi-Trailer
										<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
											<i class="fi fi-sr-interrogation fs-5"></i>
											<template #content>
												<div>{{ carTypeModel.semi_trailor }}</div>
											</template>
										</VPopover>
									</td>
									<td>
										<VNumberInput
											v-model="store.data.over_ten_wheel.semi_trailor"
											:min="0"
											:precision="2"
											name="semi_trailor"
										/>
									</td>
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
						<BtnCancel @click="onCancle()" />
						<BtnSubmit label="บันทึก" :disabled="store.loading" />
					</div>
				</div>
			</div>
		</VSkeletonLoader>
	</form>
</template>

<style scoped>
@media only screen and (max-width: 768px) {
	thead.thead-hidden {
		display: none;
	}

	table.bd-top-responsive {
		border-top-left-radius: 0px !important;
		border-top-right-radius: 0px !important;
	}

	table.bd-bottom-responsive {
		border-bottom-left-radius: 0px !important;
		border-bottom-right-radius: 0px !important;
	}
}

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

tbody tr td.fw-semibold {
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
