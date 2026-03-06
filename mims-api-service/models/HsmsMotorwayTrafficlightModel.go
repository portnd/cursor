package models

import "time"

// Todo ...
type HsmsMotorwayTrafficlight struct {
	Id            int       `json:"id"`
	RoadId        int       `json:"road_id"`
	RoadAssetId   int       `json:"road_asset_id"`
	SurveyedDate  time.Time `json:"surveyed_date"`
	TheGeom       string    `json:"ithe_geomd"`
	IdParent      int       `json:"id_parent"`
	HashData      string    `json:"hash_data"`
	IsDeleted     bool      `json:"is_deleted"`
	Km            float64   `json:"km"`
	RoadCode      string    `json:"road_code"`
	SectionCode   string    `json:"section_code"`
	Location      string    `json:"location"`
	LocationType  string    `json:"location_type"`
	LampType      string    `json:"lamp_type"`
	SystemType    string    `json:"system_type"`
	PhaseType     string    `json:"phase_type"`
	NumLight      float64   `json:"num_light"`
	NumPole       float64   `json:"num_pole"`
	ControlType   string    `json:"control_type"`
	ExpireDate    time.Time `json:"expire_date"`
	Contractor    string    `json:"contractor"`
	Budget        float64   `json:"budget"`
	DepotName     string    `json:"depot_name"`
	ApproveStatus string    `json:"approve_status"`
	UpdateBy      string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayTrafficlight) TableName() string {
	return "hsms_motorway_trafficlight"
}
