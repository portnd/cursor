//รายงานสินทรัยย์
export interface IReportChildAssetModel {
	department_id: number
	road_group_id: number | null
	road_section_id: number | null
	asset_group_id: number
	asset_id: number | null
	type: string
}

export interface IRoadsAssetFilter {
	filter_road: IRoadsAssetFilterRoad[]
	filter_asset: IRoadsAssetFilterAssets[]
}

export interface IRoadsAssetFilterRoad {
	id: number
	name: string
	road_group: IRoadsAssetFilterRoadGroup[]
}

export interface IRoadsAssetFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsAssetFilterRoadSections[]
}

export interface IRoadsAssetFilterRoadSections {
	id: number
	name: string
}

export interface IRoadsAssetFilterAssets {
	id: number
	name: string
	asset: IRoadsAssetFilterAssetItems[]
}

export interface IRoadsAssetFilterAssetItems {
	id: number
	name: string
}
//รายงานสินทรัยย์

//รายงานแผนที่สินทรัยย์
export interface IReportChildAssetMapModel {
	department_id: number
	road_group_id: number | null
	road_section_id: number | null
	asset_group_id: number
	asset_id: number | null
	type: string
}

export interface IRoadsAssetMapFilter {
	filter_road: IRoadsAssetMapFilterRoad[]
	filter_asset: IRoadsAssetMapFilterAssets[]
}

export interface IRoadsAssetFilterMapRoad {
	id: number
	name: string
	road_group: IRoadsAssetMapFilterRoadGroup[]
}

export interface IRoadsAssetMapFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsAssetMapFilterRoadSections[]
}

export interface IRoadsAssetMapFilterRoadSections {
	id: number
	name: string
}

export interface IRoadsAssetMapFilterAssets {
	id: number
	name: string
	asset: IRoadsAssetFilterAssetItems[]
}

export interface IRoadsAssetMapFilterAssetItems {
	id: number
	name: string
}
//รายงานแผนที่สินทรัยย์

//รายงานสรุปสินทรัยย์
export interface IReportChildRoadConditionAssetModel {
	department_id: number
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IRoadsRoadConditionAssetFilter {
	filter_road: IRoadsRoadConditionAssetFilterRoad[]
}

export interface IRoadsAssetFilterMapRoad {
	id: number
	name: string
	road_group: IRoadsRoadConditionAssetFilterRoadGroup[]
}

export interface IRoadsRoadConditionAssetFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadConditionAssetFilterRoadSections[]
}

export interface IRoadsRoadConditionAssetFilterRoadSections {
	id: number
	name: string
}
//รายงานสรุปสินทรัยย์

//รายงานสรุปข้อมูลสายทาง
export interface IReportChildRoadSummaryModel {
	road_group_ids: number[]
	type: string
}

export interface IRoadsRoadSummaryFilter {
	filter_road: IRoadsRoadSummaryFilterRoad[]
}

export interface IRoadsRoadSummaryFilterRoad {
	id: number
	name: string
}
//รายงานสรุปข้อมูลสายทาง

//รายงานสรุปรายละเอียดชนิดผิวทาง
export interface IReportChildRoadSurfaceModel {
	year: number
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IRoadsRoadSurfaceFilter {
	filter_road: IRoadsRoadSurfaceFilterRoad[]
}

export interface IRoadsRoadSurfaceFilterRoad {
	year: number
	depot: IRoadsRoadSurfaceDepotFilterRoad[]
}

export interface IRoadsRoadSurfaceDepotFilterRoad {
	id: number
	name: string
	road_group: IRoadsRoadSurfaceFilterRoadGroup[]
}

export interface IRoadsRoadSurfaceFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadSurfaceFilterRoadSections[]
}

export interface IRoadsRoadSurfaceFilterRoadSections {
	id: number
	name: string
}
//รายงานสรุปรายละเอียดชนิดผิวทาง

//รายงานข้อมูลสภาพทาง
export interface IReportChildRoadConditionModel {
	year: number
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	km: number
	type: string
}

