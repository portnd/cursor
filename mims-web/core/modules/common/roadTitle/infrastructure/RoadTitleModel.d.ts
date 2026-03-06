export interface IRoadTitle {
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
	name: Name
}
export interface IRoadType {
	id: number
	name: Name
}
