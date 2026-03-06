<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { useMaintenanceRoadsCreateStore } from "../../store"
import { IValidate } from "../../../../../shared/types/Validate"

const store = useMaintenanceRoadsCreateStore()
useStoreLifecycle(store)

const route = useRoute()
const id = Number(route.params.id)
const sumDistance = ref(0.0)
const popoverVisible = reactive<{ [key: string]: boolean }>({})

function togglePopover(lane: number) {
	popoverVisible[`${lane}`] = !popoverVisible[`${lane}`]
}

onMounted(async () => {
	store.loading = true
	await store.getIsShowMethod(id)
	await store.getRoadGroupOptions()
	await store.getInterventionCriteria()
	store.loading = false
})

const toggleLane = (laneNumber: number) => {
	store.lane_no = laneNumber
}
const toggleGridNumber = (lane: number, gridIndex: number) => {
	toggleLane(lane)
	store.grid_no = gridIndex
}

const updateSumDistance = () => {
	const kmEnd = store.km_end ? convertStringToKm(store.km_end) : 0
	const kmStart = store.km_start ? convertStringToKm(store.km_start) : 0
	const sum = Math.abs(kmEnd - kmStart) / 1000

	store.sum_distance = isNaN(sum) ? 0.0 : parseFloat(sum.toFixed(3))
	sumDistance.value = isNaN(sum) ? 0.0 : parseFloat(sum.toFixed(3))
}

const validates = computed(() => {
	const validate: IValidate = {}
	validate.km_start = "km|required|km-start-range"
	validate.km_end = "km|required|km-end-range"
	validate.road_id = "required"
	validate.intervention_criteria_id = store.is_show_method ? "required" : ""

	return validate
})

defineRule("km-start-range", (value: string) => {
	const direction = store.refDirectionId
	const kmStart = store.based_km_start ?? 0
	const kmValue = convertStringToKm(value)

	const message = "กม.เริ่มต้นไม่อยู่ในช่วงสายทาง"

	if ((direction === 1 && kmValue < kmStart) || (direction === 2 && kmValue > kmStart)) {
		return message
	}

	return ""
})

defineRule("km-end-range", (value: string) => {
	const direction = store.refDirectionId
	const kmEnd = store.based_km_end ?? 0
	const kmValue = convertStringToKm(value)

	const message = "กม.สิ้นสุดไม่อยู่ในช่วงสายทาง"

	if ((direction === 1 && kmValue > kmEnd) || (direction === 2 && kmValue < kmEnd)) {
		return message
	}

	return ""
})

const { handleReset, errors, isSubmitting, handleSubmit, resetField } = useForm({ validationSchema: validates })

watch(
	() => store.road_id,
	() => {
		if (!store.road_id) {
			resetField("km_start", { errors: undefined })
			resetField("km_end", { errors: undefined })
		}
	}
)

const onSubmit = handleSubmit(async (_, action) => {
	useAction(action)

	if (store.maintenance_type === 2 && store.grid_no === null) {
		showAlert({ type: "warning", title: "แจ้งเตือน", message: "กรุณาเลือกกริด" })
	} else if (store.maintenance_type === 1 && store.lane_no === null) {
		showAlert({ type: "warning", title: "แจ้งเตือน", message: "กรุณาเลือกช่องจราจร" })
	} else {
		const res = await store.createMaintenanceRoadsInfo(id)

		if (res?.status) {
			useHandlerSuccess(res.code, {
				showAlert: true,
				fn: () => {
					navigateTo(`/maintenances/history/${id}/info`)
					handleReset()
				},
			})
		}
	}
})

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

onUnmounted(() => {
	handleReset()
	store.$reset()
})

