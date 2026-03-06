<script setup lang="ts">
import { useRoadDamageStore } from "../store/RoadDamageStore"
import { IRoadDamageRange } from "../infrastructure/RoadDamageModel"

const store = useRoadDamageStore()

defineProps({
	data: {
		type: Array as PropType<IRoadDamageRange[]>,
	},
})

</script>

<template>
	<div class="row mt-5">
		<div class="col-12">
			<div class="card card-rounded p-5 pt-2">
				<ul class="nav nav-tabs nav-line-tabs mb-5 me-8 justify-content-start justify-content-md-start">
					<li class="nav-item active" data-bs-toggle="tab" data-bs-target="#detail-ac" role="tab" aria-selected="true">
						<span class="nav-link cursor-pointer">ลาดยาง</span>
						<span class="line"></span>
					</li>
					<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-concrete" role="tab" aria-selected="false">
						<span class="nav-link cursor-pointer">คอนกรีต</span>
						<span class="line"></span>
					</li>
				</ul>

				<div class="tab-content">
					<div id="detail-ac" class="tab-pane fade active show" role="tabpanel">
						<div class="mt-1">
							<div class="table-responsive">
								<table class="table customize-basic-table fixed-head mb-0">
									<thead>
										<tr>
											<th>
												<span>&nbsp;</span>
											</th>
											<th class="text-center text-nowrap">ช่วง กม.</th>
											<th class="text-center text-nowrap">รอยแตกต่อเนื่อง (ตร.ม.)</th>
											<th class="text-center text-nowrap">รอยแตกไม่ต่อเนื่อง (ม.)</th>
											<th class="text-center text-nowrap">ผิวหลุดร่อน (ตร.ม.)</th>
											<th class="text-center text-nowrap">รอยปะซ่อม (ตร.ม.)</th>
											<th class="text-center text-nowrap">หลุมบ่อ (ตร.ม.)</th>
											<th class="text-center text-nowrap">การเยิ้มของลาดยาง (ตร.ม.)</th>
											<th class="text-center text-nowrap">หลุมบ่อ (หลุม)</th>
										</tr>
									</thead>
									<tbody v-if="data?.length === 0">
										<tr class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
									<tbody v-else>
										<template v-for="parent in data" :key="parent.id">
											<tr
												class="cursor-pointer expand collapsed"
												data-bs-toggle="collapse"
												:data-bs-target="`#collapse${parent.id}`"
												aria-expanded="false"
											>
												<td class="text-center">
													<i class="fi fi-br-caret-down fs-4 lh-0 text-gray-600"></i>
												</td>
												<td class="text-center text-nowrap">
													{{ convertMeterToKm(parent.km_start) }} - {{ convertMeterToKm(parent.km_end) }}
												</td>
												<td class="text-end">
													{{ parent.ac_icrack === null ? "-" : parent.ac_icrack }}
												</td>
												<td class="text-end">{{ parent.ac_ucrack === null ? "-" : parent.ac_ucrack }}</td>
												<td class="text-end">{{ parent.ac_ravelling === null ? "-" : parent.ac_ravelling }}</td>
												<td class="text-end">{{ parent.ac_patching === null ? "-" : parent.ac_patching }}</td>
												<td class="text-end">{{ parent.ac_pothole_area === null ? "-" : parent.ac_pothole_area }}</td>
												<td class="text-end">
													{{ parent.ac_bleeding === null ? "-" : parent.ac_bleeding }}
												</td>
												<td class="text-end">{{ parent.ac_pothole_count === null ? "-" : parent.ac_pothole_count }}</td>
											</tr>

											<tr
												v-for="child in parent.road_damage_m"
												:id="`collapse${parent.id}`"
												:key="child.id"
												class="collapse hover cursor-pointer"
												@click="store.setDetails(child)"
											>
												<td class="text-center">
													<span class="fs-5 text-gray-600">|</span>
												</td>
												<td class="text-center">{{ convertMeterToKm(child.km) }}</td>
												<td class="text-end">{{ child.ac_icrack === null ? "-" : child.ac_icrack }}</td>
												<td class="text-end">{{ child.ac_ucrack === null ? "-" : child.ac_ucrack }}</td>
												<td class="text-end">{{ child.ac_ravelling === null ? "-" : child.ac_ravelling }}</td>
												<td class="text-end">{{ child.ac_patching === null ? "-" : child.ac_patching }}</td>
												<td class="text-end">{{ child.ac_pothole_area === null ? "-" : child.ac_pothole_area }}</td>
												<td class="text-end">{{ child.ac_bleeding === null ? "-" : child.ac_bleeding }}</td>
												<td class="text-end">{{ child.ac_pothole_count === null ? "-" : child.ac_pothole_count }}</td>
											</tr>
										</template>
									</tbody>
								</table>
							</div>
						</div>
					</div>
					<div id="detail-concrete" class="tab-pane fade" role="tabpanel">
						<div class="mt-1">
							<div class="table-responsive">
								<table class="table customize-basic-table fixed-head mb-0">
									<thead>
										<tr>
											<th>
												<span>&nbsp;</span>
											</th>
											<th class="text-center text-nowrap">ช่วง กม.</th>
											<th class="text-center text-nowrap">รอยแตกตามขวาง (ม.)</th>
											<th class="text-center text-nowrap">รอยแตกตามยาวและแนวทแยง (ม.)</th>
											<th class="text-center text-nowrap">รอยแตกที่มุม (จุด)</th>
											<th class="text-center text-nowrap">วัสดุยาแนวรอยต่อเสียหาย (ม.)</th>
											<th class="text-center text-nowrap">รอยปะซ่อม (ตร.ม.)</th>
											<th class="text-center text-nowrap">รอยบิ่นกะเทาะ (ตร.ม.)</th>
											<th class="text-center text-nowrap">ผิวทางหลุดร่อน (ตร.ม.)</th>
										</tr>
									</thead>
									<tbody v-if="data?.length === 0">
										<tr class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
									<tbody v-else>
										<template v-for="parent in data" :key="parent.id">
											<tr
												class="cursor-pointer expand collapsed"
												data-bs-toggle="collapse"
												:data-bs-target="`#collapse${parent.id}`"
												aria-expanded="false"
											>
												<td class="text-center">
													<i class="fi fi-br-caret-down fs-4 lh-0 text-gray-600"></i>
												</td>
												<td class="text-end text-nowrap">
													{{ convertMeterToKm(parent.km_start) }} - {{ convertMeterToKm(parent.km_end) }}
												</td>
												<td class="text-end">
													{{ parent.cc_transverse_crack === null ? "-" : parent.cc_transverse_crack }}
												</td>
												<td class="text-end">
													{{ parent.cc_non_transverse_crack === null ? "-" : parent.cc_non_transverse_crack }}
												</td>
												<td class="text-end">{{ parent.cc_corner_break === null ? "-" : parent.cc_corner_break }}</td>
												<td class="text-end">
													{{ parent.cc_joint_seal_damage === null ? "-" : parent.cc_joint_seal_damage }}
												</td>
												<td class="text-end">{{ parent.cc_patching === null ? "-" : parent.cc_patching }}</td>
												<td class="text-end">
													{{ parent.cc_spalling === null ? "-" : parent.cc_spalling }}
												</td>
												<td class="text-end">{{ parent.cc_scaling === null ? "-" : parent.cc_scaling }}</td>
											</tr>

											<tr
												v-for="child in parent.road_damage_m"
												:id="`collapse${parent.id}`"
												:key="child.id"
												class="collapse hover cursor-pointer"
												@click="store.setDetails(child)"
											>
												<td class="text-center">
													<span class="fs-5 text-gray-600">|</span>
												</td>
												<td class="text-center">{{ convertMeterToKm(child.km) }}</td>
												<td class="text-end">
													{{ child.cc_transverse_crack === null ? "-" : child.cc_transverse_crack }}
												</td>
												<td class="text-end">
													{{ child.cc_non_transverse_crack === null ? "-" : child.cc_non_transverse_crack }}
												</td>
												<td class="text-end">{{ child.cc_corner_break === null ? "-" : child.cc_corner_break }}</td>
												<td class="text-end">
													{{ child.cc_joint_seal_damage === null ? "-" : child.cc_joint_seal_damage }}
												</td>
												<td class="text-end">{{ child.cc_patching === null ? "-" : child.cc_patching }}</td>
												<td class="text-end">
													{{ child.cc_spalling === null ? "-" : child.cc_spalling }}
												</td>
												<td class="text-end">{{ child.cc_scaling === null ? "-" : child.cc_scaling }}</td>
											</tr>
										</template>
									</tbody>
								</table>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.table.fixed-head {
	border-radius: 8px;
	padding: 0.25rem;
	th {
		position: sticky;
		top: 0;
		&:first-of-type span {
			width: 40px;
			display: block;
		}
		padding: 1rem 1.25rem;

		font-weight: 500;
		background-color: #fff;
	}
	.expand {
		background-color: #f4f4f4f4 !important;
	}
}
.table-responsive {
	overflow-y: auto;
}

.hover:hover {
	background: #d9d9d9;
}

.active_child {
	background-color: #d9d9d9;
}
</style>
