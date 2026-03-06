export interface IMaintenanceRequest {
	[key: string]: IItemIntervention
}

interface IItemIntervention {
	[key: string]: IMaintenanceItemList[]
}
export interface IMaintenanceItemList {
	id?: number
	maintenance_method: string
	maintenance_cost_per_unit: number | null
	maintenance_description: string
	maintenance_scraping: number | null
	maintenance_sequence?: number
	maintenance_standard_name: string
	maintenance_surface_type_id: number
	maintenance_thickness: number | null
	maintenance_type: string
	maintenance_condition: ICondition[]
}

export interface IConditionList {
	id?: number
	condition_criterion: string
	condition_link: string
	condition_operation_1: string
	condition_operation_2: string
	condition_value_1: number
	condition_value_2: number
}

export interface IAnalysisRuleItemList {
	id?: number
	maintenance_method: string
	maintenance_cost_per_unit: number | null
	maintenance_description: string
	maintenance_scraping: number | null
	maintenance_sequence?: number
	maintenance_standard_name: string
	maintenance_surface_type_id: number
	maintenance_thickness: number | null
	maintenance_type: string
	maintenance_condition: ICondition[]
}
