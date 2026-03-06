export interface IDashboardMapAssetRequest {
	depot_code: string[]
	road_id: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
	ref_asset_id: number[]
	left: number | null
	bottom: number | null
	right: number | null
	top: number | null
	zoom: number | null
}

export interface IDashboardAssetRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}

export interface IDashboardAssetDetailsRequest {
	road_id: string[]
	ref_asset_id: number[]
	asset_type: string
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}

export interface IDashboardSurfaceRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}

export interface IDashboardSurfaceMapRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
	display: number
}

export interface IDashboardMaintenanceRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}
export interface IDashboardMaintenanceTableRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}

export interface IDashboardMaintenanceMapRequest {
	road_id: string[]
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	year: number | null
}

export interface IDashboardConditionMapRequest {
	condition_type: number
	condition_owner_id: number | null
	road_id: string[]
	year: number | null
	depot_code: string[]
	km_start: number | null
	km_end: number | null
	limit: number
	left: number
	bottom: number
	right: number
	top: number
	page: number
}
