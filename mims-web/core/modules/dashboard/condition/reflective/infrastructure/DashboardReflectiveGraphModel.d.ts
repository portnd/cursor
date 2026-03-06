import { IReflective } from "~core/modules/setting/surveyRule/infrastructure/SurveyRuleModel"
export interface IRoadReflective {
	status: boolean
	code: number
	data: IReflectiveData
}

export interface IReflectiveData {
	id: number
	id_parent: number
	line_no: number
	surveyed_date: string
	remarks: string
	csv_file: string
	img_filepath: string
	direction: IDirection
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

export interface IRoadLaneList {
	status: boolean
	code: number
	data: ILane[]
}

export interface ILine {
	line_no: number
	// lane_name: string
}

export interface IReflectiveGraph {
	status: boolean
	code: number
	data: IGraphData
}

export interface IGraphData {
	id: number
	id_parent: number
	updated_date: Date
	updated_by: IGraphUpdatedBy
	status: string
	status_code: string
	reject_reason: string
	permissions: IPermissions
	direction: IDirection
	road_type_id: number
	items: IGraphDataItem[]
}

export interface IDirection {
	id: number
	name: string
}

export interface IGraphDataItem {
	km_start: number
	km_end: number
	value: number
	items: IGraphItem[]
}

export interface IGraphItem {
	km_start: number
	km_end: number
	value: number
	grade: number
	geom_cl: string
	img_filepath: string
}

export interface IPermissions {
	can_edit: boolean
	can_delete: boolean
	can_approve: boolean
	can_send: boolean
	can_reject: boolean
}

export interface IGraphUpdatedBy {
	uid: string
	user_name: string
	full_name: string
	department: IDirection
	profile_picture: string
}

export interface IGraphCompareYear {
	status: boolean
	code: number
	data: ICompareYearData[]
}

export interface ICompareYearData {
	line: number
	items: IYearItem[]
}

export interface IYearItem {
	km_end: number
	km_start: number
	value: number
	year: number
}

export interface IGraphCompareLane {
	status: boolean
	code: number
	data: ICompareLaneData[]
}

export interface ICompareLaneData {
	year: number
	items: ILineItem[]
}

export interface ILineItem {
	line_no: number
	km_start: number
	km_end: number
	value: number
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

export interface IHighChartSeries {
	name: string
	data: number[]
	zoneAxis: string
	zones: Zones[]
}

export interface Zones {
	value: number
	color: string
}

export interface Chart {
	width?: number
	height: number
	type: string
	zoomType: string
	events: any
	stacked: boolean
	toolbar: Toolbar
}

export interface Toolbar {
	show: boolean
}

export interface Credits {
	enabled: boolean
}

export interface Fill {
	opacity: number
}

export interface Grid {
	borderColor: string
	strokeDashArray: number
	xaxis: Axis
	yaxis: Axis
}

export interface Axis {
	lines: Toolbar
}

export interface PlotOptions {
	series: Series
}

export interface Series {
	turboThreshold: number
	marker: Credits
}

export interface Stroke {
	curve: string
}

export interface IHighChartTitle {
	text: string
}

export interface Tooltip {
	useHTML: boolean
	borderRadius: number
	backgroundColor: string
	borderColor: string
	padding: number
	formatter: any
}

export interface XAxis {
	categories: any[]
	tickInterval: number
	margin: number
	title: XAxisTitle
	labels: Labels
}

export interface Labels {
	enable: boolean
	style: Style
}

export interface Style {
	fontSize: string
	color: string
}

export interface XAxisTitle {
	text: null | string
	style: Style
}

export interface YAxis {
	title: XAxisTitle
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

export interface IReflectivityData {
	color: string
	items: IReflectivityItem[]
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

export interface IReflectivityRefStripe {
	id: number
	name: string
	name_th: string
}

export interface IReflectivityDirection {
	id: number
	name: string
}

export interface IReflectivityUpdatedBy {
	uid: string
	user_name: string
	full_name: string
	department: IReflectivityDirection
	profile_picture: string
}
