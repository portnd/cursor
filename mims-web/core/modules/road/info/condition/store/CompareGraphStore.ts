import { ICompareYearData, ICompareLaneData, ILane } from "../infrastructure/RoadConditionModel"
import { RoadConditionService } from "../infrastructure/RoadConditionService"

interface Params {
	yearsInput: number[]
	laneInput: number[]
	years: number[]
	lanes: number[]
	condition_type: string
}

interface IState {
	roadID: number
	loading: boolean
	conditionType: string
	lane: ICompareLaneData[]
	year: ICompareYearData[]
	laneParams: Params
	yearParams: Params
	laneList: ILane[]
	min: number
	max: number
}

export const useGraphCompareStore = defineStore("condition/compare", {
	state: (): IState => ({
		roadID: 1,
		loading: false,
		conditionType: "",
		lane: [],
		year: [],
		laneList: [],
		laneParams: {
			yearsInput: [],
			laneInput: [],
			years: [],
			lanes: [],
			condition_type: "",
		},
		yearParams: {
			yearsInput: [],
			laneInput: [],
			years: [],
			lanes: [],
			condition_type: "",
		},
		min: 0,
		max: 0,
	}),
	actions: {
		async getCompareYear(roadId: number) {
			this.loading = true

			this.yearParams.condition_type = this.conditionType
			const service = new RoadConditionService()
			const res = await service.compareYear(roadId, this.yearParams)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.year = res.data
			}
		},
		async getCompareLane(roadId: number) {
			this.loading = true
			this.laneParams.condition_type = this.conditionType
			const service = new RoadConditionService()
			const res = await service.compareLane(roadId, this.laneParams)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.lane = res.data
			}
		},
		async getLaneList(id: number) {
			const service = new RoadConditionService()
			const res = await service.getLaneList(id)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.laneList = res.data
			}
		},
	},
	getters: {
		getLaneListOptions(state) {
			const options = state.laneList.map((lane) => ({ label: lane.lane_no.toString(), value: lane.lane_no }))
			return options || []
		},
	},
})
