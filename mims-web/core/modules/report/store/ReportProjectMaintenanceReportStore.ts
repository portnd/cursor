import {
	ChildReportSerivce,
	IReportProjectMaintenanceModel,
	IReportProjectMaintenanceFilter,
	IReportProjectMaintenanceFilterSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportProjectMaintenanceModel
	dataReportProjectMaintenanceFilter: IReportProjectMaintenanceFilter
}

export const useReportProjectMaintenanceReportStore = defineStore("report/report-Project-maintenance", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportProjectMaintenanceModel,
		dataReportProjectMaintenanceFilter: {} as IReportProjectMaintenanceFilter,
	}),
	actions: {
		async getReportProjectMaintenanceFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getReportProjectMaintenanceFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataReportProjectMaintenanceFilter = res.data
			}
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const roadSectionsList = state.dataReportProjectMaintenanceFilter?.filter_road?.map(
				(roadItem: IReportProjectMaintenanceFilterSections) => {
					return { value: roadItem.id, label: roadItem.name }
				}
			)
			options = roadSectionsList
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportProjectMaintenanceFilter?.filter_road?.find((e) => e.id === state.data.road_group_id)
			if (data) {
				const roadSectionsList = data?.road_section?.map((roadItem: IReportProjectMaintenanceFilterSections) => {
					return { value: roadItem.id, label: roadItem.name }
				})
				options = roadSectionsList
			}
			return options
		},
		getYearStartOptions(state) {
			let options: IOption[] = []
			const yearStart = state.dataReportProjectMaintenanceFilter?.filter_Year?.start_year?.map((year: number) => {
				return { value: year, label: `${year + 543}` }
			})
			options = yearStart
			return options
		},
		getYearEndOptions(state) {
			let options: IOption[] = []
			const yearEnd = state.dataReportProjectMaintenanceFilter?.filter_Year?.end_year?.map((year: number) => {
				return { value: year, label: `${year + 543}` }
			})
			options = yearEnd
			return options
		},
	},
})