const onCanel = () => {
	navigateTo(`/maintenances/history/${id}/info`)
}
</script>
<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="row card p-5 py-8">
			<div class="col-12 mb-3">
				<h5>ข้อมูลการซ่อมบำรุง</h5>
			</div>
			<form class="row pe-0" @submit="onSubmit">
				<div class="col-md-4 col-12">
					{{ store.road_id }}
					<VTree
						v-model="store.road_id"
						label="สายทาง"
						name="road_id"
						placeholder="เลือก"
						mode="LEAF_PRIORITY"
						:disable-branch-nodes="true"
						:options="store.roadOptions"
						:required="true"
						@update:model-value="store.onUpdateRoadId"
					/>
				</div>
				<div class="col-md-4 col-12" :style="store.is_show_method === false ? 'visibility: hidden !important' : ''">
					<VTree
						v-model="store.intervention_criteria_id"
						label="วิธีการซ่อมบำรุง"
						name="intervention_criteria_id"
						placeholder="เลือก"
						:disable-branch-nodes="true"
						:options="store.interventionCriteria"
						:required="store.is_show_method ? true : false"
					/>
				</div>
				<div class="col-md-4-col-12"></div>
				<div class="col-md-4 col-12 mt-4">
					<div class="row align-item-end">
						<div class="col">
							<VTextInput
								v-model="store.km_start"
								name="km_start"
								label="ช่วง กม."
								:validate-position="true"
								:required="true"
								:disabled="!store.road_id"
								@update:model-value="updateSumDistance"
							/>
						</div>
						<span class="col-1 align-self-center mt-9 text-center">-</span>
						<div class="col align-self-end">
							<VTextInput
								v-model="store.km_end"
								name="km_end"
								:disabled="!store.road_id"
								:validate-position="true"
								:required="true"
								@update:model-value="updateSumDistance"
							/>
						</div>
					</div>
				</div>
				<div class="col-md-4 col-12 mt-4">
					<VNumberInput
						v-model="store.sum_distance"
						label="ระยะทาง (กม.)"
						name="distance"
						:precision="3"
						:disabled="true"
					/>
				</div>
				<div class="col-md-4 col-12 mt-4 mb-10">
					<VRadio
						v-model="store.maintenance_type"
						label="ประเภทการซ่อมบำรุง"
						:options="[
							{ label: 'ซ่อมบำรุงแบบเลน', value: 1 },
							{ label: 'ซ่อมบำรุงแบบกริด', value: 2 },
						]"
						name="maintenance_type"
						:required="true"
						@update:model-value="store.onUpdateMaintenanceType"
					/>
				</div>
				<!-- <div class="col-md-8 col-12 mt-4">
          <VUploadFile
            ref="upLoadFile"
            label="เอกสารขอการเบิกจ่าย (รองรับ 10 ไฟล์)"
            max-file-size="10MB"
            name="attrachment"
            :required="true"
            :multiple="true"
            aspect-ratio="0.4"
            :max-files="10"
            :accepted-file-types="[
              'image/png',
              'image/jpg',
              'image/jpeg',
              'application/pdf',
              'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
              'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
              '.dwg',
            ]"
            :file-validate-type-detect-type="handleFileValidation"
          />
        </div> -->
				<div v-show="store.refDirectionId" class="col-md-12 mt-10 ms-8">
					<div
						class="col-12 text-center d-flex"
						:class="store.refDirectionId === 1 ? 'roads-rescale-left' : 'roads-rescale-right'"
					>
						<div class="col-md-1"></div>
						<div v-show="store.refDirectionId === 1" class="road-container">
							<div class="fs-4 km-top-left">{{ store.km_end ?? "0+000" }}</div>
							<div class="section-line section-top-left"></div>
							<h3
								:style="{
									position: 'absolute',
									top: '-20px',
									left: '47.6%',
								}"
							>
								LT
							</h3>
							<table>
								<thead>
									<tr>
										<th
											v-for="(item, headKey) of store.getGenerateLane"
											:key="headKey"
											class="road-header border-left-3"
											:class="item === store.totalLane ? 'border' : 'border border-right-3'"
										></th>
									</tr>
								</thead>
								<tbody>
									<tr>
										<td
											v-for="(lane, laneKey) of store.getGenerateLane"
											:key="laneKey"
											class="road-lanes p-0 cursor-pointer border-right-3"
											:class="lane === store.totalLane ? 'border' : 'border border-right-3'"
											:style="{
												backgroundColor: store.lane_no === lane ? '#BD8C33' : '',
											}"
											@mouseenter="() => togglePopover(lane)"
											@mouseleave="() => togglePopover(lane)"
											@click="toggleLane(lane)"
										>
											<div
												v-if="popoverVisible[`${lane}`] && store.maintenance_type === 2"
												class="popover row py-2 px-2"
											>
												<div class="col-12 border-bottom border-2 text-center">
													<span class="fs-2">ช่องจราจรที่ {{ lane }}</span>
												</div>
												<div class="col-12 mt-2">
													<div class="d-flex justify-content-between align-items-end">
														<span class="fs-4">กริดที่ {{ store.getGenerateGrid[0] }}</span>
														<span class="fs-4"
															>กริดที่ {{ store.getGenerateGrid[store.getGenerateGrid.length - 1] }}</span
														>
													</div>
												</div>
												<div class="col-12 text-center">
													<span>
														<img src="/images/icons/svg/left-arrow.svg" alt="left-arrow-icon" />
													</span>
													<div class="fs-7">
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
													</div>
												</div>
											</div>
											<table v-show="store.maintenance_type === 2">
												<tr>
													<td
														v-for="(grid, gridIndex) of store.getGenerateGrid"
														:key="`${lane}${gridIndex}`"
														class="lane-grid cursor-pointer"
														:class="grid === 1 ? 'border-left-1' : ' border-right-1'"
														:style="{
															borderColor: '#BD8C33 !important',
															backgroundColor: `${lane}${grid}` === `${store.lane_no}${store.grid_no}` ? '#BD8C33' : '',
														}"
														@click.stop="toggleGridNumber(lane, grid)"
													></td>
												</tr>
											</table>
										</td>
									</tr>
								</tbody>
								<tfoot class="mt-1">
									<tr>
										<td
											v-for="(lane, footKey) of store.getGenerateLane"
											:key="footKey"
											class="roads-footer m-0 p-0 pt-1 border-right-3 border-light"
										>
											<Icon id="icon" class="arrow-icon" name="typcn:arrow-up-outline" />
											<p class="text-white fs-2 mb-0">
												{{ `L${lane}` }}
											</p>
										</td>
									</tr>
								</tfoot>
							</table>
							<div class="section-line section-end-left"></div>
							<div class="fs-4 km-bottom-left">{{ store.km_start ?? "0+000" }}</div>
						</div>
						<div class="">
							<div id="center" class="p-2 center border-right-3">
								<div class="section-center-line-top"></div>
								<div id="box" class="box"></div>
								<div class="section-center-line-bottom"></div>
							</div>
						</div>
						<div v-show="store.refDirectionId === 2" class="road-container">
							<div class="fs-4 km-top-right">{{ store.km_start ?? "0+000" }}</div>
							<div class="section-line section-top-right"></div>
							<h3
								:style="{
									position: 'absolute',
									top: '-20px',
									left: '47.6%',
								}"
							>
								RT
							</h3>
							<table>
								<thead>
									<tr>
										<th
											v-for="(item, headKey) of store.getGenerateLane"
											:key="headKey"
											class="road-header border-left-3"
											:class="item === store.totalLane ? 'border' : 'border border-right-3'"
										></th>
									</tr>
								</thead>
								<tbody>
									<tr>
										<td
											v-for="(lane, laneKey) of store.getGenerateLane"
											:key="laneKey"
											class="road-lanes p-0 cursor-pointer border-left-3"
											:class="lane === store.totalLane ? '' : 'border border-right-3'"
											:style="{
												borderLeft: lane === 1 ? '1px solid #fff !important' : '',
												backgroundColor: store.lane_no === lane ? '#BD8C33' : '',
											}"
											@mouseenter="() => togglePopover(lane)"
											@mouseleave="() => togglePopover(lane)"
											@click="toggleLane(lane)"
										>
											<div
												v-if="popoverVisible[`${lane}`] && store.maintenance_type === 2"
												class="popover row py-2 px-2"
											>
												<div class="col-12 border-bottom border-2 text-center">
													<span class="fs-2">ช่องจราจรที่ {{ lane }}</span>
												</div>
												<div class="col-12 mt-2">
													<div class="d-flex justify-content-between align-items-end">
														<span class="fs-4">กริดที่ {{ store.getGenerateGrid[0] }}</span>
														<span class="fs-4"
															>กริดที่ {{ store.getGenerateGrid[store.getGenerateGrid.length - 1] }}</span
														>
													</div>
												</div>
												<div class="col-12 text-center">
													<span style="display: inline-block; transform: rotate(180deg)">
														<img src="/images/icons/svg/left-arrow.svg" alt="left-arrow-icon" />
													</span>
													<div class="fs-7">
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
														<span class="me-1">25%</span>
													</div>
												</div>
											</div>
											<table v-show="store.maintenance_type === 2">
												<tr>
													<td
														v-for="(grid, gridIndex) of store.getGenerateGrid"
														:key="`${lane}${gridIndex}`"
														class="lane-grid cursor-pointer"
														:class="grid === 4 ? 'border-left-1' : ' border-right-1'"
														:style="{
															borderColor: '#BD8C33 !important',
															backgroundColor: `${lane}${grid}` === `${store.lane_no}${store.grid_no}` ? '#BD8C33' : '',
														}"
														@click.stop="toggleGridNumber(lane, grid)"
													></td>
												</tr>
											</table>
										</td>
									</tr>
								</tbody>
								<tfoot class="mt-1">
									<tr>
										<td
											v-for="(lane, footKey) of store.getGenerateLane"
											:key="footKey"
											class="roads-footer m-0 p-0 pt-1 border-left-3 border-light"
										>
											<Icon id="icon" class="arrow-icon" name="typcn:arrow-down-outline" />
											<p class="text-white fs-2 mb-0">
												{{ `R${lane}` }}
											</p>
										</td>
									</tr>
								</tfoot>
							</table>
							<div class="section-line section-end-right"></div>
							<div class="fs-4 km-bottom-right">{{ store.km_end ?? "0+000" }}</div>
						</div>
					</div>
				</div>
				<div class="col-12 text-end p-0 mt-10">
					<BtnCancel @click="onCanel" />
					<BtnSubmit label="บันทึก" />
				</div>
			</form>
		</div>
	</VSkeletonLoader>
