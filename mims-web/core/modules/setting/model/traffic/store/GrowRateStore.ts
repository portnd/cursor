import { GrowthRateService, IGrowRateRequest } from "../infrastructure"
interface IGrowRateState {
	road_group_id: number
	r: string
	code: string
	road_group_name: string
}
interface IState {
	loading: boolean
	data: IGrowRateState[]
}

export const useGrowthRateStore = defineStore("setting/models/aadt/growrate", {
	state: (): IState => ({
		data: [],
		loading: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true
			const growRateService = new GrowthRateService()
			const res = await growRateService.get()
			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = []
				res.data.forEach((e, index) => {
					this.data.push({
						road_group_id: 0,
						r: "",
						code: "",
						road_group_name: "",
					})
					this.data[index].code = e.code
					this.data[index].r = e.r.toString()
					this.data[index].road_group_id = e.road_group_id
					this.data[index].road_group_name = e.road_group_name
				})
				return res
			}
		},
		async post() {
			// Loading
			this.loading = true

			const params: IGrowRateRequest[] = this.data.map((e) => {
				return { r: Number(e.r), road_group_id: e.road_group_id }
			})
			const growRateService = new GrowthRateService()
			const res = await growRateService.post(params)
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
			// Loading
		},
	},
	getters: {},
})
