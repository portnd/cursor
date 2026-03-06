<script setup lang="ts">
import { AnnualSummaryList } from "~/core/modules/analysis/annual/ui"

useHead({
	title: "สรุปรวม",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: 'วิเคราะห์การซ่อมบำรุง', to: '/analyses' }]" title="บำรุงรักษาประจำปี" />
	<AnnualSummaryList />
</template>

<style scoped></style>
