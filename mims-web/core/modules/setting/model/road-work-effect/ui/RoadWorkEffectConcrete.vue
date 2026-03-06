<script setup lang="ts">
import RoadWorkEffectTab from "./RoadWorkEffectTab.vue"
import TableConcrete from "./table/TableConcrete.vue"
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
					<RoadWorkEffectTab />
					<TableConcrete />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
