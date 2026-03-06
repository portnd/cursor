import { IRequestUsers, IUserRolesItem, UsersService } from "../infrastructure"
import { IUsersDepartmentItems } from "../infrastructure/UsersModel"
import { IFile } from "~/core/shared/types/File"

interface IStateParams {
	created_by: number | null

	email: string | null
	firstname: string | null
	lastname: string | null
	profile_img_part: string | null
	updated_by: number | null
	tel: string | null
	role: number[]
	status: number
	username: string
	password: string
	ref_user_owner_id: number | null
	ref_depot_id: number | null
}

interface IState {
	loading: boolean
	file: IFile
	filePath: string
	params: IStateParams
	roles: IUserRolesItem[]
	department: IUsersDepartmentItems[]
}

export const useUsersCreateStore = defineStore("setting/users/create", {
	state: (): IState => ({
		loading: false,
		file: {} as IFile,
		filePath: "",
		params: {
			created_by: null,

			email: null,
			firstname: null,
			lastname: null,
			profile_img_part: null,
			role: [],
			updated_by: null,
			tel: "",
			status: 0,
			username: "",
			password: "",
			ref_depot_id: null,
			ref_user_owner_id: null,
		},
		roles: [],
		department: [],
	}),
	actions: {
		// async getDepartments() {
		// 	this.loading = true

		// 	const service = new UsersService()
		// 	const res = await service.getDepartmentsList()

		// 	if (!res.status) {
		// 		useHandlerError(res.code, res.error, { showToast: true })
		// 	} else {
		// 		this.department = res.data?.items

		// 		await this.getRoles()
		// 	}

		// 	this.loading = false
		// },
		async getRoles() {
			const service = new UsersService()
			const res = await service.getRolesList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roles = res.data?.items
				this.params.status = 1
			}
		},
		async createUser() {
			// Loading
			this.loading = true

			const params = this.generateParams()

			const usersService = new UsersService()
			const res = await usersService.createUser(params)

			// Loading
			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				// this.name = ""
				return res
			}
		},
		generateParams() {
			const newParams: IRequestUsers = {
				created_by: this.params.created_by ? this.params.created_by : 0,
				ref_user_owner_id: this.params.ref_user_owner_id,
				ref_depot_id: this.params.ref_depot_id,
				email: this.params.email!,
				firstname: this.params.firstname!,
				lastname: this.params.lastname!,
				profile_img_path: this.file?.data?.base64 ? this.file?.data?.base64 : "",
				updated_by: this.params.updated_by ? this.params.updated_by : 0,
				roles: this.params.role,
				tel: this.params.tel ? this.params.tel : "",
				status: this.params.status === 1,
				username: this.params.username,
				password: this.params.password,
			}

			return newParams
		},
	},
	getters: {
		// getDepartmentOptions(state) {
		// 	const departments = state.department
		// 	const options = [{ label: "ไม่มีฝ่าย", value: 0 }]

		// 	departments.forEach((item) => {
		// 		options.push({ label: item.name, value: item.id })
		// 	})

		// 	return options || []
		// },
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
					return e.id === state.params.ref_user_owner_id
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
