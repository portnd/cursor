import {
	IMaintenanceRoadGroup,
	IMaintenanceRoadOptions,
	IMaintenanceWarrantyData,
	IMaintenanceWarrantyUpdateRequest,
	MaintenanceHistoryService,
} from "../infrastructure"
import { ITree } from "~/core/shared/types/Tree"

interface IState {
	loading: boolean
	submitLoading: boolean
	interventionCriteria: ITree[]
	roadOptions: ITree[]
	defaultData: IMaintenanceWarrantyData
	distance: number
	interventionId: number | null
	kmStart: string
	kmEnd: string
	gridNo: number | null
	roadGroupData: IMaintenanceRoadGroup[]
	matchedRoad: IMaintenanceRoadOptions
	is_show_method: boolean
}

export const useMaintenanceHistoryWarrantyEditStore = defineStore("maintenance-history/warranty/edit", {
	state: (): IState => ({
		loading: false,
		submitLoading: false,
		interventionCriteria: [],
		roadOptions: [],
		defaultData: {} as IMaintenanceWarrantyData,
		distance: 0,
		interventionId: null,
		kmStart: "",
		kmEnd: "",
		gridNo: null,
		roadGroupData: [],
		matchedRoad: {} as IMaintenanceRoadOptions,
		is_show_method: false,
	}),
	actions: {
		async getIsShowMethod(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceDefault(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.is_show_method = res.data.is_show_method
			}
		},
		async getMaintenanceWarrantyInfo(id: number, mRoadId: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceWarranty(id, mRoadId)

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.defaultData = res.data
				this.interventionId = res.data.intervention_criteria.id
				this.kmEnd = convertMeterToKm(res.data.km_end)
				this.kmStart = convertMeterToKm(res.data.km_start)
				this.gridNo = res.data.grid_no
				this.onUpdateKmStart()
			}
		},
		async updateMaintenanceWarrantyInfo(id: number, mRoadId: number) {
			this.submitLoading = true

			const params: IMaintenanceWarrantyUpdateRequest = {
				grid_no: this.gridNo,
				intervention_criteria_id: Number(this.interventionId),
				km_end: convertStringToKm(this.kmEnd),
				km_start: convertStringToKm(this.kmStart),
				lane_no: this.defaultData.lane_no,
				maintenance_type: this.defaultData.maintenance_type,
				road_id: Number(this.defaultData.road_id),
			}

			const service = new MaintenanceHistoryService()
			const res = await service.updateMaintenanceWarranty(id, mRoadId, params)

			this.submitLoading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getInterventionCriteria() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceMethodList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				const options = res.data.map((data) => ({
					id: `parent${data.id}`,
					label: data.label,
					children: data.children.map((child) => ({
						id: `${child.id}`,
						label: child.label,
					})),
				}))

				this.interventionCriteria = options || []
			}
		},
		async getRoadGroupOptions() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceRoadGroup()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadGroupData = res.data

				const options = res.data?.map((group) => ({
					id: `parent${group.id}`,
					label: `${group.road_number} ${group.short_name}`,
					children: group.road_sections?.map((section) => ({
						id: `section${section.id}`,
						label: `${section.number} ${section?.name_origin} - ${section?.name_destination}`,
						children: section.roads?.map((road) => ({ id: `${road.id}`, label: road.name })),
					})),
				}))

				this.roadOptions = options || []
			}
		},
		onUpdateRoadId(id: number) {
			const { roadGroupData } = this
			const targetId = Number(id)

			if (!id) {
				this.kmEnd = ""
				this.kmStart = ""
			} else {
				const matchedData = roadGroupData.reduce((acc: IMaintenanceRoadOptions, parent) => {
					parent.road_sections?.forEach((section) => {
						section.roads?.forEach((road) => {
							if (road.id === targetId) {
								acc = road
							}
						})
					}, {})

					return acc
				}, {} as IMaintenanceRoadOptions)

				this.matchedRoad = matchedData
				this.defaultData.ref_direction_id = matchedData.ref_direction_id
				this.defaultData.lane_total = matchedData.lane_total
				this.kmEnd = convertMeterToKm(matchedData.km_end)
				this.kmStart = convertMeterToKm(matchedData.km_start)
				this.updateSumDistance()
			}
		},
		updateSumDistance() {
			const kmEnd = this.kmEnd ? convertStringToKm(this.kmEnd) : 0
			const kmStart = this.kmStart ? convertStringToKm(this.kmStart) : 0
			const sum = Math.abs(kmEnd - kmStart) / 1000

			this.defaultData.distance = isNaN(sum) ? 0.0 : parseFloat(sum.toFixed(3))
		},
		onUpdateKmStart() {
			const sum = Math.abs(convertStringToKm(this.kmEnd) - convertStringToKm(this.kmStart)) / 1000
			this.distance = parseFloat(sum?.toFixed(3))
		},
		onUpdateMaintenanceCriteria() {
			if (this.defaultData.maintenance_type === 1) {
				this.gridNo = null
			}
		},
	},
	getters: {
		getGenerateLane(state) {
			const { defaultData } = state

			const lanes = Array.from({ length: defaultData.lane_total }, (_, i) => i + 1) || []

			return defaultData.ref_direction_id === 1 ? lanes.reverse() : lanes
		},
		getGenerateGrid(state) {
			const { defaultData } = state

			const grids = Array.from({ length: 4 }, (_, i) => i + 1)

			return defaultData.ref_direction_id === 1 ? grids.reverse() : grids || []
		},
	},
})
