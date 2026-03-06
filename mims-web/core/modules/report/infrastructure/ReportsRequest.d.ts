export interface IReportConditionParams {
	year: number
	road_id: number
	distance: number
	type: string
}

export interface IReportDamageParams {
	year: number
	road_id: number
	type: string
}

export interface IReportMaintenanceHistoryParams {
	year_start: number
	year_end: number
	road_group_id: number[]
	type: string
}

export interface IReportMaintenanceTrackingParams {
	year: number
	road_group_id: number[]
	type: string
}

export interface IReportTrafficParams {
	year: number
	road_group_id: number
	type: string
}
export interface IReportAccidentParams {
	year: number
	road_group_id: number
	type: string
}

export interface IReportAsset {
	road_id: number
	asset_id: number
	type: string
}

export interface IReportAssetSummary {
	road_id: number
	type: string
}

export interface IReportAssetImprove {
	year: number
	month: string
	road_id: number
	type: string
}

export interface IReportRoadConditionSummary {
	factor: string
	road_id: number
	year: number
	measure_id: number
	type: string
}

export interface IReportRoadSurfaceSummary {
	year: number
	road_id: number
	type: sting
}
