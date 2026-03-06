export interface IAnnualRoadsTree {
	id: number
	road_number: string
	short_name: string
	road_sections: IAnnualRoadsSectionTree[]
}

export interface IAnnualRoadsSectionTree {
	id: number
	road_group_id: number
	number: string
	name_origin: string
	name_destination: string
	roads: IAnnualRoadsTreeChild[]
}

export interface IAnnualRoadsTreeChild {
	id: number
	name: string
	km_end: number
	km_start: number
	lane_total: number
	ref_direction_id: number
}

export interface IAnnualRoadsTree {
	id: number
	label: string
	children: IAnnualRoadsTreeChild[]
}

export interface IAnnualRoadsTreeChild {
	id: number
	label: string
}

export interface IAnnualStrategicsList {
	id: number
	name: string
	budget: IAnnualStrategicsBudget[]
}

export interface IAnnualStrategicsBudget {
	id: number
	maintenance_analysis_strategic_type_id: number
	name: string
	target: IAnnualStrategicsTarget[]
}

export interface IAnnualStrategicsTarget {
	id: number
	maintenance_analysis_strategic_budget_type_id: number
	name: string
}

export interface IAnnualAnalyzeData {
	id: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	road_group_id: number
	is_iri: boolean
	iri1: number
	iri2: number
	is_aadt: boolean
	aadt1: number
	aadt2: number
	is_age: boolean
	age1: number
	age2: number
	is_gn: boolean
	ifi1: number
	ifi2: number
	is_group_km: boolean
	group_km: number
	percentage: number
	status: number
	is_favorite: boolean
	is_latest: boolean
	is_deleted: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	prepare_data: IAnnualAnalyzePrepareData[]
}

