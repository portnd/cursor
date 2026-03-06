import {
	ChildReportSerivce,
	IReportChildRoadConditionSummaryModel,
	IRoadsRoadConditionSummaryDepotFilterRoad,
	IRoadsRoadConditionSummaryFilter,
	IRoadsRoadConditionSummaryFilterRoadGroup,
	IRoadsRoadConditionSummaryFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadConditionSummaryModel
	dataChildRoadConditionSummaryFilter: IRoadsRoadConditionSummaryFilter
}

export const useChildRoadConditionSummaryReportStore = defineStore("report/child-road-condition-type7", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadConditionSummaryModel,
		dataChildRoadConditionSummaryFilter: {} as IRoadsRoadConditionSummaryFilter,
	}),
	actions: {
		async getChildRoadConditionSummaryFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadConditionSummaryFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadConditionSummaryFilter = res.data
			}
		},
	},
	getters: {
		getConditionOptions(state) {
			const options: IOption[] = state.dataChildRoadConditionSummaryFilter?.filter_condition?.map((e) => {
				return { value: e, label: `${e}` }
			})
			return options
		},
		getCriteriaOptions(state) {
			const options: IOption[] = state.dataChildRoadConditionSummaryFilter?.filter_criteria?.map((e) => {
				return { value: e.id, label: `${e.name}` }
			})
			return options
		},
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadConditionSummaryFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadConditionSummaryDepotFilterRoad) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadConditionSummaryDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map(
						(roadItem: IRoadsRoadConditionSummaryFilterRoadGroup) => {
							return { value: roadItem.id, label: roadItem.name }
						}
					)
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadConditionSummaryDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadConditionSummaryFilterRoadSections) => {
								return { value: roadItem.id, label: roadItem.name }
							}
						)
						options = roadSectionsList
					}
				}
			}
			return options
		},
	},
})
