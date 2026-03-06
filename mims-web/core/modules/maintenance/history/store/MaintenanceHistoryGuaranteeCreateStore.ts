import { IMaintenanceHistoryGuaranteeCreateParams } from "../infrastructure/MaintenanceHistoryRequest"
import {
	IRoadChildListData,
	IMaintenanceMethodListData,
	IMaintenanceBudgetCriteria,
} from "../infrastructure/MaintenanceHistoryModel"
import { MaintenanceHistoryService } from "../infrastructure/MaintenanceHistoryService"
import { IOption } from "~/core/shared/types/Option"
import { IFile } from "~/core/shared/types/File"

interface IParamsState {
	id: number
	intervention_criteria_id: number | null
	km_end: string
	km_start: string
	lane: number | null
	road_group_id: string
	road_id: string
}

interface IState {
	loading: boolean
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

export const useMaintenanceHistoryGuaranteeCreateStore = defineStore("maintenance/history/guarantee-create", {
	state: (): IState => ({
		loading: false,
		mantenanceMethodList: [],
		roadList: {} as IRoadChildListData,
		kmControlOptions: [],
		laneOptions: [],
		interventionCriteriaOptions: [],
		params: {
			id: 0,
			intervention_criteria_id: null,
			km_end: "",
			km_start: "",
			lane: null,
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
		async createHistoryGuarantee(id: number) {
			this.loading = true

			const params: IMaintenanceHistoryGuaranteeCreateParams = {
				id: this.params.id,
				intervention_criteria_id:
					this.params.intervention_criteria_id !== null ? Number(this.params.intervention_criteria_id) : null,
				km_end: convertStringToKm(this.params.km_end),
				km_start: convertStringToKm(this.params.km_start),
				lane: this.params.lane ? this.params.lane : 0,
				road_group_id: Number(this.params.road_group_id),
				road_id: Number(this.params.road_id),
				attachments: toFiles([], this.files),
			}

			const service = new MaintenanceHistoryService()
			const res = await service.createMaintenanceHistoryGuarantee(id, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
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
		async getRoadList(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getRoadList(id)

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadList = res.data
				await this.getBudgetsList()
			}
		},
		async getMaintenanceBudgets() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.budgetCriteriaList = res.data
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
			if (Object.keys(this.roadList).length > 0) {
				const main = this.roadList.id === Number(this.params.road_group_id) ? this.roadList.roads : []
				let options: IOption[] = []

				main.forEach((item) => {
					if (item.id === kmControlId) {
						options = item.lanes.map((child) => {
							const kmStart = convertMeterToKm(child.km_start)
							const kmEnd = convertMeterToKm(child.km_end)
							return { label: `${child.lane} (กม.ที่ ${kmStart} - ${kmEnd})`, value: child.lane }
						})
					}
				})

				this.laneOptions = options.length > 0 ? options : []
			}
		},
		createInterventionOptions() {
			const data = this.mantenanceMethodList.flatMap((parent) => {
				return parent
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
	getters: {},
})
