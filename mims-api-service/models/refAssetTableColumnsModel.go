package models

type RefAssetTableColumns struct {
	ID                  int    `json:"id" example:"14"`
	RefAssetTableID     int    `json:"ref_asset_table_id" example:"1"`
	ColumnName          string `json:"column_name" example:"id"`
	ColumnSeq           int    `json:"column_seq" example:"1"`
	ColumnDataType      string `json:"column_data_type" example:"integer"`
	ComponentType       string `json:"component_type" example:"hidden"`
	ComponentTitle      string `json:"component_title" example:"คีย์หลัก"`
	ComponentAttributes string `gorm:"default:null" json:"component_attributes" example:""`
	IsRequired          bool   `json:"is_required" example:"true"`
	IsVisibleView       bool   `json:"is_visible_view" example:"false"`
	IsVisibleEdit       bool   `json:"is_visible_edit" example:"true"`
	TableNameRef        string `gorm:"default:null" json:"table_name_ref" example:""`
	IgnoreRefItemList   string `gorm:"default:null" json:"ignore_ref_item_list" example:""`
	IsMandatory         bool   `json:"is_mandatory" example:"true"`
	IsVisibleReport     bool   `json:"is_visible_report" example:"true"`
}

func (ratc *RefAssetTableColumns) TableName() string {
	return "ref_asset_table_columns"
}
