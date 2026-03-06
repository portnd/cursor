import { IFile } from "~/core/shared/types/File"

export interface IMaintenanceHistoryList {
	status: boolean
	code: number
	data: IMaintenanceHistoryListData
}

export interface IMaintenanceHistoryListData {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IMaintenanceHistoryListItem[]
}

export interface IMaintenanceHistoryDevisionList {
	id: number
	name: string
	owner_code_key: string
	division_code: string
	districts: IMaintenanceHistoryDevisionDistricts[]
}
export interface IMaintenanceHistoryDevisionDistricts {
	id: number
	name: string
	depots: IMaintenanceHistoryDevisionDistrictsDepots[]
	district_code: string
	owner_code_key: string
}
export interface IMaintenanceHistoryDevisionDistrictsDepots {
	id: number
	name: string
	depot_code: string
	owner_code_key: string
}

export interface IMaintenanceHistoryRoadDropdownList {
	id: number
	road_number: string
	road_sections: IMaintenanceHistoryRoadDropdownRoadSectionsList[]
}
export interface IMaintenanceHistoryRoadDropdownRoadSectionsList {
	id: number
	name_destination: string
	name_origin: string
	number: string
	road_group_id: number
	roads: IMaintenanceHistoryRoadDropdownRoadSectionsRoadsList[]
}
export interface IMaintenanceHistoryRoadDropdownRoadSectionsRoadsList {
	id: number
	name: string
	km_end: number
	km_start: number
	ref_direction_id: number
}

export interface IMaintenanceHistoryListItem {
	id: number
	id_parent: number
	name: string
	budget_year: number
	contract_number: string
	budget_maintenance: number
	budget_id: number
	budget_method_id: number
	contractor_name: string
	budget_procurement: number
	adviser_name: string
	project_end_date: Date
	project_guarantee_expiration_date: Date
	middle_price: number
	contract_work_value: number
	project_secretary_name: string
	project_details: string
	is_complete: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	last_inspection_date: Date
	guarantee_expiration_date: Date
	budget: IMaintenanceHistoryListBudget
	budget_method: IMaintenanceHistoryListBudgetMethod
	road_group: IMaintenanceHistoryListRoadGroup
	km_total: number
	percent_progress: number
	percent_pay: number
	color: string
	remaining_time_header: string
	remaining_time: string
	// maintenance_roads: IMaintenanceHistoryDetailMaintenanceRoad[]
	roads: IMaintenanceHistoryDetailMaintenanceRoad[]
	// road_group_names: string
	road_group_names: Array<string | null>
	road_name: string
	lane_no: number
}

export interface IMaintenanceHistoryListBudget {
	id: number
	name: string
}

export interface IMaintenanceHistoryBudgetCriteriaData {
	id: number
	name: string
	can_delete: boolean
	is_deleted: boolean
	updated_by: number
	created_by: number
	updated_at: Date
	created_at: Date
	budget_methods: IMaintenanceHistoryBudgetCriteriaBudgetMethods
}

export interface IMaintenanceHistoryBudgetCriteriaBudgetMethods {
	id: number
	method_name: string
	budget_id: number
	cost_per_unit: number
	is_show_method: boolean
	is_deleted: boolean
	updated_by: number
	created_by: number
	updated_at: Date
	created_at: Date
}

export interface IMaintenanceHistoryRoadGroupList {
	status: boolean
	code: number
	data: IMaintenanceHistoryRoadGroupListData[]
}

export interface IMaintenanceHistoryRoadGroupListData {
	code: string
	id: number
	name: string
}

export interface IMaintenanceHistoryDetails {
	status: boolean
	code: number
	data: IMaintenanceHistoryDetailData
}