export interface IAnnualAnalyzePrepareData {
	id: number
	maintenance_analysis_id: number
	is_selected: boolean
	group_name: string
	road_name: string
	road_group_id: number
	road_id: number
	lane_no: number
	lane_km_start: number
	lane_km_end: number
	lane_length: number
	lane_width: number
	km_start: number
	km_end: number
	length: number
	area: number
	type: number
	analyst_year: number
	year_road_begin: number
	year_last_overlay: number
	year_last_seal: number
	year_last_mol_rcl: number
	year_last_reconstruction: number
	age: number
	rut: number
	iri: number
	gn: number
	number_of_pothole: number
	area_ac_icrack: number
	percent_ac_icrack: number
	area_ac_ucrack: number
	percent_ac_ucrack: number
	percent_ac_ravelling: number
	cc_transverse_crack: number
	cc_faulting: number
	cc_spalling: number
	current_surface_id: number
	current_surface_name: string
	current_surface_type: string
	current_surface_surface_group: string
	current_surface_layer_coefficient: number
	current_surface_drainage: number
	current_surface_a: number
	current_surface_b: number
	current_surface_c_base: number
	current_surface_c_exp: number
	current_surface_crt: number
	current_surface_rrf: number
	hsold: number
	hsnew: number
	snp_surface: number
	snp_base: number
	snp_subbase: number
	snp: number
	aadt: number
	truck_factor: number
	esal: number
	yax: number
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export interface IAnnualStep2 {
	budget?: number
	comment: string
	condition_id: number
	ifi?: number
	ifi_avg: number | null
	iri?: number
	iri_avg: number | null
	prepare_data_id: number[]
	surface_type: string
	target: number | null
	total_km: number
	name: string
}
export interface IAnnualDefaultDataStep2 {
	name: string
	budget?: number
	comment: string
	condition_id: number
	ifi?: number
	ifi_avg: number | null
	iri?: number
	iri_avg: number | null
	prepare_data_id: number[]
	surface_type: string
	target: number
	total_km: number
	discount: number | null
}

export interface IAnnualAnalyzeDataDefault {
	id: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	roads: IAnnualAnalyzeRoads[]
	condition: number
	is_iri: boolean
	iri1: number | null
	iri2: number | null
	is_aadt: boolean
	aadt1: number | null
	aadt2: number | null
	is_age: boolean
	age1: number | null
	age2: number | null
	is_gn: boolean
	ifi1: number | null
	ifi2: number | null
	is_group_km: boolean
	group_km: number
	percentage: number
	status: number
	is_favorite: boolean
	is_latest: boolean
	is_deleted: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	prepare_data: IAnnualAnalyzePrepareData[]
	name: string
}

export interface IAnnualAnalyzeRoads {
	created_at: string
	created_by: number
	id: number
	maintenance_analysis_id: number
	road_group_id: number
	road_id: number
	updated_at: string
	updated_by: number
}

export interface IAnnualAnalyzePrepareData {
	id: number
	maintenance_analysis_id: number
	is_selected: boolean
	group_name: string
	road_name: string
	road_group_id: number
	road_id: number
	lane_no: number
	lane_km_start: number
	lane_km_end: number
	lane_length: number
	lane_width: number
	km_start: number
	km_end: number
	length: number
	area: number
	type: number
	analyst_year: number
	year_road_begin: number
	year_last_overlay: number
	year_last_seal: number
	year_last_mol_rcl: number
	year_last_reconstruction: number
	age: number
	rut: number
	iri: number
	gn: number
	number_of_pothole: number
	area_ac_icrack: number
	percent_ac_icrack: number
	area_ac_ucrack: number
	percent_ac_ucrack: number
	percent_ac_ravelling: number
	cc_transverse_crack: number
	cc_faulting: number
	cc_spalling: number
	current_surface_id: number
	current_surface_name: string
	current_surface_type: string
	current_surface_surface_group: string
	current_surface_layer_coefficient: number
	current_surface_drainage: number
	current_surface_a: number
	current_surface_b: number
	current_surface_c_base: number
	current_surface_c_exp: number
	current_surface_crt: number
	current_surface_rrf: number
	hsold: number
	hsnew: number
	snp_surface: number
	snp_base: number
	snp_subbase: number
	snp: number
	aadt: number
	truck_factor: number
	esal: number
	yax: number
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export interface IAnnualCopy {
	id: number
	name: string
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	road_group_id: number[]
	iri1: null
	iri2: null
	aadt1: null
	aadt2: null
	ifi1: null
	ifi2: null
	group_km: number
	percentage: number
	status: string
	is_favorite: boolean
	condition: number
	discount: null
	year: number
	target: number
	number_plan: null
	comment: string
	budget: null
	iri: null
	ifi: null
	intervention_criteria_parmas_id: number
	road_work_effect_parmas_id: number
	road_user_cost_parmas_id: number
	deteration_parmas_id: number
	optimization_parmas_id: number
	aadt_parmas_id: number
	filter_data: string
	previous_id: number
	is_latest: boolean
	is_deleted: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	prepare_data_status: boolean
}

export interface IAnnualDashboard {
	road: string[]
	filter: IAnnualDashboardFilter
	condition: IAnnualDashboardCondition
	comment: string
	bar1: IAnnualDashboardBar1
	bar2: IAnnualDashboardBar2
	analyst_year: number
	table: IAnnualDashboardTable
}

export interface IAnnualDashboardBar1 {
	name: string
	lable: string[]
	data: number[]
	area: number[]
	color: string[]
}
export interface IAnnualDashboardBar2 {
	name: string
	lable: string[]
	data: number[]
	budget: number[]
	color: string[]
}

export interface IAnnualDashboardCondition {
	condition: string
	target: string
	discount: number
}

export interface IAnnualDashboardFilter {
	surface_type: string
	lane: string
	km: number
	filter: string[]
}

export interface IAnnualDashboardTable {
	table1: IAnnualDashboardTableTable1
	table2: IAnnualDashboardTableTable2[]
}

export interface IAnnualDashboardTableTable1 {
	budget: number
	iri_after: number
	iri_before: number
}

export interface IAnnualDashboardTableTable2 {
	name: string
	aera: number
	iri_after: number
	range: number
}

export interface IAnnualMapFilter {
	year: number[]
	display: IAnnualMapFilterDisplay[]
	criteria: IAnnualMapFilterCriterion[]
	plan: number[]
	method: IAnnualMapFilterMethod[]
}

export interface IAnnualMapFilterCriterion {
	id: number
	name: string
	grade: IAnnualMapFilterCriterionGrade[]
	grade_cc: IAnnualMapFilterCriterionGrade[]
}

export interface IAnnualMapFilterCriterionGrade {
	name: string
	color: string
	left_value_ac: number
	left_condition_ac: string
	right_value_ac: number
	right_condition_ac: string
	left_value_cc: number
	left_condition_cc: string
	right_value_cc: number
	right_condition_cc: string
}

export interface IAnnualMapFilterDisplay {
	id: number
	name: string
}

export interface IAnnualMapFilterMethod {
	id: number
	name: string
	color: string
}

export interface IAnnualMapData {
	criteria_method: IAnnualMapCriteriaMethod[]
	items: IAnnualMapItem[]
}

export interface IAnnualMapCriteriaMethod {
	name: string
	color: string
}

export interface IAnnualMapItem {
	title: string
	display: number
	road_name: string
	km_start: number
	km_end: number
	the_geom: IAnnualMapGeom
	iri_before: number
	iri_after: number
	year: number
	color: string
}

export interface IAnnualMapGeom {
	type: string
	coordinates: Array<number[]>
}

export interface IAnnualRoadGroup {
	id: number
	road_number: string
	short_name: string
	road_sections: IAnnualRoadSection[]
}

export interface IAnnualRoadSection {
	id: number
	road_group_id: number
	number: string
	name_origin: string
	name_destination: string
	roads: IAnnualRoadOptions[]
}

export interface IAnnualRoadOptions {
	id: number
	name: string
	km_end: number
	km_start: number
	lane_total: number
	ref_direction_id: number
}

export interface IAnnualRoadsTree {
	id: number
	label: string
	children: IAnnualRoadsTreeChild[]
}

export interface IAnnualRoadsTreeChild {
	id: number
	label: string
}

export interface IAnnualInterventionCriteria {
	id: number
	label: string
	children: IAnnualChildCriteria[]
}

export interface IAnnualChildCriteria {
	id: number
	label: string
}

export interface IAnnualAnalyzeData {
	id: number
	name: string
	condition: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	roads: IAnnualAnalyzeRoads[]
	is_iri: boolean
	iri1: number | null
	iri2: number | null
	is_aadt: boolean
	aadt1: number | null
	aadt2: number | null
	is_age: boolean
	age1: number | null
	age2: number | null
	is_gn: boolean
	ifi1: number | null
	ifi2: number | null
	is_group_km: boolean
	group_km: number
	percentage: number
	status: number
	is_favorite: boolean
	is_latest: boolean
	is_deleted: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
	prepare_data: IAnnualAnalyzePrepareData[]
}

export interface IAnnualAnalyzeRoads {
	created_at: string
	created_by: number
	id: number
	maintenance_analysis_id: number
	road_group_id: number
	road_id: number
	updated_at: string
	updated_by: number
}

export interface IAnnualAnalyzePrepareData {
	id: number
	maintenance_analysis_id: number
	is_selected: boolean
	group_name: string
	road_name: string
	road_group_id: number
	road_id: number
	lane_no: number
	lane_km_start: number
	lane_km_end: number
	lane_length: number
	lane_width: number
	km_start: number
	km_end: number
	length: number
	area: number
	type: number
	analyst_year: number
	year_road_begin: number
	year_last_overlay: number
	year_last_seal: number
	year_last_mol_rcl: number
	year_last_reconstruction: number
	age: number
	rut: number
	iri: number
	gn: number
	number_of_pothole: number
	area_ac_icrack: number
	percent_ac_icrack: number
	area_ac_ucrack: number
	percent_ac_ucrack: number
	percent_ac_ravelling: number
	cc_transverse_crack: number
	cc_faulting: number
	cc_spalling: number
	current_surface_id: number
	current_surface_name: string
	current_surface_type: string
	current_surface_surface_group: string
	current_surface_layer_coefficient: number
	current_surface_drainage: number
	current_surface_a: number
	current_surface_b: number
	current_surface_c_base: number
	current_surface_c_exp: number
	current_surface_crt: number
	current_surface_rrf: number
	hsold: number
	hsnew: number
	snp_surface: number
	snp_base: number
	snp_subbase: number
	snp: number
	aadt: number
	truck_factor: number
	esal: number
	yax: number
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export type IAnnualHandleCheckBox = "pie-chart-index" | "tree-map-chart"
