<script setup lang="ts">
// import { useReportStore } from "../store/ReportStore"
import { ReportSelect, ReportDownload } from "./index"

// const store = useReportStore()

const data = ref()

// onMounted(async () => {
// await store.getRoadTree()
// })

// onUnmounted(() => {
// store.$dispose()
// store.$reset()
// })
</script>

<template>
	<div class="row">
		<div class="col-md-7 col-12">
			<ReportSelect @select="(e) => (data = e)" />
		</div>
		<div class="col-md-5 col-12">
			<ReportDownload :access-key="data" />
		</div>
	</div>
</template>

<style scoped></style>
