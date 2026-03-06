import { InitDataService, IInitData } from "../infrastructure"

interface IState {
	data: IInitData | null
	dateTime: Date | null
}

export const useInitDataStore = defineStore("init-data", {
	state: (): IState => ({
		data: null,
		dateTime: null,
	}),
	actions: {
		async initData() {
			const initDataService = new InitDataService()
			const res = await initDataService.initData()

			if (res.status === false) {
				useHandlerError(res.code, res.error)
			} else {
				this.data = res.data
				this.dateTime = new Date()

				console.log("[InitData] Store updated ✅")
			}
		},
	},
	getters: {},
	persist: true,
})
