import {
	MaintenanceService,
	IMaintenance,
	ICondition,
	IAnalysiSurfacesRule,
	IAnalysisRuleItem,
	IAnalysisRule,
	IAnalysisRuleItemList,
	IAnalysisMethod,
} from "../infrastructure"

interface IState {
	category: string
	value: string
	type: string
	activeStandardIndex: number
	data: IAnalysisRule
	prevData: IMaintenance
	loading: boolean

	methodId: number
	methods: IAnalysiSurfacesRule
	interventionCriteriasSelected: IAnalysisRuleItem[]
	standardSelected?: IAnalysisRuleItem
}

export const useMaintenanceEditStore = defineStore("setting/model/maintenance/edit", {
	state: (): IState => ({
		category: "",
		value: "",
		type: "",
		activeStandardIndex: 0,
		data: {} as IAnalysisRule,
		prevData: {} as IMaintenance,
		loading: false,
		methodId: 0,
		methods: {} as IAnalysiSurfacesRule,
		interventionCriteriasSelected: [],
		standardSelected: undefined,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true
			const maintenanceService = new MaintenanceService()
			const res = await maintenanceService.get()

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				return res
			}
		},
		async fetchRule() {
			const analysisRuleService = new MaintenanceService()
			const res = await analysisRuleService.getRule()

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.methods = res.data
				return res
			}
		},

		async getData() {
			// Loading
			this.loading = true
			await this.fetchRule()
			await this.get()
			// Loading
			this.loading = false
			const final = this.methods.asphalt[0]
			this.methodId = final.id
			this.category = final.name
			this.type = "asphalt"
		},

		async edit() {
			this.loading = true
			const results: any[] = []

			console.log("params", this.data)
			const res = await new MaintenanceService().create(this.data)
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				results.push(res)
			}
			console.log(2, results)

			return res
		},
		getAnalysisRuleParams(item: any, method: string, type: string): IAnalysisRuleItemList {
			if (item.maintenance_condition !== null) {
				return {
					maintenance_condition: item.maintenance_condition.map((e: any) => {
						return {
							condition_criterion: e.condition_criterion,
							condition_link: e.condition_link,
							condition_operation_1: e.condition_operation_1,
							condition_operation_2: e.condition_operation_2,
							condition_value_1: Number(e.condition_value_1),
							condition_value_2: Number(e.condition_value_2),
							id: e.is_new ? null : e.id,
						}
					}),
					maintenance_cost_per_unit: Number(item.maintenance_cost_per_unit),
					maintenance_description: item.maintenance_description,
					maintenance_method: method,
					maintenance_scraping: Number(item.maintenance_scraping),
					// maintenance_sequence: Number(item.maintenance_sequence),
					maintenance_standard_name: item.maintenance_standard_name,
					maintenance_surface_type_id: item.maintenance_surface_type_id,
					maintenance_thickness: Number(item.maintenance_thickness),
					maintenance_type: type,
					id: item.is_new ? null : item.id,
				}
			} else {
				return {
					maintenance_condition: [],
					maintenance_cost_per_unit: Number(item.maintenance_cost_per_unit),
					maintenance_description: item.maintenance_description,
					maintenance_method: method,
					maintenance_scraping: Number(item.maintenance_scraping),
					// maintenance_sequence: Number(item.maintenance_sequence),
					maintenance_standard_name: item.maintenance_standard_name,
					maintenance_surface_type_id: item.maintenance_surface_type_id,
					maintenance_thickness: Number(item.maintenance_thickness),
					maintenance_type: type,
					id: item.is_new ? null : item.id,
				}
			}
		},

		async delete(id: number) {
			// Loading
			this.loading = true
			const maintenanceService = new MaintenanceService()
			const res = await maintenanceService.delete(id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				return res
			}
		},
		addStandard() {
			const dataType = this.data[this.type]
			const index = dataType?.findIndex((e: IAnalysisMethod) => {
				return e.id === this.methodId
			})

			if (index === -1) {
				const item = {
					id: this.methodId,
					intervention_criterias: [],
					name: this.category,
				} as IAnalysisMethod
				dataType.push(item)
			} else if (!dataType[index].intervention_criterias) {
				dataType[index].intervention_criterias = []
			}
			this.updateInterventionCriteriasSelected()

			const id = Math.max(
				...this.interventionCriteriasSelected.map((item: any) => {
					return item.id ?? 0
				})
			)

			this.interventionCriteriasSelected.push({
				id: id + 1,
				is_new: true,
				maintenance_method: this.category,
				maintenance_cost_per_unit: null,
				maintenance_description: "",
				maintenance_scraping: null,
				maintenance_sequence: null,
				maintenance_standard_name: this.category + "-" + new Date().valueOf(),
				maintenance_surface_type_id: null,
				maintenance_thickness: null,
				maintenance_type: this.type,
				maintenance_condition: [
					{
						condition_criterion: "Cd.IRI",
						condition_link: "",
						condition_operation_1: "<=",
						condition_operation_2: "<",
						condition_value_1: 0.0,
						condition_value_2: 0.0,
					} as ICondition,
				],
			} as IAnalysisRuleItem)

			this.updateStandartsSelected()
			this.activeStandardIndex = this.interventionCriteriasSelected.length - 1
			console.log("addStandard ===> standardSelected", this.standardSelected)
			console.log("addStandard=", this.interventionCriteriasSelected)
		},
		deleteStandard(index: any) {
			this.getInterventionCriterias.splice(index, 1)
			this.activeStandardIndex = this.getInterventionCriterias.length - 1
			if (this.activeStandardIndex < 0) {
				this.activeStandardIndex = 0
			}
		},
		duplicateStandard(index: any) {
			const clonedObject = JSON.parse(JSON.stringify(this.getInterventionCriterias[index]))
			const id = Math.max(...this.getInterventionCriterias.map((item: IAnalysisRuleItem) => item.id ?? 0))
			const sequence = Math.max(
				...this.getInterventionCriterias.map((item: IAnalysisRuleItem) => item.maintenance_sequence ?? 0)
			)
			clonedObject.id = id + 1
			clonedObject.is_new = true
			clonedObject.maintenance_sequence = sequence + 1
			clonedObject.maintenance_standard_name = this.category + "-copy-" + new Date().valueOf()
			this.getInterventionCriterias[index].maintenance_condition.forEach((item: any, i) => {
				const data = item

				data.id = null

				clonedObject.maintenance_condition[i] = data
			})
			this.getInterventionCriterias.push(clonedObject)
			this.activeStandardIndex = this.getInterventionCriterias.length - 1
		},
		addCondition() {
			this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition.push({
				condition_link: "AND",
				condition_value_1: 0.0,
				condition_operation_1: "<=",
				condition_criterion: "Rutting",
				condition_operation_2: "<",
				condition_value_2: 0.0,
				is_new: true,
				id: this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition.length,
			} as ICondition)
		},
		deleteCondition(index: any) {
			this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition.splice(index, 1)
		},
		switchCondition(index1: any, index2: any, option: any) {
			// index1 ตำแหน่งของ row นั้น
			// index2 ตำแหน่งของ row ที่จะเอามาเปรียบเทียบ
			const conditionCurrent =
				this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[index1].condition_link
			const conditionNext =
				this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[index2].condition_link
			const data1 = Object.assign(
				{},
				this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[index1]
			)
			const data2 = Object.assign(
				{},
				this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[index2]
			)
			if (option === "up") {
				if (index1 === 0) {
					return false
				} else if (data2.id === this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[0].id) {
					data1.condition_link = ""
					data2.condition_link = conditionCurrent
				} else {
					data1.condition_link = conditionCurrent
					data2.condition_link = conditionNext
				}
			} else if (option === "down") {
				if (index2 === this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition.length) {
					return false
				} else if (data1.id === this.getInterventionCriterias[this.activeStandardIndex].maintenance_condition[0].id) {
					data1.condition_link = conditionNext
					data2.condition_link = ""
				} else {
					data1.condition_link = conditionCurrent
					data2.condition_link = conditionNext
				}
			}

			if (this.getICSelected) {
				this.getICSelected.maintenance_condition[index1] = data2
				this.getICSelected.maintenance_condition[index2] = data1
			}
		},
		async restore() {
			const lastIndex = this.getInterventionCriterias.length - 1

			if (this.getInterventionCriterias[lastIndex].is_new ?? false) {
				const data = this.getInterventionCriterias[lastIndex]
				const dataType = this.data[this.type]
				const dataSelected = dataType?.find((e: IAnalysisMethod) => {
					return e.id === this.methodId
				})

				if (dataSelected) {
					dataSelected.intervention_criterias[this.activeStandardIndex] = {
						id: data.id,
						is_new: true,
						maintenance_method: data.maintenance_method,
						maintenance_cost_per_unit: null,
						maintenance_description: "",
						maintenance_scraping: null,
						maintenance_sequence: data.maintenance_sequence,
						maintenance_standard_name: data.maintenance_standard_name,
						maintenance_surface_type_id: null,
						maintenance_thickness: null,
						maintenance_type: data.maintenance_type,
						maintenance_condition: [
							{
								condition_criterion: "Cd.IRI",
								condition_link: "",
								condition_operation_1: "<=",
								condition_operation_2: "<",
								condition_value_1: 0.0,
								condition_value_2: 0.0,
							} as ICondition,
						],
					} as IAnalysisRuleItem
				}
				this.updateInterventionCriteriasSelected()
				this.updateStandartsSelected()
			} else {
				const type = this.type
				const value = this.value
				const methodId = this.methodId
				const activeStandard = this.activeStandardIndex
				const category = this.category

				this.$reset()
				this.loading = true
				await this.fetchRule()
				await this.get()
				// Loading
				this.loading = false
				this.category = category
				this.value = value
				this.methodId = methodId
				this.type = type
				this.activeStandardIndex = activeStandard
			}
		},

		updateInterventionCriteriasSelected() {
			const dataType = this.data[this.type]
			const dataSelected = dataType?.find((e: IAnalysisMethod) => {
				return e.id === this.methodId
			})
			this.interventionCriteriasSelected = dataSelected?.intervention_criterias ?? []
		},
		updateStandartsSelected() {
			const dataType = this.data[this.type]
			const dataSelected = dataType?.find((e: IAnalysisMethod) => {
				return e.id === this.methodId
			})
			this.standardSelected = dataSelected?.intervention_criterias[this.activeStandardIndex]
		},
	},

	getters: {
		generateSummary(state) {
			const sum = ref<Array<any>>([])
			if (state) {
				// const types = state.data[state.type]
				const dataType = state.data[state.type]
				const dataSelected = dataType?.find((e: IAnalysisMethod) => {
					return e.id === state.methodId
				})
				const iCSelected = dataSelected?.intervention_criterias[state.activeStandardIndex]

				if (iCSelected) {
					const conditions = iCSelected?.maintenance_condition
					for (let i = 0; i < conditions?.length; i++) {
						const condition = conditions[i]
						if (i === 0) {
							sum.value.push("<br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;")
							sum.value.push("= ")
							sum.value.push("(")
						}
						if (condition.condition_link !== "") {
							sum.value.push(` ${condition.condition_link} `)
							sum.value.push(
								`${toNumber(condition.condition_value_1)} ${condition.condition_operation_1} ${
									condition.condition_criterion
								} ${condition.condition_operation_2} ${toNumber(condition.condition_value_2)}`
							)
						} else {
							sum.value.push(
								`${toNumber(condition.condition_value_1)} ${condition.condition_operation_1} ${
									condition.condition_criterion
								} ${condition.condition_operation_2} ${toNumber(condition.condition_value_2)}`
							)
						}
						if (i === conditions.length - 1) {
							sum.value.push(")")
						}
					}
					const final = sum.value.flatMap((item: any) => {
						// if (sum.value[1] === " AND ") {
						// 	sum.value[1] = ""
						// }
						if (item === " OR ") {
							return [")", "<br>", item, "("]
						} else {
							return [item]
						}
					})

					return iCSelected.maintenance_standard_name + " " + final.join("")
				}
			}
		},
		getDataSelected(state) {
			const dataType = state.data[state.type]

			const dataSelected = dataType?.find((e: IAnalysisMethod) => {
				return e.id === state.methodId
			})

			return dataSelected
		},
		getInterventionCriterias(state) {
			const dataType = state.data[state.type]
			const dataSelected = dataType?.find((e: IAnalysisMethod) => {
				return e.id === state.methodId
			})
			return dataSelected?.intervention_criterias ?? []
		},
		getICSelected(state) {
			const dataType = state.data[state.type]
			const dataSelected = dataType?.find((e: IAnalysisMethod) => {
				return e.id === state.methodId
			})
			if (dataSelected && dataSelected?.intervention_criterias) {
				return dataSelected.intervention_criterias[state.activeStandardIndex]
			}

			return {} as IAnalysisRuleItem
		},
		getMethodIndex(state) {
			const dataType = state.data[state.type]
			const index =
				dataType?.findIndex((e: IAnalysisMethod) => {
					return e.id === state.methodId
				}) ?? -1
			return index
		},
	},
})
