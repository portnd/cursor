import { ICompareYearData, ICompareLaneData } from "../infrastructure/RoadReflectiveModel"
import { RoadReflectiveService } from "../infrastructure/RoadReflectiveService"

interface Params {
	yearsInput: Array<number>
	lineInput: Array<number>
	years: Array<number>
	lines: Array<number>
}

interface IState {
	roadID: number
	loading: boolean
	lane: ICompareLaneData[]
	year: ICompareYearData[]
	laneParams: Params
	yearParams: Params
	min: number
	max: number
}

export const useReflectGraphCompareStore = defineStore("reflectivity/compare", {
	state: (): IState => ({
		roadID: 1,
		loading: false,
		lane: [],
		year: [],
		laneParams: {
			yearsInput: [],
			lineInput: [],
			years: [],
			lines: [],
		},
		yearParams: {
			yearsInput: [],
			lineInput: [],
			years: [],
			lines: [],
		},
		min: 0,
		max: 0,
	}),
	actions: {
		async getCompareYear(roadId: number) {
			this.loading = true
			const service = new RoadReflectiveService()
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
			const service = new RoadReflectiveService()
			const res = await service.compareLane(roadId, this.laneParams)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.lane = res.data
			}
		},
	},
	getters: {},
})
