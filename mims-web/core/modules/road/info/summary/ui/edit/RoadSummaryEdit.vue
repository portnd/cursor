<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useRoadSummaryEditStore } from "../../store/RoadSummaryEdit"
import RoadTitle from "~~/core/modules/common/roadTitle/ui"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useRoadSummaryEditStore()
const roadTitleStore = useRoadTitleStore()
const route = useRoute()

useStoreLifecycle([store, roadTitleStore])
const roadId = route.params.roadId

onBeforeMount(() => {
	setTimeout(async () => {
		await store.getSurface(Number(roadId))
	}, 200)
})
onMounted(() => {
	roadTitleStore.tab = route.query.tab!.toString()
	// Clear All Errors
	setErrors({})
	setTimeout(() => {
		setErrors({})
	}, 300)
})

watch(
	() => roadTitleStore.data,
	(_) => {
		store.setDataRoad(
			roadTitleStore.data.road_info?.ref_direction_id,
			roadTitleStore.data.road_info?.km_start,
			roadTitleStore.data.road_info?.km_end
		)
	}
)

// ลบช่องจราจร
const deleteItem = (item: any, items: any) => {
	useDeleteItem({
		name: items.direction,
		callBack: function () {
			store.deleteLane(item, items)
		},
	})
}

// ลบคอลัมน์
const deleteItemColumn = (item: any) => {
	useDeleteItem({
		name: "",
		callBack: function () {
			store.deleteColumn(item)
		},
	})
}

// Validate Input
const handleValidate = computed(() => {
	const validations: IValidate = {}
	store.data.forEach((item, index: number) => {
		validations[`km_start${index}`] = `required|km|kmOverlapStart:${index}|kmStart:${index}|kmOverlap`
		validations[`km_end${index}`] = `required|km|kmOverlapEnd:${index}|kmOverlap`
		validations[`surface_cross_section_code${index}`] = "required"
		item.lane.forEach((_, i: number) => {
			validations[`direction${i}${index}`] = "direction"
		})
		validations[`surface_shoulder_left${index}`] = "required"
		validations[`surface_shoulder_right${index}`] = "required"
		validations[`width_surface${index}`] = "required"
		validations[`thickness_surface${index}`] = "required"
		validations[`thickness_surface_concrete${index}`] = "required"
		validations[`width_shoulder_left${index}`] = "required"
		validations[`width_shoulder_right${index}`] = "required"
		validations[`material_base${index}`] = "required"
		validations[`thickness_base${index}`] = "required"
	})

	return validations
})

defineRule("direction", (value: any) => {
	if (value === null || value === undefined) {
		return "โปรดระบุ"
	} else {
		return true
	}
})

