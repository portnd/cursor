package models

import "time"

// Todo ...
type HsmsMotorwayUTurn struct {
	Id                 int       `json:"id"`
	RoadId             int       `json:"road_id"`
	RoadAssetId        int       `json:"road_asset_id"`
	SurveyedDate       time.Time `json:"surveyed_date"`
	TheGeom            string    `json:"the_geom"`
	IdParent           int       `json:"id_parent"`
	HashData           string    `json:"hash_data"`
	IsDeleted          bool      `json:"is_deleted"`
	Km                 float64   `json:"km"`
	RoadCode           string    `json:"road_code"`
	SectionCode        string    `json:"section_code"`
	DirectionType      string    `json:"direction_type"`
	ConnectedBuilding  string    `json:"connected_building"`
	IslandType         string    `json:"island_type"`
	IslandWidth        float64   `json:"island_width"`
	IslandWidthArea    string    `json:"island_width_area"`
	LaneWidth          float64   `json:"lane_width"`
	AscentLength       float64   `json:"ascent_length"`
	AscentSlope        float64   `json:"ascent_slope"`
	CurveLength        float64   `json:"curve_length"`
	CurveRadias        float64   `json:"curve_radias"`
	DescentLength      float64   `json:"descent_length"`
	DescentSlope       float64   `json:"descent_slope"`
	SignHeight         string    `json:"sign_height"`
	SignDistance       float64   `json:"sign_distance"`
	FlashlightHeight   string    `json:"flashlight_height"`
	FlashlightDistance float64   `json:"flashlight_distance"`
	SpeedLimit         string    `json:"speed_limit"`
	HeightLimit        string    `json:"height_limit"`
	HasLight           string    `json:"has_light"`
	RailBridgeHeight   float64   `json:"rail_bridge_height"`
	DepotName          string    `json:"depot_name"`
	ApproveStatus      string    `json:"approve_status"`
	UpdateBy           string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayUTurn) TableName() string {
	return "hsms_motorway_u_turn"
}
