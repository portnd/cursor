import { IRequestLogin, LoginService } from "../infrastructure"

interface IState {
	username: string
	password: string
	loading: boolean
}

export const useLoginStore = defineStore("login", {
	state: (): IState => ({
		username: "",
		password: "",
		loading: false,
	}),
	actions: {
		async login() {
			// Loading
			this.loading = true

			const params: IRequestLogin = {
				username: this.username,
				password: this.password,
			}

			const loginService = new LoginService()
			const res = await loginService.login(params)

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
