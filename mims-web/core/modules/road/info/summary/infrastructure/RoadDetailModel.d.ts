export interface IRoad {
	id: number
	seq: number
	parent_road_id: number
	road_level: number
	is_active: number
	ref_direction_id: number
	road_group_id: number
	is_init: boolean
	road_section_id: number
	created_by: number
	created_at: Date
	road_code: string
	road_section_name_th: string
	road_section_name_en: string
	province: string
	responsible_code: string
	origin_to_destination: string
	km_range: string
	distance: number
	road_info: IRoadInfo
	road_geom: IRoadGeom[]
	road_surface_icon: IRoadSurface[]
	ref_depot: IRoadRefDepot
}

export interface IRoadInfo {
	id: hnumber
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
	created_at: string
	created_by: number
	updated_at: string
	updated_by: number
	remark: string
	ref_road_type_id: number
	center_lane_shape_file_path: string
	center_line_shape_file_path: string
	origin_to_destination: string
	road_code: string
	responsible_code: string
	km_range: string
	user: IUser
	ref_road_type: IRoadType
	year_construction_completed: number
}

export interface IUser {
	id: number
	username: string
	firstname: string
	lastname: string
}

export interface IRoadType {
	id: number
	name: string
	icon: string
}

export interface IRoadGeom {
	id: number
	road_id: number
	lane_no: number
	km_start: number
	km_end: number
	the_geom: string
	revision: number
	status: string
	remark: string
	created_at: number
	created_by: Date
	updated_at: number
	updated_by: Date
}

export interface IRoadSurface {
	color_code: string
	id: number
	name: string
}

export interface IRoadRefDepot {
	id: number
	depot_code: string
	name: string
	the_geom: string
}
