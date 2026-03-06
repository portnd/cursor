<script setup lang="ts">
import { useTrafficParameterRoadGroupStore, useTrafficParameterStore } from "../store"
import TrafficTab from "./TrafficTab.vue"
import { TableParameter } from "./index"
import { ModelTab } from "~/core/modules/setting/model/tab/ui"

const trafficRoadGroupStore = useTrafficParameterRoadGroupStore()
const trafficParamterStore = useTrafficParameterStore()
useStoreLifecycle([trafficRoadGroupStore, trafficParamterStore])
onBeforeMount(async () => {
	await trafficRoadGroupStore.getAadtRoadGroupData()
	trafficParamterStore.roadGroupID = trafficRoadGroupStore.data[0].road_group_id
})

</script>

<template>
	<div class="row">
		<div class="col-12 mb-5">
			<div class="card shadow">
				<div class="card-body p-5">
					<ModelTab />
				</div>
			</div>
		</div>
		<div class="col-12 mb-5">
			<div class="card shadow">
				<div class="card-body p-5">
					<TrafficTab />
					<TableParameter />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
