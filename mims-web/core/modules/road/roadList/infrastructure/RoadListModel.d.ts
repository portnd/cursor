export interface IRoadData {
	id: number
	name: string
	code: string
	count_waiting: number
	count_rejected: number
	roads: IRoad[]
}

export interface IRoad {
	road_id: number
	code: string
	seq: number
	name: string
	km_start: number
	km_end: number
	road_level: number
	parent_road_id: number
	status_latest: string
	lane_count: number
	direction: IDirection
	road_type: IRoadType
	geom_cl: string
	count_waiting: number
	count_rejected: number
	road_color_code: string
	road_type_icon_id: number
	child_road: IChildRoad[]
}

export interface IChildRoad {
	road_id: number
	code: string
	seq: number
	name: string
	km_start: number
	km_end: number
	road_level: number
	parent_road_id: number
	status_latest: string
	lane_count: number
	direction: IDirection
	road_type: IRoadType
	geom_cl: string
	count_waiting: number
	count_rejected: number
	road_color_code: string
	road_type_icon_id: number
}

export interface IDirection {
	id: number
	name: string
}
export interface IRoadType {
	id: number
	name: string
}

export interface IRoadList {
	id: number
	number: string
	name: string
	short_name: string
	km_start: number
	km_end: number
	distance: number
	sections: IRoadListSection[]
}

export interface IRoadListSection {
	id: number
	road_group_id: number
	number: string
	name_origin_th: string
	name_destination_th: string
	name_origin_en: string
	name_destination_en: string
	km_start: number
	km_end: number
	distance: number
	province: string[]
	ref_division: IRefDi
	ref_district: IRefDi
	ref_depot: IRoadListRefDepot
	roads: IRoadListRoads[]
}

export interface IRoadListRefDepot {
	id: number
	depot_code: string
	name: string
	the_geom: string
}

export interface IRefDi {
	id: number
	district_code: string
	name: string
	name_en: string
	the_geom: string
	division_code: string
}

export interface IRoadListRoads {
	id: number
	seq: number
	parent_road_id: number
	road_level: number
	road_code: string
	is_active: boolean
	condition_status: boolean
	survey_status: boolean
	retro_status: boolean
	damage_status: boolean
	condition_status_color: string
	survey_status_color: string
	retro_status_color: string
	damage_status_color: string
	road_group_id: number
	road_section_id: number
	created_by: number
	createdAt: Date
	road_info: IRoadListRoadInfo
	road_geom: IRoadListRoadGeom[]
	child_roads: IRoadListRoads[]
	road_surface_icon: IRoadListSurface[]
}

export interface IRoadListRoadGeom {
	id: number
	road_id: number
	lane_no: number
	km_start: number
	km_end: number
	the_geom: string
	revision: number
	status: Status
	remark: string
	created_by: number
	created_date: Date
	updated_by: number
	updated_date: Date
}

export interface IRoadListRoadInfo {
	id: number
	road_id: number
	year: number
	ref_direction_id: number
	name: string
	km_start: number
	km_end: number
	the_geom: string
	revision: number
	status: string
	ramp_id: string
	road_color_code: string
	created_by: number
	created_At: Date
	updated_by: number
	updated_At: Date
	ref_road_type: IRefRoadType
	origin_to_destination: string
	road_code: string
	responsible_code: string
	km_range: string
	direction: IRoadListRoadInfoDirection
}

export interface IRoadListRoadInfoDirection {
	id: number
	name: string
}

export interface IRefRoadType {
	id: number
	name: string
	icon: string
}

export interface IRoadListSurface {
	id: number
	name: string
	color_code: string
}
