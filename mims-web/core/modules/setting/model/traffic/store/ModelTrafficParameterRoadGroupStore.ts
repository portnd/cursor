import { ITrafficRoadGroupData, ITrafficVehicleData } from "../infrastructure/ModelTrafficParameterModel"
import { TrafficParameterService } from "../infrastructure/ModelTrafficParameterService"
import { IOption } from "../../../../../shared/types/Option"

interface IState {
	isTruckFactor: boolean
	roadGroupId: number
	sixAxleId: number
	tenAxleId: number
	sixLValue: number
	tenLValue: number
	loading: boolean
	data: ITrafficRoadGroupData[]
	vehicleData: ITrafficVehicleData[]
}

export const useTrafficParameterRoadGroupStore = defineStore("traffic/parameter/road-group-volume", {
	state: (): IState => ({
		isTruckFactor: false,
		roadGroupId: 0,
		sixAxleId: 0,
		tenAxleId: 0,
		sixLValue: 0,
		tenLValue: 0,
		loading: false,
		data: [],
		vehicleData: [],
	}),
	actions: {
		async getAadtRoadGroupData() {
			this.loading = true

			const service = new TrafficParameterService()
			const res = await service.getAadtParameterRoadGroupVolume()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
			}
		},
		async getAadtVehicleType(roadGroupID: number) {
			this.roadGroupId = roadGroupID
			this.loading = true

			const service = new TrafficParameterService()
			const res = await service.getAadtParametersVehicleType(roadGroupID)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.vehicleData = res.data
			}
		},
		setLValue(wheel: number) {
			if (wheel === 6) {
				const six = this.vehicleData.find((item) => item.id === this.sixAxleId)
				this.sixLValue = six?.load_equivalent ? six.load_equivalent : 0
			}

			if (wheel === 10) {
				const ten = this.vehicleData.find((item) => item.id === this.tenAxleId)
				this.tenLValue = ten?.load_equivalent ? ten.load_equivalent : 0
			}
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			if (state.data.length === 0) {
				return []
			}

			const data = state.data.sort((a, b) => a.road_group_id - b.road_group_id)
			const options = data.map((item) => {
				return { label: item.road_group_name, value: item.road_group_id }
			})

			return options
		},
		getVolumeAadt(state) {
			if (state.data.length === 0) {
				return []
			}

			return state.data.find((item) => item.road_group_id === state.roadGroupId)?.volume_aadt
		},
		getCalculate(state) {
			const data = state.data.find((item) => item.road_group_id === state.roadGroupId)

			const calculateData = data?.volume_aadt?.calculate

			return (
				calculateData || {
					four_wheel_total: 0,
					six_to_ten_wheel_percentage: 0,
					six_to_ten_wheel_total: 0,
					ten_wheel_percentage: 0,
					ten_wheel_total: 0,
				}
			)
		},
		getOptions(state) {
			if (state.vehicleData.length === 0) {
				return { six: [], ten: [] }
			}

			const sixOptions: IOption[] = []
			const tenOptions: IOption[] = []

			state.vehicleData.forEach((item) => {
				if (item.num_wheel === 6) {
					sixOptions.push({ label: item.name, value: item.id })
				} else if (item.num_wheel === 10) {
					tenOptions.push({ label: item.name, value: item.id })
				}
			})

			return { six: sixOptions, ten: tenOptions }
		},

		getVehicleImagePath(state) {
			const sixWheelImage = state.vehicleData.find((item) => item.id === state.sixAxleId)?.image_path
			const tenWheelImage = state.vehicleData.find((item) => item.id === state.tenAxleId)?.image_path

			return { six: sixWheelImage || "", ten: tenWheelImage || "" }
		},
		getLoadEquivalent(state) {
			let sixLValue
			let tenLValue
			if (state.sixAxleId || state.tenAxleId) {
				sixLValue = state.vehicleData.find((item) => item.id === state.sixAxleId)?.load_equivalent
				tenLValue = state.vehicleData.find((item) => item.id === state.tenAxleId)?.load_equivalent
			}

			return { six: sixLValue || 0, ten: tenLValue || 0 }
		},
	},
})
