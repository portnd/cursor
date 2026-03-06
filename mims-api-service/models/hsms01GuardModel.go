package models

import "time"

// Todo ...
type Hsms01Guard struct {
	Id                int       `gorm:"column:id"`
	MainId            int       `gorm:"column:main_id"`
	RoadCode          string    `gorm:"column:road_code"`
	SectionCode       string    `gorm:"column:section_code"`
	SectionGid        int       `gorm:"column:section_gid"`
	SectionPartId     int       `gorm:"column:section_part_id"`
	LocationTypeId    int       `gorm:"column:location_type_id"`
	LocationTypeText  string    `gorm:"column:location_type_text"`
	Geom              string    `gorm:"column:geom"`
	LocationTypeOther string    `gorm:"column:location_type_other"`
	Location          string    `gorm:"column:location"`
	GuardTypeId       int       `gorm:"column:guard_type_id"`
	GuardTypeText     string    `gorm:"column:guard_type_text"`
	GuardLeft         int       `gorm:"column:guard_left"`
	GuardLeftLength   int       `gorm:"column:guard_left_length"`
	GuardRight        int       `gorm:"column:guard_right"`
	GuardRightLength  int       `gorm:"column:guard_right_length"`
	Contractor        string    `gorm:"column:contractor"`
	GuardCenter       int       `gorm:"column:guard_center"`
	GuardCenterLength int       `gorm:"column:guard_center_length"`
	KmStart           string    `gorm:"column:km_start"`
	KmEnd             string    `gorm:"column:km_end"`
	SetupDate         time.Time `gorm:"column:setup_date"`
	PlanYear          int       `gorm:"column:plan_year"`
	Budget            int       `gorm:"column:budget"`
	Status            string    `gorm:"column:status"`
	ApproveStatus     string    `gorm:"column:approve_status"`
	Depot             int       `gorm:"column:depot"`
	DepotName         string    `gorm:"column:depot_name"`
	NeedUpdate        bool      `gorm:"column:need_update"`
	Revision          int       `gorm:"column:revision"`
	UpdateDate        time.Time `gorm:"column:update_date"`
	UpdateBy          string    `gorm:"column:update_by"`
	StatusText        string    `gorm:"column:status_text"`
	RoadId            int       `gorm:"column:road_id"`
	MimsKmStart       float64   `gorm:"column:mims_km_start"`
	MimsKmEnd         float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Guard) TableName() string {
	return "hsms_01_guard"
}
