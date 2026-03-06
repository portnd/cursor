<script setup lang="ts">
import { useAnnualAnalyseCopyStore } from "../../store"
import { CopyAnnualSelectRoadForm, CopyAnnualSelectRoadTable, CopyAnnualSetConditionForm } from "."

const route = useRoute()
const id = Number(route.params.id)
const store = useAnnualAnalyseCopyStore()
useStoreLifecycle(store)

onMounted(() => {
	store.getDefaultData(id)
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
				<CopyAnnualSelectRoadForm />
				<CopyAnnualSelectRoadTable />
			</div>
			<div v-else-if="store.step === 2" class="card p-5 mb-5">
				<CopyAnnualSetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
