<script setup lang="ts">
import { useSurveyRuleCreateStore } from "~/core/modules/setting/surveyRule/store"
import { SurveyRuleCreate } from "~~/core/modules/setting/surveyRule/ui"

useHead({
	title: "เพิ่มเกณฑ์การจำแนกสภาพทาง",
})

const store = useSurveyRuleCreateStore()
useStoreLifecycle(store)

const reRenderkey = ref(0)

const requestReload = () => {
	reRenderkey.value++
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'เกณฑ์การจำแนกสภาพทาง', to: '/settings/survey-rules' }]"
		title="เพิ่มข้อมูล"
	/>
	<SurveyRuleCreate :key="reRenderkey" @on-request-reload="requestReload" />
</template>

<style scoped></style>
