import { IVehicleType, IVehicleTypeRequest, VehicleTypeService } from "../infrastructure"

interface IState {
	loading: boolean
	id: number
	data: IVehicleType
}

export const useVehicleTypeStore = defineStore("setting/models/aadt/vehicle-type", {
	state: (): IState => ({
		data: {
			road_group_id: 0,
			four_wheel: {
				car_less_than_equal_seven: null,
				car_over_than_seven: null,
				light_bus: null,
				light_truck: null,
			},
			six_to_ten_wheel: {
				medium_bus: null,
				medium_truck: null,
			},
			over_ten_wheel: {
				heavy_bus: null,
				heavy_truck: null,
				full_trailor: null,
				semi_trailor: null,
			},
		} as IVehicleType,
		id: 0,
		loading: false,
	}),
	actions: {
		async get(id: number) {
			if (id !== null) {
				// Loading
				this.loading = true
				const vehicleTypeService = new VehicleTypeService()
				const res = await vehicleTypeService.get(id)
				// Loading
				this.loading = false
				if (res.status === false) {
					useHandlerError(res.code, res.error, { showAlert: true })
				} else {
					this.data = {} as IVehicleType
					this.data = res.data
					return res.data
				}
			}
		},
		async post() {
			// Loading
			this.loading = true
			const params: IVehicleTypeRequest = {
				road_group_id: Number(this.data.road_group_id),
				four_wheel: {
					car_less_than_equal_seven: this.data.four_wheel.car_less_than_equal_seven,
					car_over_than_seven: this.data.four_wheel.car_over_than_seven,
					light_bus: this.data.four_wheel.light_bus,
					light_truck: this.data.four_wheel.light_truck,
				},
				six_to_ten_wheel: {
					medium_bus: this.data.six_to_ten_wheel.medium_bus,
					medium_truck: this.data.six_to_ten_wheel.medium_truck,
				},
				over_ten_wheel: {
					heavy_bus: this.data.over_ten_wheel.heavy_bus,
					heavy_truck: this.data.over_ten_wheel.heavy_truck,
					full_trailor: this.data.over_ten_wheel.full_trailor,
					semi_trailor: this.data.over_ten_wheel.semi_trailor,
				},
			}
			const vehicleTypeService = new VehicleTypeService()
			const res = await vehicleTypeService.post(params)
			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
