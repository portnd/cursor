<script setup lang="ts">
import { StrategicSummaryList } from "~/core/modules/analysis/strategic/ui"

useHead({
	title: "สรุปรวม",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' }]" title="บำรุงรักษาเชิงกลยุทธ์" />
	<StrategicSummaryList />
</template>

<style scoped></style>
