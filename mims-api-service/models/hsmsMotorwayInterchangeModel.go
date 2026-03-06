package models

import "time"

// Todo ...
type HsmsMotorwayInterchange struct {
	Id              int       `json:"id"`
	RoadId          int       `json:"road_id"`
	RoadAssetId     int       `json:"road_asset_id"`
	SurveyedDate    time.Time `json:"surveyed_date"`
	TheGeom         string    `json:"the_geom"`
	IdParent        int       `json:"id_parent"`
	HashData        string    `json:"hash_data"`
	IsDeleted       bool      `json:"is_deleted"`
	Km              float64   `json:"km"`
	RoadCode        string    `json:"road_code"`
	SectionCode     string    `json:"section_code"`
	Route2          float64   `json:"route2"`
	Control2        float64   `json:"control2"`
	Km2             float64   `json:"km2"`
	Route3          float64   `json:"route3"`
	Control3        float64   `json:"control3"`
	Km3             float64   `json:"km3"`
	Overpass        string    `json:"overpass"`
	Underpass       string    `json:"underpass"`
	InterchangeType string    `json:"interchange_type"`
	DepotName       string    `json:"depot_name"`
	ApproveStatus   string    `json:"approve_status"`
	UpdateBy        string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayInterchange) TableName() string {
	return "hsms_motorway_interchange"
}
