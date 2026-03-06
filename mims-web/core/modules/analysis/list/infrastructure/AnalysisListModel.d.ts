export interface IStrategicsList {
	id: number
	name: string
	budget: IStrategicsBudget[]
}

export interface IStrategicsBudget {
	id: number
	maintenance_analysis_strategic_type_id: number
	name: string
	target: IStrategicsTarget[]
}

export interface IStrategicsTarget {
	id: number
	maintenance_analysis_strategic_budget_type_id: number
	name: string
}

export interface IMaintenanceAnalysis {
	id: number
	name: string
	type_analysis: string
	maintenance_analysis_type_id: number
	maintenance_condition_type_id: number
	comment: string
	analysis_date: Date
	percentage: number
	status: string
	is_favorite: boolean
}

export interface IMaintenanceAnalysisCondition {
	road_name: string
	filter: string
	lane: string
	km_group: string
	discount: number
	condition: string
	condition_filter: string
	target: string
}
