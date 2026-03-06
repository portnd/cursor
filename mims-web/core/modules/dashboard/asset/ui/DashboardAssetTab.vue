<script setup lang="ts">
const asset = ref("IN")
const emit = defineEmits(["asset"])

watch(asset, () => {
	emit("asset", asset.value)
})
onMounted(() => {
	emit("asset", asset.value)
})
</script>

<template>
	<div class="card p-5 mb-5">
		<nav>
			<div id="nav-tab" class="nav nav-tabs" role="tablist">
				<button
					id="assetIn-tab"
					class="nav-link col fs-5 fw-semibold rounded-3 py-3 me-2 active"
					data-bs-toggle="tab"
					type="button"
					role="tab"
					aria-selected="true"
					@click="asset = 'IN'"
				>
					สินทรัพย์
				</button>
				<!-- <button
					id="assetOut-tab"
					class="nav-link col fs-5 fw-semibold rounded-3 py-3"
					data-bs-toggle="tab"
					type="button"
					role="tab"
					aria-selected="false"
					@click="asset = 'OUT'"
				>
					สินทรัพย์นอกเขตทาง
				</button> -->
			</div>
		</nav>
	</div>
</template>

<style scoped>
#nav-tab .active {
	background-color: #fff0d9 !important;
	color: #fdb833 !important;
}

#nav-tab .nav-link {
	background: var(--kt-gray-100);
}
</style>
