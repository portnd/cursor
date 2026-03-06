export interface IRucList {
	status: boolean
	code: number
	data: IRucListData[]
}

export interface IRucListData {
	id: number
	name: string
	name_en: string
}

export interface IRucData {
	status: boolean
	code: number
	data: IRucDataTable
}

export interface IRucDataTable {
	car_less_than_equal_seven: IRucItem
	car_over_than_seven: IRucItem
	full_trailor: IRucItem
	heavy_bus: IRucItem
	heavy_truck: IRucItem
	light_bus: IRucItem
	light_truck: IRucItem
	medium_bus: IRucItem
	medium_track: IRucItem
	semi_trailor: IRucItem
}

export interface IRucItem {
	vehicle_name?: string
	[key: string]: number
}
