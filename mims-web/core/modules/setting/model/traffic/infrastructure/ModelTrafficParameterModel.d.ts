export interface ITrafficRoadGroup {
	status: boolean
	code: number
	data: ITrafficRoadGroupData[]
}

export interface ITrafficRoadGroupData {
	road_group_id: number
	road_group_name: string
	volume_aadt: IVolumeTraffic
}

export interface IVolumeTraffic {
	Veh2: number
	Veh3: number
	Veh4: number
	calculate: ITrafficCalculate
	veh1: number
}

export interface ITrafficCalculate {
	four_wheel_total: number
	six_to_ten_wheel_percentage: number
	six_to_ten_wheel_total: number
	ten_wheel_percentage: number
	ten_wheel_total: number
}

export interface IGetTraffic {
	status: boolean
	code: number
	data: IGetTrafficData
}

export interface IGetTrafficData {
	road_group_id: number
	elane: number
	four_wheel_axle_number: number
	four_wheel_vehicle_volume: number
	six_wheel_axle_number_id: number
	six_wheel_vehicle_volume: number
	six_wheel_percentage_truck: number
	six_wheel_factor_result: number
	ten_wheel_axle_number_id: number
	ten_wheel_vehicle_volume: number
	ten_wheel_percentage_truck: number
	ten_wheel_factor_result: number
	is_truck_factor: boolean
	speed_average: number
	speed_heavy_truck: number
	lane_distribution_factor: number
	directional_distribution_factor: number
}

export interface ITrafficVehicle {
	status: boolean
	code: number
	data: ITrafficVehicleData[]
}

export interface ITrafficVehicleData {
	id: number
	num_wheel: number
	name: string
	num_axle: number
	load_equivalent: number
	image_path: string
}
