<script setup lang="ts">
import { useForm } from "vee-validate"
import { useStrategicEditStore } from "../../../store"
import { CopySelectRoadTable, CopyStrategicStep } from ".."
import { IValidate } from "~/core/shared/types/Validate"

const store = useStrategicEditStore()
const route = useRoute()
const id = Number(route.params.id)
const dataTable = ref()

interface IPrepareDataTable {
	loadData: Function
	searchData: Function
}

onMounted(async () => {
	store.isCopy = true
	if (Object.keys(store.default).length === 0) {
		handleReset()
		store.loading = true

		await store.getRoadsOptions()
		await store.getRefMaintenanceOptions()
		await store.getDefaultData(id)
		if (!store.isCopy) {
			await store.getSelectedId()
		}
		await store.getAllId()

		store.loading = false
	}

	await dataTable.value?.searchData()

	console.log(store.isCopy)
})

const handleDataTable = (table: IPrepareDataTable) => {
	dataTable.value = table
}

const surfaceOptions = [
	{ label: "ลาดยาง", value: 1 },
	{ label: "คอนกรีต", value: 2 },
]

const laneOptions = [
	{ label: "ทั้งหมด", value: 0 },
	{ label: "1", value: 1 },
	{ label: "2", value: 2 },
	{ label: "3", value: 3 },
]

const groupKmOptions = [
	{ label: "1 กม.", value: 1 },
	{ label: "2 กม.", value: 2 },
	{ label: "5 กม.", value: 5 },
	{ label: "10 กม.", value: 10 },
]

const validate = computed(() => {
	const validation: IValidate = {}

	validation.name = "required"
	validation.road_id = "required"
	validation.surface_type_id = "required"
	validation.lane_type_id = "required"
	validation.group_km = "required"

	return validation
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.createPrepareData()

	if (res?.status) {
		await store.checkPrepareDataStatus()
		await store.getAllId()
		await dataTable.value?.searchData()

		store.wasStep2 = false
	}
})
</script>

<template>
	<!-- ทำเป็นการ์ดเดียว -->
	<div class="row">
		<CopyStrategicStep />
		<div class="col-md-6 col-12 mb-2">
			<VTextInput v-model="store.default.name" label="ชื่อรายการ" name="name" :required="true" />
		</div>
		<div class="col-md-6 col-12 mb-2">
			<VTree
				v-model="store.road_id"
				label="สายทาง"
				:multiple="true"
				:options="store.roadGroupOptions"
				placeholder="เลือก"
				name="road_id"
				:required="true"
				:limit="0"
				:mode="'LEAF_PRIORITY'"
			/>
		</div>
	</div>
	<div class="row">
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.surface_type_id"
				:options="surfaceOptions"
				name="surface_type_id"
				label="ชนิดผิวทาง"
				placeholder="เลือก"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.lane_type_id"
				:options="laneOptions"
				name="lane_type_id"
				label="ช่องจราจร"
				placeholder="เลือก"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.group_km"
				:options="groupKmOptions"
				name="group_km"
				label="จัดกลุ่ม"
				placeholder="เลือก"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div class="col-12">
			<VLabel label="กรองค่า" />
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.iri1" :precision="2" :required="true" name="iri1" />
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-2 col-2 text-center align-self-md-end align-self-center mb-md-4">
					<span>IRI</span>
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.iri2" :precision="2" :required="true" name="iri2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.aadt1" :precision="2" :required="true" name="aadt1" />
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-2 col-2 text-center align-self-md-end align-self-center mb-md-4">
					<span>AADT</span>
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.aadt2" :precision="2" :required="true" name="aadt2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.ifi1" :precision="2" :required="true" name="ifi1" />
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-2 col-2 text-center align-self-md-end align-self-center mb-md-4">
					<span>IFI</span>
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.ifi2" :precision="2" :required="true" name="ifi2" />
				</div>
			</div>
		</div>
		<div class="col-12 align-self-end text-end mb-md-0 mb-3">
			<!-- <NuxtLink
				class="btn btn-outline btn-outline-primary rounded-4 px-8 py-3 me-3 mt-md-0 mt-2 fw-semibold"
				@click="store.exportDamage(id)"
			>
				Export ความเสียหาย
			</NuxtLink> -->
			<BtnSearch :disabled="store.loading" @click="onSubmit" />
		</div>
	</div>
	<CopySelectRoadTable @data-table="handleDataTable" />
</template>

<style scoped></style>
