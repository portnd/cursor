import {
	ChildReportSerivce,
	IReportChildAssetMapModel,
	IRoadsAssetMapFilter,
	IRoadsAssetMapFilterRoadGroup,
	IRoadsAssetMapFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildAssetMapModel
	dataChildAssetMapFilter: IRoadsAssetMapFilter
}

export const useChildAssetMapReportStore = defineStore("report/child-asset-map", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildAssetMapModel,
		dataChildAssetMapFilter: {} as IRoadsAssetMapFilter,
	}),
	actions: {
		async getChildAssetMapFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildAssetMapFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildAssetMapFilter = res.data
			}
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetMapFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroupList = data.road_group?.map((roadItem: IRoadsAssetMapFilterRoadGroup) => {
					return { value: roadItem.id, label: roadItem.name }
				})
				options = roadGroupList
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetMapFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroup = data.road_group.find(
					(item: IRoadsAssetMapFilterRoadGroup) => item.id === state.data.road_group_id
				)
				if (roadGroup) {
					const roadSectionsList = roadGroup?.road_section?.map((roadItem: IRoadsAssetMapFilterRoadSections) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadSectionsList
				}
			}
			return options
		},
		getAssetOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildAssetMapFilter?.filter_asset?.find((e) => e.id === state.data.asset_group_id)
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
