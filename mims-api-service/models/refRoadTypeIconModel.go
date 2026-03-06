package models

// Todo ...
type RefRoadTypeIcon struct {
	ID          int    `json:"id"`
	RoadTypeID  int    `json:"road_type_id"`
	DirectionID int    `json:"direction_id"`
	Icon        string `json:"icon"`
}

// TableName use to specific table
func (b *RefRoadTypeIcon) TableName() string {
	return "ref_road_type_icon"
}