export interface IMaintenanceHistoryDetailData {
	id: number
	name: string
	budget_year: number
	contract_number: string
	budget_maintenance: number
	budget_id: number
	budget_method_id: number
	contractor_name: string
	budget_procurement: number
	adviser_name: string
	project_end_date: Date
	project_guarantee_expiration_date: Date
	middle_price: number
	contract_work_value: number
	project_secretary_name: string
	project_details: string
	is_complete: boolean
	created_by: number
	updated_by: IMaintenanceHistoryDetailUpdatedBy
	created_at: Date
	updated_at: Date
	last_inspection_date: Date
	guarantee_expiration_date: Date
	budget: IMaintenanceHistoryDetailBudget
	budget_method: IMaintenanceHistoryDetailBudgetMethod
	is_show_method: boolean
	maintenance_roads: IMaintenanceHistoryDetailMaintenanceRoad[]
	maintenance_road_histories: IMaintenanceHistoryDetailMaintenanceRoadHistory[]
	owner_name: string
	advisor_name: string
	attachments: IMaintenanceHistoryDetailAttrachment[]
	roads: IMaintenanceHistoryDetailMaintenanceRoad[]
	road_histories: IMaintenanceHistoryDetailMaintenanceRoadHistory[]
	road_group_names: string
	color: string
	remaining_time: string
	ref_depot: IMaintenanceRefDepot
}

export interface IMaintenanceHistoryDetailUpdatedBy {
	id: number
	email: string
	department_id: number
	firstname: string
	lastname: string
	profile_img_path: string
	status: string
	tel: string
	created_by: number
	updated_by: number
	department: IMaintenanceHistoryDetailDepartment
}

export interface IMaintenanceHistoryDetailDepartment {
	id: number
	name: string
	can_delete: boolean
}

export interface IMaintenanceHistoryDetailBudget {
	id: number
	name: string
}

export interface IMaintenanceHistoryDetailBudgetMethod {
	id: number
	method_name: string
}

export interface IMaintenanceHistoryDetailMaintenanceRoad {
	id: number
	maintenance_id: number
	road_group_id: number
	road_id: number
	lane: number
	km_start: number
	km_end: number
	maintenance_method_id: number
	ref_surface_id: number
	ref_surface_params_id: number
	intervention_criteria_id: number
	intervention_criteria_id_params: number
	is_show_method: boolean
	// the_geom: string
	the_geom: IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom
	road_group: IMaintenanceHistoryDetailRoadGroup
	road_info: IMaintenanceHistoryDetailRoadInfo
	intervention_criteria: IMaintenanceHistoryDetailInterventionCriteria
	road: IMaintenanceHistoryDetailRoadInfo
	road_group_name: string
	road_name: string
	lane_no: string
	color: string
	distance: number
	ref_direction_id: number
	lane_total: number
	grid_no: number
	maintenance_type: number
}
export interface IMaintenanceHistoryDetailMaintenanceRoadHistory {
	id: number
	maintenance_id: number
	road_group_id: number
	road_id: number
	lane: number
	km_start: number
	km_end: number
	maintenance_method_id: number
	ref_surface_id: number
	ref_surface_params_id: number
	intervention_criteria_id: number
	intervention_criteria_id_params: number
	ref_direction_id: number
	is_show_method: boolean
	// the_geom: string
	the_geom: IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom
	road_group: IMaintenanceHistoryDetailRoadGroup
	road_info: IMaintenanceHistoryDetailRoadInfo
	intervention_criteria: IMaintenanceHistoryDetailInterventionCriteria
	attacchments: IMaintenanceHistoryDetailAttrachment[]
	road: IMaintenanceHistoryDetailRoadInfo
	color: string
	road_name: string
	road_group_name: string
	lane_no: string
	distance: number
	lane_total: number
	grid_no: number
	maintenance_type: number
}

export interface IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom {
	type: string
	coordinates: number[][]
	// coordinates: [number, number];
}

