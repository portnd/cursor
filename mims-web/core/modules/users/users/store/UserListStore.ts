import { IUserRolesItem, IUsersDepartmentItems, UsersService } from "../infrastructure"
import { IRequestUserSearch } from "../infrastructure/UsersRequest"

interface IStateParams {
	fullname: string
	username: string

	ref_user_owner_id: number | null
	ref_depot_id: number | null
	permission: string
	status: number | null
}

interface IState {
	loading: boolean
	roles: IUserRolesItem[]
	department: IUsersDepartmentItems[]
	params: IStateParams
}

export const useUserListStore = defineStore("setting/user", {
	state: (): IState => ({
		loading: false,
		roles: [],
		department: [],
		params: {
			fullname: "",
			username: "",

			ref_user_owner_id: null,
			ref_depot_id: null,
			permission: "",
			status: null,
		},
	}),
	actions: {
		async getRoles() {
			const service = new UsersService()
			const res = await service.getRolesList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roles = res.data?.items
				// this.params.status = 1
			}
		},
		updateParams() {
			let status: boolean | string = ""

			if (typeof this.params.status === "number") {
				if (this.params.status === 1) {
					status = true
				} else if (this.params.status === 2) {
					status = false
				}
			}

			const newParams: IRequestUserSearch = {
				fullname: this.params.fullname,
				username: this.params.username,
				ref_user_owner_id: this.params.ref_user_owner_id,
				ref_depot_id: this.params.ref_user_owner_id === 3 ? this.params.ref_depot_id : null,
				permission: this.params.permission ?? "",
				status,
			}

			return newParams
		},
	},
	getters: {
		getUserOwnersOption() {
			const options = useInitData()
				?.refUserOwner()
				?.map((e: any) => {
					return { label: e.name, value: e.id }
				})

			return options || []
		},
		getDepotOption() {
			const owner = useInitData()
				?.refUserOwner()
				?.find((e) => {
					return e.id === 3
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
