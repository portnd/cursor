import { InitUserService, IUser } from "../infrastructure"

interface IAccessPermissions {
	[key: string]: boolean
}
interface IState {
	data: IUser | null
	userAccessKey: string[]
	accessPermissions: IAccessPermissions
}

export const useInitUserStore = defineStore("init-user", {
	state: (): IState => ({
		data: null,
		userAccessKey: [],
		accessPermissions: {},
	}),
	actions: {
		async initUser() {
			const initUserService = new InitUserService()
			const res = await initUserService.initUser()

			if (res.status === false) {
				useHandlerError(res.code, res.error)
			} else {
				// mock
				// res.data.ref_owner = {
				// 	name: "องค์การบริหารส่วนจังหวัดขอนแก่น",
				// 	prefix: "Khonkaen-ams",
				// 	application_name_th: "ระบบบริหารจัดการสินทรัพย์องค์การบริหารส่วนจังหวัดขอนแก่น",
				// 	application_name_en: "Khonkaen-AMS",
				// 	logo: "https://bma-ams-api-stg.infra-corp.co/storages/branch/logo/5_Khonkaen-ams",
				// 	color_code: "#048c46",
				// } as IRefOwner

				this.$patch((state) => {
					state.data = res.data
					state.userAccessKey = res.data.access_control.map((e) => e.access_key)

					const obj: { [key: string]: boolean } = {}
					Object.values(IUserRolesAccess).forEach((e) => {
						obj[e] = state.userAccessKey.includes(e)
					})
					state.accessPermissions = obj
				})

				console.log("[User] Store this.accessPermissions ✅ :")
				console.log("[User] Store updated ✅ : ", this.data)
				console.log("[User] Store updated ✅")
			}
		},
		getIsOwnerManagePermission(isOwnerAccess: boolean, depotId: number) {
			console.log("[User] Store depotId ", depotId)
			console.log("[User] Store this.data?.ref_depot_id ", this.data?.ref_depot_id)
			return isOwnerAccess && depotId === this.data?.ref_depot_id
		},
		clearUser() {
			this.$patch((state) => {
				state.data = null
				state.userAccessKey = []
				state.accessPermissions = {}
				console.log("[User] Store clear ✅")
			})
		},
	},
	getters: {},
	persist: true,
})
