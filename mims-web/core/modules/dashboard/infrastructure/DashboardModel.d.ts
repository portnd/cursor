export interface IDashboardSurface {
	surface_dashboard_table: ISurfaceDetail[]
	geom_list: ISurfaceGeomList[]
	summary: ISurfaceSummary[]
}

export interface ISurfaceDetail {
	surface_name: string
	surface_lane_type: ISurfaceLane
}

export interface ISurfaceLane {
	one_lane: number
	two_lane: number
	three_lane: number
	four_lane: number
	more_than_four: number
}

export interface ISurfaceId {
	id: number
	name: string
	color_code: string
}

export interface ISurfaceGeomList {
	code: string
	geom_cl: Array<String>
	id: number
	surface: ISurfaceId
}

export interface ISurfaceSummary {
	surface: ISurfaceId
	value: number
	road_id: string
}

export interface ISurfaceMap {
	title: string
	color: string
	road_group_name: string
	contract_number: string
	year: string
	last_inspection_date: string
	road_name: string
	km_start: number
	km_end: number
	km_total: number
	surface_name: string
	ref_surface_id: number
	age: number
	the_geom: ITheGeomSurfaceMap
}

export interface ITheGeomSurfaceMap {
	coordinates: Array<number[]>
	type: string
}

export interface TheGeom {
	type: string
	coordinates: Array<number[] | number>
}

export interface IDataMart {
	id: number
	road_id: number
	road_surface_id: number
	lane_count: number
	surface_year: number
	year: number

	contract_number: string

	km_start: number
	km_end: number
	lane_no: number
	ref_surface_id: number
	age: number
	last_inspection_date: string
	the_geom: string
	created_by: string
	created_at: string
	updated_by: string
	updated_at: string
}

export interface IDataMartCheck {
	stauts: boolean
	percent: number
	updated_by: string
	updated_at: string
}

export interface IHighChart {
	title: IHighChartTitle
	chart: Chart
	dataLabels: Credits
	stroke: Stroke
	xAxis: XAxis
	yAxis: YAxis
	tooltip: Tooltip
	plotOptions: PlotOptions
	grid: Grid
	legend: Credits
	fill: Fill
	credits: Credits
	series: Array<IHighChartSeries>
}

export interface IHighChartTitle {
	text: string
}

export interface IHighChartSeries {
	name: string
	data: number[]
	zoneAxis: string
	zones: Zones[]
}

export interface IMapBound {
	maxLat: number
	maxLon: number
	minLat: number
	minLon: number
}

export interface IDashboardAsset {
	id: number
	name: string
	icon_filepath: string
	thumbnail_icon_filepath: string
	line_color_code: string
	is_cluster: boolean
	cluster: number
	the_geom: TheGeom
	detail: string
	asset_table_id: number
}

export interface IDashboardAssetGeom {
	type: string
	coordinates: Array<number[] | number>
}

export interface IDashboardAssetData {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IDashboardAssetItems[]
}

export interface IDashboardAssetItems {
	id: number
	name: string
	label: string[]
	data: number[]
}

export interface IDashboardAssetDetail {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IDashboardAssetDetailItem[]
}

export interface IDashboardAssetDetailItem {
	asset_group: IDasboardAssetDetailGroup
	asset_list: IDasboardAssetDetailList[]
}

export interface IDasboardAssetDetailGroup {
	id: number
	name: string
}

export interface IDasboardAssetDetailList {
	asset: IDasboardAssetDetailAsset
	is_range: boolean
	value: number
}

export interface IDasboardAssetDetailAsset {
	default_color: string
	default_icon_url: string
	id: number
	name: string
	thumbnail_icon_url: string
}

export interface IDashboardStateAssetDetail {
	id: number
	name: string
	is_active: boolean
	colors: string
	items: IDashboardStateAssetItems[]
}

export interface IDashboardStateAssetItems {
	id: number
	name: string
	icon_path: string
	thumbnail_path: string
	is_active: boolean
	value: number
	color: string
	is_range: boolean
}

