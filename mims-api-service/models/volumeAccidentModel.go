package models

import "time"

type VolumeAccident struct {
	ID          int       `json:"id"`
	RoadGroupID int       `json:"road_group_id"`
	Year        int       `json:"year"`
	CreatedBy   int       `json:"created_by"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedDate time.Time `json:"updated_date"`
	// TheGeom      string    `json:"the_geom"`
	Revision     int       `json:"revision"`
	IDParent     int       `json:"id_parent"`
	RejectReason string    `json:"reject_reason"`
	Acc1         int       `json:"acc1"`
	Acc2         int       `json:"acc2"`
	Acc3         int       `json:"acc3"`
	Acc4         int       `json:"acc4"`
	Total        int       `json:"total"`
	SurveyedDate time.Time `json:"surveyed_date"`
	HashData     string    `json:"hash_data"`
	Status       string    `json:"status"`
}

type VolumeAccidentList struct {
	VolumeAccident
	StatusCode     string         `json:"status_code"`
	RoadInfo       []RoadInfo     `json:"road_info"`
	RoadGroup      RoadGroup      `json:"road_group" gorm:"ForeignKey:RoadGroupID;AssociationForeignKey:ID"`
	UserDepartment UserDepartment `json:"updated_by" gorm:"ForeignKey:Id;AssociationForeignKey:UpdatedBy"`
}

type VolumeAccidentRevision struct {
	ID           int       `json:"id"`
	Status       string    `json:"status"`
	Year         int       `json:"year"`
	Revision     int       `json:"revision"`
	IDParent     int       `json:"id_parent"`
	SurveyedDate time.Time `json:"surveyed_date"`
}

// TableName use to specific table
func (b *VolumeAccident) TableName() string {
	return "volume_accident"
}

func (b *VolumeAccidentRevision) TableName() string {
	return "volume_accident"
}
