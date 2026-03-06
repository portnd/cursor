<script setup lang="ts">
import ChildAssetType1 from "./reportType/ChildAssetType1.vue"
import ChildAssetMapType2 from "./reportType/ChildAssetMapType2.vue"
import ChildRoadConditionAssetType3 from "./reportType/ChildRoadConditionAssetType3.vue"
import ChildRoadSummaryType4 from "./reportType/ChildRoadSummaryType4.vue"
import ChildRoadSurfaceType5 from "./reportType/ChildRoadSurfaceType5.vue"
import ChildRoadConditionType6 from "./reportType/ChildRoadConditionType6.vue"
import ChildRoadConditionSummaryType7 from "./reportType/ChildRoadConditionSummaryType7.vue"
import ChildRoadReflectLightType8 from "./reportType/ChildRoadReflectLightType8.vue"
import ChildRoadReflectLightSummaryType9 from "./reportType/ChildRoadReflectLightSummaryType9.vue"
import ChildRoadDamageType10 from "./reportType/ChildRoadDamageType10.vue"
import ChildRoadDamageSummaryType11 from "./reportType/ChildRoadDamageSummaryType11.vue"
import ReportMaintenanceType12 from "./reportType/ReportMaintenanceType12.vue"
import ReportProjectMaintenanceType13 from "./reportType/ReportProjectMaintenanceType13.vue"
import ReportVolumeAadtType14 from "./reportType/ReportVolumeAadtType14.vue"
// import { useReportStore } from "../store"

// const store = useReportStore()

const props = defineProps({
	accessKey: {
		type: String,
		default: "child_asset",
	},
})
</script>

<template>
	<div class="card p-6 mb-8">
		<h4 class="fw-semibold mb-5">
			<div class="d-flex align-items-center gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none">
					<g clip-path="url(#clip0_742_47662)">
						<path
							d="M15 7V0H5C3.343 0 2 1.343 2 3V24H22V7H15ZM13.406 19.417C13.018 19.805 12.509 19.998 12 19.998C11.491 19.998 10.981 19.805 10.594 19.417L7.295 16.118L8.709 14.704L11 16.995V10.998H13V17.006L15.291 14.704L16.705 16.118L13.406 19.417ZM21.414 5H17V0.586L21.414 5Z"
							fill="#FDB833"
						/>
					</g>
					<defs>
						<clipPath id="clip0_742_47662">
							<rect width="24" height="24" fill="white" />
						</clipPath>
					</defs>
				</svg>
				<span>ขั้นตอนที่ 2 : กำหนดเงื่อนไข และ ดาวน์โหลด</span>
			</div>
		</h4>
		<!-- Begin:: รายงานสินทรัพย์ -->
		<template v-if="props.accessKey === 'child_asset'">
			<ChildAssetType1 />
		</template>
		<template v-else-if="props.accessKey === 'child_asset_map'">
			<ChildAssetMapType2 />
		</template>
		<template v-else-if="props.accessKey === 'child_road_condition_asset'">
			<ChildRoadConditionAssetType3 />
		</template>

		<!-- End:: รายงานสินทรัพย์ -->

		<!-- Begin:: รายงานข้อมูลสภาพทาง -->
		<template v-else-if="props.accessKey === 'child_road_summary'">
			<ChildRoadSummaryType4 />
		</template>
		<template v-else-if="props.accessKey === 'child_road_surface'">
			<ChildRoadSurfaceType5 />
		</template>
		<template v-else-if="props.accessKey === 'child_road_condition'">
			<ChildRoadConditionType6 />
		</template>
		<template v-else-if="props.accessKey === 'child_road_condition_summary'">
			<ChildRoadConditionSummaryType7 />
		</template>
		<template v-else-if="props.accessKey === 'child_reflect_light'">
			<ChildRoadReflectLightType8 />
		</template>
		<template v-else-if="props.accessKey === 'child_reflect_light_summary'">
			<ChildRoadReflectLightSummaryType9 />
		</template>
		<!-- End:: รายงานข้อมูลสภาพทาง -->

		<!-- Begin:: รายงานข้อมูลความเสียหาย -->
		<template v-else-if="props.accessKey === 'child_road_damage'">
			<ChildRoadDamageType10 />
		</template>
		<template v-else-if="props.accessKey === 'child_road_damage_summary'">
			<ChildRoadDamageSummaryType11 />
		</template>
		<!-- End:: รายงานข้อมูลความเสียหาย -->

		<!-- Begin:: รายงานการซ่อมบำรุงตามเกณฑ์ KPI -->
		<template v-else-if="props.accessKey === 'maintenance'">
			<ReportMaintenanceType12 />
		</template>
		<!-- End:: รายงานการซ่อมบำรุงตามเกณฑ์ KPI -->

		<!-- Begin:: รายงานประวัติการซ่อมบำรุง -->
		<template v-else-if="props.accessKey === 'project_maintenance'">
			<ReportProjectMaintenanceType13 />
		</template>
		<!-- End:: รายงานประวัติการซ่อมบำรุง -->

		<!-- Begin:: รายงานปริมาณจราจร -->
		<template v-else>
			<ReportVolumeAadtType14 />
		</template>
		<!-- End:: รายงานปริมาณจราจร -->
	</div>
</template>

<style lang="scss" scoped>
.btn-code {
	border-radius: 16px;
	background-color: var(--kt-white);
	border: 1px solid #1f70f3 !important;
	color: #1f70f3;

	&:hover {
		color: var(--kt-white);
		background-color: #1f70f3;
	}
}
</style>
