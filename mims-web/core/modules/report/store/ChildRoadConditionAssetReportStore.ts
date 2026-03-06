import {
	ChildReportSerivce,
	IReportChildRoadConditionAssetModel,
	IRoadsRoadConditionAssetFilter,
	IRoadsRoadConditionAssetFilterRoadGroup,
	IRoadsRoadConditionAssetFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadConditionAssetModel
	dataChildRoadConditionAssetFilter: IRoadsRoadConditionAssetFilter
}

export const useChildRoadConditionAssetReportStore = defineStore("report/child-road-condition-asset", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadConditionAssetModel,
		dataChildRoadConditionAssetFilter: {} as IRoadsRoadConditionAssetFilter,
	}),
	actions: {
		async getChildRoadConditionAssetFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadConditionAssetFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadConditionAssetFilter = res.data
			}
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionAssetFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroupList = data.road_group?.map((roadItem: IRoadsRoadConditionAssetFilterRoadGroup) => {
					return { value: roadItem.id, label: roadItem.name }
				})
				options = roadGroupList
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionAssetFilter?.filter_road?.find((e) => e.id === state.data.department_id)
			if (data) {
				const roadGroup = data.road_group?.find(
					(item: IRoadsRoadConditionAssetFilterRoadGroup) => item.id === state.data.road_group_id
				)
				if (roadGroup) {
					const roadSectionsList = roadGroup?.road_section?.map(
						(roadItem: IRoadsRoadConditionAssetFilterRoadSections) => {
							return { value: roadItem.id, label: roadItem.name }
						}
					)
					options = roadSectionsList
				}
			}
			return options
		},
	},
})
