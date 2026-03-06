import { IRequestAssetGroup, AssetGroupService } from "../infrastructure"

interface IState {
	name: string
	loading: boolean
}

export const useAssetGroupCreateStore = defineStore("setting/asset-group/create", {
	state: (): IState => ({
		name: "",
		loading: false,
	}),
	actions: {
		async create() {
			// Loading
			this.loading = true

			const params: IRequestAssetGroup = {
				name: this.name,
			}

			const assetgroupService = new AssetGroupService()
			const res = await assetgroupService.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.name = ""
				return res
			}
		},
	},
	getters: {},
})
