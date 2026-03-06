export interface IRucParentParams {
	car_less_than_equal_seven: IRucChildParams
	car_over_than_seven: IRucChildParams
	full_trailor: IRucChildParams
	heavy_bus: IRucChildParams
	heavy_truck: IRucChildParams
	light_bus: IRucChildParams
	light_truck: IRucChildParams
	medium_bus: IRucChildParams
	medium_track: IRucChildParams
	semi_trailor: IRucChildParams
}

export interface IRucChildParams {
	[key: string]: number | string
}
