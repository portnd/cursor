package models

import "time"

// Todo ...
type HsmsMotorwayFootbridge struct {
	ID            int       `json:"id"`
	RoadAssetID   int       `json:"road_asset_id"`
	RoadId        int       `json:"road_id"`
	SurveyedDate  time.Time `json:"surveyed_date,omitempty"` // omitempty is used since a pointer to time.Time can be nil
	TheGeom       string    `json:"the_geom"`
	IDParent      int       `json:"id_parent"`
	HashData      string    `json:"hash_data"`
	IsDeleted     bool      `json:"is_deleted"`
	KM            float64   `json:"km"`
	RoadCode      string    `json:"road_code"`
	SectionCode   string    `json:"section_code"`
	TypeBridge    int       `json:"type_bridge"`
	BridgeLength  float64   `json:"bridge_length"`
	SpanNum       int       `json:"span_num"`
	BridgeWidth   float64   `json:"bridge_width"`
	BridgeHeight  float64   `json:"bridge_height"`
	SpanWidth     float64   `json:"span_width"`
	PlateHeight   float64   `json:"plate_height"`
	FinishDate    time.Time `json:"finish_date,omitempty"` // omitempty again for possible nil value
	BudgetOwner   int       `json:"budget_owner"`
	Contractor    string    `json:"contractor"`
	Budget        float64   `json:"budget"`
	DepotName     string    `json:"depot_name"`
	ApproveStatus string    `json:"approve_status"`
	// UpdateBy      string    `json:"update_by"`
}

// TableName use to specific table
func (b *HsmsMotorwayFootbridge) TableName() string {
	return "hsms_motorway_footbridge"
}
