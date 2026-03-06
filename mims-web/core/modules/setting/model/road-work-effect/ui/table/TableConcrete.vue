<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadWorkEffectCCStore } from "../../store/RoadWorkEffectCCStore"
import { IValidate } from "~/core/shared/types/Validate"

const store = useRoadWorkEffectCCStore()
useStoreLifecycle(store)
const params = store.params

onMounted(async () => {
	await store.getCCValue()
})

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

	const res = await store.postCCValue()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getCCValue()
			},
		})
	}
})

const onCancel = async () => {
	store.setDefaultParams()
	await store.getCCValue()
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
								<td colspan="2" class="fw-semibold py-5">Full Depth Repair (FDR)</td>
							</tr>
							<tr>
								<td>IRI after FDR</td>
								<td>
									<VNumberInput v-model="params.cc_fdr_iri_after_fdr" :precision="2" name="cc_fdr_iri_after_fdr" />
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Faulting</td>
								<td><VNumberInput v-model="params.cc_fdr_faulting" :precision="2" name="cc_fdr_faulting" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Cracking</td>
								<td><VNumberInput v-model="params.cc_fdr_cracking" :precision="2" name="cc_fdr_cracking" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Spalling</td>
								<td><VNumberInput v-model="params.cc_fdr_spalling" :precision="2" name="cc_fdr_spalling" /></td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">Bonded Concrete Overlay (BCO)</td>
							</tr>
							<tr>
								<td>IRI after BCO</td>
								<td>
									<VNumberInput v-model="params.cc_bco_iri_after_bco" :precision="2" name="cc_bco_iri_after_bco" />
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Faulting</td>
								<td><VNumberInput v-model="params.cc_bco_faulting" :precision="2" name="cc_bco_faulting" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Cracking</td>
								<td><VNumberInput v-model="params.cc_bco_cracking" :precision="2" name="cc_bco_cracking" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Spalling</td>
								<td><VNumberInput v-model="params.cc_bco_spalling" :precision="2" name="cc_bco_spalling" /></td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">Rehabilitation of Concrete Pavement (RB)</td>
							</tr>
							<tr>
								<td>IRI after RB</td>
								<td>
									<VNumberInput v-model="params.cc_rbc_iri" :precision="2" name="cc_rbc_iri" />
								</td>
							</tr>
							<tr>
								<td>slab thickness default</td>
								<td><VNumberInput v-model="params.cc_rbc_slabthk" :precision="2" name="cc_rbc_slabthk" /></td>
							</tr>
							<tr>
								<td>เปอร์เซนต์คงเหลือของค่า Faulting</td>
								<td>
									<VNumberInput
										v-model="params.cc_rbc_percent_faulting"
										:precision="2"
										name="cc_rbc_percent_faulting"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซนต์คงเหลือของค่า Spalling</td>
								<td>
									<VNumberInput
										v-model="params.cc_rbc_percent_spalling"
										:precision="2"
										name="cc_rbc_percent_spalling"
									/>
								</td>
							</tr>
							<tr>
								<td>เปอร์เซนต์คงเหลือของค่า Cracking</td>
								<td>
									<VNumberInput
										v-model="params.cc_rbc_percent_cracking"
										:precision="2"
										name="cc_rbc_percent_cracking"
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
								<th class="text-center fw-semibold">ค่า Default หลังการซ่อม <span class="required"></span></th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td colspan="2" class="fw-semibold py-5">Concrete Mil&Overlay (M-OL)</td>
							</tr>
							<tr>
								<td>IRI after M-OL</td>
								<td>
									<VNumberInput v-model="params.cc_mol_iri_after_mol" :precision="2" name="cc_mol_iri_after_mol" />
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Faulting</td>
								<td><VNumberInput v-model="params.cc_mol_faulting" :precision="2" name="cc_mol_faulting" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Cracking</td>
								<td><VNumberInput v-model="params.cc_mol_cracking" :precision="2" name="cc_mol_cracking" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Spalling</td>
								<td><VNumberInput v-model="params.cc_mol_spalling" :precision="2" name="cc_mol_spalling" /></td>
							</tr>
							<tr>
								<td colspan="2" class="fw-semibold py-5">งานฉาบผิว (Seal)</td>
							</tr>
							<tr>
								<td>IRI after Seal</td>
								<td>
									<VNumberInput v-model="params.cc_seal_iri_after_seal" :precision="2" name="cc_seal_iri_after_seal" />
								</td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Faulting</td>
								<td><VNumberInput v-model="params.cc_seal_faulting" :precision="2" name="cc_seal_faulting" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Cracking</td>
								<td><VNumberInput v-model="params.cc_seal_cracking" :precision="2" name="cc_seal_cracking" /></td>
							</tr>
							<tr>
								<td>เปอร์เซ็นต์คงเหลือของค่า Spalling</td>
								<td><VNumberInput v-model="params.cc_seal_spalling" :precision="2" name="cc_seal_spalling" /></td>
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
