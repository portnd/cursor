import { ICompareYearData, ICompareLaneData } from "../infrastructure/DashboardReflectiveGraphModel.d"
import { DashboardReflectiveGraphService } from "../infrastructure/DashboardReflectiveGraphService"
import { useDashboardStore } from "~/core/modules/dashboard/store/DashboardStore"

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

export const useDashboardReflectiveGraphStore = defineStore("dashboard/reflective/graph", {
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
		async getCompareYear() {
			this.loading = true

			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id[0])

			const service = new DashboardReflectiveGraphService()
			const res = await service.compareYear(roadId, this.yearParams)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.year = res.data
			}
		},
		async getCompareLane() {
			this.loading = true

			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id[0])

			const service = new DashboardReflectiveGraphService()
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
