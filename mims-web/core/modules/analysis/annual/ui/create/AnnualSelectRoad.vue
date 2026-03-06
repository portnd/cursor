<script setup lang="ts">
import { useAnnualCreateStore } from "../../store"
import AnnualSelectRoadForm from "./selectRoad/AnnualSelectRoadForm.vue"
import AnnualSetConditionForm from "./setCondition/AnnualSetConditionForm.vue"

const store = useAnnualCreateStore()
useStoreLifecycle(store)

// onMounted(async () => {
// 	await store.getRoadTrees()
// 	if (Object.keys(store.copyData).length === 0) {
// 		store.setDefaultParam1()
// 	} else {
// 		store.setCopyParams()
// 	}
// })

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<!-- <AnnualStep /> -->
			<div v-if="store.step === 1" class="card p-5 mb-5">
				<AnnualSelectRoadForm />
			</div>
			<div v-else-if="store.step === 2" class="card p-5 mb-5">
				<AnnualSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
