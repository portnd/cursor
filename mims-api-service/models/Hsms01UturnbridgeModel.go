package models

import "time"

// Todo ...
type Hsms01Uturnbridge struct {
	Id                     int       `gorm:"column:id"`
	UturnBridgeId          int       `gorm:"column:uturn_bridge_id"`
	SectionPartId          int       `gorm:"column:section_part_id"`
	SectionGid             int       `gorm:"column:section_gid"`
	SubsectionId           int       `gorm:"column:subsection_id"`
	Depot                  int       `gorm:"column:depot"`
	Km                     string    `gorm:"column:km"`
	UturnTypeId            int       `gorm:"column:uturn_type_id"`
	UturnTypeText          string    `gorm:"column:uturn_type_text"`
	DirectionTypeId        int       `gorm:"column:direction_type_id"`
	DirectionTypeText      string    `gorm:"column:direction_type_text"`
	IslandTypeId           int       `gorm:"column:island_type_id"`
	IslandTypeText         string    `gorm:"column:island_type_text"`
	IslandWidth            float64   `gorm:"column:island_width"`
	IslandWidthAreaId      int       `gorm:"column:island_width_area_id"`
	IslandWidthAreaText    string    `gorm:"column:island_width_area_text"`
	ConnectedBuildingId    int       `gorm:"column:connected_building_id"`
	ConnectedBuildingText  string    `gorm:"column:connected_building_text"`
	Lane                   int       `gorm:"column:lane"`
	LaneWidth              float64   `gorm:"column:lane_width"`
	AscentLength           float64   `gorm:"column:ascent_length"`
	AscentSlope            float64   `gorm:"column:ascent_slope"`
	CurveLength            float64   `gorm:"column:curve_length"`
	CurveRadias            float64   `gorm:"column:curve_radias"`
	DescentLength          float64   `gorm:"column:descent_length"`
	DescentSlope           float64   `gorm:"column:descent_slope"`
	HasSign                bool      `gorm:"column:has_sign"`
	SignHeight             float64   `gorm:"column:sign_height"`
	SignDistance           float64   `gorm:"column:sign_distance"`
	HasSpeedSign           bool      `gorm:"column:has_speed_sign"`
	SpeedLimit             float64   `gorm:"column:speed_limit"`
	HasHeightSign          bool      `gorm:"column:has_height_sign"`
	HeightLimit            float64   `gorm:"column:height_limit"`
	HasFlashlight          bool      `gorm:"column:has_flashlight"`
	FlashlightHeight       float64   `gorm:"column:flashlight_height"`
	FlashlightDistance     float64   `gorm:"column:flashlight_distance"`
	HasLight               bool      `gorm:"column:has_light"`
	RailBridgeHeight       float64   `gorm:"column:rail_bridge_height"`
	Remark                 string    `gorm:"column:remark"`
	Row                    string    `gorm:"column:row"`
	LaneCount              int       `gorm:"column:lane_count"`
	Status                 string    `gorm:"column:status"`
	ApproveStatus          string    `gorm:"column:approve_status"`
	Revision               int       `gorm:"column:revision"`
	UpdateBy               string    `gorm:"column:update_by"`
	UpdateDate             time.Time `gorm:"column:update_date"`
	ApproveBy              int       `gorm:"column:approve_by"`
	ApproveDate            time.Time `gorm:"column:approve_date"`
	Year                   int       `gorm:"column:year"`
	Geom                   string    `gorm:"column:geom"`
	SectionPartIdRevision  int       `gorm:"column:section_part_id_revision"`
	SubsectionIdRevision   int       `gorm:"column:subsection_id_revision"`
	NeedUpdate             bool      `gorm:"column:need_update"`
	ApproveComment         string    `gorm:"column:approve_comment"`
	MainId                 int       `gorm:"column:main_id"`
	RoadRevision           int       `gorm:"column:road_revision"`
	ConnectedBuildingOther string    `gorm:"column:connected_building_other"`
	IslandOther            string    `gorm:"column:island_other"`
	HasSignSideway         bool      `gorm:"column:has_sign_sideway"`
	HasFlashlightSideway   bool      `gorm:"column:has_flashlight_sideway"`
	BridgeHigh             float64   `gorm:"column:bridge_high"`
	CleaningArea           string    `gorm:"column:cleaning_area"`
	LightCount             int       `gorm:"column:light_count"`
	HasGuard               bool      `gorm:"column:has_guard"`
	GuardType              int       `gorm:"column:guard_type"`
	SignCount              int       `gorm:"column:sign_count"`
	VentilatorCount        int       `gorm:"column:ventilator_count"`
	PumpCount              int       `gorm:"column:pump_count"`
	StartYear              int       `gorm:"column:start_year"`
	Budget                 int       `gorm:"column:budget"`
	LastMaintenanceYear    int       `gorm:"column:last_maintenance_year"`
	Damage                 string    `gorm:"column:damage"`
	RoadCode               string    `gorm:"column:road_code"`
	SectionCode            string    `gorm:"column:section_code"`
	Latitude               float64   `gorm:"column:latitude"`
	Longitude              float64   `gorm:"column:longitude"`
	DepotName              string    `gorm:"column:depot_name"`
	AreaText               string    `gorm:"column:area_text"`
	SectionKmStart         string    `gorm:"column:section_km_start"`
	SectionKmEnd           string    `gorm:"column:section_km_end"`
	StatusText             string    `gorm:"column:status_text"`
	RoadId                 int       `gorm:"column:road_id"`
	MimsKmStart            float64   `gorm:"column:mims_km_start"`
	MimsKmEnd              float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Uturnbridge) TableName() string {
	return "hsms_01_uturnbridge"
}
