<script setup lang="ts">
import DeteriorationTab from "./DeteriorationTab.vue"
import { TableAC } from "./index"
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
					<DeteriorationTab />
					<TableAC />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
