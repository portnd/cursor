package models

import "time"

type SettingRoadWorkEffectParams struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	CreatedBy int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	IsLatest  bool      `json:"is_latest"`
}

type SettingRoadWorkEffect struct {
	Id                                      int       `json:"id"`
	AsOlOverlayA0                           float64   `json:"as_ol_overlay_a0"`
	AsOlArVb                                float64   `json:"as_ol_ar_vb"`
	AsOlPoTb                                float64   `json:"as_ol_po_tb"`
	AsOlAcAb                                float64   `json:"as_ol_ac_ab"`
	AsOlRdMb                                float64   `json:"as_ol_rd_mb"`
	AsSsRweSsModelA0                        float64   `json:"as_ss_rwe_ss_model_a0"`
	AsSsDefaultLowerBoundIriAfterSlurrySeal float64   `json:"as_ss_default_lower_bound_iri_after_slurry_seal"`
	AsSsArVb                                float64   `json:"as_ss_ar_vb"`
	AsSsApoTb                               float64   `json:"as_ss_apo_tb"`
	AsSsAcAb                                float64   `json:"as_ss_ac_ab"`
	AsSsRdMb                                float64   `json:"as_ss_rd_mb"`
	AsMolIriAfterMillOverlay                float64   `json:"as_mol_iri_after_mill_overlay"`
	AsMolArVb                               float64   `json:"as_mol_ar_vb"`
	AsMolApoTb                              float64   `json:"as_mol_apo_tb"`
	AsMolAcAb                               float64   `json:"as_mol_ac_ab"`
	AsMolRdMb                               float64   `json:"as_mol_rd_mb"`
	AsRclSnc                                float64   `json:"as_rcl_snc"`
	AsRclIriAfterRecycling                  float64   `json:"as_rcl_iri_after_recycling"`
	AsRclArVb                               float64   `json:"as_rcl_ar_vb"`
	AsRclApoTb                              float64   `json:"as_rcl_apo_tb"`
	AsRclAcAb                               float64   `json:"as_rcl_ac_ab"`
	AsRclRdMb                               float64   `json:"as_rcl_rd_mb"`
	AsRclDefaultHsOld                       float64   `json:"as_rcl_default_hs_old"`
	AsRcSnc                                 float64   `json:"as_rc_snc"`
	AsRcIriAfterReconstruction              float64   `json:"as_rc_iri_after_reconstruction"`
	AsRcArVb                                float64   `json:"as_rc_ar_vb"`
	AsRcApoTb                               float64   `json:"as_rc_apo_tb"`
	AsRcAcAb                                float64   `json:"as_rc_ac_ab"`
	AsRcRdMb                                float64   `json:"as_rc_rd_mb"`
	CcFdrIriAfterFdr                        float64   `json:"cc_fdr_iri_after_fdr"`
	CcFdrFaulting                           float64   `json:"cc_fdr_faulting"`
	CcFdrCracking                           float64   `json:"cc_fdr_cracking"`
	CcFdrSpalling                           float64   `json:"cc_fdr_spalling"`
	CcBcoIriAfterBco                        float64   `json:"cc_bco_iri_after_bco"`
	CcBcoFaulting                           float64   `json:"cc_bco_faulting"`
	CcBcoCracking                           float64   `json:"cc_bco_cracking"`
	CcBcoSpalling                           float64   `json:"cc_bco_spalling"`
	CcMolIriAfterMol                        float64   `json:"cc_mol_iri_after_mol"`
	CcMolFaulting                           float64   `json:"cc_mol_faulting"`
	CcMolCracking                           float64   `json:"cc_mol_cracking"`
	CcMolSpalling                           float64   `json:"cc_mol_spalling"`
	CcSealIriAfterSeal                      float64   `json:"cc_seal_iri_after_seal"`
	CcSealFaulting                          float64   `json:"cc_seal_faulting"`
	CcSealCracking                          float64   `json:"cc_seal_cracking"`
	CcSealSpalling                          float64   `json:"cc_seal_spalling"`
	CcRbcIri                                float64   `json:"cc_rbc_iri"`
	CcRbcSlabthk                            float64   `json:"cc_rbc_slabthk"`
	CcRbcPercentFaulting                    float64   `json:"cc_rbc_percent_faulting"`
	CcRbcPercentSpalling                    float64   `json:"cc_rbc_percent_spalling"`
	CcRbcPercentCracking                    float64   `json:"cc_rbc_percent_cracking"`
	IsLatest                                bool      `json:"is_latest"`
	IsDeleted                               bool      `json:"is_deleted"`
	UpdatedBy                               int       `json:"updated_by"`
	CreatedBy                               int       `json:"created_by"`
	UpdatedAt                               time.Time `json:"updated_at"`
	CreatedAt                               time.Time `json:"created_at"`
}

