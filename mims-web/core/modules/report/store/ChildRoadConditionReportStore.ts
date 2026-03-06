import {
	ChildReportSerivce,
	IReportChildRoadConditionModel,
	IRoadsRoadConditionDepotFilterRoad,
	IRoadsRoadConditionFilter,
	IRoadsRoadConditionFilterRoadGroup,
	IRoadsRoadConditionFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadConditionModel
	dataChildRoadConditionFilter: IRoadsRoadConditionFilter
}

export const useChildRoadConditionReportStore = defineStore("report/child-road-condition--type6", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadConditionModel,
		dataChildRoadConditionFilter: {} as IRoadsRoadConditionFilter,
	}),
	actions: {
		async getChildRoadConditionFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadConditionFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadConditionFilter = res.data
			}
		},
	},
	getters: {
		getKmOptions(state) {
			const options: IOption[] = state.dataChildRoadConditionFilter?.filter_range?.map((e) => {
				return { value: e, label: `${e}` }
			})
			return options
		},
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadConditionFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadConditionDepotFilterRoad) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadConditionDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map((roadItem: IRoadsRoadConditionFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadConditionFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadConditionDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadConditionFilterRoadSections) => {
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
