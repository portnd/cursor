<script setup lang="ts">
import { useRoadDamageStore } from "../store"
const store = useRoadDamageStore()

</script>

<template>
	<div class="row">
		<div class="col-md-6">
			<div class="card card-rounded p-4 mt-5">
				<table class="table customize-basic-table mb-0">
					<thead>
						<tr>
							<th colspan="3" class="text-center">ลาดยาง</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>รอยแตกต่อเนื่อง</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_icrack }}</td>
							<td class="text-center">ตร.ม.</td>
						</tr>
						<tr>
							<td>รอยแตกไม่ต่อเนื่อง</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_ucrack }}</td>
							<td class="text-center">ม.</td>
						</tr>
						<tr>
							<td>ผิวหลุดร่อน</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_ravelling }}</td>
							<td class="text-center">ตร.ม.</td>
						</tr>
						<tr>
							<td>รอยปะซ่อม</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_patching }}</td>
							<td class="text-center">ตร.ม.</td>
						</tr>
						<tr>
							<td>หลุมบ่อ</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_pothole_area }}</td>
							<td class="text-center">ตร.ม.</td>
						</tr>
						<tr>
							<td>การเยิ้มของลาดยาง</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_bleeding }}</td>
							<td class="text-center">ตร.ม.</td>
						</tr>
						<tr>
							<td>หลุมบ่อ</td>
							<td class="text-center">{{ store.data?.road_damage?.ac_pothole_count }}</td>
							<td class="text-center">หลุม</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<div class="col-md-6">
			<div class="card card-rounded p-4 mt-5">
				<div>
					<table class="table customize-basic-table mb-0">
						<thead>
							<tr>
								<th colspan="3" class="text-center">คอนกรีต</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>รอยแตกตามขวาง</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_transverse_crack }}</td>
								<td class="text-center">ม.</td>
							</tr>
							<tr>
								<td>รอยแตกตามยาวและแนวทแยง</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_non_transverse_crack }}</td>
								<td class="text-center">ม.</td>
							</tr>
							<tr>
								<td>รอยแตกที่มุม</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_corner_break }}</td>
								<td class="text-center">จุด</td>
							</tr>
							<tr>
								<td>วัสดุยาแนวรอยต่อเสียหาย</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_joint_seal_damage }}</td>
								<td class="text-center">ม.</td>
							</tr>
							<tr>
								<td>รอยปะซ่อม</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_patching }}</td>
								<td class="text-center">ตร.ม.</td>
							</tr>
							<tr>
								<td>รอยบิ่นกะเทาะ</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_spalling }}</td>
								<td class="text-center">ตร.ม.</td>
							</tr>
							<tr>
								<td>ผิวหลุดร่อน</td>
								<td class="text-center">{{ store.data?.road_damage?.cc_scaling }}</td>
								<td class="text-center">ตร.ม.</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
