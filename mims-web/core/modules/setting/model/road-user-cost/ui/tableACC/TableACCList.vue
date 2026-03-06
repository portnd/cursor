<script setup lang="ts">
import { useRoadUserCostAccStore } from "../../store/RoadUserCostAccStore"
import TableLossValue from "./TableLossValue.vue"
import TableAccident from "./TableAccident.vue"
import { useRoadListStore } from "~/core/modules/road/roadList/store"
import { IRoadData } from "~/core/modules/road/roadList/infrastructure"
import { IOption } from "~/core/shared/types/Option"
const route = useRoute()
const tab = route.path.split("/")[5]
const roadListStore = useRoadListStore()

const store = useRoadUserCostAccStore()
useStoreLifecycle([store, roadListStore])

onMounted(async () => {
	await roadListStore.getData()

	store.roadGroupID = roadListStore.roads[0].id
	if (tab === "loss-value") {
		store.accId = 1
		await store.getLossAccidentData()
	} else {
		store.accId = 2
		await store.getAccChanceData(store.roadGroupID)
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

const switchPage = (id: number) => {
	switch (id) {
		case 1:
			navigateTo("/settings/models/road-user-cost/acc/loss-value")
			break
		case 2:
			navigateTo("/settings/models/road-user-cost/acc/chance-accident")
			break
	}
}

</script>

<template>
	<div class="row">
		<div class="col-md-4 col-12 mb-5">
			<VSelect
				v-model="store.accId"
				:options="store.getAccListOptions"
				:can-clear="false"
				:can-deselect="false"
				:required="true"
				label="ตัวแปร"
				name=""
				placeholder="เลือกตัวแปร"
				@update:model-value="(e: number) => switchPage(e)"
			/>
		</div>
		<template v-if="tab !== 'loss-value'">
			<div class="col-md-4 col-12 mb-5">
				<VSelect
					v-model="store.roadGroupID"
					name="road_group_id"
					:can-clear="false"
					:can-deselect="false"
					:searchable="true"
					:options="generateOptionTable(roadListStore.roads)"
					label="สายทาง"
					placeholder="เลือกสายทาง"
					:required="true"
					@update:model-value="(e: number) => store.getAccChanceData(e)"
				/>
			</div>
		</template>
		<template v-if="tab === 'loss-value'">
			<TableLossValue />
		</template>
		<template v-else>
			<TableAccident />
		</template>
	</div>
</template>

<style scoped></style>
