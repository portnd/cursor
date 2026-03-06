<script setup lang="ts">
import { useForm } from "vee-validate"
import { useStrategicCreateStore } from "../../../store"
import StrategicStep from "../StrategicStep.vue"
import { StrategicSelectRoadTable } from "../.."
import { IValidate } from "~/core/shared/types/Validate"

interface IPrepareDataTable {
	loadData: Function
	searchData: Function
}

const store = useStrategicCreateStore()
const dataTable = ref<IPrepareDataTable>()

onMounted(async () => {
	resetField("road_id", { errors: undefined, value: null })
	await store.getRoadsOptions()
	await store.getStrategicList()
})

const handleDataTable = (table: IPrepareDataTable) => {
	dataTable.value = table
}

const validate = computed(() => {
	const validation: IValidate = {}

	const keys = Object.keys(store.params1)
	keys.forEach((key) => {
		validation[key] = ["iri1", "iri2", "aadt1", "aadt2", "ifi1", "ifi2"].includes(key) ? "" : "required"
	})

	return validation
})

const { handleSubmit, resetField } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	store.prepare_data_loading = true

	await store.createPrepareData()
	await store.checkPrepareDataStatus()
	await store.getPrepareDataId()
	await dataTable.value?.searchData()

	store.wasStep2 = false

	store.prepare_data_loading = false
})

onMounted(() => {
	if (store.wasStep2 === true) {
		dataTable.value?.searchData()
	}
})
</script>

<template>
	<!-- ทำเป็นการ์ดเดียว -->
	<div class="row">
		<StrategicStep />
		<div class="col-md-6 col-12 mb-2">
			<VTextInput v-model="store.params1.name" label="ชื่อรายการ" name="name" :required="true" />
		</div>
		<div class="col-md-6 col-12 mb-2">
			<VTree
				v-model="store.params1.road_id"
				label="สายทาง"
				placeholder="เลือก"
				name="road_id"
				:options="store.roadGroupOption"
				:mode="'LEAF_PRIORITY'"
				:multiple="true"
				:limit="0"
				:required="true"
				@update:model-value="resetField('lane_type_id', { errors: undefined, value: null })"
			/>
		</div>
	</div>
	<div class="row">
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params1.surface_type_id"
				:options="store.surfaceOptions"
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
				v-model="store.params1.lane_type_id"
				:options="store.getLaneOptions"
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
				v-model="store.params1.group_km"
				:options="store.groupKmOptions"
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
					<VNumberInput v-model="store.params1.iri1" :precision="2" :required="true" name="iri1" />
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
					<VNumberInput v-model="store.params1.iri2" :precision="2" :required="true" name="iri2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.params1.aadt1" :precision="2" :required="true" name="aadt1" />
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
					<VNumberInput v-model="store.params1.aadt2" :precision="2" :required="true" name="aadt2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.params1.ifi1" :precision="2" :required="true" name="ifi1" />
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
					<VNumberInput v-model="store.params1.ifi2" :precision="2" :required="true" name="ifi2" />
				</div>
			</div>
		</div>
		<div class="col-12 align-self-end text-end mb-md-0 mb-3">
			<!-- <NuxtLink
				class="btn btn-outline btn-outline-primary rounded-4 px-8 py-3 me-3 mt-md-0 mt-2 fw-semibold"
				@click="store.exportDamage(store.prepareData.id)"
			>
				Export ความเสียหาย
			</NuxtLink> -->
			<BtnSearch :disabled="store.prepare_data_loading" @click="onSubmit" />
		</div>
	</div>
	<StrategicSelectRoadTable @data-table="handleDataTable" />
</template>

<style scoped></style>
