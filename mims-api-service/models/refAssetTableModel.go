package models

type RefAssetTable struct {
	ID              int    `json:"id" example:"1" gorm:"column:id"`
	RefAssetID      int    `json:"ref_asset_id" example:"0"`
	TableNameColumn string `json:"table_name" gorm:"column:table_name" example:"road_asset_in_sign"`
	TableLabel      string `json:"table_label" example:"ป้ายจราจร"`
	IconFilepath    string `gorm:"column:icon_filepath;default:null" json:"icon_filepath" example:"public://icons/icon32/sign.png"`
	LineColorCode   string `gorm:"default:null" json:"line_color_code" example:""`
	Seq             int    `json:"seq" gorm:"column:seq" example:"1"`
	IsInRoad        bool   `json:"is_in_road" example:"true" gorm:"column:is_in_road"`
	IsActive        bool   `json:"is_active" example:"true"`
	GeomType        int    `json:"geom_type" example:"1"`
	CanDelete       bool   `json:"can_delete" example:"false"`
}

type AssetTable struct {
	RefAssetTable
	RefAsset          RefAsset               `gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
	AssetTableColumns []RefAssetTableColumns `gorm:"ForeignKey:RefAssetTableID;AssociationForeignKey:ID"`
}

func (rat *RefAssetTable) TableName() string {
	return "ref_asset_table"
}
