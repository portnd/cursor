import { InitMenuService, IInitMenu } from "../infrastructure"

interface IState {
	data: IInitMenu[] | null
	dateTime: Date | null
}

export const useInitMenuStore = defineStore("init-menu", {
	state: (): IState => ({
		data: null,
		dateTime: null,
	}),
	actions: {
		async initMenu() {
			const initMenuService = new InitMenuService()
			const res = await initMenuService.initMenu()

			if (res.status === false) {
				useHandlerError(res.code, res.error)
			} else {
				this.data = res.data
				this.dateTime = new Date()

				console.log("[Menu] Store updated ✅")
			}
		},
	},
	getters: {},
	persist: true,
})
