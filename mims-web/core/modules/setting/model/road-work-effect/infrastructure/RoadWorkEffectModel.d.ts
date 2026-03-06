export interface IRoadWorkEffectAC {
	status: boolean
	code: number
	data: IRoadWorkEffectACData
}

export interface IRoadWorkEffectACData {
	as_mol_ac_ab: number
	as_mol_apo_tb: number
	as_mol_ar_vb: number
	as_mol_iri_after_mill_overlay: number
	as_mol_rd_mb: number
	as_ol_ac_ab: number
	as_ol_ar_vb: number
	as_ol_overlay_a0: number
	as_ol_po_tb: number
	as_ol_rd_mb: number
	as_rc_ac_ab: number
	as_rc_apo_tb: number
	as_rc_ar_vb: number
	as_rc_iri_after_reconstruction: number
	as_rc_rd_mb: number
	as_rc_snc: number
	as_rcl_snc: number
	as_rcl_ac_ab: number
	as_rcl_apo_tb: number
	as_rcl_ar_vb: number
	as_rcl_default_hs_old: number
	as_rcl_iri_after_recycling: number
	as_rcl_rd_mb: number
	as_ss_ac_ab: number
	as_ss_apo_tb: number
	as_ss_ar_vb: number
	as_ss_default_lower_bound_iri_after_slurry_seal: number
	as_ss_rd_mb: number
	as_ss_rwe_ss_model_a0: number
}

export interface IRoadWorkEffectCC {
	status: boolean
	code: number
	data: IRoadWorkEffectCCData
}

export interface IRoadWorkEffectCCData {
	cc_fdr_cracking: number
	cc_fdr_faulting: number
	cc_fdr_iri_after_fdr: number
	cc_fdr_spalling: number
	cc_mol_cracking: number
	cc_mol_faulting: number
	cc_mol_iri_after_mol: number
	cc_mol_spalling: number
	cc_ovl_cracking: number
	cc_ovl_faulting: number
	cc_ovl_iri_after_ovl: number
	cc_ovl_spalling: number
	cc_seal_cracking: number
	cc_seal_faulting: number
	cc_seal_iri_after_seal: number
	cc_seal_spalling: number
}
