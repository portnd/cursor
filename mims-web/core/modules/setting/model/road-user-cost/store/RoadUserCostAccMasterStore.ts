import { IRucAccMasterData } from "../infrastructure/RoadUserCostAccModel"
import { RoadUserCostAccService } from "./../infrastructure/RoadUserCostAccService"

interface IFilter {
	masterId: number
}

interface IState {
	loading: boolean
	data: IRucAccMasterData[]
	filter: IFilter
}

export const useRoadUserCostMasterDataStore = defineStore("ruc/masterData", {
	state: (): IState => ({
		loading: false,
		data: [],
		filter: {
			masterId: 0,
		},
	}),
	actions: {
		async getMasterData() {
			this.loading = true

			const service = new RoadUserCostAccService()
			const res = await service.getMasterAccident()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: false })
			} else {
				this.data = res.data
			}
		},
		toggleSettingTable(id: number) {
			switch (id) {
				case 1:
					navigateTo("/settings/models/road-user-cost/acc/loss-value")
					break
				case 2:
					navigateTo("/settings/models/road-user-cost/acc/chance-accident")
					break
			}
		},
	},
	getters: {
		getAccMasterOptions(state) {
			if (state.data.length === 0) {
				return []
			}

			const options = state.data.map((item) => {
				return { label: item.name, value: item.id }
			})

			return options
		},
	},
})
