export interface IRequestReflectivityRule {
	name: string
	ref_reflectivity_range_id: number
	road_line_list: IRequestConditionReflectivityRule[]
}

export interface IRequestConditionReflectivityRule {
	grade_id: number
	left_value_white: number
	right_value_white: number
	left_value_yellow: number
	right_value_yellow: number
}
