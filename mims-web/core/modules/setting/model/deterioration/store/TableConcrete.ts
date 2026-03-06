import { ITableConcrete, ITableConcreteRequest, TableConcreteService } from "../infrastructure"

interface IState {
	loading: boolean
	id: number
	data: ITableConcrete
}

export const useTableConcreteStore = defineStore("setting/models/deteriotion/concrete", {
	state: (): IState => ({
		data: {} as ITableConcrete,
		id: 0,
		loading: false,
	}),
	actions: {
		async get(id: number) {
			if (id === null) {
				this.data = {} as ITableConcrete
			} else {
				// Loading
				this.loading = true
				const tableConcreteService = new TableConcreteService()
				const res = await tableConcreteService.get(id)
				// Loading
				this.loading = false
				if (res.status === false) {
					useHandlerError(res.code, res.error, { showAlert: true })
				} else {
					this.data = res.data
					return res
				}
			}
		},
		async post() {
			// Loading
			this.loading = true
			const params: ITableConcreteRequest = {
				b_stress: Number(this.data.b_stress),
				dwl_cor: Number(this.data.dwl_cor),
				ec: Number(this.data.ec),
				fi: Number(this.data.fi),
				jt_space: Number(this.data.jt_space),
				kjrc: Number(this.data.kjrc),
				kjrf: Number(this.data.kjrf),
				kjrr: Number(this.data.kjrr),
				road_group_id: Number(this.id),
				kjrs: Number(this.data.kjrs),
				mi: Number(this.data.mi),
				p_steel: Number(this.data.p_steel),
				pred_seal: Number(this.data.pred_seal),
				widened: Number(this.data.widened),
			}
			const tableConcreteService = new TableConcreteService()
			const res = await tableConcreteService.post(params)
			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
