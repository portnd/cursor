export interface IReflectivityRule {
	id: number
	owner_name: string
	ref_reflectivity_range_id: number
	road_line: IRoadLine
}

export interface IRoadLine {
	yellow: IYellowLine[]
	white: IWhiteLine[]
}

export interface IYellowLine {
	grade: Grade
	left_value_yellow: number
	left_condition_yellow: string
	right_value_yellow: number
	right_condition_yellow: string
}

export interface IWhiteLine {
	grade: Grade
	left_value_white: number
	left_condition_white: string
	right_value_white: number
	right_condition_white: string
}

export interface Grade {
	id: number
	name?: string
	color?: string
}
