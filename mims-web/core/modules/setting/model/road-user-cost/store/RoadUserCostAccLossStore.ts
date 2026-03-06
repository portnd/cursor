import { IRucAccLossData } from "../infrastructure/RoadUserCostAccModel"
import { RoadUserCostAccService } from "../infrastructure/RoadUserCostAccService"
import { IRucAccLossParams } from "../infrastructure/RoadUserCostAccRequest"

interface IState {
	loading: boolean
	data: IRucAccLossData
}

export const useRoadUserCostAccLossStore = defineStore("ruc/loss_value", {
	state: (): IState => ({
		loading: false,
		data: {} as IRucAccLossData,
	}),
	actions: {
		async getLossAccidentData() {
			this.loading = true

			const service = new RoadUserCostAccService()
			const res = await service.getLossValueAccident()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
			}
		},
		async postAccidentLossValue() {
			const params = this.checkParams(this.data)
			this.loading = true

			const service = new RoadUserCostAccService()
			const res = await service.postLossValueAccident(params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		checkParams(params: IRucAccLossParams) {
			const result = {} as IRucAccLossParams
			for (const key in params) {
				if (
					typeof params[key as keyof IRucAccLossParams] !== "number" &&
					params[key as keyof IRucAccLossParams] !== null
				) {
					result[key as keyof IRucAccLossParams] = Number(params[key as keyof IRucAccLossParams])
				} else {
					result[key as keyof IRucAccLossParams] = params[key as keyof IRucAccLossParams]
				}
			}

			return result
		},
	},
	getters: {},
})
