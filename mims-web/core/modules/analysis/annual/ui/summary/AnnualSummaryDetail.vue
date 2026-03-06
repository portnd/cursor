<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAnnualAnalyseDetailStore } from "../../store"
import { ISearchModelReq } from "../../infrastructure"
import { THeader } from "~/core/shared/types/Datatable"
import ServerSideDataTable from "~/core/modules/common/datatable/ui"

const store = useAnnualAnalyseDetailStore()
useStoreLifecycle(store)
const route = useRoute()
const id = Number(route.params.id)
const dataTable = ref()

onMounted(async () => {
	handleReset()

	store.loading = true
	await store.getRoadsTree()
	await store.getStrategicList()
	await store.getDetails(id)
	store.loading = false
})

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 80 },
	{ text: "ปี", value: "year", width: 100 },
	{ text: "แผน", value: "plan", width: 100 },
	{ text: "สายทาง", value: "road_group_name", width: 200 },
	{ text: "ช่วง", value: "road_name", width: 100 },
	{ text: "กม. เริ่มต้น", value: "km_start", width: 100 },
	{ text: "กม. สิ้นสุด", value: "km_end", width: 100 },
	{ text: "ระยะทาง (กม.)", value: "distance", width: 100 },
	{ text: "ช่องจราจร", value: "lane", width: 100 },
	{ text: "วิธีการซ่อมบำรุง", value: "intervention_criteria", width: 250 },
	{ text: "ปริมาณงาน (ตร.ม)", value: "area", width: 100 },
	{ text: "ค่าซ่อมบำรุง (บาท)", value: "cost", width: 100 },
	{ text: "B/C", value: "bc", width: 100 },
	{ text: "AADT (คัน/วัน)", value: "volume_aadt", width: 100 },
	{ text: "IRI เมื่อไม่มีการซ่อมบำรุง", value: "iri_before", width: 100 },
	{ text: "IRI เมื่อมีการซ่อมบำรุง", value: "iri_after", width: 100 },
]

const { handleSubmit, handleReset } = useForm()

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.updateModel(id)

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/analyses")
			},
		})
	}
})

const onSearch = async () => {
	store.search_loading = true

	const params: ISearchModelReq = {
		aadt1: store.data.aadt1,
		aadt2: store.data.aadt2,
		age1: store.data.age1,
		age2: store.data.age2,
		group_km: store.data.group_km,
		ifi1: store.data.ifi1,
		ifi2: store.data.ifi2,
		iri1: store.data.iri1,
		iri2: store.data.iri2,
		lane_type_id: store.data.lane_type_id,
		maintenance_analysis_type_id: store.data.maintenance_analysis_type_id,
		name: store.data.name,
		roads: store.params.road_id.map(Number),
		surface_type_id: store.data.surface_type_id,
		intervention_criteria_id:
			(store.params.intervention_criteria_id ?? []).length > 0 ? store.params.intervention_criteria_id![0] : null,
	}

	await dataTable.value?.searchData(params)
	store.search_loading = false
}

const onReset = async () => {
	store.search_loading = true
	await store.getDetails(id)

	store.params.road_id = []
	store.params.intervention_criteria_id = undefined

	const params: ISearchModelReq = {
		aadt1: store.data.aadt1,
		aadt2: store.data.aadt2,
		age1: store.data.age1,
		age2: store.data.age2,
		group_km: store.data.group_km,
		ifi1: store.data.ifi1,
		ifi2: store.data.ifi2,
		iri1: store.data.iri1,
		iri2: store.data.iri2,
		lane_type_id: store.data.lane_type_id,
		maintenance_analysis_type_id: store.data.maintenance_analysis_type_id,
		name: store.data.name,
		roads: store.params.road_id.map(Number),
		surface_type_id: store.data.surface_type_id,
		intervention_criteria_id: null,
	}

	dataTable.value?.searchData(params)
	store.search_loading = false
}

