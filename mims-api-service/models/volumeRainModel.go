package models

import (
	"time"
)

type VolumeRain struct {
	ID          int       `json:"id"`
	RoadGroupID int       `json:"road_group_id"`
	Year        int       `json:"year"`
	Revision    int       `json:"revision"`
	IDParent    int       `json:"id_parent"`
	ProvinceID  int       `json:"province_id"`
	MinRain     float64   `json:"min_rain"`
	MaxRain     float64   `json:"max_rain"`
	AvgRain     float64   `json:"avg_rain"`
	Source      string    `json:"source"`
	Status      string    `json:"status"`
	CreatedBy   int       `json:"created_by"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedDate time.Time `json:"updated_date"`
}

func (b *VolumeRain) TableName() string {
	return "volume_rain"
}
