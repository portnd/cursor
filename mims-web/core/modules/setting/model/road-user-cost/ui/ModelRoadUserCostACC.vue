<script setup lang="ts">
import { RoadUserCostTab, TableACCList } from "./index"
import { ModelTab } from "~/core/modules/setting/model/tab/ui"

defineProps({
	type: {
		type: null,
		required: true,
	},
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
					<RoadUserCostTab :type="type" />
					<TableACCList />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
