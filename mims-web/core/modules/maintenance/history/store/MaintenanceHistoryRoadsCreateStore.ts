import {
	IMaintenanceRoadGroup,
	IMaintenanceRoadOptions,
	IMaintenanceWarrantyCreateRequest,
	MaintenanceHistoryService,
} from "../infrastructure"
import { ITree } from "~/core/shared/types/Tree"

interface IState {
	loading: boolean
	submitLoading: boolean
	roadOptions: ITree[]
	interventionCriteria: ITree[]
	grid_no: number | null
	intervention_criteria_id: number | null
	km_end: string | null
	km_start: string | null
	lane_no: number | null
	maintenance_type: number | null
	road_id: number | null
	roadGroupData: IMaintenanceRoadGroup[]
	refDirectionId: number
	totalLane: number
	matchedRoad: IMaintenanceRoadOptions
	based_km_start: number | null
	based_km_end: number | null
	is_show_method: boolean
	sum_distance: number
}

export const useMaintenanceRoadsCreateStore = defineStore("maintenance-history/roads/create", {
	state: (): IState => ({
		loading: false,
		submitLoading: false,
		roadOptions: [],
		interventionCriteria: [],
		grid_no: null,
		intervention_criteria_id: null,
		km_end: null,
		km_start: null,
		lane_no: null,
		maintenance_type: 1,
		road_id: null,
		roadGroupData: [],
		refDirectionId: 1,
		totalLane: 2,
		matchedRoad: {} as IMaintenanceRoadOptions,
		based_km_start: null,
		based_km_end: null,
		is_show_method: false,
		sum_distance: 0,
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
		async createMaintenanceRoadsInfo(id: number) {
			this.submitLoading = true

			const params: IMaintenanceWarrantyCreateRequest = {
				grid_no: this.grid_no,
				intervention_criteria_id: Number(this.intervention_criteria_id),
				km_end: this.km_end ? convertStringToKm(this.km_end) : null,
				km_start: this.km_start ? convertStringToKm(this.km_start) : null,
				lane_no: this.lane_no,
				maintenance_type: this.maintenance_type,
				road_id: Number(this.road_id),
			}

			const service = new MaintenanceHistoryService()
			const res = await service.createMaintenanceRoad(id, params)

			this.submitLoading = false
			if (!res.data) {
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

				console.log("res.data =", res.data)

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
				// this.refDirectionId = 1
				this.km_end = ""
				this.km_start = ""
			} else {
				const matchedData = roadGroupData.reduce((acc: IMaintenanceRoadOptions, parent) => {
					parent.road_sections.forEach((section) => {
						section.roads?.forEach((road) => {
							if (road.id === targetId) {
								acc = road
							}
						})
					})

					return acc
				}, {} as IMaintenanceRoadOptions)

				this.matchedRoad = matchedData
				this.refDirectionId = matchedData.ref_direction_id
				this.totalLane = matchedData.lane_total
				this.based_km_end = matchedData.km_end
				this.based_km_start = matchedData.km_start
				this.km_end = convertMeterToKm(matchedData.km_end)
				this.km_start = convertMeterToKm(matchedData.km_start)
				this.updateSumDistance()
			}
		},
		updateSumDistance() {
			const kmEnd = this.km_end ? convertStringToKm(this.km_end) : 0
			const kmStart = this.km_start ? convertStringToKm(this.km_start) : 0
			const sum = Math.abs(kmEnd - kmStart) / 1000

			this.sum_distance = isNaN(sum) ? 0.0 : parseFloat(sum.toFixed(3))
		},
		onUpdateMaintenanceType(id: number) {
			if (id === 1) {
				this.grid_no = null
			}
		},
	},
	getters: {
		// getSumDistance(state) {
		// 	const kmEnd = state.km_end ? convertStringToKm(state.km_end) : 0
		// 	const kmStart = state.km_start ? convertStringToKm(state.km_start) : 0
		// 	const sum = Math.abs(kmEnd - kmStart) / 1000
		// 	return isNaN(sum) ? 0.0 : sum.toFixed(2)
		// },
		getGenerateLane(state) {
			const { totalLane, refDirectionId } = state

			const lanes = Array.from({ length: totalLane }, (_, i) => i + 1)

			return refDirectionId === 1 ? lanes.reverse() : lanes
		},
		getGenerateGrid(state) {
			const { refDirectionId } = state

			const grids = Array.from({ length: 4 }, (_, i) => i + 1)

			return refDirectionId === 1 ? grids.reverse() : grids
		},
	},
})
