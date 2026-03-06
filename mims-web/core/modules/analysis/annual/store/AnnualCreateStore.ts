import { IPrepareData } from "../../strategic/infrastructure"
import {
	AnnualService,
	IAnnualDefaultDataStep2,
	IAnnualParams,
	IAnnualRoadsTree,
	IAnnualStepParams2,
} from "../infrastructure"
import { IStrategicsList } from "../../list/infrastructure"
import { ITree } from "~/core/shared/types/Tree"
import { IOption } from "~/core/shared/types/Option"

interface IStateParams {
	road_id: string[]
	surface_type_id: number | null
	group_km: number | null
	lane_type_id: number | null
	iri1: number | null
	iri2: number | null
	aadt1: number | null
	aadt2: number | null
	ifi1: number | null
	ifi2: number | null
	name: string
}

interface IState {
	loading: boolean
	submit_loading: boolean
	prepare_data_loading: boolean
	step: number
	roadGroupOption: ITree[]
	surfaceOptions: IOption[]
	laneOptions: IOption[]
	roadGroup: IAnnualRoadsTree[]
	groupKmOptions: IOption[]
	params1: IStateParams
	prepare_data_id: number | null
	all_ids: number[]
	selected_item: IPrepareData[]
	selected_id: number[]
	step2Data: IAnnualDefaultDataStep2
	strategicList: IStrategicsList[]
	wasStep2: boolean
}

