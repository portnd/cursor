<script setup lang="ts">
import { useRoadUserCostRucStore } from "../../store/RoadUserCostRucStore"
import TableDefaultData from "./TableDefaultData.vue"
import TableDriving from "./TableDriving.vue"
import TableEngineSpeed from "./TableEngineSpeed.vue"
import TableFuelConsumption from "./TableFuelConsumption.vue"
import TableMaintenance from "./TableMaintenance.vue"
import TableTrafficData from "./TableTrafficData.vue"
import TableTravelTime from "./TableTravelTime.vue"
import TableVehicleSpeedCalc from "./TableVehicleSpeedCalc.vue"
import TableWastedLubricant from "./TableWastedLubricant.vue"
import TableWastedTires from "./TableWastedTires.vue"
const route = useRoute()
const tab = route.path.split("/")[5]
const rucStore = useRoadUserCostRucStore()

onMounted(() => {
	switch (tab) {
		case "default-data":
			rucStore.rucListId = 1
			break
		case "driving":
			rucStore.rucListId = 2
			break
		case "engine-speed":
			rucStore.rucListId = 3
			break
		case "fuel-consumption":
			rucStore.rucListId = 4
			break
		case "wasted-lubricant":
			rucStore.rucListId = 5
			break
		case "wasted-tires":
			rucStore.rucListId = 6
			break
		case "maintenance":
			rucStore.rucListId = 7
			break
		case "travel-time":
			rucStore.rucListId = 8
			break
		case "vehicle-speed-calc":
			rucStore.rucListId = 9
			break
		case "traffic-data":
			rucStore.rucListId = 10
			break
	}
	rucStore.getRucData(rucStore.rucListId)
})

const selectRucData = (id: number) => {
	rucStore.switchPath(id)
}
</script>

<template>
	<div class="row">
		<div class="col-md-4 col-12 mb-5">
			<VSelect
				v-model="rucStore.rucListId"
				:options="rucStore.getRucListOptions"
				:can-clear="false"
				name="rucListId"
				:can-deselect="false"
				:required="true"
				label="ตัวแปร"
				@update:model-value="(e: number) => selectRucData(e)"
			/>
		</div>
		<template v-if="rucStore.rucListId === 1">
			<TableDefaultData />
		</template>
		<template v-else-if="rucStore.rucListId === 2">
			<TableDriving />
		</template>
		<template v-else-if="rucStore.rucListId === 3">
			<TableEngineSpeed />
		</template>
		<template v-else-if="rucStore.rucListId === 4">
			<TableFuelConsumption />
		</template>
		<template v-else-if="rucStore.rucListId === 5">
			<TableWastedLubricant />
		</template>
		<template v-else-if="rucStore.rucListId === 6">
			<TableWastedTires />
		</template>
		<template v-else-if="rucStore.rucListId === 7">
			<TableMaintenance />
		</template>
		<template v-else-if="rucStore.rucListId === 8">
			<TableTravelTime />
		</template>
		<template v-else-if="rucStore.rucListId === 9">
			<TableVehicleSpeedCalc />
		</template>
		<template v-else-if="rucStore.rucListId === 10">
			<TableTrafficData />
		</template>
	</div>
</template>

<style scoped></style>
