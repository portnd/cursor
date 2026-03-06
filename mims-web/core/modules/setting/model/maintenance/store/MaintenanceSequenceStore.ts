import { MaintenanceSequenceService, IMaintenanceSequence } from "../infrastructure"

interface IState {
	data: IMaintenanceSequence
	loading: boolean
}

export const useMaintenanceSequenceStore = defineStore("setting/model/maintenance/sequence", {
	state: (): IState => ({
		data: {} as IMaintenanceSequence,
		loading: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true
			const maintenanceSequenceService = new MaintenanceSequenceService()
			const res = await maintenanceSequenceService.get()

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				return res
			}
		},
		async post(params: any) {
			// Loading
			this.loading = true

			const maintenanceSequenceService = new MaintenanceSequenceService()
			const res = await maintenanceSequenceService.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			}

			return res
		},
	},
	getters: {},
})
