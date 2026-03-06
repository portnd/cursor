<script setup lang="ts">
import TrafficTab from "./TrafficTab.vue"
import { TableVehicleType } from "./index"
import { ModelTab } from "~/core/modules/setting/model/tab/ui"
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
					<TableVehicleType />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
