import {
	ChildReportSerivce,
	IReportVolumeAadtModel,
	IRoadsVolumeAadtFilter,
	IRoadsVolumeAadtFilterRoadGroup,
	IRoadsVolumeAadtFilterSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportVolumeAadtModel
	dataReportVolumeAadtFilter: IRoadsVolumeAadtFilter
}

export const useReportVolumeAadtReportStore = defineStore("report/report-volume-aadt", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportVolumeAadtModel,
		dataReportVolumeAadtFilter: {} as IRoadsVolumeAadtFilter,
	}),
	actions: {
		async getReportVolumeAadtFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getReportVolumeAadtFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataReportVolumeAadtFilter = res.data
			}
		},
	},
	getters: {
		getYearOptions(state) {
			const options: IOption[] = state.dataReportVolumeAadtFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},

		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportVolumeAadtFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const roadGroupList = data?.road_group?.map((roadItem: IRoadsVolumeAadtFilterRoadGroup) => {
					return { value: roadItem.id, label: roadItem.name }
				})
				options = roadGroupList
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataReportVolumeAadtFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const roadGroupList = data?.road_group?.find((e) => e.id === state.data.road_group_id)
				if (roadGroupList) {
					const roadSectionsList = roadGroupList?.road_section?.map((roadItem: IRoadsVolumeAadtFilterSections) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadSectionsList
				}
			}
			return options
		},
	},
})
