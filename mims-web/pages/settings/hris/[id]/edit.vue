<script setup lang="ts">
import { useSurveyRuleEditStore } from "~/core/modules/setting/surveyRule/store/SurveyRuleEditStore"
import { SurveyRuleEdit } from "~~/core/modules/setting/surveyRule/ui"

useHead({
	title: "แก้ไขเกณฑ์การจำแนกสภาพทาง",
})

const store = useSurveyRuleEditStore()
useStoreLifecycle(store)

const reRenderkey = ref(0)

const requestReload = () => {
	reRenderkey.value++
	console.log(reRenderkey.value)
}

onUnmounted(() => {
	console.log("onUnmounted")
	store.$reset()
})
</script>

<template>
	<TheBreadcrumb
		:breadcrumbs="[{ name: 'ตั้งค่าระบบ' }, { name: 'เกณฑ์การจำแนกสภาพทาง', to: '/settings/survey-rules' }]"
		title="แก้ไขข้อมูล"
	/>
	<SurveyRuleEdit :key="reRenderkey" @on-request-reload="requestReload" />
</template>

<style scoped></style>