export interface IRoadsRoadConditionFilter {
	filter_road: IRoadsRoadConditionFilterRoad[]
	filter_range: number[]
}

export interface IRoadsRoadConditionFilterRoad {
	year: number
	depot: IRoadsRoadConditionDepotFilterRoad[]
}

export interface IRoadsRoadConditionDepotFilterRoad {
	id: number
	name: string
	road_group: IRoadsRoadConditionFilterRoadGroup[]
}

export interface IRoadsRoadConditionFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadConditionFilterRoadSections[]
}

export interface IRoadsRoadConditionFilterRoadSections {
	id: number
	name: string
}
//รายงานข้อมูลสภาพทาง

//รายงานสรุปข้อมูลสภาพทาง
export interface IReportChildRoadConditionSummaryModel {
	year: number
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	condition: string
	criteria_id: number
	type: string
}

export interface IRoadsRoadConditionSummaryFilter {
	filter_road: IRoadsRoadConditionSummaryFilterRoad[]
	filter_condition: string[]
	filter_criteria: IRoadsRoadConditionSummaryFilterCriteria[]
}

export interface IRoadsRoadConditionSummaryFilterRoad {
	year: number
	depot: IRoadsRoadConditionSummaryDepotFilterRoad[]
}

export interface IRoadsRoadConditionSummaryDepotFilterRoad {
	id: number
	name: string
	road_group: IRoadsRoadConditionSummaryFilterRoadGroup[]
}

export interface IRoadsRoadConditionSummaryFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadConditionSummaryFilterRoadSections[]
}

export interface IRoadsRoadConditionSummaryFilterRoadSections {
	id: number
	name: string
}

export interface IRoadsRoadConditionSummaryFilterCriteria {
	id: number
	name: string
}
//รายงานสรุปข้อมูลสภาพทาง

//รายงานค่าการสะท้อนแสงของเส้นจราจร
export interface IReportChildRoadReflectLightModel {
	year: number
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	km: number
	type: string
}

export interface IRoadsRoadReflectLightFilter {
	filter_road: IRoadsRoadReflectLightFilterRoad[]
	filter_criteria: IRoadsRoadReflectLightFilterCriteria[]
}

export interface IRoadsRoadReflectLightFilterRoad {
	year: number
	depot: IRoadsRoadReflectLightDepotFilterRoad[]
}

export interface IRoadsRoadReflectLightDepotFilterRoad {
	id: number
	name: string
	road_group: IRoadsRoadReflectLightFilterRoadGroup[]
}

export interface IRoadsRoadReflectLightFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadReflectLightFilterRoadSections[]
}

export interface IRoadsRoadReflectLightFilterRoadSections {
	id: number
	name: string
}

export interface IRoadsRoadReflectLightFilterCriteria {
	id: number
	name: string
}
//รายงานค่าการสะท้อนแสงของเส้นจราจร

//รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร
export interface IReportChildRoadReflectLightSummaryModel {
	year: number
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	criteria_id: number
	type: string
}

export interface IRoadsRoadReflectLightSummaryFilter {
	filter_road: IRoadsRoadReflectLightSummaryFilterRoad[]
	filter_criteria: IRoadsRoadReflectLightSummaryFilterCriteria[]
}

export interface IRoadsRoadReflectLightSummaryFilterRoad {
	year: number
	depot: IRoadsRoadReflectLightSummaryDepotFilterRoad[]
}

export interface IRoadsRoadReflectLightSummaryDepotFilterRoad {
	id: number
	name: string
	road_group: IRoadsRoadReflectLightSummaryFilterRoadGroup[]
}

export interface IRoadsRoadReflectLightSummaryFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadReflectLightSummaryFilterRoadSections[]
}

export interface IRoadsRoadReflectLightSummaryFilterRoadSections {
	id: number
	name: string
}

export interface IRoadsRoadReflectLightSummaryFilterCriteria {
	id: number
	name: string
}
//รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร

