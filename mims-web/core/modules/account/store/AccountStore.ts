import { AccountService, IAccount, IRequestAccount } from "../infrastructure"
import { IFile } from "~/core/shared/types/File"

interface IStateParams {
	file: IFile
	file_path: string
	roles: number[]
}

interface IState {
	loading: boolean
	account: IAccount
	params: IStateParams
}

export const useAccountStore = defineStore("account", {
	state: (): IState => ({
		loading: false,
		account: { ref_user_owner: { email: "" }, ref_depot: { name: "" } } as IAccount,
		params: {
			file: {} as IFile,
			file_path: "",
			roles: [],
		},
	}),
	actions: {
		async getAccount() {
			this.loading = true

			const service = new AccountService()
			const res = await service.getAccountData()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.account = res.data
				this.params.file_path = this.account.profile_img_path
				this.params.roles = this.account.roles.map((item) => item.id)
			}
		},
		async updateAccount() {
			this.loading = true

			const params: IRequestAccount = {
				firstname: this.account.firstname,
				lastname: this.account.lastname,
				email: this.account.email,
				profile_img_path: this.params.file.data?.base64 || "",
				department_id: this.account.department_id,
				tel: this.account.tel,
			}

			const service = new AccountService()
			const res = await service.updateAccount(params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {
		getDepartment(state) {
			const data = state.account
			const department = data.department
			const departmentName = department?.name

			return departmentName || ""
		},
	},
})
