import { IOptimization, IOptimizationRequest, OptizationService } from "../infrastructure"

interface IState {
	data: IOptimization
	loading: boolean
}

export const useOptimizationStore = defineStore("setting/model/optimization", {
	state: (): IState => ({
		data: {} as IOptimization,
		loading: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true
			const optizationService = new OptizationService()
			const res = await optizationService.get()

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				return res
			}
		},
		async post() {
			// Loading
			this.loading = true

			const params: IOptimizationRequest = {
				bc_ratio_constraint: Number(this.data.bc_ratio_constraint),
				default_design_life: Number(this.data.default_design_life),
			}

			const optizationService = new OptizationService()
			const res = await optizationService.post(params)

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
