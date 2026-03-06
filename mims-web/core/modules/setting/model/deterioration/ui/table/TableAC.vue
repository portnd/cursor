<script setup lang="ts">
import { useForm } from "vee-validate"
import { useTableAcStore } from "../../store"
import { ITableAC } from "../../infrastructure"
import { useRoadListStore } from "~/core/modules/road/roadList/store"
import { IRoadData } from "~/core/modules/road/roadList/infrastructure"
import { IOption } from "~/core/shared/types/Option"
import { IValidate } from "~/core/shared/types/Validate"

const roadListStore = useRoadListStore()
const store = useTableAcStore()

onBeforeMount(async () => {
	store.$reset()

	await roadListStore.getData()
	store.id = roadListStore.roads[0].id
	store.data = {} as ITableAC
	store.get(store.id)
})

const handleValidate = computed(() => {
	const validations: IValidate = {}
	validations.cdb = `required`
	validations.cds = `required`
	validations.comp = `required`
	validations.kcia = `required`
	validations.kciw = `required`
	validations.kcpa = `required`
	validations.kcpw = `required`
	validations.kgm = `required`
	validations.kgp = `required`
	validations.kpi = `required`
	validations.kpp = `required`
	validations.krid = `required`
	validations.krpd = `required`
	validations.krst = `required`
	validations.kvi = `required`
	validations.kvp = `required`
	validations.tlf = `required`
	validations.cmod = `required`
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
									<td class="text-center">TLF</td>
									<td>ปัจจัยด้านเวลา</td>
									<td><VNumberInput v-model="store.data.tlf" name="tlf" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">CDB</td>
									<td>ตัวบ่งชี้ข้อบกพร่องของการก่อสร้างชั้นพื้นทาง</td>
									<td><VNumberInput v-model="store.data.cdb" name="cdb" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">CDS</td>
									<td>ตัวบ่งชี้ข้อบกพร่องของการก่อสร้างพื้นผิวบิทูมินัส</td>
									<td><VNumberInput v-model="store.data.cds" name="cds" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">COMP</td>
									<td>ค่าบดอัดสัมพัทธ์</td>
									<td><VNumberInput v-model="store.data.comp" name="comp" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>vi</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kvi" name="kvi" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>vp</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kvp" name="kvp" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>pi</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kpi" name="kpi" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>pp</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kpp" name="kpp" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>rid</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.krid" name="krid" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>rst</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.krst" name="krst" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>rpd</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.krpd" name="krpd" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>gp</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kgp" name="kgp" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>gm</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kgm" name="kgm" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>cia</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kcia" name="kcia" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">CMOD</td>
									<td>ค่าโมดูลลัสการคืนตัวของดินผสมซีเมนต์</td>
									<td><VNumberInput v-model="store.data.cmod" name="cmod" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>ciw</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kciw" name="kciw" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>cpa</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kcpa" name="kcpa" :precision="4" /></td>
								</tr>
								<tr>
									<td class="text-center">K<sub>cpw</sub></td>
									<td>ค่าปรับแก้</td>
									<td><VNumberInput v-model="store.data.kcpw" name="kcpw" :precision="4" /></td>
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
