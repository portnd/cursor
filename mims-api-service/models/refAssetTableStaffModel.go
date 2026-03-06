package models

type RefAssetTableStaff struct {
	ID              int  `json:"id" example:"50"`
	RefAssetTableID int  `json:"ref_asset_table_id" example:"16"`
	RefDepartmentID int  `json:"ref_department_id" example:"1"`
	IsApprover      bool `json:"is_approver" example:"false"`
}

type AssetTableStaff struct {
	RefAssetTableStaff
	RefDepartment RefDepartment `gorm:"ForeignKey:RefDepartmentID;AssociationForeignKey:ID"`
}

func (rats *RefAssetTableStaff) TableName() string {
	return "ref_asset_table_staff"
}
