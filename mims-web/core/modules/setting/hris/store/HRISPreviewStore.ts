import { HRISService, IHRISPreview } from "../infrastructure"

interface IState {
	data: IHRISPreview
	loading: boolean
	importing: boolean
}

export const useHRISPreviewStore = defineStore("setting/preview/edit", {
	state: (): IState => ({
		data: { road_group: [], road_section: [] } as IHRISPreview,

		loading: false,
		importing: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true

			const service = new HRISService()
			const res = await service.preview()

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data

				return res
			}
		},
		async import() {
			// Loading
			this.importing = true

			const service = new HRISService()
			const res = await service.import()

			// Loading
			this.importing = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
