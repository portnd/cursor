package models

import "time"

// Todo ...
type RoadAsset struct {
	Id              int       `json:"id" gorm:"column:id"`
	RoadId          int       `json:"road_id"`
	RefAssetTableId int       `json:"ref_asset_table_id"`
	Revision        int       `json:"revision"`
	Status          string    `json:"status"`
	IdParent        int       `json:"id_parent"`
	IsExclusiveLock bool      `json:"is_exclusive_lock"`
	RejectReason    string    `json:"reject_reason"`
	IsActiveData    bool      `json:"is_active_data"`
	IdTemp          int       `json:"id_temp"`
	CreatedBy       int       `json:"created_by"`
	CreatedDate     time.Time `json:"created_date"`
	UpdatedBy       int       `json:"updated_by"`
	UpdatedDate     time.Time `json:"updated_date"`
}

type RoadAssetForSummary struct {
	Id              int       `json:"id" gorm:"column:id"`
	RoadId          int       `json:"road_id"`
	RefAssetTableId int       `json:"ref_asset_table_id"`
	Revision        int       `json:"revision"`
	Status          string    `json:"status"`
	IdParent        int       `json:"id_parent"`
	IsExclusiveLock bool      `json:"is_exclusive_lock"`
	RejectReason    string    `json:"reject_reason"`
	IsActiveData    bool      `json:"is_active_data"`
	IdTemp          int       `json:"id_temp"`
	CreatedBy       int       `json:"created_by"`
	CreatedDate     time.Time `json:"created_date"`
	UpdatedBy       int       `json:"updated_by"`
	UpdatedDate     time.Time `json:"updated_date"`
	AssetCount      int       `json:"asset_count"`
}

type RoadAssetSummary struct {
	RefAsset
	RefAssetTable []RefAssetTableSummary `json:"ref_asset_table" gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
	// RoadAsset	[]RoadAsset	`json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:ID"`
}

type RefAssetTableSummary struct {
	RefAssetTable
	RoadAsset []RoadAssetForSummary `json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:ID"`
}

type RaData struct {
	ID              int       `json:"id"`
	RoadID          int       `json:"road_id"`
	CreatedDate     time.Time `json:"created_date"`
	CreatedBy       int       `json:"created_by"`
	UpdatedDate     time.Time `json:"updated_date"`
	UpdatedBy       int       `json:"updated_by"`
	IDParent        int       `json:"id_parent"`
	Status          string    `json:"status"`
	Revision        int       `json:"revision"`
	IsExclusiveLock bool      `json:"is_exclusive_lock"`
	RefAssetTableID int       `json:"ref_asset_table_id"`
}

type RoadAssetForCount struct {
	Id     int    `json:"id"`
	RoadId int    `json:"road_id"`
	Status string `json:"status"`
}

// type RoadAssetTableColumnStaff struct {
// 	RefAssetTable
// 	// RoadAsset RoadAsset `json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:Id"`
// 	AssetTableColumn []RefAssetTableColumns `json:"asset_table_column" gorm:"ForeignKey:RefAssetTableID;AssociationForeignKey:Id"`
// 	AssetTableStaff  []RefAssetTableStaff   `json:"asset_table_staff" gorm:"ForeignKey:RefAssetTableID;AssociationForeignKey:Id"`
// }

type RoadAssetRefDataStatus struct {
	RoadAsset
	DataStatus     RefDataStatus  `json:"asset_table_staff" gorm:"ForeignKey:StatusCode;AssociationForeignKey:Status"`
	UserDepartment UserDepartment `json:"updated_by" gorm:"ForeignKey:Id;AssociationForeignKey:UpdatedBy"`
}

type RoadGroupAsset struct {
	RefAsset
	RefAssetTable []RefAssetTableSummary `json:"ref_asset_table" gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
	// RoadAsset	[]RoadAsset	`json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:ID"`
}

// TableName use to specific table
func (b *RoadAsset) TableName() string {
	return "road_asset"
}

// TableName use to specific table
func (b *RoadAssetForCount) TableName() string {
	return "road_asset"
}

// TableName use to specific table
func (b *RaData) TableName() string {
	return "road_asset"
}

// TableName use to specific table
func (b *RoadAssetForSummary) TableName() string {
	return "road_asset"
}
