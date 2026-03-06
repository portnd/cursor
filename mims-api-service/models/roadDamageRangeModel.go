package models

// Todo ...
type RoadDamageRange struct {
	Id                   int     `json:"id"`
	RoadDamageId         int     `json:"road_damage_id"`
	KmStart              float64 `json:"km_start"`
	KmEnd                float64 `json:"km_end"`
	TheGeom              string  `json:"the_geom"`
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
}

// TableName use to specific table
func (b *RoadDamageRange) TableName() string {
	return "road_damage_range"
}
