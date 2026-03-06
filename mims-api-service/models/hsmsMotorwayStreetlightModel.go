package models

import "time"

// Todo ...
type HsmsMotorwayStreetlight struct {
	Id            int       `json:"id"`
	RoadId        int       `json:"road_id"`
	RoadAssetId   int       `json:"road_asset_id"`
	SurveyedDate  time.Time `json:"surveyed_date"`
	TheGeom       string    `json:"the_geom"`
	IdParent      int       `json:"id_parent"`
	HashData      string    `json:"hash_data"`
	IsDeleted     bool      `json:"is_deleted"`
	KmStart       float64   `json:"km_start"`
	KmEnd         float64   `json:"km_end"`
	RoadCode      string    `json:"road_code"`
	SectionCode   string    `json:"section_code"`
	LocationType  string    `json:"location_type"`
	LampType      string    `json:"lamp_type"`
	Watt          float64   `json:"watt"`
	PoleType      string    `json:"pole_type"`
	SetupDate     time.Time `json:"setup_date"`
	PlanYear      string    `json:"plan_year"`
	Contractor    string    `json:"contractor"`
	Budget        float64   `json:"budget"`
	DepotName     string    `json:"depot_name"`
	ApproveStatus string    `json:"approve_status"`
	UpdateBy      string    `json:"iupdate_byd"`
}

// TableName use to specific table
func (b *HsmsMotorwayStreetlight) TableName() string {
	return "hsms_motorway_streetlight"
}
