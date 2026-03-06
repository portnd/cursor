import { IFile } from "~~/core/shared/types/File"
export interface IRequestRoad {
	road_code: string
	name: string
	road_group_id: number
	ref_owner_code: number
	ref_district_code: number
	origin: string
	destination: string
	km_start: number
	km_end: number
	road_color_code: string
	ref_road_type_id: number
	completion_year: number
	register_date: string
	center_line_shape_file: IFile
	center_line_shape_filepath?: string
	center_line_shape_file_status: string
	center_lane_shape_file: IFile
	center_lane_shape_filepath?: string
	center_lane_shape_file_status: string
	remark: string
	ramp_id: number
	year_construction_completed: number
}
