export interface IVehicleTypeRequest {
	road_group_id: number
	four_wheel: IFourWheel
	six_to_ten_wheel: ISixWheel
	over_ten_wheel: IWheelWheel
}

interface IFourWheelRequest {
	car_less_than_equal_seven: number
	car_over_than_seven: number
	light_bus: number
	light_truck: number
}
interface ISixWheelRequest {
	medium_bus: number
	medium_truck: number
}
interface IWheelWheelRequest {
	heavy_bus: number
	heavy_truck: number
	full_trailor: number
	semi_trailor: number
}
