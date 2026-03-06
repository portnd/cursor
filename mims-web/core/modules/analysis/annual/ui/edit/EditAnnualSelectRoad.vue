<script setup lang="ts">
import { useAnnualAnalyseEditStore } from "../../store"
import EditAnnualSelectRoadForm from "./selectRoad/EditAnnualSelectRoadForm.vue"
import EditAnnualSetConditionForm from "./setCondition/EditAnnualSetConditionForm.vue"

// const route = useRoute()
// const id = Number(route.params.id)
const store = useAnnualAnalyseEditStore()
useStoreLifecycle(store)

onMounted(() => {
	// store.getDefaultData(id)
})

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<!-- <AnnualStep /> -->
			<div v-if="store.step === 1" class="card p-5 mb-5">
				<EditAnnualSelectRoadForm />
			</div>
			<div v-else-if="store.step === 2" class="card p-5 mb-5">
				<EditAnnualSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
