import { ChildReportSerivce, IReportChildAssetModel, IRoadsAssetFilter } from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildAssetModel
	dataChildAssetFilter: IRoadsAssetFilter
}

export const useChildAssetReportStore = defineStore("report/child-asset", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildAssetModel,
		dataChildAssetFilter: {} as IRoadsAssetFilter,
	}),
	actions: {
		async getChildAssetFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildAssetFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildAssetFilter = res.data
			}
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroupList = data.road_group?.map((roadItem) => {
					return { value: roadItem.id, label: roadItem.name }
				})
				options = roadGroupList
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroup = data.road_group.find((item) => item.id === state.data.road_group_id)
				if (roadGroup) {
					const roadSectionsList = roadGroup?.road_section?.map((roadItem) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadSectionsList
				}
			}
			return options
		},
		getAssetOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetFilter?.filter_asset?.find((e) => e.id === state.data.asset_group_id)
			if (data) {
				const roadGroupList = data.asset?.map((assetItem) => {
					return { value: assetItem.id, label: assetItem.name }
				})
				options = roadGroupList
			}
			return options
		},
	},
})
