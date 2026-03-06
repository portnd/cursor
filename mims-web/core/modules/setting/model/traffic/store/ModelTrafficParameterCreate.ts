import { ITrafficCalculate } from "../infrastructure/ModelTrafficParameterModel"
import { ITrafficParams } from "../infrastructure/ModelTrafficParameterRequest"
import { TrafficParameterService } from "../infrastructure/ModelTrafficParameterService"

interface IState {
	loading: boolean
	params: ITrafficParams
}

interface ILoadEquivalent {
	six: number
	ten: number
}

export const useTrafficParameterCreateStore = defineStore("traffic/parameter/create", {
	state: (): IState => ({
		loading: false,
		params: {
			directional_distribution_factor: 0,
			elane: 0,
			four_wheel_axle_number: 0,
			four_wheel_vehicle_volume: 0,
			is_truck_factor: false,
			lane_distribution_factor: 0,
			road_group_id: 0,
			six_wheel_axle_number_id: 0,
			six_wheel_factor_result: 0,
			six_wheel_percentage_truck: 0,
			six_wheel_vehicle_volume: 0,
			speed_average: 0,
			speed_heavy_truck: 0,
			ten_wheel_axle_number_id: 0,
			ten_wheel_factor_result: 0,
			ten_wheel_percentage_truck: 0,
			ten_wheel_vehicle_volume: 0,
		},
	}),
	actions: {
		async postAadtParameter() {
			const param = this.checkParams(this.params)
			this.loading = true

			const service = new TrafficParameterService()
			const res = await service.postAadtParameter(param)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		checkParams(params: ITrafficParams) {
			for (const key in params) {
				if (typeof params[key as keyof ITrafficParams] === "string") {
					// @ts-ignore
					params[key as keyof ITrafficParams] = Number(params[key as keyof ITrafficParams])
				}
			}
			return params
		},
		updateParams(calculateData: ITrafficCalculate) {
			this.params.four_wheel_vehicle_volume = calculateData.four_wheel_total
			this.params.six_wheel_vehicle_volume = calculateData.six_to_ten_wheel_total
			this.params.ten_wheel_vehicle_volume = calculateData.ten_wheel_total
			this.params.six_wheel_percentage_truck = calculateData.six_to_ten_wheel_percentage
			this.params.ten_wheel_percentage_truck = calculateData.ten_wheel_percentage
		},
		calculateTruckFactor(loadEquivalent: ILoadEquivalent) {
			this.params.six_wheel_factor_result = Number((this.params.six_wheel_percentage_truck / 100) * loadEquivalent.six)
			this.params.ten_wheel_factor_result = Number((this.params.ten_wheel_percentage_truck / 100) * loadEquivalent.ten)
		},
	},
	getters: {
		getVehiclePercentage(state) {
			const sumVehicle = state.params.six_wheel_vehicle_volume + state.params.ten_wheel_vehicle_volume
			const sixTruckPercent = state.params.six_wheel_vehicle_volume / sumVehicle
			const tenTruckPercent = state.params.ten_wheel_vehicle_volume / sumVehicle

			return {
				six: sixTruckPercent ? Number(sixTruckPercent) : 0,
				ten: tenTruckPercent ? Number(tenTruckPercent) : 0,
			}
		},
	},
})
