export interface IStrategicCreatePrepareDataReq {
	aadt1: number | null
	aadt2: number | null
	ifi1: number | null
	ifi2: number | null
	group_km: number | null
	iri1: number | null
	iri2: number | null
	lane_type_id: number | null
	maintenance_analysis_type_id: number
	roads: number[]
	surface_type_id: number | null
	name: string
}

export interface IStrategicCreateAnalyzeParams {
	budget?: number
	comment: string
	condition_id: number
	discount: number
	gn?: number
	iri?: number
	number_plan: number
	plans: IStrategicCreateParamsPlan[]
	prepare_data_id: number[]
	surface_type: string
	target: number
	year: number
	name: string
}
export interface IStrategicUpdateAnalyzeParams {
	budget?: number
	comment: string
	condition_id: number
	discount: number
	gn?: number
	iri?: number
	number_plan: number
	plans: IStrategicCreateParamsPlan[]
	prepare_data_id: number[]
	surface_type: string
	target: number
	year: number
	name: string
}
export interface IStrategicParams {
	budget?: number
	comment: string
	condition: number
	discount: number
	gn?: number
	iri?: number
	number_plan: number
	plans?: IStrategicParamsPlan[]
	prepare_data_id: number[]
	surface_type: string
	target: number
	year: number
}

export interface IStrategicCreateParamsPlan {
	id?: number | null
	plan_1?: number
	plan_2?: number
	plan_3?: number
	plan_year?: number
}
export interface IStrategicParamsPlan {
	id: number
	plan_1: number
	plan_2: number
	plan_3: number
	plan_year: number
}

export interface IStrategicUpdatePrepareData {
	aadt1: number | null
	aadt2: number | null
	ifi1: number | null
	ifi2: number | null
	group_km: number
	iri1: number | null
	iri2: number | null
	lane_type_id: number
	maintenance_analysis_type_id: number
	roads: number[]
	surface_type_id: number
	name: string
}

export interface IReportParams {
	type: string
	plan?: string
}

export interface ISearchModelReq {
	aadt1: number | null
	aadt2: number | null
	age1: number | null
	age2: number | null
	group_km: number | null
	ifi1: number | null
	ifi2: number | null
	iri1: number | null
	iri2: number | null
	lane_type_id: number | null
	maintenance_analysis_type_id: number | null
	name: string
	roads: number[]
	surface_type_id: number | null
	intervention_criteria_id: string | null
}

export interface IMapDataReq {
	plan: number
	year: number
	display: number
	criteria?: number
	method?: number
}

export interface IUpdateModelReq {
	id: number
	intervention_criteria_id: string
}
