<script setup lang="ts">
import { StrategicSelectRoadForm, StrategicSetConditionForm } from ".."
import { useStrategicCreateStore } from "../../store"

const store = useStrategicCreateStore()
useStoreLifecycle(store)

// onBeforeMount(async () => {
// 	// await store.checkCopyData()
// 	await store.getRoadTrees()
// 	if (Object.keys(store.copyData).length === 0) {
// 		store.setDefaultParam1()
// 	}
// })

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div v-if="store.step === 1" class="card p-5 mb-5">
				<StrategicSelectRoadForm />
			</div>
			<div v-else-if="store.step === 2" class="card p-5 mb-5">
				<StrategicSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
