import { IGetTrafficData } from "../infrastructure/ModelTrafficParameterModel"
import { TrafficParameterService } from "../infrastructure/ModelTrafficParameterService"

interface IState {
	loading: boolean
	data: IGetTrafficData
	roadGroupID: number
}

export const useTrafficParameterStore = defineStore("traffic/paramter", {
	state: (): IState => ({
		loading: false,
		data: {} as IGetTrafficData,
		roadGroupID: 0,
	}),
	actions: {
		async getAadtData() {
			this.loading = true

			const service = new TrafficParameterService()
			const res = await service.getAadtParameter(this.roadGroupID)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
			}
		},
	},
	getters: {},
})
