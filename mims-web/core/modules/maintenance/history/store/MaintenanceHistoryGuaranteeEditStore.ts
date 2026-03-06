import { IMaintenanceHistoryEditParams, MaintenanceHistoryService } from "../infrastructure"
import {
	IRoadChildListData,
	IMaintenanceHistoryDetailMaintenanceRoad,
	IMaintenanceMethodListData,
	IMaintenanceBudgetCriteria,
	IMaintenanceHistoryDetailData,
} from "../infrastructure/MaintenanceHistoryModel"
import { IOption } from "~/core/shared/types/Option"
import { IFile, IMultiFile } from "~/core/shared/types/File"

interface IParamsState {
	attachments: IMultiFile[]
	id: number
	intervention_criteria_id: number | null
	km_end: string
	km_start: string
	lane: string
	road_group_id: string
	road_id: string
}

interface IState {
	loading: boolean
	sumInput: number
	projectData: IMaintenanceHistoryDetailData
	data: IMaintenanceHistoryDetailMaintenanceRoad
	mantenanceMethodList: IMaintenanceMethodListData[]
	roadList: IRoadChildListData
	kmControlOptions: IOption[]
	laneOptions: IOption[]
	interventionCriteriaOptions: IOption[]
	params: IParamsState
	sumDistance: number
	isShowMethod: boolean | null
	budgetCriteriaList: IMaintenanceBudgetCriteria[]
	filePaths: string[]
	files: IFile[]
}

export const useMaintenanceHistoryGuaranteeEditStore = defineStore("maintenance/history/guarantee-edit", {
	state: (): IState => ({
		loading: false,
		sumInput: 0,
		projectData: {} as IMaintenanceHistoryDetailData,
		data: {} as IMaintenanceHistoryDetailMaintenanceRoad,
		roadList: {} as IRoadChildListData,
		mantenanceMethodList: [],
		kmControlOptions: [],
		laneOptions: [],
		interventionCriteriaOptions: [],
		params: {
			attachments: [],
			id: 0,
			intervention_criteria_id: null,
			km_end: "",
			km_start: "",
			lane: "",
			road_group_id: "",
			road_id: "",
		},
		sumDistance: 0,
		isShowMethod: null,
		budgetCriteriaList: [],
		filePaths: [],
		files: [],
	}),
	actions: {
		async getHistoryDetail(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryDetails(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.projectData = res.data
			}
		},
		async putHistoryGuaranteeEdit(id: number) {
			this.loading = true

			const params: IMaintenanceHistoryEditParams = {
				id: this.params.id,
				intervention_criteria_id:
					this.params.intervention_criteria_id === null || this.params.intervention_criteria_id === 0
						? null
						: Number(this.params.intervention_criteria_id),
				km_end: convertStringToKm(this.params.km_end),
				km_start: convertStringToKm(this.params.km_start),
				lane: Number(this.params.lane),
				road_group_id: Number(this.params.road_group_id),
				road_id: Number(this.params.road_id),
				attachments: [],
			}

			this.projectData.maintenance_road_histories.forEach((item) => {
				if (item.id === params.id) {
					params.attachments = toFiles(item.attacchments, this.files)
				}
			})

			const service = new MaintenanceHistoryService()
			const res = await service.updateMaintenanceHistoryGuarantee(id, this.params.id, params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getRoadList(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getRoadList(id)

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadList = res.data
				await this.getMaintenanceBudgetCriteria()
				await this.getBudgetsList()
			}
		},
		createKmControlOptions(roadGroupId: number) {
			if (Object.keys(this.roadList).length > 0) {
				const main = this.roadList.id === roadGroupId ? this.roadList.roads : []
				const kmOptions =
					main.length > 0
						? main?.map((item) => {
								return { label: item.name, value: item.id }
						  })
						: []

				this.kmControlOptions = kmOptions
			}
		},
		createLaneOptions(kmControlId: number) {
			const main = this.roadList.id === Number(this.params.road_group_id) ? this.roadList.roads : []
			let options: IOption[] = []
			main.forEach((item) => {
				if (item.id === kmControlId) {
					options = item.lanes.map((child) => {
						return { label: child.lane.toString(), value: child.lane }
					})
				}
			})

			this.laneOptions = options.length > 0 ? options : []
		},
		async getBudgetsList() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceMethodList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.mantenanceMethodList = res.data
			}
		},
		async getMaintenanceBudgetCriteria() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.budgetCriteriaList = res.data
			}
		},
		createInterventionOptions() {
			const data = this.mantenanceMethodList.flatMap((item) => {
				return item
			})

			const options = data.map((parent) => {
				return {
					id: "m" + parent.id,
					label: parent.label,
					children: parent.children.map((child) => {
						return { id: child.id, label: child.label }
					}),
				}
			})

			this.interventionCriteriaOptions = options
		},
		setDefaultParams(itemId: number) {
			const data = this.projectData.maintenance_road_histories

			data.forEach((item) => {
				if (item.id === itemId) {
					this.params.id = itemId
					this.params.intervention_criteria_id = item.intervention_criteria_id
					this.params.km_end = convertMeterToKm(item.km_end)
					this.params.km_start = convertMeterToKm(item.km_start)
					this.params.road_id = item.road_id.toString()
					this.createLaneOptions(item.road_id)
					this.params.lane = this.getLane.toString()

					this.filePaths = item.attacchments?.map((file) => file.path)
				}
			})
		},
		calculateSumDistance() {
			if (this.params.km_start === "" && this.params.km_end === "") {
				this.sumInput = 0
			} else {
				const kmEnd = this.params.km_end ? convertStringToKm(this.params.km_end) : 0
				const kmStart = this.params.km_start ? convertStringToKm(this.params.km_start) : 0

				if (kmStart === kmEnd) {
					return 0
				}

				const sum = Math.abs(kmEnd - kmStart) / 1000
				this.sumInput = sum
			}
		},
		sumDistanceInput() {
			if (this.params.km_start || this.params.km_end) {
				const start = convertStringToKm(this.params.km_start)
				const end = convertStringToKm(this.params.km_end)
				this.sumDistance = parseFloat(((start + end) / 1000).toFixed(2))
			} else {
				this.sumDistance = 0
			}
		},
		checkIsShowMethod(methodId: number) {
			this.budgetCriteriaList.forEach((item) => {
				item.budget_methods.forEach((child) => {
					if (methodId === child.id) {
						this.isShowMethod = child.is_show_method
					}
				})
			})
		},
	},
	getters: {
		getInterventionId(state) {
			if (Object.keys(state.data).length === 0) {
				return 0
			}

			const data = state.data
			return data.intervention_criteria_id
		},
		getKmEndKmStart(state) {
			if (Object.keys(state.data).length === 0) {
				return { km_end: "", km_start: "" }
			}

			const data = state.data

			return { km_end: data.km_end, km_start: data.km_start }
		},
		getLane(state) {
			if (Object.keys(state.projectData).length === 0) {
				return ""
			}

			let lane = 0

			const data = state.projectData
			data.maintenance_road_histories.forEach((item) => {
				if (item.id === state.params.id) {
					lane = item.lane
				}
			})

			return lane
		},
	},
})
