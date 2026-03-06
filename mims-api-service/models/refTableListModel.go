package models

type RefTableList struct {
	RefName        string `json:"ref_name"`
	RefDescription string `gorm:"column:ref_desc" json:"ref_desc"`
}
