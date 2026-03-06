import { IFile } from "~/core/shared/types/File"
export interface IRequestRoadList {
	keyword: string
	road_group_id: string
	road_section_id: string
	km_start: string
	km_end: string
	depot_code: string
	road_id: string
	ref_surface_id: string
	is_iri_1000?: boolean
	is_iri_100?: boolean
	is_rut_100?: boolean
	is_ifi_100?: boolean
	is_g7_100?: boolean
}

export interface IRequestRoadInit {
	road_code: string
	road_section_name_th: string
	road_section_name_en: string
	province: string
	district: string
	depot: string
	origin: string
	destination: string
}

export interface ICreateRoad {
	name: string
	road_section_id: number
	road_id: number
	road_level: number
	ramp_id: number
	km_start: number
	km_end: number
	road_color_code: string
	ref_road_type_id: number
	register_date: string
	center_lane_shape_file: file
	center_line_shape_file: file
	remark: string
	year_construction_completed: number
}

// export interface IRequestRoadList {
// 	keyword: string
// 	direction: string
// 	road_type: string
// }
