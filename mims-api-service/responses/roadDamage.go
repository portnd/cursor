package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type RoadDamageDetailRespond struct {
	RoadDamageRespond
	RoadDamageStatus models.RefDataStatus          `json:"road_damage_status" gorm:"ForeignKey:StatusCode;AssociationForeignKey:Status"`
	RoadDamageRange  []ChildRoadDamageRangeRespond `json:"road_damage_range" gorm:"ForeignKey:RoadDamageId;AssociationForeignKey:Id"`
}

type RoadDamageRespond struct {
	Id           int       `json:"id"`
	RoadId       int       `json:"road_id"`
	LaneNo       int       `json:"lane_no"`
	Year         int       `json:"year"`
	KmStart      float64   `json:"km_start"`
	KmEnd        float64   `json:"km_end"`
	SurveyedDate time.Time `json:"surveyed_date"`
	// AcCracks             float32   `json:"ac_cracks"`
	Revision             int     `json:"revision"`
	Status               string  `json:"status"`
	IdParent             int     `json:"id_parent"`
	ImgFilepath          string  `json:"img_filepath"`
	DamageInputFilepath  string  `json:"damage_input_filepath"`
	RejectReason         string  `json:"reject_reason"`
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
	CreatedBy            int     `json:"created_by"`
	CreatedDate          string  `json:"created_date"`
	UpdatedBy            int     `json:"updated_by"`
	UpdatedDate          string  `json:"updated_date"`
}

type ChildRoadDamageRangeRespond struct {
	models.RoadDamageRange
	RoadDamageM []RoadDamageMRespond `json:"road_damage_m" gorm:"ForeignKey:RoadDamageRangeId;AssociationForeignKey:Id"`
	// ChildRoadM []ChildRoadDamageRange `json:"road_damage_range" gorm:"ForeignKey:RoadDamageRangeId;AssociationForeignKey:Id"`
}

type RoadDamageMRespond struct {
	Id                int     `json:"id"`
	RoadDamageRangeId int     `json:"road_damage_range_id"`
	Km                float32 `json:"km"`
	// AcCracks              float32   `json:"ac_cracks"`
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
	CreatedBy            int     `json:"created_by"`
	CreatedDate          string  `json:"created_date"`
	UpdatedBy            int     `json:"updated_by"`
	UpdatedDate          string  `json:"updated_date"`
}
