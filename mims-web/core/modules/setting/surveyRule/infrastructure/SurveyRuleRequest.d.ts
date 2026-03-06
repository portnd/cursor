export interface IRequestSurveyRule {
	name: string
	ref_condition_range_id: number
	condition_list: IRequestConditionSurveyRule[]
}

export interface IRequestConditionSurveyRule {
	grade_id: number
	left_value_ac: number
	right_value_ac: number
	left_value_cc: number
	right_value_cc: number
	condition_type: string
}
