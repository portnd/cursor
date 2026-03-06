export interface IMaintenance {
	[key: string]: IMaintenanceType
}

export interface IMaintenanceType {
	[key: string]: IMaintenanceItem[]
}

export interface IMaintenanceItem {
	id: number
	maintenance_method: string
	maintenance_cost_per_unit: number | null
	maintenance_description: string
	maintenance_scraping: number | null
	maintenance_sequence: number | null
	maintenance_standard_name: string
	maintenance_surface_type_id: number | null
	maintenance_thickness: number | null
	maintenance_type: string
	maintenance_condition: ICondition[]
	isNew?: boolean
}

export interface ICondition {
	id: number
	condition_criterion: string
	condition_link: string
	condition_operation_1: string
	condition_operation_2: string
	condition_value_1: number
	condition_value_2: number
}

export interface IAnalysiSurfacesRule {
	asphalt: ISurfacesRuleItem[]
	concrete: ISurfacesRuleItem[]
}

export interface ISurfacesRuleItem {
	id: number
	name: string
}

export interface IAnalysisRule {
	[key: string]: IAnalysisMethod[]
}

export interface IAnalysisMethod {
	id: number
	intervention_criterias: IAnalysisRuleItem[]
	name: string
}

export interface IAnalysisRuleItem {
	id: number
	maintenance_method: string
	maintenance_cost_per_unit: number | null
	maintenance_description: string
	maintenance_scraping: number | null
	maintenance_sequence: number | null
	maintenance_standard_name: string
	maintenance_surface_type_id: number | null
	maintenance_thickness: number | null
	maintenance_type: string
	maintenance_condition: ICondition[]
	is_new?: boolean
}

export interface ICondition {
	id?: number | null
	condition_criterion: string
	condition_link: string
	condition_operation_1: string
	condition_operation_2: string
	condition_value_1: number
	condition_value_2: number
}
