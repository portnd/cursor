<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadWorkEffectACStore } from "../../store/RoadWorkEffectACStore"
import { IValidate } from "~/core/shared/types/Validate"

onMounted(async () => {
	await store.getACValue()
})

const store = useRoadWorkEffectACStore()
useStoreLifecycle(store)
const params = store.params

const handleValidate = computed(() => {
	const validations: IValidate = {}
	for (const key in params) {
		validations[key] = "required"
	}

	return validations
})

const { handleSubmit } = useForm({ validationSchema: handleValidate })
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.postACValue()
	if (res?.status) {
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getACValue()
			},
		})
	}
})

const onCancel = async () => {
	store.setDefaultParams()
	await store.getACValue()
}

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="row mt-3">
			<div class="col-md-6 col-12">
				<div class="table-responsive">
					<table class="table customize-basic-table bd-bottom-responsive mb-0">
						<thead>
							<tr>
								<th class="text-center fw-semibold">ตัวแปร</th>
								<th class="text-center fw-semibold">ค่า Default หลังการซ่อม <span class="required"></span></th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td colspan="2" class="fw-semibold py-5">งานเสริมผิวทาง (OL-Overlay)</td>
							</tr>
							<tr>
								<td>OVERLAY A0</td>
								<td><VNumberInput v-model="params.as_ol_overlay_a0" :precision="2" name="as_ol_overlay_a0" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ARVb (Raveling)</td>
								<td><VNumberInput v-model="params.as_ol_ar_vb" :precision="2" name="as_ol_ar_vb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า APOTb (Pothole)</td>
								<td><VNumberInput v-model="params.as_ol_po_tb" :precision="2" name="as_ol_po_tb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ACAb (Cracking)</td>
								<td><VNumberInput v-model="params.as_ol_ac_ab" :precision="2" name="as_ol_ac_ab" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า RDMb (Rut depth)</td>
								<td><VNumberInput v-model="params.as_ol_rd_mb" :precision="2" name="as_ol_rd_mb" /></td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">ฉาบผิวลาดยาง (SS-Slurry Seal)</td>
							</tr>
							<tr>
								<td>RWE SS Model A0</td>
								<td>
									<VNumberInput v-model="params.as_ss_rwe_ss_model_a0" :precision="2" name="as_ss_rwe_ss_model_a0" />
								</td>
							</tr>
							<tr>
								<td>Default lower bound IRI after slurry seal</td>
								<td>
									<VNumberInput
										v-model="params.as_ss_default_lower_bound_iri_after_slurry_seal"
										:precision="2"
										name="as_ss_default_lower_bound_iri_after_slurry_seal"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ARVb (Raveling)</td>
								<td><VNumberInput v-model="params.as_ss_ar_vb" :precision="2" name="as_ss_ar_vb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า APOTb (Pothole)</td>
								<td><VNumberInput v-model="params.as_ss_apo_tb" :precision="2" name="as_ss_apo_tb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ACAb (Cracking)</td>
								<td><VNumberInput v-model="params.as_ss_ac_ab" :precision="2" name="as_ss_ac_ab" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า RDMb (Rut depth)</td>
								<td><VNumberInput v-model="params.as_ss_rd_mb" :precision="2" name="as_ss_rd_mb" /></td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">งานขูดไสและปูทับด้วยแอสฟัลต์คอนกรีต (M&OL-Mill&Overlay)</td>
							</tr>
							<tr>
								<td>IRI after Mill&Overlay</td>
								<td>
									<VNumberInput
										v-model="params.as_mol_iri_after_mill_overlay"
										:precision="2"
										name="as_mol_iri_after_mill_overlay"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ARVb (Raveling)</td>
								<td><VNumberInput v-model="params.as_mol_ar_vb" :precision="2" name="as_mol_ar_vb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า APOTb (Pothole)</td>
								<td><VNumberInput v-model="params.as_mol_apo_tb" :precision="2" name="as_mol_apo_tb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ACAb (Cracking)</td>
								<td><VNumberInput v-model="params.as_mol_ac_ab" :precision="2" name="as_mol_ac_ab" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า RDMb (Rut depth)</td>
								<td><VNumberInput v-model="params.as_mol_rd_mb" :precision="2" name="as_mol_rd_mb" /></td>
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
								<th class="text-center fw-semibold">ค่า Default หลังการซ่อม <span class="required"></span></th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td colspan="2" class="fw-semibold py-5">
									การหมุนเวียนวัสดุชั้นทางเดิมและปูผิวทางใหม่ (RCL-Recycling)
								</td>
							</tr>
							<tr>
								<td>SNC</td>
								<td><VNumberInput v-model="params.as_rcl_snc" :precision="2" name="as_rcl_snc" /></td>
							</tr>
							<tr>
								<td>IRI after recycling</td>
								<td>
									<VNumberInput
										v-model="params.as_rcl_iri_after_recycling"
										:precision="2"
										name="as_rcl_iri_after_recycling"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ARVb (Raveling)</td>
								<td><VNumberInput v-model="params.as_rcl_ar_vb" :precision="2" name="as_rcl_ar_vb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า APOTb (Pothole)</td>
								<td><VNumberInput v-model="params.as_rcl_apo_tb" :precision="2" name="as_rcl_apo_tb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ACAb (Cracking)</td>
								<td><VNumberInput v-model="params.as_rcl_ac_ab" :precision="2" name="as_rcl_ac_ab" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า RDMb (Rut depth)</td>
								<td><VNumberInput v-model="params.as_rcl_rd_mb" :precision="2" name="as_rcl_rd_mb" /></td>
							</tr>
							<tr>
								<td>Default HS OLD</td>
								<td>
									<VNumberInput v-model="params.as_rcl_default_hs_old" :precision="2" name="as_rcl_default_hs_old" />
								</td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">Rehabilitation of Asphalt Pavement</td>
							</tr>
							<tr>
								<td>SNC</td>
								<td><VNumberInput v-model="params.as_rc_snc" :precision="2" name="as_rc_snc" /></td>
							</tr>
							<tr>
								<td>IRI after reconstruction</td>
								<td>
									<VNumberInput
										v-model="params.as_rc_iri_after_reconstruction"
										:precision="2"
										name="as_rc_iri_after_reconstruction"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ARVb (Raveling)</td>
								<td><VNumberInput v-model="params.as_rc_ar_vb" :precision="2" name="as_rc_ar_vb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า APOTb (Pothole)</td>
								<td><VNumberInput v-model="params.as_rc_apo_tb" :precision="2" name="as_rc_apo_tb" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า ACAb (Cracking)</td>
								<td><VNumberInput v-model="params.as_rc_ac_ab" :precision="2" name="as_rc_ac_ab" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า RDMb (Rut depth)</td>
								<td><VNumberInput v-model="params.as_rc_rd_mb" :precision="2" name="as_rc_rd_mb" /></td>
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
					<BtnCancel @click="onCancel()" />
					<BtnSubmit label="บันทึก" :disabled="store.loading" @click="onSubmit" />
				</div>
			</div>
		</div>
	</VSkeletonLoader>
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
