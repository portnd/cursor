export interface IResSurveyRule {
	id: number
	ref_condition_range_id: number
	owner_name: string
	condition_list: ICondition[]
}

export interface ICondition {
	condition_type: string
	surface_type: ISurfaceType
}

export interface ISurfaceType {
	ac: ConditionAC[]
	cc: ConditionCC[]
}

export interface ConditionAC {
	grade: Grade
	left_value_ac: number
	left_condition_ac?: String
	right_value_ac: number
	right_condition_ac?: String
}

export interface ConditionCC {
	grade: Grade
	left_value_cc: number
	left_condition_cc?: String
	right_value_cc: number
	right_condition_cc?: String
}

export interface Grade {
	id: number
	name?: string
	color?: string
}
