<script setup lang="ts">
import { useAnnualAnalyseCopyStore } from "../../../store"
import type { THeader } from "~~/core/shared/types/Datatable"

const store = useAnnualAnalyseCopyStore()

const submitNext = async () => {
	if (store.selectedPrepareData.length === 0) {
		useHandlerError(0, { message: "โปรดเลือกช่วงที่ต้องการวิเคราะห์" }, { showAlert: true })
	} else {
		await store.createAnalyzeStep2()
		store.step = 2
	}
}

// watch(
// 	() => store.selectedPrepareData,
// 	(newPrepareData) => {
// 		if (newPrepareData.length > 0) {
// 			store.selectedPrepareId = newPrepareData.filter((item) => item.is_selected).map((item) => item.id)
// 		}
// 	}
// )

const cancel = () => {
	navigateTo("/analyses")
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
	{ text: "GN", value: "gn", width: 100 },
]
</script>

<template>
	<div class="row">
		<div class="col-12">
			<VDatatable
				:rows-per-page="50"
				:headers="headers"
				:items="store.prepareData.prepare_data"
				items-selected
				:loading="store.loading"
				:selected="store.selectedPrepareId"
				@selected="store.handlePreparedData"
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
				<template #item-gn="{ item }">
					<div class="text-center">{{ toNumber(item.gn, 2) }}</div>
				</template>
			</VDatatable>
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
