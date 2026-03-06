import {
	IDefaultUsersData,
	IRequestUsers,
	IUserRolesItem,
	IUsersDepartmentItems,
	UsersService,
} from "../infrastructure"
import { IFile } from "~/core/shared/types/File"

interface IStateParams {
	status: number
	password: string
	roles: number[]
}

interface IState {
	loading: boolean
	data: IDefaultUsersData
	file: IFile
	roles: IUserRolesItem[]
	department: IUsersDepartmentItems[]
	params: IStateParams
}

export const useUsersEditStore = defineStore("setting/users/edit", {
	state: (): IState => ({
		loading: false,
		file: {} as IFile,
		data: {} as IDefaultUsersData,
		roles: [],
		department: [],
		params: {
			status: 0,
			password: "",
			roles: [],
		},
	}),
	actions: {
		async getDefaultUsers(id: number) {
			this.loading = true

			const service = new UsersService()
			const res = await service.getDefaultUser(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data

				this.params.status = this.data.status ? 1 : 2
				this.params.roles = this.data.roles.filter((item) => item.is_checked).map((item) => item.id)

				await this.getRoles()
			}
			this.loading = false
		},
		// async getDepartments() {
		// 	const service = new UsersService()
		// 	const res = await service.getDepartmentsList()

		// 	if (!res.status) {
		// 		useHandlerError(res.code, res.error, { showToast: true })
		// 	} else {
		// 		this.department = res.data?.items

		// 		await this.getRoles()
		// 	}
		// },
		async getRoles() {
			const service = new UsersService()
			const res = await service.getRolesList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roles = res.data?.items
			}
		},
		async updateUser(id: number) {
			this.loading = true
			const params = this.generateParams()

			const service = new UsersService()
			const res = await service.updateUser(id, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		updateStatus() {
			this.data.status = this.params.status === 1
		},
		generateParams() {
			const newParams: IRequestUsers = {
				created_by: this.data.created_by ? this.data.created_by : 0,
				ref_user_owner_id: this.data.ref_user_owner_id,
				ref_depot_id: this.data.ref_depot_id,
				email: this.data.email!,
				firstname: this.data.firstname!,
				lastname: this.data.lastname!,
				profile_img_path: this.file?.data?.base64 ? this.file?.data?.base64 : "",
				updated_by: this.data.updated_by ? this.data.updated_by : 0,
				roles: this.params.roles,
				tel: this.data.tel ? this.data.tel : "",
				status: this.params.status === 1,
				username: this.data.username,
				password: this.params.password,
			}

			return newParams
		},
	},
	getters: {
		getDepartmentOptions(state) {
			const departments = state.department
			const options = [{ label: "ไม่มีฝ่าย", value: 0 }]

			departments.forEach((item) => {
				options.push({ label: item.name, value: item.id })
			})

			return options || []
		},
		getUserOwnersOption() {
			const options = useInitData()
				?.refUserOwner()
				?.map((e: any) => {
					return { label: e.name, value: e.id }
				})

			return options || []
		},
		getDepotOption(state) {
			const owner = useInitData()
				?.refUserOwner()
				?.find((e) => {
					return e.id === state.data.ref_user_owner_id
				})

			console.log("owner =", owner)

			if (!owner) {
				return []
			}

			const options = owner.ref_depot?.map((e: any) => {
				return { label: e.name, value: e.id }
			})

			return options || []
		},
	},
})
