<script setup lang="ts">
import { useForm } from "vee-validate"
import { useTableConcreteStore } from "../../store"
import { ITableConcrete } from "../../infrastructure"
import { IOption } from "~/core/shared/types/Option"
import { IRoadData } from "~/core/modules/road/roadList/infrastructure"
import { useRoadListStore } from "~/core/modules/road/roadList/store"
import { IValidate } from "~/core/shared/types/Validate"

const roadListStore = useRoadListStore()
const store = useTableConcreteStore()

onBeforeMount(async () => {
	store.$reset()

	await roadListStore.getData()
	store.id = roadListStore.roads[0].id
	store.data = {} as ITableConcrete
	store.get(store.id)
})

const handleValidate = computed(() => {
	const validations: IValidate = {}
	validations.b_stress = `required`
	validations.dwl_cor = `required`
	validations.ec = `required`
	validations.fi = `required`
	validations.jt_space = `required`
	validations.kjrc = `required`
	validations.kjrf = `required`
	validations.kjrr = `required`
	validations.kjrs = `required`
	validations.mi = `required`
	validations.p_steel = `required`
	validations.pred_seal = `required`
	validations.widened = `required`
	return validations
})

const { handleSubmit } = useForm({ validationSchema: handleValidate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.post()
	if (res) {
		useHandlerSuccess(res.code, {
			showAlert: true,
		})
	}
})

const generateOptionTable = (item: IRoadData[]) => {
	const data: IOption[] = []
	const selectedForm = item
	if (Array.isArray(selectedForm)) {
		selectedForm.map((e) => data.push({ value: e.id, label: e.name }))
	}
	return data
}

const onCancel = () => {
	store.get(store.id)
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
							@update:model-value="($e: number) => store.get($e)"
						/>
					</div>
				</div>
			</div>
		</div>
		<VSkeletonLoader :loading="store.loading">
			<div class="row">
				<div class="col-12">
					<div class="table-responsive">
						<table class="table customize-basic-table mb-0">
							<thead>
								<tr>
									<th width="10%" class="text-center fw-semibold">ตัวแปร</th>
									<th class="text-center fw-semibold">คำอธิบาย</th>
									<th class="text-center fw-semibold">ค่าคงที่ <span class="required"></span></th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td class="text-center">PSTEEL</td>
									<td>เปอร์เซ็นต์การเสริมเหล็กแนวยาว</td>
									<td><VNumberInput v-model="store.data.p_steel" name="p_steel" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">Ec</td>
									<td>ค่าโมดูลัสความยืดหยุ่นของคอนกรีต</td>
									<td><VNumberInput v-model="store.data.ec" name="ec" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">MI</td>
									<td>ดัชนีความชื้น Thornthwaite</td>
									<td><VNumberInput v-model="store.data.mi" name="mi" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">FI</td>
									<td>ค่าดัชนีแช่แข็ง</td>
									<td><VNumberInput v-model="store.data.fi" name="fi" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>jrc</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kjrc" name="kjrc" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">BSTRESS</td>
									<td>ความเค้นรับน้ำหนักคอนกรีตสูงสุด ในระบบเดือยคอนกรีต</td>
									<td><VNumberInput v-model="store.data.b_stress" name="b_stress" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">JTSPACE</td>
									<td>ระยะห่างระหว่างรอยต่อตามขวางเฉลี่ย</td>
									<td><VNumberInput v-model="store.data.jt_space" name="jt_space" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>jrf</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kjrf" name="kjrf" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">WIDENED</td>
									<td>การขยายช่องจราจร</td>
									<td><VNumberInput v-model="store.data.widened" name="widened" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">PREFSEAL</td>
									<td>มี pre-formed sealant ในรอยต่อ = 1, ไม่มี= 0</td>
									<td><VNumberInput v-model="store.data.pred_seal" name="pred_seal" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">DWLCOR</td>
									<td>มีการป้องกันการกร่อยของเหล็กเดือย = 1, ไม่มี =0</td>
									<td><VNumberInput v-model="store.data.dwl_cor" name="dwl_cor" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>jrs</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kjrs" name="kjrs" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>jrr</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kjrr" name="kjrr" :precision="4" /></td>
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
						<BtnSubmit label="บันทึก" :disabled="store.loading" />
					</div>
				</div>
			</div>
		</VSkeletonLoader>
	</form>
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
	padding-bottom: 0.5rem;
	padding-top: 0.5rem;
	vertical-align: middle;
}

tbody tr td[colspan="2"] {
	background-color: var(--kt-gray-300);
}
/* @media only screen and (max-width: 768px) {
	.customize-basic-table {
		width: max-content;
	}
} */

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
