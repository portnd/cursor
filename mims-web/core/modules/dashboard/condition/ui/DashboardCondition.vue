<script setup lang="ts">
import { useDashboardStore } from "../../store"
import { useDashboardConditionStore, useDashboardRoadConditionStore, useDashboardReflectiveStore } from "../store"
import {
	DashboardConditionChart,
	DashboardConditionTable,
	DashboardReflectiveChart,
	DashboardReflectiveTable,
	DashboardRoadConditionChart,
	DashboardRoadConditionTable,
} from "./index"

const props = defineProps({
	mapCollapsed: {
		type: Boolean,
		default: false,
	},
	searchCall: {
		type: Boolean,
		default: false,
	},
})

const store = useDashboardConditionStore()
const storeRoadCondition = useDashboardRoadConditionStore()
const storeDashboardReflectiveStore = useDashboardReflectiveStore()
const dashboardStore = useDashboardStore()


enum EConditionType {
	IRI = 1,
	MPD = 2,
	RUT = 3,
	IFI = 4,
	Reflect = 5,
}

const conditionOptions = [
	{
		label: "IRI",
		value: EConditionType.IRI,
	},
	{
		label: "MPD",
		value: EConditionType.MPD,
	},
	{
		label: "RUT",
		value: EConditionType.RUT,
	},
	{
		label: "IFI",
		value: EConditionType.IFI,
	},
	{
		label: "แถบสะท้อนแสง",
		value: EConditionType.Reflect,
	},
]

const onUpdateConditionType = async () => {
	const result = conditionOptions.find((option) => option.value === store.conditionType)

	store.conditionTypeString = result?.label as string
	storeRoadCondition.params.condition_type = result?.label as string
	if (dashboardStore.params.road_id.length === 1) {
		if (store.conditionType !== 5) {
			// set default owner_id and related params
			storeRoadCondition.setDaultConditionType()
			storeRoadCondition.getConditionType()

			// setDaultConditionType() always resets owner_id to conditionGrade()[0].id which may be
			// IRI-only. If checkConditionType() couldn't correct it (e.g. refTypeRangeIds not yet
			// seeded for this type), the wrong owner reaches setConditionMapParams() as a truthy
			// value and the fallback ownerId is never used → wrong colors on map.
			// Safe fix: apply the same getSurveyRangeOptions guard used for multiple roads.
			const surveyOptions = storeRoadCondition.getSurveyRangeOptions
			if (
				surveyOptions.length > 0 &&
				!surveyOptions.some((o) => Number(o.value) === storeRoadCondition.params.owner_id)
			) {
				storeRoadCondition.params.owner_id = Number(surveyOptions[0].value)
			}
		} else {
			// set default owner_id
			await storeDashboardReflectiveStore.setDaultConditionType()
		}
		// Refresh chart and map for single road so map shows correct condition type (e.g. MPD colors)
		store.getCondition()
		await dashboardStore.getConditionMap()
	} else {
		if (store.conditionType !== 5) {
			// Use the first owner that actually supports the current condition type.
			// conditionGrade()[0] may only have IRI rules — using it for IFI/MPD/RUT
			// causes the API to apply IRI thresholds on IFI values (e.g., IFI 0.3 < 2.75 → green)
			// which produces wrong colors on the map.
			const surveyOptions = storeRoadCondition.getSurveyRangeOptions
			storeRoadCondition.params.owner_id =
				surveyOptions.length > 0
					? Number(surveyOptions[0].value)
					: useInitData()?.conditionGrade()[0]?.id
		} else {
			storeDashboardReflectiveStore.params.owner_id = useInitData()?.reflectivityGrade()
				? useInitData()?.reflectivityGrade()?.[0]?.id ?? null
				: null
		}

		store.getCondition()
		await dashboardStore.getConditionMap()
	}
	// dashboardStore.searchCall = !dashboardStore.searchCall
}

const onConditionChange = async () => {
	store.loading = true
	await onUpdateConditionType()

	store.loading = false
}

