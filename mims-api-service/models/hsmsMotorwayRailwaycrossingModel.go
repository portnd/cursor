package models

import (
	"time"
)

// Todo ...
type HsmsMotorwayRailwaycrossing struct {
	Id                  int       `json:"id"`
	RoadId              int       `json:"road_id"`
	RoadAssetId         int       `json:"road_asset_id"`
	SurveyedDate        time.Time `json:"surveyed_date"`
	TheGeom             string    `json:"the_geom"`
	IdParent            int       `json:"id_parent"`
	HashData            string    `json:"hash_data"`
	IsDeleted           bool      `json:"is_deleted"`
	Km                  float64   `json:"km"`
	RoadCode            string    `json:"road_code"`
	SectionCode         string    `json:"section_code"`
	CrossType           string    `json:"cross_type"`
	HighwayType         string    `json:"highway_type"`
	SurfaceType         string    `json:"surface_type"`
	Width               float64   `json:"width"`
	ShoulderSurfaceType string    `json:"shoulder_surface_type"`
	ShoulderWidth       float64   `json:"shoulder_width"`
	IslandWidth         float64   `json:"island_width"`
	RailwayDivision     string    `json:"railway_division"`
	RailwayDistrict     string    `json:"railway_district"`
	RailwayKm           float64   `json:"railway_km"`
	RailwayWidth        float64   `json:"railway_width"`
	RailwayLandWidth    float64   `json:"railway_land_width"`
	RailwayAadt         float64   `json:"railway_aadt"`
	DohInventory        string    `json:"doh_inventory"`
	RailwayInventory    string    `json:"railway_inventory"`
	DepotName           string    `json:"depot_name"`
	ApproveStatus       string    `json:"approve_status"`
	UpdateBy            string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayRailwaycrossing) TableName() string {
	return "hsms_motorway_railwaycrossing"
}
