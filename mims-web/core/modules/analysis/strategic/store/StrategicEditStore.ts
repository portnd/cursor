import {
	IStrategicAnalyzeData,
	IStrategicCreatePrepareDataReq,
	IStrategicRoadGroup,
	IStrategicStep2,
	IStrategicStep2plan,
	IStrategicUpdateAnalyzeParams,
	IStrategicUpdatePrepareData,
	StrategicsService,
} from "../infrastructure"
import { IStrategicsList } from "../../list/infrastructure"
import { ITree } from "~/core/shared/types/Tree"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	submit_loading: boolean
	prepare_loading: boolean
	roadGroup: IStrategicRoadGroup[]
	roadGroupOptions: ITree[]
	default: IStrategicAnalyzeData
	road_id: number[]
	surfaceOptions: IOption[]
	groupKmOptions: IOption[]
	strategicList: IStrategicsList[]
	selected_id: number[]
	allPrepareId: number[]
	step2Data: IStrategicStep2
	step: number
	isCopy: boolean
	wasStep2: boolean
	prepare_data_id: number
}

export const useStrategicEditStore = defineStore("strategic/new-edit", {
	state: (): IState => ({
		loading: false,
		submit_loading: false,
		prepare_loading: false,
		roadGroup: [],
		roadGroupOptions: [],
		default: {} as IStrategicAnalyzeData,
		road_id: [],
		surfaceOptions: [
			{ label: "ลาดยาง", value: 1 },
			{ label: "คอนกรีต", value: 2 },
		],
		groupKmOptions: [
			{ label: "1 กม.", value: 1 },
			{ label: "2 กม.", value: 2 },
			{ label: "5 กม.", value: 5 },
			{ label: "10 กม.", value: 10 },
		],
		strategicList: [],
		selected_id: [],
		allPrepareId: [],
		step2Data: {} as IStrategicStep2,
		step: 1,
		isCopy: false,
		wasStep2: false,
		prepare_data_id: 0,
	}),
	actions: {
		async getRoadsOptions() {
			const service = new StrategicsService()
			const res = await service.getRoadTree()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadGroup = res.data

				const options: ITree[] = res.data.map((parent) => ({
					label: parent.short_name,
					id: `parent${parent.id}`,
					children: parent.road_sections.map((section) => ({
						label: `${section.name_origin} - ${section.name_destination}`,
						id: `section${section?.id}`,
						children: section.roads?.map((road) => ({
							label: road.name,
							id: `${road.id}`,
						})),
					})),
				}))

				this.roadGroupOptions = options ?? []
			}
		},
		async copy(id: number) {
			const service = new StrategicsService()
			const res = await service.copy(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// console.log(res.data)
				this.isCopy = true
				return res
			}
		},
		async getRefMaintenanceOptions() {
			const service = new StrategicsService()
			const res = await service.getStrategicList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.strategicList = res.data
			}
		},
		async getDefaultData(id: number) {
			const service = new StrategicsService()
			const res = await service.getAnalyzeDefaultDetails(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.default = res.data
				this.road_id = res.data.roads.map((item) => item.road_id)
			}
		},
		async checkPrepareDataStatus() {
			this.prepare_loading = true

			const service = new StrategicsService()
			const res = await service.getCheckPrepareData(this.prepare_data_id)

			this.prepare_loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else if (res.data?.status === false) {
				await this.checkPrepareDataStatus()
			} else {
				return res
			}
		},
		async getAllId() {
			const service = new StrategicsService()
			const res = await service.getPrepareDataId(this.prepare_data_id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.allPrepareId = res.data
			}
		},
		async getSelectedId() {
			const service = new StrategicsService()
			const res = await service.getPrepareDataSelectedId(this.default.id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.selected_id = res.data
			}
		},
		handleSelectedPrepareData(itemId: number[]) {
			this.selected_id = itemId
		},
		async createPrepareData() {
			// const { default } = this
			// reset ค่า
			this.selected_id = []

			const params: IStrategicCreatePrepareDataReq = {
				aadt1: this.default.aadt1,
				aadt2: this.default.aadt2,
				group_km: this.default.group_km,
				ifi1: this.default.ifi1,
				ifi2: this.default.ifi2,
				iri1: this.default.iri1,
				iri2: this.default.iri2,
				lane_type_id: this.default.lane_type_id,
				roads: this.default.roads.map(Number),
				surface_type_id: this.default.surface_type_id,
				maintenance_analysis_type_id: 1,
				name: this.default.name,
			}

			const service = new StrategicsService()
			const res = await service.createPreparingData(params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.prepare_data_id = res.data.id
				return res
			}
		},
		async updatePrepareData() {
			this.prepare_loading = true

			this.selected_id = []

			const params: IStrategicUpdatePrepareData = {
				aadt1: this.default.aadt1,
				aadt2: this.default.aadt2,
				iri1: this.default.iri1,
				iri2: this.default.iri2,
				ifi1: this.default.ifi1,
				ifi2: this.default.ifi2,
				group_km: this.default.group_km,
				lane_type_id: this.default.lane_type_id,
				surface_type_id: this.default.surface_type_id,
				name: this.default.name,
				roads: this.road_id.map(Number),
				maintenance_analysis_type_id: this.default.maintenance_analysis_type_id,
			}

			const service = new StrategicsService()
			const res = await service.updatePrepareData(this.default.id, params)

			this.prepare_loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.prepare_data_id = res.data.id
				return res
			}
		},
		async nextStep() {
			this.loading = true

			const id = this.default.id

			const service = new StrategicsService()
			const res = await service.creaeteAnalyseStep2(id, this.selected_id)

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.step2Data = res.data
				this.step2Data.plans = res.data.plans?.flatMap((plan) => {
					return Object.keys(plan).reduce((acc, key) => {
						if (plan[key as keyof typeof plan] !== null) {
							acc[key as keyof typeof plan] = plan[key as keyof typeof plan] as any
						}
						return acc
					}, {} as Partial<IStrategicStep2plan>)
				})

				console.log(this.step2Data)

				if (res.data.number_plan > 0) {
					this.createTable()
				}
			}
		},

		onUpdateConditionId(id: number) {
			if (id !== 1) {
				this.step2Data.number_plan = 1
				this.step2Data.year = 1
				this.step2Data.target = this.getTargetOptions[0].value
				this.createTable()
			} else {
				this.step2Data.number_plan = 0
			}
		},

		createTable() {
			const { step2Data } = this

			if (step2Data.condition_id !== 1) {
				const createplan = (planId: number): Partial<IStrategicStep2plan> => {
					return {
						id: planId,
						plan_year: planId + 1,
						...Array(step2Data.number_plan)
							.fill(0)
							.reduce((a, v, i) => ({ ...a, [`plan_${i + 1}`]: v }), {}),
						isNew: true,
					}
				}

				let existingPlan: Partial<IStrategicStep2plan>[] = step2Data.plans.map((plan) => {
					return { ...createplan(plan.id!), ...plan }
				})

				const diffLength = existingPlan.length - step2Data.year

				if (diffLength < 0) {
					existingPlan.push(
						...Array(-diffLength)
							.fill(0)
							.map((_, i) => createplan(existingPlan.length + i))
					)
				} else if (diffLength > 0) {
					existingPlan = existingPlan.slice(0, -diffLength)
				}

				step2Data.plans = existingPlan
			}
		},
		checkParamsStep2() {
			const newPlans = this.step2Data.plans.map((item) => {
				return {
					...item,
					id: !item.isNew ? item.id : null,
				}
			})

			const newObject: IStrategicUpdateAnalyzeParams = {
				comment: this.step2Data.comment,
				condition_id: this.step2Data.condition_id,
				discount: this.step2Data.discount ? this.step2Data.discount : 0,
				number_plan: this.step2Data.number_plan ? this.step2Data.number_plan : 0,
				plans: this.step2Data.condition_id === 1 ? [] : newPlans,
				prepare_data_id: this.selected_id,
				surface_type: this.step2Data.surface_type,
				target: this.step2Data.target ? this.step2Data.target : 0,
				year: this.step2Data.year ? this.step2Data.year : 0,
				name: this.default.name,
			}

			return newObject
		},
		removePlan(plan: number) {
			if (this.step2Data.number_plan !== null) {
				this.step2Data.plans = this.step2Data.plans.map((item) => {
					// สลับค่าตัวที่ ต้องการลบก่อน
					if (plan !== this.step2Data.number_plan) {
						for (let i = plan; i < this.step2Data.number_plan!; i++) {
							item[`plan_${i}` as keyof IStrategicStep2plan] = item[`plan_${i + 1}` as keyof IStrategicStep2plan] as any
						}
					}

					// ลบ key-value pair ของแผนออก
					delete item[`plan_${this.step2Data.number_plan}` as keyof IStrategicStep2plan]

					return item
				})

				// ลบจำนวนแผนของ dropdown
				if (this.step2Data.number_plan > 1) {
					this.step2Data.number_plan--
				}
			}
		},
		async startAnalyze() {
			this.submit_loading = true
			const params = this.checkParamsStep2()

			const service = new StrategicsService()
			const res = await service.createAnalyzing(this.default.id, params)

			this.submit_loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {
		getLaneOptions(state) {
			const { roadGroup } = state

			if (state.road_id.length === 0) {
				return []
			}

			const roadIds = new Set(state.road_id.map(Number))

			const lanes = roadGroup.flatMap((parent) =>
				parent.road_sections.flatMap(
					(section) => section.roads?.filter((road) => roadIds.has(road.id)).map((road) => road.lane_total) || []
				)
			)

			const validLanes = lanes.filter((lane) => lane != null)
			const maxLanes = validLanes.length > 0 ? Math.max(...validLanes) : 0

			const options = [
				{ label: "ทั้งหมด", value: 0 },
				...Array.from({ length: maxLanes }, (_, index) => ({ label: `${index + 1}`, value: index + 1 })),
			]

			return options
		},
		getConditionOptions(state) {
			const { strategicList } = state

			const matchedBudget = strategicList.find((strategic) => strategic.id === 1)?.budget || []
			const options = matchedBudget?.map((budget) => ({ label: budget.name, value: budget.id }))

			return options || []
		},
		getTargetOptions(state) {
			const { strategicList, step2Data } = state
			const conditionId = step2Data.condition_id

			const matchedBudget = strategicList.find((strategic) => strategic.id === 1)?.budget
			const targets = matchedBudget?.find((budget) => budget.id === conditionId)?.target || []
			const options = targets.map((target) => ({ label: target.name, value: target.id })) || []

			return options
		},
		getPath(state) {
			const analysisMapping: Record<number, string> = {
				1: "strategic",
				2: "annual",
			}

			const conditionMapping: Record<number, string> = {
				1: "no-budget-limit",
				2: "budget-limit",
				3: "iri-target",
			}

			const group = analysisMapping[state.default?.maintenance_analysis_type_id] ?? ""
			const criteria = conditionMapping[state.default?.condition] ?? ""

			return {
				group,
				criteria,
			}
		},
	},
})