export interface IMaintenanceHistoryDetailAttrachment {
	id: number
	maintenance_id: number
	maintenance_road_his_id: number
	path: string
	file_name: string
	file_type: string
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export interface IMaintenanceHistoryDetailInterventionCriteria {
	id: number
	maintenance_standard_name: string
}

export interface IMaintenanceHistoryDetailRoadGroup {
	id: number
	code: string
	name: string
}

export interface IMaintenanceHistoryDetailRoadInfo {
	id: number
	road_id: number
	ref_direction_id: number
	name: string
	road_color_code: string
	direction: IMaintenanceHistoryDetailDirections
	color: string
	km_start: number
	km_end: number
	road_name: string
	km_total: number
	lane_no: number
	distance: number
	project_end_date: string
}

export interface IMaintenanceHistoryDetailDirections {
	id: number
	name: string
}

export interface IMaintenanceHistoryBudgets {
	status: boolean
	code: number
	data: IMaintenanceHistoryBudgetsData[]
}

export interface IMaintenanceHistoryBudgetsData {
	id: number
	name: string
}

export interface IRoadChildList {
	status: boolean
	code: number
	data: IRoadChildListData[]
}

export interface IRoadChildListData {
	id: number
	roads: IRoadChildListItem[]
	roads_group_name: string
}

export interface IRoadChildListItem {
	id: number
	lanes: IRoadChildListLane[]
	name: string
}

export interface IRoadChildListLane {
	km_end: number
	km_start: number
	lane: number
}

export interface IMaintenanceMethodList {
	status: boolean
	code: number
	data: IMaintenanceMethodListData[]
}

export interface IMaintenanceMethodListData {
	id: number
	label: string
	children: IMaintenanceMethodListChild[]
}

export interface IMaintenanceMethodListChild {
	id: number
	label: string
}

export interface IMaintenancePlanList {
	status: boolean
	code: number
	data: IMaintenancePlanListData[]
}

export interface IMaintenancePlanListData {
	created_at: string
	created_by: number
	end_date: string
	id: number
	is_current: boolean
	maintenance_id: number
	name: string
	start_date: string
	updated_at: string
	updated_by: number
}

export interface IMaintenanceHistoryPlanStatus {
	status: boolean
	code: number
	data: IMaintenanceHistoryPlanStatusData[]
}

export interface IMaintenanceHistoryPlanStatusData {
	is_current: boolean
	name: string
	schedules: IMaintenanceHistoryPlanStatusSchedule[]
}

export interface IMaintenanceHistoryPlanStatusSchedule {
	disbursement_date: Date
	is_checked: boolean
	schedule: string
	status: string
}

export interface IMaintenanceHistoryAttrachments {
	created_at: string
	created_by: IMaintenanceHistoryAttrachmentsCreatedBy
	file_name: string
	file_type: string
	id: number
	maintenance_id: number
	maintenance_plan_detail_progress_id: number
	path: string
	updated_at: string
	updated_by: number
}

export interface IMaintenanceHistoryAttrachmentsCreatedBy {
	created_by: number
	department: IMaintenanceHistoryAttrachmentsDepartment
	department_id: number
	email: string
	firstname: string
	id: number
	lastname: string
	profile_img_path: string
	status: string
	tel: string
	updated_by: number
}

export interface IMaintenanceHistoryAttrachmentsDepartment {
	id: number
	name: string
	can_delete: boolean
}

export interface IMaintenanceHistoryPlanGraph {
	data: IMaintenanceHistoryPlanGraphData[]
	name: string
	schedule: string[]
}

export interface IMaintenanceHistoryPlanGraphData {
	name: string
	data: Array<number | null>
	color: string
}

// export interface IMaintenanceHistoryPlanGraphData {
// 	color: string
// 	data: string
// 	name: string
// }

export interface IMaintenanceHistoryPlanTable {
	plan_name: string
	problems: string[]
	value: IMaintenanceHistoryPlanTableValue[]
}

export interface IMaintenanceHistoryPlanTableValue {
	schedule: string
	plan: number
	plan_total: number
	progress_plan: number
	progress_plan_total: number
	disbursement_plan: number
	disbursement_plan_total: number
	disbursement_progress: number
	disbursement_progress_total: number
}

export interface IMaintenanceBudgetCriteria {
	id: number
	name: string
	can_delete: boolean
	is_deleted: boolean
	updated_by: number
	created_by: number
	updated_at: Date
	created_at: Date
	budget_methods: IMaintenanceBudgetCriteriaMethod[]
}

export interface IMaintenanceBudgetCriteriaMethod {
	id: number
	method_name: string
	budget_id: number
	cost_per_unit: number
	is_show_method: boolean
	is_deleted: boolean
	updated_by: number
	created_by: number
	updated_at: Date
	created_at: Date
}

export interface IPlanProgressGraphReportHistTableData {
	problems: string[]
	maintenance_plan: IMaintenanceReportHistTablePlan[]
}

export interface IMaintenanceReportHistTablePlan {
	plan_name: string
	value: IMaintenanceReportHistTablePlanValue[]
}

export interface IMaintenanceReportHistTablePlanValue {
	schedule: string
	plan: number
	plan_total: number
	progress_plan: number | null
	progress_plan_total: number
	disbursement_plan: number
	disbursement_plan_total: number
	disbursement_progress: number | null
	disbursement_progress_total: number
}

export interface IMaintenanceDefaultData {
	attachments: IMaintenanceAttrachments[]
	advisor_name: string
	budget: IMaintenanceBudget
	budget_maintenance: number
	budget_method: IMaintenanceBudgetMethod
	budget_procurement: number
	budget_year: number
	color: string
	contract_number: string
	contract_work_value: number
	contractor_name: string
	created_at: string
	created_by: IMaintenanceUpdateBy
	guarantee_expiration_date: string
	id: number
	id_parent: number
	km_total: number
	middle_price: number
	name: string
	owner_code: string
	owner_name: string
	is_show_method: boolean
	project_details: string
	project_end_date: string
	project_secretary_name: string
	ref_depot_code: string
	ref_district_code: string
	ref_division_code: string
	remaining_time: string
	revision: number
	road_group_names: string
	roads: IMaintenanceRoad[]
	status: string
	update_by: IMaintenanceUpdateBy
	updated_at: string
}

export interface IMaintenanceAttrachments {
	id: number
	path: string
	file_name: string
}

export interface IMaintenanceBudget {
	id: number
	name: string
}

export interface IMaintenanceBudgetMethod {
	id: number
	method_name: string
}

export interface IMaintenanceUpdateBy {
	depart_name: string
	id: number
	name: string
	profile_pic: string
}

export interface IMaintenanceRoad {
	color: string
	distance: number
	grid_no: number
	id: number
	id_parent: number
	intervention_criteria: InterventionCriteria
	km_end: number
	km_start: number
	lane_no: number
	lane_total: number
	ref_direction_id: number
	ref_direction_name: string
	road_group_name: string
	road_id: number
	road_name: string
	the_geom: IMaintenanceTheGeom
}

export interface IMaintenanceInterventionCriteria {
	id: number
	maintenance_standard_name: string
}

export interface IMaintenanceTheGeom {
	coordinates: Array<number[]>
	type: string
}

export interface IMaintenanceDivision {
	id: number
	division_code: string
	owner_code_key: string
	name: string
	districts: IMaintenanceDivisionDistrict[]
}

export interface IMaintenanceDivisionDistrict {
	id: number
	district_code: string
	name: string
	owner_code_key: string
	depots: IMaintenanceDivisionDepot[]
}

export interface IMaintenanceDivisionDepot {
	id: number
	depot_code: string
	name: string
	owner_code_key: string
}

export interface IMaintenanceRoadGroup {
	id: number
	road_number: string
	short_name: string
	road_sections: IMaintenanceRoadSection[]
}

export interface IMaintenanceRoadSection {
	id: number
	road_group_id: number
	number: string
	name_origin: string
	name_destination: string
	roads?: IMaintenanceRoadOptions[]
}

export interface IMaintenanceRoadOptions {
	id: number
	name: string
	km_end: number
	km_start: number
	lane_total: number
	ref_direction_id: number
}

export interface IMaintenanceRoadData {
	color: string
	distance: number
	grid_no: number
	id: number
	id_parent: number
	intervention_criteria: IMaintenanceRoadInterventionCriteria
	km_end: number
	km_start: number
	lane_no: number | null
	maintenance_type: number
	lane_total: number
	ref_direction_id: number
	ref_direction_name: string
	road_group_name: string
	road_id: number
	road_name: string
	the_geom: IMaintenanceRoadTheGeom
}

export interface IMaintenanceRoadInterventionCriteria {
	id: number
	maintenance_standard_name: string
}

export interface IMaintenanceRoadTheGeom {
	coordinates: Array<number[]>
	type: string
}
export interface IMaintenanceWarrantyData {
	color: string
	distance: number
	grid_no: number
	id: number
	id_parent: number
	intervention_criteria: IMaintenanceWarrantyInterventionCriteria
	km_end: number
	km_start: number
	lane_no: number
	maintenance_type: number
	lane_total: number
	ref_direction_id: number
	ref_direction_name: string
	road_group_name: string
	road_id: number
	road_name: string
	the_geom: IMaintenanceWarrantyTheGeom
}

export interface IMaintenanceWarrantyInterventionCriteria {
	id: number
	maintenance_standard_name: string
}

export interface IMaintenanceRoadTheGeom {
	coordinates: Array<number[]>
	type: string
}

export interface IMaintenanceRefDepot {
	id: number
	depot_code: string
	name: string
}
