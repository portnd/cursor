import {
	AnnualService,
	IAnnualAnalyzePrepareData,
	IAnnualParams,
	IAnnualStep2,
	IAnnualStepParams2,
} from "../infrastructure"
import {
	IAnnualRoadsTree,
	IAnnualAnalyzeData,
	IAnnualStrategicsList,
	IAnnualCopy,
	IAnnualDefaultDataStep2,
} from "../infrastructure/AnnualModel"
import { IOption } from "~/core/shared/types/Option"

interface IStateParamsStep1 {
	name: string
	aadt1: number | null
	aadt2: number | null
	ifi1: number | null
	ifi2: number | null
	group_km: number | null
	iri1: number | null
	iri2: number | null
	lane_type_id: number | null
	road_group_id: number[] | null
	surface_type_id: number | null
}

interface IState {
	loading: boolean
	roadTrees: IAnnualRoadsTree[]
	params1: IStateParamsStep1
	prepareData: IAnnualAnalyzeData
	copyData: IAnnualCopy
	step: number
	prepareDataId: number[]
	selectedPrepareData: IAnnualAnalyzePrepareData[]
	step2Data: IAnnualDefaultDataStep2
	params2: IAnnualStep2
	strategicList: IAnnualStrategicsList[]
	targetOptions: IOption[]
	selectedId: number[]
	surfaceOptions: IOption[]
	laneOptions: IOption[]
	groupKmOptions: IOption[]
	wasStep2: boolean
}

