import { IRequestAssetGroup, AssetGroupService } from "../infrastructure"

interface IState {
	id: number
	name: string
	loading: boolean
}

export const useAssetGroupEditStore = defineStore("setting/asset-group/edit", {
	state: (): IState => ({
		id: 0,
		name: "",
		loading: false,
	}),
	actions: {
		async get(id: number) {
			this.id = id
			// Loading
			this.loading = true

			const assetGroupService = new AssetGroupService()
			const res = await assetGroupService.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.name = res.data.name
				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			const params: IRequestAssetGroup = {
				name: this.name,
			}

			const assetGroupService = new AssetGroupService()
			const res = await assetGroupService.put(this.id, params)

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