defineRule("kmOverlapStart", (value: any, index: number) => {
	const a = store.convertStringToKm(value)
	const b = store.convertStringToKm(store.data[index].km_end)

	if (store.directionId === 2) {
		if (a < b) {
			return "กม. เริ่มต้น ต้องมีค่ามากกว่า กม. สิ้นสุด"
		} else if (a === b) {
			return "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		} else {
			return true
		}
	} else if (store.directionId === 1) {
		if (a > b) {
			return "กม. เริ่มต้น ต้องมีค่าน้อยกว่า กม. สิ้นสุด"
		} else if (a === b) {
			return "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		}
	}
	return true
})

defineRule("kmOverlapEnd", (value: any, index: number) => {
	const a = store.convertStringToKm(value)
	const b = store.convertStringToKm(store.data[index].km_start)

	if (store.directionId === 2) {
		if (a > b) {
			return "กม. สิ้นสุด ต้องมีค่าน้อยกว่า กม. เริ่มต้น"
		} else if (b === a) {
			return "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		}
	} else if (store.directionId === 1) {
		if (a < b) {
			return "กม. สิ้นสุด ต้องมีค่ามากกว่า กม. เริ่มต้น"
		} else if (b === a) {
			return "กม. เริ่มต้น และกม. สิ้นสุด ต้องมีค่าไม่เท่ากัน"
		}
	}
	return true
})

defineRule("kmOverlap", (value: any) => {
	const a = store.convertStringToKm(value)
	let overLap = null
	store.data.forEach((e) => {
		const kmStart = store.convertStringToKm(e.km_start)
		const kmEnd = store.convertStringToKm(e.km_end)
		if (store.directionId === 1) {
			if (kmStart < a && a < kmEnd) {
				overLap = true
			}
		} else if (kmStart > a && a > kmEnd) {
			overLap = true
		}
	})
	if (overLap === true) {
		return "มีช่วง กม. ที่ทับซ้อน (overlap) กัน"
	} else {
		return true
	}
})

defineRule("kmStart", (value: any, index: any) => {
	if (index[0] === "0") {
		const a = store.convertStringToKm(value)
		const b = store.roadKmStart
		if (store.directionId === 1) {
			if (a < b) {
				return "กม. เริ่มต้นไม่ได้อยู่ในช่วงของสายทาง"
			} else {
				return true
			}
		} else if (a > b) {
			return "กม. เริ่มต้นไม่ได้อยู่ในช่วงของสายทาง"
		} else {
			return true
		}
	}
	return true
})

defineRule("kmEnd", (value: any, index: any) => {
	if (index[0] === `${store.data.length - 1}`) {
		const a = store.convertStringToKm(value)
		const b = store.roadKmEnd

		if (b === 0) {
			return true
		}

		if (store.directionId === 1) {
			if (a > b) {
				return "กม. สิ้นสุดไม่ได้อยู่ในช่วงของสายทาง"
			} else {
				return true
			}
		} else if (a < b) {
			return "กม. สิ้นสุดไม่ได้อยู่ในช่วงของสายทาง"
		} else {
			return true
		}
	}
	return true
})

const { handleSubmit, isSubmitting, errors, handleReset, setErrors } = useForm({
	validationSchema: handleValidate,
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const kmTotalCheck = store.checkKmTotal()
	// บันทึก
	if (kmTotalCheck.status === true) {
		const res = await store.updateRoad()

		if (res?.status) {
			useHandlerSuccess(res.code, {
				showAlert: true,
				fn: function () {
					handleReset()
					store.getSurface(Number(roadId))
					navigateTo(`/roads/${roadId}/summary?tab=${route.query.tab}`)
				},
			})
		}
	} else {
		useHandlerError(0, { message: kmTotalCheck.Message }, { showAlert: true })
	}
})

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const surfaceOptions = computed(() => {
	const options = [{ value: -1, label: "ไม่มีผิวทาง" }]
	options.push(...toOptions(useInitData().refSurface()))
	return options
})

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12 mb-10">
			<form @submit.prevent="onSubmit">
				<div class="card p-5 pb-1 pt-2 mb-5">
					<div class="row">
						<div class="col-xl-8">
							<RoadTitle :road-id="Number(roadId)" :is-left-sum-position="true" />
						</div>
					</div>
					<VSkeletonLoader :loading="store.loading">
						<div class="row mb-5 mt-3">
							<div v-for="(item, index) in store.data" :key="index" class="col-12">
								<div class="table-responsive mb-5">
									<table class="table customize-basic-table mb-0">
										<thead>
											<tr>
												<th class="text-center" style="width: 18%">กม. เริ่มต้น-สิ้นสุด</th>
												<th class="text-center" style="width: 18%">ช่องจราจร</th>
												<th class="text-center" style="width: 14%">ผิวทาง</th>
												<th class="text-center" style="width: 14%">ไหล่ทาง</th>
												<th class="text-center" style="width: 26%">
													ชั้นทาง
													<BtnDelete
														v-show="index !== 0"
														class="float-end me-5 lh-0"
														@click="deleteItemColumn(item.id)"
													/>
												</th>
											</tr>
										</thead>
										<tbody>
											<tr>
												<td>
													<div class="row">
														<div class="col-12 mb-5">
															<VTextInput
																v-model="item.km_start"
																:required="true"
																placeholder="00+000"
																label="เริ่มต้น"
																:name="`km_start${index}`"
															/>
														</div>
														<div class="col-12 mb-5">
															<VTextInput
																v-model="item.km_end"
																:required="true"
																placeholder="00+000"
																label="สิ้นสุด"
																:name="`km_end${index}`"
															/>
														</div>
														<div class="col-12 mb-5">
															<VSelect
																v-model="item.surface_cross_section_code"
																:options="toOptions(useInitData().refStructureSurface())"
																label="หน้าตัดผิวทาง"
																:required="true"
																:name="`surface_cross_section_code${index}`"
																placeholder="เลือก"
																:auto-height="true"
															/>
														</div>
													</div>
												</td>
												<td>
													<div v-for="(laneItem, i) in item.lane" :key="i" class="row gx-2 mb-5">
														<div class="col-10">
															<VSelect
																v-model="laneItem.surface.id"
																:options="surfaceOptions"
																:label="`ช่องจราจร ${(i + 1).toString()}`"
																:name="`direction${i}${index}`"
																placeholder="เลือก"
																:required="true"
															/>
														</div>
														<div class="col-2 align-self-center" style="margin-top: 2.3em">
															<BtnDelete v-show="i !== 0" @click="() => deleteItem(index, laneItem)" />
														</div>
													</div>
													<div
														v-show="store.data[index].lane.length < roadTitleStore.data.road_geom?.length"
														class="row"
													>
														<div class="col-10 mb-5 mt-10">
															<button
																type="button"
																class="btn btn-primary rounded-4 px-3 py-3 fw-semibold btn-outline btn-outline-primary w-100 h-auto"
																@click="store.addLane(index)"
															>
																เพิ่มช่องจราจร
															</button>
														</div>
													</div>
												</td>
												<td>
													<div class="row">
														<div class="col-12 mb-5">
															<VNumberInput
																v-model="item.width_surface"
																:required="true"
																label="กว้าง"
																align="start"
																text-end="ม."
																:name="`width_surface${index}`"
																:precision="3"
															/>
														</div>
														<div class="col-12 mb-5">
															<VNumberInput
																v-model="item.thickness_surface"
																:required="true"
																label="ความหนาผิวทางลาดยาง"
																align="start"
																text-end="ซม."
																:name="`thickness_surface${index}`"
																:precision="3"
															/>
														</div>
														<div class="col-12 mb-5">
															<VNumberInput
																v-model="item.thickness_surface_concrete"
																:required="true"
																label="ความหนาผิวคอนกรีต"
																align="start"
																text-end="ซม."
																:name="`thickness_surface_concrete${index}`"
																:precision="3"
															/>
														</div>
													</div>
												</td>
												<td>
													<div class="row">
														<div class="col-12 mb-5">
															<VNumberInput
																v-model="item.width_shoulder_left"
																:required="true"
																label="ไหล่ทางซ้าย กว้าง"
																align="start"
																text-end="ม."
																:name="`width_shoulder_left${index}`"
																:precision="3"
															/>
														</div>
														<div class="col-12 mb-5">
															<VSelect
																v-model="item.surface_shoulder_left.id"
																:options="toOptions(useInitData().refSurface())"
																label="ไหล่ทางซ้าย ชนิดผิว"
																:required="true"
																:name="`surface_shoulder_left${index}`"
																placeholder="เลือก"
															/>
														</div>
														<div class="col-12 mb-5">
															<VNumberInput
																v-model="item.width_shoulder_right"
																:required="true"
																label="ไหล่ทางขวา กว้าง"
																align="start"
																text-end="ม."
																:name="`width_shoulder_right${index}`"
																:precision="3"
															/>
														</div>
														<div class="col-12 mb-5">
															<VSelect
																v-model="item.surface_shoulder_right.id"
																:options="toOptions(useInitData().refSurface())"
																label="ไหล่ทางขวา ชนิดผิว"
																:required="true"
																:name="`surface_shoulder_right${index}`"
																placeholder="เลือก"
															/>
														</div>
													</div>
												</td>
												<td>
													<div class="row mb-5">
														<div class="col">
															<VNumberInput
																v-model="item.thickness_concrete_slab"
																:required="false"
																label="ความหนา Concrete Slab"
																align="start"
																text-end="ซม."
																:name="`thickness_concrete_slab${index}`"
																:precision="3"
															/>
														</div>
													</div>
													<div class="row mb-5">
														<div class="col">
															<VSelect
																v-model="item.material_base.id"
																:options="toOptions(useInitData().refMaterialBase())"
																label="Base"
																:required="true"
																:name="`material_base${index}`"
																placeholder="เลือก"
															/>
														</div>
														<div class="col">
															<VNumberInput
																v-model="item.thickness_base"
																:required="true"
																label="ความหนา"
																align="start"
																text-end="ซม."
																:name="`thickness_base${index}`"
																:precision="3"
															/>
														</div>
													</div>
													<div class="row mb-5">
														<div class="col">
															<VSelect
																v-model="item.material_subbase.id"
																:options="toOptions(useInitData().refMaterialSubbase())"
																label="Subbase"
																:required="false"
																:name="`material_subbase${index}`"
																placeholder="เลือก"
															/>
														</div>
														<div class="col">
															<VNumberInput
																v-model="item.thickness_subbase"
																:required="false"
																label="ความหนา"
																align="start"
																text-end="ซม."
																:name="`thickness_subbase${index}`"
																:precision="3"
															/>
														</div>
													</div>
													<div class="row mb-5">
														<div class="col">
															<VSelect
																v-model="item.material_subgrade.id"
																:options="toOptions(useInitData().refMaterialSubgrade())"
																label="Subgrade"
																:required="false"
																:name="`material_subgrade${index}`"
																placeholder="เลือก"
															/>
														</div>
														<div class="col">
															<VNumberInput
																v-model="item.thickness_subgrade"
																:required="false"
																label="ความหนา"
																align="start"
																text-end="ซม."
																:name="`thickness_subgrade${index}`"
																:precision="3"
															/>
														</div>
													</div>
												</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>
							<div class="col-12 col-md-5 mt-3">
								<button
									type="button"
									class="btn btn-outline btn-outline-primary rounded-4 px-5 py-2 fw-semibold"
									@click="store.addRow()"
								>
									<i class="fi fi-rr-plus align-middle fs-8"></i>
									เพิ่มช่วง กม.
								</button>
							</div>
							<div class="col-12 col-md-7 text-end mt-3">
								<NuxtLink
									:to="`/roads/${roadId}/summary?tab=${route.query.tab}`"
									class="btn btn-light rounded-4 px-8 py-3 me-5 fw-semibold text-black"
								>
									ยกเลิก
								</NuxtLink>
								<BtnSubmit :loading="store.loading" :disabled="store.loading" label="บันทึก" />
							</div>
						</div>
					</VSkeletonLoader>
				</div>
			</form>
		</div>
	</div>
</template>

<style scoped>
@media only screen and (max-width: 768px) {
	.table th {
		min-width: 150px;
	}

	.table th:last-of-type {
		min-width: 300px;
	}
}

thead tr th {
	padding: 20px 0px;
	border: 0;
	background-color: var(--kt-gray-300);
}

thead tr th:first-of-type {
	border-top-left-radius: 7px;
}

thead tr th:last-of-type {
	border-top-right-radius: 7px;
}

tbody tr td {
	padding-left: 1rem !important;
	padding-right: 1rem !important;
	padding-bottom: 0.75rem;
	padding-top: 0.5rem;
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
