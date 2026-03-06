import {
	ChildReportSerivce,
	IReportChildRoadReflectLightModel,
	IRoadsRoadReflectLightDepotFilterRoad,
	IRoadsRoadReflectLightFilter,
	IRoadsRoadReflectLightFilterRoadGroup,
	IRoadsRoadReflectLightFilterRoadSections,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	data: IReportChildRoadReflectLightModel
	dataChildRoadReflectLightFilter: IRoadsRoadReflectLightFilter
}

export const useChildRoadReflectLightReportStore = defineStore("report/child-road-reflect-light", {
	state: (): IState => ({
		loading: false,
		data: {} as IReportChildRoadReflectLightModel,
		dataChildRoadReflectLightFilter: {} as IRoadsRoadReflectLightFilter,
	}),
	actions: {
		async getChildRoadReflectLightFilter() {
			this.loading = true

			const service = new ChildReportSerivce()
			const res = await service.getChildRoadReflectLightFilter()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataChildRoadReflectLightFilter = res.data
			}
		},
	},
	getters: {
		getKmOptions(state) {
			const options: IOption[] = state.dataChildRoadReflectLightFilter?.filter_criteria?.map((e) => {
				return { value: e.id, label: e.name }
			})
			return options
		},
		getYearOptions(state) {
			const options: IOption[] = state.dataChildRoadReflectLightFilter?.filter_road?.map((e) => {
				return { value: e.year, label: `${e.year + 543}` }
			})
			return options
		},
		getDepartmentOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadReflectLightFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const depotList = data?.depot?.map((depotItem: IRoadsRoadReflectLightDepotFilterRoad) => {
					return { value: depotItem.id, label: depotItem.name }
				})
				options = depotList
			}
			return options
		},
		getRoadGroupOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadReflectLightFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadReflectLightDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.map((roadItem: IRoadsRoadReflectLightFilterRoadGroup) => {
						return { value: roadItem.id, label: roadItem.name }
					})
					options = roadGroupList
				}
			}
			return options
		},
		getRoadSectionOptions(state) {
			let options: IOption[] = []
			const data = state.dataChildRoadReflectLightFilter?.filter_road?.find((e) => e.year === state.data.year)
			if (data) {
				const departmentList = data.depot?.find(
					(item: IRoadsRoadReflectLightDepotFilterRoad) => item.id === state.data.department_id
				)
				if (departmentList) {
					const roadGroupList = departmentList?.road_group?.find((e) => e.id === state.data.road_group_id)
					if (roadGroupList) {
						const roadSectionsList = roadGroupList?.road_section?.map(
							(roadItem: IRoadsRoadReflectLightFilterRoadSections) => {
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
