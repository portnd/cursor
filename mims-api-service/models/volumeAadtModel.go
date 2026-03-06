package models

import (
	"time"
)

type VolumeAadt struct {
	ID          int       `json:"id"`
	RoadId      int       `json:"road_id"`
	Year        int       `json:"year"`
	CreatedBy   int       `json:"created_by"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedDate time.Time `json:"updated_date"`
	// TheGeom      string    `json:"the_geom"`
	Revision     int       `json:"revision"`
	IDParent     int       `json:"id_parent"`
	RejectReason string    `json:"reject_reason"`
	Veh1         int       `json:"veh1"`
	Veh2         int       `json:"veh2"`
	Veh3         int       `json:"veh3"`
	Total        int       `json:"total"`
	Aadt         int       `json:"aadt"`
	Esal         float64   `json:"esal"`
	Yax          float64   `json:"yax"`
	SurveyedDate time.Time `json:"surveyed_date"`
	HashData     string    `json:"hash_data"`
	Status       string    `json:"status"`
}
type VolumeAadtList struct {
	VolumeAadt
	StatusCode string `json:"status_code"`
	// RoadInfo       RoadInfo       `json:"road_info" gorm:"ForeignKey:RoadId;references:RoadId"`
	// Road           Road           `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:ID"`
	UserDepartment UserDepartment `json:"updated_by" gorm:"ForeignKey:Id; references:UpdatedBy"`
}

type VolumeAadtRevision struct {
	ID           int       `json:"id"`
	Status       string    `json:"status"`
	Year         int       `json:"year"`
	Revision     int       `json:"revision"`
	IDParent     int       `json:"id_parent"`
	SurveyedDate time.Time `json:"surveyed_date"`
}

// TableName use to specific table
func (b *VolumeAadt) TableName() string {
	return "volume_aadt"
}

func (b *VolumeAadtRevision) TableName() string {
	return "volume_aadt"
}
