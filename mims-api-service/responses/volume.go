package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type VolumeAadtRevision struct {
	Year int                         `json:"year" extensions:"x-order=0"`
	Item []models.VolumeAadtRevision `json:"items" extensions:"x-order=1"`
}

type Volume struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}

type VolumeAccidentRevision struct {
	Year int                             `json:"year" extensions:"x-order=0"`
	Item []models.VolumeAccidentRevision `json:"items" extensions:"x-order=1"`
}

type VolumeAccident struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}
type VolumeAccidentRespond struct {
	ID          int    `json:"id"`
	RoadGroupID int    `json:"road_group_id"`
	Year        int    `json:"year"`
	CreatedBy   int    `json:"created_by"`
	CreatedDate string `json:"created_date"`
	UpdatedBy   int    `json:"updated_by"`
	UpdatedDate string `json:"updated_date"`
	// TheGeom      string    `json:"the_geom"`
	Revision     int       `json:"revision"`
	Status       string    `json:"status"`
	IDParent     int       `json:"id_parent"`
	RejectReason string    `json:"reject_reason"`
	Acc1         int       `json:"acc1"`
	Acc2         int       `json:"acc2"`
	Acc3         int       `json:"acc3"`
	Acc4         int       `json:"acc4"`
	Total        int       `json:"total"`
	SurveyedDate time.Time `json:"surveyed_date"`
	HashData     string    `json:"hash_data"`
}

type VolumeAadtRespond struct {
	ID          int    `json:"id"`
	RoadGroupID int    `json:"road_group_id"`
	Year        int    `json:"year"`
	CreatedBy   int    `json:"created_by"`
	CreatedDate string `json:"created_date"`
	UpdatedBy   int    `json:"updated_by"`
	UpdatedDate string `json:"updated_date"`
	// TheGeom      string    `json:"the_geom"`
	Revision     int       `json:"revision"`
	IDParent     int       `json:"id_parent"`
	RejectReason string    `json:"reject_reason"`
	Veh1         int       `json:"veh1"`
	Veh2         int       `json:"veh2"`
	Veh3         int       `json:"veh3"`
	Veh4         int       `json:"veh4"`
	Aadt         int       `json:"aadt"`
	Esal         float64   `json:"esal"`
	Yax          float64   `json:"yax"`
	SurveyedDate time.Time `json:"surveyed_date"`
	HashData     string    `json:"hash_data"`
	Status       string    `json:"status"`
}
type VolumTheGeom struct {
	TheGeom string `json:"the_geom"`
}
