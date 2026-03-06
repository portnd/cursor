//รายงานสินทรัยย์
export interface IReportChildAssetRequest {
	department_id: number
	road_group_id: number
	road_section_id: number
	asset_group_id: number
	asset_id: number
	type: string
}
//รายงานสินทรัยย์

//รายงานแผนที่สินทรัยย์
export interface IReportChildAssetMapRequest {
	department_id: number
	road_group_id: number
	road_section_id: number
	asset_group_id: number
	asset_id: number
	type: string
}
//รายงานแผนที่สินทรัยย์

//รายงานสรุปสินทรัยย์
export interface IReportChildRoadConditionAssetRequest {
	department_id: number
	road_group_id: number
	road_section_id: number
	type: string
}
//รายงานสรุปสินทรัยย์

//รายงานสรุปข้อมูลสายทาง
export interface IReportChildRoadSummaryRequest {
	road_group_ids: number[]
	type: string
}
//รายงานสรุปข้อมูลสายทาง

//รายงานสรุปรายละเอียดชนิดผิวทาง
export interface IReportChildRoadSurfaceRequest {
	year: number
	department_id: number
	road_group_id: number
	road_section_id: number
	type: string
}
//รายงานสรุปรายละเอียดชนิดผิวทาง

//รายงานสรุปรายละเอียดชนิดผิวทาง
export interface IReportChildRoadConditionRequest {
	year: number
	department_id: number
	road_group_id: number
	road_section_id: number
	km: number
	type: string
}
//รายงานสรุปรายละเอียดชนิดผิวทาง

//รายงานสรุปข้อมูลสภาพทาง
export interface IReportChildRoadConditionSummaryRequest {
	year: number
	department_id: number
	road_group_id: number
	road_section_id: number
	condition_id: string
	criteria_id: number
	type: string
}
//รายงานสรุปข้อมูลสภาพทาง

//รายงานค่าการสะท้อนแสงของเส้นจราจร
export interface IReportChildRoadReflectLightRequest {
	year: number
	department_id: number
	road_group_id: number
	road_section_id: number
	km: number
	type: string
}
//รายงานค่าการสะท้อนแสงของเส้นจราจร

//รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร
export interface IReportChildRoadReflectLightSummaryRequest {
	year: number
	department_id: number
	road_group_id: number
	road_section_id: number
	criteria_id: number
	type: string
}
//รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร

//รายงานข้อมูลความเสียหาย
export interface IReportVolumeAadtRequest {
	year: number | null
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}
//รายงานข้อมูลความเสียหาย

//รายงานสรุปข้อมูลความเสียหาย
export interface IReportVolumeAadtRequest {
	year: number | null
	department_id: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}
//รายงานสรุปข้อมูลความเสียหาย

//รายงานการซ่อมบำรุงตามเกณฑ์ KPI
export interface IReportVolumeAadtRequest {
	condition_id: number | null
	year: number | null
	road_group_id: number | null
	road_section_id: number | null
	type: string
}
//รายงานการซ่อมบำรุงตามเกณฑ์ KPI

//รายงานประวัติการซ่อมบำรุง
export interface IReportVolumeAadtRequest {
	road_group_id: number | null
	road_section_id: number | null
	year_start: number | null
	year_end: number | null
	type: string
}
//รายงานประวัติการซ่อมบำรุง

//รายงานสรุปข้อมูลปริมาณจราจร
export interface IReportVolumeAadtRequest {
	year: number
	road_group_id: number | null
	road_section_id: number | null
	type: string
}
//รายงานสรุปข้อมูลปริมาณจราจร
