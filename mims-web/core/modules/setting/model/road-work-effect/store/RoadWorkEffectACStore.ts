import { IRoadWorkEffectACParams } from "../infrastructure/RoadWorkEffectRequest"
import { RoadWorkEffectService } from "../infrastructure/RoadWorkEffectService"
import { IRoadWorkEffectACData } from "./../infrastructure/RoadWorkEffectModel.d"

interface IState {
	loading: boolean
	acData: IRoadWorkEffectACData
	params: IRoadWorkEffectACParams
}

export const useRoadWorkEffectACStore = defineStore("road_work_effect/ac", {
	state: (): IState => ({
		loading: false,
		acData: {} as IRoadWorkEffectACData,
		params: {
			as_mol_ac_ab: 0,
			as_mol_apo_tb: 0,
			as_mol_ar_vb: 0,
			as_mol_iri_after_mill_overlay: 0,
			as_mol_rd_mb: 0,
			as_ol_ac_ab: 0,
			as_ol_ar_vb: 0,
			as_ol_overlay_a0: 0,
			as_ol_po_tb: 0,
			as_ol_rd_mb: 0,
			as_rc_ac_ab: 0,
			as_rc_apo_tb: 0,
			as_rc_ar_vb: 0,
			as_rc_iri_after_reconstruction: 0,
			as_rc_rd_mb: 0,
			as_rc_snc: 0,
			as_rcl_snc: 0,
			as_rcl_ac_ab: 0,
			as_rcl_apo_tb: 0,
			as_rcl_ar_vb: 0,
			as_rcl_default_hs_old: 0,
			as_rcl_iri_after_recycling: 0,
			as_rcl_rd_mb: 0,
			as_ss_ac_ab: 0,
			as_ss_apo_tb: 0,
			as_ss_ar_vb: 0,
			as_ss_default_lower_bound_iri_after_slurry_seal: 0,
			as_ss_rd_mb: 0,
			as_ss_rwe_ss_model_a0: 0,
		},
	}),
	actions: {
		async getACValue() {
			this.loading = true

			const service = new RoadWorkEffectService()
			const res = await service.getAsphalt()

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.acData = res.data
				this.setDefaultParams()
			}
		},
		async postACValue() {
			this.loading = true

			const params = this.convertParams(this.params)

			const service = new RoadWorkEffectService()
			const res = await service.postAsphalt(params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		setDefaultParams() {
			if (this.acData) {
				for (const key in this.params) {
					if (key in this.acData) {
						this.params[key as keyof IRoadWorkEffectACParams] = this.acData[key as keyof IRoadWorkEffectACData]
					}
				}
			}
		},
		convertParams(params: IRoadWorkEffectACParams) {
			for (const key in params) {
				if (typeof params[key as keyof IRoadWorkEffectACParams] !== "number") {
					params[key as keyof IRoadWorkEffectACParams] = Number(params[key as keyof IRoadWorkEffectACParams])
				}
			}
			return params
		},
	},
	getters: {},
})
