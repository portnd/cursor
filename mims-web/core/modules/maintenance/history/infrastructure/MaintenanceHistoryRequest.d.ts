import { IMultiFile } from "~/core/shared/types/File"

export interface IMaintenanceHistorySearch {
	road_group_id: number
	budget_year: number
	budget_type_id: number
	budget_method_id: number
	budget_maintenance: number
	name: string
	owner_code: string
	page: number
	limit: number
}

export interface IMaintenanceHistoryCreateRequest {
	advisor_name: string
	attachments: IFile[]
	budget_id: number
	budget_maintenance: number | null
	budget_method_id: number
	budget_procurement: number | null
	budget_year: number
	contract_number: string
	contract_work_value: number | null
	contractor_name: string
	guarantee_expiration_date: string
	middle_price: number | null
	name: string
	owner_code: string | null
	project_details: string
	project_end_date: string
	project_secretary_name: string
}
export interface IMaintenanceHistoryUpdateRequest {
	advisor_name: string
	attachments: IFile[]
	budget_id: number | null
	budget_maintenance: number | null
	budget_method_id: number | null
	budget_procurement: number | null
	budget_year: number
	contract_number: string
	contract_work_value: number | null
	contractor_name: string
	guarantee_expiration_date: string
	middle_price: number | null
	name: string
	owner_code: string | null
	project_details: string
	project_end_date: string
	project_secretary_name: string
}

export interface IMaintenanceHistoryGuaranteeCreateParams {
	attachments: IMultiFile[]
	id: number
	intervention_criteria_id: number | null
	km_end: number
	km_start: number
	lane: number
	road_group_id: number
	road_id: number
}

export interface IMaintenanceHistoryEditParams {
	attachments: IMultiFile[]
	id: number
	intervention_criteria_id: number | null
	km_end: number
	km_start: number
	lane: number
	road_group_id: number
	road_id: number
}

export interface IMaintenanceHistoryFileParams {
	file_type: string
	order: string
}

export interface IMaintenanceHistoryRoadsUpdateRequest {
	grid_no: number | null
	intervention_criteria_id: number | null
	km_end: number | null
	km_start: number | null
	lane_no: number | null
	maintenance_type: number | null
	road_id: number | null
}

export interface IMaintenanceWarrantyCreateRequest {
	grid_no: number | null
	intervention_criteria_id: number | null
	km_end: number | null
	km_start: number | null
	lane_no: number | null
	maintenance_type: number | null
	road_id: number | null
}
export interface IMaintenanceWarrantyUpdateRequest {
	grid_no: number | null
	intervention_criteria_id: number | null
	km_end: number | null
	km_start: number | null
	lane_no: number | null
	maintenance_type: number | null
	road_id: number | null
}
