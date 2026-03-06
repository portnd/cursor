import {
	IMaintenanceBudgetCriteria,
	IMaintenanceDefaultData,
	IMaintenanceHistoryUpdateRequest,
	MaintenanceHistoryService,
} from "../infrastructure"
import { IFile } from "~/core/shared/types/File"
import { ITree } from "~/core/shared/types/Tree"

interface IState {
	loading: boolean
	submitLoading: boolean
	maintenanceBudget: IMaintenanceBudgetCriteria[]
	ownerOptions: ITree[]
	defaultData: IMaintenanceDefaultData
	budgetId: number | null
	budgetMethodId: number | null
	filesPath: string[]
	files: IFile[]
}

export const useMaintenanceHistoryEditStore = defineStore("maintenance-history/edit", {
	state: (): IState => ({
		loading: false,
		submitLoading: false,
		maintenanceBudget: [],
		ownerOptions: [],
		defaultData: {} as IMaintenanceDefaultData,
		budgetId: null,
		budgetMethodId: null,
		filesPath: [],
		files: [],
	}),
	actions: {
		async getMaintenanceBudget() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.maintenanceBudget = res.data
			}
		},

		async getDefault(id: number) {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceDefault(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.defaultData = res.data
				this.budgetId = res.data.budget.id
				this.budgetMethodId = res.data.budget_method.id
				this.filesPath = res.data.attachments.map((item) => item.path)
			}
		},
		async getDivisiontOptions() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceDivision()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				const options = res.data.map((parent) => {
					return {
						id: parent.owner_code_key,
						label: parent.name,
						children: parent.districts.map((district) => {
							return {
								id: district.owner_code_key,
								label: district.name,
								children: district.depots.map((depot) => {
									return {
										id: depot.owner_code_key,
										label: depot.name,
									}
								}),
							}
						}),
					}
				})

				this.ownerOptions = options || []
			}
		},
		async updateMaintenanceData(id: number) {
			this.submitLoading = true

			this.files.forEach((item, index) => {
				const data = item.data
				const attachment = this.defaultData.attachments[index]

				if (data && attachment && attachment.path.includes(data.name)) {
					data.id = String(attachment.id)
				}
			})

			const params: IMaintenanceHistoryUpdateRequest = {
				advisor_name: this.defaultData.advisor_name,
				attachments: this.files.map((item) => ({
					id: Number(item.data!.id),
					file_name: item.data?.name,
					file: item.data?.base64,
					status: item.status,
				})),
				budget_id: this.budgetId,
				budget_maintenance: this.defaultData.budget_maintenance,
				budget_method_id: this.budgetMethodId,
				budget_procurement: this.defaultData.budget_procurement,
				budget_year: this.defaultData.budget_year,
				contract_number: this.defaultData.contract_number,
				contract_work_value: this.defaultData.contract_work_value,
				contractor_name: this.defaultData.contractor_name,
				guarantee_expiration_date: formatDate(this.defaultData.guarantee_expiration_date),
				middle_price: this.defaultData.middle_price,
				name: this.defaultData.name,
				owner_code: this.defaultData.owner_code,
				project_details: this.defaultData.project_details,
				project_end_date: formatDate(this.defaultData.project_end_date),
				project_secretary_name: this.defaultData.project_secretary_name,
			}

			const service = new MaintenanceHistoryService()
			const res = await service.updateMaintenanceData(id, params)

			this.submitLoading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		getDivisionOption() {
			const refDivision = useInitData().refDivision()

			if (refDivision?.length === 0) {
				return []
			}

			const options = useInitData()
				.refDivision()
				?.map((item) => {
					if (item.districts.length !== 0) {
						const districtsArray = item.districts.map((subDistricts) => {
							const depotsArray = subDistricts.depots.map((subDepots) => ({
								label: subDepots.name,
								id: subDepots.owner_code_key,
							}))
							return {
								label: subDistricts.name,
								id: subDistricts.owner_code_key,
								children: depotsArray,
							}
						})
						return {
							label: item.name,
							id: item.owner_code_key,
							children: districtsArray,
						}
					}
					return {
						label: item.name,
						id: item.owner_code_key,
						children: [],
					}
				})

			return options
		},
	},
	getters: {
		getBudgetOptions(state) {
			const { maintenanceBudget } = state
			const options = maintenanceBudget.map((item) => ({ label: item.name, value: item.id }))

			return options ?? []
		},
		getMaintenanceCriteriaOptions(state) {
			const { maintenanceBudget, budgetId } = state
			const matchedBudget = maintenanceBudget.find((item) => item.id === Number(budgetId))?.budget_methods || []
			const options = matchedBudget.map((item) => ({ label: item.method_name, value: item.id }))

			return options ?? []
		},
	},
})
