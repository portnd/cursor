<script setup lang="ts">
import EditAnnualStep from "./EditAnnualStep.vue"
import EditAnnualSetConditionForm from "./setCondition/EditAnnualSetConditionForm.vue"
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<EditAnnualStep />
			<div class="card p-5 mb-5">
				<EditAnnualSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
