package models

import "time"

// Todo ...
type RoadDamageM struct {
	Id                   int     `json:"id"`
	RoadDamageRangeId    int     `json:"road_damage_range_id"`
	Km                   float32 `json:"km"`
	TheGeom              string  `json:"the_geom"`
	ImgFilepath          string  `json:"img_filepath"`
	HashData             string  `json:"hash_data"`
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

type RoadDamageMPrepareData struct {
	Id                int     `json:"id"`
	RoadDamageRangeId int     `json:"road_damage_range_id"`
	Km                float64 `json:"km"`
	// AcCracks              float32   `json:"ac_cracks"`
	TheGeom               string    `json:"the_geom"`
	ImgFilepath           string    `json:"img_filepath"`
	HashData              string    `json:"hash_data"`
	AcIcrack              *float64  `json:"ac_icrack"`
	AcUcrack              *float64  `json:"ac_ucrack"`
	AcRavelling           *float64  `json:"ac_ravelling"`
	AcPatching            *float64  `json:"ac_patching"`
	AcPothole             *float64  `json:"ac_pothole"`
	AcSurfaceDeform       *float64  `json:"ac_surface_deform"`
	AcBleeding            *float64  `json:"ac_bleeding"`
	CcTransverseCrack     *float64  `json:"cc_transverse_crack"`
	CcNon_transverseCrack *float64  `json:"cc_non_transverse_crack"`
	CcFaulting            *float64  `json:"cc_faulting"`
	CcSpalling            *float64  `json:"cc_spalling"`
	CcCornerbreaks        *float64  `json:"cc_cornerbreaks"`
	CcJointSealDamage     *float64  `json:"cc_joint_seal_damage"`
	CcPatching            *float64  `json:"cc_patching"`
	CreatedBy             int       `json:"created_by"`
	CreatedDate           time.Time `json:"created_date"`
	UpdatedBy             int       `json:"updated_by"`
	UpdatedDate           time.Time `json:"updated_date"`
}

// TableName use to specific table
func (b *RoadDamageM) TableName() string {
	return "road_damage_m"
}

func (b *RoadDamageMPrepareData) TableName() string {
	return "road_damage_m"
}