</template>
<style lang="scss" scoped>
.popover {
	position: absolute;
	top: -15%;
	width: 184px;
	background-color: white;
	border: 1px solid #ccc;
	border-radius: 16px;
	z-index: 100;
	box-shadow: none !important;
	transform: scale(0.8);
}

.road-container {
	position: relative;
}

.road-header {
	background-color: #181c32;
	width: 6em;
	height: 5em;
	border-bottom: 6px solid #ff0000 !important;
}

.box {
	height: 100%;
	width: 0.8em;
	background-color: #fdb833;
}

.center {
	position: relative;
	height: 100%;
	background-color: #181c32;
	border-top: 1px solid #fff;
	/* border-bottom: 3px solid #fff; */
}

.road-lanes {
	height: 180px;
	background-color: #181c32;
	border-bottom: 6px solid #ff0000 !important;
}

.roads-footer {
	background-color: #181c32;
	width: 6em;
	height: 5em;
}

.arrow-icon {
	color: #fff;
}

svg.icon {
	font-size: 30px;
}

.lane-grid {
	height: 174px;
	width: 10em;
	background-color: #181c32;
}

.section-center-line-top {
	position: absolute;
	background-color: #ff0000;
	width: 26px;
	top: 61.8px;
	right: -2.1px;
	height: 6px;
}
.section-center-line-bottom {
	position: absolute;
	background-color: #ff0000;
	width: 26px;
	bottom: 62.5px;
	right: -2.1px;
	height: 6px;
}

