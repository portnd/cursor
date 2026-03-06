package models

import "time"

// Todo ...
type RoadGeom struct {
	Id        int       `json:"id"`
	RoadId    int       `json:"road_id"`
	LaneNo    int       `json:"lane_no"`
	KmStart   float64   `json:"km_start"`
	KmEnd     float64   `json:"km_end"`
	TheGeom   string    `json:"the_geom"`
	Revision  int       `json:"revision"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy int       `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoadGeomData struct {
	Id         int     `json:"id"`
	RoadID     int     `json:"road_id"`
	LaneNo     int     `json:"lane_no"`
	KmStart    float64 `json:"km_start"`
	KmEnd      float64 `json:"km_end"`
	TheGeom    string  `json:"the_geom"`
	Revision   int     `json:"revision"`
	Status     string  `json:"status"`
	LineString string  `json:"line_string"`
}

type RoadLanes struct {
	Id               int    `json:"-"`
	RoadId           int    `json:"road_id"`
	RefDirectionId   int    `json:"ref_direction_id"`
	RefDirectionName string `json:"ref_direction_name"`
	LaneNo           int    `json:"lane_no"`
}

// TableName use to specific table
func (b *RoadGeom) TableName() string {
	return "road_geom"
}

func (b *RoadGeomData) TableName() string {
	return "road_geom"
}

func (b *RoadLanes) TableName() string {
	return "road_geom"
}
