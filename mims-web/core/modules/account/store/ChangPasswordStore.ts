import { ChangePasswordService, IRequestChangePassword } from "../infrastructure"

interface IStateParams {
	current_password: string
	new_password: string
	confirm_new_password: string
}

interface IState {
	loading: boolean
	params: IStateParams
}

export const useChangePasswordStore = defineStore("change-password", {
	state: (): IState => ({
		loading: false,
		params: {
			current_password: "",
			new_password: "",
			confirm_new_password: "",
		},
	}),
	actions: {
		async updatePassword() {
			this.loading = true

			const params: IRequestChangePassword = {
				current_password: this.params.current_password,
				new_password: this.params.new_password,
				confirm_new_password: this.params.confirm_new_password,
			}

			const service = new ChangePasswordService()
			const res = await service.updatePassword(params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
