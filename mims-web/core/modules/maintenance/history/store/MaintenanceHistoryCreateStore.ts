import {
	IMaintenanceBudgetCriteria,
	IMaintenanceHistoryCreateRequest,
	MaintenanceHistoryService,
} from "../infrastructure"
import { IFile } from "~/core/shared/types/File"
import { ITree } from "~/core/shared/types/Tree"

interface IStateParams {
	advisor_name: string
	attachments: IFile[]
	budget_id: number
	budget_maintenance: number | null
	budget_method_id: number
	budget_procurement: number | null
	budget_year: number
	contract_number: string
	contract_work_value: number | null
	contractor_name: string
	guarantee_expiration_date: string
	middle_price: number | null
	name: string
	owner_code: string | null
	project_details: string
	project_end_date: string
	project_secretary_name: string
}

interface IState {
	loading: boolean
	submitLoading: boolean
	maintenanceBudget: IMaintenanceBudgetCriteria[]
	ownerOptions: ITree[]
	params: IStateParams
}

export const useMaintenanceHistoryCreateStore = defineStore("maintenance-history/create", {
	state: (): IState => ({
		loading: false,
		submitLoading: false,
		maintenanceBudget: [],
		ownerOptions: [],
		params: {
			advisor_name: "",
			attachments: [],
			budget_id: 0,
			budget_maintenance: null,
			budget_method_id: 0,
			budget_procurement: null,
			budget_year: 0,
			contract_number: "",
			contract_work_value: null,
			contractor_name: "",
			guarantee_expiration_date: "",
			middle_price: null,
			name: "",
			owner_code: null,
			project_details: "",
			project_end_date: "",
			project_secretary_name: "",
		},
	}),
	actions: {
		async getMaintenanceBudget() {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.maintenanceBudget = res.data
			}
		},
		async createMaintenance() {
			this.submitLoading = true

			const params: IMaintenanceHistoryCreateRequest = {
				advisor_name: this.params.advisor_name,
				attachments: this.params.attachments.map((item) => ({
					file_name: item.data?.name,
					file: item.data?.base64,
					status: item.status,
				})),
				budget_id: this.params.budget_id,
				budget_maintenance: this.params.budget_maintenance,
				budget_method_id: this.params.budget_method_id,
				budget_procurement: this.params.budget_procurement,
				budget_year: this.params.budget_year,
				contract_number: this.params.contract_number,
				contract_work_value: this.params.contract_work_value,
				contractor_name: this.params.contractor_name,
				guarantee_expiration_date: formatDate(this.params.guarantee_expiration_date),
				middle_price: this.params.middle_price,
				name: this.params.name,
				owner_code: this.params.owner_code,
				project_details: this.params.project_details,
				project_end_date: formatDate(this.params.project_end_date),
				project_secretary_name: this.params.project_secretary_name,
			}

			const service = new MaintenanceHistoryService()
			const res = await service.createMaintenance(params)

			this.submitLoading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
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
			const { maintenanceBudget, params } = state
			const matchedBudget = maintenanceBudget.find((item) => item.id === Number(params.budget_id))?.budget_methods || []
			const options = matchedBudget.map((item) => ({ label: item.method_name, value: item.id }))

			return options ?? []
		},
	},
})
