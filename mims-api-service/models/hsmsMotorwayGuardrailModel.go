package models

import "time"

// Todo ...
type HsmsMotorwayGuardrail struct {
	Id                int       `json:"id"`
	RoadId            int       `json:"road_id"`
	RoadAssetId       int       `json:"road_asset_id"`
	SurveyedDate      time.Time `json:"surveyed_date"`
	TheGeom           string    `json:"the_geom"`
	IdParent          int       `json:"id_parent"`
	HashData          string    `json:"hash_data"`
	IsDeleted         bool      `json:"is_deleted"`
	KmStart           float64   `json:"km_start"`
	KmEnd             float64   `json:"km_end"`
	RoadCode          string    `json:"road_code"`
	SectionCode       string    `json:"section_code"`
	LocationType      string    `json:"location_type"`
	GuardType         string    `json:"guard_type"`
	GuardLeft         string    `json:"guard_left"`
	GuardLeftLength   float64   `json:"guard_left_length"`
	GuardRight        string    `json:"guard_right"`
	GuardRightLength  float64   `json:"guard_right_length"`
	GuardCenter       string    `json:"guard_center"`
	GuardCenterLength float64   `json:"guard_center_length"`
	SetupDate         time.Time `json:"setup_date"`
	PlanYear          string    `json:"plan_year"`
	Contractor        string    `json:"contractor"`
	Budget            float64   `json:"budget"`
	DepotName         string    `json:"depot_name"`
	ApproveStatus     string    `json:"approve_status"`
	UpdateBy          string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayGuardrail) TableName() string {
	return "hsms_motorway_guardrail"
}
