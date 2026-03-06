export interface IVehicleType {
	road_group_id: number
	four_wheel: IFourWheel
	six_to_ten_wheel: ISixWheel
	over_ten_wheel: IWheelWheel
}

interface IFourWheel {
	car_less_than_equal_seven: number | null
	car_over_than_seven: number | null
	light_bus: number | null
	light_truck: number | null
}
interface ISixWheel {
	medium_bus: number | null
	medium_truck: number | null
}
interface IWheelWheel {
	heavy_bus: number | null
	heavy_truck: number | null
	full_trailor: number | null
	semi_trailor: number | null
}
