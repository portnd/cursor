package models

import "time"

type RefAadtParameterVehicleType struct {
	Id             int       `json:"id"`
	NumWheel       int       `json:"num_wheel"`
	Name           string    `json:"name"`
	NumAxle        int       `json:"num_axle"`
	LoadEquivalent float64   `json:"load_equivalent"`
	ImagePath      string    `json:"image_path"`
	ForRoadGroupId int       `json:"for_road_group_id"`
	IsLatest       bool      `json:"is_latest"`
	IsDeleted      bool      `json:"is_deleted"`
	UpdatedBy      int       `json:"updated_by"`
	CreatedBy      int       `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

func (pc *RefAadtParameterVehicleType) TableName() string {
	return "ref_aadt_parameter_vehicle_type"
}
