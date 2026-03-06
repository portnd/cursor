import { IRoadWorkEffectCCData, IRoadWorkEffectCCParams } from "../infrastructure"
import { RoadWorkEffectService } from "../infrastructure/RoadWorkEffectService"

interface IState {
	loading: boolean
	ccData: IRoadWorkEffectCCData
	params: IRoadWorkEffectCCParams
}

export const useRoadWorkEffectCCStore = defineStore("road_work_effect/concrete", {
	state: (): IState => ({
		loading: false,
		ccData: {} as IRoadWorkEffectCCData,
		params: {
			cc_fdr_cracking: 0,
			cc_fdr_faulting: 0,
			cc_fdr_iri_after_fdr: 0,
			cc_fdr_spalling: 0,
			cc_mol_cracking: 0,
			cc_mol_faulting: 0,
			cc_mol_iri_after_mol: 0,
			cc_mol_spalling: 0,
			cc_bco_cracking: 0,
			cc_bco_faulting: 0,
			cc_bco_iri_after_bco: 0,
			cc_bco_spalling: 0,
			cc_seal_cracking: 0,
			cc_seal_faulting: 0,
			cc_seal_iri_after_seal: 0,
			cc_seal_spalling: 0,
			cc_rbc_iri: 0,
			cc_rbc_slabthk: 0,
			cc_rbc_percent_faulting: 0,
			cc_rbc_percent_spalling: 0,
			cc_rbc_percent_cracking: 0,
		},
	}),
	actions: {
		async getCCValue() {
			this.loading = true

			const service = new RoadWorkEffectService()
			const res = await service.getConcrete()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.ccData = res.data
				this.setDefaultParams()
			}
		},
		async postCCValue() {
			this.loading = true
			const params = this.convertParams(this.params)

			const service = new RoadWorkEffectService()
			const res = await service.postConcrete(params)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		setDefaultParams() {
			if (this.ccData) {
				for (const key in this.params) {
					if (key in this.ccData) {
						this.params[key as keyof IRoadWorkEffectCCParams] = this.ccData[key as keyof IRoadWorkEffectCCData]
					}
				}
			}
		},
		convertParams(params: IRoadWorkEffectCCParams) {
			for (const key in params) {
				if (typeof params[key as keyof IRoadWorkEffectCCParams] !== "number") {
					params[key as keyof IRoadWorkEffectCCParams] = Number(params[key as keyof IRoadWorkEffectCCParams])
				}
			}
			return params
		},
	},
	getters: {},
})
