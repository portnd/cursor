<script setup lang="ts">
import { useStrategicEditStore } from "../../store"
import { CopySelectRoadForm, CopySetConditionForm } from "."

const store = useStrategicEditStore()
useStoreLifecycle(store)

const route = useRoute()
const id = Number(route.params.id)

onBeforeMount(async () => {
	await store.getRoadsOptions()
	await store.getDefaultData(id)
})

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div v-if="store.step === 1" class="card p-5 mb-5">
				<CopySelectRoadForm />
			</div>
			<div v-else-if="store.step === 2" class="card p-5 mb-5">
				<CopySetConditionForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
