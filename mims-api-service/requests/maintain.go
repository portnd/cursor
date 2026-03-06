package requests

type CalDistance struct	{
	RoadID int `json:"road_id"`
	DirectionID int `json:"direction_id"`
	KmStart float64 `json:"km_start"`
	KmEnd float64	`json:"km_end"`
}