.section-line {
	width: 26px;
	height: 6px;
	background-color: #ff0000;
}

.section-top-left {
	position: absolute;
	/* bottom: 21%; */
	top: 62.8px;
	left: -20px;
}

.km-top-left {
	position: absolute;
	top: 17%;
	left: -80px;
}

.km-bottom-left {
	position: absolute;
	bottom: 17%;
	left: -80px;
}

.km-top-right {
	position: absolute;
	top: 17%;
	right: -80px;
}

.km-bottom-right {
	position: absolute;
	bottom: 17%;
	right: -80px;
}

.section-end-left {
	position: absolute;
	bottom: 62.4px;
	left: -20px;
}

.section-top-right {
	position: absolute;
	top: 62.8px;
	right: -20px;
}

.section-end-right {
	position: absolute;
	bottom: 62.5px;
	right: -20px;
}

.section-end-left,
.km-bottom-left,
.section-center-line-bottom,
.section-end-right {
	position: absolute;
}

@media (max-width: 1400px) {
	.roads-rescale-left,
	.roads-rescale-right {
		position: relative;
		transform: scale(0.85);
	}

	.roads-rescale-left {
		left: -3%;
	}
	.roads-rescale-right {
		left: -10%;
	}
	.section-end-left {
		bottom: 63.1px;
	}
	.section-center-line-bottom {
		bottom: 62.5px;
	}
	.km-bottom-left {
		bottom: 16.5%;
	}
}

@media (max-width: 1378px) {
	.section-end-left {
		bottom: 62.5px;
	}
}
@media (max-width: 1120px) {
	.section-end-left {
		bottom: 63.1px;
	}
}

@media (max-width: 1180px) {
	.roads-rescale-right {
		left: -10%;
	}
	.section-end-left {
		bottom: 62.1px;
	}
}

@media (width: 1180px) {
	.section-end-left {
		bottom: 63.1px;
	}
}

@media (max-width: 950px) {
	.roads-rescale-left,
	.roads-rescale-right {
		position: relative;
		transform: scale(0.85);
	}
	.roads-rescale-left {
		left: 3%;
	}
	.roads-rescale-right {
		left: -13%;
	}
	.section-end-left,
	.section-center-line-bottom {
		bottom: 62.1px;
	}
	.section-end-right {
		bottom: 62px;
	}
}

@media (max-width: 890px) {
	.roads-rescale-right {
		left: -14%;
	}
	.section-center-line-bottom {
		bottom: 61.5px;
	}
	.section-end-left {
		bottom: 62px;
	}
}

@media (max-width: 768px) {
	.roads-rescale-right {
		left: -5%;
	}
	.section-center-line-bottom {
		bottom: 62px;
	}
	.section-end-left {
		bottom: 61.1px;
	}
}

@media (max-width: 700px) {
	.roads-rescale-right {
		left: -5%;
	}
	.section-end-left,
	.km-bottom-left,
	.section-center-line-bottom {
		bottom: 61.1px;
	}
}

@media (max-width: 450px) {
	.roads-rescale-right {
		left: -14%;
	}
	.section-end-right {
		bottom: 62px;
	}
}
</style>
