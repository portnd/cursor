package models

import "time"

// Todo ...
type Hsms01Signal struct {
	Id                    int       `gorm:"column:id"`
	MainId                int       `gorm:"column:main_id"`
	RoadCode              string    `gorm:"column:road_code"`
	SectionCode           string    `gorm:"column:section_code"`
	SectionGid            int       `gorm:"column:section_gid"`
	SectionPartId         int       `gorm:"column:section_part_id"`
	Latitude              float64   `gorm:"column:latitude"`
	Longitude             float64   `gorm:"column:longitude"`
	LocationTypeId        int       `gorm:"column:location_type_id"`
	LocationTypeText      string    `gorm:"column:location_type_text"`
	LocationTypeOther     string    `gorm:"column:location_type_other"`
	Location              string    `gorm:"column:location"`
	LampTypeId            int       `gorm:"column:lamp_type_id"`
	LampTypeText          string    `gorm:"column:lamp_type_text"`
	SystemTypeId          int       `gorm:"column:system_type_id"`
	SystemTypeText        string    `gorm:"column:system_type_text"`
	PhaseTypeId           int       `gorm:"column:phase_type_id"`
	PhaseTypeText         string    `gorm:"column:phase_type_text"`
	NumLight              int       `gorm:"column:num_light"`
	NumPole               int       `gorm:"column:num_pole"`
	ControlTypeId         int       `gorm:"column:control_type_id"`
	ControlTypeText       string    `gorm:"column:control_type_text"`
	Contractor            string    `gorm:"column:contractor"`
	Km                    string    `gorm:"column:km"`
	Budget                string    `gorm:"column:budget"`
	ExpireDate            time.Time `gorm:"column:expire_date"`
	Status                string    `gorm:"column:status"`
	ApproveStatus         string    `gorm:"column:approve_status"`
	PhaseOther            string    `gorm:"column:phase_other"`
	Depot                 string    `gorm:"column:depot"`
	DepotName             string    `gorm:"column:depot_name"`
	NeedUpdate            bool      `gorm:"column:need_update"`
	UpdateDate            time.Time `gorm:"column:update_date"`
	UpdateBy              string    `gorm:"column:update_by"`
	LocationTypeTextother string    `gorm:"column:location_type_textother"`
	StatusText            string    `gorm:"column:status_text"`
	RoadId                int       `gorm:"column:road_id"`
	MimsKmStart           float64   `gorm:"column:mims_km_start"`
	MimsKmEnd             float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Signal) TableName() string {
	return "hsms_01_signal"
}
