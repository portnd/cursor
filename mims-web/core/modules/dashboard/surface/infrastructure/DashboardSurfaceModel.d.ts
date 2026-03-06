export interface IDashboardSurface {
	surface_dashboard_table: ISurfaceDetail[]
	geom_list: ISurfaceGeomList[]
	summary: ISurfaceSummary[]
}

export interface ISurfaceDetail {
	surface_name: string
	surface_lane_type: ISurfaceLane
}

export interface ISurfaceLane {
	one_lane: number
	two_lane: number
	three_lane: number
	four_lane: number
	more_than_four: number
}

export interface ISurfaceId {
	id: number
	name: string
	color_code: string
}

export interface ISurfaceGeomList {
	code: string
	geom_cl: Array<String>
	id: number
	surface: ISurfaceId
}

export interface ISurfaceSummary {
	surface: ISurfaceId
	value: number
}
