package models

import "time"

// Todo ...
type Hsms01Interchange struct {
	Id                    int       `gorm:"column:id"`
	InterchangeId         int       `gorm:"column:interchange_id"`
	SectionPartId         int       `gorm:"column:section_part_id"`
	SectionGid            int       `gorm:"column:section_gid"`
	SubsectionId          int       `gorm:"column:subsection_id"`
	Depot                 int       `gorm:"column:depot"`
	Row                   string    `gorm:"column:row"`
	LaneCount             int       `gorm:"column:lane_count"`
	Km                    string    `gorm:"column:km"`
	Route                 string    `gorm:"column:route"`
	Control               string    `gorm:"column:control"`
	Km2                   string    `gorm:"column:km2"`
	Route2                string    `gorm:"column:route2"`
	Control2              string    `gorm:"column:control2"`
	Km3                   string    `gorm:"column:km3"`
	Route3                string    `gorm:"column:route3"`
	Control3              string    `gorm:"column:control3"`
	Overpass              int       `gorm:"column:overpass"`
	Underpass             int       `gorm:"column:underpass"`
	Interchange           int       `gorm:"column:interchange"`
	Status                string    `gorm:"column:status"`
	ApproveStatus         string    `gorm:"column:approve_status"`
	Revision              int       `gorm:"column:revision"`
	UpdateBy              string    `gorm:"column:update_by"`
	UpdateDate            time.Time `gorm:"column:update_date"`
	ApproveBy             int       `gorm:"column:approve_by"`
	Approvedate           time.Time `gorm:"column:approve_date"`
	Year                  int       `gorm:"column:year"`
	Geom                  string    `gorm:"column:geom"`
	SectionPartIdRevision int       `gorm:"column:section_part_id_revision"`
	SubsectionIdRevision  int       `gorm:"column:subsection_id_revision"`
	NeedUpdate            bool      `gorm:"column:need_update"`
	ApproveComment        string    `gorm:"column:approve_comment"`
	MainId                int       `gorm:"column:main_id"`
	RoadRevision          int       `gorm:"column:road_revision"`
	Location              string    `gorm:"column:location"`
	RoadCode              string    `gorm:"column:road_code"`
	SectionCode           string    `gorm:"column:section_code"`
	Latitude              float64   `gorm:"column:latitude"`
	Longitude             float64   `gorm:"column:longitude"`
	DepotName             string    `gorm:"column:depot_name"`
	SectionKmStart        string    `gorm:"column:section_km_start"`
	SectionKmEnd          string    `gorm:"column:section_km_end"`
	InterchangeTypeId     int       `gorm:"column:interchange_type_id"`
	InterchangeTypeText   string    `gorm:"column:interchange_type_text"`
	StatusText            string    `gorm:"column:status_text"`
	RoadId                int       `gorm:"column:road_id"`
	MimsKmStart           float64   `gorm:"column:mims_km_start"`
	MimsKmEnd             float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Interchange) TableName() string {
	return "hsms_01_interchange"
}