const onConditionCriteriaChange = async () => {
	store.loading = true
	storeRoadCondition.onUpdateOwner()
	if (dashboardStore.params.road_id.length !== 1) {
		store.getCondition()
		await dashboardStore.getConditionMap()
	}

	store.loading = false
}
const onReflectCriteriaChange = async () => {
	store.loading = true

	storeDashboardReflectiveStore.onUpdateOwner()
	if (dashboardStore.params.road_id.length !== 1) {
		store.getCondition()
		await dashboardStore.getConditionMap()
	}

	store.loading = false
}
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
				:disabled="store.loading"
				on-update-condition-type
				@update:model-value="onConditionChange"
			/>
		</div>
		<div
			v-if="store.conditionType !== EConditionType.Reflect"
			class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12"
		>
			<VSelect
				v-model="storeRoadCondition.params.owner_id"
				:options="storeRoadCondition.getSurveyRangeOptions"
				:can-clear="false"
				:can-deselect="false"
				:disabled="store.loading"
				label="เลือกเกณฑ์การจำแนกสภาพทาง"
				name="ref_condition_range_id"
				placeholder="เลือก"
				@update:model-value="onConditionCriteriaChange"
			/>
		</div>
		<div
			v-else-if="store.conditionType === EConditionType.Reflect"
			class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12"
		>
			<VSelect
				v-model="storeDashboardReflectiveStore.params.owner_id"
				:options="storeDashboardReflectiveStore.getOwnerOptions"
				:can-clear="false"
				:can-deselect="false"
				:disabled="store.loading"
				label="เลือกเกณฑ์การจำแนกสภาพทาง"
				name="ref_condition_range_id1"
				placeholder="เลือก"
				@update:model-value="onReflectCriteriaChange"
			/>
		</div>
		<div
			v-show="store.conditionType !== EConditionType.Reflect && dashboardStore.isSingleRoad"
			class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12"
		>
			<VSelect
				v-model="storeRoadCondition.params.id_parent"
				:options="storeRoadCondition.getSurveyLaneOptions"
				label="ช่องจราจร"
				name="lane"
				placeholder="เลือก"
				:disabled="store.loading"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="() => storeRoadCondition.onUpdateIdParent()"
			/>
		</div>
		<div
			v-show="store.conditionType === EConditionType.Reflect && dashboardStore.isSingleRoad"
			class="col-xxl-4 col-xl-4 col-lg-6 col-md-4 col-sm-6 col-12"
		>
			<VSelect
				v-model="storeDashboardReflectiveStore.params.id_parent"
				:options="storeDashboardReflectiveStore.getLineOptions"
				label="ช่องจราจร"
				name="lane"
				placeholder="เลือก"
				:disabled="store.loading"
				:can-clear="false"
				:can-deselect="false"
				@update:model-value="() => storeDashboardReflectiveStore.onUpdateIdParent()"
			/>
		</div>
	</div>
	<!-- {{ store.loading }} -->
	<template v-if="store.loading === false">
		<div class="col-12 pt-5">
			<div v-if="store.conditionType === EConditionType.Reflect && dashboardStore.isSingleRoad">
				<DashboardReflectiveChart
					:collapsed="mapCollapsed"
					:search-call="props.searchCall"
					:reload-data="store.reloadData"
				/>
				<DashboardReflectiveTable />
			</div>
			<div v-if="store.conditionType !== EConditionType.Reflect && dashboardStore.isSingleRoad">
				<DashboardRoadConditionChart
					:collapsed="mapCollapsed"
					:search-call="props.searchCall"
					:reload-data="store.reloadData"
					:road-id="dashboardStore.params.road_id.map(Number)"
				/>

				<DashboardRoadConditionTable />
			</div>
		</div>
		<div class="col-12 pt-5">
			<DashboardConditionChart :collapsed="mapCollapsed" :search-call="props.searchCall" />
			<DashboardConditionTable />
		</div>
	</template>
	<template v-else>
		<div class="col-12 d-flex justify-content-between placeholder-glow mt-5">
			<div class="placeholder me-4 graph-loading"></div>
			<div class="placeholder graph-loading"></div>
		</div>
		<div class="col-12 placeholder-glow mt-6">
			<span class="placeholder text-loading" style="width: 30% !important"></span>
			<div class="table-loading placeholder mt-3"></div>
			<div class="placeholder-glow d-flex justify-content-between mt-4">
				<div class="placeholder-glow d-flex gap-2 align-items-center">
					<span class="placeholder w-50 text-loading"></span>
					<div class="placeholder select-loading"></div>
				</div>
				<div class="placeholder box-loading"></div>
				<!-- <div class="table-loading placeholder mt-3"></div> -->
			</div>
		</div>
	</template>
</template>

<style lang="scss" scoped>
.graph-loading {
	height: 300px;
	width: 50%;
	border-radius: 10px;
}

.text-loading {
	height: 15px;
	border-radius: 4px;
	width: 6rem !important;
}

.table-loading {
	height: 400px;
	width: 100%;
	border-radius: 8px;
}

.select-loading {
	height: 35px;
	width: 60px;
	border-radius: 8px;
}

.box-loading {
	height: 35px;
	width: 70px;
	border-radius: 8px;
}
</style>
