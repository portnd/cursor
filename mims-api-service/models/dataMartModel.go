package models

import "time"

type DataMart struct {
	ID                 int       `json:"id"`
	RoadID             int       `json:"road_id"`
	RoadSurfaceID      int       `json:"road_surface_id"`
	LaneCount          int       `json:"lane_count"`
	SurfaceYear        int       `json:"surface_year"`
	Year               int       `json:"year"`
	ContractNumber     string    `json:"contract_number"`
	KmStart            float64   `json:"km_start"`
	KmEnd              float64   `json:"km_end"`
	LaneNo             int       `json:"lane_no"`
	RefSurfaceID       int       `json:"ref_surface_id"`
	Age                int       `json:"age"`
	LastInspectionDate *string   `json:"last_inspection_date"`
	TheGeom            string    `json:"the_geom"`
	CreatedBy          int       `json:"created_by"`
	UpdatedBy          int       `json:"updated_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (m *DataMart) TableName() string {
	return "data_mart"
}
