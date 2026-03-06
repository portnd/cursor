package models

import (
	"time"
)

type MaintenanceAnalysisRoad struct {
	ID                    int       `json:"id"`
	MaintenanceAnalysisID int       `json:"maintenance_analysis_id"`
	RoadGroupID           int       `json:"road_group_id"`
	RoadID                int       `json:"road_id"`
	CreatedBy             int       `json:"created_by"`
	UpdatedBy             int       `json:"updated_by"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type MaintenanceAnalysisRoadData struct {
	MaintenanceAnalysisRoad
	RoadInfo RoadInfo `json:"road_info" gorm:"ForeignKey:RoadID;AssociationForeignKey:RoadID"`
}

type MaintenanceAnalysisRoads struct {
	ID                    int `json:"id"`
	MaintenanceAnalysisID int `json:"maintenance_analysis_id"`
	RoadGroupID           int `json:"road_group_id"`
	RoadID                int `json:"road_id"`
}

func (b *MaintenanceAnalysisRoad) TableName() string {
	return "maintenance_analysis_road"
}

func (b *MaintenanceAnalysisRoads) TableName() string {
	return "maintenance_analysis_road"
}