type SettingRoadWorkEffectAsphalt struct {
	AsOlOverlayA0                           float64 `json:"as_ol_overlay_a0"`
	AsOlArVb                                float64 `json:"as_ol_ar_vb"`
	AsOlPoTb                                float64 `json:"as_ol_po_tb"`
	AsOlAcAb                                float64 `json:"as_ol_ac_ab"`
	AsOlRdMb                                float64 `json:"as_ol_rd_mb"`
	AsSsRweSsModelA0                        float64 `json:"as_ss_rwe_ss_model_a0"`
	AsSsDefaultLowerBoundIriAfterSlurrySeal float64 `json:"as_ss_default_lower_bound_iri_after_slurry_seal"`
	AsSsArVb                                float64 `json:"as_ss_ar_vb"`
	AsSsApoTb                               float64 `json:"as_ss_apo_tb"`
	AsSsAcAb                                float64 `json:"as_ss_ac_ab"`
	AsSsRdMb                                float64 `json:"as_ss_rd_mb"`
	AsMolIriAfterMillOverlay                float64 `json:"as_mol_iri_after_mill_overlay"`
	AsMolArVb                               float64 `json:"as_mol_ar_vb"`
	AsMolApoTb                              float64 `json:"as_mol_apo_tb"`
	AsMolAcAb                               float64 `json:"as_mol_ac_ab"`
	AsMolRdMb                               float64 `json:"as_mol_rd_mb"`
	AsRclSnc                                float64 `json:"as_rcl_snc"`
	AsRclIriAfterRecycling                  float64 `json:"as_rcl_iri_after_recycling"`
	AsRclArVb                               float64 `json:"as_rcl_ar_vb"`
	AsRclApoTb                              float64 `json:"as_rcl_apo_tb"`
	AsRclAcAb                               float64 `json:"as_rcl_ac_ab"`
	AsRclRdMb                               float64 `json:"as_rcl_rd_mb"`
	AsRclDefaultHsOld                       float64 `json:"as_rcl_default_hs_old"`
	AsRcSnc                                 float64 `json:"as_rc_snc"`
	AsRcIriAfterReconstruction              float64 `json:"as_rc_iri_after_reconstruction"`
	AsRcArVb                                float64 `json:"as_rc_ar_vb"`
	AsRcApoTb                               float64 `json:"as_rc_apo_tb"`
	AsRcAcAb                                float64 `json:"as_rc_ac_ab"`
	AsRcRdMb                                float64 `json:"as_rc_rd_mb"`
}

type SettingRoadWorkEffectConcrete struct {
	CcFdrIriAfterFdr     float64 `json:"cc_fdr_iri_after_fdr"`
	CcFdrFaulting        float64 `json:"cc_fdr_faulting"`
	CcFdrCracking        float64 `json:"cc_fdr_cracking"`
	CcFdrSpalling        float64 `json:"cc_fdr_spalling"`
	CcBcoIriAfterBco     float64 `json:"cc_bco_iri_after_bco"`
	CcBcoFaulting        float64 `json:"cc_bco_faulting"`
	CcBcoCracking        float64 `json:"cc_bco_cracking"`
	CcBcoSpalling        float64 `json:"cc_bco_spalling"`
	CcMolIriAfterMol     float64 `json:"cc_mol_iri_after_mol"`
	CcMolFaulting        float64 `json:"cc_mol_faulting"`
	CcMolCracking        float64 `json:"cc_mol_cracking"`
	CcMolSpalling        float64 `json:"cc_mol_spalling"`
	CcSealIriAfterSeal   float64 `json:"cc_seal_iri_after_seal"`
	CcSealFaulting       float64 `json:"cc_seal_faulting"`
	CcSealCracking       float64 `json:"cc_seal_cracking"`
	CcSealSpalling       float64 `json:"cc_seal_spalling"`
	CcRbcIri             float64 `json:"cc_rbc_iri"`
	CcRbcSlabthk         float64 `json:"cc_rbc_slabthk"`
	CcRbcPercentFaulting float64 `json:"cc_rbc_percent_faulting"`
	CcRbcPercentSpalling float64 `json:"cc_rbc_percent_spalling"`
	CcRbcPercentCracking float64 `json:"cc_rbc_percent_cracking"`
}

type GetRoadWorkEffectParams struct {
	Concrete SettingRoadWorkEffectConcrete `json:"concrete"`
	Asphalt  SettingRoadWorkEffectAsphalt  `json:"asphalt"`
}

// TableName use to specific table
func (b *SettingRoadWorkEffectParams) TableName() string {
	return "setting_road_work_effect_params"
}

// TableName use to specific table
func (b *SettingRoadWorkEffect) TableName() string {
	return "setting_road_work_effect"
}