onUnmounted(() => {
	store.$reset()
})
</script>
<template>
	<div class="row">
		<div class="col-12">
			<div class="card p-5 mb-5">
				<div class="row mb-md-5">
					<div class="col-12 mb-2">
						<VTree
							v-model="store.params.road_id"
							label="สายทาง"
							:multiple="true"
							:searchable="true"
							:options="store.getRoadGroupOptions"
							placeholder="เลือก"
							name="road"
							:required="true"
							:limit="0"
							:mode="'LEAF_PRIORITY'"
						/>
					</div>
					<div class="col-12">
						<VLabel label="กรองค่า" />
					</div>
					<div class="col-md-4 col-12 mb-2">
						<div class="row mb-5">
							<div class="col-md-4 col-4 mb-2 mb-md-0">
								<VNumberInput v-model="store.data.iri1" :precision="2" :required="true" name="iri1" />
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
								<VNumberInput v-model="store.data.iri2" :precision="2" :required="true" name="iri2" />
							</div>
						</div>
					</div>
					<div class="col-md-4 col-12 mb-2">
						<div class="row mb-5">
							<div class="col-md-4 col-4 mb-2 mb-md-0">
								<VNumberInput v-model="store.data.aadt1" :precision="2" :required="true" name="aadt1" />
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
								<VNumberInput v-model="store.data.aadt2" :precision="2" :required="true" name="aadt2" />
							</div>
						</div>
					</div>
					<div class="col-md-4 col-12 mb-2">
						<div class="row mb-5">
							<div class="col-md-4 col-4 mb-2 mb-md-0">
								<VNumberInput v-model="store.data.ifi1" :precision="2" :required="true" name="ifi1" />
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
								<VNumberInput v-model="store.data.ifi2" :precision="2" :required="true" name="ifi2" />
							</div>
						</div>
					</div>
					<div class="col-lane mb-2">
						<VSelect
							v-model="store.data.lane_type_id"
							:options="store.getLaneOptions"
							name="lane_type_id"
							label="ช่องจราจร"
							placeholder="เลือก"
							:can-clear="false"
							:can-deselect="false"
							:required="true"
						/>
					</div>
					<div class="col-maintenance mb-2">
						<VTree
							v-model="store.params.intervention_criteria_id"
							:options="store.getInterventionOptions"
							name="maintenance"
							label="วิธีการซ่อมบำรุง"
							placeholder="เลือก"
							:disable-branch-nodes="true"
							:required="true"
							:limit="0"
							:mode="'LEAF_PRIORITY'"
						/>
					</div>
					<div class="col-tool d-flex gap-2 align-items-end my-2">
						<BtnSearch :disabled="store.search_loading" @click="onSearch" />
						<button
							type="button"
							class="btn btn-outline-primary rounded-4 fw-semibold text-gray-700"
							@click="onReset()"
						>
							รีเซ็ต
						</button>
						<div style="flex-grow: 1"></div>
						<NuxtLink class="btn btn-outline btn-outline-primary rounded-4" @click="onSubmit">
							วิเคราะห์การซ่อมบำรุงใหม่
						</NuxtLink>
					</div>
				</div>
				<div class="row mt-2">
					<div class="col-12">
						<ServerSideDataTable
							ref="dataTable"
							:url="`analyze/${id}/model`"
							:loading="store.search_loading"
							:headers="headers"
							@get-data="store.handleDataTable"
						>
							<template #item-no="{ item }">
								<div class="text-center">{{ item.no }}</div>
							</template>
							<template #item-year="{ item }">
								<div class="text-center">{{ item.year + 543 }}</div>
							</template>
							<template #item-plan="{ item }">
								<div class="text-center">{{ item.plan }}</div>
							</template>
							<template #item-road_group_name="{ item }">
								<div class="text-center">{{ item.road_group_name }}</div>
							</template>
							<template #item-road_name="{ item }">
								<div class="text-center">{{ item.road_name }}</div>
							</template>
							<template #item-km_start="{ item }">
								<div class="text-center">{{ convertMeterToKm(item.km_start) }}</div>
							</template>
							<template #item-km_end="{ item }">
								<div class="text-center">{{ convertMeterToKm(item.km_end) }}</div>
							</template>
							<template #item-distance="{ item }">
								<div class="text-center">{{ item.distance }}</div>
							</template>
							<template #item-lane="{ item }">
								<div class="text-center">{{ item.lane }}</div>
							</template>
							<template #item-intervention_criteria="{ item }">
								<div class="text-center d-flex">
									<VTree
										v-model="item.intervention_criteria"
										:options="store.getInterventionOptions"
										name="maintenance_id"
										label=""
										placeholder="เลือก"
										:disable-branch-nodes="true"
										:default-expand-level="1"
										:reposition="true"
									/>
								</div>
							</template>
							<template #item-area="{ item }">
								<div class="text-center">{{ toNumber(item.area, 2) }}</div>
							</template>
							<template #item-cost="{ item }">
								<div class="text-center">{{ toNumber(item.cost, 2) }}</div>
							</template>
							<template #item-bc="{ item }">
								<div class="text-center">{{ toNumber(item.bc, 2) }}</div>
							</template>
							<template #item-volume_aadt="{ item }">
								<div class="text-center">{{ toNumber(item.volume_aadt) }}</div>
							</template>
							<template #item-iri_before="{ item }">
								<div class="text-center">{{ toNumber(item.iri_before, 2) }}</div>
							</template>
							<template #item-iri_after="{ item }">
								<div class="text-center">{{ toNumber(item.iri_after, 2) }}</div>
							</template>
						</ServerSideDataTable>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
@media (max-width: 9999px) and (min-width: 1350px) {
	.col-lane {
		width: calc(100% / 3);
	}
	.col-maintenance {
		width: calc(100% / 3);
	}
	.col-tool {
		width: calc(100% / 3);
	}
}
@media (max-width: 1340px) and (min-width: 1000px) {
	.col-lane {
		width: calc(100% / 4);
	}
	.col-maintenance {
		width: calc(100% / 4);
	}
	.col-tool {
		width: calc(100% / 2);
	}
}
</style>
