export interface ILine {
	line_no: number
	// lane_name: string
}

export interface IReflectiveListData {
	year: number
	items: Item[]
}

export interface Item {
	id: number
	id_parent: number
	direction: IDirection
	line_no: number
	surveyed_date: Date
}

export interface IDirection {
	id: number
	name: string
}

export interface IReflectivityItem {
	km_start: number
	km_end: number
	retro_avg: number
	geom_cl: string
	ref_stripe_color_id: number
	ref_stripe_color: IReflectivityRefStripe
	ref_stripe_type_id: number
	ref_stripe_type: IReflectivityRefStripe
	items: IReflectivityItem[]
}

export interface IReflectivityDetails {
	id: number
	id_parent: number
	updated_date: string
	updated_by: IReflectivityUpdatedBy
	status: string
	direction: IReflectivityDirection
	has_white_line: boolean
	has_yellow_line: boolean
	datas: IReflectivityData[]
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

export interface IRoadConditionList {
	year: number
	items: IRoadConditionListItem[]
}

export interface IRoadConditionListItem {
	id: number
	id_parent: number
	direction: IRoadConditionDirection
	lane_no: number
	surveyed_date: string
}

export interface IRoadConditionDirection {
	id: number
	name: string
}

export interface IRoadConditionDetails {
	id: number
	id_parent: number
	updated_date: Date
	updated_by: IRoadConditionUpdatedBy
	status: string
	direction: IRoadConditionDirection
	condition_datas: IRoadConditionData[]
}

export interface IRoadConditionData {
	condition_type: string
	items: IRoadConditionItem[]
}

export interface ICompareYearData {
	lane: number
	items: IYearItem[]
}

export interface IYearItem {
	km_end: number
	km_start: number
	value: number
	year: number
}

export interface ICompareLaneData {
	year: number
	items: ILaneItem[]
}

export interface ILaneItem {
	lane_no: number
	km_start: number
	km_end: number
	value: number
}

export interface ILane {
	lane_no: number
	lane_name: string
}

export interface IDashboardCondition {
	chart: IDashboardConditionChart
	has_mutiple_road: boolean
	table: IDashboardConditionTable[]
}

export interface IDashboardConditionChart {
	// color: string[]
	// data: number[]
	// lable: string[]
	color: Array[]
	data: Array[]
	lable: Array[]
	name: string
	road_id: string[]
}

export interface IDashboardConditionTable {
	avg_value: number
	detail_km: IDashboardConditionTableDetail[]
	lane_no: number
	total_km: number
}

export interface IDashboardConditionTableDetail {
	ref_grade_id: number
	ref_grade_name: string
	value: number
	value_percent: number
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

// export interface IDashboardMap {
// 	items: IDashboardMapItems
// }

export interface IDashboardMapItems {
	color: string
	the_geom: IDashboardMapItemsGeom
}

export interface IDashboardMapItemsGeom {
	coordinates: number[][]
	type: string
}
