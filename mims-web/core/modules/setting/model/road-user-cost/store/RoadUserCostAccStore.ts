import { IRucAccLossData, IRucAccChanceData, IRucAccMasterData } from "../infrastructure/RoadUserCostAccModel"
import { RoadUserCostAccService } from "../infrastructure/RoadUserCostAccService"

interface IState {
	loading: boolean
	dataLoss: IRucAccLossData
	dataAcc: IRucAccChanceData
	accList: IRucAccMasterData[]
	accId: number
	roadGroupID: number
}

export const useRoadUserCostAccStore = defineStore("ruc/acc", {
	state: (): IState => ({
		loading: false,
		dataLoss: {} as IRucAccLossData,
		dataAcc: {} as IRucAccChanceData,
		accList: [
			{
				id: 1,
				name: "มูลค่าความสูญเสีย",
				name_en: "loss_value",
			},
			{
				id: 2,
				name: "โอกาสเกิดอุบัติเหตุ",
				name_en: "chance_of_accident",
			},
		],
		accId: 1,
		roadGroupID: 0,
	}),
	actions: {
		async getAccChanceData(roadGroupId: number) {
			this.loading = true
			const service = new RoadUserCostAccService()
			const res = await service.getChanceOfAccident(roadGroupId)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataAcc = res.data
			}
		},
		async postAccidentChanceParams() {
			this.loading = true
			this.dataAcc.road_group_id = this.roadGroupID
			const params = this.checkParamsChance(this.dataAcc)

			const service = new RoadUserCostAccService()
			const res = await service.postChanceOfAccident(params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getLossAccidentData() {
			this.loading = true

			const service = new RoadUserCostAccService()
			const res = await service.getLossValueAccident()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dataLoss = res.data
			}
		},
		async postAccidentLossValue() {
			const params = this.checkParamsLoss(this.dataLoss)
			this.loading = true

			const service = new RoadUserCostAccService()
			const res = await service.postLossValueAccident(params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		checkParamsLoss(params: IRucAccLossData) {
			const result = {} as IRucAccLossData
			for (const key in params) {
				if (typeof params[key as keyof IRucAccLossData] !== "number") {
					result[key as keyof IRucAccLossData] = Number(params[key as keyof IRucAccLossData])
				} else {
					result[key as keyof IRucAccLossData] = params[key as keyof IRucAccLossData]
				}
			}

			return result
		},
		checkParamsChance(params: IRucAccChanceData) {
			const result = {} as IRucAccChanceData
			for (const key in params) {
				if (typeof params[key as keyof IRucAccChanceData] !== "number") {
					result[key as keyof IRucAccChanceData] = Number(params[key as keyof IRucAccChanceData])
				} else {
					result[key as keyof IRucAccChanceData] = params[key as keyof IRucAccChanceData]
				}
			}

			return result
		},
	},
	getters: {
		getAccListOptions(state) {
			return state.accList.map((item) => ({ label: item.name, value: item.id }))
		},
	},
})
