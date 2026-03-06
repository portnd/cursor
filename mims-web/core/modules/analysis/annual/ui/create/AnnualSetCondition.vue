<script setup lang="ts">
import AnnualStep from "../AnnualStep.vue"
import AnnualSetConditionForm from "./setCondition/AnnualSetConditionForm.vue"
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<AnnualStep />
			<div class="card p-5 mb-5">
				<AnnualSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
