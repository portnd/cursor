import {
	ICompareYearData,
	ICompareLaneData,
	ILane,
	IRoadConditionList,
} from "../infrastructure/DashboardConditionModel"
import { DashboardConditionService } from "../infrastructure/DashboardConditionService"
import { useDashboardConditionStore } from "./index"
import { useDashboardStore } from "~/core/modules/dashboard/store/DashboardStore"

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
	conditionList: IRoadConditionList[]
}

export const useDashboardConditionCompareStore = defineStore("dashboard/conditioncompare", {
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
		conditionList: [],
	}),
	actions: {
		async getConditionType() {},
		async getCompareYear() {
			this.loading = true
			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id[0])

			this.yearParams.condition_type = this.conditionType
			const service = new DashboardConditionService()
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

			this.laneParams.condition_type = this.conditionType
			const service = new DashboardConditionService()
			const res = await service.compareLane(roadId, this.laneParams)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.lane = res.data
			}
		},
		async getParamsData() {},
		async getLaneList() {
			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id)

			const dashboardConditionStore = useDashboardConditionStore()
			this.conditionType = dashboardConditionStore.conditionTypeString

			const service = new DashboardConditionService()
			const res = await service.getLaneLists(roadId)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.laneList = res.data
			}
		},
		async getConditionList() {
			this.loading = true

			const dashboard = useDashboardStore()
			const dashboardConditionStore = useDashboardConditionStore()
			const roadId = Number(dashboard.params.road_id[0])

			if (dashboardConditionStore.conditionType !== 5 && dashboard.params.road_id.length > 0) {
				const service = new DashboardConditionService()
				const res = await service.getConditionList(roadId)

				this.loading = false

				if (!res.status) {
					useHandlerError(res.code, res.error, { showToast: true })
				} else {
					this.conditionList = res.data
				}
			}
		},
	},
	getters: {
		getLaneListOptions(state) {
			const options = state.laneList.map((lane) => ({ label: lane.lane_no.toString(), value: lane.lane_no }))
			return options || []
		},
		getYearOptions(state) {
			const { conditionList } = state

			const options = conditionList?.map((condition) => ({ label: `${condition.year + 543}`, value: condition.year }))

			return options.sort((a, b) => b.value - a.value) || []
		},
	},
})
