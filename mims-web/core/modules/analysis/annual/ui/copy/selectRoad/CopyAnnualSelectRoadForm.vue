<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAnnualAnalyseCopyStore } from "../../../store"
import { CopyAnnualStep } from ".."
import { IValidate } from "~/core/shared/types/Validate"

const route = useRoute()
const id = Number(route.params.id)
const store = useAnnualAnalyseCopyStore()

onBeforeMount(() => {
	store.checkFilterParams()
})

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

	const keys = Object.keys(store.params1)
	keys.forEach((key) => {
		switch (key) {
			case "iri1":
				validation[key] = ""
				break
			case "iri2":
				validation[key] = ""
				break
			case "aadt1":
				validation[key] = ""
				break
			case "aadt2":
				validation[key] = ""
				break
			case "gn1":
				validation[key] = ""
				break
			case "gn2":
				validation[key] = ""
				break
			default:
				validation[key] = `required`
				break
		}
	})

	return validation
})

const { handleSubmit } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit((_, actions) => {
	useAction(actions)
	store.updatePrepareData(id)
})
</script>

<template>
	<div class="row">
		<CopyAnnualStep />
		<div class="col-12 mb-2">
			<VTree
				v-model="store.params1.road_group_id"
				label="สายทาง"
				:multiple="true"
				:options="store.getRoadTreeOptions"
				name="road_group_id"
				:limit="1"
				:mode="'LEAF_PRIORITY'"
				:required="true"
			/>
		</div>
	</div>
	<div class="row">
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params1.surface_type_id"
				:options="surfaceOptions"
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
				v-model="store.params1.lane_type_id"
				:options="laneOptions"
				name="lane_type_id"
				placeholder="เลือก"
				label="ช่องจราจร"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params1.group_km"
				:options="groupKmOptions"
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
					<VNumberInput v-model="store.params1.iri1" :precision="2" name="iri1" label="" />
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
					<VNumberInput v-model="store.params1.iri2" :precision="2" name="iri2" label="" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.params1.aadt1" :precision="2" name="aadt1" />
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
					<VNumberInput v-model="store.params1.aadt2" :precision="2" name="aadt2" />
				</div>
			</div>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<div class="row mb-5">
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.params1.gn1" :precision="2" name="gn1" />
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-2 col-2 text-center align-self-md-end align-self-center mb-md-4">
					<span>GN</span>
				</div>
				<div class="col-md-1 col-1 text-center align-self-md-end align-self-center mb-md-4">
					<span>&lt;</span>
				</div>
				<div class="col-md-4 col-4 mb-2 mb-md-0">
					<VNumberInput v-model="store.params1.gn2" :precision="2" name="gn2" />
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
</template>

<style scoped></style>
