package models

import "time"

// type SurfaceInfo struct {
// 	RoadID        int64   `json:"road_id"`
// 	RoadCode      string  `json:"road_code"`
// 	RoadSurfaceID int     `json:"road_surface_id" gorm:"column:road_surface_id"`
// 	LaneNo        int     `json:"lane_no"`
// 	LaneCount     int     `json:"lane_count"`
// 	SurfaceID     int64   `json:"surface_id"`
// 	SurfaceName   string  `json:"surface_name"`
// 	ColorCode     string  `json:"color_code"`
// 	KmStart       float64 `json:"km_start"`
// 	KmEnd         float64 `json:"km_end"`
// 	Geometry      string  `json:"geometry" gorm:"column:geom_cl"`
// }

type SurfaceInfo struct {
	RoadID             int64     `json:"road_id"`
	RoadGroupName      string    `json:"road_group_name"`
	RoadName           string    `json:"road_name"`
	RoadCode           string    `json:"road_code"`
	RoadSurfaceID      int       `json:"road_surface_id" gorm:"column:road_surface_id"`
	LaneNo             int       `json:"lane_no"`
	LaneCount          int       `json:"lane_count"`
	SurfaceID          int64     `json:"surface_id"`
	SurfaceName        string    `json:"surface_name"`
	SurfaceGroup       string    `json:"surface_group"`
	ColorCode          string    `json:"color_code"`
	KmStart            float64   `json:"km_start"`
	KmEnd              float64   `json:"km_end"`
	Geometry           string    `json:"geometry" gorm:"column:geom_cl"`
	ContractNumber     string    `json:"contract_number"`
	Year               int       `json:"year"`
	Age                int       `json:"age"`
	LastInspectionDate time.Time `json:"last_inspection_date"`
}

type SurfaceRespond struct {
	Summary  []Summary  `json:"summary"`
	Detail   Detail     `json:"detail"`
	GeomList []GeomList `json:"geom_list"`
}
type GeomList struct {
	Title              string  `json:"title"`
	Color              string  `json:"color"`
	RoadGroupName      string  `json:"road_group_name"`
	ContractNumber     string  `json:"contract_number"`
	Year               string  `json:"year"`
	LastInspectionDate *string `json:"last_inspection_date"`
	RoadName           string  `json:"road_name"`
	KmStart            float64 `json:"km_start"`
	KmEnd              float64 `json:"km_end"`
	KmTotal            float64 `json:"km_total"`
	SurfaceName        string  `json:"surface_name"`
	// SurfaceGroup       string      `json:"surface_group"`
	RefSurfaceID int         `json:"ref_surface_id"`
	Age          int         `json:"age"`
	TheGeom      interface{} `json:"the_geom"`
}
type Detail struct {
	// LaneCountList []int         `json:"lane_count_list"`
	DetailKm []SubDetailKm `json:"detail_km"`
}

type SubDetailKm struct {
	Surface  Surface `json:"surface"`
	LaneNo   int     `json:"lane_no"`
	Value    float64 `json:"value"`
	LaneType string  `json:"lane_type"`
}
type Summary struct {
	Summary Surface `json:"surface"`
	Value   float64 `json:"value"`
	RoadID  string  `json:"road_id"`
}

type Surface struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ColorCode string `json:"color_code" gorm:"column:color"`
}

func (b *Surface) TableName() string {
	return "ref_surface"
}
