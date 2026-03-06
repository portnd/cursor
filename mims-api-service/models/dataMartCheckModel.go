package models

import "time"

type DataMartCheck struct {
	Stauts    bool      `json:"stauts"`
	Percent   float64   `json:"percent"`
	UpdatedBy int       `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rc *DataMartCheck) TableName() string {
	return "data_mart_check"
}
