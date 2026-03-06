import {
	ChildReportSerivce,
	IReportChildRoadReflectLightSummaryModel,
	IRoadsRoadReflectLightSummaryDepotFilterRoad,
	IRoadsRoadReflectLightSummaryFilter,
	IRoadsRoadReflectLightSummaryFilterRoadGroup,
	IRoadsRoadReflectLightSummaryFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadReflectLightSummaryModel
	dataChildRoadReflectLightSummaryFilter: IRoadsRoadReflectLightSummaryFilter
}

export const useChildRoadReflectLightSummaryReportStore = defineStore("report/child-road-reflect-light-summary", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadReflectLightSummaryModel,
		dataChildRoadReflectLightSummaryFilter: {} as IRoadsRoadReflectLightSummaryFilter,
	}),
	actions: {
		async getChildRoadReflectLightSummaryFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadReflectLightSummaryFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadReflectLightSummaryFilter = res.data
			}
		},
	},
	getters: {
		getCriteriaOptions(state) {
			const options: IOption[] = state.dataChildRoadReflectLightSummaryFilter?.filter_criteria?.map((e) => {
				return { value: e.id, label: `${e.name}` }
			})
			return options
		},
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadReflectLightSummaryFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadReflectLightSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadReflectLightSummaryDepotFilterRoad) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadReflectLightSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadReflectLightSummaryDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map(
						(roadItem: IRoadsRoadReflectLightSummaryFilterRoadGroup) => {
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
			const data = state.dataChildRoadReflectLightSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadReflectLightSummaryDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadReflectLightSummaryFilterRoadSections) => {
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
