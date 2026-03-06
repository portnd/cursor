import {
	ChildReportSerivce,
	IReportMaintenanceModel,
	IReportMaintenanceFilter,
	IReportMaintenanceFilterSections,
	IReportMaintenanceFilterCondition,
	IReportMaintenanceFilterYear,
	IReportMaintenanceFilterRoadGroup,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportMaintenanceModel
	dataReportMaintenanceFilter: IReportMaintenanceFilter
}

export const useReportMaintenanceReportStore = defineStore("report/report-maintenance", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportMaintenanceModel,
		dataReportMaintenanceFilter: {} as IReportMaintenanceFilter,
	}),
	actions: {
		async getReportMaintenanceFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getReportMaintenanceFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataReportMaintenanceFilter = res.data
			}
		},
	},
	getters: {
		getConditionOptions(state) {
			let options: IOption[] = []
			const condition = state.dataReportMaintenanceFilter?.filter_condition?.map(
				(conditionItem: IReportMaintenanceFilterCondition) => {
					return { value: conditionItem.name, label: conditionItem.name }
				}
			)
			options = condition
			return options
		},
		getYearOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportMaintenanceFilter?.filter_condition?.find(
				(e) => e.name === state.data.condition_name
			)
			if (data) {
				const year = data?.year?.map((yearItem: IReportMaintenanceFilterYear) => {
					return { value: yearItem.year, label: `${yearItem.year + 543}` }
				})
				options = year
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportMaintenanceFilter?.filter_condition?.find(
				(e) => e.name === state.data.condition_name
			)
			if (data) {
				const year = data?.year?.find((e) => e.year === state.data.year)
				if (year) {
					const roadGroupList = year.road_group?.map((roadItem: IReportMaintenanceFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportMaintenanceFilter?.filter_condition?.find(
				(e) => e.name === state.data.condition_name
			)
			if (data) {
				const year = data?.year?.find((e) => e.year === state.data.year)
				if (year) {
					const roadGroupList = year.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map((roadItem: IReportMaintenanceFilterSections) => {
							return { value: roadItem.id, label: roadItem.name }
						})
						options = roadSectionsList
					}
				}
			}
			return options
		},
	},
})
