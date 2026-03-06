<script setup lang="ts">
import {
	IMaintenanceHistoryDetailMaintenanceRoad,
	IMaintenanceHistoryDetailMaintenanceRoadHistory,
} from "../../../infrastructure"

// Modal
const { $bootstrap }: any = useNuxtApp()
const modal = ref()
const roads = ref<IMaintenanceHistoryDetailMaintenanceRoadHistory | IMaintenanceHistoryDetailMaintenanceRoad>()

const generateLane = computed(() => {
	const lanes = Array.from({ length: Number(roads.value?.lane_total) }, (_, i) => i + 1)

	return roads.value?.ref_direction_id === 1 ? lanes.reverse() : lanes
})

const popoverVisible = reactive<{ [key: string]: boolean }>({})

function togglePopover(lane: number) {
	popoverVisible[`${lane}`] = true
}

const generateGrid = computed(() => {
	const grids = Array.from({ length: 4 }, (_, i) => i + 1)

	return roads.value?.ref_direction_id === 1 ? grids.reverse() : grids
})

const onCancel = () => {
	hideModal()
}

const showModal = (
	item: IMaintenanceHistoryDetailMaintenanceRoadHistory | IMaintenanceHistoryDetailMaintenanceRoad
) => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)

	// reset
	roads.value = {} as IMaintenanceHistoryDetailMaintenanceRoadHistory

	roads.value = item

	if (item.maintenance_type === 2) {
		togglePopover(Number(item.lane_no))
	}
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal-project-delete" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered" :class="generateLane.length >= 4 ? 'modal-lg' : 'modal-'">
			<div class="modal-content p-8">
				<div class="modal-header border-0 p-0 mb-8">
					<h3 class="modal-title fw-semibold fs-3">ตำแหน่งงานซ่อมบำรุง</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="d-flex justify-content-center ms-8 mt-8">
					<div
						class="text-center d-flex"
						:class="roads?.ref_direction_id === 1 ? 'roads-rescale-left' : 'roads-rescale-right'"
					>
						<div :style="{ display: roads?.ref_direction_id === 2 ? 'none' : '' }" class="road-container">
							<h3
								:style="{
									position: 'absolute',
									top: '-20px',
									left: '48%',
								}"
							>
								LT
							</h3>
							<div class="fs-4 km-top-left">{{ convertMeterToKm(roads?.km_end ?? 0) ?? "0+000" }}</div>
							<div class="section-line section-top-left"></div>
							<table>
								<thead>
									<tr>
										<th
											v-for="(item, headKey) of generateLane"
											:key="headKey"
											class="road-header border-left-3"
											:class="item === roads?.lane_total ? 'border' : 'border border-right-3'"
										></th>
									</tr>
								</thead>
								<tbody>
									<tr>
										<td
											v-for="(lane, laneKey) of generateLane"
											:key="laneKey"
											class="road-lanes p-0 border-right-3"
											:class="lane === roads?.lane_total ? 'border' : 'border border-right-3'"
											:style="{
												backgroundColor: Number(roads?.lane_no) === lane ? '#BD8C33' : '',
											}"
										>
											<div
												v-if="popoverVisible[`${lane}`] && roads?.maintenance_type === 2"
												class="popover row py-2 px-2"
											>
												<div class="col-12 border-bottom border-2 text-center">
													<span class="fs-2">ช่องจราจรที่ {{ lane }}</span>
												</div>
												<div class="col-12 mt-2">
													<div class="d-flex justify-content-between align-items-end">
														<span class="fs-4">กริดที่ {{ generateGrid[0] }}</span>
														<span class="fs-4">กริดที่ {{ generateGrid[generateGrid.length - 1] }}</span>
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
											<table v-show="roads?.maintenance_type === 2">
												<tr>
													<td
														v-for="(grid, gridIndex) of generateGrid"
														:key="`${lane}${gridIndex}`"
														class="lane-grid"
														:class="grid === 1 ? 'border-left-1' : ' border-right-1'"
														:style="{
															borderColor: '#BD8C33 !important',
															backgroundColor:
																`${lane}${grid}` === `${roads?.lane_no}${roads?.grid_no}` ? '#BD8C33' : '',
														}"
													></td>
												</tr>
											</table>
										</td>
									</tr>
								</tbody>
								<tfoot class="mt-1">
									<tr>
										<td
											v-for="(lane, footKey) of generateLane"
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
							<div class="fs-4 km-bottom-left">{{ convertMeterToKm(roads?.km_start ?? 0) ?? "0+000" }}</div>
						</div>
						<div>
							<div id="center" class="p-2 center border-right-3">
								<div class="section-center-line-top"></div>
								<div id="box" class="box"></div>
								<div class="section-center-line-bottom"></div>
							</div>
						</div>
						<div :style="{ display: roads?.ref_direction_id === 1 ? 'none' : '' }" class="road-container">
							<div class="fs-4 km-top-right">{{ convertMeterToKm(roads?.km_start ?? 0) ?? "0+000" }}</div>
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
											v-for="(item, headKey) of generateLane"
											:key="headKey"
											class="road-header border-left-3"
											:class="item === roads?.lane_total ? 'border' : 'border border-right-3'"
										></th>
									</tr>
								</thead>
								<tbody>
									<tr>
										<td
											v-for="(lane, laneKey) of generateLane"
											:key="laneKey"
											class="road-lanes p-0 border-left-3"
											:class="lane === roads?.lane_total ? '' : 'border border-right-3'"
											:style="{
												borderLeft: lane === 1 ? '1px solid #fff !important' : '',
												backgroundColor: Number(roads?.lane_no) === lane ? '#BD8C33' : '',
											}"
										>
											<div
												v-if="popoverVisible[`${lane}`] && roads?.maintenance_type === 2"
												class="popover row py-2 px-2"
											>
												<div class="col-12 border-bottom border-2 text-center">
													<span class="fs-2">ช่องจราจรที่ {{ lane }}</span>
												</div>
												<div class="col-12 mt-2">
													<div class="d-flex justify-content-between align-items-end">
														<span class="fs-4">กริดที่ {{ generateGrid[0] }}</span>
														<span class="fs-4">กริดที่ {{ generateGrid[generateGrid.length - 1] }}</span>
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
											<table v-show="roads?.maintenance_type === 2">
												<tr>
													<td
														v-for="(grid, gridIndex) of generateGrid"
														:key="`${lane}${gridIndex}`"
														class="lane-grid"
														:class="grid === 4 ? 'border-left-1' : ' border-right-1'"
														:style="{
															borderColor: '#BD8C33 !important',
															backgroundColor:
																`${lane}${grid}` === `${roads?.lane_no}${roads?.grid_no}` ? '#BD8C33' : '',
														}"
													></td>
												</tr>
											</table>
										</td>
									</tr>
								</tbody>
								<tfoot class="mt-1">
									<tr>
										<td
											v-for="(lane, footKey) of generateLane"
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
							<div class="fs-4 km-bottom-right">{{ convertMeterToKm(roads?.km_end ?? 0) ?? "0+000" }}</div>
						</div>
					</div>
				</div>
				<div class="col-12 text-end m-0 mt-10">
					<BtnCancel label="ปิด" class="m-0" @click="onCancel" />
				</div>
			</div>
		</div>
	</div>
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

.roads-rescale-left,
.roads-rescale-right {
	position: relative;
	transform: scale(0.85);
}

.section-end-left,
.km-bottom-left,
.section-center-line-bottom,
.section-end-right {
	position: absolute;
}

@media (max-width: 1400px) {
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