export const useAnnualAnalyseCreateStore = defineStore("annual/create", {
	state: (): IState => ({
		loading: false,
		roadTrees: [],
		prepareData: {} as IAnnualAnalyzeData,
		copyData: {} as IAnnualCopy,
		params1: {
			name: "",
			aadt1: null,
			aadt2: null,
			ifi1: null,
			ifi2: null,
			group_km: null,
			iri1: null,
			iri2: null,
			lane_type_id: null,
			road_group_id: null,
			surface_type_id: null,
		},
		params2: {
			comment: "",
			condition_id: 1,
			prepare_data_id: [],
			surface_type: "",
			target: 1,
			total_km: 0,
			ifi_avg: 0,
			iri_avg: 0,
			budget: 14,
			iri: 14,
			name: "",
		},
		step: 1,
		prepareDataId: [],
		selectedPrepareData: [],
		step2Data: {} as IAnnualDefaultDataStep2,
		strategicList: [],
		targetOptions: [],
		selectedId: [],
		surfaceOptions: [
			{ label: "ลาดยาง", value: 1 },
			{ label: "คอนกรีต", value: 2 },
		],
		laneOptions: [
			{ label: "ทั้งหมด", value: 0 },
			{ label: "1", value: 1 },
			{ label: "2", value: 2 },
			{ label: "3", value: 3 },
		],
		groupKmOptions: [
			{ label: "1 กม.", value: 1 },
			{ label: "2 กม.", value: 2 },
			{ label: "5 กม.", value: 5 },
			{ label: "10 กม.", value: 10 },
		],
		wasStep2: false,
	}),
	actions: {
		async getRoadTrees() {
			this.loading = true

			const service = new AnnualService()
			const res = await service.getRoadTrees()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadTrees = res.data

				if (Object.keys(this.copyData).length === 0) {
					this.params1.road_group_id = this.roadTrees.flatMap((item) =>
						item.children.sort((a, b) => a.id - b.id).map((child) => child.id)
					)
				}

				await this.getStrategicList()
			}

			this.loading = false
		},
		async getPrepareData() {
			this.loading = true

			const params: IAnnualParams = {
				name: this.params1.name,
				aadt1: this.params1.aadt1 !== null ? this.params1.aadt1 : null,
				aadt2: this.params1.aadt2 !== null ? this.params1.aadt2 : null,
				ifi1: this.params1.ifi1 !== null ? this.params1.ifi1 : null,
				ifi2: this.params1.ifi2 !== null ? this.params1.ifi2 : null,
				group_km: this.params1.group_km!,
				iri1: this.params1.iri1 !== null ? this.params1.iri1 : null,
				iri2: this.params1.iri2 !== null ? this.params1.iri2 : null,
				lane_type_id: this.params1.lane_type_id!,
				roads: this.params1.road_group_id!,
				surface_type_id: this.params1.surface_type_id!,
				maintenance_analysis_type_id: this.strategicList[1].id,
			}

			const service = new AnnualService()
			const res = await service.createPreparingData(params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.prepareData = res.data
			}
		},
		async createAnalyzeStep2() {
			this.loading = true

			const id = Object.keys(this.copyData).length > 0 ? this.copyData.id : this.prepareData.id

			const service = new AnnualService()
			const res = await service.createAnnualStep2(id, this.getPrepareDataId)

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.step2Data = res.data

				if (Object.keys(this.copyData).length === 0) {
					if (!this.wasStep2) {
						this.params2.condition_id = this.getConditionOptions[0].value
						this.createTargetOptions(this.params2.condition_id)
					}
				} else if (!this.wasStep2) {
					this.params2 = this.step2Data
				}
				this.step2Data.total_km = Number((this.step2Data.total_km / 1000).toFixed(2))
				this.wasStep2 = true
			}
		},
		handlePreparedData(item: IAnnualAnalyzePrepareData[]) {
			this.selectedPrepareData = item
			this.prepareDataId = this.selectedPrepareData?.map((item) => item?.id)
		},
		async getStrategicList() {
			const service = new AnnualService()
			const res = await service.getStrategicList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.strategicList = res.data
			}
		},
		createTargetOptions(conditionId: number) {
			const strategicItem = this.strategicList.find((item) => item.id === 1)
			const strategicBudget = strategicItem?.budget.find((item) => item.id === conditionId)

			this.targetOptions = strategicBudget?.target.map((item) => ({ label: item.name, value: item.id })) || []
			if (Object.keys(this.copyData).length === 0) {
				this.params2.target = Number(this.targetOptions[0].value)
			}
		},
		checkParamsStep2() {
			const newParams: IAnnualStepParams2 = {
				name: this.params2.name ? this.params2.name : "",
				budget: this.params2.budget ? this.params2.budget : 0,
				comment: this.params2.comment ? this.params2.comment : "",
				condition_id: this.params2.condition_id,
				ifi_avg: this.step2Data.ifi_avg,
				iri_avg: this.step2Data.iri_avg,
				surface_type: this.step2Data.surface_type,
				prepare_data_id: this.getPrepareDataId,
				target: this.params2.target,
				total_km: this.step2Data.total_km,
				discount: 0,
			}

			switch (this.params2.condition_id) {
				case 2:
					newParams.budget = this.params2.budget
					break
				case 3:
					newParams.iri = this.params2.iri
					break
				case 4:
					newParams.ifi = this.params2.ifi
					break
			}

			return newParams
		},
		async createAnalyse() {
			this.loading = true
			const strategicId = this.prepareData.id
			const params = this.checkParamsStep2()

			const service = new AnnualService()
			const res = await service.createAnalye(strategicId, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		setCopyParams() {
			this.params1.aadt1 = this.copyData.aadt1
			this.params1.aadt2 = this.copyData.aadt2
			this.params1.ifi1 = this.copyData.ifi1
			this.params1.ifi2 = this.copyData.ifi2
			this.params1.group_km = this.copyData.group_km
			this.params1.iri1 = this.copyData.iri1
			this.params1.iri2 = this.copyData.iri2
			this.params1.lane_type_id = this.copyData.lane_type_id
			this.params1.road_group_id = this.copyData.road_group_id
			this.params1.surface_type_id = this.copyData.surface_type_id
		},
		setDefaultParam1() {
			const surfaceValue = this.surfaceOptions[0]?.value
			this.params1.surface_type_id = typeof surfaceValue === "number" ? surfaceValue : Number(surfaceValue)
			const laneValue = this.laneOptions[0]?.value
			this.params1.lane_type_id = typeof laneValue === "number" ? laneValue : Number(laneValue)
			const groupKmValue = this.groupKmOptions[0]?.value
			this.params1.group_km = typeof groupKmValue === "number" ? groupKmValue : Number(groupKmValue)
		},
		exportDamage(prepareDataId: number) {
			const id = prepareDataId

			if (!id && id !== 0) {
				useHandlerError(0, { message: "โปรดค้นหาข้อมูล Prepare Data" }, { showAlert: true })
			} else {
				useDownloadFile("ดาวน์โหลดรายงานความเสียหาย", `analyze/${id}/export_data`)
			}
		},
	},
	getters: {
		getRoadTreeOptions(state) {
			const roadTrees = state.roadTrees

			if (roadTrees.length === 0) {
				return []
			}

			const options = roadTrees.map((parent) => {
				return {
					id: `parent_${parent.id}`,
					label: parent.label,
					children: parent.children.map((child) => {
						return {
							id: child.id,
							label: child.label,
						}
					}),
				}
			})

			return [{ id: "select_all", label: "เลือกทั้งหมด", children: options }]
		},
		getPrepareDataId(state) {
			const prepareData = state.selectedPrepareData
			if (prepareData.length === 0) {
				return []
			}

			const ids = prepareData.map((item) => item?.id)

			return ids
		},
		getConditionOptions(state) {
			if (state.strategicList.length === 0) {
				return []
			}

			type Option = {
				label: string
				value: number
			}

			const options: Option[] = []

			state.strategicList.forEach((item) => {
				if (item.name === "บำรุงรักษาเชิงกลยุทธ์") {
					item.budget.forEach((budget) => {
						options.push({ label: budget.name, value: budget.id })
					})
				}
			})

			return options
		},
	},
})
