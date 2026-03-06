package models

import "time"

// Todo ...
type HsmsMotorwayIntersection struct {
	Id                          int       `json:"id"`
	RoadId                      int       `json:"road_id"`
	RoadAssetId                 int       `json:"road_asset_id"`
	SurveyedDate                time.Time `json:"surveyed_date"`
	TheGeom                     string    `json:"the_geom"`
	IdParent                    int       `json:"id_parent"`
	HashData                    string    `json:"hash_data"`
	IsDeleted                   bool      `json:"is_deleted"`
	Km                          float64   `json:"km"`
	RoadCode                    string    `json:"road_code"`
	SectionCode                 string    `json:"section_code"`
	IntersectType               string    `json:"intersect_type"`
	Junction                    string    `json:"junction"`
	SurfaceType                 string    `json:"surface_type"`
	Width                       float64   `json:"width"`
	ShoulderWidth               float64   `json:"shoulder_width"`
	ResponsibleTrafficlightType string    `json:"responsible_trafficlight_type"`
	ResponsibleFlashlightType   string    `json:"responsible_flashlight_type"`
	ResponsibleTrafficSignType  string    `json:"responsible_traffic_sign_type"`
	ResponsibleLightType        string    `json:"responsible_light_type"`
	ResponsibleLightFlType      string    `json:"responsible_light_fl_type"`
	ResponsibleLightHmType      string    `json:"responsible_light_hm_type"`
	ResponsibleLightOtherType   string    `json:"responsible_light_other_type"`
	FeatureType                 string    `json:"feature_type"`
	DepotName                   string    `json:"depot_name"`
	ApproveStatus               string    `json:"approve_status"`
	UpdateBy                    string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayIntersection) TableName() string {
	return "hsms_motorway_intersection"
}
