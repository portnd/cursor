<script setup lang="ts">
import { useStrategicEditStore } from "../../store"
import { EditSelectRoadForm, EditSetConditionForm } from "."

const newStore = useStrategicEditStore()
useStoreLifecycle(newStore)

onUnmounted(() => {
	newStore.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div v-if="newStore.step === 1" class="card p-5 mb-5">
				<EditSelectRoadForm />
			</div>
			<div v-else-if="newStore.step === 2" class="card p-5 mb-5">
				<EditSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
