package models

type RefRoadType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type RefRoadTypeInit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rrt *RefRoadType) TableName() string {
	return "ref_road_type"
}

func (rrt *RefRoadTypeInit) TableName() string {
	return "ref_road_type"
}
