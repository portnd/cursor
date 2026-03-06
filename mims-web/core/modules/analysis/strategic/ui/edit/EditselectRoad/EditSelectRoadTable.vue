<script setup lang="ts">
import { useStrategicEditStore } from "../../../store"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~/core/modules/common/datatable/ui"

const store = useStrategicEditStore()

const route = useRoute()
const dataTable = ref()

const emit = defineEmits(["dataTable"])

onMounted(() => {
	store.prepare_data_id = Number(useRoute().params.id)
	emit("dataTable", dataTable.value)
})

const submitNext = async () => {
	if (store.selected_id.length === 0) {
		showAlert({ type: "warning", title: "แจ้งเตือน", message: "โปรดเลือกช่วงที่ต้องการวิเคราะห์" })
	} else {
		if (store.wasStep2 === false) {
			await store.nextStep()
		}
		store.step = 2
	}
}

const cancel = () => {
	const path = store.isCopy
		? "/analyses"
		: `/analyses/${store.getPath?.group}/summary/${route.params.id}/${store.getPath?.criteria}`
	const prevUrl = path
	navigateTo(prevUrl)
}

const headers: THeader[] = [
	{ text: "สายทาง", value: "group_name", width: 200 },
	{ text: "ช่วง", value: "road_name", width: 100 },
	{ text: "กม. เริ่มต้น", value: "km_start", width: 100 },
	{ text: "กม. สิ้นสุด", value: "km_end", width: 100 },
	{ text: "ระยะทาง (กม.)", value: "length", width: 100 },
	{ text: "ช่องจราจร", value: "lane_no", width: 100 },
	{ text: "IRI", value: "iri", width: 100 },
	{ text: "AADT", value: "aadt", width: 100 },
	{ text: "IFI", value: "ifi", width: 100 },
]
</script>

<template>
	<div class="row">
		<div class="col-12">
			<ServerSideDataTable
				ref="dataTable"
				:url="`/analyze/${store.prepare_data_id}/prepare_data`"
				:rows-per-page="50"
				:headers="headers"
				:is-init="false"
				items-selected
				:loading="store.prepare_loading"
				:all-ids="store.allPrepareId"
				:selected="store.selected_id"
				@selected="store.handleSelectedPrepareData"
			>
				<template #item-group_name="{ item }">
					<div class="text-center">{{ item.group_name }}</div>
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
				<template #item-length="{ item }">
					<div class="text-center">{{ (item.length / 1000).toFixed(2) }}</div>
				</template>
				<template #item-lane_no="{ item }">
					<div class="text-center">{{ item.lane_no }}</div>
				</template>
				<template #item-iri="{ item }">
					<div class="text-center">{{ toNumber(item.iri, 2) }}</div>
				</template>
				<template #item-aadt="{ item }">
					<div class="text-center">{{ toNumber(item.aadt, 2) }}</div>
				</template>
				<template #item-ifi="{ item }">
					<div class="text-center">{{ toNumber(item.ifi, 2) }}</div>
				</template>
			</ServerSideDataTable>
		</div>
	</div>
	<div class="row mt-5">
		<div class="col-12 text-end">
			<BtnCancel @click="cancel" />
			<BtnSubmit label="ถัดไป" @click="submitNext" />
		</div>
	</div>
</template>

<style scoped></style>
