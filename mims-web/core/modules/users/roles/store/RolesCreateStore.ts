import {
	IRequestUpdateRoles,
	RolesService,
	IAccessRoles,
	IRequestRolesAccessDetail,
	IRolesAccessGroup,
	IRolesAccessControl,
	IRolesAccessDetail,
	IAccessRolesAccessControl,
} from "../infrastructure"

interface IState {
	name: string
	control: IAccessRoles[]
	controlFlat: IAccessRolesAccessControl[]
	currentItemChecked?: IAccessRolesAccessControl
	loading: boolean
}
export const useRolesCreateStore = defineStore("user/roles/create", {
	state: (): IState => ({
		name: "",
		control: [],
		controlFlat: [],
		currentItemChecked: undefined,
		loading: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true

			const rolesService = new RolesService()
			const res = await rolesService.getAccess()

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.control = res.data
				this.controlFlat = this.control
					.map((summary) => summary.menu.flatMap((menu) => menu.access_control))
					.reduce((acc, val) => acc.concat(val), [])

				return res
			}
		},
		async create() {
			// Loading
			this.loading = true

			const result: IRequestRolesAccessDetail[] = this.control
				.flatMap((el: IRolesAccessGroup) => el.menu)
				.flatMap((el: IRolesAccessControl) => el.access_control)
				.filter((el: IRolesAccessDetail) => el.is_check === true)
				.map((el: IRolesAccessDetail) => ({ access_control_id: el.id }))

			const params: IRequestUpdateRoles = {
				name: this.name,
				access_control: result,
			}

			const roleService = new RolesService()
			const res = await roleService.createRoles(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		clearInput() {
			this.name = ""
		},
		updateCheckedByRelation(isChecked: boolean) {
			const relationAccess: string[] = []

			this.controlFlat
				.filter((e) => e.is_check === isChecked)
				.forEach((acc) => {
					const relationAcc = getPermissionRelation(acc.access_key, isChecked)
					relationAccess.push(...relationAcc)
				})

			let viewOwnerMaintenanceAnalysis = false
			let viewAllMaintenanceAnalysis = false
			let manageMyselfMaintenanceAnalysis = false

			this.controlFlat.forEach((acc) => {
				if (acc.access_key === IUserRolesAccess.view_all_maintenance_analysis) {
					viewAllMaintenanceAnalysis = acc.is_check
				}
				if (acc.access_key === IUserRolesAccess.view_owner_maintenance_analysis) {
					viewOwnerMaintenanceAnalysis = acc.is_check
				}
				if (acc.access_key === IUserRolesAccess.manage_myself_maintenance_analysis) {
					manageMyselfMaintenanceAnalysis = acc.is_check
				}
			})

			if (isChecked) {
				if (manageMyselfMaintenanceAnalysis && !viewAllMaintenanceAnalysis && !viewOwnerMaintenanceAnalysis) {
					relationAccess.push(IUserRolesAccess.view_owner_maintenance_analysis)
				}
			} else if (!viewAllMaintenanceAnalysis && !viewOwnerMaintenanceAnalysis) {
				relationAccess.push(IUserRolesAccess.manage_myself_maintenance_analysis)
			}

			this.controlFlat.forEach((acc) => {
				if (relationAccess.includes(acc.access_key)) {
					acc.is_check = isChecked
				}
			})
		},
		updateCurrentItemChecked(item: IAccessRolesAccessControl) {
			this.currentItemChecked = item
		},
	},
	getters: {},
})
