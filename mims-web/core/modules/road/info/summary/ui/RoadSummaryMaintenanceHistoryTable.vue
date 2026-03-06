<script setup lang="ts">
// import { api as viewerApi } from "v-viewer"
import {
	IMaintenanceProjectsMaintenanceRoad,
	// IMaintenanceProjectsAttacchment,
} from "../infrastructure/RoadSummaryModel"
import { useRoadSummaryStore } from "../store"

const store = useRoadSummaryStore()

const props = defineProps({
	maintenanceData: {
		type: Array<IMaintenanceProjectsMaintenanceRoad>,
		default: [],
	},
	maintenanceHistoryData: {
		type: Array<IMaintenanceProjectsMaintenanceRoad>,
		default: [],
	},
})

// const showImages = (attacchments: IMaintenanceProjectsAttacchment[]) => {
// 	const images: string[] = []
// 	attacchments.forEach((attacchment) => {
// 		images.push(attacchment.path)
// 	})

// 	viewerApi({
// 		options: { title: true, toolbar: true, navbar: true },
// 		images,
// 	})
// }

const maintenanceSum = computed(() => {
	return props.maintenanceData.reduce((acc, item) => acc + Math.abs(item.km_start - item.km_end), 0) / 1000
})

const maintenanceHistSum = computed(() => {
	return props.maintenanceHistoryData.reduce((acc, item) => acc + Math.abs(item.km_start - item.km_end), 0) / 1000
})
</script>

<template>
	<div class="row">
		<div class="col-12 mb-10">
			<h5 class="fw-semibold">ข้อมูลการซ่อมบำรุง</h5>
			<div class="table-responsive mt-3">
				<table class="table customize-basic-table mb-0 text-truncate table-hover">
					<thead>
						<tr>
							<th class="text-center">ลำดับ</th>
							<th class="text-center">ช่วง</th>
							<th class="text-center">ช่องจราจร</th>
							<th class="text-center">วิธีการซ่อมบำรุง</th>
							<th class="text-center">ช่วง กม.</th>
							<th class="text-center">ระยะทาง (กม.)</th>
						</tr>
					</thead>
					<tbody v-if="maintenanceData.length > 0">
						<template v-for="(item, index) of maintenanceData" :key="item.id">
							<tr class="cursor-pointer cursor-hover" @click="store.toLocation(item.the_geom)">
								<td class="text-center">{{ index + 1 }}</td>
								<td class="text-center">{{ item.road_name && item.road_name !== "" ? item.road_name : "-" }}</td>
								<td class="text-center">{{ item.lane_no }}</td>
								<td class="text-center">
									{{
										item.intervention_criteria.maintenance_standard_name
											? item.intervention_criteria.maintenance_standard_name
											: "-"
									}}
								</td>
								<td class="text-center">{{ convertMeterToKm(item.km_start) }} - {{ convertMeterToKm(item.km_end) }}</td>
								<td class="text-center">{{ (Math.abs(item.km_end - item.km_start) / 1000).toFixed(2) }}</td>
							</tr>
						</template>
						<tr>
							<td colspan="4"></td>
							<td class="text-center">ระยะทาง</td>
							<td class="text-center">
								{{ toNumber(maintenanceSum, 2) }}
							</td>
						</tr>
					</tbody>
					<tbody v-else>
						<tr class="text-center">
							<td colspan="6">ไม่พบข้อมูล</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<div class="col-12">
			<h5 class="fw-semibold">ประวัติการซ่อมบำรุง ในช่วงค้ำประกันโครงการ</h5>
			<div class="table-responsive">
				<table class="table customize-basic-table mb-0 text-truncate table-hover">
					<thead>
						<tr>
							<th class="text-center">ลำดับ</th>
							<th class="text-center">ช่วง</th>
							<th class="text-center">ช่องจราจร</th>
							<th class="text-center">วิธีการซ่อมบำรุง</th>
							<th class="text-center">ช่วง กม.</th>
							<th class="text-center">ระยะทาง (กม.)</th>
							<!-- <th class="text-center">รูปภาพ</th> -->
						</tr>
					</thead>
					<tbody v-if="maintenanceHistoryData?.length > 0">
						<tr v-for="(item, index) of maintenanceHistoryData" :key="index" class="align-middle cursor-pointer">
							<td class="text-center" @click="store.toLocation(item.the_geom)">{{ index + 1 }}</td>
							<td class="text-center" @click="store.toLocation(item.the_geom)">
								{{ item.road_name }}
							</td>
							<td class="text-center" @click="store.toLocation(item.the_geom)">{{ item.lane_no }}</td>
							<td class="text-center" @click="store.toLocation(item.the_geom)">
								{{
									item.intervention_criteria.maintenance_standard_name
										? item.intervention_criteria.maintenance_standard_name
										: "-"
								}}
							</td>
							<td class="text-center" @click="store.toLocation(item.the_geom)">
								{{ convertMeterToKm(item.km_start) }} - {{ convertMeterToKm(item.km_end) }}
							</td>
							<td class="text-center" @click="store.toLocation(item.the_geom)">
								{{ Math.abs((item.km_start - item.km_end) / 1000).toFixed(2) }}
							</td>
							<!-- <td class="text-center">
								<div
									v-show="item.attacchments?.length > 0"
									class="cursor-pointer align-middle lh-0"
									@click="showImages(item.attacchments)"
								>
									<i class="fi fi-sr-picture fs-4 text-gray-600"></i>
								</div>
							</td> -->
						</tr>
						<tr>
							<td colspan="4"></td>
							<td class="text-center">ระยะทาง</td>
							<td class="text-center">
								{{ toNumber(maintenanceHistSum, 2) }}
							</td>
						</tr>
					</tbody>
					<tbody v-else>
						<tr class="text-center">
							<td colspan="7">ไม่พบข้อมูล</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
