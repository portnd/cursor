import { IUpdataAnalys } from "./StrategicModel.d"

export interface IStrategicRoadGroup {
	id: number
	road_number: string
	short_name: string
	road_sections: IStrategicRoadSection[]
}

export interface IStrategicRoadSection {
	id: number
	road_group_id: number
	number: string
	name_origin: string
	name_destination: string
	roads: IStrategicRoadOptions[]
}

export interface IStrategicRoadOptions {
	id: number
	name: string
	km_end: number
	km_start: number
	lane_total: number
	ref_direction_id: number
}

export interface IStrategicRoadsTree {
	id: number
	label: string
	children: IStrategicRoadsTreeChild[]
}

export interface IStrategicRoadsTreeChild {
	id: number
	label: string
}

export interface IStrategicStep1 {
	id: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	road_group_id: number
	iri1: number
	iri2: number
	aadt1: number
	aadt2: number
	age1: number
	age2: number
	gn1: number
	gn2: number
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
	prepare_data: IPrepareData[]
}

export interface IPrepareData {
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

export interface IStrategicStep2 {
	budget: number
	comment: string
	condition_id: number
	discount: number
	gn?: number
	ifi_avg: number
	iri?: number
	iri_avg: number
	number_plan: number
	plans: Partial<IStrategicStep2plan>[]
	prepare_data_id: number[]
	surface_type: string
	target: number
	total_km: number
	year: number
}

export interface IStrategicStep2plan {
	id: number
	isNew: boolean
	plan_1: number
	plan_2: number
	plan_3: number
	plan_year: number
}

export interface IStrategicAnalyzeData {
	id: number
	name: string
	condition: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	roads: IStrategicAnalyzeRoads[]
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
	prepare_data: IStrategicAnalyzePrepareData[]
}

export interface IStrategicAnalyzeRoads {
	created_at: string
	created_by: number
	id: number
	maintenance_analysis_id: number
	road_group_id: number
	road_id: number
	updated_at: string
	updated_by: number
}

export interface IStrategicAnalyzePrepareData {
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

export interface IUpdataAnalyseResponseData {
	id: number
	maintenance_analysis_type_id: number
	roads: IUpdataAnalyseRoad[]
	surface_type_id: number
	lane_type_id: number
	iri1: number
	iri2: number
	aadt1: number
	aadt2: number
	age1: number
	age2: number
	gn1: number
	gn2: number
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
	prepare_data: IUpdataAnalysePrepareData[]
}

export interface IUpdataAnalysePrepareData {
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

export interface IUpdataAnalyseRoad {
	id: number
	maintenance_analysis_id: number
	road_group_id: number
	road_id: number
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export interface ICopy {
	oldId?: number
	id: number
	maintenance_analysis_type_id: number
	surface_type_id: number
	lane_type_id: number
	road_group_id: number[]
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
	gn1: number
	gn2: number
	is_group_km: boolean
	group_km: number
	percentage: number
	status: number
	is_favorite: boolean
	condition: number
	discount: number
	year: number
	target: number
	number_plan: number
	comment: string
	budget: null
	iri: null
	gn: null
	is_latest: boolean
	is_deleted: boolean
	created_by: number
	updated_by: number
	created_at: Date
	updated_at: Date
}

export interface IStrategicInterventionCriteria {
	id: number
	label: string
	children: IStrategicChildCriteria[]
}

export interface IStrategicChildCriteria {
	id: number
	label: string
}

export interface IStrategicDashboard {
	road: any[]
	number_plan: number
	filter: IStrategicDashboardFilter
	condition: IStrategicDashboardCondition
	comment: string
	graph1: IStrategicDashboardGraph1
	bar1: IStrategicDashboardBar1
	bar2: IStrategicDashboardBar2
	table: IStrategicDashboardTable
}

export interface IStrategicDashboardBar1 {
	name: string
	lable: number[]
	datasets: IStrategicDashboardBar1Dataset[]
	color: string[]
}

export interface IStrategicDashboardBar1Dataset {
	lable: string
	value: number[]
}

export interface IStrategicDashboardBar2 {
	name: string
	lable: number[]
	datasets: IStrategicDashboardBar2Dataset[]
}

export interface IStrategicDashboardBar2Dataset {
	plan: string
	data: IStrategicDashboardDataset[]
}

export interface IStrategicDashboardDataset {
	lable: string[]
	value: number[]
	budget: number[]
}

export interface IStrategicDashboardCondition {
	condition: string
	target: string
	discount: null
}

export interface IStrategicDashboardFilter {
	surface_type: string
	lane: string
	km: number
	filter: string[]
}

export interface IStrategicDashboardGraph1 {
	name: string
	lable: number[]
	value: Array<number[]>
	line: string[]
	color: string[]
}

export interface IStrategicDashboardTable {
	summary: IStrategicDashboardTableSummary[]
	plan_1: IStrategicDashboardTableUnlimitedPlan[]
	plan_2: IStrategicDashboardTableUnlimitedPlan[]
	plan_3: IStrategicDashboardTableUnlimitedPlan[]
	unlimited_plan: IStrategicDashboardTableUnlimitedPlan[]
}

export interface IStrategicDashboardTableSummary {
	name: string
	data: IStrategicDashboardTableSummaryDatum[]
}

export interface IStrategicDashboardTableSummaryDatum {
	name: string
	value: number[]
}

export interface IStrategicDashboardTableUnlimitedPlan {
	method_name: string
	data: IStrategicDashboardTableUnlimitedPlanDatum[]
}

export interface IStrategicDashboardTableUnlimitedPlanDatum {
	year: number
	km: number
	budget: number
}

export interface IStrtegicMapFilter {
	year: number[]
	display: IStrtegicMapFilterDisplay[]
	criteria: IStrtegicMapFilterCriterion[]
	plan: IStrtegicMapFilterDisplay[]
	method: IStrtegicMapFilterMethod[]
}

export interface IStrtegicMapFilterCriterion {
	id: number
	name: string
	grade: IStrtegicMapFilterCriterionGrade[]
	grade_cc: IStrtegicMapFilterCriterionGrade[]
}

export interface IStrtegicMapFilterCriterionGrade {
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

export interface IStrtegicMapFilterDisplay {
	id: number
	name: string
}

export interface IStrtegicMapFilterMethod {
	id: number
	name: string
	color: string
}

export interface IStrategicMapData {
	criteria_method: IStrategicMapDataCriteriaMethod[]
	items: IStrategicMapDataItem[]
}

export interface IStrategicMapDataCriteriaMethod {
	name: string
	color: string
}

export interface IStrategicMapDataItem {
	title: string
	display: number
	road_name: RoadName
	km_start: number
	km_end: number
	the_geom: IStrategicMapDataGeom
	iri_before: number
	iri_after: number
	year: number
	color: string
}

export interface IStrategicMapDataGeom {
	type: Type
	coordinates: Array<number[]>
}
