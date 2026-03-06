package models

import "time"

// Todo ...
type RoadDamage struct {
	Id                   int       `json:"id"`
	RoadId               int       `json:"road_id"`
	LaneNo               int       `json:"lane_no"`
	Year                 int       `json:"year"`
	KmStart              float64   `json:"km_start"`
	KmEnd                float64   `json:"km_end"`
	SurveyedDate         time.Time `json:"surveyed_date"`
	Revision             int       `json:"revision"`
	Status               string    `json:"status"`
	IdParent             int       `json:"id_parent"`
	ImgFilepath          string    `json:"img_filepath"`
	DamageInputFilepath  string    `json:"damage_input_filepath"`
	RejectReason         string    `json:"reject_reason"`
	SurveyType           string    `json:"survey_type"`
	AcIcrack             float64   `json:"ac_icrack"`
	AcUcrack             float64   `json:"ac_ucrack"`
	AcRavelling          float64   `json:"ac_ravelling"`
	AcPatching           float64   `json:"ac_patching"`
	AcPotholeArea        float64   `json:"ac_pothole_area"`
	AcBleeding           float64   `json:"ac_bleeding"`
	AcPotholeCount       float64   `json:"ac_pothole_count"`
	CcTransverseCrack    float64   `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64   `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64   `json:"cc_corner_break"`
	CcJointSealDamage    float64   `json:"cc_joint_seal_damage"`
	CcPatching           float64   `json:"cc_patching"`
	CcSpalling           float64   `json:"cc_spalling"`
	CcScaling            float64   `json:"cc_scaling"`
	CreatedBy            int       `json:"created_by"`
	CreatedDate          time.Time `json:"created_date"`
	UpdatedBy            int       `json:"updated_by"`
	UpdatedDate          time.Time `json:"updated_date"`
}
type RoadDamageList struct {
	RoadDamage
	DirectionId   int    `json:"direction_id"`
	DirectionName string `json:"direction_name"`
}

type RoadDamageListResponse struct {
	Id           int            `json:"id"`
	IdParent     int            `json:"id_parent"`
	Direction    []RefDirection `json:"direction"`
	LaneNo       int            `json:"lane_no"`
	SurveyedDate time.Time      `json:"surveyed_date"`
	Revision     int            `json:"revision"`
}

type RoadDamageForCount struct {
	Id     int    `json:"id"`
	RoadId int    `json:"road_id"`
	Status string `json:"status"`
}

type RoadDamageDetail struct {
	RoadDamage
	RoadDamageStatus RefDataStatus          `json:"road_damage_status" gorm:"ForeignKey:StatusCode;references:Status"`
	RoadDamageRange  []ChildRoadDamageRange `json:"road_damage_range" gorm:"ForeignKey:RoadDamageId;AssociationForeignKey:Id"`
}

type ChildRoadDamageRange struct {
	RoadDamageRange
	RoadDamageM []RoadDamageM `json:"road_damage_m" gorm:"ForeignKey:RoadDamageRangeId;AssociationForeignKey:Id"`
	// ChildRoadM []ChildRoadDamageRange `json:"road_damage_range" gorm:"ForeignKey:RoadDamageRangeId;AssociationForeignKey:Id"`
}

// type ChildRoadM struct {

// 	// RoadDamageRange []RoadDamageRange `json:"road_damage_range" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
// }

// ChildRoad
// TableName use to specific table
func (b *RoadDamage) TableName() string {
	return "road_damage"
}

func (b *RoadDamageForCount) TableName() string {
	return "road_damage"
}
