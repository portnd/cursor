import { LogoutService } from "../infrastructure"

export const useLogoutStore = defineStore("logout", {
	state: () => ({}),
	actions: {
		async logout() {
			const logoutService = new LogoutService()
			const res = await logoutService.logout()

			if (res.status === false) {
				useHandlerError(res.code, res.error)
			} else {
				return true
			}
		},
	},
	getters: {},
})
