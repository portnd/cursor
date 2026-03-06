package models

import "time"

// Todo ...
type Hsms01Bridge struct {
	ID                  int       `gorm:"column:id;primaryKey"`
	MainID              int       `gorm:"column:main_id"`
	RoadCode            string    `gorm:"column:road_code;size:10"`
	SectionCode         string    `gorm:"column:section_code;size:10"`
	SectionGID          int       `gorm:"column:section_gid"`
	SectionPartID       int       `gorm:"column:section_part_id"`
	Latitude            float64   `gorm:"column:latitude"`
	Longitude           float64   `gorm:"column:longitude"`
	Location            string    `gorm:"column:location;size:255"`
	KM                  string    `gorm:"column:km;size:10"`
	Budget              int       `gorm:"column:budget"`
	Contractor          string    `gorm:"column:contractor;size:255"`
	FinishDate          time.Time `gorm:"column:finish_date"`
	BridgeLength        float64   `gorm:"column:bridge_length"`
	BridgeType          int       `gorm:"column:bridge_type"`
	BridgeTypeOther     string    `gorm:"column:bridge_type_other;size:255"`
	BridgeWidth         float64   `gorm:"column:bridge_width"`
	BridgeHeight        float64   `gorm:"column:bridge_height"`
	HasPlate            bool      `gorm:"column:has_plate"`
	PlateHeight         float64   `gorm:"column:plate_height"`
	SpanNum             int       `gorm:"column:span_num"`
	SpanWidth           float64   `gorm:"column:span_width"`
	BudgetOwner         int       `gorm:"column:budget_owner"`
	Status              string    `gorm:"column:status;size:1"`
	ApproveStatus       string    `gorm:"column:approve_status;size:1"`
	Depot               int       `gorm:"column:depot"`
	DepotName           string    `gorm:"column:depot_name;size:255"`
	NeedUpdate          bool      `gorm:"column:need_update"`
	UpdateDate          time.Time `gorm:"column:update_date"`
	UpdateBy            string    `gorm:"column:update_by;size:255"`
	StatusText          string    `gorm:"column:status_text;size:255"`
	TypeBridgeID        int       `gorm:"column:type_bridge_id"`
	TypeBridgeText      string    `gorm:"column:type_bridge_text;size:255"`
	TypeBridgeTextOther string    `gorm:"column:type_bridge_textother;size:255"`
	RoadID              int       `gorm:"column:road_id"`
	MimsKMStart         float64   `gorm:"column:mims_km_start"`
	MimsKMEnd           float64   `gorm:"column:mims_km_end"`
}

// TableName use to specific table
func (b *Hsms01Bridge) TableName() string {
	return "hsms_01_bridge"
}
