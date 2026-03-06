import {
	ChildReportSerivce,
	IReportChildRoadSummaryModel,
	IReportChildRoadSummaryRequest,
	IRoadsRoadSummaryFilter,
} from "../infrastructure"

interface IState {
	loading: boolean
	data: IReportChildRoadSummaryModel
	dataChildRoadSummaryFilter: IRoadsRoadSummaryFilter
}

export const useChildRoadSummaryReportStore = defineStore("report/child-road-summary", {
	state: (): IState => ({
		loading: false,
		data: { road_group_ids: [], type: "" } as IReportChildRoadSummaryModel,
		dataChildRoadSummaryFilter: {} as IRoadsRoadSummaryFilter,
	}),
	actions: {
		clearData() {
			this.data = {} as IReportChildRoadSummaryModel
		},
		async getChildRoadSummaryFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadSummaryFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadSummaryFilter = res.data
			}
		},
		async getExport() {
			this.loading = true

			const params: IReportChildRoadSummaryRequest = {
				road_group_ids: this.data.road_group_ids,
				type: this.data.type,
			}

			const service = new ChildReportSerivce()
			const res = await service.getExportRoadSummary(params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
