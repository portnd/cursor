import {
	IRequestCheckPasswordToken,
	IRequestResetPassword,
	ResetPasswordService,
	CheckResetPasswordTokenService,
} from "../infrastructure"

interface IState {
	email: string
	resetPasswordToken: string
	newPassword: string
	confirmNewPassword: string
	loading: boolean
	status: boolean
	statusResetPasswordToken: boolean | null
}

export const useResetPasswordStore = defineStore("reset-password", {
	state: (): IState => ({
		email: "",
		resetPasswordToken: "",
		newPassword: "",
		confirmNewPassword: "",
		loading: false,
		status: false,
		statusResetPasswordToken: null,
	}),
	actions: {
		setResetPasswordToken(token: string) {
			this.resetPasswordToken = token
		},
		async checkToken() {
			// Loading
			this.loading = true

			const params: IRequestCheckPasswordToken = {
				reset_password_token: this.resetPasswordToken,
			}

			const resetPasswordService = new CheckResetPasswordTokenService()
			const res = await resetPasswordService.checkResetPasswordToken(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error)
			} else {
				this.email = res.data.email
				this.statusResetPasswordToken = true
				return true
			}
		},
		async resetPassword() {
			// Loading
			this.loading = true

			const params: IRequestResetPassword = {
				reset_password_token: this.resetPasswordToken,
				new_password: this.newPassword,
				confirm_new_password: this.confirmNewPassword,
			}

			const resetPasswordService = new ResetPasswordService()
			const res = await resetPasswordService.resetPassword(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				const self = this
				useHandlerError(res.code, res.error, {
					showAlert: true,
					fn: function () {
						self.status = false
					},
				})
			} else {
				this.status = true
				return true
			}
		},
	},
	getters: {},
})
