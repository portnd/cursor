import {
	ChildReportSerivce,
	IReportChildRoadSurfaceModel,
	IRoadsRoadSurfaceDepotFilterRoad,
	IRoadsRoadSurfaceFilter,
	IRoadsRoadSurfaceFilterRoadGroup,
	IRoadsRoadSurfaceFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadSurfaceModel
	dataChildRoadSurfaceFilter: IRoadsRoadSurfaceFilter
}

export const useChildRoadSurfaceReportStore = defineStore("report/child-road-surface", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadSurfaceModel,
		dataChildRoadSurfaceFilter: {} as IRoadsRoadSurfaceFilter,
	}),
	actions: {
		clearData() {
			this.data = {} as IReportChildRoadSurfaceModel
			this.dataChildRoadSurfaceFilter = {} as IRoadsRoadSurfaceFilter
		},
		async getChildRoadSurfaceFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadSurfaceFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadSurfaceFilter = res.data
			}
		},
	},
	getters: {
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadSurfaceFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadSurfaceFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadSurfaceDepotFilterRoad) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadSurfaceFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadSurfaceDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map((roadItem: IRoadsRoadSurfaceFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadSurfaceFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadSurfaceDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadSurfaceFilterRoadSections) => {
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
