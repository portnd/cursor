import { IRucAccChanceData } from "../infrastructure/RoadUserCostAccModel"
import { RoadUserCostAccService } from "../infrastructure/RoadUserCostAccService"
import { IRucAccChanceParams } from "../infrastructure/RoadUserCostAccRequest"
import { IRoadData } from "~/core/modules/road/roadList/infrastructure"

interface IState {
	loading: boolean
	data: IRucAccChanceData
	roadGroup: IRoadData[]
	roadGroupId: number
	params: IRucAccChanceParams
}

export const useRoadUserCostAccChanceStore = defineStore("ruc/chance_of_accident", {
	state: (): IState => ({
		loading: false,
		data: {} as IRucAccChanceData,
		roadGroup: [],
		roadGroupId: 0,
		params: {} as IRucAccChanceParams,
	}),
	actions: {
		async getAccChanceData(roadGroupId: number) {
			const service = new RoadUserCostAccService()
			const res = await service.getChanceOfAccident(roadGroupId)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
				this.params = this.data
			}
		},
		async postAccidentChanceParams() {
			this.loading = true
			this.data.road_group_id = this.roadGroupId
			const params = this.checkParams(this.data)

			const service = new RoadUserCostAccService()
			const res = await service.postChanceOfAccident(params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		checkParams(params: IRucAccChanceParams) {
			const result = {} as IRucAccChanceParams
			for (const key in params) {
				if (
					typeof params[key as keyof IRucAccChanceParams] !== "number" &&
					params[key as keyof IRucAccChanceParams] !== null
				) {
					result[key as keyof IRucAccChanceParams] = Number(params[key as keyof IRucAccChanceParams])
				} else {
					result[key as keyof IRucAccChanceParams] = params[key as keyof IRucAccChanceParams]
				}
			}

			return result
		},
	},
	getters: {},
})
