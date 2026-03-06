package models

type RoadDatebegin struct {
	Id        int     `json:"id"`
	RoadId    int     `json:"road_id"`
	BeginYear int     `json:"begin_year"`
	KmStart   float32 `json:"km_start"`
	KmEnd     float32 `json:"km_end"`
}

func (rt *RoadDatebegin) TableName() string {
	return "road_datebegin"
}
