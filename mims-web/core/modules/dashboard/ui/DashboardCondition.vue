<script setup lang="ts">
import { useDashboardStore } from "../store"
import { DashboardConditionChart, DashboardConditionTable } from "./index"

defineProps({
	mapCollapsed: {
		type: Boolean,
		default: false,
	},
})

const store = useDashboardStore()
const conditionOptions = [
	{
		label: "IRI",
		value: 1,
	},
	{
		label: "MPD",
		value: 2,
	},
	{
		label: "RUT",
		value: 3,
	},
	{
		label: "IFI",
		value: 4,
	},
	{
		label: "แถบสะท้อนแสง",
		value: 5,
	},
]

const surveyRuleOptions = [
	{ label: "25 เมตร", value: 1 },
	{ label: "100 เมตร", value: 2 },
	{ label: "1 กิโลเมตร", value: 3 },
]

const laneOptions = [
	{ label: "1 (สำรวจ: 10 ม.ค. 2567)", value: 1 },
	{ label: "2 (สำรวจ: 10 ม.ค. 2567)", value: 2 },
	{ label: "3 (สำรวจ: 10 ม.ค. 2567)", value: 3 },
]
</script>
<template>
	<div class="row">
		<div class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12">
			<VSelect
				v-model="store.conditionType"
				:options="conditionOptions"
				label="ประเภทสภาพทาง"
				name="select"
				placeholder="ทั้งหมด"
				:close-on-select="true"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="(e: any) => console.log(e)"
			/>
		</div>
		<div class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12">
			<VSelect
				v-model="store.surveyRule"
				:options="surveyRuleOptions"
				label="เลือกเกณฑ์การจำแนกสภาพทาง"
				name="select"
				placeholder="ทั้งหมด"
				:close-on-select="true"
				:can-clear="false"
				:can-deselect="false"
			/>
		</div>
		<div v-if="store.roads.length === 1" class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12">
			<VSelect
				v-model="store.lane"
				:options="laneOptions"
				label="ช่องจราจร"
				name="select"
				placeholder="ทั้งหมด"
				:close-on-select="true"
				:can-clear="false"
			/>
		</div>
	</div>
	<div class="col-12 pt-5">
		<DashboardConditionChart :collapsed="mapCollapsed" />
		<DashboardConditionTable />
	</div>
</template>

<style lang="scss" scoped></style>
