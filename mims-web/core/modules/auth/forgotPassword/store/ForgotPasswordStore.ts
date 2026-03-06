import { IRequestForgotPassword, ForgotPasswordService } from "../infrastructure"

interface IState {
	email: string
	callbackUrl: string
	loading: boolean
}

export const useForgotPasswordStore = defineStore("forgotPassword", {
	state: (): IState => ({
		email: "",
		callbackUrl: useRuntimeConfig().siteUrl + "/auth/reset-password",
		loading: false,
	}),
	actions: {
		async forgotPassword() {
			// Loading
			this.loading = true

			const params: IRequestForgotPassword = {
				email: this.email,
				callback_url: this.callbackUrl,
			}

			const forgotPasswordService = new ForgotPasswordService()
			const res = await forgotPasswordService.forgotPassword(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return true
			}
		},
	},
	getters: {},
})
