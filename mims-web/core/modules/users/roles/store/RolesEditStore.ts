import { IAccessRolesAccessControl } from "../infrastructure/RolesAccessModel"
import {
	IRequestUpdateRoles,
	RolesService,
	IRequestRolesAccessDetail,
	IRolesAccessGroup,
	IRolesAccessControl,
	IRolesAccessDetail,
} from "../infrastructure"

interface IState {
	id: number
	name: string
	access_group: IRolesAccessGroup[]
	loading: boolean
	controlFlat: IAccessRolesAccessControl[]
	currentItemChecked?: IAccessRolesAccessControl
}

export const useRolesEditStore = defineStore("user/roles/edit", {
	state: (): IState => ({
		id: 0,
		name: "",
		access_group: [],
		loading: false,
		controlFlat: [],
		currentItemChecked: undefined,
	}),
	actions: {
		async get(id: number) {
			this.id = id

			// Loading
			this.loading = true

			const rolesService = new RolesService()
			const res = await rolesService.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.name = res.data.role
				this.access_group = res.data.access_group
				this.controlFlat = this.access_group
					.map((summary) => summary.menu.flatMap((menu) => menu.access_control))
					.reduce((acc, val) => acc.concat(val), [])

				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			const result: IRequestRolesAccessDetail[] = this.access_group
				.flatMap((el: IRolesAccessGroup) => el.menu)
				.flatMap((el: IRolesAccessControl) => el.access_control)
				.filter((el: IRolesAccessDetail) => el.is_check === true)
				.map((el: IRolesAccessDetail) => ({ access_control_id: el.id }))

			const params: IRequestUpdateRoles = {
				name: this.name,
				access_control: result,
			}

			const rolesService = new RolesService()
			const res = await rolesService.put(this.id, params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
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

			if (!isChecked) {
				if (!viewAllMaintenanceAnalysis && !viewOwnerMaintenanceAnalysis) {
					relationAccess.push(IUserRolesAccess.manage_myself_maintenance_analysis)
				}
			} else if (manageMyselfMaintenanceAnalysis && !viewAllMaintenanceAnalysis && !viewOwnerMaintenanceAnalysis) {
				relationAccess.push(IUserRolesAccess.view_owner_maintenance_analysis)
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
