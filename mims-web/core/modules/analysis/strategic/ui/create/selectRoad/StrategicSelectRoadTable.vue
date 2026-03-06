<script setup lang="ts">
import { useStrategicCreateStore } from "../../../store"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

const store = useStrategicCreateStore()

const dataTable = ref()
const emit = defineEmits(["dataTable"])

onMounted(() => {
	emit("dataTable", dataTable.value)
})

const submitNext = async () => {
	if (store.selected_id.length === 0) {
		showAlert({ type: "warning", title: "แจ้งเตือน", message: "โปรดเลือกช่วงที่ต้องการวิเคราะห์" })
	} else {
		if (store.wasStep2 === false) {
			await store.submitToStep2()
		}
		store.step = 2
	}
}

const cancel = () => {
	const prevUrl = "/analyses"

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
				:loading-top="true"
				:items="items"
				items-selected
				:all-ids="store.all_ids"
				:selected="store.selected_id"
				:loading="store.prepare_data_loading"
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
					<div class="text-center">{{ (item.length / 1000).toFixed(3) }}</div>
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
			<BtnSubmit label="ถัดไป" :disabled="store.loading" :loading="store.loading" @click="submitNext" />
		</div>
	</div>
</template>

<style scoped></style>
