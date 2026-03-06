import {
	ChildReportSerivce,
	IReportChildRoadDamageModel,
	IReportChildRoadDamageFilter,
	IReportChildRoadDamageFilterDepot,
	IReportChildRoadDamageFilterRoadGroup,
	IReportChildRoadDamageFilterSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadDamageModel
	dataChildRoadDamageFilter: IReportChildRoadDamageFilter
}

export const useChildRoadDamageReportStore = defineStore("report/child-road-damage", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadDamageModel,
		dataChildRoadDamageFilter: {} as IReportChildRoadDamageFilter,
	}),
	actions: {
		async getChildRoadDamageFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getReportRoadDamageFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadDamageFilter = res.data
			}
		},
	},
	getters: {
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadDamageFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IReportChildRoadDamageFilterDepot) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IReportChildRoadDamageFilterDepot) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map((roadItem: IReportChildRoadDamageFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IReportChildRoadDamageFilterDepot) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IReportChildRoadDamageFilterSections) => {
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
