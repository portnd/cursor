import {
	ChildReportSerivce,
	IReportChildRoadDamageSummaryModel,
	IRoadsRoadDamageSummaryFilter,
	IRoadsRoadDamageSummaryFilterDepot,
	IRoadsRoadDamageSummaryFilterRoadGroup,
	IRoadsRoadDamageSummaryFilterSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadDamageSummaryModel
	dataChildRoadDamageSummaryFilter: IRoadsRoadDamageSummaryFilter
}

export const useChildRoadDamageSummaryReportStore = defineStore("report/child-road-damage-summary", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadDamageSummaryModel,
		dataChildRoadDamageSummaryFilter: {} as IRoadsRoadDamageSummaryFilter,
	}),
	actions: {
		async getChildRoadDamageSummaryFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getReportRoadDamageSummaryFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadDamageSummaryFilter = res.data
			}
		},
	},
	getters: {
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadDamageSummaryFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadDamageSummaryFilterDepot) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find((item) => item.id === state.data.department_id)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map((roadItem: IRoadsRoadDamageSummaryFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadDamageSummaryFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find((item) => item.id === state.data.department_id)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadDamageSummaryFilterSections) => {
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