export const useAnnualCreateStore = defineStore("annual-analysis/create", {
	state: (): IState => ({
		loading: false,
		submit_loading: false,
		prepare_data_loading: false,
		step: 1,
		roadGroupOption: [],
		surfaceOptions: [
			{ label: "ลาดยาง", value: 1 },
			{ label: "คอนกรีต", value: 2 },
		],
		laneOptions: [],
		groupKmOptions: [
			{ label: "1 กม.", value: 1 },
			{ label: "2 กม.", value: 2 },
			{ label: "5 กม.", value: 5 },
			{ label: "10 กม.", value: 10 },
		],
		roadGroup: [],
		strategicList: [],
		params1: {
			road_id: [],
			surface_type_id: null,
			group_km: null,
			lane_type_id: null,
			iri1: null,
			iri2: null,
			aadt1: null,
			aadt2: null,
			ifi1: null,
			ifi2: null,
			name: "",
		},
		prepare_data_id: null,
		all_ids: [],
		selected_item: [],
		selected_id: [],
		step2Data: {} as IAnnualDefaultDataStep2,
		wasStep2: false,
	}),
	actions: {
		async getRoadsOptions() {
			this.loading = true

			const service = new AnnualService()
			const res = await service.getRoadTrees()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadGroup = res.data

				const options = res.data?.map((group) => ({
					id: `parent${group.id}`,
					label: `${group.road_number} ${group.short_name}`,
					children: group.road_sections?.map((section) => ({
						id: `section${section.id}`,
						label: `${section.number} ${section?.name_origin} - ${section?.name_destination}`,
						children: section.roads?.map((road) => ({ id: `${road.id}`, label: road.name })),
					})),
				}))

				this.roadGroupOption = options ?? []
			}
		},
		onRoadIdUpdate() {
			this.params1.lane_type_id = null
		},
		async getStrategicList() {
			const service = new AnnualService()
			const res = await service.getStrategicList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.strategicList = res.data
			}
		},
		async createPrepareData() {
			const { params1 } = this
			// reset ค่า
			this.selected_id = []

			const params: IAnnualParams = {
				aadt1: params1.aadt1,
				aadt2: params1.aadt2,
				group_km: params1.group_km,
				ifi1: params1.ifi1,
				ifi2: params1.ifi2,
				iri1: params1.iri1,
				iri2: params1.iri2,
				lane_type_id: params1.lane_type_id,
				roads: params1.road_id.map(Number),
				surface_type_id: params1.surface_type_id,
				maintenance_analysis_type_id: 2,
				name: params1.name,
			}

			const service = new AnnualService()
			const res = await service.createPreparingData(params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.prepare_data_id = res.data.id
			}
		},
		async checkPrepareDataStatus() {
			if (!this.prepare_data_id) {
				return
			}

			const service = new AnnualService()
			const res = await service.getCheckPrepareData(this.prepare_data_id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else if (res.data.status === false) {
				await this.checkPrepareDataStatus()
			} else {
				return res.data.status
			}
		},
		async getPrepareDataId() {
			if (!this.prepare_data_id) {
				return
			}

			const service = new AnnualService()
			const res = await service.getPrepareDataId(this.prepare_data_id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.all_ids = res.data
			}
		},
		handleSelectedPrepareData(itemId: number[]) {
			this.selected_id = itemId
		},
		async submitToStep2() {
			if (!this.prepare_data_id) {
				return
			}

			this.loading = true

			const service = new AnnualService()
			const res = await service.createAnnualStep2(this.prepare_data_id, this.selected_id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else if (this.wasStep2 === false) {
				this.step2Data = res.data
				this.setStep2Default()
			} else {
				this.step2Data.total_km = res.data.total_km
				this.step2Data.iri_avg = res.data.iri_avg
				this.step2Data.ifi_avg = res.data.ifi_avg
				this.step2Data.surface_type = res.data.surface_type
				// this.step2Data.target = this.strategicList.find(item => item.id === this.step2Data.condition_id)?.budget
				// this.step2Data.target = res.data.target
			}

			console.log(this.step2Data.condition_id)
			console.log(this.step2Data.target)
		},
		setStep2Default() {
			if (this.getConditionOptions.length > 0) {
				this.step2Data.condition_id = this.getConditionOptions[0]?.value
				this.step2Data.target = this.getTargetOptions[0]?.value
			}
		},
		onUpdateConditionId(id: number) {
			if (id !== 1) {
				this.step2Data.target = this.getTargetOptions[0].value
			}
		},
		async analyseData() {
			if (!this.prepare_data_id) {
				return
			}

			this.loading = true

			const params: IAnnualStepParams2 = this.checkParamsStep2()

			const service = new AnnualService()
			const res = await service.createAnalye(this.prepare_data_id, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		checkParamsStep2() {
			const newObject: IAnnualStepParams2 = {
				comment: this.step2Data.comment,
				condition_id: this.step2Data.condition_id,
				prepare_data_id: this.selected_id,
				surface_type: this.step2Data.surface_type,
				target: this.step2Data.target ? this.step2Data.target : 0,
				name: this.params1.name,
				ifi_avg: this.step2Data.ifi_avg,
				iri_avg: this.step2Data.iri_avg,
				total_km: this.step2Data.total_km,
				budget: this.step2Data.budget,
				discount: this.step2Data.discount ?? 0,
			}

			return newObject
		},
	},
	getters: {
		getLaneOptions(state) {
			const { roadGroup, params1 } = state

			if (params1.road_id?.length === 0) {
				return []
			}

			const roadIds = new Set(params1.road_id?.map(Number))

			const lanes = roadGroup.flatMap((parent) =>
				parent.road_sections?.flatMap(
					(section) => section.roads?.filter((road) => roadIds.has(road.id))?.map((road) => road.lane_total) || []
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

			const matchedBudget = strategicList.find((strategic) => strategic.id === 2)?.budget || []
			const options = matchedBudget?.map((budget) => ({ label: budget.name, value: budget.id }))

			return options || []
		},
		getTargetOptions(state) {
			const { strategicList, step2Data } = state
			const conditionId = step2Data.condition_id

			const matchedBudget = strategicList.find((strategic) => strategic.id === 2)?.budget
			const targets = matchedBudget?.find((budget) => budget.id === conditionId)?.target || []
			const options = targets.map((target) => ({ label: target.name, value: target.id })) || []

			return options
		},
	},
})
