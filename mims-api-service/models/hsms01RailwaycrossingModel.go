package models

import "time"

// Todo ...
type Hsms01Railwaycrossing struct {
	Id                      int       `gorm:"column:id"`
	RailwayCrossingId       int       `gorm:"column:railway_crossing_id"`
	SectionPartId           int       `gorm:"column:section_part_id"`
	Section_Gid             int       `gorm:"column:section_gid"`
	SubsectionId            int       `gorm:"column:subsection_id"`
	Depot                   int       `gorm:"column:depot"`
	Km                      string    `gorm:"column:km"`
	CrossTypeId             int       `gorm:"column:cross_type_id"`
	CrossTypeText           string    `gorm:"column:cross_type_id"`
	HighwayTypeId           int       `gorm:"column:highway_type_id"`
	HighwayTypeText         string    `gorm:"column:highway_type_text"`
	Surface                 int       `gorm:"column:surface"`
	Width                   float64   `gorm:"column:width"`
	ShoulderSurface         int       `gorm:"column:shoulder_surface"`
	ShoulderWidth           float64   `gorm:"column:shoulder_width"`
	Row                     string    `gorm:"column:row"`
	LaneCount               int       `gorm:"column:lane_count"`
	IslandWidth             float64   `gorm:"column:island_width"`
	HasTrafficsign          bool      `gorm:"column:has_trafficsign"`
	HasTrafficpavement      bool      `gorm:"column:has_trafficpavement"`
	HasLight                bool      `gorm:"column:has_light"`
	HasTrafficlight         bool      `gorm:"column:has_trafficlight"`
	HasOtherInventory       bool      `gorm:"column:has_other_inventory"`
	OtherInventory          string    `gorm:"column:other_inventory"`
	HasLiftBeam             bool      `gorm:"column:has_lift_beam"`
	HasStraightBeam         bool      `gorm:"column:has_straight_beam"`
	HasHaul                 bool      `gorm:"column:has_haul"`
	HasFhb                  bool      `gorm:"column:has_fhb"`
	HasFlb                  bool      `gorm:"column:has_flb"`
	HasWarningsign          bool      `gorm:"column:has_warningsign"`
	HasStopsign             bool      `gorm:"column:has_stopsign"`
	HasWinch                bool      `gorm:"column:has_winch"`
	HasPanelup              bool      `gorm:"column:has_panelup"`
	HasOtherSign            bool      `gorm:"column:has_other_sign"`
	OtherSign               string    `gorm:"column:other_sign"`
	RailwayKm               string    `gorm:"column:railway_km"`
	RailwayWidth            float64   `gorm:"column:railway_width"`
	RailwayLandWidth        float64   `gorm:"column:railway_land_width"`
	RailwayDivision         string    `gorm:"column:railway_division"`
	RailwayDistrict         string    `gorm:"column:railway_district"`
	RailwayAadt             int       `gorm:"column:railway_aadt"`
	Remark                  string    `gorm:"column:remark"`
	Status                  string    `gorm:"column:status"`
	ApproveStatus           string    `gorm:"column:approve_status"`
	Revision                int       `gorm:"column:revision"`
	UpdateBy                string    `gorm:"column:update_by"`
	UpdateDate              time.Time `gorm:"column:update_date"`
	ApproveBy               int       `gorm:"column:approve_by"`
	ApproveDate             time.Time `gorm:"column:approve_date"`
	Year                    int       `gorm:"column:year"`
	Geom                    string    `gorm:"column:geom"`
	SectionPartIdRevision   int       `gorm:"column:section_part_id_revision"`
	SubsectionIdRevision    int       `gorm:"column:subsection_id_revision"`
	NeedUpdate              bool      `gorm:"column:need_update"`
	ApproveComment          string    `gorm:"column:approve_comment"`
	MainId                  int       `gorm:"column:road_revision"`
	RoadRevision            int       `gorm:"column:road_revision"`
	Location                string    `gorm:"column:location"`
	RailwayLocation         string    `gorm:"column:railway_location"`
	Tambon                  string    `gorm:"column:tambon"`
	Amphoe                  string    `gorm:"column:amphoe"`
	RoadCode                string    `gorm:"column:road_code"`
	SectionCode             string    `gorm:"column:section_code"`
	Latitude                float64   `gorm:"column:latitude"`
	Longitude               float64   `gorm:"column:longitude"`
	DepotName               string    `gorm:"column:depot_name"`
	SectionKmStart          string    `gorm:"column:section_km_start"`
	SectionKmEnd            string    `gorm:"column:section_km_end"`
	SurfaceTypeId           int       `gorm:"column:surface_type_id"`
	SurfaceTypeText         string    `gorm:"column:surface_type_text"`
	ShoulderSurfaceTypeId   int       `gorm:"column:shoulder_surface_type_id"`
	ShoulderSurfaceTypeText string    `gorm:"column:shoulder_surface_type_text"`
	StatusText              string    `gorm:"column:status_text"`
	RoadId                  int       `gorm:"column:road_id"`
	MimsKmStart             float64   `gorm:"column:mims_km_start"`
	MimsKmEnd               float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Railwaycrossing) TableName() string {
	return "hsms_01_railwaycrossing"
}
