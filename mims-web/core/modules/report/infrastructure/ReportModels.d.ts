export interface IReportsYears {
	year: number[]
}

export interface IReportsHistMaintenance {
	road_group_id: number
	road_group_name: string
	year: number[]
}

export interface IReportTrackingMaintenance {
	year: number
	road_group: IReportTrackingRoadGroup[]
}

export interface IReportTrackingRoadGroup {
	road_group_id: number
	road_group_name: string
}

export interface IReportRoadsTree {
	id: number
	label: string
	children: IReportRoadsTreeChild[]
}

export interface IReportRoadsTreeChild {
	id: number
	label: string
}

export interface IReportTrafficData {
	year: number[]
	road_group: IReportTrafficRoadGroup[]
}

export interface IReportTrafficRoadGroup {
	id: number
	name: string
}
export interface IReportAccidentData {
	year: number[]
	road_group: IReportAccidentRoadGroup[]
}

export interface IReportAccidentRoadGroup {
	id: number
	name: string
}

export interface IReportAssetData {
	road: IReportAssetRoad[]
	group: IReportAssetGroup[]
}

export interface IReportAssetGroup {
	id: number
	name: string
	asset: IReportAssetItem[]
}

export interface IReportAssetRoad {
	id: number
	name: string
}
export interface IReportAssetItem {
	id: number
	name: string
}

export interface IReportAssetSummaryData {
	budgetYear: number
	id: number
	name: string
}

export interface IReportAssetImproveData {
	month: string[]
	road: IReportAssetImproveRoadData[]
	year: number[]
}

interface IReportAssetImproveRoadData {
	id: number
	name: string
}

export interface IReportRoadConditionSummaryData {
	measure: IReportRoadConditionSummaryMeasure[]
	road: IReportRoadConditionSummaryRoad[]
	type: string[]
	year: number[]
}

interface IReportRoadConditionSummaryMeasure {
	id: number
	name: string
}

interface IReportRoadConditionSummaryRoad {
	id: number
	name: string
}

export interface IReportSurfaceSummary {
	year: number[]
	road: IReportSurfaceSummaryRoad[]
}

interface IReportSurfaceSummaryRoad {
	id: number
	name: string
}
