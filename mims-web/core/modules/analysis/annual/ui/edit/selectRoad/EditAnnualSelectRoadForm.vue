<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAnnualAnalyseEditStore } from "../../../store"
import EditAnnualStep from "../EditAnnualStep.vue"
import EditAnnualSelectRoadTable from "./EditAnnualSelectRoadTable.vue"
import { IValidate } from "~/core/shared/types/Validate"

interface IPrepareDataTable {
	loadData: Function
	searchData: Function
}

const route = useRoute()
const id = Number(route.params.id)
const store = useAnnualAnalyseEditStore()
const dataTable = ref()

onMounted(async () => {
	if (store.wasStep2 === false) {
		handleReset()
		store.loading = true

		store.prepare_data_id = id
		await store.getRoadsOptions()
		await store.getRefMaintenanceOptions()
		await store.getDefaultData(id)
		await store.getSelectedId()
		await store.getAllId()

		store.loading = false
	}
	await dataTable.value?.searchData()
})

const validate = computed(() => {
	const validation: IValidate = {}
	validation.road_id = "required"
	validation.surface_type_id = "required"
	validation.lane_type_id = "required"
	validation.group_km = "required"

	return validation
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const handleDataTable = (table: IPrepareDataTable) => {
	dataTable.value = table
}

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.updatePrepareData()

	if (res?.status) {
		store.default.id = res.data.id
		console.log(res.data.id)

		await store.checkPrepareDataStatus()
		await store.getAllId()
		await dataTable.value?.searchData()

		store.wasStep2 = false
	}
})
</script>

<template>
	<div class="row">
		<EditAnnualStep />
		<div class="col-md-6 col-12 mb-2">
			<VTextInput v-model="store.default.name" label="ชื่อรายการ" name="name" :required="true" />
		</div>
		<div class="col-md-6 col-12 mb-2">
			<VTree
				v-model="store.road_id"
				label="สายทาง"
				:multiple="true"
				:options="store.roadGroupOptions"
				name="road_id"
				:limit="0"
				:mode="'LEAF_PRIORITY'"
				:required="true"
			/>
		</div>
	</div>
	<div class="row">
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.surface_type_id"
				:options="store.surfaceOptions"
				name="surface_type_id"
				placeholder="เลือก"
				label="ชนิดผิวทาง"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.lane_type_id"
				:options="store.getLaneOptions"
				name="lane_type_id"
				placeholder="เลือก"
				label="ช่องจราจร"
				:can-clear="false"
				:can-deselect="false"
				:disabled="store.road_id.length === 0"
				:required="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.default.group_km"
				:options="store.groupKmOptions"
				name="group_km"
				placeholder="เลือก"
				label="จัดกลุ่ม"
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
					<VNumberInput v-model="store.default.iri1" :precision="2" name="iri1" label="" />
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
					<VNumberInput v-model="store.default.iri2" :precision="2" name="iri2" label="" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.aadt1" :precision="2" name="aadt1" />
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
					<VNumberInput v-model="store.default.aadt2" :precision="2" name="aadt2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.default.ifi1" :precision="2" name="ifi1" />
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
					<VNumberInput v-model="store.default.ifi2" :precision="2" name="ifi2" />
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
	<EditAnnualSelectRoadTable @data-table="handleDataTable" />
</template>

<style scoped></style>
