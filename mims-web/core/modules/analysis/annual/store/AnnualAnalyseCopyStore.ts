import { AnnualService } from "../infrastructure/AnnualService"
import { IAnnualAnalyzeDataDefault, IAnnualDefaultDataStep2 } from "../infrastructure/AnnualModel"
import { IAnnualStepParams2, IAnnualUpdatePrepareDataParams } from "../infrastructure/AnnualRequest"
import {
	IAnnualAnalyzeData,
	IAnnualAnalyzePrepareData,
	IAnnualRoadsTree,
	IAnnualStrategicsList,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IStateParamsStep1 {
	aadt1: number | null
	aadt2: number | null
	gn1: number | null
	gn2: number | null
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
	defaultData: IAnnualAnalyzeDataDefault
	step: number
	prepareDataId: number[]
	selectedPrepareData: IAnnualAnalyzePrepareData[]
	selectedPrepareId: number[]
	step2Data: IAnnualDefaultDataStep2
	strategicList: IAnnualStrategicsList[]
	targetOptions: IOption[]
	wasCalled: boolean
}
export const useAnnualAnalyseCopyStore = defineStore("annual/edit", {
	state: (): IState => ({
		loading: false,
		roadTrees: [],
		prepareData: {} as IAnnualAnalyzeData,
		defaultData: {} as IAnnualAnalyzeDataDefault,
		params1: {
			aadt1: null,
			aadt2: null,
			gn1: null,
			gn2: null,
			group_km: null,
			iri1: null,
			iri2: null,
			lane_type_id: null,
			road_group_id: null,
			surface_type_id: null,
		},
		step: 1,
		prepareDataId: [],
		selectedPrepareData: [],
		selectedPrepareId: [],
		step2Data: {} as IAnnualDefaultDataStep2,
		strategicList: [],
		targetOptions: [],
		wasCalled: false,
	}),
	actions: {
		async getDefaultData(id: number) {
			this.loading = true

			const service = new AnnualService()
			const res = await service.getDefaultData(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.defaultData = res.data
				this.getRoadsTree()
			}

			this.loading = false

			this.setDefaultParams1(this.defaultData)
		},
		async getRoadsTree() {
			const service = new AnnualService()
			const res = await service.getRoadTrees()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadTrees = res.data

				this.getStrategic()
			}
		},
		async getStrategic() {
			const service = new AnnualService()
			const res = await service.getStrategicList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.strategicList = res.data
			}
		},
		async updatePrepareData(id: number) {
			this.loading = true

			const params: IAnnualUpdatePrepareDataParams = {
				aadt1: this.params1.aadt1 !== null ? this.params1.aadt1 : null,
				aadt2: this.params1.aadt2 !== null ? this.params1.aadt2 : null,
				gn1: this.params1.gn1 !== null ? this.params1.gn1 : null,
				gn2: this.params1.gn2 !== null ? this.params1.gn2 : null,
				group_km: this.params1.group_km!,
				iri1: this.params1.iri1 !== null ? this.params1.iri1 : null,
				iri2: this.params1.iri2 !== null ? this.params1.iri2 : null,
				lane_type_id: this.params1.lane_type_id!,
				roads: this.params1.road_group_id!,
				surface_type_id: this.params1.surface_type_id!,
				maintenance_analysis_type_id: this.strategicList[1].id,
			}

			const service = new AnnualService()
			const res = await service.updatePrepareData(id, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.prepareData = res.data
				console.log(res.data)
			}
		},
		async createAnalyzeStep2() {
			this.loading = true

			const service = new AnnualService()
			const res = await service.createAnnualStep2(
				this.defaultData.id,
				this.selectedPrepareId.filter((item) => item)
			)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				if (!this.wasCalled) {
					this.step2Data = res.data
					this.wasCalled = true
				}

				this.createTargetOptions(this.step2Data.condition_id)
				this.step2Data.total_km = parseFloat(toNumber(this.step2Data.total_km / 1000, 2))
			}
		},
		setDefaultParams1(data: IAnnualAnalyzeDataDefault) {
			this.params1.aadt1 = data.aadt1
			this.params1.aadt2 = data.aadt2
			this.params1.gn1 = data.gn1
			this.params1.gn2 = data.gn2
			this.params1.group_km = data.group_km
			this.params1.iri1 = data.iri1
			this.params1.iri2 = data.iri2
			this.params1.lane_type_id = data.lane_type_id
			this.params1.road_group_id = data.roads.map((item) => item.road_id)
			this.params1.surface_type_id = data.surface_type_id
			this.prepareData.prepare_data = this.defaultData.prepare_data
			this.selectedPrepareData = this.prepareData.prepare_data
			this.selectedPrepareId = this.selectedPrepareData.filter((item) => item.is_selected).map((item) => item.id)
		},
		checkParamsStep2() {
			const newParams: IAnnualStepParams2 = {
				comment: this.step2Data.comment ? this.step2Data.comment : "",
				condition_id: this.step2Data.condition_id,
				gn_avg: this.step2Data.gn_avg,
				iri_avg: this.step2Data.iri_avg,
				surface_type: this.step2Data.surface_type,
				prepare_data_id: this.selectedPrepareId,
				target: this.step2Data.target,
				total_km: this.step2Data.total_km,
				discount: 0,
			}

			switch (this.step2Data.condition_id) {
				case 2:
					newParams.budget = this.step2Data.budget
					break
				case 3:
					newParams.iri = this.step2Data.iri
					break
				case 4:
					newParams.gn = this.step2Data.gn
					break
			}

			return newParams
		},
		async createAnalyse(id: number) {
			this.loading = true
			// const strategicId = this.prepareData.id
			const params = this.checkParamsStep2()

			const service = new AnnualService()
			const res = await service.createAnalye(id, params)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		createTargetOptions(conditionId: number) {
			const strategicItem = this.strategicList.find((item) => item.id === 1)
			const strategicBudget = strategicItem?.budget.find((item) => item.id === conditionId)

			this.targetOptions = strategicBudget?.target.map((item) => ({ label: item.name, value: item.id })) || []
			this.step2Data.target = Number(this.targetOptions[0].value)
		},
		handlePreparedData(item: IAnnualAnalyzePrepareData[]) {
			this.selectedPrepareData = item
			this.selectedPrepareId = this.selectedPrepareData.map((item) => item?.id)
		},
		toDetailsPage() {
			const analysisId = this.defaultData.maintenance_analysis_type_id
			const condition = this.defaultData.condition
			const id = this.defaultData.id

			if (analysisId === 1 || analysisId === 2) {
				const basePath = analysisId === 1 ? "strategic" : "annual"

				switch (condition) {
					case 1:
						navigateTo(`/analyses/${basePath}/summary/${id}/no-budget-limit`)
						break
					case 2:
						navigateTo(`/analyses/${basePath}/summary/${id}/budget-limit`)
						break
					case 3:
						navigateTo(`/analyses/${basePath}/summary/${id}/iri-target`)
						break
					case 4:
						navigateTo(`/analyses/${basePath}/summary/${id}/gn-target`)
						break
				}
			}
		},
		exportDamage(prepareDataId: number) {
			const id = prepareDataId

			if (!id && id !== 0) {
				useHandlerError(0, { message: "โปรดค้นหาข้อมูล Prepare Data" }, { showAlert: true })
			} else {
				useDownloadFile("ดาวน์โหลดรายงานความเสียหาย", `analyze/${id}/export_data`)
			}
		},
		checkFilterParams() {
			this.params1.iri1 = Number.isInteger(this.params1.iri1) ? parseInt(`${this.params1.iri1}`) : this.params1.iri1
			this.params1.iri2 = Number.isInteger(this.params1.iri2) ? parseInt(`${this.params1.iri2}`) : this.params1.iri2
			this.params1.aadt1 = Number.isInteger(this.params1.aadt1) ? parseInt(`${this.params1.aadt1}`) : this.params1.aadt1
			this.params1.aadt2 = Number.isInteger(this.params1.aadt2) ? parseInt(`${this.params1.aadt2}`) : this.params1.aadt2
			this.params1.gn1 = Number.isInteger(this.params1.gn1) ? parseInt(`${this.params1.gn1}`) : this.params1.gn1
			this.params1.gn2 = Number.isInteger(this.params1.gn2) ? parseInt(`${this.params1.gn2}`) : this.params1.gn2
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
