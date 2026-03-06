package models

import "time"

// Todo ...
type Hsms01Light struct {
	Id                    int       `gorm:"column:id"`
	RoadCode              string    `gorm:"column:road_code"`
	SectionCode           string    `gorm:"column:section_code"`
	MainId                int       `gorm:"column:main_id"`
	SectionGid            int       `gorm:"column:section_gid"`
	SectionPartId         int       `gorm:"column:section_part_id"`
	Polelight             string    `gorm:"column:polelight"`
	PoleTypeText          string    `gorm:"column:pole_type_text"`
	LocationTypeId        int       `gorm:"column:location_type_id"`
	LocationTypeText      string    `gorm:"column:location_type_text"`
	Geom                  string    `gorm:"column:geom"`
	LocationTypeOther     string    `gorm:"column:location_type_other"`
	Location              string    `gorm:"column:location"`
	LampTypeId            int       `gorm:"column:lamp_type_id"`
	LampTypeText          string    `gorm:"column:lamp_type_text"`
	Watt                  string    `gorm:"column:watt"`
	WattOther             string    `gorm:"column:watt_other"`
	Contractor            string    `gorm:"column:contractor"`
	KmStart               string    `gorm:"column:km_start"`
	KmEnd                 string    `gorm:"column:km_end"`
	SetupDate             time.Time `gorm:"column:setup_date"`
	PlanYear              int       `gorm:"column:plan_year"`
	Budget                string    `gorm:"column:budget"`
	Status                string    `gorm:"column:status"`
	ApproveStatus         string    `gorm:"column:approve_status"`
	Depot                 string    `gorm:"column:depot"`
	DepotName             string    `gorm:"column:depot_name"`
	NeedUpdate            bool      `gorm:"column:need_update"`
	Revision              int       `gorm:"column:revision"`
	UpdateDate            time.Time `gorm:"column:update_date"`
	UpdateBy              string    `gorm:"column:update_by"`
	WattTypeId            int       `gorm:"column:watt_type_id"`
	WattTypeText          string    `gorm:"column:watt_type_text"`
	StatusText            string    `gorm:"column:status_text"`
	LocationTypeTextother string    `gorm:"column:location_type_textother"`
	RoadId                int       `gorm:"column:road_id"`
	MimsKmStart           float64   `gorm:"column:mims_km_start"`
	MimsKmEnd             float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Light) TableName() string {
	return "hsms_01_light"
}
