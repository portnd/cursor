import { ITableAC, ITableACRequest, TableACService } from "../infrastructure"

interface IState {
	loading: boolean
	id: number
	data: ITableAC
}

export const useTableAcStore = defineStore("setting/models/deteriotion/AC", {
	state: (): IState => ({
		data: {} as ITableAC,
		id: 0,
		loading: false,
	}),
	actions: {
		async get(id: number) {
			this.data = {} as ITableAC
			if (id !== null) {
				// Loading
				this.loading = true
				const tableACService = new TableACService()
				const res = await tableACService.get(id)
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
			const params: ITableACRequest = {
				cdb: Number(this.data.cdb),
				cds: Number(this.data.cds),
				comp: Number(this.data.comp),
				kcia: Number(this.data.kcia),
				kciw: Number(this.data.kciw),
				kcpa: Number(this.data.kcpa),
				kcpw: Number(this.data.kcpw),
				kgm: Number(this.data.kgm),
				road_group_id: Number(this.id),
				kgp: Number(this.data.kgp),
				kpi: Number(this.data.kpi),
				kpp: Number(this.data.kpp),
				krid: Number(this.data.krid),
				krpd: Number(this.data.krpd),
				krst: Number(this.data.krst),
				kvi: Number(this.data.kvi),
				kvp: Number(this.data.kvp),
				tlf: Number(this.data.tlf),
				cmod: Number(this.data.cmod),
			}
			const tableACService = new TableACService()
			const res = await tableACService.post(params)
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
