export interface IRoadSummary {
	id: number
	items: IRoadSummaryItem
	permissions: IRoadSummaryPermissions
	status: string
	status_code: string
	reject_reason: string
	update_by: IRoadSummaryUpdateBy
	update_date: string
}

export interface IRoadSummaryItem {
	id: number
	no?: number
	isNew?: boolean
	km_end: string
	km_start: string
	surface_cross_section_code: number
	width_surface: number | null
	thickness_surface: number | null
	thickness_surface_concrete: number | null
	width_shoulder_left: number | null
	surface_shoulder_left: IRoadItem
	width_shoulder_right: number | null
	surface_shoulder_right: IRoadItem
	thickness_base: number | null
	material_base: IRoadMaterial
	thickness_subbase: number | null
	material_subbase: IRoadMaterial
	thickness_subgrade: number | null
	material_subgrade: IRoadMaterial
	lane_count: number
	lane: IRoadLane[]
	thickness_concrete_slab: number | null
	// no?: number
	// isNew?: boolean
	//
}

export interface IRoadMaterial {
	id: number
	name: string
	is_initial: boolean
	layer_coefficient: number
	drainage: number
	type: string
}

export interface IRoadLane {
	direction: string
	geom_cl: string
	lane_no: number
	surface: IRoadItem
}

interface IRoadSummaryPermissions {
	can_approve: boolean
	can_delete: boolean
	can_edit: boolean
	can_reject: boolean
	can_send: boolean
}

export interface IRoadSummaryUpdateBy {
	id: number
	email: string
	full_name: string
	department: IRoadItem
	profile_picture: string
}

export interface ISurfaceIcon {
	id: number
	name: string
	color_code: string
}

export interface IRoadDetail {
	code: string
	direction: IRoadItem
	id: number
	name: string
	geom_cl: string
	km_end: number
	km_start: number
	lane_count: number
	name: string
	parent_road_id: number
	road_color_code: string
	road_id: number
	road_level: number
	road_type: IRoadItem
	road_type_icon_id: number
	seq: number
}

interface IRoadItem {
	id: number
	name: string
	surface_group: string
	color_code: string
}

export interface ILaneList {
	lane_no: number
	lane_name: string
}

export interface IConditionCompareAverage {
	lane: number
	items: IItemConditionCompareAverage[]
}

export interface IItemConditionCompareAverage {
	year: number
	km_start: number
	km_end: number
	iri: number
	mpd: number
	rut: number
	ifi: number
	gn: number
}

export interface IMaintenanceProjectsData {
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
	color: string
	is_complete: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	last_inspection_date: Date
	guarantee_expiration_date: Date
	budget: IMaintenanceProjectsBudget
	budget_method: IMaintenanceProjectsBudgetMethod
	roads: IMaintenanceProjectsMaintenanceRoad[]
	road_histories: IMaintenanceProjectsMaintenanceRoad[]
	percent_progress: number
	percent_pay: number
}

export interface IMaintenanceProjectsBudget {
	id: number
	name: string
}

export interface IMaintenanceProjectsBudgetMethod {
	id: number
	method_name: string
}

export interface IMaintenanceProjectsMaintenanceRoad {
	id: number
	maintenance_id: number
	road_group_id: number
	road_id: number
	lane_no: number
	km_start: number
	km_end: number
	maintenance_method_id: number
	ref_surface_id: number
	ref_surface_params_id: number
	intervention_criteria_id: number
	intervention_criteria_id_params: number
	the_geom: IMaintenanceProjectGeoms
	road_group: IMaintenanceProjectsRoadGroup
	road_name: string
	color: string
	// road_info: IMaintenanceProjectsRoadInfo
	intervention_criteria: IMaintenanceProjectsInterventionCriteria
	attacchments: IMaintenanceProjectsAttacchment[]
}

interface IMaintenanceProjectGeoms {
	coordinates: number[][]
	type: string
}

export interface IMaintenanceProjectsAttacchment {
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

export interface IMaintenanceProjectsInterventionCriteria {
	id: number
	maintenance_standard_name: string
}

export interface IMaintenanceProjectsRoadGroup {
	id: number
	code: string
	name: string
}

export interface IMaintenanceProjectsRoadInfo {
	id: number
	road_id: number
	ref_direction_id: number
	name: string
	road_color_code: string
	direction: IMaintenanceProjectsDirection
}

export interface IMaintenanceProjectsDirection {
	id: number
	name: string
}