//รายงานข้อมูลความเสียหาย
export interface IReportChildRoadDamageModel {
	year: number | null
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IReportChildRoadDamageFilter {
	filter_road: IReportChildRoadDamageFilterCondition[]
}

export interface IReportChildRoadDamageFilterCondition {
	year: numbrt
	depot: IReportChildRoadDamageFilterDepot[]
}

export interface IReportChildRoadDamageFilterDepot {
	id: number
	name: string
	road_group: IReportChildRoadDamageFilterRoadGroup[]
}

export interface IReportChildRoadDamageFilterRoadGroup {
	id: number
	name: string
	road_section: IReportChildRoadDamageFilterSections[]
}

export interface IReportChildRoadDamageFilterSections {
	id: number
	name: string
}
//รายงานข้อมูลความเสียหาย

//รายงานสรุปข้อมูลความเสียหาย
export interface IReportChildRoadDamageSummaryModel {
	year: number | null
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IRoadsRoadDamageSummaryFilter {
	filter_road: IRoadsRoadDamageSummaryFilterCondition[]
}

export interface IRoadsRoadDamageSummaryFilterCondition {
	year: number
	depot: IRoadsRoadDamageSummaryFilterDepot[]
}

export interface IRoadsRoadDamageSummaryFilterDepot {
	id: number
	name: string
	road_group: IRoadsRoadDamageSummaryFilterRoadGroup[]
}

export interface IRoadsRoadDamageSummaryFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsRoadDamageSummaryFilterSections[]
}

export interface IRoadsRoadDamageSummaryFilterSections {
	id: number
	name: string
}
//รายงานสรุปข้อมูลความเสียหาย

//รายงานการซ่อมบำรุงตามเกณฑ์ KPI
export interface IReportMaintenanceModel {
	condition_name: string | null
	year: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IReportMaintenanceFilter {
	filter_condition: IReportMaintenanceFilterCondition[]
}

export interface IReportMaintenanceFilterCondition {
	name: string
	year: IReportMaintenanceFilterYear[]
}

export interface IReportMaintenanceFilterYear {
	year: number
	road_group: IReportMaintenanceFilterRoadGroup[]
}

export interface IReportMaintenanceFilterRoadGroup {
	id: number
	name: string
	road_section: IReportMaintenanceFilterSections[]
}

export interface IReportMaintenanceFilterSections {
	id: number
	name: string
}
//รายงานการซ่อมบำรุงตามเกณฑ์ KPI

//รายงานประวัติการซ่อมบำรุง
export interface IReportProjectMaintenanceModel {
	road_group_id: number | null
	road_section_id: number | null
	year_start: number | null
	year_end: number | null
	type: string
}

export interface IReportProjectMaintenanceFilter {
	filter_road: IReportProjectMaintenanceFilterRoadGroup[]
	filter_Year: IYearFilter
}

export interface IReportProjectMaintenanceFilterRoadGroup {
	id: number
	name: string
	road_section: IReportProjectMaintenanceFilterSections[]
}

export interface IReportProjectMaintenanceFilterSections {
	id: number
	name: string
}

export interface IYearFilter {
	start_year: number[]
	end_year: number[]
}
//รายงานประวัติการซ่อมบำรุง

//รายงานสรุปข้อมูลปริมาณจราจร
export interface IReportVolumeAadtModel {
	year: number
	road_group_id: number | null
	road_section_id: number | null
	type: string
}

export interface IRoadsVolumeAadtFilter {
	filter_road: IRoadsVolumeAadtFilterRoad[]
}

export interface IRoadsVolumeAadtFilterRoad {
	year: number
	road_group: IRoadsVolumeAadtFilterRoadGroup[]
}

export interface IRoadsVolumeAadtFilterRoadGroup {
	id: number
	name: string
	road_section: IRoadsVolumeAadtFilterSections[]
}

export interface IRoadsVolumeAadtFilterSections {
	id: number
	name: string
}
//รายงานสรุปข้อมูลปริมาณจราจร
