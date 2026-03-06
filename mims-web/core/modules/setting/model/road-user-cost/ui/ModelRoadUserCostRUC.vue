<script setup lang="ts">
import { useRoadUserCostRucStore } from "../store/RoadUserCostRucStore"
import { TableRUCList, RoadUserCostTab } from "./index"
import { ModelTab } from "~/core/modules/setting/model/tab/ui"

defineProps({
	type: {
		type: null,
		required: true,
	},
})

const rucStore = useRoadUserCostRucStore()
useStoreLifecycle(rucStore)
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
					<RoadUserCostTab :type="type" />
					<TableRUCList />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