export interface IDashboardRoadOptions {
	id: number
	road_number: string
	road_sections: IDashboardRoadSection[]
	short_name: string
}

export interface IDashboardRoadSection {
	id: number
	name_destination: string
	name_origin: string
	number: string
	road_group_id: number
	roads: IDashboardRoads[]
}

export interface IDashboardRoads {
	id: number
	km_end: number
	km_start: number
	lane_total: number
	ref_direction_id: number
	name: string
}
export interface IDashboardRoad {
	road: IDashboardRoadData
	length_roads: IDashboardRoadLength[]
	aadt_roads: IDashboardRoadAadt[]
}

export interface IDashboardRoadAadt {
	name: string
	aadt: number
	year1: number
	year2: number
	percent: number
	growth_rate: string
}

export interface IDashboardRoadLength {
	name: string
	total: number
	asphalt: number
	concrete: number
}

export interface IDashboardRoadData {
	name: string
	label: string[]
	data: number[]
	color: string[]
	road_group_id: number[]
}

export interface IDashboardSurface {
	summary: IDashboardSurfaceChart[]
	surface_dashboard_table: IDashboardSurfaceTable[]
}

export interface IDashboardSurfaceChart {
	surface: IDashboardSurfaceData
	value: number
}

export interface IDashboardSurfaceData {
	id: number
	name: string
	color_code: string
}

export interface IDashboardSurfaceTable {
	id: number
	surface_name: string
	surface_lane_type: IDashboardSurfaceLaneType
}

export interface IDashboardSurfaceLaneType {
	one_lane: number
	two_lane: number
	three_lane: number
	four_lane: number
	more_than_four: number
}

export enum IDashboardMenu {
	Asset = "asset",
	Condition = "condition",
	Surface = "surface",
	Maintenance = "maintenance",
}

export interface IDashboardSurfaceMap {
	title: string
	color: string
	road_group_name: string
	contract_number: string
	year: string
	last_inspection_date: string
	road_name: string
	km_start: number
	km_end: number
	km_total: number
	surface_name: string
	ref_surface_id: number
	age: number
	the_geom: IDashboardSurfaceMapGeoms
}

export interface IDashboardSurfaceMapGeoms {
	type: string
	coordinates: Array<number[]>
}

export interface IDashboardMaintenance {
	updated_at: Date
	number_maintenance_chart: IDashboardMaintenanceChart
	maintenance_budget_chart: IDashboardMaintenanceChart
	top_ten_maintenance_budget_chart: IDashboardTopTenMaintenanceChart
	// table: IDashboardMaintenanceTable
}

export interface IDashboardMaintenanceChart {
	name: string
	lable: string[]
	data: number[]
	color: string[]
	road_group_id: string[]
}
export interface IDashboardTopTenMaintenanceChart {
	name: string
	lable: string[]
	data: number[]
	color: string[]
	maintenance_id: number[]
}

export interface IDashboardMaintenanceTable {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IDashboardMaintenanceTableItem[]
}

export interface IDashboardMaintenanceTableItem {
	id: number
	contract_number: string
	road_name: string[]
	section_name: string[]
	ref_depot_name: string[]
	budget: number
	guarantee_expiration_date: string
	remain_date: number
}

export interface IDashboardMaintenanceMap {
	id_parent: number
	title: string
	road_name: string
	section_name: string
	ref_depot_name: string
	contract_number: string
	name: string
	km_start: string
	km_end: string
	km_total: number
	lane_no: number
	color: string
	the_geom: IDashboardMaintenanceMapGeom
}

export interface IDashboardMaintenanceMapGeom {
	coordinates: number[][]
	type: string
}

export interface IDashboardConditionMap {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IDashboardMapItems[]
}

export interface IDashboardMapItems {
	color: string
	the_geom: IDashboardMapItemsGeom
}

export interface IDashboardMapItemsGeom {
	coordinates: number[][]
	type: string
}

export type IHandleCheckboxMode = "maintenance-quantity" | "maintenance-budget"
